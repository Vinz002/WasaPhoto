package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// uncommentPhoto is the handler for DELETE /user/:userId/photos/:photoId/comment/:commentId
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Get the token from the Authorization header
	token := r.Header.Get("Authorization")
	if !isAuth(rt, w, token) {
		return
	}

	// Parse the request
	userId := ps.ByName("userId")
	photoId := ps.ByName("photoId")
	commentId := ps.ByName("commentId")

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

	if commentId == "" {
		// handle the error
		http.Error(w, "Missing commentId parameter", http.StatusBadRequest)
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

	// Check if the comment comment_Id has been commented by the user userId
	var commentUserId bool
	commentUserId, err = rt.db.CheckCommentUserId(commentId, userId)
	if err != nil {
		http.Error(w, "Unable to query the database", http.StatusInternalServerError)
		return
	}

	if !commentUserId {
		http.Error(w, "The comment has not been posted by the selected user", http.StatusBadRequest)
		return
	}

	// Uncomment the photo
	err = rt.db.UncommentPhoto(commentId, photoId)
	if err != nil {
		http.Error(w, "Unable to uncomment photo", http.StatusInternalServerError)
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

	_, err = w.Write([]byte(`{"message": "Comment removed successfully"}`))
	if err != nil {
		http.Error(w, "Unable to write the response", http.StatusInternalServerError)
		return
	}
}
