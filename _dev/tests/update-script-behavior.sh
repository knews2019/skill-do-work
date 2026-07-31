#!/usr/bin/env bash
# update-script-behavior.sh — behavioral probes for tools/do-work-update.sh.
#
# The updater overwrites a live install in place and keeps NO rollback copy: git is the undo.
# That contract only holds if two things are true at runtime, and neither is greppable from
# prose — the happy path must leave no `.preupdate-*.bak` behind, and a failure inside the
# destructive region must announce the partial install with runnable recovery commands
# instead of exiting quietly (there is nothing to restore from automatically any more).
#
# The real upstream fetch is stubbed by putting a fake `curl` first on PATH, so these probes
# never touch the network and pin an exact upstream tree to assert against.
#
# Invoked by _dev/tests/contract-regressions.sh — there is no auto-discovery of test files,
# so a probe added here only runs because that file calls this one.
#
# Exit 0: every probe passed (or the whole suite was skipped for a missing prerequisite).
# Exit 1: at least one probe failed; each failure prints a FAIL line naming the probe.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
update_script="$repo_root/tools/do-work-update.sh"
fail_count=0

if ! command -v git >/dev/null 2>&1; then
  printf 'SKIP: git unavailable — update-script behavior probes not run.\n'
  exit 0
fi
if [ ! -x "$update_script" ]; then
  printf 'FAIL: tools/do-work-update.sh must exist and be executable.\n' >&2
  exit 1
fi

# Same reasoning as the record-commit-hash fixture: a developer's global git config
# (core.hooksPath, commit.gpgsign, init.templateDir) must not change what these measure.
export GIT_CONFIG_GLOBAL=/dev/null
export GIT_CONFIG_SYSTEM=/dev/null
export GIT_TERMINAL_PROMPT=0

fixture_root="$(mktemp -d)"
cleanup_fixture() {
  # The failure probe deliberately makes a directory unwritable; restore it or rm -rf fails.
  chmod -R u+rwX "$fixture_root" 2>/dev/null || true
  rm -rf "$fixture_root"
}
trap cleanup_fixture EXIT

probe_output=""
probe_status=0

assert_status() {
  local expected_status="$1" probe_name="$2"
  if [ "$probe_status" -ne "$expected_status" ]; then
    printf 'FAIL: %s — expected exit %s, got %s. Output:\n%s\n' \
      "$probe_name" "$expected_status" "$probe_status" "$probe_output" >&2
    fail_count=$((fail_count + 1))
  fi
}

assert_status_nonzero() {
  local probe_name="$1"
  if [ "$probe_status" -eq 0 ]; then
    printf 'FAIL: %s — expected a non-zero exit, got 0. Output:\n%s\n' \
      "$probe_name" "$probe_output" >&2
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

assert_output_lacks() {
  local pattern_text="$1" probe_name="$2"
  if printf '%s' "$probe_output" | grep -Eq -- "$pattern_text"; then
    printf 'FAIL: %s — output unexpectedly matched /%s/. Output:\n%s\n' \
      "$probe_name" "$pattern_text" "$probe_output" >&2
    fail_count=$((fail_count + 1))
  fi
}

assert_path_exists() {
  local candidate_path="$1" probe_name="$2"
  if [ ! -e "$candidate_path" ]; then
    printf 'FAIL: %s — expected %s to exist.\n' "$probe_name" "$candidate_path" >&2
    fail_count=$((fail_count + 1))
  fi
}

assert_path_absent() {
  local candidate_path="$1" probe_name="$2"
  if [ -e "$candidate_path" ]; then
    printf 'FAIL: %s — expected %s to be absent.\n' "$probe_name" "$candidate_path" >&2
    fail_count=$((fail_count + 1))
  fi
}

# No rollback copy, anywhere under the project — the whole point of the no-copy contract.
assert_no_backup_copy() {
  local project_path="$1" probe_name="$2"
  local stray_backups
  stray_backups="$(find "$project_path" -name '*.preupdate-*.bak' -maxdepth 4 2>/dev/null)"
  if [ -n "$stray_backups" ]; then
    printf 'FAIL: %s — the updater left a rollback copy:\n%s\n' "$probe_name" "$stray_backups" >&2
    fail_count=$((fail_count + 1))
  fi
}

# --- Fake upstream ----------------------------------------------------------------------
# A minimal but structurally real skill tree at v0.0.2, tarred with a top-level directory
# because the updater extracts with --strip-components=1.
upstream_tree="$fixture_root/upstream-src/do-work-upstream"
mkdir -p "$upstream_tree/actions" "$upstream_tree/prompts" "$upstream_tree/interviews" \
  "$upstream_tree/docs" "$upstream_tree/tools"
printf '# do-work\n\nFake upstream SKILL.md at v0.0.2.\n' > "$upstream_tree/SKILL.md"
printf '# Version Action\n\n**Current version**: 0.0.2\n' > "$upstream_tree/actions/version.md"
printf '# Prompts\n' > "$upstream_tree/prompts/README.md"
printf '# Fresh prompt\n' > "$upstream_tree/prompts/fresh-prompt.md"
printf '# Fresh interview\n' > "$upstream_tree/interviews/fresh-interview.md"
printf '# Guide\n' > "$upstream_tree/docs/guide.md"
cp "$update_script" "$upstream_tree/tools/do-work-update.sh"
chmod +x "$upstream_tree/tools/do-work-update.sh"
upstream_tarball="$fixture_root/upstream.tar.gz"
tar czf "$upstream_tarball" -C "$fixture_root/upstream-src" do-work-upstream

# --- Stub curl --------------------------------------------------------------------------
# The updater calls `curl -fsSL -o <dest> <url>`; the stub honours -o and serves the fixture
# tarball, so the probes are hermetic and the upstream tree is known exactly.
stub_bin="$fixture_root/stub-bin"
mkdir -p "$stub_bin"
cat > "$stub_bin/curl" <<STUB
#!/usr/bin/env bash
destination_path=''
while [ "\$#" -gt 0 ]; do
  case "\$1" in
    -o) destination_path="\$2"; shift 2 ;;
    *) shift ;;
  esac
done
[ -n "\$destination_path" ] || exit 2
cp "$upstream_tarball" "\$destination_path"
STUB
chmod +x "$stub_bin/curl"
export PATH="$stub_bin:$PATH"

# --- Fixture installs ---------------------------------------------------------------------
# A v0.0.1 nested install, plus a runtime queue file that must survive, plus stale globbed
# files that the pre-clean must remove.
build_install() {
  local project_path="$1"
  local install_path="$project_path/.claude/skills/do-work"
  mkdir -p "$install_path/actions" "$install_path/prompts" "$install_path/interviews" \
    "$install_path/docs" "$install_path/tools" "$project_path/do-work/queue"
  printf '# do-work\n\nInstalled SKILL.md at v0.0.1.\n' > "$install_path/SKILL.md"
  printf '# Version Action\n\n**Current version**: 0.0.1\n' > "$install_path/actions/version.md"
  printf '# Prompts\n' > "$install_path/prompts/README.md"
  printf '# Stale prompt\n' > "$install_path/prompts/stale-prompt.md"
  printf '# Stale interview\n' > "$install_path/interviews/stale-interview.md"
  printf '# Guide\n' > "$install_path/docs/guide.md"
  cp "$update_script" "$install_path/tools/do-work-update.sh"
  chmod +x "$install_path/tools/do-work-update.sh"
  printf 'REQ-001 sentinel\n' > "$project_path/do-work/queue/sentinel.txt"
}

run_updater() {
  local project_path="$1" answer_text="$2"
  probe_output="$(printf '%s\n' "$answer_text" \
    | bash "$project_path/.claude/skills/do-work/tools/do-work-update.sh" \
        --project-root "$project_path" 2>&1)"
  probe_status=$?
}

# --- Probe 1: happy path leaves no rollback copy ------------------------------------------
git_project="$fixture_root/git-project"
build_install "$git_project"
git -c init.defaultBranch=main init -q "$git_project"
git -C "$git_project" config user.name "Fixture Runner"
git -C "$git_project" config user.email "fixture@example.invalid"
git -C "$git_project" add -A
git -C "$git_project" commit -qm 'install do-work v0.0.1'

run_updater "$git_project" y
git_install="$git_project/.claude/skills/do-work"
assert_status 0 'happy path: updater exits 0'
assert_output_matches 'Updated to v0\.0\.2' 'happy path: reports the new version'
assert_output_lacks 'Rollback copy' 'happy path: no rollback-copy pointer is printed'
assert_no_backup_copy "$git_project" 'happy path'
assert_output_matches 'no rollback copy is kept' 'happy path: the prompt states the no-copy contract'
if ! grep -q '0\.0\.2' "$git_install/actions/version.md"; then
  printf 'FAIL: happy path — the install was not advanced to v0.0.2.\n' >&2
  fail_count=$((fail_count + 1))
fi
assert_path_absent "$git_install/prompts/stale-prompt.md" 'happy path: stale globbed prompt removed'
assert_path_absent "$git_install/interviews/stale-interview.md" 'happy path: stale interview removed'
assert_path_exists "$git_install/prompts/fresh-prompt.md" 'happy path: upstream prompt extracted'
assert_path_exists "$git_project/do-work/queue/sentinel.txt" 'happy path: runtime queue preserved'

# --- Probe 2: failure in the destructive region reports, never restores --------------------
# An unwritable shipped directory makes the extraction fail after the pre-clean has already
# deleted files — the exact half-old/half-new state the removed rollback copy used to undo.
failing_project="$fixture_root/failing-project"
build_install "$failing_project"
git -c init.defaultBranch=main init -q "$failing_project"
git -C "$failing_project" config user.name "Fixture Runner"
git -C "$failing_project" config user.email "fixture@example.invalid"
git -C "$failing_project" add -A
git -C "$failing_project" commit -qm 'install do-work v0.0.1'
failing_install="$failing_project/.claude/skills/do-work"
chmod 500 "$failing_install/docs"

run_updater "$failing_project" y
chmod 700 "$failing_install/docs"
assert_status_nonzero 'mid-update failure: updater exits non-zero'
assert_output_matches 'Update did not complete' 'mid-update failure: announces the incomplete update'
assert_output_matches 'may be partially updated' 'mid-update failure: names the partial state'
assert_output_matches 'git -C .* checkout --' 'mid-update failure: prints a runnable git restore command'
assert_output_matches 'git -C .* clean -nd --' 'mid-update failure: prints the added-file cleanup command'
assert_no_backup_copy "$failing_project" 'mid-update failure'
# No automatic restore: what the pre-clean deleted stays deleted, and the operator is told so.
assert_path_absent "$failing_install/prompts/stale-prompt.md" 'mid-update failure: nothing is silently restored'

# --- Probe 3: a non-git install is warned before the prompt --------------------------------
# Without the copy, an untracked install has no undo at all — the updater must say so rather
# than let the operator infer a safety net that no longer exists.
plain_project="$fixture_root/plain-project"
build_install "$plain_project"

run_updater "$plain_project" n
assert_status 0 'non-git cancel: declining exits 0'
assert_output_matches 'not tracked in git here' 'non-git install: warns that nothing can be restored'
assert_output_matches 'Update cancelled; no files were changed' 'non-git cancel: reports the clean cancel'
assert_no_backup_copy "$plain_project" 'non-git cancel'
if ! grep -q '0\.0\.1' "$plain_project/.claude/skills/do-work/actions/version.md"; then
  printf 'FAIL: non-git cancel — the declined update modified the install.\n' >&2
  fail_count=$((fail_count + 1))
fi

if [ "$fail_count" -gt 0 ]; then
  printf 'update-script behavior probes: %s failure(s).\n' "$fail_count" >&2
  exit 1
fi

printf 'update-script behavior probes passed.\n'
