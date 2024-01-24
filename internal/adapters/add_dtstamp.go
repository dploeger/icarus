package adapters

import (
	"github.com/akamensky/argparse"
	"github.com/dploeger/icarus/v2/pkg/processors"
	"github.com/emersion/go-ical"
	"time"
)

// The AddDTStampAdapter adds a DTSTAMP field to all selected events
type AddDTStampAdapter struct {
	timestamp *string
	overwrite *bool
	toolbox   processors.Toolbox
}

func (t *AddDTStampAdapter) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
	command := parser.NewCommand("addDTStamp", "Adds a DTStamp field to all selected events")
	t.timestamp = command.String("T", "timestamp", &argparse.Options{
		Help: "Set DTSTAMP to this timestamp. Defaults to the current timestamp.",
	})
	t.overwrite = command.Flag("O", "overwrite", &argparse.Options{
		Help:    "Overwrite DTSTAMP if event already has one",
		Default: true,
	})
	return command, nil
}

func (t *AddDTStampAdapter) SetToolbox(toolbox processors.Toolbox) {
	t.toolbox = toolbox
}

func (t *AddDTStampAdapter) Process(input ical.Calendar, output *ical.Calendar) error {
	var parsedTimestamp time.Time
	if t.timestamp == nil || *t.timestamp == "" {
		parsedTimestamp = time.Now().In(time.UTC)
	} else {
		if parsed, err := time.Parse("20060102T150405Z", *t.timestamp); err != nil {
			return err
		} else {
			parsedTimestamp = parsed
		}
	}
	p := processors.AddDTStampProcessor{
		Timestamp: parsedTimestamp,
		Overwrite: *t.overwrite,
	}
	p.SetToolbox(t.toolbox)
	return p.Process(input, output)
}

var _ Adapter = &AddDTStampAdapter{}
