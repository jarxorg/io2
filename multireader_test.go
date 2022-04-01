package io2

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func testMultiReaderStrings(t *testing.T, ss ...string) io.ReadSeekCloser {
	rs := make([]io.ReadSeekCloser, len(ss))
	for i, s := range ss {
		rs[i] = NopReadSeekCloser(strings.NewReader(s))
	}
	mr, err := MultiReadSeekCloser(rs...)
	if err != nil {
		t.Fatal(err)
	}
	return mr
}

func TestMultiRead(t *testing.T) {
	tests := []struct {
		reader func() io.ReadCloser
		n      int
		errstr string
		want   string
	}{
		{
			reader: func() io.ReadCloser {
				return testMultiReaderStrings(t, "abc", "def")
			},
			n:    6,
			want: "abcdef",
		}, {
			reader: func() io.ReadCloser {
				return MultiReadCloser(&Delegator{
					ReadFunc: func(p []byte) (int, error) {
						return 0, errors.New("failed to read for coverage")
					},
				})
			},
			errstr: "failed to read for coverage",
		},
	}

	var deferClose func() error
	close := func() {
		if deferClose != nil {
			deferClose()
		}
	}
	defer close()

	for i, test := range tests {
		close()
		r := test.reader()
		deferClose = r.Close

		p := make([]byte, test.n)
		n, err := r.Read(p)
		if test.errstr != "" {
			if err == nil {
				t.Fatalf("tests[%d] no error", i)
			}
			if err.Error() != test.errstr {
				t.Errorf("tests[%d] error %s; want %s", i, err.Error(), test.errstr)
			}
			continue
		}
		if err != nil {
			t.Fatalf("tests[%d] error %v", i, err)
		}
		if n != test.n {
			t.Errorf("tests[%d] n is %d; want %d", i, n, test.n)
		}
		got := string(p)
		if got != test.want {
			t.Errorf("tests[%d] got %s; want %s", i, got, test.want)
		}
	}
}

func TestMultiSeek(t *testing.T) {
	newErrSeekReader := func() io.ReadSeekCloser {
		return &multiReader{
			rs: []*reader{
				{
					ReadSeekCloser: &Delegator{
						SeekFunc: func(offset int64, whence int) (int64, error) {
							return 0, errors.New("failed to seek for coverage")
						},
					},
				},
			},
		}
	}

	tests := []struct {
		reader func() io.ReadSeekCloser
		offset int64
		whence int
		n      int64
		errstr string
		after  string
	}{
		{
			reader: func() io.ReadSeekCloser {
				// NOTE: Check single reader.
				return NopReadSeekCloser(strings.NewReader("abcdef"))
			},
			offset: int64(2),
			whence: io.SeekStart,
			n:      int64(2),
			after:  "cdef",
		}, {
			reader: func() io.ReadSeekCloser {
				return testMultiReaderStrings(t, "abc", "def")
			},
			offset: int64(2),
			whence: io.SeekStart,
			n:      int64(2),
			after:  "cdef",
		}, {
			reader: func() io.ReadSeekCloser {
				return testMultiReaderStrings(t, "abc", "def")
			},
			offset: int64(4),
			whence: io.SeekStart,
			n:      int64(4),
			after:  "ef",
		}, {
			reader: func() io.ReadSeekCloser {
				r := testMultiReaderStrings(t, "abc", "def")
				ioutil.ReadAll(r)
				return r
			},
			offset: int64(0),
			whence: io.SeekStart,
			n:      int64(0),
			after:  "abcdef",
		}, {
			reader: func() io.ReadSeekCloser {
				r0 := NopReadSeekCloser(strings.NewReader("abc"))
				r1 := strings.NewReader("def")
				d1 := Delegate(r1)
				r, err := MultiReadSeekCloser(r0, d1)
				if err != nil {
					t.Fatal(err)
				}
				d1.SeekFunc = func(offset int64, whence int) (int64, error) {
					return 0, errors.New("failed to reset tails")
				}
				return r
			},
			offset: int64(0),
			whence: io.SeekStart,
			errstr: "failed to reset tails",
		}, {
			reader: func() io.ReadSeekCloser {
				// NOTE: Check single reader.
				return NopReadSeekCloser(strings.NewReader("abcdefghi"))
			},
			offset: int64(-4),
			whence: io.SeekEnd,
			n:      int64(5),
			after:  "fghi",
		}, {
			reader: func() io.ReadSeekCloser {
				return testMultiReaderStrings(t, "abc", "def", "ghi")
			},
			offset: int64(-4),
			whence: io.SeekEnd,
			n:      int64(5),
			after:  "fghi",
		}, {
			reader: func() io.ReadSeekCloser {
				return testMultiReaderStrings(t, "abc", "def", "ghi")
			},
			offset: int64(-9),
			whence: io.SeekEnd,
			n:      int64(0),
			after:  "abcdefghi",
		}, {
			reader: func() io.ReadSeekCloser {
				// NOTE: Check single reader.
				r := NopReadSeekCloser(strings.NewReader("abcdefghi"))
				if _, err := r.Seek(4, io.SeekStart); err != nil {
					t.Fatal(err)
				}
				return r
			},
			offset: int64(3),
			whence: io.SeekCurrent,
			n:      int64(7),
			after:  "hi",
		}, {
			reader: func() io.ReadSeekCloser {
				r := testMultiReaderStrings(t, "abc", "def", "ghi")
				if _, err := r.Seek(4, io.SeekStart); err != nil {
					t.Fatal(err)
				}
				return r
			},
			offset: int64(3),
			whence: io.SeekCurrent,
			n:      int64(7),
			after:  "hi",
		}, {
			reader: func() io.ReadSeekCloser {
				r := testMultiReaderStrings(t, "abc", "def", "ghi")
				if _, err := r.Seek(4, io.SeekStart); err != nil {
					t.Fatal(err)
				}
				return r
			},
			offset: int64(-1),
			whence: io.SeekCurrent,
			n:      int64(3),
			after:  "defghi",
		}, {
			reader: func() io.ReadSeekCloser {
				r := testMultiReaderStrings(t, "abc", "def", "ghi")
				if _, err := r.Seek(4, io.SeekStart); err != nil {
					t.Fatal(err)
				}
				return r
			},
			offset: int64(-4),
			whence: io.SeekCurrent,
			n:      int64(0),
			after:  "abcdefghi",
		}, {
			reader: func() io.ReadSeekCloser {
				r0 := NopReadSeekCloser(strings.NewReader("abc"))
				r1 := strings.NewReader("def")
				d1 := Delegate(r1)
				r, err := MultiReadSeekCloser(r0, d1)
				if err != nil {
					t.Fatal(err)
				}
				d1.SeekFunc = func(offset int64, whence int) (int64, error) {
					return 0, errors.New("failed to reset tails")
				}
				return r
			},
			offset: int64(0),
			whence: io.SeekCurrent,
			errstr: "failed to reset tails",
		}, {
			reader: func() io.ReadSeekCloser {
				// NOTE: Check single reader.
				return NopReadSeekCloser(strings.NewReader("abcdef"))
			},
			offset: int64(-1),
			whence: io.SeekStart,
			errstr: "strings.Reader.Seek: negative position",
		}, {
			reader: func() io.ReadSeekCloser {
				return testMultiReaderStrings(t, "abc", "def")
			},
			offset: int64(-1),
			whence: io.SeekStart,
			errstr: "strings.Reader.Seek: negative position",
		}, {
			reader: func() io.ReadSeekCloser {
				// NOTE: Check single reader.
				return NopReadSeekCloser(strings.NewReader("abcdef"))
			},
			offset: int64(10),
			whence: io.SeekStart,
			n:      int64(10),
		}, {
			reader: func() io.ReadSeekCloser {
				return testMultiReaderStrings(t, "abc", "def")
			},
			offset: int64(10),
			whence: io.SeekStart,
			n:      int64(10),
		}, {
			reader: func() io.ReadSeekCloser {
				return testMultiReaderStrings(t, "abc", "def")
			},
			whence: -1,
			errstr: "invalid whence",
		}, {
			reader: newErrSeekReader,
			whence: io.SeekStart,
			errstr: "failed to seek for coverage",
		}, {
			reader: newErrSeekReader,
			whence: io.SeekCurrent,
			errstr: "failed to seek for coverage",
		}, {
			reader: newErrSeekReader,
			whence: io.SeekEnd,
			errstr: "failed to seek for coverage",
		},
	}

	var deferClose func() error
	close := func() {
		if deferClose != nil {
			deferClose()
		}
	}
	defer close()

	for i, test := range tests {
		close()
		r := test.reader()
		deferClose = r.Close

		n, err := r.Seek(test.offset, test.whence)
		if test.errstr != "" {
			if err == nil {
				t.Fatalf("tests[%d] no error", i)
			}
			if err.Error() != test.errstr {
				t.Errorf("tests[%d] error %s; want %s", i, err.Error(), test.errstr)
			}
			continue
		}
		if err != nil {
			t.Fatalf("tests[%d] error %v", i, err)
		}
		if n != test.n {
			t.Errorf("tests[%d] n is %d; want %d", i, n, test.n)
		}

		p, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatalf("tests[%d] error %v", i, err)
		}
		after := string(p)
		if after != test.after {
			t.Errorf("tests[%d] after %s; want %s", i, after, test.after)
		}
	}

}

func TestMultiReadSeeker_Errors(t *testing.T) {
	tests := []struct {
		r      io.ReadSeeker
		errstr string
	}{
		{
			r: &Delegator{
				SeekFunc: func(offset int64, whence int) (int64, error) {
					if whence == io.SeekStart {
						return 0, errors.New("failed to seek start for coverage")
					}
					return 0, nil
				},
			},
			errstr: "failed to seek start for coverage",
		}, {
			r: &Delegator{
				SeekFunc: func(offset int64, whence int) (int64, error) {
					if whence == io.SeekEnd {
						return 0, errors.New("failed to seek end for coverage")
					}
					return 0, nil
				},
			},
			errstr: "failed to seek end for coverage",
		},
	}
	for i, test := range tests {
		_, err := MultiReadSeeker(test.r)
		if err == nil {
			t.Fatalf("tests[%d] no error", i)
		}
		if err.Error() != test.errstr {
			t.Errorf("tests[%d] error %s; want %s", i, err.Error(), test.errstr)
		}
	}
}

func TestMultiClose(t *testing.T) {
	errCloseReader := &Delegator{
		CloseFunc: func() error {
			return errors.New("close error")
		},
	}

	tests := []struct {
		r      io.Closer
		errstr string
	}{
		{
			r: MultiReadCloser(NopReadCloser(strings.NewReader("no error"))),
		}, {
			r:      MultiReadCloser(errCloseReader),
			errstr: "failed to close: close error",
		}, {
			r:      MultiReadCloser(errCloseReader, errCloseReader),
			errstr: "failed to close: close error; close error",
		},
	}
	for i, test := range tests {
		err := test.r.Close()
		if test.errstr != "" {
			if err == nil {
				t.Fatalf("tests[%d] no error", i)
			}
			if err.Error() != test.errstr {
				t.Errorf("tests[%d] error %s; want %s", i, err.Error(), test.errstr)
			}
			continue
		}
		if err != nil {
			t.Fatalf("tests[%d] error %v", i, err)
		}
	}
}
