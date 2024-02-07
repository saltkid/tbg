# Table of Contents
- [Overview](#tbg-config-[arg])
- [Walkthroughs](#walkthroughs)
    - [Printing currently used config](#printing-currently-used-config)
    - [Setting default config as the currently used config](#setting-default-config-as-the-currently-used-config)
    - [Setting another config as the currently used config](#setting-a-config-as-the-currently-used-config)

---

# `tbg config [arg]`
#### args: `path/to/config.yaml`, `default`, no arg

`config` command prints the currently used config if no arg is specified. If an arg is specified, it will set the arg as the currently used config after validation, then print it.

The way **tbg** keeps track of the currently used config is by using `tbg_profile.yaml` that was auto generated on initial execution. It is not recommended to edit this. See [`tbg_profile.yaml`](https://github.com/saltkid/tbg/blob/main/docs/tbg_profile.yaml.md) for more information.

# Walkthroughs
#### Printing currently used config
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
| default_alignment:     center
| default_stretch:       uniform
| default_opacity:       0.1
------------------------------------------------------------------------------------
```

#### Setting default config as the currently used config
To set the default config as the currently used config, just do
```
tbg config default
```
The default config is the `config.yaml` in the same path as the **tbg** executable. It was auto generated on initial execution of **tbg**.

This will edit the `tbg_profile.yaml`'s `used_config` field to the path of the default config. Output on console should look similar to the previous example
```
------------------------------------------------------------------------------------
| abs/path/to/default/config.yaml
------------------------------------------------------------------------------------
| image_col_paths:
|                        abs/path/to/images/dir1
|                        abs/path/to/images/dir2
|                        abs/path/to/images/dir3
|
| profile:               default
| interval:              30
| default_alignment:     center
| default_stretch:       uniform
| default_opacity:       0.1
------------------------------------------------------------------------------------
```

#### Setting another config as the currently used config
This is pretty much the same as the previous example, except instead of `default`, you are passing in the abs path of another `config.yaml` that you want to set as the currently used config
```
tbg config path/to/another/config.yaml
```
Again, it edits the `tbg_profile.yaml`'s `used_config` field to the path of the specified config. Output looks similar too
```
------------------------------------------------------------------------------------
| abs/path/to/another/config.yaml
------------------------------------------------------------------------------------
| image_col_paths:
|                        abs/path/to/images/dir1
|                        abs/path/to/images/dir2
|
| profile:               default
| interval:              30
| default_alignment:     center
| default_stretch:       uniform
| default_opacity:       0.1
------------------------------------------------------------------------------------
```
