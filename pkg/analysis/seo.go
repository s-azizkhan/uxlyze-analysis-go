package analysis

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

func AnalyzeSEO(ctx context.Context) (map[string]interface{}, error) {
	var result map[string]interface{}
	fmt.Println("Analyzing SEO...")
	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
		(function() {
			const metaTags = document.querySelectorAll('meta');
			const metaInfo = {};

			metaTags.forEach(tag => {
				const name = tag.getAttribute('name') || tag.getAttribute('property');
				const content = tag.getAttribute('content');
				if (name && content) {
					metaInfo[name] = content;
				}
			});

			return metaInfo;
		})()
		`, &result),
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
