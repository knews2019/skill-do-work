#!/usr/bin/env bash
# shellcheck disable=SC2016
# record-commit-hash-guards.sh — behavioral probes for tools/checks/record-commit-hash.sh.
#
# Unlike the rest of _dev/tests, these do not grep prose: they build a throwaway git
# repository, run the real script against real REQ fixtures, and assert on exit codes and
# stdout. The script exists to stop a data-loss bug (six archived REQ files truncated to
# 0 bytes by the Step 9 metadata commit in a consumer repo), so a grep-level ratchet would
# not tell us whether the guards actually fire.
#
# Invoked by _dev/tests/contract-regressions.sh — there is no auto-discovery of test files,
# so a probe added here only runs because that file calls this one.
#
# Exit 0: every probe passed (or the whole suite was skipped for a missing git).
# Exit 1: at least one probe failed; each failure prints a FAIL line naming the probe.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
record_script="$repo_root/tools/checks/record-commit-hash.sh"
fail_count=0

if ! command -v git >/dev/null 2>&1; then
  printf 'SKIP: git unavailable — record-commit-hash guard probes not run.\n'
  exit 0
fi
if [ ! -x "$record_script" ]; then
  printf 'FAIL: tools/checks/record-commit-hash.sh must exist and be executable.\n' >&2
  exit 1
fi

# The fixture repo must not read the developer's global/system git config: a global
# core.hooksPath, commit.gpgsign, or init.templateDir would change what these probes
# measure (or hang them waiting for a passphrase).
export GIT_CONFIG_GLOBAL=/dev/null
export GIT_CONFIG_SYSTEM=/dev/null
export GIT_TERMINAL_PROMPT=0

fixture_root="$(mktemp -d)"
cleanup_fixture() { rm -rf "$fixture_root"; }
trap cleanup_fixture EXIT

probe_output=""
probe_status=0

# Runs the script from inside the fixture repo and captures both streams, so a probe can
# assert on a usage error (stderr) exactly like it asserts on a FAIL line (stdout).
run_record_script() {
  probe_output="$(cd "$fixture_root" && "$record_script" "$@" 2>&1)"
  probe_status=$?
}

assert_status() {
  local expected_status="$1" probe_name="$2"
  if [ "$probe_status" -ne "$expected_status" ]; then
    printf 'FAIL: %s — expected exit %s, got %s. Output:\n%s\n' \
      "$probe_name" "$expected_status" "$probe_status" "$probe_output" >&2
    fail_count=$((fail_count + 1))
  fi
}

assert_output_matches() {
  local pattern_text="$1" probe_name="$2"
  if ! printf '%s' "$probe_output" | grep -Eq -- "$pattern_text"; then
    printf 'FAIL: %s — output did not match /%s/. Output:\n%s\n' \
      "$probe_name" "$pattern_text" "$probe_output" >&2
    fail_count=$((fail_count + 1))
  fi
}

assert_output_not_matches() {
  local pattern_text="$1" probe_name="$2"
  if printf '%s' "$probe_output" | grep -Eq -- "$pattern_text"; then
    printf 'FAIL: %s — output unexpectedly matched /%s/. Output:\n%s\n' \
      "$probe_name" "$pattern_text" "$probe_output" >&2
    fail_count=$((fail_count + 1))
  fi
}

assert_equals() {
  local expected_value="$1" actual_value="$2" probe_name="$3"
  if [ "$expected_value" != "$actual_value" ]; then
    printf 'FAIL: %s — expected %s, got %s.\n' "$probe_name" "$expected_value" "$actual_value" >&2
    fail_count=$((fail_count + 1))
  fi
}

git -c init.defaultBranch=main init -q "$fixture_root"
git -C "$fixture_root" config user.name "Fixture Runner"
git -C "$fixture_root" config user.email "fixture@example.invalid"

# A realistic archived REQ: frontmatter, then a body that QUOTES the schema inside a fenced
# YAML block. That quoted `commit:` is the trap — several real archived REQs contain one, and
# a file-wide sed would rewrite it. Padded past 8 KB so the HEAD size floor has room to work,
# mirroring the 9 KB–26 KB files the incident destroyed.
write_request_fixture() {
  local fixture_path="$fixture_root/$1" commit_field="$2"
  mkdir -p "$(dirname "$fixture_path")"
  {
    printf -- '---\n'
    printf 'id: REQ-1282\n'
    printf 'title: Fixture request\n'
    printf 'status: completed\n'
    printf 'user_request: UR-900\n'
    printf 'completed_at: 2026-07-30T10:00:00Z\n'
    [ -n "$commit_field" ] && printf 'commit: %s\n' "$commit_field"
    printf -- '---\n\n'
    printf '# Fixture request\n\n## What\n\nBody prose that quotes the schema:\n\n'
    printf '```yaml\ncommit: deadbee   # body prose, MUST NOT be rewritten\n```\n\n'
    printf '## Detailed Requirements\n\n'
    local padding_index=0
    while [ "$padding_index" -lt 120 ]; do
      printf -- '- Requirement line %s: prose padding so the fixture is a realistic size.\n' "$padding_index"
      padding_index=$((padding_index + 1))
    done
  } > "$fixture_path"
}

commit_fixture() {
  git -C "$fixture_root" add -A
  # Committing unconditionally would print "nothing to commit" and exit non-zero whenever a
  # previous probe's seed commit already swept the edit in — noise that hides a real failure.
  if ! git -C "$fixture_root" diff --cached --quiet; then
    git -C "$fixture_root" commit -q -m "$1"
  fi
}

body_commit_line_count() {
  grep -c '^commit: deadbee' "$fixture_root/$1"
}

# --- Probe 1: replace an existing commit: line -----------------------------------------
# The upstream acceptance criterion: a correct write-back is a 1-insertion/1-deletion diff.
write_request_fixture "do-work/archive/UR-900/REQ-1282-replace.md" "0000000"
commit_fixture "seed replace fixture"
implementation_hash="$(git -C "$fixture_root" rev-parse --short HEAD)"
run_record_script "do-work/archive/UR-900/REQ-1282-replace.md" "$implementation_hash"
assert_status 0 "replace: exits 0"
assert_output_matches '^OK:' "replace: reports OK"
numstat_line="$(git -C "$fixture_root" diff --numstat HEAD -- do-work/archive/UR-900/REQ-1282-replace.md | awk '{print $1, $2}')"
assert_equals "1 1" "$numstat_line" "replace: produces a 1-insertion/1-deletion diff"
assert_equals "commit: $implementation_hash" \
  "$(grep -m1 '^commit:' "$fixture_root/do-work/archive/UR-900/REQ-1282-replace.md")" \
  "replace: frontmatter records the hash"
assert_equals "1" "$(body_commit_line_count do-work/archive/UR-900/REQ-1282-replace.md)" \
  "replace: the body's quoted commit: line is untouched"

# --- Probe 2: insert a missing commit: line --------------------------------------------
write_request_fixture "do-work/archive/UR-900/REQ-1282-insert.md" ""
commit_fixture "seed insert fixture"
run_record_script "do-work/archive/UR-900/REQ-1282-insert.md" "$implementation_hash"
assert_status 0 "insert: exits 0"
numstat_line="$(git -C "$fixture_root" diff --numstat HEAD -- do-work/archive/UR-900/REQ-1282-insert.md | awk '{print $1, $2}')"
assert_equals "1 0" "$numstat_line" "insert: produces a 1-insertion/0-deletion diff"
assert_equals "commit: $implementation_hash" \
  "$(sed -n '7p' "$fixture_root/do-work/archive/UR-900/REQ-1282-insert.md")" \
  "insert: the new field lands immediately after completed_at:"

# --- Probe 3: idempotency, and the stranded-edit exception ------------------------------
run_record_script "do-work/archive/UR-900/REQ-1282-insert.md" "$implementation_hash"
assert_status 0 "rerun: exits 0"
assert_output_matches '^OK:.*not in HEAD' "rerun: an uncommitted edit is reported as stranded, not NOOP"
commit_fixture "record the insert"
run_record_script "do-work/archive/UR-900/REQ-1282-insert.md" "$implementation_hash"
assert_status 0 "noop: exits 0"
assert_output_matches '^NOOP:' "noop: a committed, already-recorded hash is a no-op"

# The stranded-edit exception is metadata-specific, not a blanket permission for any dirty
# archived request. Normalize only the first frontmatter block's top-level commit field; body
# prose and fenced examples remain evidence and must reject the staging instruction.
idempotency_body_path="do-work/archive/UR-900/REQ-1282-idempotency-body.md"
write_request_fixture "$idempotency_body_path" "$implementation_hash"
commit_fixture "seed idempotency body fixture"
sed 's/Body prose that quotes the schema:/Body prose changed after metadata write:/' \
  "$fixture_root/$idempotency_body_path" > "$fixture_root/idempotency-body.tmp"
mv "$fixture_root/idempotency-body.tmp" "$fixture_root/$idempotency_body_path"
run_record_script "$idempotency_body_path" "$implementation_hash"
assert_status 1 "idempotency body change: exits 1"
assert_output_matches 'content beyond.*frontmatter.*commit' \
  "idempotency body change: explains that non-metadata content differs"
assert_output_not_matches 'Stage and commit it now' \
  "idempotency body change: does not print the staging instruction"

idempotency_fence_path="do-work/archive/UR-900/REQ-1282-idempotency-fence.md"
write_request_fixture "$idempotency_fence_path" "$implementation_hash"
commit_fixture "seed idempotency fenced-example fixture"
sed 's/^commit: deadbee/commit: feed123/' \
  "$fixture_root/$idempotency_fence_path" > "$fixture_root/idempotency-fence.tmp"
mv "$fixture_root/idempotency-fence.tmp" "$fixture_root/$idempotency_fence_path"
run_record_script "$idempotency_fence_path" "$implementation_hash"
assert_status 1 "idempotency fenced commit change: exits 1"
assert_output_matches 'content beyond.*frontmatter.*commit' \
  "idempotency fenced commit change: keeps body commit text significant"
assert_output_not_matches 'Stage and commit it now' \
  "idempotency fenced commit change: does not print the staging instruction"

# --- Probe 4: the incident — an already-blanked REQ ------------------------------------
# This is the state the six destroyed files were left in. Recording a hash here would commit
# the loss, so it must be refused before anything is written.
write_request_fixture "do-work/archive/UR-900/REQ-1282-blank.md" "0000000"
commit_fixture "seed blank fixture"
: > "$fixture_root/do-work/archive/UR-900/REQ-1282-blank.md"
run_record_script "do-work/archive/UR-900/REQ-1282-blank.md" "$implementation_hash"
assert_status 2 "blanked: exits 2"
assert_output_matches '0 bytes' "blanked: the message says the file is empty"
assert_output_matches 'git show' "blanked: the message names a recovery command"

# --- Probe 5: the incident's partial form — a truncated REQ -----------------------------
write_request_fixture "do-work/archive/UR-900/REQ-1282-truncated.md" "0000000"
commit_fixture "seed truncated fixture"
head -c 400 "$fixture_root/do-work/archive/UR-900/REQ-1282-truncated.md" > "$fixture_root/truncated.tmp"
mv "$fixture_root/truncated.tmp" "$fixture_root/do-work/archive/UR-900/REQ-1282-truncated.md"
truncated_bytes_before="$(wc -c < "$fixture_root/do-work/archive/UR-900/REQ-1282-truncated.md" | tr -d '[:space:]')"
run_record_script "do-work/archive/UR-900/REQ-1282-truncated.md" "$implementation_hash"
assert_status 1 "truncated: exits 1"
assert_output_matches 'bytes in HEAD' "truncated: the message compares on-disk size against HEAD"
assert_equals "$truncated_bytes_before" \
  "$(wc -c < "$fixture_root/do-work/archive/UR-900/REQ-1282-truncated.md" | tr -d '[:space:]')" \
  "truncated: the file is left exactly as found"

# --- Probe 6: ambiguous frontmatter (two commit: keys) ---------------------------------
write_request_fixture "do-work/archive/UR-900/REQ-1282-duplicate.md" "0000000"
sed -e 's/^commit: 0000000$/commit: 0000000\ncommit: 1111111/' \
  "$fixture_root/do-work/archive/UR-900/REQ-1282-duplicate.md" > "$fixture_root/duplicate.tmp"
mv "$fixture_root/duplicate.tmp" "$fixture_root/do-work/archive/UR-900/REQ-1282-duplicate.md"
commit_fixture "seed duplicate fixture"
duplicate_bytes_before="$(wc -c < "$fixture_root/do-work/archive/UR-900/REQ-1282-duplicate.md" | tr -d '[:space:]')"
run_record_script "do-work/archive/UR-900/REQ-1282-duplicate.md" "$implementation_hash"
assert_status 1 "duplicate: exits 1"
assert_equals "$duplicate_bytes_before" \
  "$(wc -c < "$fixture_root/do-work/archive/UR-900/REQ-1282-duplicate.md" | tr -d '[:space:]')" \
  "duplicate: the file is left exactly as found"

# --- Probe 7: unterminated frontmatter --------------------------------------------------
# The awk buffers the frontmatter block; without an END flush it would emit nothing here and
# the script would become the truncation it exists to prevent. Since REQ-276 the refusal
# happens EARLIER and harder: `require_closed_frontmatter` rejects the shape as bad input
# (exit 2) before any reader runs, so the writer's END-flush guard is no longer the thing
# that catches it — same class as the CRLF refusal in Probe 8 directly below, and asserted
# the same way. The property this probe pins is unchanged: nonzero, and the file byte-for-
# byte as found.
printf -- '---\nid: REQ-1283\nstatus: completed\ncompleted_at: 2026-07-30T10:00:00Z\n\n# No closing delimiter\n' \
  > "$fixture_root/do-work/archive/UR-900/REQ-1283-unterminated.md"
commit_fixture "seed unterminated fixture"
unterminated_bytes_before="$(wc -c < "$fixture_root/do-work/archive/UR-900/REQ-1283-unterminated.md" | tr -d '[:space:]')"
run_record_script "do-work/archive/UR-900/REQ-1283-unterminated.md" "$implementation_hash"
assert_status 2 "unterminated: refused as bad input before any reader runs"
assert_equals "$unterminated_bytes_before" \
  "$(wc -c < "$fixture_root/do-work/archive/UR-900/REQ-1283-unterminated.md" | tr -d '[:space:]')" \
  "unterminated: the file is left exactly as found"

# --- Probe 8: CRLF line endings ---------------------------------------------------------
# `$0 == "---"` cannot see a CRLF delimiter, so the edit would silently land nowhere.
write_request_fixture "do-work/archive/UR-900/REQ-1284-crlf.md" "0000000"
awk '{printf "%s\r\n", $0}' "$fixture_root/do-work/archive/UR-900/REQ-1284-crlf.md" > "$fixture_root/crlf.tmp"
mv "$fixture_root/crlf.tmp" "$fixture_root/do-work/archive/UR-900/REQ-1284-crlf.md"
commit_fixture "seed crlf fixture"
run_record_script "do-work/archive/UR-900/REQ-1284-crlf.md" "$implementation_hash"
assert_status 2 "crlf: exits 2"
assert_output_matches 'CRLF' "crlf: the message names the line endings"

# --- Probe 9: argument hygiene ----------------------------------------------------------
run_record_script "do-work/archive/UR-900/REQ-1282-replace.md" "<hash>"
assert_status 2 "placeholder hash: exits 2"
assert_output_matches 'placeholder' "placeholder hash: the message calls out the literal placeholder"

run_record_script "do-work/archive/UR-900/REQ-1282-replace.md"
assert_status 2 "one argument: exits 2"

run_record_script "do-work/archive/UR-900/REQ-1282-replace.md" "$implementation_hash" "extra"
assert_status 2 "three arguments: exits 2"

run_record_script "do-work/archive/UR-900/REQ-1282-missing.md" "$implementation_hash"
assert_status 2 "missing file: exits 2"

# --- Probe 10: a hash this repository cannot resolve -------------------------------------
run_record_script "do-work/archive/UR-900/REQ-1282-replace.md" "abcdef1234567"
assert_status 1 "unresolvable hash: exits 1"
assert_output_matches 'does not resolve' "unresolvable hash: the message says so"

# --- Probe 11: a file with no trailing newline -------------------------------------------
# awk adds one; without accounting for that byte the size arithmetic misfires on a legitimate file.
write_request_fixture "do-work/archive/UR-900/REQ-1285-nonewline.md" "0000000"
printf '%s' "$(cat "$fixture_root/do-work/archive/UR-900/REQ-1285-nonewline.md")" \
  > "$fixture_root/nonewline.tmp"
mv "$fixture_root/nonewline.tmp" "$fixture_root/do-work/archive/UR-900/REQ-1285-nonewline.md"
commit_fixture "seed no-trailing-newline fixture"
run_record_script "do-work/archive/UR-900/REQ-1285-nonewline.md" "$implementation_hash"
assert_status 0 "no trailing newline: exits 0"
assert_output_matches 'trailing newline' "no trailing newline: the added byte is reported"

# --- Probe 12: --verify against a genuine one-line metadata commit -------------------------
# A content-mutating pre-commit hook (formatter, lint --fix) is a sufficient cause of the
# original incident and is invisible to every pre-commit guard. Only reading back what landed
# sees it. This probe establishes the clean baseline: write the hash, commit exactly that one
# path (what the script's own printed instructions do), then verify.
verify_path="do-work/archive/UR-900/REQ-1287-verify.md"
write_request_fixture "$verify_path" "0000000"
commit_fixture "seed verify fixture"
run_record_script "$verify_path" "$implementation_hash"
assert_status 0 "verify setup: the write-back exits 0"
git -C "$fixture_root" commit -q -m "[REQ-1287] record commit hash $implementation_hash" -- "$verify_path"
run_record_script --verify "$verify_path" "$implementation_hash"
assert_status 0 "verify clean: exits 0 on a genuine one-line metadata commit"
assert_output_matches 'exactly the one line' \
  "verify clean: reports that HEAD's patch was the single commit: line"

# --- Probe 13: --verify catches a re-staging hook that rewrote the body ---------------------
# THE case a blob read-back cannot see. `lint-staged` and friends rewrite the worktree file and
# re-stage it, so the committed blob and the worktree agree — same bytes, same single
# `commit:` line — while body content is silently gone. Only inspecting the committed PATCH
# fails this, which is why --verify asserts on the patch and not just on sizes.
sed 's/^- Requirement line 42:.*/- Requirement line 42: REWRITTEN BY A PRE-COMMIT HOOK./' \
  "$fixture_root/$verify_path" > "$fixture_root/hooked.tmp"
mv "$fixture_root/hooked.tmp" "$fixture_root/$verify_path"
git -C "$fixture_root" commit -q -m "simulate a hook that rewrote the body and re-staged it" -- "$verify_path"
# Precondition: the read-back layer genuinely cannot tell — worktree and HEAD are identical.
assert_equals "$(git -C "$fixture_root" cat-file -s "HEAD:$verify_path")" \
  "$(wc -c < "$fixture_root/$verify_path" | tr -d '[:space:]')" \
  "verify hook-rewrite: the committed blob and the worktree agree, so only the patch check can fire"
run_record_script --verify "$verify_path" "$implementation_hash"
assert_status 1 "verify hook-rewrite: exits 1"
assert_output_matches "changed more than the 'commit:' line" \
  "verify hook-rewrite: the message names what the metadata commit actually changed"

# --- Probe 14: --verify run later than the procedure calls for ------------------------------
# HEAD no longer touches the file, so there is no metadata-commit patch to inspect. That is
# neither evidence of damage nor evidence of correctness: it must degrade with a stated reason
# rather than false-fail or pass silently.
git -C "$fixture_root" commit -q --allow-empty -m "an unrelated later commit"
run_record_script --verify "do-work/archive/UR-900/REQ-1282-insert.md" "$implementation_hash"
assert_status 0 "verify late: exits 0"
assert_output_matches 'committed-patch check was skipped' "verify late: says the patch check did not run"
assert_output_matches 'read-back only' "verify late: labels the weaker guarantee it did make"

# --- Probe 15: --verify catches an outright truncation --------------------------------------
: > "$fixture_root/do-work/archive/UR-900/REQ-1282-replace.md"
git -C "$fixture_root" add -A
git -C "$fixture_root" commit -q -m "simulate a hook that truncated the file"
run_record_script --verify "do-work/archive/UR-900/REQ-1282-replace.md" "$implementation_hash"
assert_status 1 "verify tampered: exits 1"

# --- Probe 16: --verify catches a hook that deleted a BODY commit: line ---------------------
# The narrow gap left by matching the removed line on SHAPE (`commit:*`) instead of against the
# parent's real frontmatter. The fixture body quotes the schema — as real archived REQs do — so
# a hook that drops that quoted line while the write-back INSERTS the frontmatter one nets
# +1/-1: one added, one removed, both starting `commit:`. That read as a legitimate replace and
# passed with a message claiming the patch was a single line, while a body line was gone.
verify_body_path="do-work/archive/UR-900/REQ-1288-body-commit.md"
write_request_fixture "$verify_body_path" ""
commit_fixture "seed body-commit fixture"
run_record_script "$verify_body_path" "$implementation_hash"
assert_status 0 "verify body-commit setup: the write-back exits 0"
grep -v '^commit: deadbee' "$fixture_root/$verify_body_path" > "$fixture_root/body-hooked.tmp"
mv "$fixture_root/body-hooked.tmp" "$fixture_root/$verify_body_path"
git -C "$fixture_root" commit -q -m "[REQ-1282] record commit hash $implementation_hash" -- "$verify_body_path"
assert_equals "0" "$(body_commit_line_count "$verify_body_path")" \
  "verify body-commit: the simulated hook really did remove the body's commit: line"
# Precondition: as in Probe 13, the read-back layer cannot tell — the hook's edit was committed.
assert_equals "$(git -C "$fixture_root" cat-file -s "HEAD:$verify_body_path")" \
  "$(wc -c < "$fixture_root/$verify_body_path" | tr -d '[:space:]')" \
  "verify body-commit: the committed blob and the worktree agree, so only the patch check can fire"
run_record_script --verify "$verify_body_path" "$implementation_hash"
assert_status 1 "verify body-commit: exits 1"
assert_output_matches 'HEAD\^ has no frontmatter commit: field' \
  "verify body-commit: the message states what the removal side was measured against"
assert_output_matches 'came from the BODY' \
  "verify body-commit: the message names the content actually lost"

# --- Probe 13: outside a git repository ---------------------------------------------------
# Content guards must still run; only the git-dependent ones degrade, and with a stated reason.
non_git_root="$(mktemp -d)"
mkdir -p "$non_git_root/do-work/archive/UR-900"
fixture_root_backup="$fixture_root"
fixture_root="$non_git_root"
write_request_fixture "do-work/archive/UR-900/REQ-1286-nogit.md" "0000000"
run_record_script "do-work/archive/UR-900/REQ-1286-nogit.md" "1234567"
assert_status 0 "non-git: exits 0"
assert_output_matches 'not a git repository' "non-git: the skipped git checks are stated"
assert_equals "commit: 1234567" \
  "$(grep -m1 '^commit:' "$non_git_root/do-work/archive/UR-900/REQ-1286-nogit.md")" \
  "non-git: the edit still happens"
fixture_root="$fixture_root_backup"
rm -rf "$non_git_root"

# =========================================================================================
# tools/checks/blanked-req-scan.sh — detection of REQ/UR files whose content was destroyed.
# Shares this fixture harness deliberately: the scanner's whole job is finding the aftermath
# of the write-back failure the probes above prevent, so the two are tested against the same
# repo shapes.
# =========================================================================================
scan_script="$repo_root/tools/checks/blanked-req-scan.sh"
if [ ! -x "$scan_script" ]; then
  printf 'FAIL: tools/checks/blanked-req-scan.sh must exist and be executable.\n' >&2
  exit 1
fi

scan_root="$(mktemp -d)"
cleanup_scan() { rm -rf "$scan_root"; }
trap 'cleanup_fixture; cleanup_scan' EXIT

# Takes no arguments on purpose: every scan probe exercises the bare detector. The restore
# probes below use run_restore_script, which does forward its arguments.
run_scan_script() {
  probe_output="$(cd "$scan_root" && "$scan_script" 2>&1)"
  probe_status=$?
}

git -c init.defaultBranch=main init -q "$scan_root"
git -C "$scan_root" config user.name "Fixture Runner"
git -C "$scan_root" config user.email "fixture@example.invalid"

# Reproduce the incident end to end: a complete archived REQ, then a metadata commit whose
# message claims success while replacing the whole file with nothing.
mkdir -p "$scan_root/do-work/archive/UR-900"
scan_target="$scan_root/do-work/archive/UR-900/REQ-1282-incident.md"
{
  printf -- '---\nid: REQ-1282\ntitle: Incident fixture\nstatus: completed\ncompleted_at: 2026-07-30T10:00:00Z\ncommit: 0000000\n---\n\n# Incident fixture\n\n'
  padding_index=0
  while [ "$padding_index" -lt 150 ]; do
    printf -- '- decision trail line %s that must be recoverable from git history.\n' "$padding_index"
    padding_index=$((padding_index + 1))
  done
} > "$scan_target"
# A healthy neighbour: the scan must not report it.
printf -- '---\nid: REQ-1283\ntitle: Healthy\nstatus: completed\ncompleted_at: 2026-07-30T10:00:00Z\ncommit: 0000000\n---\n\n# Healthy\n\nIntact body.\n' \
  > "$scan_root/do-work/archive/UR-900/REQ-1283-healthy.md"
mkdir -p "$scan_root/do-work/archive/UR-900/assets"
printf -- '# Screenshot descriptions\n\nThis is intact prose, not a REQ record.\n' \
  > "$scan_root/do-work/archive/UR-900/assets/REQ-1284-screenshot-descriptions.md"
git -C "$scan_root" add -A
git -C "$scan_root" commit -q -m "[REQ-1282] Incident fixture (Route C)"
recovery_source_sha="$(git -C "$scan_root" rev-parse HEAD)"
recorded_hash="$(git -C "$scan_root" rev-parse --short HEAD)"
intact_bytes="$(wc -c < "$scan_target" | tr -d '[:space:]')"

: > "$scan_target"
git -C "$scan_root" add -A
git -C "$scan_root" commit -q -m "[REQ-1282] record commit hash $recorded_hash"

run_scan_script
assert_status 1 "scan: exits 1 when a blanked file is found"
assert_output_matches 'REQ-1282-incident\.md' "scan: names the blanked file"
assert_output_matches '0 bytes' "scan: reports the file as empty"
assert_output_matches "$recorded_hash" "scan: recovers the hash from the blanking commit's message"
# The recoverable size is what tells an operator whether the loss is worth acting on, and it
# must be the PRE-blanking size read out of the recovery commit — not the 0 bytes on disk.
assert_output_matches "Recoverable: $intact_bytes bytes" \
  "scan: reports the recoverable byte count from the pre-blanking commit"
if printf '%s' "$probe_output" | grep -q 'REQ-1283-healthy'; then
  printf 'FAIL: scan: reported the intact neighbour REQ-1283 as blanked.\n' >&2
  fail_count=$((fail_count + 1))
fi
if printf '%s' "$probe_output" | grep -q 'REQ-1284-screenshot-descriptions'; then
  printf 'FAIL: scan: reported intact prose under archive/assets as blanked.\n' >&2
  fail_count=$((fail_count + 1))
fi

# The machine-readable line is what REQ-064's --restore consumes; assert its exact shape so
# the two sides cannot drift.
scan_record="$(printf '%s\n' "$probe_output" | grep '^BLANKED' | head -1)"
assert_equals "BLANKED	do-work/archive/UR-900/REQ-1282-incident.md	$recovery_source_sha	$recorded_hash" \
  "$scan_record" "scan: emits a machine-readable BLANKED record with path, recovery source and hash"

# A REQ blanked with no non-empty ancestor cannot be recovered from history — that must be
# reported as such, never skipped into silence.
printf '' > "$scan_root/do-work/archive/UR-900/REQ-1284-neverhad.md"
git -C "$scan_root" add -A
git -C "$scan_root" commit -q -m "add an always-empty REQ"
run_scan_script
assert_output_matches 'REQ-1284-neverhad\.md' "scan: names a file with no recovery source"
assert_output_matches 'No recoverable content in git history' "scan: says the file has no recoverable content in history"

# A clean tree must exit 0 and say so, so forensics can report "no findings" confidently.
clean_root="$(mktemp -d)"
mkdir -p "$clean_root/do-work/archive/UR-900"
printf -- '---\nid: REQ-1290\nstatus: completed\n---\n\nIntact.\n' \
  > "$clean_root/do-work/archive/UR-900/REQ-1290-fine.md"
probe_output="$(cd "$clean_root" && "$scan_script" 2>&1)"; probe_status=$?
assert_status 0 "scan clean: exits 0 when nothing is blanked"
assert_output_matches 'No blanked' "scan clean: says nothing was found"
rm -rf "$clean_root"

# --- --restore: the repair path (actions/cleanup.md Pass 6) -------------------------------
# This is the recovery that had to be performed by hand for six files in the incident. The
# assertion that matters is byte-identity with the pre-blanking blob apart from the one
# commit: line — a restore that "mostly" recovers the content is not a recovery.
restore_root="$(mktemp -d)"
cleanup_restore() { rm -rf "$restore_root"; }
trap 'cleanup_fixture; cleanup_scan; cleanup_restore' EXIT

git -c init.defaultBranch=main init -q "$restore_root"
git -C "$restore_root" config user.name "Fixture Runner"
git -C "$restore_root" config user.email "fixture@example.invalid"
mkdir -p "$restore_root/do-work/archive/UR-900"
restore_target="$restore_root/do-work/archive/UR-900/REQ-1282-restore.md"
{
  printf -- '---\nid: REQ-1282\ntitle: Restore fixture\nstatus: completed\ncompleted_at: 2026-07-30T10:00:00Z\ncommit: 0000000\n---\n\n# Restore fixture\n\n'
  padding_index=0
  while [ "$padding_index" -lt 150 ]; do
    printf -- '- irreplaceable decision trail line %s.\n' "$padding_index"
    padding_index=$((padding_index + 1))
  done
} > "$restore_target"
# A healthy file the restore must not touch.
printf -- '---\nid: REQ-1288\nstatus: completed\n---\n\nHealthy body.\n' \
  > "$restore_root/do-work/archive/UR-900/REQ-1288-healthy.md"
git -C "$restore_root" add -A
git -C "$restore_root" commit -q -m "[REQ-1282] Restore fixture (Route C)"
expected_restored_content="$(cat "$restore_target")"
healthy_before="$(cat "$restore_root/do-work/archive/UR-900/REQ-1288-healthy.md")"
# The recorded hash must be a REAL commit, as it is in the field: it is the implementation
# commit Step 9 just made. record-commit-hash.sh refuses a hash the repo cannot resolve, so a
# fictional one would (correctly) be rejected and the restore would leave commit: untouched.
restore_recorded_hash="$(git -C "$restore_root" rev-parse --short HEAD)"

: > "$restore_target"
git -C "$restore_root" add -A
git -C "$restore_root" commit -q -m "[REQ-1282] record commit hash $restore_recorded_hash"

run_restore_script() {
  probe_output="$(cd "$restore_root" && "$scan_script" "$@" 2>&1)"
  probe_status=$?
}

# --dry-run must report the plan and write nothing.
run_restore_script --restore --dry-run
assert_output_matches 'REQ-1282-restore\.md' "restore dry-run: names the file it would restore"
assert_equals "0" "$(wc -c < "$restore_target" | tr -d '[:space:]')" \
  "restore dry-run: writes nothing"

run_restore_script --restore
assert_status 0 "restore: exits 0 after a successful repair"
# Byte-identical to the pre-blanking content except the one commit: line, which must now
# carry the hash recorded in the blanking commit's own message.
expected_after_restore="$(printf '%s' "$expected_restored_content" | sed "s/^commit: 0000000$/commit: $restore_recorded_hash/")"
assert_equals "$expected_after_restore" "$(cat "$restore_target")" \
  "restore: content matches the pre-blanking blob byte for byte, with the recorded hash applied"
assert_equals "$healthy_before" "$(cat "$restore_root/do-work/archive/UR-900/REQ-1288-healthy.md")" \
  "restore: the healthy neighbour is untouched"

# Re-running finds nothing left to repair.
run_restore_script --restore
assert_status 0 "restore rerun: exits 0"
assert_output_matches 'No blanked' "restore rerun: nothing left to restore"

# --- --restore: malformed non-empty blobs stay inside the damage chain -------------------
# A failed metadata write can leave bytes behind without leaving parseable frontmatter. The
# newest non-empty blob is therefore not necessarily recoverable content: selection must walk
# through malformed damage to the earlier valid request, while the same malformed commit's
# subject remains eligible to supply the recorded implementation hash.
malformed_root="$(mktemp -d)"
cleanup_malformed() { rm -rf "$malformed_root"; }
trap 'cleanup_fixture; cleanup_scan; cleanup_restore; cleanup_malformed' EXIT

git -c init.defaultBranch=main init -q "$malformed_root"
git -C "$malformed_root" config user.name "Fixture Runner"
git -C "$malformed_root" config user.email "fixture@example.invalid"
mkdir -p "$malformed_root/do-work/archive/UR-900"
malformed_target="$malformed_root/do-work/archive/UR-900/REQ-1291-malformed-chain.md"
printf -- '---\nid: REQ-1291\ntitle: Malformed chain fixture\nstatus: completed\ncompleted_at: 2026-07-30T10:00:00Z\ncommit: 0000000\n---\n\n# Valid recovery body\n\nIrreplaceable context.\n' \
  > "$malformed_target"
git -C "$malformed_root" add -A
git -C "$malformed_root" commit -q -m "[REQ-1291] Valid recovery fixture (Route B)"
malformed_valid_sha="$(git -C "$malformed_root" rev-parse HEAD)"
malformed_recorded_hash="$(git -C "$malformed_root" rev-parse --short HEAD)"
malformed_expected_content="$(cat "$malformed_target")"

printf 'not frontmatter\npartial decision trail\n' > "$malformed_target"
git -C "$malformed_root" add -A
git -C "$malformed_root" commit -q -m "[REQ-1291] record commit hash $malformed_recorded_hash"
malformed_damage_sha="$(git -C "$malformed_root" rev-parse HEAD)"

probe_output="$(cd "$malformed_root" && "$scan_script" 2>&1)"; probe_status=$?
assert_status 1 "malformed chain scan: exits 1 before repair"
malformed_scan_record="$(printf '%s\n' "$probe_output" | grep '^BLANKED' | head -1)"
malformed_expected_record="$(printf 'BLANKED\tdo-work/archive/UR-900/REQ-1291-malformed-chain.md\t%s\t%s' \
  "$malformed_valid_sha" "$malformed_recorded_hash")"
assert_equals "$malformed_expected_record" \
  "$malformed_scan_record" "malformed chain scan: reports the earlier valid source and the damage commit's recorded hash"
if printf '%s\n' "$malformed_scan_record" | grep -qF "$malformed_damage_sha"; then
  printf 'FAIL: malformed chain scan: selected the malformed damage blob as its recovery source.\n' >&2
  fail_count=$((fail_count + 1))
fi

probe_output="$(cd "$malformed_root" && "$scan_script" --restore 2>&1)"; probe_status=$?
assert_status 0 "malformed chain restore: exits 0 only after valid content is repaired"
malformed_expected_after_restore="$(printf '%s' "$malformed_expected_content" | sed "s/^commit: 0000000$/commit: $malformed_recorded_hash/")"
assert_equals "$malformed_expected_after_restore" "$(cat "$malformed_target")" \
  "malformed chain restore: restores the earlier valid content and recorded hash"

probe_output="$(cd "$malformed_root" && "$scan_script" 2>&1)"; probe_status=$?
assert_status 0 "malformed chain rescan: exits 0 after repair"
assert_output_matches 'No blanked' "malformed chain rescan: reports a clean tree"

# --- --restore: a PARTIAL repair is a finding, not a success -------------------------------
# Content back but the recorded commit: hash rejected leaves a file that looks healthy and is
# committable with provenance pointing at nothing. Reporting that as repaired would put a
# fresh falsehood into the very trail this pass exists to save, so it must exit non-zero.
# The blanking commit here names a hash this repo cannot resolve — the shape the guarded
# write-back refuses — which is how a real recorded hash goes bad (history rewritten,
# recovered from a different clone).
partial_root="$(mktemp -d)"
cleanup_partial() { rm -rf "$partial_root"; }
trap 'cleanup_fixture; cleanup_scan; cleanup_restore; cleanup_malformed; cleanup_partial' EXIT

git -c init.defaultBranch=main init -q "$partial_root"
git -C "$partial_root" config user.name "Fixture Runner"
git -C "$partial_root" config user.email "fixture@example.invalid"
mkdir -p "$partial_root/do-work/archive/UR-900"
partial_target="$partial_root/do-work/archive/UR-900/REQ-1289-partial.md"
printf -- '---\nid: REQ-1289\ntitle: Partial fixture\nstatus: completed\ncompleted_at: 2026-07-30T10:00:00Z\ncommit: 0000000\n---\n\n# Partial fixture\n\nDecision trail that must come back even when the hash cannot.\n' \
  > "$partial_target"
git -C "$partial_root" add -A
git -C "$partial_root" commit -q -m "[REQ-1289] Partial fixture (Route C)"
partial_intact_bytes="$(wc -c < "$partial_target" | tr -d '[:space:]')"

: > "$partial_target"
git -C "$partial_root" add -A
git -C "$partial_root" commit -q -m "[REQ-1289] record commit hash abcdef1234567"

probe_output="$(cd "$partial_root" && "$scan_script" --restore 2>&1)"; probe_status=$?
assert_status 1 "restore partial: exits 1 when the recorded hash cannot be re-applied"
assert_output_matches 'content restored, but re-applying' \
  "restore partial: names the file as content-restored-but-hash-failed"
assert_output_matches 'does not resolve' \
  "restore partial: passes the write-back's own diagnosis through instead of swallowing it"
assert_output_matches 'Fully repaired 0 of 1' \
  "restore partial: does not count a partial repair as a full one"
# The content half must still have landed — a partial repair is reported, never rolled back.
assert_equals "$partial_intact_bytes" "$(wc -c < "$partial_target" | tr -d '[:space:]')" \
  "restore partial: the recovered content is written even though the hash failed"

# --- Probe 20: a fence that never closes is refused on the write path -----------------------
# The rewrite awk buffers the frontmatter and emits it only after the closing `---`, so the
# WRITE path could never corrupt such a file. The three READERS scan to end of file, so on
# this shape they reach the body's `commit: deadbee` line and take it for the schema field —
# and the readers run first, at startup, deciding the id, the duplicate count and the
# pre-edit status before the writer is ever consulted.
#
# What this fixture demonstrates, precisely: it has a frontmatter `commit:` AND the body's
# fenced `commit: deadbee`, so before the fix the duplicate-count reader saw TWO and refused
# with "frontmatter has 2 'commit:' lines" — the readers counting a body line as schema is
# the misread itself, wearing a different error message. The assertions below discriminate
# the two eras by exit code and reason, not merely by failing: exit 2 and "never closed"
# versus the old exit 1 and "ambiguous". Probe 21 is the one that exercises the silent
# misattribution end to end, on the path where it can pass something wrong (REQ-276).
unclosed_fence_path="do-work/archive/UR-900/REQ-1288-unclosed-fence.md"
write_request_fixture "$unclosed_fence_path" "0000000"
# Remove ONLY the closing fence, leaving the opening one and the whole body intact.
awk 'NR == 1 { print; next } $0 == "---" && !dropped { dropped = 1; next } { print }' \
  "$fixture_root/$unclosed_fence_path" > "$fixture_root/unclosed-fence.tmp"
mv "$fixture_root/unclosed-fence.tmp" "$fixture_root/$unclosed_fence_path"
commit_fixture "seed unclosed-fence fixture"
# Precondition: the shape really is what the probe claims, and the body bait is really there.
assert_equals "1" "$(grep -c '^---$' "$fixture_root/$unclosed_fence_path")" \
  "unclosed fence: the fixture opens a fence and never closes it"
assert_equals "1" "$(body_commit_line_count "$unclosed_fence_path")" \
  "unclosed fence: the body still carries the column-0 commit: line a scanning reader would find"
unclosed_fence_bytes="$(wc -c < "$fixture_root/$unclosed_fence_path" | tr -d '[:space:]')"
run_record_script "$unclosed_fence_path" "$implementation_hash"
assert_status 2 "unclosed fence: refused as bad input rather than acted on"
assert_output_matches 'never closed' \
  "unclosed fence: names the fence as the reason instead of failing somewhere downstream"
assert_equals "$unclosed_fence_bytes" "$(wc -c < "$fixture_root/$unclosed_fence_path" | tr -d '[:space:]')" \
  "unclosed fence: the refused file is byte-identical"

# --- Probe 21: --verify refuses a PARENT blob whose fence never closes ----------------------
# The parent blob is the one reader input that is not $request_file, so the startup guard
# above cannot cover it. --verify reads HEAD^'s copy to decide what the metadata commit was
# allowed to remove; if that read scans past a missing fence into the body, the expected
# removal is a guess about prose and a corrupted file can verify clean. This is the last
# check every REQ in the pipeline passes through (REQ-276).
verify_fence_path="do-work/archive/UR-900/REQ-1289-verify-fence.md"
write_request_fixture "$verify_fence_path" ""
awk 'NR == 1 { print; next } $0 == "---" && !dropped { dropped = 1; next } { print }' \
  "$fixture_root/$verify_fence_path" > "$fixture_root/verify-fence.tmp"
mv "$fixture_root/verify-fence.tmp" "$fixture_root/$verify_fence_path"
commit_fixture "seed verify parent-fence fixture"
# The parent commit now holds the unclosed-fence content. Give HEAD a closed-fence version so
# the startup guard passes and --verify gets as far as reading the parent.
write_request_fixture "$verify_fence_path" "$implementation_hash"
git -C "$fixture_root" commit -q -m "[REQ-1289] record commit hash $implementation_hash" -- "$verify_fence_path"
assert_equals "1" \
  "$(git -C "$fixture_root" show "HEAD^:$verify_fence_path" | grep -c '^---$')" \
  "verify parent fence: the parent blob really is the unclosed-fence shape"
run_record_script --verify "$verify_fence_path" "$implementation_hash"
assert_status 1 "verify parent fence: refuses rather than verifying against an unreadable parent"
assert_output_matches 'never closed' \
  "verify parent fence: names the parent's fence as the reason"

if [ "$fail_count" -gt 0 ]; then
  printf '%s guard probe(s) failed.\n' "$fail_count" >&2
  exit 1
fi

printf 'record-commit-hash and blanked-req-scan guard probes passed.\n'
