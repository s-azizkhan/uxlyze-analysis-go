package api

import (
	"encoding/json"
	"net/http"
	"uxlyze/analyzer/pkg/report"
)

func HandleVersionRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"version": "1.0.0"})
}

func HandleAnalyzeRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		URL               string `json:"url"`
		IncludePreview    bool   `json:"includePreview"`
		IncludePSI        bool   `json:"includePSI"`
		IncludeAIAnalysis bool   `json:"includeAIAnalysis"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	report, err := report.Generate(
		request.URL,
		request.IncludePreview,
		request.IncludePSI,
		request.IncludeAIAnalysis,
	)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Error generating report: " + err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
