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
    - args: `topRight`, `top`, `topLeft`, `left`, `center`, `right`, `bottomLeft`, `bottom`, `bottomRight`
    - it will remove the alignment option of the path specified
2. `--stretch [arg]`
    - args: `none`, `fill`, `uniform`, `uniformToFill`
    - it will remove the stretch option of the path specified
3. `--opacity [arg]`
    - args: any float between 0 and 1 (inclusive)
    - it will remove the opacity option of the path specified

# Usage
### Removing a path
Let's say this is the config:
```yml
paths:
- path: path/to/images/dir1

interval: 30
profile: default
port: 9545
```
If we run:
```bash
tbg remove path/to/images/dir1
```
It will remove `path/to/images/dir1` from the config's `paths`
field like this:
```yml
paths: []

interval: 30
profile: default
port: 9545
```
### Removing options from a path
Let's use this config
```yml
paths:
- path: path/to/images/dir1
  alignment: left
  stretch: fill
  opacity: 0.5

interval: 30
profile: default
port: 9545
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

interval: 30
profile: default
port: 9545
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

interval: 30
profile: default
port: 9545
```
The path will now inherit the default `stretch` value (`uniformToFill`).

Now let's see if we remove opacity:
```bash
tbg remove path/to/images/dir1 --opacity
```
```yml
paths:
- path: path/to/images/dir1

interval: 30
profile: default
port: 9545
```
Now all the path will inehrit the default alignment (`center`), stretch
(`uniformToFill`), and opacity (`1.0`).
