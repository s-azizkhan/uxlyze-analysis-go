package types

import "github.com/google/generative-ai-go/genai"

var GeminiResponseSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"total_score": {
			Type:        genai.TypeNumber, // Overall score based on the analysis
			Description: "The aggregate score of the website based on all evaluated categories",
		},
		"website_category": {
			Type:        genai.TypeString, // Category of the website (e.g., e-commerce, blog)
			Description: "The classification of the website based on its visible structure and design",
		},
		"website_category_score": {
			Type:        genai.TypeNumber, // Score for how well the website fits its category
			Description: "Score for the appropriateness and execution of the website's design according to its category",
		},
		"color_scheme": {
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"primary_colors":   createColorArraySchema(), // Array of primary colors
				"secondary_colors": createColorArraySchema(), // Array of secondary colors
				"accent_colors":    createColorArraySchema(), // Array of accent colors
			},
			Description: "Color scheme extracted from the website, including primary, secondary, and accent colors",
		},
		"usability":     createCategorySchema("Analyzes the organization of elements and user paths in the visible layout"),
		"visual_design": createCategorySchema("Evaluates visual aesthetics, color schemes, and brand coherence"),
		"typography":    createCategorySchema("Assesses font legibility, sizing, and consistency"),
		"cta_design":    createCategorySchema("Examines the design, prominence, and effectiveness of CTAs"),
		"navigation":    createCategorySchema("Reviews visible navigation structure and intuitiveness"),
		"accessibility": createCategorySchema("Identifies potential accessibility issues based on visual elements"),
		"user_flow":     createCategorySchema("Analyzes content organization and flow of visible information"),
		"interactivity": createCategorySchema("Evaluates visual cues for interactivity, such as CTA, buttons or forms"),
	},
}

// Helper function for color schema
func createColorArraySchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeString, // Hex codes or color names
		},
		Description: "Array of colors (as hex codes or color names)",
	}
}

func createCategorySchema(description string) *genai.Schema {
	return &genai.Schema{
		Type:        genai.TypeObject,
		Description: description,
		Properties: map[string]*genai.Schema{
			"score": {
				Type:        genai.TypeNumber, // Score for this specific category
				Description: "Score representing the performance of the website in this category",
			},
			"issues":      createIssuesSchema(),      // Issues found based on the screenshot
			"suggestions": createSuggestionsSchema(), // Actionable suggestions for improvements
		},
	}
}

func createIssuesSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"description": {Type: genai.TypeString}, // Issue description
				"location":    {Type: genai.TypeString}, // Location on the page
				"severity":    {Type: genai.TypeString}, // Severity of the issue
			},
		},
	}
}

func createSuggestionsSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"description":     {Type: genai.TypeString}, // Suggested fix or improvement
				"expected_impact": {Type: genai.TypeString}, // Impact of the suggested improvement
			},
		},
	}
}
