package gozon

import (
	"net/url"
	"strings"
)

func ParseProxy(urlToParse, userName, password string) (*url.URL, error) {
	urlToUse, err := url.Parse(urlToParse)
	if err != nil {
		return nil, err
	}
	urlToUse.User = url.UserPassword(userName, password)
	return urlToUse, nil
}
func RemoveSpace(value string) string {
	for strings.Contains(value, " ") {
		value = strings.ReplaceAll(value, " ", " ")
	}
	value = strings.ReplaceAll(value, "\t", " ")
	value = strings.ReplaceAll(value, "\n", " ")
	for strings.Contains(value, "  ") {
		value = strings.ReplaceAll(value, "  ", " ")
	}
	value = strings.TrimSpace(value)
	return value
}
