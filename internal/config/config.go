package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type Config struct {
	DbHost     string `toml:"db_host"`
	DbPort     string `toml:"db_port"`
	DbUser     string `toml:"db_user"`
	DbPassword string `toml:"db_password"`
	DbName     string `toml:"db_name"`

	//Logging
	LogLevel string `toml:"log_level"`

	HTTPAddr string `toml:"bind_addr"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}
	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	_, err := toml.DecodeFile(configPath, &cfg)
	if err != nil {
		log.Panic(err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
