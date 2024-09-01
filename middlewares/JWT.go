package middlewares

import (
	"go-jwt-mux/config"
	"go-jwt-mux/helper"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				respon := map[string]string{"message": "Unauthorized"}
				helper.ResponseJSON(w, http.StatusUnauthorized, respon)
				return
			}
		}

		//ambil token
		tokenStr := c.Value

		claims := &config.JWTClaim{}
		//parse token
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid: //token invalid
				respon := map[string]string{"message": "Unauthorized"}
				helper.ResponseJSON(w, http.StatusUnauthorized, respon)
				return
			case jwt.ValidationErrorExpired: //token expired
				respon := map[string]string{"message": "Unauthorized, token expired"}
				helper.ResponseJSON(w, http.StatusUnauthorized, respon)
				return
			default:
				respon := map[string]string{"message": "Unauthorized"}
				helper.ResponseJSON(w, http.StatusUnauthorized, respon)

			}
		}
		if !token.Valid {
			respon := map[string]string{"message": "Unauthorized"}
			helper.ResponseJSON(w, http.StatusUnauthorized, respon)
			return
		}
		//next
		next.ServeHTTP(w, r)

	})
}
