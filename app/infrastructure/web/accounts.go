package web

import (
	"bitbucket.org/marketingx/upvideo/app/accounts"
	"bitbucket.org/marketingx/upvideo/utils"
	"bitbucket.org/marketingx/upvideo/validator"
	"bitbucket.org/marketingx/upvideo/youtubeauth"
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
)

var (
	uploadSuccessfulRegexp = regexp.MustCompile("Upload successful! Video ID: ([a-zA-Z0-9]+)")
	uploadErrorRegexp      = regexp.MustCompile("Error making YouTube API call:([a-zA-Z0-9,.!_\\-\\s]+)")
)

type AccountResponse struct {
	Items []*accounts.Account
	Total int
}

func (this *WebServer) accountIndex(c *gin.Context) {
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 32)
	if err != nil {
		limit = 0
	}
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 32)
	if err != nil {
		offset = 0
	}
	items, err := this.AccountService.FindAll(accounts.Params{
		UserId: this.getUser(c).Id,
		Limit:  uint64(limit),
		Offset: uint64(offset),
	})
	if err != nil {
		_ = c.AbortWithError(400, err)
		return
	}

	//log.Println(items[0])
	c.JSON(200, AccountResponse{Items: items, Total: len(items)})
}

func (this *WebServer) accountCreate(c *gin.Context) {
	// auth to youtube
	type AccountCreateRequest struct {
		ChannelName string `validate:"required,title"`
		ChannelUrl  string `validate:"required,url"`
		// ClientSecrets string // file
		OTPCode string `validate:"required,alphanumunicode"`
		Note    string `validate:"required,text"`
	}

	req := &AccountCreateRequest{
		ChannelName: c.PostForm("ChannelName"),
		ChannelUrl:  c.PostForm("ChannelURL"),
		OTPCode:     c.PostForm("OTPCode"),
		Note:        c.PostForm("Note"),
	}

	err := validator.GetValidatorInstance().Struct(req)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	file, err := c.FormFile("ClientSecrets")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	src, err := file.Open()
	if err != nil {
		fmt.Println("\n Can not open attached ClientSecrets file. Error: ", err.Error())
		c.String(http.StatusInternalServerError, fmt.Sprintf("get err: %s", err.Error()))
		return
	}
	defer src.Close()

	// check client_id exists for current user
	clientSecretsConf, err := youtubeauth.ParseConfig(src)
	if err != nil {
		fmt.Println("\n Can not parse ClientSecrets. Error: ", err.Error())
		c.String(http.StatusInternalServerError, fmt.Sprintf("get err: %s", err.Error()))
		return
	}

	items, err := this.AccountService.FindAll(accounts.Params{
		UserId:   this.getUser(c).Id,
		ClientId: clientSecretsConf.ClientID,
	})
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("\n AccountService.FindAll Error: ", err.Error())
		c.String(http.StatusInternalServerError, fmt.Sprintf("get err: %s", err.Error()))
		return
	} else if len(items) > 0 {
		c.String(http.StatusBadRequest, fmt.Sprintf("This clent secrets client_id already exists"))
		return
	}

	// store client_secrets file to disk
	var (
		operationId       string
		clientSecretsPath string
	)
	targetDir := path.Join(this.Config.YoutubeUploaderPath, "/secrets")

	// look for unique operationId
	for {
		operationId = utils.RandomString(32)
		clientSecretsPath = path.Join(targetDir, "client_secrets_"+operationId+".json")

		if _, err := os.Stat(clientSecretsPath); os.IsNotExist(err) {
			break
		}
	}

	if err = os.MkdirAll(targetDir, os.ModePerm); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err = c.SaveUploadedFile(file, clientSecretsPath); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	authUrl, err := youtubeauth.GetAuthURL(operationId, clientSecretsPath)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	_account := &accounts.Account{}
	_account.UserId = this.getUser(c).Id
	_account.ChannelName = req.ChannelName
	_account.ChannelUrl = req.ChannelUrl
	_account.ClientId = clientSecretsConf.ClientID
	_account.ClientSecrets = clientSecretsPath
	_account.AuthUrl = authUrl
	_account.OTPCode = req.OTPCode
	_account.Note = req.Note
	_account.OperationId = operationId

	err = this.AccountService.Insert(_account)
	if err != nil {
		fmt.Println("\n AccountService.Insert Error: ", err.Error())
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// response
	type AccountCreateResponse struct {
		OperationId string `json:"operation_id"`
		AuthUrl     string `json:"AuthUrl"`
	}

	res := &AccountCreateResponse{
		OperationId: operationId,
		AuthUrl:     _account.AuthUrl,
	}

	c.JSON(200, res)
}

func (this *WebServer) accountConfirm(c *gin.Context) {
	type AccountConfirmRequest struct {
		OperationId string `json:"operation_id" validate:"required,alphanumunicode"`
		ConfirmCode string `json:"confirm_code" validate:"required,confirm_code"`
	}
	req := &AccountConfirmRequest{
		OperationId: c.PostForm("operation_id"),
		ConfirmCode: c.PostForm("confirm_code"),
	}

	err := validator.GetValidatorInstance().Struct(req)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	_account, err := this.AccountService.FindByOperation(this.getUser(c).Id, req.OperationId)
	if err != nil {
		fmt.Println("\n this.AccountService.FindByOperation Error: ", err.Error())
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _account.RequestToken != "" {
		c.String(http.StatusBadRequest, "Your account already confirmed")
		return
	}

	tokensDir := path.Join(this.Config.YoutubeUploaderPath, "/tokens")
	tokenPath := path.Join(tokensDir, "request_"+req.OperationId+".token.json")
	if err = os.MkdirAll(tokensDir, os.ModePerm); err != nil {
		fmt.Println("\n create tokenPath Error: ", err.Error())
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	secretsDir := path.Join(this.Config.YoutubeUploaderPath, "/secrets")
	clientSecretsPath := path.Join(secretsDir, "client_secrets_"+req.OperationId+".json")

	err = youtubeauth.VerifyCode(req.ConfirmCode, tokenPath, clientSecretsPath)
	if err != nil {
		fmt.Println("\n this.AccountService.FindByOperation Error: ", err.Error())
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// cmd youtubeuploader upload test video
	cmd := exec.Command(this.Config.YoutubeUploaderCmd, "-headlessAuth", "-secrets", clientSecretsPath, "-cache", tokenPath, "-filename", this.Config.TestVideoPath, "-metaJSON", this.Config.TestVideoMetaPath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Youtube uploader Error: " + out.String())
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		c.Status(http.StatusInternalServerError)
		return
	}
	if !uploadSuccessfulRegexp.Match(out.Bytes()) {
		fmt.Println("Youtube uploader Result: " + out.String())
		c.Status(http.StatusInternalServerError)
		return
	}

	matches := uploadSuccessfulRegexp.FindStringSubmatch(out.String())
	url := "https://www.youtube.com/watch?v=" + matches[1]

	_account.RequestToken = tokenPath

	err = this.AccountService.Update(_account)
	if err != nil {
		fmt.Println("\n this.AccountService.Update Error: ", err.Error())
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, gin.H{"uploaded": "succcessfull", "url": url})
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
