package google

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"bytes"

	"golang.org/x/net/html"
)

// Cleanup a fetched document from the engine-specific markup
func Cleanup(source string) string {
	doc, _ := html.Parse(strings.NewReader(source))

	// ...
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				for _, a := range n.Attr {
					if a.Key == "href" {
						fmt.Println(a.Val)
						break
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	// ...

	// iframe(&doc)
	// onload(&doc)
	var resBytes []byte
	res := bytes.NewBuffer(resBytes)
	resWriter := bufio.NewWriter(res)
	html.Render(resWriter, doc)
	resWriter.Flush()

	return string(resBytes)
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
