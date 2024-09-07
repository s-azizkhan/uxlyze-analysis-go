package types

import "github.com/google/generative-ai-go/genai"

var GeminiResponseSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"usability":             createCategorySchema(),
		"visual_design":         createCategorySchema(),
		"typography":            createCategorySchema(),
		"button_design":         createCategorySchema(),
		"navigation":            createCategorySchema(),
		"accessibility":         createCategorySchema(),
		"mobile_responsiveness": createCategorySchema(),
		"user_flow":             createCategorySchema(),
		"interactivity":         createCategorySchema(),
	},
}

func createCategorySchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"issues":      createIssuesSchema(),
			"suggestions": createSuggestionsSchema(),
		},
	}
}

func createIssuesSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"description": {Type: genai.TypeString},
				"location":    {Type: genai.TypeString},
				"impact":      {Type: genai.TypeString},
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
				"description":     {Type: genai.TypeString},
				"expected_impact": {Type: genai.TypeString},
			},
		},
	}
}
