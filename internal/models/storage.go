package models

type ConsumptionStorage struct {
	Electricity   []Electricity  `json:"Electricity"`
	Oil           []Oil          `json:"Oil"`
	OilFillLevels []OilFillLevel `json:"OilFillLevels"`
	Water         []Water        `json:"Water"`
}

type ChartDataset struct {
	Label string    `json:"label"`
	Data  []float64 `json:"data"`
	Unit  string    `json:"unit,omitempty"`
}

type ChartModel struct {
	Labels []string       `json:"labels"`
	Sets   []ChartDataset `json:"sets"`
}
