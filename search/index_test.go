package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name     string
	text     string
	query    string
	found    bool
	expected []Posting
}

func TestIndexDNSSubdomainNames(t *testing.T) {
	cases := []testCase{
		{
			name:     "empty text",
			text:     "",
			query:    "",
			found:    false,
			expected: nil,
		},
		{
			name:     "simple text",
			text:     "simple",
			query:    "simple",
			found:    true,
			expected: []Posting{fakePosting()},
		},
		{
			name:     "long text is truncated",
			text:     longText(254),
			query:    longText(253),
			found:    true,
			expected: []Posting{fakePosting()},
		},
		{
			name:     "exact dotted",
			text:     "dns.subdomain.name",
			query:    "dns.subdomain.name",
			found:    true,
			expected: []Posting{fakePosting()},
		},
		{
			name:     "first part of dotted",
			text:     "dns.subdomain.name",
			query:    "dns",
			found:    true,
			expected: []Posting{fakePosting()},
		},
		{
			name:     "middle part of dotted",
			text:     "dns.subdomain.name",
			query:    "subdomain",
			found:    true,
			expected: []Posting{fakePosting()},
		},
		{
			name:     "last part of dotted",
			text:     "dns.subdomain.name",
			query:    "name",
			found:    true,
			expected: []Posting{fakePosting()},
		},
		{
			name:     "exact hyphenated",
			text:     "dns-subdomain-name",
			query:    "dns-subdomain-name",
			found:    true,
			expected: []Posting{fakePosting()},
		},
		{
			name:     "first part of hyphenated",
			text:     "dns-subdomain-name",
			query:    "dns",
			found:    true,
			expected: []Posting{fakePosting()},
		},
		{
			name:     "middle part of hyphenated",
			text:     "dns-subdomain-name",
			query:    "subdomain",
			found:    true,
			expected: []Posting{fakePosting()},
		},
		{
			name:     "last part of hyphenated",
			text:     "dns-subdomain-name",
			query:    "name",
			found:    true,
			expected: []Posting{fakePosting()},
		},
		{
			name:     "exact mixed",
			text:     "dns.sub-domain.name",
			query:    "dns.sub-domain.name",
			found:    true,
			expected: []Posting{fakePosting()},
		},
		{
			name:     "first part of mixed",
			text:     "dns.sub-domain.name",
			query:    "dns",
			found:    true,
			expected: []Posting{fakePosting()},
		},
		{
			name:     "middle part of mixed",
			text:     "dns.sub-domain.name",
			query:    "sub-domain",
			found:    true,
			expected: []Posting{fakePosting()},
		},
		{
			name:     "first part of middle split of mixed",
			text:     "dns.sub-domain.name",
			query:    "sub",
			found:    true,
			expected: []Posting{fakePosting()},
		},
		{
			name:     "last part of middle split of mixed",
			text:     "dns.sub-domain.name",
			query:    "domain",
			found:    true,
			expected: []Posting{fakePosting()},
		},
		{
			name:     "last part of mixed",
			text:     "dns.sub-domain.name",
			query:    "name",
			found:    true,
			expected: []Posting{fakePosting()},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			testIndexDNSSubdomainNames(t, c)
		})
	}
}

func testIndexDNSSubdomainNames(t *testing.T, c testCase) {
	index := NewIndex()

	IndexDNSSubdomainNames(index, c.text, fakePosting())

	result, found := index.Get(c.query)

	assert.Equal(t, c.found, found)
	assert.Equal(t, c.expected, result)
}

func fakePosting() Posting {
	return Posting{Key: "flargle", Kind: "blargle"}
}

func longText(size int) (text string) {
	for i := 0; i < size; i++ {
		text += "a"
	}
	return
}
