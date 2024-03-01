package outputTypes

import (
	csv2 "encoding/csv"
	"github.com/akamensky/argparse"
	"github.com/emersion/go-ical"
	"io"
)

type CSVOutputType struct {
	separator *string
}

func (t *CSVOutputType) Initialize(parser *argparse.Parser) error {
	t.separator = parser.String("S", "separator", &argparse.Options{
		Help:    "Separator to use",
		Default: ",",
	})
	return nil
}

func (t *CSVOutputType) Generate(calendar *ical.Calendar, writer io.Writer, options OutputOptions) error {
	csv := csv2.NewWriter(writer)
	csv.Comma = []rune(*t.separator)[0]
	if err := csv.Write(options.Columns); err != nil {
		return err
	}
	for _, event := range calendar.Events() {
		var row []string
		for _, column := range options.Columns {
			if event.Props.Get(column) == nil {
				row = append(row, "")
			} else {
				row = append(row, event.Props.Get(column).Value)
			}

		}
		if err := csv.Write(row); err != nil {
			return err
		}
	}
	csv.Flush()
	return nil
}

func (*CSVOutputType) GetHelp() string {
	return "outputs a list of events from the processed calendar in CSV format"
}

var _ BaseOutputType = &CSVOutputType{}
