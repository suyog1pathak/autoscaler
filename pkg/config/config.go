package config

import (
	"github.com/spf13/viper"
	"log"
	"strings"
	"sync"
	"time"
)

type AppConfig struct {
	TargetCPU        float64       `mapstructure:"targetCPU"`
	StatusEndpoint   string        `mapstructure:"statusEndpoint"`
	ReplicasEndpoint string        `mapstructure:"replicasEndpoint"`
	PollInterval     time.Duration `mapstructure:"pollInterval"`
	CoolDownPeriod   time.Duration `mapstructure:"coolDownPeriod"`
	ApiTimeOut       time.Duration `mapstructure:"apiTimeOut"`
	LogLevel         string        `mapstructure:"logLevel"`
}

var C *AppConfig
var once sync.Once
var err error

func CreateConfig() error {
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	viper.AutomaticEnv()
	viper.SetDefault("targetCPU", 0.80)
	viper.SetDefault("statusEndpoint", "http://localhost:8123/app/status")
	viper.SetDefault("replicasEndpoint", "http://localhost:8123/app/replicas")
	viper.SetDefault("pollInterval", 10*time.Second)
	viper.SetDefault("coolDownPeriod", 20*time.Second)
	viper.SetDefault("logLevel", "DEBUG")
	viper.SetDefault("apiTimeOut", 2*time.Second)

	if err := viper.Unmarshal(&C); err != nil {
		log.Fatal("Error unmarshalling config:", err)
		return err
	}

	return nil
}

func GetConfig() *AppConfig {
	once.Do(func() {
		err = CreateConfig()
	})
	if err != nil {
		log.Fatal("Server Shutdown: unable to read config", err)
	}
	return C
}
