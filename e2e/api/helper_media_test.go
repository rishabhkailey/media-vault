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
)

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
	body        map[string]any
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
	url, err := url.JoinPath(BASE_URL, req.url)
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
	reqBodyBytes, _ := io.ReadAll(req.bodyReader)
	reqBody := string(reqBodyBytes)
	req.bodyReader = bytes.NewReader(reqBodyBytes)
	r, err := http.NewRequest(req.method, url, req.bodyReader)
	if err != nil {
		return
	}
	_ = reqBody
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
	return c.sendHttpRequest(httpRequest{
		method: "POST",
		url:    "/v1/finishChunkUpload",
		body: map[string]any{
			"requestID": requestID,
			"checksum":  "",
		},
		bearerToken: bearerToken,
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
