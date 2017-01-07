package de.dieploegers.icarus;

import de.dieploegers.icarus.exceptions.ProcessException;
import de.dieploegers.icarus.modifier.Modifier;
import net.fortuna.ical4j.data.CalendarBuilder;
import net.fortuna.ical4j.data.ParserException;
import net.fortuna.ical4j.model.Calendar;
import net.fortuna.ical4j.model.component.CalendarComponent;
import net.fortuna.ical4j.model.component.VEvent;
import org.apache.commons.cli.*;
import org.reflections.Reflections;

import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.InputStream;
import java.net.URISyntaxException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

public class Main {
    public static void main(String[] args) throws java.text.ParseException,
        IOException, URISyntaxException {
        Options options = new Options();

        options.addOption(
            Option.builder("h")
            .longOpt("help")
            .desc("Show this help")
            .build()
        );

        options.addOption(
            Option.builder()
                .longOpt("query")
                .hasArg()
                .desc("A regular expression, that appointment names have to match")
                .build()
        );
        options.addOption(
            Option.builder()
                .longOpt("from")
                .hasArg()
                .desc("Match all events starting at this timestamp in the " +
                    "format yyyyMMddHHmmss")
                .build()
        );
        options.addOption(
            Option.builder()
                .longOpt("until")
                .hasArg()
                .desc("Match all events ending at this timestamp in the " +
                    "format yyyyMMddHHmmss")
                .build()
        );

        // Search for modifiers

        Reflections reflections = new Reflections(
            "de.dieploegers.icarus.modifier"
        );

        List<Modifier> modifiers = new ArrayList<>();

        for (
            Class<? extends Modifier> modifier :
            reflections.getSubTypesOf(Modifier.class)
            ) {
            try {
                modifiers.add(
                    (Modifier) Class.forName(modifier.getName()).newInstance()
                );
            } catch (InstantiationException | IllegalAccessException |
                ClassNotFoundException e) {
                System.out.println(
                    "Error scanning modifiers: " + e.toString()
                );
                System.exit(1);
            }
        }

        for (Modifier modifier : modifiers) {
            for (Option option : modifier.getOptions()) {
                options.addOption(option);
            }
        }

        CommandLine commandLine = null;
        try {
            commandLine = new DefaultParser().parse(
                options,
                args
            );
        } catch (ParseException e) {
            usage(options);
            System.exit(1);
        }

        if (commandLine == null) {
            usage(options);
            System.exit(1);
        }

        if (commandLine.hasOption("help")) {
            usage(options);
            System.exit(0);
        }

        InputStream stream = null;
        if (commandLine.getArgList().size() == 0) {
            stream = System.in;
        } else {
            String icalFilename = commandLine.getArgList().get(0);

            try {
                stream = new FileInputStream(icalFilename);
            } catch (FileNotFoundException e) {
                System.out.println("File not found");
                usage(options);
                System.exit(1);
            }
        }

        CalendarBuilder builder = new CalendarBuilder();

        Calendar calendar = null;
        try {
            calendar = builder.build(stream);
        } catch (IOException e) {
            System.out.println("Error reading file.");
            usage(options);
        } catch (ParserException e) {
            System.out.println("Error parsing file.");
            usage(options);
        }

        if (calendar == null) {
            System.out.println("Error loading calendar");
            System.exit(1);
        }

        Date dateFrom = null;

        if (commandLine.hasOption("from")) {
            dateFrom = new SimpleDateFormat("yyyyMMddHHmmss").parse(
                commandLine.getOptionValue("from")
            );
        }

        Date dateUntil = null;

        if (commandLine.hasOption("until")) {
            dateUntil = new SimpleDateFormat("yyyyMMddHHmmss").parse(
                commandLine.getOptionValue("until")
            );
        }

        List<VEvent> matchedEvents = new ArrayList<>();

        for (CalendarComponent component : calendar.getComponents()) {
            if (component.getClass().getSimpleName().equals("VEvent")) {
                VEvent event = (VEvent) component;

                if (
                    !commandLine.hasOption("query") ||
                        event.getName().matches(
                            commandLine.getOptionValue("query")
                        )
                    ) {

                    // Skip events not in the requested timeline

                    boolean dateFiltered = false;

                    if (
                        dateFrom != null &&
                            (
                                event.getStartDate().getDate().before(dateFrom) ||
                                    event.getEndDate().getDate().before(dateFrom)
                            )
                        ) {
                        dateFiltered = true;
                    }

                    if (
                        dateUntil != null &&
                            (
                                event.getEndDate().getDate().after(dateUntil) ||
                                    event.getEndDate().getDate().after(dateUntil)
                            )
                        ) {
                        dateFiltered = true;
                    }

                    if (!dateFiltered) {

                        matchedEvents.add(event);

                        for (Modifier modifier : modifiers) {
                            try {
                                modifier.process(
                                    commandLine,
                                    calendar,
                                    event
                                );
                            } catch (ProcessException e) {
                                System.out.println(
                                    "Error modifying event: " + e.toString()
                                );
                                System.exit(1);
                            }
                        }

                    }

                }
            }
        }

        for (Modifier modifier : modifiers) {
            try {
                modifier.finalize(
                    commandLine,
                    calendar,
                    matchedEvents
                );
            } catch (ProcessException e) {
                System.out.println(
                    "Error finalizing events: " + e.toString()
                );
                System.exit(1);
            }
        }

        System.out.println(calendar.toString());

        System.exit(0);

    }

    private static void usage(Options options) {
        HelpFormatter formatter = new HelpFormatter();

        formatter.printHelp(
            "Main",
            "Modify an ical file",
            options,
            "Please visit https://github.com/dploeger/icarus for further " +
                "information",
            true
        );
    }
}
