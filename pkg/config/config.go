package config

import (
	"github.com/spf13/viper"
	"log"
	"strings"
	"sync"
	"time"
)

// AppConfig holds configuration settings for the auto-scaler application.
type AppConfig struct {
	TargetCPU              float64       `mapstructure:"targetCPU"`              // TargetCPU is the target CPU utilization percentage (e.g., 0.80 for 80%).
	StatusEndpoint         string        `mapstructure:"statusEndpoint"`         // StatusEndpoint is the URL endpoint for fetching current application status.
	ReplicasEndpoint       string        `mapstructure:"replicasEndpoint"`       // ReplicasEndpoint is the URL endpoint for updating the number of application replicas.
	PollInterval           time.Duration `mapstructure:"pollInterval"`           // PollInterval is the interval between consecutive status checks.
	CoolDownPeriod         time.Duration `mapstructure:"coolDownPeriod"`         // CoolDownPeriod is the duration to wait after scaling replicas before making another scaling decision.
	ApiTimeOut             time.Duration `mapstructure:"apiTimeOut"`             // ApiTimeOut is the timeout duration for API requests.
	LogLevel               string        `mapstructure:"logLevel"`               // LogLevel sets the logging level for the application (e.g., "debug", "info", "warn").
	DownscaleAfterAttempts int           `mapstructure:"downscaleAfterAttempts"` // DownscaleAfterAttempts specifies the number of retry attempts after which the system initiates downscaling of resources.
}

var C *AppConfig
var once sync.Once
var err error

// createConfig initializes configuration settings using environment variables and defaults.
func createConfig() error {
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	viper.AutomaticEnv()
	// Set default values for configuration parameters.
	viper.SetDefault("targetCPU", 0.80)
	viper.SetDefault("statusEndpoint", "http://localhost:8123/app/status")
	viper.SetDefault("replicasEndpoint", "http://localhost:8123/app/replicas")
	viper.SetDefault("pollInterval", 10*time.Second)
	viper.SetDefault("coolDownPeriod", 20*time.Second)
	viper.SetDefault("logLevel", "DEBUG")
	viper.SetDefault("apiTimeOut", 2*time.Second)
	viper.SetDefault("downscaleAfterAttempts", 3)
	// Unmarshal the configuration into the global AppConfig variable C.
	if err := viper.Unmarshal(&C); err != nil {
		log.Fatal("Error unmarshalling config:", err)
		return err
	}

	return nil
}

// GetConfig retrieves the application configuration, initializing it if necessary.
func GetConfig() *AppConfig {
	once.Do(func() {
		err = createConfig()
	})
	if err != nil {
		log.Fatal("Server Shutdown: unable to read config", err)
	}
	return C
}
