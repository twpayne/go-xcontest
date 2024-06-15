package xcontest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/twpayne/go-xctrack"
)

// A SavedTask is a saved task.
type SavedTask struct {
	Task         *xctrack.Task
	Author       string
	LastModified time.Time
}

// A SavedTaskV2 is a saved task in QR code format.
type SavedTaskV2 struct {
	Task         *xctrack.QRCodeTask
	Author       string
	LastModified time.Time
}

// XCTskLoad loads the saved task with the given task code.
func (c *Client) XCTaskLoad(ctx context.Context, taskCode string) (*SavedTask, error) {
	var task xctrack.Task
	author, lastModified, err := c.doRequestTask(ctx, "/xctsk/load/"+taskCode, &task)
	if err != nil {
		return nil, err
	}
	return &SavedTask{
		Task:         &task,
		Author:       author,
		LastModified: lastModified,
	}, nil
}

// XCTskLoadV2 loads the saved task with the given task code in QR code format.
func (c *Client) XCTaskLoadV2(ctx context.Context, taskCode string) (*SavedTaskV2, error) {
	var task xctrack.QRCodeTask
	author, lastModified, err := c.doRequestTask(ctx, "/xctsk/loadV2/"+taskCode, &task)
	if err != nil {
		return nil, err
	}
	return &SavedTaskV2{
		Task:         &task,
		Author:       author,
		LastModified: lastModified,
	}, nil
}

// XCTskQR returns SVG data for the given task.
func (c *Client) XCTskQR(ctx context.Context, task *xctrack.QRCodeTask) ([]byte, error) {
	request, err := c.newPostQRCodeTaskRequest(ctx, "/xctsk/qr", task)
	if err != nil {
		return nil, err
	}
	_, body, err := c.doRequest(request)
	return body, err
}

// XCTskSave saves the given task and returns the task code.
func (c *Client) XCTskSave(ctx context.Context, task *xctrack.QRCodeTask, author string) (string, error) {
	request, err := c.newPostQRCodeTaskRequest(ctx, "/xctsk/save", task)
	if err != nil {
		return "", err
	}
	if author != "" {
		request.Header.Set("Author", author)
	}
	var saveResponse struct {
		TaskCode string `json:"taskCode"`
	}
	if _, err := c.doRequestJSON(request, &saveResponse); err != nil {
		return "", err
	}
	return saveResponse.TaskCode, nil
}

// TaskURL returns the URL for the give task code.
func (c *Client) TaskURL(taskCode string) string {
	return c.baseURL + "/xctsk/load?taskCode=" + taskCode
}

func (c *Client) newPostQRCodeTaskRequest(ctx context.Context, relURL string, task *xctrack.QRCodeTask) (*http.Request, error) {
	bodyData, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+relURL, bytes.NewReader(bodyData))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-type", "application/json")
	return request, nil
}

func (c *Client) doRequest(request *http.Request) (*http.Response, []byte, error) {
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return response, nil, errExpectedHTTPStatusOK(response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return response, nil, err
	}
	return response, body, nil
}

func (c *Client) doRequestJSON(request *http.Request, value any) (*http.Response, error) {
	response, body, err := c.doRequest(request)
	if err != nil {
		return response, err
	}
	if err := json.Unmarshal(body, value); err != nil {
		return response, err
	}
	return response, nil
}

func (c *Client) doRequestTask(ctx context.Context, relURL string, task any) (string, time.Time, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+relURL, nil)
	if err != nil {
		return "", time.Time{}, err
	}
	response, err := c.doRequestJSON(request, task)
	if err != nil {
		return "", time.Time{}, err
	}
	author := response.Header.Get("Author")
	lastModified, err := http.ParseTime(response.Header.Get("Last-modified"))
	if err != nil {
		return "", time.Time{}, err
	}
	return author, lastModified, nil
}
