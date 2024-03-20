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
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"github.com/mattn/go-isatty"

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
		repr := setRepresenatation()

		// Get the last part of the path
		// Since only Linux (or *nix) is expected, use '/' as a seperator
		currDir := config.Path.WorkingDir

		pathArr := strings.Split(currDir, "/")
		lastPart := pathArr[len(pathArr)-1]

		dirName := cmd.Flag("dir")
		fileName := cmd.Flag("file")
		noInitdata := cmd.Flag("no-initdata").Changed
		overwrite := cmd.Flag("overwrite").Changed
		noReadme := cmd.Flag("no-readme").Changed
		noColour := cmd.Flag("no-colour").Changed

		dispColour = isColour(noColour)

		// Do something if the directory is passed
		if dirName.Changed {
			handleDirCreation(repr, dirName.Value.String(), noReadme, overwrite, lastPart)
		}

		// Do something if the file is passed
		if fileName.Changed {

			handleFileCreation(repr, fileName.Value.String(), noInitdata, overwrite, lastPart)
		}
	},
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
	// Hardcode the filename for now
	yamlData := assets.AssetsString("dir")

	var projects Projects
	err := yaml.Unmarshal([]byte(yamlData), &projects)
	if err != nil {
		log.Error(err, dispColour)
	}

	return projects
}

// Create the file when the user requests
func handleFileCreation(repr Projects, fileID string, noInitdata bool, overwrite bool, lastPart string) {
	// Use the lastPart as the directory name and then match it with the id that the user passed
	for _, project := range repr.Projects {
		for _, dir := range project.Directories {
			if dir.Name == lastPart {
				for _, file := range dir.Files {
					if file.ID == fileID {
						pfn := file.Name
						// Create the file with the content inside
						var filedata string
						if !noInitdata {
							filedata = file.InitData
						}

						config.CreateFile(pfn, filedata, overwrite)

					}
				}
				break
			} else {
				errMsg := fmt.Sprintf("%s: Couldn't find where the file belongs to, please go to the correct location", config.Path.WorkingDir)
				log.Error(errMsg, dispColour)
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
func handleDirCreation(repr Projects, dirID string, noReadme bool, overwrite bool, lastPart string) {
	for _, project := range repr.Projects {
		if project.Name == lastPart {
			for _, dir := range project.Directories {
				if dir.ID == dirID {
					dn := dir.Name
					config.CreateDir(dn, overwrite)
					// Create The README.md file
					if !noReadme {
						readmePath := filepath.Join(config.Path.WorkingDir, dn, "README.md")
						content := fmt.Sprintf("# %s\n\n%s\n", dir.Name, dir.Description)
						config.CreateFile(readmePath, content, overwrite)
					}
				}
			}
		}
		errMsg := fmt.Sprintf("%s: Couldn't find where the directory belongs to, please go to the correct location", config.Path.WorkingDir)
		log.Error(errMsg, dispColour)
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
