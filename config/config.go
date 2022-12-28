package config

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/SultanKs4/ytarchiver/pkg/fileutils"
	"github.com/spf13/viper"
)

// config struct
type Config struct {
	APIKey  Apikey `yaml:"apikey"`
	Youtube Youtube
	Save    string `yaml:"save"`
}

type Apikey struct {
	Youtube string `yaml:"youtube"`
	Drive   string `yaml:"drive"`
}

type Youtube struct {
	Audio    string `yaml:"audio"`
	Height   int    `yaml:"height"`
	Mimetype string `yaml:"mimetype"`
}

func NewConfig() *Config {
	// set default value
	return &Config{
		Youtube: Youtube{
			Audio:    "AUDIO_QUALITY_MEDIUM",
			Height:   720,
			Mimetype: "mp4",
		},
	}
}

// load from yml files to struct config
func (c *Config) LoadConfig() error {
	basePath := fileutils.BasePath()
	configPath := filepath.Join(basePath, "config")
	// main viper config
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// read the file
	if err := viper.ReadInConfig(); err != nil {
		err = fmt.Errorf("error reading config file: %s", err)
		return err
	}

	// map to app
	if err := viper.Unmarshal(c); err != nil {
		err = fmt.Errorf("unable to decode into struct: %s", err)
		return err
	}

	// done
	log.Printf("config loaded successfully")
	return nil
}
