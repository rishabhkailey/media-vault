package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// todo test pagination

// upload few files and search for those files sorted by date, ideally we should get those files back
func TestSearch(t *testing.T) {
	testFileNamePrefix := uuid.New().String()
	testFileType := strings.Split(uuid.New().String(), "-")[0]
	testClient, err := newTestHttpClient()
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	files, err := testClient.GenerateAndUploadTestFiles(t, 6, testFileNamePrefix, testFileType, 100, 10, time.Second*10)
	if err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	// wait for search indexing
	time.Sleep(time.Second * 2)
	// sort by date desc
	sortedFilesDesc := make([]testFile, len(files))
	copy(sortedFilesDesc, files)
	sort.Slice(sortedFilesDesc, func(i, j int) bool {
		return sortedFilesDesc[i].date > sortedFilesDesc[j].date
	})
	// sort by date asc
	sortedFilesAsc := make([]testFile, len(files))
	copy(sortedFilesAsc, files)
	sort.Slice(sortedFilesAsc, func(i, j int) bool {
		return sortedFilesAsc[i].date < sortedFilesAsc[j].date
	})

	testCases := []struct {
		name                string
		requestQuery        url.Values
		exptectedStatusCode int
		expectedResponse    []testFile
	}{
		{
			name: "search: by name: desc: single file",
			requestQuery: url.Values{
				"query":   {sortedFilesDesc[0].name},
				"page":    {"1"},
				"perPage": {"1"},
				"order":   {"date"},
				"sort":    {"desc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesDesc[0:1],
		},
		{
			name: "search: by name: asc: single file",
			requestQuery: url.Values{
				"query":   {sortedFilesAsc[0].name},
				"page":    {"1"},
				"perPage": {"1"},
				"order":   {"date"},
				"sort":    {"desc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesAsc[0:1],
		},
		{
			name: "search: by name prefix: desc: 3 files",
			requestQuery: url.Values{
				"query":   {testFileNamePrefix},
				"page":    {"1"},
				"perPage": {"3"},
				"order":   {"date"},
				"sort":    {"desc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesDesc[0:3],
		},
		{
			name: "search: by name prefix: asc: 3 files",
			requestQuery: url.Values{
				"query":   {testFileNamePrefix},
				"page":    {"1"},
				"perPage": {"3"},
				"order":   {"date"},
				"sort":    {"asc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesAsc[0:3],
		},
		{
			name: "search: by name prefix: desc: all files",
			requestQuery: url.Values{
				"query":   {testFileNamePrefix},
				"page":    {"1"},
				"perPage": {"10"},
				"order":   {"date"},
				"sort":    {"desc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesDesc,
		},
		{
			name: "search: by name prefix: asc: all files",
			requestQuery: url.Values{
				"query":   {testFileNamePrefix},
				"page":    {"1"},
				"perPage": {"10"},
				"order":   {"date"},
				"sort":    {"asc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesAsc,
		},
		{
			name: "search: by type: desc: single file",
			requestQuery: url.Values{
				"query":   {testFileType},
				"page":    {"1"},
				"perPage": {"1"},
				"order":   {"date"},
				"sort":    {"desc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesDesc[0:1],
		},
		{
			name: "search: by type: desc: single file",
			requestQuery: url.Values{
				"query":   {testFileType},
				"page":    {"1"},
				"perPage": {"1"},
				"order":   {"date"},
				"sort":    {"asc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesAsc[0:1],
		},
		{
			name: "search: by type: desc: 3 files",
			requestQuery: url.Values{
				"query":   {testFileType},
				"page":    {"1"},
				"perPage": {"3"},
				"order":   {"date"},
				"sort":    {"desc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesDesc[0:3],
		},
		{
			name: "search: by type: desc: 3 files",
			requestQuery: url.Values{
				"query":   {testFileType},
				"page":    {"1"},
				"perPage": {"3"},
				"order":   {"date"},
				"sort":    {"asc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesAsc[0:3],
		},
		{
			name: "search: by type: desc: all files",
			requestQuery: url.Values{
				"query":   {testFileType},
				"page":    {"1"},
				"perPage": {"10"},
				"order":   {"date"},
				"sort":    {"desc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesDesc,
		},
		{
			name: "search: by type: asc: all files",
			requestQuery: url.Values{
				"query":   {testFileType},
				"page":    {"1"},
				"perPage": {"10"},
				"order":   {"date"},
				"sort":    {"asc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesAsc,
		},
		// pagination asc
		{
			name: "search: by name prefix: asc: page=1: perPage=2",
			requestQuery: url.Values{
				"query":   {testFileType},
				"page":    {"1"},
				"perPage": {"2"},
				"order":   {"date"},
				"sort":    {"asc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesAsc[0:2],
		},
		{
			name: "search: by name prefix: asc: page=2: perPage=2",
			requestQuery: url.Values{
				"query":   {testFileType},
				"page":    {"2"},
				"perPage": {"2"},
				"order":   {"date"},
				"sort":    {"asc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesAsc[2:4],
		},
		{
			name: "search: by name prefix: asc: page=2: perPage=2",
			requestQuery: url.Values{
				"query":   {testFileType},
				"page":    {"3"},
				"perPage": {"2"},
				"order":   {"date"},
				"sort":    {"asc"},
			},
			exptectedStatusCode: 200,
			expectedResponse:    sortedFilesAsc[4:6],
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
				url:         "/v1/search",
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
			if !assert.True(t, len(resultFiles) >= len(test.expectedResponse), "search result length") {
				return
			}

			index := 0
			for _, expectedFile := range test.expectedResponse {
				var resultFile resultFile
				for ; index < len(resultFiles); index++ {
					if expectedFile.name == resultFiles[index].Name {
						resultFile = resultFiles[index]
						break
					}
				}
				assert.Equal(t, expectedFile.name, resultFile.Name)
				assert.Equal(t, expectedFile.mediaType, resultFile.Type)
				assert.Equal(t, expectedFile.size, resultFile.Size)
				assert.True(t, resultFile.Thumbnail)
			}
		})
	}
}

type resultFile struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Size      int64  `json:"size"`
	Thumbnail bool   `json:"thumbnail"`
}

type resultFiles []resultFile

func unmarshalSearchFilesResponse(r io.Reader) (files resultFiles, err error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &files)
	return
}
