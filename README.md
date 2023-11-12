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
