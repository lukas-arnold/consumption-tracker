package models

type ElectricityInput struct {
	TimeFrom    string  `json:"timeFrom"`
	TimeTo      string  `json:"timeTo"`
	Consumption float64 `json:"consumption"`
	Costs       float64 `json:"costs"`
	Retailer    string  `json:"retailer"`
	Payments    float64 `json:"payments"`
	Note        string  `json:"note"`
}

type Electricity struct {
	Id int64 `json:"id"`
	ElectricityInput
}

type ElectricityCharts struct {
	Consumption ChartModel
	Price       ChartModel
}
