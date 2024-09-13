package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
	"uxlyze/analyzer/pkg/report"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Starting UI/UX analysis server...")

	AnalyzeWebsite("https://logicwind.com")
}

// func startServer() {
// 	http.HandleFunc("/version", logExecutionTime(api.HandleVersionRequest))
// 	http.HandleFunc("/analyze", logExecutionTime(api.HandleAnalyzeRequest))
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func logExecutionTime(handler http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		startTime := time.Now()
// 		defer func() {
// 			if err := recover(); err != nil {
// 				log.Printf("Panic occurred: %v", err)
// 				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 			}
// 			duration := time.Since(startTime)
// 			log.Printf("%s %s %v", r.Method, r.URL.Path, duration)
// 		}()
// 		handler.ServeHTTP(w, r)
// 	}
// }

func AnalyzeWebsite(url string) {
	startTime := time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	startAlloc := m.Alloc

	rep, err := report.Generate(url, false, false, false)
	if err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	// save resport to json
	jsonData, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal report to JSON: %v", err)
	}

	// save json to file
	err = os.WriteFile(fmt.Sprintf("%s_ui_and_ux_analysis_report.json", strings.Split(url, "://")[1]), jsonData, 0o644)
	if err != nil {
		log.Fatalf("Failed to save report: %v", err)
	}
	// timestamp := time.Now().Format("2006-01-02_15-04-05")
	// filename := fmt.Sprintf("%s_%s_ui_and_ux_analysis_report.html", strings.Split(url, "://")[1], timestamp)
	// err = report.Save(rep, filename)
	// if err != nil {
	// 	log.Fatalf("Failed to save report: %v", err)
	// } else {
	// 	log.Printf("Report saved to %s\n", filename)
	// }

	runtime.ReadMemStats(&m)
	endAlloc := m.Alloc
	duration := time.Since(startTime)

	fmt.Println("UI/UX report generated successfully!")
	fmt.Printf("Time taken: %v\n", duration)
	fmt.Printf("Memory allocated: %v bytes\n", endAlloc-startAlloc)
	fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
}
