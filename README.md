# Lapitar [![Build Status](https://travis-ci.org/LapisBlue/Lapitar.svg?branch=master)](https://travis-ci.org/LapisBlue/Lapitar)
[Lapitar] is a new **open source** Minecraft avatar service, providing 2D and 3D avatars with a simple API. If you just want to test it out or use it on your website you can check out how to use Lapitar [on our website][Lapitar]. This readme describes how to setup your own [Lapitar] installation.

## Introduction
[Lapitar] is written in a mixture of [Go] and C. Most parts like the web server or caching are written in [Go], but the actual 3D rendering is done with OpenGL in C and referenced from [Go]. Due to its dependency on [OSMesa], a library used to render images with OpenGL off-screen, it works only on Linux right now. We're looking forward to changing this in the future.

## Installation

Lapitar is primarily tested on the latest Ubuntu LTS and Go release. Other distributions or older Go versions may work, but they're not tested. You may have to compile [OSMesa] yourself if the version supplied by your distribution isn't working properly.

- Install the latest version of Go, refer to this guide for the steps to install it: https://golang.org/doc/install

    - On Debian/Ubuntu you can use [godeb] to install the latest Go version easily. Follow these steps:

      ```
      wget https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz
      tar xzvf godeb-amd64.tar.gz
      ./godeb install
      ```

      Now you only need to set a permanent environment variable for your `GOPATH`, for example to `~/go`.

- Install the native dependencies, the steps required for this depend on the distribution you're using:

    - Ubuntu: `sudo apt-get install build-essential pkg-config libosmesa6-dev libglu1-mesa-dev`

- Install and compile Lapitar by executing the following command: `go get github.com/LapisBlue/lapitar/lapitar`
- The executable will be created in `$GOPATH/bin/lapitar`.
- If you want to update Lapitar later you can execute `go get -u github.com/LapisBlue/lapitar/lapitar`.

[Lapitar]: https://lapitar.lapis.blue
[Go]: https://golang.org
[OSMesa]: http://www.mesa3d.org/osmesa.html
[godeb]: https://github.com/niemeyer/godeb
