// Package osfs provides a filesystem for the OS.
package osfs

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jarxorg/io2"
)

const defaultFileMode = 0755

// DirFS returns a filesystem for the tree of files rooted at the directory dir.
// The filesystem can write using io2.WriteFile(fsys fs.FS, name string, p []byte).
func DirFS(dir string) fs.FS {
	return DirFSDelegator(dir)
}

// DirFSDelegator returns a io2.FSDelegator delegate the functions of files rooted at the directory dir.
func DirFSDelegator(dir string) *io2.FSDelegator {
	dirFsys := os.DirFS(dir)
	fsys := io2.DelegateFS(dirFsys)
	fsys.SubFunc = subFunc(dirFsys, dir)
	fsys.WriteFileFunc = writeFileFunc(dir)
	fsys.RemoveFileFunc = removeFileFunc(dir)
	fsys.RemoveAllFunc = removeAllFunc(dir)
	return fsys
}

// NOTE: copy from os package.
func containsAny(s, chars string) bool {
	for i := 0; i < len(s); i++ {
		for j := 0; j < len(chars); j++ {
			if s[i] == chars[j] {
				return true
			}
		}
	}
	return false
}

func containsDenyWin(name string) bool {
	return containsAny(name, `\:`)
}

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

func subFunc(fsys fs.FS, dir string) func(dir string) (fs.FS, error) {
	return func(name string) (fs.FS, error) {
		return DirFS(filepath.Join(dir, name)), nil
	}
}

func writeFileFunc(dir string) func(name string, p []byte) (int, error) {
	return func(name string, p []byte) (int, error) {
		if isInvalidPath(name) {
			return 0, &os.PathError{Op: "create", Path: name, Err: os.ErrInvalid}
		}

		path := filepath.Join(dir, name)
		err := osMkdirAllFunc(filepath.Dir(path), defaultFileMode)
		if err != nil {
			return 0, err
		}

		f, err := osCreateFunc(path)
		if err != nil {
			return 0, err
		}
		defer f.Close()

		return f.Write(p)
	}
}

func removeFileFunc(dir string) func(name string) error {
	return func(name string) error {
		if isInvalidPath(name) {
			return &os.PathError{Op: "remove", Path: name, Err: os.ErrInvalid}
		}
		return osRemoveFunc(filepath.Join(dir, name))
	}
}

func removeAllFunc(dir string) func(path string) error {
	return func(path string) error {
		if isInvalidPath(path) {
			return &os.PathError{Op: "removeAll", Path: path, Err: os.ErrInvalid}
		}
		return osRemoveAllFunc(filepath.Join(dir, path))
	}
}
