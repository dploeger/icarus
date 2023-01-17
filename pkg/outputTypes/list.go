package outputTypes

import (
	"github.com/akamensky/argparse"
	"github.com/emersion/go-ical"
	"github.com/olekukonko/tablewriter"
	"io"
)

type ListOutputType struct {
	columns *[]string
}

func (t *ListOutputType) Initialize(parser *argparse.Parser) error {
	t.columns = parser.StringList("c", "columns", &argparse.Options{
		Help: "Columns to display",
	})
	return nil
}

func (t *ListOutputType) Generate(calendar *ical.Calendar, writer io.Writer) error {
	if t.columns == nil {
		t.columns = &[]string{"SUMMARY", "DTSTART", "DTEND", "DESCRIPTION"}
	}
	table := tablewriter.NewWriter(writer)
	table.SetHeader(*t.columns)
	for _, event := range calendar.Events() {
		var row []string
		for _, column := range *t.columns {
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
