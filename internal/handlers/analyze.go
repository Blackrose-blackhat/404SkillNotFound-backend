// internal/handlers/analyze.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"strconv"
	"strings"
	"time"

	"github.com/Blackrose-blackhat/404SkillNotFound/internal/parser"
	"github.com/Blackrose-blackhat/404SkillNotFound/internal/types"
	"github.com/Blackrose-blackhat/404SkillNotFound/services"
)

var ipRequestCount = make(map[string]int)
var lastRequestTime = make(map[string]time.Time)

const requestLimit = 10
const requestWindow = time.Minute

func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Cache-Control", "no-store")

	// Require internal secret header for extra security
	// internalSecret := os.Getenv("INTERNAL_SECRET")
	// fmt.Sprintf(internalSecret);
	// if internalSecret != "" && r.Header.Get("X-Internal-Secret") != internalSecret {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
	// 	return
	// }

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Only POST requests allowed"})
		return
	}

	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	now := time.Now()
	if t, ok := lastRequestTime[clientIP]; ok && now.Sub(t) < requestWindow {
		if ipRequestCount[clientIP] >= requestLimit {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": "Rate limit exceeded. Try again later."})
			return
		}
		ipRequestCount[clientIP]++
	} else {
		ipRequestCount[clientIP] = 1
		lastRequestTime[clientIP] = now
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // Max 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid form data"})
		return
	}

	file, _, err := r.FormFile("resume")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Resume file missing"})
		return
	}
	defer file.Close()

	resumeText, err := parser.ExtractResume(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to extract resume"})
		return
	}

	githubUsername := r.FormValue("github_username")
	roastMode, _ := strconv.ParseBool(r.FormValue("roast_mode"))

	githubProfile, err := parser.FetchGithubProfile(githubUsername)
	if err != nil {
		log.Println("[⚠️ GitHub]", err)
	} else {
		log.Printf("✅ GitHub user: %s, repos: %d\n", githubUsername, githubProfile.TotalRepos)
	}

	prompt := services.BuildPrompt(resumeText, githubProfile, roastMode)
	responseText, err := services.GenerateContent(prompt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "AI generation failed"})
		return
	}

	var profile types.TwitterProfile
	if err := json.Unmarshal([]byte(responseText), &profile); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid AI output format"})
		log.Println("❌ Unmarshal error:", err)
		return
	}

	json.NewEncoder(w).Encode(profile)
}

// RootHandler responds with 'hello world' at the root path
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello world"))
}
