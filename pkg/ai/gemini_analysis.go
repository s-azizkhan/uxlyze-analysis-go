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
		log.Printf("Error opening file: %v", err)
		return ""
	}
	defer file.Close()

	options := genai.UploadFileOptions{
		DisplayName: path,
		MIMEType:    mimeType,
	}
	fileData, err := client.UploadFile(ctx, "", file, &options)
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return ""
	}

	log.Printf("Uploaded file %s as: %s", fileData.DisplayName, fileData.URI)
	return fileData.URI
}

// loadPrompt loads the prompt from the given JSON file.
func loadPrompt(filePath string) (string, error) {
	promptData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading %s: %v", filePath, err)
	}

	var promptsJson map[string]interface{}
	err = json.Unmarshal(promptData, &promptsJson)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling %s: %v", filePath, err)
	}

	geminiPrompt, ok := promptsJson["gemini"].(map[string]interface{})["v2"].(string)
	if !ok || geminiPrompt == "" {
		return "", fmt.Errorf("gemini prompt not found in %s", filePath)
	}

	return geminiPrompt, nil
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
		log.Print("Environment variable GEMINI_API_KEY not set")
		return nil, fmt.Errorf("GEMINI_API_KEY is not set: %s", apiKey)
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("Error creating client: %v", err)
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

	// Use prompts.json to get the prompt
	geminiPrompt, err := loadPrompt("pkg/ai/prompts.json")

	if err != nil {
		log.Printf("Error loading gemini prompt: %v", err)
		return nil, fmt.Errorf("Error loading gemini prompt: %v", err)
	}

	if geminiPrompt == "" {
		log.Printf("Error getting gemini prompt: %v", err)
		return nil, fmt.Errorf("Error getting gemini prompt: %v", err)
	}

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
				genai.Text(string(geminiPrompt)),
			},
		},
	}

	resp, err := session.SendMessage(ctx, genai.Text("Keep the response short and concise"))
	// resp, err := session.SendMessage(ctx)
	if err != nil {
		log.Printf("Error sending message: %v", err)
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
