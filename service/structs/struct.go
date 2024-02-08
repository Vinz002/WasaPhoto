package structs

type User struct {
	Id              uint64
	Username        string
	Photo_count     int
	Follower_count  int
	Following_count int
}

type Photo struct {
	Id           uint64
	UserID       uint64
	Username     string
	Filename     string
	ImageData    []byte
	DateUploaded string
	NumLikes     int
	NumComments  int
	Comments     []UserComment
}

type UserComment struct {
	Id       uint64
	UserID   uint64
	Username string
	PhotoId  uint64
	Comment  string
}

type UserProfile struct {
	Username        string
	Photo_count     int
	Follower_count  int
	Following_count int
	Photos          []Photo
}
