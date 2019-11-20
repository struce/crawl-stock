
package models

import (
"fmt"
    "sync"
)

// IP struct
type IP struct {
    ID    int64  `xorm:"pk autoincr" json:"-"`
    Data  string `xorm:"NOT NULL" json:"ip"`
    Type1 string `xorm:"NOT NULL" json:"type1"`
    Type2 string `xorm:"NULL" json:"type2,omitempty"`
    Speed int64  `xorm:"NOT NULL" json:"speed,omitempty"`
}

var IpMap *sync.Map

func init()  {
    IpMap = new(sync.Map)
}
// NewIP .
func NewIP() *IP {
    //init the speed to 100 Sec
    return &IP{Speed: 100}
}


func countIps() int64 {
    // set id >= 0, fix bug: when this is nothing in the database
    count := int64(0)
    IpMap.Range(func(key, value interface{}) bool {
        count+=1
        return true
    })
    return count
}

// CountIPs .
func CountIPs() int64 {
    return countIps()
}

func deleteIP(ip *IP) error {
    IpMap.Delete(ip.Data)
    return nil
}

// DeleteIP .
func DeleteIP(ip *IP) error {
    return deleteIP(ip)
}

func getOne(ip string) *IP {
    tmpIp,ok := IpMap.Load(ip)
    if ok {
        return tmpIp.(*IP)
    }

    return NewIP()
}

// GetOne .
func GetOne(ip string) *IP {
    return getOne(ip)
}

func getAll() ([]*IP, error) {
    tmpIp := make([]*IP, 0)
    IpMap.Range(func(_, value interface{}) bool {
        tmpIp = append(tmpIp, value.(*IP))
        return true
    })
    return tmpIp, nil
}

// GetAll .
func GetAll() ([]*IP, error) {
    return getAll()
}
