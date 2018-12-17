package web

import (
	"bitbucket.org/marketingx/upvideo/app/domain/usr"
	"bitbucket.org/marketingx/upvideo/validator"
	_ "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

	type UserRequest struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}

	req := &UserRequest{
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	// validate user
	err := validator.GetValidatorInstance().Struct(req)
	if err != nil {
		c.String(http.StatusBadRequest, "username or password not valid")
		return
	}

	user := &usr.User{
		Email:        req.Email,
		PasswordHash: this.UserService.PasswordHash(req.Password),
	}

	// check user exists
	userExists, err := this.UserService.CheckUserExists(user)
	if err != nil {
		fmt.Println("UserService.CheckUserExists: " + err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	if userExists {
		c.String(http.StatusBadRequest, "Username already exists")
		return
	}

	err = this.UserService.Insert(user)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	fmt.Println(user)
	sess := this.SessionService.Create()
	sess.Data["userId"] = strconv.Itoa(user.Id)
	err = this.SessionService.Update(sess)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
	} else {
		// err = this.InviteService.ClearInvite(inviteCode)
		// if this.Params.InviteOnly && err != nil {
		// 	this.UserService.Delete(user)
		// 	c.AbortWithError(500, err)
		// 	return
		// }
		c.JSON(http.StatusOK, gin.H{"token": sess.Id})
	}

}

func (this *WebServer) profile(c *gin.Context) {
	user := this.getUser(c)
	c.JSON(200, ProfileResponse{Email: user.Email, APIKey: user.APIKey})
}
