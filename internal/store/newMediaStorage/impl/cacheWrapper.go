package newmediastoragestorageimpl

import (
	"context"
	"io/fs"

	lru "github.com/hashicorp/golang-lru/v2"
	mediastorage "github.com/rishabhkailey/media-vault/internal/services/mediaStorage"
	newmediastorage "github.com/rishabhkailey/media-vault/internal/store/newMediaStorage"
)

type FileCacheWrapper struct {
	mediastorage.File
	stat *fs.FileInfo
}

func WrapFile(file mediastorage.File) FileCacheWrapper {
	return FileCacheWrapper{
		File: file,
	}
}

func (f FileCacheWrapper) Stat() (stat fs.FileInfo, err error) {
	if f.stat == nil {
		stat, err = f.File.Stat()
		if err != nil {
			return stat, err
		}
		f.stat = &stat
	}
	return *f.stat, err
}

var _ mediastorage.File = (*FileCacheWrapper)(nil)

type StoreCacheWrapper struct {
	cache *lru.ARCCache[string, FileCacheWrapper]
	store newmediastorage.Store
}

func NewStoreCacheWrapper(store newmediastorage.Store) (*StoreCacheWrapper, error) {
	// todo object size/memory usage?
	arcCache, err := lru.NewARC[string, FileCacheWrapper](1000)
	if err != nil {
		return nil, err
	}
	return &StoreCacheWrapper{
		store: store,
		cache: arcCache,
	}, err
}

var _ newmediastorage.Store = (*StoreCacheWrapper)(nil)

func (s *StoreCacheWrapper) SaveFile(ctx context.Context, cmd mediastorage.StoreSaveFileCmd) (int64, error) {
	return s.store.SaveFile(ctx, cmd)
}

func (s *StoreCacheWrapper) GetByFileName(ctx context.Context, fileName string) (file mediastorage.File, err error) {
	file, ok := s.cache.Get(fileName)
	if ok {
		return file, nil
	}
	// todo way to validate this?
	file, err = s.store.GetByFileName(ctx, fileName)
	if err != nil {
		return file, err
	}
	s.cache.Add(fileName, WrapFile(file))
	return file, err
}

func (s *StoreCacheWrapper) DeleteOne(ctx context.Context, fileName string) error {
	err := s.store.DeleteOne(ctx, fileName)
	if err != nil {
		return err
	}
	if s.cache.Contains(fileName) {
		s.cache.Remove(fileName)
	}
	return nil
}
