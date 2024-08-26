// Package converteradapters holds CLI adapters that connect the icarus CLI to the converters
package converteradapters

import (
	"github.com/akamensky/argparse"
	"github.com/emersion/go-ical"
	"os"
)

// The ConverterAdapter connects the Icarus CLI with a processor
type ConverterAdapter interface {
	// Initialize creates a new subcommand for the argparse parser.
	Initialize(parser *argparse.Parser) (*argparse.Command, error)
	// Convert converts the incoming calendar and fills the output calendar
	Convert(input *os.File, output *ical.Calendar) error
}

// GetConverterAdapters returns a list of enabled processor adapters
func GetConverterAdapters() []ConverterAdapter {
	return []ConverterAdapter{&CSVConverterAdapter{}}
}
