package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
	"uxlyze/analyzer/api"
	"uxlyze/analyzer/pkg/report"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Starting UI/UX analysis server...")

	http.HandleFunc("/version", logExecutionTime(api.HandleVersionRequest))
	http.HandleFunc("/analyze", logExecutionTime(api.HandleAnalyzeRequest))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func logExecutionTime(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		handler.ServeHTTP(w, r)
		duration := time.Since(startTime)
		log.Printf("%s %s %v", r.Method, r.URL.Path, duration)
	}
}

func AnalyzeWebsite(url string) {
	startTime := time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	startAlloc := m.Alloc

	rep, err := report.Generate(url, true, true, true)
	if err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s_ui_and_ux_analysis_report.html", strings.Split(url, "://")[1], timestamp)
	err = report.Save(rep, filename)
	if err != nil {
		log.Fatalf("Failed to save report: %v", err)
	} else {
		log.Printf("Report saved to %s\n", filename)
	}

	runtime.ReadMemStats(&m)
	endAlloc := m.Alloc
	duration := time.Since(startTime)

	fmt.Println("UI/UX report generated successfully!")
	fmt.Printf("Time taken: %v\n", duration)
	fmt.Printf("Memory allocated: %v bytes\n", endAlloc-startAlloc)
	fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
}
