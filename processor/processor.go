package processor

import (
	"context"
	"fmt"
	"sort"
	"time"

	"carbon_intensity/models"
)

var defaultResponseSlot = 30 * time.Minute

func GetContinuousSlots(_ context.Context, requiredDuration time.Duration, data []models.CarbonIntensityPeriod) ([]models.CarbonIntensityPeriod, error) {
	response := make([]models.CarbonIntensityPeriod, 0)
	requiredSlots := int(requiredDuration / defaultResponseSlot)

	if requiredSlots < 1 {
		return response, fmt.Errorf("needed duration must be at least %v", defaultResponseSlot)
	}

	if len(data) < requiredSlots {
		return response, fmt.Errorf("not enough data for %d slots, got %d", requiredSlots, len(data))
	}

	sum := 0
	for i := 0; i < requiredSlots; i++ {
		sum += data[i].Intensity.Forecast
	}

	average := sum / requiredSlots
	minIndex, maxIndex := 0, requiredSlots-1

	for i := requiredSlots; i < len(data); i++ {
		sum -= data[i-requiredSlots].Intensity.Forecast
		sum += data[i].Intensity.Forecast
		movingAverage := sum / requiredSlots

		if movingAverage < average {
			maxIndex = i
			minIndex = i - requiredSlots + 1
		}
	}

	for i := minIndex; i <= maxIndex; i++ {
		response = append(response, data[i])
	}

	return response, nil
}

func GetNonContinuousSlots(_ context.Context, requiredDuration time.Duration, data []models.CarbonIntensityPeriod) ([]models.CarbonIntensityPeriod, error) {
	response := make([]models.CarbonIntensityPeriod, 0)
	requiredSlots := int(requiredDuration / defaultResponseSlot)

	if requiredSlots < 1 {
		return response, fmt.Errorf("needed duration must be at least %v", defaultResponseSlot)
	}

	if len(data) < requiredSlots {
		return response, fmt.Errorf("not enough data for %d slots, got %d", requiredSlots, len(data))
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Intensity.Forecast < data[j].Intensity.Forecast
	})

	for i := 0; i < requiredSlots; i++ {
		response = append(response, data[i])
	}

	return response, nil
}
