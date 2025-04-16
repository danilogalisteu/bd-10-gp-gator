package config


import (
    "encoding/json"
	"os"
	"path/filepath"
)


const configFileName = ".gatorconfig.json"


type Config struct {
	DbUrl string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}


func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configFileName), nil
}


func Read() (Config, error) {
	cfg := Config{}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return cfg, err
	}

    if err = json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
    }

    return cfg, nil
}


func write(cfg Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

    data, err := json.Marshal(cfg)
    if err != nil {
		return  err
    }

	err = os.WriteFile(configFilePath, data, 0666)
    if err != nil {
		return  err
    }

	return nil
}


func SetUser(user string) error {
	cfg, err := Read()
	if err != nil {
		return err
	}

	cfg.CurrentUserName = user

	err = write(cfg)
	if err != nil {
		return err
	}

	return nil
}
