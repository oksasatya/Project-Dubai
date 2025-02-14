package utils

import (
	"encoding/json"
)

// SendErrorResponse sends a success response to the response channel.
func SendErrorResponse(responseChannel chan string, message string, statusCode int) {
	response := map[string]interface{}{
		"meta": map[string]interface{}{
			"message": message,
			"code":    statusCode,
			"status":  "fail",
		},
		"data": nil,
	}
	responseMsg, _ := json.Marshal(response)
	responseChannel <- string(responseMsg)
}

// SendSuccessResponse sends a success response to the response channel.
func SendSuccessResponse(responseChannel chan string, response interface{}) {
	responseMsg, _ := json.Marshal(response)
	responseChannel <- string(responseMsg)
}
