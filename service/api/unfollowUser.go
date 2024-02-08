package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// unfollowAnUser is the handler for DELETE /user/:userId/follow/:followId
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Get the token from the Authorization header
	token := r.Header.Get("Authorization")
	if !isAuth(rt, w, token) {
		return
	}

	// Get the userId and the fluid parameter from the URL
	userId := ps.ByName("userId")
	fluid := ps.ByName("fluid")

	if userId == "" {
		// handle the error
		http.Error(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}
	if fluid == "" {
		// handle the error
		http.Error(w, "Missing fluid parameter", http.StatusBadRequest)
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

	// Check if the fluid exists
	exists, err = rt.db.FindUserByUserId(fluid)
	if err != nil {
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Fluid does not exist", http.StatusNotFound)
		return
	}

	// Check if the user is following the fluid
	Id, _ := strconv.ParseUint(userId, 10, 64)
	fluidId, _ := strconv.ParseUint(fluid, 10, 64)

	exists, err = rt.db.FindFollow(Id, fluidId)
	if err != nil {
		// handle the error
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "User is not following", http.StatusBadRequest)
		return
	}

	// Make the first user unfollow the second one
	err = rt.db.UnfollowUser(Id, fluidId)
	if err != nil {
		http.Error(w, "Unable to unfollow user", http.StatusInternalServerError)
		return
	}

	// Get the username of the fluid
	var username_f string
	username_f, err = rt.db.GetUsername(fluid)
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

	// Send the response
	w.WriteHeader(http.StatusNoContent)
	_, err = w.Write([]byte(`{"message": "` + username_u + ` is not following ` + username_f + ` anymore ` + `"}` + "\n"))
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
		return
	}
}
