package models

type WaterInput struct {
	Year             int     `json:"Year"`
	VolumeWater      float64 `json:"VolumeWater"`
	VolumeWastewater float64 `json:"VolumeWastewater"`
	VolumeRainwater  float64 `json:"VolumeRainwater"`
	CostsWater       float64 `json:"CostsWater"`
	CostsWastewater  float64 `json:"CostsWastewater"`
	CostsRainwater   float64 `json:"CostsRainwater"`
	Payments         float64 `json:"Payments"`
	FixedPrice       float64 `json:"FixedPrice"`
	Note             string  `json:"Note"`
}

type Water struct {
	Id int64 `json:"Id"`
	WaterInput
}

type WaterCharts struct {
	Consumption ChartModel
	Price       ChartModel
}
