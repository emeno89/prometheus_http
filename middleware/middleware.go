package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) StatusString() string {
	return strconv.Itoa(rw.status)
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

type Handler struct {
	duration    *prometheus.HistogramVec
	reqCounter  *prometheus.CounterVec
	respCounter *prometheus.CounterVec
	projectName string
	pathFn      func(r *http.Request) string
}

// NewHandler creates middleware handler, gets projectName for metrics label and pathFunc,
// which should return query path
func NewHandler(projectName string, pathFunc func(r *http.Request) string) *Handler {
	return &Handler{
		duration:    newDuration(),
		reqCounter:  newRequestCounter(),
		respCounter: newResponseCounter(),
		projectName: projectName,
		pathFn:      pathFunc,
	}
}

//Middleware gets http.Handler for execution and wraps it, saving metrics
func (s *Handler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := s.pathFn(r)

		s.reqCounter.WithLabelValues(s.projectName, path, r.Method).Inc()

		timer := prometheus.NewTimer(s.duration.WithLabelValues(s.projectName, path, r.Method))

		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)

		timer.ObserveDuration()

		s.respCounter.WithLabelValues(s.projectName, path, r.Method, wrapped.StatusString()).Inc()
	})
}
