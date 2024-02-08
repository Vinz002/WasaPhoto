package api

import (
	"net/http"
)

func isAuth(rt *_router, w http.ResponseWriter, token string) bool {
	if token == "" || token == "null" {
		http.Error(w, "Missing Authorization header", http.StatusBadRequest)
		return false
	} else {
		// Check if the token is valid
		_, err := rt.db.FindUserByUserId(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return false
		}
	}
	return true
}
