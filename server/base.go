package server

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// index basic request to /
//nolint: revive // do not check for unused param required by interface
func index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	responseData := map[string]interface{}{"message": "Welcome to the Paymail Server ✌(◕‿-)✌"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(responseData)

	if err != nil {
		ErrorResponse(w, ErrorEncodingResponse, err.Error(), http.StatusInternalServerError)
		return
	}
}

// health is a basic request to return a health response
func health(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

// notFound handles all 404 requests
//nolint: revive // do not check for unused param required by interface
func notFound(w http.ResponseWriter, req *http.Request) {
	ErrorResponse(w, ErrorRequestNotFound, "request not found", http.StatusNotFound)
}

// methodNotAllowed handles all 405 requests
func methodNotAllowed(w http.ResponseWriter, req *http.Request) {
	ErrorResponse(w, ErrorMethodNotFound, "method "+req.Method+" not allowed", http.StatusMethodNotAllowed)
}
