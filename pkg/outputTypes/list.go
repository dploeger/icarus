package outputTypes

import (
	"github.com/akamensky/argparse"
	"github.com/emersion/go-ical"
	"github.com/olekukonko/tablewriter"
	"io"
)

// The ListOutputType converts the internal calendar into a table format
type ListOutputType struct {
}

func (t *ListOutputType) Initialize(_ *argparse.Parser) error {
	return nil
}

func (t *ListOutputType) Generate(calendar *ical.Calendar, writer io.Writer, options OutputOptions) error {
	table := tablewriter.NewWriter(writer)
	table.SetHeader(options.Columns)
	for _, event := range calendar.Events() {
		var row []string
		for _, column := range options.Columns {
			if event.Props.Get(column) == nil {
				row = append(row, "")
			} else {
				row = append(row, event.Props.Get(column).Value)
			}

		}
		table.Append(row)
	}
	table.Render()
	return nil
}

func (t *ListOutputType) GetHelp() string {
	return "outputs a list of events from the processed calendar"
}

var _ BaseOutputType = &ListOutputType{}
