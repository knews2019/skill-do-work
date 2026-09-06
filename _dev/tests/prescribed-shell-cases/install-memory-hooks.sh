#!/usr/bin/env bash
# Fixture execution proofs for install-memory-hooks.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

# install-memory-hooks: a partial prior install gets only its missing sibling and preserves foreign hooks.
hooks_project="$fixture_root/hooks-project"
mkdir -p "$hooks_project/.claude"
printf '%s\n' '{"hooks":{"SessionStart":[{"hooks":[{"type":"command","command":"foreign.sh"},{"type":"command","command":".claude/skills/do-work-knowledge/hooks/memory-session-start.sh"}]}]}}' > "$hooks_project/.claude/settings.json"
if command -v jq >/dev/null 2>&1; then
  "$knowledge_scripts/install-memory-hooks.sh" "$hooks_project" "$repo_root/skills/do-work-knowledge/hooks/memory-hooks.json" >/dev/null || fail_case 'install-memory-hooks partial-merge case returned nonzero'
  [ "$(grep -o 'memory-session-start.sh' "$hooks_project/.claude/settings.json" | wc -l | tr -d ' ')" -eq 1 ] || fail_case 'install-memory-hooks partial-merge case duplicated the present hook'
  grep -q 'memory-stop-capture.sh' "$hooks_project/.claude/settings.json" || fail_case 'install-memory-hooks partial-merge case omitted the missing hook'
  grep -q 'foreign.sh' "$hooks_project/.claude/settings.json" || fail_case 'install-memory-hooks partial-merge case clobbered a foreign hook'
else
  manual_output="$(PATH=/usr/bin:/bin "$knowledge_scripts/install-memory-hooks.sh" "$hooks_project" "$repo_root/skills/do-work-knowledge/hooks/memory-hooks.json")" || fail_case 'install-memory-hooks no-jq case returned nonzero'
  grep -q 'MANUAL STEP' <<<"$manual_output" || fail_case 'install-memory-hooks no-jq case omitted manual status'
fi

prescribed_shell_finish
