# nMage

[![build](https://github.com/bloeys/nmage/actions/workflows/build-nmage.yml/badge.svg)](https://github.com/bloeys/nmage/actions/workflows/build-nmage.yml)

nMage is a (hopefully!) high performance 3D Game Engine written in Go being developed [live](https://twitch.tv/bloeys), with recordings posted on [YouTube](https://www.youtube.com/channel/UCCf4qyNGPVwpj1HYFGahs_A).

This project is being built with the goals being (in no particular order):

* To build a good game engine that can be used to develop full fledged games (as opposed to just demos and proofs-of-concept)
* Sharing knowledge about less popular/less taught (e.g. compared to web dev), yet very powerful computing topics by building things and explaining as we go
* Showing the development process of large, high performance software, including things like: learning unfamiliar topics, reading docs, fixing bugs and profiling and optimizing
* Have fun through the entire thing!

## Running the code

To run the project you need:

* A recent version of [Go](https://golang.org/) installed
* A C/C++ compiler installed and in your path
  * Windows: [MingW](https://www.mingw-w64.org/downloads/#mingw-builds) or similar
  * Mac/Linux: Should be installed by default, but if not try [GCC](https://gcc.gnu.org/) or [Clang](https://releases.llvm.org/download.html)
* Install SDL2 by following their [requirements](https://github.com/veandco/go-sdl2#requirements).
* Get the required [assimp-go](https://github.com/bloeys/assimp-go) DLLs/DyLibs and place them correctly by following the assimp-go [README](https://github.com/bloeys/assimp-go#using-assimp-go).

Then you can start nMage with `go run .`

> Note: It *might* take a while to clone/run the first time because of downloading/compiling dependencies.
