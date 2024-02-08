package api

import (
	"encoding/json"
	"github.com/Vinz002/WASAPhoto/service/structs"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// getMyStream is the handler for GET /user/stream/:userId
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Get the token from the Authorization header
	token := r.Header.Get("Authorization")
	if !isAuth(rt, w, token) {
		return
	}

	// Parse the request
	userId := ps.ByName("userId")

	if userId == "" {
		http.Error(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}

	// Check if the user exists
	exists, err := rt.db.FindUserByUserId(userId)
	if err != nil {
		http.Error(w, "Unable to query the database 1", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}

	user, err := rt.db.GetUser(userId)
	if err != nil {
		http.Error(w, "Unable to query the database 2", http.StatusInternalServerError)
		return
	}

	// Get the user's stream
	photos, err := rt.db.GetMyStream(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := structs.UserProfile{
		Username:        user.Username,
		Photo_count:     user.Photo_count,
		Follower_count:  user.Follower_count,
		Following_count: user.Following_count,
		Photos:          photos,
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Unable to encode the response", http.StatusInternalServerError)
		return
	}
}
