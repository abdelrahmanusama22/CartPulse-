package httpserver

import "github.com/gin-gonic/gin"

type ErrorDetail struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func WriteError(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, ErrorResponse{Error: ErrorDetail{
		Code: code, Message: message, RequestID: c.GetString("request_id"),
	}})
}
