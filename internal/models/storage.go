package models

type UsageStorage struct {
	Electricity   []Electricity  `json:"Electricity"`
	Oil           []Oil          `json:"Oil"`
	OilFillLevels []OilFillLevel `json:"OilFillLevels"`
	Water         []Water        `json:"Water"`
}

type ElectricityForChart struct {
	Electricities []Electricity
	Labels        []string
	Usages        []float64
	Costs         []float64
	Prices        []float64
}

type OilForChart struct {
	Oils    []Oil
	Labels  []string
	Volumes []float64
	Costs   []float64
	Prices  []float64
}

type OilFillLevelForChart struct {
	OilFillLevels []OilFillLevel
	Dates         []string
	Levels        []float64
}

type WaterForChart struct {
	Waters  []Water
	Labels  []string
	Volumes []float64
	VolumesWater []float64
	VolumesWastewater []float64
	VolumesRainwater []float64
	Costs   []float64
	CostsWater []float64
	CostsWastewater []float64
	CostsRainwater []float64
	Prices  []float64
	PricesWater []float64
	PricesWastewater []float64
	PricesRainwater []float64
	FixedPrices []float64
}
