package api

import (
	"encoding/json"
	"net/http"
	"uxlyze/analyzer/pkg/report"
	"uxlyze/analyzer/pkg/types"
)

func HandleAnalyzeRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		URL                   string               `json:"url"`
		IncludeScreenshots    bool                 `json:"includeScreenshots"`
		ScreenshotMode        types.ScreenshotMode `json:"screenshotMode"`
		IncludePSI            bool                 `json:"includePSI"`
		IncludeGeminiAnalysis bool                 `json:"includeGeminiAnalysis"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	report, err := report.Generate(
		request.URL,
		request.IncludeScreenshots,
		request.ScreenshotMode,
		request.IncludePSI,
		request.IncludeGeminiAnalysis,
	)

	if err != nil {
		http.Error(w, "Error generating report: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
