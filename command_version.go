package main

import (
	"fmt"
)

var TbgVersion = "dev"

type VersionCommand struct{}

func (cmd *VersionCommand) Type() CommandType { return VersionCommandType }
func (cmd *VersionCommand) String() {
	fmt.Println("Version Command")
	fmt.Println("version:", TbgVersion)
}
func (cmd *VersionCommand) ValidateValue(val *string) error {
	if val == nil || *val == "" {
		return nil
	}
	return fmt.Errorf("'version' takes no args. got: '%s'", *val)
}
func (cmd *VersionCommand) ValidateFlag(f Flag) error {
	return fmt.Errorf("'version' takes no flags. got: '%s'", f.Type)
}
func (cmd *VersionCommand) ValidateSubCommand(sc Command) error {
	switch sc.Type() {
	case NoCommandType:
		return nil
	default:
		return fmt.Errorf("'version' takes no sub commands. got: '%s'", sc.Type())
	}
}
func (cmd *VersionCommand) Execute() error {
	fmt.Println(TbgVersion)
	return nil
}
