# github.com/jarxorg/io2

[![PkgGoDev](https://pkg.go.dev/badge/github.com/jarxorg/io2)](https://pkg.go.dev/github.com/jarxorg/io2)
[![Report Card](https://goreportcard.com/badge/github.com/jarxorg/io2)](https://goreportcard.com/report/github.com/jarxorg/io2)

Go "io" package utilities.

## Delegator

Delegator implements io.Reader, io.Writer, io.Seeker, io.Closer.
Delegator can override the io functions that is useful for unit tests.

```go
package main

import (
  "bytes"
  "errors"
  "fmt"
  "io/ioutil"

  "github.com/jarxorg/io2"
)

func main() {
  org := bytes.NewReader([]byte(`original`))

  r := io2.DelegateReader(org)
  r.ReadFunc = func(p []byte) (int, error) {
    return 0, errors.New("custom")
  }

  var err error
  _, err = ioutil.ReadAll(r)
  fmt.Printf("Error: %v\n", err)

  // Output: Error: custom
}
```

### No-op Closer using Delegator

```go
// NopReadWriteCloser returns a ReadWriteCloser with a no-op Close method wrapping the provided interface.
func NopReadWriteCloser(rw io.ReadWriter) io.ReadWriteCloser {
  return DelegateReadWriter(rw)
}

// NopReadSeekCloser returns a ReadSeekCloser with a no-op Close method wrapping the provided interface.
func NopReadSeekCloser(r io.ReadSeeker) io.ReadSeekCloser {
  return DelegateReadSeeker(r)
}

// NopWriteCloser returns a WriteCloser with a no-op Close method wrapping the provided interface.
func NopWriteCloser(w io.Writer) io.WriteCloser {
  return DelegateWriter(w)
}
```

## FSDelegator and FileDelegator

FSDelegator implements FS, ReadDirFS, ReadFileFS, StatFS, SubFS of [io/fs](https://github.com/golang/go/tree/master/src/io/fs) package.
FSDelegator can override the FS functions that is useful for unit tests.

```go
package main

import (
  "errors"
  "fmt"
  "io/fs"
  "os"

  "github.com/jarxorg/io2"
)

func main() {
  fsys := io2.DelegateFS(os.DirFS("."))
  fsys.ReadDirFunc = func(name string) ([]fs.DirEntry, error) {
    return nil, errors.New("custom")
  }

  var err error
  _, err = fs.ReadDir(fsys, ".")
  fmt.Printf("Error: %v\n", err)

  // Output: Error: custom
}
```

## WriteSeekBuffer

WriteSeekBuffer implements io.Writer, io.Seeker and io.Closer.
NewWriteSeekBuffer(capacity int) returns the buffer.

```go
// WriteSeekCloser is the interface that groups the basic Write, Seek and Close methods.
type WriteSeekCloser interface {
  io.Writer
  io.Seeker
  io.Closer
}
```

```go
package main

import (
  "fmt"
  "io"

  "github.com/jarxorg/io2"
)

func main() {
  o := io2.NewWriteSeekBuffer(16)
  o.Write([]byte(`Hello!`))
  o.Truncate(o.Len() - 1)
  o.Write([]byte(` world!`))

  fmt.Println(string(o.Bytes()))

  o.Seek(-1, io.SeekEnd)
  o.Write([]byte(`?`))

  fmt.Println(string(o.Bytes()))

  // Output:
  // Hello world!
  // Hello world?
}
```
