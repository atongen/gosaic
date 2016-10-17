# Gosaic

Create your own image mosaics.

## Installation

Download the latest release from the github releases page. Unzip and put
the binary somewhere on your PATH.

## Usage

### Overview

### 1: Index

### 2: Build Mosaic

#### Aspect

#### Quad

## Building

### Requirements

* golang 1.6.x or later
* [gb](https://getgb.io/)
* sqlite3

```shell
$ go get github.com/constabulary/gb/...
$ git clone git@github.com/atongen/gosaic.git
$ cd gosaic
$ ./build.sh
```

## TODO

* persist project completion status
* optionally cleanup after project completion (default true)
  - delete cover for project
* resume messaging for complete/incomplete project
* mosaic build option: destroy partial comparisons after mosaic partial created
* exiftool wrapper
* makefile
* cross compile
  - linux amd64
  - darwin amd64
  - windows amd64
* move code to github
* travis ci

## References

* http://www.easyrgb.com/index.php
* http://en.wikipedia.org/wiki/Color_difference
* fogleman quad trees
* lab
* other mosaic examples

## Contributing

1. Fork it
1. Create your feature branch (`git checkout -b my-new-feature`)
1. Commit your changes (`git commit -am 'Add some feature'`)
1. Push to the branch (`git push origin my-new-feature`)
1. Create new Pull Request

## Acknowledgement

This project began as a hack day project at [Leadpages](https://www.leadpages.net).
