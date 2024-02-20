package mediastorageimpl

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	mediastorage "github.com/rishabhkailey/media-vault/internal/services/mediaStorage"
)

type onGoingUploadRequest struct {
	Reader     io.ReadCloser
	Writer     io.WriteCloser
	err        error // in case of failure
	completed  bool
	index      int64
	size       int64
	ctx        context.Context
	cancelFunc context.CancelFunc
	userID     string
	// checksum   string
}

// todo session affinity required till all the browsers support http/2 protocol (which support stream upload)
// https://caniuse.com/http2, right now android browser's don't have good support
// requestID -> uploadRequest
type onGoingUploadRequestsStore struct {
	requests map[string]*onGoingUploadRequest
	mutex    sync.Mutex
}

func newOnGoingUploadRequestsStore() onGoingUploadRequestsStore {
	return onGoingUploadRequestsStore{
		requests: make(map[string]*onGoingUploadRequest),
		mutex:    sync.Mutex{},
	}
}

func (s *onGoingUploadRequestsStore) add(ctx context.Context, cancelFunc context.CancelFunc, cmd mediastorage.InitChunkUploadCmd) (*onGoingUploadRequest, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.requests[cmd.RequestID]; ok {
		return nil, fmt.Errorf("request with ID %v already exist", cmd.RequestID)
	}
	reader, writer := io.Pipe()
	uploadRequst := &onGoingUploadRequest{
		Reader:     reader,
		Writer:     writer,
		err:        nil,
		completed:  false,
		index:      0,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		size:       cmd.FileSize,
		userID:     cmd.UserID,
	}
	s.requests[fmt.Sprintf("%s:%s", cmd.RequestID, cmd.UserID)] = uploadRequst
	return uploadRequst, nil
}

func (s *onGoingUploadRequestsStore) get(requestID, userID string) (*onGoingUploadRequest, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	uploadRequest, ok := s.requests[fmt.Sprintf("%s:%s", requestID, userID)]
	if !ok {
		return nil, fmt.Errorf("request with ID %v doesn't Exist", requestID)
	}
	return uploadRequest, nil
}

// todo instead of this add expire in get and immediately call this
// todo max size of the store if it reaches we return some status code to mention to many upload requests in progress
func (s *onGoingUploadRequestsStore) deleteUploadRequestAfter(t time.Duration, requestID, userID string) error {
	key := fmt.Sprintf("%s:%s", requestID, userID)
	if _, ok := s.requests[key]; !ok {
		return fmt.Errorf("request with ID %v doesn't Exist", requestID)
	}
	go func() {
		<-time.NewTicker(t).C
		s.mutex.Lock()
		defer s.mutex.Unlock()
		delete(s.requests, key)
	}()
	return nil
}
