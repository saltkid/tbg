package cmd

import (
	"fmt"
)

var version string

func VersionExecute() error {
	fmt.Println(version)
	return nil
}
