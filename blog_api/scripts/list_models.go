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
// 	ctx := context.Background()
// 	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Close()

// 	models, err := client.ListModels(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Available Models:")
// 	for _, m := range models {
// 		fmt.Printf("Name: %s\n", m.Name)
// 		fmt.Printf("Supports GenerateContent: %v\n", 
// 			contains(m.SupportedGenerationMethods, "generateContent"))
// 		fmt.Println("------")
// 	}
// }

// func contains(methods []string, target string) bool {
// 	for _, m := range methods {
// 		if m == target {
// 			return true
// 		}
// 	}
// 	return false
// }
