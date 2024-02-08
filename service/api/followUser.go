package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// followAnUser is the handler for PUT /user/:userId/follow/:followId
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

	Id, _ := strconv.ParseUint(userId, 10, 64)
	fluidId, _ := strconv.ParseUint(fluid, 10, 64)

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

	// Check if the user is banned by the fluid
	banned, err := rt.db.CheckBan(fluidId, Id)
	if err != nil {
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}
	if banned {
		http.Error(w, "Unable to follow an User that has banned you", http.StatusBadRequest)
		return
	}

	// Check if the user is already following the fluid
	exists, err = rt.db.FindFollow(Id, fluidId)
	if err != nil {
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "User is already following", http.StatusBadRequest)
		return
	}

	// Make the first user follow the second one
	err = rt.db.FollowUser(Id, fluidId)
	if err != nil {
		http.Error(w, "Unable to follow user", http.StatusInternalServerError)
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

	// Return the response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write([]byte(`{"status": "success", "message": "User ` + username_u + ` is now following user ` + username_f + `"}` + "\n"))
	if err != nil {
		http.Error(w, "Unable to write the response", http.StatusInternalServerError)
		return
	}
}
