package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/owl-93/snapper/model"
	"github.com/owl-93/snapper/service"
	"github.com/owl-93/snapper/utils"
)

var (
	options *model.SnapperConfig
)

func InitRoutes(router *gin.Engine, config *model.SnapperConfig) {
	options = config
	router.POST("/", snapper)
}

/*
Main route that accepts a URL as a query param and attempts to extract the
metadata from the page meta tags
*/
func snapper(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "unable to process request body")
		return
	}
	var request = &model.SnapperRequest{}
	marshalError := json.Unmarshal([]byte(body), request)
	if marshalError != nil {
		log.Println(marshalError.Error())
		c.String(http.StatusInternalServerError, "unable to process request body")
		return
	}

	if len(request.Page) == 0 {
		c.String(http.StatusBadRequest, "must include site query param")
	} else {
		tags, err := service.GetMetaTagsForPage(request.Page, options.DisableCache || request.Refresh, options.CacheTTL)
		handleTagClientResponse(c, tags, err, request.Raw)
	}
}

/*
Helper function to avoid duplicating this code in each Handler
*/
func handleTagClientResponse(c *gin.Context, tags *[]model.MetaTag, err error, raw bool) {
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	} else if tags == nil {
		c.Status(http.StatusNoContent)
	} else if len(*tags) == 0 {
		c.Status(http.StatusNoContent)
	} else {
		if raw {
			c.IndentedJSON(http.StatusOK, tags)
		} else {
			if converted, err := utils.ToSnapperResult(tags); err == nil {
				c.IndentedJSON(http.StatusOK, converted)
			} else {
				c.String(http.StatusInternalServerError, "unable to convert meta tags to snapper response")
			}
		}
	}
}
