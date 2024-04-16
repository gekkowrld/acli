package git

import (
	"fmt"
	"github.com/gekkowrld/acli/src/config"
	git "github.com/go-git/go-git/v5"
	"path/filepath"
	"regexp"
	"strings"
)

type gi struct {
	RepoName string
}

var GitInfo gi

func IsGitRepo() bool {
	repoDir := config.Path.GitRoot

	_, err := git.PlainOpen(repoDir)
	if err != nil {
		fmt.Printf("%v\n", err);
		return false
	}

	return true
}

func SetGitInfo() {
	GitInfo.RepoName = getGitRepoName()
}

func getGitRepoName() string {
	repoPath := config.Path.GitRoot
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		fmt.Printf("Error opening repository: %v\n", err)
		return ""
	}

	cfg, err := r.Config()
	if err != nil {
		fmt.Printf("Error getting repository config: %v\n", err)
		return ""
	}

	var remoteURL string
	for _, remote := range cfg.Remotes {
		if remote.Name == "origin" {
			remoteURL = remote.URLs[0]
			break
		}
	}

	if remoteURL == "" {
		cwd := filepath.Base(config.Path.GitRoot)
		return cwd
	}

	repoName, err := extractRepoName(remoteURL)
	if err != nil {
		fmt.Printf("Error extracting repository name: %v\n", err)
		return ""
	}

	return repoName
}

func extractRepoName(url string) (string, error) {
	// Pattern to match the repository name in the URL
	pattern := regexp.MustCompile(`(?:https://github\.com/|git@github\.com:)(.*?)(?:\.git)`)
	matches := pattern.FindStringSubmatch(url)

	if len(matches) < 2 {
		return "", fmt.Errorf("could not extract repository name from URL")
	}
	reponame := matches[1]
	reponame = strings.Split(reponame, "/")[1]

	// The repository name is the last part of the match
	return reponame, nil
}
