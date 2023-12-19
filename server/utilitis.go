package server

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"net/http"
)

func writeJsonResponse(w http.ResponseWriter, req *http.Request, log *zerolog.Logger, response any) {
	if response == nil {
		panic("writeJsonRespone: empty response data")
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		ErrorResponse(w, req, ErrorFailedMarshalJSON, "failed to marshal JSON response", http.StatusInternalServerError, log)
		return
	}

	writeResponse(w, req, log, http.StatusOK, "application/json", jsonData)
}

func writeResponse(w http.ResponseWriter, req *http.Request, log *zerolog.Logger, statusCode int, contentType string, responseData []byte) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	log.Debug().Msgf("%d | %s | %s | %s ", statusCode, req.RemoteAddr, req.Method, req.URL)

	if responseData != nil {
		_, err := w.Write(responseData)
		if err != nil {
			panic("writeResponse: " + err.Error())
		}
	}
}
