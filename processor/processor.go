package processor

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"carbon_intensity/models"
)

type DataFetchAdapter interface {
	GetCarbonIntensityForecast(ctx context.Context, fromTime time.Time) ([]models.CarbonIntensityPeriod, error)
}

type Processor struct {
	DataClient DataFetchAdapter
}

var defaultResponseSlot = 30

func (p *Processor) getContinuousSlots(ctx context.Context, requiredDuration int, data []models.CarbonIntensityPeriod) (models.Response, error) {
	response := models.Response{}
	windowAverage := make([]int, 0)
	requiredSlots := int(math.Ceil(float64(requiredDuration) / float64(defaultResponseSlot)))

	if len(data) < requiredSlots {
		return response, fmt.Errorf("not enough data for %d slots, got %d", requiredSlots, len(data))
	}

	offset := 0
	minIndex := offset
	for i := 0; i < requiredSlots; i++ {
		windowAverage = append(windowAverage, data[i].Intensity.Forecast)
	}

	minAverage := p.getWeightedAverage(windowAverage, requiredDuration)

	for i := requiredSlots; i < len(data); i++ {
		offset += 1
		windowAverage = append(windowAverage, data[i].Intensity.Forecast)
		movingAverage := p.getWeightedAverage(windowAverage[offset:], requiredDuration)

		if movingAverage < minAverage {
			minIndex = offset
			minAverage = movingAverage
		}
	}

	slots := make([]models.CarbonIntensityPeriod, 0)
	for i := minIndex; i < minIndex+requiredSlots; i++ {
		slots = append(slots, models.CarbonIntensityPeriod{
			From:      data[i].From,
			To:        data[i].To,
			Intensity: models.Intensity{Forecast: data[i].Intensity.Forecast},
		})
	}

	response, err := p.transformResponse(ctx, slots, requiredDuration)
	if err != nil {
		return response, fmt.Errorf("failed to transform response: %w", err)
	}

	response.AverageForecast = minAverage
	return response, nil
}

func (p *Processor) getNonContinuousSlots(ctx context.Context, requiredDuration int, data []models.CarbonIntensityPeriod) (models.Response, error) {
	response := models.Response{}
	windowAverage := make([]int, 0)
	requiredSlots := int(math.Ceil(float64(requiredDuration) / float64(defaultResponseSlot)))

	if len(data) < requiredSlots {
		return response, fmt.Errorf("not enough data for %d slots, got %d", requiredSlots, len(data))
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Intensity.Forecast < data[j].Intensity.Forecast
	})

	for i := 0; i < requiredSlots; i++ {
		windowAverage = append(windowAverage, data[i].Intensity.Forecast)
	}

	averageForecast := p.getWeightedAverage(windowAverage, requiredDuration)

	slots := make([]models.CarbonIntensityPeriod, 0)
	for i := 0; i < requiredSlots; i++ {
		slots = append(slots, models.CarbonIntensityPeriod{
			From:      data[i].From,
			To:        data[i].To,
			Intensity: models.Intensity{Forecast: data[i].Intensity.Forecast},
		})
	}

	response, err := p.transformResponse(ctx, slots, requiredDuration)
	if err != nil {
		return response, fmt.Errorf("failed to transform response: %w", err)
	}

	response.AverageForecast = averageForecast
	return response, nil
}

func (p *Processor) GetSlots(ctx context.Context, requiredDuration int, isContinuous bool) (models.Response, error) {
	data, err := p.DataClient.GetCarbonIntensityForecast(ctx, time.Now())
	if err != nil {
		return models.Response{}, fmt.Errorf("failed to fetch carbon intensity data: %w", err)
	}

	if isContinuous {
		return p.getContinuousSlots(ctx, requiredDuration, data)
	}
	return p.getNonContinuousSlots(ctx, requiredDuration, data)
}

// assumption 1 - average returns int value
// assumption 2 - partial slot is always searched at the end of the data
func (p *Processor) getWeightedAverage(data []int, requiredDuration int) int {
	sum := 0
	for i := 0; i < len(data)-1; i++ {
		sum += data[i] * 30
	}

	remainingDuration := requiredDuration % defaultResponseSlot
	if remainingDuration == 0 {
		sum += data[len(data)-1] * defaultResponseSlot
		return sum / requiredDuration
	}

	sum += data[len(data)-1] * remainingDuration
	return sum / requiredDuration
}

func (p *Processor) transformResponse(_ context.Context, slots []models.CarbonIntensityPeriod, requiredDuration int) (models.Response, error) {
	response := models.Response{}
	timeLayout := "2006-01-02T15:04Z"
	isPartialSlot := (requiredDuration % defaultResponseSlot) != 0
	requiredSlots := int(math.Ceil(float64(requiredDuration) / float64(defaultResponseSlot)))
	toTime, err := time.Parse(timeLayout, slots[requiredSlots-1].From)
	if err != nil {
		return response, fmt.Errorf("failed to parse time: %w", err)
	}

	for i := 0; i < requiredSlots-1; i++ {
		response.Slots = append(response.Slots, slots[i])
	}

	if isPartialSlot {
		slot := models.CarbonIntensityPeriod{
			From:      slots[requiredSlots-1].From,
			To:        toTime.Add(time.Duration(requiredDuration%defaultResponseSlot) * time.Minute).Format(timeLayout),
			Intensity: models.Intensity{Forecast: slots[requiredSlots-1].Intensity.Forecast},
		}
		response.Slots = append(response.Slots, slot)
		return response, nil
	}

	response.Slots = append(response.Slots, models.CarbonIntensityPeriod{
		From:      slots[requiredSlots-1].From,
		To:        slots[requiredSlots-1].To,
		Intensity: models.Intensity{Forecast: slots[requiredSlots-1].Intensity.Forecast},
	})
	return response, nil
}
