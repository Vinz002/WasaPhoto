package api

import (
	"encoding/json"
	"github.com/Vinz002/WASAPhoto/service/structs"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
)

// ...

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get the token from the Authorization header
	token := r.Header.Get("Authorization")
	if !isAuth(rt, w, token) {
		return
	}

	// Get the user ID from the URL
	userId := ps.ByName("userId")
	if userId == "" {
		http.Error(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}

	// Check if the user exists
	exists, err := rt.db.FindUserByUserId(userId)
	if err != nil {
		http.Error(w, "Unable to query database", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}

	// Parse multipart form data
	err = r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file content into a byte slice
	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading the file content", http.StatusInternalServerError)
		return
	}

	// Now you have the file content in the 'fileContent' variable
	// You can do further processing or save it to your database

	// Create a new photo struct
	var photo structs.Photo
	photo.UserID, _ = strconv.ParseUint(userId, 10, 64)
	photo.Filename = handler.Filename
	photo.ImageData = fileContent

	// Upload the photo (you might want to implement your own logic here)
	photo, err = rt.db.UplaodPhoto(photo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the photo details
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		http.Error(w, "Unable to encode photo to JSON", http.StatusInternalServerError)
		return
	}
}
