# Table of Contents
- [Log Types](#log-types)
  1. [tbg initialization](#tbg-initialization)
  2. [Starting tbg server](#starting-tbg-server)
  3. [Edited Windows Terminal's `settings.json`](#edited-windows-terminals-settingsjson-to-change-the-background-image)
  4. [Automatic image change at every n-interval](#automatic-image-change-at-every-n-interval)
  5. [Changing image through `tbg next-image`](#changing-image-through-tbg-next-image)
  6. [Setting a specific image as the background image through `tbg set-image`](#setting-a-specific-image-as-the-background-image-through-tbg-set-image)
  7. [Quit server through `tbg quit`](#quit-server-through-tbg-quit)

---
# Log Types
All logs have a `"time": "TIMESTAMP"` and `"level": "INFO"` so I will leave
these two out.

---
### tbg initialization
```json
{
  "msg": "Start tbg...Config used",
  "paths": [
    {
      "opacity": 0.1,
      "path": "/path/to/images/dir1"
    },
    {
      "opacity": 0.5,
      "path": "/path/to/images/dir2"
      "stretch": "uniform"
    },
    {
      "alignment": "right",
      "opacity": 1.0,
      "path": "/path/to/images/dir3"
      "stretch": "fill"
    }
  ],
  "interval": 20,
  "port": 8000,
  "profile": "default"
}
```

---
### Starting tbg server
    - override-[alignment,opacity,stretch] is set through their respective
    flags
    - e.g. `tbg run --alignment center`
```json
{
  "msg": "Starting server...",
  "interval": 20,
  "port": ":8000",
  "profile": "default",
  "override-alignment": "center",
  "override-opacity": "no override",
  "override-stretch": "no override"
}
```

---
### Edited Windows Terminal's `settings.json` to change the background image
```json
{
  "msg": "Starting server...",
  "image": "/path/to/image/file.png",
  "profile": "default",
  "alignment": "center",
  "opacity": "0.25",
  "stretch": "uniformToFill"
}

```

---
### Automatic image change at every n-interval
```json
{ "msg": "Image change tick" }
```

---
### Changing image through `tbg next-image`
...or by making a POST request to the `next-image` endpoint
    - there are multiple logs during this action
```json
{ "msg": "Recieved next-image request" }
{ "msg": "next-image body decoded" }

```
_below may or may not be logged, depending on whether the option is set_
```json
{
  "msg": "alignment",
  "value": "center"
}
{
  "msg": "opacity",
  "value": "0.1"
}
{
  "msg": "stretch",
  "value": "fill"
}
```

---
### Setting a specific image as the background image through `tbg set-image`
...or by making a POST request to the `set-image` endpoint
    - there are multiple logs during this action
```json
{ "msg": "Recieved set-image request" }
{ "msg": "set-image body decoded" }
{
  "msg": "path",
  "value": "/path/to/image/file.png"
}
```
_below may or may not be logged, depending on whether the option is set_
```json
{
  "msg": "alignment",
  "value": "center"
}
{
  "msg": "opacity",
  "value": "0.1"
}
{
  "msg": "stretch",
  "value": "fill"
}
```

---
### Quit server through `tbg quit`
...or by makign a POST request to the `quit` endpoint
```json
{ "msg": "Recieved quit request" }
```
_below logs after cleaning up_
```json
{ "msg": "Goodbye!" }
```
