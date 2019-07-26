package entities

import (
	"github.com/parrotmac/go-modemmanager/pkg/modem"
)

type CellInfoProvider interface {
	Setup() error
	FetchReport() (*ModemReport, error)
}

type ModemReport struct {
	Modem  *modem.Modem  `json:"modem"`
	Bearer *modem.Bearer `json:"bearer"`
	Sim    *modem.Sim    `json:"sim"`
}
