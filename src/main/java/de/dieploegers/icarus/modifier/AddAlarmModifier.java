package de.dieploegers.icarus.modifier;

import de.dieploegers.icarus.ModifierOption;
import de.dieploegers.icarus.OptionStore;
import net.fortuna.ical4j.model.Calendar;
import net.fortuna.ical4j.model.Dur;
import net.fortuna.ical4j.model.component.VAlarm;
import net.fortuna.ical4j.model.component.VEvent;
import net.fortuna.ical4j.model.property.Action;
import net.fortuna.ical4j.model.property.Description;

import java.util.ArrayList;
import java.util.List;

/**
 * Add alarm to all filtered events
 */
public class AddAlarmModifier implements Modifier {
    @Override
    public List<ModifierOption> getOptions() {
        final List<ModifierOption> options = new ArrayList<>();
        options.add(
            new ModifierOption(
                "alarmBefore",
                "Add an alarm before all appointments (hours)",
                true
            )
        );
        options.add(
            new ModifierOption(
                "alarmDuration",
                "The concrete duration string used in ics for the alarm (see https://www.kanzaki.com/docs/ical/duration-t.html)",
                true
            )
        );
        options.add(
            new ModifierOption(
                "alarmMessage",
                "The alarm message to use.",
                true
            )
        );
        return options;
    }

    @Override
    public void process(
        final OptionStore options, final Calendar calendar, final VEvent event
    )
    {
        if (options.isSet("alarmBefore") || options.isSet("alarmDuration")) {
            final VAlarm alarm;

            if (options.isSet("alarmBefore")) {
                alarm = new VAlarm(
                    new Dur(
                        0,
                        Integer.valueOf(
                            options.get(
                                "alarmBefore"
                            )
                        ) * -1,
                        0,
                        0
                    )
                );
            } else {
                alarm = new VAlarm(
                    new Dur(options.get("alarmDuration"))
                );
            }

            alarm.getProperties().add(Action.DISPLAY);
            if (options.isSet("alarmMessage")) {
                alarm.getProperties().add(new Description(
                    options.get("alarmMessage")
                ));
            } else {
                alarm.getProperties().add(new Description(
                    "Alarm provided by icarus"
                ));
            }

            event.getAlarms().add(
                alarm
            );
        }

    }

    @Override
    public void finalize(final OptionStore options, final Calendar calendar, final List<VEvent> matchedEvents)
    {

    }

}
