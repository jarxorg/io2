package memfs_test

import (
	"fmt"
	"log"
	"io/fs"

	"github.com/jarxorg/io2"
	"github.com/jarxorg/io2/memfs"
)

func ExampleNew() {
	name := "path/to/example.txt"
	content := []byte(`Hello`)

	fsys := memfs.New()
	var err error
	_, err = io2.WriteFile(fsys, name, content, fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	wrote, err := fs.ReadFile(fsys, name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(wrote))

	// Output: Hello
}
