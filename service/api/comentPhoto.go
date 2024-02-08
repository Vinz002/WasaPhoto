package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// CommentPhoto is the handler for PUT /user/:userId/photos/:photoId
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Get the token from the Authorization header
	token := r.Header.Get("Authorization")
	if !isAuth(rt, w, token) {
		return
	}

	// Parse the request
	userId := ps.ByName("userId")
	photoId := ps.ByName("photoId")

	// Decode the comment from the request
	var comment string
	err := json.NewDecoder(r.Body).Decode(&comment)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	if comment == "" {
		// handle the error
		http.Error(w, "Missing comment parameter", http.StatusBadRequest)
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

	// Comment the photo
	err = rt.db.CommentPhoto(userId, photoId, comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the photo
	photo, err := rt.db.GetPhoto(photoId)
	if err != nil {
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		http.Error(w, "Unable to encode the response", http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte("The user has commented the photo with the comment: " + comment + "\n"))
	if err != nil {
		http.Error(w, "Unable to write the response", http.StatusInternalServerError)
		return
	}
}
