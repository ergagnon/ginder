# ginder

## Windows
Mingw
```bash
choco install mingw
```

Pkg-config
```bash
choco install pkgconfiglite
```

## Build docker
```bash
docker build -t ginder -f docker/ginder.dockerfile
```

## WSL
```bash
sudo apt-get update
sudo apt-get install -yq libhyperscan-dev libpcap-dev
sudo apt install gcc
sudo apt install g++
```

## Documentation
* https://github.com/flier/gohs
* https://intel.github.io/hyperscan/dev-reference/getting_started.html
* https://sourceforge.net/projects/pcre/files/pcre/8.41/
* https://github.com/intel/hyperscan/issues/86                           

## VCPKG
### Install 
[Install documentation](https://learn.microsoft.com/en-us/vcpkg/get_started/get-started?pivots=shell-powershell#1---set-up-vcpkg)

### Build Hyperscan
Install Visual studio Community with C++

```powershell
$env:VCPKG_LIBRARY_LINKAGE = "dynamic"
vcpkg install hyperscan
```

### Install pkg-config-lite
https://community.chocolatey.org/packages/pkgconfiglite

Create a folder containing pkg-config files
Set env var to the previous variable
PKG_CONFIG_PATH=C:\Dev\pkg-config


Create a file libhs.pc in C:\Dev\pkg-config

```
prefix=C:/Dev/vcpkg/packages/hyperscan_x64-windows
exec_prefix=${prefix}
includedir=${prefix}/include
libdir=${exec_prefix}/lib

Name: libhs
Description: Intel(R) Hyperscan Library
Version: 5.4.2
Libs: -L${libdir} -lhs
Libs.private:
Cflags: -I${includedir}/hs
```





