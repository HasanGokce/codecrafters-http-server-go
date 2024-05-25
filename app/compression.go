package main

import (
	"bytes"
	"compress/gzip"
)

func GzipCompress(text string) string {
	buffer := new(bytes.Buffer)
	writer := gzip.NewWriter(buffer)
	writer.Write([]byte(text))
	writer.Close()
	return buffer.String()
}
