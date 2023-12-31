package server

import (
	"encoding/json"
	"errors"
	"github.com/bitcoin-sv/go-paymail/logging"
	"github.com/rs/zerolog"
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
)

// Error codes for server response errors
const (
	ErrorFindingPaymail                = "error-finding-paymail"
	ErrorInvalidDt                     = "invalid-dt"
	ErrorInvalidParameter              = "invalid-parameter"
	ErrorInvalidPubKey                 = "invalid-pubkey"
	ErrorInvalidSenderHandle           = "invalid-sender-handle"
	ErrorInvalidSignature              = "invalid-signature"
	ErrorMethodNotFound                = "method-405"
	ErrorMissingField                  = "missing-field"
	ErrorPaymailNotFound               = "not-found"
	ErrorRecordingTx                   = "error-recording-tx"
	ErrorRequestNotFound               = "request-404"
	ErrorScript                        = "script-error"
	ErrorUnknownDomain                 = "unknown-domain"
	ErrorFailedMarshalJSON             = "failed-marshal-json"
	ErrorEncodingResponse              = "error-encoding-response"
	ErrorNotImplmented                 = "not-implemented"
	ErrorSimplifiedPaymentVerification = "spv-failed"
)

var (
	// ErrDomainMissing is the error for missing domain
	ErrDomainMissing = errors.New("domain is missing")

	// ErrServiceProviderNil is the error for having a nil service provider
	ErrServiceProviderNil = errors.New("service provider is nil")

	// ErrPortMissing is when the port is not found
	ErrPortMissing = errors.New("missing a port")

	// ErrServiceNameMissing is when the service name is not found
	ErrServiceNameMissing = errors.New("missing service name")

	// ErrCapabilitiesMissing is when the capabilities struct is nil or not set
	ErrCapabilitiesMissing = errors.New("missing capabilities struct")

	// ErrBsvAliasMissing is when the bsv alias version is missing
	ErrBsvAliasMissing = errors.New("missing bsv alias version")

	// ErrFailedMarshalJSON is when the JSON marshal fails
	ErrFailedMarshalJSON = errors.New("failed to marshal JSON response")
)

// ErrorResponse is a standard way to return errors to the client
//
// Specs: http://bsvalias.org/99-01-recommendations.html
func ErrorResponse(w http.ResponseWriter, req *http.Request, code, message string, statusCode int, log *zerolog.Logger) {
	if log == nil {
		log = logging.GetDefaultLogger()
	}

	srvErr := &paymail.ServerError{Code: code, Message: message}
	jsonData, err := json.Marshal(srvErr)

	if err != nil {
		log.Debug().
			Str("logger", "http-error").
			Msgf("%d | %s | %s | %s | %s", http.StatusInternalServerError, req.RemoteAddr, req.Method, req.URL, message)
		http.Error(w, ErrorFailedMarshalJSON, http.StatusInternalServerError)
		return
	}

	errorLogger := log.With().
		Str("logger", "http-error").
		Str("code", code).
		Str("msg", message).
		Logger()

	writeResponse(w, req, &errorLogger, statusCode, "application/json", jsonData)
}
