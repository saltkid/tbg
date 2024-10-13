package main

import (
	"fmt"
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

func HelpValidateFlag(f *Flag) error {
	return nil // accept all flags
}

func HelpValidateSubCmd(c *Cmd) error {
	return nil // accept all subcommands
}

func HelpExecute(c *Cmd) error {
	if len(c.SubCmds) == 0 && len(c.Flags) == 0 {
		fmt.Printf("%-37s%s\n", *Decorate("tbg").Bold().Underline(),
			*Decorate("Terminal Background Gallery").Italic())
		fmt.Printf("%-30s%s\n", "", fmt.Sprintf("%s (%s) allows the user to have and manage multiple background", *Decorate("tbg").Bold(), *Decorate("teabag").Italic()))
		fmt.Printf("%-30s%s\n", "", fmt.Sprintf("images, that rotate at a set amount of time, for Windows Terminal."))
		fmt.Printf("%-37s%s\n", *Decorate("Version").Bold().Underline(),
			*Decorate(TbgVersion).Italic())
		fmt.Printf("\n%-37s%s\n", *Decorate("Usage").Bold().Underline(),
			*Decorate("tbg run").Italic())
		fmt.Printf("\n%s:\n", *Decorate("Commands").Bold().Underline())
		RunHelp(false)
		AddHelp(false)
		RemoveHelp(false)
		ConfigHelp(false)
		HelpHelp(false)
		VersionHelp(false)
		fmt.Printf("\n%s:\n", *Decorate("Flags").Bold().Underline())
		ProfileHelp(false)
		IntervalHelp(false)
		AlignmentHelp(false)
		StretchHelp(false)
		OpacityHelp(false)
		RandomHelp(false)
		fmt.Printf("\n%s\n", *Decorate("Not all flags are applicable to all commands. See help <command> for more info").Italic())
		return nil
	}

	// verbose messages
	fmt.Println("------------------------------------------------------------------------------------")
	for subCmd := range c.SubCmds {
		switch subCmd {
		case Run:
			RunHelp(true)
		case Add:
			AddHelp(true)
		case Remove:
			RemoveHelp(true)
		case Config:
			ConfigHelp(true)
		case Help:
			HelpHelp(true)
		case Version:
			VersionHelp(true)
		}
		fmt.Println("------------------------------------------------------------------------------------")
	}
	for f := range c.Flags {
		switch f {
		case Profile:
			ProfileHelp(true)
		case Interval:
			IntervalHelp(true)
		case Alignment:
			AlignmentHelp(true)
		case Stretch:
			StretchHelp(true)
		case Opacity:
			OpacityHelp(true)
		case Random:
			RandomHelp(true)
		}
		fmt.Println("------------------------------------------------------------------------------------")
	}
	return nil
}

func RunHelp(verbose bool) {
	fmt.Printf("%-33s%s", *Decorate("  run").Bold(),
		"reads the used config and edits Windows Terminal's settings.json to change background images\n")
	if verbose {
		fmt.Printf("\n  %s: run takes no args\n", *Decorate("Args").Bold())
		fmt.Printf("\n  %s:\n", *Decorate("Subcommands").Bold())
		fmt.Println("  1. config [arg]")
		fmt.Println("     [default, path/to/a/config.yaml]")
		fmt.Println("     You can specify which config to read from using the 'config' subcommand.")
		fmt.Println("     If you do not specify a config, the currently used config will be used.")
		fmt.Printf("\n  %s:\n", *Decorate("Flags").Bold())
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
		fmt.Printf("\n  %s:\n", *Decorate("Key Events").Bold())
		fmt.Println("  while tbg is running, it accepts optional key events.")
		fmt.Println("  Press a key to execute the command")
		fmt.Println("  1. q: [q]uit tbg")
		fmt.Println("  2. n: goes to [n]ext image")
		fmt.Println("  3. p: goes to [p]revious image")
		fmt.Println("  4. f: goes [f]orward to next image collection dir")
		fmt.Println("  5. b: goes [b]ack to previous image collection dir")
		fmt.Println("  6. c: list all [c]ommands")
		fmt.Printf("\n  %s:\n", *Decorate("Examples").Bold())
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
	fmt.Printf("%-33s%s", *Decorate("  add").Bold(),
		"Adds a path containing images to currently used config\n")
	if verbose {
		fmt.Printf("\n  %s:\n", *Decorate("Args").Bold())
		fmt.Println("  1. path/to/images/dir")
		fmt.Println("     Path to images dir should have at least one image")
		fmt.Println("     file under it. All subdirectories will be ignored.")
		fmt.Printf("\n  %s:\n", *Decorate("Subcommands").Bold())
		fmt.Println("  1. config [arg]")
		fmt.Println("     [default, path/to/a/config.yaml]")
		fmt.Println("     You can specify which config to add to using the 'config' subcommand.")
		fmt.Println("     If you do not specify a config, the currently used config will be used.")
		fmt.Printf("\n  %s:\n", *Decorate("Flags").Bold())
		fmt.Println("  You can specify alignment, stretch, and opacity using flags. See example 2 and 3")
		fmt.Println("  1. -a, --alignment [arg]")
		fmt.Println("         [top, topLeft, topRight, left, center, right, bottomLeft, bottom, bottomRight]")
		fmt.Println("  2. -o, --opacity   [arg]")
		fmt.Println("         [any float between 0 and 1 (inclusive)]")
		fmt.Println("  3. -s, --stretch   [arg]")
		fmt.Println("         [fill, none, uniform, uniformToFill]")
		fmt.Printf("\n  %s:\n", *Decorate("Examples").Bold())
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
		fmt.Println("      |   - /path/to/images/dir | top _ _")
		fmt.Println("      |")
		fmt.Println("      | default_alignment: right")
		fmt.Println("      | default_stretch: fill")
		fmt.Println("      | default_alignment: 0.1")
		fmt.Println("      |")
		fmt.Println("      | other fields...")
		fmt.Println("      ----------------------")
		fmt.Println("     The other flags that were not specified will be left as a blank undersocre (_)")
		fmt.Println("     This indicates that this blank flag will inherit its respective default")
		fmt.Println("\n  4. tbg add path/to/images/dir --alignment top")
		fmt.Println("      before:                                       after:")
		fmt.Println("      --------------------------------------   --------------------------------------------")
		fmt.Println("      | image_col_paths:                       | image_col_paths:")
		fmt.Println("      |   - /path/to/images/dir | center _ _   |   - /path/to/images/dir | top _ _")
		fmt.Println("      |                                        |")
		fmt.Println("      | default_alignment: right               | default_alignment: right")
		fmt.Println("      | default_stretch: fill                  | default_stretch: fill")
		fmt.Println("      | default_alignment: 0.2                 | default_alignment: 0.2")
		fmt.Println("      |                                        |")
		fmt.Println("      | other fields...                        | other fields...")
		fmt.Println("      --------------------------------------   ---------------------------------------------")
		fmt.Print("     You can also change flags of paths if they already exist\n\n")
	}
}

func RemoveHelp(verbose bool) {
	fmt.Printf("%-33s%s", *Decorate("  remove").Bold(),
		"Removes a path from the currently used config\n")
	if verbose {
		fmt.Printf("\n  %s:\n", *Decorate("Args").Bold())
		fmt.Println("  1. path/to/images/dir")
		fmt.Printf("\n  %s:\n", *Decorate("Subcommands").Bold())
		fmt.Println("  1. config [arg]")
		fmt.Println("     [default, path/to/a/config.yaml]")
		fmt.Println("     You can specify which config to remove from using the 'config' subcommand.")
		fmt.Println("     If you do not specify a config, the currently used config will be used.")
		fmt.Printf("\n  %s:\n", *Decorate("Flags").Bold())
		fmt.Println("  You can remove alignment, stretch, and opacity flags from a path by specifying flags")
		fmt.Println("  See example 2 and 3")
		fmt.Println("  1. -a, --alignment [arg]")
		fmt.Println("         [top, topLeft, topRight, left, center, right, bottomLeft, bottom, bottomRight]")
		fmt.Println("  2. -o, --opacity   [arg]")
		fmt.Println("         [any float between 0 and 1 (inclusive)]")
		fmt.Println("  3. -s, --stretch   [arg]")
		fmt.Println("         [fill, none, uniform, uniformToFill]")
		fmt.Printf("\n  %s:\n", *Decorate("Examples").Bold())
		fmt.Println("  1. tbg remove path/to/images/dir")
		fmt.Println("      before:                        after:")
		fmt.Println("      --------------------------     ---------------------")
		fmt.Println("      | image_col_paths: []          | image_col_paths: []")
		fmt.Println("      |   - /path/to/images/dir      |")
		fmt.Println("      |                              | other fields...")
		fmt.Println("      | other fields...              ---------------------")
		fmt.Println("      ----------------------")
		fmt.Println("      This is to remove a single path")
		fmt.Println("\n  2. tbg remove path/to/images/dir --alignment --stretch --opacity")
		fmt.Println("      before:                                       after:")
		fmt.Println("      -------------------------------------------   --------------------------------------------")
		fmt.Println("      | image_col_paths:                            | image_col_paths:")
		fmt.Println("      |   - /path/to/images/dir | center none 0.1   |   - /path/to/images/dir")
		fmt.Println("      |                                             |")
		fmt.Println("      | default_alignment: right                    | default_alignment: right")
		fmt.Println("      | default_stretch: fill                       | default_stretch: fill")
		fmt.Println("      | default_alignment: 0.2                      | default_alignment: 0.2")
		fmt.Println("      |                                             |")
		fmt.Println("      | other fields...                             | other fields...")
		fmt.Println("      -------------------------------------------   ---------------------------------------------")
		fmt.Println("      This is to remove all flags from a path: by specifying all 3 flags")
		fmt.Println("\n  3. tbg remove path/to/images/dir --alignment")
		fmt.Println("      before:                                       after:")
		fmt.Println("      -------------------------------------------   --------------------------------------------")
		fmt.Println("      | image_col_paths:                            | image_col_paths:")
		fmt.Println("      |   - /path/to/images/dir | center none 0.1   |   - /path/to/images/dir | _ none 0.1")
		fmt.Println("      |                                             |")
		fmt.Println("      | default_alignment: right                    | default_alignment: right")
		fmt.Println("      | default_stretch: fill                       | default_stretch: fill")
		fmt.Println("      | default_alignment: 0.2                      | default_alignment: 0.2")
		fmt.Println("      |                                             |")
		fmt.Println("      | other fields...                             | other fields...")
		fmt.Println("      -------------------------------------------   --------------------------------------------")
		fmt.Println("      If you don't specify all 3 flags to be removed, the flag you specified (--alingment in this example),")
		fmt.Println("      will be blanked out. This indicates that this blank flag will inherit its respective default flag")
		fmt.Println("      field value (default_alignment in this example)")
		fmt.Println("\n  4. tbg remove path/to/images/dir --stretch --opacity")
		fmt.Println("      before:                                       after:")
		fmt.Println("      -------------------------------------------   --------------------------------------------")
		fmt.Println("      | image_col_paths:                            | image_col_paths:")
		fmt.Println("      |   - /path/to/images/dir | _ none 0.1        |   - /path/to/images/dir")
		fmt.Println("      |                                             |")
		fmt.Println("      | default_alignment: right                    | default_alignment: right")
		fmt.Println("      | default_stretch: fill                       | default_stretch: fill")
		fmt.Println("      | default_alignment: 0.2                      | default_alignment: 0.2")
		fmt.Println("      |                                             |")
		fmt.Println("      | other fields...                             | other fields...")
		fmt.Println("      -------------------------------------------   --------------------------------------------")
		fmt.Println("      All flags were removed from the path because alignment was already blank and you specified")
		fmt.Println("      --stretch and --opacity. This means stretch and opacity will be blanked too. So instead of")
		fmt.Println("      keeping all 3 blank, tbg will just remove them as if you specified removing all 3 flags in")
		fmt.Print("      the first place. See example 2\n\n")
	}
}

func ConfigHelp(verbose bool) {
	fmt.Printf("%-33s%s", *Decorate("  config").Bold(),
		"Prints the currently used config if no arg.\n")
	fmt.Printf("%-25s%s", "",
		"If an arg is specified, it sets that arg as the currently used config, then prints it.\n")
	if verbose {
		fmt.Printf("\n  %s:\n", *Decorate("Args").Bold())
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
		fmt.Printf("\n  %s: config takes no subcommands\n", *Decorate("Subcommands").Bold())
		fmt.Printf("  %s: config takes no flags\n", *Decorate("Flags").Bold())
		fmt.Printf("\n  %s:\n", *Decorate("Examples").Bold())
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

func HelpHelp(verbose bool) {
	fmt.Printf("%-33s%s", *Decorate("  help").Bold(),
		"Prints this help message. If a command or flag is specified,\n")
	fmt.Printf("%-25s%s", "", "prints verbose help for that command or flag\n")
	if verbose {
		fmt.Printf("\n  %s:\n", *Decorate("Args").Bold())
		fmt.Println("  1. command names or flag names")
		fmt.Println("     allows multiple commands or flags to be specified")
		fmt.Printf("\n  %s:\n", *Decorate("Examples").Bold())
		fmt.Println("  1. tbg help")
		fmt.Println("     Prints this help message")
		fmt.Println("  2. tbg help run")
		fmt.Println("     Prints verbose help for the 'run' command")
		fmt.Println("  3. tbg help --alignment add --profile run")
		fmt.Println("     Prints verbose help for the '--alingment' flag,")
		fmt.Println("     'add' command, '--profile' flag and the 'run' command")
	}
}

func VersionHelp(verbose bool) {
	fmt.Printf("%-33s%s", *Decorate("  version").Bold(),
		"Prints the version of tbg\n")
	if verbose {
		fmt.Printf("\n  %s: --version does not take args\n", *Decorate("Args").Bold())
		fmt.Printf("\n  %s:\n", *Decorate("Examples").Bold())
		fmt.Println("  1. tbg version")
		fmt.Println("     Prints the version of tbg")
	}
}

func ProfileHelp(verbose bool) {
	fmt.Printf("%-33s%s", *Decorate("  -p, --profile").Bold(),
		"Specifies which Windows Terminal profile to use in a command.\n")
	if verbose {
		fmt.Printf("\n  %s:\n", *Decorate("Args").Bold())
		fmt.Println("  1. default")
		fmt.Println("     Sets the default profile as the profile to be used")
		fmt.Println("     by the parent command")
		fmt.Println("  2. list-n")
		fmt.Println("     where n is the list index Windows Terminal uses to identify the profile")
		fmt.Printf("\n  %s:\n", *Decorate("Examples").Bold())
		fmt.Println("  1. tbg run --profile default")
		fmt.Println("     whatever value the \"profile\" field in the currently used config is will")
		fmt.Println("     be ignored and tbg will edit the default Windows Terminal profile instead")
		fmt.Println("\n  2. tbg edit --profile list-2 config /path/to/a/config.yaml")
		fmt.Println("     this will change the \"profile\" field on the config /path/to/a/config.yaml")
		fmt.Println("     to \"list-2\"")
	}
}

func IntervalHelp(verbose bool) {
	fmt.Printf("%-33s%s", *Decorate("  -i, --interval").Bold(),
		"The interval of image change in minutes to use in a command.\n")
	if verbose {
		fmt.Printf("\n  %s:\n", *Decorate("Args").Bold())
		fmt.Println("  1. any positive integer")
		fmt.Printf("\n  %s:\n", *Decorate("Examples").Bold())
		fmt.Println("  1. tbg run --interval 30")
		fmt.Println("     whatever value the \"interval\" field in the currently used config is will")
		fmt.Println("     be ignored and tbg change images every 30 minutes instead.")
		fmt.Println("\n  2. tbg edit --interval 30 config /path/to/a/config.yaml")
		fmt.Println("     this will change the \"interval\" field on the config /path/to/a/config.yaml")
		fmt.Println("     to 30")
	}
}

func AlignmentHelp(verbose bool) {
	fmt.Printf("%-33s%s", *Decorate("  -a, --alignment").Bold(),
		"The alignment of the image to use in a command.\n")
	if verbose {
		fmt.Printf("\n  %s:\n", *Decorate("Args").Bold())
		fmt.Println("  1. topLeft,    top,    topRight")
		fmt.Println("  2. left,       center, right")
		fmt.Println("  3. bottomLeft, bottom, bottomRight")
		fmt.Printf("\n  %s:\n", *Decorate("Examples").Bold())
		fmt.Println("  1. tbg run --alignment center")
		fmt.Println("     whatever value the \"alignment\" field in the currently used config is will")
		fmt.Println("     be ignored and tbg will center the image instead")
		fmt.Println("\n  2. tbg edit --alignment center config /path/to/a/config.yaml")
		fmt.Println("     this will change the \"alignment\" field on the config /path/to/a/config.yaml")
		fmt.Println("     to \"center\"")
	}
}

func StretchHelp(verbose bool) {
	fmt.Printf("%-33s%s", *Decorate("  -s, --stretch").Bold(),
		"The stretch of the image to use in a command.\n")
	if verbose {
		fmt.Printf("\n  %s:\n", *Decorate("Args").Bold())
		fmt.Println("  1. fill, none, uniform, uniformToFill")
		fmt.Printf("\n  %s:\n", *Decorate("Examples").Bold())
		fmt.Println("  1. tbg run --stretch fill")
		fmt.Println("     whatever value the \"stretch\" field in the currently used config is will")
		fmt.Println("     be ignored and tbg will upscale the image to exactly fill the screen instead")
		fmt.Println("\n  2. tbg edit --stretch fill config /path/to/a/config.yaml")
		fmt.Println("     this will change the \"stretch\" field on the config /path/to/a/config.yaml")
		fmt.Println("     to \"fill\"")
	}
}

func OpacityHelp(verbose bool) {
	fmt.Printf("%-33s%s", *Decorate("  -o, --opacity").Bold(),
		"The opacity of the image to use in a command.\n")
	if verbose {
		fmt.Printf("\n  %s:\n", *Decorate("Args").Bold())
		fmt.Println("  1. any float between 0 and 1 (inclusive)")
		fmt.Printf("\n  %s:\n", *Decorate("Examples").Bold())
		fmt.Println("  1. tbg run --opacity 0.5")
		fmt.Println("     whatever value the \"opacity\" field in the currently used config is will")
		fmt.Println("     be ignored and tbg will set the image opacity to 0.5")
		fmt.Println("\n  2. tbg edit --opacity 0.5 config /path/to/a/config.yaml")
		fmt.Println("     this will change the \"opacity\" field on the config /path/to/a/config.yaml")
		fmt.Println("     to 0.5")
	}
}

func RandomHelp(verbose bool) {
	fmt.Printf("%-33s%s", *Decorate("  -r, --random").Bold(),
		"Randomize image collections and images. Specific to `run` command\n")
	if verbose {
		fmt.Printf("\n  %s: --random does not take args\n", *Decorate("Args").Bold())
		fmt.Printf("\n  %s:\n", *Decorate("Examples").Bold())
		fmt.Println("  1. tbg run --random")
		fmt.Println("     This will randomize the order of the image collections read")
		fmt.Println("     from the currently used config. It randomizes the order the")
		fmt.Println("     images in each collection.")
		fmt.Println()
		fmt.Println("     When image collections are exhausted and tbg wraps around, the")
		fmt.Println("     order of the image collections will be randomized again. This")
		fmt.Println("     behavior applies to images too.")
	}
}
