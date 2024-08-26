// Package converters includes all converters from foreign formats to iCal
//
// Icarus parses the incoming file data and produces an iCal for output
package converters

import (
	"github.com/emersion/go-ical"
	"io"
)

// The BaseConverter is the interface for all Icarus converters
type BaseConverter interface {
	// Convert converts the incoming file data and produces a calendar object
	Convert(input io.Reader, output *ical.Calendar) error
}
