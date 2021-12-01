package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	cases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "an empty string",
			text:     "",
			expected: nil,
		},
		{
			name:     "just a term",
			text:     "simple",
			expected: []string{"simple"},
		},
		{
			name:     "just a couple of terms",
			text:     "multiple terms",
			expected: []string{"multiple", "terms"},
		},
		{
			name:     "extranenous characters",
			text:     ":::@@@---...   multiple :::@@@---...  terms  :::@@@---...",
			expected: []string{"multiple", "terms"},
		},
		{
			name:     "a domain name",
			text:     "blargle.example.com",
			expected: []string{"blargle", "example", "com"},
		},
		{
			name:     "a hyphenated string",
			text:     "blargle-example-com",
			expected: []string{"blargle", "example", "com"},
		},
		{
			name:     "an image name",
			text:     "foo.com/blargle:flargle@sha1234",
			expected: []string{"foo", "com", "blargle", "flargle", "sha1234"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, Tokenize(c.text))
		})
	}
}
