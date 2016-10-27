# Gosaic

Create your own image mosaics.

## Installation

Download the latest release from the github releases page. Unzip and put
the binary somewhere on your PATH.

## Usage

### Overview

There are two steps to creating an image mosaic with gosaic. First, add images to
an index. These are the small images that will make up the tiles in the mosaic. Second,
build the mosaic from an image using the images added to the index.

Gosaic will populate a local sqlite3 database with metadata used to build the photo mosaic.
It will not modify any images added to the index, and it will not modify the aspect ratio
of any of the partial index images rendered in the mosaic. The final mosaic jpeg
image is created at 300 dpi.

The best mosaics will be created from large indexes with a wide variety of high quality images.
However, large indexes will increase the amount of time required to generate mosaics.
It's usually best to start small and do some experimentation and work your way up to a
large mosaic. Gosaic does its best to calculate sane defaults for all values, so in most
cases the most basic command will result in a good result.

### 1: Index

Use the `index` sub-command to add images to an index. For example:

```shell
$ gosaic index path/to/image/directory
$ gosaic index path/to/file1.jpg path/to/file2.jpg
```

Or, if you have a file with a list of files and/or directories that you want to
add to the index, you can do something like:

```shell
$ gosaic index < path/to/list.txt
```

`index` sub-command help:

```shell
$ gosaic index -h
Manage index images

Usage:
  gosaic index [PATHS...] [flags]

Flags:
  -c, --clean   Clean the index
  -l, --list    List the index
  -r, --rm      Remove entries from the index

Global Flags:
      --db string     Path to project database (default "$HOME/.gosaic.sqlite3")
      --workers int   Number of workers to use (default 8)
```

### 2: Build Mosaic

#### Aspect

`mosaic aspect` sub-command help:

```
$ gosaic mosaic aspect -h
Create an aspect mosaic from image at PATH

Usage:
  gosaic mosaic aspect PATH [flags]

Flags:
  -a, --aspect string      Aspect of mosaic partials (CxR)
      --cleanup            Delete mosaic metadata after completion
      --cover-out string   File to write cover partial pattern image
  -d, --destructive        Delete mosaic metadata during creation
  -f, --fill-type string   Mosaic fill to use, either 'random' or 'best' (default "random")
      --height int         Pixel height of mosaic, 0 maintains aspect from width
      --macro-out string   File to write resized macro image
      --max-repeats int    Number of times an index image can be repeated, 0 is unlimited, -1 is the minimun number (default -1)
  -n, --name string        Name of mosaic
      --out string         File to write final mosaic image
  -s, --size int           Number of mosaic partials in smallest dimension, 0 auto-calculates
  -t, --threashold float   How similar aspect ratios must be (default -1)
  -w, --width int          Pixel width of mosaic, 0 maintains aspect from image height

Global Flags:
      --db string     Path to project database (default "$HOME/.gosaic.sqlite3")
      --workers int   Number of workers to use (default 8)
```

#### Quad

`mosaic quad` sub-command help:

```
$ gosaic mosaic quad -h
Create quad-tree mosaic from image at PATH

Usage:
  gosaic mosaic quad PATH [flags]

Flags:
      --cleanup            Delete mosaic metadata after completion
      --cover-out string   File to write cover partial pattern image
  -d, --destructive        Delete mosaic metadata during creation
  -f, --fill-type string   Mosaic fill to use, either 'random' or 'best' (default "random")
      --height int         Pixel height of mosaic, 0 maintains aspect from width
      --macro-out string   File to write resized macro image
      --max-depth int      Number of times a partial can be split into quads (default -1)
      --max-repeats int    Number of times an index image can be repeated, 0 is unlimited, -1 is the minimun number (default -1)
      --min-area int       The smallest an partial can get before it can't be split (default -1)
  -n, --name string        Name of mosaic
  -o, --out string         File to write final mosaic image
  -s, --size int           Number of times to split the partials into quads (default -1)
  -t, --threashold float   How similar aspect ratios must be (default -1)
  -w, --width int          Pixel width of mosaic, 0 maintains aspect from image height

Global Flags:
      --db string     Path to project database (default "$HOME/.gosaic.sqlite3")
      --workers int   Number of workers to use (default 8)
```

## Tips

If you want to maintain multiple indexes of images, possibly with different themes,
it can be helpful to create a bash alias that specifies the global db option for that
theme. For example:

```shell
alias gosaic_wedding="gosaic --db $HOME/.gosaic_wedding.sqlite3"
alias gosaic_africa="gosaic --db $HOME/.gosaic_africa.sqlite3"
alias gosaic_albums="gosaic --db $HOME/.gosaic_albums.sqlite3"
```

Then you can generate themed mosaics, like so:

```shell
$ gosaic_wedding index < path/to/wedding/photos
$ gosaic_wedding mosaic aspect path/to/wedding/photo.jpg
```

## Building

### Requirements

* golang 1.5.x or later

```shell
$ git clone git@github.com/atongen/gosaic.git
$ cd gosaic
$ make
```

## References

* http://tools.medialab.sciences-po.fr/iwanthue/
* http://en.wikipedia.org/wiki/Color_difference
* https://en.wikipedia.org/wiki/Lab_color_space
* https://github.com/fogleman/Quads
* https://www.flickr.com/photos/tsevis/

## Contributing

1. Fork it
1. Create your feature branch (`git checkout -b my-new-feature`)
1. Commit your changes (`git commit -am 'Add some feature'`)
1. Push to the branch (`git push -u origin my-new-feature`)
1. Create new Pull Request
