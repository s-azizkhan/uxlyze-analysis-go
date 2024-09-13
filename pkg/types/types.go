package types

type Report struct {
	Title              string
	URL                string
	Description        string
	Navigation         string
	MobileFriendliness string
	Readability        string
	Screenshots        map[string]string
	ColorUsage         map[string]interface{}
	FontUsage          map[string]interface{}
	ImageUsage         map[string]interface{}
	GeminiAnalysis     *GeminiUXAnalysisResult `json:"geminiAnalysis,omitempty"`
	PageSpeedInsights  *PageSpeedInsights      `json:"pageSpeedInsights,omitempty"`
}

type SectionAnalysis struct {
	Name        string
	FontSizes   map[string]int
	CtaStyles   map[string]string
	ColorScheme map[string]string
	Score       int
	Details     string
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

// GeminiUXAnalysisResult represents the complete analysis result, including scores and category-specific analysis.
type GeminiUXAnalysisResult struct {
	TotalScore           float64     `json:"total_score"`            // Overall score based on the entire analysis
	WebsiteCategory      string      `json:"website_category"`       // Category of the website (e.g., e-commerce, blog)
	WebsiteCategoryScore float64     `json:"website_category_score"` // Score for how well the website fits its category
	ColorScheme          ColorScheme `json:"color_scheme"`           // Color scheme information (primary, secondary, accent colors)

	Usability     CategoryAnalysis `json:"usability"`     // Usability category analysis with score
	VisualDesign  CategoryAnalysis `json:"visual_design"` // Visual design analysis with score
	Typography    CategoryAnalysis `json:"typography"`    // Typography analysis with score
	CtaDesign     CategoryAnalysis `json:"cta_design"`    // Button & CTA design analysis with score
	Navigation    CategoryAnalysis `json:"navigation"`    // Navigation analysis with score
	Accessibility CategoryAnalysis `json:"accessibility"` // Accessibility analysis with score
	UserFlow      CategoryAnalysis `json:"user_flow"`     // User flow & information architecture analysis with score
	Interactivity CategoryAnalysis `json:"interactivity"` // Interactivity & feedback analysis with score
}

// ColorScheme represents the color palette used in the website.
type ColorScheme struct {
	PrimaryColors   []string `json:"primary_colors"`   // Primary colors (hex or color names)
	SecondaryColors []string `json:"secondary_colors"` // Secondary colors (hex or color names)
	AccentColors    []string `json:"accent_colors"`    // Accent colors (hex or color names)
}

// CategoryAnalysis represents the analysis for a specific category, including score, issues, and suggestions.
type CategoryAnalysis struct {
	Score       float64      `json:"score"`       // Score for the specific category
	Issues      []Issue      `json:"issues"`      // Issues found during the analysis
	Suggestions []Suggestion `json:"suggestions"` // Suggestions for improvement
}

// Issue represents individual issues found during the analysis.
type Issue struct {
	Description string `json:"description"` // Description of the issue
	Location    string `json:"location"`    // Location of the issue on the page
	Impact      string `json:"impact"`      // Impact or severity of the issue
}

// Suggestion represents actionable suggestions to improve the issues found.
type Suggestion struct {
	Description    string `json:"description"`     // Description of the suggestion
	ExpectedImpact string `json:"expected_impact"` // Expected impact of the suggestion
}
