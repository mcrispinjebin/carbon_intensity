package processor

import (
	"testing"
	"time"

	"carbon_intensity/models"
)

func TestGetContinuousSlots(t *testing.T) {
	type args struct {
		requiredDuration time.Duration
		data             []models.CarbonIntensityPeriod
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected []int
	}{
		{
			name: "valid 1 hour slot",
			args: args{
				requiredDuration: 1 * time.Hour,
				data: []models.CarbonIntensityPeriod{
					{Intensity: models.Intensity{Forecast: 10}},
					{Intensity: models.Intensity{Forecast: 20}},
					{Intensity: models.Intensity{Forecast: 30}},
					{Intensity: models.Intensity{Forecast: 40}},
				},
			},
			wantErr:  false,
			expected: []int{10, 20},
		},
		{
			name: "valid 1 hour slot at the end",
			args: args{
				requiredDuration: 1 * time.Hour,
				data: []models.CarbonIntensityPeriod{
					{Intensity: models.Intensity{Forecast: 10}},
					{Intensity: models.Intensity{Forecast: 20}},
					{Intensity: models.Intensity{Forecast: 10}},
					{Intensity: models.Intensity{Forecast: 10}},
				},
			},
			wantErr:  false,
			expected: []int{10, 10},
		},
		{
			name: "valid 30 min slot at the end",
			args: args{
				requiredDuration: 30 * time.Minute,
				data: []models.CarbonIntensityPeriod{
					{Intensity: models.Intensity{Forecast: 10}},
					{Intensity: models.Intensity{Forecast: 20}},
					{Intensity: models.Intensity{Forecast: 10}},
					{Intensity: models.Intensity{Forecast: 5}},
				},
			},
			wantErr:  false,
			expected: []int{5},
		},
		{
			name: "duration less than slot",
			args: args{
				requiredDuration: 10 * time.Minute,
				data: []models.CarbonIntensityPeriod{
					{Intensity: models.Intensity{Forecast: 10}},
				},
			},
			wantErr:  true,
			expected: nil,
		},
		{
			name: "not enough data",
			args: args{
				requiredDuration: 2 * time.Hour,
				data: []models.CarbonIntensityPeriod{
					{Intensity: models.Intensity{Forecast: 10}},
				},
			},
			wantErr:  true,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetContinuousSlots(
				nil,
				tt.args.requiredDuration,
				tt.args.data,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.expected) {
					t.Errorf("got len = %v, want %v", len(got), len(tt.expected))
				}
				for i := range tt.expected {
					if got[i].Intensity.Forecast != tt.expected[i] {
						t.Errorf("got[%d] = %v, want %v", i, got[i].Intensity.Forecast, tt.expected[i])
					}
				}
			}
		})
	}
}

func TestGetNonContinuousSlots(t *testing.T) {
	type args struct {
		requiredDuration time.Duration
		data             []models.CarbonIntensityPeriod
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected []int
	}{
		{
			name: "valid 1 hour slot",
			args: args{
				requiredDuration: 1 * time.Hour,
				data: []models.CarbonIntensityPeriod{
					{Intensity: models.Intensity{Forecast: 30}},
					{Intensity: models.Intensity{Forecast: 10}},
					{Intensity: models.Intensity{Forecast: 20}},
					{Intensity: models.Intensity{Forecast: 40}},
				},
			},
			wantErr:  false,
			expected: []int{10, 20},
		},
		{
			name: "not enough data",
			args: args{
				requiredDuration: 2 * time.Hour,
				data: []models.CarbonIntensityPeriod{
					{Intensity: models.Intensity{Forecast: 10}},
				},
			},
			wantErr:  true,
			expected: nil,
		},
		{
			name: "duration less than slot",
			args: args{
				requiredDuration: 10 * time.Minute,
				data: []models.CarbonIntensityPeriod{
					{Intensity: models.Intensity{Forecast: 10}},
				},
			},
			wantErr:  true,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNonContinuousSlots(
				nil,
				tt.args.requiredDuration,
				tt.args.data,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.expected) {
					t.Errorf("got len = %v, want %v", len(got), len(tt.expected))
				}
				for i := range tt.expected {
					if got[i].Intensity.Forecast != tt.expected[i] {
						t.Errorf("got[%d] = %v, want %v", i, got[i].Intensity.Forecast, tt.expected[i])
					}
				}
			}
		})
	}
}
