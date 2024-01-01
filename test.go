package gozon

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func test1() {
	productURL := "https://www.amazon.com/Fruit-Loom-Recycled-Underwear-Greystone/dp/B08D2YRH6K"
	client := DefaulClient()
	client.Currency = "EUR"
	data, _, err := client.GetFromURL(productURL)
	if err != nil {
		log.Println("test1 1 -> err: ", err)
		return
	}
	f, err := os.Create("./globaldata.json")
	if err != nil {
		log.Println("test1 2 -> err: ", err)
		return
	}
	json.NewEncoder(f).Encode(data)
	f.Close()
}
func Test() {
	if err := os.MkdirAll("./test", 0644); err != nil {
		log.Println("test 1 -> err: ", err)
		return
	}
	client := DefaulClient()
	productsURLs, _, err := client.GetMainURLs()
	if err != nil {
		log.Println("test 2 -> err: ", err)
		return
	}
	var datas []Data
	for i, productURL := range productsURLs {
		folderPath := fmt.Sprintf("./test/%d/images", i)
		os.MkdirAll(folderPath, 0644)
		data, _, err := client.GetFromURL(productURL)
		if err != nil {
			log.Println("test 3 -> err: ", err)
			continue
		}
		if err := data.SetImages(client.ProxyURL); err != nil {
			log.Println("test2 4 -> err: ", err)
		}
		for j, imgs := range data.Images {
			if len(imgs.Large.Content) != 0 {
				f_name1 := fmt.Sprintf("./test/%d/images/%d_large%s", i, j, imgs.Large.Extension)
				os.WriteFile(f_name1, imgs.Large.Content, 0644)
			}
			if len(imgs.Thumb.Content) != 0 {
				f_name2 := fmt.Sprintf("./test/%d/images/%d_thumb%s", i, j, imgs.Thumb.Extension)
				os.WriteFile(f_name2, imgs.Thumb.Content, 0644)
			}
			if len(imgs.HiRes.Content) != 0 {
				f_name3 := fmt.Sprintf("./test/%d/images/%d_hires%s", i, j, imgs.HiRes.Extension)
				os.WriteFile(f_name3, imgs.HiRes.Content, 0644)
			}
		}
		filePath := fmt.Sprintf("./test/%d/data.json", i)
		f, err := os.Create(filePath)
		if err != nil {
			log.Println("test 5 -> err: ", err)
			continue
		}
		json.NewEncoder(f).Encode(data)
		f.Close()
		log.Printf("Progress: %d/%d\n", i+1, len(productsURLs))
		datas = append(datas, data)
	}
	f, err := os.Create("./test/globaldata.json")
	if err != nil {
		log.Println("test 6 -> err: ", err)
		return
	}
	json.NewEncoder(f).Encode(datas)
	f.Close()
}
func TestWithProxy() {
	if err := os.MkdirAll("./test", 0644); err != nil {
		log.Println("test 1 -> err: ", err)
		return
	}
	proxyURL, err := ParseProxy("http://[IP | domain]:[port]", "username", "password")
	if err != nil {
		log.Println("test 2 -> err: ", err)
		return
	}
	client := NewClient("US", proxyURL, 6)
	productsURLs, _, err := client.GetMainURLs()
	if err != nil {
		log.Println("test 3 -> err: ", err)
		return
	}
	var datas []Data
	for i, productURL := range productsURLs {
		folderPath := fmt.Sprintf("./test/%d/images", i)
		os.MkdirAll(folderPath, 0644)
		data, _, err := client.GetFromURL(productURL)
		if err != nil {
			log.Println("test 4 -> err: ", err)
			continue
		}
		if err := data.SetImages(client.ProxyURL); err != nil {
			log.Println("test2 5 -> err: ", err)
		}
		for j, imgs := range data.Images {
			if len(imgs.Large.Content) != 0 {
				f_name1 := fmt.Sprintf("./test/%d/images/%d_large%s", i, j, imgs.Large.Extension)
				os.WriteFile(f_name1, imgs.Large.Content, 0644)
			}
			if len(imgs.Large.Content) != 0 {
				f_name2 := fmt.Sprintf("./test/%d/images/%d_thumb%s", i, j, imgs.Thumb.Extension)
				os.WriteFile(f_name2, imgs.Thumb.Content, 0644)
			}
			if len(imgs.Large.Content) != 0 {
				f_name3 := fmt.Sprintf("./test/%d/images/%d_hires%s", i, j, imgs.HiRes.Extension)
				os.WriteFile(f_name3, imgs.HiRes.Content, 0644)
			}
		}
		filePath := fmt.Sprintf("./test/%d/data.json", i)
		f, err := os.Create(filePath)
		if err != nil {
			log.Println("test 6 -> err: ", err)
			continue
		}
		json.NewEncoder(f).Encode(data)
		f.Close()
		log.Printf("Progress: %d/%d\n", i+1, len(productsURLs))
		datas = append(datas, data)
	}
	f, err := os.Create("./test/globaldata.json")
	if err != nil {
		log.Println("test 7 -> err: ", err)
		return
	}
	json.NewEncoder(f).Encode(datas)
	f.Close()
}
