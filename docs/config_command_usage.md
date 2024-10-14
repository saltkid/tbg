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
    - args: `default`, `list-1`, `list-2`, ..., `list-<n>`
    - edits: `profile`
2. `--interval [arg]`
    - args: `topRight`, `top`, `topLeft`, `left`, `center`, `right`, `bottomLeft`, `bottom`, `bottomRight`
    - edits: `interval`
3. `--alignment [arg]`
    - args: `topRight`, `top`, `topLeft`, `left`, `center`, `right`, `bottomLeft`, `bottom`, `bottomRight`
    - edits: `alignment` (the top level alignment, not the per path alignment)
4. `--stretch [arg]`
    - args: `none`, `fill`, `uniform`, `uniformToFill`
    - edits: `stretch` (top level stretch)
5. `--opacity [arg]`
    - args: any float between 0 and 1 (inclusive)
    - edits: `opacity` (top level opacity)

# Usage
#### Printing config
To print the currently used config, just do
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
| alignment:             center
| stretch:               uniform
| opacity:               0.100000
------------------------------------------------------------------------------------
```

#### Editing fields of config
To edit fields of config, specify any fields you want to edit with flags like this this:
```bash
tbg config --alignment topRight
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
| profile:               default
| interval:              30
| alignment:             topRight # used to be center in the above example
| stretch:               uniform
| opacity:               0.100000
------------------------------------------------------------------------------------
```
You can do this with the other four fields as well
```bash
tbg config --profile list-1 --interval 5 --stretch fill --opacity 0.35
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
| profile:               list-1   # used to be default
| interval:              5        # used to be 30
| alignment:             topRight
| stretch:               fill     # used to be uniform
| opacity:               0.350000 # used to be 0.100000
------------------------------------------------------------------------------------
```
