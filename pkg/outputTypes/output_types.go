// Package outputTypes includes all valid output types that Icarus supports.
// outputtypes use the internal processed calendar structure and format it for the output.
package outputTypes

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/emersion/go-ical"
	"github.com/thoas/go-funk"
	"io"
	"strings"
)

// The BaseOutputType is the basic interface for all output types
type BaseOutputType interface {
	// Initialize can add arguments to the argument parser.
	Initialize(parser *argparse.Parser) error
	// Generate generates the output based on the given calendar and writes it to writer
	Generate(calendar *ical.Calendar, writer io.Writer) error
	// GetHelp returns a help string about what the output produces
	GetHelp() string
}

// GetOutputTypes returns a list of enabled output types
func GetOutputTypes() map[string]BaseOutputType {
	m := make(map[string]BaseOutputType)
	m["ics"] = &ICSOutputType{}
	m["list"] = &ListOutputType{}
	return m
}

// GetOutputHelp returns a formatted help string for all available output types
func GetOutputHelp() string {
	outputTypes := GetOutputTypes()
	return strings.Join(
		funk.Reduce(
			funk.Keys(outputTypes),
			func(helps []string, outputType string) []string {
				return append(helps, fmt.Sprintf("\t\t\t%s: %s", outputType, outputTypes[outputType].GetHelp()))
			},
			[]string{},
		).([]string), "\n")
}
