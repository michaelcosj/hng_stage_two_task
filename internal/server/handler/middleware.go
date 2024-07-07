package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/michaelcosj/hng-task-two/internal/app"
)

func StripSlashes(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > 1 && path[len(path)-1] == '/' {
			newPath := path[:len(path)-1]
			r.URL.Path = newPath
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			errResp := map[string]any{
				"status":     "Unauthorized",
				"statusCode": http.StatusUnauthorized,
				"message":    "Not authorised to access this resource",
			}
			writeJSON(w, http.StatusUnauthorized, errResp)
			return
		}

		tokenString = tokenString[len("Bearer "):]
		token, err := app.VerifyToken(tokenString)
		if err != nil {
			log.Printf("error validating jwt token: %v", err)

			errResp := map[string]any{
				"status":     "Invalid jwt token",
				"statusCode": http.StatusUnauthorized,
				"message":    "JWT token is invalid or expired",
			}
			writeJSON(w, http.StatusUnauthorized, errResp)
			return
		}

		mapClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			errResp := map[string]any{
				"status":     "Invalid jwt token",
				"statusCode": http.StatusUnauthorized,
				"message":    "JWT token is invalid or expired",
			}
			writeJSON(w, http.StatusUnauthorized, errResp)
			return
		}

		userId := mapClaims["id"].(string)
		ctx := context.WithValue(r.Context(), "userId", userId)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
