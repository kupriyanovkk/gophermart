package luhn

import "testing"

func TestValid(t *testing.T) {
	tests := []struct {
		input    int
		expected bool
	}{
		{4111111111111111, true},
		{5555555555554444, true},
		{1234567890123456, false},
		{33333, false},
		{0, true},
		{1, false},
	}

	for _, test := range tests {
		result := Valid(test.input)
		if result != test.expected {
			t.Errorf("For input %d, expected %t, but got %t", test.input, test.expected, result)
		}
	}
}
