# Table of Contents
- [Overview](#tbg-remove-[arg])
- [Valid Flags](#valid-flags)
- [Usage](#usage)
    - [Removing a path](#removing-a-path)
    - [Removing options from a path](#removing-options-from-a-path)
---

# `tbg remove [arg] [--flags]`
#### args: `path/to/images/dir`
`remove` command removes a path to **tbg**'s config.
You can remove flags from a path using `--`flags

# Valid Flags
1. `--alignment [arg]`
    - args: `topRight`, `top`, `topLeft`, `left`, `center`, `right`,
    `bottomLeft`, `bottom`, `bottomRight`
    - it will remove the alignment option of the path specified
2. `-c, --config [arg]`
    - args: `/path/to/custom/config.yml`
    - the path will be removed from the custom config instead of the default
    one
3. `--opacity [arg]`
    - args: any float between 0 and 1 (inclusive)
    - it will remove the opacity option of the path specified
4. `--stretch [arg]`
    - args: `none`, `fill`, `uniform`, `uniformToFill`
    - it will remove the stretch option of the path specified

# Usage
### Removing a path
Let's say this is the config:
```yml
paths:
- path: path/to/images/dir1

# other fields...
```
If we run:
```bash
tbg remove path/to/images/dir1
```
It will remove `path/to/images/dir1` from the config's `paths`
field like this:
```yml
paths: []

# other fields...
```

To remove a path to a another config, use `-c, --config`. Let's use this custom
config:
```yml
paths:
- path: /path/to/images/dir3
  opacity: 0.1
- path: /path/to/images/dir2
  stretch: none
- path: /path/to/images/dir1

# other fields...
```
To remove `/path/to/images/dir1` from this custom config, do:
```bash
tbg remove /path/to/images/dir1 --config /path/to/custom/config.yml
```
Resulting file:
```yml
paths:
- path: /path/to/images/dir1
  opacity: 0.1
- path: /path/to/images/dir1
  stretch: none

# other fields...
```
---
### Removing options from a path
Let's use this config
```yml
paths:
- path: path/to/images/dir1
  alignment: left
  stretch: fill
  opacity: 0.5

# other fields...
```
Let's "remove" the alignment flag of `path/to/images/dir1`:
```bash
tbg remove path/to/images/dir1 --alignment
```
```yml
paths:
- path: path/to/images/dir2 
  stretch: fill
  opacity: 0.5

# other fields...
```
The path will now inherit the default `alignment` value (`center`).
Let's remove stretch next.
```bash
tbg remove path/to/images/dir1 --stretch
```
```yml
paths:
- path: path/to/images/dir2 
  stretch: fill
  opacity: 0.5

# other fields...
```
The path will now inherit the default `stretch` value (`uniformToFill`).

Now let's see if we remove opacity:
```bash
tbg remove path/to/images/dir1 --opacity
```
```yml
paths:
- path: path/to/images/dir1

# other fields...
```
Now all the path will inehrit the default alignment (`center`), stretch
(`uniformToFill`), and opacity (`1.0`).
