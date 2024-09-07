package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"uxlyze/analyzer/pkg/types"
)

func Save(report *types.Report, filename string, psi *types.PageSpeedInsights) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	htmlContent := generateHTMLContent(report, psi)

	_, err = file.WriteString(htmlContent)
	return err
}

func generateHTMLContent(report *types.Report, psi *types.PageSpeedInsights) string {
	// HTML content generation code here
	htmlContent := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>UI/UX Analysis Report</title>
		<style>
			body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 800px; margin: 0 auto; padding: 20px; }
			h1, h2 { color: #2c3e50; }
			.card { background: #f9f9f9; border-radius: 5px; padding: 20px; margin-bottom: 20px; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
			img { max-width: 100%; height: auto; display: block; margin: 10px 0; }
		</style>
	</head>
	<body>
		<h1>UI/UX Analysis Report</h1>
		<div class="card">
			<h2>Summary</h2>
			<p>` + report.Summary + `</p>
		</div>
		<div class="card">
			<h2>Visual Hierarchy</h2>
			<p>` + report.VisualHierarchy + `</p>
			<img src="data:image/png;base64,` + report.Screenshots["VisualHierarchy"] + `" alt="Visual Hierarchy Screenshot">
		</div>
		<div class="card">
			<h2>Navigation</h2>
			<p>` + report.Navigation + `</p>
			<img src="data:image/png;base64,` + report.Screenshots["Navigation"] + `" alt="Navigation Screenshot">
		</div>
		<div class="card">
			<h2>Mobile Friendliness</h2>
			<p>` + report.MobileFriendliness + `</p>
			<img src="data:image/png;base64,` + report.Screenshots["MobileFriendliness"] + `" alt="Mobile Friendliness Screenshot">
		</div>
		<div class="card">
			<h2>Readability</h2>
			<p>` + report.Readability + `</p>
			<img src="data:image/png;base64,` + report.Screenshots["Readability"] + `" alt="Readability Screenshot">
		</div>`

	if psi != nil {
		htmlContent += fmt.Sprintf(`
		<div class="card">
			<h2>PageSpeed Insights</h2>
			<p>Performance Score: %.2f</p>
			<p>Accessibility Score: %.2f</p>
			<p>Best Practices Score: %.2f</p>
			<p>SEO Score: %.2f</p>
		</div>`,
			psi.LighthouseResult.Categories.Performance.Score*100,
			psi.LighthouseResult.Categories.Accessibility.Score*100,
			psi.LighthouseResult.Categories.BestPractices.Score*100,
			psi.LighthouseResult.Categories.SEO.Score*100)
	}

	htmlContent += `
	</body>
	</html>`

	return htmlContent
}

func GetPageSpeedInsights(url string) (*types.PageSpeedInsights, error) {
	apiKey := os.Getenv("GOOGLE_PAGE_SPEED_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_PAGE_SPEED_API_KEY is not set")
	}
	apiURL := fmt.Sprintf("https://www.googleapis.com/pagespeedonline/v5/runPagespeed?url=%s&key=%s", url, apiKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var psi types.PageSpeedInsights
	err = json.Unmarshal(body, &psi)
	if err != nil {
		return nil, err
	}

	return &psi, nil
}
