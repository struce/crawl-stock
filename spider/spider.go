package spider

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)
import "encoding/json"
import "crawl/models"

func Post(url string, data *models.RequestBody, headers map[string]string) (content []byte, err error) {
	if url == "" {
		return nil, errors.New("url is null")
	}

	if data == nil {
		return nil, errors.New("data is null")
	}

	requestBody, err := json.Marshal(data)
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	if err != nil {
		return
	}

	for k, v := range headers{
		req.Header.Set(k, v)
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if req.Header.Get("Content-Length") != "" {
		req.Header.Del("Content-Length")
	}

	req.Header.Set("Content-Length", string(len(string(requestBody))))

	response, err := client.Do(req)
	if err != nil {
		return
	}

	defer response.Body.Close()

	if content, err = ioutil.ReadAll(response.Body); err != nil {
		return
	}
	return
}

func registered(rep *models.ResponseBody) (registered bool) {
	if rep.DESCRIPTION != "" {
		return true
	}

	if rep.PHOTO != "" {
		return true
	}

	if rep.TIME[4] == '/' {
		return true
	}

	return false
}


func IsRegisterPhone(url, phone string) (reged bool, err error) {
	httpHeaders := models.NewHttpHeader().HEADERS
	requestbody := models.NewRequestBody(phone)
	response, err := Post(url, requestbody, httpHeaders)
	if err != nil {
		return
	}

	respBody := models.NewResponseBody()
	err = json.Unmarshal(response, respBody)
	if err != nil {
		return
	}

	reged = registered(respBody)
	return
}