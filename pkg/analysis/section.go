package analysis

import (
	"context"
	"fmt"

	"uxlyze/analyzer/pkg/types"

	"github.com/chromedp/chromedp"
)

func AnalyzeSection(ctx context.Context, selector string) (*types.SectionAnalysis, error) {
	analysis := &types.SectionAnalysis{
		Name:        selector,
		FontSizes:   make(map[string]int),
		CtaStyles:   make(map[string]string),
		ColorScheme: make(map[string]string),
	}

	// Analyze font sizes
	err := chromedp.Run(ctx,
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
		`, &analysis.CtaStyles),
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

func calculateScore(analysis *types.SectionAnalysis) int {
	score := 0

	if len(analysis.FontSizes) >= 3 {
		score += 20
	} else if len(analysis.FontSizes) == 2 {
		score += 10
	}

	if analysis.CtaStyles["borderRadius"] != "0px" {
		score += 10
	}
	if analysis.CtaStyles["backgroundColor"] != "" && analysis.CtaStyles["backgroundColor"] != "transparent" {
		score += 10
	}

	if len(analysis.ColorScheme) >= 3 {
		score += 20
	} else if len(analysis.ColorScheme) == 2 {
		score += 10
	}

	return score
}

func generateDetails(analysis *types.SectionAnalysis) string {
	details := fmt.Sprintf("Font sizes used: %v\n", analysis.FontSizes)
	details += fmt.Sprintf("Button styles: Border Radius: %s, Background Color: %s, Text Color: %s\n",
		analysis.CtaStyles["borderRadius"],
		analysis.CtaStyles["backgroundColor"],
		analysis.CtaStyles["color"])
	details += fmt.Sprintf("Color scheme: %v\n", analysis.ColorScheme)
	return details
}
