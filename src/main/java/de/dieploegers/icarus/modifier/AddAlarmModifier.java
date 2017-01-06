package de.dieploegers.icarus.modifier;

import de.dieploegers.icarus.exceptions.ProcessException;
import net.fortuna.ical4j.model.Calendar;
import net.fortuna.ical4j.model.Dur;
import net.fortuna.ical4j.model.component.VAlarm;
import net.fortuna.ical4j.model.component.VEvent;
import org.apache.commons.cli.CommandLine;
import org.apache.commons.cli.Option;

import java.util.ArrayList;
import java.util.List;

/**
 * Add alarm to all filtered events
 */
public class AddAlarmModifier implements Modifier {
    @Override
    public List<Option> getOptions() {
        List<Option> options = new ArrayList<>();
        options.add(
            Option.builder()
                .longOpt("alarmBefore")
                .hasArg()
                .desc("Add an alarm before all appointments (hours)")
                .build()
        );
        return options;
    }

    @Override
    public void process(
        CommandLine commandLine, Calendar calendar, VEvent event
    ) {
        if (commandLine.hasOption("alarmBefore")) {
            event.getAlarms().add(
                new VAlarm(
                    new Dur(
                        0,
                        Integer.valueOf(
                            commandLine.getOptionValue(
                                "alarmBefore"
                            )
                        ),
                        0,
                        0
                    )
                )
            );
        }
    }

    @Override
    public void finalize(
        CommandLine commandLine, Calendar calendar, List<VEvent> matchedEvents
    ) throws ProcessException {
        // Nothing to do
    }
}
