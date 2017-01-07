package de.dieploegers.icarus;

import de.dieploegers.icarus.modifier.Modifier;

/**
 * A pojo describing options needed by a modifier
 */
public class ModifierOption {

    /**
     * The key of this Modifier (used as a long option in the CLI)
     */

    private String key;

    /**
     * Wether this modifier expects a value (or if it just has to be set or not)
     */

    private boolean hasValue;

    /**
     * The value (set, when used in an optionstore)
     */

    private String value;

    /**
     * A description for this option (used in the usage information in the CLI)
     */

    private String description;

    /**
     * Create a new modifier option
     *
     * @param key key
     * @param description A short description
     * @param hasValue Does the option require a value?
     */

    public ModifierOption(String key, String description, boolean hasValue) {
        this.key = key;
        this.description = description;
        this.hasValue = hasValue;
        this.value = "";
    }

    /**
     * Create a new ModifierOption
     *
     * @param key Key
     * @param description Short description
     */

    public ModifierOption(String key, String description) {
        this.key = key;
        this.description = description;
        this.hasValue = false;
        this.value = "";
    }

    /**
     * Create a rather empty Modifier Option
     * @param key The key of the modifier
     */

    public ModifierOption(String key) {
        this.key = key;
    }

    public String getKey() {
        return key;
    }

    public void setKey(String key) {
        this.key = key;
    }

    public boolean getHasValue() {
        return hasValue;
    }

    public void setHasValue(boolean hasValue) {
        this.hasValue = hasValue;
    }

    public String getValue() {
        return value;
    }

    public void setValue(String value) {
        this.value = value;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    static public ModifierOption withValue(String key, String value) {
        ModifierOption option = new ModifierOption(key);
        option.setValue(value);
        return option;
    }

}
