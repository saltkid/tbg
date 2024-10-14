# Table of Contents
- [tbg](#tbg-Terminal-Background-Gallery)
- [Installation](#installation)
- [Usage](#usage)
    - [Building from source](#building-from-source)
- [Config](#config)
    - [Fields](#fields)
- [Commands](#commands)
    - [Flags](#flags)
- [Credits](#credits)

---

# tbg (Terminal Background Gallery)
**tbg** (*teabag*) allows the user to have and manage multiple background images
that rotate at a set amount of time for Windows Terminal

This edits the `settings.json` used by *Windows Terminal*; specifically, the
`backgroundImage` on the default profile by default but user can override it.
It overwrites the `backgroundImage` value every 30 minutes by default but the
user can override that too (as well as image alignment, stretch, and opacity)

# Installation
Download the latest release of [tbg](https://github.com/saltkid/tbg/releases)

or

## Building from source
- clone the repo
```
git clone git@github.com:saltkid/tbg.git
cd tbg
```
- build it:
```
go mod tidy
go build -ldflags '-s'
```
**Optionally** add the `tbg` executable to your path

# [Usage](https://github.com/saltkid/tbg/blob/main/docs/run_command_usage.md)
```
tbg run
```
On initial execution of **tbg**, it will create a [default config](#config)
(`.tbg.yml`) in the same directory as the executable. You can edit this
manually **or** use [tbg commands](#commands) to edit these with input
validation

While **tbg** is running, it takes in optional user input through key presses.
Here's a list of commands:
- `q`: quit
- `c`: shows the available commands
- `n`: goes to next image in the current images path
- `N`: goes to next images path
- `p`: goes to previous image in the current images path
- `P`: goes to previous images path
- `r`: randomize images
- `R`: randomize images paths

See [fields](#fields) for more information about images paths.

The images paths wrap around so if you go past the last image in the last path,
it will go back to the first image in the first images path. Same goes for the
reverse direction.

# Config
`.tbg.yml` is what is used by **tbg** to edit the `settings.json` *Windows
Terminal* uses.

## Fields
Although you can edit the fields in the config directly, it is recommended to
use the [`config`](https://github.com/saltkid/tbg/blob/main/docs/config_command_usage.md) command to edit them.
1. **profile**
    - *args*: `default`, `list-0`, `list-1`, ...
    - target profile in *Windows Terminal*.
2. **interval**
    - *args*: any positive integer 
    - time in minutes between each image change.
3. **paths** 
    - *args*:
        - `[]`
        - `- path: path/to/dir1` 
        - ```yaml
            - path: path/to/dir2
              alignment: center      # optional
              stretch: uniformToFill # optional
              opacity: 1.0           # optional
    - paths containing images used in changing the background image of Windows Terminal
4. **alignment**
    - ( args ): `top`, `topLeft`, `topRight`, `left`, `center`, `right`, `bottom`, `bottomLeft`, `bottomRight` 
    - image alignment in Windows Terminal.
    - Can be overriden on a per-path basis
5. `stretch` 
    - *args*: `uniform`, `fill`, `uniformToFill`, `none` 
    - image stretch in Windows Terminal. Can be overriden on a per-path basis
6. `opacity` 
    - *args*: inclusive range between `0` and `1` 
    - image opacity of background images in Windows Terminal.
    - Can be overriden on a per-path basis

# Commands
For a more detailed explanation on each command, follow the command name links

1. [run](https://github.com/saltkid/tbg/blob/main/docs/run_command_usage.md) 
    - *args*: none 
    - *flags*: `-p, --profile`, `-i, --interval`, `-a, --alignment`, `-o, --opacity`, `-s, --stretch`, `-r, --random`
    - edit `settings.json` used by *Windows Terminal* using settings from
    `.tbg.yml`. If any of the flags are specified, it will use those values in
    editing `settings.json` instead of what's specified in the currently used
    config
2. [config](https://github.com/saltkid/tbg/blob/main/docs/config_command_usage.md) 
    - *args*: none 
    - *flags*: `-p, --profile`, `-i, --interval`, `-a, --alignment`, `-o, --opacity`, `-s, --stretch` 
    - If no flags are present, it will print out the currently used config to
    console. If any of the flags are present, it will edit the fields of the
    config based on the flags and values passed
3. [add](https://github.com/saltkid/tbg/blob/main/docs/add_command_usage.md) 
    - `path/to/dir` 
    - *flags*: `-a, --alignment`, `-o, --opacity`, `-s, --stretch` 
    - Add a path containing images to the currently used config.
    - If any flags are present, it will add those options to that path,
    regardless of whether the path exists or not
4. [remove](https://github.com/saltkid/tbg/blob/main/docs/remove_command_usage.md) 
    - `path/to/dir` 
    - *flags*: `-a, --alignment`, `-o, --opacity`, `-s, --stretch` 
    - Remove a path from the currently used config.
    - If any flags are present, it will remove only those options of that path
5. help
    - args: no arg, or any command or any flag (can be multiple)
    - Prints the general help message when no arg is given
    - Prints the help message/s of command/s and/or flag/s if specified
## Flags
Flags are used to override [field entries in a config](#config), which are then
passed to the parent command.

Flags behave differently based on the parent command so for more detailed
explanation, go to the documentation of the command instead.
| Flag | Field Overriden |
| --- | --- |
| `-p, --profile` | `profile` |
| `-i, --interval` | `interval` |
| `-a, --alignment` | `alignment` |
| `-o, --opacity` | `opacity` |
| `-s, --stretch` | `stretch` |

---
# Credits
- [Windows Terminal](https://github.com/microsoft/terminal)
- [keyboard](https://github.com/eiannone/keyboard) for handling key events
- [saltkid](https://github.com/saltkid)
