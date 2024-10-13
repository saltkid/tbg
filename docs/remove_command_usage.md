# Table of Contents
- [Overview](#tbg-remove-[arg])
- [Valid Flags](#valid-flags)
- [Usage](#usage)
    - [Removing a path](#removing-a-path)
    - [Removing flags from a path](#removing-flags-from-a-path)
    - [Removing all flags from a path](#removing-all-flags-from-a-path)
---

# `tbg remove [arg] [--flags]`
#### args: `path/to/images/dir`
`remove` command removes a path to **tbg**'s currently used config.
You can remove flags from a path using `--`flags

# Valid Flags
1. `--alignment [arg]`
    - args: `topRight`, `top`, `topLeft`, `left`, `center`, `right`, `bottomLeft`, `bottom`, `bottomRight`
    - it will remove the alignment flag of the path specified and replace it
    with the `alignment` value
2. `--stretch [arg]`
    - args: `none`, `fill`, `uniform`, `uniformToFill`
    - it will remove the stretch flag of the path specified and replace it
    with the `stretch` value
3. `--opacity [arg]`
    - args: any float between 0 and 1 (inclusive)
    - it will remove the opacity flag of the path specified and replace it
    with the `opacity` value

# Usage
### Removing a path
Let's say this is the currently used config:
```
paths:
- path: path/to/images/dir1

alignment: center
stretch: uniform
opacity: 0.5

other fields...
```
If we run:
```
tbg remove path/to/images/dir1
```
It will remove `path/to/images/dir1` from the currently used config's `image_col_paths` field like this:
```
paths: []

alignment: center
stretch: uniform
opacity: 0.5

other fields...
```
### Removing flags from a path
Let's use this config
```
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
```
tbg remove path/to/images/dir1 --alignment
```
```
paths:
- path: path/to/images/dir2 
  stretch: fill
  opacity: 0.5

other fields...
```
The path will now inherit the top level `alignment` value (`center`).

Let's remove stretch next.
```
tbg remove path/to/images/dir1 --stretch
```
```
paths:
- path: path/to/images/dir2 
  stretch: fill
  opacity: 0.5

other fields...
```
The stretch flag of `path/to/images/dir1` is blanked out as well. Again, this will inherit the `default_stretch` value: `uniform`.

Now let's see if we remove opacity:
```
tbg remove path/to/images/dir1 --opacity
```
```
image_col_paths:
- path/to/images/dir1

profile: default
interval: 30

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```
It removed the flags of `path/to/images/dir1`. This is because if we removed opacity, all three would've been blank. In this case, **tbg** will just remove the flags of `path/to/images/dir`.

But that seems a lot of prompts just to get rid of all the flags of one path. See next walkthrough for a *"shortcut"*

### Removing all flags from a path
Let's use this config:
```
image_col_paths:
- path/to/images/dir1 | left fill 0.5

profile: default
interval: 30

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```
Instead of removing the flags of `path/to/images/dir1` one by one, we can specify all flags in the prompt
```
tbg remove path/to/images/dir1 -s -a -o
```
Now all 3 flags are specified. This means it will remove all flags of `path/to/images/dir1` like this:
```
image_col_paths:
- path/to/images/dir1

profile: default
interval: 30

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```
That's it. That's the *"shortcut"*
