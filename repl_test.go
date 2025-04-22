package main

import (
	"testing"
)


func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: " hello     world ",
			expected: []string{"hello", "world"},
		},
		{
			input: " WINDBREAKERS            IS            LOWKEY             GOOD ",
			expected: []string{"windbreakers", "is", "lowkey", "good"},
		},
		{
			input: "                  WhY              ",
			expected: []string{"why"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("length of slice does not match expected length: got %d, expected %d (input: %q)", len(actual), len(c.expected), c.input)
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("word in slice does not match expected: got %q, expected %q (input: %q)", word, expectedWord, c.input)
			}
		}
	}

}

