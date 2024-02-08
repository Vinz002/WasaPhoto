package database

import ()

// BanUser makes the user with the id userId ban the user with the id banId
func (db *appdbimpl) BanUser(userId uint64, banId uint64) error {
	_, err := db.c.Exec("INSERT INTO banned (user_id, banned_id) VALUES (?, ?)", userId, banId)
	if err != nil {
		return err
	}
	return nil
}

// Check ban makes the user with the id userId check if the user with the id banId is banned from him
func (db *appdbimpl) CheckBan(userId uint64, banId uint64) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM banned WHERE user_id=? AND banned_id=?)", userId, banId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
