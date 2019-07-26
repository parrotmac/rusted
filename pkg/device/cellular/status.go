package cellular

import (
	"errors"

	"github.com/godbus/dbus"
	"github.com/parrotmac/go-modemmanager/pkg/modem"
	"go.uber.org/zap"

	"github.com/parrotmac/rusted/pkg/central/entities"
)

type StatusProvider struct {
	Logger  *zap.Logger
	manager *modem.Manager
}

func (s *StatusProvider) Setup() error {
	systemDbusConn, err := dbus.SystemBus()
	if err != nil {
		s.Logger.Warn("modem.system_dbus.connection_failure", zap.Error(err))
		return err
	}

	s.manager = &modem.Manager{
		Logger:    s.Logger,
		SystemBus: systemDbusConn,
	}
	return nil
}

func (s *StatusProvider) FetchReport() (*entities.ModemReport, error) {
	if s.manager == nil {
		return nil, errors.New("manager hasn't been Setup() yet")
	}

	modemPaths, err := s.manager.GetManagedModems()
	if err != nil {
		return nil, err
	}
	if len(modemPaths) == 0 {
		return nil, errors.New("no modems found")
	}

	report := &entities.ModemReport{}

	// Note - this only deals with the first modem returned
	// Other modems could exist but won't show up
	defaultModem, err := s.manager.GetModem(modemPaths[0])
	if err != nil {
		return nil, err
	}
	report.Modem = &defaultModem

	if len(defaultModem.Bearers) > 0 {
		bearer, err := s.manager.GetBearer(defaultModem.Bearers[0])
		if err != nil {
			s.Logger.Warn("modem.get_bearer.failure", zap.Error(err))
		} else {
			report.Bearer = &bearer
		}
	}

	if defaultModem.Sim.IsValid() && string(defaultModem.Sim) != "" {
		sim, err := s.manager.GetSim(defaultModem.Sim)
		if err != nil {
			s.Logger.Warn("modem.get_sim.failure", zap.Error(err))
		} else {
			report.Sim = &sim
		}
	}
	return report, nil
}
