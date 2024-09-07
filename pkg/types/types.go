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
	GeminiAnalysis     *GeminiUXAnalysisResult `json:"geminiAnalysis,omitempty"`
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

// enum for screenShot mode
type ScreenshotMode string

const (
	Mobile  ScreenshotMode = "mobile"
	Desktop ScreenshotMode = "desktop"
	Both    ScreenshotMode = "both"
)

// Gemini return type
type GeminiUXAnalysisResult struct {
	Usability            CategoryAnalysis `json:"usability"`
	VisualDesign         CategoryAnalysis `json:"visual_design"`
	Typography           CategoryAnalysis `json:"typography"`
	ButtonDesign         CategoryAnalysis `json:"button_design"`
	Navigation           CategoryAnalysis `json:"navigation"`
	Accessibility        CategoryAnalysis `json:"accessibility"`
	MobileResponsiveness CategoryAnalysis `json:"mobile_responsiveness"`
	UserFlow             CategoryAnalysis `json:"user_flow"`
	Interactivity        CategoryAnalysis `json:"interactivity"`
}

// Issue represents individual issues found during the analysis.
type Issue struct {
	Description string `json:"description"`
	Location    string `json:"location"`
	Impact      string `json:"impact"`
}

// Suggestion represents actionable suggestions to improve the issues found.
type Suggestion struct {
	Description    string `json:"description"`
	ExpectedImpact string `json:"expected_impact"`
}

type CategoryAnalysis struct {
	Issues      []Issue      `json:"issues"`
	Suggestions []Suggestion `json:"suggestions"`
}
