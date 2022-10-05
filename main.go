//go:build !test

package main

import (
	_ "embed"

	"github.com/mauroalderete/weasel/cmd"
)

//go:generate build/get_version.sh
//go:embed version.txt
var version string

func main() {
	cmd.Execute()
}

func init() {
	if len(version) == 0 {
		version = "beta"
	}
	cmd.SetVersion(version)
}
