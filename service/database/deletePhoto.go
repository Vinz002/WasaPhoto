package database

// DeletePhoto removes a photo from the database
func (db *appdbimpl) DeletePhoto(userId string, photoId string) error {
	deleteSQL := "DELETE FROM photos WHERE id = ?"
	_, err := db.c.Exec(deleteSQL, photoId)
	if err != nil {
		return err
	}
	update_counter := "UPDATE users SET photo_count = photo_count - 1 WHERE id = ?"
	_, err = db.c.Exec(update_counter, userId)
	if err != nil {
		return err
	}
	return nil
}
