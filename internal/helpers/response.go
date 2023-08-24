package helpers

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   any    `json:"errors,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func ResponseBuilder(c *gin.Context, code int, what string, err, data any) {
	response := Response{
		Data:    data,
		Error:   err,
		Message: what,
	}

	if code <= 299 && code >= 200 {
		response.Success = true
		response.Message += " successful"
	} else {
		response.Message += " failed"
	}

	c.JSON(code, response)
}
