package api_test

import (
	"bytes"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// media list
// download media
// get media range
// get thumbnail

func TestMediaList(t *testing.T) {
	testFileNamePrefix := uuid.New().String()
	testFileType := strings.Split(uuid.New().String(), "-")[0]
	testClient, err := newTestHttpClient()
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	files, err := testClient.GenerateAndUploadTestFiles(t, 6, testFileNamePrefix, testFileType, 100, 10, time.Microsecond)
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	// wait for search indexing
	// time.Sleep(time.Second * 2)
	// sort by date desc
	sortedFilesDesc := make([]testFile, len(files))
	copy(sortedFilesDesc, files)
	sort.Slice(sortedFilesDesc, func(i, j int) bool {
		return sortedFilesDesc[i].date > sortedFilesDesc[j].date
	})
	// we can only test using desc order
	testCases := []struct {
		name                string
		requestQuery        url.Values
		exptectedStatusCode int
		expectedResponse    []testFile
	}{
		{
			name: "media list: desc: single file",
			requestQuery: url.Values{
				"page":    {"1"},
				"perPage": {"1"},
				"order":   {"date"},
				"sort":    {"desc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesDesc[0:1],
		},
		{
			name: "media list: desc: 3 files",
			requestQuery: url.Values{
				"page":    {"1"},
				"perPage": {"3"},
				"order":   {"date"},
				"sort":    {"desc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesDesc[0:3],
		},
		{
			name: "media list: desc: all files",
			requestQuery: url.Values{
				"page":    {"1"},
				"perPage": {"6"},
				"order":   {"date"},
				"sort":    {"desc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesDesc,
		},
		// pagination
		{
			name: "media list: desc: page=2: perpage=2",
			requestQuery: url.Values{
				"page":    {"1"},
				"perPage": {"2"},
				"order":   {"date"},
				"sort":    {"desc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesDesc[0:2],
		},
		{
			name: "media list: desc: page=2: perpage=2",
			requestQuery: url.Values{
				"page":    {"2"},
				"perPage": {"2"},
				"order":   {"date"},
				"sort":    {"desc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesDesc[2:4],
		},
		{
			name: "media list: desc: page=3: perpage=1",
			requestQuery: url.Values{
				"page":    {"3"},
				"perPage": {"2"},
				"order":   {"date"},
				"sort":    {"desc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesDesc[4:6],
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			testClient, err := newTestHttpClient()
			if err != nil {
				t.Fail()
				t.Error(err)
				return
			}

			resp, err := testClient.sendHttpRequest(httpRequest{
				method:      "GET",
				query:       test.requestQuery,
				url:         "/v1/mediaList",
				bearerToken: AUTH_TOKEN,
			}, false)
			if !assert.NoError(t, err, "search request") {
				return
			}

			respBody, ok := resp.body.([]byte)
			if !assert.True(t, ok, "type cast response body to []byte") {
				return
			}

			resultFiles, err := unmarshalSearchFilesResponse(bytes.NewReader(respBody))
			if !assert.NoError(t, err, "search request") {
				t.Error(err)
				return
			}
			if !assert.Equal(t, len(test.expectedResponse), len(resultFiles), "media list result length") {
				return
			}

			for index, expectedFile := range test.expectedResponse {
				resultFile := resultFiles[index]
				assert.Equal(t, expectedFile.name, resultFile.Name)
				assert.Equal(t, expectedFile.mediaType, resultFile.Type)
				assert.Equal(t, expectedFile.size, resultFile.Size)
				assert.True(t, resultFile.Thumbnail)
			}
		})
	}
}

func BenchmarkGetMedia(b *testing.B) {
	testClient, err := newTestHttpClient()
	if err != nil {
		b.Errorf("newTestHttpClient failed: %v", err)
		return
	}
	file := testFile{
		name:          "test.txt",
		date:          time.Now().UnixMilli(),
		data:          randomString(100_000_000),
		size:          100_000_000,
		thumbnialData: randomString(100),
		mediaType:     "txt",
	}
	resp, err := testClient.UploadTest(b, file, int64(1_000_000), AUTH_TOKEN)
	if err != nil || resp.status != http.StatusOK {
		b.Errorf("upload failed: err=%v, status=%d", err, resp.status)
		return
	}
	var mediaUrl string
	{
		body := resp.body.(map[string]any)
		value := body["url"]
		mediaUrl = value.(string)
	}

	for i := 0; i < b.N; i++ {
		testClient.GetMediaRangeTest(b, mediaUrl, AUTH_TOKEN, int64(1_000_000), file)
	}
}
