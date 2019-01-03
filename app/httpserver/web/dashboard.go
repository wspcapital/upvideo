package web

import (
	"bitbucket.org/marketingx/upvideo/app/domain/usr"
	"bitbucket.org/marketingx/upvideo/app/storage/topusers"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type DashboardResponse struct {
	Items []*usr.User `json:"items"`
	Total int         `json:"total"`
}

func (this *WebServer) dashboardView(c *gin.Context) {

	queryParams := *new(topusers.Params)
	timeFormat := "2006-01-02 15:04:05"
	switch c.Param("period") {
	case "today":
		now := time.Now()
		queryParams.LessThan = now.Format(timeFormat)

		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		queryParams.GreaterThan = startOfDay.Format(timeFormat)
		break
	case "week":
		now := time.Now()
		queryParams.LessThan = now.Format(timeFormat)

		weekAgo := now.AddDate(0, 0, -7)
		queryParams.GreaterThan = weekAgo.Format(timeFormat)
		break
	case "month":
		now := time.Now()
		queryParams.LessThan = now.Format(timeFormat)

		weekAgo := now.AddDate(0, -1, 0)
		queryParams.GreaterThan = weekAgo.Format(timeFormat)
		break
	default:
		c.String(http.StatusBadRequest, fmt.Sprintf("validation error: %s", "Parameter period not valid. Please use: today,week,month."))
		return
	}

	fmt.Println(queryParams)

	items, err := this.TopUserService.FindAll(queryParams)
	if err != nil {
		fmt.Println("\n this.CampaignService.FindAll Error: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": len(items)})

}
