package google

import (
	"bufio"
	"regexp"
	"strings"

	"bytes"

	"net/url"

	"fmt"

	"golang.org/x/net/html"
)

// Handler for the DOM node, located by element and attr
type Handler struct {
	element string
	attr    string
	attrVal string
	handler func(n *html.Node, attrIndex int)
}

func traverse(n *html.Node, handlers *[]Handler) {
	// if n.Data == "span" {
	// 	fmt.Println("YAY", n.Type)
	// }

	// TODO: our manipulations somehow brakes the traverse (incorrect Next on changes/removal?)
	// if false && n.Type == html.ElementNode {
	if n.Type == html.ElementNode {
		for _, handler := range *handlers {

			if handler.element == "" || n.Data == handler.element {
				skipHandler := false
				attrIndex := -1
				if handler.attr != "" {
					for i, a := range n.Attr {
						if a.Key == handler.attr {
							if handler.attrVal == "" || a.Val == handler.attrVal {
								attrIndex = i
							} else {
								skipHandler = true
							}

							break
						}
					}
				}

				if !skipHandler {
					handler.handler(n, attrIndex)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(c, handlers)
	}
}

func removeAttr(s []html.Attribute, i int) []html.Attribute {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

var nodesForRemoval []*html.Node

// schedules node removal to prevent skips on traverse
func scheduleNodeRemoval(n *html.Node) {
	nodesForRemoval = append(nodesForRemoval, n)
}

func removeScheduledNodes() {
	for i := 0; i < len(nodesForRemoval); i++ {
		fmt.Println(i, len(nodesForRemoval))
		n := nodesForRemoval[i]
		n.Parent.RemoveChild(n)
	}
}

// Cleanup a fetched document from the engine-specific markup
func Cleanup(source string) string {
	nodesForRemoval = nil
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
		Handler{element: "base", handler: func(n *html.Node, attrIndex int) {
			// TODO: base removal breaks view heavily without dead requests - to confirm
			// scheduleNodeRemoval(n)
		}},

		// href fix
		// url base replacement
		Handler{element: "a", attr: "href", handler: func(n *html.Node, attrIndex int) {
			a := &n.Attr[attrIndex]

			// href fix
			if strings.Contains(a.Val, "https://translate.googleusercontent.com/translate") {
				hrefURL, _ := url.Parse(a.Val)
				a.Val = hrefURL.Query()["u"][0]
			}

			// url base replacement
			a.Val = regexp.MustCompile("https?://([^./]*.?)"+sourceDomain+"/?").ReplaceAllLiteralString(a.Val, "/")
		}},

		// ENGINE
		// iframe
		Handler{element: "iframe", attr: "src", handler: func(n *html.Node, attrIndex int) {
			if strings.Contains(n.Attr[attrIndex].Val, "translate.google.com") {
				scheduleNodeRemoval(n)
			}
		}},

		// onload
		Handler{attr: "onload", handler: func(n *html.Node, attrIndex int) {
			if attrIndex == -1 {
				return
			}
			n.Attr = removeAttr(n.Attr, attrIndex)
		}},

		// css
		Handler{element: "style", handler: func(n *html.Node, attrIndex int) {
			if strings.Contains(n.FirstChild.Data, ".google-src-text") {
				scheduleNodeRemoval(n)
			}
		}},

		// scripts
		Handler{element: "script", attr: "src", handler: func(n *html.Node, attrIndex int) {
			if attrIndex != -1 {
				if strings.Contains(n.Attr[attrIndex].Val, "translate_c") {
					scheduleNodeRemoval(n)
				}
			} else if regexp.MustCompile("(_intlStrings|function ti_|_setupIW|performance)").FindStringIndex(n.FirstChild.Data) != nil {
				scheduleNodeRemoval(n)
			}
		}},

		// span wrappers
		// TODO: "span" not located, why?!
		Handler{element: "span", attr: "class", attrVal: "notranslate", handler: func(n *html.Node, attrIndex int) {
			// Handler{element: "span", handler: func(n *html.Node, attrIndex int) {
			// Handler{element: "span", handler: func(n *html.Node, attrIndex int) {
			// return
			fmt.Println("!!!!!", n.LastChild)
			scheduleNodeRemoval(n.FirstChild)
			n.Attr = []html.Attribute{}
			// content := html.Node{
			// 	Data: n.LastChild.Data,
			// }
			// n.Parent.AppendChild(&content)
		}},
	}

	traverse(doc, &handlers)
	removeScheduledNodes()

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
