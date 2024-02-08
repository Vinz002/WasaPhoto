package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// checkLiked is the handler for GET /users/:userId/likes/:photoId
func (rt *_router) checkLiked(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Get the token from the Authorization header
	token := r.Header.Get("Authorization")
	if !isAuth(rt, w, token) {
		return
	}

	// Parse the request
	userId := ps.ByName("userId")
	photoId := ps.ByName("photoId")

	if userId == "" {
		// handle the error
		http.Error(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}

	if photoId == "" {
		// handle the error
		http.Error(w, "Missing photoId parameter", http.StatusBadRequest)
		return
	}

	// Check if the user exists
	exists, err := rt.db.FindUserByUserId(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	// Check if the photo exists
	ph_exists, err := rt.db.FindPhotoByPhotoId(photoId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ph_exists {
		http.Error(w, "Photo does not exist", http.StatusNotFound)
		return
	}

	// Check if the user has already liked the photo
	liked, err := rt.db.CheckLike(userId, photoId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if liked {
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
