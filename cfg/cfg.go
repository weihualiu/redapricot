package cfg

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/toolkits/file"
)

type SocketConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
	Timeout int    `json:"timeout"`
}

type GlobalConfig struct {
	Debug         bool              `json:"debug"`
	Hostname      string            `json:"hostname"`
	Socket        *SocketConfig     `json:"socket"`
	DefaultTags   map[string]string `json:"default_tags"`
	IgnoreMetrics map[string]bool   `json:"ignore"`
	Database      string            `json:"database"`
	MaxConns  int         `json:"maxConns"`
	MaxIdle   int         `json:"maxIdle"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func Hostname() (string, error) {
	hostname := Config().Hostname
	if hostname != "" {
		return hostname, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ERROR: os.Hostname() fail", err)
	}
	return hostname, err
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
}
