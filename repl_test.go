package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "	hello 	world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "	Nay 	Kg Lah? ",
			expected: []string{"Nay", "Kg", "Lah?"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Expected %v, got %v", c.expected, actual)
		}

		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("Expected %v, got %v", c.expected[i], actual[i])
			}
		}
	}

}
