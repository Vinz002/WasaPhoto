package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// UpdateUsername is the handler for PUT /user/:userId
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Get the token from the Authorization header
	token := r.Header.Get("Authorization")
	if !isAuth(rt, w, token) {
		return
	}

	// Get the userId parameter from the URL
	userId := ps.ByName("userId")
	if userId == "" {
		// handle the error
		http.Error(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}
	// Check if the user id exists in the database
	Id, _ := strconv.ParseUint(userId, 10, 64)

	exists, err := rt.db.FindUserByUserId(userId)
	if err != nil {
		http.Error(w, "Unable to query database", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}

	// Get the new username parameter from the request body
	var newUsername string
	err = json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Check if the new username already exists in the database
	exists, err = rt.db.FindUserByUsername(newUsername)
	if err != nil {
		http.Error(w, "Unable to query database", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}
	// Update the username in the database
	err = rt.db.UpdateUsername(Id, newUsername)
	if err != nil {
		http.Error(w, "Unable to update username", http.StatusInternalServerError)
		return
	}
	// Return a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("Username updated successfully\n"))
	if err != nil {
		http.Error(w, "Unable to write response body", http.StatusInternalServerError)
		return
	}
}
