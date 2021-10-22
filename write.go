package io2

import (
	"io"
)

// WriteSeekCloser is the interface that groups the basic Write, Seek and Close methods.
type WriteSeekCloser interface {
	io.Writer
	io.Seeker
	io.Closer
}

// WriteSeekBuffer implements io.WriteSeeker that using in-memory byte buffer.
type WriteSeekBuffer struct {
	buf []byte
	off int
	len int
}

var _ WriteSeekCloser = (*WriteSeekBuffer)(nil)

// NewWriteSeekBuffer returns an WriteSeekBuffer with the initial capacity.
func NewWriteSeekBuffer(capacity int) *WriteSeekBuffer {
	return &WriteSeekBuffer{
		buf: make([]byte, capacity),
	}
}

// NewWriteSeekBufferBytes returns an WriteSeekBuffer with the initial buffer.
func NewWriteSeekBufferBytes(buf []byte) *WriteSeekBuffer {
	off := len(buf)
	return &WriteSeekBuffer{
		buf: buf,
		off: off,
		len: off,
	}
}

// Write appends the contents of p to the buffer, growing the buffer as needed.
// The return value n is the length of p; err is always nil.
func (b *WriteSeekBuffer) Write(p []byte) (int, error) {
	capacity := len(b.buf)
	n := len(p)
	noff := b.off + n

	// NOTE: Expand inner buffer simply.
	if noff > capacity {
		ncapacity := noff
		nbuf := make([]byte, int(ncapacity))
		copy(nbuf, b.buf)
		b.buf = nbuf
	}

	copy(b.buf[b.off:], p)

	b.off = noff
	if b.len < noff {
		b.len = noff
	}

	return n, nil
}

// Seek sets the offset for the next Write to offset, interpreted according to whence:
//   SeekStart means relative to the start of the file,
//   SeekCurrent means relative to the current offset,
//   SeekEnd means relative to the end.
// Seek returns the new offset relative to the start of the file and an error, if any.
func (b *WriteSeekBuffer) Seek(offset int64, whence int) (int64, error) {
	off := int(offset)
	noff := 0
	switch whence {
	case io.SeekStart:
		noff = off
	case io.SeekCurrent:
		noff = b.off + off
	case io.SeekEnd:
		noff = b.len + off
	}
	if noff < 0 {
		noff = 0
	}
	b.off = noff
	return int64(noff), nil
}

// Close calls b.Truncate(0).
func (b *WriteSeekBuffer) Close() error {
	b.Truncate(0)
	return nil
}

// Offset returns the offset.
func (b *WriteSeekBuffer) Offset() int {
	return b.off
}

// Len returns the number of bytes of the buffer; b.Len() == len(b.Bytes()).
func (b *WriteSeekBuffer) Len() int {
	return b.len
}

// Bytes returns a slice of length b.Len() of the buffer.
func (b *WriteSeekBuffer) Bytes() []byte {
	return b.buf[:b.len]
}

// Truncate changes the size of the buffer with offset.
func (b *WriteSeekBuffer) Truncate(n int) {
	l := len(b.buf)
	if n < 0 {
		n = l + n
	}
	if n >= l {
		n = l
	}
	b.buf = b.buf[0:n]
	b.off = n
	b.len = n
}
