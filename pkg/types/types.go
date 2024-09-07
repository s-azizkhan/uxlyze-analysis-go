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
				Score float64 `json:"best-practices"`
			} `json:"best-practices"`
			SEO struct {
				Score float64 `json:"score"`
			} `json:"seo"`
		} `json:"categories"`
	} `json:"lighthouseResult"`
}
