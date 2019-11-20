package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/henson/proxypool/pkg/setting"
    "github.com/henson/proxypool/pkg/storage"
    "github.com/lwl1989/crawl-stock/getter"
    "sync"
    "github.com/lwl1989/crawl-stock/models"
)

// VERSION for this program
const VERSION = "/v2"


func run(ipChan chan<- *models.IP) {
    var wg sync.WaitGroup
    funs := []func() []*models.IP{
        //getter.Data5u,
        getter.Feiyi,
        getter.IP66, //need to remove it
        getter.KDL,
        //getter.GBJ,	//因为网站限制，无法正常下载数据
        //getter.Xici,
        //getter.XDL,
        //getter.IP181,  // 已经无法使用
        //getter.YDL,	//失效的采集脚本，用作系统容错实验
        getter.PLP,   //need to remove it
        getter.IP89,
    }
    for _, f := range funs {
        wg.Add(1)
        go func(f func() []*models.IP) {
            temp := f()
            //log.Println("[run] get into loop")
            for _, v := range temp {
                //log.Println("[run] len of ipChan %v",v)
                ipChan <- v
            }
            wg.Done()
        }(f)
    }
    wg.Wait()
    log.Println("All getters finished.")
}

// Run for request
func main() {

    mux := http.NewServeMux()
    mux.HandleFunc(VERSION+"/ip", ProxyHandler)
    mux.HandleFunc(VERSION+"/https", FindHandler)
    log.Println("Starting server", setting.AppAddr+":"+setting.AppPort)
    http.ListenAndServe(setting.AppAddr+":"+setting.AppPort, mux)
}

// ProxyHandler .
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        w.Header().Set("content-type", "application/json")
        b, err := json.Marshal(storage.ProxyRandom())
        if err != nil {
            return
        }
        w.Write(b)
    }
}

// FindHandler .
func FindHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        w.Header().Set("content-type", "application/json")
        b, err := json.Marshal(storage.ProxyFind("https"))
        if err != nil {
            return
        }
        w.Write(b)
    }
}
