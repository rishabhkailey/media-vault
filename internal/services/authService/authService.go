package authservice

type Service interface {
	ValidateUserAccess(query ValidateUserAccessQuery, scopes []string) (userID string, err error)
	TerminateSession(cmd TerminateSessionCmd) error
	ValidateUserMediaAccess(query ValidateUserMediaAccessQuery) error
	GetSessionExpireTime(query GetSessionExpireTimeQuery) (int64, error)
	RefreshSession(RefreshSessionQuery) (expireTime int64, err error)
}
