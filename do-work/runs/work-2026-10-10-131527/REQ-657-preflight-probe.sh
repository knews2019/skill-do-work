#!/usr/bin/env bash
# REQ-657 pre-build baseline: the toolbox command tests the change sits beside
# (verb registration lock-in and the architecture verb, the closest existing
# toolbox verb) must pass at the base before the builder starts.
set -euo pipefail

repository_root="$(git rev-parse --show-toplevel)"
module_directory="$repository_root/skills/do-work/tools/do-work-cli"

test_output="$(go test -C "$module_directory" -count=1 -v \
  -run '^(TestHandlersRegisterCanonicalToolboxCommands|TestArchitecture.*|TestRemediationArchitecture.*)$' \
  ./internal/toolboxcommands/ 2>&1)" || {
  printf '%s\n' "$test_output"
  exit 1
}
printf '%s\n' "$test_output"

if grep -q -- '--- SKIP:' <<<"$test_output"; then
  printf 'REQ-657 preflight: a baseline test skipped; a skip is not a pass\n' >&2
  exit 1
fi
grep -q -- '--- PASS: TestHandlersRegisterCanonicalToolboxCommands' <<<"$test_output"
grep -q -- '--- PASS: TestArchitecturePublishUsesFirstFreeNumericBundle' <<<"$test_output"
printf 'REQ-657 preflight: baseline toolbox command tests pass\n'
