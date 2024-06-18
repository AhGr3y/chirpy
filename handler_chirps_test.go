package main

import (
	"testing"
)

func TestCleanBody(t *testing.T) {

	profanities := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "This is a normal sentence.",
			expected: "This is a normal sentence.",
		},
		{
			input:    "This is a kerfuffle sentence.",
			expected: "This is a **** sentence.",
		},
		{
			input:    "This is a sharbert? sentence.",
			expected: "This is a sharbert? sentence.",
		},
	}

	for _, c := range cases {
		cleanBody := cleanBody(c.input, profanities)

		if cleanBody != c.expected {
			t.Errorf("%v != %v", cleanBody, c.expected)
		}
	}
}
