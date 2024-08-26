package converteradapters

import (
	"github.com/akamensky/argparse"
	"github.com/dploeger/icarus/v2/pkg/converters"
	"github.com/emersion/go-ical"
	"os"
	"strings"
	"time"
)

type CSVConverterAdapter struct {
	fieldMapOption  *[]string
	separator       *string
	timestampFormat *string
	hasHeaders      *bool
	headers         *[]string
}

var _ ConverterAdapter = &CSVConverterAdapter{}

func (c *CSVConverterAdapter) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
	command := parser.NewCommand("convertCSV", "Convert a CSV file to a calendar")
	c.separator = command.String("S", "separator", &argparse.Options{
		Help:    "The field separator used in the file",
		Default: ",",
	})
	c.hasHeaders = command.Flag("H", "hasHeaders", &argparse.Options{
		Help: "The CSV file has headers in the first row",
	})
	c.headers = command.StringList("E", "headers", &argparse.Options{
		Help:    "A list of headers if the CSV file has no headers in the first row",
		Default: []string{},
	})
	c.fieldMapOption = command.StringList("F", "field", &argparse.Options{
		Help:    "A map of headers to iCal field names in the form of header:field (e.g. startDate:DTSTART)",
		Default: []string{},
	})
	c.timestampFormat = command.String("T", "timestamp", &argparse.Options{
		Help:    "The format of the timestamps in the CSV file. It uses the golang time format (https://go.dev/src/time/format.go)",
		Default: time.RFC3339,
	})
	return command, nil
}

func (c *CSVConverterAdapter) Convert(input *os.File, output *ical.Calendar) error {
	fieldMap := make(map[string]string)
	for _, o := range *c.fieldMapOption {
		header, field := func() (string, string) {
			a := strings.Split(o, ":")
			return a[0], a[1]
		}()
		fieldMap[header] = field
	}

	converter := converters.CSVConverter{
		Separator:       *c.separator,
		TimestampFormat: *c.timestampFormat,
		FieldMap:        fieldMap,
		HasHeaders:      *c.hasHeaders,
		Headers:         *c.headers,
	}

	return converter.Convert(input, output)
}
