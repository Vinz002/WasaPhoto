package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) FindFollow(userId uint64, fluid uint64) (bool, error) {
	query := "SELECT * FROM followers WHERE user_id = ? AND follower_id = ?"
	err := db.c.QueryRow(query, userId, fluid).Scan(&userId, &fluid)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *appdbimpl) FollowUser(userId uint64, fluid uint64) error {
	_, err := db.c.Exec("INSERT INTO followers (user_id, follower_id) VALUES (?, ?)", userId, fluid)
	if err != nil {
		return err
	}
	// Update the number of followers of the fluid
	_, err = db.c.Exec("UPDATE users SET follower_count = follower_count + 1 WHERE id = ?", fluid)
	if err != nil {
		return err
	}
	// Update the number of following of the user
	_, err = db.c.Exec("UPDATE users SET following_count = following_count + 1 WHERE id = ?", userId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetFollowers(userId uint64) ([]string, error) {
	var followers []string
	query := "SELECT u.username AS follower_username FROM users u JOIN followers f ON u.id = f.follower_id WHERE f.user_id = ?"
	rows, err := db.c.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var follower string
		err := rows.Scan(&follower)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return followers, nil
}

func (db *appdbimpl) GetFollowing(userId uint64) ([]string, error) {
	var following []string
	query := "SELECT u.username AS following_username FROM users u JOIN followers f ON u.id = f.user_id WHERE f.follower_id = ?"
	rows, err := db.c.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var follow string
		err := rows.Scan(&follow)
		if err != nil {
			return nil, err
		}
		following = append(following, follow)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return following, nil
}

func (db *appdbimpl) GetUsername(userId string) (string, error) {
	var username string
	err := db.c.QueryRow("SELECT username FROM users WHERE id = ?", userId).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
