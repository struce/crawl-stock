package test

import (
	"fmt"
	"testing"
)
import "crawl/models"
import "crawl/spider"

func TestPost(t *testing.T) {
	headers := make(map[string]string)

	headers["User-Agent"] = "Mozilla/5.0 (Linux; Android 4.4.4; SM-G9350 Build/KTU84P) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/33.0.0.0 Mobile Safari/537.36"
	headers["Host"] = "405mtf.mitake.com.tw:8516"
	headers["Connection"] = "close"
	headers["Accept-Encoding"] = "gzip, deflate"
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	requestbody := models.NewRequestBody("999999999")

	url := "http://127.0.0.1/1.php"
	response, err := spider.Post(url, requestbody, headers)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(response))
}
