package database

func (db *appdbimpl) UnfollowUser(userId uint64, fluid uint64) error {
	_, err := db.c.Exec("DELETE FROM followers WHERE user_id = ? AND follower_id = ?", userId, fluid)
	if err != nil {
		return err
	}
	// Update the number of followers of the fluid
	_, err = db.c.Exec("UPDATE users SET follower_count = follower_count - 1 WHERE id = ?", fluid)
	if err != nil {
		return err
	}
	// Update the number of following of the user
	_, err = db.c.Exec("UPDATE users SET following_count = following_count - 1 WHERE id = ?", userId)
	if err != nil {
		return err
	}
	return nil
}
