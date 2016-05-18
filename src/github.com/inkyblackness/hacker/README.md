[![Build Status][drone-image]][drone-url]
[![Coverage Status][coveralls-image]][coveralls-url]

# InkyBlackness Hacker

This is a tool as part of the [InkyBlackness](https://inkyblackness.github.io) project, written in [Go](http://golang.org/). It provides low-level access to the data files of System Shock via a command line interface.

The content of the supported files is documented in the [ss-specs](https://github.com/inkyblackness/ss-specs) sub-project of InkyBlackness. This tool is meant to aid in verification of known and unknown data fields.

## Usage

Hacker can be run from any location. After startup, it provides a prompt where you can enter commands.

### Command Line Interface

```
Usage:
  hacker [--run <file>...]
  hacker -h | --help
  hacker --version

Options:
  -h --help     Show this screen.
  --version     Show version.
  --run <file>  Run the specified file. Can be repeated to run several in sequence.
```

### Command Reference

#### Quit
```
quit
```

This command exits the program. Quitting is an instant action, even if there is unsaved data.

#### Load
```
load "/path/to/hd/data/files" "/path/to/cd/data/files"
```

On startup, Hacker has no files loaded. With this command, the names of the data files from the given directory, or directories, are loaded.
Hacker does a reference check on the files in the given directories to determine which release they relate to.

The second directory is optional, in case a HD-only release is to be loaded.

For quicker access, it is recommended to have the load command in a text file, which is then passed with the ```--run``` parameter at startup.

Example:
```
> load "/tmp/dosgames/SystemShock/dosbox/SSHOCK/DATA" "/tmp/dosgames/SystemShock/original_cd/cdrom/data"
Loaded release [DOS CD Release]
> _
```

#### Change Directory
```
cd path
```
When a release was successfully loaded, the data is available in a directory-like hierarchy. The ```cd``` command changes the current directory. The ```path``` parameter can be any combination of ```..```, ```/``` and a directory name. A ```/``` at the beginning refers to the root, otherwise the provided path is relative to the current node.

For example,
```
> cd /hd/mfdart.res/0026
/hd/mfdart.res/0026> cd ../002D
/hd/mfdart.res/002D> _
```
changes to the HD file ```mfdart.res```, chunk 0x0026 ; The second then switches to 0x002D, also in mfdart.res

For nodes containing serial items (such as chunks with blocks or texture properties), the subnodes are simply 0 up to the number of children (-1). Object property nodes are identified by ```class-subclass-type```, e.g.: ```1-3-0```, below which the nodes ```common```, ```generic``` and ```specific``` are available.

File nodes will only be loaded into memory when they are entered. So, although a ```load``` command may succeed, loading specific files may not. Once a file has been loaded into memory, changes outside of Hacker are not reflected (e.g. overwriting a save-game file). Restart the application to load the newest file(s) - this is where the ```--run``` parameter is useful.

#### Node Info
```
info
```
This command prints out some information on the current node. Parent nodes may also give information on which children are available.

Example:
```
/cd/cutspal.res> info
ResourceFile: cutspal.res
IDs: 0003 0004 0005 0006 0007 0008 0009 000A 000B 000C 000D 000E 000F 0010 0011 0012 0013 0014
/cd/cutspal.res> _
```

#### Dump
```
dump
```
For data nodes, this command returns the raw data content as a hexdump.

Example:
```
/cd/textprop.dat/34> dump
0000  00 00 00 00 22 22 0A 00  00 00 00                 ...."".. ...
/cd/textprop.dat/34> _
```

The first column is the offset (in hexadecimal), then up to 16 bytes as 2-digit hex values and the ASCII representation of these bytes on the right (if possible, ```.``` otherwise)

#### Diff
```
diff path
```
This command compares the raw data of the current node against that of another. The other node is referenced with the given path (see the ```cd``` command for a reference on the path).

The result is a dump of both the other node's data (first) and then this node's data. Any difference is highlighted with color.

#### Put
```
put offset bytes...
```
This command, on data nodes, will modify the raw data. The ```offset``` parameter is a hexadecimal number, starting with 0, specifying the offset where to put the following bytes. ```bytes``` is a blank separated list of 2-digit hexadecimal numbers.

The result is a double-dump of data, both the old and the new, with the changes highlighted in color (similar to the ```dump``` command).

This command modifies the data only in memory. To commit the changes to disk, use the ```save``` command.

#### Save
```
save
```
This command iterates through all currently loaded files and saves them to disk. It will rewrite the complete files and returns all names of the files that were saved.

*This command changes your data files! Remember to keep backups!*

## License

The project is available under the terms of the **New BSD License** (see LICENSE file).

[drone-url]: https://drone.io/github.com/inkyblackness/hacker/latest
[drone-image]: https://drone.io/github.com/inkyblackness/hacker/status.png
[coveralls-url]: https://coveralls.io/r/inkyblackness/hacker
[coveralls-image]: https://coveralls.io/repos/inkyblackness/hacker/badge.svg
