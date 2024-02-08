package database

// UnBanUser makes the user with the id userId unban the user with the id banId
func (db *appdbimpl) UnBanUser(userId uint64, banId uint64) error {
	_, err := db.c.Exec("DELETE FROM banned WHERE user_id=? AND banned_id=?", userId, banId)
	if err != nil {
		return err
	}
	return nil
}
