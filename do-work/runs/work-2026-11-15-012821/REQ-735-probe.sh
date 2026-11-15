#!/usr/bin/env bash
# REQ-735 GREEN probe: R1 and R2 live once, in the always-on communication-style.md; R3 stays only
# in the three background-agents.md copies (REQ-724 D-05); the F3 precedence clauses and the F5
# commit rule are shipped; lock-step copies match; Aliases and Examples are unchanged.
# Fails at the base commit by design (the moved lines and new clauses are absent there).
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"
base_commit=547d4cc5
do_work_crew="skills/do-work/crew-members"
toolbox_crew="skills/do-work-toolbox/crew-members"
knowledge_crew="skills/do-work-knowledge/crew-members"
failure_count=0

fail() {
  printf 'FAIL: %s\n' "$1"
  failure_count=$((failure_count + 1))
}

# A phrase must appear only in the expected files (sorted, space-separated), the changelog excluded.
expect_only_in() {
  local rule_label="$1" phrase="$2" expected_files="$3" matching_files
  matching_files="$(grep -rliF --exclude=CHANGELOG.md -- "$phrase" skills | LC_ALL=C sort | tr '\n' ' ' | sed 's/ $//' || true)"
  if [ "$matching_files" != "$expected_files" ]; then
    fail "$rule_label: '$phrase' expected only in $expected_files, found: ${matching_files:-nothing}"
  else
    printf 'PASS: %s only in %s\n' "$rule_label" "$expected_files"
  fi
}

background_copy_paths=("$knowledge_crew/background-agents.md" "$toolbox_crew/background-agents.md" "$do_work_crew/background-agents.md")
background_copies="${background_copy_paths[*]}"
expect_only_in R1 'unless the user says strict' "$do_work_crew/communication-style.md"
expect_only_in R2 'while redesigning' "$do_work_crew/communication-style.md"
expect_only_in R3 'cheapest model tier' "$background_copies"
expect_only_in F3-ask 'unless the action in use says how to ask' "$do_work_crew/communication-style.md"
expect_only_in F5 'change you did not make' "skills/do-work/actions/commit.md"

# R1 and R2 moved byte-identical: each old line appears exactly once in communication-style.md.
communication_text="$(cat "$do_work_crew/communication-style.md")"
for moved_line in \
  '- Read a number in the request (a size, a count, a time) as approximate unless the user says strict.' \
  '- Never remove existing UI while redesigning unless the user asked for it.'; do
  moved_count="$(grep -cxF -- "$moved_line" <<<"$communication_text" || true)"
  if [ "$moved_count" = "1" ]; then
    printf 'PASS: moved line present once: %s\n' "${moved_line:0:40}"
  else
    fail "moved line expected once in communication-style.md, found $moved_count: $moved_line"
  fi
done

# Lock-step: the model paragraph is identical in all three background-agents.md copies.
extract_model_paragraph() {
  awk '/^\*\*Pick each agent.s model by its work\.\*\*/{inside=1} inside && /^$/{exit} inside' "$1"
}
do_work_paragraph="$(extract_model_paragraph "$do_work_crew/background-agents.md")"
if [ -z "$do_work_paragraph" ]; then
  fail "model paragraph not found in $do_work_crew/background-agents.md"
fi
for copy_path in "$toolbox_crew/background-agents.md" "$knowledge_crew/background-agents.md"; do
  if [ "$(extract_model_paragraph "$copy_path")" = "$do_work_paragraph" ]; then
    printf 'PASS: model paragraph identical in %s\n' "$copy_path"
  else
    fail "model paragraph in $copy_path differs from the do-work copy"
  fi
done

# F3 model clause: the paragraph is hard-wrapped, so read it joined into one line.
joined_model_paragraph="$(tr '\n' ' ' <<<"$do_work_paragraph")"
if grep -qF 'keep the session model unless the action names another model' <<<"$joined_model_paragraph"; then
  printf 'PASS: F3-model clause in the model paragraph (identical in all three copies)\n'
else
  fail "F3-model: the model paragraph lacks 'keep the session model unless the action names another model'"
fi

# R3 stays out of the always-on file (REQ-724 D-05): no copy and no pointer, read across line wraps.
joined_communication_text="$(tr '\n' ' ' <<<"$communication_text")"
if grep -qiE 'model tier|background-agents|model by its work' <<<"$joined_communication_text"; then
  fail "R3: communication-style.md carries a copy of or pointer to the model rule"
else
  printf 'PASS: R3 has no copy or pointer in communication-style.md\n'
fi

# Aliases and Examples of communication-style.md are unchanged against the base commit.
extract_section() {
  awk -v heading="$1" '$0 == heading {inside=1; print; next} inside && /^## /{exit} inside'
}
for section_heading in '## Aliases' '## Examples'; do
  base_section="$(git show "$base_commit:$do_work_crew/communication-style.md" | extract_section "$section_heading")"
  current_section="$(extract_section "$section_heading" <<<"$communication_text")"
  if [ -n "$base_section" ] && [ "$base_section" = "$current_section" ]; then
    printf 'PASS: %s unchanged\n' "$section_heading"
  else
    fail "$section_heading of communication-style.md changed or is missing"
  fi
done

# No machine, vendor, model or tool name, and no citation of the maintainer doc, in the changed files.
changed_files_text="$(cat "$do_work_crew/communication-style.md" "$do_work_crew/coding-guardrails.md" \
  "$toolbox_crew/coding-guardrails.md" "${background_copy_paths[@]}" skills/do-work/actions/commit.md)"
if grep -qE 'zsh|AskUserQuestion|haiku|Sonnet|Opus|EnterWorktree|CLAUDE\.md' <<<"$changed_files_text"; then
  fail "a machine, vendor, model or tool name, or CLAUDE.md, is in a changed file"
else
  printf 'PASS: no machine, vendor, model or tool names in the seven changed files\n'
fi

if bash _dev/tests/shipped-package-reference-contract.sh >/dev/null 2>&1; then
  printf 'PASS: shipped-package-reference-contract\n'
else
  fail "shipped-package-reference-contract.sh failed"
fi

if [ "$failure_count" -ne 0 ]; then
  printf 'REQ-735 probe FAIL (%s)\n' "$failure_count"
  exit 1
fi
printf 'REQ-735 probe PASS\n'
