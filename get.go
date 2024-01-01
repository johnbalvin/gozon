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

func GetFromID(label, id, countryCode, currency string, proxyURL *url.URL, tries int) (Data, int, error) {
	url_to_use := fmt.Sprintf(ep, label, id)
	data, tries, err := GetFromURL(url_to_use, countryCode, currency, proxyURL, tries)
	if err != nil {
		return Data{}, 0, trace.NewOrAdd(1, "main", "GetFromID", err, "")
	}
	return data, tries, nil
}

// GetFromURL gets the information with the product url
// //the code is optimize to work on this format:
// https://www.amazon.com/[label]/dp/[id]?th=1&psc=1
// proxyURL is the proxy you would like to use, its recommended to use residential proxies
// countryCode, set to any country code, this will change the end result, probably the pricing or if available or not
// -->>>iF COUNTRY CODE IS NOT PRESENT AMAZON WILL CONSIDER THE IP LOCATION <-----
// tries is how many tries you want to make before throwing an error, this is usefull because sometimes amazom will require a captcha
// make soure your url contains "?th=1&psc=1" for better results, sometimes if you don't add it it wont' show the price
// it will return the data along with how many tries did for extracting the data
func GetFromURL(productURL, countryCode, currency string, proxyURL *url.URL, tries int) (Data, int, error) {
	maxTry := 5
	if tries != 0 {
		maxTry = tries
	}
	for try := 0; ; try++ {
		if try == maxTry {
			return Data{}, 0, trace.NewOrAdd(1, "main", "GetFromURL", trace.ErrCaptcha, "")
		}
		data, err := getFromURL(productURL, countryCode, currency, proxyURL)
		if err == nil {
			return data, try + 1, nil
		}
		if trace.GetMainErr(err) != trace.ErrCaptcha {
			return Data{}, 0, trace.NewOrAdd(2, "main", "GetFromURL", err, "")
		}
	}
}
func getFromURL(productURL, countryCode, currency string, proxyURL *url.URL) (Data, error) {
	req, err := http.NewRequest("GET", productURL, nil)
	if err != nil {
		return Data{}, trace.NewOrAdd(1, "main", "getFromURL", err, "")
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Accept-Language", "en")
	req.Header.Add("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Add("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "none")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
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
		return Data{}, trace.NewOrAdd(2, "main", "getFromURL", err, "")
	}
	if resp.StatusCode != 200 {
		errData := fmt.Sprintf("status: %d headers: %+v", resp.StatusCode, resp.Header)
		return Data{}, trace.NewOrAdd(3, "main", "getFromURL", trace.ErrStatusCode, errData)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Data{}, trace.NewOrAdd(4, "main", "getFromURL", err, "")
	}
	if strings.Contains(string(body), "Sorry, we just need to make sure you're not a robot") {
		return Data{}, trace.NewOrAdd(5, "main", "getFromURL", trace.ErrCaptcha, "")
	}
	data, err := ParseBodyDetails(body)
	if err != nil {
		return Data{}, trace.NewOrAdd(6, "main", "getFromURL", err, "")
	}
	return data, nil
}
