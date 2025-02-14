package config

import (
	"api-gateway/webResponse"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// ForwardProxy function is a utility function that forwards the request to the appropriate service
func ForwardProxy(c echo.Context, serviceURL string) error {
	req, err := http.NewRequest(c.Request().Method, serviceURL, c.Request().Body)
	if err != nil {
		return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to create request")
	}

	req.Header = c.Request().Header

	// Create a new HTTP client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to send request")
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}

	return c.Stream(resp.StatusCode, resp.Header.Get(contentType), resp.Body)
}
