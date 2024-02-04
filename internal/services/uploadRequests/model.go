package uploadrequests

type Status string

const (
	COMPLETED_UPLOAD_STATUS   Status = "completed"
	FAILED_UPLOAD_STATUS      Status = "failed"
	IN_PROGRESS_UPLOAD_STATUS Status = "inProgress"
)

// will be used by upload request service
type CreateUploadRequestCommand struct {
	UserID string
}

type GetByIDQuery struct {
	ID string
}

type UpdateStatusCommand struct {
	Status Status
	ID     string
}

type DeleteOneCommand struct {
	ID string
}

type DeleteManyCommand struct {
	IDs []string
}
