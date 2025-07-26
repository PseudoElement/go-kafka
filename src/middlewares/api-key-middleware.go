package middlewares

import (
	"net/http"
	"os"

	"github.com/pseudoelement/go-kafka/src/shared"
)

func XApiTokenMiddleware(next http.Handler) http.Handler {
	xApiToken := os.Getenv("X_API_TOKEN")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-token") != xApiToken {
			shared.FailedResp(w, "Unauthorized. Header 'x-api-token' is not provided.", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
