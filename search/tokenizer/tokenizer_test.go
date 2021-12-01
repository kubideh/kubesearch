package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDNSSubdomainNamesTokenizer(t *testing.T) {
	cases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "just a term",
			text:     "simple",
			expected: []string{"simple"},
		},
		{
			name:     "hyphenated text",
			text:     "dns-subdomain-name",
			expected: []string{"dns", "subdomain", "name", "dns-subdomain-name"},
		},
		{
			name:     "dotted text",
			text:     "dns.subdomain.name",
			expected: []string{"dns", "subdomain", "name", "dns.subdomain.name"},
		},
		{
			name:     "dotted and hyphenated text",
			text:     "dns.sub-domain.name",
			expected: []string{"dns", "sub", "domain", "name", "dns.sub-domain.name"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, DNSSubdomainNamesTokenizer(c.text))
		})
	}
}
