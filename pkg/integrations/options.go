package integrations

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Every integration follows the same pattern: build its options with
// New<Integration>Options()...Build(), then create the sender with
// New<Integration>(opts), which validates them.

var validate = validator.New()

// OptionsError reports options an integration cannot deliver with, such as a
// missing token or an invalid URL. Retrying the delivery does not help, so it
// reports DeliveryRejected.
type OptionsError struct {
	Integration string
	Err         error
}

func (err OptionsError) Error() string {
	return fmt.Sprintf("%s options: %v", err.Integration, err.Err)
}

func (err OptionsError) Unwrap() error { return err.Err }

// DeliveryRejected reports that the delivery cannot succeed as configured.
func (err OptionsError) DeliveryRejected() bool { return true }

func validateOptions[T any](integration string, opts *T) error {
	if opts == nil {
		return OptionsError{Integration: integration, Err: errors.New("options are required")}
	}
	if err := validate.Struct(opts); err != nil {
		return OptionsError{Integration: integration, Err: err}
	}
	return nil
}
