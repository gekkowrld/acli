package git

import (
    git "github.com/go-git/go-git/v5"
    "github.com/gekkowrld/acli/src/config"
)

func IsGitRepo() bool {
        repoDir := config.Path.WorkingDir

        _, err := git.PlainOpen(repoDir)
        if err != nil {
                return false
        }

        return true
}
