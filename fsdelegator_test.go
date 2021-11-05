package io2

import (
	"errors"
	"io/fs"
	"os"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestOpenFSDelegator_TestFS(t *testing.T) {
	d := DelegateOpenFS(os.DirFS("osfs/testdata"))
	if err := fstest.TestFS(d, "dir0/file01.txt"); err != nil {
		t.Errorf("Error testing/fstest: %+v", err)
	}
}

func TestOpenFSDelegator_ErrNotImplemented(t *testing.T) {
	d := &OpenFSDelegator{}
	var err error
	_, err = d.Open("")
	if !errors.Is(err, ErrNotImplemented) {
		t.Errorf("Error unknown: %v", err)
	}
}

func TestFSDelegator_TestFS(t *testing.T) {
	d := DelegateFS(os.DirFS("osfs/testdata"))
	if err := fstest.TestFS(d, "dir0/file01.txt"); err != nil {
		t.Errorf("Error testing/fstest: %+v", err)
	}
}

func testFSDelegatorErrors(t *testing.T, d *FSDelegator, wantErr error) {
	var err error
	if _, err = d.Open(""); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
	if _, err = d.ReadDir(""); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
	if _, err = d.ReadFile(""); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
	if _, err = d.Glob(""); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
	if _, err = d.Stat(""); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
	if _, err = d.Sub(""); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
	if err = d.MkdirAll("", fs.ModePerm); err != nil {
		t.Errorf("Error unknown: %v", err)
	}
	if _, err = d.CreateFile("", fs.ModePerm); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
	if _, err = d.WriteFile("", []byte{}, fs.ModePerm); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
	if err = d.RemoveFile(""); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
	if err = d.RemoveAll(""); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
}

func TestFSDelegator_ErrNotImplemented(t *testing.T) {
	testFSDelegatorErrors(t, &FSDelegator{}, ErrNotImplemented)
}

func TestFSDelegator(t *testing.T) {
	wantErr := errors.New("test")

	testFSDelegatorErrors(t, &FSDelegator{
		OpenFunc: func(_ string) (fs.File, error) {
			return nil, wantErr
		},
		ReadDirFunc: func(_ string) ([]fs.DirEntry, error) {
			return nil, wantErr
		},
		ReadFileFunc: func(_ string) ([]byte, error) {
			return nil, wantErr
		},
		GlobFunc: func(_ string) ([]string, error) {
			return nil, wantErr
		},
		StatFunc: func(_ string) (fs.FileInfo, error) {
			return nil, wantErr
		},
		SubFunc: func(_ string) (fs.FS, error) {
			return nil, wantErr
		},
		MkdirAllFunc: func(_ string, _ fs.FileMode) error {
			return nil
		},
		CreateFileFunc: func(_ string, _ fs.FileMode) (WriterFile, error) {
			return nil, wantErr
		},
		WriteFileFunc: func(_ string, _ []byte, _ fs.FileMode) (int, error) {
			return 0, wantErr
		},
		RemoveFileFunc: func(_ string) error {
			return wantErr
		},
		RemoveAllFunc: func(_ string) error {
			return wantErr
		},
	}, wantErr)
}

func TestDelegateFS(t *testing.T) {
	DelegateFS(&FSDelegator{})
}

func TestDelegateFS_ReadDir(t *testing.T) {
	fsys := os.DirFS("osfs/testdata")
	path := "dir0"
	want, err := fs.ReadDir(fsys, path)
	if err != nil {
		t.Fatal(err)
	}
	d := DelegateFS(fsys)
	got, err := d.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error ReadDir returns %v; want %v", got, want)
	}
}

func TestDelegateFS_ReadFile(t *testing.T) {
	fsys := os.DirFS("osfs/testdata")
	path := "dir0/file01.txt"
	want, err := fs.ReadFile(fsys, path)
	if err != nil {
		t.Fatal(err)
	}
	d := DelegateFS(fsys)
	got, err := d.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error ReadFile returns %v; want %v", got, want)
	}
}

func TestDelegateFS_Glob(t *testing.T) {
	fsys := os.DirFS("osfs/testdata")
	pattern := "dir0/*.txt"
	want, err := fs.Glob(fsys, pattern)
	if err != nil {
		t.Fatal(err)
	}
	d := DelegateFS(fsys)
	got, err := d.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error ReadFile returns %v; want %v", got, want)
	}
}

func TestDelegateFS_Stat(t *testing.T) {
	fsys := os.DirFS("osfs/testdata")
	path := "dir0/file01.txt"
	want, err := fs.Stat(fsys, path)
	if err != nil {
		t.Fatal(err)
	}
	d := DelegateFS(fsys)
	got, err := d.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error Stat returns %v; want %v", got, want)
	}
}

func TestDelegateFS_Sub(t *testing.T) {
	fsys := os.DirFS("osfs/testdata")
	path := "dir0"
	want, err := fs.Sub(fsys, path)
	if err != nil {
		t.Fatal(err)
	}
	d := DelegateFS(fsys)
	got, err := d.Sub(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error Sub returns %v; want %v", got, want)
	}
}

func TestDelegateFile(t *testing.T) {
	DelegateFile(&FileDelegator{})
}

func testFileDelegatorErrors(t *testing.T, d *FileDelegator, wantErr error) {
	var err error
	if _, err = d.Stat(); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
	if _, err = d.Read([]byte{}); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
	if err = d.Close(); err != nil {
		t.Errorf("Error unknown: %v", err)
	}
	if _, err = d.ReadDir(-1); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
	if _, err = d.Write([]byte{}); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
}

func TestFileDelegator_ErrNotImplemented(t *testing.T) {
	testFileDelegatorErrors(t, &FileDelegator{}, ErrNotImplemented)
}

func TestFileDelegator(t *testing.T) {
	wantErr := errors.New("test")

	testFileDelegatorErrors(t, &FileDelegator{
		StatFunc: func() (fs.FileInfo, error) {
			return nil, wantErr
		},
		ReadFunc: func(p []byte) (int, error) {
			return 0, wantErr
		},
		CloseFunc: func() error {
			return nil
		},
		ReadDirFunc: func(n int) ([]fs.DirEntry, error) {
			return nil, wantErr
		},
		WriteFunc: func(p []byte) (int, error) {
			return 0, wantErr
		},
	}, wantErr)
}
