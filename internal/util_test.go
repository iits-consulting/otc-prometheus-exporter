package internal

import (
	"os/user"
	"strings"
	"testing"
)

func TestExpandUserHome(t *testing.T) {
	NoTildePrefixTestCases := []string{
		"/hello",
		"/helloworld",
		"hello/~/world",
	}

	for _, inputValue := range NoTildePrefixTestCases {
		resultValue := expandUserHome(inputValue)
		if resultValue != inputValue {
			t.Errorf("Expected was \"%s\", but got \"%s\"", resultValue, inputValue)
		}
	}

	WithTildePrefixTestCases := []string{
		"~/hello",
		"~/helloworld",
	}

	usr, _ := user.Current()
	testUserHome := usr.HomeDir

	for _, inputValue := range WithTildePrefixTestCases {
		resultValue := expandUserHome(inputValue)
		if !strings.HasPrefix(resultValue, testUserHome) {
			t.Errorf("Expected was \"%s\", but got \"%s\"", resultValue, inputValue)
		}
	}

}
