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
  # alignment: right # override default alignment of center
  # stretch: fill    # override default stretch of uniformToFill
  # opacity: 0.25    # override default opacity of 1.0

profile: default
port: 9545
interval: 30

#------------------------------------------
# Fields:
#   paths: list of image collection paths
#     - path: directory that contain images
#       alignment: (optional) image alignment in Windows Terminal
#                  valid values: topLeft, top, topRight, left, center, right, bottomLeft, bottom, bottomRight
#                  https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-alignment
#
#       opacity:   (optional) image opacity of background images in Windows Terminal
#                  valid values: 0.0 - 1.0 (inclusive)
#                  https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-opacity
#
#       stretch:   (optional) image stretch in Windows Terminal
#                  valid values: fill, none, uniform, uniformToFill
#                  https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-image-stretch-mode 
#
#   profile: profile in Windows Terminal. can be "default" profile to select the
#            defaults field in wt settings, numbers greater than 0 to select a
#            specific profile by index, or can be any string to select a profile by name
#      valid values: default, 1, 2, ..., n, any string
#      https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general
#
#   interval: time in minutes between each image change
#------------------------------------------
```
## Fields
Although you can edit the fields in the config directly, it is recommended to use the [`config` command](https://github.com/saltkid/tbg/blob/main/docs/config_command_usage.md) to edit them.
1. **profile**
    - *args*: `default`, `1`, `2`, ..., `n`, any string
    - target profile in *Windows Terminal*.
    - To change background images in user created profiles, set `profile` to
    the name of a profile
    - Profiles can also be selected through profile number by setting `profile`
    to `<n>` where n is the index used by *Windows Terminal* to identify the
    profile. This is most useful when multiple profiles share the same name.
    - See [Microsoft's documentation](https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-general)
    for more information
2. **interval**
    - *args*: any positive integer 
    - time in minutes between each image change.
3. **port**
    - *args*: any positive integer
    - port that the tbg server uses to listen to POST requests
4. **paths** 
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
    1. **alignment**
        - *args*: `top`, `topLeft`, `topRight`, `left`, `center`, `right`, `bottom`, `bottomLeft`, `bottomRight` 
        - image alignment in Windows Terminal.
        - Can be overriden on a per-path basis
    2. `stretch` 
        - *args*: `uniform`, `fill`, `uniformToFill`, `none` 
        - image stretch in Windows Terminal. Can be overriden on a per-path basis |
    3. `opacity` 
        - *args*: inclusive range between `0` and `1` 
        - image opacity of background images in Windows Terminal.
        - Can be overriden on a per-path basis

For the default flag fields (`alignment`, `stretch`, and `opacity`), see
[Mircrosoft's documentation](https://learn.microsoft.com/en-us/windows/terminal/customize-settings/profile-appearance#background-images-and-icons)
for more information
