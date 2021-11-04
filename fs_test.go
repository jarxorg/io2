package io2

import (
	"io/fs"
	"reflect"
	"testing"
)

func TestCreateFile(t *testing.T) {
	want := &FileDelegator{}
	called := false
	fsys := &FSDelegator{
		CreateFileFunc: func(name string) (WriterFile, error) {
			called = true
			return want, nil
		},
	}

	got, err := CreateFile(fsys, "test.txt")
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("Error CreateFile is not called")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error CreateFile returns %v; want %v", got, want)
	}
}

func TestCreateFile_ErrNotImplemented(t *testing.T) {
	fsys := &OpenFSDelegator{}

	name := "test.txt"
	wantErr := &fs.PathError{Op: "CreateFile", Path: name, Err: ErrNotImplemented}

	var err error
	_, err = CreateFile(fsys, name)
	if err == nil {
		t.Errorf("Error CreateFile returns no error")
	}
	gotErr, ok := err.(*fs.PathError)
	if !ok {
		t.Errorf("Error CreateFile returns unknown error %v", err)
	}
	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Errorf("Error CreateFile returns unknown error %v; want %v", gotErr, wantErr)
	}
}

func TestWriteFile(t *testing.T) {
	want := 1
	called := false
	fsys := &FSDelegator{
		WriteFileFunc: func(name string, p []byte) (int, error) {
			called = true
			return want, nil
		},
	}

	got, err := WriteFile(fsys, "", []byte{})
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("Error WriteFile is not called")
	}
	if got != want {
		t.Errorf("Error WriteFile returns %d; want %d", got, want)
	}
}

func TestWriteFile_ErrNotImplemented(t *testing.T) {
	fsys := &OpenFSDelegator{}

	name := "test.txt"
	wantErr := &fs.PathError{Op: "WriteFile", Path: name, Err: ErrNotImplemented}

	var err error
	_, err = WriteFile(fsys, name, []byte{})
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
