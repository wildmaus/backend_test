package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Address string `yaml:"address"`
		Port    string `yaml:"port"`
	} `yaml:"server"`
	Storage StorageConfig `yaml:"storage"`
}

type StorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var instanse *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instanse = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instanse); err != nil {
			help, _ := cleanenv.GetDescription(instanse, nil)
			log.Println(help)
			log.Fatal(err)
		}
	})
	return instanse
}
