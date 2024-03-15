package api_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 401 vs 403
// 403 when token is correct but the user doesn't have access to the resource.
// 		e.g. file doesn't belong to user
// 401 invalid, missing or expired token

func TestNormalEndpointsAuth(t *testing.T) {
	testCases := []struct {
		name               string
		requestUri         string
		requestBody        any
		requestQuery       url.Values
		requestMethod      string
		requestBearerToken string
		expectedStatusCode int
	}{
		// we don't want to upload any actual file
		// 400 status code for the correct bearer token and 401 for incorrect
		{
			name:       "correct bearer token: initChunkUpload",
			requestUri: "/v1/upload",
			requestBody: map[string]any{
				"fileName123": "test.txt",
				"size":        100,
				"mediaType":   "txt",
				"date":        time.Now().UnixMilli(),
			},
			requestMethod:      "POST",
			requestBearerToken: AUTH_TOKEN,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:       "wrong bearer token: initChunkUpload",
			requestUri: "/v1/upload",
			requestBody: map[string]any{
				"fileName":  "test.txt",
				"size":      100,
				"mediaType": "txt",
				"date":      time.Now().UnixMilli(),
			},
			requestMethod:      "POST",
			requestBearerToken: "blabla",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:       "correct bearer token: search",
			requestUri: "/v1/search",
			requestQuery: url.Values{
				"order":    {"date"},
				"sort":     {"desc"},
				"page":     {"1"},
				"per_page": {"10"},
				"query":    {"test"},
			},
			requestMethod:      "GET",
			requestBearerToken: AUTH_TOKEN,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:       "wrong bearer token: search",
			requestUri: "/v1/search",
			requestQuery: url.Values{
				"order":    {"date"},
				"sort":     {"desc"},
				"page":     {"1"},
				"per_page": {"10"},
				"query":    {"test"},
			},
			requestMethod:      "GET",
			requestBearerToken: "blabla",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:       "correct bearer token: media list",
			requestUri: "/v1/media-list",
			requestQuery: url.Values{
				"order":    {"date"},
				"sort":     {"desc"},
				"page":     {"1"},
				"per_page": {"10"},
			},
			requestMethod:      "GET",
			requestBearerToken: AUTH_TOKEN,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:       "wrong bearer token: media list",
			requestUri: "/v1/media-list",
			requestQuery: url.Values{
				"order":    {"date"},
				"sort":     {"desc"},
				"page":     {"1"},
				"per_page": {"10"},
			},
			requestMethod:      "GET",
			requestBearerToken: "blabla",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "correct bearer token: upload chunk",
			requestUri:         "/v1/upload/invalid-upload-request-id/chunk",
			requestBody:        nil,
			requestMethod:      "POST",
			requestBearerToken: AUTH_TOKEN,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "wrong bearer token: upload chunk",
			requestUri:         "/v1/upload/invalid-upload-request-id/chunk",
			requestBody:        nil,
			requestMethod:      "POST",
			requestBearerToken: "blabla",
			expectedStatusCode: http.StatusUnauthorized,
		},
		// todo(fix): status code should be 400 or 404
		{
			name:               "correct bearer token: upload finish",
			requestUri:         "/v1/upload/invalid-upload-request-id/finish",
			requestBody:        nil,
			requestMethod:      "POST",
			requestBearerToken: AUTH_TOKEN,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "wrong bearer token: upload finish",
			requestUri:         "/v1/upload/invalid-upload-request-id/finish",
			requestBody:        nil,
			requestMethod:      "POST",
			requestBearerToken: "blabla",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// new test client for each request no session persistence between tests
			testClient, err := newTestHttpClient()
			if err != nil {
				t.Error(err)
				t.Fail()
				return
			}
			resp, err := testClient.sendHttpRequest(httpRequest{
				url:         test.requestUri,
				method:      test.requestMethod,
				query:       test.requestQuery,
				body:        test.requestBody,
				bearerToken: test.requestBearerToken,
			}, false)
			if !assert.NoError(t, err, test.name) {
				return
			}
			assert.Equal(t, test.expectedStatusCode, resp.status, test.name)
		})
	}
}

func TestMediaEndpointsAuth(t *testing.T) {
	testCases := []struct {
		name                             string
		requestUri                       string
		requestMethod                    string
		requestBearerToken               string
		refreshSession                   bool
		expectedGetMediaStatusCode       int
		expectedRefreshSessionStatusCode int
	}{
		{
			name:                             "Get Media: with refreshSession: correct auth token",
			requestUri:                       "/v1/file/test123321",
			requestMethod:                    "GET",
			requestBearerToken:               AUTH_TOKEN,
			refreshSession:                   true,
			expectedRefreshSessionStatusCode: http.StatusOK,
			expectedGetMediaStatusCode:       http.StatusForbidden, // file doesn't exist
		},
		{
			name:                             "Get Media: with refreshSession: empty auth token",
			requestUri:                       "/v1/file/test123321",
			requestMethod:                    "GET",
			requestBearerToken:               "",
			refreshSession:                   true,
			expectedRefreshSessionStatusCode: http.StatusUnauthorized,
			expectedGetMediaStatusCode:       http.StatusUnauthorized,
		},
		{
			name:                             "Get Media: with refreshSession: incorrect auth token",
			requestUri:                       "/v1/media/test123321",
			requestMethod:                    "GET",
			requestBearerToken:               "wrong auth token",
			refreshSession:                   true,
			expectedRefreshSessionStatusCode: http.StatusUnauthorized,
			expectedGetMediaStatusCode:       http.StatusUnauthorized,
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			testClient, err := newTestHttpClient()
			if err != nil {
				t.Error(err)
				t.Fail()
				return
			}
			if test.refreshSession {
				resp, err := testClient.sendRefreshSessionRequest(test.requestBearerToken)
				if !assert.NoError(t, err, test.name, "refresh token request") {
					return
				}
				if !assert.Equal(t, test.expectedRefreshSessionStatusCode, resp.status, test.name, "refresh token request") {
					return
				}
			}
			resp, err := testClient.sendHttpRequest(httpRequest{
				method: test.requestMethod,
				url:    test.requestUri,
			}, false)
			if !assert.NoError(t, err, test.name) {
				return
			}
			assert.Equal(t, test.expectedGetMediaStatusCode, resp.status, test.name)
		})
	}
}
