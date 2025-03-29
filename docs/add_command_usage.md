# Table of Contents
- [Overview](#tbg-add-[arg])
- [Valid Flags](#valid-flags)
- [Usage](#Usage)
    - [Adding a path](#adding-a-path)
    - [Adding a path with options](#adding-a-path-with-options)
    - [Adding an option to an existing path](#adding-an-option-to-an-existing-path)
    - [Changing an option of an existing path](#changing-an-option-of-an-existing-path)
---

# `tbg add [arg] [--flags]`
#### args: `path/to/images/dir`
`add` command adds a path to `.tbg.yml`.
You can add options to a path to be added using flags

# Valid Flags
1. `-a, --alignment [arg]`
    - args: `topRight`, `top`, `topLeft`, `left`, `center`, `right`, `bottomLeft`, `bottom`, `bottomRight`
    - it will add alignment option to the path being added
2. `-s, --stretch [arg]`
    - args: `none`, `fill`, `uniform`, `uniformToFill`
    - it will add stretch option to the path being added
3. `-o, --opacity [arg]`
    - args: any float between 0 and 1 (inclusive)
    - it will add opacity option to the path being added

# Usage
### Adding a path
Let's say this is `.tbg.yml`:
```yml
paths: []

interval: 30
profile: default
port: 9545
```
If we run:
```bash
tbg add path/to/images/dir1
```
It will add `path/to/images/dir1` to `.tbg.yml`'s `paths` field
like this:
```yml
paths:
- path: path/to/images/dir1

interval: 30
profile: default
port: 9545
```
### Adding a path with options
Let's continue with our config and run this:
```bash
tbg add path/to/images/dir2 --alignment left --stretch fill --opacity 0.5
```
It will add `path/to/images/dir2` in `paths` field and add options to it like this:
```yml
paths:
- path: path/to/images/dir1
- path: path/to/images/dir2 
  alignment: left
  stretch: fill
  opacity: 0.5

interval: 30
profile: default
port: 9545
```
Let's add another one:
```bash
tbg add path/to/images/dir3 --alignment right
```
```yml
paths:
- path: path/to/images/dir1
- path: path/to/images/dir2 
  alignment: left
  stretch: fill
  opacity: 0.5
- path: path/to/images/dir3
  alignment: right
```
Options that were not specified will inherit their respective default value
(`stretch` and `opacity` will be `center` and `1.0` in this example)

#### Adding an option to an existing path
Let's continue with our config and run this:
```bash
tbg add /path/to/images/dir3 --stretch fill
```
This will find if `path/to/images/dir3` is already in `paths` field and add a
stretch option `fill` to it.
```yml
paths:
- other paths ...
- path: path/to/images/dir3
  alignment: right
  stretch: fill # this path only had alignment in the example above

interval: 30
profile: default
port: 9545
```
Let's fill assign the opacity too
```bash
tbg add /path/to/images/dir3 --opacity 0.25
```
```yml
paths:
- other paths ...
- path: path/to/images/dir3
  alignment: right
  stretch: fill
  opacity: 0.25

interval: 30
profile: default
port: 9545
```
#### Changing an option of an existing path
Let's change the opacity and stretch:
```bash
tbg add /path/to/images/dir3 --opacity 1 --stretch none
```
This will find if `path/to/images/dir3` is already in `paths` field and set
the opacity to `1` and the stretch to `none`.
```yml
paths:
- other paths ...
- path: path/to/images/dir3
  alignment: right
  stretch: none # used to be fill
  opacity: 1    # used to be 0.25

interval: 30
profile: default
port: 9545
```
Let's change the alignment too
```bash
tbg add /path/to/images/dir3 --alignment bottom
```
```yml
paths:
- other paths ...
- path: path/to/images/dir3
  alignment: bottom # used to be right
  stretch: none 
  opacity: 1

interval: 30
profile: default
port: 9545
```
