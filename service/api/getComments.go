package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// get the token from the Authorization header
	token := r.Header.Get("Authorization")
	if !isAuth(rt, w, token) {
		return
	}

	photoId := p.ByName("photoId")
	if photoId == "" {
		http.Error(w, "Missing photoId parameter", http.StatusBadRequest)
		return
	}

	comments, err := rt.db.GetUsersAndComments(photoId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
