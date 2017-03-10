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
				attr := html.Attribute{}
				if handler.attr != "" {
					for _, a := range n.Attr {
						if a.Key == handler.attr {
							if handler.attrVal == "" || a.Val == handler.attrVal {
								attr = a
							} else {
								skipHandler = true
							}

							break
						}
					}
				}

				if !skipHandler {
					handler.handler(n, &attr)
				}
			}
		}
		// switch n.Data {
		// case "a":
		// 	for _, a := range n.Attr {
		// 		if a.Key == "href" {
		// 			fmt.Println(a.Val)
		// 			break
		// 		}
		// 	}
		// }
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(c, handlers)
	}
}

// Cleanup a fetched document from the engine-specific markup
func Cleanup(source string) string {
	// TODO: provide io.Reader instead
	doc, _ := html.Parse(strings.NewReader(source))

	handlers := []Handler{
		Handler{element: "a", attr: "href", handler: func(n *html.Node, a *html.Attribute) { fmt.Println("BOO>> :", n.Data, a.Val) }},
		Handler{element: "a", attr: "href", attrVal: "/sign_in", handler: func(n *html.Node, a *html.Attribute) { fmt.Println(">> YAY << :", n.Data, a.Val) }},
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
