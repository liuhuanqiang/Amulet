package config

import (
	"encoding/json"
	"sync"
	"io/ioutil"
	"github.com/golang/glog"
)

type GlobalConfig struct {
	UserName	string 		`json:"user_name"`
	Password	string		`json:"password"`
	Port  		string		`json:"port"`
	Table       string		`json:"table"`
}

var (
	ConfigFile string
	cfg        *GlobalConfig
	mu         = &sync.RWMutex{}
)

func Config() *GlobalConfig {
	mu.RLock()
	defer mu.RUnlock()

	return cfg
}

func LoadConfig(f string) bool {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		glog.Error("read config file failed f:", f, " err:", err)
		return false
	}

	config := &GlobalConfig{}
	err = json.Unmarshal(b, config)
	if err != nil {
		glog.Error("unmarshal config file failed err:", err, " content:", string(b))
		return false
	}

	mu.Lock()
	defer mu.Unlock()

	ConfigFile = f
	cfg = config
	return true
}
