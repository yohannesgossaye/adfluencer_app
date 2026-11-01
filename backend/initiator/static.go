package initiator

import (
	"io/fs"
	"net/http"
)

// filesOnlyFS wraps an http.FileSystem so that directories cannot be opened.
// http.FileServer would otherwise render a listing of every uploaded file.
type filesOnlyFS struct {
	fs http.FileSystem
}

func (f filesOnlyFS) Open(name string) (http.File, error) {
	file, err := f.fs.Open(name)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	if info.IsDir() {
		_ = file.Close()
		return nil, fs.ErrNotExist
	}

	return file, nil
}
