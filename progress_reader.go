package main

import "io"

type ProgressReader struct {
	reader   io.Reader
	callback func(int)
}

var _ io.Reader = ProgressReader{}

func NewProgressReader(reader io.Reader, callback func(int)) ProgressReader {
	return ProgressReader{
		reader:   reader,
		callback: callback,
	}
}
func (r ProgressReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	r.callback(n)
	return
}
