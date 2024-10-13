# Table of Contents
- [Overview](#tbg-add-[arg])
- [Valid Flags](#valid-flags)
- [Usage](#Usage)
    - [Adding a path](#adding-a-path)
    - [Adding a path with flags](#adding-a-path-with-flags)
    - [Adding a flag to an existing path](#adding-a-flag-to-an-existing-path)
    - [Changing a flag of an existing path](#changing-a-flag-of-an-existing-path)
---

# `tbg add [arg] [--flags]`
#### args: `path/to/images/dir`
`add` command adds a path to **tbg**'s currently used config.
You can add flags to a path to be added using `--`flags

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

# Usage
### Adding a path
Let's say this is the currently used config:
```
paths: []

alignment: center
stretch: uniform
opacity: 0.5

other fields...
```
If we run:
```
tbg add path/to/images/dir1
```
It will add `path/to/images/dir1` to the currently used config's `paths` field like this:
```
paths:
- path: path/to/images/dir1

alignment: center
stretch: uniform
opacity: 0.5

other fields...
```
### Adding a path with flags
Let's continue with our config and run this:
```
tbg add path/to/images/dir2 --alignment left --stretch fill --opacity 0.5
```
It will add `path/to/images/dir2` in `paths` field and add flags to it like this:
```
paths:
- path: path/to/images/dir1
- path: path/to/images/dir2 
  alignment: left
  stretch: fill
  opacity: 0.5

alignment: center
stretch: uniform
opacity: 0.5

other fields...
```
Let's add another one:
```
tbg add path/to/images/dir3 --alignment right
```
```
paths:
- path: path/to/images/dir1
- path: path/to/images/dir2 
  alignment: left
  stretch: fill
  opacity: 0.5
- path: path/to/images/dir3
  alignment: right

alignment: center
stretch: uniform
opacity: 0.5
```
Flags that were not specified will inherit their respective default flag field value
(`stretch` and `opacity` in this example)

#### Adding a flag to an existing path
Let's continue with our config and run this:
```
tbg add /path/to/images/dir3 --stretch fill
```
This will find if `path/to/images/dir3` is already in `paths` field and add a stretch flag `fill` to it.
```
paths:
- other paths ...
- path: path/to/images/dir3
  alignment: right
  stretch: fill # this path only had alignment in the example above

alignment: center
stretch: uniform
opacity: 0.5
```
Let's fill assign the opacity too
```
tbg add /path/to/images/dir3 --opacity 0.25
```
```
paths:
- other paths ...
- path: path/to/images/dir3
  alignment: right
  stretch: fill
  opacity: 0.25

alignment: center
stretch: uniform
opacity: 0.5
```
#### Changing a flag of an existing path
Let's change the opacity and stretch:
```
tbg add /path/to/images/dir3 --opacity 1 --stretch none
```
This will find if `path/to/images/dir3` is already in `paths` field and set
the opacity to `1` and the stretch to `none`.
```
paths:
- other paths ...
- path: path/to/images/dir3
  alignment: right
  stretch: none # used to be fill
  opacity: 1    # used to be 0.25

alignment: center
stretch: uniform
opacity: 0.5
```
Let's change the alignment too
```
tbg add /path/to/images/dir3 --alignment bottom
```
```
paths:
- other paths ...
- path: path/to/images/dir3
  alignment: bottom # used to be right
  stretch: none 
  opacity: 1

alignment: center
stretch: uniform
opacity: 0.5
```
