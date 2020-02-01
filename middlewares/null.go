package middlewares

import "net/http"

type nullHandler struct{}

func (h *nullHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("Reached null handler")
}

func NullHandler(dir string) http.Handler {
	return http.FileServer(http.Dir(dir))
}
