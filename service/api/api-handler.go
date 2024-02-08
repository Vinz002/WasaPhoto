package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.doLogin)
	rt.router.PUT("/user/:userId", rt.setMyUserName)
	rt.router.GET("/users/:userId/profile/:profileId", rt.getUserProfile)
	rt.router.GET("/user/stream/:userId", rt.getMyStream)
	rt.router.GET("/user/profile/:userId/search/:search", rt.searchUser)
	rt.router.POST("/user/:userId", rt.uploadPhoto)
	rt.router.DELETE("/user/:userId/photos/:photoId", rt.deletePhoto)
	rt.router.PUT("/user/:userId/follow/:fluid", rt.followUser)
	rt.router.DELETE("/user/:userId/follow/:fluid", rt.unfollowUser)
	rt.router.GET("/users/:userId/follows/:fluid", rt.checkFollower)
	rt.router.POST("/user/:userId/ban/:banId", rt.banUser)
	rt.router.DELETE("/user/:userId/ban/:banId", rt.unbanUser)
	rt.router.GET("/users/:userId/bans/:banId", rt.checkBanned)
	rt.router.POST("/user/:userId/likes/:photoId", rt.likePhoto)
	rt.router.DELETE("/user/:userId/likes/:photoId", rt.unlikePhoto)
	rt.router.GET("/users/:userId/likes/:photoId", rt.checkLiked)
	rt.router.POST("/user/:userId/photos/:photoId", rt.commentPhoto)
	rt.router.DELETE("/user/:userId/photos/:photoId/comment/:commentId", rt.uncommentPhoto)
	rt.router.GET("/photos/:photoId/comments", rt.getComments)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	return rt.router
}
