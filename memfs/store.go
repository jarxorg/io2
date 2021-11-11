package memfs

import (
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Value works as fs.DirEntry or fs.FileInfo.
type value struct {
	name    string
	data    []byte
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
}

var (
	_ fs.DirEntry = (*value)(nil)
	_ fs.FileInfo = (*value)(nil)
)

func (v *value) Name() string {
	return filepath.Base(v.name)
}

func (v *value) Size() int64 {
	if v.isDir {
		return 0
	}
	return int64(len(v.data))
}

func (v *value) Mode() fs.FileMode {
	return v.mode
}

func (v *value) ModTime() time.Time {
	return v.modTime
}

func (v *value) IsDir() bool {
	return v.isDir
}

func (v *value) Sys() interface{} {
	return nil
}

func (v *value) Type() fs.FileMode {
	return v.mode
}

func (v *value) Info() (fs.FileInfo, error) {
	return v, nil
}

// Store represents an in-memory key value store.
// store.keys is always sorted.
// All functions of the store are not thread safety.
type store struct {
	keys   []string
	values map[string]*value
}

func newStore() *store {
	return &store{
		values: map[string]*value{},
	}
}

func (s *store) get(k string) *value {
	return s.values[k]
}

func (s *store) put(k string, v *value) *value {
	if _, ok := s.values[k]; !ok {
		s.keys = append(s.keys, k)
		sort.Strings(s.keys)
	}

	s.values[k] = v
	return v
}

func (s *store) remove(key string) *value {
	i := s.keyIndex(key)
	if i == -1 {
		return nil
	}
	v := s.values[key]
	s.keys = append(s.keys[0:i], s.keys[i+1:]...)
	delete(s.values, key)
	return v
}

func (s *store) removeAll(prefix string) {
	from := s.keyIndex(prefix)
	if from == -1 {
		return
	}

	max := len(s.keys)
	to := -1
	for i := from; i < max; i++ {
		key := s.keys[i]
		if !strings.HasPrefix(key, prefix) {
			break
		}
		delete(s.values, key)
		to = i
	}
	s.keys = append(s.keys[0:from], s.keys[to+1:]...)
}

func (s *store) keyIndex(key string) int {
	i := sort.SearchStrings(s.keys, key)
	if i < len(s.keys) && s.keys[i] == key {
		return i
	}
	return -1
}

func (s *store) prefixKeys(prefix string) []string {
	i := s.keyIndex(prefix)
	if i == -1 {
		return nil
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	var keys []string
	max := len(s.keys)
	for i++; i < max; i++ {
		key := s.keys[i]
		if !strings.HasPrefix(key, prefix) {
			break
		}
		if strings.Contains(key[len(prefix):], "/") {
			break
		}
		keys = append(keys, key)
	}
	return keys
}

func (s *store) prefixGlobKeys(prefix, pattern string) ([]string, error) {
	i := s.keyIndex(prefix)
	if i == -1 {
		return nil, nil
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	var keys []string
	max := len(s.keys)
	for i++; i < max; i++ {
		key := s.keys[i]
		if !strings.HasPrefix(key, prefix) {
			break
		}
		ok, err := path.Match(pattern, key[len(prefix):])
		if err != nil {
			return nil, err
		}
		if ok {
			keys = append(keys, key)
		}
	}
	return keys, nil
}
