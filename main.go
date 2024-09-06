package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

// Report structure to store UI/UX analysis
type Report struct {
	URL                string
	Summary            string
	VisualHierarchy    string
	Navigation         string
	MobileFriendliness string
	Readability        string
	Screenshots        map[string]string
}

// Generate summary of the website
func generateSummary(ctx context.Context) (string, error) {
	var title, description string

	err := chromedp.Run(ctx,
		chromedp.Title(&title),
		chromedp.EvaluateAsDevTools(`document.querySelector("meta[name='description']")?.getAttribute("content") || "No description found"`, &description),
	)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Website: %s\nDescription: %s", title, description), nil
}

// Analyze the visual hierarchy based on heading tags
func analyzeVisualHierarchy(ctx context.Context) (string, error) {
	var h1Count, h2Count, h3Count int

	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`document.querySelectorAll("h1").length || 0`, &h1Count),
		chromedp.EvaluateAsDevTools(`document.querySelectorAll("h2").length || 0`, &h2Count),
		chromedp.EvaluateAsDevTools(`document.querySelectorAll("h3").length || 0`, &h3Count),
	)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("H1: %d, H2: %d, H3: %d", h1Count, h2Count, h3Count), nil
}

// Analyze the navigation structure
func analyzeNavigation(ctx context.Context) (string, error) {
	var linkCount, navCount int

	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`document.querySelectorAll("a").length || 0`, &linkCount),
		chromedp.EvaluateAsDevTools(`document.querySelectorAll("nav").length || 0`, &navCount),
	)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Total links: %d, Navigation elements: %d", linkCount, navCount), nil
}

// Check mobile friendliness (basic check for viewport meta tag)
func analyzeMobileFriendliness(ctx context.Context) (string, error) {
	var viewportContent string
	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`document.querySelector("meta[name='viewport']")?.getAttribute("content") || ""`, &viewportContent),
	)

	if err != nil {
		return "", err
	}

	if viewportContent != "" {
		return "Mobile-friendly: Yes", nil
	}
	return "Mobile-friendly: No", nil
}

// Analyze readability based on the length of paragraphs
func analyzeReadability(ctx context.Context) (string, error) {
	var paragraphCount int
	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`document.querySelectorAll("p").length || 0`, &paragraphCount),
	)

	if err != nil {
		return "", err
	}

	if paragraphCount > 20 {
		return "Readability: High (many paragraphs)", nil
	}
	return "Readability: Low (few paragraphs)", nil
}

// Capture screenshot of a specific selector if it exists and is visible
func captureScreenshot(ctx context.Context, selector string) (string, error) {
	var buf []byte
	var nodeFound bool

	// Check if the element exists
	err := chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf(`document.querySelector('%s') !== null`, selector), &nodeFound),
	)

	if err != nil {
		return "", err
	}

	if !nodeFound {
		return "", nil // Element not found; return empty string
	}

	// Capture screenshot of the element
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Screenshot(selector, &buf, chromedp.NodeVisible, chromedp.ByQuery),
	)

	if err != nil {
		return "", err
	}

	// Convert the screenshot to base64 for embedding in the HTML report
	return base64.StdEncoding.EncodeToString(buf), nil
}

// Scrape website and generate report using headless browser
func generateReport(url string) (*Report, error) {
	log.Println("Starting report generation for", url)
	startTime := time.Now()

	// Start a new headless Chrome context
	log.Println("Creating Chrome context...")
	ctxStart := time.Now()
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	log.Printf("Chrome context created in %v\n", time.Since(ctxStart))

	// Create a timeout context
	log.Println("Setting up timeout context...")
	timeoutStart := time.Now()
	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
	defer cancel()
	log.Printf("Timeout context set up in %v\n", time.Since(timeoutStart))

	// Navigate to the URL
	log.Println("Navigating to URL...")
	navStart := time.Now()
	err := chromedp.Run(ctx, chromedp.Navigate(url))
	if err != nil {
		return nil, err
	}
	log.Printf("Navigation completed in %v\n", time.Since(navStart))

	var report Report
	report.URL = url
	report.Screenshots = make(map[string]string)

	// Generate various sections of the report
	log.Println("Generating summary...")
	summaryStart := time.Now()
	report.Summary, err = generateSummary(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("Summary generated in %v\n", time.Since(summaryStart))

	log.Println("Analyzing visual hierarchy...")
	vhStart := time.Now()
	report.VisualHierarchy, err = analyzeVisualHierarchy(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("Visual hierarchy analyzed in %v\n", time.Since(vhStart))

	log.Println("Analyzing navigation...")
	navAnalysisStart := time.Now()
	report.Navigation, err = analyzeNavigation(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("Navigation analyzed in %v\n", time.Since(navAnalysisStart))

	log.Println("Analyzing mobile friendliness...")
	mobileStart := time.Now()
	report.MobileFriendliness, err = analyzeMobileFriendliness(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("Mobile friendliness analyzed in %v\n", time.Since(mobileStart))

	log.Println("Analyzing readability...")
	readStart := time.Now()
	report.Readability, err = analyzeReadability(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("Readability analyzed in %v\n", time.Since(readStart))

	// Capture screenshots
	log.Println("Capturing Visual Hierarchy screenshot...")
	vhScreenStart := time.Now()
	report.Screenshots["VisualHierarchy"], err = captureScreenshot(ctx, "body")
	if err != nil {
		log.Printf("Failed to capture Visual Hierarchy screenshot: %v\n", err)
		report.Screenshots["VisualHierarchy"] = ""
	}
	log.Printf("Visual Hierarchy screenshot captured in %v\n", time.Since(vhScreenStart))

	log.Println("Capturing Navigation screenshot...")
	navScreenStart := time.Now()
	report.Screenshots["Navigation"], err = captureScreenshot(ctx, "nav")
	if err != nil {
		log.Printf("Failed to capture Navigation screenshot: %v\n", err)
		report.Screenshots["Navigation"] = ""
	}
	log.Printf("Navigation screenshot captured in %v\n", time.Since(navScreenStart))

	log.Println("Emulating mobile viewport...")
	mobileEmulateStart := time.Now()
	err = chromedp.Run(ctx, chromedp.EmulateViewport(375, 812, chromedp.EmulateScale(2.0)))
	if err != nil {
		log.Printf("Failed to emulate mobile viewport: %v\n", err)
	}
	log.Printf("Mobile viewport emulated in %v\n", time.Since(mobileEmulateStart))

	log.Println("Capturing Mobile Friendliness screenshot...")
	mobileScreenStart := time.Now()
	report.Screenshots["MobileFriendliness"], err = captureScreenshot(ctx, "body")
	if err != nil {
		log.Printf("Failed to capture Mobile Friendliness screenshot: %v\n", err)
		report.Screenshots["MobileFriendliness"] = ""
	}
	log.Printf("Mobile Friendliness screenshot captured in %v\n", time.Since(mobileScreenStart))

	log.Println("Resetting viewport...")
	resetViewportStart := time.Now()
	err = chromedp.Run(ctx, chromedp.EmulateViewport(0, 0))
	if err != nil {
		log.Printf("Failed to reset viewport: %v\n", err)
	}
	log.Printf("Viewport reset in %v\n", time.Since(resetViewportStart))

	log.Println("Capturing Readability screenshot...")
	readScreenStart := time.Now()
	report.Screenshots["Readability"], err = captureScreenshot(ctx, "main")
	if err != nil || report.Screenshots["Readability"] == "" {
		log.Println("Failed to capture main, trying body...")
		report.Screenshots["Readability"], err = captureScreenshot(ctx, "body")
		if err != nil {
			log.Printf("Failed to capture Readability screenshot: %v\n", err)
			report.Screenshots["Readability"] = ""
		}
	}
	log.Printf("Readability screenshot captured in %v\n", time.Since(readScreenStart))

	log.Printf("Total report generation time: %v\n", time.Since(startTime))
	return &report, nil
}

// Print the report to an HTML file
func saveReport(report *Report, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	htmlContent := fmt.Sprintf(`
<html>
<head>
    <title>UI/UX Report for %s</title>
</head>
<body>
    <h1>UI/UX Report for %s</h1>
    <hr />
    <h2>Summary</h2>
    <p>%s</p>

    <h2>Visual Hierarchy</h2>
    <p>%s</p>
    %s

    <h2>Navigation</h2>
    <p>%s</p>
    %s

    <h2>Mobile Friendliness</h2>
    <p>%s</p>
    %s

    <h2>Readability</h2>
    <p>%s</p>
    %s

</body>
</html>
`, report.URL, report.URL, report.Summary, report.VisualHierarchy,
		imageTag(report.Screenshots["VisualHierarchy"], "Visual Hierarchy Screenshot"),
		report.Navigation, imageTag(report.Screenshots["Navigation"], "Navigation Screenshot"),
		report.MobileFriendliness, imageTag(report.Screenshots["MobileFriendliness"], "Mobile Friendliness Screenshot"),
		report.Readability, imageTag(report.Screenshots["Readability"], "Readability Screenshot"))

	_, err = file.WriteString(htmlContent)
	return err
}

// Helper function to create an image tag if the base64 string is not empty
func imageTag(base64Str, altText string) string {
	if base64Str == "" {
		return "<p>Screenshot not available.</p>"
	}
	return fmt.Sprintf(`<img src="data:image/png;base64,%s" alt="%s"/>`, base64Str, altText)
}

func main() {
	// URL of the website to analyze
	url := "https://www.logicwind.com" // Replace with actual URL

	report, err := generateReport(url)
	if err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	// Save report to an HTML file
	err = saveReport(report, "uiux_report.html")
	if err != nil {
		log.Fatalf("Failed to save report: %v", err)
	}

	fmt.Println("UI/UX report generated successfully!")
}
