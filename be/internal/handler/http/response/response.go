package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shortlyst-ai/go-helper"
	"net/http"
)

type (
	ErrorResponse struct {
		Code     string      `json:"code,omitempty"`
		Error    string      `json:"error,omitempty"`
		Message  string      `json:"error_message,omitempty"`
		Payload  interface{} `json:"payload,omitempty"`
		HttpCode int         `json:"-"`
	}

	SuccessResponse struct {
		Success bool `json:"success" default:"true"`
	}

	Message struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

const (
	ErrorDuplicateKey = "duplicate_key"
	ErrorNotFound     = "not_found"
	ErrorUnauthorized = "unauthorized"
	ErrorBadRequest   = "bad_request"
	ErrorServerError  = "server_error"
	ErrorForbidden    = "forbidden"
)

var (
	SuccessOK = SuccessResponse{
		Success: true,
	}

	ErrNotFound = ErrorResponse{
		Error:    ErrorNotFound,
		Message:  "Entry not found",
		HttpCode: http.StatusNotFound,
	}
	ErrBadRequest = ErrorResponse{
		Error:    ErrorBadRequest,
		Message:  "Bad request",
		HttpCode: http.StatusBadRequest,
	}
	ErrUnauthorized = ErrorResponse{
		Error:    ErrorUnauthorized,
		Message:  "Unauthorized, please login",
		HttpCode: http.StatusUnauthorized,
	}
	ErrForbidden = ErrorResponse{
		Error:    ErrorForbidden,
		Message:  "You are unauthorized for this request",
		HttpCode: http.StatusForbidden,
	}
	ErrDuplicate = ErrorResponse{
		Error:    ErrorDuplicateKey,
		Message:  "Created value already exists",
		HttpCode: http.StatusConflict,
	}
	ErrValidation = ErrorResponse{
		Error:    ErrorBadRequest,
		Message:  "Invalid parameters or payload",
		HttpCode: http.StatusUnprocessableEntity,
	}
	ErrServerError = ErrorResponse{
		Error:    ErrorServerError,
		Message:  "Something bad happened",
		HttpCode: http.StatusInternalServerError,
	}
)

func SendErrorResponse(c *gin.Context, err ErrorResponse, msg string) {
	ErrorWithPayload(c, err, msg, nil)
}

func JSONSuccessWithPayload(c *gin.Context, payload interface{}) {
	if payload != nil {
		c.JSON(http.StatusOK, payload)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}
func JSONError(c *gin.Context, err error) {
	var errType helper.Error
	switch {
	case errors.As(err, &errType):
		c.JSON(errType.Code, ErrorResponse{
			Error: errType.Message,
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
	}

	return
}

func ErrorWithPayload(c *gin.Context, err ErrorResponse, msg string, payload interface{}) {
	c.Writer.Header().Del("content-type")
	if msg != "" {
		err.Message = msg
	}
	if payload != nil {
		err.Payload = payload
	}
	status := http.StatusBadRequest
	if err.HttpCode != 0 {
		status = err.HttpCode
	}
	c.JSON(status, err)
}
func Success(c *gin.Context) {
	SuccessWithPayload(c, nil)
}

func SuccessWithPayload(c *gin.Context, payload interface{}) {
	if payload != nil {
		c.JSON(http.StatusOK, payload)
		return
	}
	c.JSON(http.StatusOK, SuccessOK)
}
