package google

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"bytes"

	"net/url"

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
	if n.Data == "span" {
		fmt.Println("YAY", n.Type)
	}

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
		Handler{element: "base", handler: func(n *html.Node, attrIndex int) { n.Parent.RemoveChild(n) }},

		// href fix
		// url base replacement
		Handler{element: "a", attr: "href", handler: func(n *html.Node, attrIndex int) {
			a := &n.Attr[attrIndex]

			// href fix
			if strings.Contains(a.Val, "https://translate.googleusercontent.com/translate") {
				hrefURL, _ := url.Parse(a.Val)
				fmt.Println(">>>>", hrefURL.Query())
				a.Val = hrefURL.Query()["u"][0]
			}

			// url base replacement
			a.Val = regexp.MustCompile("https?://([^./]*.?)"+sourceDomain+"/?").ReplaceAllLiteralString(a.Val, "/")
		}},

		// ENGINE
		// iframe
		Handler{element: "iframe", attr: "src", handler: func(n *html.Node, attrIndex int) {
			if strings.Contains(n.Attr[attrIndex].Val, "translate.google.com") {
				n.Parent.RemoveChild(n)
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
				n.Parent.RemoveChild(n)
			}
		}},

		// scripts
		Handler{element: "script", attr: "src", handler: func(n *html.Node, attrIndex int) {
			if attrIndex != -1 {
				if strings.Contains(n.Attr[attrIndex].Val, "translate_c") {
					n.Parent.RemoveChild(n)
				}
			}

			if regexp.MustCompile("(_intlStrings|function ti_|_setupIW|performance)").FindStringIndex(n.FirstChild.Data) != nil {
				n.Parent.RemoveChild(n)
			}
		}},

		// span wrappers
		// TODO: "span" not located, why?!
		// Handler{element: "span", attr: "class", attrVal: "notranslate", handler: func(n *html.Node, attrIndex int) {
		// Handler{element: "span", handler: func(n *html.Node, attrIndex int) {
		Handler{element: "span", handler: func(n *html.Node, attrIndex int) {
			fmt.Println("!!!!!", n.Data, n)
			n.Parent.RemoveChild(n)
			// n.RemoveChild(n.LastChild)
			n.Parent.AppendChild(n.LastChild)
		}},
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
