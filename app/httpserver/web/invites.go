package web

import (
	"github.com/gin-gonic/gin"
	"bitbucket.org/marketingx/upvideo/app/storage/invites"
	"strconv"
)

type InviteResponse struct {
	Items []*invites.Invite
	Total int
}

func (this *WebServer) inviteIndex(c *gin.Context) {
	limit, err := strconv.ParseInt(c.Param("limit"), 10, 32)
	if err != nil {
		limit = 0
	}
	offset, err := strconv.ParseInt(c.Param("offset"), 10, 32)
	if err != nil {
		offset = 0
	}
	items, err := this.InviteService.FindAll(invites.Params{
		UserId: this.getUser(c).Id,
		Limit:  uint64(limit),
		Offset: uint64(offset),
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	c.JSON(200, InviteResponse{Items: items, Total: len(items)})
}

func (this *WebServer) inviteCreate(c *gin.Context) {
	unsafe := &invites.Invite{}
	c.BindJSON(unsafe)
	invite := &invites.Invite{}
	invite.Title = unsafe.Title
	//invite.Code = unsafe.Code
	invite.UserId = this.getUser(c).Id
	err := this.InviteService.Insert(invite)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, invite)
}

func (this *WebServer) inviteUpdate(c *gin.Context) {
	var err error
	var invite *invites.Invite
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	invite, err = this.InviteService.FindOne(invites.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	unsafe := &invites.Invite{}
	c.BindJSON(unsafe)
	invite.Title = unsafe.Title
	invite.Code = unsafe.Code
	err = this.InviteService.Update(invite)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, invite)
}

func (this *WebServer) inviteView(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	invite, err := this.InviteService.FindOne(invites.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	c.JSON(200, invite)
}

func (this *WebServer) inviteDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	invite, err := this.InviteService.FindOne(invites.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	this.InviteService.Delete(invite)
	c.Status(200)
}
