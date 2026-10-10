#!/usr/bin/env bash
# shellcheck disable=SC2016  # backticks in the patterns are literal Markdown, not expansions
# GREEN checks for REQ-691 (integrator test gate): the trace action and its reference exist and carry the
# required contract words, SKILL.md routes coverage questions to trace above the capture fallback and lists
# trace in the argument hint, help lists it, the board's read-only request-commits subcommand is wired and its
# two focused tests plus the existing correlation tests pass, and the two shipped-prose contract tests pass.
# Reads paths relative to the tree it runs in. Expected to FAIL at base (the trace files do not exist yet).
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
actions_dir="$root/skills/do-work/actions"
board_dir="$root/skills/do-work-board/tools/queue-kanban"
failures=0
fail() { echo "FAIL: $1"; failures=$((failures + 1)); }

for action_file in trace.md trace-reference.md; do
  [ -f "$actions_dir/$action_file" ] || { echo "FAIL: skills/do-work/actions/$action_file is missing"; exit 1; }
done
trace_text="$(<"$actions_dir/trace.md")"
reference_text="$(<"$actions_dir/trace-reference.md")"
trace_pair_text="$trace_text
$reference_text"
skill_text="$(<"$root/skills/do-work/SKILL.md")"
help_text="$(<"$actions_dir/help.md")"
board_main_text="$(<"$board_dir/main.go")"
board_prime_text="$(<"$board_dir/prime-do-kanban.md")"

# Action contract: untrusted source first, questions through clear-questions, the board evidence call, four verdicts.
grep -qF -- 'crew-members/prompt-injection.md' <<<"$trace_text" || fail "trace.md does not load crew-members/prompt-injection.md"
grep -qF -- 'crew-members/clear-questions.md' <<<"$trace_text" || fail "trace.md does not load crew-members/clear-questions.md"
grep -qF -- 'request-commits' <<<"$trace_pair_text" || fail "trace files do not call queue-kanban request-commits"
for verdict_word in '`complete`' '`partial`' '`not started`' '`unverified`'; do
  grep -qF -- "$verdict_word" <<<"$trace_pair_text" || fail "trace files lack the verdict $verdict_word"
done
grep -qF -- 'git status --porcelain' <<<"$trace_pair_text" || fail "trace files do not state the read-only-until-go check"

# Routing: a trace row above the capture fallback; capture-request still routes to capture; hint lists trace.
trace_row="$(grep -F -- '`./actions/trace.md`' <<<"$skill_text" | head -n 1 || true)"
[ -n "$trace_row" ] || fail "SKILL.md has no row routing to ./actions/trace.md"
for trigger_phrase in 'trace' 'is this captured' 'how much of it is implemented' 'already implemented'; do
  grep -qF -- "\`$trigger_phrase\`" <<<"$trace_row" || fail "SKILL.md trace row lacks the trigger: $trigger_phrase"
done
trace_row_line="$(grep -nF -- '`./actions/trace.md`' <<<"$skill_text" | head -n 1 | cut -d: -f1 || true)"
capture_row_line="$(grep -nF -- '`./actions/capture.md`' <<<"$skill_text" | head -n 1 | cut -d: -f1 || true)"
if [ -z "$trace_row_line" ] || [ -z "$capture_row_line" ]; then
  fail "SKILL.md: trace row or capture row not found"
elif [ "$trace_row_line" -ge "$capture_row_line" ]; then
  fail "SKILL.md: trace row is not above the capture fallback row"
fi
capture_row="$(grep -F -- '`./actions/capture.md`' <<<"$skill_text" | head -n 1 || true)"
grep -qF -- 'capture-request:' <<<"$capture_row" || fail "SKILL.md capture row no longer routes capture-request:"
hint_line="$(grep -- '^argument-hint:' <<<"$skill_text" || true)"
grep -qF -- 'trace <url|image|text|UR-NNN>' <<<"$hint_line" || fail "SKILL.md argument hint lacks trace <url|image|text|UR-NNN>"
grep -qF -- 'do-work trace' <<<"$help_text" || fail "help.md menu lacks do-work trace"

# Board subcommand wiring.
grep -qF -- 'case "request-commits":' <<<"$board_main_text" || fail "queue-kanban main.go does not dispatch request-commits"
board_prime_first_lines="$(head -n 5 <<<"$board_prime_text")"
grep -qF -- 'request-commits' <<<"$board_prime_first_lines" || fail "prime-do-kanban.md subcommand list lacks request-commits"

if [ "$failures" -gt 0 ]; then echo "$failures check(s) failed"; exit 1; fi

# Focused Go tests: the two new request-commits tests and the existing correlation tests the helper extraction touches.
go_tests=(
  TestRequestCommitsListsOnlyCommitsTheBoardCredits
  TestRequestCommitsRefusesAnArgumentThatIsNotARequestId
  TestRequestActivityAttributesACommitThatTouchesTheRequestFileWithoutAPrefix
  TestRequestActivityAttributesEveryPrefixTokenInASubject
  TestRequestPathPatternMatchesOnlyRequestFilesAndRunArtifacts
)
output="$(go test -C "$board_dir" -count=1 -v -run "^($(IFS='|'; echo "${go_tests[*]}"))\$" . 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for go_test in "${go_tests[@]}"; do
  grep -q -- "^--- PASS: $go_test " <<<"$output" || { printf '%s\n' "$output"; echo "$go_test did not pass (skipped or missing)"; exit 1; }
done

bash "$root/_dev/tests/shipped-package-reference-contract.sh"
bash "$root/_dev/tests/action-shell-blocks.sh"
echo "REQ-691 GREEN probe: all checks passed"
