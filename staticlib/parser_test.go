package staticlib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	tests := [][]interface{}{
		{"name: Michael", "name", "Michael", true},
		{`"quoted": "true"`, "quoted", "true", true},
		{"name:", "name", "", true},
		{"name", "", "", false},
		{"- Colorado", "", "Colorado", true},
		{`- "quoted"`, "", "quoted", true},
	}

	for i, test := range tests {
		line, _ := test[0].(string)
		key, value, err := parseLine(line)
		assert.Equal(t, test[1], key, "Test", i)
		assert.Equal(t, test[2], value, "Test", i)
		assert.Equal(t, test[3], err == nil)
	}
}

func TestParseSequenceGetElementByKey(t *testing.T) {
	seq := sequence{
		element{key: "first", value: "Michael"},
		element{key: "last", value: "Dwan"},
	}

	el, ok := seq.elForKey("first")
	assert.True(t, ok)
	assert.Equal(t, "Michael", el.value)
	el, ok = seq.elForKey("last")
	assert.True(t, ok)
	assert.Equal(t, "Dwan", el.value)
	el, ok = seq.elForKey("missing")
	assert.False(t, ok)
	assert.Equal(t, "", el.value)
}

func TestParseStringSliceForKeySeqValues(t *testing.T) {
	el := element{key: "list", sequence: sequence{
		element{value: "first"},
		element{value: "second"},
		element{value: "third"},
	}}
	assert.Equal(t, []string{"first", "second", "third"}, el.stringSliceForSeqValues())
}

func TestParseExtractNodes(t *testing.T) {
	reader := strings.NewReader(`

# a key & value
key: value

# a list
list:
  - a
  - b
# a map
map:

  m1: hello

`)
	nodes := extractNodes(reader)
	assert.Equal(t, 3, len(nodes))
	assert.Equal(t, "key: value", nodes[0].line)
	assert.Equal(t, 4, nodes[0].lineNo)
	assert.Equal(t, "list:", nodes[1].line)
	assert.Equal(t, 7, nodes[1].lineNo)
	assert.Equal(t, 2, len(nodes[1].subnodes))
	assert.Equal(t, "- a", nodes[1].subnodes[0].line)
	assert.Equal(t, 8, nodes[1].subnodes[0].lineNo)
	assert.Equal(t, "- b", nodes[1].subnodes[1].line)
	assert.Equal(t, 9, nodes[1].subnodes[1].lineNo)

	assert.Equal(t, "map:", nodes[2].line)
	assert.Equal(t, 11, nodes[2].lineNo)
	assert.Equal(t, 1, len(nodes[2].subnodes))
	assert.Equal(t, "m1: hello", nodes[2].subnodes[0].line)
	assert.Equal(t, 13, nodes[2].subnodes[0].lineNo)
}

func TestParse(t *testing.T) {
	reader := strings.NewReader(`
# this is a comment
name: Michael Dwan
company: Highrise
website: https://highrisehq.com
with_underscore: yes
with-hyphen: yes
with.dot: yes
list:
  # this is a comment
  - first
  - second
map:
  key1: value1
  key2: value2
	`)
	seq, err := parseConfig(reader)
	assert.Nil(t, err)
	assert.Equal(t, "Michael Dwan", seq.valForKey("name"))
	assert.Equal(t, "Highrise", seq.valForKey("company"))
	assert.Equal(t, "https://highrisehq.com", seq.valForKey("website"))
	assert.Equal(t, "yes", seq.valForKey("with_underscore"))
	assert.Equal(t, "yes", seq.valForKey("with-hyphen"))
	assert.Equal(t, "yes", seq.valForKey("with.dot"))
	el, ok := seq.elForKey("list")
	assert.True(t, ok)
	assert.Equal(t, "", el.value)
	assert.Equal(t, []string{"first", "second"}, el.stringSliceForSeqValues())

	el, ok = seq.elForKey("map")
	assert.True(t, ok)
	assert.Equal(t, "", el.value)
	assert.Equal(t, "key1", el.sequence[0].key)
	assert.Equal(t, "value1", el.sequence[0].value)
	assert.Equal(t, "key2", el.sequence[1].key)
	assert.Equal(t, "value2", el.sequence[1].value)

	el, ok = seq.elForKey("missing")
	assert.False(t, ok)
}
