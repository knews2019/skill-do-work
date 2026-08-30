package settingshooks

import (
	"bytes"
	"strings"
	"testing"
)

const coreHookFragment = `{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/session-start.sh\""
          }
        ]
      }
    ]
  }
}
`

// Go's encoding/json sorts map keys. Marshalling a consumer's settings through a map would
// silently reorder unrelated user state on every install, which the byte-preservation
// requirement forbids and no content assertion would catch.
func TestComposeSettingsPreservesEveryUnrelatedKeyInItsOriginalOrder(t *testing.T) {
	settings := `{
  "zebra": 1,
  "alpha": {
    "nested-last": true,
    "nested-first": "keep"
  },
  "middle": [3, 2, 1]
}
`
	composed, err := ComposeSettings([]byte(settings), []byte(coreHookFragment))
	if err != nil {
		t.Fatalf("ComposeSettings: %v", err)
	}
	output := string(composed)
	for _, ordered := range [][2]string{
		{`"zebra"`, `"alpha"`},
		{`"alpha"`, `"middle"`},
		{`"nested-last"`, `"nested-first"`},
	} {
		if strings.Index(output, ordered[0]) > strings.Index(output, ordered[1]) {
			t.Errorf("key %s must still precede %s:\n%s", ordered[0], ordered[1], output)
		}
	}
	if !strings.Contains(output, `"middle": [`) || !strings.Contains(output, "3,") {
		t.Errorf("array order or numeric literals were rewritten:\n%s", output)
	}
	if !strings.HasSuffix(output, "}\n") {
		t.Errorf("composed settings must end with one trailing newline:\n%q", output)
	}
}

// The retired pipeline guard is removed from wherever it sits, but only that command object.
// The wrapper around it survives when it still holds a custom hook, and a wrapper that was
// already empty is preserved because nothing was removed from it.
func TestRetiredPipelineGuardRemovalKeepsEveryOtherEntry(t *testing.T) {
	settings := `{
  "custom": {"keep": [1, 2, 3]},
  "hooks": {
    "SessionStart": [
      {"hooks": [{"type": "command", "command": "echo custom-start"}]}
    ],
    "Stop": [
      {"hooks": [
        {"type": "command", "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/pipeline-guard.sh\""},
        {"type": "command", "command": "echo custom-stop"}
      ]},
      {"hooks": [
        {"type": "command", "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/pipeline-guard.sh\""}
      ]},
      {"matcher": "preserve-empty", "hooks": []},
      {"hooks": [{"type": "command", "command": "echo memory-capture"}]}
    ]
  }
}
`
	composed, err := ComposeSettings([]byte(settings), []byte(coreHookFragment))
	if err != nil {
		t.Fatalf("ComposeSettings: %v", err)
	}
	output := string(composed)
	if strings.Contains(output, "pipeline-guard.sh") {
		t.Errorf("the retired pipeline guard survived composition:\n%s", output)
	}
	if strings.Count(output, "echo custom-stop") != 1 {
		t.Errorf("a custom hook sharing the guard's wrapper was dropped:\n%s", output)
	}
	if strings.Count(output, "preserve-empty") != 1 {
		t.Errorf("a deliberately empty unrelated wrapper was dropped:\n%s", output)
	}
	if strings.Count(output, "echo memory-capture") != 1 {
		t.Errorf("an unrelated Stop wrapper was dropped:\n%s", output)
	}
	if strings.Count(output, "echo custom-start") != 1 {
		t.Errorf("a custom SessionStart hook was dropped:\n%s", output)
	}
	if strings.Count(output, "do-work/hooks/session-start.sh") != 1 {
		t.Errorf("the core SessionStart hook was not composed exactly once:\n%s", output)
	}
	if !strings.Contains(output, `"keep"`) {
		t.Errorf("an unrelated top-level key was dropped:\n%s", output)
	}
}

// A guard-only Stop array must lose the Stop key entirely rather than keep an empty array.
func TestStopEventIsDeletedWhenTheGuardWasItsOnlyContent(t *testing.T) {
	settings := `{
  "hooks": {
    "Stop": [
      {"hooks": [
        {"type": "command", "command": "bash \".claude/skills/do-work/hooks/pipeline-guard.sh\""}
      ]}
    ]
  }
}
`
	composed, err := ComposeSettings([]byte(settings), []byte(coreHookFragment))
	if err != nil {
		t.Fatalf("ComposeSettings: %v", err)
	}
	if strings.Contains(string(composed), `"Stop"`) {
		t.Errorf("an emptied Stop event was retained:\n%s", composed)
	}
}

// Reinstalling must be byte-idempotent, so a fragment entry already present is never
// appended a second time.
func TestFragmentEntriesAreAppendedOnlyOnce(t *testing.T) {
	firstPass, err := ComposeSettings([]byte("{}\n"), []byte(coreHookFragment))
	if err != nil {
		t.Fatalf("first ComposeSettings: %v", err)
	}
	secondPass, err := ComposeSettings(firstPass, []byte(coreHookFragment))
	if err != nil {
		t.Fatalf("second ComposeSettings: %v", err)
	}
	if string(firstPass) != string(secondPass) {
		t.Errorf("composition is not idempotent:\nfirst:\n%s\nsecond:\n%s", firstPass, secondPass)
	}
	if strings.Count(string(firstPass), "do-work/hooks/session-start.sh") != 1 {
		t.Errorf("the core hook was composed more than once:\n%s", firstPass)
	}
}

// jq — the branch this port replaces — emits raw UTF-8 rather than \uXXXX escapes, so a
// consumer's non-ASCII settings values keep their original bytes.
func TestNonAsciiValuesStayRawUtf8(t *testing.T) {
	composed, err := ComposeSettings([]byte(`{"greeting":"héllo→"}`), []byte(coreHookFragment))
	if err != nil {
		t.Fatalf("ComposeSettings: %v", err)
	}
	if !strings.Contains(string(composed), "héllo→") {
		t.Errorf("non-ASCII value was escaped:\n%s", composed)
	}
	if strings.Contains(string(composed), `\u00e9`) {
		t.Errorf("non-ASCII value was \\u-escaped:\n%s", composed)
	}
}

// jq — the branch this port preferred whenever it was installed — accepts a leading UTF-8
// byte-order mark and drops it, so a settings.json written by a Windows editor or a
// PowerShell redirect has always installed. Refusing it here would lock those consumers out
// of the installer entirely.
func TestLeadingByteOrderMarkIsStrippedLikeJq(t *testing.T) {
	settings := "\uFEFF" + `{
  "hooks": {
    "SessionStart": [
      {"hooks": [{"type": "command", "command": "echo custom-start"}]}
    ]
  }
}
`
	composed, err := ComposeSettings([]byte(settings), []byte(coreHookFragment))
	if err != nil {
		t.Fatalf("ComposeSettings: %v", err)
	}
	if bytes.Contains(composed, []byte("\uFEFF")) {
		t.Errorf("the byte-order mark came back in the re-encoded output:\n%q", composed)
	}
	if !bytes.HasPrefix(composed, []byte("{\n")) {
		t.Errorf("composed settings must start with the object itself:\n%q", composed)
	}
	output := string(composed)
	if strings.Count(output, "echo custom-start") != 1 {
		t.Errorf("the consumer's own hook was dropped:\n%s", output)
	}
	if strings.Count(output, "do-work/hooks/session-start.sh") != 1 {
		t.Errorf("the core SessionStart hook was not composed exactly once:\n%s", output)
	}
}

// HTML-significant characters are common in hook commands; escaping them would rewrite bytes
// the consumer never asked to change.
func TestHtmlSignificantCharactersAreNotEscaped(t *testing.T) {
	composed, err := ComposeSettings([]byte(`{"command":"a < b && c > d"}`), []byte(coreHookFragment))
	if err != nil {
		t.Fatalf("ComposeSettings: %v", err)
	}
	if !strings.Contains(string(composed), "a < b && c > d") {
		t.Errorf("HTML-significant characters were escaped:\n%s", composed)
	}
}

func TestMalformedSettingsAreRefusedWithoutProducingOutput(t *testing.T) {
	tests := []struct {
		name            string
		settings        string
		expectedMessage string
	}{
		{name: "not JSON at all", settings: "{not json", expectedMessage: "settings are not valid JSON"},
		{name: "root is not an object", settings: `[1, 2, 3]`, expectedMessage: "settings root must be an object"},
		{name: "hooks is not an object", settings: `{"hooks": []}`, expectedMessage: "settings hooks must be an object"},
		{name: "Stop is not an array", settings: `{"hooks": {"Stop": {}}}`, expectedMessage: "settings Stop hook event must be an array"},
		{name: "a hook event is not an array", settings: `{"hooks": {"SessionStart": {}}}`, expectedMessage: "settings hook event must be an array"},
		// A byte-order mark is one mark at position zero. Two of them, or one anywhere else, is
		// malformed input — stripping the leading one must not become lenient JSON parsing.
		{name: "doubled byte-order mark", settings: "\uFEFF\uFEFF{}", expectedMessage: "settings are not valid JSON"},
		{name: "byte-order mark at a non-zero offset", settings: "{\uFEFF\"hooks\": {}}", expectedMessage: "settings are not valid JSON"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			composed, err := ComposeSettings([]byte(test.settings), []byte(coreHookFragment))
			if err == nil {
				t.Fatalf("malformed settings were accepted: %s", composed)
			}
			if !strings.Contains(err.Error(), test.expectedMessage) {
				t.Errorf("error = %q, want it to contain %q", err.Error(), test.expectedMessage)
			}
			if composed != nil {
				t.Errorf("a refused composition produced output: %s", composed)
			}
		})
	}
}

// The installer creates settings.json from an empty object when the consumer has none.
func TestComposingOntoAnEmptyObjectProducesJustTheCoreHooks(t *testing.T) {
	composed, err := ComposeSettings([]byte("{}\n"), []byte(coreHookFragment))
	if err != nil {
		t.Fatalf("ComposeSettings: %v", err)
	}
	expected := `{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/session-start.sh\""
          }
        ]
      }
    ]
  }
}
`
	if string(composed) != expected {
		t.Errorf("composed settings =\n%s\nwant\n%s", composed, expected)
	}
}
