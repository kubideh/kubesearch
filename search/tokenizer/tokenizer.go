package tokenizer

import (
	"bufio"
	"strings"
	"unicode/utf8"

	"k8s.io/klog/v2"
)

// Tokenize follows the rules for naming
// objects in Kubernetes (https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names).
// In addition to tokenizing on hyphens or dots, the exact name
// is also returned as the first token. For example, for the name
// `dns.sub-domain.name`, the following tokens are returned:
// `dns`, `sub`, `domain`, `name`, and `dns.sub-domain.name`.
func Tokenize(text string) (results []string) {
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(scan)

	for scanner.Scan() {
		results = append(results, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		klog.Warningln("scanner error: ", err)
	}

	if len(results) > 1 {
		results = append(results, text)
	}

	return
}

// scan is a split function for a Scanner that returns UTF-8 tokens
// split on dots or hyphens. This algorithm is taken from bufio.ScanWords.
func scan(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if r != '.' && r != '-' {
			break
		}
	}

	// Scan until dot or hyphen, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if r == '.' || r == '-' {
			return i + width, data[start:i], nil
		}
	}

	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	// Request more data.
	return start, nil, nil
}
