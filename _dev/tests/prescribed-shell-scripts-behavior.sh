#!/usr/bin/env bash
# Fixture execution proofs for every promoted prescribed-shell script.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
# shellcheck source=_dev/tests/fixture-repo.sh
source "$repo_root/_dev/tests/fixture-repo.sh"
fixture_root="$(mktemp -d "${TMPDIR:-/tmp}/prescribed-shell-scripts.XXXXXX")" || exit 1
trap 'chmod -R u+rwX "$fixture_root" 2>/dev/null || true; rm -rf "$fixture_root"' EXIT
failure_count=0

fail_case() {
  printf 'FAIL: %s\n' "$1" >&2
  failure_count=$((failure_count + 1))
}

core_scripts="$repo_root/skills/do-work/scripts"
knowledge_scripts="$repo_root/skills/do-work-knowledge/scripts"
toolbox_scripts="$repo_root/skills/do-work-toolbox/scripts"

# show-commit-diff: a real merge must expose first-parent file changes.
merge_repo="$fixture_root/merge-repo"
fixture_repo_init "$merge_repo"
printf 'base\n' > "$merge_repo/base.txt"
fixture_repo_commit_all "$merge_repo" base
git -C "$merge_repo" checkout -qb feature
printf 'feature\n' > "$merge_repo/feature.txt"
fixture_repo_commit_all "$merge_repo" feature
git -C "$merge_repo" checkout -q main
printf 'main\n' > "$merge_repo/main.txt"
fixture_repo_commit_all "$merge_repo" main
git -C "$merge_repo" merge --no-ff -qm merge feature
merge_commit="$(git -C "$merge_repo" rev-parse HEAD)"
merge_output="$(cd "$merge_repo" && "$core_scripts/show-commit-diff.sh" "$merge_commit" 2>&1)" || fail_case 'show-commit-diff real-merge case returned nonzero'
printf '%s' "$merge_output" | grep -q 'feature.txt' || fail_case 'show-commit-diff real-merge case hid the merged file'

# add-local-git-exclude: a subdirectory caller must append once to Git's actual exclude.
exclude_repo="$fixture_root/exclude-repo"
fixture_repo_init "$exclude_repo"
mkdir -p "$exclude_repo/nested" "$exclude_repo/cache/data"
(cd "$exclude_repo/nested" && "$core_scripts/add-local-git-exclude.sh" ../cache/data/file.bin '**/cache/data/' >/dev/null) || fail_case 'add-local-git-exclude subdirectory case returned nonzero'
exclude_file="$exclude_repo/$(git -C "$exclude_repo" rev-parse --git-path info/exclude)"
[ "$(grep -Fc '**/cache/data/' "$exclude_file")" -eq 1 ] || fail_case 'add-local-git-exclude subdirectory case did not append exactly once'

# atomic-download: a failed transfer may write private bytes but never change the final target.
atomic_bin="$fixture_root/atomic-bin"
mkdir -p "$atomic_bin"
printf '%s\n' '#!/usr/bin/env bash' 'output_path=""' 'while [ "$#" -gt 0 ]; do case "$1" in -o) output_path="$2"; shift 2 ;; *) shift ;; esac; done' 'printf partial > "$output_path"' 'exit 22' > "$atomic_bin/curl"
chmod +x "$atomic_bin/curl"
printf 'stable' > "$fixture_root/atomic-target"
PATH="$atomic_bin:$PATH" "$core_scripts/atomic-download.sh" https://example.invalid/fail "$fixture_root/atomic-target" >/dev/null 2>&1 && fail_case 'atomic-download partial-publication case accepted a failed transfer'
[ "$(cat "$fixture_root/atomic-target")" = stable ] || fail_case 'atomic-download partial-publication case changed the final target'
find "$fixture_root" -name 'atomic-target.download.*' -print -quit | grep -q . && fail_case 'atomic-download partial-publication case leaked private scratch'

# capture-screenshot: coordinate two writers so the loser cannot publish the winner's private copy.
capture_root="$fixture_root/capture"
mkdir -p "$capture_root/a" "$capture_root/b" "$capture_root/assets"
printf 'dispatch-a' > "$capture_root/a/source.png"
printf 'dispatch-b' > "$capture_root/b/source.png"
capture_destination="$capture_root/assets/result.png"
race_a_verified="$capture_root/a-verified"
race_b_copied="$capture_root/b-copied"
(
  cmp() { command cmp "$@"; comparison_status=$?; if [ "$comparison_status" -eq 0 ]; then : > "$race_a_verified"; while [ ! -e "$race_b_copied" ]; do sleep 0.01; done; fi; return "$comparison_status"; }
  export race_a_verified race_b_copied
  export -f cmp
  "$core_scripts/capture-screenshot.sh" --staged "$capture_root/a/source.png" "$capture_destination"
) >"$capture_root/a.out" 2>"$capture_root/a.err" & race_a_pid=$!
(
  cp() { while [ ! -e "$race_a_verified" ]; do sleep 0.01; done; command cp "$@"; copy_status=$?; [ "$copy_status" -ne 0 ] || : > "$race_b_copied"; return "$copy_status"; }
  cmp() { command cmp "$@"; comparison_status=$?; while [ ! -e "$capture_destination" ]; do sleep 0.01; done; return "$comparison_status"; }
  export race_a_verified race_b_copied capture_destination
  export -f cp cmp
  "$core_scripts/capture-screenshot.sh" --staged "$capture_root/b/source.png" "$capture_destination"
) >"$capture_root/b.out" 2>"$capture_root/b.err" & race_b_pid=$!
race_wait_ticks=0
while { kill -0 "$race_a_pid" 2>/dev/null || kill -0 "$race_b_pid" 2>/dev/null; } \
  && [ "$race_wait_ticks" -lt 1000 ]; do
  sleep 0.01
  race_wait_ticks=$((race_wait_ticks + 1))
done
if kill -0 "$race_a_pid" 2>/dev/null || kill -0 "$race_b_pid" 2>/dev/null; then
  kill "$race_a_pid" "$race_b_pid" 2>/dev/null || true
  wait "$race_a_pid" 2>/dev/null || true
  wait "$race_b_pid" 2>/dev/null || true
  fail_case 'capture-screenshot coordinated-race case timed out'
  race_a_status=1
  race_b_status=1
else
  wait "$race_a_pid"; race_a_status=$?
  wait "$race_b_pid"; race_b_status=$?
fi
[ "$race_a_status" -eq 0 ] && [ "$race_b_status" -ne 0 ] || fail_case 'capture-screenshot coordinated-race case did not install exactly one writer'
[ "$(cat "$capture_destination")" = dispatch-a ] || fail_case 'capture-screenshot coordinated-race case published cross-dispatch bytes'
[ ! -e "$capture_root/a/source.png" ] && [ "$(cat "$capture_root/b/source.png")" = dispatch-b ] || fail_case 'capture-screenshot coordinated-race case did not preserve only the loser source'
find "$capture_root/assets" -name 'result.png.copying.*' -print -quit | grep -q . && fail_case 'capture-screenshot coordinated-race case leaked private scratch'

# run-blocked-check: remove GNU timeout tools from PATH and require fallback status 124.
fallback_bin="$fixture_root/fallback-bin"
mkdir -p "$fallback_bin"
ln -s "$(command -v bash)" "$fallback_bin/bash"
ln -s "$(command -v sh)" "$fallback_bin/sh"
ln -s "$(command -v sleep)" "$fallback_bin/sleep"
printf 'sleep 5\n' > "$fixture_root/blocked-probe.sh"
PATH="$fallback_bin" "$core_scripts/run-blocked-check.sh" "$fixture_root/blocked-probe.sh" 1 >/dev/null 2>&1
blocked_status=$?
[ "$blocked_status" -eq 124 ] || fail_case "run-blocked-check portable-timeout case returned $blocked_status instead of 124"

# protected-inventory: once a secret path is quarantined it stays out of association.
inventory_repo="$fixture_root/inventory-repo"
fixture_repo_init "$inventory_repo"
mkdir -p "$inventory_repo/do-work/archive/UR-001"
printf '%s\n' '---' 'id: REQ-001' 'status: completed' '---' '' '## Implementation Summary' '' '- `safe.txt` — fixture' > "$inventory_repo/do-work/archive/UR-001/REQ-001-fixture.md"
printf 'base\n' > "$inventory_repo/safe.txt"
fixture_repo_commit_all "$inventory_repo" base
printf 'change\n' >> "$inventory_repo/safe.txt"
printf 'secret\n' > "$inventory_repo/.env.local"
inventory_output="$(cd "$inventory_repo" && "$core_scripts/protected-inventory.sh" start)" || fail_case 'protected-inventory start case returned nonzero'
printf '%s' "$inventory_output" | grep -q $'X\t.env.local' || fail_case 'protected-inventory start case did not quarantine the secret path'
association_output="$(cd "$inventory_repo" && "$core_scripts/protected-inventory.sh" associate)" || fail_case 'protected-inventory associate case returned nonzero'
printf '%s' "$association_output" | grep -q $'REQ-001\tsafe.txt' || fail_case 'protected-inventory associate case lost the safe owner'
printf '%s' "$association_output" | grep -q '.env.local' && fail_case 'protected-inventory associate case leaked the quarantined path'

# stage-exact-deletion: stage only the named pathological deletion.
deletion_repo="$fixture_root/deletion-repo"
fixture_repo_init "$deletion_repo"
printf 'delete\n' > "$deletion_repo/secret name.key"
printf 'keep\n' > "$deletion_repo/other.txt"
fixture_repo_commit_all "$deletion_repo" base
rm "$deletion_repo/secret name.key"
printf 'changed\n' >> "$deletion_repo/other.txt"
(cd "$deletion_repo" && "$core_scripts/stage-exact-deletion.sh" 'secret name.key') || fail_case 'stage-exact-deletion pathological-name case returned nonzero'
[ "$(git -C "$deletion_repo" diff --cached --name-only)" = 'secret name.key' ] || fail_case 'stage-exact-deletion pathological-name case staged another path'

# lexical-memory-recall: apostrophes and command syntax remain data while attribution is emitted.
memory_root="$fixture_root/memory with space"
mkdir -p "$memory_root/logs"
printf '%s\n' '# Memory' '## Decisions' 'Use cobalt release trains.' > "$memory_root/working-memory.md"
recall_output="$($knowledge_scripts/lexical-memory-recall.sh "$memory_root" "cobalt'; touch $fixture_root/injected #")" || fail_case 'lexical-memory-recall raw-query case returned nonzero'
[ ! -e "$fixture_root/injected" ] || fail_case 'lexical-memory-recall raw-query case executed query text'
printf '%s' "$recall_output" | grep -q $'working memory\t## Decisions\tUse cobalt' || fail_case 'lexical-memory-recall raw-query case omitted attribution'

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
  printf '%s' "$manual_output" | grep -q 'MANUAL STEP' || fail_case 'install-memory-hooks no-jq case omitted manual status'
fi

# generate-report-image: a direct backend receives inert prompt text and must create a non-empty exact path.
image_bin="$fixture_root/image-bin"
mkdir -p "$image_bin"
printf '%s\n' '#!/usr/bin/env bash' 'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; --prompt) printf "%s" "$2" > "$IMAGE_PROMPT_LOG"; shift 2 ;; *) shift ;; esac; done' 'printf png > "$output_path"' > "$image_bin/imagegen"
chmod +x "$image_bin/imagegen"
image_output="$fixture_root/report.png"
PATH="$image_bin:$PATH" IMAGE_PROMPT_LOG="$fixture_root/image-prompt" "$toolbox_scripts/generate-report-image.sh" "$image_output" 'blue style' 'diagram $(touch injected-image)' || fail_case 'generate-report-image direct-backend case returned nonzero'
[ -s "$image_output" ] && [ ! -e "$fixture_root/injected-image" ] || fail_case 'generate-report-image direct-backend case did not keep the prompt inert'

# install-last30days: install from a fixture upstream, add the sibling-core ignore, and verify.
upstream_repo="$fixture_root/last30days-upstream"
fixture_repo_init "$upstream_repo"
mkdir -p "$upstream_repo/skills/last30days/scripts"
printf '# Last30Days\n' > "$upstream_repo/skills/last30days/SKILL.md"
printf 'tool\n' > "$upstream_repo/skills/last30days/scripts/tool.py"
fixture_repo_commit_all "$upstream_repo" fixture
last_project="$fixture_root/last-project"
fixture_repo_init "$last_project"
python_bin="$fixture_root/python-bin"
mkdir -p "$python_bin"
printf '%s\n' '#!/usr/bin/env bash' 'exit 0' > "$python_bin/python3.12"
chmod +x "$python_bin/python3.12"
last_output="$(PATH="$python_bin:$PATH" "$toolbox_scripts/install-last30days.sh" install "$last_project" "$upstream_repo" 2>&1)" || fail_case "install-last30days fixture-source case returned nonzero: $last_output"
[ -s "$last_project/.claude/skills/last30days/SKILL.md" ] || fail_case 'install-last30days fixture-source case omitted the skill file'
git -C "$last_project" check-ignore -q .claude/skills/last30days/SKILL.md || fail_case 'install-last30days fixture-source case did not resolve the sibling exclude helper'

if [ "$failure_count" -gt 0 ]; then
  exit 1
fi
printf 'Prescribed shell script behavior probes passed (11 named script cases).\n'
