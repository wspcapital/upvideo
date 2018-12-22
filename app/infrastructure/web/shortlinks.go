package web

import (
	"bitbucket.org/marketingx/upvideo/app/storage/shortlinks"
	"bitbucket.org/marketingx/upvideo/validator"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"github.com/teris-io/shortid"
)

type ShortlinksResponse struct {
	Items []*shortlinks.Shortlink `json:"items"`
	Total int                     `json:"total"`
}

func (this *WebServer) shortlinksIndex(c *gin.Context) {
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 32)
	if err != nil {
		limit = 0
	}
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 32)
	if err != nil {
		offset = 0
	}
	uniqId := c.Query("uniq_id")
	if uniqId != "" {
		offset = 0
	}

	items, err := this.ShortlinksService.FindAll(shortlinks.Params{
		UserId:    this.getUser(c).Id,
		UniqId:    uniqId,
		Limit:     uint64(limit),
		Offset:    uint64(offset),
	})
	if err != nil {
		_ = c.AbortWithError(400, err)
		return
	}
	c.JSON(200, ShortlinksResponse{Items: items, Total: len(items)})
}

func (this *WebServer) shortlinkCreate(c *gin.Context) {
	// request
	type ShortlinksCreateRequest struct {
		Url          string `validate:"url"`
	}

	req := &ShortlinksCreateRequest{
		Url:          c.PostForm("url"),
	}

	err := validator.GetValidatorInstance().Struct(req)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("validation error: %s", err.Error()))
		return
	}

	user := this.getUser(c)

	_shortlink := &shortlinks.Shortlink{}
	_shortlink.UserId = user.Id
	_shortlink.UniqId, err = shortid.Generate()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("generate uniq_id error: %s", err.Error()))
		return
	}
	_shortlink.Url = req.Url
	err = this.ShortlinksService.Insert(_shortlink)
	if err != nil {
		fmt.Println("\n this.ShortlinksService.Insert Error: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, _shortlink)
}

func (this *WebServer) shortlinkUpdate(c *gin.Context) {
	var err error
	var _shortlink *shortlinks.Shortlink
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	type RedirectRequest struct {
		//UniqId          string `validate:"shortlink"`
		Url             string `validate:"url"`
	}

	req := &RedirectRequest{
		//UniqId:          c.Param("uniq_id"),
		Url:             c.PostForm("url"),
	}

	err = validator.GetValidatorInstance().Struct(req)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("validation error: %s", err.Error()))
		return
	}

	_shortlink, err = this.ShortlinksService.FindOne(shortlinks.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	//_shortlink.UniqId = req.UniqId
	_shortlink.Url    = req.Url

	err = this.ShortlinksService.Update(_shortlink)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, _shortlink)
}


func (this *WebServer) shortlinkUpdatebyUniqId(c *gin.Context) {
	var err error
	var _shortlink *shortlinks.Shortlink

	type RedirectRequest struct {
		UniqId          string `validate:"shortlink"`
		Url             string `validate:"url"`
	}

	req := &RedirectRequest{
		UniqId:          c.Param("uniqid"),
		Url:             c.PostForm("url"),
	}

	err = validator.GetValidatorInstance().Struct(req)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("validation error: %s", err.Error()))
		return
	}

	_shortlink, err = this.ShortlinksService.FindOne(shortlinks.Params{
		UniqId:     req.UniqId,
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	_shortlink.UniqId = req.UniqId
	_shortlink.Url    = req.Url

	err = this.ShortlinksService.Update(_shortlink)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, _shortlink)
}


func (this *WebServer) shortlinkView(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	_shortlink, err := this.ShortlinksService.FindOne(shortlinks.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	c.JSON(200, _shortlink)
}

func (this *WebServer) shortlinkDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	_shortlink, err := this.ShortlinksService.FindOne(shortlinks.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	this.ShortlinksService.Delete(_shortlink)
	c.Status(200)
}


func (this *WebServer) shortlinkExternal(c *gin.Context) {
	type RedirectRequest struct {
		UniqId          string `validate:"shortlink"`
	}

	req := &RedirectRequest{
		UniqId:          c.Param("uniq_id"),
	}

	err := validator.GetValidatorInstance().Struct(req)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("validation error: %s", err.Error()))
		return
	}

	_shortlink, err := this.ShortlinksService.FindOne(shortlinks.Params{
			UniqId: req.UniqId,
		})

	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, _shortlink.Url)
}