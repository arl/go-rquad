[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/aurelien-rainone/imgtools/binimg)

# binimg - Binary images in Go


`binimg` package proposes an in-memory *binary image format*, that is an image
that has only two possible values for each pixel. In this package, we refer to
those 2 colors as `binimg.On` and `binimg.Off`, that respectively convert (via
`binimg.Model`) to the standard Go colors `color.White` and `color.Black`.
`binimg.Model` implements the `color.Model` interface.

Such images are also referred to as *bi-level*, or *two-level*.

`binimg.Binary` implements the standard Go `image.Image` and `draw.Image`.

A pixel could be stored as a single bit, but as the main goal of this package
is fast manipulation of binary images, `binimg.Bit`, the underlying pixel data
type manipulated by `binimg.Binary` image, is 1 `byte` wide.

`Binary` are instantiated by the following functions:

```go
func New(r image.Rectangle, p Palette) *Binary
func NewFromImage(src image.Image, p Palette) *Binary
```

-----------------------
TODO: REMOVE (there is no more palette, but show an exemple of conversion from a color image to a binary one, with standard color model and maybe how to use a custom color model by subclassing/composition and providing a new color model... if that seems useful)
 remove all references to binimg.Black and binimg.White

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
