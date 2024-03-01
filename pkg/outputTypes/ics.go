package outputTypes

import (
	"github.com/akamensky/argparse"
	"github.com/emersion/go-ical"
	"io"
)

// The ICSOutputType converts the internal calendar into the standardizes ICS format
type ICSOutputType struct{}

func (t *ICSOutputType) Initialize(_ *argparse.Parser) error {
	// no need for more arguments
	return nil
}

func (t *ICSOutputType) Generate(calendar *ical.Calendar, writer io.Writer, _ OutputOptions) error {
	if err := ical.NewEncoder(writer).Encode(calendar); err != nil {
		return err
	}
	return nil
}

func (t *ICSOutputType) GetHelp() string {
	return "outputs an ICS file from the processed calendar"
}

var _ BaseOutputType = &ICSOutputType{}
