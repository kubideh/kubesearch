package tokenizer

import (
	"bufio"
	"strings"
	"unicode"
	"unicode/utf8"

	"k8s.io/klog/v2"
)

// Tokenize uses the default Golang word scanner as a base, and it
// applies additional separators such as colons, dots, and hyphens,
// etc.
func Tokenize(text string) (results []string) {
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(scan)

	for scanner.Scan() {
		results = append(results, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		klog.Warningln("scanner error: ", err)
	}

	return
}

// scan is based on bufio.ScanWords.
func scan(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip anything leading that isn't a digit or letter.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			break
		}
	}
	// Scan until something other than a digit or letter, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) {
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
