#!/usr/bin/env bash
# Fixture execution proofs for every promoted prescribed-shell script.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
# Absolute path to this file, so the closing count can read it from any caller's directory.
suite_file="$repo_root/_dev/tests/${BASH_SOURCE[0]##*/}"
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

# atomic-download: a rate-limited host answers 429 once and then succeeds. The fake curl
# below models curl's own internal retry loop, so it survives that 429 only if the caller
# allowed a retry — which is the whole point of the flag set. It also records the
# Authorization header it was handed, so the opt-in credential path is observable.
atomic_retry_bin="$fixture_root/atomic-retry-bin"
mkdir -p "$atomic_retry_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'output_path=""' \
  'retry_limit=0' \
  'authorization_header=""' \
  'while [ "$#" -gt 0 ]; do' \
  '  case "$1" in' \
  '    -o) output_path="$2"; shift 2 ;;' \
  '    --retry) retry_limit="$2"; shift 2 ;;' \
  '    -H) authorization_header="$2"; shift 2 ;;' \
  '    *) shift ;;' \
  '  esac' \
  'done' \
  'printf "%s" "$authorization_header" > "$ATOMIC_HEADER_LOG"' \
  'transfer_attempt=0' \
  'while :; do' \
  '  transfer_attempt=$((transfer_attempt + 1))' \
  '  printf "%s" "$transfer_attempt" > "$ATOMIC_ATTEMPT_LOG"' \
  '  if [ "$transfer_attempt" -gt 1 ]; then' \
  '    printf complete-payload > "$output_path"' \
  '    exit 0' \
  '  fi' \
  '  if [ "$transfer_attempt" -gt "$retry_limit" ]; then' \
  '    exit 22' \
  '  fi' \
  'done' \
  > "$atomic_retry_bin/curl"
chmod +x "$atomic_retry_bin/curl"

printf 'stale before retry\n' > "$fixture_root/atomic-retry-target"
GH_TOKEN='' GITHUB_TOKEN='' PATH="$atomic_retry_bin:$PATH" \
  ATOMIC_ATTEMPT_LOG="$fixture_root/atomic-attempts" \
  ATOMIC_HEADER_LOG="$fixture_root/atomic-header" \
  "$core_scripts/atomic-download.sh" https://example.invalid/rate-limited "$fixture_root/atomic-retry-target" >/dev/null 2>&1 \
  || fail_case 'atomic-download retry case did not survive a transient 429'
[ "$(cat "$fixture_root/atomic-retry-target")" = complete-payload ] \
  || fail_case 'atomic-download retry case did not publish the successful attempt'
[ "$(cat "$fixture_root/atomic-attempts")" = 2 ] \
  || fail_case 'atomic-download retry case did not let curl retry the rate-limited transfer'
[ -z "$(cat "$fixture_root/atomic-header")" ] \
  || fail_case 'atomic-download retry case sent an Authorization header with no token configured'
find "$fixture_root" -name 'atomic-retry-target.download.*' -print -quit | grep -q . \
  && fail_case 'atomic-download retry case leaked private scratch'

# atomic-download: an opt-in token becomes a bearer credential; GH_TOKEN wins over GITHUB_TOKEN.
printf 'stale before credential\n' > "$fixture_root/atomic-credential-target"
GH_TOKEN=primary-token GITHUB_TOKEN=fallback-token PATH="$atomic_retry_bin:$PATH" \
  ATOMIC_ATTEMPT_LOG="$fixture_root/atomic-credential-attempts" \
  ATOMIC_HEADER_LOG="$fixture_root/atomic-credential-header" \
  "$core_scripts/atomic-download.sh" https://example.invalid/private "$fixture_root/atomic-credential-target" >/dev/null 2>&1 \
  || fail_case 'atomic-download credential case returned nonzero'
[ "$(cat "$fixture_root/atomic-credential-header")" = 'Authorization: Bearer primary-token' ] \
  || fail_case 'atomic-download credential case did not send GH_TOKEN as a bearer credential'
printf 'stale before fallback credential\n' > "$fixture_root/atomic-fallback-target"
GH_TOKEN='' GITHUB_TOKEN=fallback-token PATH="$atomic_retry_bin:$PATH" \
  ATOMIC_ATTEMPT_LOG="$fixture_root/atomic-fallback-attempts" \
  ATOMIC_HEADER_LOG="$fixture_root/atomic-fallback-header" \
  "$core_scripts/atomic-download.sh" https://example.invalid/private "$fixture_root/atomic-fallback-target" >/dev/null 2>&1 \
  || fail_case 'atomic-download fallback-credential case returned nonzero'
[ "$(cat "$fixture_root/atomic-fallback-header")" = 'Authorization: Bearer fallback-token' ] \
  || fail_case 'atomic-download fallback-credential case did not fall back to GITHUB_TOKEN'

# atomic-download: a target occupied by a DIRECTORY must fail closed. `mv` treats a
# directory operand as a container rather than a collision, so the download nests
# inside it and exits zero — and the caller reads that status as proof the file
# landed. The canonical statement of the rule is the shipped guide's
# "Verified exact publication" section.
atomic_success_bin="$fixture_root/atomic-success-bin"
mkdir -p "$atomic_success_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'output_path=""' \
  'while [ "$#" -gt 0 ]; do case "$1" in -o) output_path="$2"; shift 2 ;; *) shift ;; esac; done' \
  'printf complete-payload > "$output_path"' \
  'exit 0' \
  > "$atomic_success_bin/curl"
chmod +x "$atomic_success_bin/curl"

atomic_occupied_target="$fixture_root/atomic-occupied-target"
mkdir -p "$atomic_occupied_target"
printf 'occupant\n' > "$atomic_occupied_target/pre-existing.txt"
PATH="$atomic_success_bin:$PATH" \
  "$core_scripts/atomic-download.sh" https://example.invalid/payload "$atomic_occupied_target" >/dev/null 2>&1 \
  && fail_case 'atomic-download occupied-target case reported success for a publication that nested'
[ -d "$atomic_occupied_target" ] \
  || fail_case 'atomic-download occupied-target case did not leave the occupying directory in place'
[ "$(cat "$atomic_occupied_target/pre-existing.txt")" = occupant ] \
  || fail_case 'atomic-download occupied-target case disturbed the occupying directory contents'
find "$atomic_occupied_target" -name '*.download.*' -print -quit | grep -q . \
  && fail_case 'atomic-download occupied-target case abandoned its private file inside the occupant'

# capture-screenshot: a destination occupied by a DIRECTORY must fail closed and keep
# the staged source. `ln` refuses an occupied FILE — which is where the no-clobber
# guarantee comes from — but nests on a directory and exits zero, and under --staged
# that zero is read as permission to delete the staged source. The dispatch holds the
# only copy, so a false success there destroys it.
capture_occupied_root="$fixture_root/capture-occupied"
mkdir -p "$capture_occupied_root/stage" "$capture_occupied_root/assets/result.png"
printf 'the only copy' > "$capture_occupied_root/stage/source.png"
printf 'occupant\n' > "$capture_occupied_root/assets/result.png/pre-existing.txt"
"$core_scripts/capture-screenshot.sh" --staged \
  "$capture_occupied_root/stage/source.png" "$capture_occupied_root/assets/result.png" >/dev/null 2>&1 \
  && fail_case 'capture-screenshot occupied-destination case reported success for a publication that nested'
[ -f "$capture_occupied_root/stage/source.png" ] \
  || fail_case 'capture-screenshot occupied-destination case destroyed the staged source it never published'
[ -d "$capture_occupied_root/assets/result.png" ] \
  || fail_case 'capture-screenshot occupied-destination case did not leave the occupying directory in place'
[ "$(cat "$capture_occupied_root/assets/result.png/pre-existing.txt")" = occupant ] \
  || fail_case 'capture-screenshot occupied-destination case disturbed the occupying directory contents'
find "$capture_occupied_root/assets/result.png" -name '*.copying.*' -print -quit | grep -q . \
  && fail_case 'capture-screenshot occupied-destination case abandoned its private copy inside the occupant'

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
# kill -0 counts zombies as alive; a killed-but-unreaped descendant (reparent target
# reaps lazily, e.g. container PID 1) must read as dead, so check the process state.
process_runs_unreaped_excluded() {
  local process_state
  process_state="$(ps -o stat= -p "$1" 2>/dev/null | tr -d '[:space:]')" || return 1
  [ -n "$process_state" ] || return 1
  case "$process_state" in Z*) return 1 ;; esac
  return 0
}
process_runs_unreaped_excluded "$blocked_wrapper_pid" && fail_case 'run-blocked-check process-tree case left the wrapper alive'
process_runs_unreaped_excluded "$blocked_descendant_pid" && fail_case 'run-blocked-check process-tree case left the descendant alive'
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

portfolio_root="$fixture_root/portfolio"
mkdir -p "$portfolio_root/deliverables/portfolio-snapshots"
portfolio_source="$portfolio_root/retained-summary.md"
portfolio_canonical="$portfolio_root/deliverables/portfolio-summary.md"
portfolio_candidate="$portfolio_root/deliverables/portfolio-snapshots/portfolio-summary-20260815T120000Z.md"
printf 'new portfolio bytes\n' > "$portfolio_source"

# publish-portfolio-summary: canonical-only publication atomically refreshes only
# the canonical path and retains the source used for publication.
printf 'old canonical\n' > "$portfolio_canonical"
portfolio_canonical_output="$($toolbox_scripts/publish-portfolio-summary.sh --canonical-only "$portfolio_source" "$portfolio_canonical" 2>/dev/null)" \
  || fail_case 'publish-portfolio-summary canonical-only case returned nonzero'
cmp -s "$portfolio_source" "$portfolio_canonical" \
  || fail_case 'publish-portfolio-summary canonical-only case changed the retained bytes'
[ -f "$portfolio_source" ] \
  || fail_case 'publish-portfolio-summary canonical-only case consumed the retained source'
[ "$portfolio_canonical_output" = "$portfolio_canonical" ] \
  || fail_case 'publish-portfolio-summary canonical-only case did not report the canonical path'
find "$portfolio_root/deliverables/portfolio-snapshots" -type f -print -quit | grep -q . \
  && fail_case 'publish-portfolio-summary canonical-only case created a snapshot'

# publish-portfolio-summary: the preservation branch publishes a snapshot first, then
# refreshes canonical from the same verified bytes — but as an independent file. Same
# bytes is the requirement; same inode is the defect, because a shared inode makes the
# immutable snapshot follow every later in-place edit of canonical (REQ-205, was asserted
# the other way round by REQ-199 before durable immutability was tested).
printf 'prior canonical\n' > "$portfolio_canonical"
portfolio_snapshot_output="$($toolbox_scripts/publish-portfolio-summary.sh --with-snapshot "$portfolio_source" "$portfolio_canonical" "$portfolio_candidate" 2>/dev/null)" \
  || fail_case 'publish-portfolio-summary snapshot-success case returned nonzero'
[ "$portfolio_snapshot_output" = "$(printf '%s\n%s' "$portfolio_canonical" "$portfolio_candidate")" ] \
  || fail_case 'publish-portfolio-summary snapshot-success case did not report both published paths'
cmp -s "$portfolio_source" "$portfolio_candidate" && cmp -s "$portfolio_source" "$portfolio_canonical" \
  || fail_case 'publish-portfolio-summary snapshot-success case did not preserve byte identity'
[ -f "$portfolio_candidate" ] && [ -f "$portfolio_canonical" ] \
  || fail_case 'publish-portfolio-summary snapshot-success case did not publish two regular files'
[ "$portfolio_candidate" -ef "$portfolio_canonical" ] \
  && fail_case 'publish-portfolio-summary snapshot-success case aliased the snapshot to the canonical inode'
printf 'canonical mutated after publication\n' > "$portfolio_canonical"
[ "$(cat "$portfolio_candidate")" = 'new portfolio bytes' ] \
  || fail_case 'publish-portfolio-summary snapshot-success case let a later canonical edit rewrite the snapshot'
cp "$portfolio_source" "$portfolio_canonical"

# publish-portfolio-summary: an occupied candidate remains immutable and advances
# to the first numeric suffix without cleaning any prior snapshot.
portfolio_collision_candidate="$portfolio_root/deliverables/portfolio-snapshots/portfolio-summary-20260815T130000Z.md"
portfolio_collision_suffix="$portfolio_root/deliverables/portfolio-snapshots/portfolio-summary-20260815T130000Z-2.md"
printf 'occupied collision\n' > "$portfolio_collision_candidate"
portfolio_collision_output="$($toolbox_scripts/publish-portfolio-summary.sh --with-snapshot "$portfolio_source" "$portfolio_canonical" "$portfolio_collision_candidate" 2>/dev/null)" \
  || fail_case 'publish-portfolio-summary collision case returned nonzero'
[ "$(cat "$portfolio_collision_candidate")" = 'occupied collision' ] \
  || fail_case 'publish-portfolio-summary collision case changed the occupant'
cmp -s "$portfolio_source" "$portfolio_collision_suffix" \
  || fail_case 'publish-portfolio-summary collision case did not publish the numeric suffix'
[ "$portfolio_collision_output" = "$(printf '%s\n%s' "$portfolio_canonical" "$portfolio_collision_suffix")" ] \
  || fail_case 'publish-portfolio-summary collision case reported the wrong suffix'
[ "$(cat "$portfolio_candidate")" = 'new portfolio bytes' ] \
  || fail_case 'publish-portfolio-summary collision case cleaned an unrelated prior snapshot'

# publish-portfolio-summary: exclusive snapshot failure leaves the prior canonical
# unchanged and never leaks the private verified copy.
portfolio_ln_failure_bin="$fixture_root/portfolio-ln-failure-bin"
mkdir -p "$portfolio_ln_failure_bin"
printf '%s\n' '#!/usr/bin/env bash' 'exit 1' > "$portfolio_ln_failure_bin/ln"
chmod +x "$portfolio_ln_failure_bin/ln"
portfolio_failure_candidate="$portfolio_root/deliverables/portfolio-snapshots/portfolio-summary-20260815T140000Z.md"
printf 'stable before snapshot failure\n' > "$portfolio_canonical"
PATH="$portfolio_ln_failure_bin:$PATH" \
  "$toolbox_scripts/publish-portfolio-summary.sh" --with-snapshot "$portfolio_source" "$portfolio_canonical" "$portfolio_failure_candidate" >/dev/null 2>&1 \
  && fail_case 'publish-portfolio-summary snapshot-failure case returned success'
[ "$(cat "$portfolio_canonical")" = 'stable before snapshot failure' ] \
  || fail_case 'publish-portfolio-summary snapshot-failure case changed the prior canonical'
[ ! -e "$portfolio_failure_candidate" ] \
  || fail_case 'publish-portfolio-summary snapshot-failure case left a published snapshot'
find "$portfolio_root/deliverables" -name '.portfolio-summary.md.publishing.*' -print -quit | grep -q . \
  && fail_case 'publish-portfolio-summary snapshot-failure case leaked private bytes'

# publish-portfolio-summary: a later canonical replacement failure retains the
# already-published snapshot while preserving the prior canonical.
portfolio_mv_failure_bin="$fixture_root/portfolio-mv-failure-bin"
mkdir -p "$portfolio_mv_failure_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'if [ "$2" = "$PORTFOLIO_FAIL_CANONICAL" ]; then exit 1; fi' \
  'exec "$PORTFOLIO_REAL_MV" "$@"' \
  > "$portfolio_mv_failure_bin/mv"
chmod +x "$portfolio_mv_failure_bin/mv"
portfolio_late_failure_candidate="$portfolio_root/deliverables/portfolio-snapshots/portfolio-summary-20260815T150000Z.md"
printf 'stable before canonical failure\n' > "$portfolio_canonical"
PORTFOLIO_FAIL_CANONICAL="$portfolio_canonical" \
PORTFOLIO_REAL_MV="$(command -v mv)" \
PATH="$portfolio_mv_failure_bin:$PATH" \
  "$toolbox_scripts/publish-portfolio-summary.sh" --with-snapshot "$portfolio_source" "$portfolio_canonical" "$portfolio_late_failure_candidate" >/dev/null 2>&1 \
  && fail_case 'publish-portfolio-summary canonical-failure case returned success'
[ "$(cat "$portfolio_canonical")" = 'stable before canonical failure' ] \
  || fail_case 'publish-portfolio-summary canonical-failure case changed the prior canonical'
cmp -s "$portfolio_source" "$portfolio_late_failure_candidate" \
  || fail_case 'publish-portfolio-summary canonical-failure case did not retain the published snapshot'
find "$portfolio_root/deliverables" -name '.portfolio-summary.md.publishing.*' -print -quit | grep -q . \
  && fail_case 'publish-portfolio-summary canonical-failure case leaked private bytes'

# publish-portfolio-summary: `ln` links *into* a directory operand instead of colliding
# with it, so a snapshot candidate occupied by a directory must advance to the numeric
# suffix and leave no private file nested inside the directory.
portfolio_directory_candidate="$portfolio_root/deliverables/portfolio-snapshots/portfolio-summary-20260815T160000Z.md"
portfolio_directory_suffix="$portfolio_root/deliverables/portfolio-snapshots/portfolio-summary-20260815T160000Z-2.md"
mkdir -p "$portfolio_directory_candidate"
printf 'occupant\n' > "$portfolio_directory_candidate/occupant.txt"
printf 'stable before directory candidate\n' > "$portfolio_canonical"
portfolio_directory_output="$($toolbox_scripts/publish-portfolio-summary.sh --with-snapshot "$portfolio_source" "$portfolio_canonical" "$portfolio_directory_candidate" 2>/dev/null)" \
  || fail_case 'publish-portfolio-summary snapshot-directory case returned nonzero instead of advancing'
[ -d "$portfolio_directory_candidate" ] && [ "$(ls -A "$portfolio_directory_candidate")" = occupant.txt ] \
  || fail_case 'publish-portfolio-summary snapshot-directory case nested a private file inside the occupying directory'
cmp -s "$portfolio_source" "$portfolio_directory_suffix" \
  || fail_case 'publish-portfolio-summary snapshot-directory case did not advance to the numeric suffix'
[ "$portfolio_directory_output" = "$(printf '%s\n%s' "$portfolio_canonical" "$portfolio_directory_suffix")" ] \
  || fail_case 'publish-portfolio-summary snapshot-directory case reported the wrong published paths'

# publish-portfolio-summary: `mv` moves *into* a directory operand, so a canonical path
# occupied by a directory must fail closed — never advance, never publish inside it, and
# never leave the private copy nested there.
portfolio_directory_canonical="$portfolio_root/deliverables/canonical-as-directory.md"
mkdir -p "$portfolio_directory_canonical"
printf 'canonical occupant\n' > "$portfolio_directory_canonical/occupant.txt"
"$toolbox_scripts/publish-portfolio-summary.sh" --canonical-only "$portfolio_source" "$portfolio_directory_canonical" >/dev/null 2>&1 \
  && fail_case 'publish-portfolio-summary canonical-directory case reported success'
[ -d "$portfolio_directory_canonical" ] && [ "$(ls -A "$portfolio_directory_canonical")" = occupant.txt ] \
  || fail_case 'publish-portfolio-summary canonical-directory case did not leave the occupying directory unchanged'
find "$portfolio_root/deliverables" -name '.canonical-as-directory.md.publishing.*' -print -quit | grep -q . \
  && fail_case 'publish-portfolio-summary canonical-directory case leaked private bytes'

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

# generate-report-image-batch: the batch joins each target name to its own private
# staging directory, so a target carrying a path separator would write outside that
# boundary, and a pair with no colon would silently make the prompt the filename. Both
# are usage errors that must be refused before any staging directory exists.
batch_arguments_root="$fixture_root/batch-arguments"
mkdir -p "$batch_arguments_root/ai-reports/report"
"$toolbox_scripts/generate-report-image-batch.sh" "$batch_arguments_root/ai-reports/report" style \
  '../escaped.png:<prompt 1>' >/dev/null 2>&1
batch_escape_status=$?
[ "$batch_escape_status" -eq 2 ] \
  || fail_case "generate-report-image-batch path-separator target case returned $batch_escape_status instead of the usage status 2"
"$toolbox_scripts/generate-report-image-batch.sh" "$batch_arguments_root/ai-reports/report" style \
  '01-architecture.png' >/dev/null 2>&1
batch_unpaired_status=$?
[ "$batch_unpaired_status" -eq 2 ] \
  || fail_case "generate-report-image-batch unpaired-argument case returned $batch_unpaired_status instead of the usage status 2"
find "$batch_arguments_root/ai-reports/report" -name '.generated.staging.*' -print -quit | grep -q . \
  && fail_case 'generate-report-image-batch usage-error cases allocated invocation-private staging'

# generate-report-image-batch: the shipped batch owns staging, launch, wait-all,
# freshness, and publication. An all-failed batch must leave no public or private
# directory and still return zero so the caller falls back to SVG/Mermaid; a mixed
# batch must publish only the status-zero, non-empty current image after waiting
# every job, and report the published directory on stdout.
run_ai_report_batch_replay() {
  replay_name="$1"
  replay_bin="$2"
  replay_root="$fixture_root/$replay_name"
  mkdir -p "$replay_root/ai-reports/<report-slug>"
  (
    cd "$replay_root" || exit 2
    PATH="$replay_bin:$PATH" \
      REPLAY_FIRST_DONE="$replay_root/first.done" \
      REPLAY_SECOND_DONE="$replay_root/second.done" \
      "$toolbox_scripts/generate-report-image-batch.sh" \
        'ai-reports/<report-slug>' \
        'replay style' \
        '01-architecture.png:<prompt 1>' \
        '02-dataflow.png:<prompt 2>'
  ) > "$fixture_root/$replay_name.stdout" 2> "$fixture_root/$replay_name.stderr"
}

image_all_failed_bin="$fixture_root/image-all-failed-bin"
mkdir -p "$image_all_failed_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; --prompt) image_prompt="$2"; shift 2 ;; *) shift ;; esac; done' \
  'case "$image_prompt" in *"<prompt 1>"*) : > "$REPLAY_FIRST_DONE" ;; *) : > "$REPLAY_SECOND_DONE" ;; esac' \
  'printf partial > "$output_path"' \
  'exit 9' \
  > "$image_all_failed_bin/imagegen"
chmod +x "$image_all_failed_bin/imagegen"
run_ai_report_batch_replay image-all-failed "$image_all_failed_bin" \
  || fail_case 'ai-report all-failed batch replay returned nonzero instead of falling back'
[ ! -e "$fixture_root/image-all-failed/ai-reports/<report-slug>/generated" ] \
  || fail_case 'ai-report all-failed batch replay published an empty generated/ directory'
find "$fixture_root/image-all-failed/ai-reports/<report-slug>" -name '.generated.staging.*' -print -quit | grep -q . \
  && fail_case 'ai-report all-failed batch replay leaked invocation-private staging'
[ -e "$fixture_root/image-all-failed/first.done" ] && [ -e "$fixture_root/image-all-failed/second.done" ] \
  || fail_case 'ai-report all-failed batch replay did not wait for every launched job'
[ ! -s "$fixture_root/image-all-failed.stdout" ] \
  || fail_case 'ai-report all-failed batch replay reported a published directory on stdout'

image_batch_mixed_bin="$fixture_root/image-batch-mixed-bin"
mkdir -p "$image_batch_mixed_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; --prompt) image_prompt="$2"; shift 2 ;; *) shift ;; esac; done' \
  'case "$image_prompt" in *"<prompt 1>"*) printf current-success > "$output_path"; : > "$REPLAY_FIRST_DONE"; exit 0 ;; *) printf partial > "$output_path"; : > "$REPLAY_SECOND_DONE"; exit 9 ;; esac' \
  > "$image_batch_mixed_bin/imagegen"
chmod +x "$image_batch_mixed_bin/imagegen"
run_ai_report_batch_replay image-batch-mixed "$image_batch_mixed_bin" \
  || fail_case 'ai-report mixed batch replay returned nonzero'
[ "$(cat "$fixture_root/image-batch-mixed/ai-reports/<report-slug>/generated/01-architecture.png" 2>/dev/null)" = current-success ] \
  && [ ! -e "$fixture_root/image-batch-mixed/ai-reports/<report-slug>/generated/02-dataflow.png" ] \
  || fail_case 'ai-report mixed batch replay did not publish only the status-backed successful image'
find "$fixture_root/image-batch-mixed/ai-reports/<report-slug>" -name '.generated.staging.*' -print -quit | grep -q . \
  && fail_case 'ai-report mixed batch replay leaked invocation-private staging'
[ -e "$fixture_root/image-batch-mixed/first.done" ] && [ -e "$fixture_root/image-batch-mixed/second.done" ] \
  || fail_case 'ai-report mixed batch replay did not wait for every launched job'
# The caller learns that generated/ was published from stdout — a script cannot set the
# caller's own GEN variable, so the published path is the batch's success signal.
mixed_published_directory="$(cat "$fixture_root/image-batch-mixed.stdout" 2>/dev/null)"
[ -n "$mixed_published_directory" ] \
  && [ "$mixed_published_directory" -ef "$fixture_root/image-batch-mixed/ai-reports/<report-slug>/generated" ] \
  || fail_case 'ai-report mixed batch replay did not report the published generated/ directory on stdout'

# generate-report-image-batch, interrupted: the batch owns the process tree it started, so
# signalling the caller must leave no helper — and no helper's descendant — alive, and
# must reap them before staging is removed. A file-only cleanup passes the directory
# assertions below while backends keep running against the deleted stage.
image_interrupt_batch_bin="$fixture_root/image-interrupt-batch-bin"
mkdir -p "$image_interrupt_batch_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; --prompt) image_prompt="$2"; shift 2 ;; *) shift ;; esac; done' \
  'printf partial > "$output_path"' \
  '( while :; do sleep 0.2; done ) &' \
  'helper_descendant_pid=$!' \
  'case "$image_prompt" in' \
  '  *"<prompt 1>"*) printf "%s\n" "$$" > "$REPLAY_FIRST_PID"; printf "%s\n" "$helper_descendant_pid" > "$REPLAY_FIRST_CHILD_PID"; : > "$REPLAY_FIRST_DONE" ;;' \
  '  *) printf "%s\n" "$$" > "$REPLAY_SECOND_PID"; printf "%s\n" "$helper_descendant_pid" > "$REPLAY_SECOND_CHILD_PID"; : > "$REPLAY_SECOND_DONE" ;;' \
  'esac' \
  'while :; do sleep 0.2; done' \
  > "$image_interrupt_batch_bin/imagegen"
chmod +x "$image_interrupt_batch_bin/imagegen"

interrupt_batch_root="$fixture_root/image-interrupt-batch"
mkdir -p "$interrupt_batch_root/ai-reports/<report-slug>"
(
  cd "$interrupt_batch_root" \
    && exec env PATH="$image_interrupt_batch_bin:$PATH" \
      REPLAY_FIRST_DONE="$interrupt_batch_root/first.done" \
      REPLAY_SECOND_DONE="$interrupt_batch_root/second.done" \
      REPLAY_FIRST_PID="$interrupt_batch_root/first.pid" \
      REPLAY_SECOND_PID="$interrupt_batch_root/second.pid" \
      REPLAY_FIRST_CHILD_PID="$interrupt_batch_root/first.child.pid" \
      REPLAY_SECOND_CHILD_PID="$interrupt_batch_root/second.child.pid" \
      "$toolbox_scripts/generate-report-image-batch.sh" \
        'ai-reports/<report-slug>' \
        'replay style' \
        '01-architecture.png:<prompt 1>' \
        '02-dataflow.png:<prompt 2>'
) &
interrupt_batch_pid=$!
interrupt_ready_ticks=0
while [ "$interrupt_ready_ticks" -lt 200 ]; do
  [ -s "$interrupt_batch_root/first.child.pid" ] && [ -s "$interrupt_batch_root/second.child.pid" ] && break
  sleep 0.05
  interrupt_ready_ticks=$((interrupt_ready_ticks + 1))
done
[ -s "$interrupt_batch_root/first.child.pid" ] && [ -s "$interrupt_batch_root/second.child.pid" ] \
  || fail_case 'ai-report interrupted batch replay never reached both running helpers'
for interrupt_pid_file in first.pid second.pid first.child.pid second.child.pid; do
  background_process_ids="$background_process_ids $(cat "$interrupt_batch_root/$interrupt_pid_file" 2>/dev/null || true)"
done
kill -TERM "$interrupt_batch_pid" 2>/dev/null || true
interrupt_batch_status=0
wait "$interrupt_batch_pid" || interrupt_batch_status=$?
[ "$interrupt_batch_status" -eq 143 ] \
  || fail_case "ai-report interrupted batch replay exited $interrupt_batch_status instead of the TERM status 143"
interrupt_survivor_ticks=0
while [ "$interrupt_survivor_ticks" -lt 40 ]; do
  interrupt_survivors=0
  for interrupt_pid_file in first.pid second.pid first.child.pid second.child.pid; do
    interrupt_recorded_pid="$(cat "$interrupt_batch_root/$interrupt_pid_file" 2>/dev/null || true)"
    [ -n "$interrupt_recorded_pid" ] || continue
    kill -0 "$interrupt_recorded_pid" 2>/dev/null && interrupt_survivors=$((interrupt_survivors + 1))
  done
  [ "$interrupt_survivors" -eq 0 ] && break
  sleep 0.05
  interrupt_survivor_ticks=$((interrupt_survivor_ticks + 1))
done
[ "$interrupt_survivors" -eq 0 ] \
  || fail_case "ai-report interrupted batch replay left $interrupt_survivors helper process(es) or descendant(s) alive"
find "$interrupt_batch_root/ai-reports/<report-slug>" -name '.generated.staging.*' -print -quit | grep -q . \
  && fail_case 'ai-report interrupted batch replay leaked invocation-private staging'
[ ! -e "$interrupt_batch_root/ai-reports/<report-slug>/generated" ] \
  || fail_case 'ai-report interrupted batch replay published generated/'

# generate-report-image-batch, destination appears at the final boundary: `mv` treats an
# existing directory as a container, so the check-then-rename window can nest the
# private stage inside a colliding generated/ and still exit zero. The mv shim below
# creates the destination in exactly that window.
image_publish_collision_bin="$fixture_root/image-publish-collision-bin"
mkdir -p "$image_publish_collision_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; --prompt) image_prompt="$2"; shift 2 ;; *) shift ;; esac; done' \
  'case "$image_prompt" in *"<prompt 1>"*) : > "$REPLAY_FIRST_DONE" ;; *) : > "$REPLAY_SECOND_DONE" ;; esac' \
  'printf current-success > "$output_path"' \
  > "$image_publish_collision_bin/imagegen"
chmod +x "$image_publish_collision_bin/imagegen"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'if [ "$#" -eq 2 ]; then' \
  '  case "$2" in */generated) mkdir -p "$2"; printf owned-by-someone-else > "$2/keep.txt" ;; esac' \
  'fi' \
  'exec /bin/mv "$@"' \
  > "$image_publish_collision_bin/mv"
chmod +x "$image_publish_collision_bin/mv"
run_ai_report_batch_replay image-publish-collision "$image_publish_collision_bin" \
  && fail_case 'ai-report publish-collision replay reported success after the destination appeared'
publish_collision_generated="$fixture_root/image-publish-collision/ai-reports/<report-slug>/generated"
[ "$(cat "$publish_collision_generated/keep.txt" 2>/dev/null)" = owned-by-someone-else ] \
  || fail_case 'ai-report publish-collision replay did not preserve the colliding destination byte-for-byte'
[ "$(ls -A "$publish_collision_generated" 2>/dev/null)" = keep.txt ] \
  || fail_case 'ai-report publish-collision replay left its staged batch nested inside the colliding destination'
find "$fixture_root/image-publish-collision/ai-reports/<report-slug>" -name '.generated.staging.*' -print -quit | grep -q . \
  && fail_case 'ai-report publish-collision replay leaked invocation-private staging'

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

# generate-report-image, interrupted directly: the helper owns the process tree it
# started, so signalling it must leave neither the backend nor the backend's own
# descendant alive, and must reap them before the private stage is removed. A
# bare-PID kill passes the file assertions below while the descendant keeps running.
image_tree_bin="$fixture_root/image-tree-bin"
mkdir -p "$image_tree_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; *) shift ;; esac; done' \
  'printf "%s" "$output_path" > "$IMAGE_TREE_STAGE_LOG"' \
  'printf partial > "$output_path"' \
  '( while :; do sleep 0.2; done ) &' \
  'printf "%s\n" "$!" > "$IMAGE_TREE_DESCENDANT_PID"' \
  'printf "%s\n" "$$" > "$IMAGE_TREE_BACKEND_PID"' \
  ': > "$IMAGE_TREE_READY"' \
  'while :; do sleep 0.2; done' \
  > "$image_tree_bin/imagegen"
chmod +x "$image_tree_bin/imagegen"
image_tree_target="$fixture_root/process-tree.png"
printf stable-tree > "$image_tree_target"
PATH="$image_tree_bin:$PATH" \
  IMAGE_TREE_STAGE_LOG="$fixture_root/process-tree-stage" \
  IMAGE_TREE_BACKEND_PID="$fixture_root/process-tree-backend.pid" \
  IMAGE_TREE_DESCENDANT_PID="$fixture_root/process-tree-descendant.pid" \
  IMAGE_TREE_READY="$fixture_root/process-tree-ready" \
  "$toolbox_scripts/generate-report-image.sh" "$image_tree_target" style process-tree >/dev/null 2>&1 &
image_tree_helper_pid=$!
background_process_ids="$background_process_ids $image_tree_helper_pid"
image_tree_ready_ticks=0
while [ "$image_tree_ready_ticks" -lt 200 ]; do
  [ -e "$fixture_root/process-tree-ready" ] \
    && [ -s "$fixture_root/process-tree-backend.pid" ] \
    && [ -s "$fixture_root/process-tree-descendant.pid" ] && break
  sleep 0.05
  image_tree_ready_ticks=$((image_tree_ready_ticks + 1))
done
if [ ! -s "$fixture_root/process-tree-backend.pid" ] || [ ! -s "$fixture_root/process-tree-descendant.pid" ]; then
  fail_case 'generate-report-image process-tree case never reached a running backend and descendant'
else
  image_tree_backend_pid="$(cat "$fixture_root/process-tree-backend.pid")"
  image_tree_descendant_pid="$(cat "$fixture_root/process-tree-descendant.pid")"
  background_process_ids="$background_process_ids $image_tree_backend_pid $image_tree_descendant_pid"
  kill -TERM "$image_tree_helper_pid" 2>/dev/null || true
  wait "$image_tree_helper_pid" 2>/dev/null
  image_tree_status=$?
  [ "$image_tree_status" -eq 143 ] \
    || fail_case "generate-report-image process-tree case returned $image_tree_status instead of the TERM status 143"
  image_tree_survivor_ticks=0
  while [ "$image_tree_survivor_ticks" -lt 40 ]; do
    image_tree_survivors=0
    for image_tree_recorded_pid in "$image_tree_backend_pid" "$image_tree_descendant_pid"; do
      kill -0 "$image_tree_recorded_pid" 2>/dev/null && image_tree_survivors=$((image_tree_survivors + 1))
    done
    [ "$image_tree_survivors" -eq 0 ] && break
    sleep 0.05
    image_tree_survivor_ticks=$((image_tree_survivor_ticks + 1))
  done
  [ "$image_tree_survivors" -eq 0 ] \
    || fail_case "generate-report-image process-tree case left $image_tree_survivors backend process(es) or descendant(s) alive"
  image_tree_stage_path="$(cat "$fixture_root/process-tree-stage")"
  [ ! -e "$image_tree_stage_path" ] || fail_case 'generate-report-image process-tree case leaked private staging'
  [ "$(cat "$image_tree_target")" = stable-tree ] || fail_case 'generate-report-image process-tree case changed the old target'
fi

# generate-report-image: `mv` treats an existing destination directory as a container,
# so an output path occupied by a directory nests the staged image inside it and still
# exits zero. Publication must fail closed and leave that directory untouched.
image_directory_bin="$fixture_root/image-directory-bin"
mkdir -p "$image_directory_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'while [ "$#" -gt 0 ]; do case "$1" in --output) output_path="$2"; shift 2 ;; *) shift ;; esac; done' \
  'printf new-png > "$output_path"' \
  > "$image_directory_bin/imagegen"
chmod +x "$image_directory_bin/imagegen"
image_directory_parent="$fixture_root/directory-target"
image_directory_target="$image_directory_parent/report.png"
mkdir -p "$image_directory_target"
printf owned-by-someone-else > "$image_directory_target/keep.txt"
PATH="$image_directory_bin:$PATH" DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND=0 \
  "$toolbox_scripts/generate-report-image.sh" "$image_directory_target" style directory-collision >/dev/null 2>&1 \
  && fail_case 'generate-report-image output-is-a-directory case reported success'
[ "$(cat "$image_directory_target/keep.txt" 2>/dev/null)" = owned-by-someone-else ] \
  || fail_case 'generate-report-image output-is-a-directory case did not preserve the occupying directory byte-for-byte'
[ "$(ls -A "$image_directory_target" 2>/dev/null)" = keep.txt ] \
  || fail_case 'generate-report-image output-is-a-directory case left its staged image nested inside the occupying directory'
find "$image_directory_parent" -name '.report.png.generating.*' -print -quit | grep -q . \
  && fail_case 'generate-report-image output-is-a-directory case leaked private staging'

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

# install-last30days: the target can reappear between the backup `mv` and the
# publication `mv`, and `mv` then nests the staging tree inside it and exits zero. The
# mv shim below recreates the target in exactly that window; publication must fail
# closed, leave the reappeared tree byte-for-byte, and keep the prior tree recoverable.
publication_collision_bin="$fixture_root/publication-collision-bin"
mkdir -p "$publication_collision_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'case "${2:-}" in */.claude/skills/last30days) mkdir -p "$2"; printf owned-by-someone-else > "$2/keep.txt" ;; esac' \
  'exec /bin/mv "$@"' \
  > "$publication_collision_bin/mv"
chmod +x "$publication_collision_bin/mv"
publication_collision_project="$fixture_root/last-publication-collision-project"
publication_collision_target="$publication_collision_project/.claude/skills/last30days"
mkdir -p "$publication_collision_target/legacy"
printf collision-sentinel > "$publication_collision_target/SKILL.md"
printf collision-byte > "$publication_collision_target/legacy/data.bin"
TMPDIR="$fixture_root" PATH="$publication_collision_bin:$python_bin:$PATH" \
  "$toolbox_scripts/install-last30days.sh" install "$publication_collision_project" "$upstream_repo" >/dev/null 2>&1 \
  && fail_case 'install-last30days publication-collision case returned success'
[ "$(cat "$publication_collision_target/keep.txt" 2>/dev/null)" = owned-by-someone-else ] \
  || fail_case 'install-last30days publication-collision case did not preserve the reappeared target byte-for-byte'
[ "$(ls -A "$publication_collision_target" 2>/dev/null)" = keep.txt ] \
  || fail_case 'install-last30days publication-collision case left its staging tree nested inside the reappeared target'
find "$publication_collision_project/.claude/skills" -name '.last30days.staging.*' -print -quit | grep -q . \
  && fail_case 'install-last30days publication-collision case leaked private staging'
[ -s "$(find "$publication_collision_project/.claude/skills" -path '*/.last30days.backup.*/previous/SKILL.md' -print -quit)" ] \
  || fail_case 'install-last30days publication-collision case did not leave the prior tree recoverable at its backup path'

# cleanup-req-reservations: a marker whose REQ file is committed is removed at
# any zero-padding width and archive depth; non-marker entries stay untouched.
reservation_project="$fixture_root/reservation-project"
fixture_repo_init "$reservation_project"
mkdir -p "$reservation_project/do-work/.req-reservations" \
  "$reservation_project/do-work/queue" \
  "$reservation_project/do-work/archive/UR-001"
printf 'req\n' > "$reservation_project/do-work/queue/REQ-203-fixture.md"
printf 'req\n' > "$reservation_project/do-work/archive/UR-001/REQ-7-archived-fixture.md"
: > "$reservation_project/do-work/.req-reservations/REQ-000203"
: > "$reservation_project/do-work/.req-reservations/REQ-000007"
: > "$reservation_project/do-work/.req-reservations/README"
fixture_repo_commit_all "$reservation_project" 'captured fixtures'
ln -s REQ-000203 "$reservation_project/do-work/.req-reservations/REQ-000777"
reservation_output="$("$core_scripts/cleanup-req-reservations.sh" "$reservation_project")" \
  || fail_case 'cleanup-req-reservations redundant-marker case returned nonzero'
printf '%s' "$reservation_output" | grep -q 'removed 2 stale REQ reservation marker' \
  || fail_case 'cleanup-req-reservations redundant-marker case did not report exactly two removals'
[ ! -e "$reservation_project/do-work/.req-reservations/REQ-000203" ] \
  || fail_case 'cleanup-req-reservations redundant-marker case kept the queue-claimed marker'
[ ! -e "$reservation_project/do-work/.req-reservations/REQ-000007" ] \
  || fail_case 'cleanup-req-reservations redundant-marker case kept the archive-claimed marker'
[ -e "$reservation_project/do-work/.req-reservations/README" ] \
  || fail_case 'cleanup-req-reservations redundant-marker case deleted a non-marker file'
[ -L "$reservation_project/do-work/.req-reservations/REQ-000777" ] \
  || fail_case 'cleanup-req-reservations redundant-marker case deleted a symlinked marker'

# cleanup-req-reservations: in a git work tree, a REQ file merely present on
# disk is a capture still staging — its marker must survive until the capture
# commits, then be reaped. Deleting early breaks capture's prescribed
# `git add do-work/.req-reservations/REQ-NNNNNN`.
printf 'req\n' > "$reservation_project/do-work/queue/REQ-500-inflight-fixture.md"
: > "$reservation_project/do-work/.req-reservations/REQ-000500"
reservation_uncommitted_output="$("$core_scripts/cleanup-req-reservations.sh" "$reservation_project")" \
  || fail_case 'cleanup-req-reservations uncommitted-capture case returned nonzero'
[ -z "$reservation_uncommitted_output" ] \
  || fail_case 'cleanup-req-reservations uncommitted-capture case printed output while the capture was mid-flight'
[ -f "$reservation_project/do-work/.req-reservations/REQ-000500" ] \
  || fail_case 'cleanup-req-reservations uncommitted-capture case deleted a mid-capture marker'
fixture_repo_commit_all "$reservation_project" 'captured REQ-500'
reservation_committed_output="$("$core_scripts/cleanup-req-reservations.sh" "$reservation_project")" \
  || fail_case 'cleanup-req-reservations committed-capture case returned nonzero'
printf '%s' "$reservation_committed_output" | grep -q 'removed 1 stale REQ reservation marker' \
  || fail_case 'cleanup-req-reservations committed-capture case did not reap the landed marker'
[ ! -e "$reservation_project/do-work/.req-reservations/REQ-000500" ] \
  || fail_case 'cleanup-req-reservations committed-capture case kept the landed marker'

# cleanup-req-reservations: a young marker with no REQ file is a capture still in
# flight — it must survive, and a run that removes nothing must print nothing.
: > "$reservation_project/do-work/.req-reservations/REQ-000999"
reservation_inflight_output="$("$core_scripts/cleanup-req-reservations.sh" "$reservation_project")" \
  || fail_case 'cleanup-req-reservations in-flight-marker case returned nonzero'
[ -z "$reservation_inflight_output" ] \
  || fail_case 'cleanup-req-reservations in-flight-marker case printed output for a no-op run'
[ -f "$reservation_project/do-work/.req-reservations/REQ-000999" ] \
  || fail_case 'cleanup-req-reservations in-flight-marker case deleted a young unmatched marker'

# cleanup-req-reservations: the same unmatched marker aged past two days is an
# abandoned capture and is removed.
touch -m -t 202001010000 "$reservation_project/do-work/.req-reservations/REQ-000999"
reservation_timeout_output="$("$core_scripts/cleanup-req-reservations.sh" "$reservation_project")" \
  || fail_case 'cleanup-req-reservations timeout-marker case returned nonzero'
printf '%s' "$reservation_timeout_output" | grep -q 'removed 1 stale REQ reservation marker' \
  || fail_case 'cleanup-req-reservations timeout-marker case did not report the timeout removal'
[ ! -e "$reservation_project/do-work/.req-reservations/REQ-000999" ] \
  || fail_case 'cleanup-req-reservations timeout-marker case kept the abandoned marker'

# cleanup-req-reservations: a repo without a reservation directory is a silent no-op.
reservation_absent_project="$fixture_root/reservation-absent-project"
mkdir -p "$reservation_absent_project/do-work/queue"
reservation_absent_output="$("$core_scripts/cleanup-req-reservations.sh" "$reservation_absent_project")" \
  || fail_case 'cleanup-req-reservations missing-directory case returned nonzero'
[ -z "$reservation_absent_output" ] \
  || fail_case 'cleanup-req-reservations missing-directory case printed output with nothing to clean'

# cleanup-req-reservations: a symlinked reservation directory is refused, so the
# automatic hook can never delete files outside the project through the link —
# a regular child reached through a symlinked parent passes a per-file -L check.
reservation_symlink_project="$fixture_root/reservation-symlink-project"
reservation_external_store="$fixture_root/reservation-external-store"
mkdir -p "$reservation_symlink_project/do-work/queue" "$reservation_external_store"
printf 'req\n' > "$reservation_symlink_project/do-work/queue/REQ-42-fixture.md"
: > "$reservation_external_store/REQ-000042"
touch -m -t 202001010000 "$reservation_external_store/REQ-000042"
ln -s "$reservation_external_store" "$reservation_symlink_project/do-work/.req-reservations"
reservation_symlink_output="$("$core_scripts/cleanup-req-reservations.sh" "$reservation_symlink_project")" \
  || fail_case 'cleanup-req-reservations symlinked-directory case returned nonzero'
[ -z "$reservation_symlink_output" ] \
  || fail_case 'cleanup-req-reservations symlinked-directory case printed output for a refused store'
[ -f "$reservation_external_store/REQ-000042" ] \
  || fail_case 'cleanup-req-reservations symlinked-directory case deleted through the symlinked store'

# repair-req-timestamps: a future created_at in the queue is rewritten to the
# file's own mtime — the actual write instant — and the correction is logged.
repair_mtime_project="$fixture_root/repair-mtime-project"
mkdir -p "$repair_mtime_project/do-work/queue"
printf -- '---\nid: REQ-801\nstatus: pending\ncreated_at: 2093-01-01T00:00:00Z\n---\n\nbody\n' \
  > "$repair_mtime_project/do-work/queue/REQ-801-future.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_mtime_project/do-work/queue/REQ-801-future.md"
repair_mtime_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_mtime_project")" \
  || fail_case 'repair-req-timestamps future-stamp case returned nonzero'
grep -q '^created_at: 2026-08-10T12:00:00Z$' "$repair_mtime_project/do-work/queue/REQ-801-future.md" \
  || fail_case 'repair-req-timestamps future-stamp case did not rewrite the stamp to the file mtime'
printf '%s' "$repair_mtime_output" \
  | grep -q 'REQ-801-future.md created_at: 2093-01-01T00:00:00Z -> 2026-08-10T12:00:00Z (file mtime)' \
  || fail_case 'repair-req-timestamps future-stamp case did not log the correction'

# repair-req-timestamps: impossible orderings in working/ are repaired and
# clamped so created_at <= claimed_at <= completed_at <= now — here the derived
# mtime precedes created_at, so both later fields land exactly on the clamp floor.
repair_order_project="$fixture_root/repair-order-project"
mkdir -p "$repair_order_project/do-work/working"
printf -- '---\nid: REQ-802\nstatus: completed\ncreated_at: 2026-08-10T12:00:00Z\nclaimed_at: 2026-08-01T09:00:00Z\ncompleted_at: 2026-08-03T10:00:00Z\n---\nbody\n' \
  > "$repair_order_project/do-work/working/REQ-802-order.md"
TZ=UTC touch -m -t 202608050800.00 "$repair_order_project/do-work/working/REQ-802-order.md"
repair_order_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_order_project")" \
  || fail_case 'repair-req-timestamps ordering case returned nonzero'
grep -q '^claimed_at: 2026-08-10T12:00:00Z$' "$repair_order_project/do-work/working/REQ-802-order.md" \
  || fail_case 'repair-req-timestamps ordering case did not clamp claimed_at up to created_at'
grep -q '^completed_at: 2026-08-10T12:00:00Z$' "$repair_order_project/do-work/working/REQ-802-order.md" \
  || fail_case 'repair-req-timestamps ordering case did not clamp completed_at up to the repaired claimed_at'
printf '%s' "$repair_order_output" | grep -q 'clamped to 2026-08-10T12:00:00Z' \
  || fail_case 'repair-req-timestamps ordering case did not log the clamp'

# repair-req-timestamps: a committed file that matches HEAD is repaired to the
# author time of the commit that introduced the stamp line, not to a fresh clone's
# meaningless mtime.
repair_blame_project="$fixture_root/repair-blame-project"
fixture_repo_init "$repair_blame_project"
mkdir -p "$repair_blame_project/do-work/queue"
printf -- '---\nid: REQ-803\nstatus: pending\ncreated_at: 2026-08-14T09:00:00Z\nclaimed_at: "2093-02-02T02:02:02Z"\n---\nbody\n' \
  > "$repair_blame_project/do-work/queue/REQ-803-committed.md"
git -C "$repair_blame_project" add -A
GIT_AUTHOR_DATE='2026-08-15T14:00:00Z' GIT_COMMITTER_DATE='2026-08-15T14:05:00Z' \
  git -C "$repair_blame_project" commit -qm fixture
repair_blame_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_blame_project")" \
  || fail_case 'repair-req-timestamps committed-file case returned nonzero'
grep -q '^claimed_at: 2026-08-15T14:00:00Z$' "$repair_blame_project/do-work/queue/REQ-803-committed.md" \
  || fail_case 'repair-req-timestamps committed-file case did not use the introducing commit author time (quoted stamps must be repairable too)'
printf '%s' "$repair_blame_output" | grep -q 'author time' \
  || fail_case 'repair-req-timestamps committed-file case did not name the commit author time as the replacement source'

# repair-req-timestamps: a clean fixture passes through byte-identical — including
# the shapes the repairer must not touch: a nested (indented) calculated_at, a
# numeric-offset value it refuses permanently because the timezone arithmetic is
# the risk and not the obstacle (REQ-257), and an archive-scope directory it must
# never scan.
repair_clean_project="$fixture_root/repair-clean-project"
mkdir -p "$repair_clean_project/do-work/queue" "$repair_clean_project/do-work/archive"
printf -- '---\nid: REQ-804\nstatus: pending\ncreated_at: 2026-08-10T12:00:00Z   # trailing comment\nblocked_at: 2026-08-11T09:00:00+09:00\nestimate:\n  calculated_at: 2093-06-06T06:06:06Z\n---\nbody\n' \
  > "$repair_clean_project/do-work/queue/REQ-804-clean.md"
printf -- '---\nid: REQ-805\nstatus: completed\ncreated_at: 2093-03-03T03:03:03Z\n---\nbody\n' \
  > "$repair_clean_project/do-work/archive/REQ-805-archived.md"
cp "$repair_clean_project/do-work/queue/REQ-804-clean.md" "$fixture_root/repair-clean-before.md"
cp "$repair_clean_project/do-work/archive/REQ-805-archived.md" "$fixture_root/repair-archive-before.md"
repair_clean_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_clean_project")" \
  || fail_case 'repair-req-timestamps clean-fixture case returned nonzero'
[ -z "$repair_clean_output" ] \
  || fail_case 'repair-req-timestamps clean-fixture case printed output for a no-op run'
cmp -s "$fixture_root/repair-clean-before.md" "$repair_clean_project/do-work/queue/REQ-804-clean.md" \
  || fail_case 'repair-req-timestamps clean-fixture case changed a file with nothing provably wrong'
cmp -s "$fixture_root/repair-archive-before.md" "$repair_clean_project/do-work/archive/REQ-805-archived.md" \
  || fail_case 'repair-req-timestamps clean-fixture case wrote into archive scope'

# repair-req-timestamps: a tripped guard leaves the file byte-identical and exits
# nonzero — here the truncation floor: a file at less than half its committed size
# lost content before the run, and repairing a stamp in the remains would help
# commit the loss.
repair_guard_project="$fixture_root/repair-guard-project"
fixture_repo_init "$repair_guard_project"
mkdir -p "$repair_guard_project/do-work/queue"
{
  printf -- '---\nid: REQ-806\nstatus: pending\ncreated_at: 2093-04-04T04:04:04Z\n---\n'
  awk 'BEGIN { for (line_index = 1; line_index <= 200; line_index++) print "ballast decision-trail line " line_index }'
} > "$repair_guard_project/do-work/queue/REQ-806-truncated.md"
fixture_repo_commit_all "$repair_guard_project" fixture
head -n 5 "$repair_guard_project/do-work/queue/REQ-806-truncated.md" > "$fixture_root/repair-truncated.tmp"
mv "$fixture_root/repair-truncated.tmp" "$repair_guard_project/do-work/queue/REQ-806-truncated.md"
cp "$repair_guard_project/do-work/queue/REQ-806-truncated.md" "$fixture_root/repair-guard-before.md"
repair_guard_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_guard_project")" \
  && fail_case 'repair-req-timestamps tripped-guard case exited zero on a truncated file'
cmp -s "$fixture_root/repair-guard-before.md" "$repair_guard_project/do-work/queue/REQ-806-truncated.md" \
  || fail_case 'repair-req-timestamps tripped-guard case modified the truncated file'
printf '%s' "$repair_guard_output" | grep -q 'content was lost before this run' \
  || fail_case 'repair-req-timestamps tripped-guard case did not name the truncation as the reason'

# repair-req-timestamps: an unquoted space-separated future instant is repaired
# whole — never half-rewritten into a date plus a phantom time-of-day — and the
# audit line reports the full old value (REQ-255 I1, the corrupting shape).
repair_space_project="$fixture_root/repair-space-project"
mkdir -p "$repair_space_project/do-work/queue"
printf -- '---\nid: REQ-807\nstatus: pending\ncreated_at: 2093-01-01 00:00:00\n---\nbody\n' \
  > "$repair_space_project/do-work/queue/REQ-807-space.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_space_project/do-work/queue/REQ-807-space.md"
repair_space_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_space_project")" \
  || fail_case 'repair-req-timestamps space-separated case returned nonzero'
grep -q '^created_at: 2026-08-10T12:00:00Z$' "$repair_space_project/do-work/queue/REQ-807-space.md" \
  || fail_case 'repair-req-timestamps space-separated case did not rewrite the whole value to the canonical form'
grep -q '00:00:00$' "$repair_space_project/do-work/queue/REQ-807-space.md" \
  && fail_case 'repair-req-timestamps space-separated case left a phantom time-of-day suffix behind'
printf '%s' "$repair_space_output" \
  | grep -q 'REQ-807-space.md created_at: 2093-01-01 00:00:00 -> 2026-08-10T12:00:00Z (file mtime)' \
  || fail_case 'repair-req-timestamps space-separated case did not report the full old value in the audit line'

# repair-req-timestamps: a quoted space-separated future instant is repaired to
# the canonical unquoted form instead of silently truncating at the unmatched
# quote and passing through (REQ-255 I1's quoted sibling).
repair_quoted_space_project="$fixture_root/repair-quoted-space-project"
mkdir -p "$repair_quoted_space_project/do-work/queue"
printf -- '---\nid: REQ-808\nstatus: pending\ncreated_at: "2093-01-01 00:00:00"\n---\nbody\n' \
  > "$repair_quoted_space_project/do-work/queue/REQ-808-quoted-space.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_quoted_space_project/do-work/queue/REQ-808-quoted-space.md"
repair_quoted_space_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_quoted_space_project")" \
  || fail_case 'repair-req-timestamps quoted-space case returned nonzero'
grep -q '^created_at: 2026-08-10T12:00:00Z$' "$repair_quoted_space_project/do-work/queue/REQ-808-quoted-space.md" \
  || fail_case 'repair-req-timestamps quoted-space case did not repair the quoted instant to the canonical unquoted form'
printf '%s' "$repair_quoted_space_output" | grep -q 'REQ-808-quoted-space.md created_at' \
  || fail_case 'repair-req-timestamps quoted-space case did not log the correction'

# repair-req-timestamps: a CRLF-fenced file is scanned like the board scans it,
# and a repair preserves every line's CRLF ending — Windows agents are the
# likeliest source of both CRLF files and wrong local-time stamps (REQ-255 I2).
repair_crlf_project="$fixture_root/repair-crlf-project"
mkdir -p "$repair_crlf_project/do-work/queue"
printf -- '---\r\nid: REQ-809\r\nstatus: pending\r\ncreated_at: 2093-03-03T03:03:03Z\r\n---\r\nbody\r\n' \
  > "$repair_crlf_project/do-work/queue/REQ-809-crlf.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_crlf_project/do-work/queue/REQ-809-crlf.md"
repair_crlf_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_crlf_project")" \
  || fail_case 'repair-req-timestamps CRLF case returned nonzero'
grep -q $'^created_at: 2026-08-10T12:00:00Z\r$' "$repair_crlf_project/do-work/queue/REQ-809-crlf.md" \
  || fail_case 'repair-req-timestamps CRLF case did not repair the stamp behind the CRLF fence (or dropped the CR)'
[ "$(grep -c $'\r$' "$repair_crlf_project/do-work/queue/REQ-809-crlf.md")" -eq 6 ] \
  || fail_case 'repair-req-timestamps CRLF case did not preserve every CRLF line ending'
printf '%s' "$repair_crlf_output" | grep -q 'REQ-809-crlf.md created_at: 2093-03-03T03:03:03Z -> 2026-08-10T12:00:00Z' \
  || fail_case 'repair-req-timestamps CRLF case did not log the correction'

# repair-req-timestamps: a BOM-prefixed file is scanned like the board scans it
# (the board strips the BOM before the fence match), and a repair keeps the BOM
# bytes in place (REQ-255 I2).
repair_bom_project="$fixture_root/repair-bom-project"
mkdir -p "$repair_bom_project/do-work/queue"
printf -- '\xef\xbb\xbf---\nid: REQ-810\nstatus: pending\ncreated_at: 2093-04-04T04:04:04Z\n---\nbody\n' \
  > "$repair_bom_project/do-work/queue/REQ-810-bom.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_bom_project/do-work/queue/REQ-810-bom.md"
repair_bom_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_bom_project")" \
  || fail_case 'repair-req-timestamps BOM case returned nonzero'
grep -q '^created_at: 2026-08-10T12:00:00Z$' "$repair_bom_project/do-work/queue/REQ-810-bom.md" \
  || fail_case 'repair-req-timestamps BOM case did not repair the stamp behind the BOM-prefixed fence'
[ "$(head -c 3 "$repair_bom_project/do-work/queue/REQ-810-bom.md")" = "$(printf '\xef\xbb\xbf')" ] \
  || fail_case 'repair-req-timestamps BOM case did not keep the BOM bytes in place'
printf '%s' "$repair_bom_output" | grep -q 'REQ-810-bom.md created_at: 2093-04-04T04:04:04Z -> 2026-08-10T12:00:00Z' \
  || fail_case 'repair-req-timestamps BOM case did not log the correction'

# repair-req-timestamps: a shape-valid but calendar-impossible stamp is left
# byte-identical for diagnosis — the board's parser rejects it, so erasing it
# to a derived instant would destroy the malformed evidence while claiming
# parity (REQ-255, PR #145 external review). The range check must match the
# read side's real calendar: month, day-in-that-month, leap years, and time
# components — a real leap-day future stamp must still be repaired.
repair_calendar_project="$fixture_root/repair-calendar-project"
mkdir -p "$repair_calendar_project/do-work/queue"
printf -- '---\nid: REQ-811\nstatus: pending\ncreated_at: 9999-99-99T99:99:99Z\n---\nbody\n' \
  > "$repair_calendar_project/do-work/queue/REQ-811-impossible.md"
printf -- '---\nid: REQ-812\nstatus: pending\ncreated_at: 2093-04-31T10:00:00Z\n---\nbody\n' \
  > "$repair_calendar_project/do-work/queue/REQ-812-april-31.md"
printf -- '---\nid: REQ-813\nstatus: pending\ncreated_at: 2093-02-29T10:00:00Z\n---\nbody\n' \
  > "$repair_calendar_project/do-work/queue/REQ-813-not-a-leap-year.md"
printf -- '---\nid: REQ-814\nstatus: pending\ncreated_at: 2092-02-29T10:00:00Z\n---\nbody\n' \
  > "$repair_calendar_project/do-work/queue/REQ-814-real-leap-day.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_calendar_project/do-work/queue/"REQ-81*.md
for impossible_fixture in REQ-811-impossible REQ-812-april-31 REQ-813-not-a-leap-year; do
  cp "$repair_calendar_project/do-work/queue/$impossible_fixture.md" "$fixture_root/$impossible_fixture-before.md"
done
repair_calendar_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_calendar_project")" \
  || fail_case 'repair-req-timestamps calendar case returned nonzero'
for impossible_fixture in REQ-811-impossible REQ-812-april-31 REQ-813-not-a-leap-year; do
  cmp -s "$fixture_root/$impossible_fixture-before.md" "$repair_calendar_project/do-work/queue/$impossible_fixture.md" \
    || fail_case "repair-req-timestamps calendar case erased the impossible stamp in $impossible_fixture instead of leaving it for diagnosis"
done
grep -q '^created_at: 2026-08-10T12:00:00Z$' "$repair_calendar_project/do-work/queue/REQ-814-real-leap-day.md" \
  || fail_case 'repair-req-timestamps calendar case refused a real leap-day future stamp the board parses'
printf '%s' "$repair_calendar_output" | grep -q 'REQ-811\|REQ-812\|REQ-813' \
  && fail_case 'repair-req-timestamps calendar case logged a correction for a value it must not touch'

# repair-req-timestamps: a numeric UTC offset or fractional seconds is refused
# permanently — the decided answer to the residual, not a to-do (REQ-257). Every
# such shape the read side parses and future-badges is left byte-identical, is
# never logged, and does not fail the run. This case is what fails if someone
# quietly teaches comparison_key_for offset arithmetic: a refusal that lived only
# in a header comment pinned nothing, and the risk being refused is real —
# repairing here would rewrite an instant the read side already reads correctly.
repair_refused_shape_project="$fixture_root/repair-refused-shape-project"
mkdir -p "$repair_refused_shape_project/do-work/queue"
printf -- '---\nid: REQ-817\nstatus: pending\ncreated_at: 2093-01-01T00:00:00+02:00\n---\nbody\n' \
  > "$repair_refused_shape_project/do-work/queue/REQ-817-offset-ahead.md"
printf -- '---\nid: REQ-818\nstatus: pending\ncreated_at: 2093-01-01T00:00:00-05:00\n---\nbody\n' \
  > "$repair_refused_shape_project/do-work/queue/REQ-818-offset-behind.md"
printf -- '---\nid: REQ-819\nstatus: pending\ncreated_at: "2093-01-01T00:00:00+02:00"\n---\nbody\n' \
  > "$repair_refused_shape_project/do-work/queue/REQ-819-quoted-offset.md"
printf -- '---\nid: REQ-820\nstatus: pending\ncreated_at: 2093-01-01T00:00:00.500Z\n---\nbody\n' \
  > "$repair_refused_shape_project/do-work/queue/REQ-820-fractional-zulu.md"
printf -- '---\nid: REQ-821\nstatus: pending\ncreated_at: 2093-01-01T00:00:00.5\n---\nbody\n' \
  > "$repair_refused_shape_project/do-work/queue/REQ-821-fractional-bare.md"
printf -- '---\nid: REQ-822\nstatus: pending\ncreated_at: 2093-01-01 00:00:00.5\n---\nbody\n' \
  > "$repair_refused_shape_project/do-work/queue/REQ-822-fractional-space.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_refused_shape_project/do-work/queue/"REQ-8[12]*.md
for refused_fixture in "$repair_refused_shape_project/do-work/queue/"REQ-8[12]*.md; do
  cp "$refused_fixture" "$fixture_root/refused-$(basename "$refused_fixture")"
done
repair_refused_shape_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_refused_shape_project")" \
  || fail_case 'repair-req-timestamps refused-shape case returned nonzero for shapes it deliberately does not repair'
[ -z "$repair_refused_shape_output" ] \
  || fail_case 'repair-req-timestamps refused-shape case logged a correction for an offset or fractional stamp'
for refused_fixture in "$repair_refused_shape_project/do-work/queue/"REQ-8[12]*.md; do
  cmp -s "$fixture_root/refused-$(basename "$refused_fixture")" "$refused_fixture" \
    || fail_case "repair-req-timestamps refused-shape case rewrote $(basename "$refused_fixture") — the offset/fractional refusal is a decided permanent answer; changing it means re-deciding it, not widening comparison_key_for"
done

# repair-req-timestamps: the read-side layout list the refusal is scoped against
# stays in lock-step — the residual is the GAP between the board's parser and
# this script's, so pinning only the refusing side would let a widened board grow
# the gap silently. A new layout here means the offset/fractional decision (and
# the header paragraph stating it) gets re-made, not inherited.
board_timestamp_layouts="$(awk '/^func parseTimestamp\(/, /^}/' \
  "$repo_root/skills/do-work-board/tools/queue-kanban/model.go" \
  | sed -n -E 's/^[[:space:]]*(time\.RFC3339|"2006[^"]*"),$/\1/p' | tr '\n' ' ')"
[ "$board_timestamp_layouts" = 'time.RFC3339 "2006-01-02T15:04:05" "2006-01-02 15:04:05" "2006-01-02" ' ] \
  || fail_case "repair-req-timestamps read-side-layout case: the board's parseTimestamp layouts are now [$board_timestamp_layouts] — re-decide what repair-req-timestamps.sh refuses before widening the read side"

# repair-req-timestamps: a duplicated anchor key follows the last occurrence,
# exactly like the read side (the board's YAML dedup keeps the LAST value of a
# repeated top-level key) — a later-then-earlier claimed_at pair is a real
# ordering defect on the board and must not be reported clean; and a future
# FIRST occurrence shadowed by a clean last one is invisible to every YAML
# reader and must stay untouched (REQ-255, PR #145 external review).
repair_duplicate_project="$fixture_root/repair-duplicate-project"
mkdir -p "$repair_duplicate_project/do-work/working"
printf -- '---\nid: REQ-815\nstatus: claimed\ncreated_at: 2026-08-10T12:00:00Z\nclaimed_at: 2026-08-11T12:00:00Z\nclaimed_at: 2026-08-01T09:00:00Z\n---\nbody\n' \
  > "$repair_duplicate_project/do-work/working/REQ-815-duplicate-anchor.md"
printf -- '---\nid: REQ-816\nstatus: pending\ncreated_at: 2026-08-01T09:00:00Z\nblocked_at: 2093-01-01T00:00:00Z\nblocked_at: 2026-08-02T09:00:00Z\n---\nbody\n' \
  > "$repair_duplicate_project/do-work/working/REQ-816-shadowed-first.md"
TZ=UTC touch -m -t 202608121200.00 "$repair_duplicate_project/do-work/working/"REQ-81*.md
cp "$repair_duplicate_project/do-work/working/REQ-816-shadowed-first.md" "$fixture_root/repair-shadowed-before.md"
repair_duplicate_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_duplicate_project")" \
  || fail_case 'repair-req-timestamps duplicate-anchor case returned nonzero'
grep -q '^claimed_at: 2026-08-12T12:00:00Z$' "$repair_duplicate_project/do-work/working/REQ-815-duplicate-anchor.md" \
  || fail_case 'repair-req-timestamps duplicate-anchor case reported clean instead of repairing the effective (last) occurrence'
grep -q '^claimed_at: 2026-08-11T12:00:00Z$' "$repair_duplicate_project/do-work/working/REQ-815-duplicate-anchor.md" \
  || fail_case 'repair-req-timestamps duplicate-anchor case rewrote the shadowed first occurrence'
cmp -s "$fixture_root/repair-shadowed-before.md" "$repair_duplicate_project/do-work/working/REQ-816-shadowed-first.md" \
  || fail_case 'repair-req-timestamps duplicate-anchor case touched a future first occurrence no YAML reader can see'
printf '%s' "$repair_duplicate_output" | grep -q 'REQ-815-duplicate-anchor.md claimed_at: 2026-08-01T09:00:00Z -> 2026-08-12T12:00:00Z' \
  || fail_case 'repair-req-timestamps duplicate-anchor case did not log the effective-occurrence correction'

# repair-req-timestamps: a file whose opening fence is never closed is refused
# whole, because the board's splitFrontmatter reports NO frontmatter for that
# shape and renders every line as body (REQ-267 I1). The second fixture is the
# shape that could wedge the repair permanently: the file ends on the defective
# stamp with no trailing newline, so a last-line rewrite can never produce the
# final-newline diff pair the changed-line guard expects — the repair was
# rejected and the script exited 1 on EVERY run, with nothing able to heal it.
# The run must therefore be clean and silent, not merely non-destructive.
repair_unterminated_project="$fixture_root/repair-unterminated-project"
mkdir -p "$repair_unterminated_project/do-work/queue"
printf -- '---\nid: REQ-823\nstatus: pending\n\n# Body prose, and the fence above was never closed\n\ncreated_at: 2093-01-01T00:00:00Z\n' \
  > "$repair_unterminated_project/do-work/queue/REQ-823-unterminated-body.md"
printf -- '---\nid: REQ-824\nstatus: pending\ncreated_at: 2093-01-01T00:00:00Z' \
  > "$repair_unterminated_project/do-work/queue/REQ-824-unterminated-eof.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_unterminated_project/do-work/queue/"REQ-82*.md
for unterminated_fixture in REQ-823-unterminated-body REQ-824-unterminated-eof; do
  cp "$repair_unterminated_project/do-work/queue/$unterminated_fixture.md" "$fixture_root/$unterminated_fixture-before.md"
done
repair_unterminated_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_unterminated_project")" \
  || fail_case 'repair-req-timestamps unterminated-fence case exited nonzero — a shape that fails every session with no self-heal is exactly the defect'
[ -z "$repair_unterminated_output" ] \
  || fail_case 'repair-req-timestamps unterminated-fence case printed output for a file the read side sees no frontmatter in'
for unterminated_fixture in REQ-823-unterminated-body REQ-824-unterminated-eof; do
  cmp -s "$fixture_root/$unterminated_fixture-before.md" "$repair_unterminated_project/do-work/queue/$unterminated_fixture.md" \
    || fail_case "repair-req-timestamps unterminated-fence case rewrote $unterminated_fixture — the board reads that whole file as body text"
done

# repair-req-timestamps: a stamp padded INSIDE its quotes is repaired, because
# the read side unquotes and then trims and so parses and future-badges it
# (REQ-267 I2). The non-ASCII-padded sibling pins the stated boundary of that
# trim — this script matches bytes under LC_ALL=C, so a U+00A0-padded value the
# read side still parses stays refused byte-identical, and the header says so.
repair_padded_quote_project="$fixture_root/repair-padded-quote-project"
mkdir -p "$repair_padded_quote_project/do-work/queue"
printf -- '---\nid: REQ-825\nstatus: pending\ncreated_at: "2093-01-01 00:00:00 "\n---\nbody\n' \
  > "$repair_padded_quote_project/do-work/queue/REQ-825-padded-quote.md"
printf -- '---\nid: REQ-826\nstatus: pending\ncreated_at: "2093-01-01T00:00:00Z\xc2\xa0"\n---\nbody\n' \
  > "$repair_padded_quote_project/do-work/queue/REQ-826-non-ascii-pad.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_padded_quote_project/do-work/queue/"REQ-82*.md
cp "$repair_padded_quote_project/do-work/queue/REQ-826-non-ascii-pad.md" "$fixture_root/repair-non-ascii-pad-before.md"
repair_padded_quote_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_padded_quote_project")" \
  || fail_case 'repair-req-timestamps padded-quote case returned nonzero'
grep -q '^created_at: 2026-08-10T12:00:00Z$' "$repair_padded_quote_project/do-work/queue/REQ-825-padded-quote.md" \
  || fail_case 'repair-req-timestamps padded-quote case refused a padded quoted instant the board parses and future-badges'
cmp -s "$fixture_root/repair-non-ascii-pad-before.md" "$repair_padded_quote_project/do-work/queue/REQ-826-non-ascii-pad.md" \
  || fail_case 'repair-req-timestamps padded-quote case repaired a non-ASCII-padded value — the header states that residual is refused'
printf '%s' "$repair_padded_quote_output" | grep -q 'REQ-825-padded-quote.md created_at: "2093-01-01 00:00:00 " -> 2026-08-10T12:00:00Z' \
  || fail_case 'repair-req-timestamps padded-quote case did not report the full padded old value in the audit line'
printf '%s' "$repair_padded_quote_output" | grep -q 'REQ-826' \
  && fail_case 'repair-req-timestamps padded-quote case logged a correction for the refused non-ASCII-padded value'

# audit-archive-timestamps: under --fix, a future stamp in a committed archived REQ
# (inside an archived UR folder, proving the recursive scan) is rewritten to the
# introducing commit's author time and the correction logs the sourcing commit hash.
audit_fix_project="$fixture_root/audit-fix-project"
fixture_repo_init "$audit_fix_project"
mkdir -p "$audit_fix_project/do-work/archive/UR-901"
printf -- '---\nid: REQ-901\nstatus: completed\ncreated_at: 2026-08-10T09:00:00Z\nclaimed_at: 2026-08-10T10:00:00Z\ncompleted_at: 2093-05-05T05:05:05Z\n---\nbody\n' \
  > "$audit_fix_project/do-work/archive/UR-901/REQ-901-future.md"
git -C "$audit_fix_project" add -A
GIT_AUTHOR_DATE='2026-08-12T10:00:00Z' GIT_COMMITTER_DATE='2026-08-12T10:05:00Z' \
  git -C "$audit_fix_project" commit -qm fixture
audit_fix_output="$("$core_scripts/audit-archive-timestamps.sh" --fix "$audit_fix_project")" \
  || fail_case 'audit-archive-timestamps fixing case returned nonzero'
grep -q '^completed_at: 2026-08-12T10:00:00Z$' "$audit_fix_project/do-work/archive/UR-901/REQ-901-future.md" \
  || fail_case 'audit-archive-timestamps fixing case did not rewrite the stamp to the introducing commit author time'
printf '%s' "$audit_fix_output" \
  | grep -Eq 'REQ-901-future\.md completed_at: 2093-05-05T05:05:05Z -> 2026-08-12T10:00:00Z \(commit [0-9a-f]{7} author time\)' \
  || fail_case 'audit-archive-timestamps fixing case did not log the correction with its sourcing commit hash'

# audit-archive-timestamps: without --fix the default run only reports — the pending
# correction prints as `would repair`, the archived file keeps its bytes, and the
# exit code is nonzero so a caller can gate on findings.
audit_report_project="$fixture_root/audit-report-project"
fixture_repo_init "$audit_report_project"
mkdir -p "$audit_report_project/do-work/archive"
printf -- '---\nid: REQ-902\nstatus: completed\ncreated_at: 2026-08-10T09:00:00Z\ncompleted_at: 2093-05-05T05:05:05Z\n---\nbody\n' \
  > "$audit_report_project/do-work/archive/REQ-902-future.md"
git -C "$audit_report_project" add -A
GIT_AUTHOR_DATE='2026-08-12T10:00:00Z' GIT_COMMITTER_DATE='2026-08-12T10:05:00Z' \
  git -C "$audit_report_project" commit -qm fixture
cp "$audit_report_project/do-work/archive/REQ-902-future.md" "$fixture_root/audit-report-before.md"
audit_report_output="$("$core_scripts/audit-archive-timestamps.sh" "$audit_report_project")" \
  && fail_case 'audit-archive-timestamps report-only case exited zero with a correction pending'
printf '%s' "$audit_report_output" \
  | grep -q 'would repair do-work/archive/REQ-902-future.md completed_at: 2093-05-05T05:05:05Z -> 2026-08-12T10:00:00Z' \
  || fail_case 'audit-archive-timestamps report-only case did not print the pending correction'
cmp -s "$fixture_root/audit-report-before.md" "$audit_report_project/do-work/archive/REQ-902-future.md" \
  || fail_case 'audit-archive-timestamps report-only case changed bytes without --fix'

# audit-archive-timestamps: a clean committed archive passes through byte-identical
# with exit zero — and queue scope is never scanned: a wrong queue stamp belongs to
# the hook-run repairer, not the archive audit.
audit_clean_project="$fixture_root/audit-clean-project"
fixture_repo_init "$audit_clean_project"
mkdir -p "$audit_clean_project/do-work/archive" "$audit_clean_project/do-work/queue"
printf -- '---\nid: REQ-903\nstatus: completed\ncreated_at: 2026-08-01T09:00:00Z\nclaimed_at: 2026-08-02T09:00:00Z\ncompleted_at: 2026-08-03T09:00:00Z\n---\nbody\n' \
  > "$audit_clean_project/do-work/archive/REQ-903-clean.md"
printf -- '---\nid: REQ-904\nstatus: pending\ncreated_at: 2093-06-06T06:06:06Z\n---\nbody\n' \
  > "$audit_clean_project/do-work/queue/REQ-904-queue.md"
git -C "$audit_clean_project" add -A
GIT_AUTHOR_DATE='2026-08-04T09:00:00Z' GIT_COMMITTER_DATE='2026-08-04T09:00:00Z' \
  git -C "$audit_clean_project" commit -qm fixture
cp "$audit_clean_project/do-work/archive/REQ-903-clean.md" "$fixture_root/audit-clean-before.md"
cp "$audit_clean_project/do-work/queue/REQ-904-queue.md" "$fixture_root/audit-queue-before.md"
audit_clean_output="$("$core_scripts/audit-archive-timestamps.sh" --fix "$audit_clean_project")" \
  || fail_case 'audit-archive-timestamps clean-archive case returned nonzero'
printf '%s' "$audit_clean_output" | grep -q 'archive audit clean' \
  || fail_case 'audit-archive-timestamps clean-archive case did not report a clean audit'
cmp -s "$fixture_root/audit-clean-before.md" "$audit_clean_project/do-work/archive/REQ-903-clean.md" \
  || fail_case 'audit-archive-timestamps clean-archive case changed a clean archived file'
cmp -s "$fixture_root/audit-queue-before.md" "$audit_clean_project/do-work/queue/REQ-904-queue.md" \
  || fail_case 'audit-archive-timestamps clean-archive case wrote into queue scope'

# audit-archive-timestamps: an impossible ordering in a committed archived REQ is
# clamped so created_at <= claimed_at <= completed_at — the introducing commit's
# author time precedes the anchor here, so both later fields land on the clamp floor.
audit_order_project="$fixture_root/audit-order-project"
fixture_repo_init "$audit_order_project"
mkdir -p "$audit_order_project/do-work/archive"
printf -- '---\nid: REQ-905\nstatus: completed\ncreated_at: 2026-08-10T12:00:00Z\nclaimed_at: 2026-08-01T09:00:00Z\ncompleted_at: 2026-08-03T10:00:00Z\n---\nbody\n' \
  > "$audit_order_project/do-work/archive/REQ-905-order.md"
git -C "$audit_order_project" add -A
GIT_AUTHOR_DATE='2026-08-05T08:00:00Z' GIT_COMMITTER_DATE='2026-08-05T08:00:00Z' \
  git -C "$audit_order_project" commit -qm fixture
audit_order_output="$("$core_scripts/audit-archive-timestamps.sh" --fix "$audit_order_project")" \
  || fail_case 'audit-archive-timestamps ordering case returned nonzero'
grep -q '^claimed_at: 2026-08-10T12:00:00Z$' "$audit_order_project/do-work/archive/REQ-905-order.md" \
  || fail_case 'audit-archive-timestamps ordering case did not clamp claimed_at up to created_at'
grep -q '^completed_at: 2026-08-10T12:00:00Z$' "$audit_order_project/do-work/archive/REQ-905-order.md" \
  || fail_case 'audit-archive-timestamps ordering case did not clamp completed_at up to the repaired claimed_at'
printf '%s' "$audit_order_output" | grep -q 'clamped to 2026-08-10T12:00:00Z' \
  || fail_case 'audit-archive-timestamps ordering case did not log the clamp'

# audit-archive-timestamps: an archived defect whose introducing commit cannot be
# blamed (an uncommitted file) is reported and left byte-identical — replacements
# derive from git alone, and the file mtime is never consulted as a fallback.
audit_blameless_project="$fixture_root/audit-blameless-project"
fixture_repo_init "$audit_blameless_project"
mkdir -p "$audit_blameless_project/do-work/archive"
printf -- '---\nid: REQ-906\nstatus: completed\ncreated_at: 2093-07-07T07:07:07Z\n---\nbody\n' \
  > "$audit_blameless_project/do-work/archive/REQ-906-untracked.md"
TZ=UTC touch -m -t 202608101200.00 "$audit_blameless_project/do-work/archive/REQ-906-untracked.md"
cp "$audit_blameless_project/do-work/archive/REQ-906-untracked.md" "$fixture_root/audit-blameless-before.md"
audit_blameless_output="$("$core_scripts/audit-archive-timestamps.sh" --fix "$audit_blameless_project")" \
  && fail_case 'audit-archive-timestamps blameless case exited zero on an unrepairable defect'
printf '%s' "$audit_blameless_output" | grep -q 'FAILED to repair' \
  || fail_case 'audit-archive-timestamps blameless case did not report the unrepairable defect'
printf '%s' "$audit_blameless_output" | grep -q 'file mtime' \
  && fail_case 'audit-archive-timestamps blameless case offered the file mtime as a source'
cmp -s "$fixture_root/audit-blameless-before.md" "$audit_blameless_project/do-work/archive/REQ-906-untracked.md" \
  || fail_case 'audit-archive-timestamps blameless case changed the file it could not derive for'

# audit-archive-timestamps: the widened shapes (space-separated, quoted, CRLF,
# BOM) repair through the archive scan too — the fix lives in the sourced
# library, and this pins the shared-fix-reaches-both-tools property instead of
# assuming it (REQ-255; REQ-247 review).
audit_shapes_project="$fixture_root/audit-shapes-project"
fixture_repo_init "$audit_shapes_project"
mkdir -p "$audit_shapes_project/do-work/archive"
printf -- '---\nid: REQ-907\nstatus: completed\ncreated_at: 2093-01-01 00:00:00\n---\nbody\n' \
  > "$audit_shapes_project/do-work/archive/REQ-907-space.md"
printf -- '---\nid: REQ-908\nstatus: completed\ncreated_at: "2093-01-01 00:00:00"\n---\nbody\n' \
  > "$audit_shapes_project/do-work/archive/REQ-908-quoted-space.md"
printf -- '---\r\nid: REQ-909\r\nstatus: completed\r\ncreated_at: 2093-03-03T03:03:03Z\r\n---\r\nbody\r\n' \
  > "$audit_shapes_project/do-work/archive/REQ-909-crlf.md"
printf -- '\xef\xbb\xbf---\nid: REQ-910\nstatus: completed\ncreated_at: 2093-04-04T04:04:04Z\n---\nbody\n' \
  > "$audit_shapes_project/do-work/archive/REQ-910-bom.md"
git -C "$audit_shapes_project" add -A
GIT_AUTHOR_DATE='2026-08-12T10:00:00Z' GIT_COMMITTER_DATE='2026-08-12T10:05:00Z' \
  git -C "$audit_shapes_project" commit -qm fixture
audit_shapes_output="$("$core_scripts/audit-archive-timestamps.sh" --fix "$audit_shapes_project")" \
  || fail_case 'audit-archive-timestamps widened-shapes case returned nonzero'
for widened_fixture in REQ-907-space REQ-908-quoted-space REQ-910-bom; do
  grep -q '^created_at: 2026-08-12T10:00:00Z$' "$audit_shapes_project/do-work/archive/$widened_fixture.md" \
    || fail_case "audit-archive-timestamps widened-shapes case did not repair $widened_fixture through the archive scan"
done
grep -q $'^created_at: 2026-08-12T10:00:00Z\r$' "$audit_shapes_project/do-work/archive/REQ-909-crlf.md" \
  || fail_case 'audit-archive-timestamps widened-shapes case did not repair the CRLF file (or dropped the CR)'
[ "$(head -c 3 "$audit_shapes_project/do-work/archive/REQ-910-bom.md")" = "$(printf '\xef\xbb\xbf')" ] \
  || fail_case 'audit-archive-timestamps widened-shapes case did not keep the BOM bytes in place'
printf '%s' "$audit_shapes_output" \
  | grep -q 'REQ-907-space.md created_at: 2093-01-01 00:00:00 -> 2026-08-12T10:00:00Z' \
  || fail_case 'audit-archive-timestamps widened-shapes case did not report the full old value in the audit line'

# audit-archive-timestamps: the refused and duplicate-key shapes hold through
# the archive scan too — a calendar-impossible stamp and a numeric-offset stamp
# stay byte-identical (and are not defects), while a duplicated anchor repairs
# on its effective (last) occurrence from the introducing commit's author time
# (REQ-255; REQ-257). The refusals belong to the sourced library, so this is
# what fails if the auditor ever grows its own recognizer beside the shared one.
audit_parity_project="$fixture_root/audit-parity-project"
fixture_repo_init "$audit_parity_project"
mkdir -p "$audit_parity_project/do-work/archive"
printf -- '---\nid: REQ-911\nstatus: completed\ncreated_at: 9999-99-99T99:99:99Z\n---\nbody\n' \
  > "$audit_parity_project/do-work/archive/REQ-911-impossible.md"
printf -- '---\nid: REQ-912\nstatus: completed\ncreated_at: 2026-08-10T12:00:00Z\nclaimed_at: 2026-08-11T12:00:00Z\nclaimed_at: 2026-08-01T09:00:00Z\n---\nbody\n' \
  > "$audit_parity_project/do-work/archive/REQ-912-duplicate-anchor.md"
printf -- '---\nid: REQ-913\nstatus: completed\ncreated_at: 2093-01-01T00:00:00+02:00\n---\nbody\n' \
  > "$audit_parity_project/do-work/archive/REQ-913-offset.md"
git -C "$audit_parity_project" add -A
GIT_AUTHOR_DATE='2026-08-12T10:00:00Z' GIT_COMMITTER_DATE='2026-08-12T10:05:00Z' \
  git -C "$audit_parity_project" commit -qm fixture
cp "$audit_parity_project/do-work/archive/REQ-911-impossible.md" "$fixture_root/audit-impossible-before.md"
cp "$audit_parity_project/do-work/archive/REQ-913-offset.md" "$fixture_root/audit-offset-before.md"
audit_parity_output="$("$core_scripts/audit-archive-timestamps.sh" --fix "$audit_parity_project")" \
  || fail_case 'audit-archive-timestamps refusal-parity case returned nonzero'
cmp -s "$fixture_root/audit-impossible-before.md" "$audit_parity_project/do-work/archive/REQ-911-impossible.md" \
  || fail_case 'audit-archive-timestamps refusal-parity case erased a calendar-impossible stamp in the archive'
cmp -s "$fixture_root/audit-offset-before.md" "$audit_parity_project/do-work/archive/REQ-913-offset.md" \
  || fail_case 'audit-archive-timestamps refusal-parity case repaired a numeric-offset stamp in the archive — the offset refusal is the sourced library one and must reach every tool built on it'
grep -q '^claimed_at: 2026-08-12T10:00:00Z$' "$audit_parity_project/do-work/archive/REQ-912-duplicate-anchor.md" \
  || fail_case 'audit-archive-timestamps refusal-parity case did not repair the effective (last) anchor occurrence'
grep -q '^claimed_at: 2026-08-11T12:00:00Z$' "$audit_parity_project/do-work/archive/REQ-912-duplicate-anchor.md" \
  || fail_case 'audit-archive-timestamps refusal-parity case rewrote the shadowed first occurrence'
printf '%s' "$audit_parity_output" | grep -qE 'REQ-911|REQ-913' \
  && fail_case 'audit-archive-timestamps refusal-parity case logged a refused stamp as a correction'

# audit-archive-timestamps: the fence and padding shapes reach the archive scan
# too, because both live in the sourced extractor and recognizer rather than in
# the queue/working scan (REQ-267, both instances through the second scope). An
# unterminated fence is refused whole here as well; a stamp padded inside its
# quotes is repaired from the introducing commit's author time.
audit_shape_project="$fixture_root/audit-shape-project"
fixture_repo_init "$audit_shape_project"
mkdir -p "$audit_shape_project/do-work/archive"
printf -- '---\nid: REQ-914\nstatus: completed\ncreated_at: 2093-01-01T00:00:00Z' \
  > "$audit_shape_project/do-work/archive/REQ-914-unterminated.md"
printf -- '---\nid: REQ-915\nstatus: completed\ncreated_at: "2093-01-01 00:00:00 "\n---\nbody\n' \
  > "$audit_shape_project/do-work/archive/REQ-915-padded-quote.md"
git -C "$audit_shape_project" add -A
GIT_AUTHOR_DATE='2026-08-12T10:00:00Z' GIT_COMMITTER_DATE='2026-08-12T10:05:00Z' \
  git -C "$audit_shape_project" commit -qm fixture
cp "$audit_shape_project/do-work/archive/REQ-914-unterminated.md" "$fixture_root/audit-unterminated-before.md"
audit_shape_output="$("$core_scripts/audit-archive-timestamps.sh" --fix "$audit_shape_project")" \
  || fail_case 'audit-archive-timestamps shape-parity case returned nonzero'
cmp -s "$fixture_root/audit-unterminated-before.md" "$audit_shape_project/do-work/archive/REQ-914-unterminated.md" \
  || fail_case 'audit-archive-timestamps shape-parity case rewrote an archived file whose fence never closes'
grep -q '^created_at: 2026-08-12T10:00:00Z$' "$audit_shape_project/do-work/archive/REQ-915-padded-quote.md" \
  || fail_case 'audit-archive-timestamps shape-parity case refused a padded quoted instant in the archive'
printf '%s' "$audit_shape_output" | grep -q 'REQ-914' \
  && fail_case 'audit-archive-timestamps shape-parity case logged a finding for the refused unterminated file'

# repair-req-timestamps: the 2-minute future-skew constant stays in lock-step
# with the board's futureTimestampSkewAllowance — a fourth hand-kept copy of
# the same allowance, pinned the way the repo pins cause-clause pairs
# (REQ-255 rider; REQ-246 review nit).
grep -q '^future_stamp_skew_seconds=120$' "$core_scripts/repair-req-timestamps.sh" \
  || fail_case 'repair-req-timestamps skew-constant case: the repairer no longer declares future_stamp_skew_seconds=120'
grep -q '^const futureTimestampSkewAllowance = 2 \* time\.Minute$' \
  "$repo_root/skills/do-work-board/tools/queue-kanban/model.go" \
  || fail_case 'repair-req-timestamps skew-constant case: the board constant moved or changed — keep the two allowances in lock-step'

# qualify: an output line added to a file that owns its process exit is the file's own
# reporting, not a debug artifact — the scan passes it and names the reason (REQ-254; the
# scan once FAILed a checker's only success line, REQ-244's remediation). Both output
# primitives (print(, console.log) get the same treatment: the condition is the class,
# not the token. Serial and range modes are both exercised — each reads a different diff.
core_checks="$repo_root/skills/do-work/tools/checks"
qualify_repo="$fixture_root/qualify-repo"
fixture_repo_init "$qualify_repo"
mkdir -p "$qualify_repo/src" "$qualify_repo/do-work"
printf '%s\n' '#!/usr/bin/env bash' 'python3 - <<PY' 'import sys' 'if sites_missing:' \
  '    raise SystemExit("site check failed")' 'PY' > "$qualify_repo/site-checker.sh"
printf '%s\n' 'if (writeFailed) {' '  process.exit(1);' '}' > "$qualify_repo/report-writer.js"
printf '%s\n' 'def parse_value(raw_text):' '    return raw_text.strip()' > "$qualify_repo/src/value_parser.py"
fixture_repo_commit_all "$qualify_repo" base
qualify_base_commit="$(git -C "$qualify_repo" rev-parse HEAD)"
printf '%s\n' '## Implementation Summary' \
  '- `site-checker.sh` (modified) — success line' \
  '- `report-writer.js` (modified) — success line' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-950-reporting.md"
printf 'print("site check: all sites cited")\n' >> "$qualify_repo/site-checker.sh"
printf 'console.log("report written");\n' >> "$qualify_repo/report-writer.js"
qualify_reporting_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-950-reporting.md 2>&1)" \
  || fail_case 'qualify reporter-output serial case FAILed a check file on its own success line'
printf '%s' "$qualify_reporting_output" | grep -q 'reporting' \
  || fail_case 'qualify reporter-output serial case did not name the pass reason (reporting)'
fixture_repo_commit_all "$qualify_repo" reporting
qualify_reporting_commit="$(git -C "$qualify_repo" rev-parse HEAD)"
(cd "$qualify_repo" && DO_WORK_DIFF_RANGE="$qualify_base_commit..$qualify_reporting_commit" \
  "$core_checks/qualify.sh" do-work/REQ-950-reporting.md >/dev/null 2>&1) \
  || fail_case 'qualify reporter-output range case FAILed a check file on its own success line'

# qualify: an output primitive added to a file that never ends its own process still FAILs
# as leftover instrumentation, and the FAIL names the file and the reason (REQ-254).
printf '%s\n' '## Implementation Summary' \
  '- `src/value_parser.py` (modified) — parsing tweak' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-951-instrumentation.md"
printf 'print(raw_text)\n' >> "$qualify_repo/src/value_parser.py"
qualify_instrumentation_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-951-instrumentation.md 2>&1)" \
  && fail_case 'qualify instrumentation serial case passed a debug print in library code'
printf '%s' "$qualify_instrumentation_output" | grep -q 'instrumentation' \
  || fail_case 'qualify instrumentation serial case did not name the FAIL reason (instrumentation)'
printf '%s' "$qualify_instrumentation_output" | grep -q 'src/value_parser.py' \
  || fail_case 'qualify instrumentation serial case did not name the offending file'
fixture_repo_commit_all "$qualify_repo" instrumentation
qualify_instrumentation_commit="$(git -C "$qualify_repo" rev-parse HEAD)"
(cd "$qualify_repo" && DO_WORK_DIFF_RANGE="$qualify_reporting_commit..$qualify_instrumentation_commit" \
  "$core_checks/qualify.sh" do-work/REQ-951-instrumentation.md >/dev/null 2>&1) \
  && fail_case 'qualify instrumentation range case passed a debug print in library code'

# qualify: the reporter exemption covers output primitives only — an unfinished-work
# marker (TODO) added to the same reporter file still FAILs (REQ-254; pins the hole a
# file-level exemption would open).
printf '%s\n' '## Implementation Summary' \
  '- `site-checker.sh` (modified) — regex note' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-952-marker.md"
printf '# TODO: tighten the site regex\n' >> "$qualify_repo/site-checker.sh"
qualify_marker_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-952-marker.md 2>&1)" \
  && fail_case 'qualify unfinished-marker case let a reporter file carry a fresh TODO'
printf '%s' "$qualify_marker_output" | grep -q 'debug artifacts' \
  || fail_case 'qualify unfinished-marker case did not report the TODO as a debug artifact'
git -C "$qualify_repo" checkout -q -- site-checker.sh

if [ "$failure_count" -gt 0 ]; then
  exit 1
fi
# One case is one fixture block, and every block opens with a header comment of the shape
# `<script-name>: <what it proves>` at column zero. That shape is the definition, and the
# count below is that shape grepped out of this file at run time — so the reported number
# and the file cannot disagree, and nothing here is a remembered figure.
named_case_count="$(grep -cE '^# [a-z0-9][a-z0-9-]*: ' "$suite_file")"
printf 'Prescribed shell script behavior probes passed (%s named script cases).\n' "$named_case_count"
