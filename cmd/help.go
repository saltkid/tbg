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
	if len(c.SubCmds) == 0 && len(c.Flags) == 0 {
		fmt.Println("tbg (Terminal Background Gallery)")
		fmt.Print("Version: ")
		VersionExecute()
		fmt.Println("Usage: tbg run")
		fmt.Println("\nCommands:")
		RunHelp(false)
		AddHelp(false)
		RemoveHelp(false)
		EditHelp(false)
		ConfigHelp(false)
		fmt.Println("\nFlags:")
		fmt.Println("\nNot all flags are applicable to all commands. See help <command> for more info")
	}

	for subCmd := range c.SubCmds {
		switch subCmd {
		case Run:
			RunHelp(true)
		case Add:
			AddHelp(true)
		case Remove:
			RemoveHelp(true)
		case Edit:
			EditHelp(true)
		case Config:
			ConfigHelp(true)
		}
		fmt.Println("------------------------------------------------------------------------------------")
	}

	for f := range c.Flags {
		fmt.Println("\t", f.ToString())
	}

	return nil
}

func RunHelp(verbose bool) {
	fmt.Printf("%-30s%s", "  run",
		"reads the used config and edits Windows Terminal's settings.json to change background images\n")
	if verbose {
		fmt.Println("\n  Args: run takes no args")
		fmt.Println("\n  Subcommands:")
		fmt.Println("  1. config [arg]")
		fmt.Println("     [default, path/to/a/config.yaml]")
		fmt.Println("     You can specify which config to read from using the 'config' subcommand.")
		fmt.Println("     If you do not specify a config, the currently used config will be used.")
		fmt.Println("\n  Flags:")
		fmt.Println("  You can specify alignment, stretch, and opacity using flags.")
		fmt.Println("  These will override the values in the used config (not edit)")
		fmt.Println("  1. -a, --alignment [arg]")
		fmt.Println("         [top, topLeft, topRight, left, center, right, bottomLeft, bottom, bottomRight]")
		fmt.Println("  2. -o, --opacity   [arg]")
		fmt.Println("         [any float between 0 and 1 (inclusive)]")
		fmt.Println("  3. -s, --stretch   [arg]")
		fmt.Println("         [fill, none, uniform, uniformToFill]")
		fmt.Println("  4. -p, --profile   [arg]")
		fmt.Println("         [default, list-n]")
		fmt.Println("         where n is the list index Windows Terminal uses to identify the profile")
		fmt.Println("  5. -i, --interval  [arg]")
		fmt.Println("         [any positive integer]")
		fmt.Println("         note that this is in minutes")
		fmt.Println("\n  Examples:")
		fmt.Println("  1. tbg run")
		fmt.Println("     This will use the currently used config's values to edit Windows Terminal's settings.json")
		fmt.Println("\n  2. tbg run config path/to/a/config.yaml ")
		fmt.Println("     tbg run config default ")
		fmt.Println("     These two are similar in the sense that this will have tbg use whatever")
		fmt.Println("     config was specified to edit Windows Terminal's settings.json")
		fmt.Println("\n  3. tbg run --profile list-2 --interval 5 --alignment center")
		fmt.Println("      used_config                      values used to edit settings.json")
		fmt.Println("      --------------------------       ------------------------------------------------")
		fmt.Println("      | image_col_paths:               | image_col_paths:")
		fmt.Println("      |   - /path/to/images/dir1       |   - /path/to/images/dir1 | center fill 0.1")
		fmt.Println("      |   - /path/to/images/dir2       |   - /path/to/images/dir2 | center fill 0.1")
		fmt.Println("      |                                |")
		fmt.Println("      | profile: default               | profile: list-2")
		fmt.Println("      | interval: 30                   | interval: 5")
		fmt.Println("      |                                ------------------------------------------------")
		fmt.Println("      | default_alignment: right")
		fmt.Println("      | default_stretch: fill")
		fmt.Println("      | default_alignment: 0.1")
		fmt.Println("      --------------------------")
		fmt.Println("     This means that instead of editing the default profile, it will edit the")
		fmt.Println("     2nd profile in Windows Terminal's list. The interval will be 5 minutes")
		fmt.Println("     instead of 30 minutes.")
		fmt.Println("     The dirs's alignment is set to center instead of inheriting the.")
		fmt.Println("     default_alignment. The stretch and opacity though are inherited")
		fmt.Println("     from the default values since it was not specified by the user.")
		fmt.Println("\n     Also note that the values on the right are not the 'edited' version but only")
		fmt.Print("     exist in the current execution. The values in the config stays the same\n\n")
	}
}

func AddHelp(verbose bool) {
	fmt.Printf("%-30s%s", "  add",
		"Adds a path containing images to currently used config\n")
	if verbose {
		fmt.Println("\n  Args:")
		fmt.Println("  1. path/to/images/dir")
		fmt.Println("     Path to images dir should have at least one image")
		fmt.Println("     file under it. All subdirectories will be ignored.")
		fmt.Println("\n  Subcommands:")
		fmt.Println("  1. config [arg]")
		fmt.Println("     [default, path/to/a/config.yaml]")
		fmt.Println("     You can specify which config to add to using the 'config' subcommand.")
		fmt.Println("     If you do not specify a config, the currently used config will be used.")
		fmt.Println("\n  Flags:")
		fmt.Println("  You can specify alignment, stretch, and opacity using flags. See example 2 and 3")
		fmt.Println("  1. -a, --alignment [arg]")
		fmt.Println("         [top, topLeft, topRight, left, center, right, bottomLeft, bottom, bottomRight]")
		fmt.Println("  2. -o, --opacity   [arg]")
		fmt.Println("         [any float between 0 and 1 (inclusive)]")
		fmt.Println("  3. -s, --stretch   [arg]")
		fmt.Println("         [fill, none, uniform, uniformToFill]")
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

func RemoveHelp(verbose bool) {
	fmt.Printf("%-30s%s", "  remove",
		"Removes a path from the currently used config\n")
	if verbose {
		fmt.Println("\n  Args:")
		fmt.Println("  1. path/to/images/dir")
		fmt.Println("\n  Subcommands:")
		fmt.Println("  1. config [arg]")
		fmt.Println("     [default, path/to/a/config.yaml]")
		fmt.Println("     You can specify which config to remove from using the 'config' subcommand.")
		fmt.Println("     If you do not specify a config, the currently used config will be used.")
		fmt.Println("\n  Flags: remove takes no flags")
		fmt.Println("\n  Examples:")
		fmt.Println("  1. tbg remove path/to/images/dir")
		fmt.Println("      before:                        after:")
		fmt.Println("      --------------------------     ---------------------")
		fmt.Println("      | image_col_paths: []          | image_col_paths: []")
		fmt.Println("      |   - /path/to/images/dir      |")
		fmt.Println("      |                              | other fields...")
		fmt.Println("      | other fields...              ---------------------")
		fmt.Println("      ----------------------")
	}
}

func EditHelp(verbose bool) {
	fmt.Printf("%-30s%s", "  edit",
		"Edits a path, all paths, or just the fields from the currently used config\n")
	if verbose {
		fmt.Println("\n  Args:")
		fmt.Println("  1. path/to/images/dir")
		fmt.Println("  2. all")
		fmt.Println("     edit all paths with specified flags")
		fmt.Println("  3. fields")
		fmt.Println("     edit only the default flag fields with the specified flags")
		fmt.Println("\n  Subcommands:")
		fmt.Println("  1. config [arg]")
		fmt.Println("     [default, path/to/a/config.yaml]")
		fmt.Println("     You can specify which config to remove from using the 'config' subcommand.")
		fmt.Println("     If you do not specify a config, the currently used config will be used.")
		fmt.Println("\n  Flags:")
		fmt.Println("  You can specify alignment, stretch, and opacity using flags. See examples below.")
		fmt.Println("  1. -a, --alignment [arg]")
		fmt.Println("         [top, topLeft, topRight, left, center, right, bottomLeft, bottom, bottomRight]")
		fmt.Println("  2. -o, --opacity   [arg]")
		fmt.Println("         [any float between 0 and 1 (inclusive)]")
		fmt.Println("  3. -s, --stretch   [arg]")
		fmt.Println("         [fill, none, uniform, uniformToFill]")
		fmt.Println("  4. -p, --profile   [arg]")
		fmt.Println("         [default, list-n]")
		fmt.Println("         where n is the list index Windows Terminal uses to identify the profile")
		fmt.Println("         See example 3.")
		fmt.Println("  5. -i, --interval  [arg]")
		fmt.Println("         [any positive integer]")
		fmt.Println("         note that this is in minutes")
		fmt.Println("\n  Note that profile and interval are always edited on the config level,")
		fmt.Println("  never on per path level even if user specified 'all' or a path. Only,")
		fmt.Println("  alignment, stretch, and opacity can be both per path and config level")
		fmt.Println("  See example 3.")
		fmt.Println("\n  Examples:")
		fmt.Println("  1. tbg edit path/to/images/dir --alignment center --stretch none")
		fmt.Println("      before:                          after:")
		fmt.Println("      --------------------------       --------------------------------------------")
		fmt.Println("      | image_col_paths:               | image_col_paths:")
		fmt.Println("      |   - /path/to/images/dir        |   - /path/to/images/dir | center none 0.1")
		fmt.Println("      |                                |")
		fmt.Println("      | default_alignment: right       | default_alignment: right")
		fmt.Println("      | default_stretch: fill          | default_stretch: fill")
		fmt.Println("      | default_alignment: 0.1         | default_alignment: 0.1")
		fmt.Println("      |                                |")
		fmt.Println("      | other fields...                | other fields...")
		fmt.Println("      --------------------------       ---------------------------------------------")
		fmt.Println("     Notice that even though there are only 2 flags specified, there are 3 after |")
		fmt.Println("     This is because it inerited the default value for opacity since it was not specified")
		fmt.Println("\n  2. tbg edit fields --alignment center --stretch none")
		fmt.Println("      before:                          after:")
		fmt.Println("      --------------------------       -----------------------------")
		fmt.Println("      | image_col_paths:               | image_col_paths:")
		fmt.Println("      |   - /path/to/images/dir        |   - /path/to/images/dir")
		fmt.Println("      |                                |")
		fmt.Println("      | default_alignment: right       | default_alignment: center")
		fmt.Println("      | default_stretch: fill          | default_stretch: none")
		fmt.Println("      | default_alignment: 0.1         | default_alignment: 0.1")
		fmt.Println("      |                                |")
		fmt.Println("      | other fields...                | other fields...")
		fmt.Println("      --------------------------       -----------------------------")
		fmt.Println("     Notice that only the default fields were edited, not the paths")
		fmt.Println("\n  3. tbg edit all --opacity 0.5 --profile list-1 --interval 5")
		fmt.Println("      before:                          after:")
		fmt.Println("      --------------------------       ------------------------------------------------")
		fmt.Println("      | image_col_paths:               | image_col_paths:")
		fmt.Println("      |   - /path/to/images/dir1       |   - /path/to/images/dir1 | center uniform 0.5")
		fmt.Println("      |   - /path/to/images/dir2       |   - /path/to/images/dir2 | center uniform 0.5")
		fmt.Println("      |                                |")
		fmt.Println("      | profile: default               | profile: list-1")
		fmt.Println("      | interval: 30                   | interval: 5")
		fmt.Println("      |                                |")
		fmt.Println("      | default_alignment: right       | default_alignment: right")
		fmt.Println("      | default_stretch: fill          | default_stretch: fill")
		fmt.Println("      | default_alignment: 0.1         | default_alignment: 0.1")
		fmt.Println("      --------------------------       -------------------------------------------------")
		fmt.Println("     Notice that only the paths were edited, not the default fields")
		fmt.Println("     Also, even though the edit arg was 'all', profile and interval")
		fmt.Print("     were edited on the config level\n\n")
	}
}

func ConfigHelp(verbose bool) {
	fmt.Printf("%-30s%s", "  config",
		"Prints the currently used config if no arg.\n")
	fmt.Printf("%-30s%s", "",
		"If an arg is specified, it sets that arg as the currently used config, then prints it.\n")
	if verbose {
		fmt.Println("\n  Args:")
		fmt.Println("  1. path/to/a-config.yaml")
		fmt.Println("     Sets this path as the currently used config.")
		fmt.Println("     It does this by editing the 'used_config' field on tbg_profile.yaml")
		fmt.Println("     in the same directory as the tbg executable (auto generated)")
		fmt.Println("  2. default")
		fmt.Println("     Sets the default config as the currently used config.")
		fmt.Println("     Default config is the config.yaml on the same path as the tbg")
		fmt.Println("     executable. It was auto generated by tbg on initial execution.")
		fmt.Println("     This also edits the 'used_config' field in tbg_profile.yaml")
		fmt.Println("  3. no arg")
		fmt.Println("     Prints the currently used config.")
		fmt.Println("\n  Subcommands: config takes no subcommands")
		fmt.Println("  Flags: config takes no flags")
		fmt.Println("\n  Examples:")
		fmt.Println("  1. tbg config")
		fmt.Println("      print currently used config:")
		fmt.Println("      --------------------------")
		fmt.Println("      | image_col_paths:")
		fmt.Println("      |   - /path/to/images/dir")
		fmt.Println("      |")
		fmt.Println("      | profile: default")
		fmt.Println("      | interval: 30")
		fmt.Println("      |")
		fmt.Println("      | default_alignment: right")
		fmt.Println("      | default_stretch: fill")
		fmt.Println("      | default_alignment: 0.1")
		fmt.Println("      --------------------------")
		fmt.Println("\n  2. tbg config default")
		fmt.Println("      Sets the default config as the currently used config then prints it.")
		fmt.Println("      It does this by editing the 'used_config' field on tbg_profile.yaml")
		fmt.Println("      before:")
		fmt.Println("      ----------------------------------")
		fmt.Println("      | # tbg_profile.yaml")
		fmt.Println("      | used_config: path/to/config.yaml")
		fmt.Println("      ----------------------------------")
		fmt.Println("      after:")
		fmt.Println("      ----------------------------------")
		fmt.Println("      | # tbg_profile.yaml")
		fmt.Println("      | used_config: path/to/default/config.yaml")
		fmt.Println("      ----------------------------------")
		fmt.Println("      Note that specifying a path/to/config.yaml instead of 'default' will do essentially")
		fmt.Print("      the same thing: set the used_config field to the specified path.\n\n")
	}
}
