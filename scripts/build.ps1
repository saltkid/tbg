param($arg)

function __build_tbg {
    param($arch)
    if (-not (Test-Path -Path "./bin")) {
        New-Item -ItemType Directory -Path "./bin" | Out-Null
    }
    $version = git describe --tags --dirty --always
    if ($arch -ne "arm" -and $arch -ne "arm64" -and $arch -ne "amd64" -and $arch -ne "386" -and $arch -ne "dev") {
        Write-Host "Invalid GOARCH value: $arch"
        Write-Host "Please provide one of: arm, arm64, amd64, or 386."
        return
    }
    $env:GOOS = "windows"
    if ($arch -ne "dev") {
        $env:GOARCH = "$arch"
    }
    go mod tidy
    if ($arch -eq "dev") {
        go build -ldflags "-X main.TbgVersion=$version-dev"
        Move-Item -Path "./tbg.exe" -Destination "./bin/tbg.exe"
    } else {
        go build -ldflags "-s -w -X main.TbgVersion=$version"
        Move-Item -Path "./tbg.exe" -Destination "./bin/tbg-$version.windows.$arch.exe"
        $checksum = Get-FileHash -Algorithm SHA256 -Path "./bin/tbg-$version.windows.$arch.exe"
        $checksum.Hash | Out-File -FilePath "./bin/tbg-$version.windows.$arch.exe.sha256"
    }

    $env:GOOS = $null
    $env:GOARCH = $null
}

if ($arg -eq "release") {
    __build_tbg "arm"
    __build_tbg "arm64"
    __build_tbg "amd64"
    __build_tbg "386"
} elseif ($arg -eq $null) {
    __build_tbg "dev"
} else {
    Write-Host "Invalid argument: `"$arg`". Please use either `"release`" or leave the argument empty."
    exit 1
}
