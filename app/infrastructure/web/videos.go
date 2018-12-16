package web

import (
	"bitbucket.org/marketingx/upvideo/app/videos"
	"bitbucket.org/marketingx/upvideo/app/videos/titles"
	"bitbucket.org/marketingx/upvideo/aws"
	"bitbucket.org/marketingx/upvideo/validator"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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
	type VideoCreateRequest struct {
		Title       string `validate:"title"`
		Description string `validate:"text"`
		Tags        string `validate:"text"`
		Category    string
		Language    string
		Playlist    string
		IpAddress   string
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

	userId := this.getUser(c).Id
	fileName := filepath.Base(file.Filename)

	req := &videos.Video{
		Title:       c.PostForm("Title"),
		Description: c.PostForm("Description"),
		Tags:        c.PostForm("Tags"),
		Category:    c.PostForm("Category"),
		Language:    c.PostForm("Language"),
		Playlist:    c.PostForm("Playlist"),
		IpAddress:   c.PostForm("IpAddress"),
	}

	err = validator.GetValidatorInstance().Struct(req)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	video := &videos.Video{}
	video.Title = req.Title
	video.Description = req.Description
	video.Tags = req.Tags
	video.Category = req.Category
	video.Language = req.Language
	video.File = fileName
	video.Playlist = req.Playlist
	video.IpAddress = req.IpAddress
	video.UserId = userId
	err = this.VideoService.Insert(video)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	targetPath := fmt.Sprintf("/%d/%d/%s", userId, video.Id, fileName)
	location, err := aws.UploadS3File(targetPath, fileReader)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	video.File = location
	err = this.VideoService.Update(video)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
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

func (this *WebServer) VideoGenerateTitles(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	video, err := this.VideoService.FindOne(videos.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// check that titles is already generated
	if video.TitleGen {
		c.String(http.StatusBadRequest, "Titles already generated.")
		return
	}

	keywords, err := this.KeywordtoolService.GetKeywords(video.Title)
	if err != nil {
		c.String(http.StatusInternalServerError, "Can not get keywords, try again later.")
		return
	}

	//var items []*titles.Title
	titlesCreated := 0
	for _, keyword := range keywords {
		if keyword == "" {
			continue
		}

		tags, err := this.RapidtagsService.GetTags(keyword)
		if err != nil {
			fmt.Printf("Insert title '%s' err: \n%v\n", keyword, err)
			continue
		}

		_title := &titles.Title{
			UserId:    this.getUser(c).Id,
			VideoId:   video.Id,
			Title:     keyword,
			Tags:      strings.Join(tags, ","),
			File:      "",
			Posted:    false,
			Converted: false,
			Pending:   false,
			IpAddress: c.ClientIP(),
		}

		err = this.TitleService.Insert(_title)
		if err != nil {
			fmt.Printf("Insert title '%s' err: \n%v\n", keyword, err)
			continue
		}

		//items = append(items, _title)
		titlesCreated++
	}

	if titlesCreated > 0 {
		video.TitleGen = true
		err = this.VideoService.Update(video)
		if err != nil {
			c.String(http.StatusInternalServerError, "Can not update video.")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"total": titlesCreated})
}

func (this *WebServer) VideoGetTitles(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	limit, err := strconv.ParseInt(c.Query("limit"), 10, 32)
	if err != nil {
		limit = 0
	}
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 32)
	if err != nil {
		offset = 0
	}
	items, err := this.TitleService.FindAll(titles.Params{
		UserId:  this.getUser(c).Id,
		VideoId: int(id),
		Limit:   uint64(limit),
		Offset:  uint64(offset),
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": len(items)})
}
