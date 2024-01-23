// Package processors includes all calendar processors available in Icarus.
//
// Icarus parses the incoming calendar data, hands it over to the processor to process data in it and
// formats the resulting data using an output type
package processors

import (
	"github.com/akamensky/argparse"
	"github.com/emersion/go-ical"
)

// The BaseProcessor is the interface for all Icarus processors
type BaseProcessor interface {
	// Initialize creates a new subcommand for the argparse parser.
	Initialize(parser *argparse.Parser) (*argparse.Command, error)
	// SetToolbox sets the toolbox that can be used by the processor
	SetToolbox(toolbox Toolbox)
	// Process processes the incoming calendar and fills the output calendar
	Process(input ical.Calendar, output *ical.Calendar) error
}

// GetProcessors returns a list of enabled processors
func GetProcessors() []BaseProcessor {
	return []BaseProcessor{
		&FilterProcessor{},
		&PrintProcessor{},
		&ConvertAllDayProcessor{},
		&AddDTStampProcessor{},
		&AddAlarmProcessor{},
		&AddPropertyProcessor{},
		&DeletePropertyProcessor{},
	}
}
