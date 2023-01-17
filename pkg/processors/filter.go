package processors

import (
	"github.com/akamensky/argparse"
	"github.com/emersion/go-ical"
)

type FilterProcessor struct {
	toolbox Toolbox
	inverse *bool
}

func (f *FilterProcessor) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
	command := parser.NewCommand("filter", "Output only events that match the selector")
	f.inverse = command.Flag("I", "inverse", &argparse.Options{
		Help:    "Inverse the function. Output events that DO NOT match the selector",
		Default: false,
	})
	return command, nil
}

func (f *FilterProcessor) Process(inputCalendar ical.Calendar, outputCalendar *ical.Calendar) error {
	for _, event := range inputCalendar.Events() {
		if (*f.inverse && !f.toolbox.EventMatchesSelector(event)) ||
			(!*f.inverse && f.toolbox.EventMatchesSelector(event)) {
			outputCalendar.Children = append(outputCalendar.Children, event.Component)
			continue
		}
	}
	return nil
}

func (f *FilterProcessor) SetToolbox(toolbox Toolbox) {
	f.toolbox = toolbox
}
