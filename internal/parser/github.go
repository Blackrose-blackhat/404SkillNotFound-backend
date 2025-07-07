package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

type GithubRepo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stargazers  int    `json:"stargazers_count"`
	Language    string `json:"language"`
	PushedAt    string `json:"pushed_at"` // ISO8601 format
}

type GithubProfile struct {
	TotalRepos int
	TopRepos   []GithubRepo
}

func FetchGithubProfile(username string) (*GithubProfile, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100", username)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GitHub API error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub user not found (%d)", resp.StatusCode)
	}

	var repos []GithubRepo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("failed to decode GitHub repos: %w", err)
	}

	// Sort repos by pushed_at descending
	sort.Slice(repos, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339, repos[i].PushedAt)
		t2, _ := time.Parse(time.RFC3339, repos[j].PushedAt)
		return t1.After(t2)
	})

	// Pick top 5 most recently updated
	topN := 5
	if len(repos) < topN {
		topN = len(repos)
	}

	return &GithubProfile{
		TotalRepos: len(repos),
		TopRepos:   repos[:topN],
	}, nil
}
