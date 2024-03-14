// Set the path configurations of the project

package config

import (
	"os"
	"path/filepath"
)

type pathStruct struct {
	HomeDir    string // User home directory
	WorkingDir string // Where the code is executed from
	GitRoot    string // If a git repo set it, else blank
}

var Path pathStruct

func SetPath() {
	// Set the user's home directory
	Path.HomeDir, _ = os.UserHomeDir()

	// Set the working directory
	// This is the directory from where the code is executed
	workingDir, _ := os.Getwd()
	Path.WorkingDir = workingDir

	// Set the Git root directory
	// This is a placeholder for where you would check if the current directory is a Git repo
	// and set the GitRoot accordingly. For simplicity, we'll leave it blank.

	Path.GitRoot = getGitRoot()

}

func getGitRoot() string {
	wd := Path.WorkingDir

	// Now check recursively by appending .git
	// If the lookup reaches the root directory
	// where the current directory and the previous directory are
	// the same and .git is not found then bail out with an error.
	for {
		gitDir := filepath.Join(wd, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return ""
		}
		wd = parent
	}
}
