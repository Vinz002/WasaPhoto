package database

import (
	"database/sql"
	"errors"
)

// LikePhoto add a like to a photo
func (db *appdbimpl) LikePhoto(userId string, photoId string) error {

	// Insert the like
	_, err := db.c.Exec(`INSERT INTO likes (user_id, photo_id) VALUES (?, ?)`, userId, photoId)
	if err != nil {
		return err
	}
	// Update the number of likes
	_, err = db.c.Exec(`UPDATE photos SET num_likes = num_likes + 1 WHERE id = ?`, photoId)
	if err != nil {
		return err
	}
	return nil
}

// CheckLike checks if a user has already liked a photo
func (db *appdbimpl) CheckLike(userId string, photoId string) (bool, error) {
	var id string
	err := db.c.QueryRow("SELECT user_id FROM likes WHERE user_id = ? AND photo_id = ?", userId, photoId).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, nil
}

// GetUsersAndLikes returns the usernames of the users that liked a photo
func (db *appdbimpl) GetUsersAndLikes(photoId string) ([]string, error) {

	// Get the likes
	rows, err := db.c.Query(`
		SELECT u.username FROM likes l
		INNER JOIN users u ON l.user_id = u.user_id
		WHERE l.photo_id = ?`, photoId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Build the likes list
	var likes []string
	for rows.Next() {
		var like string
		err = rows.Scan(&like)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return likes, nil
}
