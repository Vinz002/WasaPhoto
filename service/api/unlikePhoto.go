package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// unlikePhoto is the handler for DELETE /user/:userId/like/:photoId
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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
	exist, err := rt.db.FindUserByUserId(userId)
	if err != nil {
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}
	if !exist {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	// Check if the photo exists
	ph_exist, err := rt.db.FindPhotoByPhotoId(photoId)
	if err != nil {
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}
	if !ph_exist {
		http.Error(w, "Photo does not exist", http.StatusNotFound)
		return
	}

	// Check if the user has already liked the photo
	liked, err := rt.db.CheckLike(userId, photoId)
	if err != nil {
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}

	if !liked {
		http.Error(w, "User has not liked the photo", http.StatusBadRequest)
		return
	}

	// Unlike the photo
	err = rt.db.UnlikePhoto(userId, photoId)
	if err != nil {
		http.Error(w, "Unable to unlike photo", http.StatusInternalServerError)
		return
	}

	// Get the photo
	photo, err := rt.db.GetPhoto(photoId)
	if err != nil {
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusNoContent)
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		http.Error(w, "Unable to encode the response", http.StatusInternalServerError)
		return
	}

}
