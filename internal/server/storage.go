package server

import (
	"os"
	"path/filepath"
)

const dataDir = "data"

func ensureDataDir() {
	os.MkdirAll(dataDir, 0755)
}

func OpenForWrite(name string) (*os.File, int64, error) {
	ensureDataDir()

	path := filepath.Join(dataDir, name)

	var offset int64
	if st, err := os.Stat(path); err == nil {
		offset = st.Size()
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, 0, err
	}

	_, err = f.Seek(offset, 0)
	return f, offset, err
}

func OpenForRead(name string) (*os.File, int64, error) {
	path := filepath.Join(dataDir, name)
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	st, _ := f.Stat()
	return f, st.Size(), nil
}
