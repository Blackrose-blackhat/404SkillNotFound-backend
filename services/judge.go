package services

import (
	"fmt"
	"strings"

	"github.com/Blackrose-blackhat/404SkillNotFound/internal/parser"
)

func BuildPrompt(resumeText string, github *parser.GithubProfile, roast bool) string {
	var b strings.Builder

	b.WriteString("You're an AI roast master. Your job is to analyze a developer's resume and GitHub profile, and generate a fake Twitter profile that sarcastically reflects their skillset, achievements (or lack thereof), and career ambitions. Your tone should be witty, slightly savage, yet humorous — think parody, not insult.\n\n")

	b.WriteString("### 📄 Resume:\n")
	if resumeText == "" {
		b.WriteString("No resume provided.\n")
	} else {
		b.WriteString(resumeText + "\n")
	}

	b.WriteString("\n### 📦 GitHub Summary:\n")
	if github == nil {
		b.WriteString("No GitHub profile available.\n")
	} else {
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
	}

	b.WriteString("\n---\n\n")
	b.WriteString("Now, respond ONLY with a valid JSON object in this exact structure (no markdown, no explanation):\n\n")
	b.WriteString(`{
  "handle": "@wannabeDev420",
  "display_name": "Musharraf: Bug Breeder 🪲",
  "verified_reason": "Verified for consistently shipping broken builds.",
  "bio": "My GitHub has more repos than stars. Claims 'full-stack' but mostly googles. Still waiting for a job that isn't 'remote, unpaid, blockchain-for-good'.",
  "location": "Remote (Mom’s Basement)",
  "followers": 42,
  "following": 1337,
  "joined": "Joined Oct 2021",
  "profile_image": "https://robohash.org/wannabeDev420.png",
  "banner_caption": "Currently building something no one asked for",
  "pinned_tweet": "My GitHub handle is 'Blackrose-blackhat' but my repos have 0 stars. Currently coding my next *undescribed* project.",
  "pinned_skills": [
    "#NextJS enthusiast (still figuring out dynamic routing 🤷‍♂️)",
    "#TypeScript lover, any  questions?",
    "#Web3 builder, token balance: ₹0.69"
  ],
  "tweets": [
    "My 51 GitHub repos have more 404s in their descriptions than stars. Progress?",
    "Just finished another 'AI-powered, blockchain-first, Next.js' project. Now just waiting for the 'users' to appear.",
    "Proudly using Web3.js to connect to smart contracts. My crypto balance, however, remains stubbornly dumb.",
    "Claims 'full-stack' on resume. In reality, I'm the guy who breaks both the front *and* back end.",
    "Got my Networking Essentials cert. My personal network still consists of my router and a half-eaten pizza."
  ],
  "endorsements": [
    "‘Once debugged for 6 hours only to realize it was a missing semicolon.’ – Anonymous Reviewer",
    "‘Makes up for broken code with even more broken comments.’ – Past Teammate"
  ],
  "spaces_hosted": [
    "🔥 How to Deploy Crashes at Scale",
    "💸 Monetizing Zero-Star Repos",
    "😤 Burnout or Just Bad at Time Management?"
  ],
  "replies": [
    "@bugHunter92: Bro really built StackOverflow 2.0 – it's just his own error logs.",
    "@mom: Please clean your room before starting another side project.",
    "@GPT4: Sorry, I can’t debug this either."
  ],
  "latest_tweet_stats": {
    "likes": 1,
    "retweets": 0,
    "replies": 21,
    "ratioed": true
  },
  "score": 38
}
`)

	return b.String()
}

func nonEmpty(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}
