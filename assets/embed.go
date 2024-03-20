package assets

import (
	"embed"
	"fmt"
	"os"
)

//go:embed alx-low_level_programming.json
//go:embed dir.yaml

var structure embed.FS

func AssetsString(name string) string {
	file, err := structure.ReadFile(fmt.Sprintf("%s.yaml", name))
	if err != nil {
		fmt.Printf("Sorry, the repo that you are requesting is not yet supported\n")
		os.Exit(1)
	}
	return string(file)
}
