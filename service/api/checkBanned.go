package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// checkBanned is the API handler to check if a user is banned by another user
func (rt *_router) checkBanned(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// check if the user is authenticated
	token := r.Header.Get("Authorization")
	if !isAuth(rt, w, token) {
		return
	}

	// Get the user ID from the URL
	userId := p.ByName("userId")
	banId := p.ByName("banId")

	if userId == "" {
		// handle the error
		http.Error(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}
	if banId == "" {
		// handle the error
		http.Error(w, "Missing fluid parameter", http.StatusBadRequest)
		return
	}

	exists, err := rt.db.FindUserByUserId(userId)
	if err != nil {
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	exists, err = rt.db.FindUserByUserId(banId)
	if err != nil {
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "banId does not exist", http.StatusNotFound)
		return
	}

	us, _ := strconv.ParseUint(userId, 10, 64)
	ban, _ := strconv.ParseUint(banId, 10, 64)
	// Check if the user is banned by the other user
	banned, err := rt.db.CheckBan(us, ban)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if banned {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
