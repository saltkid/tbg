# Table of Contents
- [Server Commands](#server-commands)
  - [Commands](#commands)
- [Keybind Examples](#keybind-examples)
---

# Server Commands
```zsh
tbg <server-command>
```
These commands only work when there's a **tbg** server active. See [run](https://github.com/saltkid/tbg/blob/main/docs/run_command_usage.md)
to learn more about starting a **tbg** server. Usage is the same as other
commands.

## Commands
All available APIs have an associated command. If there is no command for an
action, there is no API for it.
1. next-image
    - triggers an image change in the currently running **tbg** server
    - if no server is found, this will fail
    - curl equivalent: `curl -X POST localhost:9545/next-image`
2. quit
    - stops the currently running **tbg** server
    - if no server is found, this will fail
    - curl equivalent: `curl -X POST localhost:9545/quit`

These are useful when integrating it with the shell through keybinds.
# Keybind Examples
1. powershell

In your `$PROFILE`, do:
```powershell
# Auto start tbg server at every new pwsh instance.
# Since tbg uses the one port defined in the config, starting another server
# through `tbg run` will fail quietly in the background.
Start-Job -Name tbg-server -ScriptBlock { tbg.exe run -p pwsh } | Out-Null

# change image through ctrl+i
Set-PSReadLineKeyHandler -Key "Ctrl+i" -ScriptBlock {
  tbg.exe next-image &
}

# quit server through Ctrl+Alt+i
Set-PSReadLineKeyHandler -Key "Ctrl+Alt+i" -ScriptBlock {
  tbg.exe quit &
}
```
2. zsh (in wsl)

Assuming the wsl distro is Debian, in your `.zshrc`, do:
```zsh
# Auto start tbg server at every new wsl Debian instance.
# Reference the exe built for windows since windows have the proper environment
# variables needed to edit wt's settings.json
tbg.exe run -p Debian &>/dev/null &!

# register functions as zle widget
function __tbg_next_image() {
  tbg.exe next-image &>/dev/null &!
}
function __tbg_quit() {
  tbg.exe quit &>/dev/null &!
}
zle -N __tbg_next_image
zle -N __tbg_quit

# change image through ctrl+i
bindkey '^I' __tbg_next_image

# quit server through Ctrl+Alt+i
bindkey '^[^I' __tbg_quit
```
