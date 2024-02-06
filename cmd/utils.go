package cmd

import (
	"github.com/saltkid/tbg/flag"
)

// returns nil if not found in map of subcommands
func ExtractSubCmd(val CmdType, subCmds map[CmdType]*Cmd) *Cmd {
	subCmd, ok := subCmds[val]
	if !ok {
		return nil
	}
	return subCmd
}

// returns None if not found in map of subcommands
func ExtractSubCmdType(val CmdType, subCmds map[CmdType]*Cmd) CmdType {
	sc := ExtractSubCmd(val, subCmds)
	if sc == nil {
		return None
	}
	return sc.Type
}

// returns empty string if not found in map of subcommands
func ExtractSubCmdValue(val CmdType, subCmds map[CmdType]*Cmd) string {
	subCmd := ExtractSubCmd(val, subCmds)
	if subCmd == nil {
		return ""
	}
	return subCmd.Value
}

// returns nil if not found in map of flags
func ExtractFlag(val flag.FlagType, flags map[flag.FlagType]*flag.Flag) *flag.Flag {
	flag, ok := flags[val]
	if !ok {
		return nil
	}
	return flag
}

// returns None if not found in map of flags
func ExtractFlagType(val flag.FlagType, flags map[flag.FlagType]*flag.Flag) flag.FlagType {
	f := ExtractFlag(val, flags)
	if f == nil {
		return flag.None
	}
	return f.Type
}

// returns empty string if not found in map of flags
func ExtractFlagValue(val flag.FlagType, flags map[flag.FlagType]*flag.Flag) string {
	flag := ExtractFlag(val, flags)
	if flag == nil {
		return ""
	}
	return flag.Value
}
