package database

// UncommentPhoto removes a comment from a photo
func (db *appdbimpl) UncommentPhoto(comment_id string, photoId string) error {

	// Remove the comment
	_, err := db.c.Exec(`DELETE FROM comments WHERE comment_id=?`, comment_id)
	if err != nil {
		return err
	}

	// Update the photo
	_, err = db.c.Exec(`UPDATE photos SET num_comments = num_comments - 1 WHERE id = ?`, photoId)
	if err != nil {
		return err
	}
	return nil
}

// CheckCommentUserId checks if the comment comment_id has been commented by the user user_id
func (db *appdbimpl) CheckCommentUserId(comment_id string, user_id string) (bool, error) {

	// Get the comment
	var commentUserId string
	err := db.c.QueryRow(`SELECT user_id FROM comments WHERE comment_id=?`, comment_id).Scan(&commentUserId)
	if err != nil {
		return false, err
	}

	if commentUserId == user_id {
		return true, nil
	}

	return false, nil
}
