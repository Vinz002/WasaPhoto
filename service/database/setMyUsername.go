package database

import (
	"database/sql"
	"errors"
	"github.com/Vinz002/WASAPhoto/service/structs"
)

func (db *appdbimpl) FindUserByUserId(userId string) (bool, error) {
	query := "SELECT id FROM users WHERE id = ?"
	var userID uint64
	err := db.c.QueryRow(query, userId).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, err
}

func (db *appdbimpl) UpdateUsername(userId uint64, newUsername string) error {
	query := "UPDATE users SET username = ? WHERE id = ?"
	_, err := db.c.Exec(query, newUsername, userId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) FindUserByUsername(username string) (bool, error) {
	query := "SELECT id FROM users WHERE username = ?"
	var userID int
	err := db.c.QueryRow(query, username).Scan(&userID)
	if err != nil {
		return false, nil
	}
	return true, err
}

// GetUser returns the user with the given id
func (db *appdbimpl) GetUser(userId string) (structs.User, error) {
	var user structs.User
	err := db.c.QueryRow("SELECT * FROM users WHERE id = ?", userId).Scan(&user.Id, &user.Username, &user.Photo_count, &user.Follower_count, &user.Following_count)
	if err != nil {
		return structs.User{}, err
	}
	return user, nil
}
