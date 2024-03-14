// Copyright 2022 Naohiro CHIKAMATSU
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package git

import (
	"fmt"
	"github.com/gekkowrld/acli/src/config"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var osExit = os.Exit

const cmdName string = "gal"
const version string = "1.2.1"

type user struct {
	name string
	mail string
}

const (
	exitSuccess int = iota // 0
	exitFailure
)

func AuthorsList(projectname string) string {
	header := getHeader(projectname)
	authors := getAuthors()

	return strAuthors(header, authors)
}

// getAuthors returns authors in this project.
// This method get authors name and mail address from git log.
func getAuthors() []string {
	return getAuthorsAlphabeticalOrder()
}

// getAuthorsAlphabeticalOrder returs authors name and authors mail address in alphabetical order
func getAuthorsAlphabeticalOrder() []string {
	out, err := exec.Command("git", "-C", config.Path.GitRoot, "log", "--pretty=format:%an <%ae>").Output()
	if err != nil {
		die(err.Error())
	}

	list := strings.Split(string(out), "\n")
	list = removeDuplicate(list)
	sort.Strings(list)
	return list
}
func getHeader(pn string) string {
	// Header as shown by Google:
	// opensource.google/documentation/reference/releasing/authors/
	// and Moby:
	// github.com/moby/moby
	rs := fmt.Sprintf("# This is the list of %s's contributors.\n", pn)
	rs += fmt.Sprintf("# This does not necessarily list everyone who has contributed code.\n")
	rs += fmt.Sprintf("# To see the full list of contributors, see the revision history in\n# source control.\n")
	return rs
}

// removeDuplicate removes duplicates in the slice.
func removeDuplicate(list []string) []string {
	results := make([]string, 0, len(list))
	encountered := map[string]bool{}
	for i := 0; i < len(list); i++ {
		if !encountered[list[i]] {
			encountered[list[i]] = true
			results = append(results, list[i])
		}
	}
	return results
}

// exists check whether file or directory exists.
func exists(path string) bool {
	_, err := os.Stat(path)
	return (err == nil)
}

// die exit program with message.
func die(msg string) {
	fmt.Fprintln(os.Stderr, cmdName+": "+msg)
	osExit(exitFailure)
}

// create Author File string representation
func strAuthors(header string, authors []string) string {
	var rt string
	var au string
	for _, txt := range authors {
		au += fmt.Sprintf("%s\n", txt)
	}
	rt = fmt.Sprintf("%s\n%s", header, au)

	return rt
}
