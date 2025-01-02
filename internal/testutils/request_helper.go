package testutils

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

type TestRequest struct {
	Name            string
	Method          string
	Headers         map[string]string
	Path            string
	Body            []byte
	ExpectedStatus  int
	ExpectedBody    string
	ExpectedHeaders map[string]string
}

// only accepts JSON as body
func testReqSingle(t *testing.T, app *fiber.App, tc TestRequest) *http.Response {
	t.Helper()

	req, err := http.NewRequest(tc.Method, tc.Path, bytes.NewReader(tc.Body))
	req.Header.Set("Content-Type", fiber.MIMEApplicationJSONCharsetUTF8)
	assert.NoError(t, err)

	// Set any custom headers
	for key, value := range tc.Headers {
		req.Header.Set(key, value)
	}

	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Check status code
	assert.Equal(t, tc.ExpectedStatus, resp.StatusCode)

	// Check response body if expected
	if tc.ExpectedBody != "" {
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, tc.ExpectedBody, string(body))
	}

	// Check headers if expected
	for key, value := range tc.ExpectedHeaders {
		assert.Equal(t, value, resp.Header.Get(key))
	}

	return resp
}

// only accepts JSON as body
func TestReqs(t *testing.T, app *fiber.App, requests []TestRequest) []*http.Response {
	t.Helper()
	var results = make([]*http.Response, len(requests))

	for i, tc := range requests {
		results[i] = testReqSingle(t, app, tc)
	}

	return results
}
