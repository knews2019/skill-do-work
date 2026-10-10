#!/usr/bin/env bash
# REQ-696 GREEN probe: the five stale wording sites now carry the reviewers' text.
# Reads paths from the tree it runs in, so it checks a builder worktree or main alike.
# No Go and no test run: the staged-skills contract is a heavy-only lane the integrator drains.
# The expected texts hold literal backticks and $names, so single quotes are deliberate.
# shellcheck disable=SC2016
set -euo pipefail

tree_root="$(git rev-parse --show-toplevel)"
failure_count=0
fail() {
  printf 'FAIL: %s\n' "$1" >&2
  failure_count=$((failure_count + 1))
}

forensics_file="$tree_root/skills/do-work/actions/forensics.md"
readme_file="$tree_root/README.md"
clarify_file="$tree_root/skills/do-work/actions/clarify.md"
capture_file="$tree_root/skills/do-work/actions/capture.md"
fixture_file="$tree_root/_dev/tests/fixtures/retired-core-moved-command-triggers.tsv"
contract_file="$tree_root/_dev/tests/staged-skills-contract.sh"

# F9 (REQ-690 review): "stuck" routes to status, forensics keeps broken or failed work.
if ! grep -qxF -- '- User suspects something is broken or producing confusing results (for how long in-flight work has been quiet, see `actions/status.md`)' "$forensics_file"; then
  fail 'forensics.md When-to-Use line is not the F9 text'
fi
if grep -qF -- 'User suspects something is stuck' "$forensics_file"; then
  fail 'forensics.md still routes stuck questions to forensics'
fi
if ! grep -qF -- 'Run `do-work status` to see how long in-flight work has been quiet, and `do-work forensics` to diagnose failed or broken work.' "$readme_file"; then
  fail 'README.md failure FAQ does not carry the F9 sentence'
fi
if grep -qF -- 'to diagnose stuck or failed work' "$readme_file"; then
  fail 'README.md still sends stuck work to forensics'
fi

# F2 (REQ-688 review): clarify states the three-backtick minimum and no info string.
if ! grep -qF -- 'open a code fence one backtick longer than the longest backtick run anywhere in the text, never shorter than three, with no info string.' "$clarify_file"; then
  fail 'clarify.md fence rule lacks the minimum of three or the no-info-string clause'
fi
if grep -qF -- 'open a code fence longer than the longest backtick run' "$clarify_file"; then
  fail 'clarify.md still carries the old fence rule'
fi

# F3 (REQ-688 review): the addendum example uses a bare three-backtick fence.
if grep -qE -- '^> ````' "$capture_file"; then
  fail 'capture.md still has a four-backtick quoted fence'
fi
bare_fence_example_count="$(awk '
  previous_line == "> ```" && $0 == "> dark mode should also affect the sidebar" { opened = NR }
  opened && NR == opened + 1 && $0 == "> ```" { found++ }
  { previous_line = $0 }
  END { print found + 0 }
' "$capture_file")"
if [ "$bare_fence_example_count" -ne 1 ]; then
  fail "capture.md addendum example is not wrapped in bare three-backtick fences (found $bare_fence_example_count)"
fi

# M2 (REQ-692 review): test-side texts name the forward-row exception.
fixture_first_line="$(head -n 1 "$fixture_file")"
case "$fixture_first_line" in
  '#'*' The one exception is the core forward row for validate-feedback / triage feedback (UR-153, 2026-10-10).') ;;
  *) fail 'fixture header line 1 does not end with the forward-row exception (or lost its # comment prefix)' ;;
esac
if ! grep -qF -- 'fail "core must not route sibling-owned action $public_action through ./actions/ (a ../$sibling_owner/ forward row is allowed)"' "$contract_file"; then
  fail 'staged-skills-contract.sh fail message does not name the forward-row exception'
fi
if grep -qxF -- '    fail "core must not route sibling-owned action $public_action"' "$contract_file"; then
  fail 'staged-skills-contract.sh still carries the old fail message'
fi
if ! bash -n "$contract_file"; then
  fail 'staged-skills-contract.sh does not parse'
fi

if [ "$failure_count" -ne 0 ]; then
  printf 'REQ-696 GREEN probe: %s failure(s)\n' "$failure_count" >&2
  exit 1
fi
printf 'REQ-696 GREEN probe: ok\n'
