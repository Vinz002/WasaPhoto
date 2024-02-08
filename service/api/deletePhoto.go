package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// deletePhoto is the handler for DELETE /user/:userId/photos/:photoId
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	// Check if the photo exists
	exists, err = rt.db.FindPhotoByPhotoId(photoId)
	if err != nil {
		http.Error(w, "Unable to query database", http.StatusNotFound)
		return
	}

	if !exists {
		http.Error(w, "Photo does not exist", http.StatusNotFound)
		return
	}

	// Delete the photo
	err = rt.db.DeletePhoto(userId, photoId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusNoContent)
	_, err = w.Write([]byte("Photo deleted successfully"))
	if err != nil {
		http.Error(w, "Unable to write the response", http.StatusInternalServerError)
		return
	}
}
