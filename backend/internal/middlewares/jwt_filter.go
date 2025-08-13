package middlewares

import (
	"backend/global"
	"backend/internal/runtime_errors"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

func JwtFilter(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		authHeader := req.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader,"Bearer ") {
			http.Error(w,"Bearer token not found.",http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader,"Bearer ")

		keySecret,ok := os.LookupEnv("KEY_SECRET")

		if !ok {
			http.Error(w,"Key secret not found",http.StatusInternalServerError)
			return
		}

		jwtToken,err := jwt.Parse(token,func(t *jwt.Token) (any, error) {
			if _,ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])

			}
			return []byte(keySecret),nil
		})


		if err!=nil{
			global.HandleError(&runtime_errors.UnauthorizedError{
				Message: "Invalid token: "+err.Error(),
			},w)
			return
		}

		claims,ok := jwtToken.Claims.(jwt.MapClaims)

		if !ok {
			http.Error(w, "Failed to extract claims", http.StatusUnauthorized)
			return
		}

		exp,ok := claims["exp"].(float64)
		if !ok {
			http.Error(w,"Couldnt parse expiration time",http.StatusInternalServerError)
			return
		}

		if time.Now().Unix() > int64(exp) {
			http.Error(w,"Token expired",http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), ClaimsKey, claims)

		next.ServeHTTP(w,req.WithContext(ctx))
	})
}