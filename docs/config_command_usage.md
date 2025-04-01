# Table of Contents
- [Overview](#tbg-config-[arg])
- [Walkthroughs](#walkthroughs)
- [Usage](#usage)
    - [Printing config](#printing-config)
    - [Editing fields of config](#editing-fields-of-config)
    - [Printing and editing fields of a custom config](#printing-and-editing-fields-of-a-custom-config)
---
# `tbg config [--flags]`
#### args: no arg, `path/to/config/file.yml`

`config` command prints the `config.yaml` if no flags are specified. Pass in a
path to a custom config file to use that instead. If any of the flags are
specified, it will edit the fields of the config based on flags passed to
it. 

# Valid Flags
1. `--interval [arg]`
    - edits: interval field
    - args: `topRight`, `top`, `topLeft`, `left`, `center`, `right`, `bottomLeft`, `bottom`, `bottomRight`
2. `--port [arg]`
    - edits: port field
    - args: any positive integer up to 65535 (unsigned 16 byte int)
3. `--profile [arg]`
    - edits: profile field
    - args: `default`, `1`, `2`, ..., `<n>`, `"profile name"`

# Usage
#### Printing config
To print the config, just do
```bash
tbg config
```
Output on console should look something like this
```bash
## /path/to/config.yml
paths:
    - path: path/to/images/dir
      alignment: right
      stretch: fill
      opacity: 0.25
profile:  default
port:     9545
interval: 1800
```

---
#### Editing fields of config
To edit fields of config, specify any fields you want to edit with flags like this this:
```bash
tbg config --profile pwsh
```
output:
```bash
## /path/to/config.yml
paths:
    - path: path/to/images/dir
      alignment: right
      stretch: fill
      opacity: 0.25
profile:  pwsh
port:     9545
interval: 1800
## EDITED
# profile    default --> pwsh

```
You can do this with interval and port as well
```bash
tbg config --interval 300 --port 9000
```
output:
```bash
## /path/to/config.yml
paths:
    - path: path/to/images/dir
      alignment: right
      stretch: fill
      opacity: 0.25
profile:  1        
port:     9545
interval: 300
## EDITED
# interval   1800 --> 300
# port       9545 --> 9000
```

---
#### Printing and editing fields of a custom config
Just specify a path to a custom config to print/edit its fields:
```bash
tbg config /path/to/custom/config.yml
```
output:
```bash
## /path/to/custom/config.yml
paths:
    - path: path/to/images/dir1
      opacity: 0.250000
    - path: path/to/images/dir2
      alignment: right
profile:  Arch
port:     8000
interval: 60
```
Editing fields is just the same as before:
```bash
tbg config /path/to/custom/config.yml -p Debian -P 8090 -i 6000
```
output:
```bash
## /path/to/custom/config.yml
paths:
    - path: path/to/images/dir1
      opacity: 0.250000
    - path: path/to/images/dir2
      alignment: right
profile:  Arch
port:     8000
interval: 6000
## EDITED
# interval   60 --> 6000
# port       8000 --> 8090
# profile    Arch --> Debian
```
