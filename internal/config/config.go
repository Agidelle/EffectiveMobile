package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strconv"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	AppPort    string `mapstructure:"APP_PORT"`
}

func LoadCfg() (*Config, error) {
	//Конфиг для разработки из .env
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println(".env not found")
		} else {
			fmt.Println(err)
		}
	}

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config not loaded, error: %v", err)
	}

	//Проверки конфига
	if cfg.DBHost == "" {
		return nil, fmt.Errorf("DB host not specified")
	}
	if cfg.DBName == "" {
		return nil, fmt.Errorf("DB name not specified")
	}
	if cfg.DBUser == "" {
		return nil, fmt.Errorf("DB user not specified")
	}
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB Password not specified")
	}

	port, err := strconv.Atoi(cfg.DBPort)
	if err != nil {
		return nil, fmt.Errorf("incorrect port db: %s", cfg.DBPort)
	}
	if port <= 0 || port > 65535 {
		return nil, fmt.Errorf("incorrect port db: %d", port)
	}

	return &cfg, nil
}
