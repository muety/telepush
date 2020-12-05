package middleware

import "net/http"

type StatusRecorderWriter struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorderWriter) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *StatusRecorderWriter) IsSuccess() bool {
	return r.Status >= 200 && r.Status <= 299
}
