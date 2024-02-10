# Table of Contents
- [Overview](#tbg-remove-[arg])
- [Valid Flags](#valid-flags)
- [Walkthroughs](#walkthroughs)
    - [Removing a path](#removing-a-path)
    - [Removing flags from a path](#removing-flags-from-a-path)
    - [Removing all flags from a path](#removing-all-flags-from-a-path)
---

# `tbg remove [arg]`
#### args: `path/to/images/dir`
`remove` command removes a path to **tbg**'s currently used config.

You can remove flags from a path using `--`flags (see [valid flags](#valid-flags))

# Valid Flags
1. `--alignment [arg]`
    - args: `topRight`, `top`, `topLeft`, `left`, `center`, `right`, `bottomLeft`, `bottom`, `bottomRight`
    - it will remove the alignment flag of the path specified and replace it with the `default_alignment` value
2. `--stretch [arg]`
    - args: `none`, `fill`, `uniform`, `uniformToFill`
    - it will remove the stretch flag of the path specified and replace it with the `default_stretch` value
3. `--opacity [arg]`
    - args: any float between 0 and 1 (inclusive)
    - it will remove the opacity flag of the path specified and replace it with the `default_opacity` value

# Walkthroughs
### Removing a path
Let's say this is the currently used config:
```
image_col_paths:
- path/to/images/dir1

profile: default
interval: 30

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```
If we run:
```
tbg remove path/to/images/dir1
```
It will remove `path/to/images/dir1` from the currently used config's `image_col_paths` field like this:
```
image_col_paths: []

profile: default
interval: 30

default_alignment: center
default_stretch: fill
default_opacity: 0.1
```
### Removing flags from a path
Let's use this config
```
image_col_paths:
- path/to/images/dir1 | left fill 0.5

profile: default
interval: 30

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```
Let's "remove" the alignment flag of `path/to/images/dir1`:
```
tbg remove path/to/images/dir1 --alignment
```
```
image_col_paths:
- path/to/images/dir1 | _ fill 0.5

profile: default
interval: 30

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
```
When you remove a single flag, it will be blanked out. This just means the blanked out alignment flag will inherit the `default_alignment` value (`center`).

Let's remove stretch next.
```
tbg remove path/to/images/dir1 --stretch
```
```
image_col_paths:
- path/to/images/dir1 | _ _ 0.5

profile: default
interval: 30

default_alignment: center
default_stretch: uniform
default_opacity: 0.1
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
