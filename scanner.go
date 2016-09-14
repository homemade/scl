package scl

import (
	"bufio"
	"io"
	"strings"
)

type scannerTree []*scannerLine

type scanner struct {
	file   string
	reader io.Reader
	lines  scannerTree
}

func newScanner(reader io.Reader, filename ...string) *scanner {

	file := "<no file>"

	if len(filename) > 0 {
		file = filename[0]
	}

	s := scanner{
		file:   file,
		reader: reader,
		lines:  make(scannerTree, 0),
	}

	return &s
}

func (s *scanner) scan() (lines scannerTree, err error) {

	// Split to lines
	scanner := bufio.NewScanner(s.reader)
	scanner.Split(bufio.ScanLines)

	line_number := 0
	raw_lines := make(scannerTree, 0)

	for scanner.Scan() {
		line_number++

		text := strings.TrimRight(scanner.Text(), " \t{}")

		if text == "" {
			continue
		}

		raw_lines = append(raw_lines, newLine(s.file, line_number, 0, text))
	}

	// Make sure the first line has no indent
	if len(raw_lines) > 0 {
		index := 0
		s.indentLines(&index, raw_lines, &lines, raw_lines[0].content.indent())
	}

	return
}

func (s *scanner) indentLines(index *int, input scannerTree, output *scannerTree, indent int) {

	// Ends when there are no more lines
	if *index >= len(input) {
		return
	}

	var line_to_add *scannerLine

	for ; *index < len(input); *index++ {

		line_indent := input[*index].content.indent()

		if line_indent == indent {
			line_to_add = input[*index].branch()
			*output = append(*output, line_to_add)

		} else if line_indent > indent {
			s.indentLines(index, input, &line_to_add.children, line_indent)

		} else if line_indent < indent {
			*index--
			return
		}

	}

	return
}
