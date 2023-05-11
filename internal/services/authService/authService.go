package authservice

type Service interface {
	// GetSession(c *gin.Context) (UserSession, error)
	ValidateUserAccess(query ValidateUserAccessQuery, scopes []string) (userID string, err error)
	TerminateSession(cmd TerminateSessionCmd) error
	ValidateUserMediaAccess(query ValidateUserMediaAccessQuery) error
	GetSessionExpireTime(query GetSessionExpireTimeQuery) (int64, error)
}
