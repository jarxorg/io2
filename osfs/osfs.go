// Package osfs provides a filesystem for the OS.
package osfs

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jarxorg/io2"
)

// defaultFileMode is specified by creating directories and files.
const defaultFileMode = 0755

// DirFS returns a filesystem for the tree of files rooted at the directory dir.
// The filesystem can write using io2.WriteFile(fsys fs.FS, name string, p []byte).
func DirFS(dir string) fs.FS {
	return NewOSFS(dir)
}

// containsDenyWin reports whether any path characters of windows are within s.
func containsDenyWin(s string) bool {
	return strings.ContainsAny(s, `\:`)
}

// isInvalidPath reports whether the given path name is valid for use in a call to Create and Write.
func isInvalidPath(name string) bool {
	return !fs.ValidPath(name) || runtime.GOOS == "windows" && containsDenyWin(name)
}

var osCreateFunc = func(name string) (*os.File, error) {
	return os.Create(name)
}

var osMkdirAllFunc = func(dir string, perm os.FileMode) error {
	return os.MkdirAll(dir, perm)
}

var osRemoveFunc = func(name string) error {
	return os.Remove(name)
}

var osRemoveAllFunc = func(path string) error {
	return os.RemoveAll(path)
}

// OSFS represents a filesystem for the OS.
type OSFS struct {
	Dir  string
	osFS *io2.FSDelegator
}

var (
	_ fs.FS            = (*OSFS)(nil)
	_ fs.GlobFS        = (*OSFS)(nil)
	_ fs.ReadDirFS     = (*OSFS)(nil)
	_ fs.ReadFileFS    = (*OSFS)(nil)
	_ fs.StatFS        = (*OSFS)(nil)
	_ fs.SubFS         = (*OSFS)(nil)
	_ io2.WriteFileFS  = (*OSFS)(nil)
	_ io2.RemoveFileFS = (*OSFS)(nil)
)

// NewOSFS returns a filesystem for the tree of files rooted at the directory dir.
func NewOSFS(dir string) *OSFS {
	return &OSFS{
		Dir:  dir,
		osFS: io2.DelegateFS(os.DirFS(dir)),
	}
}

// Open opens the named file.
func (fsys *OSFS) Open(name string) (fs.File, error) {
	return fsys.osFS.Open(name)
}

// Glob returns the names of all files matching pattern, providing an implementation
// of the top-level Glob function.
func (fsys *OSFS) Glob(pattern string) ([]string, error) {
	return fsys.osFS.Glob(pattern)
}

// ReadDir reads the named directory and returns a list of directory entries sorted
// by filename.
func (fsys *OSFS) ReadDir(dir string) ([]fs.DirEntry, error) {
	return fsys.osFS.ReadDir(dir)
}

// ReadFile reads the named file and returns its contents.
func (fsys *OSFS) ReadFile(name string) ([]byte, error) {
	return fsys.osFS.ReadFile(name)
}

// Stat returns a FileInfo describing the file. If there is an error, it should be
// of type *PathError.
func (fsys *OSFS) Stat(name string) (fs.FileInfo, error) {
	return fsys.osFS.Stat(name)
}

// Sub returns an FS corresponding to the subtree rooted at dir.
func (fsys *OSFS) Sub(dir string) (fs.FS, error) {
	return NewOSFS(filepath.Join(fsys.Dir, dir)), nil
}

// CreateFile creates the named file.
func (fsys *OSFS) CreateFile(name string) (io2.WriterFile, error) {
	if isInvalidPath(name) {
		return nil, &os.PathError{Op: "Create", Path: name, Err: os.ErrInvalid}
	}

	path := filepath.Join(fsys.Dir, name)
	err := osMkdirAllFunc(filepath.Dir(path), defaultFileMode)
	if err != nil {
		return nil, err
	}
	return osCreateFunc(path)
}

// WriteFile writes the specified bytes to the named file.
func (fsys *OSFS) WriteFile(name string, p []byte) (int, error) {
	f, err := fsys.CreateFile(name)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return f.Write(p)
}

// RemoveFile removes the specified named file.
func (fsys *OSFS) RemoveFile(name string) error {
	if isInvalidPath(name) {
		return &os.PathError{Op: "remove", Path: name, Err: os.ErrInvalid}
	}
	return osRemoveFunc(filepath.Join(fsys.Dir, name))
}

// RemoveAll removes path and any children it contains.
func (fsys *OSFS) RemoveAll(path string) error {
	if isInvalidPath(path) {
		return &os.PathError{Op: "removeAll", Path: path, Err: os.ErrInvalid}
	}
	return osRemoveAllFunc(filepath.Join(fsys.Dir, path))
}
