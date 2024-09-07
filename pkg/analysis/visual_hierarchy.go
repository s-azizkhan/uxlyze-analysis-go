package analysis

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

func AnalyzeVisualHierarchy(ctx context.Context) (string, error) {
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
