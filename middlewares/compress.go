package middlewares

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
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

func requestAcceptsBrotli(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "br")
}

func requestAcceptsGzip(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

func requestAcceptsDeflate(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "deflate")
}

func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var out io.WriteCloser = nil
		var cType string = ""
		if requestAcceptsBrotli(r) {
			out = brotli.NewWriter(w)
			cType = "br"
		} else if requestAcceptsGzip(r) {
			out = gzip.NewWriter(w)
			cType = "gzip"
		} else if requestAcceptsDeflate(r) {
			var err error
			out, err = flate.NewWriter(w, 0)
			if err != nil {
				panic(err)
			}
			cType = "deflate"
		}
		if out == nil {
			next.ServeHTTP(w, r)
			return
		}
		defer out.Close()
		compressWriter := CompressWriter{w: w, out: out, total: 0}
		w.Header().Set("Content-Encoding", cType)
		next.ServeHTTP(compressWriter, r)
	})
}
