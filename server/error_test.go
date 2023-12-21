package server

import (
	"github.com/bitcoin-sv/go-paymail/logging"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestErrorResponse will test the method ErrorResponse()
func TestErrorResponse(t *testing.T) {
	t.Run("placeholder test", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		log := logging.GetDefaultLogger()
		ErrorResponse(w, req, ErrorMethodNotFound, "test message", http.StatusBadRequest, log)

		// todo: actually test the error response
	})
}
