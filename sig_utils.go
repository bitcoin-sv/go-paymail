package paymail

import "encoding/base64"

func EncodeSignature(sigBytes []byte) string {
	return base64.StdEncoding.EncodeToString(sigBytes)
}

func DecodeSignature(signature string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(signature)
}
