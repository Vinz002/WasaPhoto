package database

import (
	"github.com/Vinz002/WASAPhoto/service/structs"
)

// CommentPhoto adds a comment to a photo
func (db *appdbimpl) CommentPhoto(userId string, photoId string, comment string) error {

	// Insert the comment
	_, err := db.c.Exec(`INSERT INTO comments (user_id, photo_id, text) VALUES (?, ?, ?)`, userId, photoId, comment)
	if err != nil {
		return err
	}

	// Update the photo
	_, err = db.c.Exec(`UPDATE photos SET num_comments = num_comments + 1 WHERE id = ?`, photoId)
	if err != nil {
		return err
	}

	return nil
}

// GetUsersAndComments returns the usernames and comments of a photo
// rows.err must be checked
func (db *appdbimpl) GetUsersAndComments(photoId string) ([]structs.UserComment, error) {
	// Get the comments
	rows, err := db.c.Query(`
		SELECT comments.comment_id, comments.user_id, users.username, comments.photo_id, comments.text FROM comments 
		INNER JOIN users ON comments.user_id = users.id
		WHERE comments.photo_id = ?`, photoId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Build the comments list
	var comments []structs.UserComment
	for rows.Next() {
		var comment structs.UserComment
		err = rows.Scan(&comment.Id, &comment.UserID, &comment.Username, &comment.PhotoId, &comment.Comment)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return comments, nil
}
