package core_http_response

import "net/http"

var statusCodeUninitialized = -1

type ResponceWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponceWriter {
	return &ResponceWriter{
		ResponseWriter: w,
		statusCode:     statusCodeUninitialized,
	}
}

func (w *ResponceWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponceWriter) GetStatusCodeOrPanic() int {
	if w.statusCode == statusCodeUninitialized {
		panic("status code is uninitialized")
	}
	return w.statusCode
}
