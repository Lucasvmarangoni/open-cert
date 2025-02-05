package middlewares

import (
	"net/http"

	"github.com/go-chi/jwtauth"
)

func getSub(w http.ResponseWriter, r *http.Request) string {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "Failed to get JWT claims", http.StatusInternalServerError)
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		http.Error(w, "Sub claim is missing or not a string", http.StatusInternalServerError)
	}
	return sub
}