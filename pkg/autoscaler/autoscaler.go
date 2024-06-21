package autoscaler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/suyog1pathak/autoscaler/api/v1/model"
	c "github.com/suyog1pathak/autoscaler/pkg/config"
	"github.com/suyog1pathak/autoscaler/pkg/rest"
	"github.com/suyog1pathak/autoscaler/pkg/util"
	"log/slog"
	"math"
	"net/http"
	"os"
	"time"
)

var log *slog.Logger          // Logger instance for application events.
var config *c.AppConfig       // AppConfig instance holding configuration settings.
var downscaleReqCount int = 0 // Counter for consecutive downscale attempts.

// Start initializes the auto-scaler and begins monitoring and adjusting replicas.
func Start() {
	// Initialize logger with JSON format.
	log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// Retrieve application configuration.
	config = c.GetConfig()
	// Adjust logger with configured log level.
	log = slog.New(slog.NewJSONHandler(os.Stdout, util.LogLevelHandler(config.LogLevel)))
	// Log startup information.
	log.Info("Auto-scaler started", "upstream", config.StatusEndpoint, "targetCpu", config.TargetCPU, "pollingInterval", config.PollInterval.Seconds(), "logLevel", config.LogLevel)
	// Start main loop for monitoring and adjusting replicas.
	for {
		// Fetch current application status.
		status, err := getStatus()
		if err != nil {
			log.Error("Failed to fetch status", "error", err.Error())
			time.Sleep(config.PollInterval)
			continue
		}

		// Calculate new number of replicas.
		newReplicas := calculateNewReplicas(status.CPU.HighPriority, status.Replicas, config.DownscaleAfterAttempts)

		if newReplicas != status.Replicas {
			err = updateReplicas(newReplicas)
			// Log successful replica update.
			log.Debug("Replica count updated, waiting for cooldown period")
			time.Sleep(config.CoolDownPeriod)
			if err != nil {
				// Log error updating replicas.
				log.Error("Failed to update replicas", "error", err.Error())
			}

		}
		// Log skip replica update when recommendation matches current replicas.
		log.Debug("No change in recommended replicas")

		// Wait for the next polling interval.
		log.Debug("Waiting for next polling interval")
		time.Sleep(config.PollInterval)
	}
}

// getStatus retrieves the current application status from the upstream service.
func getStatus() (*model.StatusResponse, error) {
	// Make HTTP GET request to fetch status.
	r, responseBody, err := rest.Client(http.MethodGet, config.StatusEndpoint, map[string]string{"Accept": "application/json"}, []byte(""), config.ApiTimeOut)
	if err != nil {
		return nil, err
	}
	// Handle non-OK status code from upstream service.
	if r.StatusCode != http.StatusOK {
		log.Warn("Failed to fetch status from upstream service", "retryInterval", config.PollInterval.Seconds(), "statusCode", r.StatusCode)
		return nil, errors.New(fmt.Sprintf("Unexpected status code %d", r.StatusCode))
	}

	var status model.StatusResponse
	// Unmarshal JSON response into status model.
	log.Info("Current status retrieved", "response", util.ToJson(responseBody))
	err = json.Unmarshal(responseBody, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// calculateNewReplicas computes the new number of replicas based on current CPU utilization.
func calculateNewReplicas(currentCPU float64, currentReplicas int, downscaleAfterAttempts int) int {
	if currentCPU > config.TargetCPU {
		// Log adding replicas due to high CPU utilization.
		log.Debug("High CPU utilization detected, increasing replicas")
		return int(math.Ceil(float64(currentReplicas) * (currentCPU / config.TargetCPU)))
	} else {
		// Handle downscaled recommendation based on downscale attempts.
		if downscaleReqCount == downscaleAfterAttempts {
			log.Debug("Received downscaled recommendation")
			downscaleReqCount = 0
			return int(math.Floor(float64(currentReplicas) * (currentCPU / config.TargetCPU)))
		}
		// Log received downscale request.
		log.Debug("Received downscale request", "downscaleReqCount", downscaleReqCount)
		downscaleReqCount++
		return currentReplicas
	}
}

// updateReplicas sends a request to update the number of replicas to the upstream service.
func updateReplicas(newReplicas int) error {
	// Log updating new replica count.
	log.Info("Updating replica count", "count", newReplicas)
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	requestBody, err := json.Marshal(model.ReplicasRequest{Replicas: newReplicas})
	if err != nil {
		return err
	}

	r, _, err := rest.Client(http.MethodPut, config.ReplicasEndpoint, headers, requestBody, config.ApiTimeOut)

	if err != nil {
		// Log error updating replicas.
		log.Warn("Failed to update replica count on upstream service", "retryInterval", config.PollInterval.Seconds(), "error", err.Error())
		return err
	}

	if r.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("Unexpected status code %d", r.StatusCode))
	}

	return nil
}
