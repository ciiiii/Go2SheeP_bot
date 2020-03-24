package config

import (
	"os"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Cos struct {
		Bucket    string `env:"COS_BUCKET"`
		Region    string `env:"COS_REGION"`
		SecretId  string `env:"COS_SECRET_ID"`
		SecretKey string `env:"COS_SECRET_KEY"`
	}
	Translate struct {
		AppId string `env:"TRANSLATE_APP_ID"`
		Key   string `env:"TRANSLATE_KEY"`
	}
	Bot struct {
		Token     string `env:"TOKEN"`
		PublicUrl string `env:"PUBLIC_URL"`
	}
	Deploy struct {
		Port string `env:"PORT"`
	}
}

var (
	c    Config
	once sync.Once
)

func Parser() Config {
	once.Do(func() {
		if os.Getenv("PORT") != "" {
			if err := env.Parse(&c); err != nil {
				panic(err)
			}
		} else {
			rootPath, _ := os.Getwd()
			confPath := strings.Join([]string{rootPath, "conf.toml"}, "/")
			if _, err := toml.DecodeFile(confPath, &c); err != nil {
				panic(err)
			}
		}
	})
	return c
}
