package api_test

// //go:build e2e
// // +build e2e
import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"time"

	"github.com/stretchr/testify/assert"
)

// copy it from the UI
const AUTH_TOKEN = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJzcGEtdGVzdCIsImV4cCI6MTY4NDMzOTQ0Miwic3ViIjoiMSJ9.HiHiTfZ7UojmHHZdH2crB3I8SqxXPC4SIVjRhHMur_kaYsgYYG3uQKn7Q9bm6ZyZxqhYp32yZ8eejIkmViqCi9k5IIgcISmJVIS-XpPccI5HOH07R0R77AXt-S3kGrssHsnsPs3IqwJnwF-1evO2_UTfiqrmu2Ca5L-tI6kxFSoUdX44tLcwuHrHcb4If1xojquB9pNwYuSSFPtZCDwaCQIUj7tpQzXdmmVISlOKKdK9BMEQRBuMR36kdQbWxbN9JhydhGZEBT_NLOasjmQwbKh-pHKcfDHu6kH9tZpXLs1lakQk8al-zaX44vvsXrF1fJR88e3RqTLmGZpS1IUmMQ"
const BASE_URL = "http://localhost:8090"

type testHttpClient struct {
	http.Client
}

func newTestHttpClient() (*testHttpClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}
	return &testHttpClient{
		Client: *client,
	}, nil
}

type httpRequest struct {
	method      string
	url         string
	body        any
	query       url.Values
	bodyReader  io.Reader // either bodyReader or body. bodyReader take precedence over body
	bearerToken string
	headers     http.Header
}

type httpResponse struct {
	status  int
	body    any
	headers http.Header
}

func (c *testHttpClient) sendHttpRequest(req httpRequest, jsonResponse bool) (resp httpResponse, err error) {
	requestUrl, err := url.JoinPath(BASE_URL, req.url)
	if req.query != nil {
		parsedUrl, err := url.Parse(requestUrl)
		if err != nil {
			return resp, err
		}
		parsedUrl.RawQuery = req.query.Encode()
		requestUrl = parsedUrl.String()
	}
	if err != nil {
		return
	}
	if req.bodyReader == nil {
		var bodyBytes []byte
		bodyBytes, err = json.Marshal(req.body)
		if err != nil {
			return
		}
		req.bodyReader = bytes.NewReader(bodyBytes)
	}
	r, err := http.NewRequest(req.method, requestUrl, req.bodyReader)
	if err != nil {
		return
	}
	if req.headers != nil {
		r.Header = req.headers
	}
	if len(req.bearerToken) != 0 {
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", req.bearerToken))
	}
	response, err := c.Do(r)
	if err != nil {
		return
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	if len(responseData) > 0 && jsonResponse {
		err = json.Unmarshal(responseData, &resp.body)
		if err != nil {
			return
		}
	} else {
		resp.body = responseData
	}
	resp.status = response.StatusCode
	resp.headers = response.Header
	return
}

func (c *testHttpClient) sendInitChunkUploadRequest(fileName string, fileSize int64, mediaType string, date int64, bearerToken string) (httpResponse, error) {
	return c.sendHttpRequest(httpRequest{
		method: "POST",
		url:    "/v1/initChunkUpload",
		body: map[string]any{
			"fileName":  fileName,
			"Size":      fileSize,
			"MediaType": mediaType,
			"Date":      date,
		},
		bearerToken: bearerToken,
	}, true)
}

func (c *testHttpClient) sendUploadChunkRequest(requestID string, index int64, chunkSize int64, chunkData string, fileName string, bearerToken string) (resp httpResponse, err error) {
	var reqBody bytes.Buffer
	writer := multipart.NewWriter(&reqBody)
	writer.WriteField("requestID", requestID)
	writer.WriteField("index", fmt.Sprintf("%d", index))
	writer.WriteField("chunkSize", fmt.Sprintf("%d", chunkSize))
	part, err := writer.CreateFormFile("chunkData", fileName)
	if err != nil {
		return
	}
	_, err = part.Write([]byte(chunkData))
	if err != nil {
		return
	}
	writer.Close()
	headers := http.Header{}
	headers.Add("Content-Type", writer.FormDataContentType())
	resp, err = c.sendHttpRequest(httpRequest{
		method:      "POST",
		url:         "/v1/uploadChunk",
		bodyReader:  &reqBody,
		bearerToken: bearerToken,
		headers:     headers,
	}, true)
	return
}
func (c *testHttpClient) sendUploadThumbnailRequest(requestID string, thumbnailData string, fileName string, bearerToken string) (resp httpResponse, err error) {
	var reqBody bytes.Buffer
	writer := multipart.NewWriter(&reqBody)
	writer.WriteField("requestID", requestID)
	writer.WriteField("size", fmt.Sprintf("%d", len(thumbnailData)))
	part, err := writer.CreateFormFile("thumbnail", fileName)
	if err != nil {
		return
	}
	_, err = part.Write([]byte(thumbnailData))
	if err != nil {
		return
	}
	writer.Close()
	// writer.WriteField()
	headers := http.Header{}
	headers.Add("Content-Type", writer.FormDataContentType())

	return c.sendHttpRequest(httpRequest{
		method:      "POST",
		url:         "/v1/uploadThumbnail",
		bodyReader:  &reqBody,
		bearerToken: bearerToken,
		headers:     headers,
	}, true)
}

func (c *testHttpClient) sendFinishUploadRequest(requestID string, bearerToken string) (resp httpResponse, err error) {
	headers := http.Header{}
	headers.Add("content-type", "application/json")
	return c.sendHttpRequest(httpRequest{
		method: "POST",
		url:    "/v1/finishChunkUpload",
		body: map[string]any{
			"requestID": requestID,
			"checksum":  "",
		},
		bearerToken: bearerToken,
		headers:     headers,
	}, true)
}

func (c *testHttpClient) sendRefreshSessionRequest(bearerToken string) (resp httpResponse, err error) {
	return c.sendHttpRequest(httpRequest{
		method:      "POST",
		bearerToken: bearerToken,
		url:         "/v1/refreshSession",
	}, true)

}
func (c *testHttpClient) sendGetMediaRequest(url string) (resp httpResponse, err error) {
	resp, err = c.sendHttpRequest(httpRequest{
		method: "GET",
		url:    url,
	}, false)
	if err == nil {
		responseBytes, ok := resp.body.([]byte)
		if !ok {
			return resp, fmt.Errorf("invalid response type from sendHttpRequest. expected []byte")
		}
		resp.body = string(responseBytes)
	}
	return
}

func (c *testHttpClient) sendGetMediaRangeRequest(url string, startRange, endRange int64) (resp httpResponse, err error) {
	headers := http.Header{}
	headers.Add("Range", fmt.Sprintf("bytes=%d-%d", startRange, endRange))
	resp, err = c.sendHttpRequest(httpRequest{
		method:  "GET",
		url:     url,
		headers: headers,
	}, false)
	if err == nil {
		responseBytes, ok := resp.body.([]byte)
		if !ok {
			return resp, fmt.Errorf("invalid response type from sendHttpRequest. expected []byte")
		}
		resp.body = string(responseBytes)
	}
	return
}

// it will panic if any error so only use this in tests
func randomString(n int64) string {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return string(bytes)
}

type testFile struct {
	name          string
	size          int64
	date          int64
	data          string
	mediaType     string
	thumbnialData string
}
type ByDate []testFile

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].date < a[j].date }

var _ sort.Interface = (ByDate)(nil)

func (c *testHttpClient) UploadTest(t assert.TestingT, file testFile, chunkSize int64, bearerToken string) (httpResponse, error) {
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
func (c *testHttpClient) DownloadTest(t assert.TestingT, url string, bearerToken string, expectedFile testFile) (httpResponse, error) {
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

func (c *testHttpClient) DownloadThumbnailTest(t assert.TestingT, url string, bearerToken string, expectedFile testFile) (httpResponse, error) {
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

func (c *testHttpClient) GetMediaRangeTest(t assert.TestingT, url string, bearerToken string, rangeSize int64, expectedFile testFile) (httpResponse, error) {
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
	for index < expectedFile.size-1 {
		startRange := index
		// end range inclusive
		endRange := index + rangeSize
		if endRange > expectedFile.size-1 {
			endRange = expectedFile.size - 1
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
		if !assert.Equal(t, endRange-startRange+1, int64(len(respData)), "GetMediaRangeRequest file size") {
			return resp, fmt.Errorf("GetMediaRangeRequest file size: expected %d got %d", endRange-startRange, len(respData))
		}
		// end range inclusive
		if !assert.Equal(t, expectedFile.data[startRange:endRange+1], respData, "GetMediaRangeRequest file data") {
			return resp, fmt.Errorf("GetMediaRangeRequest file data did not match")
		}
		index = endRange
	}
	return resp, err
}

func (c *testHttpClient) DeleteMediaTest(t assert.TestingT, mediaID uint, bearerToken string) (httpResponse, error) {
	resp, err := c.sendHttpRequest(
		httpRequest{
			method:      "DELETE",
			url:         fmt.Sprintf("/v1/media/%d", mediaID),
			bearerToken: bearerToken,
		}, false,
	)
	if !assert.NoError(t, err, "delete media request") {
		return resp, fmt.Errorf("[testHttpClient.DeleteMediaTest] deleting %d media failed: %w", mediaID, err)
	}
	if !assert.Equal(t, http.StatusOK, resp.status, "delete media request status code") {
		return resp, fmt.Errorf("[testHttpClient.DeleteMediaTest] deleting %d media failed with status %d", mediaID, resp.status)
	}
	return resp, err
}

func (c *testHttpClient) GenerateAndUploadTestFiles(t assert.TestingT, n int, fileNamePrefix string, fileType string, fileSize int64, thumbnailSize int64, timeDifference time.Duration) (files []testFile, mediaIDs []uint, err error) {
	for index := 0; index < n; index++ {
		file := testFile{
			name:          fmt.Sprintf("%s-%02d", fileNamePrefix, index+1),
			size:          fileSize,
			data:          randomString(fileSize),
			date:          time.Now().Add(timeDifference * time.Duration(-1*index)).UnixMilli(),
			mediaType:     fileType,
			thumbnialData: randomString(thumbnailSize),
		}
		var resp httpResponse
		resp, err = c.UploadTest(t, file, 10_000_000, AUTH_TOKEN)
		if err != nil {
			return
		}
		var id uint
		{
			body := resp.body.(map[string]any)
			value := body["id"]
			id = uint(value.(float64))
		}
		mediaIDs = append(mediaIDs, id)
		files = append(files, file)
	}
	return
}

func (c *testHttpClient) DeleteTestMediaFiles(t assert.TestingT, mediaIDs []uint, bearerToken string) (errs []error) {
	for _, mediaId := range mediaIDs {
		_, err := c.DeleteMediaTest(t, mediaId, bearerToken)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return
}
