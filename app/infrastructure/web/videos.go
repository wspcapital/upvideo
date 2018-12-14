package web

import (
	"bitbucket.org/marketingx/upvideo/app/videos"
	"github.com/gin-gonic/gin"
	_ "log"
	"strconv"
)

type VideoResponse struct {
	Items []*videos.Video
	Total int
}

func (this *WebServer) videoIndex(c *gin.Context) {
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 32)
	if err != nil {
		limit = 0
	}
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 32)
	if err != nil {
		offset = 0
	}
	items, err := this.VideoService.FindAll(videos.Params{
		UserId: this.getUser(c).Id,
		Limit:  uint64(limit),
		Offset: uint64(offset),
	})
	if err != nil {
		_ = c.AbortWithError(400, err)
		return
	}
	c.JSON(200, VideoResponse{Items: items, Total: len(items)})
}

func (this *WebServer) videoCreate(c *gin.Context) {
	// request
	//{
	//	"Id": 7,
	//	"UserId": 6,
	//	"Title": "dsadasdsa",
	//	"Description": "dsadasda",
	//	"Tags": "dsadasdasda",
	//	"Category": "TT",
	//	"Language": "AA",
	//	"File": "dsadsadsa",
	//	"Playlist": "xadasdsa",
	//	"IpAddress": "dasdsadsa"
	//}

	unsafe := &videos.Video{}
	c.BindJSON(unsafe)
	video := &videos.Video{}
	video.Title = unsafe.Title
	video.Description = unsafe.Description
	video.Tags = unsafe.Tags
	video.Category = unsafe.Category
	video.Language = unsafe.Language
	video.File = unsafe.File
	video.Playlist = unsafe.Playlist
	video.IpAddress = unsafe.IpAddress
	video.UserId = this.getUser(c).Id
	err := this.VideoService.Insert(video)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, video)
}

func (this *WebServer) videoUpdate(c *gin.Context) {
	var err error
	var video *videos.Video
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	video, err = this.VideoService.FindOne(videos.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	unsafe := &videos.Video{}
	c.BindJSON(unsafe)

	video.Title = unsafe.Title
	video.Description = unsafe.Description
	video.Tags = unsafe.Tags
	video.Category = unsafe.Category
	video.Language = unsafe.Language
	video.File = unsafe.File
	video.Playlist = unsafe.Playlist
	video.IpAddress = unsafe.IpAddress

	err = this.VideoService.Update(video)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, video)
}

func (this *WebServer) videoView(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	video, err := this.VideoService.FindOne(videos.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	c.JSON(200, video)
}

func (this *WebServer) videoDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	video, err := this.VideoService.FindOne(videos.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	this.VideoService.Delete(video)
	c.Status(200)
}
