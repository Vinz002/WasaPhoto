package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// banAnUser is the handler for DELETE /user/:userId/ban/:banId
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

	// Check if the userId is already following the banId or vice versa and if so, make them unfollow each other
	Id, _ := strconv.ParseUint(userId, 10, 64)
	BanId, _ := strconv.ParseUint(banId, 10, 64)

	exists_f1, err := rt.db.FindFollow(Id, BanId)
	if err != nil {
		// handle the error
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
	}

	exists_f2, err := rt.db.FindFollow(BanId, Id)
	if err != nil {
		// handle the error
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}

	if exists_f1 {
		err = rt.db.UnfollowUser(Id, BanId)
		if err != nil {
			http.Error(w, "Unable to unfollow user", http.StatusInternalServerError)
			return
		}
	}

	if exists_f2 {
		err = rt.db.UnfollowUser(BanId, Id)
		if err != nil {
			http.Error(w, "Unable to unfollow user", http.StatusInternalServerError)
			return
		}
	}

	// Make the first user ban the second one
	err = rt.db.BanUser(Id, BanId)
	if err != nil {
		http.Error(w, "Unable to ban user", http.StatusInternalServerError)
		return
	}

	// Get the username of the ban
	var username_b string
	username_b, err = rt.db.GetUsername(banId)
	if err != nil {
		http.Error(w, "Unable to get username", http.StatusInternalServerError)
		return
	}
	// Get the username of the user
	var username_u string
	username_u, err = rt.db.GetUsername(userId)
	if err != nil {
		http.Error(w, "Unable to get username", http.StatusInternalServerError)
		return
	}

	// Return a 201 status created
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte("User " + username_u + " has banned " + username_b + "\n"))
	if err != nil {
		http.Error(w, "Unable to write body", http.StatusInternalServerError)
		return
	}
}
