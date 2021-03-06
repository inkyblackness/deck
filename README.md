# InkyBlackness Deck

**Obsolescence Notice: Due to the release of [InkyBlackness - HackEd](https://github.com/inkyblackness/hacked), this package here has become obsolete: HackEd is the new editor to be used and is fully self-contained. Furthermore, this package is based on pre-source-release information and will no longer be maintained. Use it at your own risk.**

This is the release package of the [InkyBlackness](https://inkyblackness.github.io) project.
It contains the collection of the dependencies as submodules, to bind them to specific versions.

Releases of InkyBlackness are created through this project.

## Building
The base system for building the binaries is a Linux system. The scripts will cross-compile to MS Windows and require
the presence of a mingw64 compiler.

### Updating
The script ```update.sh``` removes the src directory and re-downloads all InkyBlackness components and their dependencies.

### Compiling
The script ```build.sh``` then compiles all binaries and places the package contents in the subfolders of ```dist```.


## License

The project is available under the terms of the **New BSD License** (see LICENSE file).
