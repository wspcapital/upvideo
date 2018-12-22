package web

import (
	"bitbucket.org/marketingx/upvideo/app/storage/videos"
	"bitbucket.org/marketingx/upvideo/app/utils/aws"
	"bitbucket.org/marketingx/upvideo/app/utils/validator"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	_ "log"
	"net/http"
	"path/filepath"
	"strconv"
)

type VideoResponse struct {
	Items []*videos.Video `json:"items"`
	Total int             `json:"total"`
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
	accountId, err := strconv.ParseInt(c.Query("account"), 10, 32)
	if err != nil {
		offset = 0
	}
	items, err := this.VideoService.FindAll(videos.Params{
		UserId:    this.getUser(c).Id,
		AccountId: int(accountId),
		Limit:     uint64(limit),
		Offset:    uint64(offset),
	})
	if err != nil {
		_ = c.AbortWithError(400, err)
		return
	}
	c.JSON(200, VideoResponse{Items: items, Total: len(items)})
}

func (this *WebServer) videoCreate(c *gin.Context) {
	// request
	type VideoCreateRequest struct {
		Title       string `validate:"required,title"`
		Description string `validate:"text"`
		Tags        string `validate:"required,text"`
		Category    string
		Language    string
		Playlist    string
	}

	file, err := c.FormFile("File")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	fileReader, err := file.Open()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	defer fileReader.Close()

	user := this.getUser(c)
	fileName := filepath.Base(file.Filename)

	req := &VideoCreateRequest{
		Title:       c.PostForm("Title"),
		Description: c.PostForm("Description"),
		Tags:        c.PostForm("Tags"),
		Category:    c.PostForm("Category"),
		Language:    c.PostForm("Language"),
		Playlist:    c.PostForm("Playlist"),
	}

	err = validator.GetValidatorInstance().Struct(req)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	videoUuid, err := uuid.NewV4()
	if err != nil {
		fmt.Println("\n UUID generation Error: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	targetPath := fmt.Sprintf("/%d/%s/%s", user.Id, videoUuid.String(), fileName)
	location, err := aws.UploadS3File(targetPath, fileReader)
	if err != nil {
		fmt.Println("\n aws.UploadS3File Error: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	video := &videos.Video{}
	video.UserId = user.Id
	video.AccountId = user.AccountId
	video.Title = req.Title
	video.Description = req.Description
	video.Tags = req.Tags
	video.Category = req.Category
	video.Language = req.Language
	video.File = location
	video.Playlist = req.Playlist
	video.IpAddress = c.ClientIP()
	err = this.VideoService.Insert(video)
	if err != nil {
		fmt.Println("\n this.VideoService.Insert Error: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, video)
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
	video.IpAddress = c.ClientIP()

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
