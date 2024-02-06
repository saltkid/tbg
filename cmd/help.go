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
		fmt.Printf("help takes no args. got '%s'\n", val)
		return nil
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
		fmt.Println("tbg (Terminal Background Gallery)")
		fmt.Print("Version: ")
		VersionExecute()
		fmt.Println("Usage: tbg run")
		fmt.Println("\nCommands:")
		AddHelp(false)
		fmt.Println("\nFlags:")
		fmt.Println("\nNot all flags are applicable to all commands. See help <command> for more info")
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

func AddHelp(verbose bool) {
	fmt.Printf("%-60s%s", "  add path/to/images/dir",
		"Adds a path containing images to currently used config\n")
	if verbose {
		fmt.Println("\n\n  Path to images dir should have at least one image")
		fmt.Println("  file under it. All subdirectories will be ignored.")
		fmt.Println("\n  You can specify which config to add to using the 'config' subcommand.")
		fmt.Println("  If you do not specify a config, the currently used config will be used.")
		fmt.Println("\n  You can specify alignment, stretch, and opacity using flags.")
		fmt.Println("\n  Examples:")
		fmt.Println("  1. tbg add path/to/images/dir")
		fmt.Println("     This is how it would look like in the config:")
		fmt.Println("      ----------------------")
		fmt.Println("      | image_col_paths:")
		fmt.Println("      |   - /path/to/images/dir")
		fmt.Println("      |")
		fmt.Println("      | other fields...")
		fmt.Println("      ----------------------")
		fmt.Println("\n  2. tbg add path/to/images/dir --alignment center --opacity 0.5 --stretch uniform")
		fmt.Println("     This is how it would look like in the config:")
		fmt.Println("      ----------------------")
		fmt.Println("      | image_col_paths:")
		fmt.Println("      |   - /path/to/images/dir | center uniform 0.5")
		fmt.Println("      |")
		fmt.Println("      | other fields...")
		fmt.Println("      ----------------------")
		fmt.Println("\n  3. tbg add path/to/images/dir --alignment top")
		fmt.Println("     This is how it would look like in the config:")
		fmt.Println("      ----------------------")
		fmt.Println("      | image_col_paths:")
		fmt.Println("      |   - /path/to/images/dir | top fill 0.1")
		fmt.Println("      |")
		fmt.Println("      | default_alignment: right")
		fmt.Println("      | default_stretch: fill")
		fmt.Println("      | default_alignment: 0.1")
		fmt.Println("      |")
		fmt.Println("      | other fields...")
		fmt.Println("      ----------------------")
		fmt.Println("     Notice that even though there is only one flag specified, there are 3 after |")
		fmt.Print("     This is because it inerited the default values for flags that were not specified\n\n")
	}
}
