package de.dieploegers.icarus;

import de.dieploegers.icarus.modifier.Modifier;
import net.fortuna.ical4j.data.ParserException;
import org.apache.commons.cli.*;
import org.apache.commons.io.IOUtils;

import java.io.*;

public class Main {
    public static void main(final String[] args) {
        final Options options = new Options();

        // Add CLI options

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

        // Instantiate processor

        Processor processor = null;
        try {
            processor = new Processor();
        } catch (final InstantiationException | IllegalAccessException |
            ClassNotFoundException e) {
            System.out.println(
                "Error scanning modifiers: " + e.toString()
            );
            System.exit(1);
        }

        // Add modifier options

        for (final Modifier modifier : processor.getModifiers()) {
            for (final ModifierOption option : modifier.getOptions()) {
                final Option.Builder optionBuilder = Option.builder();

                optionBuilder.longOpt(option.getKey());
                optionBuilder.desc(option.getDescription());

                if (option.getHasValue()) {
                    optionBuilder.hasArg();
                }

                options.addOption(
                    optionBuilder.build()
                );
            }
        }

        // Parse commandline

        CommandLine commandLine = null;
        try {
            commandLine = new DefaultParser().parse(
                options,
                args
            );
        } catch (final ParseException e) {
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

        // Build up optionstore

        final OptionStore optionStore = new OptionStore();

        for (final Modifier modifier: processor.getModifiers()) {
            for (final ModifierOption option: modifier.getOptions()) {
                if (commandLine.hasOption(option.getKey())) {
                    final ModifierOption commandOption = new ModifierOption(
                        option.getKey()
                    );
                    if (option.getHasValue()) {
                        commandOption.setValue(
                            commandLine.getOptionValue(option.getKey())
                        );
                    }
                    optionStore.addOption(commandOption);
                }
            }
        }

        // Load ICS file/stream

        InputStream stream = null;
        if (commandLine.getArgList().size() == 0) {
            stream = System.in;
        } else {
            final String icalFilename = commandLine.getArgList().get(0);

            try {
                stream = new FileInputStream(icalFilename);
            } catch (final FileNotFoundException e) {
                System.out.println("File not found");
                usage(options);
                System.exit(1);
            }
        }

        final StringWriter writer = new StringWriter();
        try {
            IOUtils.copy(stream, writer, "UTF-8");
        } catch (final IOException e) {
            System.out.println("Can not read calendar");
            usage(options);
            System.exit(1);
        }
        final String iCalData = writer.toString();

        // Set from/until/query

        if (commandLine.hasOption("from")) {
            optionStore.addOption(
                ModifierOption.withValue(
                    "from",
                    commandLine.getOptionValue("from")
                )
            );
        }

        if (commandLine.hasOption("until")) {
            optionStore.addOption(
                ModifierOption.withValue(
                    "until",
                    commandLine.getOptionValue("until")
                )
            );
        }

        if (commandLine.hasOption("query")) {
            optionStore.addOption(
                ModifierOption.withValue(
                    "query",
                    commandLine.getOptionValue("query")
                )
            );
        }

        // Do the processing

        try {
            System.out.println(processor.process(optionStore, iCalData));
            System.exit(0);
        } catch (final IOException | ParserException e) {
            System.out.println("Error parsing calendar: " + e.toString());
            usage(options);
        } catch (final java.text.ParseException e) {
            System.out.println("Error parsing from or until dates: " + e
                .toString());
            usage(options);
        }

    }

    /**
     * Show usage information
     * @param options The generated CLI options
     */

    private static void usage(final Options options) {
        final HelpFormatter formatter = new HelpFormatter();

        formatter.printHelp(
            "icarus.jar",
            "iCal batch processor",
            options,
            "Please visit https://github.com/dploeger/icarus for further " +
                "information",
            true
        );
    }
}
