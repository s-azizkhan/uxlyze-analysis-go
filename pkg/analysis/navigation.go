package analysis

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

func AnalyzeNavigation(ctx context.Context) (string, error) {
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

func AnalyzeMobileFriendliness(ctx context.Context) (string, error) {
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
