package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/Blackrose-blackhat/404SkillNotFound/internal/handlers"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, using system environment variables.")
	}

	// Check critical env vars
	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Fatal("❌ Missing GEMINI_API_KEY in environment")
	}

	// Define routes
	http.HandleFunc("/api/analyze", handlers.AnalyzeHandler)

	// Add root handler
	http.HandleFunc("/", handlers.RootHandler)

	// Get port from env or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Start server
	log.Printf("🚀 AI Roast Judge backend running at http://localhost:%s/api/analyze\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
