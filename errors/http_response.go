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
	if err == nil {
		return ResponseError{Code: UnknownErrorCode, Message: "No error information available"}, 500
	}

	var errDetails SPVError
	ok := errors.As(err, &errDetails)
	if !ok {
		return ResponseError{Code: UnknownErrorCode, Message: "Unable to get information about error"}, 500
	}

	return ResponseError{Code: errDetails.Code, Message: errDetails.Message}, errDetails.StatusCode
}
