# InkyBlackness - Deck

This package project contains the public distribution of the [InkyBlackness](https://inkyblackness.github.io) binaries.

The executables are modding tools for the DOS game [System Shock (1994)](http://en.wikipedia.org/wiki/System_Shock).

## Disclaimer
For all tools, a general warning: Keep backups of your files!
These are executables downloaded from the internet. If you are unsure what they do or how they work, refer to the sources [on GitHub](https://github.com/inkyblackness).


## Package Contents

### InkyBlackness - Shocked (Server & Client)
This tool is the primary application. It is a GUI editor for the resources, including a map editor. The tool is setup as a client/server application, with the client available as a browser-page and as a native application.
The browser client uses WebGL while the native uses OpenGL 3.2 .

#### Running the server
There are two ways to run the server:
* Project based, where original files are used as a reference and modifications are stored in dedicated folders.
* In-Place, where any modification is directly applied to the provided files.

Either way, you need the original resource files of the game to start it.

##### Project based
Example:

```
shocked-server project --source <path to resources> --projects <a projects path> --client client

```

The ```--source``` parameter must point to a resource base. This can be the root directory of the original CD or a similar package. The server will look for the various necessary files below this path.

The ```--projects``` parameter must point to a directory where all the projects will be stored. Each project will have its own subdirectory there.

If you want to use the browser-client, then the ```client``` parameter must point to a directory where the HTML client is located. When started from the distribution package, this is simply the ```client``` subdirectory.
When running, point a modern browser to ```http://localhost:8080/client/index.html``` and the client should load.


A few seconds after something is changed in a project, the affected file(s) are written in the subdirectory of the project. The files in the ```source``` directories are not changed. (A clever setup might allow pointing both directories at the same files, though this is not intended.)
As soon as a resource file exists in a project's directory, that file is used in favour to the original source for future edits.

##### In-Place
Example:

```
shocked-server inplace --path <path to resources> --client client

```

The ```--path``` parameter can be provided several times and point to one directory with resources files. If all resources are in the same directory, only one is enough. Two might be used for a vanillay System Shock installation, with one directory pointing to the resources on the CD-ROM and the other to the files on harddisk.
Note that the provided path(s) must point to an actual data directory, within which the .res files are stored.

The ```--client``` parameter behaves as described above.

##### Extra parameters
The server, by default, listens on ```localhost:8080```. It is possible to specify a different address with the ```--address``` parameter.

#### Native client application

Instead of using the browser, it is also possible to use a native application: ```shocked-client-console```. This also connects to the server (using HTTP) and presents a similar interface as the browser. Due to technical limitations, the control interface is text-based.

Should the server listen on a different address than default, provide a ```--address``` parameter as well.


### InkyBlackness - Construct
This console tool is ideally used in combination with the map editor. It can generate a small ```archive.dat``` file with a single, empty level.


### InkyBlackness - Chunkie
This is a resource import/export tool for the console. Some raw media formats are supported with their respective file formats.


### InkyBlackness - Hacker
This is a low-level file access tool for the console; A better hex-editor if you will. One can modify all bytes of the resources in any way. One can also crash the game with it.
