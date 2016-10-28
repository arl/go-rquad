# imgtools

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/aurelien-rainone/imgtools) [![Build Status](https://travis-ci.org/aurelien-rainone/imgtools.svg?branch=master)](https://travis-ci.org/aurelien-rainone/imgtools) [![Coverage Status](https://coveralls.io/repos/github/aurelien-rainone/imgtools/badge.svg?branch=master)](https://coveralls.io/github/aurelien-rainone/imgtools?branch=master)


`imgtools` package contains some utilities for working with 2D images in Go,
completing the standard Go `image` package.

- [`imgtools/binimg`](./binimg) : binary image implementation of the `image.Image`
interface. that is an image that has only two possible values for each pixel.

- [`imgtools/imgscan`](./imgscan) : fast scanning of rectangular regions of `image.Image`.
