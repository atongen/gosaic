# Gosaic

Create your own image mosaics.

Requirements
------------
* golang 1.6.x or later
* [gb](https://getgb.io/)
* postgres 9.6


Building
--------

```shell
$ go get github.com/constabulary/gb/...
$ git clone git@github.com/atongen/gosaic.git
$ cd gosaic
$ gb build all
```

## Usage

TODO: Write usage instructions here

## Background

* http://www.easyrgb.com/index.php
* http://en.wikipedia.org/wiki/Color_difference
* http://en.wikipedia.org/wiki/Dithering
* http://en.wikipedia.org/wiki/Color_quantization
* http://zejn.net/b/2016/06/10/postgresql-tutorial-color-similarity-search/
* https://raonyguimaraes.com/how-to-install-postgresql-9-6-on-ubuntudebianlinux-mint/

## Libraries

* https://github.com/fogleman/gg
* https://github.com/disintegration/imaging
* https://github.com/rwcarlsen/goexif
* https://github.com/spf13/cobra

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
