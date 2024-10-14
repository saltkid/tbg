# Table of Contents
- [Overview](#tbg-run)
- [Key events](#key-events)
- [Executing with flags](#executing-with-flags)
    - [Using `--random` flag](#using---random-flag)
- [Usage](#usage)
    - [Normal Execution, key events, and path specific options](#normal-execution)
    - [Overriding `profile` and `interval` fields](#overriding-profile-and-interval-fields)
    - [Overriding default option fields](#overriding-default-option-fields)
---

# `tbg run`

`run` command edits the `settings.json` used by *Windows Terminal* using
settings from the currently used config. **tbg** will keep running, editing the
`settings.json` of *Windows Terminal*, replacing the background image. You can
quit by pressing `q` or `ctrl+c`

On initial execution of **tbg**, it will create a `config.yaml` in the same
directory as the **tbg** executable if it does not exist already. **There can
only be one `config.yaml`**. For more information, see documentation on
[config.yaml](https://github.com/saltkid/tbg/blob/main/docs/config.yaml.md).

# Key events
**tbg** takes optional commands during execution:
- `q`: quits
- `c`: shows the available commands
- `n`: goes to next image
- `p`: goes to previous image
- `N`: goes to next image collection
- `P`: goes to previous image collection
- `r`: randomizes the images in the current image collection starting from the
current image
    - this does not affect the order of the previous images
- `R`: randomizes the images in the current image collection starting from the
current collection
    - this does not affect the order of the previous collections

**tbg** will continue running until you press `q` or `ctrl+c`.
This means even if all images are exhausted, **tbg** will wrap back around.
For an example, see [usage with key events](#normal-execution)

# Executing with flags
#### Valid Flags: `--profile`, `--interval`, `--alignment`, `--opacity`, `--stretch`, `--random`

`--`flags can be used to override these fields in the config:
`profile`, `interval`, `alignment`, `stretch`, `opacity`.

The flags that override default options fields will override the options set
per path as well. So if there is a `path/to/dir | center fill 0.1`, **tbg**
will use the flags instead of that or the default options fields.

The order of importance is:
1. flags (`--alignment`, `--opacity`, `--stretch`)
2. per path options in config
3. default options fields in config

For an example, see [overriding default flags walkthrough](#overriding-default-option-fields)

#### Using `--random` flag
The `--random` flag will randomize the order of image collections dirs **and**
the images. Even if you consumed all paths and wrap around to the first path
again, the paths will be re-randomized. Even the images: when you go to next
dir B, then go to previous dir A, the order will be different from the first
time you go to dir A, since images are also re-randomized every time you enter
a dir.

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
- path: path/to/dir1
- path: path/to/dir2
  alignment: right
  stretch: fill
  opacity: 0.35

profile: default
interval: 30

alignment: center
stretch: fill
opacity: 0.1
```
This just means that when we do `tbg run`, we want to change the background
image of the **default** *Windows Terminal* profile every **30 minutes**. The
first few images will be from `path/to/dir1`. The image will be at the
**center**, **fill** the entire screen without regard of the original aspect
ratio, with an opacity of **10%**. 

When I press `n`, it goes to the next images without waiting for 30 minutes. I
can go back by pressing `p`.

When I press `N`, it goes to the next image collection dir. This means we are
now in `path/to/dir2`. This path has flags specific to it so these values will
be used instead of the default flag fields. This means instead of the image
being at the **center**, it will be at the **right**. - Instead of having an
opacity of **10%**, it the images will have **35%** opacity. However, since
stretch is the same, it will still **fill** the screen without regard of the
orignal aspect ratio.

When I press `P`, it goes back to the previous image collection dir
(`path/to/dir1`). If i press `P` again, it will wrap around and go to the last
image collection dir (`path/to/dir2`). This wrap around behavior also applies
to `N`.

Now let's quit **tbg** by pressing `q` or `ctrl+c`.

---
### Overriding `profile` and `interval` fields

Instead of `tbg run`, let's do:
```bash
tbg run --profile list-1 --interval 5
```
```yml
paths:
- path: path/to/dir1

profile: default
interval: 30

alignment: center
stretch: fill
opacity: 0.1
```

The `--profile` and `--interval` flags will override the values in the config.
Again, not edit them. This means instead of changing the background image of
the `default` profile every 30 minutes, it will change the background image of
the first profile under `list` field in `settings.json` every 5 minutes instead

---
### Overriding default option fields
This will delve on overriding default options fields on the config using
flags. This will also override the per-path options.

Let's use this config:
```yml
paths:
- path: path/to/dir1
- path: path/to/dir2 
  alignment: left
  stretch: uniformToFill
  opacity: 0.25

profile: default
interval: 30

alignment: center
stretch: fill
opacity: 0.1
```
Let's do:
```bash
tbg run --alignment right --opacity 0.35 --stretch none
```

The `--alignment`, `--opacity`, and `--stretch` flags will override the values
in `.tbg.yml`. This means instead of `path/to/dir`'s images being centered,
filling entire screen at 10% opacity, the images will be on the right, with no
upscaling/downscaling, at 35% opacity.

Notice that `path/to/dir2` has options that should override the default options
fields. However, since we specified `--alignment right --opacity 0.35 --stretch
none`, tbg will use these value instead, like what we did with `path/to/dir1`.
