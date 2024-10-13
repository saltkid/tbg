package cmd

import (
	"fmt"
)

var TbgVersion = "dev"

func VersionExecute() error {
	fmt.Println(TbgVersion)
	return nil
}
