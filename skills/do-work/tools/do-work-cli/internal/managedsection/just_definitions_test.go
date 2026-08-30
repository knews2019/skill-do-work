package managedsection

import (
	"sort"
	"strings"
	"testing"
)

func definitionNamesOf(data string) []string {
	names := make([]string, 0)
	for name := range JustDefinitionNames([]byte(data)) {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// The multiline-literal state machine decides which lines can define a recipe at all. Its
// boundaries are proved end to end today by contract-regressions.sh driving the whole
// installer; these cases pin them where a failure names the offending construct directly.
func TestJustDefinitionScannerClassifiesEveryLiteralForm(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected []string
	}{
		{
			name:     "a plain recipe and an alias",
			source:   "run-kanban:\n    echo x\nalias rk := run-kanban\n",
			expected: []string{"rk", "run-kanban"},
		},
		{
			name:     "an at-prefixed recipe keeps its bare name",
			source:   "@quiet-recipe:\n    echo x\n",
			expected: []string{"quiet-recipe"},
		},
		{
			name:     "an assignment is not a recipe",
			source:   "some-setting := \"value\"\n",
			expected: []string{},
		},
		{
			name:     "a comment line defines nothing",
			source:   "# run-kanban:\n",
			expected: []string{},
		},
		{
			name:     "an indented body line defines nothing",
			source:   "outer:\n    inner-looking:\n",
			expected: []string{"outer"},
		},
		{
			name:     "a colon inside a cooked string does not end the header",
			source:   "recipe param=\"a : b\":\n    echo x\n",
			expected: []string{"recipe"},
		},
		{
			name:     "an escaped quote inside a cooked string keeps the literal open",
			source:   "recipe param=\"a \\\" : b\":\n    echo x\n",
			expected: []string{"recipe"},
		},
		{
			name:     "a raw single-quoted string does not honour backslash escapes",
			source:   "recipe param='a \\':\n    echo x\n",
			expected: []string{"recipe"},
		},
		{
			name:     "a triple-quoted default spanning lines defers classification",
			source:   "recipe param=\"\"\"\nstill inside : not a recipe\n\"\"\":\n    echo x\n",
			expected: []string{"recipe"},
		},
		{
			name:     "a triple-single-quoted default spanning lines defers classification",
			source:   "recipe param='''\nstill inside : not a recipe\n''':\n    echo x\n",
			expected: []string{"recipe"},
		},
		{
			name:     "a backtick command spanning lines defers classification",
			source:   "recipe param=`\nstill inside : not a recipe\n`:\n    echo x\n",
			expected: []string{"recipe"},
		},
		{
			name:     "a triple-backtick command spanning lines defers classification",
			source:   "recipe param=```\nstill inside : not a recipe\n```:\n    echo x\n",
			expected: []string{"recipe"},
		},
		{
			name:     "a quadruple backtick is not a triple-backtick delimiter",
			source:   "thing := ````\nrun-kanban:\n    echo x\n",
			expected: []string{"run-kanban"},
		},
		{
			name:     "a trailing comment ends the scan without opening a literal",
			source:   "setting := 1 # a \" quote in a comment\nrun-kanban:\n",
			expected: []string{"run-kanban"},
		},
		{
			name:     "a BOM on line zero does not hide the first definition",
			source:   "\xef\xbb\xbfrun-kanban:\r\n    echo x\r\n",
			expected: []string{"run-kanban"},
		},
		{
			name:     "a BOM below line zero is not stripped",
			source:   "first:\n\xef\xbb\xbfrun-kanban:\n",
			expected: []string{"first"},
		},
		{
			name:     "lone carriage returns still separate definitions",
			source:   "first:\rrun-kanban:\r",
			expected: []string{"first", "run-kanban"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := definitionNamesOf(test.source)
			if strings.Join(actual, ",") != strings.Join(test.expected, ",") {
				t.Errorf("definitions = %v, want %v", actual, test.expected)
			}
		})
	}
}
