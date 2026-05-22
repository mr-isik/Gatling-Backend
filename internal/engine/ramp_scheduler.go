package engine

import (
	"math"
	"time"

	"github.com/mr-isik/gatling-backend/internal/domain"
)

type RampScheduler struct {
	config domain.RunConfig
}

func NewRampScheduler(config domain.RunConfig) *RampScheduler {
	return &RampScheduler{config: config}
}

func (s *RampScheduler) CalculateVUs(elapsed time.Duration) int {
	totalVUs := float64(s.config.VUs)

	rampUp := s.config.RampUpDuration
	steady := s.config.Duration
	rampDown := s.config.RampDownDuration

	if elapsed < 0 {
		return 0
	}

	if elapsed <= rampUp {
		if rampUp == 0 {
			return int(totalVUs)
		}
		// Linear ramp up from 1 to VUs
		// fraction = elapsed / rampUp
		fraction := float64(elapsed) / float64(rampUp)
		vus := 1.0 + (totalVUs-1.0)*fraction
		return int(math.Round(vus))
	}

	if elapsed <= rampUp+steady {
		return int(totalVUs)
	}

	if elapsed <= rampUp+steady+rampDown {
		if rampDown == 0 {
			return 0
		}
		// Linear ramp down from VUs to 0
		remaining := (rampUp + steady + rampDown) - elapsed
		fraction := float64(remaining) / float64(rampDown)
		vus := totalVUs * fraction
		return int(math.Round(vus))
	}

	return 0
}
