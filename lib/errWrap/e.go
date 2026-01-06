package errWrap

import "fmt"

func Wrap(descrip string, err error) error {
	return fmt.Errorf("%s: %w", descrip, err)
}

func WrapIfErr(descrip string, err error) error {
	if err != nil {
		return Wrap(descrip, err)
	}
	return nil
}
