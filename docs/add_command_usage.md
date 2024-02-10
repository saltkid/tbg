# Table of Contents
- [Overview](#tbg-add-[arg])
- [Valid Flags](#valid-flags)
- [Walkthroughs](#walkthroughs)
    - [Adding a path](#adding-a-path)
    - [Adding a path with flags](#adding-a-path-with-flags)
    - [Adding a flag to an existing path](#adding-a-flag-to-an-existing-path)
    - [Changing a flag of an existing path](#changing-a-flag-of-an-existing-path)
---

# `tbg add [arg]`
#### args: `path/to/images/dir`
`add` command adds a path to **tbg**'s currently used config.

You can add flags to a path to be added using `--`flags (see [valid flags](#valid-flags))

# Valid Flags
1. `--alignment [arg]`
    - args: `topRight`, `top`, `topLeft`, `left`, `center`, `right`, `bottomLeft`, `bottom`, `bottomRight`
    - it will add flags to the path being added after a pipe `|`
        - example: `path/to/images/dir | center fill 0.5`
2. `--stretch [arg]`
    - args: `none`, `fill`, `uniform`, `uniformToFill`
    - it will add flags to the path being added after a pipe `|`
3. `--opacity [arg]`
    - args: any float between 0 and 1 (inclusive)
    - it will add flags to the path being added after a pipe `|`

# Walkthroughs
### Adding a path
Let's say this is the currently used config:
```
image_col_paths: []

other fields...
```
If we run:
```
tbg add path/to/images/dir1
```
It will add `path/to/images/dir1` to the currently used config's `image_col_paths` field like this:
```
image_col_paths:
- path/to/images/dir1

other fields...
```
### Adding a path with flags
Let's continue with our config and run this:
```
tbg add path/to/images/dir2 --alignment left --stretch uniform --opacity 0.5
```
It will add `path/to/images/dir2` in `image_col_paths` field and add flags to it like this:
```
image_col_paths:
- path/to/images/dir1
- path/to/images/dir2 | left uniform 0.5

other fields...
```
Let's add another one:
```
tbg add path/to/images/dir3 --alignment center
```
```
image_col_paths:
- path/to/images/dir1
- path/to/images/dir2 | left uniform 0.5
- path/to/images/dir3 | center _ _

other fields...
```
Since only `--alignment` was specified, the other two flags are blanked out. This just means the blanked out flags will inherit their respective default flag field value (`default_stretch` and `default_opacity` in this example)

#### Adding a flag to an existing path
Let's continue with our config and run this:
```
tbg add /path/to/images/dir3 --stretch fill
```
This will find if `path/to/images/dir3` is already in `image_col_paths` field and add a stretch flag `fill` to it.
```
image_col_paths:
- path/to/images/dir1
- path/to/images/dir2 | left uniform 0.5
- path/to/images/dir3 | center fill _

other fields...
```
Let's fill assign the opacity too
```
tbg add /path/to/images/dir3 --opacity 0.5
```
```
image_col_paths:
- path/to/images/dir1
- path/to/images/dir2 | left uniform 0.5
- path/to/images/dir3 | center fill 0.5

other fields...
```
#### Changing a flag of an existing path
Let's continue with our config:
```
image_col_paths:
- path/to/images/dir1
- path/to/images/dir2 | left uniform 0.5
- path/to/images/dir3 | center fill 0.5

other fields...
```
Let's change the opacity and stretch:
```
tbg add /path/to/images/dir3 --opacity 1 --stretch none
```
This will find if `path/to/images/dir3` is already in `image_col_paths` field and set the opacity to `1` and the stretch to `none`.
```
image_col_paths:
- path/to/images/dir1
- path/to/images/dir2 | left uniform 0.5
- path/to/images/dir3 | center none 1

other fieds...
```
Let's change the alignment too
```
tbg add /path/to/images/dir3 --alignment bottom
```
```
image_col_paths:
- path/to/images/dir1
- path/to/images/dir2 | left uniform 0.5
- path/to/images/dir3 | bottom center none 1

other fields...
```
