package gozon

import (
	"log"
	"path/filepath"
)

func (dataRaw dataRaw) Standardize() Data { //In case a field changed, it won't messup a project that uses this library
	data := Data{
		Available:       dataRaw.Available,
		Rating:          dataRaw.Rating,
		URL:             dataRaw.URL,
		ShippingAddress: dataRaw.ShippingAddress,
		Asin:            dataRaw.Asin,
		Title:           dataRaw.Title,
		MerchandID:      dataRaw.MerchandID,
		Price:           dataRaw.Price,
		MainPanelDesc:   dataRaw.MainPanelDesc,
	}

	for label, values := range dataRaw.VariationDisplayLabels {
		labelValues := LabelValues{
			Label:  label,
			Values: values,
		}
		data.Variations.Labels = append(data.Variations.Labels, labelValues)
	}
	for label, values := range dataRaw.Variations {
		combination := Combination{
			Asin:   label,
			Values: values,
		}
		data.Variations.Combinations = append(data.Variations.Combinations, combination)
		if len(values) != len(dataRaw.Dimensions) {
			log.Println("Dimension don't match with values len, please report this bug", data.URL)
			continue
		}
	}
	for _, imgRaw := range dataRaw.Images {
		img := Imgs{
			Large: URLImg{
				URL:       imgRaw.Large,
				Extension: filepath.Ext(imgRaw.Large),
			},
			Thumb: URLImg{
				URL:       imgRaw.Thumb,
				Extension: filepath.Ext(imgRaw.Thumb),
			},
			HiRes: URLImg{
				URL:       imgRaw.HiRes,
				Extension: filepath.Ext(imgRaw.HiRes),
			},
			Variant: imgRaw.Variant,
		}
		data.Images = append(data.Images, img)
	}
	return data
}
