package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func GzipCompress(text string) string {
	fmt.Println("Text: " + text)
	buffer := new(bytes.Buffer)
	writer := gzip.NewWriter(buffer)
	writer.Write([]byte(text))
	writer.Close()
	resultString := buffer.String()
	fmt.Println("Compressed: " + resultString + "\r\n")
	return resultString
}
