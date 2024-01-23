package processors

import (
	"github.com/akamensky/argparse"
	"github.com/emersion/go-ical"
)

// The PrintProcessor simply outputs the input calendar
type PrintProcessor struct {
	toolbox Toolbox
}

func (f *PrintProcessor) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
	command := parser.NewCommand("print", "Outputs the calendar")
	return command, nil
}

func (f *PrintProcessor) Process(inputCalendar ical.Calendar, outputCalendar *ical.Calendar) error {
	for _, event := range inputCalendar.Events() {
		outputCalendar.Children = append(outputCalendar.Children, event.Component)
	}
	return nil
}

func (f *PrintProcessor) SetToolbox(toolbox Toolbox) {
	f.toolbox = toolbox
}

var _ BaseProcessor = &PrintProcessor{}
