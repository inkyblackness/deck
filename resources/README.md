# InkyBlackness - Deck

This package project contains the public distribution of the [InkyBlackness](https://inkyblackness.github.io) binaries.

The executables are modding tools for the DOS game [System Shock (1994)](http://en.wikipedia.org/wiki/System_Shock).

## Disclaimer
For all tools, a general warning: Keep backups of your files!
These are executables downloaded from the internet. If you are unsure what they do or how they work, refer to the sources [on GitHub](https://github.com/inkyblackness).


## Package Contents

### InkyBlackness - Shocked (Server & Client)
This tool is the primary application. It is a GUI editor for the resources, including a map editor. The tool is setup as a client/server application, with the client being a browser-page.

#### Running the server
You need the original resource files of the game to start it. Example:

```
shocked-server --source <path to resources> --projects <a projects path> --client client

```

The ```--source``` parameter must point to a resource base. This can be the root directory of the original CD or a similar package. The server will look for the various necessary files below this path.

The ```--projects``` parameter must point to a directory where all the projects will be stored. Each project will have its own subdirectory there.

The ```client``` parameter must point to a directory where the HTML client is located. When started from the distribution package, this is simply the ```client``` subdirectory.

When running, point a modern browser to ```http://localhost:8080/client/index.html``` and the client should load.


A few seconds after something is changed in a project, the affected file(s) are written in the subdirectory of the project. The files in the ```source``` directories are not changed. (A clever setup might allow pointing both directories at the same files, though this is not intended.)
As soon as a resource file exists in a project's directory, that file is used in favour to the original source for future edits.


### InkyBlackness - Construct
This console tool is ideally used in combination with the map editor. It can generate a small ```archive.dat``` file with a single, empty level.


### InkyBlackness - Chunkie
This is a resource import/export tool for the console. Some raw media formats are supported with their respective file formats.


### InkyBlackness - Hacker
This is a low-level file access tool for the console; A better hex-editor if you will. One can modify all bytes of the resources in any way. One can also crash the game with it.
