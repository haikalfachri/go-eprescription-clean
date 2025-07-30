package usecase_test

import "errors"

var ErrInternalServErr = errors.New("internal server error")

type TestCase struct {
	Name string
	Mock func()
	Res  interface{}
	Err  error
}
