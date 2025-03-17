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
#      valid values: default, 0, 1, ..., n
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
1. **profile**
    - *args*: `default`, `0`, `1`, ...
    - target profile in *Windows Terminal*.
    - To change background images in user created profiles, set `profile` to
    `<n>` where n is the index used by *Windows Terminal* to identify the
    profile.
    - See [Microsoft's documentation](https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general)
    for more information
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
    - Each path can override the default fields below.
    - default values for per-path options if not specified are:
        - `alignment: center`
        - `stretch: uniformToFill`
        - `opacity: 1.0`
4. **alignment**
    - ( args ): `top`, `topLeft`, `topRight`, `left`, `center`, `right`, `bottom`, `bottomLeft`, `bottomRight` 
    - image alignment in Windows Terminal.
    - Can be overriden on a per-path basis
5. `stretch` 
    - *args*: `uniform`, `fill`, `uniformToFill`, `none` 
    - image stretch in Windows Terminal. Can be overriden on a per-path basis |
6. `opacity` 
    - *args*: inclusive range between `0` and `1` 
    - image opacity of background images in Windows Terminal.
    - Can be overriden on a per-path basis

For the default flag fields (`alignment`, `stretch`, and `opacity`), see
[Mircrosoft's documentation](https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-images-and-icons)
for more information
