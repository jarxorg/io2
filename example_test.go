package io2_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/jarxorg/io2"
)

func Example_DelegateReader_Error() {
	org := bytes.NewReader([]byte(`original`))

	r := io2.DelegateReader(org)
	r.Delegate.Read = func(p []byte) (int, error) {
		return 0, errors.New("custom")
	}

	var err error
	_, err = ioutil.ReadAll(r)
	fmt.Printf("Error: %v\n", err)

	// Output: Error: custom
}

func Example_WriteSeeker() {
	o := io2.NewWriteSeekBuffer(0)
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
