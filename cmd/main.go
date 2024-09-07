package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
	"uxlyze/analyzer/pkg/report"
	"uxlyze/analyzer/pkg/types"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Starting UI/UX analysis...")

	AnalyzeWebsite("https://justaziz.com")
	// resp, err := report.AnalyzeUXWithGemini("test.png")
	// if err != nil {
	// 	log.Fatalf("Failed to analyze UX with Gemini: %v", err)
	// }
	log.Printf("UX analysis completed")
}

func AnalyzeWebsite(url string) {

	startTime := time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	startAlloc := m.Alloc

	// capturing SC took 4-5sec more
	rep, err := report.Generate(url, false, types.Both)
	if err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	log.Println("Fetching PageSpeed Insights...")
	psiStart := time.Now()
	psi, err := report.GetPageSpeedInsights(url)
	if err != nil {
		log.Printf("Failed to get PageSpeed Insights: %v\n", err)
	} else {
		log.Printf("PageSpeed Insights fetched in %v\n", time.Since(psiStart))
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s_ui_and_ux_analysis_report.html", strings.Split(url, "://")[1], timestamp)
	err = report.Save(rep, filename, psi)
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
