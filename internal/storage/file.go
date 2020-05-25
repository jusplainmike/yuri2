package storage

import (
	"io"
	"time"
)

type File struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size,omitempty"`
	ContentType string    `json:"contenttype,omitempty"`
	Modified    time.Time `json:"modified,omitempty"`
	ETag        string    `json:"etag,omitempty"`

	Err    error         `json:"-"`
	Reader io.ReadCloser `json:"-"`
}

func (f *File) Close() error {
	return f.Reader.Close()
}
