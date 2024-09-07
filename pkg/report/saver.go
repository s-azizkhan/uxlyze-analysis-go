package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"uxlyze/analyzer/pkg/types"
)

func Save(report *types.Report, filename string, psi *types.PageSpeedInsights) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	htmlContent, err := generateHTMLContent(report, psi)
	if err != nil {
		return err
	}

	_, err = file.WriteString(htmlContent)
	return err
}

func generateHTMLContent(report *types.Report, psi *types.PageSpeedInsights) (string, error) {
	templatePath := filepath.Join("pkg", "report", "report_template.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	data := struct {
		*types.Report
		PageSpeedInsights *types.PageSpeedInsights
	}{
		Report:            report,
		PageSpeedInsights: psi,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return buf.String(), nil
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
