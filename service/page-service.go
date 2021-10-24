package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"snapper/cache"
	"snapper/model"
	"snapper/utils"
)

/*
	handles fetching and webpages HTML and parsing and extracting metadata tags from it.
	Interfaces with a caching layer as well
*/
func GetMetaTagsForPage(address string) (*[]model.MetaTag, error) {
	log.Printf("get meta data for %s\n", address)
	//check the cache if configured
	if cache.IsInitialized() {
		tags, e := cache.CheckCacheForPage(address)
		if e == nil && tags != nil {
			return tags, nil
		}
	}
	//Cache missed or errored, fetch and add page to cache
	response, err := http.Get(address)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to fetch page %s", address))
	}
	parsed, err := getHTML(response)
	if err != nil {
		return nil, errors.New("unable to read html")
	}
	tags, err := utils.ExtractMetaTags(parsed)
	if err == nil && cache.IsInitialized() {
		cacheError := cache.CachePageMetaData(tags, address)
		if cacheError != nil {
			return tags, cacheError
		}
	}
	return tags, err
}

/*
	Extract meta tags for test html
*/
func GetMetaTagsForTest(htmlContent string) (*[]model.MetaTag, error) {
	return utils.ExtractMetaTags(htmlContent)
}

/*
	fetches a requested page and returns the result as a string of HTML
*/
func getHTML(response *http.Response) (string, error) {
	if response == nil {
		return "", errors.New("no response")
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("can't parse html from body")
	}
	return string(bytes), nil
}
