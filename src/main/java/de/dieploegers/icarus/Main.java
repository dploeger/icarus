package de.dieploegers.icarus;

import de.dieploegers.icarus.modifier.Modifier;
import net.fortuna.ical4j.data.ParserException;
import org.apache.commons.cli.*;
import org.apache.commons.io.IOUtils;

import java.io.*;

public class Main {
    public static void main(String[] args) {
        Options options = new Options();

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
        } catch (InstantiationException | IllegalAccessException |
            ClassNotFoundException e) {
            System.out.println(
                "Error scanning modifiers: " + e.toString()
            );
            System.exit(1);
        }

        // Add modifier options

        for (Modifier modifier : processor.getModifiers()) {
            for (ModifierOption option : modifier.getOptions()) {
                Option.Builder optionBuilder = Option.builder();

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

        // Build up optionstore

        OptionStore optionStore = new OptionStore();

        for (Modifier modifier: processor.getModifiers()) {
            for (ModifierOption option: modifier.getOptions()) {
                if (commandLine.hasOption(option.getKey())) {
                    ModifierOption commandOption = new ModifierOption(
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
            String icalFilename = commandLine.getArgList().get(0);

            try {
                stream = new FileInputStream(icalFilename);
            } catch (FileNotFoundException e) {
                System.out.println("File not found");
                usage(options);
                System.exit(1);
            }
        }

        StringWriter writer = new StringWriter();
        try {
            IOUtils.copy(stream, writer, "UTF-8");
        } catch (IOException e) {
            System.out.println("Can not read calendar");
            usage(options);
            System.exit(1);
        }
        String iCalData = writer.toString();

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
        } catch (IOException | ParserException e) {
            System.out.println("Error parsing calendar: " + e.toString());
            usage(options);
        } catch (java.text.ParseException e) {
            System.out.println("Error parsing from or until dates: " + e
                .toString());
            usage(options);
        }

    }

    /**
     * Show usage information
     * @param options The generated CLI options
     */

    private static void usage(Options options) {
        HelpFormatter formatter = new HelpFormatter();

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
