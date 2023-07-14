package helpers

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseSuccess(c *gin.Context, message string, data interface{}, httpCode int) {

	if message == "" {
		message = "Success"
	}

	response := Response{
		Status:  "Success",
		Message: message,
		Data:    data,
	}

	c.JSON(httpCode, response)

}

func ResponseError(c *gin.Context, httpCode int, err error) {

	errMsg := err.Error()

	response := Response{
		Status:  "Error",
		Message: errMsg,
	}

	c.JSON(httpCode, response)
}
