package analysis

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

func AnalyzeVisualHierarchy(ctx context.Context) (string, error) {
	var h1Count, h2Count, h3Count, imgCount int
	var h1FontSizes, h2FontSizes, h3FontSizes, pFontSizes map[string]int

	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`document.querySelectorAll("h1").length || 0`, &h1Count),
		chromedp.EvaluateAsDevTools(`document.querySelectorAll("h2").length || 0`, &h2Count),
		chromedp.EvaluateAsDevTools(`document.querySelectorAll("h3").length || 0`, &h3Count),
		chromedp.EvaluateAsDevTools(`document.querySelectorAll("img").length || 0`, &imgCount),
		chromedp.EvaluateAsDevTools(`
			function getFontSizes(elementType) {
				var elements = document.querySelectorAll(elementType);
				var fontSizes = {};
				elements.forEach(function(el) {
					var fontSize = window.getComputedStyle(el).fontSize;
					fontSizes[fontSize] = (fontSizes[fontSize] || 0) + 1;
				});
				return fontSizes;
			}
			[
				getFontSizes('h1'),
				getFontSizes('h2'),
				getFontSizes('h3'),
				getFontSizes('p')
			]
		`, &[]*map[string]int{&h1FontSizes, &h2FontSizes, &h3FontSizes, &pFontSizes}),
	)

	if err != nil {
		return "", err
	}

	formatFontSizes := func(count int, sizes map[string]int) string {
		if count == 0 {
			return fmt.Sprintf("%d", count)
		}
		var result string
		for size, num := range sizes {
			result += fmt.Sprintf("%s: %d, ", size, num)
		}
		return fmt.Sprintf("%d (%s)", count, result[:len(result)-2])
	}

	return fmt.Sprintf("H1: %s\nH2: %s\nH3: %s\nP: %s\nImages: %d",
		formatFontSizes(h1Count, h1FontSizes),
		formatFontSizes(h2Count, h2FontSizes),
		formatFontSizes(h3Count, h3FontSizes),
		formatFontSizes(len(pFontSizes), pFontSizes),
		imgCount), nil
}
