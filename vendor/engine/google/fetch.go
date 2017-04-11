package google

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
)

func translationURL(sourceLangCode string, resultLangCode string, url string) string {
	return fmt.Sprintf("http://translate.google.com/translate?hl=en&sl=%s&tl=%s&u=%s&prev=hp", sourceLangCode, resultLangCode, url)
}

func fetch(url string) string {
	resp, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(resp.Body)
	body := string(bytes)
	return body
}

// Fetch a translated document for the URI given
func Fetch(sourceLangCode, resultLangCode, sourceUri string) string {
	// TODO: hey, we just testing!
	// return fetch(sourceUri)

	targetURL := translationURL(sourceLangCode, resultLangCode, sourceUri)

	re1 := regexp.MustCompile(`<iframe sandbox="allow-same-origin allow-forms allow-scripts" src="(http://translate.googleusercontent.com/translate_p\?[^"]+)`)
	matched1 := re1.FindStringSubmatch(fetch(targetURL))[1]

	re2 := regexp.MustCompile(`<meta http-equiv="refresh" content="0;URL=([^"]+)`)
	matched2 := re2.FindStringSubmatch(fetch(matched1))[1]

	finalURL := html.UnescapeString(matched2)

	return fetch(finalURL)
}
