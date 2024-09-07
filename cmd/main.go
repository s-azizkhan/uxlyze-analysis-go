package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
	"uxlyze/analyzer/pkg/report"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	url := "https://www.logicwind.com" // Replace with actual URL

	startTime := time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	startAlloc := m.Alloc

	rep, err := report.Generate(url)
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

	err = report.Save(rep, "uiux_report.html", psi)
	if err != nil {
		log.Fatalf("Failed to save report: %v", err)
	}

	runtime.ReadMemStats(&m)
	endAlloc := m.Alloc
	duration := time.Since(startTime)

	fmt.Println("UI/UX report generated successfully!")
	fmt.Printf("Time taken: %v\n", duration)
	fmt.Printf("Memory allocated: %v bytes\n", endAlloc-startAlloc)
	fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
}
