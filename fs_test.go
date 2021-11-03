package io2

import (
	"errors"
	"io/fs"
	"reflect"
	"testing"
)

func TestWriteFile(t *testing.T) {
	want := 1
	writeFileCalled := false
	fsys := &FSDelegator{
		WriteFileFunc: func(name string, p []byte) (int, error) {
			writeFileCalled = true
			return want, nil
		},
	}

	got, err := WriteFile(fsys, "", []byte{})
	if err != nil {
		t.Fatal(err)
	}
	if !writeFileCalled {
		t.Error("Error WriteFile is not called")
	}
	if got != want {
		t.Errorf("Error WriteFile returns %d; want %d", got, want)
	}
}

type openOnlyFsTest struct {
	file fs.File
}

func (fsys *openOnlyFsTest) Open(name string) (fs.File, error) {
	return fsys.file, nil
}

func TestWriteFile_OpenWriteFile(t *testing.T) {
	want := 1
	writeCalled := false
	fsys := &openOnlyFsTest{
		file: &FileDelegator{
			WriteFunc: func(p []byte) (int, error) {
				writeCalled = true
				return want, nil
			},
		},
	}
	got, err := WriteFile(fsys, "", []byte{})
	if err != nil {
		t.Fatal(err)
	}
	if !writeCalled {
		t.Error("Error Write is not called")
	}
	if got != want {
		t.Errorf("Error Write returns %d; want %d", got, want)
	}
}

func TestOpenWriteFile_OpenError(t *testing.T) {
	wantErr := errors.New("test")
	fsys := &FSDelegator{
		OpenFunc: func(name string) (fs.File, error) {
			return nil, wantErr
		},
	}

	var gotErr error
	_, gotErr = OpenWriteFile(fsys, "", []byte{})
	if gotErr == nil {
		t.Errorf("Error WriteFile returns no error")
	}
	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Errorf("Error WriteFile returns unknown error %v; want %v", gotErr, wantErr)
	}
}

type readOnlyFileTest struct {
}

func (f *readOnlyFileTest) Stat() (fs.FileInfo, error) {
	return nil, nil
}

func (f *readOnlyFileTest) Read(p []byte) (int, error) {
	return 0, nil
}

func (f *readOnlyFileTest) Close() error {
	return nil
}

func TestOpenWriteFile_ErrNotImplemented(t *testing.T) {
	fsys := &openOnlyFsTest{
		file: &readOnlyFileTest{},
	}

	name := "test.txt"
	wantErr := &fs.PathError{Op: "Write", Path: name, Err: ErrNotImplemented}

	var err error
	_, err = OpenWriteFile(fsys, name, []byte{})
	if err == nil {
		t.Errorf("Error WriteFile returns no error")
	}
	gotErr, ok := err.(*fs.PathError)
	if !ok {
		t.Errorf("Error WriteFile returns unknown error %v", err)
	}
	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Errorf("Error WriteFile returns unknown error %v; want %v", gotErr, wantErr)
	}
}

func TestRemoveFile(t *testing.T) {
	called := false
	fsys := &FSDelegator{
		RemoveFileFunc: func(name string) error {
			called = true
			return nil
		},
	}

	err := RemoveFile(fsys, "")
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("Error RemoveFile is not called")
	}
}

func TestRemoveAll(t *testing.T) {
	called := false
	fsys := &FSDelegator{
		RemoveAllFunc: func(name string) error {
			called = true
			return nil
		},
	}

	err := RemoveAll(fsys, "")
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("Error RemoveAll is not called")
	}
}

func TestRemoveFile_ErrNotImplemented(t *testing.T) {
	fsys := &OpenFSDelegator{}

	name := "test.txt"
	wantErr := &fs.PathError{Op: "RemoveFile", Path: name, Err: ErrNotImplemented}

	err := RemoveFile(fsys, name)
	if err == nil {
		t.Errorf("Error RemoveFile returns no error")
	}
	gotErr, ok := err.(*fs.PathError)
	if !ok {
		t.Errorf("Error RemoveFile returns unknown error %v", err)
	}
	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Errorf("Error RemoveFile returns unknown error %v; want %v", gotErr, wantErr)
	}
}

func TestRemoveAll_ErrNotImplemented(t *testing.T) {
	fsys := &OpenFSDelegator{}

	path := "path/to/dir"
	wantErr := &fs.PathError{Op: "RemoveAll", Path: path, Err: ErrNotImplemented}

	err := RemoveAll(fsys, path)
	if err == nil {
		t.Errorf("Error RemoveAll returns no error")
	}
	gotErr, ok := err.(*fs.PathError)
	if !ok {
		t.Errorf("Error RemoveAll returns unknown error %v", err)
	}
	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Errorf("Error RemoveAll returns unknown error %v; want %v", gotErr, wantErr)
	}
}
