package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DB DBConfig
	Twitch TwitchConfig
}

func YAML() (config *Config, err error) {
	return YAMLfromFile("./config.yml")
}

func YAMLfromFile(file string) (config *Config, err error) {
	f, err := os.Open(file)
	if err != nil {
		return &Config{}, err
	}
	err = yaml.NewDecoder(f).Decode(&config)
	return config, err
}
func (config *Config) SprintfDBConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.Database)
}
