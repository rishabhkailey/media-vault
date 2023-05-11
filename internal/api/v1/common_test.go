package v1_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rishabhkailey/media-service/internal/api"
	v1 "github.com/rishabhkailey/media-service/internal/api/v1"
	"github.com/rishabhkailey/media-service/internal/config"
	"github.com/rishabhkailey/media-service/internal/services"
	authserviceimpl "github.com/rishabhkailey/media-service/internal/services/authService/authServiceImpl"
	"github.com/rishabhkailey/media-service/internal/services/media/mediaimpl"
	mediametadataimpl "github.com/rishabhkailey/media-service/internal/services/mediaMetadata/mediaMetadataImpl"
	mediasearchimpl "github.com/rishabhkailey/media-service/internal/services/mediaSearch/mediaSearchimpl"
	"github.com/rishabhkailey/media-service/internal/services/mediaStorage/mediastorageimpl"
	"github.com/rishabhkailey/media-service/internal/services/uploadRequests/uploadrequestsimpl"
	usermediabindingsimpl "github.com/rishabhkailey/media-service/internal/services/userMediaBindings/userMediaBindingsimpl"
	"github.com/stretchr/testify/assert"
)

func NewTestServer() *v1.Server {
	return &v1.Server{
		Services: services.Services{
			Media:             mediaimpl.NewFakeService(),
			AuthService:       authserviceimpl.NewFakeService(),
			MediaMetadata:     mediametadataimpl.NewFakeService(),
			UserMediaBindings: usermediabindingsimpl.NewFakeService(),
			UploadRequests:    uploadrequestsimpl.NewFakeService(),
			MediaSearch:       mediasearchimpl.NewFakeService(),
			MediaStorage:      mediastorageimpl.NewFakeService(),
		},
	}
}

func NewTestRouter(v1ApiServer *v1.Server) (*gin.Engine, error) {
	return api.NewRouter(v1ApiServer, config.Config{
		Session: config.Session{
			Secret: "test",
		},
	})
}

// if only want to check that key exists in the response then set the value to nil
// for int int64 response use float64 type in expected response, we losses type when parsing json
func verifyMapResponse(t *testing.T, response httptest.ResponseRecorder, expectedResponse map[string]any) error {
	var responseBody map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
		return err
	}
	for key, expectedValue := range expectedResponse {
		value, ok := responseBody[key]
		assert.True(t, ok)
		if expectedValue != nil {
			assert.Equal(t, expectedValue, value)
		}
	}
	return nil
}

func verifyResponse(t *testing.T, response httptest.ResponseRecorder, expectedResponse any) error {
	var responseBody any
	if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
		return err
	}
	assert.Equal(t, expectedResponse, responseBody)
	return nil
}

func verifyResponseBytes(t *testing.T, response httptest.ResponseRecorder, expectedResponse []byte) error {
	responseBytes, _ := io.ReadAll(response.Body)
	assert.Equal(t, expectedResponse, responseBytes)
	return nil
}
