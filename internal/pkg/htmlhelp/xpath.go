package htmlhelp

import (
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xmlquery"
	"golang.org/x/net/html"
)

func GetXpathValue(doc *html.Node, expression string) string {
	node := htmlquery.FindOne(doc, expression)
	if node == nil {
		return ""
	}
	return htmlquery.InnerText(node)
}

func GetXpathXmlValue(doc *xmlquery.Node, expression string) string {
	node := xmlquery.FindOne(doc, expression)
	if node == nil {
		return ""
	}
	return node.InnerText()
}
