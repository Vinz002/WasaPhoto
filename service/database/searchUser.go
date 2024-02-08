package database

import (
	"database/sql"
	"errors"
	"github.com/Vinz002/WASAPhoto/service/structs"
)

// SearchUser searches for users by name and is going to be used in a search bar in the app but dont show user that has banned the userId

func (db *appdbimpl) SearchUser(userId string, search string) ([]structs.User, error) {

	// Get the users
	rows, err := db.c.Query(`
		SELECT id, username, photo_count, follower_count, following_count FROM users
		WHERE username LIKE ? AND id NOT IN (SELECT user_id FROM banned WHERE banned_id = ?) `, "%"+search+"%", userId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	defer rows.Close()
	// Build the users list
	var users []structs.User
	for rows.Next() {
		var user structs.User
		err = rows.Scan(&user.Id, &user.Username, &user.Photo_count, &user.Follower_count, &user.Following_count)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}
