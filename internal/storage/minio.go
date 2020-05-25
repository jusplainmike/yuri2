package storage

import (
	"errors"
	"io"

	"github.com/minio/minio-go/v6"
)

const bucketName = "yuriv3"

type Minio struct {
	client *minio.Client
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	Secret    string
	Location  string
	UseSSL    bool
}

func (m *Minio) Init(params interface{}) (err error) {
	cfg, ok := params.(*MinioConfig)
	if !ok {
		return errors.New("invalid parameter type")
	}

	m.client, err = minio.New(
		cfg.Endpoint, cfg.AccessKey, cfg.Secret, cfg.UseSSL)
	if err != nil {
		return
	}

	ok, err = m.client.BucketExists(bucketName)
	if err != nil {
		return
	}
	if !ok {
		if err = m.client.MakeBucket(bucketName, cfg.Location); err != nil {
			return
		}
	}

	return
}

func (m *Minio) Close() {}

func (m *Minio) IsNotFoundErr(err error) bool {
	return err != nil && err.Error() == "The specified key does not exist."
}

func (m *Minio) Put(name string, r io.Reader, fileSize int64, contentType string) (err error) {
	_, err = m.client.PutObject(bucketName, name, r, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return
}

func (m *Minio) ListAsync(cdone <-chan struct{}) <-chan *File {
	cout := make(chan *File)
	c := m.client.ListObjects(bucketName, "", false, cdone)

	go func() {
		for obj := range c {
			cout <- fileFromMinioObjectInfo(&obj)
		}
		close(cout)
	}()

	return cout
}

func (m *Minio) Get(name string) (*File, error) {
	obj, err := m.client.GetObject(bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	info, err := obj.Stat()
	if err != nil {
		return nil, err
	}

	f := fileFromMinioObjectInfo(&info)
	f.Reader = obj

	return f, nil
}

func (m *Minio) Remove(name string) error {
	return m.client.RemoveObject(bucketName, name)
}

func fileFromMinioObjectInfo(info *minio.ObjectInfo) *File {
	return &File{
		Name:        info.Key,
		ContentType: info.ContentType,
		ETag:        info.ETag,
		Modified:    info.LastModified,
		Size:        info.Size,
		Err:         info.Err,
	}
}
