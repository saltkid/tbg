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
**tbg** (*teabag*) allows the user to have and manage multiple background images that rotate at a set amount of time for Windows Terminal

This edits the `settings.json` used by *Windows Terminal*; specifically, the `backgroundImage` on the default profile by default but user can override it. It overwrites the `backgroundImage` value every 30 minutes by default but the user can override that too (as well as image alignment, stretch, and opacity)

# Installation
Download the latest release of [tbg](https://github.com/saltkid/tbg/releases)

or

## Building from source
- clone the repo
```
git clone git@github.com:saltkid/tbg.git
```
- build it:
```
cd tbg && go mod tidy && go build
```
**Optionally** add the `tbg` executable to your path

# [Usage](https://github.com/saltkid/tbg/blob/main/docs/run_command_usage.md)
```
tbg run
```
On initial execution of **tbg**, it will create a [default config](#config) (`config.yaml`) in the same directory as the executable. You can edit this manually **or** use [tbg commands](#commands) to edit these with input validation

While **tbg** is running, it takes in optional user input through key presses. Here's a list of commands:
- `q`: quit
- `c`: shows the available commands
- `n`: goes to next image in the current image collection
- `p`: goes to previous image in the current image collection
- `r`: randomize images
- `N`: goes to next image collection
- `P`: goes to previous image collection
- `R`: randomize image collections

See [fields](#fields) for more information about image collection dirs.

The image collection dirs wrap around so if you go past the last image in the last image collection dir, it will go back to the first image in the first image collection dir. Same goes with the reverse direction.

When changing image collection dirs, it always starts at the first image to avoid disorientation.

# Config
This is what is used by **tbg** to edit the `settings.json` *Windows Terminal* uses. As stated earlier, **tbg** creates a `config.yaml` in the same path as the **tbg** executable on intial execution.

## Fields
Although you can edit the fields in the config directly, it is recommended to use the `config` command to edit them.
| Field | Valid Values | Description |
| --- | --- | --- |
| `profile` | `default`, `list-0`, `list-1` | target profile in *Windows Terminal*.<br><br>See [Microsoft's documentation](https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general) for more information |
| `interval` | any positive integer | time in minutes between each image change. |
| `image_col_paths` | `[]`<br>`- path/to/dir1`<br>`- path/to/dir2 \| center uniform 0.1` | list of image collection paths. Must be directories containing images, not image paths.<br><br>Each dir can override the default fields by putting all 3 options after a pipe `\|`. Example:<br>`path/to/dir \| center fill 0.2` |
| `default_alignment` | `top`, `top-left`, `top-right`, `left`, `center`, `right`, `bottom`, `bottom-left`, `bottom-right` | image alignment in Windows Terminal. Can be overriden on a per-dir basis. |
| `default_stretch` | `uniform`, `fill`, `uniform-fill`, `none` | image stretch in Windows Terminal. Can be overriden on a per-dir basis |
| `default_opacity` | inclusive range between `0` and `1` | image opacity of background images in Windows Terminal. Can be overriden on a per-dir basis |

For the default flag fields (`default_alignment`, `default_stretch`, and `default_opacity`), see [Mircrosoft's documentation](https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-images-and-icons) for more information

# Commands
The commands other than `run` are used to safely edit the config with input validation. For a more detailed explanation on each types, follow the command name links

| Command | Valid Args | Valid Flags/Subcommands | Description |
| --- | --- | --- | --- |
| [run](https://github.com/saltkid/tbg/blob/main/docs/run_command_usage.md) | none | `config`, `--profile`, `--interval`, `--alignment`, `--opacity`, `--stretch` | edit `settings.json` used by *Windows Terminal* using settings from the currently used config.<br><br>If `config` is specified, it will use that instead of the currently used config.<br><br>If any of the `--`flags are specified, it will use those values in editing `settings.json` instead of what's specified in the currently used config|
| [config](https://github.com/saltkid/tbg/blob/main/docs/config_command_usage.md) | `path/to/user-config-name.yaml` `default` | none | If no arg is given to `config`, it will print out the currently used config to console.<br><br>If an arg is given, it will set the arg as the currently used config after validation.<br><br>If used as a flag, it will let the parent command use te config specified instead of the currently set config|
| [add](https://github.com/saltkid/tbg/blob/main/docs/add_command_usage.md) | `path/to/dir` | `config`, `--alignment`, `--opacity`, `--stretch` | Add a dir containing images to the currently used config.<br><br>If `config` is specified, it will add the dir to the specific config instead.<br><br>If any of the `--`flags are present, it will add those field values after the dir: eg `path/to/dir \| center fill 0.5`<br>This will override the default flag fields in the config.
| [remove](https://github.com/saltkid/tbg/blob/main/docs/remove_command_usage.md) | `path/to/dir` | `config` | Remove a dir from the currently used config.<br><br>If `config` is specified, it will remove the dir in the specific config instead. |

## Flags
Flags are used to override [field entries in a config](#config), which are then passed to the parent command.

Flags behave differently based on the parent command so for more detailed explanation, go to the documentation of the command instead.
| Flag | Field Overriden |
| --- | --- |
| `--profile`<br>`-p` | `profile` |
| `--interval`<br>`-i` | `interval` |
| `--alignment`<br>`-a` | `default_alignment` |
| `--opacity`<br>`-o` | `default_opacity` |
| `--stretch`<br>`-s` | `default_stretch` |

---
# Credits
- [Windows Terminal](https://github.com/microsoft/terminal)
- [keyboard](https://github.com/eiannone/keyboard) for handling key events
- [saltkid](https://github.com/saltkid)
