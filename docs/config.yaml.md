# Table of Contents
- [Config](#config)
- [Fields](#fields)

# Config
This is what is used by **tbg** to edit the `settings.json` *Windows Terminal* uses. **tbg** creates a `config.yaml` in the same path as the **tbg** executable on initial execution. This is what it should look like:
```
#------------------------------------------
# this is a tbg config. Whenver tbg is ran, it will load this config file
# if it's the config in tbg's profile and use the fields below to control
# the behavior of tbg when changing background images of Windows Terminal
#
# to use your own config file, use the 'config' command:
#   tbg config path/to/config.yaml
#------------------------------------------

image_col_paths: []

profile: default
interval: 30

default_alignment: center
default_stretch: uniform
default_opacity: 0.1

#------------------------------------------
# Fields:
#   image_col_paths: list of image collection paths
#      notes:
#        - put directories that contain images, not image filepaths
#        - can override default options for a specific path by putting a | after the path
#          and putting "alignment", "stretch", and "opacity" after the |
#          eg. abs/path/to/images/dir | right uniform 0.1
#
#   profile: profile profile in Windows Terminal (default, list-0, list-1, etc...)
#      see https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general for more information
#
#   interval: time in minutes between each image change
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
Although you can edit the fields in the config directly, it is recommended to use the [`config` command](#link) to edit them.
| Field | Valid Values | Description |
| --- | --- | --- |
| `profile` | `default`, `list-0`, `list-1` | target profile in *Windows Terminal*.<br><br>To change background images in user created profiles, set `profile` to `list-<n>` where n is the index used by *Windows Terminal* to identify the profile.<br><br>See [Microsoft's documentation](https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general) for more information |
| `interval` | any positive integer | time in minutes between each image change. |
| `image_col_paths` | `[]`<br>`- path/to/dir1`<br>`- path/to/dir2 \| center uniform 0.1` | list of image collection paths. Must be directories containing images, not image paths.<br><br>Each dir can override the default fields by putting all 3 options after a pipe `\|`. See [add command](#link) and [edit command](#link)<br><br>Example:<br>`path/to/dir \| center fill 0.2` |
| `default_alignment` | `top`, `top-left`, `top-right`, `left`, `center`, `right`, `bottom`, `bottom-left`, `bottom-right` | image alignment in Windows Terminal.|
| `default_stretch` | `uniform`, `fill`, `uniform-fill`, `none` | image stretch in Windows Terminal. Can be overriden on a per-dir basis |
| `default_opacity` | inclusive range between `0` and `1` | image opacity of background images in Windows Terminal. Can be overriden on a per-dir basis |

For the default flag fields, see [Mircrosoft's documentation](https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-images-and-icons)
