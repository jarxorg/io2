package io2

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

var errSkip = errors.New("skip")

type reader struct {
	io.ReadSeekCloser
	off    int64
	length int64
}

type multiReader struct {
	rs     []*reader
	off    int
	length int64
}

var _ io.ReadSeekCloser = (*multiReader)(nil)

func MultiReadCloser(rs ...io.ReadCloser) io.ReadCloser {
	ds := make([]*reader, len(rs))
	for i, r := range rs {
		ds[i] = &reader{
			ReadSeekCloser: Delegate(r),
		}
	}
	return &multiReader{rs: ds}
}

func MultiReadSeekCloser(rs ...io.ReadSeekCloser) (io.ReadSeekCloser, error) {
	length := int64(0)
	ds := make([]*reader, len(rs))
	for i, r := range rs {
		n, err := r.Seek(0, io.SeekEnd)
		if err != nil {
			return nil, err
		}
		if _, err = r.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		ds[i] = &reader{
			ReadSeekCloser: Delegate(r),
			length:         n,
		}
		length += n
	}
	return &multiReader{rs: ds, length: length}, nil
}

func (mr *multiReader) each(i int, offset int64, fn func(r *reader) error) error {
	mr.off = i
	if offset >= 0 {
		for ; mr.off < len(mr.rs); mr.off++ {
			if err := fn(mr.rs[mr.off]); err != nil {
				if err == errSkip {
					err = nil
				}
				return err
			}
		}
		if mr.off >= len(mr.rs) {
			mr.off = len(mr.rs) - 1
		}
		return nil
	}
	for ; mr.off >= 0; mr.off-- {
		if err := fn(mr.rs[mr.off]); err != nil {
			if err == errSkip {
				err = nil
			}
			return err
		}
	}
	if mr.off < 0 {
		mr.off = 0
	}
	return nil
}

func (mr *multiReader) Read(p []byte) (n int, err error) {
	if mr.off >= len(mr.rs) {
		return 0, io.EOF
	}
	off := 0
	for ; mr.off < len(mr.rs); mr.off++ {
		r := mr.rs[mr.off]
		n, err = r.Read(p[off:])
		if err != nil {
			return 0, err
		}
		r.off += int64(n)
		off += n
		if off >= len(p) {
			return off, nil
		}
	}
	return off, nil
}

func (mr *multiReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		return mr.seekStart(offset, 0)
	case io.SeekCurrent:
		return mr.seekCurrent(offset)
	case io.SeekEnd:
		return mr.seekEnd(offset, len(mr.rs)-1)
	}
	return 0, errors.New("invalid whence")
}

func (mr *multiReader) seekStart(offset int64, start int) (int64, error) {
	off := int64(0)
	for i := 0; i < start; i++ {
		off += mr.rs[i].length
	}
	err := mr.each(start, offset, func(r *reader) error {
		safeOffset := offset
		if mr.off != len(mr.rs)-1 && safeOffset > (r.length-r.off) {
			safeOffset = (r.length - r.off)
		}
		n, err := r.Seek(safeOffset, io.SeekStart)
		if err != nil {
			return err
		}
		r.off = n
		if offset < r.length {
			off += offset
			return errSkip
		}
		offset -= safeOffset
		off += safeOffset
		return nil
	})
	if err != nil {
		return 0, err
	}
	return off, nil
}

func (mr *multiReader) seekCurrent(offset int64) (int64, error) {
	off := int64(0)
	for i := 0; i < mr.off; i++ {
		off += mr.rs[i].length
	}

	r := mr.rs[mr.off]
	safeOffset := offset
	if offset >= 0 {
		if mr.off != len(mr.rs)-1 && safeOffset > (r.length-r.off) {
			safeOffset = (r.length - r.off)
		}
	} else {
		if mr.off != 0 && -safeOffset > r.off {
			safeOffset = -r.off
		}
	}
	n, err := r.Seek(safeOffset, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	r.off = n
	if diffset := offset - safeOffset; diffset != 0 {
		if offset >= 0 {
			return mr.seekStart(diffset, mr.off+1)
		}
		return mr.seekEnd(diffset, mr.off-1)
	}
	return off + n, nil
}

func (mr *multiReader) seekEnd(offset int64, end int) (int64, error) {
	off := mr.length
	for i := end + 1; i < len(mr.rs); i++ {
		off -= mr.rs[i].length
	}
	err := mr.each(end, offset, func(r *reader) error {
		safeOffset := offset
		if mr.off != 0 && -safeOffset > r.length {
			safeOffset = -r.length
		}
		n, err := r.Seek(safeOffset, io.SeekEnd)
		if err != nil {
			return err
		}
		r.off = n
		if -offset < r.length {
			off += offset
			return errSkip
		}
		offset += r.length
		off -= r.length
		return nil
	})
	if err != nil {
		return 0, err
	}
	return off, nil
}

func (mr *multiReader) Close() error {
	var errs []string
	mr.each(0, 1, func(r *reader) error {
		if err := r.Close(); err != nil {
			errs = append(errs, err.Error())
		}
		return nil
	})
	if len(errs) > 0 {
		return fmt.Errorf("failed to close: %s", strings.Join(errs, "; "))
	}
	mr.rs = nil
	mr.off = 0
	mr.length = 0
	return nil
}
