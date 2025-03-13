package gozon

import (
	"bytes"
	"encoding/json"
	"html"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/johnbalvin/gozon/trace"

	"github.com/PuerkitoBio/goquery"
)

func ParseBodyDetails(body []byte) (Data, error) {
	dataRaw, err := parseBodyDetails(body)
	if err != nil {
		return Data{}, trace.NewOrAdd(1, "main", "ParseBodyDetails", err, "")
	}
	datFormated := dataRaw.Standardize()
	return datFormated, nil
}
func parseBodyDetails(body []byte) (dataRaw, error) {
	reader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return dataRaw{}, trace.NewOrAdd(1, "main", "parseBodyDetails", err, "")
	}
	var data_raw string
	var imgs_data string

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		htmlValue, _ := s.Html()
		htmlValue = html.UnescapeString(htmlValue)
		if strings.Contains(htmlValue, "twister-js-init-dpx-data") {
			value := regxpData.FindString(htmlValue)
			value = strings.ReplaceAll(value, "dataToReturn = ", "")
			data_raw = strings.TrimSuffix(value, ";")
		} else if strings.Contains(htmlValue, "ImageBlockBTF") {
			/*value := regxImgsData.FindString(htmlValue)
			value = strings.ReplaceAll(value, `jQuery.parseJSON('`, "")
			value = strings.ReplaceAll(value, "');", "")
			//yeah, it works this way, trust me
			value = strings.ReplaceAll(value, `\\`, `\`)
			value = strings.ReplaceAll(value, `\`, `\\`)
			os.WriteFile("./imgs2.json", []byte(value), 0644)
			fmt.Println("-----2")
			//------------
			imgs_data = value*/
		} else if strings.Contains(htmlValue, "ImageBlockATF") {
			htmlValue = RemoveSpace(htmlValue)
			value := regxImgsData2.FindString(htmlValue)
			value = strings.TrimPrefix(value, `'initial':`)
			value = strings.TrimSuffix(value, `}, 'colorToAsin'`)
			//------------
			imgs_data = value
		}
	})
	data := dataRaw{
		Variations: make(map[string][]string),
	}
	if imgs_data == "" {
		return dataRaw{}, trace.NewOrAdd(2, "main", "parseBodyDetails", trace.ErrEmpty, "")
	}
	if err := json.Unmarshal([]byte(imgs_data), &data.Images); err != nil {
		return dataRaw{}, trace.NewOrAdd(11, "main", "parseBodyDetails", err, "")
	}
	//parsing address
	addr1 := RemoveSpace(doc.Find("#glow-ingress-line1").Text())
	addr2 := RemoveSpace(doc.Find("#glow-ingress-line2").Text())
	if strings.Contains(addr2, "Update location") {
		data.ShippingAddress = strings.ReplaceAll(addr1, "Delivering to ", "")
	} else {
		data.ShippingAddress = addr2
	}
	//
	data.Alterts = RemoveSpace(doc.Find("#product-alert-grid_feature_div").Text())
	data.Title = RemoveSpace(doc.Find("#productTitle").Text())
	data.Asin.Me = RemoveSpace(doc.Find("#ASIN").AttrOr("value", ""))
	data.MerchandID = RemoveSpace(doc.Find("#merchantID").AttrOr("value", ""))
	data.URL = RemoveSpace(doc.Find(`[rel="canonical"]`).AttrOr("href", ""))
	if data_raw != "" {
		data_raw_dimensions := regxdimensionsValues.FindString(data_raw)
		data_raw_dimensions = strings.ReplaceAll(data_raw_dimensions, `dimensionValuesDisplayData" : `, "")
		data_raw_dimensions = strings.TrimSuffix(data_raw_dimensions, ",")
		if err := json.Unmarshal([]byte(data_raw_dimensions), &data.Variations); err != nil {
			return dataRaw{}, trace.NewOrAdd(3, "main", "parseBodyDetails", err, "")
		}
		parentAsin := regexParentAsin.FindString(data_raw)
		parentAsin = strings.ReplaceAll(parentAsin, `"parentAsin" : "`, "")
		parentAsin = strings.ReplaceAll(parentAsin, `"`, "")
		data.Asin.Parent = parentAsin
		variations := regxVariations.FindString(data_raw)
		variations = RemoveSpace(strings.ReplaceAll(variations, `variationValues" :`, ""))
		if err := json.Unmarshal([]byte(variations), &data.VariationDisplayLabels); err != nil {
			return dataRaw{}, trace.NewOrAdd(4, "main", "parseBodyDetails", err, "")
		}
		dimestionsRaw := regxdimensions.FindString(data_raw)
		dimestionsRaw = strings.ReplaceAll(dimestionsRaw, `dimensions" : `, "")
		dimestionsRaw = strings.TrimSuffix(dimestionsRaw, ",")
		var dimensions []string
		if err := json.Unmarshal([]byte(dimestionsRaw), &dimensions); err != nil {
			return dataRaw{}, trace.NewOrAdd(5, "main", "parseBodyDetails", err, "")
		}
		data.Dimensions = dimensions
	}
	ratingS := RemoveSpace(doc.Find("#acrPopover a .a-size-base").First().Text())
	if ratingS != "" {
		rating, err := strconv.ParseFloat(ratingS, 32)
		if err != nil {
			return dataRaw{}, trace.NewOrAdd(6, "main", "parseBodyDetails", err, "")
		}
		data.Rating = float32(rating)
	}
	doc.Find("#productFactsDesktopExpander ul li").Each(func(i int, s *goquery.Selection) {
		value := RemoveSpace(s.Text())
		if value == "" {
			return
		}
		data.MainPanelDesc.AboutThis = append(data.MainPanelDesc.AboutThis, value)
	})
	doc.Find("#feature-bullets ul li").Each(func(i int, s *goquery.Selection) {
		value := RemoveSpace(s.Text())
		if value == "" {
			return
		}
		data.MainPanelDesc.AboutThis = append(data.MainPanelDesc.AboutThis, value)
	})
	doc.Find("#glance_icons_div table").First().Find("tr tr").Each(func(i int, s *goquery.Selection) {
		var iconPanel IconPanel
		s.Find("td").Each(func(j int, s2 *goquery.Selection) {
			switch j {
			case 0:
				src := RemoveSpace(s2.Find(".a-image-wrapper").AttrOr("data-a-image-source", ""))
				iconPanel.URL = src
			case 1:
				s2.Find("span").Each(func(k int, s3 *goquery.Selection) {
					value := RemoveSpace(s3.Text())
					switch k {
					case 0:
						iconPanel.Label = value
					case 1:
						iconPanel.Value = value
					}
				})
			}
		})
		if iconPanel.Label != "" {
			data.MainPanelDesc.IconsPanel = append(data.MainPanelDesc.IconsPanel, iconPanel)
		}
	})
	doc.Find(".product-facts-detail").Each(func(i int, s *goquery.Selection) {
		var details LabelValue
		s.Find(".a-color-base").Each(func(j int, s2 *goquery.Selection) {
			value := RemoveSpace(s2.Text())
			switch j {
			case 0:
				details.Label = value
			case 1:
				details.Value = value
			}
		})
		if details.Label != "" {
			data.MainPanelDesc.Details = append(data.MainPanelDesc.Details, details)
		}
	})
	doc.Find("#productOverview_feature_div table").First().Find("tr").Each(func(i int, s *goquery.Selection) {
		var labelValue LabelValue
		s.Find("td").Each(func(j int, s2 *goquery.Selection) {
			switch j {
			case 0:
				value := RemoveSpace(s2.Text())
				labelValue.Label = value
			case 1:
				valueFull := RemoveSpace(s2.Find(".a-truncate-full").Text())
				valueNormal := RemoveSpace(s2.Find(".a-size-base").Text())
				if valueFull == "" {
					labelValue.Value = valueNormal
				} else {
					labelValue.Value = valueFull
				}
			}
		})
		if labelValue.Label != "" {
			data.MainPanelDesc.Details = append(data.MainPanelDesc.Details, labelValue)
		}
	})
	available := RemoveSpace(doc.Find("#availability_feature_div").Text())
	if strings.Contains(available, "Currently unavailable") || strings.Contains(available, "Temporarily out of stock") {
		return data, nil
	}
	data.Available = true
	doc.Find(".a-price-range .a-offscreen").Each(func(i int, s *goquery.Selection) {
		priceRaw := RemoveSpace(s.Text())
		splited := strings.Split(priceRaw, "")
		if len(splited) < 2 {
			return
		}
		currencySymbol := splited[0]
		priceRaw = strings.Join(splited[1:], "")
		price, err := strconv.ParseFloat(priceRaw, 32)
		if err != nil {
			errData := trace.NewOrAdd(7, "main", "parseBodyDetails", err, "")
			log.Println(errData)
			return
		}
		priceConverted := float32(price)
		switch i {
		case 0:
			data.Price.Low = priceConverted
		case 1:
			data.Price.High = priceConverted
		}
		data.Price.CurrencySymbol = currencySymbol
	})
	if data.Price.CurrencySymbol == "" {
		priceRaw := doc.Find("#corePrice_feature_div .a-price .a-offscreen").First().Text()
		symbol, price, err := parsePriceSymbol(priceRaw)
		if err != nil {
			return dataRaw{}, trace.NewOrAdd(8, "main", "parseBodyDetails", err, "")
		}
		data.Price.Low = price
		data.Price.High = price
		data.Price.CurrencySymbol = symbol
	}
	basisPriceS := doc.Find(".basisPrice .a-price span").First().Text()
	if basisPriceS != "" {
		_, price, err := parsePriceSymbol(basisPriceS)
		if err != nil {
			return dataRaw{}, trace.NewOrAdd(9, "summary", "parseBodyDetails", err, "")
		}
		data.Price.Base = price
	}
	savingsPercentageS := doc.Find(".savingsPercentage").First().Text()
	if savingsPercentageS != "" {
		savingsPercentageS = strings.ReplaceAll(savingsPercentageS, "%", "")
		savingsPercentageS = strings.ReplaceAll(savingsPercentageS, "-", "")
		savingsPercentage, err := strconv.Atoi(savingsPercentageS)
		if err != nil {
			return dataRaw{}, trace.NewOrAdd(10, "main", "parseBodyDetails", err, "")
		}
		data.Price.SavingsPercentage = savingsPercentage
	}
	return data, nil
}

func ParseBodyMainList(body []byte) ([]string, error) {
	reader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, trace.NewOrAdd(1, "main", "ParseBodyMainList", err, "")
	}
	urlsMap := make(map[string]bool)
	doc.Find(".a-link-normal").Each(func(i int, s *goquery.Selection) {
		value := s.AttrOr("href", "")
		if !strings.Contains(value, "/dp/") {
			return
		}
		urlParsed, err := url.Parse(value)
		if err != nil {
			errData := trace.NewOrAdd(2, "main", "ParseBodyMainList", err, "")
			log.Println(errData)
			return
		}
		urlParsed.Host = "www.amazon.com"
		urlParsed.Scheme = "https"
		urlParsed.RawQuery = "th=1&psc=1"
		urlParsed.Path = strings.TrimSuffix(urlParsed.Path, "/")
		urlsMap[urlParsed.String()] = true
	})
	var urls []string
	for url := range urlsMap {
		urls = append(urls, url)
	}
	return urls, nil
}
