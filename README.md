# icarus - iCal batch processor
[![Travis](https://img.shields.io/travis/dploeger/icarus.svg)](https://travis-ci.org/dploeger/icarus)

![icarus](design/logo.png)

## Introduction

icarus does batch processing on [iCal/ICS](https://en.wikipedia.org/wiki/ICalendar)
files. It's easy to use and extend.

It is written in Java 8, uses the wonderful [iCal4j library](http://ical4j.github.io/) for ICS
file processing, the [reflection library](https://github.com/ronmamo/reflections)
for easy extensibility and [Apache Commons CLI](http://commons.apache.org/proper/commons-cli/index.html) for CLI processing.

Currently, the following processing features are available:

* Add an alarm
* Convert all day events to date time events
* Remove events

The ICS file used in the test suite is taken from [schulferien.org](http://www.schulferien.org/deutschland/ical/). They do a great job regularly providing holidays and other events as ical files.

## Disclaimer

icarus is currently beta software. It will not trash your ics files, but
 you might experience problems when importing the resulting ics data.

## Download

Check the [releases page](https://github.com/dploeger/icarus/releases)
for a release and simply grab the attached JAR file.

## Requirements

icarus is based on Java 8, so a corresponding JRE is needed. All
dependent libraries are packed into the JAR file.

## Usage

icarus is a command line application, that is called with
arguments and the source ics file:

    java -jar icarus.jar [options] <icsfile.ics>

If you omit the source file, it checks STDIN for ICS data.

Use

    java -jar icarus.jar --help

to get detailed help.

Basically, icarus expects a query and/or a time range to search
events.

The query is a regular expression that is matched against the title of
the appointment. The arguments "from" and "until" describe a time range,
that the event has to be in.

The desired processor is configured using the corresponding
argument. For example, the argument "--removeEvent" will remove all events
in the filtered range or with the filtered query.

The resulting ICS data is written to stdout and can be piped to a new
ICS file to be used in your favorite calendar app.

## Using icarus as a library

You can also use icarus as a library. Just instantiate a [Processor](https://dploeger.github.io/icarus/apidocs/de/dieploegers/icarus/Processor.html)
and call "process" with an OptionStore and the source ical data and
you'll get the processed data back.

Refer to the [API-Docs](https://dploeger.github.io/icarus/apidocs/index.html)
for more information.

You can use the icarus.lib.jar for library use. I will be
providing icarus over Maven some time later.

## Extending icarus

Currently, icarus does all the processing for which
[it was designed for](http://dennis.dieploegers.de/flying-high-on-ical-files/)

However, it is designed to be open to new modifiers and if you'd like
icarus to do something special with your ical events, you can just implement
the [Modifier interface](https://github.com/dploeger/icarus/blob/master/src/main/java/de/dieploegers/icarus/modifier/Modifier.java)
and you're done.

The interface is quite easy. The following methods have to be implemented:

* getOptions: Provide a list of [ModifierOptions](https://dploeger.github.io/icarus/apidocs/de/dieploegers/icarus/ModifierOption.html), that indicate which
  options your modifier expects
* process: This is called for all events that match the filter. It is
  called with the [OptionStore](https://dploeger.github.io/icarus/apidocs/de/dieploegers/icarus/OptionStore.html) holding all given options
  so you can check for your options, the current event and the
  complete calendar object for reference. Just directly modify the event.
  Do **not** modify the calendar object here, or else you're getting
  concurrent modification exceptions.
* finalize: This is called after all events have been processed and this
  is the point where you can modify the complete calendar, if you'd
  like to. You'll get the OptionStore, the calendar and a list
  of events, that have been filtered before

See the [API-Docs](https://dploeger.github.io/icarus/apidocs/index.html) for more information
and use the existing modifiers as a reference.

Go for it! I'm happy to accept pull requests for new features.
