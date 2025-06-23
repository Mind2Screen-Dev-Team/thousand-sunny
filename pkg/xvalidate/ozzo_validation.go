package xvalidate

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Errors map[string]any

func (e Errors) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func WrapperValidation(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(validation.InternalError); ok {
		return e.InternalError()
	}

	if e, ok := err.(validation.Errors); ok {
		err = e.Filter()
	}

	var errs Errors
	if es, ok := err.(validation.Errors); ok {
		errs = make(Errors)
		for field, validationErr := range es {
			errs[field] = UnwrapErrorJSON(validationErr.Error())
		}

		return errs
	}

	return nil
}

func UnwrapErrorJSON(errStr string) any {
	var result map[string]any

	if err := json.Unmarshal([]byte(errStr), &result); err == nil {
		return result
	}

	return errStr
}

func IsErrors(err error) (errs Errors, ok bool) {
	errs, ok = err.(Errors)
	return
}
