package models

type ConsumptionStorage struct {
	Electricity   []Electricity  `json:"electricity"`
	Oil           []Oil          `json:"oil"`
	OilFillLevels []OilFillLevel `json:"oilFillLevels"`
	Water         []Water        `json:"water"`
}

type ChartDataset struct {
	Label  string    `json:"label"`
	Data   []float64 `json:"data"`
	Unit   string    `json:"unit,omitempty"`
	YAxis  string    `json:"yAxis,omitempty"`
	Hidden bool      `json:"hidden,omitempty"`
}

type ChartModel struct {
	Labels []string       `json:"labels"`
	Sets   []ChartDataset `json:"sets"`
}
