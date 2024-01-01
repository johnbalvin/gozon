package gozon

import "regexp"

const (
	epList = "https://www.amazon.com"
)

var (
	ep                   = "https://www.amazon.com/%s/dp/%s?th=1&psc=1"
	regxpData            = regexp.MustCompile(`dataToReturn = \{(.|\n)+?\};`)
	regxdimensionsValues = regexp.MustCompile(`dimensionValuesDisplayData" : {.+?},`)
	regxdimensions       = regexp.MustCompile(`dimensions" : \[.+?\],`)
	regexParentAsin      = regexp.MustCompile(`"parentAsin" : ".+?"`)
	regxVariations       = regexp.MustCompile(`variationValues" : {.+?}`)
	regxImgsData2        = regexp.MustCompile(`'initial':.+?'colorToAsin'`)
	regxPrice            = regexp.MustCompile(`\d.+`)
)
