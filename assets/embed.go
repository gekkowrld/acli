package assets

import (
	"embed"
	"fmt"
	"os"
)

// The embed "comment" is a special kind of comment
// https://pkg.go.dev/embed@go1.22.1
// The path are relative
// https://github.com/golang/go/issues/51987

//go:embed alx-low_level_programming.yaml

var structure embed.FS

// Lookup and retrieve the relevant information of a file.
func AssetsString(name string) string {
	file, err := structure.ReadFile(fmt.Sprintf("%s.yaml", name))
	if err != nil {
		fmt.Printf("%s: Sorry, the repo that you are requesting is not yet supported\n", name)
		os.Exit(1)
	}
	return string(file)
}
