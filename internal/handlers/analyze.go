package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Blackrose-blackhat/404SkillNotFound/backend/internal/parser"
	"github.com/Blackrose-blackhat/404SkillNotFound/backend/internal/types"
	"github.com/Blackrose-blackhat/404SkillNotFound/backend/services"
)

func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB max memory
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, "Resume file missing", http.StatusBadRequest)
		return
	}
	defer file.Close()

	resumeText, err := parser.ExtractResume(file)
	if err != nil {
		http.Error(w, "Could not extract resume: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Resume extracted:\n", resumeText[:300])

	resumeBytes, _ := io.ReadAll(file)

	githubUsername := r.FormValue("github_username")
	roastMode, _ := strconv.ParseBool(r.FormValue("roast_mode"))

	// Optional: fetch GitHub profile
	githubProfile, err := parser.FetchGithubProfile(githubUsername)
	if err != nil {
		fmt.Println("⚠️ GitHub error:", err)
	} else {
		fmt.Printf("✅ Fetched %d repos for %s\n", githubProfile.TotalRepos, githubUsername)
	}

	// Build the prompt for Gemini
	prompt := services.BuildPrompt(resumeText, githubProfile, roastMode)
	fmt.Println("Prompt for AI model:\n", prompt)

	// 🔥 Get actual Gemini response
	genOutput, err := services.GenerateContent(prompt)
	if err != nil {
		http.Error(w, "AI generation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Raw AI Output:\n", genOutput[:500])

	// Try to parse AI output as JSON
	var response types.JudgeOutput
	if err := json.Unmarshal([]byte(genOutput), &response); err != nil {
		// AI didn't return valid JSON — fail gracefully
		fmt.Println("❌ Failed to parse AI response:", err)
		http.Error(w, "AI returned invalid output", http.StatusInternalServerError)
		return
	}

	// Respond with structured output
	fmt.Printf("✅ Analyzed resume for GitHub: %s | Roast: %v | Size: %d\n", githubUsername, roastMode, len(resumeBytes))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
