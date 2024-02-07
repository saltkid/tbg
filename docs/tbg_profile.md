# `tbg_profile.yaml`
**Note that it is not recommended to edit this file.**

`tbg_profile.yaml` is the file **tbg** uses to keep track of what is the currently used config so you don't have to keep specifying a `config` subcommand everytime you do `run`, `add`, `remove`, and `config`. This is auto generated on initial execution of **tbg**. It should look like this
```
#---------------------------------------------
# this is a tbg profile. Whenver tbg is ran, it will
# load this profile to get the currently used config
#
# currently, it only has one field: used_config
# I'll add more if the need arises
#---------------------------------------------

used_config: path/to/default/config.yaml 

#---------------------------------------------
# Fields:
#   used_config: path to the config used by tbg
#------------------------------------------
```
It only has one field: `used_config` which keeps track of the currently used config.
Whenever `tbg config [arg]` is ran, which sets the currently used config to whatever the `[arg]` is, **tbg** edits `tbg_profile.yaml`'s `used_config` field.
