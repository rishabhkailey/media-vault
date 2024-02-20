package newmediastoragestorageimpl

import (
	"context"
	"io/fs"
	"time"

	"github.com/minio/minio-go/v7"
	mediastorage "github.com/rishabhkailey/media-vault/internal/services/mediaStorage"
	newmediastorage "github.com/rishabhkailey/media-vault/internal/store/newMediaStorage"
	"github.com/rishabhkailey/media-vault/internal/utils"
)

type minioStore struct {
	cli        *minio.Client
	bucketName string
}

func NewMinioStore(cli *minio.Client) (newmediastorage.Store, error) {
	err := utils.CreateBucketIfMissing(context.Background(), *cli, newmediastorage.DIRECTORY_NAME)
	if err != nil {
		return nil, err
	}
	return &minioStore{
		cli:        cli,
		bucketName: newmediastorage.DIRECTORY_NAME,
	}, nil
}

var _ newmediastorage.Store = (*minioStore)(nil)

type MinioFileWrapper struct {
	fileName string
	minio.Object
}

func NewMinioFileWrapper(obj minio.Object, fileName string) *MinioFileWrapper {
	return &MinioFileWrapper{
		fileName: fileName,
		Object:   obj,
	}
}

var _ fs.File = (*MinioFileWrapper)(nil)
var _ mediastorage.File = (*MinioFileWrapper)(nil)

func (m *MinioFileWrapper) Stat() (fs.FileInfo, error) {
	stat, err := m.Object.Stat()
	if err != nil {
		return nil, err
	}
	return &MinioFileStat{
		name:    m.fileName,
		size:    stat.Size,
		modTime: stat.LastModified,
	}, nil
}

type MinioFileStat struct {
	name    string
	size    int64
	modTime time.Time
}

func (ms *MinioFileStat) Name() string {
	return ms.name
}

func (ms *MinioFileStat) Size() int64 {
	return ms.size
}

func (ms *MinioFileStat) ModTime() time.Time {
	return ms.modTime
}

func (ms *MinioFileStat) IsDir() bool {
	return false
}

func (ms *MinioFileStat) Sys() any {
	return nil
}

func (ms *MinioFileStat) Mode() fs.FileMode {
	return fs.ModePerm
}

var _ fs.FileInfo = (*MinioFileStat)(nil)

func (s *minioStore) SaveFile(ctx context.Context, cmd mediastorage.StoreSaveFileCmd) (int64, error) {
	uploadedInfo, err := s.cli.PutObject(ctx, s.bucketName, cmd.FileName, cmd.FileReader, cmd.FileSize, minio.PutObjectOptions{})
	return uploadedInfo.Size, err
}

func (s *minioStore) GetByFileName(ctx context.Context, fileName string) (mediastorage.File, error) {
	object, err := s.cli.GetObject(context.Background(), s.bucketName, fileName, minio.GetObjectOptions{})
	object.Stat()
	return NewMinioFileWrapper(*object, fileName), err
}

func (s *minioStore) DeleteOne(ctx context.Context, fileName string) error {
	return s.cli.RemoveObject(ctx, s.bucketName, fileName, minio.RemoveObjectOptions{})
}
