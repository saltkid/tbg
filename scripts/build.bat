@echo off
setlocal enabledelayedexpansion

set "arg=%1"
if "%arg%" == "release" (
    call :__build_tbg arm
    call :__build_tbg arm64
    call :__build_tbg amd64
    call :__build_tbg 386
) else if "%arg%" == "" (
    call :__build_tbg dev
) else (
    echo Invalid argument: "%arg%". Please use either "release" or leave the argument empty.
    exit /b 1
)
goto :__end_build_tbg

:__build_tbg
set arch=%1
if not exist .\bin (
    mkdir .\bin
)
for /f "delims=" %%v in ('git describe --tags --dirty --always') do set version=%%v
if "%arch%" neq "arm" if "%arch%" neq "arm64" if "%arch%" neq "amd64" if "%arch%" neq "386" if "%arch%" neq "dev" (
    echo Invalid GOARCH value: %arch%
    echo Please provide one of: arm, arm64, amd64, or 386.
    exit /b 1
)
set GOOS=windows
if not "%arch%" == "dev" (
    set GOARCH=%arch%
)
go mod tidy
if "%arch%" == "dev" (
    go build -ldflags "-X main.TbgVersion=%version%-dev"
    move /y tbg.exe .\bin\tbg.exe >nul 2>&1
) else (
    go build -ldflags "-s -w -X main.TbgVersion=%version%"
    move /y tbg.exe .\bin\tbg-%version%.windows.%arch%.exe >nul 2>&1
    certutil -hashfile .\bin\tbg-%version%.windows.%arch%.exe SHA256 | findstr /v "SHA256" | findstr /v "CertUtil" > .\bin\tbg-%version%.windows.%arch%.exe.sha256
)
set GOOS=
set GOARCH=
set arch=

:__end_build_tbg
set arg=
exit /b

endlocal
