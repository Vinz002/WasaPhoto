package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// unbanUser is the handler for DELETE /user/:userId/ban/:banId
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Get the token from the Authorization header
	token := r.Header.Get("Authorization")
	if !isAuth(rt, w, token) {
		return
	}

	// Get the userId and the banId parameter from the URL
	userId := ps.ByName("userId")
	banId := ps.ByName("banId")

	if userId == "" {
		// handle the error
		http.Error(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}
	if banId == "" {
		// handle the error
		http.Error(w, "Missing banId parameter", http.StatusBadRequest)
		return
	}

	// Check if the user exists
	exists, err := rt.db.FindUserByUserId(userId)
	if err != nil {
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	// Check if the ban exists
	exists, err = rt.db.FindUserByUserId(banId)
	if err != nil {
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Ban does not exist", http.StatusNotFound)
		return
	}

	// Check if the user is already following the ban
	Id, _ := strconv.ParseUint(userId, 10, 64)
	BanId, _ := strconv.ParseUint(banId, 10, 64)

	exists, err = rt.db.CheckBan(Id, BanId)
	if err != nil {
		// handle the error
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "User is not banned", http.StatusBadRequest)
		return
	}

	// Make the first user unban the second one
	err = rt.db.UnBanUser(Id, BanId)
	if err != nil {
		http.Error(w, "Unable to unban user", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusNoContent)
	_, err = w.Write([]byte("User unbanned successfully"))
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
		return
	}
}
