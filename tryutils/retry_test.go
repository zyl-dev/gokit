package tryutils

import (
	"errors"
	"fmt"
	"log"
	"testing"
)

func TestTryExample(t *testing.T) {
	MaxRetries = 20
	SomeFunction := func() (string, error) {
		return "", nil
	}
	var value string
	err := Do(func(attempt int) (bool, error) {
		var err error
		value, err = SomeFunction()
		return attempt < 5, err // try 5 times
	})
	if err != nil {
		log.Fatalln("error:", err)
	}
	log.Fatalln("value:", value)
}
func TestTryExamplePanic(t *testing.T) {
	SomeFunction := func() (string, error) {
		panic("something went badly wrong")
	}
	var value string
	err := Do(func(attempt int) (retry bool, err error) {
		retry = attempt < 5 // try 5 times
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("panic: %v", r))
			}
		}()
		value, err = SomeFunction()
		return
	})
	if err != nil {
		//log.Fatalln("error:", err)
	}
	log.Fatalln("value:", value)
}
