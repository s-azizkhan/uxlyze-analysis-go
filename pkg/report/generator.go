package report

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"uxlyze/analyzer/pkg/ai"
	"uxlyze/analyzer/pkg/analysis"
	"uxlyze/analyzer/pkg/screenshot"
	"uxlyze/analyzer/pkg/types"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

func SaveBase64ToLocal(base64String string, pathName string) {
	// Split the Base64 string into data URI parts (if it contains metadata)
	parts := strings.SplitN(base64String, ",", 2)
	if len(parts) != 2 {
		log.Fatal("Invalid Base64 data URI format")
	}

	// Detect the MIME type from the metadata
	rawBase64 := parts[1]

	// Decode the Base64 string into file bytes
	fileData, err := base64.StdEncoding.DecodeString(rawBase64)
	if err != nil {
		log.Fatal("Failed to decode Base64 string:", err)
	}

	// Save the file to the local system
	err = os.WriteFile(pathName, fileData, 0644)
	if err != nil {
		log.Fatal("Failed to save the file:", err)
	}

	fmt.Println("File saved successfully as:", pathName)
}

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
func Generate(url string, takeScreenshots bool, includePSI bool, includeAIAnalysis bool) (*types.Report, error) {
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

	// // Step: Analyze Visual Hierarchy
	// stepStart = time.Now()
	// report.FontSizes, err = analysis.AnalyzeFontSizes(ctx)
	// if err != nil {
	// 	log.Printf("Error analyzing visual hierarchy: %v\n", err)
	// }
	// log.Printf("Analyzing visual hierarchy took: %v\n", time.Since(stepStart))

	// Step: Analyze Mobile Friendliness
	stepStart = time.Now()
	report.MobileFriendly, err = analysis.AnalyzeMobileFriendly(ctx)
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

	// Step: Analyze Navigation
	stepStart = time.Now()
	report.Navigation, err = analysis.AnalyzeNavigation(ctx)
	if err != nil {
		log.Printf("Error analyzing navigation: %v\n", err)
	}
	log.Printf("Analyzing navigation took: %v\n", time.Since(stepStart))

	// Step: Analyze Color Usage
	stepStart = time.Now()
	report.ColorUsage, err = analysis.AnalyzeColorUsage(ctx)
	if err != nil {
		log.Printf("Error analyzing color usage: %v\n", err)
	}
	log.Printf("Analyzing color usage took: %v\n", time.Since(stepStart))

	// Step: Analyze Font Usage
	stepStart = time.Now()
	report.FontUsage, err = analysis.AnalyzeFontUsage(ctx)
	if err != nil {
		log.Printf("Error analyzing font usage: %v\n", err)
	}
	log.Printf("Analyzing font usage took: %v\n", time.Since(stepStart))

	// Step: Analyze SEO
	stepStart = time.Now()
	report.SEO, err = analysis.AnalyzeSEO(ctx)
	if err != nil {
		log.Printf("Error analyzing SEO : %v\n", err)
	}
	log.Printf("Analyzing SEO took: %v\n", time.Since(stepStart))

	// Step: Capture Screenshots
	if takeScreenshots {

		stepStart = time.Now()
		report.Screenshots["Desktop"], err = screenshot.Capture(ctx, "body")
		if err != nil {
			log.Printf("Error capturing Desktop screenshot: %v\n", err)
		}
		log.Printf("Capturing Desktop screenshot took: %v\n", time.Since(stepStart))

		stepStart = time.Now()
		report.Screenshots["Navigation"], err = screenshot.Capture(ctx, "nav")
		if err != nil {
			log.Printf("Error capturing navigation screenshot: %v\n", err)
		}
		log.Printf("Capturing Navigation screenshot took: %v\n", time.Since(stepStart))

		// Emulate mobile view and capture mobile friendliness screenshot.
		stepStart = time.Now()
		// calculate the page height
		var pageHeight int64
		err = chromedp.Run(ctx, chromedp.EvaluateAsDevTools(`document.body.scrollHeight`, &pageHeight))
		if err != nil {
			log.Printf("Error calculating page height: %v\n", err)
		}
		_ = chromedp.Run(ctx, chromedp.EmulateViewport(350, pageHeight, chromedp.EmulateScale(2.0)))
		report.Screenshots["Mobile"], err = screenshot.Capture(ctx, "body")
		if err != nil {
			log.Printf("Error capturing mobile friendliness screenshot: %v\n", err)
		}
		log.Printf("Capturing Mobile screenshot took: %v\n", time.Since(stepStart))

		// Reset to default desktop viewport.
		_ = chromedp.Run(ctx, chromedp.EmulateViewport(0, 0))

	}
	// Perform Gemini UX analysis
	if includeAIAnalysis {

		if report.Screenshots["Desktop"] == "" {
			log.Println("No screenshot available for Gemini analysis")
			//  take the screenshot
			report.Screenshots["Desktop"], err = screenshot.Capture(ctx, "body")
			if err != nil {
				log.Printf("Error capturing Desktop screenshot: %v\n", err)
			}
		}
		tempUuid := uuid.New()
		tempImagePath := "temp_screenshot_" + tempUuid.String() + ".png"
		SaveBase64ToLocal("data:image/png;base64,"+report.Screenshots["Desktop"], tempImagePath)
		if err != nil {
			log.Printf("Error saving temporary screenshot: %v\n", err)
		} else {
			geminiAnalysis, err := ai.AnalyzeUXWithGemini(tempImagePath)
			if err != nil {
				log.Printf("Error analyzing UX with Gemini: %v\n", err)
			} else {
				report.GeminiAnalysis = geminiAnalysis
			}
			os.Remove(tempImagePath)
		}

	}

	if !takeScreenshots {
		report.Screenshots["Desktop"] = ""
		report.Screenshots["Mobile"] = ""
		report.Screenshots["Navigation"] = ""
	}

	if includePSI {
		stepStart = time.Now()
		psi, err := GetPageSpeedInsights(url)
		if err != nil {
			log.Printf("Error getting PageSpeed Insights: %v\n", err)
		} else {
			report.PageSpeedInsights = psi
			log.Printf("Getting PageSpeed Insights took: %v\n", time.Since(stepStart))
		}
	}

	report.Title = "UI/UX Analysis Report for " + strings.Split(report.URL, "://")[1]

	// Log total time taken for report generation.
	log.Printf("Total report generation time: %v\n", time.Since(startTime))

	// Return the generated report.
	return &report, nil
}
