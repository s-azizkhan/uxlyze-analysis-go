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

func Generate(url string) (*types.Report, error) {
	log.Println("Starting report generation for", url)
	startTime := time.Now()

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	err := chromedp.Run(ctx, chromedp.Navigate(url))
	if err != nil {
		return nil, err
	}

	var report types.Report
	report.URL = url
	report.Screenshots = make(map[string]string)

	report.Summary, _ = generateSummary(ctx)
	report.VisualHierarchy, _ = analysis.AnalyzeVisualHierarchy(ctx)
	report.Navigation, _ = analysis.AnalyzeNavigation(ctx)
	report.MobileFriendliness, _ = analysis.AnalyzeMobileFriendliness(ctx)
	report.Readability, _ = analysis.AnalyzeReadability(ctx)

	report.Screenshots["VisualHierarchy"], _ = screenshot.Capture(ctx, "body")
	report.Screenshots["Navigation"], _ = screenshot.Capture(ctx, "nav")

	_ = chromedp.Run(ctx, chromedp.EmulateViewport(375, 812, chromedp.EmulateScale(2.0)))
	report.Screenshots["MobileFriendliness"], _ = screenshot.Capture(ctx, "body")
	_ = chromedp.Run(ctx, chromedp.EmulateViewport(0, 0))

	report.Screenshots["Readability"], _ = screenshot.Capture(ctx, "main")
	if report.Screenshots["Readability"] == "" {
		report.Screenshots["Readability"], _ = screenshot.Capture(ctx, "body")
	}

	log.Printf("Total report generation time: %v\n", time.Since(startTime))
	return &report, nil
}

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
