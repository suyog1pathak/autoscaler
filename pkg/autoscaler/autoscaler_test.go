package autoscaler

import (
	c "github.com/suyog1pathak/autoscaler/pkg/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
)

func TestCalculateNewReplicas(t *testing.T) {
	// Mock logger initialization
	log = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Mock config values
	config = &c.AppConfig{
		TargetCPU:              0.8,
		DownscaleAfterAttempts: 3,
	}

	tests := []struct {
		currentCPU        float64
		currentReplicas   int
		downscaleReqCount int
		expectedReplicas  int
	}{
		{0.9, 10, 0, 12}, // Scale up
		{0.7, 10, 3, 8},  // Scale down
		{0.7, 10, 2, 10}, // Request to downscale but not yet reached the threshold
		{0.8, 10, 0, 10}, // CPU equals Target CPU, no change
	}

	for _, tt := range tests {
		downscaleReqCount = tt.downscaleReqCount
		actualReplicas := calculateNewReplicas(tt.currentCPU, tt.currentReplicas, config.DownscaleAfterAttempts)
		assert.Equal(t, tt.expectedReplicas, actualReplicas)
	}
}
