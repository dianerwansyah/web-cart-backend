package helper

import (
	"io/ioutil"
	"log"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	config *Config
	once   sync.Once
)

type Config struct {
	Server struct {
		IamPort   string `yaml:"iam_port"`
		TrxPort   string `yaml:"trx_port"`
		SetupPort string `yaml:"setup_port"`
		JwtSecret string `yaml:"jwt_secret"`
		MongoURI  string `yaml:"mongo_uri"`
		MongoDB   string `yaml:"mongo_db"`
	} `yaml:"server"`
}

func GetConfig() *Config {
	once.Do(func() {
		config = &Config{}
		data, err := ioutil.ReadFile("config/config.yaml")
		if err != nil {
			log.Fatalf("Failed to read config file: %v", err)
		}
		err = yaml.Unmarshal(data, config)
		if err != nil {
			log.Fatalf("Failed to unmarshal config file: %v", err)
		}
	})
	return config
}
