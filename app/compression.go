package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func GzipCompress(text string) string {
	buffer := new(bytes.Buffer)
	writer := gzip.NewWriter(buffer)
	writer.Write([]byte(text))
	writer.Close()
	resultString := buffer.String()
	fmt.Println("Compressed: " + resultString)
	return resultString
}
