package handler

import (
	"errors"
)

var (
	errCannotReadKey    = errors.New("cannot read key")
	errCannotReadValue  = errors.New("cannot read value")
	errCannotReadAction = errors.New("cannot read action")
	errCannotReadTTL    = errors.New("cannot read time to live")

	errActionDoesNotExists = errors.New("action does not exists")
	errActionIsInvalid     = errors.New("action is invalid")

	errCannotSet    = errors.New("cannot set")
	errCannotErase  = errors.New("cannot erase")
	errCannotGetKey = errors.New("cannot get key")

	ErrCannotReadHeader                         = errors.New("cannot read header")
	ErrNumberOfReadBytesIsDifferentThanExpected = errors.New("number of read bytes is different than expected")
)
