package apierr

import "fmt"

// ===== VALIDATION ERROR MESSAGES =====

const emptyFieldErrorMessage = "Field is required."
const minCharLengthErrorMessage = "Must be at least %d characters long."
const maxCharLengthErrorMessage = "Must not be longer than %d characters."

func EmptyFieldMessage() string {
	return emptyFieldErrorMessage
}

func MinCharLengthMessage(min int) string {
	return fmt.Sprintf(minCharLengthErrorMessage, min)
}

func MaxCharLengthMessage(max int) string {
	return fmt.Sprintf(maxCharLengthErrorMessage, max)
}
