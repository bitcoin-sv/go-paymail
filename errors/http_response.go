package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// ErrorResponse is a standard way to return errors to the client
func ErrorResponse(c *gin.Context, err error, log *zerolog.Logger) {
	response, statusCode := mapAndLog(err, log)
	c.JSON(statusCode, response)
}

func mapAndLog(err error, log *zerolog.Logger) (ResponseError, int) {
	var res ResponseError
	res.Code = UnknownErrorCode
	statusCode := 500

	logLevel := zerolog.WarnLevel
	exposedInternalError := false
	var extendedErr ExtendedError
	if errors.As(err, &extendedErr) {
		res.Code = extendedErr.GetCode()
		res.Message = extendedErr.GetMessage()
		statusCode = extendedErr.GetStatusCode()
		if statusCode >= http.StatusInternalServerError {
			logLevel = zerolog.ErrorLevel
		}
	} else {
		// we should wrap all internal errors into SPVError (with proper code, message and status code)
		// if you find out that some endpoint produces this warning, feel free to fix it
		exposedInternalError = true
	}

	if log != nil {
		logInstance := log.WithLevel(logLevel).Str("module", "errors")
		if exposedInternalError {
			logInstance.Str("warning", "internal error returned as HTTP response")
		}
		logInstance.Err(err).Msgf("Error HTTP response, returning %d", statusCode)
	}

	return res, statusCode
}
