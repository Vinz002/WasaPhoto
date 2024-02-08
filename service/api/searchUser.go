package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// SearchUser is a the handler for /user/profile/:userId/search/:search API.
// It searches for users by name.

func (rt *_router) searchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get user ID
	userId := ps.ByName("userId")
	// check if user is authenticated
	if !isAuth(rt, w, userId) {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	// Get search string
	search := ps.ByName("search")
	// Get users
	users, err := rt.db.SearchUser(userId, search)
	if err != nil {
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}

	// Build the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "Unable to encode JSON", http.StatusInternalServerError)
		return
	}
}
