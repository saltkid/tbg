package main

import (
	"fmt"
)

type HelpCommand struct {
	Flags    []Flag
	Commands []Command
}

func (cmd *HelpCommand) Type() CommandType { return HelpCommandType }

func (cmd *HelpCommand) Debug() {
	fmt.Println("Help Command")
	if len(cmd.Commands) > 0 {
		fmt.Println("Commands:")
	}
	for _, c := range cmd.Commands {
		fmt.Println(" ", c.Type())
	}
	if len(cmd.Flags) > 0 {
		fmt.Println("Flags:")
	}
	for _, f := range cmd.Flags {
		fmt.Println(" ", f.Type)
	}
}

func (cmd *HelpCommand) ValidateValue(val *string) error {
	if val == nil || *val == "" {
		return nil
	}
	return fmt.Errorf("'help' takes no args. got: '%s'", *val)
}

func (cmd *HelpCommand) ValidateFlag(f Flag) error {
	cmd.Flags = append(cmd.Flags, f)
	return nil
}

func (cmd *HelpCommand) ValidateSubCommand(sc Command) error {
	cmd.Commands = append(cmd.Commands, sc)
	return nil
}

func (cmd *HelpCommand) Execute() error {
	notVerbose := len(cmd.Commands) == 0 && len(cmd.Flags) == 0
	if notVerbose {
		fmt.Printf("%-37s%s\n",
			Decorate("tbg").Bold().Underline(),
			Decorate("Terminal Background Gallery").Italic(),
		)
		fmt.Printf("%-30s%s\n", "",
			fmt.Sprintf("%s (%s) allows the user to have and manage multiple background",
				Decorate("tbg").Bold(),
				Decorate("teabag").Italic(),
			),
		)
		fmt.Printf("%-30s%s\n", "",
			"images, that rotate at a set amount of time, for Windows Terminal.",
		)
		fmt.Printf("%-37s%s\n",
			Decorate("Version").Bold().Underline(),
			Decorate(TbgVersion).Italic(),
		)
		fmt.Printf("\n%-37s%s\n",
			Decorate("Usage").Bold().Underline(),
			Decorate("tbg run").Italic(),
		)
		fmt.Printf("\n%s:\n",
			Decorate("Commands").Bold().Underline(),
		)
		RunHelp(false)
		AddHelp(false)
		RemoveHelp(false)
		ConfigHelp(false)
		HelpHelp(false)
		VersionHelp(false)
		fmt.Printf("\n%s:\n",
			Decorate("Flags").Bold().Underline(),
		)
		ProfileHelp(false)
		IntervalHelp(false)
		AlignmentHelp(false)
		StretchHelp(false)
		OpacityHelp(false)
		RandomHelp(false)
		fmt.Printf("\n%s\n",
			Decorate("Not all flags are applicable to all commands. See help <command> for more info").Italic(),
		)
		return nil
	}

	// verbose messages
	fmt.Println("------------------------------------------------------------------------------------")
	for _, subCmd := range cmd.Commands {
		switch subCmd.Type() {
		case RunCommandType:
			RunHelp(true)
		case AddCommandType:
			AddHelp(true)
		case RemoveCommandType:
			RemoveHelp(true)
		case ConfigCommandType:
			ConfigHelp(true)
		case HelpCommandType:
			HelpHelp(true)
		case VersionCommandType:
			VersionHelp(true)
		}
		fmt.Println("------------------------------------------------------------------------------------")
	}
	for _, f := range cmd.Flags {
		switch f.Type {
		case ProfileFlag:
			ProfileHelp(true)
		case IntervalFlag:
			IntervalHelp(true)
		case AlignmentFlag:
			AlignmentHelp(true)
		case StretchFlag:
			StretchHelp(true)
		case OpacityFlag:
			OpacityHelp(true)
		case RandomFlag:
			RandomHelp(true)
		}
		fmt.Println("------------------------------------------------------------------------------------")
	}
	return nil
}

func RunHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  run").Bold(),
		"reads the used config and edits Windows Terminal's settings.json to change background images\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `: run takes no args

  `, Decorate("Subcommands").Bold(), `
  1. config [arg]
     [default, path/to/a/config.yaml]
     You can specify which config to read from using the 'config' subcommand.
     If you do not specify a config, the currently used config will be used.

  `, Decorate("Flags").Bold(), `:,
  You can specify alignment, stretch, and opacity using flags.
  These will override the values in the used config (not edit)
  1. -a, --alignment [arg]
         [top, topLeft, topRight, left, center, right, bottomLeft, bottom, bottomRight]
  2. -o, --opacity   [arg]
         [any float between 0 and 1 (inclusive)]
  3. -s, --stretch   [arg]
         [fill, none, uniform, uniformToFill]
  4. -p, --profile   [arg]
         [default, list-n]
         where n is the list index Windows Terminal uses to identify the profile
  5. -i, --interval  [arg]
         [any positive integer]
         note that this is in minutes
  6. -r, --random
         randomize shown image paths and images

  `, Decorate("Key Events").Bold(), `:
  while tbg is running, it accepts optional key events.
  Press a key to execute the command
  1. q: [q]uit tbg
  2. n: goes to [n]ext image
  3. p: goes to [p]revious image
  4. f: goes [f]orward to next image collection dir
  5. b: goes [b]ack to previous image collection dir
  6. c: list all [c]ommands

  `, Decorate("Examples").Bold(), `:
  1. tbg run
     This will use the currently used config's values to edit Windows Terminal's settings.json

  2. tbg run config path/to/a/config.yaml 
     tbg run config default 
     These two are similar in the sense that this will have tbg use whatever
     config was specified to edit Windows Terminal's settings.json

  3. tbg run --profile list-2 --interval 5 --alignment center
      used_config                      values used to edit settings.json
      --------------------------       --------------------------------
      | paths:                         | paths:
      |   - path: /path/to/images/dir1 |   - path: /path/to/images/dir1
      |   - path: /path/to/images/dir2 |   - path: /path/to/images/dir2
      |                                |
      | profile: default               | profile: list-2
      | interval: 30                   | interval: 5
      |                                |
      | alignment: right               | alignment: center
      | stretch: fill                  | stretch: fill
      | alignment: 0.1                 | alignment: 0.1
      --------------------------       --------------------------------
     This means that instead of editing the default profile, it will edit the
     2nd profile in Windows Terminal's list. The interval will be 5 minutes
     instead of 30 minutes.
     The alignment is set to center instead of right which was the value
     earlier. The stretch and opacity stay the same since it was not
     specified by the user.

     Also note that the values on the right are not the 'edited' version but only
     exist in the current execution. The values in the config stays the same\n\n
`)
	}
}

func AddHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  add").Bold(),
		"Adds a path containing images to currently used config\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `:
  1. path/to/images/dir
     Path to images dir should have at least one image
     file under it. All subdirectories will be ignored.

  `, Decorate("Subcommands").Bold(), `:
  1. config [arg]
     [default, path/to/a/config.yaml]
     You can specify which config to add to using the 'config' subcommand.
     If you do not specify a config, the currently used config will be used.

  `, Decorate("Flags").Bold(), `:
  You can specify alignment, stretch, and opacity using flags. See example 2 and 3
  1. -a, --alignment [arg]
         [top, topLeft, topRight, left, center, right, bottomLeft, bottom, bottomRight]
  2. -o, --opacity   [arg]
         [any float between 0 and 1 (inclusive)]
  3. -s, --stretch   [arg]
         [fill, none, uniform, uniformToFill]

  `, Decorate("Examples").Bold(), `:
  1. tbg add path/to/images/dir
     This is how it would look like in the config:
      ----------------------
      | paths:
      |   - path: /path/to/images/dir
      |
      | other fields...
      ----------------------

  2. tbg add path/to/another/images/dir --alignment center --opacity 0.5 --stretch uniform
     This is how it would look like in the config:
      ----------------------
      | paths:
      |   - path: /path/to/images/dir
      |   - path: /path/to/another/images/dir 
      |     alignment: center
      |     opacity: 0.5
      |     stretch: uniform
      |
      | other fields...
      ----------------------

  3. tbg add path/to/the-other/images/dir --alignment top
     This is how it would look like in the config:
      ----------------------
      | paths:
      |   - path: /path/to/images/dir
      |   - path: /path/to/another/images/dir 
      |     alignment: center
      |     opacity: 0.5
      |     stretch: uniform
      |   - path: /path/to/the-other/images/dir 
      |     alignment: top
      |
      | alignment: right
      | stretch: fill
      | alignment: 0.1
      |
      | other fields...
      ----------------------
     The other flags that were not specified will inherit its respective default

  4. tbg add path/to/images/dir --alignment top
      before:                                       after:
      -------------------------        -------------------------------
      | paths:                         | paths:
      |   - path: /path/to/images/dir  |   - path: /path/to/images/dir
      |     alignment: center          |     alignment: top
      |                                |
      | other fields...                | other fields...
      -------------------------        -------------------------------
     You can also change flags of paths if they already exist\n\n
`)
	}
}

func RemoveHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  remove").Bold(),
		"Removes a path from the currently used config\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `:
  1. path/to/images/dir

  `, Decorate("Subcommands").Bold(), `:
  1. config [arg]
     [default, path/to/a/config.yaml]
     You can specify which config to remove from using the 'config' subcommand.
     If you do not specify a config, the currently used config will be used.

  `, Decorate("Flags").Bold(), `:
  You can remove alignment, stretch, and opacity flags from a path by specifying flags
  See example 2 and 3
  1. -a, --alignment [arg]
         [top, topLeft, topRight, left, center, right, bottomLeft, bottom, bottomRight]
  2. -o, --opacity   [arg]
         [any float between 0 and 1 (inclusive)]
  3. -s, --stretch   [arg]
         [fill, none, uniform, uniformToFill]

  `, Decorate("Examples").Bold(), `:
  1. tbg remove path/to/images/dir
      before:                        after:
      -------------------------------   ---------------------
      | paths:                          | paths: []
      |   - path: /path/to/images/dir   |
      |                                 | other fields...
      | other fields...                 ---------------------
      -------------------------------
      This is to remove a single path

  2. tbg remove path/to/images/dir --alignment --stretch --opacity
      before:                                       after:
      -------------------------         -------------------------
      | paths:                          | paths:
      |   - path: /path/to/images/dir   |   - path: /path/to/images/dir
	  |     alignment: top              |
      |     opacity: 0.5                |  other fields...
      |     stretch: fill               -------------------------
      |
      | other fields...
      -------------------------
      This is to remove all flags from a path: by specifying all 3 flags

  3. tbg remove path/to/images/dir --alignment
      before:                                       after:
      -------------------------------   -------------------------------
      | paths:                          | paths:
      |   - path: /path/to/images/dir   |   - path: /path/to/images/dir
      |     alignment: center           |     stretch: fill
      |     stretch: fill               |     opacity: 0.1
      |     opacity: 0.1                |
      |                                 | alignment: right
      | alignment: right                | stretch: fill
      | stretch: fill                   | alignment: 0.2
      | alignment: 0.2                  | 
      |                                 | other fields...
      | other fields...                 -------------------------------
      -------------------------------   
      If you don't specify all 3 flags to be removed, only the flag you specified
      (--alingment in this example), will be removed out. This indicates that this
      flag will inherit its respective default flag value\n\n
`)
	}
}

func ConfigHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  config").Bold(),
		"Prints the currently used config if no arg.\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `:
  1. path/to/a-config.yaml
     Sets this path as the currently used config.
     It does this by editing the 'used_config' field on tbg_profile.yaml
     in the same directory as the tbg executable (auto generated)
  2. no arg
     Prints the currently used config.

  `, Decorate("Subcommands").Bold(), `: config takes no subcommands
  `, Decorate("Flags").Bold(), `:
  1. -a, --alignment [arg]
         [top, topLeft, topRight, left, center, right, bottomLeft, bottom, bottomRight]
  2. -o, --opacity   [arg]
         [any float between 0 and 1 (inclusive)]
  3. -s, --stretch   [arg]
         [fill, none, uniform, uniformToFill]
  4. -p, --profile   [arg]
         [default, list-n]
         where n is the list index Windows Terminal uses to identify the profile.
  5. -i, --interval  [arg]
         [any positive integer]
         note that this is in minutes.

  `, Decorate("Examples").Bold(), `:
  1. tbg config
      print currently used config:
      --------------------------
      | paths:
      |   - path: /path/to/images/dir
      |
      | profile: default
      | interval: 30
      |
      | alignment: right
      | stretch: fill
      | alignment: 0.1
      --------------------------

  2. tbg config --alignment center
      replaces the config's "alignment" field with the value "center"
      --------------------------       --------------------------------
      | paths:                         | paths:
      |   - path: /path/to/images/dir1 |   - path: /path/to/images/dir1
      |                                |
      | alignment: right               | alignment: center
      |                                |
      | other fields...                | other fields...
      --------------------------       --------------------------------
      This applies to all the other flags as well. You can specify any number of flags.\n\n
`)
	}
}

func HelpHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  help").Bold(),
		"Prints this help message. If a command or flag is specified,\n",
	)
	fmt.Printf("%-25s%s", "",
		"prints verbose help for that command or flag\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `:
  1. command names or flag names
     allows multiple commands or flags to be specified

  `, Decorate("Examples").Bold(), `:
  1. tbg help
     Prints this help message
  2. tbg help run
     Prints verbose help for the 'run' command
  3. tbg help --alignment add --profile run
     Prints verbose help for the '--alingment' flag,
     'add' command, '--profile' flag and the 'run' command
`)
	}
}

func VersionHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  version").Bold(),
		"Prints the version of tbg\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `: --version does not take args

  `, Decorate("Examples").Bold(), `:
  1. tbg version
     Prints the version of tbg
`)
	}
}

func ProfileHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  -p, --profile").Bold(),
		"Specifies which Windows Terminal profile to use in a command.\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `:
  1. default
     Sets the default profile as the profile to be used
     by the parent command
  2. list-n
     where n is the list index Windows Terminal uses to identify the profile

  `, Decorate("Examples").Bold(), `:
  1. tbg run --profile default
     whatever value the "profile" field in the currently used config will
     be ignored and tbg will edit the default Windows Terminal profile instead

  2. tbg edit --profile list-2 config /path/to/a/config.yaml
     this will change the "profile" field on the config /path/to/a/config.yaml
     to "list-2"
`)
	}
}

func IntervalHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  -i, --interval").Bold(),
		"The interval of image change in minutes to use in a command.\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `:
  1. any positive integer

  `, Decorate("Examples").Bold(), `:
  1. tbg run --interval 30
     whatever value the \"interval\" field in the currently used config will
     be ignored and tbg change images every 30 minutes instead.

  2. tbg edit --interval 30 config /path/to/a/config.yaml
     this will change the \"interval\" field on the config /path/to/a/config.yaml
     to 30
`)
	}
}

func AlignmentHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  -a, --alignment").Bold(),
		"The alignment of the image to use in a command.\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `:
  1. topLeft,    top,    topRight
  2. left,       center, right
  3. bottomLeft, bottom, bottomRight

  `, Decorate("Examples").Bold(), `:
  1. tbg run --alignment center
     whatever value the "alignment" field in the currently used config will
     be ignored and tbg will center the image instead

  2. tbg edit --alignment center config /path/to/a/config.yaml
     this will change the "alignment" field on the config /path/to/a/config.yaml
     to "center"
`)
	}
}

func StretchHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  -s, --stretch").Bold(),
		"The stretch of the image to use in a command.\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `:
  1. fill, none, uniform, uniformToFill

  `, Decorate("Examples").Bold(), `:
  1. tbg run --stretch fill
     whatever value the \"stretch\" field in the currently used config will
     be ignored and tbg will upscale the image to exactly fill the screen instead

  2. tbg edit --stretch fill config /path/to/a/config.yaml
     this will change the "stretch" field on the config /path/to/a/config.yaml
     to "fill"
`)
	}
}

func OpacityHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  -o, --opacity").Bold(),
		"The opacity of the image to use in a command.\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `:
  1. any float between 0 and 1 (inclusive)

  `, Decorate("Examples").Bold(), `:
  1. tbg run --opacity 0.5
     whatever value the "opacity" field in the currently used config will
     be ignored and tbg will set the image opacity to 0.5

  2. tbg edit --opacity 0.5 config /path/to/a/config.yaml
     this will change the "opacity" field on the config /path/to/a/config.yaml
     to 0.5
`)
	}
}

func RandomHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  -r, --random").Bold(),
		"Randomize image collections and images. Specific to `run` command\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `: --random does not take args

  `, Decorate("Examples").Bold(), `:
  1. tbg run --random
     This will randomize the order of the image collections read
     from the currently used config. It randomizes the order the
     images in each collection.

     When image collections are exhausted and tbg wraps around, the
     order of the image collections will be randomized again. This
     behavior applies to images too.
`)
	}
}
