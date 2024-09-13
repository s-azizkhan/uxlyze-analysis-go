package analysis

import (
	"context"
	"fmt"

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
		res := fmt.Sprintf("High (many paragraphs (%d))", paragraphCount)
		return res, nil
	}
	res := fmt.Sprintf("Low (few paragraphs (%d))", paragraphCount)
	return res, nil
}
