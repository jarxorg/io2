package io2

import (
	"io"
	"testing"
)

func TestDelegator(t *testing.T) {
	d := &Delegator{}
	d.Delegate.Read = func(_ []byte) (int, error) {
		return 0, nil
	}
	d.Delegate.Write = func(_ []byte) (int, error) {
		return 0, nil
	}
	d.Delegate.Seek = func(_ int64, _ int) (int64, error) {
		return 0, nil
	}
	d.Delegate.Close = func() error {
		return nil
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
		t.Errorf("Error read bytes %d; want 0", z)
	}

	z, err = d.Write([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	if z != 0 {
		t.Errorf("Error read bytes %d; want 0", z)
	}

	z64, err := d.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal(err)
	}
	if z64 != 0 {
		t.Errorf("Error read bytes %d; want 0", z64)
	}

	err = d.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDelegates(t *testing.T) {
	d := &Delegator{}
	var _ io.Reader = DelegateReader(d)
	var _ io.ReadCloser = DelegateReadCloser(d)
	var _ io.ReadSeeker = DelegateReadSeeker(d)
	var _ io.ReadSeekCloser = DelegateReadSeekCloser(d)
	var _ io.ReadWriter = DelegateReadWriter(d)
	var _ io.ReadWriteCloser = DelegateReadWriteCloser(d)
	var _ io.ReadWriteSeeker = DelegateReadWriteSeeker(d)
	var _ io.Writer = DelegateWriter(d)
	var _ io.WriteCloser = DelegateWriteCloser(d)
	var _ io.WriteSeeker = DelegateWriteSeeker(d)
	var _ WriteSeekCloser = DelegateWriteSeekCloser(d)
}

func TestNops(t *testing.T) {
	d := &Delegator{}
	var _ io.ReadSeekCloser = NopReadSeekCloser(d)
	var _ io.ReadWriteCloser = NopReadWriteCloser(d)
	var _ io.WriteCloser = NopWriteCloser(d)
}
