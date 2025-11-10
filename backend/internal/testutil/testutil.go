package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// NewJSONContext builds a Gin context for a JSON request.
func NewJSONContext(t *testing.T, method, path string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	data, err := json.Marshal(body)
	require.NoError(t, err)

	return NewRawContext(t, method, path, data, "application/json")
}

// NewRequestContext builds a Gin context for a request without a body.
func NewRequestContext(t *testing.T, method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	return NewRawContext(t, method, path, nil, "")
}

// NewRawContext builds a Gin context for a request with arbitrary payload.
func NewRawContext(t *testing.T, method, path string, payload []byte, contentType string) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	reader := bytes.NewReader(payload)

	req := httptest.NewRequest(method, path, reader)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return buildContext(req)
}

// DecodeResponse unmarshals the recorder body into a map.
func DecodeResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)

	return result
}

// CookieByName searches for a cookie with the provided name.
func CookieByName(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

// PerformRequest executes a request against the provided router and returns the recorder.
func PerformRequest(t *testing.T, router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()

	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func buildContext(req *http.Request) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}
