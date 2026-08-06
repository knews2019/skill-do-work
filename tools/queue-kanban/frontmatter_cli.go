package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// The `frontmatter` subcommand: read one field out of one REQ/UR file, with
// optional Schema Read Contract normalization and set-membership testing.
//
// WHY THIS EXISTS. Until it did, main.go's dispatch exposed seven subcommands
// and none of them took a file-and-field pair — so splitFrontmatter /
// parseFrontmatterFields / lenientFrontmatterFields had no caller outside the
// board's own tree walk, and every frontmatter read in actions/ was a hand
// reimplementation *by construction*. Three parsers already shipped (this one,
// plus awk in tools/checks/record-commit-hash.sh and blanked-req-scan.sh); the
// prose copies were the fourth onward. The `status` vocabulary alone is read at
// ~35 prose sites, and five separate Red Flags in actions/ document the same
// resulting bug: filtering on the literal `completed` and silently dropping
// `completed-with-issues`.
//
// THE FLOOR STILL WINS. This is an accelerator in the same class as `next-req`,
// `next-version`, and `now`: an action may name it as the PREFERRED source for
// something it can already obtain with shell primitives, gated on the binary
// being already built, with the prose procedure documented as the fallback.
// Nothing may build the tool to read a field, and no action may lose its floor
// path — `do-work board` is the only capability allowed to *need* a compiler,
// which is why `actions/board.md` precondition-checks `go` and degrades instead
// of blocking. `actions/work-reference.md` → Timestamp rule states the same
// already-built-and-falls-back shape for `now`.
//
// READ-ONLY, deliberately. The tool has exactly two write surfaces — the board's
// Testing view and `next-version` — and a third may only be added by amending the
// sentence that states the count, in the same commit. This command adds none:
// there is no `set` verb, and adding one is out of scope rather than merely
// unimplemented.

// frontmatterFieldSets maps the --in-set names to the canonical status sets in
// actions/work-reference.md's Schema Read Contract. Membership is delegated to
// the same predicates the board buckets with, so this command cannot drift into
// a second definition of "finished".
var frontmatterFieldSets = map[string]func(string) bool{
	"terminal-success":  isCompletedStatus,
	"terminal-resolved": isTerminalResolvedStatus,
}

// readFrontmatterField returns one raw frontmatter field's value from a REQ/UR
// file, whether the field was present, and an error for an input that is not a
// usable file.
//
// Present-but-empty scalars and absent fields are both reported as found=false:
// a caller asking for a scalar cannot act on either, and the contract already
// treats an empty optional field as absent. An explicitly empty YAML sequence is
// different: it is present and renders as empty stdout at exit 0, so callers can
// distinguish `depends_on: []` from a missing field without receiving Go's `[]`
// debug representation as a false value.
//
// A file with NO frontmatter block is an error rather than an absent field. The
// caller has to be able to tell "this REQ does not set domain" from "this is not
// a REQ" — the first is routine, the second is a broken input, and collapsing
// them is how a typo'd path reads as a queue full of unset fields.
func readFrontmatterField(filePath string, fieldName string) (string, bool, error) {
	fileBytes, readError := os.ReadFile(filePath)
	if readError != nil {
		return "", false, readError
	}
	yamlText, _, _, hasFrontmatter := splitFrontmatter(string(fileBytes))
	if !hasFrontmatter {
		return "", false, fmt.Errorf("%s has no frontmatter block", filePath)
	}
	frontmatterFields, parseError := parseFrontmatterFields(yamlText)
	if parseError != nil {
		// Strict YAML failed. The lenient scan is the whole reason a prose
		// reimplementation cannot match this: it recovers the flat top-level
		// fields from a block one bad line would otherwise take down.
		frontmatterFields = lenientFrontmatterFields(yamlText)
	}
	rawValue, isPresent := frontmatterFields[fieldName]
	if !isPresent {
		return "", false, nil
	}
	if _, isSequence := rawValue.([]any); isSequence {
		return strings.Join(coerceToStringList(rawValue), "\n"), true, nil
	}
	scalarValue := strings.TrimSpace(coerceScalarToString(rawValue))
	if scalarValue == "" {
		return "", false, nil
	}
	return scalarValue, true, nil
}

// frontmatterCommandArguments is the parsed argument set, split out so the parse
// can be asserted directly — runFrontmatterCommand's callers exit on its code,
// and an inline parse is how next-version shipped discarding its flags.
type frontmatterCommandArguments struct {
	FilePath   string
	FieldName  string
	Normalize  bool
	InSetName  string
	InSetGiven bool
	UsageError string
}

// parseFrontmatterCommandArguments accepts `get <file> <field>` plus --normalize
// and --in-set in any order after the positionals. Leftover tokens are rejected
// rather than ignored, matching every other subcommand in this tool.
func parseFrontmatterCommandArguments(args []string) frontmatterCommandArguments {
	if len(args) == 0 {
		return frontmatterCommandArguments{UsageError: "missing verb (want: get)"}
	}
	if args[0] != "get" {
		return frontmatterCommandArguments{UsageError: fmt.Sprintf("unknown verb %q (want: get)", args[0])}
	}
	remaining := args[1:]
	var positionals []string
	parsed := frontmatterCommandArguments{}
	for index := 0; index < len(remaining); index++ {
		switch remaining[index] {
		case "--normalize":
			parsed.Normalize = true
		case "--in-set":
			if index+1 >= len(remaining) {
				parsed.UsageError = "--in-set needs a set name"
				return parsed
			}
			index++
			parsed.InSetGiven = true
			parsed.InSetName = remaining[index]
		default:
			if strings.HasPrefix(remaining[index], "-") {
				parsed.UsageError = fmt.Sprintf("unknown flag %q", remaining[index])
				return parsed
			}
			positionals = append(positionals, remaining[index])
		}
	}
	if len(positionals) < 2 {
		parsed.UsageError = "want: frontmatter get <file> <field>"
		return parsed
	}
	if len(positionals) > 2 {
		parsed.UsageError = fmt.Sprintf("unexpected extra argument %q", positionals[2])
		return parsed
	}
	parsed.FilePath = positionals[0]
	parsed.FieldName = positionals[1]
	if parsed.InSetGiven {
		if parsed.FieldName != "status" {
			parsed.UsageError = fmt.Sprintf(
				"--in-set %s applies only to the status field, not %s",
				parsed.InSetName, parsed.FieldName)
			return parsed
		}
		if parsed.InSetName == "" {
			parsed.UsageError = "--in-set needs a non-empty set name"
			return parsed
		}
		if _, isKnownSet := frontmatterFieldSets[parsed.InSetName]; !isKnownSet {
			parsed.UsageError = fmt.Sprintf("unknown --in-set %q (want: terminal-success | terminal-resolved)", parsed.InSetName)
		}
	}
	return parsed
}

// runFrontmatterCommand executes the subcommand against the given writers and
// returns the process exit code. Writers are injected so the whole command is
// testable without a subprocess.
//
// Exit codes: 0 success (or set membership satisfied), 1 the field is absent or
// membership is not satisfied, 2 a usage error or an unreadable file. The
// value-versus-diagnostic split is load-bearing — the value goes to stdout and
// nothing else does, so a caller can capture stdout as a single clean value even
// when the contract emits a warning.
func runFrontmatterCommand(args []string, standardOut io.Writer, standardErr io.Writer) int {
	parsed := parseFrontmatterCommandArguments(args)
	if parsed.UsageError != "" {
		fmt.Fprintf(standardErr, "queue-kanban frontmatter: %s\n", parsed.UsageError)
		return 2
	}

	rawValue, isPresent, readError := readFrontmatterField(parsed.FilePath, parsed.FieldName)
	if readError != nil {
		fmt.Fprintf(standardErr, "queue-kanban frontmatter: %v\n", readError)
		return 2
	}

	// An absent field is a routine answer, not a usage error: exit 1 with empty
	// stdout so `if value=$(... ); then` reads naturally at a call site.
	if !isPresent {
		fmt.Fprintf(standardErr, "queue-kanban frontmatter: %s has no %s field\n", parsed.FilePath, parsed.FieldName)
		return 1
	}

	// A field the contract does not govern has nothing to normalize against, so
	// --normalize is a no-op on it rather than a violation to report: stdout gets
	// the value, stderr stays empty, and the call is observably identical to the
	// same get without the flag. Before REQ-118 the branch below fired here too,
	// because isKnownSchemaFieldValue returns false for a contract-less field just
	// as it does for a bad value — so every `--normalize created_at` printed a
	// warning claiming the timestamp was "not recognized", contradicting the
	// contract's own statement that such fields are outside it and read verbatim.
	//
	resolvedValue := rawValue
	if !hasSchemaFieldContract(parsed.FieldName) {
		writeFrontmatterValue(standardOut, resolvedValue)
		return 0
	}
	if parsed.Normalize || parsed.InSetGiven {
		// --in-set normalizes too: a membership test on a raw `done` has to
		// resolve the alias before asking, or the alias map the contract
		// promises would apply to `get` and not to the set check.
		normalizedValue := normalizeSchemaField(parsed.FieldName, rawValue)
		if !isKnownSchemaFieldValue(parsed.FieldName, normalizedValue) {
			// Every field with a contract row warns here, `status` and
			// `testing_status` included. An earlier version force-suppressed the
			// warning for those two because they are not alias-map rows, which
			// left a typo'd status printing to stdout at exit 0 with no feedback
			// — the exact path this command exists to replace.
			fmt.Fprintf(standardErr, "%s\n", schemaFieldWarningText(parsed.FieldName, rawValue))
		}
		// Prefer the contract's documented default; when a field has none
		// (`status`, `testing_status`, `route`), print what was actually found
		// rather than inventing a substitute the contract declined to define.
		resolvedValue = normalizedValue
		if defaultedValue, _ := resolveSchemaField(parsed.FieldName, rawValue); defaultedValue != "" {
			resolvedValue = defaultedValue
		}
	}

	if parsed.InSetGiven {
		isMember := frontmatterFieldSets[parsed.InSetName]
		if isMember(resolvedValue) {
			return 0
		}
		return 1
	}

	writeFrontmatterValue(standardOut, resolvedValue)
	return 0
}

// writeFrontmatterValue prints scalar values and sequence items one per line.
// An empty sequence deliberately emits no byte at all; a captured value then
// answers false to `test -n` instead of becoming the misleading string `[]`.
func writeFrontmatterValue(standardOut io.Writer, value string) {
	if value == "" {
		return
	}
	fmt.Fprintf(standardOut, "%s\n", value)
}
