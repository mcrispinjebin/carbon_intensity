package processor

import (
	"context"
	"fmt"
	"time"

	"carbon_intensity/models"
)

var defaultResponseSlot = 30 * time.Minute

func GetContinuousSlots(_ context.Context, requiredDuration time.Duration, data []models.CarbonIntensityPeriod) (models.CarbonIntensityPeriod, error) {
	requiredSlots := int(requiredDuration / defaultResponseSlot)

	if requiredSlots < 1 {
		return models.CarbonIntensityPeriod{}, fmt.Errorf("needed duration must be at least %v", defaultResponseSlot)
	}

	if len(data) < requiredSlots {
		return models.CarbonIntensityPeriod{}, fmt.Errorf("not enough data for %d slots, got %d", requiredSlots, len(data))
	}

	sum := 0
	for i := 0; i < requiredSlots; i++ {
		sum += data[i].Intensity.Forecast
	}

	average := sum / requiredSlots
	response := models.CarbonIntensityPeriod{
		From: data[0].From,
		To:   data[requiredSlots-1].To,
		Intensity: models.Intensity{
			Forecast: average,
		},
	}

	for i := requiredSlots; i < len(data); i++ {
		sum -= data[i-requiredSlots].Intensity.Forecast
		sum += data[i].Intensity.Forecast
		movingAverage := sum / requiredSlots

		if movingAverage < average {
			response = models.CarbonIntensityPeriod{
				From: data[i-requiredSlots].From,
				To:   data[i].To,
				Intensity: models.Intensity{
					Forecast: movingAverage,
				},
			}
		}
	}

	return response, nil
}
