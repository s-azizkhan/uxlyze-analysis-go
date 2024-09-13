package analysis

import (
	"context"

	"github.com/chromedp/chromedp"
)

func AnalyzeNavigation(ctx context.Context) (map[string]interface{}, error) {
	var result map[string]interface{}

	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
		(function() {
			// Select all anchor (<a>) elements on the page
			const allLinks = document.querySelectorAll('a');
			const navElements = document.querySelectorAll('nav').length || 0;
			const linksWithoutHref = document.querySelectorAll("a[href='#']").length || 0;
			const linksWithoutTarget = document.querySelectorAll("a[target='_blank']").length || 0;

			const linkStructure = {
				internalLinks: [],
				externalLinks: [],
			};

			// Helper function to check if a link is internal or external
			function isInternalLink(url) {
				return url.includes(window.location.hostname);
			}

			allLinks.forEach(link => {
				const href = link.getAttribute('href') || ''; // Get the href attribute
				const linkText = link.textContent.trim(); // Get the link text
				const absoluteLink = link.href; // Get the absolute URL of the link
				const isInternal = isInternalLink(absoluteLink); // Check if the link is internal

				// Build link info object
				const linkInfo = {
					href,
					absoluteLink,
					text: linkText.substring(0, 100), // Limit to the first 100 characters
					isInternal,
					isAbsolute: href.startsWith('http') || href.startsWith('https')
				};

				// Categorize link as internal or external
				if (isInternal) {
					linkStructure.internalLinks.push(linkInfo);
				} else {
					linkStructure.externalLinks.push(linkInfo);
				}
			});

			const result = {
				totalLinks: linkStructure.internalLinks.length + linkStructure.externalLinks.length,
				internalLinksCount: linkStructure.internalLinks.length,
				externalLinksCount: linkStructure.externalLinks.length,
				linkStructure,
				navElementCount: navElements,
				linksWithoutHref,
				linksWithTargetBlank: linksWithoutTarget
			};

			return result;
		})()
		`, &result),
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func AnalyzeMobileFriendly(ctx context.Context) (bool, error) {
	var viewportContent string
	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`document.querySelector("meta[name='viewport']")?.getAttribute("content") || ""`, &viewportContent),
	)

	if err != nil {
		return false, err
	}

	if viewportContent != "" {
		return true, nil
	}
	return false, nil
}
