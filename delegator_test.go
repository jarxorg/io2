package io2

import (
	"io"
	"io/fs"
	"testing"
)

func TestDelegator(t *testing.T) {
	d := &Delegator{
		ReadFunc: func(_ []byte) (int, error) {
			return 0, nil
		},
		WriteFunc: func(_ []byte) (int, error) {
			return 0, nil
		},
		SeekFunc: func(_ int64, _ int) (int64, error) {
			return 0, nil
		},
		CloseFunc: func() error {
			return nil
		},
	}
	testDelegatorDefaults(t, d)
}

func TestDelegatorDefaults(t *testing.T) {
	testDelegatorDefaults(t, &Delegator{})
}

func testDelegatorDefaults(t *testing.T, d *Delegator) {
	z, err := d.Read([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	if z != 0 {
		t.Errorf("Error Read returns %d; want 0", z)
	}

	z, err = d.Write([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	if z != 0 {
		t.Errorf("Error Write returns %d; want 0", z)
	}

	z64, err := d.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal(err)
	}
	if z64 != 0 {
		t.Errorf("Error Seek returns %d; want 0", z64)
	}

	err = d.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDelegates(t *testing.T) {
	d := &Delegator{}
	var (
		_ io.Reader          = DelegateReader(d)
		_ io.ReadCloser      = DelegateReadCloser(d)
		_ io.ReadSeeker      = DelegateReadSeeker(d)
		_ io.ReadSeekCloser  = DelegateReadSeekCloser(d)
		_ io.ReadWriter      = DelegateReadWriter(d)
		_ io.ReadWriteCloser = DelegateReadWriteCloser(d)
		_ io.ReadWriteSeeker = DelegateReadWriteSeeker(d)
		_ io.Writer          = DelegateWriter(d)
		_ io.WriteCloser     = DelegateWriteCloser(d)
		_ io.WriteSeeker     = DelegateWriteSeeker(d)
		_ WriteSeekCloser    = DelegateWriteSeekCloser(d)
	)
}

func TestNops(t *testing.T) {
	d := &Delegator{}
	var (
		_ io.ReadCloser      = NopReadCloser(d)
		_ io.ReadSeekCloser  = NopReadSeekCloser(d)
		_ io.ReadWriteCloser = NopReadWriteCloser(d)
		_ io.WriteCloser     = NopWriteCloser(d)
	)
}

func TestFSDelegator(t *testing.T) {
	d := &FSDelegator{
		OpenFunc: func(_ string) (fs.File, error) {
			return nil, nil
		},
		ReadDirFunc: func(_ string) ([]fs.DirEntry, error) {
			return nil, nil
		},
		ReadFileFunc: func(_ string) ([]byte, error) {
			return nil, nil
		},
		GlobFunc: func(_ string) ([]string, error) {
			return nil, nil
		},
		StatFunc: func(_ string) (fs.FileInfo, error) {
			return nil, nil
		},
		SubFunc: func(_ string) (fs.FS, error) {
			return nil, nil
		},
	}
	testFSDelegatorDefaults(t, d)
}

func TestFSDelegatorDefaults(t *testing.T) {
	testFSDelegatorDefaults(t, &FSDelegator{})
}

func testFSDelegatorDefaults(t *testing.T, d *FSDelegator) {
	f, err := d.Open("")
	if err != nil {
		t.Fatal(err)
	}
	if f != nil {
		t.Errorf("Error Open returns not nil; want nil")
	}

	entries, err := d.ReadDir("")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("Error ReadDir returns not empty; want empty")
	}

	bin, err := d.ReadFile("")
	if err != nil {
		t.Fatal(err)
	}
	if len(bin) != 0 {
		t.Errorf("Error ReadFile returns not empty; want empty")
	}

	names, err := d.Glob("")
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 0 {
		t.Errorf("Error Glob returns not empty; want empty")
	}

	info, err := d.Stat("")
	if err != nil {
		t.Fatal(err)
	}
	if info != nil {
		t.Errorf("Error Stat returns not nil; want nil")
	}

	sub, err := d.Sub("")
	if err != nil {
		t.Fatal(err)
	}
	if sub != nil {
		t.Errorf("Error Sub returns not nil; want nil")
	}
}

func TestDelegateFS(t *testing.T) {
	var _ fs.FS = DelegateFS(&FSDelegator{})
}

func TestFileDelegator(t *testing.T) {
	d := &FileDelegator{
		StatFunc: func() (fs.FileInfo, error) {
			return nil, nil
		},
		ReadFunc: func(_ []byte) (int, error) {
			return 0, nil
		},
		CloseFunc: func() error {
			return nil
		},
		ReadDirFunc: func(_ int) ([]fs.DirEntry, error) {
			return nil, nil
		},
	}
	testFileDelegatorDefaults(t, d)
}

func TestFileDelegatorDefaults(t *testing.T) {
	testFileDelegatorDefaults(t, &FileDelegator{})
}

func testFileDelegatorDefaults(t *testing.T, d *FileDelegator) {
	info, err := d.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if info != nil {
		t.Errorf("Error Stat returns not nil; want nil")
	}

	n, err := d.Read([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("Error Read returns not zero; want 0")
	}

	err = d.Close()
	if err != nil {
		t.Fatal(err)
	}

	entries, err := d.ReadDir(-1)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("Error ReadDir returns not empty; want empty")
	}
}

func TestDelegateFile(t *testing.T) {
	var _ fs.File = DelegateFile(&FileDelegator{})
}
