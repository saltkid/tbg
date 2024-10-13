package main

// returns empty string if not found in map of subcommands
func ExtractSubCmdValue(val CmdType, subCmds map[CmdType]*Cmd) *string {
	subCmd, ok := subCmds[val]
	if !ok {
		return nil
	}
	return &subCmd.Value
}

// returns empty string if not found in map of flags
func ExtractFlagValue(val FlagType, flags map[FlagType]*Flag) *string {
	flag, ok := flags[val]
	if !ok {
		return nil
	}
	return &flag.Value
}
