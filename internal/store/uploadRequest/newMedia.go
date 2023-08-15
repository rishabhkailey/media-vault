package uploadrequest

type Store interface {
	CreateUploadRequest(userID string) (uploadRequest UploadRequest, err error)
}
