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
# the script would become the truncation it exists to prevent.
printf -- '---\nid: REQ-1283\nstatus: completed\ncompleted_at: 2026-07-30T10:00:00Z\n\n# No closing delimiter\n' \
  > "$fixture_root/do-work/archive/UR-900/REQ-1283-unterminated.md"
commit_fixture "seed unterminated fixture"
unterminated_bytes_before="$(wc -c < "$fixture_root/do-work/archive/UR-900/REQ-1283-unterminated.md" | tr -d '[:space:]')"
run_record_script "do-work/archive/UR-900/REQ-1283-unterminated.md" "$implementation_hash"
assert_status 1 "unterminated: exits 1"
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

# --- Probe 12: --verify catches content that changed after the guards passed --------------
# A content-mutating pre-commit hook (formatter, lint --fix) is a sufficient cause of the
# original incident and is invisible to every pre-commit guard. Only a read-back sees it.
commit_fixture "record the replace fixture"
run_record_script --verify "do-work/archive/UR-900/REQ-1282-replace.md" "$implementation_hash"
assert_status 0 "verify clean: exits 0"

: > "$fixture_root/do-work/archive/UR-900/REQ-1282-replace.md"
git -C "$fixture_root" add -A
git -C "$fixture_root" commit -q -m "simulate a hook that truncated the file"
run_record_script --verify "do-work/archive/UR-900/REQ-1282-replace.md" "$implementation_hash"
assert_status 1 "verify tampered: exits 1"

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

if [ "$fail_count" -gt 0 ]; then
  printf '%s guard probe(s) failed.\n' "$fail_count" >&2
  exit 1
fi

printf 'record-commit-hash and blanked-req-scan guard probes passed.\n'
