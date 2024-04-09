package soundex

import (
	"testing"
)

func Test_Encoder_Encode(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{input: "", expected: ""},
		{input: "Robert", expected: "R163"},
		{input: "Rupert", expected: "R163"},
		{input: "Rubin", expected: "R150"},
		{input: "Ashcraft", expected: "A261"},
		{input: "Ashcroft", expected: "A261"},
		{input: "Tymczak", expected: "T522"},
		{input: "Pfister", expected: "P236"},
		{input: "Honeyman", expected: "H555"},
		{input: "Апельсин", expected: "0000"},
		{input: "Orange", expected: "O652"},
	}

	e := NewSoundexEncoder()
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			if c.expected != e.Encode(c.input) {
				t.Errorf("input: %s, expected: %s, got: %s", c.input, c.expected, e.Encode(c.input))
			}
		})
	}
}
