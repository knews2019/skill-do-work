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

managed_state_paths=(
  '.claude/skills/do-work'
  '.claude/skills/do-work-board'
  '.claude/skills/do-work-knowledge'
  '.claude/skills/do-work-toolbox'
  'justfile'
  '.claude/settings.json'
  'CLAUDE.md'
)

snapshot_install_state() {
  local project_root="$1" snapshot_root="$2" relative_path snapshot_path
  mkdir -p "$snapshot_root/working-tree"
  : > "$snapshot_root/existing-paths"
  for relative_path in "${managed_state_paths[@]}"; do
    if [ -e "$project_root/$relative_path" ] || [ -L "$project_root/$relative_path" ]; then
      snapshot_path="$snapshot_root/working-tree/$relative_path"
      mkdir -p "$(dirname "$snapshot_path")"
      cp -Rp "$project_root/$relative_path" "$snapshot_path"
      printf '%s\n' "$relative_path" >> "$snapshot_root/existing-paths"
    fi
  done
  git -C "$project_root" diff --binary --no-ext-diff > "$snapshot_root/git.diff"
  git -C "$project_root" diff --cached --binary --no-ext-diff > "$snapshot_root/git.cached.diff"
  git -C "$project_root" status --porcelain=v1 --untracked-files=all -z \
    > "$snapshot_root/git.status"
}

assert_install_state_unchanged() {
  local project_root="$1" snapshot_root="$2" scenario_name="$3"
  local relative_path snapshot_path existed_before=''
  for relative_path in "${managed_state_paths[@]}"; do
    snapshot_path="$snapshot_root/working-tree/$relative_path"
    existed_before=''
    grep -Fxq "$relative_path" "$snapshot_root/existing-paths" && existed_before=1
    if [ -n "$existed_before" ]; then
      if [ ! -e "$project_root/$relative_path" ] && [ ! -L "$project_root/$relative_path" ]; then
        fail "$scenario_name removed managed working-tree path $relative_path"
      elif ! diff -qr "$snapshot_path" "$project_root/$relative_path" >/dev/null; then
        fail "$scenario_name changed managed working-tree bytes at $relative_path"
      fi
    elif [ -e "$project_root/$relative_path" ] || [ -L "$project_root/$relative_path" ]; then
      fail "$scenario_name created managed working-tree path $relative_path"
    fi
  done

  git -C "$project_root" diff --binary --no-ext-diff > "$snapshot_root/git.diff.after"
  git -C "$project_root" diff --cached --binary --no-ext-diff \
    > "$snapshot_root/git.cached.diff.after"
  git -C "$project_root" status --porcelain=v1 --untracked-files=all -z \
    > "$snapshot_root/git.status.after"
  cmp -s "$snapshot_root/git.diff" "$snapshot_root/git.diff.after" \
    || fail "$scenario_name changed git diff"
  cmp -s "$snapshot_root/git.cached.diff" "$snapshot_root/git.cached.diff.after" \
    || fail "$scenario_name changed git diff --cached"
  cmp -s "$snapshot_root/git.status" "$snapshot_root/git.status.after" \
    || fail "$scenario_name changed porcelain status"
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
    if [ ! -f "$project_root/.claude/skills/$module/SKILL.md" ] \
      || [ ! -s "$project_root/.claude/skills/$module/SKILL.md" ]; then
      fail "installer did not install a non-empty regular $module/SKILL.md"
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

for retired_installer_token in \
  '--migrate-legacy-do-work' \
  '.claude/skills/do-work/hooks/memory-session-start.sh' \
  '.claude/skills/do-work/hooks/memory-stop-capture.sh'
do
  if grep -Fq -- "$retired_installer_token" "$installer"; then
    fail "suite installer retained migration-window logic: $retired_installer_token"
  fi
done

for cutover_export_path in VERSION suite skills; do
  if git -C "$repo_root" check-attr export-ignore -- "$cutover_export_path" \
    | grep -q 'export-ignore: set'; then
    fail "fresh-install archive still excludes /$cutover_export_path"
  fi
done

for retired_runtime_path in SKILL.md actions tools/do-work-update.sh tools/queue-kanban; do
  if [ -f "$repo_root/$retired_runtime_path" ] \
    || { [ -d "$repo_root/$retired_runtime_path" ] \
      && find "$repo_root/$retired_runtime_path" -type f \
        ! -path "$repo_root/tools/queue-kanban/queue-kanban" -print -quit \
        | grep -q .; }; then
    fail "fresh-install source still carries legacy root runtime: $retired_runtime_path"
  fi
done

# The installer delegates to do-work-cli, so build the binary once here. The restricted-PATH
# lanes below deliberately omit `go`, and the launcher rebuilds whenever any source is newer
# than the binary — `cp -R` copies in readdir order, so the mtimes are set explicitly rather
# than left to copy order (REQ-407 C10).
cli_module_root="$repo_root/skills/do-work/tools/do-work-cli"
if ! (cd "$cli_module_root" && go build -o do-work-cli ./cmd/do-work-cli); then
  printf 'FAIL: could not pre-build do-work-cli for the restricted-PATH lanes.\n' >&2
  exit 1
fi
touch "$cli_module_root/do-work-cli"

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
touch "$archive_root/skills/do-work/tools/do-work-cli/do-work-cli"
archive_file="$workdir/do-work-suite.tar.gz"
tar czf "$archive_file" -C "$archive_parent" skill-do-work-main

# Reject a suite whose runtime-reported action version differs before creating any
# managed path. Fresh installs have no backup to restore, so validation must precede
# the first write rather than relying on the updater's later comparison.
mismatched_archive_parent="$workdir/mismatched-archive-source"
mismatched_archive_root="$mismatched_archive_parent/skill-do-work-main"
mkdir -p "$mismatched_archive_parent"
cp -R "$archive_root" "$mismatched_archive_root"
sed 's/^\*\*Current version\*\*: .*/**Current version**: 9.9.9/' \
  "$mismatched_archive_root/skills/do-work/actions/version.md" \
  > "$mismatched_archive_root/skills/do-work/actions/version.md.next"
mv "$mismatched_archive_root/skills/do-work/actions/version.md.next" \
  "$mismatched_archive_root/skills/do-work/actions/version.md"
mismatched_archive_file="$workdir/mismatched-suite.tar.gz"
tar czf "$mismatched_archive_file" -C "$mismatched_archive_parent" skill-do-work-main
mismatched_project="$workdir/mismatched-version-project"
new_git_project "$mismatched_project"
mismatched_status=0
run_installer "$mismatched_project" "$mismatched_archive_file" \
  "$workdir/mismatched-version.out" || mismatched_status=$?
if [ "$mismatched_status" -eq 0 ]; then
  fail 'installer accepted an archive whose actions/version.md disagreed with suite VERSION'
elif [ -e "$mismatched_project/.claude" ] || [ -e "$mismatched_project/justfile" ]; then
  fail 'version-mismatched archive created managed paths before rejection'
fi

# The exact printed bootstrap must need no installed skill and fetch only this one archive.
bootstrap_command="$(bash "$installer" --print-bootstrap-command)"
assert_output_contains "$bootstrap_command" 'main\.tar\.gz' 'canonical bootstrap must fetch the upstream source archive'
assert_output_contains "$bootstrap_command" '--archive "\$archive_file"' 'canonical bootstrap must pass the downloaded archive to the installer instead of downloading a second artifact'
readme_bootstrap_command="$(awk '
  /^## Installation$/ {in_installation=1; next}
  in_installation && /^```bash$/ {in_command=1; next}
  in_command && /^```$/ {exit}
  in_command {print}
' "$repo_root/README.md")"
if [ "$readme_bootstrap_command" != "$bootstrap_command" ]; then
  fail 'README installation block must be byte-identical to --print-bootstrap-command output'
fi

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
fresh_recipe_count="$(grep -Ec '^[A-Za-z0-9_-]+([^:]*)?:$' "$fresh_project/justfile")"
[ "$fresh_recipe_count" -eq 40 ] || fail "fresh install published $fresh_recipe_count managed recipes, want 40"
if [ ! -f "$fresh_project/.claude/settings.json" ]; then
  fail 'fresh install did not create Claude settings for core hooks'
else
  assert_file_contains "$fresh_project/.claude/settings.json" 'do-work/hooks/session-start\.sh' 'fresh settings must enable the core SessionStart hook'
  if grep -q 'do-work/hooks/pipeline-guard\.sh' "$fresh_project/.claude/settings.json"; then
    fail 'fresh settings must not enable the retired pipeline Stop hook'
  fi
  if grep -q 'memory-' "$fresh_project/.claude/settings.json"; then
    fail 'fresh install must not enable memory hooks'
  fi
fi
assert_output_contains "$(cat "$bootstrap_output")" 'Installed do-work suite v[0-9]+\.[0-9]+\.[0-9]+' 'bootstrap must report the verified installed suite version'
if [ "$(grep -c '^Install this complete four-skill suite?' "$bootstrap_output")" -ne 1 ]; then
  fail 'fresh bootstrap must present exactly one confirmation boundary'
fi
if ! cmp -s "$fresh_project/CLAUDE.md" "$archive_root/skills/do-work/agent-instructions.template.md"; then
  fail 'fresh install did not create agent instructions from the core-owned template'
fi
assert_output_contains "$(cat "$bootstrap_output")" '--- managed configuration: CLAUDE\.md ---' \
  'install review must present the agent instructions surface'

# Reinstall preserves custom Just bytes and composes current modular hooks.
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
      {"hooks": [{"type": "command", "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work-knowledge/hooks/memory-session-start.sh\""}]}
    ],
    "Stop": [
      {"hooks": [
        {"type": "command", "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/pipeline-guard.sh\""},
        {"type": "command", "command": "echo custom-stop"}
      ]},
      {"hooks": [
        {"type": "command", "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/pipeline-guard.sh\""}
      ]},
      {"matcher": "preserve-empty", "hooks": []},
      {"hooks": [{"type": "command", "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work-knowledge/hooks/memory-stop-capture.sh\""}]}
    ]
  }
}
JSON
chmod 600 "$fresh_project/.claude/settings.json"
cat > "$fresh_project/CLAUDE.md" <<'MD'
# Consumer Project

Custom guidance before the managed section.

<!-- >>> do-work:communication-style >>> -->
stale managed link
<!-- <<< do-work:communication-style <<< -->

Custom guidance after the managed section.
MD
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
just_mode="$(stat -c '%a' "$fresh_project/justfile" 2>/dev/null || stat -f '%Lp' "$fresh_project/justfile")"
[ "$just_mode" = 640 ] || fail "reinstall changed Justfile mode (got $just_mode, want 640)"
settings_mode="$(stat -c '%a' "$fresh_project/.claude/settings.json" 2>/dev/null || stat -f '%Lp' "$fresh_project/.claude/settings.json")"
[ "$settings_mode" = 600 ] || fail "reinstall changed settings mode (got $settings_mode, want 600)"
assert_file_contains "$fresh_project/CLAUDE.md" '^Custom guidance before the managed section\.$' \
  'reinstall changed agent-instructions content before the managed section'
assert_file_contains "$fresh_project/CLAUDE.md" '^Custom guidance after the managed section\.$' \
  'reinstall changed agent-instructions content after the managed section'
if grep -q '^stale managed link$' "$fresh_project/CLAUDE.md"; then
  fail 'reinstall retained stale content inside the managed agent-instructions section'
fi
assert_file_contains "$fresh_project/CLAUDE.md" 'crew-members/communication-style\.md' \
  'reinstall did not link the communication-style crew member from agent instructions'
if [ "$(grep -cF '<!-- >>> do-work:communication-style >>> -->' "$fresh_project/CLAUDE.md")" -ne 1 ]; then
  fail 'reinstall must leave exactly one managed agent-instructions section'
fi
# Composition is now a byte contract, not a structural one: an order-preserving encoder
# replaced the jq and Python reconcilers, so the whole composed document is compared rather
# than probed field by field. A reordered key is exactly the silent regression this catches.
cat > "$workdir/expected-settings.json" <<'JSON'
{
  "custom": {
    "keep": [
      1,
      2,
      3
    ]
  },
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "echo custom-start"
          }
        ]
      },
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work-knowledge/hooks/memory-session-start.sh\""
          }
        ]
      },
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/session-start.sh\""
          }
        ]
      }
    ],
    "Stop": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "echo custom-stop"
          }
        ]
      },
      {
        "matcher": "preserve-empty",
        "hooks": []
      },
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work-knowledge/hooks/memory-stop-capture.sh\""
          }
        ]
      }
    ]
  }
}
JSON
if ! cmp -s "$workdir/expected-settings.json" "$fresh_project/.claude/settings.json"; then
  fail "reinstall did not compose the exact expected settings bytes: $(diff -u "$workdir/expected-settings.json" "$fresh_project/.claude/settings.json" | head -n 20)"
fi
cp "$fresh_project/justfile" "$workdir/reinstall.just.snapshot"
cp "$fresh_project/.claude/settings.json" "$workdir/reinstall.settings.snapshot"
if ! run_installer "$fresh_project" "$archive_file" "$workdir/reinstall-idempotent.out"; then
  fail 'idempotent reinstall failed'
elif ! cmp -s "$fresh_project/justfile" "$workdir/reinstall.just.snapshot" \
  || ! cmp -s "$fresh_project/.claude/settings.json" "$workdir/reinstall.settings.snapshot"; then
  fail 'reinstall is not byte-idempotent for reconciled Just/settings files'
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

# A BOM-prefixed reserved recipe is rejected without Just, before confirmation or client mutation.
no_just_path="$workdir/no-just-path"
mkdir -p "$no_just_path"
for command_name in awk bash cat cp diff dirname find git grep gzip mkdir mktemp mv rm sed tar; do
  command_path="$(command -v "$command_name" 2>/dev/null || true)"
  [ -z "$command_path" ] || ln -s "$command_path" "$no_just_path/$command_name"
done
collision_project="$workdir/reserved-recipe-collision"
new_git_project "$collision_project"
mkdir -p "$collision_project/.claude"
printf '\357\273\277run-kanban:\r\n    echo custom collision\r\n' > "$collision_project/justfile"
printf '{"custom":"unchanged"}\n' > "$collision_project/.claude/settings.json"
chmod 640 "$collision_project/justfile"
chmod 600 "$collision_project/.claude/settings.json"
git -C "$collision_project" add justfile .claude/settings.json
git -C "$collision_project" commit -qm 'collision fixture'
cp "$collision_project/justfile" "$workdir/collision.just.before"
cp "$collision_project/.claude/settings.json" "$workdir/collision.settings.before"
collision_mode="$(stat -c '%a' "$collision_project/justfile" 2>/dev/null || stat -f '%Lp' "$collision_project/justfile")"
collision_status_before="$(git -C "$collision_project" status --porcelain --untracked-files=all)"
collision_output="$workdir/reserved-recipe-collision.out"
collision_exit_status=0
printf 'y\n' | PATH="$no_just_path" bash "$installer" --project-root "$collision_project" --archive "$archive_file" >"$collision_output" 2>&1 \
  || collision_exit_status=$?
if [ "$collision_exit_status" -eq 0 ]; then
  fail 'installer accepted a BOM-prefixed reserved recipe when Just was unavailable'
else
  assert_file_contains "$collision_output" 'reserved Just recipe or alias outside managed section: run-kanban' 'installer collision error did not name the reserved recipe'
  if grep -Fq 'Install this complete four-skill suite?' "$collision_output"; then
    fail 'installer asked for confirmation before rejecting the reserved recipe collision'
  fi
  if ! cmp -s "$collision_project/justfile" "$workdir/collision.just.before" \
    || ! cmp -s "$collision_project/.claude/settings.json" "$workdir/collision.settings.before" \
    || [ -e "$collision_project/.claude/skills/do-work" ]; then
    fail 'reserved recipe rejection changed Justfile, settings, or modules'
  fi
  collision_mode_after="$(stat -c '%a' "$collision_project/justfile" 2>/dev/null || stat -f '%Lp' "$collision_project/justfile")"
  [ "$collision_mode_after" = "$collision_mode" ] || fail "reserved recipe rejection changed Justfile mode (got $collision_mode_after, want $collision_mode)"
  collision_status_after="$(git -C "$collision_project" status --porcelain --untracked-files=all)"
  [ "$collision_status_after" = "$collision_status_before" ] || fail 'reserved recipe rejection changed Git status'
fi

# The fallback scanner must also retain a recipe header whose triple-quoted default
# contains a quote of the same type. Without Just, this is the only parse boundary.
multiline_collision_project="$workdir/multiline-reserved-recipe-collision"
new_git_project "$multiline_collision_project"
printf "run-kanban value='''\npayload's\n''':\n    echo custom collision\n" \
  > "$multiline_collision_project/justfile"
cp "$multiline_collision_project/justfile" "$workdir/multiline-collision.just.before"
multiline_collision_output="$workdir/multiline-reserved-recipe-collision.out"
multiline_collision_exit_status=0
printf 'y\n' | PATH="$no_just_path" bash "$installer" \
  --project-root "$multiline_collision_project" --archive "$archive_file" \
  >"$multiline_collision_output" 2>&1 || multiline_collision_exit_status=$?
if [ "$multiline_collision_exit_status" -eq 0 ]; then
  fail 'installer accepted a delimiter-bearing reserved recipe when Just was unavailable'
else
  assert_file_contains "$multiline_collision_output" \
    'reserved Just recipe or alias outside managed section: run-kanban' \
    'installer did not name the delimiter-bearing reserved recipe collision'
  if grep -Fq 'Install this complete four-skill suite?' "$multiline_collision_output"; then
    fail 'installer asked for confirmation before rejecting the delimiter-bearing recipe collision'
  fi
  if ! cmp -s "$multiline_collision_project/justfile" "$workdir/multiline-collision.just.before" \
    || [ -e "$multiline_collision_project/.claude" ]; then
    fail 'delimiter-bearing reserved recipe rejection changed the Justfile or installed modules'
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

# A directory named SKILL.md is not an executable skill and must fail validation
# before the installer writes any managed path.
directory_skill_parent="$workdir/directory-skill-source"
cp -R "$archive_parent" "$directory_skill_parent"
rm "$directory_skill_parent/skill-do-work-main/skills/do-work-toolbox/SKILL.md"
mkdir "$directory_skill_parent/skill-do-work-main/skills/do-work-toolbox/SKILL.md"
directory_skill_archive="$workdir/directory-skill-suite.tar.gz"
tar czf "$directory_skill_archive" -C "$directory_skill_parent" skill-do-work-main
directory_skill_project="$workdir/directory-skill-project"
new_git_project "$directory_skill_project"
if run_installer "$directory_skill_project" "$directory_skill_archive" \
  "$workdir/directory-skill.out"; then
  fail 'installer accepted a suite with a directory-shaped toolbox SKILL.md'
elif [ -e "$directory_skill_project/.claude/skills" ] \
  || [ -e "$directory_skill_project/justfile" ]; then
  fail 'directory-shaped SKILL.md validation failure wrote client files'
fi

# A post-write Just validation failure restores exact module, Just, and settings originals.
rollback_project="$workdir/rollback project"
cp -R "$fresh_project" "$rollback_project"
rollback_snapshot="$workdir/rollback-snapshot"
mkdir -p "$rollback_snapshot/.claude/skills" "$rollback_snapshot/.claude"
cp -R "$rollback_project/.claude/skills/." "$rollback_snapshot/.claude/skills/"
cp "$rollback_project/justfile" "$rollback_snapshot/justfile"
cp "$rollback_project/.claude/settings.json" "$rollback_snapshot/.claude/settings.json"
rollback_state_snapshot="$workdir/rollback-state-before"
snapshot_install_state "$rollback_project" "$rollback_state_snapshot"
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

# On a case-insensitive filesystem, path probes must still record the real
# directory-entry spelling so a failed install restores Justfile as Justfile.
case_recovery_project="$workdir/case-recovery-project"
new_git_project "$case_recovery_project"
printf 'custom-case:\n    echo preserved\n' > "$case_recovery_project/Justfile"
cp "$case_recovery_project/Justfile" "$workdir/case-recovery.before"
case_recovery_status=0
printf 'y\n' | DO_WORK_TEST_JUST_COUNT="$workdir/case-recovery-just-count" \
  PATH="$flaky_bin:$PATH" bash "$installer" --project-root "$case_recovery_project" \
    --archive "$archive_file" >"$workdir/case-recovery.out" 2>&1 \
  || case_recovery_status=$?
if [ "$case_recovery_status" -eq 0 ]; then
  fail 'case-preserving recovery fixture reported success after post-write Just validation failed'
elif ! cmp -s "$case_recovery_project/Justfile" "$workdir/case-recovery.before"; then
  fail 'failed install did not restore the original Justfile bytes'
elif ! find "$case_recovery_project" -mindepth 1 -maxdepth 1 -name 'Justfile' -print -quit | grep -q . \
  || find "$case_recovery_project" -mindepth 1 -maxdepth 1 -name 'justfile' -print -quit | grep -q .; then
  fail 'failed install did not restore the real Justfile directory-entry spelling'
fi

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
assert_install_state_unchanged "$rollback_project" "$rollback_state_snapshot" \
  'post-write Just failure recovery'

# TERM during the module write phase runs the same all-or-recover path.
interrupt_project="$workdir/interruption-project"
cp -R "$fresh_project" "$interrupt_project"
interrupt_snapshot="$workdir/interruption-snapshot"
mkdir -p "$interrupt_snapshot/.claude/skills" "$interrupt_snapshot/.claude"
cp -R "$interrupt_project/.claude/skills/." "$interrupt_snapshot/.claude/skills/"
cp "$interrupt_project/justfile" "$interrupt_snapshot/justfile"
cp "$interrupt_project/.claude/settings.json" "$interrupt_snapshot/.claude/settings.json"
interrupt_state_snapshot="$workdir/interruption-state-before"
snapshot_install_state "$interrupt_project" "$interrupt_state_snapshot"
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
  assert_file_contains "$workdir/interruption.out" \
    'restored every managed path and the Git index to their exact pre-install state' \
    'interruption did not complete the filesystem and index recovery path'
fi
assert_install_state_unchanged "$interrupt_project" "$interrupt_state_snapshot" \
  'interrupted install recovery'

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
cancel_state_snapshot="$workdir/cancel-state-before"
snapshot_install_state "$cancel_project" "$cancel_state_snapshot"
cancel_status=0
printf 'n\n' | bash "$installer" --project-root "$cancel_project" --archive "$archive_file" >"$workdir/cancel.out" 2>&1 \
  || cancel_status=$?
if [ "$cancel_status" -ne 0 ]; then
  fail 'declining installation should be a successful no-op'
elif [ -e "$cancel_project/.claude" ] || [ -e "$cancel_project/justfile" ]; then
  fail 'declining installation changed the project'
fi
assert_install_state_unchanged "$cancel_project" "$cancel_state_snapshot" \
  'ordinary cancellation'

# Recovery is exact for working-tree bytes and the Git index across every dirty-state shape.
transaction_base_project="$workdir/transaction-base"
new_git_project "$transaction_base_project"
if ! run_installer "$transaction_base_project" "$archive_file" "$workdir/transaction-base.out"; then
  fail 'could not build the transaction recovery fixture'
fi
git -C "$transaction_base_project" add .
git -C "$transaction_base_project" commit -qm 'installed suite baseline'

for dirty_state in staged-only unstaged-only partially-staged; do
  state_project="$workdir/recovery-$dirty_state"
  cp -R "$transaction_base_project" "$state_project"
  state_managed_file="$state_project/.claude/skills/do-work/SKILL.md"
  case "$dirty_state" in
    staged-only)
      printf '\nSTAGED-ONLY CUSTOMIZATION\n' >> "$state_managed_file"
      git -C "$state_project" add .claude/skills/do-work/SKILL.md
      ;;
    unstaged-only)
      printf '\nUNSTAGED-ONLY CUSTOMIZATION\n' >> "$state_managed_file"
      ;;
    partially-staged)
      printf '\nSTAGED PARTIAL CUSTOMIZATION\n' >> "$state_managed_file"
      git -C "$state_project" add .claude/skills/do-work/SKILL.md
      printf 'UNSTAGED PARTIAL CUSTOMIZATION\n' >> "$state_managed_file"
      ;;
  esac
  state_snapshot="$workdir/recovery-$dirty_state-before"
  snapshot_install_state "$state_project" "$state_snapshot"
  state_status=0
  printf 'y\n' | DO_WORK_TEST_JUST_COUNT="$workdir/recovery-$dirty_state-just-count" \
    PATH="$flaky_bin:$PATH" bash "$installer" --project-root "$state_project" \
      --archive "$archive_file" > "$workdir/recovery-$dirty_state.out" 2>&1 \
    || state_status=$?
  if [ "$state_status" -eq 0 ]; then
    fail "$dirty_state recovery fixture reported success after post-write verification failed"
  else
    assert_install_state_unchanged "$state_project" "$state_snapshot" "$dirty_state recovery"
    assert_file_contains "$workdir/recovery-$dirty_state.out" \
      'restored every managed path and the Git index to their exact pre-install state' \
      "$dirty_state recovery did not report exact filesystem and index restoration"
  fi
done

# A failure inside the unstage loop must recover the index mutation that already succeeded.
unstage_failure_project="$workdir/unstage-failure"
cp -R "$transaction_base_project" "$unstage_failure_project"
printf '\nFIRST MODULE STAGED CUSTOMIZATION\n' \
  >> "$unstage_failure_project/.claude/skills/do-work/SKILL.md"
git -C "$unstage_failure_project" add .claude/skills/do-work/SKILL.md
unstage_failure_snapshot="$workdir/unstage-failure-before"
snapshot_install_state "$unstage_failure_project" "$unstage_failure_snapshot"
unstage_failure_bin="$workdir/unstage-failure-bin"
mkdir -p "$unstage_failure_bin"
real_git_path="$(command -v git)"
cat > "$unstage_failure_bin/git" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
if [[ " $* " == *" restore --staged "* ]]; then
  restore_count=0
  [ ! -f "$DO_WORK_TEST_GIT_COUNT" ] || restore_count="$(cat "$DO_WORK_TEST_GIT_COUNT")"
  restore_count=$((restore_count + 1))
  printf '%s\n' "$restore_count" > "$DO_WORK_TEST_GIT_COUNT"
  [ "$restore_count" -lt 2 ] || exit 86
fi
exec "$DO_WORK_TEST_REAL_GIT" "$@"
SH
chmod +x "$unstage_failure_bin/git"
unstage_failure_status=0
printf 'y\n' | DO_WORK_TEST_GIT_COUNT="$workdir/unstage-git-count" \
  DO_WORK_TEST_REAL_GIT="$real_git_path" PATH="$unstage_failure_bin:$PATH" \
  bash "$installer" --project-root "$unstage_failure_project" --archive "$archive_file" \
    > "$workdir/unstage-failure.out" 2>&1 || unstage_failure_status=$?
if [ "$unstage_failure_status" -eq 0 ]; then
  fail 'installer reported success after a failure inside the managed unstage loop'
else
  assert_install_state_unchanged "$unstage_failure_project" "$unstage_failure_snapshot" \
    'unstage-loop failure recovery'
  assert_file_contains "$workdir/unstage-failure.out" \
    'restored every managed path and the Git index to their exact pre-install state' \
    'unstage-loop failure did not run exact filesystem and index recovery'
fi

# Cancellation occurs before the transaction and preserves partial staging exactly.
dirty_cancel_project="$workdir/dirty-cancel"
cp -R "$transaction_base_project" "$dirty_cancel_project"
printf '\nCANCELLED STAGED CUSTOMIZATION\n' \
  >> "$dirty_cancel_project/.claude/skills/do-work/SKILL.md"
git -C "$dirty_cancel_project" add .claude/skills/do-work/SKILL.md
printf 'CANCELLED UNSTAGED CUSTOMIZATION\n' \
  >> "$dirty_cancel_project/.claude/skills/do-work/SKILL.md"
dirty_cancel_snapshot="$workdir/dirty-cancel-before"
snapshot_install_state "$dirty_cancel_project" "$dirty_cancel_snapshot"
dirty_cancel_status=0
printf 'n\n' | bash "$installer" --project-root "$dirty_cancel_project" \
  --archive "$archive_file" > "$workdir/dirty-cancel.out" 2>&1 || dirty_cancel_status=$?
if [ "$dirty_cancel_status" -ne 0 ]; then
  fail 'declining installation with a partially staged module should be a successful no-op'
else
  assert_install_state_unchanged "$dirty_cancel_project" "$dirty_cancel_snapshot" \
    'dirty cancellation'
fi

# Success still removes both staged and unstaged customizations beneath installed bytes.
successful_transaction_project="$workdir/successful-transaction"
cp -R "$transaction_base_project" "$successful_transaction_project"
printf '\nSUCCESS STAGED CUSTOMIZATION\n' \
  >> "$successful_transaction_project/.claude/skills/do-work/SKILL.md"
git -C "$successful_transaction_project" add .claude/skills/do-work/SKILL.md
printf 'SUCCESS UNSTAGED CUSTOMIZATION\n' \
  >> "$successful_transaction_project/.claude/skills/do-work/SKILL.md"
if ! run_installer "$successful_transaction_project" "$archive_file" \
  "$workdir/successful-transaction.out"; then
  fail 'ordinary installation failed while discarding confirmed managed customizations'
elif grep -q 'SUCCESS .* CUSTOMIZATION' \
  "$successful_transaction_project/.claude/skills/do-work/SKILL.md"; then
  fail 'successful installation retained a confirmed managed customization'
elif [ -n "$(git -C "$successful_transaction_project" diff --cached --name-only -- .claude/skills)" ]; then
  fail 'successful installation retained staged managed customizations in the index'
else
  assert_four_modules "$successful_transaction_project"
fi

if [ "$fail_count" -ne 0 ]; then
  exit 1
fi

printf 'suite installer behavior probes passed.\n'
