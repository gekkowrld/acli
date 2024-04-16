/*
Copyright © 2024 Gekko Wrld

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"github.com/gekkowrld/acli/assets"
	"github.com/gekkowrld/acli/src/config"
	log "github.com/gekkowrld/acli/src/errors"
	"github.com/gekkowrld/acli/src/git"
	"github.com/mattn/go-isatty"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var dispColour bool

// taskCmd represents the task command
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Create a file or a directory",
	Long:  `Create a file or a directory based on the provided "id"`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set up the path to be used later
		config.SetPath()
		git.SetGitInfo()
		if !git.IsGitRepo() {
			log.ExitError(3)
		}
		var vrt inputVar
		vrt.Repr = setRepresenatation()

		// Get the last part of the path
		// Since only Linux (or *nix) is expected, use '/' as a seperator
		currDir := config.Path.WorkingDir

		pathArr := strings.Split(currDir, "/")
		vrt.LastPart = pathArr[len(pathArr)-1]

		dirName := cmd.Flag("dir")
		fileName := cmd.Flag("file")
		vrt.InitData = cmd.Flag("no-initdata").Changed
		vrt.Overwrite = cmd.Flag("overwrite").Changed
		vrt.Readme = cmd.Flag("no-readme").Changed
		noColour := cmd.Flag("no-colour").Changed

		vrt.Colour = isColour(noColour)

		// Do something if the directory is passed
		if dirName.Changed {
			vrt.DirName = dirName.Value.String()
			handleDirCreation(vrt)
		}

		// Do something if the file is passed
		if fileName.Changed {
			vrt.FileName = fileName.Value.String()
			handleFileCreation(vrt)
		}

	},
}

type inputVar struct {
	Repr      Projects
	FileName  string
	InitData  bool
	Overwrite bool
	LastPart  string
	DirName   string
	Colour    bool
	Readme    bool
}
type Projects struct {
	Projects []Project `yaml:"projects"`
}

// Project yaml directory and file representation
type Project struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Directories []struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		ID          string `yaml:"id"`
		Files       []struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			ID          string `yaml:"id"`
			InitData    string `yaml:"init_data"`
		} `yaml:"files"`
	} `yaml:"directories"`
}

// Set the representation of the data to a struct
func setRepresenatation() Projects {
	// get the data from an embed.
	yamlData := assets.AssetsString(git.GitInfo.RepoName)

	var projects Projects
	err := yaml.Unmarshal([]byte(yamlData), &projects)
	if err != nil {
		log.Error(err, dispColour)
	}

	return projects
}

// Create the file when the user requests
func handleFileCreation(vrt inputVar) {
	// Use the lastPart as the directory name and then match it with the id that the user passed
	var found bool
	for _, project := range vrt.Repr.Projects {
		for _, dir := range project.Directories {
			if dir.Name == vrt.LastPart {
				for _, file := range dir.Files {
					if file.ID == vrt.FileName {
						pfn := file.Name
						// Create the file with the content inside
						var filedata string
						if vrt.InitData {
							filedata = file.InitData
						}

						config.CreateFile(pfn, filedata, vrt.Overwrite)
						// Add some more info into the README file
						// First check if the line '## Files' is available
						var msg string
						cont, err := os.ReadFile("README.md")
						if err != nil {
							log.Error(fmt.Sprintf("%s", err), vrt.Colour)
						}
						ct := string(cont)
						msg = ct
						if !strings.Contains(ct, "## Files") {
							msg = fmt.Sprintf("%s\n## Files\n", msg)
						}

						msg = fmt.Sprintf("%s\n- [%s](%s) - %s\n", msg, file.Name, file.Name, file.Description)
						found = true
						if !vrt.Readme {
							readmePath := filepath.Join(config.Path.WorkingDir, "README.md")
							// For readme, pass true to rewrite the previous file
							config.CreateFile(readmePath, msg, true)
						}

					}
				}
			}

			if !found {
				errMsg := fmt.Sprintf("%s: Couldn't find where the file belongs to, please go to the correct location", git.GitInfo.RepoName)
				log.Error(errMsg, vrt.Colour)
			}
		}
	}

}

func isColour(isColour bool) bool {
	// Check if the NO_COLOR environment variable is set
	noColorEnv := os.Getenv("NO_COLOR")

	// Check if the output is being redirected or if the terminal doesn't support color
	isTerminal := isatty.IsTerminal(os.Stdout.Fd())
	isRedirected := !isatty.IsTerminal(os.Stderr.Fd())

	// If the environment variable is set or the output is redirected, disable color
	if noColorEnv != "" || isRedirected {
		return false
	}

	// If the CLI flag is explicitly set to disable color, disable it
	if isColour {
		return false
	}

	// If none of the above conditions are met, enable color
	return isTerminal
}

// Create a directory when the user requests so
func handleDirCreation(vrt inputVar) {
	var found bool
	for _, project := range vrt.Repr.Projects {
		if project.Name == git.GitInfo.RepoName {
			for _, dir := range project.Directories {
				if dir.ID == vrt.DirName {
					dn := dir.Name
					config.CreateDir(dn, vrt.Overwrite)
					// Create The README.md file
					if !vrt.Readme {
						readmePath := filepath.Join(config.Path.WorkingDir, dn, "README.md")
						content := fmt.Sprintf("# %s\n\n%s\n", dir.Name, dir.Description)
						config.CreateFile(readmePath, content, true)
					}

					found = true
				}
			}
		}
		if !found {
		errMsg := fmt.Sprintf("%s: Couldn't find where the %s belongs to, please go to the correct location", config.Path.WorkingDir, git.GitInfo.RepoName)
		log.Error(errMsg, vrt.Colour)
		}
	}
}

func init() {
	rootCmd.AddCommand(taskCmd)

	taskCmd.Flags().BoolP("no-initdata", "i", false, "Don't add any data inside the file when creating it")
	taskCmd.Flags().BoolP("no-readme", "n", false, "Don't create or update the README.md file")
	taskCmd.Flags().String("dir", "", "Create a directory based on the id e.g 0x00")
	taskCmd.Flags().String("file", "", "Create a file based on the id e.g 1")
	taskCmd.Flags().BoolP("overwrite", "o", false, "Don't warn about existance of a file or directory, just overwrite it")
	taskCmd.Flags().BoolP("no-colour", "c", false, "Don't output with colours (Partial conformity to https://no-color.org/)")
}
