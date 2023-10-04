package internal

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestResolveOtcShortHandNamespace(t *testing.T) {
	type testcase struct {
		name     string
		input    []string
		expected []string
	}
	test_cases := []testcase{
		{name: "one valid element", input: []string{"ECS"}, expected: []string{"SYS.ECS"}},
		{name: "one valid element, one unknown element", input: []string{"ECS", "UNKNOWN-ELEMENT"}, expected: []string{"SYS.ECS", "UNKNOWN-ELEMENT"}},
		{name: "one resolved element", input: []string{"SYS.ECS"}, expected: []string{"SYS.ECS"}},
	}

	for _, tc := range test_cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ResolveOtcShortHandNamespace(tc.input)
			if !slices.Equal(actual, tc.expected) {
				t.Error("The result differs from the expected value")
			}
		})
	}

}
