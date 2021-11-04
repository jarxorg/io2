package osfs_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jarxorg/io2"
	"github.com/jarxorg/io2/osfs"
)

func ExampleDirFS() {
	tmpDir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	name := "example.txt"
	content := []byte(`Hello`)

	fsys := osfs.DirFS(tmpDir)
	_, err = io2.WriteFile(fsys, name, content)
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
