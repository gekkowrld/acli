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
	"github.com/spf13/cobra"
    "github.com/gekkowrld/acli/src/git"
    err "github.com/gekkowrld/acli/src/errors"
    "github.com/gekkowrld/acli/src/config"
    "fmt"
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
        if !git.IsGitRepo() {
                err.ExitError(3)
        }
	    fmt.Printf("%s",git.AuthorsList("ALX CLI"))
   },
}

func init() {
	rootCmd.AddCommand(authorCmd)
}
