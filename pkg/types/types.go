package types

type Report struct {
	Title              string
	URL                string
	Summary            string
	VisualHierarchy    string
	Navigation         string
	MobileFriendliness string
	Readability        string
	Screenshots        map[string]string
}

type SectionAnalysis struct {
	Name         string
	FontSizes    map[string]int
	ButtonStyles map[string]string
	ColorScheme  map[string]string
	Score        int
	Details      string
}

type PageSpeedInsights struct {
	LighthouseResult struct {
		Categories struct {
			Performance struct {
				Score float64 `json:"score"`
			} `json:"performance"`
			Accessibility struct {
				Score float64 `json:"score"`
			} `json:"accessibility"`
			BestPractices struct {
				Score float64 `json:"score"`
			} `json:"best-practices"`
			SEO struct {
				Score float64 `json:"score"`
			} `json:"seo"`
		} `json:"categories"`
		Audits map[string]struct {
			Score        float64 `json:"score"`
			Title        string  `json:"title"`
			DisplayValue string  `json:"displayValue"`
		} `json:"audits"`
	} `json:"lighthouseResult"`
	LoadingExperience struct {
		OverallCategory string `json:"overall_category"`
		Metrics         map[string]struct {
			Percentile int    `json:"percentile"`
			Category   string `json:"category"`
		} `json:"metrics"`
	} `json:"loadingExperience"`
}
