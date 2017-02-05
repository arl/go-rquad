# go-rquads: region quadtrees in Go
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/aurelien-rainone/go-rquad) [![Build Status](https://travis-ci.org/aurelien-rainone/go-rquad.svg?branch=master)](https://travis-ci.org/aurelien-rainone/go-rquad) [![Coverage Status](https://coveralls.io/repos/github/aurelien-rainone/go-rquad/badge.svg?branch=master)](https://coveralls.io/github/aurelien-rainone/go-rquad?branch=master)

**Quadtrees and efficient neighbour finding techniques in Go**

Package `rquad` proposes various implementations of **region quadtrees**.
The region quadtree is a special kind of quadtree that recursively
subdivides a 2 dimensional space into 4 smaller and generally equal
rectangular regions, until the wanted quadtree resolution has been reached,
or no further subdivisions can be performed.

Region quadtree may be used for image processing, in this case a node
represents a rectangular region of an image in which all pixels have the
same color.

A region quadtree may also be used as a variable resolution representation
of a data field. For example, the temperatures in an area may be stored as a
quadtree, with each leaf node storing the average temperature over the
subregion it represents.

Quadtree implementations in this package use the [`imgscan.Scanner`](https://github.com/aurelien-rainone/imgtools/tree/master/imgscan)
interface to represent the complete area and provide the quadtree with a way 
to scan over regions of this area in order to perform the subdivisions.

## API Overview

### `Node` interface
```go
type Node interface {
        Parent() Node
        Child(Quadrant) Node
        Bounds() image.Rectangle
        Color() Color
        Location() Quadrant
}
```

### `Quadtree` interface

A `Quadtree` being a hierarchical collection of `Node`s, its API is relatively
simple and gives an access to the tree root, and a way to iterate over all
the leaves.

```go
type Quadtree interface {
        ForEachLeaf(Color, func(Node))
        Root() Node
}
```

### Functions

`Locate` returns the leaf node of `q` that contains `pt`, or nil if `q` doesn't contain `pt`.
```go
func Locate(q Quadtree, pt image.Point) Node
```

`ForEachNeighbour` calls `fn` for each neighbour of `n`.
```go
func ForEachNeighbour(n Node, fn func(Node))
```

### Basic implementation: `BasicTree` and `basicNode`

`BasicTree` is in many ways the standard implementation of `Quadtree`, it just does the job.

### State of the art implementation: `CNTree` and `CNNode`

`CNTree` or **Cardinal Neighbour Quadtree** implements state of the art techniques:
 - node neighbours (of any size) are accessed in constant time *0(1)* thanks to the implementation of *Cardinal Neighbour Quadtree* technique (cf Safwan Qasem 2015). The time complexity reduction is obtained through the addition of only four pointers per leaf node in the quadtree.
 - fast point location queries (locating which leaf node contains a specific point), thanks to the *binary branching method* (cf Frisken Perry 2002). This simple and efficient method is nonrecursive, table free, and reduces the number of comparisons with
poor predictive behavior, that are otherwise required with the standard method.

## Benchmarks

![Quadtree creation benchmark](https://raw.githubusercontent.com/aurelien-rainone/go-rquad/readme-docs/Creation.png)

![Neighbour finding benchmark](https://raw.githubusercontent.com/aurelien-rainone/go-rquad/readme-docs/Neighbours.png)

![Point location benchmark](https://raw.githubusercontent.com/aurelien-rainone/go-rquad/readme-docs/PointLocation.png)

## References

 - Bottom-up neighour finding techniques  
(cf Hanan Samet 1981, *Neighbor Finding Techniques for Images Represented by
Quadtrees*)

 - Cardinal Neighbor Quadtree  
(cf Safwan Qasem 2015, *Cardinal Neighbor Quadtree: a New Quadtree-based
Structure for Constant-Time Neighbor Finding*)

 - Fast point location using binary branching method  
(cf Frisken, Perry 2002, *Simple and Efficient Traversal Methods for Quadtrees
and Octrees*)


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


