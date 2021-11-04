package io2

import (
	"io/fs"
)

// OpenFSDelegator implements fs.FS interface.
type OpenFSDelegator struct {
	OpenFunc func(name string) (fs.File, error)
}

// Open calls OpenFunc(name).
func (d *OpenFSDelegator) Open(name string) (fs.File, error) {
	if d.OpenFunc == nil {
		return nil, &fs.PathError{Op: "Open", Path: name, Err: ErrNotImplemented}
	}
	return d.OpenFunc(name)
}

// DelegateOpenFS returns a OpenFSDelegator delegate fsys.Open.
func DelegateOpenFS(fsys fs.FS) *OpenFSDelegator {
	return &OpenFSDelegator{OpenFunc: fsys.Open}
}

// FSDelegator implements all filesystem interfaces in io/fs and WriteFileFS.
type FSDelegator struct {
	OpenFunc       func(name string) (fs.File, error)
	ReadDirFunc    func(name string) ([]fs.DirEntry, error)
	ReadFileFunc   func(name string) ([]byte, error)
	GlobFunc       func(pattern string) ([]string, error)
	StatFunc       func(name string) (fs.FileInfo, error)
	SubFunc        func(dir string) (fs.FS, error)
	CreateFileFunc func(name string) (WriterFile, error)
	WriteFileFunc  func(name string, p []byte) (int, error)
	RemoveFileFunc func(name string) error
	RemoveAllFunc  func(path string) error
}

var (
	_ fs.FS         = (*FSDelegator)(nil)
	_ fs.GlobFS     = (*FSDelegator)(nil)
	_ fs.ReadDirFS  = (*FSDelegator)(nil)
	_ fs.ReadFileFS = (*FSDelegator)(nil)
	_ fs.StatFS     = (*FSDelegator)(nil)
	_ fs.SubFS      = (*FSDelegator)(nil)
	_ WriteFileFS   = (*FSDelegator)(nil)
	_ RemoveFileFS  = (*FSDelegator)(nil)
)

// Open calls OpenFunc(name).
func (d *FSDelegator) Open(name string) (fs.File, error) {
	if d.OpenFunc == nil {
		return nil, &fs.PathError{Op: "Open", Path: name, Err: ErrNotImplemented}
	}
	return d.OpenFunc(name)
}

// ReadDir calls ReadDirFunc(name).
func (d *FSDelegator) ReadDir(name string) ([]fs.DirEntry, error) {
	if d.ReadDirFunc == nil {
		return nil, &fs.PathError{Op: "ReadDir", Path: name, Err: ErrNotImplemented}
	}
	return d.ReadDirFunc(name)
}

// ReadFile calls ReadFileFunc(name).
func (d *FSDelegator) ReadFile(name string) ([]byte, error) {
	if d.ReadFileFunc == nil {
		return nil, &fs.PathError{Op: "ReadFile", Path: name, Err: ErrNotImplemented}
	}
	return d.ReadFileFunc(name)
}

// Glob calls GlobFunc(name).
func (d *FSDelegator) Glob(pattern string) ([]string, error) {
	if d.GlobFunc == nil {
		return nil, &fs.PathError{Op: "Glob", Path: pattern, Err: ErrNotImplemented}
	}
	return d.GlobFunc(pattern)
}

// Stat calls StatFunc(name).
func (d *FSDelegator) Stat(name string) (fs.FileInfo, error) {
	if d.StatFunc == nil {
		return nil, &fs.PathError{Op: "Stat", Path: name, Err: ErrNotImplemented}
	}
	return d.StatFunc(name)
}

// Sub calls SubFunc(name).
func (d *FSDelegator) Sub(name string) (fs.FS, error) {
	if d.SubFunc == nil {
		return nil, &fs.PathError{Op: "Sub", Path: name, Err: ErrNotImplemented}
	}
	return d.SubFunc(name)
}

// CreateFile calls CreateFileFunc(name).
func (d *FSDelegator) CreateFile(name string) (WriterFile, error) {
	if d.CreateFileFunc == nil {
		return nil, &fs.PathError{Op: "CreateFile", Path: name, Err: ErrNotImplemented}
	}
	return d.CreateFileFunc(name)
}

// WriteFile calls WriteFileFunc(name).
func (d *FSDelegator) WriteFile(name string, p []byte) (int, error) {
	if d.WriteFileFunc == nil {
		return 0, &fs.PathError{Op: "WriteFile", Path: name, Err: ErrNotImplemented}
	}
	return d.WriteFileFunc(name, p)
}

// RemoveFile calls RemoveFileFunc(name).
func (d *FSDelegator) RemoveFile(name string) error {
	if d.RemoveFileFunc == nil {
		return &fs.PathError{Op: "RemoveFile", Path: name, Err: ErrNotImplemented}
	}
	return d.RemoveFileFunc(name)
}

// RemoveAll calls RemoveAllFunc(name).
func (d *FSDelegator) RemoveAll(path string) error {
	if d.RemoveAllFunc == nil {
		return &fs.PathError{Op: "RemoveAll", Path: path, Err: ErrNotImplemented}
	}
	return d.RemoveAllFunc(path)
}

// DelegateFS returns a FSDelegator delegate the functions of the specified filesystem.
// If you want to delegate an open only filesystem like os.DirFS(dir string) use DelegateOpenFS instead.
func DelegateFS(fsys fs.FS) *FSDelegator {
	d := &FSDelegator{
		OpenFunc: fsys.Open,
	}
	if casted, ok := fsys.(fs.ReadDirFS); ok {
		d.ReadDirFunc = casted.ReadDir
	} else {
		d.ReadDirFunc = func(name string) ([]fs.DirEntry, error) {
			return fs.ReadDir(fsys, name)
		}
	}
	if casted, ok := fsys.(fs.ReadFileFS); ok {
		d.ReadFileFunc = casted.ReadFile
	} else {
		d.ReadFileFunc = func(name string) ([]byte, error) {
			return fs.ReadFile(fsys, name)
		}
	}
	if casted, ok := fsys.(fs.GlobFS); ok {
		d.GlobFunc = casted.Glob
	} else {
		d.GlobFunc = func(pattern string) ([]string, error) {
			return fs.Glob(fsys, pattern)
		}
	}
	if casted, ok := fsys.(fs.StatFS); ok {
		d.StatFunc = casted.Stat
	} else {
		d.StatFunc = func(name string) (fs.FileInfo, error) {
			return fs.Stat(fsys, name)
		}
	}
	if casted, ok := fsys.(fs.SubFS); ok {
		d.SubFunc = casted.Sub
	} else {
		d.SubFunc = func(dir string) (fs.FS, error) {
			return fs.Sub(fsys, dir)
		}
	}
	if casted, ok := fsys.(WriteFileFS); ok {
		d.CreateFileFunc = casted.CreateFile
		d.WriteFileFunc = casted.WriteFile
	}
	if casted, ok := fsys.(RemoveFileFS); ok {
		d.RemoveFileFunc = casted.RemoveFile
		d.RemoveAllFunc = casted.RemoveAll
	}
	return d
}

// FileDelegator implements fs.File, fs.ReadDirFile and WriterFile interface.
type FileDelegator struct {
	StatFunc    func() (fs.FileInfo, error)
	ReadFunc    func(p []byte) (int, error)
	CloseFunc   func() error
	ReadDirFunc func(n int) ([]fs.DirEntry, error)
	WriteFunc   func(p []byte) (int, error)
}

var (
	_ fs.File        = (*FileDelegator)(nil)
	_ fs.ReadDirFile = (*FileDelegator)(nil)
	_ WriterFile     = (*FileDelegator)(nil)
)

// Stat calls StatFunc().
func (f *FileDelegator) Stat() (fs.FileInfo, error) {
	if f.StatFunc == nil {
		return nil, ErrNotImplemented
	}
	return f.StatFunc()
}

// Read calls ReadFunc(p).
func (f *FileDelegator) Read(p []byte) (int, error) {
	if f.ReadFunc == nil {
		return 0, ErrNotImplemented
	}
	return f.ReadFunc(p)
}

// Close calls CloseFunc().
func (f *FileDelegator) Close() error {
	if f.CloseFunc == nil {
		// NOTE: return no error.
		return nil
	}
	return f.CloseFunc()
}

// ReadDir calls ReadDirFunc(n).
func (f *FileDelegator) ReadDir(n int) ([]fs.DirEntry, error) {
	if f.ReadDirFunc == nil {
		return nil, ErrNotImplemented
	}
	return f.ReadDirFunc(n)
}

// Write calls WriteFunc(n).
func (f *FileDelegator) Write(p []byte) (int, error) {
	if f.WriteFunc == nil {
		return 0, ErrNotImplemented
	}
	return f.WriteFunc(p)
}

// DelegateFile returns a FileDelegator delegate the functions of the specified file.
func DelegateFile(f fs.File) *FileDelegator {
	d := &FileDelegator{
		StatFunc:  f.Stat,
		ReadFunc:  f.Read,
		CloseFunc: f.Close,
	}
	if f, ok := f.(fs.ReadDirFile); ok {
		d.ReadDirFunc = f.ReadDir
	}
	if f, ok := f.(WriterFile); ok {
		d.WriteFunc = f.Write
	}
	return d
}
