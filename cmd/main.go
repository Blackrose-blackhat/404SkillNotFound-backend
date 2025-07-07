package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Blackrose-blackhat/404SkillNotFound/backend/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("⚠️  No .env file found at project root. Using system env vars.")
	}

	// Check if GEMINI_API_KEY is loaded (fail fast if missing)
	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Fatal("❌ GEMINI_API_KEY not set in environment")
	}

	// Define HTTP route
	http.HandleFunc("/analyze", handlers.AnalyzeHandler)

	// Start server
	port := "3000"
	log.Printf("🚀 Server running at http://localhost:%s/analyze\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
