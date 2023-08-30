package server

import (
	"net/http"
)

// CreateMetadata will create the base metadata using the request
func CreateMetadata(req *http.Request, alias, domain, optionalNote string) *RequestMetadata {
	ipAddress := req.Header.Get("X-Real-IP")
	if ipAddress == "" {
		ipAddress = req.Header.Get("X-Forwarded-For")
		if ipAddress == "" {
			ipAddress = req.RemoteAddr
		}
	}

	return &RequestMetadata{
		Alias:      alias,
		Domain:     domain,
		IPAddress:  ipAddress,
		Note:       optionalNote,
		RequestURI: req.RequestURI,
		UserAgent:  req.UserAgent(),
	}
}
