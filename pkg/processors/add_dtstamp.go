package processors

import (
	"github.com/akamensky/argparse"
	"github.com/emersion/go-ical"
	"time"
)

type AddDTStampProcessor struct {
	timestamp *string
	overwrite *bool
	toolbox   Toolbox
}

func (t *AddDTStampProcessor) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
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

func (t *AddDTStampProcessor) SetToolbox(toolbox Toolbox) {
	t.toolbox = toolbox
}

func (t *AddDTStampProcessor) Process(input ical.Calendar, output *ical.Calendar) error {
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
	for _, event := range input.Events() {
		if t.toolbox.EventMatchesSelector(event) {
			if event.Props.Get(ical.PropDateTimeStamp) == nil || *t.overwrite {
				event.Props.SetDateTime(ical.PropDateTimeStamp, parsedTimestamp)
			}
		}
		output.Children = append(output.Children, event.Component)
	}
	return nil
}
