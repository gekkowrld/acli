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
	"github.com/gekkowrld/acli/src/config"
	err "github.com/gekkowrld/acli/src/errors"
	"github.com/gekkowrld/acli/src/git"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

// authorCmd represents the author command
var authorCmd = &cobra.Command{
	Use:   "author",
	Short: "Generate a list of contributors in a project",
	Long: `Generate a list of contributors in a project.

The list is generated according to git contributing history.
Generation of the list is case insensitive and ordered Alphabetically`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set the path(s) before running anything
		config.SetPath()
		git.SetGitInfo()
		if !git.IsGitRepo() {
			err.ExitError(3)
		}
		toStdout := cmd.Flag("stdout").Changed
		filename := cmd.Flag("file")
		authorFile := filepath.Join(config.Path.WorkingDir, "AUTHOR")
		if filename.Changed {
			af := filename.Value.String()
			// To be extra sure but it fails if nothing is passed
			if len(af) > 0 {
				authorFile = af
			}
		}
		pn := git.GitInfo.RepoName
		pnf := cmd.Flag("name")
		if pnf.Changed {
			pn = pnf.Value.String()
		}
		fc := git.AuthorsList(pn)
		if toStdout {
			fmt.Printf("%s", fc)
		} else {
			// Save the contents to the file
			file, err := os.Create(authorFile)
			if err != nil {
				log.Fatal("%v", err)
			}
			defer file.Close()
			file.WriteString(fc)
		}
	},
}

func init() {
	rootCmd.AddCommand(authorCmd)
	authorCmd.PersistentFlags().BoolP("stdout", "o", false, "Write to the stdout (mostly the terminal) instead of a file")
	authorCmd.PersistentFlags().String("file", "$GITROOT/AUTHOR", "The file that Authors names will be written in")
	authorCmd.PersistentFlags().String("name", "REPO NAME", "The name of the project")
}
