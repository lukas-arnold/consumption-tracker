package models

type ElectricityInput struct {
	TimeFrom    string  `json:"TimeFrom"`
	TimeTo      string  `json:"TimeTo"`
	Consumption float64 `json:"Consumption"`
	Costs       float64 `json:"Costs"`
	Retailer    string  `json:"Retailer"`
	Payments    float64 `json:"Payments"`
	Note        string  `json:"Note"`
}

type Electricity struct {
	Id int64 `json:"Id"`
	ElectricityInput
}
