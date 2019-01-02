package de.dieploegers.icarus.modifier;

import de.dieploegers.icarus.ModifierOption;
import de.dieploegers.icarus.OptionStore;
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
        final List<ModifierOption> options = new ArrayList<>();
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
        final OptionStore options, final Calendar calendar, final VEvent event
    ) {

    }

    @Override
    public void finalize(
        final OptionStore options, final Calendar calendar, final List<VEvent> matchedEvents
    ) {
        if (options.isSet("removeEvent")) {
            calendar.getComponents().removeAll(matchedEvents);
        }
    }
}
