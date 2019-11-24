package storage

import (
	//"fmt"
	"sync"
	"time"
	sj "github.com/bitly/go-simplejson"
	"github.com/go-clog/clog"
	"github.com/lwl1989/crawl-stock/models"
	"github.com/parnurzeal/gorequest"
)

func init()  {
    clog.NewConsole()
}
// CheckProxy .
func CheckProxy(ip *models.IP) {
	if CheckIP(ip) {
		ProxyAdd(ip)
	}
}

// CheckIP is to check the ip work or not
func CheckIP(ip *models.IP) bool {
	var pollURL string
	var testIP string
	if ip.Type2 == "https" {
		testIP = "https://" + ip.Data
		pollURL = "https://httpbin.org/get"
	} else {
		testIP = "http://" + ip.Data
		pollURL = "http://httpbin.org/get"
	}
	clog.Info(testIP)
	begin := time.Now()
	resp, _, errs := gorequest.New().Proxy(testIP).Get(pollURL).End()
	if errs != nil {
		clog.Warn("[CheckIP] testIP = %s, pollURL = %s: Error = %v", testIP, pollURL, errs)
		return false
	}
	
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		//harrybi 20180815 判断返回的数据格式合法性
		_, err := sj.NewFromReader(resp.Body)
		if err != nil {
			clog.Warn("[CheckIP] testIP = %s, pollURL = %s: Error = %v", testIP, pollURL, err)
			return false
		}
		//harrybi 计算该代理的速度，单位毫秒
		ip.Speed = time.Now().Sub(begin).Nanoseconds() / 1000 / 1000 //ms
		if err = models.UpdateIp(ip); err != nil {
			clog.Warn("[CheckIP] Update IP = %v Error = %v", *ip, err)
		}

		return true
	}
	return false
}

// CheckProxyDB to check the ip in DB
func CheckProxyDB() {
	x := models.CountIPs()
	clog.Info("Before check, DB has: %d records.", x)
	ips, err := models.GetAll()
	if err != nil {
		clog.Warn(err.Error())
		return
	}
	var wg sync.WaitGroup
	for _, v := range ips {
		wg.Add(1)
		go func(v *models.IP) {
			if !CheckIP(v) {
				ProxyDel(v)
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
	x = models.CountIPs()
	clog.Info("After check, DB has: %d records.", x)
}

// ProxyRandom .
func ProxyRandom() (ip *models.IP) {

	ips, err := models.GetAll()
	x := len(ips)
	clog.Warn("len(ips) = %d",x)
	if (err != nil || x == 0) {
		clog.Warn(err.Error())
		return models.NewIP()
	}
	randomNum := RandInt(0, x)

	return ips[randomNum]
}


// ProxyAdd .
func ProxyAdd(ip *models.IP) {
	models.InsertIp(ip)
}

// ProxyDel .
func ProxyDel(ip *models.IP) {
	models.DeleteIP(ip)
}

