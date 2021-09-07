package middleware

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"github.com/is0405/httputil"
	"github.com/is0405/model"
	"github.com/jmoiron/sqlx"
)

const (
	bearer = "Bearer"
)

type Auth struct {
	jwtSecretKey []byte
	db           *sqlx.DB
}

func NewAuth(jwtSecretKey []byte, db *sqlx.DB) *Auth {
	return &Auth{
		jwtSecretKey: jwtSecretKey,
		db:           db,
	}
}

func (auth *Auth) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken, err := getTokenFromHeader(r)
		if err != nil {
			httputil.RespondErrorJson(w, http.StatusBadRequest, err)
			return
		}

		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(idToken, claims, func(tkn *jwt.Token) (interface{}, error) {
			return auth.jwtSecretKey, nil
		})
		if err != nil {
			httputil.RespondErrorJson(w, http.StatusBadRequest, err)
			return
		}
		if !token.Valid {
			httputil.RespondErrorJson(w, http.StatusUnauthorized, err)
			return
		}
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				httputil.RespondErrorJson(w, http.StatusUnauthorized, err)
				return
			}
			httputil.RespondErrorJson(w, http.StatusBadRequest, err)
			return
		}

		ctx := httputil.SetClaimsToContext(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTokenFromHeader(req *http.Request) (string, error) {
	header := req.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("authorization header not found")
	}

	l := len(bearer)
	if len(header) > l+1 && header[:l] == bearer {
		return header[l+1:], nil
	}

	return "", errors.New("authorization header format must be 'Bearer {token}'")
}
