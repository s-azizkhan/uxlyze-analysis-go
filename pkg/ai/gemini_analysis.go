package ai

import (
	"context"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"uxlyze/analyzer/pkg/types"

	"encoding/json"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func uploadToGemini(ctx context.Context, client *genai.Client, path, mimeType string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	options := genai.UploadFileOptions{
		DisplayName: path,
		MIMEType:    mimeType,
	}
	fileData, err := client.UploadFile(ctx, "", file, &options)
	if err != nil {
		log.Fatalf("Error uploading file: %v", err)
	}

	log.Printf("Uploaded file %s as: %s", fileData.DisplayName, fileData.URI)
	return fileData.URI
}

func AnalyzeUXWithGemini(imagePath string) (*types.GeminiUXAnalysisResult, error) {

	ctx := context.Background()
	// check if the img not exist then return
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		fmt.Println("File does not exist.")
		return nil, err
	} else if err != nil {
		fmt.Println("Error checking file:", err)
		return nil, err
	}
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatalln("Environment variable GEMINI_API_KEY not set")
		return nil, fmt.Errorf("GEMINI_API_KEY is not set: %s", apiKey)
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
		return nil, fmt.Errorf("Error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	model.SetTemperature(1)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = types.GeminiResponseSchema

	// TODO Make these files available on the local file system
	// You may need to update the file paths
	ext := filepath.Ext(imagePath)

	// Get the MIME type based on the file extension
	mimeType := mime.TypeByExtension(ext)
	fileURIs := []string{
		uploadToGemini(ctx, client, imagePath, mimeType),
	}

	session := model.StartChat()
	session.History = []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.FileData{URI: fileURIs[0]},
				genai.Text("Conduct a thorough UI & UX analysis of this website, focusing on delivering a highly detailed, insightful, crisp, straight-forward and actionable report.\n\nEvaluate key elements such as:\n\nUsability: Identify pain points in user interaction, navigation flow, and ease of use. Provide specific suggestions to streamline user paths and reduce friction.\nVisual Design: Analyze the overall aesthetics, consistency in color schemes, whitespace usage, visual hierarchy, and alignment. Suggest improvements for visual clarity, brand coherence, and engagement.\nTypography: Assess the legibility and consistency of font sizes, styles, and hierarchy. Recommend adjustments to improve readability and ensure consistent use of typography across devices and screen sizes.\nButton & CTA Design: Examine button designs, including size, color contrast, hover effects, and clarity of calls to action (CTAs). Suggest improvements for making CTAs more intuitive and visually prominent.\nNavigation: Analyze the structure and intuitiveness of the navigation menu, dropdowns, and any breadcrumb systems. Provide suggestions for improving discoverability and reducing user effort.\nAccessibility: Assess the website's accessibility for users with disabilities, including color contrast ratios, alt text, keyboard navigation, and screen reader compatibility. Provide specific fixes to ensure ADA/WCAG compliance.\nMobile Responsiveness: Evaluate how well the design adapts to various screen sizes and devices. Identify any layout issues or usability challenges on mobile and recommend improvements for a seamless experience.\nUser Flow & Information Architecture: Analyze how users move through the site, from entry points to conversion or exit. Identify any confusing steps, bottlenecks, or redundant actions, and suggest ways to optimize the flow to improve conversion rates.\nInteractivity & Feedback: Evaluate interactive elements such as forms, sliders, and hover effects. Provide suggestions for enhancing user feedback through animations, transitions, and micro-interactions to make the site feel more responsive and engaging.\nDeliver concrete, actionable suggestions for each issue found, focusing on practical improvements that can enhance the user experience, overall engagement, and conversion potential of the site, response should be very short & straight forward with minimum words."),
			},
		},
	}

	resp, err := session.SendMessage(ctx, genai.Text("INSERT_INPUT_HERE"))
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
		return nil, fmt.Errorf("Error sending message: %v", err)
	}

	// Parse the response into GeminiUXAnalysisResult
	var result types.GeminiUXAnalysisResult
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		jsonData, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
		if !ok {
			return nil, fmt.Errorf("Unexpected response format")
		}

		err = json.Unmarshal([]byte(jsonData), &result)
		if err != nil {
			return nil, fmt.Errorf("Error parsing JSON response: %v", err)
		}

		// Store the result in a JSON file
		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("Error marshaling result to JSON: %v", err)
		}

		err = os.WriteFile("gemini_result.json", resultJSON, 0644)
		if err != nil {
			return nil, fmt.Errorf("Error writing result to file: %v", err)
		}

		fmt.Println("Analysis result saved to gemini_result.json")
	} else {
		return nil, fmt.Errorf("No valid response from Gemini")
	}

	return &result, nil
}
