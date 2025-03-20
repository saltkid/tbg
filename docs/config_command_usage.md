# Table of Contents
- [Overview](#tbg-config-[arg])
- [Walkthroughs](#walkthroughs)
- [Usage](#usage)
    - [Printing config](#printing-config)
    - [Editing fields of config](#editing-fields-of-config)

---

# `tbg config [--flags]`
#### args: no arg

`config` command prints the `config.yaml` if no arg is specified.
If any of the flags are specified, it will edit the fields of `.tbg.yml`
based on flags passed to it. 

# Valid Flags
1. `--profile [arg]`
    - args: `default`, `1`, `2`, ..., `<n>`, `"profile name"`
    - edits: `profile`
2. `--interval [arg]`
    - args: `topRight`, `top`, `topLeft`, `left`, `center`, `right`, `bottomLeft`, `bottom`, `bottomRight`
    - edits: `interval`

# Usage
#### Printing config
To print the `.tbg.yml`, just do
```bash
tbg config
```
Output on console should look something like this
```bash
------------------------------------------------------------------------------------
| abs/path/to/currently/used/config.yaml
------------------------------------------------------------------------------------
| paths:
| - path: path/to/images/dir
|   alignment: right
|   stretch: fill
|   opacity: 0.250000
|
| profile:               default
| interval:              30
------------------------------------------------------------------------------------
```

#### Editing fields of config
To edit fields of config, specify any fields you want to edit with flags like this this:
```bash
tbg config --profile pwsh
```
```bash
------------------------------------------------------------------------------------
| abs/path/to/currently/used/config.yaml
------------------------------------------------------------------------------------
| paths:
| - path: path/to/images/dir
|   alignment: right
|   stretch: fill
|   opacity: 0.250000
|
| profile:               pwsh # used to be default
| interval:              30
------------------------------------------------------------------------------------
```
You can do this with interval as well
```bash
tbg config --interval 5
```
```bash
------------------------------------------------------------------------------------
| abs/path/to/currently/used/config.yaml
------------------------------------------------------------------------------------
| paths:
| - path: path/to/images/dir
|   alignment: right
|   stretch: fill
|   opacity: 0.250000
|
| profile:               1        
| interval:              5 # used to be 30
------------------------------------------------------------------------------------
```
