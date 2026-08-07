#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
installer="$repo_root/tools/install-do-work-suite.sh"
fail_count=0

fail() {
  printf 'FAIL: %s\n' "$1" >&2
  fail_count=$((fail_count + 1))
}

assert_file_contains() {
  local file_path="$1" pattern="$2" message="$3"
  grep -Eq -- "$pattern" "$file_path" || fail "$message"
}

assert_output_contains() {
  local output="$1" pattern="$2" message="$3"
  grep -Eq -- "$pattern" <<<"$output" || fail "$message"
}

new_git_project() {
  local project_root="$1"
  mkdir -p "$project_root"
  git init -q "$project_root"
  git -C "$project_root" config user.email test@example.com
  git -C "$project_root" config user.name 'Suite Installer Test'
}

run_installer() {
  local project_root="$1" archive_file="$2" output_file="$3"
  local run_status=0
  printf 'y\n' | bash "$installer" --project-root "$project_root" --archive "$archive_file" >"$output_file" 2>&1 \
    || run_status=$?
  return "$run_status"
}

assert_four_modules() {
  local project_root="$1" module module_count
  for module in do-work do-work-board do-work-knowledge do-work-toolbox; do
    if [ ! -s "$project_root/.claude/skills/$module/SKILL.md" ]; then
      fail "installer did not install a non-empty $module/SKILL.md"
    fi
  done
  module_count="$(find "$project_root/.claude/skills" -mindepth 1 -maxdepth 1 -type d -name 'do-work*' | wc -l | tr -d ' ')"
  if [ "$module_count" -ne 4 ]; then
    fail "installer must leave exactly four do-work module directories (found $module_count)"
  fi
}

if [ ! -x "$installer" ]; then
  printf 'FAIL: tools/install-do-work-suite.sh is missing or not executable.\n' >&2
  exit 1
fi

workdir="$(mktemp -d "${TMPDIR:-/tmp}/do-work-suite-install-test.XXXXXX")"
trap 'rm -rf "$workdir"' EXIT

archive_parent="$workdir/archive-source"
archive_root="$archive_parent/skill-do-work-main"
mkdir -p "$archive_root/tools"
cp "$repo_root/VERSION" "$archive_root/VERSION"
cp -R "$repo_root/suite" "$archive_root/suite"
cp -R "$repo_root/skills" "$archive_root/skills"
cp "$installer" "$archive_root/tools/install-do-work-suite.sh"
cp "$repo_root/tools/validate-suite-manifest.sh" "$archive_root/tools/validate-suite-manifest.sh"
cp "$repo_root/tools/replace-text-section.sh" "$archive_root/tools/replace-text-section.sh"
chmod +x "$archive_root/tools/"*.sh
archive_file="$workdir/do-work-suite.tar.gz"
tar czf "$archive_file" -C "$archive_parent" skill-do-work-main

# The exact printed bootstrap must need no installed skill and fetch only this one archive.
bootstrap_command="$(bash "$installer" --print-bootstrap-command)"
assert_output_contains "$bootstrap_command" 'main\.tar\.gz' 'canonical bootstrap must fetch the upstream source archive'
assert_output_contains "$bootstrap_command" '--archive "\$archive_file"' 'canonical bootstrap must pass the downloaded archive to the installer instead of downloading a second artifact'

fake_bin="$workdir/fake-bin"
mkdir -p "$fake_bin"
cat > "$fake_bin/curl" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
output_file=''
while [ "$#" -gt 0 ]; do
  case "$1" in
    -o) output_file="$2"; shift 2 ;;
    *) shift ;;
  esac
done
[ -n "$output_file" ]
count=0
[ ! -f "$DO_WORK_TEST_CURL_COUNT" ] || count="$(cat "$DO_WORK_TEST_CURL_COUNT")"
printf '%s\n' "$((count + 1))" > "$DO_WORK_TEST_CURL_COUNT"
cp "$DO_WORK_TEST_ARCHIVE" "$output_file"
SH
chmod +x "$fake_bin/curl"

fresh_project="$workdir/fresh project with spaces"
new_git_project "$fresh_project"
bootstrap_output="$workdir/bootstrap.out"
bootstrap_status=0
printf 'y\n' | (
  cd "$fresh_project"
  DO_WORK_TEST_ARCHIVE="$archive_file" \
    DO_WORK_TEST_CURL_COUNT="$workdir/curl-count" \
    PATH="$fake_bin:$PATH" \
    bash -c "$bootstrap_command"
) >"$bootstrap_output" 2>&1 || bootstrap_status=$?
if [ "$bootstrap_status" -ne 0 ]; then
  fail "canonical bootstrap failed in a clean Git project: $(tail -n 5 "$bootstrap_output")"
else
  assert_four_modules "$fresh_project"
fi
if [ "$(cat "$workdir/curl-count" 2>/dev/null || printf 0)" -ne 1 ]; then
  fail 'canonical bootstrap must download exactly one upstream artifact'
fi
if ! cmp -s "$fresh_project/justfile" "$archive_root/skills/do-work-board/justfile.template"; then
  fail 'fresh install did not create the complete Justfile from the board-owned template'
fi
if [ ! -f "$fresh_project/.claude/settings.json" ]; then
  fail 'fresh install did not create Claude settings for core hooks'
else
  assert_file_contains "$fresh_project/.claude/settings.json" 'do-work/hooks/session-start\.sh' 'fresh settings must enable the core SessionStart hook'
  assert_file_contains "$fresh_project/.claude/settings.json" 'do-work/hooks/pipeline-guard\.sh' 'fresh settings must enable the current core Stop hook'
  if grep -q 'memory-' "$fresh_project/.claude/settings.json"; then
    fail 'fresh install must not enable memory hooks'
  fi
fi
assert_output_contains "$(cat "$bootstrap_output")" 'Installed do-work suite v[0-9]+\.[0-9]+\.[0-9]+' 'bootstrap must report the verified installed suite version'
if [ "$(grep -c '^Install this complete four-skill suite?' "$bootstrap_output")" -ne 1 ]; then
  fail 'fresh bootstrap must present exactly one confirmation boundary'
fi

# Reinstall preserves custom Just bytes, composes custom hooks, and migrates only known memory paths.
cat > "$fresh_project/justfile" <<'JUST'
set shell := ["bash", "-cu"]
custom-before:
    echo before

# >>> do-work:recipes >>>
old-managed:
    echo old
# <<< do-work:recipes <<<

custom-after:
    echo after
JUST
chmod 640 "$fresh_project/justfile"
cat > "$fresh_project/.claude/settings.json" <<'JSON'
{
  "custom": {"keep": [1, 2, 3]},
  "hooks": {
    "SessionStart": [
      {"hooks": [{"type": "command", "command": "echo custom-start"}]},
      {"hooks": [{"type": "command", "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/memory-session-start.sh\""}]}
    ],
    "Stop": [
      {"hooks": [{"type": "command", "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/memory-stop-capture.sh\""}]},
      {"hooks": [{"type": "command", "command": "echo custom-stop"}]}
    ]
  }
}
JSON
chmod 600 "$fresh_project/.claude/settings.json"
reinstall_output="$workdir/reinstall.out"
if ! run_installer "$fresh_project" "$archive_file" "$reinstall_output"; then
  fail "reinstall failed: $(tail -n 5 "$reinstall_output")"
else
  assert_four_modules "$fresh_project"
fi
assert_file_contains "$fresh_project/justfile" '^custom-before:$' 'reinstall changed custom Just content before the managed section'
assert_file_contains "$fresh_project/justfile" '^custom-after:$' 'reinstall changed custom Just content after the managed section'
if grep -q '^old-managed:' "$fresh_project/justfile"; then
  fail 'reinstall retained stale content inside the managed Just section'
fi
just_mode="$(stat -f '%Lp' "$fresh_project/justfile" 2>/dev/null || stat -c '%a' "$fresh_project/justfile")"
[ "$just_mode" = 640 ] || fail "reinstall changed Justfile mode (got $just_mode, want 640)"
settings_mode="$(stat -f '%Lp' "$fresh_project/.claude/settings.json" 2>/dev/null || stat -c '%a' "$fresh_project/.claude/settings.json")"
[ "$settings_mode" = 600 ] || fail "reinstall changed settings mode (got $settings_mode, want 600)"
python3 - "$fresh_project/.claude/settings.json" <<'PY' || fail 'reinstall did not preserve/compose/migrate hook settings'
import json, pathlib, sys
data = json.loads(pathlib.Path(sys.argv[1]).read_text())
assert data["custom"] == {"keep": [1, 2, 3]}
serialized = json.dumps(data)
assert ".claude/skills/do-work/hooks/memory-session-start.sh" not in serialized
assert ".claude/skills/do-work/hooks/memory-stop-capture.sh" not in serialized
assert serialized.count(".claude/skills/do-work-knowledge/hooks/memory-session-start.sh") == 1
assert serialized.count(".claude/skills/do-work-knowledge/hooks/memory-stop-capture.sh") == 1
assert serialized.count("do-work/hooks/session-start.sh") == 1
assert serialized.count("do-work/hooks/pipeline-guard.sh") == 1
assert serialized.count("echo custom-start") == 1
assert serialized.count("echo custom-stop") == 1
PY
cp "$fresh_project/justfile" "$workdir/reinstall.just.snapshot"
cp "$fresh_project/.claude/settings.json" "$workdir/reinstall.settings.snapshot"
if ! run_installer "$fresh_project" "$archive_file" "$workdir/reinstall-idempotent.out"; then
  fail 'idempotent reinstall failed'
elif ! cmp -s "$fresh_project/justfile" "$workdir/reinstall.just.snapshot" \
  || ! cmp -s "$fresh_project/.claude/settings.json" "$workdir/reinstall.settings.snapshot"; then
  fail 'reinstall is not byte-idempotent for reconciled Just/settings files'
fi

# Exact legacy recipes migrate while surrounding custom content survives.
legacy_project="$workdir/legacy-project"
new_git_project "$legacy_project"
cat > "$legacy_project/Justfile" <<'JUST'
custom-before:
    echo before

# --- do-work board recipes (installed by `do-work install just-kanban`) ---
run-kanban $port="8090":
    echo board
run-kanban-cli:
    echo cli
kanban-static:
    echo static
kanban-summary:
    echo summary
run-do-work-update:
    echo update

custom-after:
    echo after
JUST
if ! run_installer "$legacy_project" "$archive_file" "$workdir/legacy.out"; then
  fail 'installer could not migrate the exact legacy five-recipe block'
else
  assert_file_contains "$legacy_project/Justfile" '^# >>> do-work:recipes >>>$' 'legacy migration did not add managed markers'
  assert_file_contains "$legacy_project/Justfile" '^custom-before:$' 'legacy migration changed custom prefix'
  assert_file_contains "$legacy_project/Justfile" '^custom-after:$' 'legacy migration changed custom suffix'
fi

# A marker-free custom Justfile is preserved and receives one managed append.
custom_project="$workdir/custom-project"
new_git_project "$custom_project"
printf 'custom-only:\n    echo untouched\n' > "$custom_project/justfile"
if ! run_installer "$custom_project" "$archive_file" "$workdir/custom.out"; then
  fail 'installer could not append to a marker-free custom Justfile'
else
  assert_file_contains "$custom_project/justfile" '^custom-only:$' 'custom Just recipe was overwritten'
  if [ "$(grep -c '^# >>> do-work:recipes >>>$' "$custom_project/justfile")" -ne 1 ]; then
    fail 'custom Justfile did not receive exactly one managed section'
  fi
fi

# Invalid Just ownership or invalid JSON is rejected before any module/configuration write.
invalid_just_project="$workdir/invalid-just"
new_git_project "$invalid_just_project"
printf '# >>> do-work:recipes >>>\nbroken:\n    echo broken\n' > "$invalid_just_project/.justfile"
cp "$invalid_just_project/.justfile" "$workdir/invalid-just.before"
if run_installer "$invalid_just_project" "$archive_file" "$workdir/invalid-just.out"; then
  fail 'installer accepted malformed Just ownership markers'
elif ! cmp -s "$invalid_just_project/.justfile" "$workdir/invalid-just.before" \
  || [ -e "$invalid_just_project/.claude/skills/do-work" ]; then
  fail 'malformed Just validation changed configuration or installed modules before failing'
fi

invalid_settings_project="$workdir/invalid-settings"
new_git_project "$invalid_settings_project"
mkdir -p "$invalid_settings_project/.claude"
printf '{ invalid json\n' > "$invalid_settings_project/.claude/settings.json"
cp "$invalid_settings_project/.claude/settings.json" "$workdir/invalid-settings.before"
if run_installer "$invalid_settings_project" "$archive_file" "$workdir/invalid-settings.out"; then
  fail 'installer accepted invalid existing Claude settings'
elif ! cmp -s "$invalid_settings_project/.claude/settings.json" "$workdir/invalid-settings.before" \
  || [ -e "$invalid_settings_project/.claude/skills/do-work" ] \
  || [ -e "$invalid_settings_project/justfile" ]; then
  fail 'invalid settings validation did not leave every managed target unchanged'
fi

# Missing module content is rejected before writing anything.
corrupt_parent="$workdir/corrupt-source"
cp -R "$archive_parent" "$corrupt_parent"
: > "$corrupt_parent/skill-do-work-main/skills/do-work-toolbox/SKILL.md"
corrupt_archive="$workdir/corrupt-suite.tar.gz"
tar czf "$corrupt_archive" -C "$corrupt_parent" skill-do-work-main
corrupt_project="$workdir/corrupt-project"
new_git_project "$corrupt_project"
if run_installer "$corrupt_project" "$corrupt_archive" "$workdir/corrupt.out"; then
  fail 'installer accepted a suite with an empty toolbox SKILL.md'
elif [ -e "$corrupt_project/.claude/skills" ] || [ -e "$corrupt_project/justfile" ]; then
  fail 'suite validation failure wrote client files'
fi

# Python is the JSON fallback when jq is absent.
python_path="$workdir/python-path"
mkdir -p "$python_path"
for command_name in awk bash cat chmod cmp cp diff dirname find git grep head mkdir mktemp mv python3 rm sed stat tar tr wc; do
  command_path="$(command -v "$command_name" 2>/dev/null || true)"
  [ -z "$command_path" ] || ln -s "$command_path" "$python_path/$command_name"
done
python_project="$workdir/python-fallback"
new_git_project "$python_project"
python_output="$workdir/python-fallback.out"
python_status=0
printf 'y\n' | PATH="$python_path" bash "$installer" --project-root "$python_project" --archive "$archive_file" >"$python_output" 2>&1 \
  || python_status=$?
if [ "$python_status" -ne 0 ]; then
  fail "Python JSON fallback install failed: $(tail -n 5 "$python_output")"
else
  assert_output_contains "$(cat "$python_output")" 'settings reconciler: python3' 'installer did not report/use Python when jq was unavailable'
  assert_file_contains "$python_project/.claude/settings.json" 'do-work/hooks/session-start\.sh' 'Python fallback did not compose core hooks'
fi

# With neither jq nor Python, a fresh project still gets modules/Justfile while settings stay exact and a precise manual step is printed.
no_json_path="$workdir/no-json-path"
mkdir -p "$no_json_path"
for command_name in awk bash cat chmod cmp cp diff dirname find git grep head mkdir mktemp mv rm sed stat tar tr wc; do
  command_path="$(command -v "$command_name" 2>/dev/null || true)"
  [ -z "$command_path" ] || ln -s "$command_path" "$no_json_path/$command_name"
done
manual_project="$workdir/manual-settings"
new_git_project "$manual_project"
mkdir -p "$manual_project/.claude"
printf '{"custom":"unchanged"}\n' > "$manual_project/.claude/settings.json"
cp "$manual_project/.claude/settings.json" "$workdir/manual-settings.before"
manual_output="$workdir/manual.out"
manual_status=0
printf 'y\n' | PATH="$no_json_path" bash "$installer" --project-root "$manual_project" --archive "$archive_file" >"$manual_output" 2>&1 \
  || manual_status=$?
if [ "$manual_status" -ne 0 ]; then
  fail "no-JSON-tool install failed: $(tail -n 5 "$manual_output")"
else
  assert_four_modules "$manual_project"
  cmp -s "$manual_project/.claude/settings.json" "$workdir/manual-settings.before" \
    || fail 'no-JSON-tool path changed settings instead of leaving them exact'
  assert_output_contains "$(cat "$manual_output")" '^MANUAL STEP: merge \.claude/skills/do-work/hooks/hooks\.json into \.claude/settings\.json; preserve every existing entry\.$' 'no-JSON-tool path did not print the exact manual hook instruction'
fi

# A post-write Just validation failure restores exact module, Just, and settings originals.
rollback_project="$workdir/rollback project"
cp -R "$fresh_project" "$rollback_project"
rollback_snapshot="$workdir/rollback-snapshot"
mkdir -p "$rollback_snapshot/.claude/skills" "$rollback_snapshot/.claude"
cp -R "$rollback_project/.claude/skills/." "$rollback_snapshot/.claude/skills/"
cp "$rollback_project/justfile" "$rollback_snapshot/justfile"
cp "$rollback_project/.claude/settings.json" "$rollback_snapshot/.claude/settings.json"
flaky_bin="$workdir/flaky-bin"
mkdir -p "$flaky_bin"
cat > "$flaky_bin/just" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
count=0
[ ! -f "$DO_WORK_TEST_JUST_COUNT" ] || count="$(cat "$DO_WORK_TEST_JUST_COUNT")"
count=$((count + 1))
printf '%s\n' "$count" > "$DO_WORK_TEST_JUST_COUNT"
[ "$count" -eq 1 ]
SH
chmod +x "$flaky_bin/just"
rollback_status=0
printf 'y\n' | DO_WORK_TEST_JUST_COUNT="$workdir/just-count" PATH="$flaky_bin:$PATH" \
  bash "$installer" --project-root "$rollback_project" --archive "$archive_file" >"$workdir/rollback.out" 2>&1 \
  || rollback_status=$?
if [ "$rollback_status" -eq 0 ]; then
  fail 'installer reported success after post-write Just validation failed'
elif ! diff -qr "$rollback_snapshot/.claude/skills" "$rollback_project/.claude/skills" >/dev/null \
  || ! cmp -s "$rollback_snapshot/justfile" "$rollback_project/justfile" \
  || ! cmp -s "$rollback_snapshot/.claude/settings.json" "$rollback_project/.claude/settings.json"; then
  fail 'post-write Just failure did not restore exact managed originals'
fi

# A post-write settings validation failure restores the same exact originals.
settings_rollback_project="$workdir/settings-rollback"
cp -R "$fresh_project" "$settings_rollback_project"
settings_rollback_snapshot="$workdir/settings-rollback-snapshot"
mkdir -p "$settings_rollback_snapshot/.claude/skills" "$settings_rollback_snapshot/.claude"
cp -R "$settings_rollback_project/.claude/skills/." "$settings_rollback_snapshot/.claude/skills/"
cp "$settings_rollback_project/justfile" "$settings_rollback_snapshot/justfile"
cp "$settings_rollback_project/.claude/settings.json" "$settings_rollback_snapshot/.claude/settings.json"
flaky_jq_bin="$workdir/flaky-jq-bin"
mkdir -p "$flaky_jq_bin"
cat > "$flaky_jq_bin/jq" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
count=0
[ ! -f "$DO_WORK_TEST_JQ_COUNT" ] || count="$(cat "$DO_WORK_TEST_JQ_COUNT")"
count=$((count + 1))
printf '%s\n' "$count" > "$DO_WORK_TEST_JQ_COUNT"
[ "$count" -lt 3 ] || exit 1
exec /usr/bin/jq "$@"
SH
chmod +x "$flaky_jq_bin/jq"
settings_rollback_status=0
printf 'y\n' | DO_WORK_TEST_JQ_COUNT="$workdir/jq-count" PATH="$flaky_jq_bin:$PATH" \
  bash "$installer" --project-root "$settings_rollback_project" --archive "$archive_file" >"$workdir/settings-rollback.out" 2>&1 \
  || settings_rollback_status=$?
if [ "$settings_rollback_status" -eq 0 ]; then
  fail 'installer reported success after post-write settings validation failed'
elif ! diff -qr "$settings_rollback_snapshot/.claude/skills" "$settings_rollback_project/.claude/skills" >/dev/null \
  || ! cmp -s "$settings_rollback_snapshot/justfile" "$settings_rollback_project/justfile" \
  || ! cmp -s "$settings_rollback_snapshot/.claude/settings.json" "$settings_rollback_project/.claude/settings.json"; then
  fail 'post-write settings failure did not restore exact managed originals'
fi

# TERM during the module write phase runs the same all-or-recover path.
interrupt_project="$workdir/interruption-project"
cp -R "$fresh_project" "$interrupt_project"
interrupt_snapshot="$workdir/interruption-snapshot"
mkdir -p "$interrupt_snapshot/.claude/skills" "$interrupt_snapshot/.claude"
cp -R "$interrupt_project/.claude/skills/." "$interrupt_snapshot/.claude/skills/"
cp "$interrupt_project/justfile" "$interrupt_snapshot/justfile"
cp "$interrupt_project/.claude/settings.json" "$interrupt_snapshot/.claude/settings.json"
interrupt_bin="$workdir/interrupt-bin"
mkdir -p "$interrupt_bin"
cat > "$interrupt_bin/cp" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
/bin/cp "$@"
last_argument="${!#}"
case "$last_argument" in
  */.claude/skills/do-work-board/) kill -TERM "$PPID" ;;
esac
SH
chmod +x "$interrupt_bin/cp"
interrupt_status=0
{
  printf 'y\n' | PATH="$interrupt_bin:$PATH" bash "$installer" --project-root "$interrupt_project" --archive "$archive_file"
} >"$workdir/interruption.out" 2>&1 || interrupt_status=$?
if [ "$interrupt_status" -eq 0 ]; then
  fail 'installer reported success after a TERM during module installation'
elif ! diff -qr "$interrupt_snapshot/.claude/skills" "$interrupt_project/.claude/skills" >/dev/null \
  || ! cmp -s "$interrupt_snapshot/justfile" "$interrupt_project/justfile" \
  || ! cmp -s "$interrupt_snapshot/.claude/settings.json" "$interrupt_project/.claude/settings.json"; then
  fail 'interrupted install did not recover the exact prior managed state'
else
  assert_file_contains "$workdir/interruption.out" 'restored every managed path to its exact pre-install state' 'interruption did not complete the recovery path'
fi

# --project-root must name the Git worktree root.
non_git_project="$workdir/not-git"
mkdir -p "$non_git_project"
if run_installer "$non_git_project" "$archive_file" "$workdir/non-git.out"; then
  fail 'installer accepted a non-Git project root'
fi
mkdir -p "$fresh_project/subdir"
if run_installer "$fresh_project/subdir" "$archive_file" "$workdir/subdir.out"; then
  fail 'installer accepted a subdirectory instead of the Git worktree root'
fi

# Declining the sole confirmation leaves a clean project untouched.
cancel_project="$workdir/cancel-project"
new_git_project "$cancel_project"
cancel_status=0
printf 'n\n' | bash "$installer" --project-root "$cancel_project" --archive "$archive_file" >"$workdir/cancel.out" 2>&1 \
  || cancel_status=$?
if [ "$cancel_status" -ne 0 ]; then
  fail 'declining installation should be a successful no-op'
elif [ -e "$cancel_project/.claude" ] || [ -e "$cancel_project/justfile" ]; then
  fail 'declining installation changed the project'
fi

if [ "$fail_count" -ne 0 ]; then
  exit 1
fi

printf 'suite installer behavior probes passed.\n'
