package cmd

import (
	"github.com/saltkid/tbg/flag"
)

// returns empty string if not found in map of subcommands
func ExtractSubCmdValue(val CmdType, subCmds map[CmdType]*Cmd) string {
	subCmd, ok := subCmds[val]
	if !ok {
		return ""
	}
	return subCmd.Value
}

// returns empty string if not found in map of flags
func ExtractFlagValue(val flag.FlagType, flags map[flag.FlagType]*flag.Flag) string {
	flag, ok := flags[val]
	if !ok {
		return ""
	}
	return flag.Value
}
