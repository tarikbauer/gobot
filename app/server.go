package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/tarikbauer/gobot/application"
)

type interceptorResponse struct {
	status      int
	wroteHeader bool
	http.ResponseWriter
}

func NewInterceptorResponse(w http.ResponseWriter) *interceptorResponse {
	return &interceptorResponse{ResponseWriter: w}
}

func (w *interceptorResponse) Status() int {
	return w.status
}

func (w *interceptorResponse) Write(p []byte) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}

func (w *interceptorResponse) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.status = code
	w.wroteHeader = true
}

func interceptor(logger logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "logger", logger)
			interceptorWriter := NewInterceptorResponse(w)
			defer func(begin time.Time, w *interceptorResponse) {
				if v := recover(); v != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte("internal server error"))
					logger.Error(v)
				}
				logger.WithFields(logrus.Fields{
					"path": r.URL.Path,
					"method": r.Method,
					"status": http.StatusText(w.Status()),
					"took": time.Since(begin).String(),
				}).Info("")
			}(time.Now(), interceptorWriter)
			next.ServeHTTP(interceptorWriter, r.WithContext(ctx))
		})
	}
}

func RunServer(port int, logger logrus.Logger, service application.Bot, c chan error) {
	h := serverHandler{service: service}
	router := mux.NewRouter().StrictSlash(true)
	router.Use(interceptor(logger))
	router.Path("/").HandlerFunc(h.serve).Methods("GET")
	router.NotFoundHandler = interceptor(logger)(http.HandlerFunc(h.notFound))
	router.MethodNotAllowedHandler = interceptor(logger)(http.HandlerFunc(h.notAllowed))
	c <- http.ListenAndServe(fmt.Sprint(":", port), router)
}
