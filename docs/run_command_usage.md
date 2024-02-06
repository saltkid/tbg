# `tbg run`

`run` command edits the `settings.json` used by *Windows Terminal* using settings from the currently used config.

**tbg** will keep running, editing the `settings.json` of *Windows Terminal*, replacing the background image. You can quit by pressing `q` or `ctrl+c`

On initial execution of **tbg**, it will create the `tbg_profile.yaml` and `config.yaml` in the same directory as the **tbg** executable if it does not exist already. This is so it can safely fallback on a `config.yaml`. **There can only be one `tbg_profile.yaml`**. You can have multiple `config.yaml`s if you want. `tbg_profile.yaml` will keep track of which one you want to use.

For more information, see documentation on [tbg profile](#link) and [config](#link).

## Key events
**tbg** takes optional commands during execution:
- `n`: goes to next image in the current image collection dir
- `p`: goes to previous image in the current image collection dir
- `f`: goes to next image collection dir
- `b`: goes to previous image collection dir
- `c`: shows the available commands

**tbg** will continue running until you press `q` or `ctrl+c`. This means even if all images are exhausted, **tbg** will safely wrap back around.

For an example, see [walkthrough with key events](#normal-execution,-using-used_config-specified-in-tbg_profile.yaml,-path-with-flags)

## Executing with `--`flags
`--alignment`, `--opacity`, `--stretch`

`--`flags can be used to override the default flag fields in the config: `default_alignment`, `default_stretch`, `default_opacity`.

These flags will override the flags set per path as well. So if there is a `path/to/dir | center fill 0.1`, **tbg** will use the `--`flags instead of that or the default flag fields.

The order of importance is:
1. `--`flags
2. per path flags
3. default flag fields 

For an example, see [overriding default flags walkthrough](#overriding-default-flags)

## Executing with `config` subcommand
Using the `config` command will override the `tbg_profile.yaml`'s `used_config` field. This means it will force **tbg** to use the specified config instead of the default behavior of checking `used_config` field of `tbg_profile.yaml`

Note that this will not edit the `used_config` field in `tbg_profile.yaml`, only override it for the current execution.

For an example, see [Using different config walkthrough](#using-different-config,-overriding-profile-and-interval-fields)


## Walkthroughs
#### Normal Execution, using used_config specified in tbg_profile.yaml, path with flags
Let's use this config:
```
# config.yaml
profile: default
interval: 30

image_col_paths:
- path/to/dir1
- path/to/dir2 | right fill 0.35

default_alignment: center
default_stretch: fill
default_opacity: 0.1
```
This just means that when we do `tbg run`, we want to change the background image of the **default** *Windows Terminal* profile every **30 minutes**.

The first few images will be from `path/to/dir1`. The image will be at the **center**, **fill** the entire screen without regard of the original aspect ratio, with an opacity of **10%**. 

When I press `n`, it goes to the next images without waiting for 30 minutes. I can go back by pressing `p`.

When I press `f`, it goes to the next image collection dir. This means we are now in `path/to/dir2`. This path has flags so these values will be used instead of the default flag fields. This means instead of the image being at the **center**, it will be at the **right**. Instead of having an opacity of **10%**, it the images will have **35%** opacity. However, since stretch is the same, it will still **fill** the screen without regard of the orignal aspect ratio.

When I press `b`, it goes back to the previous image collection dir (`path/to/dir1`). If i press `b` again, it will wrap around and go to the last image collection dir (`path/to/dir2`).

This wrap around behavior also applies to `f`.

---
#### Using different config, overriding profile and interval fields
Now let's quit **tbg** by pressing `q` or `ctrl+c`.

Instead of `tbg run`, let's do:
```
tbg run config path/to/config-2.yaml  --profile list-1 --interval 5
```
```
# config-2.yaml
profile: default
interval: 30

image_col_paths:
- path/to/dir1

default_alignment: center
default_stretch: fill
default_opacity: 0.1
```

The `config` subcommand tells **tbg** to use `path/to/config-2.yaml` instead of whatever `tbg_profile.yaml`'s `used_config` is pointing to. This will not edit `tbg_profile.yaml`'s `used_config` field. 

The `--profile` and `--interval` flags will override the values in `config-2.yaml`. Again, not edit them. This means instead of changing the background image of the `default` profile every 30 minutes, it will change the background image of the first profile under `list` field in `settings.json`

---
#### Overriding default flags
Let's quit **tbg** again by pressing `q` or `ctrl+c`. Let's use the default config again.
```
# config.yaml
profile: default
interval: 30

image_col_paths:
- path/to/dir1
- path/to/dir2 | right fill 0.35

default_alignment: center
default_stretch: fill
default_opacity: 0.1
```
Let's do:
```
tbg run --alignment right --opacity 0.35 --stretch none
```

The `--alignment`, `--opacity`, and `--stretch` flags will override the values in `config.yaml`. This means instead of `path/to/dir`'s images being centered, filling entire screen at 10% opacity, the images will be on the right, with no upscaling/downscaling, at 35% opacity.

Notice that `path/to/dir2` has flags that should override the default flag fields. However, since we specified `--alignment right --opacity 0.35 --stretch none`, we'll use these value instead, like what we did with `path/to/dir1`.
