package io2_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"

	"github.com/jarxorg/io2"
	"github.com/jarxorg/io2/osfs"
)

func ExampleWriteFile() {
	tmpDir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	name := "example.txt"
	content := []byte(`Hello`)

	fsys := osfs.DirFS(tmpDir)
	_, err = io2.WriteFile(fsys, name, content, fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	wrote, err := ioutil.ReadFile(tmpDir + "/" + name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(wrote))

	// Output: Hello
}

func ExampleDelegateReader() {
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

func ExampleDelegateFS() {
	fsys := io2.DelegateFS(os.DirFS("."))
	fsys.ReadDirFunc = func(name string) ([]fs.DirEntry, error) {
		return nil, errors.New("custom")
	}

	var err error
	_, err = fs.ReadDir(fsys, ".")
	fmt.Printf("Error: %v\n", err)

	// Output: Error: custom
}

func ExampleDelegateFile() {
	fsys := io2.DelegateFS(os.DirFS("."))
	fsys.OpenFunc = func(name string) (fs.File, error) {
		return &io2.FileDelegator{
			StatFunc: func() (fs.FileInfo, error) {
				return nil, errors.New("custom")
			},
		}, nil
	}

	file, _ := fsys.Open("anyfile")
	var err error
	_, err = file.Stat()
	fmt.Printf("Error: %v\n", err)

	// Output: Error: custom
}

func ExampleNewWriteSeekerBuffer() {
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
