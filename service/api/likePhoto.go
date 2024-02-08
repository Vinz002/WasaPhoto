package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// likePhoto is the handler for PUT /user/:userId/like/:photoId
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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
		http.Error(w, "Unable to query the database 1", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	// Check if the photo exists
	ph_exists, err := rt.db.FindPhotoByPhotoId(photoId)
	if err != nil {
		http.Error(w, "Unable to query the database 2", http.StatusInternalServerError)
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
		http.Error(w, "User has already liked the photo", http.StatusBadRequest)
		return
	}

	// Like the photo
	err = rt.db.LikePhoto(userId, photoId)
	if err != nil {
		http.Error(w, "Unable to like photo", http.StatusInternalServerError)
		return
	}

	// Get the photo
	photo, err := rt.db.GetPhoto(photoId)
	if err != nil {
		http.Error(w, "Unable to query the database 4", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}
