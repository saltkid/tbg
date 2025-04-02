This directory contains build scripts.

> Why not make?

Well idk. I like making my own...

# Usage
To build normally:
```bash
# bash
./scripts/build.sh

# powershell and command prompt
# powershell will auto use .ps1; cmd will auto use .bat
.\scripts\build
```
This will build a `tbg.exe` in the `./bin/` directory. To check,
`./bin/tbg.exe version` should have a `-dev` suffix

---
For releasing, do:
```bash
# bash
./scripts/build.sh release

# powershell and command prompt
# powershell will auto use .ps1; cmd will auto use .bat
.\scripts\build release
```
This will build 4 executables in the `./bin/` directory with the naming format:
```
tbg-$version.windows.$arch.exe
```
where `$arch` are amd64, 386, arm64, and arm. Each executable will have their
respective sha256 checksum.
