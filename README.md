![](logo.png "icarus Logo")

# icarus

iCal file processor

## Introduction

icarus is a command line utility that processes iCal files
as defined in [RFC5545](https://www.rfc-editor.org/rfc/rfc5545).

It utilizes the common unix philosophy of commands doing
one specific thing with one output piped to another command.

> Version 1 was developed in Java while icarus now is developed
> in Golang.

The packages that make up icarus can also be used in form
of a Golang package.

## Usage

Download the last stable release of icarus for the platform
and architecture of your choice and then run the `icarus` command 
following the  processor subcommand and relevant arguments.

Use `icarus --help` for the available subcommands and
arguments. Use `icarus <subcommand> --help` for more
information about the specific subcommand.

## Selectors

icarus supports selecting events from the incoming data
to be processed by the selected processor.

By default, all events are selected. The following selectors
are available:

* text selector: the `--selector` argument can be used
  to select events based on a regular expression.
  By default, the properties "summary" and "description"
  are searched for the pattern. This can be changed
  using the `--selector-props` list
* time selector: the `--timestamp-start` and
  `--timestamp-end` arguments can be used to select 
  events starting after and ending before the given
  timestamp respectively. The timestamp needs to be
  RFC3339-formatted

## Outputs

icarus outputs the data in iCal format by default, you
can use the `--output-type` argument to specify another
output type.

Currently, `list` is another output type that can be used
to show calendar entries in a list.

## Processors

The following processors are currently available:

### `addAlarm`

Adds an alarm before the selected events. The
argument `--alarm-before` specifies how many minutes
the alarm should be before the start of the event.

### `addDTStamp`

Adds a DTSTAMP property as specified by the `--timestamp`
argument. If an event already has a DTSTAMP property, it may
be overwritten by using the `--overwrite` argument.

### `convertAllDay`

Convert all day events into events with a start and end time
("timed events") or vice versa if the `--all-day` flag is used.

Start and end times have to be specified in a colon-separated,
24 hour format (like 13:00) in the arguments `--start` and
`--end` respectively.

If all day events span multiple dates, the `--compress` argument
can be used to make them only span the start date.

If UTC is not the expected timezone, the [IANA timezone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) 
name can be set using the `--timezone` argument.

### `filter`

Only output the events matching the selector or, if the
`--inverse` flag is used, events *not* matching the selector.

### `print`

Output *all* events from the source calendar. The selectors are
ignored for this subcommand. Can be used to make use of icarus'
output types.

## Contributing

We're happy to accept contributions to icarus. Please make sure
to use the issue tracker first before you send pull requests
to make sure the task can not be done with the current set of
options.

Since version 2 icarus is developed in [Golang](https://go.dev/).

Please make sure to supply sufficient unit tests to the changes
you submit. Always branch from the develop branch.

### Creating new processors

A processor has to ímplement the interface [*BaseProcessor*](pkg/processors/processors.go).

`Initialize` is used to create a new [argparse Command type](https://pkg.go.dev/github.com/akamensky/argparse#Command)
that runs the processor. Use *uppercase short arguments* to distinguish
them from other processors and the main and output arguments.

`SetToolbox` is called to set the [toolbox](pkg/processors/toolbox.go)
variable for the processor. The toolbox contains several useful
tools for developing a processor - most importantly the
`EventMatchesSelector` function that needs to be called for
every event so that the selector arguments are properly used.

Finally, `Process` is called with the incoming events and a reference
to the output events, both in form of a [Calendar type](https://pkg.go.dev/github.com/arran4/golang-ical#Calendar)
from the [golang-ical](https://pkg.go.dev/github.com/arran4/golang-ical) package.

To activate a processor, add it to the `GetProcessors` function in
[processors.go](pkg/processors/processors.go).

Be sure to create sufficient unit tests in a _test file and add a
documentation to this readme.

### Creating new output types

Output types implement the [*BaseOutputType* interface](pkg/outputTypes/output_types.go).

`Initialize` is used to add arguments to the main icarus argument parser.
See the [argparse package](https://pkg.go.dev/github.com/akamensky/argparse) for details.

Use *lowercase short arguments* to differentiate them from processor arguments, but
make sure they don't clash with the arguments of other output types or
main arguments.

`Generate` is provided with the processed calendar entries in form of
a [Calendar type](https://pkg.go.dev/github.com/arran4/golang-ical#Calendar) and an
[io.Writer](https://pkg.go.dev/io#Writer) variable to which the output should be written to.

`GetHelp` returns a short string about the function of the output.

Add the new output type to the return value of the `GetOutputTypes` function in 
[output_types.go](pkg/outputTypes/output_types.go)

Finally, provide sufficient unit tests in a _test file and update the documentation.
