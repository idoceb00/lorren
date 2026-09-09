// Package markdown provides low-level parsing helpers for markdown files with YAML frontmatter.
package markdown

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	// ErrNoFrontmatter is returned when the file does not open with a frontmatter delimiter.
	ErrNoFrontmatter = errors.New("no yaml frontmatter found")

	// ErrUnterminatedFrontmatter is returned when the opening delimiter has no matching closing one.
	ErrUnterminatedFrontmatter = errors.New("unterminated yaml frontmatter")
)

const delimiter = "---"

// SplitFrontmatter splits a markdown document into its YAML frontmatter and its body.
// Both are returned without the delimiter lines and without surrounding blank lines.
func SplitFrontmatter(data []byte) (frontmatter, body []byte, err error) {
	trimmed := bytes.TrimLeft(data, " \t\r\n")
	if !bytes.HasPrefix(trimmed, []byte(delimiter)) {
		return nil, nil, ErrNoFrontmatter
	}

	// Skip the opening delimiter line.
	rest := trimmed[len(delimiter):]
	nl := bytes.IndexByte(rest, '\n')
	if nl < 0 {
		return nil, nil, ErrUnterminatedFrontmatter
	}
	rest = rest[nl+1:]

	// Find the closing delimiter at the start of a line.
	closing := findDelimiterLine(rest)
	if closing < 0 {
		return nil, nil, ErrUnterminatedFrontmatter
	}

	frontmatter = bytes.TrimSpace(rest[:closing])

	after := rest[closing:]
	nl = bytes.IndexByte(after, '\n')
	if nl < 0 {
		// Closing delimiter on the last line: no body at all.
		return frontmatter, nil, nil
	}
	body = bytes.TrimSpace(after[nl+1:])

	return frontmatter, body, nil
}

// findDelimiterLine returns the offset of the first line that consists only of the delimiter, or -1.
func findDelimiterLine(data []byte) int {
	offset := 0
	for offset < len(data) {
		end := bytes.IndexByte(data[offset:], '\n')
		var line []byte
		if end < 0 {
			line = data[offset:]
		} else {
			line = data[offset : offset+end]
		}

		if string(bytes.TrimRight(line, " \t\r")) == delimiter {
			return offset
		}

		if end < 0 {
			return -1
		}
		offset += end + 1
	}
	return -1
}

// Join builds a markdown document from frontmatter and body. It is the inverse of SplitFrontmatter for well-formed input.
func Join(frontmatter, body []byte) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\n%s\n%s\n", delimiter, bytes.TrimSpace(frontmatter), delimiter)
	if b := bytes.TrimSpace(body); len(b) > 0 {
		buf.WriteString("\n")
		buf.Write(b)
		buf.WriteString("\n")
	}
	return buf.Bytes()
}
