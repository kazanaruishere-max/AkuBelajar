package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse is the standard success response format.
type SuccessResponse struct {
	Data interface{} `json:"data"`
	Meta *Meta       `json:"meta,omitempty"`
}

// Meta holds pagination metadata.
type Meta struct {
	Page      int   `json:"page"`
	PerPage   int   `json:"per_page"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

// ErrorBody is the standard error response format.
type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail holds error information.
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// OK sends a 200 response with data.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{Data: data})
}

// OKWithMeta sends a 200 response with data and pagination metadata.
func OKWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, SuccessResponse{Data: data, Meta: meta})
}

// Created sends a 201 response with data.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{Data: data})
}

// NoContent sends a 204 response.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest sends a 400 error response.
func BadRequest(c *gin.Context, code, message string) {
	c.JSON(http.StatusBadRequest, ErrorBody{
		Error: ErrorDetail{Code: code, Message: message},
	})
}

// BadRequestWithDetails sends a 400 error response with validation details.
func BadRequestWithDetails(c *gin.Context, code, message string, details interface{}) {
	c.JSON(http.StatusBadRequest, ErrorBody{
		Error: ErrorDetail{Code: code, Message: message, Details: details},
	})
}

// Unauthorized sends a 401 error response.
func Unauthorized(c *gin.Context, code, message string) {
	c.JSON(http.StatusUnauthorized, ErrorBody{
		Error: ErrorDetail{Code: code, Message: message},
	})
}

// Forbidden sends a 403 error response.
func Forbidden(c *gin.Context, code, message string) {
	c.JSON(http.StatusForbidden, ErrorBody{
		Error: ErrorDetail{Code: code, Message: message},
	})
}

// NotFound sends a 404 error response.
func NotFound(c *gin.Context, code, message string) {
	c.JSON(http.StatusNotFound, ErrorBody{
		Error: ErrorDetail{Code: code, Message: message},
	})
}

// Conflict sends a 409 error response.
func Conflict(c *gin.Context, code, message string) {
	c.JSON(http.StatusConflict, ErrorBody{
		Error: ErrorDetail{Code: code, Message: message},
	})
}

// TooManyRequests sends a 429 error response.
func TooManyRequests(c *gin.Context, code, message string) {
	c.JSON(http.StatusTooManyRequests, ErrorBody{
		Error: ErrorDetail{Code: code, Message: message},
	})
}

// InternalError sends a 500 error response.
func InternalError(c *gin.Context, code, message string) {
	c.JSON(http.StatusInternalServerError, ErrorBody{
		Error: ErrorDetail{Code: code, Message: message},
	})
}

// ServiceUnavailable sends a 503 error response.
func ServiceUnavailable(c *gin.Context, code, message string) {
	c.JSON(http.StatusServiceUnavailable, ErrorBody{
		Error: ErrorDetail{Code: code, Message: message},
	})
}
