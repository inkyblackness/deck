[![Build Status][drone-image]][drone-url]
[![Coverage Status][coveralls-image]][coveralls-url]

# Resource File Access

This is a library as part of the [InkyBlackness](https://inkyblackness.github.io) project, written in [Go](http://golang.org/), to provide basic (binary) access to the resource files of System Shock.

It supports both reading and writing of the files, in a bit-transparent manner.
Reading the extracted (and decompressed) data and then writing it again creates identical files.

## Supported Files & Data Format
The library supports the following resource files:
* chunk files (\*.res / archive.dat)
* objprop.dat (Object properties)
* textprop.dat (Texture properties)

The data format (framing) of the supported files is documented in the [ss-specs](https://github.com/inkyblackness/ss-specs) sub-project of InkyBlackness.

## License

The project is available under the terms of the **New BSD License** (see LICENSE file).

[drone-url]: https://drone.io/github.com/inkyblackness/res/latest
[drone-image]: https://drone.io/github.com/inkyblackness/res/status.png
[coveralls-url]: https://coveralls.io/r/inkyblackness/res
[coveralls-image]: https://coveralls.io/repos/inkyblackness/res/badge.svg
