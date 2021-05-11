package core

import "fmt"

func ErrorOfNonexists(t string, target string) error {
	return fmt.Errorf("can't find %s: %s", t, target)
}

func ErrorOfInvalid(t string, target string) error {
	return fmt.Errorf("invalid %s: %s", t, target)
}
