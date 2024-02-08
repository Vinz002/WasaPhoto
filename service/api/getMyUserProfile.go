package api

import (
	"encoding/json"
	"github.com/Vinz002/WASAPhoto/service/structs"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// getUserProfile is the handler for GET /user/:userId/profile/:profileId
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Get the token from the Authorization header
	token := r.Header.Get("Authorization")
	if !isAuth(rt, w, token) {
		return
	}

	// Parse the request
	userId := ps.ByName("userId")
	profileId := ps.ByName("profileId")

	if userId == "" {
		// handle the error
		http.Error(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}

	if profileId == "" {
		// handle the error
		http.Error(w, "Missing profileId parameter", http.StatusBadRequest)
		return
	}

	// Check if the user exists
	exists, err := rt.db.FindUserByUserId(profileId)
	if err != nil {
		http.Error(w, "Unable to query the database 1", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	user, err := rt.db.GetUser(profileId)
	if err != nil {
		http.Error(w, "Unable to query the database 2", http.StatusInternalServerError)
		return
	}

	// Get the user's photos
	photos, err := rt.db.GetUserPhotos(profileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the response
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
