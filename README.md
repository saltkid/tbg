# Table of Contents
- [tbg](#tbg-Terminal-Background-Gallery)
- [Installation](#installation)
    - [Building from source](#building-from-source)
- [Usage](#usage)
    - [Automatically change background image at a set interval](#automatically-change-background-image-at-a-set-interval)
    - [tbg server](#tbg-server)
- [Config](#config)
    - [Fields](#fields)
- [Commands](#commands)
    - [Server Commands](#server-commands)
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
## Automatically change background image at a set interval
```
tbg run
```
This will start the **tbg** server. This will target Windows Terminal's default
profile, changing the background image every 30 minutes, choosing an random
image from the specified `paths` in `.tbg.yml`. See [fields](#fields) for more
information about images paths. The default image properties are:

| image property | default value  |
|----------------|----------------|
| alignment      | `center`       |
| opacity        | `1.0`          |
| stretch        | `uniformToFill`|

On initial execution of `tbg run`, it will create a [default config](#config)
(`.tbg.yml`) in the same directory as the executable. You can edit this
manually **or** use [tbg commands](#commands) to edit these with input
validation

## tbg server
`tbg run` starts the **tbg** http server where POST requests can be made to
trigger certain actions. These post requests can be made through [tbg server
commands](#server-commands) for convenience

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
3. **port**
    - port that the tbg server uses
    - *args*: any positive integer
4. **paths** 
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
```zsh
tbg <command-name>
```
For a more detailed explanation on each command, follow the command name links

1. [run](https://github.com/saltkid/tbg/blob/main/docs/run_command_usage.md) 
    - edit `settings.json` used by *Windows Terminal* using settings from
    `.tbg.yml`. If any of the flags are specified, it will use those values in
    editing `settings.json` instead of what's specified in the `.tbg.yml`
    - *arg*: none 
    - *flags*: `-p, --profile`, `-i, --interval`, `-a, --alignment`,
    `-o, --opacity`, `-s, --stretch`
2. [config](https://github.com/saltkid/tbg/blob/main/docs/config_command_usage.md) 
    - If no flags are present, it will print out `.tbg.yml` to console. If any
    of the flags are present, it will edit the fields of the config based on
    the flags and values passed
    - *arg*: none 
    - *flags*: `-p, --profile`, `-P, --port`, `-i, --interval`
3. [add](https://github.com/saltkid/tbg/blob/main/docs/add_command_usage.md) 
    - Add a path containing images to `.tbg.yml`.
    - If any flags are present, this will add those options to that path,
    regardless of whether the path exists or not
    - *arg*: `/path/to/dir` 
    - *flags*: `-a, --alignment`, `-o, --opacity`, `-s, --stretch` 
4. [remove](https://github.com/saltkid/tbg/blob/main/docs/remove_command_usage.md) 
    - Remove a path from `.tbg.yml`.
    - If any flags are present, it will remove only those options of that path
    - *arg*: `/path/to/dir` 
    - *flags*: `-a, --alignment`, `-o, --opacity`, `-s, --stretch` 
5. help
    - Prints the general help message when no arg is given
    - Prints the help message/s of command/s and/or flag/s if specified
    - *arg*: no arg, or any command (can be multiple)

## [Server Commands](https://github.com/saltkid/tbg/blob/main/docs/server_commands_usage.md)
These commands only work when there's a **tbg** server active. Usage is the
same as other commands.
1. next-image
    - triggers an image change
    - *arg*: `/path/to/dir` 
    - *flags*: `-a, --alignment`, `-o, --opacity`, `-s, --stretch`
2. set-image
    - sets a specified image as the background image
    - *arg*: `/path/to/image/file` 
    - *flags*: `-a, --alignment`, `-o, --opacity`, `-s, --stretch`
3. quit
    - stops the server
    - *arg*: none
    - *flags*: none

*Tip: you can assign these commands to keybinds*

---
# Credits
- [Windows Terminal](https://github.com/microsoft/terminal)
- [levenshtein](github.com/agnivade/levenshtein) for suggesting similar profile names
- [saltkid](https://github.com/saltkid)
