# Amazon Product Information Extractor in Go

## Overview
This project is an open-source tool developed in Golang for extracting product information from Amazon. It's designed to be fast, efficient, and easy to use, making it an ideal solution for developers looking for Amazon product data.

## Features
- Extracts detailed product information from Amazon
- Implemented in Go for performance and efficiency
- Easy to integrate with existing Go projects
- The code is optimize to work on this format: ```https://www.amazon.com/[label]/dp/[id]?th=1&psc=1```

## Examples

### Quick testing

```Go
    package main

    import (
        "github.com/johnbalvin/gozon"
    )
    func main(){
        //you need to have write permissions, the result will be save inside folder "test"
        gozon.Test()
    }
```

### Basic data

```Go
    package main

    import (
        "encoding/json"
        "log"
        "os"
        "github.com/johnbalvin/gozon"
    )
    func main(){
        //you need to have write permissions, the result will be save inside folder "test"
        if err := os.MkdirAll("./test", 0644); err != nil {
            log.Println("test 1 -> err: ", err)
            return
        }
        productURL:="https://www.amazon.com/[label]/dp/[id]?th=1&psc=1"
        client := gozon.DefaulClient()
        data, _, err := client.GetFromURL(productURL)
        if err != nil {
            log.Println("test:2 -> err: ", err)
            return
        }
        f, err := os.Create("./test/data.json")
        if err != nil {
            log.Println("test:3 -> err: ", err)
            return
        }
        json.NewEncoder(f).Encode(data)
    }
```

### Basic data and images
```Go
    package main

    import (
        "encoding/json"
        "log"
        "os"
        "github.com/johnbalvin/gozon"
    )
    func main(){
        //you need to have write permissions, the result will be save inside folder "test"
        if err := os.MkdirAll("./test/images", 0644); err != nil {
            log.Println("test 1 -> err: ", err)
            return
        }
        productURL:="https://www.amazon.com/[label]/dp/[id]?th=1&psc=1"
        client := gozon.DefaulClient()
	    client.Currency = "EUR" //you can change the currency
        data, _, err := client.GetFromURL(productURL)
        if err != nil {
            log.Println("test:2 -> err: ", err)
            return
        }
        if err := data.SetImages(client.ProxyURL); err != nil {
            log.Println("test:3 -> err: ", err)
            return
        }
		for j, imgs := range data.Images {
			if len(imgs.Large.Content) != 0 {
				f_name1 := fmt.Sprintf("./test/images/%d_large%s", j, imgs.Large.Extension)
				os.WriteFile(f_name1, imgs.Large.Content, 0644)
			}
			if len(imgs.Thumb.Content) != 0 {
				f_name2 := fmt.Sprintf("./test/images/%d_thumb%s", j, imgs.Thumb.Extension)
				os.WriteFile(f_name2, imgs.Thumb.Content, 0644)
			}
			if len(imgs.HiRes.Content) != 0 {
				f_name3 := fmt.Sprintf("./test/images/%d_hires%s", j, imgs.HiRes.Extension)
				os.WriteFile(f_name3, imgs.HiRes.Content, 0644)
			}
		}
        f, err := os.Create("./test/data.json")
        if err != nil {
            log.Println("test:4 -> err: ", err)
            return
        }
        json.NewEncoder(f).Encode(data)
    }
```

### With proxy

```Go
    package main

    import (
        "encoding/json"
        "log"
        "os"
        "github.com/johnbalvin/gozon"
    )
    func main(){
        proxyURL, err := gozon.ParseProxy("http://[IP | domain]:[port]", "username", "password")
        if err != nil {
            log.Println("test:1 -> err: ", err)
            return
        }
        //You need to place the country code, otherwise it will amazon will get it from the IP
        client := gozon.NewClient("MXN", proxyURL, 6)
        productsURLs, _, err := client.GetMainURLs()
        data, _, err := client.GetFromURL(productURL)
        if err != nil {
            log.Println("test:2 -> err: ", err)
            continue
        }
        if err := data.SetImages(client.ProxyURL); err != nil {
            log.Println("test:3 -> err: ", err)
            return
        }
        for j, imgs := range data.Images {
            if len(imgs.Large.Content) != 0 {
                f_name1 := fmt.Sprintf("./test/images/%d_large%s", j, imgs.Large.Extension)
                os.WriteFile(f_name1, imgs.Large.Content, 0644)
            }
            if len(imgs.Thumb.Content) != 0 {
                f_name2 := fmt.Sprintf("./test/images/%d_thumb%s", j, imgs.Thumb.Extension)
                os.WriteFile(f_name2, imgs.Thumb.Content, 0644)
            }
            if len(imgs.HiRes.Content) != 0 {
                f_name3 := fmt.Sprintf("./test/images/%d_hires%s", j, imgs.HiRes.Extension)
                os.WriteFile(f_name3, imgs.HiRes.Content, 0644)
            }
        }
        f, err := os.Create("./test/data.json")
        if err != nil {
            log.Println("test:4 -> err: ", err)
            return
        }
        json.NewEncoder(f).Encode(data)
    }
```