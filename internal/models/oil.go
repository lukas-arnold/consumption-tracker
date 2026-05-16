package models

type OilInput struct {
	Date     string  `json:"Date"`
	Volume   float64 `json:"Volume"`
	Costs    float64 `json:"Costs"`
	Retailer string  `json:"Retailer"`
	Note     string  `json:"Note"`
}

type Oil struct {
	Id int64 `json:"Id"`
	OilInput
}

const OilTankHeightCm = 150.0

type OilFillLevelInput struct {
	Date  string  `json:"Date"`
	Level float64 `json:"Level"`
}

type OilFillLevel struct {
	Id int64 `json:"Id"`
	OilFillLevelInput
	Percentage float64 `json:"-"`
}
