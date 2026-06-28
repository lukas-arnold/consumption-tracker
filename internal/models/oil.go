package models

type OilInput struct {
	Date     string  `json:"date"`
	Volume   float64 `json:"volume"`
	Costs    float64 `json:"costs"`
	Retailer string  `json:"retailer"`
	Note     string  `json:"note"`
}

type Oil struct {
	Id int64 `json:"id"`
	OilInput
}

const OilTankHeightCm = 150.0

type OilFillLevelInput struct {
	Date  string  `json:"date"`
	Level float64 `json:"level"`
}

type OilFillLevel struct {
	Id int64 `json:"id"`
	OilFillLevelInput
	Percentage float64
}

type OilCharts struct {
	Consumption ChartModel
	Price       ChartModel
	FillLevel   ChartModel
}
