package utils

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"snapper/model"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)


const (
	MetaKey = "property"
	MetaValue = "content"
)


/*
	Extracts the meta tags from a given html string
*/
func ExtractMetaTags(htmlString string) (*[]model.MetaTag, error) {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return nil, err
	}
	head, err := getDocHead(doc)
	if err != nil {
		return nil, errors.New("error getting document head")
	}
	tagsPtr, err := getMetaTags(head)
	if err != nil {
		return nil, errors.New("error getting meta tags")
	}
	return tagsPtr, nil
}

/*
	Locate the documents head node
*/
func getDocHead(doc *html.Node) (*html.Node, error) {
	if doc.Type == html.DocumentNode {
		curNode := doc.FirstChild
		for ; curNode != nil && curNode.DataAtom != atom.Html; curNode = curNode.NextSibling {}
		if curNode != nil {
			curNode = curNode.FirstChild
			if curNode != nil && curNode.DataAtom == atom.Head {
				return curNode, nil
			}else {
				return nil, errors.New("unable to find head element in document")
			}
		}
		return nil, errors.New("unable to find html element in document")

	}
	return nil, errors.New("passed node is not an document root element")
}

/*
	Extract the meta tags from a page's head
*/
func getMetaTags(headNode *html.Node) (*[]model.MetaTag, error) {
	var tags []model.MetaTag
	curNode := headNode.FirstChild
	for ; curNode != nil; curNode = curNode.NextSibling {
		if curNode.Type == html.ElementNode && curNode.DataAtom == atom.Meta {
			var name, val string
			for _, attr := range curNode.Attr {
				switch attr.Key {
				case MetaKey:
					name = attr.Val
				case MetaValue:
					val = attr.Val
				}
			}
			if len(name) > 0 && len(val) > 0 {
				tags = append(tags, model.MetaTag{Name: name, Value: val})
			}
		}
	}
	log.Printf("found %d meta tags", len(tags))
	return &tags, nil
}


/*
	Extract key from page URL. (strip protocol and query params)
*/
func GetAddressKey(address string) (string, error) {
	parsed, e := url.Parse(address)
	if e != nil {
		return "", e
	}
	return fmt.Sprintf("%s%s", parsed.Host, parsed.Path), nil
}


/*
	converts a list of meta tags to a utility structure for JSON response
 */
func ToSnapperResult(tags *[]model.MetaTag) (*model.SnapperResult, error) {
	snapperResult := model.SnapperResult{}
	for _, tag := range *tags {
		switch tag.Name {
		case model.TitleKey:
			snapperResult.Title = tag.Value
		case model.ImageKey:
			snapperResult.Image = tag.Value
		case model.DescKey:
			snapperResult.Description = tag.Value
		case model.UrlKey:
			snapperResult.Url = tag.Value
		case model.TypeKey:
			snapperResult.Type = tag.Value
		case model.LocalKey:
			snapperResult.Locale = tag.Value
		}
	}
	return &snapperResult, nil
}