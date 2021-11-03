package io2

import (
	"errors"
	"io"
	"testing"
)

func testDelegatorErrors(t *testing.T, d *Delegator, wantErr error) {
	var err error
	if _, err = d.Read([]byte{}); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
	if _, err = d.Write([]byte{}); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
	if _, err = d.Seek(0, io.SeekStart); !errors.Is(err, wantErr) {
		t.Errorf("Error unknown: %v", err)
	}
	if err = d.Close(); err != nil {
		t.Errorf("Error unknown: %v", err)
	}
}

func TestDelegator_ErrNotImplemented(t *testing.T) {
	testDelegatorErrors(t, &Delegator{}, ErrNotImplemented)
}

func TestDelegator(t *testing.T) {
	wantErr := errors.New("test")

	testDelegatorErrors(t, &Delegator{
		ReadFunc: func(_ []byte) (int, error) {
			return 0, wantErr
		},
		WriteFunc: func(_ []byte) (int, error) {
			return 0, wantErr
		},
		SeekFunc: func(_ int64, _ int) (int64, error) {
			return 0, wantErr
		},
		CloseFunc: func() error {
			return nil
		},
	}, wantErr)
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
