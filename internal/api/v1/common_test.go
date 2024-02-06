package v1_test

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"

	cryptorand "crypto/rand"

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
	uploadrequests "github.com/rishabhkailey/media-service/internal/services/uploadRequests"
	"github.com/rishabhkailey/media-service/internal/services/uploadRequests/uploadrequestsimpl"
	usermediabindingsimpl "github.com/rishabhkailey/media-service/internal/services/userMediaBindings/userMediaBindingsimpl"
	storemodels "github.com/rishabhkailey/media-service/internal/store/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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

func verifyResponse(t *testing.T, response httptest.ResponseRecorder, expectedResponse any, msgAndArgs ...interface{}) error {
	var responseBody any
	if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
		return err
	}
	assert.Equal(t, expectedResponse, responseBody, msgAndArgs)
	return nil
}

func verifyResponseBytes(t *testing.T, response httptest.ResponseRecorder, expectedResponse []byte) error {
	responseBytes, _ := io.ReadAll(response.Body)
	assert.Equal(t, expectedResponse, responseBytes)
	return nil
}

// it will panic if any error so only use this in tests
func randomString(n int64) string {
	bytes := make([]byte, n)
	cryptorand.Read(bytes)
	var randString string
	{
		// to make string readable
		ascii := []rune("abcdefghijklmnopqrstuvwxyz ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		for _, byte := range bytes {
			randString += string(ascii[int(byte)%len(ascii)])
		}
	}
	return randString
}

func randomMediaList(n int) (mediaList []storemodels.MediaModel) {
	for i := 0; i < n; i++ {
		mediaList = append(mediaList, storemodels.MediaModel{
			Model: gorm.Model{
				ID:        uint(rand.Uint32()),
				CreatedAt: time.Now().AddDate(0, 0, -1*rand.Intn(10)),
				UpdatedAt: time.Now().AddDate(0, 0, -1*rand.Intn(10)),
			},
			FileName:        randomString(10),
			UploadRequestID: randomString(10),
			MetadataID:      uint(rand.Uint32()),
			UploadRequest: storemodels.UploadRequestsModel{
				ID:        randomString(10),
				Status:    string(uploadrequests.COMPLETED_UPLOAD_STATUS),
				CreatedAt: time.Now().AddDate(0, 0, -1*rand.Intn(10)),
				UpdatedAt: time.Now().AddDate(0, 0, -1*rand.Intn(10)),
			},
			Metadata: storemodels.MediaMetadataModel{
				Model: gorm.Model{
					ID:        uint(rand.Uint32()),
					CreatedAt: time.Now().AddDate(0, 0, -1*rand.Intn(10)),
					UpdatedAt: time.Now().AddDate(0, 0, -1*rand.Intn(10)),
				},
				MediaMetadata: storemodels.MediaMetadata{
					Name:      randomString(10),
					Date:      time.Now().AddDate(0, 0, -1*rand.Intn(10)),
					Type:      randomString(5),
					Size:      rand.Uint64(),
					Thumbnail: false,
				},
			},
		})
	}
	return
}
