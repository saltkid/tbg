# Table of Contents
- [tbg](#tbg-Terminal-Background-Gallery)
- [Installation](#installation)
- [Usage](#usage)
    - [Default Config](#default-config)
    - [User Config](#user-config)
    - [Commands](#commands)
    - [Flags](#flags)
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
On initial execution of **tbg**, it will create a [default](#default-config) `config.yaml` in the same directory as the executable. You can edit this manually **or** use [tbg commands](#commands) for input validation

While **tbg** is running, it takes in optional user input. Here's a list of commands:
- `n`: goes to next image in the current image collection dir
- `c`: goes to next image collection dir
- `h`: shows the available commands
- `q`: quits the program

After **tbg** exhausts the image collections dirs in the config, it will restart from the first image collection dir so the only way to quit is to either input `q` or press `ctrl+c` so just be aware of that.

## Default Config
This is what is used by **tbg** to edit the `settings.json` *Windows Terminal* uses. It also keeps track of the currently used config, whether it is this default config or a user config.
```
#------------------------------------------
# this is the default config. you can edit this or use your own config file by doing any of the following:
#  1. run 'tbg config <path/to/config.yaml>'
#  2. edit this file: set use_user_config to true and set user_config to the path to your config file
#------------------------------------------
user_config:

image_col_paths: []
interval: 30
profile: default

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
#------------------------------------------
# Fields:
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
### Fields
| Field | Valid Values | Description |
| --- | --- | --- |
| `user_config` | `""`, `path/to/config.yaml` | path to the user config file. Must be absolute path<br><br>If empty, currently used config will be the default config. If not empty, currently used config will be the value |
| `image_col_paths` | `[]`<br>`- path/to/dir1`<br>`- path/to/dir2 \| center uniform 0.1` | list of image collection paths. Must be directories containing images, not image paths |
| `interval` | any positive integer | time in minutes between each image change. |
| `profile` | `default`, `list-0`, `list-1` | target profile in *Windows Terminal*. To change background images in user created profiles, set `profile` to `list-<n>` where n is the index used by *Windows Terminal* to identify the profile |
| `default_alignment` | `top`, `top-left`, `top-right`, `left`, `center`, `right`, `bottom`, `bottom-left`, `bottom-right` | image alignment in Windows Terminal. Can be overriden on a per-dir basis. See valid values of `image_col_paths` |
| `default_stretch` | `uniform`, `fill`, `uniform-fill`, `none` | image stretch in Windows Terminal. Can be overriden on a per-dir basis |
| `default_opacity` | inclusive between `0` and `1` | image opacity of background images in Windows Terminal. Can be overriden on a per-dir basis |

## User Config
The default config will keep track of what config file you want to use so having many user config files are okay as long as you set `user_config` with a value.

### Fields
Other than not having `use_user_config` and `user_config` fields, User Config fields are exactly the same as Default Config

## Commands
The commands other than `run` are used to safely edit the currently used config with input validation. If your currently used config is a user config, it will automatically edit the default config as well appropriately.

| Command | Valid Args | Valid Flags/Subcommands | Description |
| --- | --- | --- | --- |
| `run` | none | `config`, `--profile`, `--interval`, `--alignment`, `--opacity`, `--stretch` | edit `settings.json` used by *Windows Terminal* using settings from the currently used config (may be default or a user config).<br><br>If `config` is specified, it will use that instead of the currently used config.<br><br>If any of the `--`flags are specified, it will use those values in editing `settings.json` instead of what's specified in the currently used config|
| `config` | `path/to/user-config-name.yaml` `default` | none | If no arg is given to `config`, it will print out the currently used config to console.<br><br>If an arg is given, it will set the arg as the currently used config after validation.<br><br>If used as a flag, it will let the parent command use te config specified instead of the currently set config|
| `add` | `path/to/dir` | `config`, `--alignment`, `--opacity`, `--stretch` | Add a dir containing images to the currently used config.<br><br>If `config` is specified, it will add the dir to the specific config instead.<br><br>If any of the `--`flags are present, it will add those field values after the dir: eg `path/to/dir \| center fill 0.5`<br>This will override the set fields in the config.
| `remove` | `path/to/dir` | `config` | Remove a dir from the currently used config.<br><br>If `config` is specified, it will remove the dir in the specific config instead. |

### Flags
Flags modify the [field entries](#default-config) in the currently used config. 
| Flag | Field Overriden |
| --- | --- |
| `--profile`<br>`-p` | `profile` |
| `--interval`<br>`-i` | `interval` |
| `--alignment`<br>`-a` | `default_alignment` |
| `--opacity`<br>`-o` | `default_opacity` |
| `--stretch`<br>`-s` | `default_stretch` |

---
For `run`, flags override field values in the currently used config, not edit them. These values will be used in editing `settings.json` of *Windows Terminal*

**example:**

start
```
# config.yaml

image_col_paths:
- /path/to/dir1
- /path/to/dir2

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```
When tbg is ran with `run` command, it will use the default values in the config to edit `settings.json` using images in `image_col_paths`

When `run --alignment top --stretch none --opacity 0.2` is called, it will use the values specified in the flags to edit `settings.json` **BUT** not edit the config.yaml file itself

---
For `add`, flags are added to a per-dir level, after `|`,  which individually can override the fields in their respective config. When that specific path is read during execution of `run` command, it will use the flags specific to that path, not the config-level default flags

**example:**

start
```
image_col_paths: []

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```
When tbg is ran with `add /path/to/dir1`, it will add `path/to/dir1` to `image_col_paths`
```
image_col_paths:
- /path/to/dir1

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```
When tbg is ran with `add /path/to/dir2 --alignment top`, it will add `path/to/dir2 | center uniform 0.1` to `image_col_paths`. When not all 3 flags are specified (`--alignment`, `stretch`, `opacity`), the omitted flags will inherit from the default fields
```
image_col_paths:
- /path/to/dir1
- /path/to/dir2 | top uniform 0.1

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```
When tbg is ran with `add /path/to/dir3 --alignment center`, it will add only `path/to/dir3` to `image_col_paths` because the `default_alignment` is already center and no other flags were specified so it will be using the default values anyway. This is why no flags were added to `path/to/dir3` even though `--alignment` was specified.
```
image_col_paths:
- /path/to/dir1
- /path/to/dir2 | top uniform 0.1
- /path/to/dir3

profile: default
interval: 30
default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```

---
For `edit`, it edits flag on a per-dir level after the `|`, same as `add`. If `--profile` or `--interval` is specified, it will edit that on a per-config level.

**example:**

start
```
image_col_paths:
- /path/to/dir1

profile: default
interval: 30
default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```
When tbg is ran with `edit --profile list-1 --interval 10`, it will edit the `profile` field to `default` and `interval` to `10`. These two fields specifically are set on a config-level and cannot be specified in a per-path level
```
image_col_paths:
- /path/to/dir1

profile: list-1
interval: 10
default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```
When tbg is ran with `edit /path/to/dir1 --alignment top`, it will edit the `/path/to/dir1` entry to have flags: `/path/to/dir1 | top uniform 0.1`. Since only alignment is specified, the rest of the flags will be inherited from the defaults
```
image_col_paths:
- /path/to/dir1 | top uniform 0.1

profile: list-1
interval: 10
default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```
When tbg is ran with `edit /path/to/dir1 --alignment center`, it will edit the `/path/to/dir1 | top uniform 0.1` entry to: `/path/to/dir1`. `default_alignment` is already center and no other flags were specified so `path/to/dir1` will be using default values anyway so it was removed.
```
image_col_paths:
- /path/to/dir1

profile: list-1
interval: 10
default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```
When tbg is ran with `edit fields --alignment right --opacity 0.5 --stretch none`, it will edit the default fields since `fields` was the argument for `edit`. The individual paths will remain unchanged
```
image_col_paths:
- /path/to/dir1

profile: list-1
interval: 10
default_alignment: right
default_stretch: none
default_opacity: 0.5
```
