# Lapitar
An avatar generator API thing.

An mixture of C and Go. Go is used to provide the web server and the actual rendering is done with OpenGL in C and referenced using cgo from Go directly.

## Installation
- Make sure you have at least Go 1.1 installed with a valid GOPATH, if not [install the latest version](http://golang.org/doc/install).
- Install the native dependencies:
  - Debian/Ubuntu: `sudo apt-get install build-essential pkg-config libosmesa6-dev libglu1-mesa-dev`
- Install and compile Lapitar by executing the following command: `go get github.com/LapisBlue/Lapitar`
- The executable will be created in `$GOPATH/bin/Lapitar`.
- If you want to update Lapitar later you can execute `go get -u github.com/LapisBlue/Lapitar`.
