# Table of Contents
- [tbg](#tbg-Terminal-Background-Gallery)
- [Installation](#installation)
    - [Building from source](#building-from-source)
- [Demo](#demo)
    - [Auto change image every 3 seconds](#auto-change-image-every-3-seconds)
    - [Manually changing image through a keybind](#manually-changing-image-through-a-keybind)
    - [Manually changing image across multiple pwsh instances](#manually-changing-image-across-multiple-pwsh-instances)
    - [Separately change image for pwsh and wsl](#separately-change-image-for-pwsh-and-wsl)
- [Usage](#usage)
    - [tbg server](#tbg-server)
        - [Logging](#logging)
- [Config](#config)
    - [Fields](#fields)
- [Commands](#commands)
    - [Server Commands](#server-commands)
- [Example Setup](#example-setup)
    - [pwsh](#pwsh)
    - [zsh (on wsl)](#zsh-on-wsl)
- [Credits](#credits)

---
# tbg (Terminal BackGround)

**tbg** (*pronounced /ˈtiːbæɡ/*) is a tool for Windows Terminal that allows you
to cycle through multiple background images at a set interval. The image
sources, interval of changing, image alignment, etc. can be configured
via a [config file](#config).

This edits the `settings.json` used by *Windows Terminal*; specifically, the
`backgroundImage` on the default profile by default but user can specify which
profile to target. 

---
# Installation
Download the latest release of [tbg](https://github.com/saltkid/tbg/releases)

or

## Building from source
- clone the repo
```
git clone https://github.com/saltkid/tbg.git
cd tbg
```
- build it:
```
go mod tidy
go build -ldflags "-X main.TbgVersion=$(git describe --tags --dirty --always)-dev"

## or build it without debug information to make the executable smaller
# go build -ldflags "-s -w -X main.TbgVersion=$(git describe --tags --dirty --always)"
```
**Optionally** add the `tbg` executable to your PATH

*oneliner for copy pasting*
```bash
git clone https://github.com/saltkid/tbg.git && cd tbg && go mod tidy && go build -ldflags "-X main.TbgVersion=$(git describe --tags --dirty --always)-dev"
```

---
# Demo
### Auto change image every 3 seconds
https://github.com/user-attachments/assets/a6ab510c-4327-4461-a617-c677781348dd
### Manually changing image through a keybind
See [example setup](#example-setup) for a guide on how to setup keybinds.

https://github.com/user-attachments/assets/d285cecc-48d5-4217-9121-cce31a978b93
### Manually changing image across multiple pwsh instances
https://github.com/user-attachments/assets/fcbfa4e4-32a1-4ae5-a9fc-aafd27d80456
### Separately change image for pwsh and wsl
See [example setup](#example-setup) for a guide on how to setup pwsh and zsh
(wsl) to set up their own keybinds and tbg servers on shell startup.

https://github.com/user-attachments/assets/3d64cf1f-c954-458c-9308-87cdf4a27d61

---
# Usage
```
tbg run
```
This will start the **tbg** server. This will target Windows Terminal's default
profile, changing the background image every 30 minutes, choosing an random
image from the specified `paths` in **tbg**'s config file. See
[fields](#fields) for more information about images paths. The default image
properties are:

| image property | default value  |
|----------------|----------------|
| alignment      | `center`       |
| opacity        | `1.0`          |
| stretch        | `uniformToFill`|

On initial execution of `tbg run`, it will create a [default
`config.yml`](#config) in `$env:LOCALAPPDATA/tbg/config.yml`. You can edit this
manually **or** use [tbg commands](#commands) to edit the fields with input
validation.

## tbg server
`tbg run` starts the **tbg** http server where POST requests can be made to
trigger certain actions. These post requests can be made through [tbg server
commands](#server-commands) for convenience.

### Logging
A log file is in the same directory as the config file:
`$env:LOCALAPPDATA/tbg/tbg.log`. It's in json so I recommend using
[`jq`](https://jqlang.org/) to parse and get useful information out of it.

See [logs](/docs/logs.md) for all log types and their structure.

---
# [Config](/docs/config.yml.md)
To edit the `settings.json` *Windows Terminal* uses, **tbg** uses `config.yml`
located at `$env:LOCALAPPDATA/tbg/config.yml` by default . On initial execution
of `tbg config`, a default config is created at that path. See example config
at [samples](./samples/config.yml). A `schema.json` is also given for some
basic autocomplete with
[yaml-language-server](https://github.com/redhat-developer/yaml-language-server)

To use a different config in a command, use the `-c, --config` flag and pass in
a path to a config file.

## Fields
Although you can edit the fields in the config directly, it is recommended to
use the [`config`](/docs/config_command_usage.md)
command to edit the file.
1. **profile**
    - target profile in *Windows Terminal*.
    - specify either the index used by Windows Terminal, or the profile name
    - *args*: `default`, `1`, `2`, ... `"profile name"`
2. **interval**
    - time in seconds between each image change.
    - *args*: any positive integer 
3. **port**
    - port that the tbg server uses
    - *args*: any positive integer
4. **paths** 
    - paths containing images used in changing the background image of Windows
    Terminal
    - *args*:
        - `[]`
        - `- path: /path/to/dir1` 
        - ```yaml
            - path: /path/to/dir2
              alignment: center      # optional
              stretch: uniformToFill # optional
              opacity: 1.0           # optional

---
# Commands
```bash
tbg <command-name>
```
For a more detailed explanation on each command, follow the command name links

1. [run](/docs/run_command_usage.md) 
    - edit `settings.json` used by *Windows Terminal* using settings from
    **tbg**'s config. If any of the flags are specified, it will use those
    values in editing `settings.json` instead of what's specified in **tbg**'s
    config
    - use `--config` to use a custom config over the default one
    - *arg*: no arg 
    - *flags*: `-a, --alignment`, `-c, --config`, `-i, --interval`, 
    `-o, --opacity`, `-p, --profile`, `-s, --stretch`
2. [config](/docs/config_command_usage.md) 
    - If no flags are present, it will print out **tbg** config  to console. If
    any of the flags are present, it will edit the fields of the config based on
    the flags and values passed
    - Uses the default config when no arg is given
    - *arg*: no arg, `/path/to/config/file.yml` 
    - *flags*: `-i, --interval`, `-p, --profile`, `-P, --port`
3. [add](/docs/add_command_usage.md) 
    - Add a directory containing images to **tbg**'s config
    - use `--config` to target a custom config
    - If any flags are present, this will add those options to that path,
    regardless of whether the path already exists in the config or not
    - *arg*: `/path/to/dir` 
    - *flags*: `-c, --config`, `-a, --alignment`, `-o, --opacity`, `-s, --stretch` 
4. [remove](/docs/remove_command_usage.md) 
    - Remove a path from **tbg**'s config
    - use `--config` to target a custom config
    - If any flags are present, it will remove only those options of that path
    - *arg*: `/path/to/dir` 
    - *flags*: `-c, --config`, `-a, --alignment`, `-o, --opacity`, `-s, --stretch` 
5. help
    - Prints the general help message when no arg is given
    - Prints the help message/s of command/s if specified
    - *arg*: no arg, or any command (can be multiple)

## [Server Commands](/docs/server_commands_usage.md)
These commands only work when there's a **tbg** server active. Usage is the
same as other commands.
1. next-image
    - triggers an image change
    - *arg*: `/path/to/dir` 
- *flags*: `-a, --alignment`, `-o, --opacity`, `-P, --port`, `-s, --stretch`
2. set-image
    - sets a specified image as the background image
    - *arg*: `/path/to/image/file` 
    - *flags*: `-a, --alignment`, `-o, --opacity`, `-P, --port`, `-s, --stretch`
3. quit
    - stops the server
    - *arg*: none
    - *flags*: `-P, --port`

*Tip: you can assign these commands to keybinds*

---
# Example Setup
I recommend to start a tbg server on a specific port for each shell instance
of the same type. For example, all `pwsh` shells start a server on port `9545`.
Since there can only be one server on port `9545`, servers spawned from other
`pwsh` instances will fail and will not overlap with the first one that spawed.
This means that as long as one `pwsh` instance exists, a **tbg** server for the
`pwsh` profile to change background images will exist.

In the following section, I'll give examples on how to:
1. start tbg server in the background on each shell instance
2. map `alt+i` to change image (`tbg next-image`)
3. map `ctrl+alt+i` to stop the server (`tbg quit`)

for both pwsh and wsl. This way, two **tbg** servers can run simultaneously
without conflict. Keybinds will target the correct **tbg** server instance too
by specifing the same port used in starting it.

## pwsh
For example, in your `$env:PROFILE`, do:
```powershell
# Set a port for all your pwsh instances
$TBG_PORT=9545

# Auto start tbg server at every new pwsh instance.
# Since tbg uses the one port defined in the config, starting another server
# through the below command will fail quietly in the background.
Start-Job -Name tbg-server -ArgumentList $TBG_PORT -ScriptBlock {
    param($port)
    tbg.exe run --profile pwsh --port $port
} | Out-Null

# change image through Alt+i (case sensitive)
Set-PSReadLineKeyHandler -Chord "Alt+i" -ScriptBlock {
    tbg.exe next-image -P $TBG_PORT &
}

# quit server through Ctrl+Alt+i (case sensitive)
Set-PSReadLineKeyHandler -Chord "Ctrl+Alt+i" -ScriptBlock {
    tbg.exe quit --port $TBG_PORT &
}
```

## zsh (on wsl)
For example, in your `~/.zshrc`, do:
```bash
# Set a port for all your wsl Debian instances
TBG_PORT=9000

# Auto start tbg server at every new wsl Debian instance.
# Reference the exe built for windows since windows have the proper environment
# variables needed to edit wt's settings.json
tbg.exe run --profile Debian --port $TBG_PORT &>/dev/null &!

# register functions as zle widget
function __tbg_next_image() {
  tbg.exe next-image --port $TBG_PORT &>/dev/null &!
}
function __tbg_quit() {
  tbg.exe quit --port $TBG_PORT &>/dev/null &!
}
zle -N __tbg_next_image
zle -N __tbg_quit

# change image through Alt+i
bindkey '^[i' __tbg_next_image

# quit server through Ctrl+Alt+i
bindkey '^[^i' __tbg_quit
```

---
_Note: as an alternative to this approach, you can pass in a custom config for
each shell instead of specifying the port and profile. This way you can even
have different paths for each shell_

---
# Credits
- [Windows Terminal](https://github.com/microsoft/terminal)
- [levenshtein](github.com/agnivade/levenshtein) for suggesting similar profile names
- [lumberjack](https://github.com/natefinch/lumberjack) for easy rotating logs
- [carnac](https://github.com/Code52/carnac) for showing key presses on demo
- [saltkid](https://github.com/saltkid)
