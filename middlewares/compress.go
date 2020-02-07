package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type CompressWriter struct {
	w     http.ResponseWriter
	out   io.Writer
	total int
}

func (c CompressWriter) Header() http.Header {
	return c.w.Header()
}

func (c CompressWriter) WriteHeader(status int) {
	c.w.WriteHeader(status)
}

func (c CompressWriter) Write(d []byte) (int, error) {
	total, err := c.out.Write(d)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func requestAcceptsGzip(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !requestAcceptsGzip(r) {
			next.ServeHTTP(w, r)
			return
		}
		out := gzip.NewWriter(w)
		defer out.Close()
		compressWriter := CompressWriter{w: w, out: out, total: 0}
		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(compressWriter, r)
	})
}
