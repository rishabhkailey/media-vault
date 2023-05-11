package api_test

// //go:build e2e
// // +build e2e
import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// copy it from the UI
const AUTH_TOKEN = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJzcGEtdGVzdCIsImV4cCI6MTY4MzE3NTY2OCwic3ViIjoiMSJ9.nW3K0c9zIn1Skl1P4BX3SSvaiyEgU863JRUyhj_q8UTcNPCp_9lTX_GQWi3ndvRGrANK6K6CxPNpBPHJZT8CHihHspMUZVgsX03N1axYqMgLj6jdOkaK5bvSoTWnQdz_f-Qy5KDInBWEpqo5SEkQ6I6iRj6GyvArv4Rn-aWnNLWQpBHdO4JCuTrlwJO_Nohtv8fNdYPQVGIo2bWHuReqJE1nTJrVFgbdMgviUM9k-R0c3JRKhWzy5-thQxFdhC8IwtfrL_c8JNG9faNfAQnf1h8wu_BMbHLvdmhqpa399vOi_SymFdiXGwUIJd0Vtok12TP4p-s8UUgPKtASAAn97A"
const BASE_URL = "http://localhost:8090"

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
				date:          time.Now().Unix(),
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
				date:          time.Now().Unix(),
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
				date:          time.Now().Unix(),
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
		})
	}
}

type testFile struct {
	name          string
	size          int64
	date          int64
	data          string
	mediaType     string
	thumbnialData string
}

func (c *testHttpClient) UploadTest(t *testing.T, file testFile, chunkSize int64, bearerToken string) (httpResponse, error) {
	if chunkSize == 0 {
		return httpResponse{}, fmt.Errorf("chunkSize should be greater than 0")
	}
	// initchunkupload
	resp, err := c.sendInitChunkUploadRequest(
		file.name,
		file.size,
		file.mediaType,
		file.date,
		bearerToken,
	)
	if err != nil {
		return resp, fmt.Errorf("init chunkupload failed: %v", err)
	}
	if ok := assert.Equal(t, 200, resp.status, "status"); !ok {
		return resp, fmt.Errorf("init chunkupload failed: expected 200 status code but got %v", resp.status)
	}

	var requestID string
	{
		body := resp.body.(map[string]any)
		value := body["requestID"]
		requestID = value.(string)
	}
	if ok := assert.True(t, len(requestID) > 0, "requestID doesn't exist in reponse"); !ok {
		return resp, fmt.Errorf("init chunkupload invalid response: requestID param is missing in response")
	}

	// uploadChunk
	{
		index := int64(0)
		for index < file.size {
			endIndex := index + chunkSize
			if endIndex > file.size {
				endIndex = file.size
				chunkSize = endIndex - index
			}
			chunkData := file.data[index:endIndex]
			resp, err := c.sendUploadChunkRequest(requestID, index, chunkSize, chunkData, file.name, bearerToken)
			if !assert.NoError(t, err, "uploadChunk request") {
				return resp, fmt.Errorf("uploadChunk request failed: %w", err)
			}
			if !assert.Equal(t, 200, resp.status, "status") {
				return resp, fmt.Errorf("uploadChunk request failed: expected 200 status code but got %v", resp.status)
			}
			index = endIndex
		}
	}

	// upload thumbnail
	{
		resp, err := c.sendUploadThumbnailRequest(requestID, file.thumbnialData, file.name, bearerToken)
		if !assert.NoError(t, err, "uploadThumbnail request") {
			return resp, fmt.Errorf("uploadChunk request failed: %w", err)
		}
		if !assert.Equal(t, 200, resp.status, "uploadThumbnail status") {
			return resp, fmt.Errorf("uploadThumbnail request failed: expected 200 status code but got %v", resp.status)
		}
	}
	// finish upload
	{
		resp, err = c.sendFinishUploadRequest(requestID, AUTH_TOKEN)
		if !assert.NoError(t, err, "finishChunkUpload request") {
			return resp, fmt.Errorf("finishChunkUpload request failed: %w", err)
		}
		if !assert.Equal(t, 200, resp.status, "finishChunkUpload status") {
			return resp, fmt.Errorf("finishChunkUpload request failed: expected 200 status code but got %v", resp.status)
		}
	}
	return resp, nil
}

// do not use this for big test files it will use memory = file size
func (c *testHttpClient) DownloadTest(t *testing.T, url string, bearerToken string, expectedFile testFile) (httpResponse, error) {
	resp, err := c.sendRefreshSessionRequest(bearerToken)
	if !assert.NoError(t, err, "sendRefreshSessionRequest request") {
		return resp, fmt.Errorf("sendRefreshSessionRequest request failed: %w", err)
	}
	if !assert.Equal(t, 200, resp.status, "sendRefreshSessionRequest status") {
		return resp, fmt.Errorf("sendRefreshSessionRequest request failed: expected 200 status code but got %v", resp.status)
	}
	resp, err = c.sendGetMediaRequest(url)
	if !assert.NoError(t, err, "GetMediaRequest request") {
		return resp, fmt.Errorf("GetMediaRequest request failed: %w", err)
	}
	if !assert.Equal(t, 200, resp.status, "GetMediaRequest status") {
		return resp, fmt.Errorf("GetMediaRequest request failed: expected 200 status code but got %v", resp.status)
	}
	respData, ok := resp.body.(string)
	if !assert.True(t, ok, "GetMediaRequest invalid response type") {
		return resp, fmt.Errorf("GetMediaRequst request failed: string type cast failed")
	}
	if !assert.Equal(t, expectedFile.size, int64(len(respData)), "GetMediaRequest file size") {
		return resp, fmt.Errorf("GetMediaRequest file size: expected %d got %d", expectedFile.size, len(respData))
	}
	if !assert.Equal(t, expectedFile.data, respData, "GetMediaRequest file data") {
		return resp, fmt.Errorf("GetMediaRequest file data did not match")
	}
	return resp, err
}

func (c *testHttpClient) DownloadThumbnailTest(t *testing.T, url string, bearerToken string, expectedFile testFile) (httpResponse, error) {

	resp, err := c.sendRefreshSessionRequest(bearerToken)
	if !assert.NoError(t, err, "sendRefreshSessionRequest request") {
		return resp, fmt.Errorf("sendRefreshSessionRequest request failed: %w", err)
	}
	if !assert.Equal(t, 200, resp.status, "sendRefreshSessionRequest status") {
		return resp, fmt.Errorf("sendRefreshSessionRequest request failed: expected 200 status code but got %v", resp.status)
	}
	resp, err = c.sendGetMediaRequest(url)
	if !assert.NoError(t, err, "GetMediaRequest request") {
		return resp, fmt.Errorf("GetMediaRequest request failed: %w", err)
	}
	if !assert.Equal(t, 200, resp.status, "GetMediaRequest status") {
		return resp, fmt.Errorf("GetMediaRequest request failed: expected 200 status code but got %v", resp.status)
	}
	respData, ok := resp.body.(string)
	if !assert.True(t, ok, "GetMediaRequest invalid response type") {
		return resp, fmt.Errorf("GetMediaRequst request failed: string type cast failed")
	}
	if !assert.Equal(t, len(expectedFile.thumbnialData), len(respData), "GetMediaRequest file size") {
		return resp, fmt.Errorf("GetMediaRequest file size: expected %d got %d", expectedFile.size, len(respData))
	}
	if !assert.Equal(t, expectedFile.thumbnialData, respData, "GetMediaRequest file data") {
		return resp, fmt.Errorf("GetMediaRequest file data did not match")
	}
	return resp, err
}

func (c *testHttpClient) GetMediaRangeTest(t *testing.T, url string, bearerToken string, rangeSize int64, expectedFile testFile) (httpResponse, error) {
	if rangeSize == 0 {
		return httpResponse{}, fmt.Errorf("rangeSize should be greater than 0")
	}
	resp, err := c.sendRefreshSessionRequest(bearerToken)
	if !assert.NoError(t, err, "sendRefreshSessionRequest request") {
		return resp, fmt.Errorf("sendRefreshSessionRequest request failed: %w", err)
	}
	if !assert.Equal(t, 200, resp.status, "sendRefreshSessionRequest status") {
		return resp, fmt.Errorf("sendRefreshSessionRequest request failed: expected 200 status code but got %v", resp.status)
	}

	var index int64 = 0
	for index < expectedFile.size {
		startRange := index
		endRange := index + rangeSize
		if endRange > expectedFile.size {
			endRange = expectedFile.size
		}
		resp, err = c.sendGetMediaRangeRequest(url, startRange, endRange)
		if !assert.NoError(t, err, "GetMediaRangeRequest request") {
			return resp, fmt.Errorf("GetMediaRangeRequest request failed: %w", err)
		}
		if !assert.Equal(t, 206, resp.status, "GetMediaRangeRequest status") {
			return resp, fmt.Errorf("GetMediaRangeRequest request failed: expected 200 status code but got %v", resp.status)
		}
		respData, ok := resp.body.(string)
		if !assert.True(t, ok, "GetMediaRangeRequest invalid response type") {
			return resp, fmt.Errorf("GetMediaRangeRequest request failed: string type cast failed")
		}
		if !assert.Equal(t, endRange-startRange, int64(len(respData)), "GetMediaRangeRequest file size") {
			return resp, fmt.Errorf("GetMediaRangeRequest file size: expected %d got %d", endRange-startRange, len(respData))
		}
		if !assert.Equal(t, expectedFile.data[startRange:endRange], respData, "GetMediaRangeRequest file data") {
			return resp, fmt.Errorf("GetMediaRangeRequest file data did not match")
		}
		index = endRange
	}
	return resp, err
}
