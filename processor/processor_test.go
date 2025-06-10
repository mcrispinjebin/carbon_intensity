package processor

import (
	"testing"

	"carbon_intensity/models"
)

func TestGetContinuousSlots(t *testing.T) {
	type args struct {
		requiredDuration int
		data             []models.CarbonIntensityPeriod
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected int
	}{
		{
			name: "valid 1 hour slot",
			args: args{
				requiredDuration: 60,
				data: []models.CarbonIntensityPeriod{
					{From: "2006-01-02T12:00Z", Intensity: models.Intensity{Forecast: 60}},
					{From: "2006-01-02T12:30Z", Intensity: models.Intensity{Forecast: 80}},
					{From: "2006-01-02T13:00Z", Intensity: models.Intensity{Forecast: 90}},
					{From: "2006-01-02T13:30Z", Intensity: models.Intensity{Forecast: 80}},
				},
			},
			wantErr:  false,
			expected: 70,
		},
		{
			name: "partial 45 min slot",
			args: args{
				requiredDuration: 45,
				data: []models.CarbonIntensityPeriod{
					{From: "2006-01-02T12:00Z", Intensity: models.Intensity{Forecast: 60}},
					{From: "2006-01-02T12:30Z", Intensity: models.Intensity{Forecast: 80}},
					{From: "2006-01-02T13:00Z", Intensity: models.Intensity{Forecast: 90}},
					{From: "2006-01-02T13:30Z", Intensity: models.Intensity{Forecast: 80}},
				},
			},
			wantErr:  false,
			expected: 66,
		},
		{
			name: "valid 1 hour slot at the end",
			args: args{
				requiredDuration: 60,
				data: []models.CarbonIntensityPeriod{
					{From: "2006-01-02T12:00Z", Intensity: models.Intensity{Forecast: 80}},
					{From: "2006-01-02T12:30Z", Intensity: models.Intensity{Forecast: 90}},
					{From: "2006-01-02T13:00Z", Intensity: models.Intensity{Forecast: 60}},
					{From: "2006-01-02T13:30Z", Intensity: models.Intensity{Forecast: 80}},
				},
			},
			wantErr:  false,
			expected: 70,
		},
		{
			name: "valid 30 min slot at the end",
			args: args{
				requiredDuration: 30,
				data: []models.CarbonIntensityPeriod{
					{From: "2006-01-02T12:00Z", Intensity: models.Intensity{Forecast: 80}},
					{From: "2006-01-02T12:30Z", Intensity: models.Intensity{Forecast: 90}},
					{From: "2006-01-02T13:00Z", Intensity: models.Intensity{Forecast: 70}},
					{From: "2006-01-02T13:30Z", Intensity: models.Intensity{Forecast: 60}},
				},
			},
			wantErr:  false,
			expected: 60,
		},
		{
			name: "not enough data",
			args: args{
				requiredDuration: 60,
				data: []models.CarbonIntensityPeriod{
					{Intensity: models.Intensity{Forecast: 10}},
				},
			},
			wantErr:  true,
			expected: 0,
		},
	}

	p := &Processor{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.getContinuousSlots(
				nil,
				tt.args.requiredDuration,
				tt.args.data,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.AverageForecast != tt.expected {
					t.Errorf("got = %v, want %v", got.AverageForecast, tt.expected)
				}
			}
		})
	}
}

func TestGetNonContinuousSlots(t *testing.T) {
	type args struct {
		requiredDuration int
		data             []models.CarbonIntensityPeriod
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected int
	}{
		{
			name: "valid 1 hour slot",
			args: args{
				requiredDuration: 60,
				data: []models.CarbonIntensityPeriod{
					{From: "2006-01-02T12:00Z", Intensity: models.Intensity{Forecast: 60}},
					{From: "2006-01-02T12:30Z", Intensity: models.Intensity{Forecast: 80}},
					{From: "2006-01-02T13:00Z", Intensity: models.Intensity{Forecast: 90}},
					{From: "2006-01-02T13:30Z", Intensity: models.Intensity{Forecast: 60}},
				},
			},
			wantErr:  false,
			expected: 60,
		},
		{
			name: "not enough data",
			args: args{
				requiredDuration: 120,
				data: []models.CarbonIntensityPeriod{
					{Intensity: models.Intensity{Forecast: 10}},
				},
			},
			wantErr:  true,
			expected: 0,
		},
	}

	p := &Processor{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.getNonContinuousSlots(
				nil,
				tt.args.requiredDuration,
				tt.args.data,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.AverageForecast != tt.expected {
					t.Errorf("got = %v, want %v", got.AverageForecast, tt.expected)
				}
			}
		})
	}
}
