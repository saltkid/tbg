# Table of Contents
- [tbg](#tbg-Terminal-Background-Gallery)
- [Installation](#installation)
    - [Building from source](#building-from-source)
- [Usage](#usage)
    - [Interactive Commands](#interactive-commands)
- [Config](#config)
    - [Fields](#fields)
- [Commands](#commands)
    - [Flags](#flags)
- [Credits](#credits)

---

# tbg (Terminal BackGround)

**tbg** (*pronounced /ˈtiːbæɡ/*) is a tool for Windows Terminal that allows you
to cycle through multiple background images at a set interval. The image
sources, interval of changing, image alignment, etc. can be configured
via **tbg**'s config file: `.tbg.yml`.

This edits the `settings.json` used by *Windows Terminal*; specifically, the
`backgroundImage` on the default profile by default but user can specify which
profile to target. 

# Installation
Download the latest release of [tbg](https://github.com/saltkid/tbg/releases)

or

## Building from source
- clone the repo
```
git clone https://github.com/saltkid/tbg.git
cd tbg
```
- build it:
```
go mod tidy
go build -ldflags '-s'
```
**Optionally** add the `tbg` executable to your PATH

*oneliner for copy pasting*
```bash
git clone https://github.com/saltkid/tbg.git && cd tbg && go mod tidy && go build -ldflags '-s'
```

# [Usage](https://github.com/saltkid/tbg/blob/main/docs/run_command_usage.md)
```
tbg run
```
This will target Windows Terminal's default profile, changing the background
image every 30 minutes, choosing an random image from the specified `paths` in
`.tbg.yml`. The default image properties are:
1. alignment: `center`
2. opacity: `1.0`
3. stretch: `uniformToFill`

On initial execution of **tbg**, it will create a [default config](#config)
(`.tbg.yml`) in the same directory as the executable. You can edit this
manually **or** use [tbg commands](#commands) to edit these with input
validation

## Interactive Commands
While **tbg** is running, it takes in optional user input through key presses.
- `q`: quit
- `c`: shows the available commands
- `n`: goes to the next image (randomly chosen)

See [fields](#fields) for more information about images paths.

# Config
**tbg** uses `.tbg.yml` to edit the `settings.json` *Windows Terminal* uses.

## Fields
Although you can edit the fields in the config directly, it is recommended to
use the [`config`](https://github.com/saltkid/tbg/blob/main/docs/config_command_usage.md)
command to edit the file.
1. **profile**
    - target profile in *Windows Terminal*.
    - specify either the index used by Windows Terminal, or the profile name
    - *args*: `default`, `1`, `2`, ... `"profile name"`
2. **interval**
    - time in minutes between each image change.
    - *args*: any positive integer 
3. **paths** 
    - paths containing images used in changing the background image of Windows
    Terminal
    - *args*:
        - `[]`
        - `- path: /path/to/dir1` 
        - ```yaml
            - path: /path/to/dir2
              alignment: center      # optional
              stretch: uniformToFill # optional
              opacity: 1.0           # optional

# Commands
For a more detailed explanation on each command, follow the command name links

1. [run](https://github.com/saltkid/tbg/blob/main/docs/run_command_usage.md) 
    - edit `settings.json` used by *Windows Terminal* using settings from
    `.tbg.yml`. If any of the flags are specified, it will use those values in
    editing `settings.json` instead of what's specified in the `.tbg.yml`
    - *args*: none 
    - *flags*: `-p, --profile`, `-i, --interval`, `-a, --alignment`,
    `-o, --opacity`, `-s, --stretch`
2. [config](https://github.com/saltkid/tbg/blob/main/docs/config_command_usage.md) 
    - If no flags are present, it will print out `.tbg.yml` to console. If any
    of the flags are present, it will edit the fields of the config based on
    the flags and values passed
    - *args*: none 
    - *flags*: `-p, --profile`, `-i, --interval`
3. [add](https://github.com/saltkid/tbg/blob/main/docs/add_command_usage.md) 
    - Add a path containing images to `.tbg.yml`.
    - If any flags are present, this will add those options to that path,
    regardless of whether the path exists or not
    - *args*: `/path/to/dir` 
    - *flags*: `-a, --alignment`, `-o, --opacity`, `-s, --stretch` 
4. [remove](https://github.com/saltkid/tbg/blob/main/docs/remove_command_usage.md) 
    - Remove a path from `.tbg.yml`.
    - If any flags are present, it will remove only those options of that path
    - *args*: `/path/to/dir` 
    - *flags*: `-a, --alignment`, `-o, --opacity`, `-s, --stretch` 
5. help
    - Prints the general help message when no arg is given
    - Prints the help message/s of command/s and/or flag/s if specified
    - *args*: no arg, or any command or any flag (can be multiple)
## Flags
Flags behave differently based on the command so for more detailed explanation,
go to the documentation of the command instead.
1. `-a, --alignment`
2. `-i, --interval`
3. `-o, --opacity`
4. `-p, --profile`
5. `-s, --stretch`

---
# Credits
- [Windows Terminal](https://github.com/microsoft/terminal)
- [keyboard](https://github.com/eiannone/keyboard) for handling key events
- [saltkid](https://github.com/saltkid)
