package v1_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authservice "github.com/rishabhkailey/media-vault/internal/services/authService"
	authserviceimpl "github.com/rishabhkailey/media-vault/internal/services/authService/authServiceImpl"
	"github.com/stretchr/testify/assert"
)

// RefreshSession and UserAuthMiddleware test
func TestRefreshSession(t *testing.T) {
	testCases := []struct {
		name             string
		expectedStatus   int
		expectedResponse map[string]any
		authService      authserviceimpl.FakeService
	}{
		{
			name:           "normal",
			expectedStatus: 200,
			expectedResponse: map[string]any{
				"expires": float64(1683518020),
			},
			authService: authserviceimpl.FakeService{
				ExpectedError:             nil,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: 1683518020,
			},
		},
		{
			name:           "internal server error",
			expectedStatus: 500,
			expectedResponse: map[string]any{
				"error": "Internal server error",
			},
			authService: authserviceimpl.FakeService{
				ExpectedError: errors.New("internal server error"),
			},
		},
		{
			name:             "unathorized",
			expectedStatus:   401,
			expectedResponse: nil,
			authService: authserviceimpl.FakeService{
				ExpectedError: authservice.ErrUnauthorized,
			},
		},
	}
	server := NewTestServer()
	router, err := NewTestRouter(server)
	if err != nil {
		t.Error(err)
		return
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			server.Services.AuthService = &test.authService
			request, _ := http.NewRequest("POST", "/v1/refreshSession", nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			if err := verifyMapResponse(t, *recorder, test.expectedResponse); err != nil {
				t.Error(err)
			}
			assert.Equal(t, test.expectedStatus, recorder.Code)
		})
	}
}
