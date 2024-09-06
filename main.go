package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
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
	SectionAnalyses    []SectionAnalysis
}

// SectionAnalysis structure to store analysis of each section
type SectionAnalysis struct {
	Name         string
	FontSizes    map[string]int
	ButtonStyles map[string]string
	ColorScheme  map[string]string
	Score        int
	Details      string
}

// PageSpeedInsights structure to store API response
type PageSpeedInsights struct {
	LighthouseResult struct {
		Categories struct {
			Performance struct {
				Score float64 `json:"score"`
			} `json:"performance"`
			Accessibility struct {
				Score float64 `json:"score"`
			} `json:"accessibility"`
			BestPractices struct {
				Score float64 `json:"best-practices"`
			} `json:"best-practices"`
			SEO struct {
				Score float64 `json:"score"`
			} `json:"seo"`
		} `json:"categories"`
	} `json:"lighthouseResult"`
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

// Analyze a specific section of the webpage
func analyzeSection(ctx context.Context, selector string) (*SectionAnalysis, error) {
	analysis := &SectionAnalysis{
		Name:         selector,
		FontSizes:    make(map[string]int),
		ButtonStyles: make(map[string]string),
		ColorScheme:  make(map[string]string),
	}

	var err error

	// Analyze font sizes
	err = chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
			const elements = document.querySelectorAll('`+selector+` *');
			const fontSizes = {};
			elements.forEach(el => {
				const size = window.getComputedStyle(el).fontSize;
				if (size) fontSizes[el.tagName] = (fontSizes[el.tagName] || 0) + 1;
			});
			return fontSizes;
		`, &analysis.FontSizes),
	)
	if err != nil {
		return nil, err
	}

	// Analyze button styles
	err = chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
			const buttons = document.querySelectorAll('`+selector+` button, `+selector+` .button, `+selector+` [role="button"]');
			const styles = {};
			buttons.forEach(btn => {
				const computed = window.getComputedStyle(btn);
				styles['borderRadius'] = computed.borderRadius;
				styles['backgroundColor'] = computed.backgroundColor;
				styles['color'] = computed.color;
			});
			return styles;
		`, &analysis.ButtonStyles),
	)
	if err != nil {
		return nil, err
	}

	// Analyze color scheme
	err = chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
			const elements = document.querySelectorAll('`+selector+` *');
			const colors = {};
			elements.forEach(el => {
				const bg = window.getComputedStyle(el).backgroundColor;
				const color = window.getComputedStyle(el).color;
				if (bg && bg !== 'rgba(0, 0, 0, 0)') colors[bg] = (colors[bg] || 0) + 1;
				if (color) colors[color] = (colors[color] || 0) + 1;
			});
			return colors;
		`, &analysis.ColorScheme),
	)
	if err != nil {
		return nil, err
	}

	// Calculate score and generate details
	analysis.Score = calculateScore(analysis)
	analysis.Details = generateDetails(analysis)

	return analysis, nil
}

func calculateScore(analysis *SectionAnalysis) int {
	score := 0

	// Score based on font size variety
	if len(analysis.FontSizes) >= 3 {
		score += 20
	} else if len(analysis.FontSizes) == 2 {
		score += 10
	}

	// Score based on button styles
	if analysis.ButtonStyles["borderRadius"] != "0px" {
		score += 10
	}
	if analysis.ButtonStyles["backgroundColor"] != "" && analysis.ButtonStyles["backgroundColor"] != "transparent" {
		score += 10
	}

	// Score based on color scheme
	if len(analysis.ColorScheme) >= 3 {
		score += 20
	} else if len(analysis.ColorScheme) == 2 {
		score += 10
	}

	return score
}

func generateDetails(analysis *SectionAnalysis) string {
	details := fmt.Sprintf("Font sizes used: %v\n", analysis.FontSizes)
	details += fmt.Sprintf("Button styles: Border Radius: %s, Background Color: %s, Text Color: %s\n",
		analysis.ButtonStyles["borderRadius"],
		analysis.ButtonStyles["backgroundColor"],
		analysis.ButtonStyles["color"])
	details += fmt.Sprintf("Color scheme: %v\n", analysis.ColorScheme)
	return details
}

func getPageSpeedInsights(url, apiKey string) (*PageSpeedInsights, error) {
	apiURL := fmt.Sprintf("https://www.googleapis.com/pagespeedonline/v5/runPagespeed?url=%s&key=%s", url, apiKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var psi PageSpeedInsights
	err = json.Unmarshal(body, &psi)
	if err != nil {
		return nil, err
	}

	return &psi, nil
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

	// Analyze specific sections
	sections := []string{".home-content-wrapper", "header", "footer", "main"}
	for _, section := range sections {
		log.Printf("Analyzing section: %s\n", section)
		sectionStart := time.Now()
		analysis, err := analyzeSection(ctx, section)
		if err != nil {
			log.Printf("Error analyzing section %s: %v\n", section, err)
			continue
		}
		log.Printf("Section %s analyzed in %v\n", section, time.Since(sectionStart))
		log.Printf("Section %s score: %d\n", section, analysis.Score)
		log.Printf("Section %s details:\n%s\n", section, analysis.Details)

		// Add section analysis to the report
		report.SectionAnalyses = append(report.SectionAnalyses, *analysis)
	}

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
func saveReport(report *Report, filename string, psi *PageSpeedInsights) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>A Journey Through %s</title>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        h1, h2 {
            color: #2c3e50;
        }
        .card {
            background-color: #fff;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            padding: 20px;
            margin-bottom: 20px;
        }
        .score {
            font-size: 24px;
            font-weight: bold;
            color: #27ae60;
        }
        img {
            max-width: 100%%;
            height: auto;
            border-radius: 8px;
            margin-top: 10px;
        }
    </style>
</head>
<body>
    <h1>Exploring the Digital Landscape of %s</h1>
    
    <div class="card">
        <h2>First Impressions</h2>
        <p>%s</p>
    </div>

    <div class="card">
        <h2>The Art of Visual Hierarchy</h2>
        <p>%s</p>
        %s
    </div>

    <div class="card">
        <h2>Navigating the Digital Seas</h2>
        <p>%s</p>
        %s
    </div>

    <div class="card">
        <h2>A Mobile-First World</h2>
        <p>%s</p>
        %s
    </div>

    <div class="card">
        <h2>The Story Unfolds: Readability</h2>
        <p>%s</p>
        %s
    </div>
`, report.URL, report.URL, report.Summary, report.VisualHierarchy,
		imageTag(report.Screenshots["VisualHierarchy"], "Visual Hierarchy Screenshot"),
		report.Navigation, imageTag(report.Screenshots["Navigation"], "Navigation Screenshot"),
		report.MobileFriendliness, imageTag(report.Screenshots["MobileFriendliness"], "Mobile Friendliness Screenshot"),
		report.Readability, imageTag(report.Screenshots["Readability"], "Readability Screenshot"))

	// Add section analyses to the HTML content
	for _, analysis := range report.SectionAnalyses {
		htmlContent += fmt.Sprintf(`
    <div class="card">
        <h2>Unveiling %s</h2>
        <p class="score">Score: %d/60</p>
        <pre>%s</pre>
    </div>
		`, analysis.Name, analysis.Score, analysis.Details)
	}

	// Add PageSpeed Insights to the HTML content
	if psi != nil {
		htmlContent += fmt.Sprintf(`
    <div class="card">
        <h2>PageSpeed Insights</h2>
        <p>Performance Score: %.2f</p>
        <p>Accessibility Score: %.2f</p>
        <p>Best Practices Score: %.2f</p>
        <p>SEO Score: %.2f</p>
    </div>
        `, psi.LighthouseResult.Categories.Performance.Score*100,
			psi.LighthouseResult.Categories.Accessibility.Score*100,
			psi.LighthouseResult.Categories.BestPractices.Score*100,
			psi.LighthouseResult.Categories.SEO.Score*100)
	}

	htmlContent += `
</body>
</html>
`

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

	startTime := time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	startAlloc := m.Alloc

	report, err := generateReport(url)
	if err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	// Get PageSpeed Insights
	log.Println("Fetching PageSpeed Insights...")
	psiStart := time.Now()
	// Get PageSpeed Insights
	psi, err := getPageSpeedInsights(url, "API_KEY_HERE")
	if err != nil {
		log.Printf("Failed to get PageSpeed Insights: %v\n", err)
	} else {
		log.Printf("PageSpeed Insights fetched in %v\n", time.Since(psiStart))
	}

	// Save report to an HTML file
	err = saveReport(report, "uiux_report.html", psi)
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
