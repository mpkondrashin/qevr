/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

progress_reader.go

Reader that callback function to report progress
*/
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
