package models

type WaterInput struct {
	Year             int     `json:"year"`
	VolumeWater      float64 `json:"volumeWater"`
	VolumeWastewater float64 `json:"volumeWastewater"`
	VolumeRainwater  float64 `json:"volumeRainwater"`
	CostsWater       float64 `json:"costsWater"`
	CostsWastewater  float64 `json:"costsWastewater"`
	CostsRainwater   float64 `json:"costsRainwater"`
	Payments         float64 `json:"payments"`
	FixedPrice       float64 `json:"fixedPrice"`
	Note             string  `json:"note"`
}

type Water struct {
	Id int64 `json:"id"`
	WaterInput
}

type WaterCharts struct {
	Consumption ChartModel
	Price       ChartModel
}
