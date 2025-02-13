package webResponse

import (
	"github.com/labstack/echo/v4"
)

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

// ResponseJson is a utility function that sends a JSON response to the client
func ResponseJson(c echo.Context, status int, payload interface{}, message string) error {
	response := Response{
		Meta: Meta{
			Message: message,
			Code:    status,
			Status:  getStatusText(status),
		},
		Data: payload,
	}
	return c.JSON(status, response)
}

func getStatusText(status int) string {
	if status >= 200 && status < 300 {
		return "success"
	} else if status >= 400 && status < 500 {
		return "fail"
	} else {
		return "error"
	}
}
