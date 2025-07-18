package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configname string = ".gatorconfig.json"

func Read() Config {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("no user directory found")
	}

	configFile, err := os.ReadFile(dir + "/" + configname)
	if err != nil {
		log.Fatalf("no config found at %s", dir+"/"+configname)
	}

	var ret Config
	err = json.Unmarshal(configFile, &ret)
	if err != nil {
		log.Fatalf("could not unmarshal config file at %s", dir+"/"+configname)
	}

	return ret
}

func write(c Config) error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(dir+"/"+configname, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	if err := write(*c); err != nil {
		return err
	}
	return nil
}
