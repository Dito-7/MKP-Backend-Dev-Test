package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func Success(c *gin.Context, httpStatus int, message string, data interface{}) {
	c.JSON(httpStatus, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(c *gin.Context, httpStatus int, message string, data interface{}, meta interface{}) {
	c.JSON(httpStatus, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(c *gin.Context, httpStatus int, message string, errs interface{}) {
	c.JSON(httpStatus, StandardResponse{
		Success: false,
		Message: message,
		Errors:  errs,
	})
}

func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error. Please try again later."
	}
	Error(c, http.StatusInternalServerError, message, nil)
}

func BadRequest(c *gin.Context, message string, errors interface{}) {
	Error(c, http.StatusBadRequest, message, errors)
}

func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized access"
	}
	Error(c, http.StatusUnauthorized, message, nil)
}

func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	Error(c, http.StatusNotFound, message, nil)
}

func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message, nil)
}
