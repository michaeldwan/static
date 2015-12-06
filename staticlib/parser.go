package staticlib

// This parses the simplified yaml-like config format. This is used instead of
// a full yaml parser in order to preserve to-to-bottom priority since the Go
// runtime randomizes map iteration order. This is also much more forgiving
// with quotes and special characters often seen in a static.yml file.

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type element struct {
	key      string
	value    string
	sequence sequence
}

type sequence []element

func (s sequence) elForKey(key string) (element, bool) {
	for _, e := range s {
		if e.key == key {
			return e, true
		}
	}
	return element{}, false
}

func (s sequence) valForKey(key string) string {
	if el, ok := s.elForKey(key); ok {
		return el.value
	}
	return ""
}

func (s sequence) seqForKey(key string) sequence {
	if el, ok := s.elForKey(key); ok {
		return el.sequence
	}
	return nil
}

func (el element) stringSliceForSeqValues() []string {
	var out []string
	for _, subEl := range el.sequence {
		out = append(out, subEl.value)
	}
	return out
}

func parseConfig(reader io.Reader) (sequence, error) {
	nodes := extractNodes(reader)
	var sequence sequence
	for _, node := range nodes {
		el, err := parseNode(*node)
		if err != nil {
			return sequence, err
		}
		sequence = append(sequence, el)
	}
	return sequence, nil
}

type node struct {
	lineNo   int
	line     string
	subnodes []node
}

type parseError struct {
	node   node
	reason string
}

func (pe parseError) Error() string {
	return fmt.Sprintf("%s (line %d near '%s')", pe.reason, pe.node.lineNo, pe.node.line)
}

func extractNodes(reader io.Reader) []*node {
	var nodes []*node
	scanner := bufio.NewScanner(reader)
	for i := 0; scanner.Scan(); i++ {
		line := strings.TrimRightFunc(scanner.Text(), unicode.IsSpace)
		if len(line) == 0 || isCommentLine(line) {
			continue
		}
		node := node{
			line:   strings.TrimSpace(line),
			lineNo: i + 1,
		}
		if strings.HasPrefix(line, "  ") {
			parent := nodes[len(nodes)-1]
			parent.subnodes = append(parent.subnodes, node)
		} else {
			nodes = append(nodes, &node)
		}
	}
	return nodes
}

func isCommentLine(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "#")
}

func parseNode(node node) (element, error) {
	el := element{}
	key, value, err := parseLine(node.line)
	if err != nil {
		return el, parseError{node, err.Error()}
	}
	el.key = key
	el.value = value
	if len(node.subnodes) > 0 {
		if el.value != "" {
			return el, parseError{node, "keys with values cannot be followed by a list"}
		}
		for _, subnode := range node.subnodes {
			seqEl, err := parseNode(subnode)
			if err != nil {
				return el, parseError{node, err.Error()}
			}
			el.sequence = append(el.sequence, seqEl)
		}
	}
	return el, nil
}

func parseLine(line string) (string, string, error) {
	if strings.HasPrefix(line, "-") {
		return "", cleanInput(line[1:]), nil
	}
	return parseKeyValueLine(line)
}

func parseKeyValueLine(line string) (string, string, error) {
	if !strings.Contains(line, ":") {
		return "", "", fmt.Errorf("expected 'key: [value]'")
	}
	match := strings.SplitN(line, ":", 2)
	key, val := cleanInput(match[0]), ""
	if len(match) == 2 {
		val = cleanInput(match[1])
	}
	return key, val, nil
}

func cleanInput(str string) string {
	str = strings.TrimSpace(str)
	if unquoted, err := strconv.Unquote(str); err == nil {
		return unquoted
	}
	return str
}
