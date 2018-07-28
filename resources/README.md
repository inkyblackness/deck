# InkyBlackness - Deck

**Obsolescence Notice: Due to the release of [InkyBlackness - HackEd](https://github.com/inkyblackness/hacked), this package here has become obsolete: HackEd is the new editor to be used and is fully self-contained. Furthermore, this package is based on pre-source-release information and will no longer be maintained. Use it at your own risk.**

This package project contains the public distribution of the [InkyBlackness](https://inkyblackness.github.io) binaries.

The executables are modding tools for the DOS game [System Shock (1994)](http://en.wikipedia.org/wiki/System_Shock).

## Disclaimer
For all tools, a general warning: Keep backups of your files!
These are executables downloaded from the internet. If you are unsure what they do or how they work, refer to the sources [on GitHub](https://github.com/inkyblackness).


## Package Contents

### InkyBlackness - Shocked ("shocked-client")
This tool is the primary application. It is a GUI editor for the resources, including a map editor.
The tool uses OpenGL 3.2 .
User documentation is available online here: https://github.com/inkyblackness/shocked-client/wiki

#### Starting the editor

Example:

```
shocked-client --path <path to resources>
```

The ```--path``` parameter can be provided several times and point to one directory with resources files each time. If all resources are in the same directory, only one is enough. Two might be used for a vanillay System Shock installation, with one directory pointing to the resources on the CD-ROM and the other to the files on harddisk.
Note that the provided path(s) must point to an actual data directory, within which the .res files are stored.


### InkyBlackness - Construct
This console tool is ideally used in combination with the map editor. It can generate a small ```archive.dat``` file with empty levels.


### InkyBlackness - Chunkie
This is a resource import/export tool for the console. Some raw media formats are supported with their respective file formats.


### InkyBlackness - Hacker
This is a low-level file access tool for the console; A better hex-editor if you will. One can modify all bytes of the resources in any way. One can also crash the game with it.


### InkyBlackness - Reactor RNG
This command-line tool randomizes the reactor code of save-game files that specific versions of the "System Shock Enhanced Edition" set to a static value.
