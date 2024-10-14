# Table of Contents
- [Overview](#tbg-remove-[arg])
- [Valid Flags](#valid-flags)
- [Usage](#usage)
    - [Removing a path](#removing-a-path)
    - [Removing options from a path](#removing-options-from-a-path)
---

# `tbg remove [arg] [--flags]`
#### args: `path/to/images/dir`
`remove` command removes a path to **tbg**'s `.tbg.yml`.
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
Let's say this is the `.tbg.yml`:
```yml
paths:
- path: path/to/images/dir1

alignment: center
stretch: uniform
opacity: 0.5

other fields...
```
If we run:
```bash
tbg remove path/to/images/dir1
```
It will remove `path/to/images/dir1` from the `.tbg.yml`'s `paths`
field like this:
```yml
paths: []

alignment: center
stretch: uniform
opacity: 0.5

other fields...
```
### Removing options from a path
Let's use this config
```yml
paths:
- path: path/to/images/dir1
  alignment: left
  stretch: fill
  opacity: 0.5

alignment: center
stretch: uniform
opacity: 0.5

other fields...
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

alignment: center
stretch: uniform
opacity: 0.5

other fields...
```
The path will now inherit the top level `alignment` value (`center`).
Let's remove stretch next.
```bash
tbg remove path/to/images/dir1 --stretch
```
```yml
paths:
- path: path/to/images/dir2 
  stretch: fill
  opacity: 0.5

alignment: center
stretch: uniform
opacity: 0.5

other fields...
```
The path will now inherit the top level `stretch` value (`uniform`).

Now let's see if we remove opacity:
```bash
tbg remove path/to/images/dir1 --opacity
```
```yml
paths:
- path: path/to/images/dir1

profile: default
interval: 30

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```
Now all the path will inehrit the top level alignment, stretch, and opacity.
