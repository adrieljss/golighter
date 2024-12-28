package testutils

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

type TestRequest struct {
	Name            string
	Method          string
	Path            string
	Body            string
	ExpectedStatus  int
	ExpectedBody    string
	ExpectedHeaders map[string]string
	SetupHeaders    map[string]string
}

func TestRequests(t *testing.T, app *fiber.App, requests []TestRequest) {
	t.Helper()

	for _, tc := range requests {
		t.Run(tc.Name, func(t *testing.T) {
			req, err := http.NewRequest(tc.Method, tc.Path, strings.NewReader(tc.Body))
			assert.NoError(t, err)

			// Set any custom headers
			for key, value := range tc.SetupHeaders {
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
		})
	}
}
