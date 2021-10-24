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


//Test HTML for development
var Tests = [3]string{
				`<html>
					<head>
						<meta property="testtag0" content="testvalue0"/>
					</head>

					<body>
					</body>
				</html>`,

				`<!DOCTYPE html/>
				<html>
					<head>
						<meta property="testtag1" content="testvalue1"/>
					</head>
					<body>
					</body>
				</html>`,

				`<html>
					<head>
						<meta property="testtag2" content="testvalue2"/>
					</head>
					<body>
						<p>p0 text</p>
						<p>p1 text</p>
						<div>
							<p>p2 text</p>
							<p>p3 text</p>
						</div>
					</body>
				</html>`,
}