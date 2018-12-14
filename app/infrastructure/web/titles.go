package web

import (
	"bitbucket.org/marketingx/upvideo/app/videos/titles"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type TitleResponse struct {
	Items []*titles.Title
	Total int
}

func (this *WebServer) titleIndex(c *gin.Context) {
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 32)
	if err != nil {
		limit = 0
	}
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 32)
	if err != nil {
		offset = 0
	}
	items, err := this.TitleService.FindAll(titles.Params{
		UserId: this.getUser(c).Id,
		//TODO: VideoID
		Limit:  uint64(limit),
		Offset: uint64(offset),
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	c.JSON(200, TitleResponse{Items: items, Total: len(items)})
}

func (this *WebServer) titleCreate(c *gin.Context) {
	unsafe := &titles.Title{}
	log.Println(unsafe)
	c.BindJSON(unsafe)
	_title := &titles.Title{}
	_title.UserId = this.getUser(c).Id
	//TODO: _title.VideoID
	_title.Title = unsafe.Title
	_title.Tags = unsafe.Tags
	_title.File = unsafe.File
	_title.Posted = unsafe.Posted
	_title.Converted = unsafe.Converted
	_title.Pending = unsafe.Pending
	_title.IpAddress = unsafe.IpAddress
	log.Println(_title)

	err := this.TitleService.Insert(_title)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, _title)
}

func (this *WebServer) titleUpdate(c *gin.Context) {
	var err error
	var _title *titles.Title
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	_title, err = this.TitleService.FindOne(titles.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	unsafe := &titles.Title{}
	c.BindJSON(unsafe)

	//TODO: _title.VideoID
	_title.Title = unsafe.Title
	_title.Tags = unsafe.Tags
	_title.File = unsafe.File
	_title.Posted = unsafe.Posted
	_title.Converted = unsafe.Converted
	_title.Pending = unsafe.Pending
	_title.IpAddress = unsafe.IpAddress

	err = this.TitleService.Update(_title)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, _title)
}

func (this *WebServer) titleView(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	_title, err := this.TitleService.FindOne(titles.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	c.JSON(200, _title)
}

func (this *WebServer) titleDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	_title, err := this.TitleService.FindOne(titles.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	this.TitleService.Delete(_title)
	c.Status(200)
}
