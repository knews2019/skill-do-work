#!/usr/bin/env bash
# Keep the shipped defensive-surface audit complete and its safe deletions deleted.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
audit_file="$repo_root/decisions/audits/2026-08-11-defensive-surface.md"
failure_count=0

fail() {
  printf 'FAIL: %s\n' "$1" >&2
  failure_count=$((failure_count + 1))
}

if [[ ! -f "$audit_file" ]]; then
  fail 'defensive-surface audit artifact is missing.'
else
  for required_heading in \
    '## Executable layers' \
    '## Prose layers — core' \
    '## Prose layers — toolbox' \
    '## Prose layers — knowledge and board' \
    '## Deletes applied' \
    '## Behavior-changing candidates retained'
  do
    grep -Fqx "$required_heading" "$audit_file" \
      || fail "defensive-surface audit is missing heading: $required_heading"
  done
fi

# Every shipped shell source is an executable defensive surface. The audit must name it,
# while REQ-165's action-shell-blocks probe owns executable parsing/lint reachability.
while IFS= read -r -d '' shell_file; do
  relative_path="${shell_file#"$repo_root/"}"
  if ! grep -Fq "\`$relative_path\`" "$audit_file"; then
    fail "$relative_path has no executable-layer disposition in the defensive-surface audit."
  fi
done < <(find "$repo_root/skills" -type f -name '*.sh' -print0)

# Explicit defensive prose sections must have a same-row location + surface disposition.
while IFS= read -r -d '' markdown_file; do
  relative_path="${markdown_file#"$repo_root/"}"
  while IFS= read -r defensive_heading; do
    [[ -n "$defensive_heading" ]] || continue
    surface_name="${defensive_heading#\#\# }"
    if ! grep -F "\`$relative_path\`" "$audit_file" | grep -Fq "$surface_name"; then
      fail "$relative_path <$defensive_heading> has no same-row audit disposition."
    fi
  done < <(grep -E '^## (Rules|Common Rationalizations|Red Flags|Warnings)$|^## Known Failure Mode & Recovery$' "$markdown_file" || true)
done < <(find "$repo_root/skills" -type f -name '*.md' -print0)

for deleted_table_file in \
  skills/do-work/actions/verify-requests.md \
  skills/do-work-toolbox/actions/code-review.md \
  skills/do-work-toolbox/actions/quick-wins.md \
  skills/do-work-toolbox/actions/ui-review.md
do
  if grep -Eq '^## (Common Rationalizations|Red Flags)$' "$repo_root/$deleted_table_file"; then
    fail "$deleted_table_file restored a generic defensive table deleted by REQ-168."
  fi
done

if grep -Fq '## Common Rationalizations' "$repo_root/skills/do-work/actions/commit.md"; then
  fail 'actions/commit.md restored its duplicate generic rationalization table.'
fi
if grep -Fq 'Single commit with >20 files' "$repo_root/skills/do-work/actions/commit.md"; then
  fail 'actions/commit.md restored the unearned file-count warning.'
fi

if grep -Fq '## Common Rationalizations' "$repo_root/skills/do-work-toolbox/actions/inspect.md"; then
  fail 'actions/inspect.md restored its duplicate generic rationalization table.'
fi
for deleted_inspect_phrase in \
  'Report produced without a `git status` / `git diff` reading actually happening' \
  'Files listed with no REQ association attempt' \
  'Debug artifacts (console.log, debugger, commented-out blocks)'
do
  if grep -Fq "$deleted_inspect_phrase" "$repo_root/skills/do-work-toolbox/actions/inspect.md"; then
    fail "actions/inspect.md restored generic warning: $deleted_inspect_phrase"
  fi
done

if [[ "$failure_count" -gt 0 ]]; then
  exit 1
fi

printf 'Defensive-surface delete-or-test audit checks passed.\n'
