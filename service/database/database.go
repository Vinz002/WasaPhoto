package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Vinz002/WASAPhoto/service/structs"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	CreateUser(user structs.User) (structs.User, error)
	FindUserByUserId(userId string) (bool, error)
	UpdateUsername(userId uint64, newUsername string) error
	FindUserByUsername(username string) (bool, error)
	SearchUser(userId string, search string) ([]structs.User, error)
	UplaodPhoto(photo structs.Photo) (structs.Photo, error)
	FindFollow(userId uint64, fluid uint64) (bool, error)
	FollowUser(userId uint64, fluid uint64) error
	UnfollowUser(userId uint64, fluid uint64) error
	GetFollowers(userId uint64) ([]string, error)
	GetFollowing(userId uint64) ([]string, error)
	GetUsername(userId string) (string, error)
	GetUser(userId string) (structs.User, error)
	BanUser(userId uint64, banId uint64) error
	CheckBan(userId uint64, banId uint64) (bool, error)
	UnBanUser(userId uint64, banId uint64) error
	LikePhoto(userId string, photoId string) error
	UnlikePhoto(userId string, photoId string) error
	CheckLike(userId string, photoId string) (bool, error)
	GetUsersAndLikes(photoId string) ([]string, error)
	FindPhotoByPhotoId(photoId string) (bool, error)
	GetPhoto(photoId string) (structs.Photo, error)
	DeletePhoto(userId string, photoId string) error
	CommentPhoto(userId string, photoId string, comment string) error
	GetUsersAndComments(photoId string) ([]structs.UserComment, error)
	CheckCommentUserId(comment_id string, user_id string) (bool, error)
	UncommentPhoto(comment_id string, photoId string) error
	GetUserPhotos(userId string) ([]structs.Photo, error)
	GetMyStream(userId string) ([]structs.Photo, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		createTableSQL := `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			photo_count INTEGER NOT NULL DEFAULT 0,
			follower_count INTEGER NOT NULL DEFAULT 0,
			following_count INTEGER NOT NULL DEFAULT 0
		);
	`
		_, err := db.Exec(createTableSQL)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='photos';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		createTableSQL := `
		CREATE TABLE photos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			filename TEXT NOT NULL,
			image BLOB NOT NULL,
			date_uploaded DATETIME NOT NULL,
			num_likes INTEGER NOT NULL DEFAULT 0,
			num_comments INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
		_, err := db.Exec(createTableSQL)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='likes';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		createTableSQL := `
		CREATE TABLE likes (
			user_id INTEGER NOT NULL,
			photo_id INTEGER NOT NULL,
			PRIMARY KEY (user_id, photo_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (photo_id) REFERENCES photos(id)
		);
	`
		_, err := db.Exec(createTableSQL)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='comments';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		createTableSQL := `
		CREATE TABLE comments (
			comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			photo_id INTEGER NOT NULL,
			text TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (photo_id) REFERENCES photos(id)
		);
	`
		_, err := db.Exec(createTableSQL)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='followers';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		createTableSQL := `
		CREATE TABLE followers (
			user_id INTEGER NOT NULL,
			follower_id INTEGER NOT NULL,
			PRIMARY KEY (user_id, follower_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (follower_id) REFERENCES users(id)
		);
	`
		_, err := db.Exec(createTableSQL)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='banned';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		createTableSQL := `
		CREATE TABLE banned (
			user_id INTEGER NOT NULL,
			banned_id INTEGER NOT NULL,
			PRIMARY KEY (user_id, banned_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (banned_id) REFERENCES users(id)
		);
	`
		_, err := db.Exec(createTableSQL)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
