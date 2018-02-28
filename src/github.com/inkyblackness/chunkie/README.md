# InkyBlackness Chunkie

This is a tool as part of the [InkyBlackness](https://inkyblackness.github.io) project, written in [Go](http://golang.org/). This tool provides import/export access to resource files for modification of media content.

## Usage

### Command Line Interface

```
Usage:
  chunkie export <resource-file> <chunk-id> [--block=<block-id>] [--raw] [--pal=<palette-file>] [--fps=<framerate>] [<folder>]
  chunkie import <resource-file> <chunk-id> [--block=<block-id>] [--data-type=<id>] <source-file>
  chunkie -h | --help
  chunkie --version

Options:
  <resource-file>       The resource file to work on.
  <chunk-id>            The chunk identifier. Defaults to decimal, use "0x" as prefix for hexadecimal. "all" for all.
  --block=<block-id>    The block identifier. Defaults to decimal, use "0x" as prefix for hexadecimal. "all" for all. [default: 0]
  --raw                 With this flag, the chunk will be exported without conversion to a common file format.
  --pal=<palette-file>  For handling bitmaps & models, use this palette file to write color information
  --fps=<framerate>     The frames per second to emulate when exporting movies. 0 names files after timestamp. [default: 0]
  --data-type=<id>      The type of the chunk to write.
  <folder>              The path of the folder to use. [default: .]
  <source-file>         The source file to import.
  -h --help             Show this screen.
  --version             Show version.
```

The base file name of files is ```XXXX_YYY.ZZZ```. XXXX is the hexadecimal presentation of the chunk number. YYY is decimal for the block number. ZZZ is the type of the file, defaulting to ```bin```.

For exporting, basic formats will be exported as known file types. Specifying --raw will export the chunk in its raw format.
Files are imported raw as well, unless a conversion is known.

The following formats are supported for import and export: .wav for audio, .png for images
The following format is supported for export only: .xml for text strings, .obj (Wavefront) for geometry, .wav/.png/.srt for movies.

### Movie handling
When movies are exported, the optional ```fps``` parameter specifies which framerate to emulate. Videos in the resource files don't follow a strict framerate and frames can't be directly used as stills. If the parameter is 0, the filename will contain the offset in ```sss.fff``` format for seconds and fractions (milliseconds). Any other value will have the export code to duplicate frames to reach the requested framerate. In this case, the filename will contain a 4-digit framenumber.

## License

The project is available under the terms of the **New BSD License** (see LICENSE file).
