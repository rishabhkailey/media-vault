package v1_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	v1models "github.com/rishabhkailey/media-service/internal/api/v1/models"
	authservice "github.com/rishabhkailey/media-service/internal/services/authService"
	authserviceimpl "github.com/rishabhkailey/media-service/internal/services/authService/authServiceImpl"
	"github.com/rishabhkailey/media-service/internal/services/media/mediaimpl"
	"github.com/rishabhkailey/media-service/internal/services/mediaStorage/mediastorageimpl"
	"github.com/stretchr/testify/assert"
)

// RefreshSession and UserAuthMiddleware test
func TestMediaList(t *testing.T) {
	testMediaListStruct := randomMediaList(10)

	var testMediaListResponseMap any
	{
		testMediaListResponse, err := v1models.NewGetMediaListResponse(testMediaListStruct)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, err := json.Marshal(testMediaListResponse)
		if err != nil {
			t.Error(err)
			t.Fail() // todo add t.Fail() at other missing places also
			return
		}
		err = json.Unmarshal(bytes, &testMediaListResponseMap)
		if err != nil {
			t.Error(err)
			t.Fail() // todo add t.Fail() at other missing places also
			return
		}
	}

	testCases := []struct {
		name               string
		requestQuery       url.Values
		expectedStatusCode int
		expectedResponse   any
		authService        authserviceimpl.FakeService
		mediaService       mediaimpl.FakeService
	}{
		{
			name: "normal",
			requestQuery: url.Values{
				"order":    {"date"},
				"sort":     {"asc"},
				"per_page": {"10"},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   testMediaListResponseMap,
			authService: authserviceimpl.FakeService{
				ExpectedError:             nil,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     nil,
			},
		},
		{
			name: "bad request: order missing",
			requestQuery: url.Values{
				"sort":     {"asc"},
				"per_page": {"10"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
			authService: authserviceimpl.FakeService{
				ExpectedError:             nil,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     nil,
			},
		},
		{
			name: "bad request: sort missing",
			requestQuery: url.Values{
				"order":    {"date"},
				"per_page": {"10"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
			authService: authserviceimpl.FakeService{
				ExpectedError:             nil,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     nil,
			},
		},
		{
			name: "bad request: per_page missing",
			requestQuery: url.Values{
				"sort":  {"asc"},
				"order": {"date"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
			authService: authserviceimpl.FakeService{
				ExpectedError:             nil,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     nil,
			},
		},
		{
			name: "bad request: negative perPage missing",
			requestQuery: url.Values{
				"sort":     {"asc"},
				"order":    {"date"},
				"per_page": {"-1"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
			authService: authserviceimpl.FakeService{
				ExpectedError:             nil,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     nil,
			},
		},
		{
			name: "bad request: invalid sort value",
			requestQuery: url.Values{
				"sort":     {"ascaa"},
				"order":    {"date"},
				"per_page": {"-1"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
			authService: authserviceimpl.FakeService{
				ExpectedError:             nil,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     nil,
			},
		},
		{
			name: "bad request: invalid order value",
			requestQuery: url.Values{
				"sort":     {"asc"},
				"order":    {"date_abc"},
				"per_page": {"-1"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
			authService: authserviceimpl.FakeService{
				ExpectedError:             nil,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     nil,
			},
		},
		{
			name: "bad request: invalid last_media_id value",
			requestQuery: url.Values{
				"sort":          {"asc"},
				"order":         {"date_abc"},
				"last_media_id": {"-1"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
			authService: authserviceimpl.FakeService{
				ExpectedError:             nil,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     nil,
			},
		},
		{
			name: "bad request: invalid last_date",
			requestQuery: url.Values{
				"sort":      {"asc"},
				"order":     {"date_abc"},
				"last_date": {"abc"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
			authService: authserviceimpl.FakeService{
				ExpectedError:             nil,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     nil,
			},
		},
		{
			name: "unathorized error",
			requestQuery: url.Values{
				"sort":     {"asc"},
				"order":    {"date"},
				"per_page": {"1"},
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   nil,
			authService: authserviceimpl.FakeService{
				ExpectedError:             authservice.ErrUnauthorized,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     nil,
			},
		},
		{
			name: "auth error",
			requestQuery: url.Values{
				"sort":     {"asc"},
				"order":    {"date"},
				"per_page": {"1"},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
			authService: authserviceimpl.FakeService{
				ExpectedError:             errors.New("auth error"),
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     nil,
			},
		},
		{
			name: "get media error",
			requestQuery: url.Values{
				"sort":     {"asc"},
				"order":    {"date"},
				"per_page": {"1"},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
			authService: authserviceimpl.FakeService{
				ExpectedError:             nil,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     errors.New("get media error"),
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
		server.Services.AuthService = &test.authService
		server.Services.Media = &test.mediaService
		url := url.URL{
			Path:     "/v1/mediaList",
			RawQuery: test.requestQuery.Encode(),
		}
		request, _ := http.NewRequest("GET", url.String(), strings.NewReader(test.requestQuery.Encode()))
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		if test.expectedResponse != nil {
			if err := verifyResponse(t, *recorder, test.expectedResponse, test.name); err != nil {
				t.Error(err)
				t.Fail()
			}
		}
		assert.Equal(t, test.expectedStatusCode, recorder.Code, test.name)
	}
}

// includes get media, media range and thumbnail tests
// todo separate these tests?
func TestGetMedia(t *testing.T) {
	testFileData := []byte(randomString(10000))
	tests := []struct {
		name             string
		url              string
		expectedStatus   int
		expectedError    error
		expectedResponse []byte
		Range            string
		authService      authserviceimpl.FakeService
		mediaStorage     mediastorageimpl.FakeService
	}{
		{
			name:             "Download",
			url:              "/v1/media/test-media-file",
			expectedStatus:   http.StatusOK,
			expectedResponse: testFileData,
			authService: authserviceimpl.FakeService{
				ExpectedError:             nil,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaStorage: mediastorageimpl.FakeService{
				ExpectedError:        nil,
				ExpectedWrittenBytes: int64(len(testFileData)),
				FileBytes:            testFileData,
			},
			// Range: bytes=0-1000
			Range: "",
		},
		{
			name:             "Thumbnail",
			url:              "/v1/thumbnail/test-media-file",
			expectedStatus:   http.StatusOK,
			expectedResponse: testFileData,
			authService: authserviceimpl.FakeService{
				ExpectedError:             nil,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaStorage: mediastorageimpl.FakeService{
				ExpectedError:        nil,
				ExpectedWrittenBytes: int64(len(testFileData)),
				FileBytes:            testFileData,
			},
			// Range: bytes=0-1000
			Range: "",
		},
		{
			name:             "Normal Range",
			url:              "/v1/media/test-media-file",
			expectedStatus:   http.StatusPartialContent,
			expectedResponse: testFileData,
			authService: authserviceimpl.FakeService{
				ExpectedError:             nil,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaStorage: mediastorageimpl.FakeService{
				ExpectedError:        nil,
				ExpectedWrittenBytes: int64(len(testFileData)),
				FileBytes:            testFileData,
			},
			Range: fmt.Sprintf("bytes=0-%d", len(testFileData)),
		},
		{
			name:             "unauthorized",
			url:              "/v1/media/test-media-file",
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: nil,
			authService: authserviceimpl.FakeService{
				ExpectedError:             authservice.ErrUnauthorized,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaStorage: mediastorageimpl.FakeService{
				ExpectedError:        nil,
				ExpectedWrittenBytes: int64(len(testFileData)),
				FileBytes:            testFileData,
			},
		},
		{
			name:             "auth error",
			url:              "/v1/media/test-media-file",
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: nil,
			authService: authserviceimpl.FakeService{
				ExpectedError:             errors.New("auth error"),
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaStorage: mediastorageimpl.FakeService{
				ExpectedError:        nil,
				ExpectedWrittenBytes: int64(len(testFileData)),
				FileBytes:            testFileData,
			},
		},
		{
			name:             "media storage error",
			url:              "/v1/media/test-media-file",
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: nil,
			authService: authserviceimpl.FakeService{
				ExpectedError:             nil,
				ExpectedUserID:            "1",
				ExpectedSessionExpireTime: time.Now().Add(time.Hour * 10).Unix(),
			},
			mediaStorage: mediastorageimpl.FakeService{
				ExpectedError:        errors.New("media storage error"),
				ExpectedWrittenBytes: int64(len(testFileData)),
				FileBytes:            testFileData,
			},
		},
	}
	server := NewTestServer()
	router, err := NewTestRouter(server)
	if err != nil {
		t.Error(err)
		return
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server.Services.AuthService = &test.authService
			server.Services.MediaStorage = &test.mediaStorage
			request, _ := http.NewRequest("GET", test.url, nil)
			if len(test.Range) != 0 {
				request.Header.Add("Range", test.Range)
			}
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			assert.Equal(t, test.expectedStatus, recorder.Code)
			if test.expectedResponse != nil {
				if err := verifyResponseBytes(t, *recorder, test.expectedResponse); err != nil {
					t.Error(err)
					t.Fail()
				}
			}
			if len(test.Range) != 0 {
				// this is just to verify that the request was forwarded to range handler
				// this is a unit controller/server test so we don't care about the service logic
				expectedRangeHeader := fmt.Sprintf("%s/%d", test.Range, len(test.expectedResponse))
				rangeHeader := recorder.Header().Get("Range")
				assert.Equal(t, expectedRangeHeader, rangeHeader)
			}
			// router.ServeHTTP()
		})
	}
}
