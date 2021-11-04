package osfs

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/jarxorg/io2"
)

func TestDirFS_TestFS(t *testing.T) {
	if err := fstest.TestFS(DirFS("testdata"), "dir0"); err != nil {
		t.Errorf("Error testing/fstest: %+v", err)
	}
	if err := fstest.TestFS(DirFS("testdata"), "dir0/file01.txt"); err != nil {
		t.Errorf("Error testing/fstest: %+v", err)
	}
}

func TestCreateFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fsys := DirFS(tmpDir)
	got, err := io2.CreateFile(fsys, "test.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer got.Close()
}

func TestCreateFile_MkdirAllError(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	orgMkdirAllFunc := osMkdirAllFunc
	defer func() { osMkdirAllFunc = orgMkdirAllFunc }()

	wantErr := errors.New("test")
	osMkdirAllFunc = func(dir string, perm os.FileMode) error {
		return wantErr
	}

	fsys := DirFS(tmpDir)
	var gotErr error
	_, gotErr = io2.CreateFile(fsys, "name.txt")

	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Errorf("Error CreateFile returns unknown error %v; want %v", gotErr, wantErr)
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	name := "test.txt"
	want := []byte(`test`)

	fsys := DirFS(tmpDir)
	n, err := io2.WriteFile(fsys, name, want)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(want) {
		t.Errorf("Error len %d; want %d", n, len(want))
	}

	got, err := ioutil.ReadFile(tmpDir + "/" + name)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error content %s; want %s", got, want)
	}
}

func TestWriteFile_InvalidError(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fsys := DirFS(tmpDir)
	_, err = io2.WriteFile(fsys, "../invalid.txt", []byte{})
	if err == nil {
		t.Fatal("Error WriteFile returns no error")
	}
}

func TestContainsDenyWin(t *testing.T) {
	testCases := []struct {
		name string
		want bool
	}{
		{
			name: `allow.txt`,
			want: false,
		}, {
			name: `path/to/allow.txt`,
			want: false,
		}, {
			name: `deny:txt`,
			want: true,
		}, {
			name: `C:/deny.txt`,
			want: true,
		}, {
			name: `path\to\deny.txt`,
			want: true,
		},
	}
	for i, testCase := range testCases {
		got := containsDenyWin(testCase.name)
		if got != testCase.want {
			t.Errorf("Error[%d] containsDenyWin(%s) %v; want %v",
				i, testCase.name, got, testCase.want)
		}
	}
}

func TestSub_WriteFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dir := "sub"
	name := "test.txt"
	want := []byte(`test`)

	fsys, err := fs.Sub(DirFS(tmpDir), dir)
	if err != nil {
		t.Fatal(err)
	}
	n, err := io2.WriteFile(fsys, name, want)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(want) {
		t.Errorf("Error len %d; want %d", n, len(want))
	}

	got, err := ioutil.ReadFile(tmpDir + "/" + dir + "/" + name)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error content %s; want %s", got, want)
	}
}

func TestRemoveFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fsys := DirFS(tmpDir)
	name := "test.txt"

	if err = ioutil.WriteFile(tmpDir+"/"+name, []byte{}, defaultFileMode); err != nil {
		t.Fatal(err)
	}

	err = io2.RemoveFile(fsys, name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveFile_InvalidError(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fsys := DirFS(tmpDir)
	err = io2.RemoveFile(fsys, "../invalid-dir")
	if err == nil {
		t.Fatal("Error RemoveFile returns no error")
	}
}

func TestRemoveAll(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fsys := DirFS(tmpDir)
	path := "dir"
	name := "test.txt"

	if err = os.Mkdir(tmpDir+"/"+path, defaultFileMode); err != nil {
		t.Fatal(err)
	}
	if err = ioutil.WriteFile(tmpDir+"/"+path+"/"+name, []byte{}, defaultFileMode); err != nil {
		t.Fatal(err)
	}

	err = io2.RemoveAll(fsys, path)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveAll_InvalidError(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fsys := DirFS(tmpDir)
	err = io2.RemoveAll(fsys, "../invalid-dir")
	if err == nil {
		t.Fatal("Error RemoveAll returns no error")
	}
}
