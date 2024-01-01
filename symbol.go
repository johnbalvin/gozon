package gozon

import (
	"strconv"
	"strings"

	"github.com/johnbalvin/gozon/trace"
)

func parsePriceSymbol(priceRaw string) (string, float32, error) {
	priceRaw = strings.ReplaceAll(priceRaw, ",", "")
	priceNumber := regxPrice.FindString(priceRaw)
	priceCurrency := RemoveSpace(strings.ReplaceAll(priceRaw, priceNumber, ""))
	splited := strings.Split(priceRaw, "")
	if len(splited) < 2 {
		return "", 0, trace.NewOrAdd(1, "summary", "parsePriceSymbol", trace.ErrParameter, priceRaw)
	}
	price, err := strconv.ParseFloat(priceNumber, 32)
	if err != nil {
		return "", 0, trace.NewOrAdd(2, "summary", "parsePriceSymbol", err, priceRaw)
	}
	priceConverted := float32(price)
	return priceCurrency, priceConverted, nil
}
