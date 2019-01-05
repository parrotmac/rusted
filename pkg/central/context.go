package central

import (
	"github.com/parrotmac/rusted/pkg/central/status"
)

type Context struct {
	StatusReportingConfig status.StatusReportingConfig
	ClientIdentifier      string
}
