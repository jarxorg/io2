package io2

import (
	"io"
)

// Delegate holds Reader, Writer, Seeker, Closer functions.
type Delegate struct {
	Read  func(p []byte) (n int, err error)
	Write func(p []byte) (n int, err error)
	Seek  func(offset int64, whence int) (int64, error)
	Close func() error
}

// Delegator implements Reader, Writer, Seeker, Closer.
type Delegator struct {
	Delegate Delegate
}

// DelegateReader returns a Delegator with the provided Read function.
func DelegateReader(i io.Reader) *Delegator {
	return &Delegator{
		Delegate: Delegate{
			Read: i.Read,
		},
	}
}

// DelegateReadCloser returns a Delegator with the provided Read and Close functions.
func DelegateReadCloser(i io.ReadCloser) *Delegator {
	return &Delegator{
		Delegate: Delegate{
			Read:  i.Read,
			Close: i.Close,
		},
	}
}

// DelegateReadSeeker returns a Delegator with the provided Read and Seek functions.
func DelegateReadSeeker(i io.ReadSeeker) *Delegator {
	return &Delegator{
		Delegate: Delegate{
			Read: i.Read,
			Seek: i.Seek,
		},
	}
}

// DelegateReadSeekCloser returns a Delegator with the provided Read, Seek and Close functions.
func DelegateReadSeekCloser(i io.ReadSeekCloser) *Delegator {
	return &Delegator{
		Delegate: Delegate{
			Read:  i.Read,
			Seek:  i.Seek,
			Close: i.Close,
		},
	}
}

// DelegateReadWriteCloser returns a Delegator with the provided Read, Write and Close functions.
func DelegateReadWriteCloser(i io.ReadWriteCloser) *Delegator {
	return &Delegator{
		Delegate: Delegate{
			Read:  i.Read,
			Write: i.Write,
			Close: i.Close,
		},
	}
}

// DelegateReadWriteSeeker returns a Delegator with the provided Read, Write and Seek functions.
func DelegateReadWriteSeeker(i io.ReadWriteSeeker) *Delegator {
	return &Delegator{
		Delegate: Delegate{
			Read:  i.Read,
			Write: i.Write,
			Seek:  i.Seek,
		},
	}
}

// DelegateReadWriter returns a Delegator with the provided Read and Write functions.
func DelegateReadWriter(i io.ReadWriter) *Delegator {
	return &Delegator{
		Delegate: Delegate{
			Read:  i.Read,
			Write: i.Write,
		},
	}
}

// DelegateWriter returns a Delegator with the provided Write function.
func DelegateWriter(i io.Writer) *Delegator {
	return &Delegator{
		Delegate: Delegate{
			Write: i.Write,
		},
	}
}

// DelegateWriteCloser returns a Delegator with the provided Write and Close functions.
func DelegateWriteCloser(i io.WriteCloser) *Delegator {
	return &Delegator{
		Delegate: Delegate{
			Write: i.Write,
			Close: i.Close,
		},
	}
}

// DelegateWriteSeeker returns a Delegator with the provided Write and Seek functions.
func DelegateWriteSeeker(i io.WriteSeeker) *Delegator {
	return &Delegator{
		Delegate: Delegate{
			Write: i.Write,
			Seek:  i.Seek,
		},
	}
}

// DelegateWriteSeekCloser returns a Delegator with the provided Write , Seek and Close functions.
func DelegateWriteSeekCloser(i WriteSeekCloser) *Delegator {
	return &Delegator{
		Delegate: Delegate{
			Write: i.Write,
			Seek:  i.Seek,
			Close: i.Close,
		},
	}
}

// Read calls Delegate.Read(p).
func (d *Delegator) Read(p []byte) (int, error) {
	if d.Delegate.Read == nil {
		return 0, nil
	}
	return d.Delegate.Read(p)
}

// Write calls Delegate.Write(p).
func (d *Delegator) Write(p []byte) (int, error) {
	if d.Delegate.Write == nil {
		return 0, nil
	}
	return d.Delegate.Write(p)
}

// Seek calls Delegate.Seek(offset, whence).
func (d *Delegator) Seek(offset int64, whence int) (int64, error) {
	if d.Delegate.Seek == nil {
		return 0, nil
	}
	return d.Delegate.Seek(offset, whence)
}

// Close calls Delegate.Close().
func (d *Delegator) Close() error {
	if d.Delegate.Close == nil {
		return nil
	}
	return d.Delegate.Close()
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
