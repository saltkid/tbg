# Table of Contents
- [Overview](#tbg-config-[arg])
- [Walkthroughs](#walkthroughs)
    - [Printing config](#printing-config)
    - [Editing fields of config](#editing-fields-of-config)

---

# `tbg config [arg]`
#### args: `edit`, no arg

`config` command prints the `config.yaml` if no arg is specified. If `edit` is specified, it will edit the fields of `config.yaml` based on flags passed to it. 

# Walkthroughs
#### Printing config
To print the currently used config, just do
```
tbg config
```
Output on console should look something like this
```
------------------------------------------------------------------------------------
| abs/path/to/currently/used/config.yaml
------------------------------------------------------------------------------------
| image_col_paths:
|                        abs/path/to/images/dir1
|                        abs/path/to/images/dir2
|                        abs/path/to/images/dir3
|
| profile:               default
| interval:              30
|
| default_alignment:     center
| default_stretch:       uniform
| default_opacity:       0.1
------------------------------------------------------------------------------------
```

#### Editing fields of config
To edit fields of config, specify `edit` as the arg and then specify the fields you want to edit with flags like this this:
```
tbg config edit --alignment topRight
```
This means we are going to edit the `default_alignment` field in `config.yaml`. This is the before:
```
------------------------------------------------------------------------------------
| abs/path/to/default/config.yaml
------------------------------------------------------------------------------------
| image_col_paths: []
|
| profile:               default
| interval:              30
|
| default_alignment:     center
| default_stretch:       uniform
| default_opacity:       0.1
------------------------------------------------------------------------------------
```
After:
```
------------------------------------------------------------------------------------
| abs/path/to/default/config.yaml
------------------------------------------------------------------------------------
| image_col_paths: []
|
| profile:               default
| interval:              30
|
| default_alignment:     topRight
| default_stretch:       uniform
| default_opacity:       0.1
------------------------------------------------------------------------------------
```
You can do this with the other four fields as well (not `image_col_paths`)
```
tbg config edit --profile list-1 --interval 5 --stretch fill --opacity 0.35
```
```
------------------------------------------------------------------------------------
| abs/path/to/default/config.yaml
------------------------------------------------------------------------------------
| image_col_paths: []
|
| profile:               list-1
| interval:              5
|
| default_alignment:     topRight
| default_stretch:       fill
| default_opacity:       0.35
------------------------------------------------------------------------------------
```
