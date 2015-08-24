bzip2
---

[GoDoc](http://godoc.org/github.com/larzconwell/bzip2)

Package bzip2 implements access to compress and decompress data in the bzip2 format.

Currently this focuses on the compressor since [compress/bzip2](http://golang.org/pkg/compress/bzip2) doesn't include one.

Hopefully this will be eventually merged into the standard library without any changes on the users part.

### Notes

References used to write the compressor since there's no specification:
- https://en.wikipedia.org/wiki/Bzip2
- https://github.com/cscott/compressjs
- https://code.google.com/p/jbzip2/source/browse/trunk/jbzip2/src/org/itadaki/bzip2
- http://lbzip2.org/

### Todo

- [ ] Block checksum wrong with bigger inputs
- [ ] Repeat count too large in Go reader with a full block during MTF/RLE2
- [ ] BWT index out of bounds with full block with random bytes
- [ ] Tons of tests writing files and reading with both Go reader, and bzip binary
- [ ] Write go-fuzz tests
- [ ] Optimize BWT, using LSD Radix sorting
- [ ] Tree selections should go from 0-5 to 5-0 and alternate, makes smaller MTF outputs
- [ ] Research ways to make smaller Huffman codes(http://fastcompression.blogspot.fr/2015/07/huffman-revisited-part-1.html)

### Install

```
go get github.com/larzconwell/bzip2
```

### License

MIT licensed, see [here](https://raw.github.com/larzconwell/bzip2/master/LICENSE)
