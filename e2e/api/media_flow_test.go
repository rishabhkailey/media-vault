package api_test

// //go:build e2e
// // +build e2e
import (
	"testing"
	"time"
)

func TestCorrectFlow(t *testing.T) {
	testCases := []struct {
		name            string
		file            testFile
		uploadChunkSize int64
		getRangeSize    int64
		bearerToken     string
	}{
		{
			name: "normal upload",
			file: testFile{
				name:          "test.txt",
				date:          time.Now().UnixMilli(),
				data:          randomString(1000),
				size:          1000,
				thumbnialData: randomString(100),
				mediaType:     "txt",
			},
			uploadChunkSize: 100,
			getRangeSize:    100,
			bearerToken:     AUTH_TOKEN,
		},
		{
			name: "uploadChunkSize > file.size",
			file: testFile{
				name:          "test.txt",
				date:          time.Now().UnixMilli(),
				data:          randomString(1000),
				size:          1000,
				thumbnialData: randomString(100),
				mediaType:     "txt",
			},
			uploadChunkSize: 10000,
			getRangeSize:    100,
			bearerToken:     AUTH_TOKEN,
		},
		{
			name: "getRangeSize > file.size",
			file: testFile{
				name:          "test.txt",
				date:          time.Now().UnixMilli(),
				data:          randomString(1000),
				size:          1000,
				thumbnialData: randomString(100),
				mediaType:     "txt",
			},
			uploadChunkSize: 100,
			getRangeSize:    10000,
			bearerToken:     AUTH_TOKEN,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			// different client and session for different tests
			testClient, err := newTestHttpClient()
			if err != nil {
				t.Errorf("newTestHttpClient failed: %v", err)
				return
			}
			resp, err := testClient.UploadTest(t, testCase.file, int64(testCase.uploadChunkSize), testCase.bearerToken)
			if err != nil {
				t.Errorf("upload failed: %v", err)
				return
			}
			var mediaUrl string
			{
				body := resp.body.(map[string]any)
				value := body["url"]
				mediaUrl = value.(string)
			}
			var thumbnailUrl string
			{
				body := resp.body.(map[string]any)
				value := body["thumbnail_url"]
				thumbnailUrl = value.(string)
			}
			var mediaID uint
			{
				body := resp.body.(map[string]any)
				value := body["id"]
				mediaID = uint(value.(float64))
			}
			resp, err = testClient.DownloadTest(t, mediaUrl, testCase.bearerToken, testCase.file)
			if err != nil {
				t.Errorf("DownloadTest failed: %v", err)
				return
			}
			resp, err = testClient.DownloadThumbnailTest(t, thumbnailUrl, testCase.bearerToken, testCase.file)
			if err != nil {
				t.Errorf("DownloadThumbnailTest failed: %v", err)
			}
			resp, err = testClient.GetMediaRangeTest(t, mediaUrl, testCase.bearerToken, testCase.getRangeSize, testCase.file)
			if err != nil {
				t.Errorf("GetMediaRangeTest failed: %v", err)
				return
			}
			resp, err = testClient.DeleteMediaTest(t, mediaID, testCase.bearerToken)
			if err != nil {
				t.Errorf("DeleteMediaTest failed: %v", err)
				return
			}
		})
	}
}
