#!/usr/bin/env bash
# Keep the exact generic guidance removed by REQ-168 from returning verbatim.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
failure_count=0

assert_phrase_absent() {
  local relative_path="$1"
  local deleted_phrase="$2"

  if grep -Fq -- "$deleted_phrase" "$repo_root/$relative_path"; then
    printf 'FAIL: %s restored generic guidance deleted by REQ-168: %s\n' \
      "$relative_path" "$deleted_phrase" >&2
    failure_count=$((failure_count + 1))
  fi
}

# One exact sentinel catches restoration of each deleted table without banning a future
# incident-backed Common Rationalizations or Red Flags section in the same action.
assert_phrase_absent \
  'skills/do-work/actions/verify-requests.md' \
  'The REQs cover everything — quick pass is fine'
assert_phrase_absent \
  'skills/do-work-toolbox/actions/code-review.md' \
  "This pattern is fine because it's used elsewhere in the codebase"
assert_phrase_absent \
  'skills/do-work-toolbox/actions/quick-wins.md' \
  '"No quick wins found" after scanning 3 files'
assert_phrase_absent \
  'skills/do-work-toolbox/actions/ui-review.md' \
  'The design looks fine to me'
assert_phrase_absent \
  'skills/do-work/actions/commit.md' \
  'No REQ matches — just commit everything together'
assert_phrase_absent \
  'skills/do-work/actions/commit.md' \
  'Single commit with >20 files'
assert_phrase_absent \
  'skills/do-work-toolbox/actions/inspect.md' \
  'This looks ready to commit'
assert_phrase_absent \
  'skills/do-work-toolbox/actions/inspect.md' \
  "Report produced without a \`git status\` / \`git diff\` reading actually happening"
assert_phrase_absent \
  'skills/do-work-toolbox/actions/inspect.md' \
  'Files listed with no REQ association attempt'
assert_phrase_absent \
  'skills/do-work-toolbox/actions/inspect.md' \
  'Debug artifacts (console.log, debugger, commented-out blocks)'

if [[ "$failure_count" -gt 0 ]]; then
  exit 1
fi

printf 'Defensive-surface exact deletion regressions passed.\n'
