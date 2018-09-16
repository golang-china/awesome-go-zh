// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	for i := 0; i < 10; i++ {
		ch := make(chan string, 32) // HL

		go func() {
			if _, err := searchByBaidu("golang"); err == nil {
				ch <- "baidu" // HL
			}
		}()
		go func() {
			if _, err := searchByBing("golang"); err == nil {
				ch <- "bing" // HL
			}
		}()

		fmt.Println(<-ch) // HL
	}
}

func searchBySoso(key string) (string, error) {
	return httpGet(fmt.Sprintf("http://www.sogou.com/web?query=%s", url.QueryEscape(key)))
}

func searchByBing(key string) (string, error) {
	return httpGet(fmt.Sprintf("http://www.bing.com/search?q=%s", url.QueryEscape(key)))
}

func searchByBaidu(key string) (string, error) {
	return httpGet(fmt.Sprintf("http://www.baidu.com/s?wd=%s", url.QueryEscape(key)))
}

func httpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
