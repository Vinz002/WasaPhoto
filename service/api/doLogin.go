package api

import (
	"encoding/json"
	"github.com/Vinz002/WASAPhoto/service/structs"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user structs.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err = rt.db.CreateUser(user)
	if err != nil {
		http.Error(w, "Unable to create User", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Unable to encode User to JSON", http.StatusInternalServerError)
		return
	}
}
