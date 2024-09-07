package report

import (
	"context"
	"fmt"
	"log"
	"time"

	"uxlyze/analyzer/pkg/analysis"
	"uxlyze/analyzer/pkg/screenshot"
	"uxlyze/analyzer/pkg/types"

	"github.com/chromedp/chromedp"
)

// Generate creates a report for the given URL, analyzing various aspects such as
// visual hierarchy, navigation, mobile friendliness, and readability. It also captures
// screenshots for each of these categories and includes a summary of the website's metadata.
//
// Parameters:
//
//	url - The URL of the website to generate the report for.
//
// Returns:
//
//	*types.Report - A pointer to the generated report containing the analysis results.
//	error - An error if any step of the report generation fails.
func Generate(url string) (*types.Report, error) {
	log.Println("Starting report generation for", url)
	startTime := time.Now()

	// Create a new context for running Chrome DevTools Protocol (chromedp) commands.
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set a timeout for all chromedp actions.
	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// Start timer for navigation.
	stepStart := time.Now()
	err := chromedp.Run(ctx, chromedp.Navigate(url))
	if err != nil {
		return nil, err
	}
	log.Printf("Navigation to URL took: %v\n", time.Since(stepStart))

	// Initialize the report.
	var report types.Report
	report.URL = url
	report.Screenshots = make(map[string]string)

	// Measure and log each step with timing.

	// Step: Generate summary
	stepStart = time.Now()
	report.Summary, err = generateSummary(ctx)
	if err != nil {
		log.Printf("Error generating summary: %v\n", err)
	}
	log.Printf("Generating summary took: %v\n", time.Since(stepStart))

	// Step: Analyze Visual Hierarchy
	stepStart = time.Now()
	report.VisualHierarchy, err = analysis.AnalyzeVisualHierarchy(ctx)
	if err != nil {
		log.Printf("Error analyzing visual hierarchy: %v\n", err)
	}
	log.Printf("Analyzing visual hierarchy took: %v\n", time.Since(stepStart))

	// Step: Analyze Navigation
	stepStart = time.Now()
	report.Navigation, err = analysis.AnalyzeNavigation(ctx)
	if err != nil {
		log.Printf("Error analyzing navigation: %v\n", err)
	}
	log.Printf("Analyzing navigation took: %v\n", time.Since(stepStart))

	// Step: Analyze Mobile Friendliness
	stepStart = time.Now()
	report.MobileFriendliness, err = analysis.AnalyzeMobileFriendliness(ctx)
	if err != nil {
		log.Printf("Error analyzing mobile friendliness: %v\n", err)
	}
	log.Printf("Analyzing mobile friendliness took: %v\n", time.Since(stepStart))

	// Step: Analyze Readability
	stepStart = time.Now()
	report.Readability, err = analysis.AnalyzeReadability(ctx)
	if err != nil {
		log.Printf("Error analyzing readability: %v\n", err)
	}
	log.Printf("Analyzing readability took: %v\n", time.Since(stepStart))

	// Step: Capture Screenshots
	stepStart = time.Now()
	report.Screenshots["VisualHierarchy"], err = screenshot.Capture(ctx, "body")
	if err != nil {
		log.Printf("Error capturing visual hierarchy screenshot: %v\n", err)
	}
	log.Printf("Capturing VisualHierarchy screenshot took: %v\n", time.Since(stepStart))

	stepStart = time.Now()
	report.Screenshots["Navigation"], err = screenshot.Capture(ctx, "nav")
	if err != nil {
		log.Printf("Error capturing navigation screenshot: %v\n", err)
	}
	log.Printf("Capturing Navigation screenshot took: %v\n", time.Since(stepStart))

	// Emulate mobile view and capture mobile friendliness screenshot.
	stepStart = time.Now()
	_ = chromedp.Run(ctx, chromedp.EmulateViewport(375, 812, chromedp.EmulateScale(2.0)))
	report.Screenshots["MobileFriendliness"], err = screenshot.Capture(ctx, "body")
	if err != nil {
		log.Printf("Error capturing mobile friendliness screenshot: %v\n", err)
	}
	log.Printf("Capturing MobileFriendliness screenshot took: %v\n", time.Since(stepStart))

	// Reset to default desktop viewport.
	_ = chromedp.Run(ctx, chromedp.EmulateViewport(0, 0))

	// Capture readability screenshot.
	stepStart = time.Now()
	report.Screenshots["Readability"], err = screenshot.Capture(ctx, "main")
	if err != nil {
		log.Printf("Error capturing readability screenshot: %v\n", err)
	}
	// Fall back to capturing body if main doesn't exist.
	if report.Screenshots["Readability"] == "" {
		report.Screenshots["Readability"], err = screenshot.Capture(ctx, "body")
		if err != nil {
			log.Printf("Error capturing fallback readability screenshot: %v\n", err)
		}
	}
	log.Printf("Capturing Readability screenshot took: %v\n", time.Since(stepStart))

	// Log total time taken for report generation.
	log.Printf("Total report generation time: %v\n", time.Since(startTime))

	// Return the generated report.
	return &report, nil
}

// generateSummary extracts the title and description metadata from the web page and formats them into a summary.
//
// Parameters:
//
//	ctx - The context controlling the chromedp browser instance.
//
// Returns:
//
//	string - A summary containing the title and description of the website.
//	error - An error if any of the extraction steps fail.
func generateSummary(ctx context.Context) (string, error) {
	var title, description string

	// Run chromedp tasks to extract the page title and meta description.
	err := chromedp.Run(ctx,
		chromedp.Title(&title),
		// This JavaScript snippet extracts the content of the "description" meta tag.
		chromedp.EvaluateAsDevTools(`document.querySelector("meta[name='description']")?.getAttribute("content") || "No description found"`, &description),
	)

	// If there was an error, return an empty string and the error.
	if err != nil {
		return "", err
	}

	// Return a formatted summary string containing the title and description.
	return fmt.Sprintf("Website: %s\nDescription: %s", title, description), nil
}
