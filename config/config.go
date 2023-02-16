package config

import (
	"flag"
	"fmt"
	"github.com/jinzhu/configor"
)

var Config = struct {
	Redis struct {
		Adders []string
		Auth   string `default:"123456789"`
		Db     int    `required:"true" env:"db" default:"0"`
		Port   uint   `default:"7001"`
	}
	Rocketmq struct {
		Producer struct {
			Group string `json:"group"`
		} `json:"producer"`
		NameServer string `yaml:"name-server"`
	} `json:"rocketmq"`

	Symbols []string

	Mongodb struct {
		Database               string `json:"database"`
		Username               string `json:"username"`
		Password               string `json:"password"`
		Address                string `json:"address"`
		AuthenticationDatabase string `yaml:"authentication-database"`
	} `json:"mongodb"`
}{}

func init() {
	configor.Load(&Config, GetProfilesConf())
}

const (
	defaultConfigFileName = "config.yml"
)

func GetProfilesConf() string {
	var profiles string
	flag.StringVar(&profiles, "p", "", "运行环境(dev,test,prod)，默认无")
	flag.Parse()
	fmt.Println("-------------------------环境" + profiles)
	configFileNameFmt := "config_%s.yml"
	switch profiles {
	case "dev":
		return fmt.Sprintf(configFileNameFmt, "dev")
	case "prod":
		return fmt.Sprintf(configFileNameFmt, "prod")
	case "test":
		return fmt.Sprintf(configFileNameFmt, "test")
	default:
		return defaultConfigFileName
	}
}
