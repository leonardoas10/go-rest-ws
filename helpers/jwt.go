package helpers

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/leonardoas10/go-rest-ws/models"
	"github.com/leonardoas10/go-rest-ws/server"
)

func GetToken(s server.Server, w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
	token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func (token *jwt.Token) (interface{}, error)  {
		return []byte(s.Config().JWTSecret), nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return nil, err
	}
	return token, nil
}