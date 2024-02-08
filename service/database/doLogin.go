package database

import (
	"database/sql"
	"errors"
	"github.com/Vinz002/WASAPhoto/service/structs"
)

func (db *appdbimpl) CreateUser(user structs.User) (structs.User, error) {
	query := "SELECT id, photo_count, follower_count, following_count FROM users WHERE username = ?"
	var id int64
	var photo_count int
	var follower_count int
	var following_count int
	err := db.c.QueryRow(query, user.Username).Scan(&id, &photo_count, &follower_count, &following_count)

	if !errors.Is(err, sql.ErrNoRows) {
		// User already exists so return the existing user in the database
		user.Id = uint64(id)
		user.Photo_count = photo_count
		user.Follower_count = follower_count
		user.Following_count = following_count
		return user, nil
	}
	insertSQL := "INSERT INTO users (username) VALUES (?)"
	result, err := db.c.Exec(insertSQL, user.Username)
	if err != nil {
		return structs.User{}, err
	}
	id, err = result.LastInsertId()
	if err != nil {
		return structs.User{}, err
	}
	user.Id = uint64(id)
	return user, nil

}
