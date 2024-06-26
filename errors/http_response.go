package errors

import (
	"errors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
)

// ErrorResponse is a standard way to return errors to the client
func ErrorResponse(c *gin.Context, err error) {
	response, statusCode := getError(err)
	c.JSON(statusCode, response)
}

func getError(err error) (models.ResponseError, int) {
	var errDetails models.SPVError
	ok := errors.As(err, &errDetails)
	if !ok {
		return models.ResponseError{Code: models.UnknownErrorCode, Message: "Unable to get information about error"}, 500
	}

	return models.ResponseError{Code: errDetails.Code, Message: errDetails.Message}, errDetails.StatusCode
}
