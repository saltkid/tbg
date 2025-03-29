# Table of Contents
- [Overview](#tbg-run)
- [Key events](#key-events)
- [Executing with flags](#executing-with-flags)
- [Usage](#usage)
    - [Normal Execution, key events, and path specific options](#normal-execution)
    - [Overriding `profile`, `port`, and `interval` fields](#overriding-profile-port-and-interval-fields)
    - [Overriding default values](#overriding-default-values)
---

# `tbg run`

`run` command edits the `settings.json` used by *Windows Terminal* using
settings from `.tbg.yml`. **tbg** will keep running, editing the
`settings.json` of *Windows Terminal*, replacing the background image. The 
background image is chosen randomly from the images under each path specified
in `.tbg.yml`. Pressing `q` or `ctrl+c` will stop execution.

On initial execution of **tbg**, it will create a `.tbg.yml` in the same
directory as the **tbg** executable if it does not exist already. **There can
only be one `.tbg.yml`**. For more information, see documentation on
[tbg.yml](https://github.com/saltkid/tbg/blob/main/docs/tbg.yml.md).

# Key events
**tbg** takes optional commands during execution:
- `q`: quits
- `c`: shows the available commands
- `n`: goes to next image

# Executing with flags
#### Valid Flags: `--profile`, `--interval`, `--port`, `--alignment`, `--opacity`, `--stretch`

The flags specified will override any per path options specified. So if there is
a `path/to/dir` with the alignment `center`, **tbg** will use whatever value
the `--alignment` flags has instead of that

The order of importance is:
1. flags (`--alignment`, `--opacity`, `--stretch`)
2. per path options in config
3. default values

For an example, see [overriding default flags walkthrough](#overriding-default-option-fields)

# Usage
### Normal Execution
This will delve on and key events and path specific options.

Let's do:
```bash
tbg run
```
Let's say that this is the config:
```yml
paths:
- path: /path/to/dir1
- path: /path/to/dir2
  alignment: right
  stretch: fill
  opacity: 0.35

profile: default
port: 9545
interval: 30
```
This just means that when we do `tbg run`, we want to change the background
image of the **default** *Windows Terminal* profile every **30 minutes**. The
image is chosen randomly from images under `/path/to/dir1` and `/path/to/dir2`.

Now let's quit **tbg** by pressing `q` or `ctrl+c`.

---
### Overriding `profile`, `port`, and `interval` fields

Instead of `tbg run`, let's do:
```bash
tbg run --profile 1 --interval 5 --port 8000
```
```yml
paths:
- path: path/to/dir1

profile: default
port: 9545
interval: 30
```

The `--profile`, `--port`, and `--interval` flags will override the values in
the config. Again, not edit them. This means instead of changing the background
image of the `default` profile every 30 minutes, it will change the background
image of the first profile under `list` field in `settings.json` every 5
minutes instead. The server will run in port 8000 instead of what is defined
in the config (9545)

---
### Overriding default values
This will delve on overriding default values on the config using flags. This
will also override the per-path options.

Let's use this config:
```yml
paths:
- path: path/to/dir1
- path: path/to/dir2 
  alignment: left
  stretch: uniformToFill
  opacity: 0.25

profile: default
port: 9545
interval: 30
```
Let's do:
```bash
tbg run --alignment right --opacity 0.35 --stretch none
```

The `--alignment`, `--opacity`, and `--stretch` flags will override the values
in `.tbg.yml`. This means instead of `path/to/dir1`'s images having the default
alignment `center`, default stretch `uniformToFill`, and default opacity `1.0`, the
images instead use the values specified by the flags (`right`, `none`, `0.35`)

Notice that `path/to/dir2` has options that should override the default options
fields. However, since we specified `--alignment right --opacity 0.35 --stretch
none`, tbg will use these value instead, just like with `path/to/dir1`.
