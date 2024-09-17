package process_worker

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"uxlyze/analyzer/pkg/report"
	"uxlyze/analyzer/pkg/types"

	"github.com/lib/pq"
)

type DbReportConfig struct {
	IncludePSI        bool `json:"includePSI"`
	SkipUrlFetch      bool `json:"skipUrlFetch"`
	IncludePreview    bool `json:"includePreview"`
	IncludeAIAnalysis bool `json:"includeAIAnalysis"`
}

// Report represents the structure of the data you are fetching
type DbReport struct {
	ID           string
	ProjectID    string
	WebURL       string
	ReportConfig DbReportConfig
	Status       string
}

func isValidURL(toCheck string) bool {
	_, err := url.ParseRequestURI(toCheck)
	if err != nil {
		log.Printf("Invalid URL for URL %s: %s\n", toCheck, err)
		return false
	}
	return true
}

// AnalyzeReportWorker fetches a report by ID from the PostgreSQL database
func AnalyzeReportWorker(id string) {
	fmt.Println("Analyzing report worker for ID:", id)

	// Get the PostgreSQL database URL from the environment
	dbURL := os.Getenv("SUPABASE_DB_URL")
	if dbURL == "" {
		log.Print("DB_URL not provided. Set DB_URL in your environment variables.")
		return
	}

	// Open a connection to the PostgreSQL database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error connecting to the database: %v\n", err)
		return
	}
	defer db.Close()

	// Check the database connection
	err = db.Ping()
	if err != nil {
		log.Printf("Could not connect to the database: %v\n", err)
		return
	}

	// Query the report by ID
	var dbReport DbReport
	var reportConfigBytes []byte
	query := `SELECT id, project_id, web_url, report_config, status FROM reports WHERE id = $1`
	err = db.QueryRow(query, id).Scan(&dbReport.ID, &dbReport.ProjectID, &dbReport.WebURL, &reportConfigBytes, &dbReport.Status)

	// Handle the result
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No report found with ID %s\n", id)
		} else if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "22P02" {
			log.Printf("Invalid UUID format for ID: %s\n", id)
		} else {
			log.Printf("Error fetching report with ID %s: %v\n", id, err)
		}
		return
	}

	// only accept when status is pending
	if dbReport.Status != "pending" {
		log.Printf("Report with ID %s is not pending\n", id)
		return
	}

	// check if the url is valid
	if !isValidURL(dbReport.WebURL) {
		log.Printf("Invalid URL for report ID %s: %s\n", id, dbReport.WebURL)
		return
	}

	// Unmarshal the reportConfigBytes into the ReportConfig struct
	err = json.Unmarshal(reportConfigBytes, &dbReport.ReportConfig)
	if err != nil {
		log.Printf("Error unmarshalling report config for report ID %s: %v\n", id, err)
		return
	}

	// start the analysis
	reportResult, err := report.Generate(
		dbReport.WebURL,
		dbReport.ReportConfig.IncludePreview,
		dbReport.ReportConfig.IncludePSI,
		dbReport.ReportConfig.IncludeAIAnalysis,
	)

	if err != nil {
		log.Printf("Error generating report for report ID %s: %v\n", id, err)
		return
	}

	storeReportResult(reportResult, &dbReport, db)
}

func storeReportResult(reportResult *types.Report, dbReport *DbReport, db *sql.DB) {
	// Serialize reportResult to JSON
	reportJSON, err := json.Marshal(reportResult)
	if err != nil {
		log.Printf("Error serializing report result for report ID %s: %v\n", dbReport.ID, err)
		return
	}

	storeQuery := `
		INSERT INTO report_results (report_id, project_id, result)
		VALUES ($1, $2, $3)
	`

	_, err = db.Exec(storeQuery, dbReport.ID, dbReport.ProjectID, reportJSON)
	if err != nil {
		log.Printf("Error storing report result for report ID %s: %v\n", dbReport.ID, err)
		return
	}

	updateQuery := `
		UPDATE reports
		SET status = 'completed'
		WHERE id = $1
	`

	_, err = db.Exec(updateQuery, dbReport.ID)
	if err != nil {
		log.Printf("Error updating report result for report ID %s: %v\n", dbReport.ID, err)
		return
	}

	log.Printf("Report's result updated for report ID %s\n", dbReport.ID)
}
