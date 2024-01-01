package gozon

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/johnbalvin/gozon/trace"
)

func GetMainURLs(countryCode, currency string, proxyURL *url.URL, tries int) ([]string, int, error) {
	maxTry := 5
	if tries != 0 {
		maxTry = tries
	}
	for try := 0; ; try++ {
		if try == maxTry {
			return nil, 0, trace.NewOrAdd(1, "main", "GetFromURL", trace.ErrCaptcha, "")
		}
		data, err := getMainURLs(countryCode, currency, proxyURL)
		if err == nil {
			return data, try + 1, nil
		}
		if trace.GetMainErr(err) != trace.ErrCaptcha {
			return nil, 0, trace.NewOrAdd(2, "main", "GetFromURL", err, "")
		}
	}
}
func getMainURLs(countryCode, currency string, proxyURL *url.URL) ([]string, error) {
	req, err := http.NewRequest("GET", epList, nil)
	if err != nil {
		return nil, trace.NewOrAdd(1, "main", "getMainURLs", err, "")
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Accept-Language", "en")
	req.Header.Add("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Add("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "none")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	var cookies []string
	if countryCode != "" {
		cookie := fmt.Sprintf(`sp-cdn="L5Z9:%s"`, countryCode)
		cookies = append(cookies, cookie)
	}
	if currency != "" {
		cookie := "i18n-prefs=" + currency
		cookies = append(cookies, cookie)
	}
	if len(cookies) != 0 {
		cookieValue := strings.Join(cookies, ";")
		req.Header.Set("Cookie", cookieValue)
	}
	transport := &http.Transport{
		MaxIdleConnsPerHost: 30,
		DisableKeepAlives:   true,
		TLSClientConfig: &tls.Config{
			Renegotiation:      tls.RenegotiateOnceAsClient,
			InsecureSkipVerify: true,
		},
	}
	if proxyURL != nil {
		transport.Proxy = http.ProxyURL(proxyURL)
	}
	client := &http.Client{
		Timeout: time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: transport,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, trace.NewOrAdd(2, "main", "getMainURLs", err, "")
	}
	if resp.StatusCode != 200 {
		errData := fmt.Sprintf("status: %d headers: %+v", resp.StatusCode, resp.Header)
		return nil, trace.NewOrAdd(3, "main", "getMainURLs", trace.ErrStatusCode, errData)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, trace.NewOrAdd(4, "main", "getMainURLs", err, "")
	}
	if strings.Contains(string(body), "Sorry, we just need to make sure you're not a robot") {
		return nil, trace.NewOrAdd(5, "main", "getMainURLs", trace.ErrCaptcha, "")
	}
	urls, err := ParseBodyMainList(body)
	if err != nil {
		return nil, trace.NewOrAdd(6, "main", "getMainURLs", err, "")
	}
	return urls, nil
}
