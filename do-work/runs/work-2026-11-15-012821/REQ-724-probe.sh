#!/usr/bin/env bash
# REQ-724 GREEN probe: the moved standing rules are shipped exactly once, in the planned homes,
# with no machine, vendor or tool name. Fails at the base commit by design (the phrases are absent).
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"
crew_directory="skills/do-work/crew-members"
failure_count=0

fail() {
  printf 'FAIL: %s\n' "$1"
  failure_count=$((failure_count + 1))
}

# Each rule-group phrase must appear in exactly one shipped file (the changelog mirror excluded)
# and that file must be the planned home.
expect_single_home() {
  local rule_group="$1" phrase="$2" expected_file="$3" matching_files
  matching_files="$(grep -rliF --exclude=CHANGELOG.md -- "$phrase" skills || true)"
  if [ "$matching_files" != "$expected_file" ]; then
    fail "$rule_group: '$phrase' expected only in $expected_file, found: ${matching_files:-nothing}"
  else
    printf 'PASS: %s found once in %s\n' "$rule_group" "$expected_file"
  fi
}

expect_single_home G1 'every status word' "$crew_directory/communication-style.md"
expect_single_home G1 'big picture in one short paragraph' "$crew_directory/communication-style.md"
expect_single_home G2 'prose list of questions' "$crew_directory/communication-style.md"
expect_single_home G3 'login shell' "$crew_directory/communication-style.md"
expect_single_home G3 'one-line launcher' "$crew_directory/communication-style.md"
expect_single_home G4 'A bug report means' "$crew_directory/communication-style.md"
expect_single_home G4 'unless the user says strict' "$crew_directory/coding-guardrails.md"
expect_single_home G4 'while redesigning' "$crew_directory/coding-guardrails.md"
expect_single_home G5 '20-minute stall check' "$crew_directory/communication-style.md"
expect_single_home G7 'cheapest model tier' "$crew_directory/background-agents.md"

# G2 option shape lives in the question-wording file.
clear_questions_text="$(cat "$crew_directory/clear-questions.md")"
if grep -qiF 'recommended option first' <<<"$clear_questions_text"; then
  printf 'PASS: G2 option shape in clear-questions.md\n'
else
  fail "G2: clear-questions.md lacks 'recommended option first'"
fi

# G5 names no path of its own: it cites the folder-structure action across the package boundary.
communication_text="$(cat "$crew_directory/communication-style.md")"
if grep -qF '../../do-work-toolbox/actions/folder-structure.md' <<<"$communication_text"; then
  printf 'PASS: G5 cites folder-structure.md\n'
else
  fail "G5: communication-style.md does not cite ../../do-work-toolbox/actions/folder-structure.md"
fi
if grep -qF 'docs/run-prompts' <<<"$communication_text"; then
  fail "G5: communication-style.md writes the run-prompts path itself"
fi

# No machine, vendor or tool name in the four changed files.
changed_files_text="$(cat "$crew_directory/communication-style.md" "$crew_directory/clear-questions.md" "$crew_directory/coding-guardrails.md" "$crew_directory/background-agents.md")"
if grep -qE 'zsh|AskUserQuestion|haiku|EnterWorktree' <<<"$changed_files_text"; then
  fail "a machine, vendor or tool name (zsh, AskUserQuestion, haiku, EnterWorktree) is in a changed crew-member file"
else
  printf 'PASS: no machine or tool names in the four crew-member files\n'
fi

# The new cross-package citation must resolve in both topologies.
if bash _dev/tests/shipped-package-reference-contract.sh >/dev/null 2>&1; then
  printf 'PASS: shipped-package-reference-contract\n'
else
  fail "shipped-package-reference-contract.sh failed"
fi

if [ "$failure_count" -ne 0 ]; then
  printf 'REQ-724 probe FAIL (%s)\n' "$failure_count"
  exit 1
fi
printf 'REQ-724 probe PASS\n'
