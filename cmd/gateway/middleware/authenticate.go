package middleware

import (
	"dtms/cmd/gateway/handler"
	"dtms/cmd/gateway/types"
	"encoding/json"
	"net/http"
)

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !handler.Authenticate(r) {
			r.Body.Close()

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(types.AuthMessage{
				Id:    "-1",
				Token: "-1",
			})
			return
		}

		h.ServeHTTP(w, r)
	})
}
