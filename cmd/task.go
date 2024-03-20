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
	"gopkg.in/yaml.v3"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

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
		log.Fatalf("error: %v", err)
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
				fmt.Printf("Couldn't find where the file belongs to, please check again")
			}
		}
	}

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
	}
}

func init() {
	rootCmd.AddCommand(taskCmd)

	taskCmd.Flags().BoolP("no-initdata", "i", false, "Don't add any data inside the file when creating it")
	taskCmd.Flags().BoolP("no-readme", "n", false, "Don't create or update the README.md file")
	taskCmd.Flags().String("dir", "", "Create a directory based on the id e.g 0x00")
	taskCmd.Flags().String("file", "", "Create a file based on the id e.g 1")
	taskCmd.Flags().BoolP("overwrite", "o", false, "Don't warn about existance of a file or directory, just overwrite it")
}
