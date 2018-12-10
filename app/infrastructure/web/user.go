package web

import (
	_"errors"
	"github.com/gin-gonic/gin"
	"bitbucket.org/marketingx/upvideo/app/domain/usr"
	"strconv"
	"fmt"
)

type ProfileResponse struct {
	Email  string
	APIKey string
}

func (this *WebServer) todoLogin(c *gin.Context) {
	user, err := this.UserService.Login(c.PostForm("login"), c.PostForm("password"))
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	sess := this.SessionService.Create()
	sess.Data["userId"] = strconv.Itoa(user.Id)
	err = this.SessionService.Update(sess)
	if err != nil {
		c.JSON(500, gin.H{"error": "User not found"})
	} else {
		c.JSON(200, gin.H{"token": sess.Id})
	}
}

func (this *WebServer) register(c *gin.Context) {

	// if this.Params.Registration == false {
	// 	c.AbortWithError(403, errors.New("Registration Unavailable"))
	// 	return
	// }
	// inviteCode := c.PostForm("code")
	// err := this.InviteService.CheckInvite(inviteCode)
	// if this.Params.InviteOnly && err != nil {
	// 	c.AbortWithError(403, err)
	// 	return
	// }

	user := &usr.User{}
	user.Email = c.PostForm("email")
	user.PasswordHash = this.UserService.PasswordHash(c.PostForm("password"))
	fmt.Println(user)
	err := this.UserService.Insert(user)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	fmt.Println(user)
	sess := this.SessionService.Create()
	sess.Data["userId"] = strconv.Itoa(user.Id)
	err = this.SessionService.Update(sess)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "Could not create user"})
	} else {
		// err = this.InviteService.ClearInvite(inviteCode)
		// if this.Params.InviteOnly && err != nil {
		// 	this.UserService.Delete(user)
		// 	c.AbortWithError(500, err)
		// 	return
		// }
		c.JSON(200, gin.H{"token": sess.Id})
	}

}

func (this *WebServer) profile(c *gin.Context) {
	user := this.getUser(c)
	c.JSON(200, ProfileResponse{Email: user.Email, APIKey: user.APIKey})
}
