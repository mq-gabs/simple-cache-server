package utils

import "fmt"

func FmtErr(err error, param any) error {
	return fmt.Errorf("%v: %v", err, param)
}
