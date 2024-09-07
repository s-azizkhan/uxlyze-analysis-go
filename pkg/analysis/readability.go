package analysis

import (
	"context"

	"github.com/chromedp/chromedp"
)

func AnalyzeReadability(ctx context.Context) (string, error) {
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
