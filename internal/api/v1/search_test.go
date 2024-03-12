package v1_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	authservice "github.com/rishabhkailey/media-vault/internal/services/authService"
	authserviceimpl "github.com/rishabhkailey/media-vault/internal/services/authService/authServiceImpl"
	"github.com/rishabhkailey/media-vault/internal/services/media/mediaimpl"
	mediasearchimpl "github.com/rishabhkailey/media-vault/internal/services/mediaSearch/mediaSearchimpl"
	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	testMediaIds := []uint{1, 2}
	testMediaListStruct := []storemodels.MediaModel{
		// {
		// 	MediaUrl:     "/test/abc",
		// 	ThumbnailUrl: "/test/def",
		// 	Metadata: mediametadata.Metadata{
		// 		Name:      "test",
		// 		Date:      time.Now(),
		// 		Type:      "txt",
		// 		Size:      100,
		// 		Thumbnail: true,
		// 	},
		// },
		// {
		// 	MediaUrl:     "/test/abc",
		// 	ThumbnailUrl: "/test/def",
		// 	Metadata: mediametadata.Metadata{
		// 		Name:      "test",
		// 		Date:      time.Now(),
		// 		Type:      "txt",
		// 		Size:      100,
		// 		Thumbnail: true,
		// 	},
		// },
	}
	var testMediaListMap any
	{
		bytes, err := json.Marshal(testMediaListStruct)
		if err != nil {
			t.Error(err)
			t.Fail() // todo add t.Fail() at other missing places also
			return
		}
		err = json.Unmarshal(bytes, &testMediaListMap)
		if err != nil {
			t.Error(err)
			t.Fail() // todo add t.Fail() at other missing places also
			return
		}
	}

	testCases := []struct {
		name               string
		authService        authserviceimpl.FakeService
		mediaService       mediaimpl.FakeService
		mediaSearchService mediasearchimpl.FakeService
		requestQuery       url.Values
		expectedStatus     int
		expectedResponse   any
	}{
		{
			name: "normal",
			authService: authserviceimpl.FakeService{
				ExpectedError:  nil,
				ExpectedUserID: "1",
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     nil,
			},
			mediaSearchService: mediasearchimpl.FakeService{
				ExpectedMediaListIDs: testMediaIds,
				ExpectedError:        nil,
			},
			requestQuery: url.Values{
				"query":    {"test query"},
				"page":     {"1"},
				"per_page": {"30"},
				"order":    {"date"},
				"sort":     {"desc"},
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: testMediaListMap,
		},
		{
			name: "auth error",
			authService: authserviceimpl.FakeService{
				ExpectedError:  authservice.ErrUnauthorized,
				ExpectedUserID: "1",
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     nil,
			},
			mediaSearchService: mediasearchimpl.FakeService{
				ExpectedMediaListIDs: testMediaIds,
				ExpectedError:        nil,
			},
			requestQuery: url.Values{
				"query":    {"test query"},
				"page":     {"1"},
				"per_page": {"30"},
				"order":    {"date"},
				"sort":     {"desc"},
			},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: nil,
		},
		{
			name: "media error",
			authService: authserviceimpl.FakeService{
				ExpectedError:  nil,
				ExpectedUserID: "1",
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     errors.New("media error"),
			},
			mediaSearchService: mediasearchimpl.FakeService{
				ExpectedMediaListIDs: testMediaIds,
				ExpectedError:        nil,
			},
			requestQuery: url.Values{
				"query":    {"test query"},
				"page":     {"1"},
				"per_page": {"30"},
				"order":    {"date"},
				"sort":     {"desc"},
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: nil,
		},
		{
			name: "bad request: empty query",
			authService: authserviceimpl.FakeService{
				ExpectedError:  nil,
				ExpectedUserID: "1",
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     errors.New("media error"),
			},
			mediaSearchService: mediasearchimpl.FakeService{
				ExpectedMediaListIDs: testMediaIds,
				ExpectedError:        nil,
			},
			requestQuery: url.Values{
				"query":    {""},
				"page":     {"1"},
				"per_page": {"30"},
				"order":    {"date"},
				"sort":     {"desc"},
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: nil,
		},
		{
			name: "bad request: negative page",
			authService: authserviceimpl.FakeService{
				ExpectedError:  nil,
				ExpectedUserID: "1",
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     errors.New("media error"),
			},
			mediaSearchService: mediasearchimpl.FakeService{
				ExpectedMediaListIDs: testMediaIds,
				ExpectedError:        nil,
			},
			requestQuery: url.Values{
				"query":    {"test"},
				"page":     {"-1"},
				"per_page": {"30"},
				"order":    {"date"},
				"sort":     {"desc"},
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: nil,
		},
		{
			name: "bad request: missing page",
			authService: authserviceimpl.FakeService{
				ExpectedError:  nil,
				ExpectedUserID: "1",
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     errors.New("media error"),
			},
			mediaSearchService: mediasearchimpl.FakeService{
				ExpectedMediaListIDs: testMediaIds,
				ExpectedError:        nil,
			},
			requestQuery: url.Values{
				"query":    {"test"},
				"per_page": {"30"},
				"order":    {"date"},
				"sort":     {"desc"},
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: nil,
		},
		{
			name: "bad request: negative per_page",
			authService: authserviceimpl.FakeService{
				ExpectedError:  nil,
				ExpectedUserID: "1",
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     errors.New("media error"),
			},
			mediaSearchService: mediasearchimpl.FakeService{
				ExpectedMediaListIDs: testMediaIds,
				ExpectedError:        nil,
			},
			requestQuery: url.Values{
				"query":    {"test"},
				"page":     {"1"},
				"per_page": {"-30"},
				"order":    {"date"},
				"sort":     {"desc"},
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: nil,
		},
		{
			name: "bad request: missing per_page",
			authService: authserviceimpl.FakeService{
				ExpectedError:  nil,
				ExpectedUserID: "1",
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     errors.New("media error"),
			},
			mediaSearchService: mediasearchimpl.FakeService{
				ExpectedMediaListIDs: testMediaIds,
				ExpectedError:        nil,
			},
			requestQuery: url.Values{
				"query": {""},
				"page":  {"1"},
				"order": {"date"},
				"sort":  {"desc"},
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: nil,
		},
		{
			name: "bad request: invalid order",
			authService: authserviceimpl.FakeService{
				ExpectedError:  nil,
				ExpectedUserID: "1",
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     errors.New("media error"),
			},
			mediaSearchService: mediasearchimpl.FakeService{
				ExpectedMediaListIDs: testMediaIds,
				ExpectedError:        nil,
			},
			requestQuery: url.Values{
				"query":    {"test"},
				"page":     {"1"},
				"per_page": {"30"},
				"order":    {"date_test"},
				"sort":     {"desc"},
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: nil,
		},
		{
			name: "bad request: missing order",
			authService: authserviceimpl.FakeService{
				ExpectedError:  nil,
				ExpectedUserID: "1",
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     errors.New("media error"),
			},
			mediaSearchService: mediasearchimpl.FakeService{
				ExpectedMediaListIDs: testMediaIds,
				ExpectedError:        nil,
			},
			requestQuery: url.Values{
				"query":    {"test"},
				"page":     {"1"},
				"per_page": {"30"},
				"sort":     {"desc"},
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: nil,
		},
		{
			name: "bad request: invalid sort",
			authService: authserviceimpl.FakeService{
				ExpectedError:  nil,
				ExpectedUserID: "1",
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     errors.New("media error"),
			},
			mediaSearchService: mediasearchimpl.FakeService{
				ExpectedMediaListIDs: testMediaIds,
				ExpectedError:        nil,
			},
			requestQuery: url.Values{
				"query":    {"test"},
				"page":     {"1"},
				"per_page": {"30"},
				"order":    {"date"},
				"sort":     {"desc_test"},
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: nil,
		},
		{
			name: "bad request: missing order",
			authService: authserviceimpl.FakeService{
				ExpectedError:  nil,
				ExpectedUserID: "1",
			},
			mediaService: mediaimpl.FakeService{
				ExpectedMediaList: testMediaListStruct,
				ExpectedError:     errors.New("media error"),
			},
			mediaSearchService: mediasearchimpl.FakeService{
				ExpectedMediaListIDs: testMediaIds,
				ExpectedError:        nil,
			},
			requestQuery: url.Values{
				"query":    {"test"},
				"page":     {"1"},
				"per_page": {"30"},
				"order":    {"date"},
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: nil,
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
			server.Services.Media = &test.mediaService
			url := url.URL{
				Path:     "/v1/search",
				RawQuery: test.requestQuery.Encode(),
			}
			request, _ := http.NewRequest("GET", url.String(), strings.NewReader(test.requestQuery.Encode()))
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			if test.expectedResponse != nil {
				if err := verifyResponse(t, *recorder, test.expectedResponse); err != nil {
					t.Error(err)
					t.Fail()
				}
			}
			assert.Equal(t, test.expectedStatus, recorder.Code)
		})
	}
}
