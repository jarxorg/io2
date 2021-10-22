package io2

import (
	"io"
	"reflect"
	"testing"
)

func TestWrite(t *testing.T) {
	testCases := []struct {
		capacity int
		p        []byte
		wantBuf  []byte
		wantOff  int
		wantLen  int
	}{
		{
			capacity: 8,
			p:        []byte(`123`),
			wantBuf:  []byte{'1', '2', '3', 0, 0, 0, 0, 0},
			wantOff:  3,
			wantLen:  3,
		}, {
			p:       []byte(`456`),
			wantBuf: []byte{'1', '2', '3', '4', '5', '6', 0, 0},
			wantOff: 6,
			wantLen: 6,
		}, {
			p:       []byte(`789`),
			wantBuf: []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'},
			wantOff: 9,
			wantLen: 9,
		},
	}

	var buf *WriteSeekBuffer
	for i, testCase := range testCases {
		if testCase.capacity > 0 {
			buf = NewWriteSeekBuffer(testCase.capacity)
			defer buf.Close()
		}
		n, err := buf.Write(testCase.p)
		if err != nil {
			t.Fatalf("Fatal[%d] write: %v", i, err)
		}
		if n != len(testCase.p) {
			t.Errorf("Error[%d] write bytes %d; want %d", i, n, len(testCase.p))
		}
		if buf.Offset() != testCase.wantOff {
			t.Errorf("Error[%d] off %d; want %d", i, buf.Offset(), testCase.wantOff)
		}
		if buf.Len() != testCase.wantLen {
			t.Errorf("Error[%d] len %d; want %d", i, buf.Len(), testCase.wantLen)
		}
		if !reflect.DeepEqual(buf.buf, testCase.wantBuf) {
			t.Errorf("Error[%d] buf %v; want %v", i, buf.buf, testCase.wantBuf)
		}
	}
}

func TestSeek(t *testing.T) {
	buf := NewWriteSeekBufferBytes([]byte(`123456789`))
	defer buf.Close()

	testCases := []struct {
		off     int64
		whence  int
		wantOff int64
	}{
		{
			off:     0,
			whence:  io.SeekStart,
			wantOff: 0,
		}, {
			off:     int64(-1),
			whence:  io.SeekStart,
			wantOff: 0,
		}, {
			off:     0,
			whence:  io.SeekEnd,
			wantOff: int64(9),
		}, {
			off:     int64(-1),
			whence:  io.SeekCurrent,
			wantOff: int64(8),
		}, {
			off:     int64(-1),
			whence:  io.SeekCurrent,
			wantOff: int64(7),
		}, {
			off:     int64(-3),
			whence:  io.SeekEnd,
			wantOff: int64(6),
		},
	}

	for i, testCase := range testCases {
		n, err := buf.Seek(testCase.off, testCase.whence)
		if err != nil {
			t.Fatalf("Fatal seek: %v", err)
		}
		if n != testCase.wantOff {
			t.Errorf("Error[%d] off %d; want %d", i, n, testCase.wantOff)
		}
	}
}

func TestSeekWrite(t *testing.T) {
	buf := NewWriteSeekBufferBytes([]byte(`123456789`))
	defer buf.Close()

	off, err := buf.Seek(int64(3), io.SeekStart)
	if err != nil {
		t.Fatalf("Fatal seek: %v", err)
	}
	if off != 3 {
		t.Errorf("Error seek off %d; want %d", off, 3)
	}
	n, err := buf.Write([]byte(`def`))
	if err != nil {
		t.Fatalf("Fatal write: %v", err)
	}
	if n != 3 {
		t.Errorf("Error write bytes %d; want %d", n, 3)
	}

	want := `123def789`
	got := string(buf.Bytes())
	if got != want {
		t.Errorf("Error bytes %s; want %s", got, want)
	}
}

func TestTruncate(t *testing.T) {
	buf := NewWriteSeekBufferBytes([]byte(`123456789`))
	defer buf.Close()

	testCases := []struct {
		n    int
		want []byte
	}{
		{
			n:    8,
			want: []byte(`12345678`),
		}, {
			n:    -5,
			want: []byte(`123`),
		}, {
			n:    100,
			want: []byte(`123`),
		},
	}

	for i, testCase := range testCases {
		buf.Truncate(testCase.n)
		got := buf.Bytes()
		if !reflect.DeepEqual(got, testCase.want) {
			t.Errorf("Error[%d] truncate bytes %v; want %v", i, got, testCase.want)
		}
	}
}
