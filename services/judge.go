package services

import (
	"fmt"
	"strings"

	"github.com/Blackrose-blackhat/404SkillNotFound/internal/parser"
)

func BuildPrompt(resumeText string, github *parser.GithubProfile, roast bool) string {
	var b strings.Builder

	b.WriteString("You're an AI judge evaluating the software engineering potential of a candidate based on their resume and GitHub profile.\n")
	b.WriteString("You give structured, brutally honest feedback. Be direct, sharp, sarcastic (if roast mode), and actionable.\n\n")

	b.WriteString("### 📄 Resume:\n")
	if resumeText == "" {
		b.WriteString("No resume provided.\n")
	} else {
		b.WriteString(resumeText + "\n")
	}

	b.WriteString("\n### 📦 GitHub Summary:\n")
	b.WriteString(fmt.Sprintf("Total public repos: %d\n", github.TotalRepos))
	for i, repo := range github.TopRepos {
		b.WriteString(fmt.Sprintf("%d. %s - %s - ⭐ %d – %s\n",
			i+1,
			repo.Name,
			nonEmpty(repo.Description, "No description"),
			repo.Stargazers,
			nonEmpty(repo.Language, "Unknown language"),
		))
	}

	if roast {
		b.WriteString("\n### Mode: ROAST\n")
		b.WriteString(`Unleash your darkest roast. Use clever sarcasm, developer memes, and brutal honesty. Be Gordon Ramsay meets a senior FAANG engineer reviewing a bootcamp portfolio. But also provide real advice.
`)
	} else {
		b.WriteString("\n### Mode: CONSTRUCTIVE\n")
		b.WriteString(`Give detailed but respectful feedback. Highlight strengths and areas for improvement. Be like a good mentor — honest, actionable, helpful.
`)
	}

	// 🧾 Output schema
	b.WriteString(`
Now respond ONLY with a valid JSON object in this exact structure:

{
  "score": (integer from 0 to 100),
  "grade": {
    "letter": "A"|"B"|"C"|"D"|"F",
    "description": "One-line description like 'Excellent', 'Needs work', etc.",
    "color": "HEX code for UI display (e.g., #52c41a for A, #faad14 for C, #ff4d4f for F)"
  },
  "summary": "One paragraph roast/overview of their skillset and red flags.",
  "feedback": [
    {
      "title": "Short label for the issue",
      "detail": "One paragraph critique. Funny, sharp, and useful."
    }
  ],
  "roast_mode": true or false,
  "recommendation": {
    "title": "Final Word",
    "detail": "Short, sharp closing advice — sarcastic or supportive depending on mode."
  }
}

DO NOT include any intro or explanation. Just valid JSON. No markdown. No bullet points. No extra commentary.
`)

	return b.String()
}

func nonEmpty(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}
