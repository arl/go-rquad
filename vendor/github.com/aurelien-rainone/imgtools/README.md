# imgtools

[![Build Status](https://travis-ci.org/aurelien-rainone/imgtools.svg?branch=master)](https://travis-ci.org/aurelien-rainone/imgtools)
[![Coverage Status](https://coveralls.io/repos/github/aurelien-rainone/imgtools/badge.svg?branch=master)](https://coveralls.io/github/aurelien-rainone/imgtools?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/aurelien-rainone/imgtools)](https://goreportcard.com/report/github.com/aurelien-rainone/imgtools)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/aurelien-rainone/imgtools)


Utilities for working with images in Go, completing the standard `image` package.

- [`imgtools/binimg`](./binimg/README.md) : binary image implementation of the `image.Image`
interface, that is an image with two possible values for each pixel.

- [`imgtools/imgscan`](./imgscan/README.md) : fast scanning of rectangular regions of `image.Image`.
