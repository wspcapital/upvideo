package web

import (
	"bitbucket.org/marketingx/upvideo/app/campaigns"
	"bitbucket.org/marketingx/upvideo/app/titles"
	"bitbucket.org/marketingx/upvideo/app/videos"
	"bitbucket.org/marketingx/upvideo/validator"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "log"
	"net/http"
	"strconv"
	"strings"
)

type CampaignResponse struct {
	Items []*campaigns.Campaign `json:"items"`
	Total int                   `json:"total"`
}

func (this *WebServer) campaignIndex(c *gin.Context) {
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 32)
	if err != nil {
		limit = 0
	}
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 32)
	if err != nil {
		offset = 0
	}
	items, err := this.CampaignService.FindAll(campaigns.Params{
		UserId: this.getUser(c).Id,
		Limit:  uint64(limit),
		Offset: uint64(offset),
	})
	if err != nil {
		fmt.Println("\n this.CampaignService.FindAll Error: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, CampaignResponse{Items: items, Total: len(items)})
}

func (this *WebServer) campaignCreate(c *gin.Context) {
	// request
	type request struct {
		AccountId int    `validate:"required,gt=0"`
		VideoId   int    `validate:"required,gt=0"`
		Title     string `validate:"required,title"`
	}

	accountId, err := strconv.ParseInt(c.PostForm("account_id"), 10, 32)
	if err != nil {
		accountId = 0
	}
	videoId, err := strconv.ParseInt(c.PostForm("video_id"), 10, 32)
	if err != nil {
		accountId = 0
	}

	user := this.getUser(c)

	req := &request{
		AccountId: int(accountId),
		VideoId:   int(videoId),
		Title:     c.PostForm("title"),
	}
	err = validator.GetValidatorInstance().Struct(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.JsonErrors(err)})
		return
	}

	campaign := &campaigns.Campaign{}
	campaign.UserId = user.Id
	campaign.AccountId = req.AccountId // user.AccountId
	campaign.VideoId = req.VideoId
	campaign.Title = req.Title
	campaign.IpAddress = c.ClientIP()
	err = this.CampaignService.Insert(campaign)
	if err != nil {
		fmt.Println("\n this.CampaignService.Insert Error: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(200, campaign)
}

func (this *WebServer) campaignUpdate(c *gin.Context) {
	type request struct {
		Title string `validate:"required,title"`
	}
	var err error
	var campaign *campaigns.Campaign
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	campaign, err = this.CampaignService.FindOne(campaigns.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	req := &request{}
	err = c.BindJSON(req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	campaign.Title = req.Title

	err = this.CampaignService.Update(campaign)
	if err != nil {
		fmt.Println("\n this.CampaignService.Update Error: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, campaign)
}

func (this *WebServer) campaignView(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	campaign, err := this.CampaignService.FindOne(campaigns.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		fmt.Println("\n this.CampaignService.FindOne Error: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, campaign)
}

func (this *WebServer) campaignDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	campaign, err := this.CampaignService.FindOne(campaigns.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err = this.CampaignService.Delete(campaign)
	if err != nil {
		fmt.Println("\n this.CampaignService.Delete Error: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func (this *WebServer) campaignGenerateTitles(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		fmt.Println("\n campaignGenerateTitles Error: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	campaign, err := this.CampaignService.FindOne(campaigns.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		fmt.Println("\n this.CampaignService.FindOne Error: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	video, err := this.VideoService.FindOne(videos.Params{Id: campaign.VideoId})
	if err != nil {
		fmt.Println("\n this.VideoService.FindOne Error: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	// check that titles is already generated
	if campaign.TitlesGenerated {
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
	for _, title := range keywords {
		if title == "" {
			continue
		}

		_title := &titles.Title{
			UserId:     this.getUser(c).Id,
			CampaignId: campaign.Id,
			Title:      title,
			IpAddress:  c.ClientIP(),
		}

		has, err := this.TitleService.Has(_title)
		if err != nil {
			fmt.Printf("Has title '%s' err: \n%v\n", title, err)
			continue
		}
		if has {
			fmt.Printf("Title already exists '%s', campaignId: %d \n", title, _title.CampaignId)
			continue
		}

		tags, err := this.RapidtagsService.GetTags(title)
		if err != nil {
			fmt.Printf("Insert title '%s' err: \n%v\n", title, err)
			continue
		}

		_title.Tags = strings.Join(tags, ",")

		err = this.TitleService.Insert(_title)
		if err != nil {
			fmt.Printf("Insert title '%s' err: \n%v\n", title, err)
			continue
		}

		//items = append(items, _title)
		titlesCreated++
	}

	if titlesCreated > 0 {
		campaign.TitlesGenerated = true
		err = this.CampaignService.CountTotalTitles(campaign)
		if err != nil {
			fmt.Println("\n this.CampaignService.CountTotalTitles Error: ", err.Error())
			c.String(http.StatusInternalServerError, "Can not update campaign.")
			return
		}
		err = this.CampaignService.Update(campaign)
		if err != nil {
			fmt.Println("\n this.CampaignService.Update Error: ", err.Error())
			c.String(http.StatusInternalServerError, "Can not update campaign.")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"count": titlesCreated})
}

func (this *WebServer) campaignGetTitles(c *gin.Context) {
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

	campaign, err := this.CampaignService.FindOne(campaigns.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		fmt.Println("\n this.CampaignService.FindOne Error: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	user := this.getUser(c)
	items, err := this.TitleService.FindAll(titles.Params{
		UserId:     user.Id,
		CampaignId: campaign.Id,
		Limit:      uint64(limit),
		Offset:     uint64(offset),
	})
	if err != nil {
		fmt.Println("\n this.CampaignService.FindAll Error: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": len(items)})
}

func (this *WebServer) campaignAddTitles(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	campaign, err := this.CampaignService.FindOne(campaigns.Params{
		Id:     int(id),
		UserId: this.getUser(c).Id,
	})
	if err != nil {
		fmt.Println("\n this.CampaignService.FindOne Error: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	type request struct {
		Titles string `validate:"required"`
	}

	req := &request{Titles: c.PostForm("titles")}
	err = validator.GetValidatorInstance().Struct(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.JsonErrors(err)})
		return
	}

	titlesCreated := 0
	_titles := strings.Split(req.Titles, "\n")
	for _, title := range _titles {
		_title := &titles.Title{
			UserId:     this.getUser(c).Id,
			CampaignId: campaign.Id,
			Title:      title,
			IpAddress:  c.ClientIP(),
		}

		has, err := this.TitleService.Has(_title)
		if err != nil {
			fmt.Printf("Has title '%s' err: \n%v\n", title, err)
			continue
		}
		if has {
			fmt.Printf("Title already exists '%s', campaignId: %d \n", title, _title.CampaignId)
			continue
		}

		tags, err := this.RapidtagsService.GetTags(title)
		if err != nil {
			fmt.Printf("Insert title '%s' err: \n%v\n", title, err)
			continue
		}

		_title.Tags = strings.Join(tags, ",")

		err = this.TitleService.Insert(_title)
		if err != nil {
			fmt.Printf("Insert title '%s' err: \n%v\n", title, err)
			continue
		}

		titlesCreated++
	}

	if titlesCreated > 0 {
		campaign.TitlesGenerated = true
		err = this.CampaignService.CountTotalTitles(campaign)
		if err != nil {
			fmt.Println("\n this.CampaignService.CountTotalTitles Error: ", err.Error())
			c.String(http.StatusInternalServerError, "Can not update campaign.")
			return
		}
		err = this.CampaignService.Update(campaign)
		if err != nil {
			fmt.Println("\n this.CampaignService.Update Error: ", err.Error())
			c.String(http.StatusInternalServerError, "Can not update campaign.")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"count": titlesCreated})
}
