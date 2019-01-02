package de.dieploegers.icarus.exceptions;

/**
 * Error processing event
 */
public class ProcessException extends Exception {
    private static final long serialVersionUID = -8336292388192367832L;

    public ProcessException(final Throwable cause) {
        super(cause);
    }
}
