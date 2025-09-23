package initialization

import (
	"arch/internal/domain/entity"
	"encoding/json"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const (
	configFilePath = "configs/config.json"
	envFilePath    = `.env`
)

var (
	ConfigService = &entity.ConfigService{}
)

func loadConfig() error {
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(file, &ConfigService); err != nil {
		return err
	}

	return nil
}

func loadEnvironment() error {
	if err := godotenv.Load(envFilePath); err != nil {
		logrus.Warning("load file not found, environment variables load from environment")
	}
	if err := env.Parse(ConfigService); err != nil {
		return err
	}
	return nil
}
