package gozon

import "net/url"

type Client struct {
	CountryCode string //ISO country code, example: USD, MX
	Currency    string //ISO currency, example: USD, EUR
	ProxyURL    *url.URL
	Tries       int
}

func DefaulClient() Client {
	client := Client{
		CountryCode: "US",
		Currency:    "USD",
		ProxyURL:    nil,
		Tries:       6,
	}
	return client
}
func NewClient(countryCode, currency string, proxyURL *url.URL, tries int) Client {
	client := Client{
		CountryCode: countryCode,
		Currency:    currency,
		ProxyURL:    proxyURL,
		Tries:       tries,
	}
	return client
}

// make sure your url contains "?th=1&psc=1" for better results, sometimes if you don't add it it won't show the price
// it will return the data along with how many tries did for extracting the data
func (cl Client) GetFromURL(productURL string) (Data, int, error) {
	return GetFromURL(productURL, cl.CountryCode, cl.Currency, cl.ProxyURL, cl.Tries)
}

func (cl Client) GetFromID(label, id string) (Data, int, error) {
	return GetFromID(label, id, cl.CountryCode, cl.Currency, cl.ProxyURL, cl.Tries)
}

func (cl Client) GetMainURLs() ([]string, int, error) {
	return GetMainURLs(cl.CountryCode, cl.Currency, cl.ProxyURL, cl.Tries)
}
