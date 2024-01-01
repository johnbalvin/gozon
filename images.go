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

func (data *Data) SetImages(proxyURL *url.URL) error {
	for i := range data.Images {
		if err := data.Images[i].SetImages(proxyURL); err != nil {
			return trace.NewOrAdd(1, "main", "Data SetImages", err, "")
		}
		fmt.Printf("Setting images: %d/%d\n", i+1, len(data.Images))
	}
	return nil
}
func (imgs *Imgs) SetImages(proxyURL *url.URL) error {
	if imgs.Large.URL != "" {
		contentType1, img1, err := GetImg(imgs.Large.URL, proxyURL)
		if err != nil {
			return trace.NewOrAdd(1, "main", "Imgs SetImages", err, "")
		}
		imgs.Large.ContentType = contentType1
		imgs.Large.Content = img1
	}
	if imgs.Thumb.URL != "" {
		contentType2, img2, err := GetImg(imgs.Thumb.URL, proxyURL)
		if err != nil {
			return trace.NewOrAdd(2, "main", "Imgs SetImages", err, "")
		}
		imgs.Thumb.ContentType = contentType2
		imgs.Thumb.Content = img2
	}
	if imgs.HiRes.URL != "" {
		contentType3, img3, err := GetImg(imgs.HiRes.URL, proxyURL)
		if err != nil {
			return trace.NewOrAdd(3, "main", "Imgs SetImages", err, "")
		}
		imgs.HiRes.ContentType = contentType3
		imgs.HiRes.Content = img3
	}
	return nil
}
func GetImg(imgURL string, proxyURL *url.URL) (string, []byte, error) {
	req, err := http.NewRequest("GET", imgURL, nil)
	if err != nil {
		return "", nil, trace.NewOrAdd(1, "main", "GetImg", err, "")
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
		return "", nil, trace.NewOrAdd(2, "main", "GetImg", err, "")
	}
	if resp.StatusCode != 200 {
		errData := fmt.Sprintf("status: %d headers: %+v", resp.StatusCode, resp.Header)
		return "", nil, trace.NewOrAdd(3, "main", "GetImg", trace.ErrStatusCode, errData)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "image") {
		errData := fmt.Sprintf("headers: %+v", resp.Header)
		return "", nil, trace.NewOrAdd(3, "main", "GetImg", err, errData)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, trace.NewOrAdd(4, "main", "GetImg", err, "")
	}
	return contentType, body, nil
}
