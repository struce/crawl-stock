package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/lwl1989/crawl-stock/getter"
    "sync"
    "github.com/lwl1989/crawl-stock/models"
    "time"
    "runtime"
    "github.com/lwl1989/crawl-stock/storage"
    "strconv"
    "fmt"
    "os"
)

// VERSION for this program


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
func Run() {

    AppAddr := "0.0.0.0"
    AppPort := "3001"
    mux := http.NewServeMux()
    mux.HandleFunc("/ip", ProxyHandler)
    mux.HandleFunc("/count", CountHandler)
    log.Println("Starting server", AppAddr+":"+AppPort)
    _ = http.ListenAndServe(AppAddr+":"+AppPort, mux)
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

func CountHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        w.Header().Set("content-type", "application/json")
        type count struct {
            Count int64
        }
        var c count
        c.Count = models.CountIPs()
        b, err := json.Marshal(c)
        if err != nil {
            return
        }
        w.Write(b)
    }
}


func main() {

    //init the database

    runtime.GOMAXPROCS(runtime.NumCPU())
    ipChan := make(chan *models.IP, 2000)
    PhonePrefixUsed = make([]string, 900)
    used = make(map[int]string)
    for i:=0;i<900 ; i++  {
       str := strconv.Itoa(i)
       if i < 10 {
           str = "00"+str
       }else if i < 100 {
           str = "0"+str
       }
        PhonePrefixUsed[i] = str
    }
    buildPrefix()
    os.Exit(1)
    // Start HTTP
    go func() {
        Run()
    }()


    // Check the IPs in channel
    for i := 0; i < 50; i++ {
        go func() {
            for {
                storage.CheckProxy(<-ipChan)
            }
        }()
    }
    // Start getters to scraper IP and put it in channel
    for {
        n := models.CountIPs()
        log.Printf("Chan: %v, IP: %v\n", len(ipChan), n)
        if len(ipChan) < 100 {
            go run(ipChan)
        }
        time.Sleep(10 * time.Minute)
    }
}


var PhonePrefixUsed []string
var used map[int]string
func buildPrefix()  {
    start:="09"
    for i:=0;i<10;i++ {
        use := storage.RandInt(0, 899)
        if v,ok := used[use]; ok {
            //fmt.Println(v,ok)
            if v == "" {
                buildMobile(start+PhonePrefixUsed[use])
                used[use] = PhonePrefixUsed[use]
            }
        }else{
            buildMobile(start+PhonePrefixUsed[use])
            used[use] = PhonePrefixUsed[use]
        }
    }
}
func buildMobile(prefix string)  {
    var start int64 = 0
    var end int64 = 99999

    phone:=""
    for i:=start; i<=end; i++ {
        suffix := strconv.FormatInt(i, 10)
        phone = prefix + suffix
    }
    fmt.Println(phone)
}