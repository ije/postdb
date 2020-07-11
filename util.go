package postdb

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/ije/gox/utils"
)

// BasicAuth imlements a simple basic auth
func BasicAuth(username string, password string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value := r.Header.Get("Authorization")
			if strings.HasPrefix(value, "Basic ") {
				authInfo, err := base64.StdEncoding.DecodeString(value[6:])
				if err == nil {
					name, pass := utils.SplitByFirstByte(string(authInfo), ':')
					if name == username && pass == password {
						h.ServeHTTP(w, r)
						return
					}
				}
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="Authorization Required"`)
			w.WriteHeader(401)
		})
	}
}

func toLowerTrim(s string) string {
	return strings.ToLower(strings.TrimSpace(string(s)))
}
