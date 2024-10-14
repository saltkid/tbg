# Table of Contents
- [Config](#config)
- [Fields](#fields)

# Config
This is what is used by **tbg** to edit the `settings.json` *Windows Terminal*
uses. **tbg** creates a `.tbg.yml` in the same path as the **tbg** executable
on initial execution. This is what it should look like:
```
#------------------------------------------
# this is a tbg config. Whenver tbg is ran, it will load this config file
# and use the fields below to control the behavior of tbg when changing
# background images of Windows Terminal
#------------------------------------------

paths:
- path: path/to/images/dir
  # alignment: right # uncomment to override the alignment field for this specific image path
  # stretch: fill    # override stretch
  # opacity: 0.25    # override opacity

profile: default
interval: 30

alignment: center
stretch: uniform
opacity: 0.1

#------------------------------------------
# Fields:
#   paths: list of image collection paths
#     - path: directory that contain images
#       alignment: optional alignment value applied only to this path.
#                  uses default alignment (see "alignment" field) if not specified
#       stretch:   optional stretch value applied only to this path
#                  uses default stretch (see "stretch" field) if not specified
#       opacity:   poptional opacity value applied only to this path
#                  uses default opacity (see "opacity" field) if not specified
#
#   profile: profile profile in Windows Terminal
#      valid values: default, list-0, list-1, ..., list-n
#      https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general
#
#   interval: time in minutes between each image change
#
#   alignment: image alignment in Windows Terminal
#     valid values: topLeft, top, topRight, left, center, right, bottomLeft, bottom, bottomRight
#     https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-alignment
#
#   opacity: image opacity of background images in Windows Terminal
#     valid values: 0.0 - 1.0 (inclusive)
#     https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-opacity
#
#   stretch: image stretch in Windows Terminal
#     valid values: fill, none, uniform, uniformToFill
#     https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-stretch-mode 
#------------------------------------------
```
## Fields
Although you can edit the fields in the config directly, it is recommended to use the [`config` command](https://github.com/saltkid/tbg/blob/main/docs/config_command_usage.md) to edit them.
| Field | Valid Values | Description |
| --- | --- | --- |
| `profile` | `default`, `list-0`, `list-1` | target profile in *Windows Terminal*.<br><br>To change background images in user created profiles, set `profile` to `list-<n>` where n is the index used by *Windows Terminal* to identify the profile.<br><br>See [Microsoft's documentation](https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general) for more information |
| `interval` | any positive integer | time in minutes between each image change. |
| `paths` | `[]`<br>`- path: path/to/dir1` | list of image collection paths. Must be directories containing images, not image paths.<br><br>Each path can override the default fields below. See [add command](https://github.com/saltkid/tbg/blob/main/docs/add_command_usage.md) and [edit command](https://github.com/saltkid/tbg/blob/main/docs/edit_command_usage.md)
| `alignment` | `top`, `top-left`, `top-right`, `left`, `center`, `right`, `bottom`, `bottom-left`, `bottom-right` | image alignment in Windows Terminal. Can be overriden on a per-dir basis |
| `stretch` | `uniform`, `fill`, `uniform-fill`, `none` | image stretch in Windows Terminal. Can be overriden on a per-dir basis |
| `opacity` | inclusive range between `0` and `1` | image opacity of background images in Windows Terminal. Can be overriden on a per-dir basis |

For the default flag fields (`alignment`, `stretch`, and `opacity`), see [Mircrosoft's documentation](https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-images-and-icons) for more information
