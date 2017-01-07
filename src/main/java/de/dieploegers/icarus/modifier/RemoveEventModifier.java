package de.dieploegers.icarus.modifier;

import de.dieploegers.icarus.ModifierOption;
import de.dieploegers.icarus.OptionStore;
import de.dieploegers.icarus.exceptions.ProcessException;
import net.fortuna.ical4j.model.Calendar;
import net.fortuna.ical4j.model.component.VEvent;

import java.util.ArrayList;
import java.util.List;

/**
 * Removes an Event
 */
public class RemoveEventModifier implements Modifier {

    @Override
    public List<ModifierOption> getOptions() {
        List<ModifierOption> options = new ArrayList<>();
        options.add(
            new ModifierOption(
                "removeEvent",
                "Remove all matching events"
            )
        );
        return options;
    }

    @Override
    public void process(
        OptionStore options, Calendar calendar, VEvent event
    ) throws ProcessException {

    }

    @Override
    public void finalize(
        OptionStore options, Calendar calendar, List<VEvent> matchedEvents
    ) throws ProcessException {
        if (options.isSet("removeEvent")) {
            calendar.getComponents().removeAll(matchedEvents);
        }
    }
}
