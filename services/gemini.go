package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

// GenerateContent streams Gemini output and returns it as one string.
func GenerateContent(prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY") // or GOOGLE_API_KEY, either works
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI, // use VertexAI backend if that’s your target
	})
	if err != nil {
		return "", fmt.Errorf("create genai client: %w", err)
	}

	var sb strings.Builder
	for resp, err := range client.Models.GenerateContentStream(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt), // use the caller-supplied prompt
		nil,                // optional *genai.GenerateContentConfig
	) {
		if err != nil {
			return "", fmt.Errorf("stream read: %w", err)
		}
		for _, cand := range resp.Candidates {
			for _, part := range cand.Content.Parts {
				sb.WriteString(part.Text)
			}
		}
	}

	// Clean up possible markdown code block wrappers and language tags
	result := strings.TrimSpace(sb.String())
	if strings.HasPrefix(result, "```") {
		result = strings.TrimPrefix(result, "```")
		result = strings.TrimSpace(result)
	}
	if strings.HasPrefix(result, "json") {
		result = strings.TrimPrefix(result, "json")
		result = strings.TrimSpace(result)
	}
	if strings.HasSuffix(result, "```") {
		result = strings.TrimSuffix(result, "```")
		result = strings.TrimSpace(result)
	}

	var response interface{}
	err = json.Unmarshal([]byte(result), &response)
	if err != nil {
		return "", fmt.Errorf("unmarshal json: %w", err)
	}

	return result, nil
}
