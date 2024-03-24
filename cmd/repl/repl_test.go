package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		{
			input: "hello world",
			want:  []string{"hello", "world"},
		},
		{
			input: "HELLO WORLD",
			want:  []string{"hello", "world"},
		},
	}
	for _, tc := range cases {
		got := cleanInput(tc.input)
		if len(tc.want) != len(got) {
			t.Errorf("cleanInput(%q) == %q, want %q", tc.input, got, tc.want)
			continue
		}
		for i := range got {
			actualWord := got[i]
			expectedWord := tc.want[i]
			if actualWord != expectedWord {
				t.Errorf("%v does not equal %v", actualWord, expectedWord)
			}
		}
	}
}
