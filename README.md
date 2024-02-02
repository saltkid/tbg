# Table of Contents
- [What is tbg?](#what-is-tbg)
- [Installation](#installation)
- [Usage](#usage)
    - [Default Config](#default-config)
    - [User Config](#user-config)
    - [Commands](#commands)
    - [Flags](#flags)
# What is tbg?
**tbg** (*Terminal Background Gallery*) allows the user to have and manage multiple background images that rotate at a set amount of time for Windows Terminal

This edits the `settings.json` used by *Windows Terminal*; specifically, the `backgroundImage` on the default profile by default but user can override it. It overwrites the `backgroundImage` value every 30 minutes by default but the user can override that to

## Installation
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

## Usage
On initial execution of **tbg**, it will create a `config.yaml` in the same directory as the executable. You can edit this manually **or** use [tbg commands](#commands)
```
#------------------------------------------
# this is the default config. you can edit this or use your own config file by doing any of the following:
#  1. run 'tbg config <path/to/config.yaml>'
#  2. edit this file: set use_user_config to true and set user_config to the path to your config file
#------------------------------------------
use_user_config: false
user_config:

image_col_paths: []
interval: 30
profile: default

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
#------------------------------------------
# Fields:
#   use_user_config: whether to use the user config set in user_config
#
#   user_config: path to the user config file
#
#   image_col_paths: list of image collection paths
#      notes:
#        - put directories that contain images, not image filepaths
#        - can override default options for a specific path by putting a | after the path
#          and putting "alignment", "stretch", and "opacity" after the |
#          eg. /abs/path/to/images/dir | right uniform 0.1
#
#   interval: time in minutes between each image change
#
#   profile: target profile in Windows Terminal (default, list-0, list-1, etc...)
#      see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general for more information
#
#   ---
#   Below are default options which can be overriden on a per-path basis by putting a | after the path
#   and putting "alignment", "stretch", and "opacity" values after the |
#   ---
#   default_alignment: image alignment in Windows Terminal (left, center, right, etc...)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-alignment for more information
#
#   default_opacity: image opacity of background images in Windows Terminal (0.0 - 1.0)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-opacity for more information
#
#   default_stretch: image stretch in Windows Terminal (uniform, fill, etc...)
#     see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-stretch-mode for more information
#------------------------------------------
```
### Default Config
This is what is used by **tbg** to edit the `settings.json` *Windows Terminal* uses. It also keeps track of the currently used config, whether it is this default config or a user config.
#### Fields
| Field | Valid Values | Description |
| --- | --- | --- |
| `use_user_config` | `true`, `false` | This determines whether **tbg** uses this default config or the config in `user_config`. If you want to use what is set in `user_config`, you have to set this to false. |
| `user_config` | `""`, `path/to/config.yaml` | path to the user config file. Must be absolute path |
| `image_col_paths` | `[]`<br>`- path/to/dir1`<br>`- path/to/dir2 \| center uniform 0.1` | list of image collection paths. Must be directories containing images, not image paths |
| `interval` | any positive integer | time in minutes between each image change. |
| `profile` | `default`, `list-0`, `list-1` | target profile in *Windows Terminal*. To change background images in user created profiles, set `profile` to `list-<n>` where n is the index used by *Windows Terminal* to identify the profile |
| `default_alignment` | `top`, `top-left`, `top-right`, `left`, `center`, `right`, `bottom`, `bottom-left`, `bottom-right` | image alignment in Windows Terminal. Can be overriden on a per-dir basis. See valid values of `image_col_paths` |
| `default_stretch` | `uniform`, `fill`, `uniform-fill`, `none` | image stretch in Windows Terminal. Can be overriden on a per-dir basis |
| `default_opacity` | inclusive between `0` and `1` | image opacity of background images in Windows Terminal. Can be overriden on a per-dir basis |

### User Config
The default config will keep track of what config file you want to use so having many user config files are okay as long as you set these two in the default config:
- `use_user_config: true`
- `user_config: /abs/path/to/user-config-name.yaml`.

#### Fields
Other than not having `use_user_config` and `user_config` fields, User Config fields are exactly the same as Default Config

### Commands
The commands other than `run` are used to safely edit the currently used config with input validation. If your currently used config is a user config, it will automatically edit the default config as well appropriately.

*Note: Flags of commands can either be a Flag or another Command*
| Command | Valid Args | Flags | Description |
| --- | --- | --- | --- |
| `run` | none | `config`, `--profile`, `--interval`, `--alignment`, `--opacity`, `--stretch` | edit `settings.json` used by *Windows Terminal* using settings from the currently used config (may be default or a user config).<br><br>If `config` is specified, it will use that instead of the currently used config.<br><br>If any of the `--`flags are specified, it will use those values in editing `settings.json` instead of what's specified in the currently used config|
| `config` | `path/to/user-config-name.yaml` `default` | none | If no arg is given to `config`, it will print out the currently used config to console.<br><br>If an arg is given, it will set the arg as the currently used config after validation.<br><br>If used as a flag, it will let the parent command use te config specified instead of the currently set config|
| `add` | `path/to/dir` | `config`, `--alignment`, `--opacity`, `--stretch` | Add a dir containing images to the currently used config.<br><br>If `config` is specified, it will add the dir to the specific config instead.<br><br>If any of the `--`flags are present, it will add those field values after the dir: eg `path/to/dir \| center fill 0.5`<br>This will override the set fields in the config.
| `remove` | `path/to/dir` | `config` | Remove a dir from the currently used config.<br><br>If `config` is specified, it will remove the dir in the specific config instead. |

### Flags
Flags override, not edit, [field entries](#default-config) in the currently used config. 

Specifically for `add` command, flags are added to a per-dir level, after `|`,  which individually can override the fields in their respective config

Specifically for `edit` command, it edits flag on a per-dir level after the `|`, same as `add`. If `--profile` or `--interval` is specified, it will edit that on a per-config level.

| Flag | Field Overriden |
| --- | --- |
| `--profile`<br>`-p` | `profile` |
| `--interval`<br>`-i` | `interval` |
| `--alignment`<br>`-a` | `default_alignment` |
| `--opacity`<br>`-o` | `default_opacity` |
| `--stretch`<br>`-s` | `default_stretch` |

