package errors

import (
	"errors"
	"github.com/gin-gonic/gin"
)

// ErrorResponse is a standard way to return errors to the client
func ErrorResponse(c *gin.Context, err error) {
	response, statusCode := getError(err)
	c.JSON(statusCode, response)
}

func getError(err error) (ResponseError, int) {
	var extendedErr ExtendedError
	if errors.As(err, &extendedErr) {
		return ResponseError{
			Code:    extendedErr.GetCode(),
			Message: extendedErr.GetMessage(),
		}, extendedErr.GetStatusCode()
	}

	return ResponseError{
		Code:    UnknownErrorCode,
		Message: "Unable to get information about error",
	}, 500
}
