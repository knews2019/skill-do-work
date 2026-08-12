#!/usr/bin/env bash
# Fixture execution proofs for every promoted prescribed-shell script.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
# shellcheck source=_dev/tests/fixture-repo.sh
source "$repo_root/_dev/tests/fixture-repo.sh"
fixture_root="$(mktemp -d "${TMPDIR:-/tmp}/prescribed-shell-scripts.XXXXXX")" || exit 1
background_process_ids=""
cleanup_prescribed_shell_fixture() {
  for background_process_id in $background_process_ids; do
    kill "$background_process_id" 2>/dev/null || true
    wait "$background_process_id" 2>/dev/null || true
  done
  chmod -R u+rwX "$fixture_root" 2>/dev/null || true
  rm -rf "$fixture_root"
}
trap cleanup_prescribed_shell_fixture EXIT
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

# run-blocked-check: force the stock-Bash fallback, then prove timeout owns the
# isolated wrapper/descendant group without touching an unrelated group member.
fallback_bin="$fixture_root/fallback-bin"
mkdir -p "$fallback_bin"
ln -s "$(command -v bash)" "$fallback_bin/bash"
ln -s "$(command -v sh)" "$fallback_bin/sh"
ln -s "$(command -v sleep)" "$fallback_bin/sleep"
sleep 30 & unrelated_process_id=$!
background_process_ids="$background_process_ids $unrelated_process_id"
printf '%s\n' \
  'trap "" TERM' \
  'printf "%s\n" "$$" > "$BLOCKED_WRAPPER_PID_FILE"' \
  '(trap "" TERM; sleep 30) &' \
  'printf "%s\n" "$!" > "$BLOCKED_DESCENDANT_PID_FILE"' \
  'wait' > "$fixture_root/blocked-probe.sh"
PATH="$fallback_bin" \
  BLOCKED_WRAPPER_PID_FILE="$fixture_root/blocked-wrapper.pid" \
  BLOCKED_DESCENDANT_PID_FILE="$fixture_root/blocked-descendant.pid" \
  "$core_scripts/run-blocked-check.sh" "$fixture_root/blocked-probe.sh" 1 >/dev/null 2>&1
blocked_status=$?
[ "$blocked_status" -eq 124 ] || fail_case "run-blocked-check portable-timeout case returned $blocked_status instead of 124"
blocked_wrapper_pid="$(cat "$fixture_root/blocked-wrapper.pid")"
blocked_descendant_pid="$(cat "$fixture_root/blocked-descendant.pid")"
background_process_ids="$background_process_ids $blocked_wrapper_pid $blocked_descendant_pid"
kill -0 "$blocked_wrapper_pid" 2>/dev/null && fail_case 'run-blocked-check process-tree case left the wrapper alive'
kill -0 "$blocked_descendant_pid" 2>/dev/null && fail_case 'run-blocked-check process-tree case left the descendant alive'
kill -0 "$unrelated_process_id" 2>/dev/null || fail_case 'run-blocked-check process-tree cleanup killed an unrelated process in the test runner group'
: > "$fixture_root/test-runner-survived-timeout"
[ -e "$fixture_root/test-runner-survived-timeout" ] || fail_case 'run-blocked-check process-tree cleanup killed the test runner group'

# run-blocked-check: the fallback must preserve an ordinary probe status.
printf 'exit 23\n' > "$fixture_root/blocked-status-probe.sh"
PATH="$fallback_bin" "$core_scripts/run-blocked-check.sh" "$fixture_root/blocked-status-probe.sh" 2 >/dev/null 2>&1
blocked_ordinary_status=$?
[ "$blocked_ordinary_status" -eq 23 ] || fail_case "run-blocked-check ordinary-status case returned $blocked_ordinary_status instead of 23"

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

# stage-exact-deletion: pathspec-looking filenames must remain literal and isolated.
magic_deletion_repo="$fixture_root/magic-deletion-repo"
fixture_repo_init "$magic_deletion_repo"
magic_deleted_path=':(glob)*'
printf 'magic\n' > "$magic_deletion_repo/$magic_deleted_path"
printf 'other\n' > "$magic_deletion_repo/other.txt"
fixture_repo_commit_all "$magic_deletion_repo" base
rm "$magic_deletion_repo/$magic_deleted_path" "$magic_deletion_repo/other.txt"
(cd "$magic_deletion_repo" && "$core_scripts/stage-exact-deletion.sh" "$magic_deleted_path")
magic_deletion_status=$?
[ "$magic_deletion_status" -eq 0 ] || fail_case "stage-exact-deletion literal-pathspec case returned $magic_deletion_status instead of 0"
magic_cached_deletions="$(git -C "$magic_deletion_repo" diff --cached --name-status --no-renames)"
expected_magic_deletion="$(printf 'D\t%s' "$magic_deleted_path")"
[ "$magic_cached_deletions" = "$expected_magic_deletion" ] || fail_case 'stage-exact-deletion literal-pathspec case staged another path'

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

# generate-report-image: a direct backend receives inert prompt text and publishes
# from a private adjacent path only after success.
image_bin="$fixture_root/image-bin"
mkdir -p "$image_bin"
printf '%s\n' '#!/usr/bin/env bash' 'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; --prompt) printf "%s" "$2" > "$IMAGE_PROMPT_LOG"; shift 2 ;; *) shift ;; esac; done' 'printf "%s" "$output_path" > "$IMAGE_STAGED_PATH_LOG"' 'printf new-png > "$output_path"' > "$image_bin/imagegen"
chmod +x "$image_bin/imagegen"
image_output="$fixture_root/report.png"
printf old-png > "$image_output"
PATH="$image_bin:$PATH" \
  IMAGE_PROMPT_LOG="$fixture_root/image-prompt" \
  IMAGE_STAGED_PATH_LOG="$fixture_root/image-staged-path" \
  "$toolbox_scripts/generate-report-image.sh" "$image_output" 'blue style' 'diagram $(touch injected-image)' \
  || fail_case 'generate-report-image direct-backend case returned nonzero'
[ "$(cat "$image_output")" = new-png ] && [ ! -e "$fixture_root/injected-image" ] \
  || fail_case 'generate-report-image direct-backend case did not atomically replace the target with inert-prompt output'
[ "$(cat "$fixture_root/image-staged-path")" != "$image_output" ] \
  || fail_case 'generate-report-image direct-backend case wrote directly to the final target'
find "$fixture_root" -name '.report.png.generating.*' -print -quit | grep -q . \
  && fail_case 'generate-report-image direct-backend case leaked private staging after success'

# generate-report-image: a failed backend may leave an old target recoverable, but
# the invocation must fail and clean its partial private output.
image_failure_bin="$fixture_root/image-failure-bin"
mkdir -p "$image_failure_bin"
printf '%s\n' '#!/usr/bin/env bash' 'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; *) shift ;; esac; done' 'printf partial > "$output_path"' 'exit 7' > "$image_failure_bin/imagegen"
chmod +x "$image_failure_bin/imagegen"
printf stable-old-png > "$image_output"
PATH="$image_failure_bin:$PATH" DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND=0 \
  "$toolbox_scripts/generate-report-image.sh" "$image_output" style failure >/dev/null 2>&1 \
  && fail_case 'generate-report-image stale-target case accepted a failed backend'
[ "$(cat "$image_output")" = stable-old-png ] \
  || fail_case 'generate-report-image stale-target case changed the recoverable old target'
find "$fixture_root" -name '.report.png.generating.*' -print -quit | grep -q . \
  && fail_case 'generate-report-image stale-target case leaked private staging after failure'

# generate-report-image caller contract: mixed results retain both PIDs/statuses and
# wait for every job before evaluating current outputs.
image_mixed_bin="$fixture_root/image-mixed-bin"
mkdir -p "$image_mixed_bin"
printf '%s\n' '#!/usr/bin/env bash' 'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; --prompt) image_prompt="$2"; shift 2 ;; *) shift ;; esac; done' 'case "$image_prompt" in *mixed-success*) printf success > "$output_path"; : > "$MIXED_SUCCESS_DONE"; exit 0 ;; *) printf partial > "$output_path"; : > "$MIXED_FAILURE_DONE"; exit 9 ;; esac' > "$image_mixed_bin/imagegen"
chmod +x "$image_mixed_bin/imagegen"
mixed_success_target="$fixture_root/mixed-success.png"
mixed_failure_target="$fixture_root/mixed-failure.png"
printf stale > "$mixed_failure_target"
PATH="$image_mixed_bin:$PATH" MIXED_SUCCESS_DONE="$fixture_root/mixed-success.done" MIXED_FAILURE_DONE="$fixture_root/mixed-failure.done" \
  "$toolbox_scripts/generate-report-image.sh" "$mixed_success_target" style mixed-success &
image_generation_pids[0]=$!
PATH="$image_mixed_bin:$PATH" MIXED_SUCCESS_DONE="$fixture_root/mixed-success.done" MIXED_FAILURE_DONE="$fixture_root/mixed-failure.done" \
  "$toolbox_scripts/generate-report-image.sh" "$mixed_failure_target" style mixed-failure &
image_generation_pids[1]=$!
image_generation_statuses=()
image_index=0
while [ "$image_index" -lt "${#image_generation_pids[@]}" ]; do
  image_status=0
  wait "${image_generation_pids[$image_index]}" || image_status=$?
  image_generation_statuses[$image_index]="$image_status"
  image_index=$((image_index + 1))
done
[ "${image_generation_statuses[0]}" -eq 0 ] && [ "${image_generation_statuses[1]}" -ne 0 ] \
  || fail_case 'generate-report-image mixed-status case did not retain success and failure independently'
[ -e "$fixture_root/mixed-success.done" ] && [ -e "$fixture_root/mixed-failure.done" ] \
  || fail_case 'generate-report-image mixed-status case did not wait for every launched job'
[ "$(cat "$mixed_success_target")" = success ] && [ "$(cat "$mixed_failure_target")" = stale ] \
  || fail_case 'generate-report-image mixed-status case did not separate current output from a stale failed target'

# generate-report-image: interruption cleans the invocation-private file and leaves
# the old target untouched.
image_interrupt_bin="$fixture_root/image-interrupt-bin"
mkdir -p "$image_interrupt_bin"
printf '%s\n' '#!/usr/bin/env bash' 'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; *) shift ;; esac; done' 'printf "%s" "$output_path" > "$IMAGE_INTERRUPT_PATH_LOG"' 'printf partial > "$output_path"' ': > "$IMAGE_INTERRUPT_READY"' 'trap "exit 143" TERM INT HUP' 'while :; do sleep 0.1; done' > "$image_interrupt_bin/imagegen"
chmod +x "$image_interrupt_bin/imagegen"
interrupt_target="$fixture_root/interrupted.png"
printf stable-interrupt > "$interrupt_target"
PATH="$image_interrupt_bin:$PATH" IMAGE_INTERRUPT_PATH_LOG="$fixture_root/interrupted-stage" IMAGE_INTERRUPT_READY="$fixture_root/interrupted-ready" \
  "$toolbox_scripts/generate-report-image.sh" "$interrupt_target" style interruption >/dev/null 2>&1 &
interrupt_helper_pid=$!
background_process_ids="$background_process_ids $interrupt_helper_pid"
interrupt_wait_ticks=0
while [ ! -e "$fixture_root/interrupted-ready" ] && [ "$interrupt_wait_ticks" -lt 200 ]; do
  sleep 0.01
  interrupt_wait_ticks=$((interrupt_wait_ticks + 1))
done
if [ ! -e "$fixture_root/interrupted-ready" ]; then
  fail_case 'generate-report-image interruption case never reached the backend wait'
else
  kill -TERM "$interrupt_helper_pid" 2>/dev/null || true
  wait "$interrupt_helper_pid" 2>/dev/null
  interrupt_status=$?
  [ "$interrupt_status" -ne 0 ] || fail_case 'generate-report-image interruption case returned success'
  interrupted_stage_path="$(cat "$fixture_root/interrupted-stage")"
  [ ! -e "$interrupted_stage_path" ] || fail_case 'generate-report-image interruption case leaked private staging'
  [ "$(cat "$interrupt_target")" = stable-interrupt ] || fail_case 'generate-report-image interruption case changed the old target'
fi

# generate-report-image: the sandbox-bypassed executable is unreachable without the
# exact opt-in value, while exact 1 still exercises the explicitly authorized branch.
agentic_bin="$fixture_root/agentic-bin"
mkdir -p "$agentic_bin"
for agentic_tool in bash chmod cp mktemp mv rm; do
  ln -s "$(command -v "$agentic_tool")" "$agentic_bin/$agentic_tool"
done
printf '%s\n' '#!/usr/bin/env bash' ': > "$AGENTIC_INVOKED_MARKER"' 'printf agentic-png > generated.png' > "$agentic_bin/codex"
chmod +x "$agentic_bin/codex"
for agentic_opt_in in unset 0 true 01; do
  agentic_marker="$fixture_root/agentic-$agentic_opt_in.invoked"
  if [ "$agentic_opt_in" = unset ]; then
    (unset DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND; PATH="$agentic_bin" AGENTIC_INVOKED_MARKER="$agentic_marker" TMPDIR="$fixture_root" \
      "$toolbox_scripts/generate-report-image.sh" "$fixture_root/agentic-$agentic_opt_in.png" style description >/dev/null 2>&1)
  else
    PATH="$agentic_bin" AGENTIC_INVOKED_MARKER="$agentic_marker" TMPDIR="$fixture_root" DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND="$agentic_opt_in" \
      "$toolbox_scripts/generate-report-image.sh" "$fixture_root/agentic-$agentic_opt_in.png" style description >/dev/null 2>&1
  fi
  agentic_status=$?
  [ "$agentic_status" -ne 0 ] || fail_case "generate-report-image agentic opt-in case accepted non-exact value $agentic_opt_in"
  [ ! -e "$agentic_marker" ] || fail_case "generate-report-image agentic opt-in case invoked codex for value $agentic_opt_in"
done
agentic_marker="$fixture_root/agentic-1.invoked"
PATH="$agentic_bin" AGENTIC_INVOKED_MARKER="$agentic_marker" TMPDIR="$fixture_root" DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND=1 \
  "$toolbox_scripts/generate-report-image.sh" "$fixture_root/agentic-1.png" style description \
  || fail_case 'generate-report-image agentic opt-in case rejected exact value 1'
[ -e "$agentic_marker" ] && [ "$(cat "$fixture_root/agentic-1.png")" = agentic-png ] \
  || fail_case 'generate-report-image agentic opt-in case did not publish the explicitly authorized output'
find "$fixture_root" \( -name '.*.generating.*' -o -name 'do-work-ai-report-image.*' \) -print -quit | grep -q . \
  && fail_case 'generate-report-image agentic opt-in case leaked private paths'

# install-last30days: a SKILL.md-only tree fails check, is repaired from a
# complete fixture, and receives the full subtree plus ignore/Python guarantees.
upstream_repo="$fixture_root/last30days-upstream"
fixture_repo_init "$upstream_repo"
mkdir -p "$upstream_repo/skills/last30days/scripts" "$upstream_repo/skills/last30days/support"
printf '# Last30Days\n' > "$upstream_repo/skills/last30days/SKILL.md"
printf 'runtime\n' > "$upstream_repo/skills/last30days/scripts/last30days.py"
printf 'support\n' > "$upstream_repo/skills/last30days/support/data.txt"
fixture_repo_commit_all "$upstream_repo" fixture
last_project="$fixture_root/last-project"
fixture_repo_init "$last_project"
mkdir -p "$last_project/.claude/skills/last30days"
printf '# Sentinel only\n' > "$last_project/.claude/skills/last30days/SKILL.md"
python_bin="$fixture_root/python-bin"
mkdir -p "$python_bin"
printf '%s\n' '#!/usr/bin/env bash' 'exit 0' > "$python_bin/python3.12"
chmod +x "$python_bin/python3.12"
PATH="$python_bin:$PATH" "$toolbox_scripts/install-last30days.sh" check "$last_project" >/dev/null 2>&1 \
  && fail_case 'install-last30days sentinel-only check accepted a missing runtime script'
last_output="$(TMPDIR="$fixture_root" PATH="$python_bin:$PATH" "$toolbox_scripts/install-last30days.sh" install "$last_project" "$upstream_repo" 2>&1)" \
  || fail_case "install-last30days complete-source repair returned nonzero: $last_output"
[ "$(cat "$last_project/.claude/skills/last30days/scripts/last30days.py")" = runtime ] \
  && [ "$(cat "$last_project/.claude/skills/last30days/support/data.txt")" = support ] \
  || fail_case 'install-last30days complete-source repair omitted part of the source subtree'
git -C "$last_project" check-ignore -q .claude/skills/last30days/SKILL.md || fail_case 'install-last30days fixture-source case did not resolve the sibling exclude helper'
PATH="$python_bin:$PATH" "$toolbox_scripts/install-last30days.sh" check "$last_project" >/dev/null 2>&1 \
  || fail_case 'install-last30days repaired tree failed the complete runtime/ignore/Python check'

# install-last30days: reject an incomplete source before publishing any final tree.
incomplete_upstream_repo="$fixture_root/last30days-incomplete-upstream"
fixture_repo_init "$incomplete_upstream_repo"
mkdir -p "$incomplete_upstream_repo/skills/last30days"
printf '# Incomplete\n' > "$incomplete_upstream_repo/skills/last30days/SKILL.md"
fixture_repo_commit_all "$incomplete_upstream_repo" fixture
incomplete_source_project="$fixture_root/last-incomplete-source-project"
mkdir -p "$incomplete_source_project"
TMPDIR="$fixture_root" PATH="$python_bin:$PATH" "$toolbox_scripts/install-last30days.sh" install "$incomplete_source_project" "$incomplete_upstream_repo" >/dev/null 2>&1 \
  && fail_case 'install-last30days incomplete-source case returned success'
[ ! -e "$incomplete_source_project/.claude/skills/last30days" ] \
  || fail_case 'install-last30days incomplete-source case published a final tree'

# install-last30days: a partial copy failure leaves no final tree or private path.
copy_failure_bin="$fixture_root/copy-failure-bin"
mkdir -p "$copy_failure_bin"
printf '%s\n' '#!/usr/bin/env bash' 'while [ "$#" -gt 1 ]; do shift; done' 'copy_destination="$1"' 'mkdir -p "$copy_destination/scripts"' 'printf partial > "$copy_destination/SKILL.md"' 'exit 1' > "$copy_failure_bin/cp"
chmod +x "$copy_failure_bin/cp"
copy_failure_project="$fixture_root/last-copy-failure-project"
mkdir -p "$copy_failure_project"
TMPDIR="$fixture_root" PATH="$copy_failure_bin:$python_bin:$PATH" "$toolbox_scripts/install-last30days.sh" install "$copy_failure_project" "$upstream_repo" >/dev/null 2>&1 \
  && fail_case 'install-last30days copy-failure case returned success'
[ ! -e "$copy_failure_project/.claude/skills/last30days" ] \
  || fail_case 'install-last30days copy-failure case left a final tree'
find "$copy_failure_project/.claude/skills" -name '.last30days.*' -print -quit | grep -q . \
  && fail_case 'install-last30days copy-failure case leaked private paths'

# install-last30days: a publication failure with no prior destination removes any
# simulated partial final tree.
publication_failure_bin="$fixture_root/publication-failure-bin"
mkdir -p "$publication_failure_bin"
printf '%s\n' '#!/usr/bin/env bash' 'publication_source="$1"' 'publication_destination="$2"' 'case "$publication_source:$publication_destination" in *.last30days.staging.*:*/.claude/skills/last30days) mkdir -p "$publication_destination"; printf partial > "$publication_destination/SKILL.md"; exit 1 ;; esac' 'exec "$LAST30DAYS_REAL_MV" "$@"' > "$publication_failure_bin/mv"
chmod +x "$publication_failure_bin/mv"
publication_failure_project="$fixture_root/last-publication-failure-project"
mkdir -p "$publication_failure_project"
TMPDIR="$fixture_root" LAST30DAYS_REAL_MV="$(command -v mv)" PATH="$publication_failure_bin:$python_bin:$PATH" \
  "$toolbox_scripts/install-last30days.sh" install "$publication_failure_project" "$upstream_repo" >/dev/null 2>&1 \
  && fail_case 'install-last30days publication-failure case returned success'
[ ! -e "$publication_failure_project/.claude/skills/last30days" ] \
  || fail_case 'install-last30days publication-failure case left a new final tree'
find "$publication_failure_project/.claude/skills" -name '.last30days.*' -print -quit | grep -q . \
  && fail_case 'install-last30days publication-failure case leaked private paths'

# install-last30days: a replacement publication failure restores an existing
# incomplete tree byte-for-byte and removes private staging/backup paths.
replacement_failure_project="$fixture_root/last-replacement-failure-project"
mkdir -p "$replacement_failure_project/.claude/skills/last30days/legacy"
printf original-sentinel > "$replacement_failure_project/.claude/skills/last30days/SKILL.md"
printf original-byte > "$replacement_failure_project/.claude/skills/last30days/legacy/data.bin"
cp -R "$replacement_failure_project/.claude/skills/last30days" "$fixture_root/last30days-original-snapshot"
TMPDIR="$fixture_root" LAST30DAYS_REAL_MV="$(command -v mv)" PATH="$publication_failure_bin:$python_bin:$PATH" \
  "$toolbox_scripts/install-last30days.sh" install "$replacement_failure_project" "$upstream_repo" >/dev/null 2>&1 \
  && fail_case 'install-last30days replacement-failure case returned success'
diff -r "$fixture_root/last30days-original-snapshot" "$replacement_failure_project/.claude/skills/last30days" >/dev/null 2>&1 \
  || fail_case 'install-last30days replacement-failure case did not restore the prior tree byte-for-byte'
find "$replacement_failure_project/.claude/skills" -name '.last30days.*' -print -quit | grep -q . \
  && fail_case 'install-last30days replacement-failure case leaked private paths'

# install-last30days: interruption after the prior tree moves to backup restores it
# and removes both staging and backup paths.
backup_interrupt_bin="$fixture_root/backup-interrupt-bin"
mkdir -p "$backup_interrupt_bin"
printf '%s\n' '#!/usr/bin/env bash' 'interrupt_destination="$2"' 'case "$interrupt_destination" in */.last30days.backup.*/previous) "$LAST30DAYS_REAL_MV" "$@" || exit $?; kill -TERM "$PPID"; exit 0 ;; esac' 'exec "$LAST30DAYS_REAL_MV" "$@"' > "$backup_interrupt_bin/mv"
chmod +x "$backup_interrupt_bin/mv"
backup_interrupt_project="$fixture_root/last-backup-interrupt-project"
mkdir -p "$backup_interrupt_project/.claude/skills/last30days/legacy"
printf interrupt-sentinel > "$backup_interrupt_project/.claude/skills/last30days/SKILL.md"
printf interrupt-byte > "$backup_interrupt_project/.claude/skills/last30days/legacy/data.bin"
cp -R "$backup_interrupt_project/.claude/skills/last30days" "$fixture_root/last30days-interrupt-snapshot"
TMPDIR="$fixture_root" LAST30DAYS_REAL_MV="$(command -v mv)" PATH="$backup_interrupt_bin:$python_bin:$PATH" \
  "$toolbox_scripts/install-last30days.sh" install "$backup_interrupt_project" "$upstream_repo" >/dev/null 2>&1
backup_interrupt_status=$?
[ "$backup_interrupt_status" -eq 143 ] \
  || fail_case "install-last30days backup-interruption case returned $backup_interrupt_status instead of 143"
diff -r "$fixture_root/last30days-interrupt-snapshot" "$backup_interrupt_project/.claude/skills/last30days" >/dev/null 2>&1 \
  || fail_case 'install-last30days backup-interruption case did not restore the prior tree byte-for-byte'
find "$backup_interrupt_project/.claude/skills" -name '.last30days.*' -print -quit | grep -q . \
  && fail_case 'install-last30days backup-interruption case leaked private paths'

if [ "$failure_count" -gt 0 ]; then
  exit 1
fi
printf 'Prescribed shell script behavior probes passed (22 named script cases).\n'
