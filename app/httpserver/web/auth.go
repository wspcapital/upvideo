package web

import (
	"bitbucket.org/marketingx/upvideo/app/domain/session"
	"bitbucket.org/marketingx/upvideo/app/domain/usr"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

func (this *WebServer) getKey(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	if strings.HasPrefix(header, "_") {
		return header[1:]
	}
	return ""
}

func (this *WebServer) getSession(c *gin.Context) *session.Entity {
	header := c.GetHeader("Authorization")
	if strings.HasPrefix(header, "Bearer") {
		entity, err := this.SessionService.FindById(header[7:])
		if err != nil {
			return nil
		}
		this.SessionService.Update(entity) // refresh session when it is used
		return entity
	}
	return nil
}

func (this *WebServer) removeSession(c *gin.Context) (err error) {
	header := c.GetHeader("Authorization")
	if strings.HasPrefix(header, "Bearer") {
		err = this.SessionService.DeleteById(header[7:])
	}
	return err
}

func (this *WebServer) getUser(c *gin.Context) *usr.User {
	sess := this.getSession(c)
	if sess != nil {
		user, err := this.UserService.FindById(sess.Data["userId"])
		if err == nil {
			return user
		}
	}
	apiKey := this.getKey(c)
	if apiKey != "" {
		user, err := this.UserService.FindByKey(apiKey)
		if err == nil {
			return user
		}
	}
	return nil
}

func (this *WebServer) requireAuth(c *gin.Context) {
	user := this.getUser(c)
	if user != nil {
		c.Next()
	} else {
		c.AbortWithError(401, errors.New("Not authorized"))
	}
}
