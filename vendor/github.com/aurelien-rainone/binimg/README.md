# binimg

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/aurelien-rainone/binimg) [![Build Status](https://travis-ci.org/aurelien-rainone/binimg.svg?branch=master)](https://travis-ci.org/aurelien-rainone/binimg) [![Coverage Status](https://coveralls.io/repos/github/aurelien-rainone/binimg/badge.svg?branch=master)](https://coveralls.io/github/aurelien-rainone/binimg?branch=master)

binimg package proposes an in-memory binary image format, implementing the
image.Image interface, alongside a set of efficient tools to scan rectangular
regions of such images. A binary image has only two possible colors for each
pixel, generally Black and White, though any two colors can be used.

Though the information represented by each pixel could be stored as a single
bit, and thus take a smaller memory footprint, choice has been made to
represent Bit pixels as byte values, that can either be 0 (Black or Off) or 255
(White or On), mostly for simplicity reasons.

Binary images are created either by calling functions such as NewFromImage and
NewBinary, or their counterparts accepting a custom binaryModel.

