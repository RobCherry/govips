### Vips for Go

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg)](http://godoc.org/github.com/RobCherry/govips)
[![Report Card](https://goreportcard.com/badge/github.com/RobCherry/govips)](https://goreportcard.com/report/github.com/RobCherry/govips)

This package is powered by the [libvips](https://github.com/jcupitt/libvips) image processing library, originally
 created in 1989 at Birkbeck College and currently maintained by [John Cupitt](https://github.com/jcupitt).

## Prerequisites

* [libvips](https://github.com/jcupitt/libvips) v8.5.0+

## Installation

```
go get github.com/RobCherry/govips
```

### Install libvips on Mac OS

```
brew install homebrew/science/vips --with-imagemagick --with-webp
```

### Install libvips on Linux

TODO

## Usage

In your own code:

```go
import "github.com/RobCherry/govips"

...
govips.Initialize();
...
```

From the command line (`go install github.com/RobCherry/govips/cli`):

```
cli -r 300x300 -q 90 -fast-resize -v path/to/input.jpg path/to/output.jpg
```

## Roadmap

- [ ] Documentation
- [ ] Tests

## Author

[Rob Cherry](https://github.com/RobCherry)

## Contributing

Contributions welcome! Please fork the repository and open a pull request with your changes.

## Notes

The provided sRGB ICC profile is from [icc-profiles-free](https://packages.debian.org/sid/all/icc-profiles-free/filelist)

The provided CMYK ICM profile is from [Argyll Color Management System](http://www.argyllcms.com/cmyk.icm)
