# InkyBlackness Reactor RNG

This is a tool as part of the [InkyBlackness](https://inkyblackness.github.io) project, written in [Go](http://golang.org/).
It randomizes the reactor code of save-game files that specific versions of the "System Shock Enhanced Edition" set to a static value.

## Usage

```
When multiple files are modified in one go, they all receive the same code.

Usage:
   reactor-rng <savefile>...
   reactor-rng -h | --help
   reactor-rng --version

Options:
   -h --help              Show this screen.
   --version              Show version.

```

### Example:

```
C:\Folder\To\Game\Data> reactor-rng SAVGAM00.DAT SAVGAM01.DAT

Processing <C:\Folder\To\Game\Data\SAVGAM00.DAT>...
Reading file.
Modifying game state.
Changing code.
Applying code where it is needed on Citadel.
Saving file.
Done.

Processing <C:\Folder\To\Game\Data\SAVGAM01.DAT>...
Reading file.
Modifying game state.
Changing code.
Applying code where it is needed on Citadel.
Saving file.
Done.

```

## License

The project is available under the terms of the **New BSD License** (see LICENSE file).
