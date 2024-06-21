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

var log *slog.Logger          // Logger instance for logging application events.
var config *c.AppConfig       // AppConfig instance holding configuration settings.
var downscaleReqCount int = 0 // Counter for tracking consecutive downscale requests.

// Start initializes the auto-scaler and begins monitoring and adjusting replicas.
func Start() {
	// Initialize logger with JSON format for structured logging.
	log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// Retrieve application configuration.
	config = c.GetConfig()
	// Adjust logger with configured log level.
	log = slog.New(slog.NewJSONHandler(os.Stdout, util.LogLevelHandler(config.LogLevel)))
	log.Info("starting autoscaler", "upstream", config.StatusEndpoint, "targetCpu", config.TargetCPU, "polling", config.PollInterval.Seconds(), "loglevel", config.LogLevel)
	// Start main loop for monitoring and adjusting replicas.
	for {
		// Fetch the current status
		status, err := getStatus()
		if err != nil {
			log.Error("error fetching status", "error", err.Error())
			time.Sleep(config.PollInterval)
			continue
		}

		// Calculate new number of replicas
		newReplicas := calculateNewReplicas(status.CPU.HighPriority, status.Replicas)

		if newReplicas != status.Replicas {
			err = updateReplicas(newReplicas)
			log.Debug("updated the number of replicas, waiting for cool down period")
			time.Sleep(config.CoolDownPeriod)
			if err != nil {
				log.Error("error updating replicas", "error", err.Error())
			}

		}
		log.Debug("recommendation is equal to the current replica, skipping replica update")

		// Wait for the next poll interval
		log.Debug("waiting for polling interval")
		time.Sleep(config.PollInterval)
	}
}

// getStatus retrieves the current status of the application from the upstream service.
func getStatus() (*model.StatusResponse, error) {
	//log.Info("getting status for current pods and cpu")
	r, responseBody, err := rest.Client(http.MethodGet, config.StatusEndpoint, map[string]string{"Accept": "application/json"}, []byte(""), config.ApiTimeOut)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		log.Warn("unable to check status from upstream service, will poll after interval", "interval", config.PollInterval.Seconds(), "response_code", r.StatusCode)
		return nil, errors.New(fmt.Sprintf("unexpected status code %d", r.StatusCode))
	}

	var status model.StatusResponse
	log.Info("current status", "response", util.ToJson(responseBody))
	err = json.Unmarshal(responseBody, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// calculateNewReplicas computes the new number of replicas based on current CPU utilization.
func calculateNewReplicas(currentCPU float64, currentReplicas int) int {
	if currentCPU > config.TargetCPU {
		log.Debug("current cpu utilization is more then threshold, adding more replicas")
		return int(math.Ceil(float64(currentReplicas) * (currentCPU / config.TargetCPU)))
	} else {
		//log.Debug("current cpu utilization is below threshold, reducing replicas")
		if downscaleReqCount == 3 {
			log.Debug("received downscaling recommendation")
			downscaleReqCount = 0
			return int(math.Floor(float64(currentReplicas) * (currentCPU / config.TargetCPU)))
		}
		log.Debug("got downscale request as", "downscaleReqCount", downscaleReqCount)
		downscaleReqCount++
		return currentReplicas
	}
}

// updateReplicas sends a request to update the number of replicas to the upstream service.
func updateReplicas(newReplicas int) error {
	log.Info("updating new replica count", "count", newReplicas)
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	requestBody, err := json.Marshal(model.ReplicasRequest{Replicas: newReplicas})
	if err != nil {
		return err
	}

	r, _, err := rest.Client(http.MethodPut, config.ReplicasEndpoint, headers, requestBody, config.ApiTimeOut)

	if err != nil {
		log.Warn("unable to update new replica count on upstream service, will try after polling interval", "interval", config.PollInterval.Seconds(), "err", err.Error())
		return err
	}

	if r.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("unexpected status code %d", r.StatusCode))
	}

	return nil
}
