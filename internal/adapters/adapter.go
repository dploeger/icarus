// Package adapters holds CLI adapters that connect the icarus CLI to the processors
package adapters

import (
	"github.com/akamensky/argparse"
	"github.com/dploeger/icarus/v2/pkg/processors"
	"github.com/emersion/go-ical"
)

// The Adapter connects the Icarus CLI with a processor
type Adapter interface {
	// Initialize creates a new subcommand for the argparse parser.
	Initialize(parser *argparse.Parser) (*argparse.Command, error)
	// SetToolbox sets the toolbox that can be used by the processor
	SetToolbox(toolbox processors.Toolbox)
	// Process processes the incoming calendar and fills the output calendar
	Process(input ical.Calendar, output *ical.Calendar) error
}

// GetAdapters returns a list of enabled processor adapters
func GetAdapters() []Adapter {
	return []Adapter{
		&FilterAdapter{},
		&PrintAdapter{},
		&ConvertAllDayAdapter{},
		&AddDTStampAdapter{},
		&AddAlarmAdapter{},
		&AddPropertyAdapter{},
		&DeletePropertyAdapter{},
	}
}
