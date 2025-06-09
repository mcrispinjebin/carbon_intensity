package models

import "time"

type Intensity struct {
	Forecast int    `json:"forecast"`
	Actual   int    `json:"actual"`
	Index    string `json:"index"`
}

type CarbonIntensityPeriod struct {
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
	Intensity Intensity `json:"intensity"`
}
