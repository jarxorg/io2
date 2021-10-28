package io2

import (
	"io"
	"io/fs"
)

// Delegator implements Reader, Writer, Seeker, Closer.
type Delegator struct {
	ReadFunc  func(p []byte) (n int, err error)
	WriteFunc func(p []byte) (n int, err error)
	SeekFunc  func(offset int64, whence int) (int64, error)
	CloseFunc func() error
}

var (
	_ io.Reader = (*Delegator)(nil)
	_ io.Writer = (*Delegator)(nil)
	_ io.Seeker = (*Delegator)(nil)
	_ io.Closer = (*Delegator)(nil)
)

// Read calls ReadFunc(p).
func (d *Delegator) Read(p []byte) (int, error) {
	if d.ReadFunc == nil {
		return 0, nil
	}
	return d.ReadFunc(p)
}

// Write calls WriteFunc(p).
func (d *Delegator) Write(p []byte) (int, error) {
	if d.WriteFunc == nil {
		return 0, nil
	}
	return d.WriteFunc(p)
}

// Seek calls SeekFunc(offset, whence).
func (d *Delegator) Seek(offset int64, whence int) (int64, error) {
	if d.SeekFunc == nil {
		return 0, nil
	}
	return d.SeekFunc(offset, whence)
}

// Close calls CloseFunc().
func (d *Delegator) Close() error {
	if d.CloseFunc == nil {
		return nil
	}
	return d.CloseFunc()
}

// DelegateReader returns a Delegator with the provided Read function.
func DelegateReader(i io.Reader) *Delegator {
	return &Delegator{
		ReadFunc: i.Read,
	}
}

// DelegateReadCloser returns a Delegator with the provided Read and Close functions.
func DelegateReadCloser(i io.ReadCloser) *Delegator {
	return &Delegator{
		ReadFunc:  i.Read,
		CloseFunc: i.Close,
	}
}

// DelegateReadSeeker returns a Delegator with the provided Read and Seek functions.
func DelegateReadSeeker(i io.ReadSeeker) *Delegator {
	return &Delegator{
		ReadFunc: i.Read,
		SeekFunc: i.Seek,
	}
}

// DelegateReadSeekCloser returns a Delegator with the provided Read, Seek and Close functions.
func DelegateReadSeekCloser(i io.ReadSeekCloser) *Delegator {
	return &Delegator{
		ReadFunc:  i.Read,
		SeekFunc:  i.Seek,
		CloseFunc: i.Close,
	}
}

// DelegateReadWriteCloser returns a Delegator with the provided Read, Write and Close functions.
func DelegateReadWriteCloser(i io.ReadWriteCloser) *Delegator {
	return &Delegator{
		ReadFunc:  i.Read,
		WriteFunc: i.Write,
		CloseFunc: i.Close,
	}
}

// DelegateReadWriteSeeker returns a Delegator with the provided Read, Write and Seek functions.
func DelegateReadWriteSeeker(i io.ReadWriteSeeker) *Delegator {
	return &Delegator{
		ReadFunc:  i.Read,
		WriteFunc: i.Write,
		SeekFunc:  i.Seek,
	}
}

// DelegateReadWriter returns a Delegator with the provided Read and Write functions.
func DelegateReadWriter(i io.ReadWriter) *Delegator {
	return &Delegator{
		ReadFunc:  i.Read,
		WriteFunc: i.Write,
	}
}

// DelegateWriter returns a Delegator with the provided Write function.
func DelegateWriter(i io.Writer) *Delegator {
	return &Delegator{
		WriteFunc: i.Write,
	}
}

// DelegateWriteCloser returns a Delegator with the provided Write and Close functions.
func DelegateWriteCloser(i io.WriteCloser) *Delegator {
	return &Delegator{
		WriteFunc: i.Write,
		CloseFunc: i.Close,
	}
}

// DelegateWriteSeeker returns a Delegator with the provided Write and Seek functions.
func DelegateWriteSeeker(i io.WriteSeeker) *Delegator {
	return &Delegator{
		WriteFunc: i.Write,
		SeekFunc:  i.Seek,
	}
}

// DelegateWriteSeekCloser returns a Delegator with the provided Write, Seek and Close functions.
func DelegateWriteSeekCloser(i WriteSeekCloser) *Delegator {
	return &Delegator{
		WriteFunc: i.Write,
		SeekFunc:  i.Seek,
		CloseFunc: i.Close,
	}
}

// NopReadCloser returns a ReadCloser with a no-op Close method wrapping the provided interface.
// This function like io.NopCloser(io.Reader).
func NopReadCloser(r io.Reader) io.ReadCloser {
	return DelegateReader(r)
}

// NopReadWriteCloser returns a ReadWriteCloser with a no-op Close method wrapping the provided interface.
func NopReadWriteCloser(rw io.ReadWriter) io.ReadWriteCloser {
	return DelegateReadWriter(rw)
}

// NopReadSeekCloser returns a ReadSeekCloser with a no-op Close method wrapping the provided interface.
func NopReadSeekCloser(r io.ReadSeeker) io.ReadSeekCloser {
	return DelegateReadSeeker(r)
}

// NopWriteCloser returns a WriteCloser with a no-op Close method wrapping the provided interface.
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return DelegateWriter(w)
}

// FSDelegator implements FS, ReadDirFS, ReadFileFS, StatFS, SubFS interface.
type FSDelegator struct {
	OpenFunc     func(name string) (fs.File, error)
	ReadDirFunc  func(name string) ([]fs.DirEntry, error)
	ReadFileFunc func(name string) ([]byte, error)
	GlobFunc     func(pattern string) ([]string, error)
	StatFunc     func(name string) (fs.FileInfo, error)
	SubFunc      func(dir string) (fs.FS, error)
}

var (
	_ fs.FS         = (*FSDelegator)(nil)
	_ fs.GlobFS     = (*FSDelegator)(nil)
	_ fs.ReadDirFS  = (*FSDelegator)(nil)
	_ fs.ReadFileFS = (*FSDelegator)(nil)
	_ fs.StatFS     = (*FSDelegator)(nil)
	_ fs.SubFS      = (*FSDelegator)(nil)
)

// Open calls OpenFunc(name).
func (d *FSDelegator) Open(name string) (fs.File, error) {
	if d.OpenFunc == nil {
		return nil, nil
	}
	return d.OpenFunc(name)
}

// ReadDir calls ReadDirFunc(name).
func (d *FSDelegator) ReadDir(name string) ([]fs.DirEntry, error) {
	if d.ReadDirFunc == nil {
		return nil, nil
	}
	return d.ReadDirFunc(name)
}

// ReadFile calls ReadFileFunc(name).
func (d *FSDelegator) ReadFile(name string) ([]byte, error) {
	if d.ReadFileFunc == nil {
		return nil, nil
	}
	return d.ReadFileFunc(name)
}

// Glob calls GlobFunc(name).
func (d *FSDelegator) Glob(pattern string) ([]string, error) {
	if d.GlobFunc == nil {
		return nil, nil
	}
	return d.GlobFunc(pattern)
}

// Stat calls StatFunc(name).
func (d *FSDelegator) Stat(name string) (fs.FileInfo, error) {
	if d.StatFunc == nil {
		return nil, nil
	}
	return d.StatFunc(name)
}

// Sub calls SubFunc(name).
func (d *FSDelegator) Sub(name string) (fs.FS, error) {
	if d.SubFunc == nil {
		return nil, nil
	}
	return d.SubFunc(name)
}

// DelegateFS returns a FSDelegator with the Open, ReadDir, ReadFile, Stat, Sub functions.
func DelegateFS(fsys fs.FS) *FSDelegator {
	d := &FSDelegator{
		OpenFunc: fsys.Open,
	}
	if fsys, ok := fsys.(fs.ReadDirFS); ok {
		d.ReadDirFunc = fsys.ReadDir
	}
	if fsys, ok := fsys.(fs.ReadFileFS); ok {
		d.ReadFileFunc = fsys.ReadFile
	}
	if fsys, ok := fsys.(fs.GlobFS); ok {
		d.GlobFunc = fsys.Glob
	}
	if fsys, ok := fsys.(fs.StatFS); ok {
		d.StatFunc = fsys.Stat
	}
	if fsys, ok := fsys.(fs.SubFS); ok {
		d.SubFunc = fsys.Sub
	}
	return d
}

// FileDelegator implements File, ReadDirFile interface.
type FileDelegator struct {
	StatFunc    func() (fs.FileInfo, error)
	ReadFunc    func([]byte) (int, error)
	CloseFunc   func() error
	ReadDirFunc func(n int) ([]fs.DirEntry, error)
}

var (
	_ fs.File        = (*FileDelegator)(nil)
	_ fs.ReadDirFile = (*FileDelegator)(nil)
)

// Stat calls StatFunc().
func (f *FileDelegator) Stat() (fs.FileInfo, error) {
	if f.StatFunc == nil {
		return nil, nil
	}
	return f.StatFunc()
}

// Read calls ReadFunc(p).
func (f *FileDelegator) Read(p []byte) (int, error) {
	if f.ReadFunc == nil {
		return 0, nil
	}
	return f.ReadFunc(p)
}

// Close calls CloseFunc().
func (f *FileDelegator) Close() error {
	if f.CloseFunc == nil {
		return nil
	}
	return f.CloseFunc()
}

// ReadDir calls ReadDirFunc(n).
func (f *FileDelegator) ReadDir(n int) ([]fs.DirEntry, error) {
	if f.ReadDirFunc == nil {
		return nil, nil
	}
	return f.ReadDirFunc(n)
}

// DelegateFile returns a FileDelegator with the Stat, Read, Close, ReadDir functions.
func DelegateFile(f fs.File) *FileDelegator {
	d := &FileDelegator{
		StatFunc:  f.Stat,
		ReadFunc:  f.Read,
		CloseFunc: f.Close,
	}
	if f, ok := f.(fs.ReadDirFile); ok {
		d.ReadDirFunc = f.ReadDir
	}
	return d
}
