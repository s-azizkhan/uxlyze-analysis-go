package analysis

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

func AnalyzeColorUsage(ctx context.Context) (map[string]interface{}, error) {
	fmt.Println("Analyzing color usage...")
	var result map[string]interface{}
	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
		(function() {
			const allElements = document.querySelectorAll('*'); // Select all elements on the page
			const colorUsage = new Set(); // To store unique colors

			allElements.forEach(el => {
				const styles = window.getComputedStyle(el);

				// Get color, background color, and border color
				const color = styles.color;
				const backgroundColor = styles.backgroundColor;
				const borderColor = styles.borderColor;

				// Add only non-transparent colors
				if (color && color !== 'rgba(0, 0, 0, 0)' && color !== 'transparent') {
					colorUsage.add(color);
				}
				if (backgroundColor && backgroundColor !== 'rgba(0, 0, 0, 0)' && backgroundColor !== 'transparent') {
					colorUsage.add(backgroundColor);
				}
				if (borderColor && borderColor !== 'rgba(0, 0, 0, 0)' && borderColor !== 'transparent') {
					colorUsage.add(borderColor);
				}
			});

			const result = {
				totalColors: colorUsage.size,
				colors: Array.from(colorUsage) // Convert the Set to an Array
			};

			return result;
		})();
		`, &result),
	)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func AnalyzeFontUsage(ctx context.Context) (map[string]interface{}, error) {
	var result map[string]interface{}
	fmt.Println("Analyzing font usage...")
	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
		(function() {
				// Select only h1-h6, span, and p elements
				const allElements = document.querySelectorAll('h1, h2, h3, h4, h5, h6, span, p');
				const fontUsage = {};
				const fontSizeDistribution = {}; // To store the distribution of font sizes grouped by tag

				allElements.forEach(el => {
					const fontFamily = window.getComputedStyle(el).fontFamily;
					const fontSize = window.getComputedStyle(el).fontSize; // Get font size
					const textContent = el.textContent.trim(); // Get trimmed text content
					const tagName = el.tagName.toLowerCase(); // Get the tag name (p, h1, h2, etc.)

					// Only proceed if the element has text content
					if (fontFamily && textContent) {
						if (!fontUsage[fontFamily]) {
							fontUsage[fontFamily] = {}; // Initialize font family group
						}

						// Initialize the tag group inside the font family if it doesn't exist
						if (!fontUsage[fontFamily][tagName]) {
							fontUsage[fontFamily][tagName] = {};
						}

						const key = textContent + " || " + fontSize; // Unique key for text and font size

						// Add this element's text content and font size under the specific tag name
						if (!fontUsage[fontFamily][tagName][key]) {
							fontUsage[fontFamily][tagName][key] = {
								text: textContent.substring(0, 100), // Limit to the first 100 characters
								fontSize, // Include font size
								count: 0 // To track the number of elements with the same text and font size
							};
						}

						// Increment the count for this combination of text and font size
						fontUsage[fontFamily][tagName][key].count += 1;

						// Add to the font size distribution grouped by tag
						if (!fontSizeDistribution[tagName]) {
							fontSizeDistribution[tagName] = {};
						}

						fontSizeDistribution[tagName][fontSize] = (fontSizeDistribution[tagName][fontSize] || 0) + 1;
					}
				});

				const result = {
					totalFonts: Object.keys(fontUsage).length,
					fontsUsed: {},
					fontSizeDistribution // Include font size distribution grouped by tag in the final result
				};

				// Format result into a more usable structure
				Object.keys(fontUsage).forEach(font => {
					result.fontsUsed[font] = {};

					Object.keys(fontUsage[font]).forEach(tag => {
						result.fontsUsed[font][tag] = Object.values(fontUsage[font][tag]).map(({ text, fontSize, count }) => ({
							text,
							fontSize,
							count
						}));
					});
				});

				return result;
			})();

		`, &result),
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}


// Deprecated
func AnalyzeFontSizes(ctx context.Context) (string, error) {
	var result string

	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
			(function() {
				// Select only h1-h6, span, and p elements
				const allElements = document.querySelectorAll('h1, h2, h3, h4, h5, h6, span, p');
				const fontUsage = {};
				const fontSizeDistribution = {}; // To store the distribution of font sizes grouped by tag

				allElements.forEach(el => {
					const fontFamily = window.getComputedStyle(el).fontFamily;
					const fontSize = window.getComputedStyle(el).fontSize; // Get font size
					const textContent = el.textContent.trim(); // Get trimmed text content
					const tagName = el.tagName.toLowerCase(); // Get the tag name (p, h1, h2, etc.)

					// Only proceed if the element has text content
					if (fontFamily && textContent) {
						if (!fontUsage[fontFamily]) {
							fontUsage[fontFamily] = {}; // Initialize font family group
						}

						// Initialize the tag group inside the font family if it doesn't exist
						if (!fontUsage[fontFamily][tagName]) {
							fontUsage[fontFamily][tagName] = {};
						}

						const key = textContent + " || " + fontSize; // Unique key for text and font size

						// Add this element's text content and font size under the specific tag name
						if (!fontUsage[fontFamily][tagName][key]) {
							fontUsage[fontFamily][tagName][key] = {
								text: textContent.substring(0, 100), // Limit to the first 100 characters
								fontSize, // Include font size
								count: 0 // To track the number of elements with the same text and font size
							};
						}

						// Increment the count for this combination of text and font size
						fontUsage[fontFamily][tagName][key].count += 1;

						// Add to the font size distribution grouped by tag
						if (!fontSizeDistribution[tagName]) {
							fontSizeDistribution[tagName] = {};
						}

						fontSizeDistribution[tagName][fontSize] = (fontSizeDistribution[tagName][fontSize] || 0) + 1;
					}
				});

				const result = {
					totalFonts: Object.keys(fontUsage).length,
					fontsUsed: {},
					fontSizeDistribution // Include font size distribution grouped by tag in the final result
				};

				// Format result into a more usable structure
				Object.keys(fontUsage).forEach(font => {
					result.fontsUsed[font] = {};

					Object.keys(fontUsage[font]).forEach(tag => {
						result.fontsUsed[font][tag] = Object.values(fontUsage[font][tag]).map(({ text, fontSize, count }) => ({
							text,
							fontSize,
							count
						}));
					});
				});

				return JSON.stringify(result, null, 2);
			})();


		`, &result),
	)
	fmt.Println(result)

	if err != nil {
		return "", err
	}

	return result, nil
}
