# Tar
An avatar generator API thing.

An mixture of C and Go. Go is used to provide the web server and the actual rendering is done with OpenGL in C and referenced using cgo from Go directly.

## Installation
- Make sure you have at least Go 1.1 installed with a valid GOPATH, if not [install the latest version](http://golang.org/doc/install).
- Install osmesa, if you don't have it already:
  - Debian / Ubuntu: `sudo apt-get install libosmesa6-dev libglu1-mesa-dev pkg-config`
- Install the resize package: `go get github.com/nfnt/resize`
- Go to `$GOPATH/src/github.com`, create a new folder `LapisBlue`, then clone the repository: `git clone https://github.com/LapisBlue/Tar.git`
- Checkout the `go` branch and install the application: `git checkout go`, `go install`.
- Run the application from `$GOPATH/bin/Tar`.
