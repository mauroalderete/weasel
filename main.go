//go:build !test

package main

import (
	"fmt"

	_ "embed"
)

//go:generate build/get_version.sh
//go:embed version.txt
var version string

func main() {
	fmt.Println("Hi world!")
}
