package utils

import "encoding/base64"

func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func DecodeBase64(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}
