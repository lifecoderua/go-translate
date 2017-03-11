package google

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"bytes"

	"golang.org/x/net/html"
)

// Handler for the DOM node, located by element and attr
type Handler struct {
	element string
	attr    string
	attrVal string
	handler func(n *html.Node, a *html.Attribute)
}

func traverse(n *html.Node, handlers *[]Handler) {
	if n.Type == html.ElementNode {
		for _, handler := range *handlers {
			if handler.element == "" || n.Data == handler.element {
				skipHandler := false
				attr := &html.Attribute{}
				if handler.attr != "" {
					for i, a := range n.Attr {
						if a.Key == handler.attr {
							if handler.attrVal == "" || a.Val == handler.attrVal {
								attr = &n.Attr[i]
							} else {
								skipHandler = true
							}

							break
						}
					}
				}

				if !skipHandler {
					// TODO: pass attrIndex || -1 instead
					handler.handler(n, attr)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(c, handlers)
	}
}

// Cleanup a fetched document from the engine-specific markup
func Cleanup(source string) string {
	// TODO: provide io.Reader instead
	doc, _ := html.Parse(strings.NewReader(source))

	// TODO: make configurable
	sourceDomain := "intercom.com"

	handlers := []Handler{
		// GENERIC
		// global replacemets
		// ? tag fixup

		// default fixup
		// base tag
		Handler{element: "base", handler: func(n *html.Node, a *html.Attribute) { n.Parent.RemoveChild(n) }},

		// url base replacement
		Handler{element: "a", attr: "href", handler: func(n *html.Node, a *html.Attribute) {
			a.Val = regexp.MustCompile("https?://([^./]*.?)"+sourceDomain+"/?").ReplaceAllLiteralString(a.Val, "/")
		}},

		// ENGINE
		// iframe
		// onload
		// css
		// scripts
		// span wrappers
		// tags
		Handler{element: "a", attr: "href", handler: func(n *html.Node, a *html.Attribute) { fmt.Println("BOO>> :", n.Data, a.Val) }},
		Handler{element: "a", attr: "href", attrVal: "/sign_in", handler: func(n *html.Node, a *html.Attribute) { fmt.Println(">> YAY << :", n.Data, a.Val) }},
		Handler{element: "style", handler: func(n *html.Node, a *html.Attribute) { n.Parent.RemoveChild(n) }},
		Handler{element: "style", handler: func(n *html.Node, a *html.Attribute) { fmt.Println(">> STYLE << :", n.Data) }},
	}

	traverse(doc, &handlers)
	// iframe(&doc)
	// onload(&doc)
	var resBytes []byte
	res := bytes.NewBuffer(resBytes)
	resWriter := bufio.NewWriter(res)
	html.Render(resWriter, doc)
	resWriter.Flush()

	return res.String()
}

func iframe(doc *string) {
	// base is "(?si:<iframe[^>]*?translate\.google\.com[^>]*?</iframe>})"
	re := regexp.MustCompile(`(?si:<iframe[^>]*?translate\.google\.com[^>].*?</iframe>})`)
	*doc = re.ReplaceAllLiteralString(*doc, "")
	// fmt.Println(*doc == xdoc)
}

func onload(doc *string) {
	// panic("NOT IMPLEMENTED")
}

func css(doc string) string {
	panic("NOT IMPLEMENTED")
	return doc
}

func scripts(doc string) string {
	panic("NOT IMPLEMENTED")
	return doc
}

func spanWrappers(doc string) string {
	panic("NOT IMPLEMENTED")
	return doc
}

func tags(doc string) string {
	panic("NOT IMPLEMENTED")
	return doc
}
