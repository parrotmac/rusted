package services

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/parrotmac/rusted/pkg/central/entities"
)

type CellularService struct {
	UpdateFrequency time.Duration
	infoProvider    entities.CellInfoProvider
	infoTransport   entities.CellInfoTransport
}

func CreateCellularService(reportingInterval time.Duration, provider entities.CellInfoProvider, transport entities.CellInfoTransport) *CellularService {
	return &CellularService{
		UpdateFrequency: reportingInterval,
		infoProvider:    provider,
		infoTransport:   transport,
	}
}

func (c *CellularService) sendReport(report *entities.ModemReport) error {
	return c.infoTransport.SendCellInfoReport(context.TODO(), report)
}

func (c *CellularService) Start() {
	ticker := time.NewTicker(c.UpdateFrequency)
	go func() {
		for range ticker.C {
			report, err := c.infoProvider.FetchReport()
			if err != nil {
				logrus.Errorln("failure fetching modem info", err)
				continue
			}
			err = c.sendReport(report)
			if err != nil {
				logrus.Errorln("unable to send cell info report:", err)
			}
		}
	}()
}
