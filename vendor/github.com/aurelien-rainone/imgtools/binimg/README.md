[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/aurelien-rainone/imgtools/binimg)

# binimg - Binary images in Go


`binimg` package proposes an in-memory *binary image format*, that is an image
that has only two possible values for each pixel. Typically, the two colors
used for a binary image are black and white, though any two colors can be used.
Such images are also referred to as *bi-level*, or *two-level*.

`Binary` implements the standard Go `image.Image` and `draw.Image` interfaces
and embeds a two colors palette of type `Palette`, that itself implements the
`color.Model` interface. `Palette` allows any `color.Color` to be converted to
`OffColor` or `OnColor`.

A pixel could be stored as a single bit, but as the main goal of this package
is fast manipulation of binary images, `Bit`, the underlying pixel data
type manipulated by `Binary` image, is 1 `byte` wide.

`Binary` are instantiated by the following functions:

```go
func New(r image.Rectangle, p Palette) *Binary
func NewFromImage(src image.Image, p Palette) *Binary
```

-----------------------

**`BlackAndWhite` predefined `Palette`**

<img src="https://github.com/aurelien-rainone/imgtools/blob/readme-images/colorgopher.png" width="128">  <img src="https://github.com/aurelien-rainone/imgtools/blob/readme-images/bwgopher.png" width="128">

**`BlackAndWhiteHighThreshold` predefined `Palette`**

<img src="https://github.com/aurelien-rainone/imgtools/blob/readme-images/colorgopher.png" width="128">  <img src="https://github.com/aurelien-rainone/imgtools/blob/readme-images/bwgopher.high.threshold.png" width="128">

**Custom `Palette`**

<img src="https://github.com/aurelien-rainone/imgtools/blob/readme-images/colorgopher.png" width="128">  <img src="https://github.com/aurelien-rainone/imgtools/blob/readme-images/redblue.gopher.png" width="128">

-----------------------

## Usage

- **Create and modify new binary image**

```go
package main

import (
	"image"
	"image/color"

	"github.com/aurelien-rainone/imgtools/binimg"
)

func main() {
	// create a new image (prefilled with OffColor: black)
	bin := binimg.New(image.Rect(0, 0, 128, 128), binimg.BlackAndWhite)

	// set a pixel to OnColor: White
	bin.SetBit(10, 0, binimg.On)

	// set a pixel, converting original color with BlackAndWhite Palette
	bin.Set(10, 0, color.RGBA{127, 23, 98, 255})

	// set rectangular region to White
	bin.SetRect(image.Rect(32, 32, 64, 64), binimg.On)
}
```

- **Convert an existing `image.Image` into black and white binary image**

```go
package main

import "github.com/aurelien-rainone/binimg"

func main() {
	// load image ("color-gopher.png")
	// ...
	bin := binimg.NewFromImage(img)

	// save image ("black&white-gopher.png")
	// ...
}
```

- **Use a custom `binimg.Palette` (i.e `color.Model`)**

```go
import (
	"image"
	"image/color"

	"github.com/aurelien-rainone/imgtools/binimg"
)

func main() {
	var img image.Image

	// ... decode image

	// create
	palette := binimg.Palette{
		OnColor:   color.RGBA{255, 0, 0, 255},
		OffColor:  color.RGBA{0, 0, 255, 255},
		Threshold: 97,
	}
	bin := binimg.NewFromImage(img, palette)

	// ... encode image
}
```
