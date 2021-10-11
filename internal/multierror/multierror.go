package multierror

import "strings"

type MultiError struct {
	Errors []error
}

func (err MultiError) Error() string {
	messages := make([]string, len(err.Errors))
	for i, innerError := range err.Errors {
		messages[i] = innerError.Error()
	}
	return strings.Join(messages, "|")
}

// Optional returns:
// - nil if the list of errors is empty
// - the first error if the list has one item
// - a MultiError instance if the list has multiple items
func Optional(errors []error) error {
	if len(errors) == 0 {
		return nil
	}

	if len(errors) == 1 {
		return errors[0]
	}

	return &MultiError{errors}
}
