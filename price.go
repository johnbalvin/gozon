package gozon

/*
func GetPrices() error {
	url_parsed, err := url.Parse("https://www.amazon.com/gp/twister/dimension")
	if err != nil {
		return trace.NewOrAdd(1, "summary", "GetPrice", err, "")
	}
	query := url.Values{}
	query.Add("isDimensionSlotsAjax", "1")
	query.Add("asinList", "B08D328GNG,B08D2NYRVN,B08D2YQXCJ,B08D2YRH6K,B08D2L5234,B08D2YSGFC,B08D2Q549T,B08D326MGF")
	query.Add("vs", "1")
	query.Add("productTypeDefinition", "SHIRT")
	query.Add("productGroupId", "apparel_display_on_website")
	query.Add("parentAsin", "B08D322MJQ")
	query.Add("isPrime", "0")
	query.Add("isOneClickEnabled", "0")
	query.Add("originalHttpReferer", "https://www.amazon.com/Fruit-Loom-Recycled-Underwear-Greystone/dp/B08D2YRH6K")
	query.Add("deviceType", "web")
	query.Add("showFancyPrice", "false")
	query.Add("twisterFlavor", "twisterPlusDesktopConfigurator")
	url_parsed.RawQuery = query.Encode()
	req, err := http.NewRequest("GET", url_parsed.String(), nil)
	if err != nil {
		return trace.NewOrAdd(1, "summary", "GetPrice", err, "")
	}
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
	client := &http.Client{
		Timeout: time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return trace.NewOrAdd(2, "summary", "GetPrice", err, "")
	}
	if resp.StatusCode != 200 {
		errData := fmt.Sprintf("status: %d headers: %+v", resp.StatusCode, resp.Header)
		return trace.NewOrAdd(3, "summary", "GetPrice", err, errData)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return trace.NewOrAdd(4, "summary", "GetPrice", err, "")
	}
	os.WriteFile("./prices.json", body, 0644)
	return nil
}
*/
