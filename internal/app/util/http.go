package util

import (
	"net/http"
	"strings"
)

func IsFilteredHttpRoute(r *http.Request) bool {
	path := r.URL.Path
	return strings.HasPrefix(path, "/health") ||
		strings.HasPrefix(path, "/debug") ||
		strings.HasPrefix(path, "/metrics")
}
