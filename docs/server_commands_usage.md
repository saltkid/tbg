# Table of Contents
- [Server Commands](#server-commands)
  - [Commands](#commands)
- [Keybind Examples](#keybind-examples)
---

# Server Commands
```zsh
tbg <server-command>
```
These commands only work when there's a **tbg** server active. See [run](/docs/run_command_usage.md)
to learn more about starting a **tbg** server. Usage is the same as other
commands.

## Commands
All available APIs have an associated command. If there is no command for an
action, there is no API for it.
1. next-image
    - valid flags: `-P, --port`, `-a, --alignment`, `-o, --opacity`, `-s, --stretch`
      - these will override the image properties of the next randomly chosen
      image
    - triggers an image change in the currently running **tbg** server at port
    9545 if no port is given
    - if no server is found, this will fail
2. set-image
    - arg: `/path/to/image/file`
    - valid flags: `-P, --port`, `-a, --alignment`, `-o, --opacity`, `-s, --stretch`
    - sets the specified image as the background image through an image change
    in the currently runing **tbg** server at port 9545 if no port is given
    - the default values for each will be used if not specified
    - if no server is found, this will fail
3. quit
    - valid flags: `-P, --port`
    - stops the currently running **tbg** server at port 9545 if no port is
    given
    - if no server is found, this will fail

These are useful when integrating it with the shell through keybinds.
# Keybind Examples
1. powershell

In your `$PROFILE`, do:
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
Set-PSReadLineKeyHandler -Key "Ctrl+Alt+i" -ScriptBlock {
    tbg.exe quit --port $TBG_PORT &
}
```
2. zsh (in wsl)

Assuming the wsl distro is Debian, in your `.zshrc`, do:
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
