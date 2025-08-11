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
// 	// Set your API key here or export it
// 	apiKey := os.Getenv("GEMINI_API_KEY")
// 	if apiKey == "" {
// 		log.Fatal(`ERROR: API key not found. 
		
// 		Either:
// 		1. Set it temporarily: 
// 		   export GEMINI_API_KEY="your-key-here"
		
// 		2. Or hardcode it temporarily:
// 		   apiKey := "your-key-here"`)
// 	}

// 	ctx := context.Background()
	
// 	// Try connecting to the API
// 	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
// 	if err != nil {
// 		log.Fatalf("🔴 Connection failed: %v\n\nTroubleshooting:\n1. Check key at https://aistudio.google.com/\n2. Verify no typos\n3. Check internet connection", err)
// 	}
// 	defer client.Close()

// 	// Test a simple request
// 	model := client.GenerativeModel("gemini-1.5-flash") // Latest free model
// 	resp, err := model.GenerateContent(ctx, genai.Text("Say 'API test successful'"))
// 	if err != nil {
// 		log.Fatalf("🔴 Generation failed: %v\n\nTry these model names:\n1. gemini-1.5-flash\n2. gemini-pro\n3. models/gemini-pro", err)
// 	}

// 	// Print result
// 	if len(resp.Candidates) > 0 {
// 		fmt.Println("\n✅ API Key Works!")
// 		fmt.Println("Response:", resp.Candidates[0].Content.Parts[0])
// 	} else {
// 		fmt.Println("⚠️  Received empty response")
// 	}
// }