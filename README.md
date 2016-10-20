# go-rquas: region quadtrees in Go
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/aurelien-rainone/go-rquad) [![Build Status](https://travis-ci.org/aurelien-rainone/go-rquad.svg?branch=master)](https://travis-ci.org/aurelien-rainone/go-rquad) [![Coverage Status](https://coveralls.io/repos/github/aurelien-rainone/go-rquad/badge.svg?branch=master)](https://coveralls.io/github/aurelien-rainone/go-rquad?branch=master)

**Quadtrees and efficient neighbour finding techniques in Go**

Package `rquad` proposes various implementations of **region quadtrees**.
The region quadtree is a special kind of quadtree that recursively
subdivides a 2D dimensional space into 4, smaller and generally equal
rectangular regions, until the wanted quadtree resolution has been reached,
or no further subdivisions can be performed.

Region quadtree may be used for image processing, in this case a node
represents a rectangular region of an image in which all pixels have the
same color.

A region quadtree may also be used as a variable resolution representation
of a data field. For example, the temperatures in an area may be stored as a
quadtree, with each leaf node storing the average temperature over the
subregion it represents.

Quadtree implementations in this package use the `binimg.Scanner` interface to
represent the complete area and provide the quadtree with a way to scan over
regions of this area in order to perform the subdivisions.


## References

 - Bottom-up neighour finding techniques  
(cf Hanan Samet 1981, *Neighbor Finding Techniques for Images Represented by
Quadtrees*)

 - Cardinal Neighbor Quadtree  
(cf Safwan Qasem 2015, *Cardinal Neighbor Quadtree: a New Quadtree-based
Structure for Constant-Time Neighbor Finding*)


## License

go-rquad is open source software distributed in accordance with the MIT
License, which says:

Copyright (c) 2016 Aurélien Rainone

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
