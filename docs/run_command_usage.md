# Table of Contents
- [Overview](#tbg-run)
- [Key events](#key-events)
    - [Order Behavior](#ordering-behavior)
    - [Order Behavior Using `--random` flag](#ordering-behavior-using-random-flag)
- [Executing with flags](#executing-with-flags)
- [Usage](#usage)
    - [Normal Execution, key events, and path specific options](#normal-execution)
    - [Overriding `profile` and `interval` fields](#overriding-profile-and-interval-fields)
    - [Overriding default values](#overriding-default-values)
---

# `tbg run`

`run` command edits the `settings.json` used by *Windows Terminal* using
settings from `.tbg.yml`. **tbg** will keep running, editing the
`settings.json` of *Windows Terminal*, replacing the background image. You can
quit by pressing `q` or `ctrl+c`

On initial execution of **tbg**, it will create a `.tbg.yml` in the same
directory as the **tbg** executable if it does not exist already. **There can
only be one `.tbg.yml`**. For more information, see documentation on
[tbg.yml](https://github.com/saltkid/tbg/blob/main/docs/tbg.yml.md).

# Key events
**tbg** takes optional commands during execution:
- `q`: quits
- `c`: shows the available commands
- `n`: goes to next image
- `p`: goes to previous image
- `N`: goes to next images path
- `P`: goes to previous images path
- `r`: randomizes the images in the current images path starting from the
current image
    - this does not affect the order of the previous images so you can go to
    the same previous image/s
- `R`: randomizes the images in the current images path starting from the
current collection
    - this does not affect the order of the previous collections

**tbg** will continue running until you press `q` or `ctrl+c`.
This means even if all images are exhausted, **tbg** will wrap back around.
For an example, see [usage with key events](#normal-execution)

## Ordering Behavior
The order of paths **and** the images in that path are randomized on
initialize. However, you'd still have to consume all the images in a path
before going to the next one. When you consumed all paths and wrap around to
the first path again, the paths will be re-randomized. Even the images: from
`path A`, when you go to next `path B`, then go to previous `path A`, the order
of images in `path A` will be different from the first time you went to it,
since images are also re-randomized every time you enter an images path.

## Ordering Behavior using `--random` flag
The `--random` flag will ensure that whenever you go to the next image, it
always will pick a random images path, then a random image from there. This
means going to the next image `[n]` is the only valid command by the user. You
cannot go to the next path `[N]`, previous image `[p]`, previous path `[P]`,
randomize images `[r]`, or randomize paths `[R]`. More info about flags in the
next section

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
- path: path/to/dir1
- path: path/to/dir2
  alignment: right
  stretch: fill
  opacity: 0.35

profile: default
interval: 30
```
This just means that when we do `tbg run`, we want to change the background
image of the **default** *Windows Terminal* profile every **30 minutes**. The
first few images will be from `path/to/dir1`. The image will be at the
**center**, with the stretch **uniform** to fill the screen while keeping
aspect ratio, with an opacity of **100%**. 

When I press `n`, it goes to the next image without waiting for 30 minutes. I
can go back by pressing `p`.

When I press `N`, it goes to the next image collection dir. This means we are
now in `path/to/dir2`. This path has flags specific to it so these values will
be used instead of the defaults. This means instead of the image being at the
**center**, it will be at the **right**. - Instead of having an opacity of
**10%**, it the images will have **35%** opacity. The image will **fill** the
screen, without regard for aspect ratio.

When I press `P`, it goes back to the previous image collection dir
(`path/to/dir1`). If i press `P` again, it will wrap around and go to the last
image collection dir (`path/to/dir2`). This wrap around behavior also applies
to `N`.

Now let's quit **tbg** by pressing `q` or `ctrl+c`.

---
### Overriding `profile` and `interval` fields

Instead of `tbg run`, let's do:
```bash
tbg run --profile 1 --interval 5
```
```yml
paths:
- path: path/to/dir1

profile: default
interval: 30
```

The `--profile` and `--interval` flags will override the values in the config.
Again, not edit them. This means instead of changing the background image of
the `default` profile every 30 minutes, it will change the background image of
the first profile under `list` field in `settings.json` every 5 minutes instead

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
interval: 30
```
Let's do:
```bash
tbg run --alignment right --opacity 0.35 --stretch none
```

The `--alignment`, `--opacity`, and `--stretch` flags will override the values
in `.tbg.yml`. This means instead of `path/to/dir1`'s images being centered,
filling entire screen while keeping aspect ratio at 100% opacity, the images
will be on the right, with no scaling, at 35% opacity.

Notice that `path/to/dir2` has options that should override the default options
fields. However, since we specified `--alignment right --opacity 0.35 --stretch
none`, tbg will use these value instead, like what we did with `path/to/dir1`.
