package io2

import (
	"io"
	"io/fs"
)

// WriterFile is a file that provides an implementation fs.File and io.Writer.
type WriterFile interface {
	fs.File
	io.Writer
}

// WriteFileFS is the interface implemented by a filesystem that provides an
// optimized implementation of WriteFile.
type WriteFileFS interface {
	fs.FS
	WriteFile(name string, p []byte) (n int, err error)
}

// WriteFile writes the specified bytes to the named file. If the filesystem implements
// WriteFileFS calls fsys.WriteFile otherwise calls OpenWriteFile.
func WriteFile(fsys fs.FS, name string, p []byte) (n int, err error) {
	if fsys, ok := fsys.(WriteFileFS); ok {
		return fsys.WriteFile(name, p)
	}
	return OpenWriteFile(fsys, name, p)
}

// OpenWriteFile opens the named file. If the file implements WriterFile calls Write
// otherwise returns a PathError.
func OpenWriteFile(fsys fs.FS, name string, p []byte) (n int, err error) {
	file, err := fsys.Open(name)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	w, ok := file.(WriterFile)
	if !ok {
		return 0, &fs.PathError{Op: "Write", Path: name, Err: ErrNotImplemented}
	}
	return w.Write(p)
}

// RemoveFileFS is the interface implemented by a filesystem that provides an
// implementation of RemoveFile.
type RemoveFileFS interface {
	fs.FS
	RemoveFile(name string) error
	RemoveAll(name string) error
}

// RemoveFile removes the specified named file. If the filesystem implements
// RemoveFileFS calls fsys.RemoveFile otherwise return a PathError.
func RemoveFile(fsys fs.FS, name string) error {
	if fsys, ok := fsys.(RemoveFileFS); ok {
		return fsys.RemoveFile(name)
	}
	return &fs.PathError{Op: "RemoveFile", Path: name, Err: ErrNotImplemented}
}

// RemoveAll removes path and any children it contains. If the filesystem
// implements RemoveFileFS calls fsys.RemoveAll otherwise return a PathError.
func RemoveAll(fsys fs.FS, path string) error {
	if fsys, ok := fsys.(RemoveFileFS); ok {
		return fsys.RemoveAll(path)
	}
	return &fs.PathError{Op: "RemoveAll", Path: path, Err: ErrNotImplemented}
}
