# InkyBlackness Construct

This is a tool as part of the [InkyBlackness](https://inkyblackness.github.io) project, written in [Go](http://golang.org/). It is meant to create a minimally valid level, with which further tests can be performed.

## Usage

### Command Line Interface

```
Usage:
  construct [--file=<file-name>] [--solid]
  construct -h | --help
  construct --version

Options:
  --file=<file-name>  specifies the target file name. [default: archive.dat]
  --solid             Creates an entirely solid map; Exception: Starting tile on level 1.
  -h --help           Show this screen.
  --version           Show version.
```

## License

The project is available under the terms of the **New BSD License** (see LICENSE file).
