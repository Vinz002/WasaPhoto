package database

import (
	"github.com/Vinz002/WASAPhoto/service/structs"
	"time"
)

// UplaodPhoto inserts a new photo into the database
func (db *appdbimpl) UplaodPhoto(photo structs.Photo) (structs.Photo, error) {
	query := "INSERT INTO photos (user_id, filename, image, date_uploaded) VALUES (?, ?, ?, ?)"
	result, err := db.c.Exec(query, photo.UserID, photo.Filename, photo.ImageData, time.Now().UTC())
	if err != nil {
		return structs.Photo{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return structs.Photo{}, err
	}
	photo.Id = uint64(id)
	// increment the number of photos uploaded by the user
	query = "UPDATE users SET photo_count = photo_count + 1 WHERE id = ?"
	_, err = db.c.Exec(query, photo.UserID)
	if err != nil {
		return structs.Photo{}, err
	}
	return photo, nil
}

// FindPhotoByPhotoId returns true if the photo with the given id exists
func (db *appdbimpl) FindPhotoByPhotoId(photoId string) (bool, error) {
	var id string
	err := db.c.QueryRow("SELECT id FROM photos WHERE id = ?", photoId).Scan(&id)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetPhoto returns the photo with the given id
func (db *appdbimpl) GetPhoto(photoId string) (structs.Photo, error) {
	var photo structs.Photo
	err := db.c.QueryRow("SELECT * FROM photos WHERE id = ?", photoId).Scan(&photo.Id, &photo.UserID, &photo.Filename, &photo.ImageData, &photo.DateUploaded, &photo.NumLikes, &photo.NumComments)
	if err != nil {
		return structs.Photo{}, err
	}
	return photo, nil
}

// GetUserPhotos returns the photos of a user sorted in reverse chronological order
func (db *appdbimpl) GetUserPhotos(userId string) ([]structs.Photo, error) {
	var photos []structs.Photo
	rows, err := db.c.Query("SELECT * FROM photos WHERE user_id = ? ORDER BY date_uploaded DESC", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var photo structs.Photo
		err := rows.Scan(&photo.Id, &photo.UserID, &photo.Filename, &photo.ImageData, &photo.DateUploaded, &photo.NumLikes, &photo.NumComments)
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
