package entities

type CellQuality struct {
	SignalStrength int `json:"signal_strength"`
}

type CellCarrier struct {
	CarrierName string `json:"carrier_name"`
}
