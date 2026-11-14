#!/usr/bin/env bash
# REQ-732 GREEN probe: the staged-skills contract accepts root docs/run-prompts/ and still
# refuses any other new folder under root docs/. It runs the contract's own docs checks (the
# export-ignore loop and the legacy-root scan), cut out of the real script by their text
# markers, so it tests the shipped lines in under a second instead of the whole 40 s contract.
# Expected at base: FAIL (the export-ignore line and the scan exclusion do not exist yet).
set -euo pipefail

tree_root="$(git rev-parse --show-toplevel)"
contract_path="$tree_root/_dev/tests/staged-skills-contract.sh"
probe_scratch="$(mktemp -d)"
trap 'rm -rf "$probe_scratch"' EXIT

probe_fail() {
  printf 'REQ-732 probe FAIL: %s\n' "$1" >&2
  exit 1
}

# The segment runs from the maintainer-docs export-ignore loop up to the board file list.
docs_checks_segment="$probe_scratch/docs-checks.sh"
awk '/^for maintainer_docs_path in /{copying=1} /^board_files=\(/{copying=0} copying' \
  "$contract_path" >"$docs_checks_segment"
segment_text="$(cat "$docs_checks_segment")"
grep -q 'docs/run-prompts' <<<"$segment_text" || probe_fail "the contract's docs checks do not name docs/run-prompts"
grep -q 'legacy_runtime_paths' <<<"$segment_text" || probe_fail "the segment markers no longer cut out the legacy-root scan"

# Run the cut-out checks against one tree; print the failure count, then the FAIL lines.
run_docs_checks() {
  local checked_root="$1"
  bash -c '
    set -euo pipefail
    repo_root="$1"
    fail_count=0
    fail() { printf "FAIL: %s\n" "$1"; fail_count=$((fail_count + 1)); }
    source "$2"
    printf "fail_count=%s\n" "$fail_count"
  ' docs-checks "$checked_root" "$docs_checks_segment"
}

# 1. The real tree, with the moved prompts present, passes.
[ ! -e "$tree_root/_dev/run-prompts" ] || probe_fail "_dev/run-prompts still exists"
moved_prompt_count="$(git -C "$tree_root" ls-files -- 'docs/run-prompts/*.md' | wc -l | tr -d ' ')"
[ "$moved_prompt_count" -ge 6 ] || probe_fail "docs/run-prompts holds $moved_prompt_count tracked prompts, expected at least 6"
attribute_line="$(git -C "$tree_root" check-attr export-ignore -- docs/run-prompts)"
grep -q 'export-ignore: set' <<<"$attribute_line" || probe_fail "docs/run-prompts is not export-ignore: $attribute_line"
real_tree_output="$(run_docs_checks "$tree_root")"
grep -q '^fail_count=0$' <<<"$real_tree_output" || probe_fail "real tree: $real_tree_output"

# 2. A fixture tree with only docs/run-prompts/probe.md passes.
fixture_root="$probe_scratch/fixture"
mkdir -p "$fixture_root/docs/run-prompts"
git -C "$fixture_root" init -q
cp "$tree_root/.gitattributes" "$fixture_root/.gitattributes"
printf 'probe\n' >"$fixture_root/docs/run-prompts/probe.md"
allowed_output="$(run_docs_checks "$fixture_root")"
grep -q '^fail_count=0$' <<<"$allowed_output" || probe_fail "docs/run-prompts/probe.md was refused: $allowed_output"

# 3. The same fixture plus docs/other/probe.md fails and names the file (the negative case,
#    so a scan that allows everything cannot pass this probe).
mkdir -p "$fixture_root/docs/other"
printf 'probe\n' >"$fixture_root/docs/other/probe.md"
refused_output="$(run_docs_checks "$fixture_root")"
grep -q '^fail_count=0$' <<<"$refused_output" && probe_fail "docs/other/probe.md was accepted"
grep -q 'legacy root runtime must be retired at modular cutover: docs (.*docs/other/probe.md)' <<<"$refused_output" \
  || probe_fail "docs/other/probe.md failed without the legacy-root message: $refused_output"

printf 'REQ-732 probe PASS\n'
