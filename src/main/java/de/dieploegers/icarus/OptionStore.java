package de.dieploegers.icarus;

import java.util.HashMap;
import java.util.List;

/**
 * A small utility class handling options
 */
public class OptionStore {
    private HashMap<String, ModifierOption> options;

    public OptionStore(List<ModifierOption> options) {
        this.options = new HashMap<>();
        for (ModifierOption option : options) {
            this.options.put(option.getKey(), option);
        }
    }

    public OptionStore() {
        this.options = new HashMap<>();
    }

    /**
     * Is the given key set as an option
     *
     * @param key Key name
     * @return true, if the option has been set
     */

    public boolean isSet(String key) {
        return this.options.containsKey(key);
    }

    /**
     * Add a new modifier option to the store
     *
     * @param option Option to add
     */

    public void addOption(ModifierOption option) {
        this.options.put(option.getKey(), option);
    }

    /**
     * Get the value of a modifier option
     *
     * @param key key of the option
     * @return The value of the option
     */

    public String get(String key) {
        return this.options.get(key).getValue();
    }
}
