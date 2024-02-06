# Table of Contents
- [tbg](#tbg-Terminal-Background-Gallery)
- [Installation](#installation)
- [Usage](#usage)
- [tbg Profile](#tbg-profile)
- [Config](#config)
    - [Fields](#fields)
- [Commands](#commands)
- [Flags](#flags)
- [Credits](#credits)
# tbg (Terminal Background Gallery)
**tbg** (*teabag*) allows the user to have and manage multiple background images that rotate at a set amount of time for Windows Terminal

This edits the `settings.json` used by *Windows Terminal*; specifically, the `backgroundImage` on the default profile by default but user can override it. It overwrites the `backgroundImage` value every 30 minutes by default but the user can override that too (as well as image alignment, stretch, and opacity)

# Installation
Download the latest release of [tbg](https://github.com/saltkid/tbg/releases)
or build it from source
- clone the repo
```
git clone git@github.com:saltkid/tbg.git
```
- build it:
```
cd tbg && go build
```
**Optionally** add the `tbg` executable to your path

# Usage
```
tbg run
```
On initial execution of **tbg**, it will create a [tbg_profile](#tbg_profile) (`tbg_profile.yaml`) and a [default config](#config) (`config.yaml`) in the same directory as the executable. You can edit this manually **or** use [tbg commands](#commands) to edit these for input validation

While **tbg** is running, it takes in optional user input through key presses. Here's a list of commands:
- `q`: quit
- `n`: goes to next image in the current image collection dir
- `p`: goes to previous image in the current image collection dir
- `f`: goes to next image collection dir
- `b`: goes to previous image collection dir
- `c`: shows the available commands

See [fields](#fields) for more information about image collection dirs.

The image collection dirs wrap around so if you go past the last image in the last image collection dir, it will go back to the first image in the first image collection dir. Same goes with the reverse direction.

When changing image collection dirs, it always starts at the first image to avoid disorientation.
# tbg Profile
This is what is used by **tbg** to keep track of what config to use when it runs. There can only be one of this and it must be in the same directory as the **tbg** executable.

`tbg_profile.yaml` will be automatically created on the initial execution of **tbg**, along with the default `config.yaml`. The used config will point to the newly created `config.yaml`

You can point to your own config by editing `tbg_profile.yaml` manually or using the [config command](#link). 

```
#---------------------------------------------
# this is a tbg profile. Whenver tbg is ran, it will
# load this profile to get the currently used config
#
# currently, it only has one field: used_config
# I'll add more if the need arises
#---------------------------------------------

used_config: ""

#---------------------------------------------
# Fields:
#   used_config: path to the config used by tbg
#------------------------------------------
```

# Config
This is what is used by **tbg** to edit the `settings.json` *Windows Terminal* uses. As stated earlier, **tbg** created a `config.yaml` in the same path as the **tbg** executable. This is what it should look like:
```
#------------------------------------------
# this is a tbg config. Whenver tbg is ran, it will load this config file
# if it's the config in tbg's profile and use the fields below to control
# the behavior of tbg when changing background images of Windows Terminal
#
# to use your own config file, use the 'config' command:
#   tbg config path/to/config.yaml
#------------------------------------------

profile: default
interval: 30

image_col_paths: []

default_alignment: center
default_stretch: uniform
default_opacity: 0.1

#------------------------------------------
# Fields:
#   profile: profile profile in Windows Terminal (default, list-0, list-1, etc...)
#      see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general for more information
#
#   interval: time in minutes between each image change
#
#   image_col_paths: list of image collection paths
#      notes:
#        - put directories that contain images, not image filepaths
#        - can override default options for a specific path by putting a | after the path
#          and putting "alignment", "stretch", and "opacity" after the |
#          eg. abs/path/to/images/dir | right uniform 0.1
#
#------------------------------------------
#   Below are default options which can be overriden on a per-path basis by putting a pipe (|)
#   after the path and putting "alignment", "stretch", and "opacity" values after the | in order
#
#   example: abs/path/to/images/dir | right uniform 0.1
#
#   whatever values the values below have, the options after the | will override
#   the values in the default values for that specific path
#------------------------------------------
#
#   default_alignment: image alignment in Windows Terminal (left, center, right, etc...)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-alignment for more information
#
#   default_opacity: image opacity of background images in Windows Terminal (0.0 - 1.0)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-opacity for more information
#
#   default_stretch: image stretch in Windows Terminal (uniform, fill, etc...)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-stretch-mode for more information
#
#------------------------------------------
```
## Fields
Although you can edit the fields in the config directly, it is recommended to use the command `config` to edit them.
| Field | Valid Values | Description |
| --- | --- | --- |
| `profile` | `default`, `list-0`, `list-1` | target profile in *Windows Terminal*. To change background images in user created profiles, set `profile` to `list-<n>` where n is the index used by *Windows Terminal* to identify the profile.<br><br>See [Microsoft's documentation](https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general) for more information |
| `interval` | any positive integer | time in minutes between each image change. |
| `image_col_paths` | `[]`<br>`- path/to/dir1`<br>`- path/to/dir2 \| center uniform 0.1` | list of image collection paths. Must be directories containing images, not image paths.<br><br>Each dir can override the default fields by putting all 3 options after a pipe `\|`. Example:<br>`path/to/dir \| center fill 0.2` |
| `default_alignment` | `top`, `top-left`, `top-right`, `left`, `center`, `right`, `bottom`, `bottom-left`, `bottom-right` | image alignment in Windows Terminal. Can be overriden on a per-dir basis. See valid values of `image_col_paths` |
| `default_stretch` | `uniform`, `fill`, `uniform-fill`, `none` | image stretch in Windows Terminal. Can be overriden on a per-dir basis |
| `default_opacity` | inclusive between `0` and `1` | image opacity of background images in Windows Terminal. Can be overriden on a per-dir basis |

For the default fields, see [Mircrosoft's documentation](https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-images-and-icons)

# Commands
The commands other than `run` are used to safely edit the currently used config with input validation. If your currently used config is a user config, it will automatically edit the default config as well appropriately.

For a more detailed explanation on each types, follow the command name links

| Command | Valid Args | Valid Flags/Subcommands | Description |
| --- | --- | --- | --- |
| [run](https://github.com/saltkid/tbg/blob/main/docs/run_command_usage.md) | none | `config`, `--profile`, `--interval`, `--alignment`, `--opacity`, `--stretch` | edit `settings.json` used by *Windows Terminal* using settings from the currently used config.<br><br>If `config` is specified, it will use that instead of the currently used config.<br><br>If any of the `--`flags are specified, it will use those values in editing `settings.json` instead of what's specified in the currently used config|
| [config](#link) | `path/to/user-config-name.yaml` `default` | none | If no arg is given to `config`, it will print out the currently used config to console.<br><br>If an arg is given, it will set the arg as the currently used config after validation.<br><br>If used as a flag, it will let the parent command use te config specified instead of the currently set config|
| [add](#link) | `path/to/dir` | `config`, `--alignment`, `--opacity`, `--stretch` | Add a dir containing images to the currently used config.<br><br>If `config` is specified, it will add the dir to the specific config instead.<br><br>If any of the `--`flags are present, it will add those field values after the dir: eg `path/to/dir \| center fill 0.5`<br>This will override the default flag fields in the config.
| [remove](#link) | `path/to/dir` | `config` | Remove a dir from the currently used config.<br><br>If `config` is specified, it will remove the dir in the specific config instead. |
| [edit](#link) | `path/to/dir`<br>`fields`<br>`all` | `config`, `--profile`, `--interval`, `--alignment`, `--opacity`, `--stretch` | Edits the flags of the path specified. A path can individually have fields to override the default values in the config.<br>`path/to/dir \| center fill 0.2`<br><br>You can also specify `all` to edit all paths to have the flags you specified.<br><br>If you want to edit the default fields (`default_alignment`, `default_stretch`, and `default_opacity`), the arg should be `fields`<br><br>`--profile` and `--interval` are always edited on a per config basis, not per path since paths only take `--alignment`, `--stretch`, and `--opacity` options. |

# Flags
Flags are used to override [field entries in a config](#config), which are then passed to the parent command.

Flags behave differently based on the main command so for more detailed explanation, go to the documentation of the command instead.
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
