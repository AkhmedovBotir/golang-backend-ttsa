package response

import "github.com/gin-gonic/gin"

type Body struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func JSON(c *gin.Context, status int, message string, data interface{}, err interface{}) {
	c.JSON(status, Body{
		Message: message,
		Data:    data,
		Error:   err,
	})
}
