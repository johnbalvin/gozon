package gozon

import (
	"net/url"
	"regexp"
	"strings"
)

var regexSpace = regexp.MustCompile(`\[\t\n\]+|[\s ]{2,}`)

func RemoveSpace(value string) string {
	return regexSpace.ReplaceAllString(strings.TrimSpace(value), " ")
}

func ParseProxy(urlToParse, userName, password string) (*url.URL, error) {
	urlToUse, err := url.Parse(urlToParse)
	if err != nil {
		return nil, err
	}
	urlToUse.User = url.UserPassword(userName, password)
	return urlToUse, nil
}
