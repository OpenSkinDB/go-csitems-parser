package modules

import (
	"bytes"
)

func RemoveBOMFromFile(fileBytes []byte) []byte {
	trimmedBytes := bytes.Trim(fileBytes, "\xef\xbb\xbf")
	return trimmedBytes
}
