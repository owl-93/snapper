package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"snapper/model"
	"snapper/service"
	"snapper/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	router.POST("/", snapper)
	router.GET("/test/:testNum", tester)
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
		tags, err := service.GetMetaTagsForPage(request.Page)
		handleTagClientResponse(c, tags, err)
	}
}


/*
This endpoint tests the service with some given HTML
*/
func tester(c *gin.Context) {
	testHtml := getTestHtml(c)
	tags, err := service.GetMetaTagsForTest(testHtml)
	handleTagClientResponse(c, tags, err)
}


/*
Extract the path and load appropriate test HTML from utils
 */
func getTestHtml(c *gin.Context) string {
	testNum, _ := c.Params.Get("testNum")
	idx, err := strconv.Atoi(testNum)
	if err != nil {
		println("could not convert", testNum, "to an int")
		return utils.Tests[0]
	}
	return utils.Tests[idx]
}


/*
Helper function to avoid duplicating this code in each Handler
*/
func handleTagClientResponse(c *gin.Context, tags *[]model.MetaTag, err error) {
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	} else if tags == nil {
		c.Status(http.StatusNoContent)
	} else {
		c.IndentedJSON(http.StatusOK, tags)
	}
}
