bzip2
---

[GoDoc](http://godoc.org/github.com/larzconwell/bzip2)

Package bzip2 implements access to compress and decompress data in the bzip2 format.

Currently this focuses on the compressor since [compress/bzip2](http://golang.org/pkg/compress/bzip2) doesn't include one.

Hopefully this will be eventually merged into the standard library without any changes on the users part.

### Notes

References used to write the compressor since there's no specification:
- https://en.wikipedia.org/wiki/Bzip2
- https://bzip.org
- https://code.google.com/p/jbzip2
- http://lbzip2.org/

### Install

```
go get github.com/larzconwell/bzip2
```

### License

MIT licensed, see [here](https://raw.github.com/larzconwell/bzip2/master/LICENSE)
