# Table of Contents
- [Overview](#tbg-run)
- [Executing with flags](#executing-with-flags)
- [Usage](#usage)
    - [Normal Execution, and path specific options](#normal-execution)
    - [Using a custom config](#using-a-custom-config)
    - [Overriding `profile`, `port`, and `interval` fields](#overriding-profile-port-and-interval-fields)
    - [Overriding per-path options](#overriding-per-path-options)
---

# `tbg run`

`run` command edits the `settings.json` used by *Windows Terminal* using
settings from **tbg**'s config. **tbg** will keep running, editing the
`settings.json` of *Windows Terminal*, replacing the background image. The 
background image is chosen randomly from the images under each `path` specified
in the config.

On initial execution of **tbg**, it will create a config at
`$env:LOCALAPPDATA/tbg/config.yml` if it does not exist already.
For more information, see documentation on [config.yml](/docs/config.yml.md).

# Executing with flags
### Valid Flags: `--profile`, `--interval`, `--port`, `--alignment`, `--opacity`, `--stretch`

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
This section will delve on path specific options.

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
The server will use the port `9545`

You can quit **tbg** by running:
```bash
tbg quit
```

---
### Using a custom config
To use a custom config instead of the default one, use the `-c, --config` flag:
```bash
tbg run --config /path/to/custom/config.yml
```
This will use the values in the custom config instead of the default one when
editing. This is useful if you want, for example, to have a separate config for
your pwsh, command prompt and various wsl distros, and you do not want to have
to specify all the flags/have different paths to choose images from.

---
### Overriding `profile`, `port`, and `interval` fields

Instead of `tbg run`, let's do:
```bash
tbg run --profile 1 --interval 5 --port 8000
```
Let's say this is the default config:
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

To quit, make sure to put the same port specified in `tbg run`:
```bash
tbg quit --port 8000
```

---
### Overriding per-path options
This will delve on overriding default alignment, opacity, and stretch values.
This will also override the per-path options.

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
in the config. This means instead of `path/to/dir1`'s images having the default
alignment of `center`, default stretch of `uniformToFill`, and default opacity
of `1.0`, the images instead use the values specified by the flags (`right`,
`none`, `0.35`)

Notice that `path/to/dir2` has options that should override the default options
fields. However, since we specified `--alignment right --opacity 0.35 --stretch
none`, tbg will use these value instead, just like with `path/to/dir1`.
