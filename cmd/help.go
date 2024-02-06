package cmd

import (
	"fmt"

	"github.com/saltkid/tbg/flag"
)

func HelpValidateValue(val string) error {
	switch val {
	case "":
		return nil
	default:
		return fmt.Errorf("help takes no args. got '%s'", val)
	}
}

func HelpValidateFlag(f *flag.Flag) error {
	return nil // accept all flags
}

func HelpValidateSubCmd(c *Cmd) error {
	return nil // accept all subcommands
}

func HelpExecute(c *Cmd) error {
	subCmds := make(map[CmdType]struct{}, 0)
	subCmds[ExtractSubCmdType(Run, c.SubCmds)] = struct{}{}
	subCmds[ExtractSubCmdType(Config, c.SubCmds)] = struct{}{}
	subCmds[ExtractSubCmdType(Add, c.SubCmds)] = struct{}{}
	subCmds[ExtractSubCmdType(Remove, c.SubCmds)] = struct{}{}
	subCmds[ExtractSubCmdType(Edit, c.SubCmds)] = struct{}{}
	subCmds[ExtractSubCmdType(Help, c.SubCmds)] = struct{}{}
	subCmds[ExtractSubCmdType(Version, c.SubCmds)] = struct{}{}

	flags := make(map[flag.FlagType]struct{}, 0)
	flags[ExtractFlagType(flag.Profile, c.Flags)] = struct{}{}
	flags[ExtractFlagType(flag.Interval, c.Flags)] = struct{}{}
	flags[ExtractFlagType(flag.Alignment, c.Flags)] = struct{}{}
	flags[ExtractFlagType(flag.Opacity, c.Flags)] = struct{}{}
	flags[ExtractFlagType(flag.Stretch, c.Flags)] = struct{}{}

	fmt.Println()

	// length of 1 means only None and flag.None types are in the map
	// meaning there's no subcmds or flags
	if len(subCmds) == 1 && len(flags) == 1 {
		AddHelp(false)
	}

	for subCmd := range subCmds {
		if subCmd == None {
			continue
		}
		switch subCmd {
		case Add:
			AddHelp(true)
		}
	}

	for f := range flags {
		if f == flag.None {
			continue
		}
		fmt.Println("\t", f.ToString())
	}

	return nil
}
