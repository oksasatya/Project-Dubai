package config

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

// ForwardProxy function is a utility function that forwards the request to the appropriate service
func ForwardProxy(c echo.Context, serviceURL string) ([]byte, error) {
	logrus.Info("Forwarding request to service: ", serviceURL)
	req, err := http.NewRequest(c.Request().Method, serviceURL, c.Request().Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header = c.Request().Header

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		logrus.Errorf("Service returned an error: %s", string(body))
		return nil, fmt.Errorf("service returned an error: %s", string(body))
	}
	logrus.Infof("Response from service: %s", string(body))
	return body, nil
}
