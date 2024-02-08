package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// this is the handler for the API GET /user/:userId//follow/:fluid
func (rt *_router) checkFollower(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// get the token from the Authorization header
	token := r.Header.Get("Authorization")
	if !isAuth(rt, w, token) {
		return
	}

	// get the userId from the URL
	userId := ps.ByName("userId")

	// get the fluid from the URL
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

	user, _ := strconv.ParseUint(userId, 10, 64)
	followid, _ := strconv.ParseUint(fluid, 10, 64)

	// check if the user is following the fluid
	isFollowing, err := rt.db.FindFollow(user, followid)
	if err != nil {
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
	}

	// return the result
	if isFollowing {
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
