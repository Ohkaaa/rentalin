package errs

import "fmt"

func WrapErr(where string, err error) error {
	return fmt.Errorf("%s error: %w", where, err)
}
