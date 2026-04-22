package errorsUtil

import (
	"errors"
	"fmt"
)

type (
	wrap struct {
		err error
		msg string
	}
)

func Wrap(err error, format string, args ...any) error {
	if err == nil {
		panic("wrapping nil error")
	}

	return wrap{
		err: err,
		msg: fmt.Sprintf(format, args...),
	}
}

func New(format string, args ...any) error {
	return errors.New(fmt.Sprintf(format, args...))
}

func (w wrap) Error() string {
	return fmt.Sprintf("%v: %v", w.msg, w.err)
}

func (w wrap) Unwrap() error {
	return w.err
}
