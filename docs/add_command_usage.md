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
`add` command adds a path to **tbg**'s config.
You can add options to a path to be added using flags

# Valid Flags
1. `-a, --alignment [arg]`
    - args: `topRight`, `top`, `topLeft`, `left`, `center`, `right`,
    `bottomLeft`, `bottom`, `bottomRight`
    - it will add alignment option to the path being added
2. `-c, --config [arg]`
    - args: `/path/to/custom/config.yml`
    - the path will be added to the custom config instead of the default one
3. `-o, --opacity [arg]`
    - args: any float between 0 and 1 (inclusive)
    - it will add opacity option to the path being added
4. `-s, --stretch [arg]`
    - args: `none`, `fill`, `uniform`, `uniformToFill`
    - it will add stretch option to the path being added

# Usage
### Adding a path
Let's say this is the default config:
```yml
paths: []

# other fields...
```
If we run:
```bash
tbg add /path/to/images/dir1
```
It will add `path/to/images/dir1` to the config's `paths` field
like this:
```yml
paths:
- path: /path/to/images/dir1 # you added this

# other fields...
```
To add the path to a another config, use `-c, --config`:
```bash
tbg add /path/to/images/dir1 --config /path/to/custom/config.yml
```
It will add `/path/to/images/dir1` to `/path/to/custom/config.yml`:
```yml
paths:
- path: /path/to/images/dir3
  opacity: 0.1
- path: /path/to/images/dir2
  stretch: none
- path: /path/to/images/dir1 # you added this

# other fields...
```

---
### Adding a path with options
Let's continue with the default config and run this:
```bash
tbg add /path/to/images/dir2 --alignment left --stretch fill --opacity 0.5
```
It will add `path/to/images/dir2` in `paths` field and add options to it like this:
```yml
paths:
- path: /path/to/images/dir1
- path: /path/to/images/dir2 
  alignment: left
  stretch: fill
  opacity: 0.5

# other fields...
```
Let's add another one:
```bash
tbg add /path/to/images/dir3 --alignment right
```
```yml
paths:
- path: /path/to/images/dir1
- path: /path/to/images/dir2 
  alignment: left
  stretch: fill
  opacity: 0.5
- path: /path/to/images/dir3
  alignment: right

# other fields...
```
Options that were not specified will inherit their respective default value.

---
#### Adding an option to an existing path
Let's continue with our config and run this:
```bash
tbg add /path/to/images/dir3 --stretch fill
```
This will find if `path/to/images/dir3` is already in `paths` field and add a
stretch option `fill` to it.
```yml
paths:
# other paths ...
- path: /path/to/images/dir3
  alignment: right
  stretch: fill # this path only had alignment in the example above

# other fields...
```
Let's fill assign the opacity too
```bash
tbg add /path/to/images/dir3 --opacity 0.25
```
```yml
paths:
- other paths ...
- path: /path/to/images/dir3
  alignment: right
  stretch: fill
  opacity: 0.25

# other fields...
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
- path: /path/to/images/dir3
  alignment: right
  stretch: none # used to be fill
  opacity: 1    # used to be 0.25

# other fields...
```
Let's change the alignment too
```bash
tbg add /path/to/images/dir3 --alignment bottom
```
```yml
paths:
- other paths ...
- path: /path/to/images/dir3
  alignment: bottom # used to be right
  stretch: none 
  opacity: 1

# other fields...
```
