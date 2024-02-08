package database

// unlikePhoto remove the like from the photo
func (db *appdbimpl) UnlikePhoto(userId string, photoId string) error {

	// Unlike the photo
	_, err := db.c.Exec(`DELETE FROM likes WHERE user_id=? AND photo_id=?`, userId, photoId)
	if err != nil {
		return err
	}

	// Decrement the number of likes from the photo
	_, err = db.c.Exec(`UPDATE photos SET num_likes=num_likes-1 WHERE id=?`, photoId)
	if err != nil {
		return err
	}

	return nil
}
