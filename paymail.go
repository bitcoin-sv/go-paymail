// Package paymail is Golang library for implementing Paymail as a client or server
// following the bsvalias specifications: http://bsvalias.org/
//
// If you have any suggestions or comments, please feel free to open an issue on
// this GitHub repository!
//
// By TonicPow Inc (https://tonicpow.com)
package paymail

// Version will return the version of the library
func Version() string {
	return version
}

// UserAgent will return the default user agent string
func UserAgent() string {
	return defaultUserAgent
}
