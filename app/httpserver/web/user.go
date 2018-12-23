package web

import (
	"bitbucket.org/marketingx/upvideo/app/domain/usr"
	"bitbucket.org/marketingx/upvideo/app/utils/validator"
	_"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProfileResponse struct {
	Email  string `json:"email"`
	APIKey string `json:"api_key"`
}

func (this *WebServer) signin(c *gin.Context) {

	if len(c.PostForm("username")) <= 3 || len(c.PostForm("password")) <= 3  {
		c.String(403, "Username or Password is Invalid.")
		return
	}

	user, err := this.UserService.Login(c.PostForm("username"), c.PostForm("password"))
	if err != nil {
		c.String(403, "Username or Password is Invalid.")
		return
	}
	sess := this.SessionService.Create()
	sess.Data["userId"] = strconv.Itoa(user.Id)
	err = this.SessionService.Update(sess)
	if err != nil {
		c.JSON(500, gin.H{"error": "User not found"})
	} else {
		c.JSON(200, gin.H{"token": sess.Id, "access": gin.H{"username": user.Email}})
	}
}

func (this *WebServer) signup(c *gin.Context) {

	if this.Params.Registration == false {
		c.String(403, "Registration for new members has ben disabled by admin.")
		return
	}
	
	inviteCode := c.PostForm("code")
	err := this.InviteService.CheckInvite(inviteCode)
	if this.Params.InviteOnly && err != nil {
		c.String(403, "Invite code requested")
		return
	}

	type UserRequest struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}

	req := &UserRequest{
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	// validate request
	err  = validator.GetValidatorInstance().Struct(req)
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

	//fmt.Println(user)
	sess := this.SessionService.Create()
	sess.Data["userId"] = strconv.Itoa(user.Id)
	err = this.SessionService.Update(sess)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
	} else {
		err = this.InviteService.ClearInvite(inviteCode)
		if this.Params.InviteOnly && err != nil {
			this.UserService.Delete(user)
			c.AbortWithError(500, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": sess.Id})
	}
}

func (this *WebServer) userForgotPassword(c *gin.Context) {
	type UserForgotPasswordRequest struct {
		Email string `validate:"required,email"`
	}

	req := &UserForgotPasswordRequest{
		Email: c.PostForm("email"),
	}

	// validate request
	err := validator.GetValidatorInstance().Struct(req)
	if err != nil {
		c.String(http.StatusBadRequest, "email is not valid")
		return
	}

	user, err := this.UserService.FindByEmail(req.Email)
	if err != nil {
		//fmt.Println("UserService.FindByEmail: " + err.Error())
		c.String(http.StatusBadRequest, "user not found")
		return
	}

	err = this.UserService.SetForgotPasswordToken(user)
	if err != nil {
		//fmt.Println("UserService.SetForgotPasswordToken: " + err.Error())
		c.String(http.StatusInternalServerError, "Try again later")
		return
	}

	err = this.EmailService.SendRestorePasswordEmail(user)
	if err != nil {
		//fmt.Println("UserService.SetForgotPasswordToken: " + err.Error())
		c.String(http.StatusInternalServerError, "Try again later")
		return
	}

	c.Status(http.StatusOK)
}

func (this *WebServer) userRestorePassword(c *gin.Context) {
	// process post request with new password
	// restore form is in frontend
	type UserRestorePasswordRequest struct {
		Email    string `validate:"required,email"`
		Token    string `validate:"required,alphanum,len=36"`
		Password string `validate:"required"`
	}

	req := &UserRestorePasswordRequest{
		Email:    c.PostForm("email"),
		Token:    c.PostForm("token"),
		Password: c.PostForm("password"),
	}

	// validate request
	err := validator.GetValidatorInstance().Struct(req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = this.UserService.SetNewForgottenPassword(&usr.User{Email: req.Email, ForgotPasswordToken: req.Token, PasswordHash: this.UserService.PasswordHash(req.Password)})
	if err != nil {
		fmt.Println("UserService.SetNewForgottenPassword: " + err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}

func (this *WebServer) userChangePassword(c *gin.Context) {
	type UserChangePasswordRequest struct {
		Password string `validate:"required"`
	}

	req := &UserChangePasswordRequest{
		Password: c.PostForm("password"),
	}

	// validate request
	err := validator.GetValidatorInstance().Struct(req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	user := this.getUser(c)
	user.PasswordHash = this.UserService.PasswordHash(req.Password)
	err = this.UserService.Update(user)
	if err != nil {
		fmt.Println("UserService.Update: " + err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (this *WebServer) userResetApikey(c *gin.Context) {
	user := this.getUser(c)
	err := this.UserService.ResetApiKey(user)
	if err != nil {
		fmt.Println("UserService.Update: " + err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	err = this.removeSession(c)
	if err != nil {
		fmt.Println("removeSession: " + err.Error())
	}

	sess := this.SessionService.Create()
	sess.Data["userId"] = strconv.Itoa(user.Id)
	err = this.SessionService.Update(sess)
	if err != nil {
		fmt.Println("SessionService.Update: " + err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"apikey": user.APIKey, "token": sess.Id})
}

func (this *WebServer) profile(c *gin.Context) {
	user := this.getUser(c)
	c.JSON(200, ProfileResponse{Email: user.Email, APIKey: user.APIKey})
}
