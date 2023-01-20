package org.nzbhydra.mapping.newznab.json.caps;

import org.assertj.core.api.AbstractObjectAssert;
import org.assertj.core.util.Objects;

/**
 * Abstract base class for {@link CapsJsonRegistration} specific assertions - Generated by CustomAssertionGenerator.
 */
@jakarta.annotation.Generated(value = "assertj-assertions-generator")
public abstract class AbstractCapsJsonRegistrationAssert<S extends AbstractCapsJsonRegistrationAssert<S, A>, A extends CapsJsonRegistration> extends AbstractObjectAssert<S, A> {

    /**
     * Creates a new <code>{@link AbstractCapsJsonRegistrationAssert}</code> to make assertions on actual CapsJsonRegistration.
     *
     * @param actual the CapsJsonRegistration we want to make assertions on.
     */
    protected AbstractCapsJsonRegistrationAssert(A actual, Class<S> selfType) {
        super(actual, selfType);
    }

    /**
     * Verifies that the actual CapsJsonRegistration's attributes is equal to the given one.
     *
     * @param attributes the given attributes to compare the actual CapsJsonRegistration's attributes to.
     * @return this assertion object.
     * @throws AssertionError - if the actual CapsJsonRegistration's attributes is not equal to the given one.
     */
    public S hasAttributes(CapsJsonRegistrationAttributes attributes) {
        // check that actual CapsJsonRegistration we want to make assertions on is not null.
        isNotNull();

        // overrides the default error message with a more explicit one
        String assertjErrorMessage = "\nExpecting attributes of:\n  <%s>\nto be:\n  <%s>\nbut was:\n  <%s>";

        // null safe check
        CapsJsonRegistrationAttributes actualAttributes = actual.getAttributes();
        if (!Objects.areEqual(actualAttributes, attributes)) {
            failWithMessage(assertjErrorMessage, actual, attributes, actualAttributes);
        }

        // return the current assertion for method chaining
        return myself;
    }

}