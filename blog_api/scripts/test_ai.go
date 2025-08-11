package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/google/generative-ai-go/genai"
// 	"google.golang.org/api/option"
// )

// func main() {
// 	// 1. Load API key from environment
// 	apiKey := os.Getenv("GEMINI_API_KEY")
// 	if apiKey == "" {
// 		log.Fatal("GEMINI_API_KEY not set in environment")
// 	}

// 	// 2. Initialize client
// 	ctx := context.Background()
// 	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
// 	if err != nil {
// 		log.Fatalf("Failed to create client: %v", err)
// 	}
// 	defer client.Close()

// 	// 3. Use the correct model name
// 	model := client.GenerativeModel("gemini-1.5-flash")  // Updated model name
// 	resp, err := model.GenerateContent(ctx, genai.Text("generate a blog post about the benefits of AI in healthcare"))
// 	if err != nil {
// 		log.Fatalf("Generation failed: %v\n"+
// 			"Troubleshooting:\n"+
// 			"1. Verify your API key at https://aistudio.google.com/\n"+
// 			"2. Check available models with: go run scripts/list_models.go", err)
// 	}

// 	// 4. Print response
// 	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
// 		fmt.Println("✅ Successful Response:")
// 		fmt.Println("----------------------")
// 		fmt.Println(resp.Candidates[0].Content.Parts[0])
// 		fmt.Println("----------------------")
// 	} else {
// 		fmt.Println("Empty response received")
// 	}
// }

