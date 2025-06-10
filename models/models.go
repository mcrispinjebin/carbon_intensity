package models

type Intensity struct {
	Forecast int    `json:"forecast"`
	Actual   int    `json:"actual"`
	Index    string `json:"index"`
}

type CarbonIntensityPeriod struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Intensity Intensity `json:"intensity"`
}

type ExternalAPIResponse struct {
	Data []CarbonIntensityPeriod `json:"data"`
}

type Response struct {
	Slots           []CarbonIntensityPeriod `json:"slots"`
	AverageForecast int                     `json:"average_forecast"`
}
