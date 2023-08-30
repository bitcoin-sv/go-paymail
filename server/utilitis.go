package server

import (
	"encoding/json"
	"net/http"
)

func writeJsonResponse(w http.ResponseWriter, statusCode int, response any) {
	if response == nil {
		panic("writeJsonRespone: empty response data")
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		ErrorResponse(w, ErrorFailedMarshalJSON, "failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	writeRespone(w, http.StatusOK, "application/json", jsonData)
}

func writeRespone(w http.ResponseWriter, statusCode int, contentType string, responseData []byte) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	if responseData != nil {
		_, err := w.Write(responseData)
		if err != nil {
			panic("writeRespone: " + err.Error())
		}
	}
}
