package database

import (
	"github.com/Vinz002/WASAPhoto/service/structs"
)

// GetMyStream returns the photos of the users that the user follows sorted in reverse chronological order and select also the username of the user that uploaded the photo.
func (db *appdbimpl) GetMyStream(userId string) ([]structs.Photo, error) {
	rows, err := db.c.Query(`
		SELECT photos.id, photos.user_id, photos.image, photos.date_uploaded, users.username, photos.num_likes, photos.num_comments
		FROM photos
		INNER JOIN users
		ON photos.user_id = users.id
		WHERE photos.user_id IN (
			SELECT follower_id
			FROM followers
			WHERE user_id = ?
		)
		ORDER BY photos.date_uploaded DESC
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	photos := []structs.Photo{}
	for rows.Next() {
		var photo structs.Photo
		err = rows.Scan(&photo.Id, &photo.UserID, &photo.ImageData, &photo.DateUploaded, &photo.Username, &photo.NumLikes, &photo.NumComments)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return photos, nil
}
