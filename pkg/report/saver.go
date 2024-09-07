package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"uxlyze/analyzer/pkg/types"
)

func Save(report *types.Report, filename string, psi *types.PageSpeedInsights) error {
	log.Printf("Starting Save function. Saving report to %s", filename)

	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}
	defer file.Close()

	log.Println("File created successfully, generating HTML content...")
	htmlContent, err := generateHTMLContent(report, psi)
	if err != nil {
		log.Printf("Error generating HTML content: %v", err)
		return err
	}

	_, err = file.WriteString(htmlContent)
	if err != nil {
		log.Printf("Error writing HTML content to file: %v", err)
		return err
	}

	log.Println("HTML content written to file successfully.")
	return nil
}

func generateHTMLContent(report *types.Report, psi *types.PageSpeedInsights) (string, error) {
	log.Println("Starting to generate HTML content...")

	templatePath := filepath.Join("pkg", "report", "report_template.html")
	log.Printf("Loading template from %s", templatePath)

	funcMap := template.FuncMap{
		"percentage": func(score float64) string {
			return fmt.Sprintf("%.0f", score*100)
		},
	}

	tmpl, err := template.New("report_template.html").Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	report.Title = "UI/UX Analysis Report for " + strings.Split(report.URL, "://")[1]
	log.Printf("Report title set to: %s", report.Title)

	data := struct {
		*types.Report
		PageSpeedInsights  *types.PageSpeedInsights
		PerformanceMetrics map[string]string
		KeyAudits          []map[string]interface{}
	}{
		Report:             report,
		PageSpeedInsights:  psi,
		PerformanceMetrics: getPerformanceMetrics(psi),
		KeyAudits:          getKeyAudits(psi),
	}

	var buf bytes.Buffer
	log.Println("Executing template with report data...")

	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		return "", fmt.Errorf("error executing template: %v", err)
	}

	log.Println("HTML content generated successfully.")
	return buf.String(), nil
}

func getPerformanceMetrics(psi *types.PageSpeedInsights) map[string]string {
	log.Println("Extracting performance metrics from PageSpeed Insights data...")
	metrics := make(map[string]string)
	for key, value := range psi.LoadingExperience.Metrics {
		metrics[key] = fmt.Sprintf("%d ms (%s)", value.Percentile, value.Category)
		log.Printf("Metric: %s - %s", key, metrics[key])
	}
	log.Println("Performance metrics extraction complete.")
	return metrics
}

func getKeyAudits(psi *types.PageSpeedInsights) []map[string]interface{} {
	log.Println("Extracting key audits from PageSpeed Insights data...")
	keyAuditNames := []string{
		"first-contentful-paint",
		"interactive",
		"speed-index",
		"total-blocking-time",
		"largest-contentful-paint",
		"cumulative-layout-shift",
	}

	var keyAudits []map[string]interface{}
	for _, name := range keyAuditNames {
		if audit, ok := psi.LighthouseResult.Audits[name]; ok {
			log.Printf("Audit found for: %s - Score: %v", name, audit.Score)
			keyAudits = append(keyAudits, map[string]interface{}{
				"name":         name,
				"score":        audit.Score,
				"title":        audit.Title,
				"displayValue": audit.DisplayValue,
			})
		} else {
			log.Printf("Audit not found for: %s", name)
		}
	}
	log.Println("Key audits extraction complete.")
	return keyAudits
}

func GetPageSpeedInsights(url string) (*types.PageSpeedInsights, error) {
	log.Printf("Fetching PageSpeed Insights for URL: %s", url)

	apiKey := os.Getenv("GOOGLE_PAGE_SPEED_API_KEY")
	if apiKey == "" {
		log.Println("GOOGLE_PAGE_SPEED_API_KEY is not set")
		return nil, fmt.Errorf("GOOGLE_PAGE_SPEED_API_KEY is not set")
	}

	apiURL := fmt.Sprintf("https://www.googleapis.com/pagespeedonline/v5/runPagespeed?url=%s&key=%s", url, apiKey)
	log.Printf("API URL: %s", apiURL)

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Error fetching PageSpeed Insights: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Println("PageSpeed Insights response received, reading body...")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	var psi types.PageSpeedInsights
	err = json.Unmarshal(body, &psi)
	if err != nil {
		log.Printf("Error unmarshaling PageSpeed Insights data: %v", err)
		return nil, err
	}

	log.Println("PageSpeed Insights data successfully fetched and unmarshaled.")
	return &psi, nil
}
