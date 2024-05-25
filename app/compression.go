package main

import (
	"bytes"
	"compress/gzip"
)

func GzipCompress(text string) string {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write([]byte(text))
	w.Close()
	return b.String()
}

