package web

import (
	"github.com/gin-gonic/gin"
	"bitbucket.org/marketingx/upvideo/app/accounts"
	"strconv"
	_"log"
)

type AccountResponse struct {
	Items []*accounts.Account
	Total int
}

func (this *WebServer) accountIndex(c *gin.Context) {
	limit, err := strconv.ParseInt(c.Param("limit"), 10, 32)
	if err != nil {
		limit = 0
	}
	offset, err := strconv.ParseInt(c.Param("offset"), 10, 32)
	if err != nil {
		offset = 0
	}
	items, err := this.AccountService.FindAll(accounts.Params{
		UserId: this.getUser(c).Id,
		Limit:  uint64(limit),
		Offset: uint64(offset),
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	//log.Println(items[0])
	c.JSON(200, AccountResponse{Items: items, Total: len(items)})
}

func (this *WebServer) accountCreate(c *gin.Context) {
	unsafe := &accounts.Account{}
	c.BindJSON(unsafe)
	_account := &accounts.Account{}
	_account.UserId = this.getUser(c).Id
	_account.Username = unsafe.Username
	_account.Password = unsafe.Password
	_account.ChannelName = unsafe.ChannelName
	_account.ChannelUrl = unsafe.ChannelUrl
	_account.ClientSecrets = unsafe.ClientSecrets
	_account.RequestToken = unsafe.RequestToken
	_account.AuthUrl = unsafe.AuthUrl
	_account.OTPCode = unsafe.OTPCode
	_account.Note = unsafe.Note
	err := this.AccountService.Insert(_account)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, _account)
}

func (this *WebServer) accountUpdate(c *gin.Context) {
	var err error
	var _account *accounts.Account
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	_account, err = this.AccountService.FindOne(accounts.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	unsafe := &accounts.Account{}
	c.BindJSON(unsafe)
	
	_account.Username = unsafe.Username
	_account.Password = unsafe.Password
	_account.ChannelName = unsafe.ChannelName
	_account.ChannelUrl = unsafe.ChannelUrl
	_account.ClientSecrets = unsafe.ClientSecrets
	_account.RequestToken = unsafe.RequestToken
	_account.AuthUrl = unsafe.AuthUrl
	_account.OTPCode = unsafe.OTPCode
	_account.Note = unsafe.Note

	err = this.AccountService.Update(_account)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, _account)
}

func (this *WebServer) accountView(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	_account, err := this.AccountService.FindOne(accounts.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	c.JSON(200, _account)
}

func (this *WebServer) accountDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	_account, err := this.AccountService.FindOne(accounts.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}


	this.AccountService.Delete(_account)
	c.Status(200)
}
