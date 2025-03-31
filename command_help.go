package main

import (
	"fmt"
)

type HelpCommand struct {
	Commands []Command
}

func (cmd *HelpCommand) Type() CommandType { return HelpCommandType }

func (cmd *HelpCommand) String() {
	fmt.Println("Help Command")
	if len(cmd.Commands) > 0 {
		fmt.Println("Commands:")
	}
	for _, c := range cmd.Commands {
		fmt.Println(" ", c.Type())
	}
}

func (cmd *HelpCommand) ValidateValue(val *string) error {
	if val == nil || *val == "" {
		return nil
	}
	return fmt.Errorf("'help' takes no args. got: '%s'", *val)
}

func (cmd *HelpCommand) ValidateFlag(f Flag) error {
	return fmt.Errorf("'help' takes no flags. got: '%s'", f.Type)
}

func (cmd *HelpCommand) ValidateSubCommand(sc Command) error {
	cmd.Commands = append(cmd.Commands, sc)
	return nil
}

func (cmd *HelpCommand) Execute() error {
	if len(cmd.Commands) == 0 {
		fmt.Printf("%-37s%s\n",
			Decorate("tbg").Bold().Underline(),
			Decorate("Terminal Background Gallery").Italic(),
		)
		fmt.Printf("%-30s%s\n", "",
			fmt.Sprintf("%s (%s) allows the user to cycle through multiple background",
				Decorate("tbg").Bold(),
				Decorate("teabag").Italic(),
			),
		)
		fmt.Printf("%-30s%s\n", "",
			"images at a set interval for Windows Terminal.",
		)
		fmt.Printf("%-37s%s\n",
			Decorate("Version").Bold().Underline(),
			Decorate(TbgVersion).Italic(),
		)
		fmt.Printf("\n%-37s%s\n",
			Decorate("Usage").Bold().Underline(),
			Decorate("tbg [command] <flags>").Italic(),
		)
		fmt.Printf("\n%s:\n",
			Decorate("Commands").Bold().Underline(),
		)
		RunHelp(false)
		NextImageHelp(false)
		SetImageHelp(false)
		QuitHelp(false)
		AddHelp(false)
		RemoveHelp(false)
		ConfigHelp(false)
		HelpHelp(false)
		VersionHelp(false)
		return nil
	}

	// verbose messages
	fmt.Println("------------------------------------------------------------------------------------")
	for _, subCmd := range cmd.Commands {
		switch subCmd.Type() {
		case AddCommandType:
			AddHelp(true)
		case ConfigCommandType:
			ConfigHelp(true)
		case HelpCommandType:
			HelpHelp(true)
		case RemoveCommandType:
			RemoveHelp(true)
		case RunCommandType:
			RunHelp(true)
		case VersionCommandType:
			VersionHelp(true)
		case NextImageCommandType:
			NextImageHelp(true)
		case SetImageCommandType:
			SetImageHelp(true)
		case QuitCommandType:
			QuitHelp(true)
		}
		fmt.Println("------------------------------------------------------------------------------------")
	}
	return nil
}

func RunHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  run").Bold(),
		"Starts the tbg server that changes the background image at an interval\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `: run takes no args

  `, Decorate("Subcommands").Bold(), `: run takes no sub-commands
  `, Decorate("Flags").Bold(), `:

  You can specify alignment, stretch, and opacity using flags.
  These will override the values in the used config (not edit)
  1. -a, --alignment [arg]
         [top, topLeft, topRight, left, center, right, bottomLeft, bottom, bottomRight]
  2. -o, --opacity   [arg]
         [any float between 0 and 1 (inclusive)]
  3. -s, --stretch   [arg]
         [fill, none, uniform, uniformToFill]
  4. -p, --profile   [arg]
         [default, n, profile name]
         Where n is the list index Windows Terminal uses to identify the profile (starting from 1).
         Can specify profile name as well: e.g. "pwsh" (case insensitive)
  5. -P, --port   [arg]
         [any positive integer]
         port to be used by tbg server to listen to POST requests
  6. -i, --interval  [arg]
         [any positive integer]
         note that this is in minutes

  `, Decorate("Key Events").Bold(), `:
  while tbg is running, it accepts optional key events.
  Press a key to execute the command
  1. q: [q]uit tbg
  2. n: goes to [n]ext image
  3. c: list all [c]ommands

  `, Decorate("Examples").Bold(), `:
  1. tbg run
     This will use tbg's config values to edit Windows Terminal's settings.json

  2. tbg run --profile 2 --interval 5 --alignment center
      used_config                      values used to edit settings.json
      --------------------------       --------------------------------
      | paths:                         | paths:
      |   - path: /path/to/images/dir1 |   - path: /path/to/images/dir1
      |   - path: /path/to/images/dir2 |   - path: /path/to/images/dir2
      |                                |
      | profile: default               | profile: 2
      | port: 9545                     | port: 9545
      | interval: 30                   | interval: 5
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
		"Adds a path containing images to tbg's config\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `:
  1. path/to/images/dir
     Path to images dir should have at least one image
     file under it. All subdirectories will be ignored.

  `, Decorate("Subcommands").Bold(), `: add takes no sub-commands
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
		"Removes a path from tbg's config\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `:
  1. path/to/images/dir

  `, Decorate("Subcommands").Bold(), `: remove takes no sub-commands
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
		"Prints tbg's config if no flags.\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `: config takes no args
  `, Decorate("Subcommands").Bold(), `: config takes no sub-commands
  `, Decorate("Flags").Bold(), `:
  1. -p, --profile   [arg]
         [default, n, profile name]
         Where n is the list index Windows Terminal uses to identify the profile (starting from 1).
         Can specify profile name as well: e.g. "pwsh" (case insensitive)
  2. -P, --port   [arg]
         [any positive integer]
         port to be used by tbg server to listen to POST requests
  3. -i, --interval  [arg]
         [any positive integer]
         note that this is in minutes.

  `, Decorate("Examples").Bold(), `:
  1. tbg config
      print tbg's config:
      --------------------------
      | paths:
      |   - path: /path/to/images/dir
      |
      | profile: default
      | port: 9545
      | interval: 30
      |
      --------------------------

  2. tbg config --profile 1
      replaces the config's "profile" field with the value "1"
      --------------------------       --------------------------------
      | paths:                         | paths:
      |   - path: /path/to/images/dir1 |   - path: /path/to/images/dir1
      |                                |
      | profile: default               | profile: 1
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
  1. command names
     allows multiple commands to be specified

  `, Decorate("Examples").Bold(), `:
  1. tbg help
     Prints this help message
  2. tbg help run
     Prints verbose help for the 'run' command
  3. tbg help set-image next-image
     Prints verbose help for the 'set-image' and 'next-image' command
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
  `, Decorate("Args").Bold(), `: version does not take args

  `, Decorate("Examples").Bold(), `:
  1. tbg version
`)
	}
}

func NextImageHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  next-image").Bold(),
		"Triggers an image change on the currently running tbg server\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `: next-image does not take args

  `, Decorate("Flags").Bold(), `:
  You can specify alignment, stretch, and opacity using flags. See example 2 and 3
  1. -a, --alignment [arg]
         [top, topLeft, topRight, left, center, right, bottomLeft, bottom, bottomRight]
  2. -o, --opacity   [arg]
         [any float between 0 and 1 (inclusive)]
  3. -s, --stretch   [arg]
         [fill, none, uniform, uniformToFill]

  `, Decorate("Examples").Bold(), `:
  1. tbg next-image
`)
	}
}

func SetImageHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  set-image").Bold(),
		"Sets the specified image as the background image\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `:
  1. path/to/image/file
     Path to the image file you want to set as the background image

  `, Decorate("Flags").Bold(), `:
  You can specify alignment, stretch, and opacity using flags. See example 2 and 3
  1. -a, --alignment [arg]
         [top, topLeft, topRight, left, center, right, bottomLeft, bottom, bottomRight]
  2. -o, --opacity   [arg]
         [any float between 0 and 1 (inclusive)]
  3. -s, --stretch   [arg]
         [fill, none, uniform, uniformToFill]

  `, Decorate("Examples").Bold(), `:
  1. tbg next-image
`)
	}
}

func QuitHelp(verbose bool) {
	fmt.Printf("%-33s%s",
		Decorate("  quit").Bold(),
		"Stops the currently running tbg server\n",
	)
	if verbose {
		fmt.Print(`
  `, Decorate("Args").Bold(), `: quit does not take args

  `, Decorate("Examples").Bold(), `:
  1. tbg quit
`)
	}
}
