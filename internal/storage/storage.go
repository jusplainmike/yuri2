package storage

import (
	"io"
)

type Storage interface {
	Init(params interface{}) error
	Close()

	IsNotFoundErr(err error) bool

	Put(name string, r io.Reader, fileSize int64, contentType string) error
	ListAsync(cdone <-chan struct{}) <-chan *File
	Get(name string) (*File, error)
	Remove(name string) error
}
