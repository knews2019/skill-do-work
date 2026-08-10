#!/usr/bin/env bash
# Hermetic behavioral probes for the current four-module suite updater.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
update_script="$repo_root/skills/do-work/tools/do-work-update.sh"
manifest_validator="$repo_root/tools/validate-suite-manifest.sh"
suite_installer="$repo_root/tools/install-do-work-suite.sh"
section_replacer="$repo_root/tools/replace-text-section.sh"
fail_count=0

for required_command in bash git tar diff; do
  if ! command -v "$required_command" >/dev/null 2>&1; then
    printf 'SKIP: %s unavailable — update-script behavior probes not run.\n' "$required_command"
    exit 0
  fi
done
if [ ! -x "$update_script" ] || [ ! -x "$manifest_validator" ] \
  || [ ! -x "$suite_installer" ] || [ ! -x "$section_replacer" ]; then
  printf 'FAIL: suite updater, installer, manifest validator, and section replacer must exist and be executable.\n' >&2
  exit 1
fi

for retired_updater_token in \
  'suite-layout-v2' \
  '--capabilities' \
  'legacy_shipped_paths' \
  'legacy all-in-one skill'
do
  if grep -Fq -- "$retired_updater_token" "$update_script"; then
    printf 'FAIL: current updater retained migration-window logic: %s\n' "$retired_updater_token" >&2
    fail_count=$((fail_count + 1))
  fi
done

export GIT_CONFIG_GLOBAL=/dev/null
export GIT_CONFIG_SYSTEM=/dev/null
export GIT_TERMINAL_PROMPT=0

fixture_root="$(mktemp -d)"
cleanup_fixture() {
  chmod -R u+rwX "$fixture_root" 2>/dev/null || true
  rm -rf "$fixture_root"
}
trap cleanup_fixture EXIT

probe_output=''
probe_status=0

record_failure() {
  printf 'FAIL: %s\n' "$1" >&2
  fail_count=$((fail_count + 1))
}

assert_status() {
  local expected_status="$1" probe_name="$2"
  if [ "$probe_status" -ne "$expected_status" ]; then
    record_failure "$probe_name — expected exit $expected_status, got $probe_status. Output: $probe_output"
  fi
}

assert_status_nonzero() {
  local probe_name="$1"
  if [ "$probe_status" -eq 0 ]; then
    record_failure "$probe_name — expected a non-zero exit. Output: $probe_output"
  fi
}

assert_output_matches() {
  local pattern_text="$1" probe_name="$2"
  if ! printf '%s' "$probe_output" | grep -Eq -- "$pattern_text"; then
    record_failure "$probe_name — output did not match /$pattern_text/. Output: $probe_output"
  fi
}

assert_output_lacks() {
  local pattern_text="$1" probe_name="$2"
  if printf '%s' "$probe_output" | grep -Eq -- "$pattern_text"; then
    record_failure "$probe_name — output unexpectedly matched /$pattern_text/. Output: $probe_output"
  fi
}

assert_file_contains() {
  local candidate_path="$1" pattern_text="$2" probe_name="$3"
  if [ ! -f "$candidate_path" ] || ! grep -Eq -- "$pattern_text" "$candidate_path"; then
    record_failure "$probe_name — $candidate_path did not contain /$pattern_text/."
  fi
}

assert_path_absent() {
  local candidate_path="$1" probe_name="$2"
  if [ -e "$candidate_path" ] || [ -L "$candidate_path" ]; then
    record_failure "$probe_name — expected $candidate_path to be absent."
  fi
}

init_project() {
  local project_path="$1"
  git -c init.defaultBranch=main init -q "$project_path"
  git -C "$project_path" config user.name 'Fixture Runner'
  git -C "$project_path" config user.email 'fixture@example.invalid'
}

commit_project() {
  local project_path="$1" message_text="$2"
  git -C "$project_path" add -A
  git -C "$project_path" commit -qm "$message_text"
}

build_suite_install() {
  local project_path="$1" module_name module_path
  mkdir -p "$project_path/do-work/queue" "$project_path/kb" "$project_path/.claude/skills"
  for module_name in do-work do-work-board do-work-knowledge do-work-toolbox; do
    module_path="$project_path/.claude/skills/$module_name"
    mkdir -p "$module_path"
    printf '# %s\n\nOld module.\n' "$module_name" > "$module_path/SKILL.md"
    printf 'old %s payload\n' "$module_name" > "$module_path/payload.txt"
  done
  mkdir -p "$project_path/.claude/skills/do-work/actions" \
    "$project_path/.claude/skills/do-work/tools"
  printf '# Version Action\n\n**Current version**: 0.0.1\n' \
    > "$project_path/.claude/skills/do-work/actions/version.md"
  cp "$update_script" "$project_path/.claude/skills/do-work/tools/do-work-update.sh"
  cp "$manifest_validator" \
    "$project_path/.claude/skills/do-work/tools/validate-suite-manifest.sh"
  cp "$suite_installer" \
    "$project_path/.claude/skills/do-work/tools/install-do-work-suite.sh"
  cp "$section_replacer" \
    "$project_path/.claude/skills/do-work/tools/replace-text-section.sh"
  chmod +x "$project_path/.claude/skills/do-work/tools/"*.sh
  printf 'queue sentinel\n' > "$project_path/do-work/queue/sentinel.txt"
  printf 'kb sentinel\n' > "$project_path/kb/sentinel.txt"
  printf 'project-recipe:\n    echo project\n' > "$project_path/Justfile"
  printf '{"hooks":{}}\n' > "$project_path/.claude/settings.json"
}

build_suite_tree() {
  local tree_root="$1" module_name module_path
  mkdir -p "$tree_root/suite" "$tree_root/tools" "$tree_root/skills"
  printf '0.0.2\n' > "$tree_root/VERSION"
  cp "$repo_root/suite/modules.tsv" "$tree_root/suite/modules.tsv"
  cp "$manifest_validator" "$tree_root/tools/validate-suite-manifest.sh"
  cp "$suite_installer" "$tree_root/tools/install-do-work-suite.sh"
  cp "$section_replacer" "$tree_root/tools/replace-text-section.sh"
  chmod +x "$tree_root/tools/"*.sh
  for module_name in do-work do-work-board do-work-knowledge do-work-toolbox; do
    module_path="$tree_root/skills/$module_name"
    mkdir -p "$module_path"
    printf '# %s\n\nNew module.\n' "$module_name" > "$module_path/SKILL.md"
    printf 'new %s payload\n' "$module_name" > "$module_path/payload.txt"
  done
  mkdir -p "$tree_root/skills/do-work/actions" "$tree_root/skills/do-work/tools" \
    "$tree_root/skills/do-work/hooks" "$tree_root/skills/do-work-board"
  printf '# Version Action\n\n**Current version**: 0.0.2\n' \
    > "$tree_root/skills/do-work/actions/version.md"
  printf '0.0.2\n' > "$tree_root/skills/do-work/VERSION"
  cp "$update_script" "$tree_root/skills/do-work/tools/do-work-update.sh"
  cp "$manifest_validator" \
    "$tree_root/skills/do-work/tools/validate-suite-manifest.sh"
  cp "$suite_installer" \
    "$tree_root/skills/do-work/tools/install-do-work-suite.sh"
  cp "$section_replacer" \
    "$tree_root/skills/do-work/tools/replace-text-section.sh"
  chmod +x "$tree_root/skills/do-work/tools/"*.sh
  cp "$repo_root/skills/do-work/hooks/hooks.json" "$tree_root/skills/do-work/hooks/hooks.json"
  cp "$repo_root/skills/do-work-board/justfile.template" \
    "$tree_root/skills/do-work-board/justfile.template"
  printf 'created during update\n' > "$tree_root/skills/do-work/new-core.txt"
}

archive_suite_tree() {
  local tree_root="$1" archive_path="$2" parent_path archive_name
  parent_path="$(dirname "$tree_root")"
  archive_name="$(basename "$tree_root")"
  tar czf "$archive_path" -C "$parent_path" "$archive_name"
}

suite_tree="$fixture_root/suite-src/do-work-upstream"
suite_tarball="$fixture_root/suite.tar.gz"
build_suite_tree "$suite_tree"
archive_suite_tree "$suite_tree" "$suite_tarball"

stub_bin="$fixture_root/stub-bin"
mkdir -p "$stub_bin"
real_cp="$(command -v cp)"
printf '%s\n' '#!/usr/bin/env bash' \
  'destination_path=""' \
  'while [ "$#" -gt 0 ]; do' \
  '  case "$1" in -o) destination_path="$2"; shift 2 ;; *) shift ;; esac' \
  'done' \
  '[ -n "$destination_path" ] || exit 2' \
  'printf "download\\n" >> "$CURL_CALL_LOG"' \
  '"$REAL_CP" "$FAKE_TARBALL" "$destination_path"' > "$stub_bin/curl"
chmod +x "$stub_bin/curl"
export PATH="$stub_bin:$PATH"
export REAL_CP="$real_cp"
export CURL_CALL_LOG="$fixture_root/curl-calls.log"

run_updater() {
  local project_path="$1" answer_text="$2" tarball_path="$3"
  : > "$CURL_CALL_LOG"
  probe_output="$(printf '%s\n' "$answer_text" \
    | FAKE_TARBALL="$tarball_path" \
      bash "$project_path/.claude/skills/do-work/tools/do-work-update.sh" \
        --project-root "$project_path" 2>&1)"
  probe_status=$?
}

# The bridge-only capability probe is no longer a public updater mode.
probe_output="$(bash "$update_script" --capabilities 2>&1)"
probe_status=$?
assert_status_nonzero 'retired capabilities: exits non-zero'

# A current modular archive installs all four modules as one reviewed transaction.
suite_project="$fixture_root/suite-project"
build_suite_install "$suite_project"
cat > "$suite_project/Justfile" <<'JUST'
custom-before:
    echo keep

# >>> do-work:recipes >>>
run-kanban $port="8090":
    .claude/skills/do-work-board/tools/queue-kanban serve --port {{port}}
run-kanban-cli:
    .claude/skills/do-work-board/tools/queue-kanban serve
kanban-static:
    .claude/skills/do-work-board/tools/queue-kanban generate
kanban-summary:
    .claude/skills/do-work-board/tools/queue-kanban summary
run-do-work-update:
    bash .claude/skills/do-work/tools/do-work-update.sh --project-root .
# <<< do-work:recipes <<<
JUST
cat > "$suite_project/.claude/settings.json" <<'JSON'
{
  "custom": "keep",
  "hooks": {
    "SessionStart": [
      {"hooks": [{"type": "command", "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work-knowledge/hooks/memory-session-start.sh\""}]}
    ],
    "Stop": [
      {"hooks": [{"type": "command", "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work-knowledge/hooks/memory-stop-capture.sh\""}]}
    ]
  }
}
JSON
init_project "$suite_project"
commit_project "$suite_project" 'current modular install'
run_updater "$suite_project" y "$suite_tarball"
assert_status 0 'suite update: exits 0'
assert_output_matches 'four-module suite' 'suite update: identifies layout'
for module_name in do-work do-work-board do-work-knowledge do-work-toolbox; do
  assert_file_contains "$suite_project/.claude/skills/$module_name/payload.txt" \
    "new $module_name payload" "suite update: installs $module_name"
done
assert_file_contains "$suite_project/do-work/queue/sentinel.txt" 'queue sentinel' \
  'suite update: preserves queue runtime'
assert_file_contains "$suite_project/kb/sentinel.txt" 'kb sentinel' \
  'suite update: preserves KB runtime'
assert_file_contains "$suite_project/Justfile" '^# >>> do-work:recipes >>>$' \
  'suite update: reconciles the managed Just section'
assert_file_contains "$suite_project/Justfile" '\.claude/skills/do-work-board/tools/queue-kanban' \
  'suite update: points board recipes at the board sibling'
assert_file_contains "$suite_project/Justfile" '^custom-before:$' \
  'suite update: preserves custom Just content'
assert_file_contains "$suite_project/.claude/settings.json" \
  '\.claude/skills/do-work-knowledge/hooks/memory-session-start\.sh' \
  'suite update: preserves the current memory SessionStart path'
assert_file_contains "$suite_project/.claude/settings.json" \
  '\.claude/skills/do-work-knowledge/hooks/memory-stop-capture\.sh' \
  'suite update: preserves the current memory Stop path'
assert_file_contains "$suite_project/.claude/settings.json" '"custom"[[:space:]]*:[[:space:]]*"keep"' \
  'suite update: preserves unrelated settings'
if [ "$(wc -l < "$CURL_CALL_LOG" | tr -d ' ')" != 1 ]; then
  record_failure 'suite update: expected exactly one archive download'
fi

# The managed Just entry point and direct core entry point converge on byte-identical state.
just_entry_project="$fixture_root/just-entry-project"
cp -R "$suite_project" "$just_entry_project"
rm -rf "$just_entry_project/.git"
init_project "$just_entry_project"
commit_project "$just_entry_project" 'modular suite baseline'
run_updater "$suite_project" y "$suite_tarball"
assert_status 0 'entry-point parity: direct updater exits 0'
: > "$CURL_CALL_LOG"
probe_output="$(cd "$just_entry_project" && printf 'y\n' \
  | FAKE_TARBALL="$suite_tarball" just run-do-work-update 2>&1)"
probe_status=$?
assert_status 0 'entry-point parity: managed Just updater exits 0'
if ! diff -qr "$suite_project/.claude/skills" "$just_entry_project/.claude/skills" >/dev/null \
  || ! cmp -s "$suite_project/Justfile" "$just_entry_project/Justfile" \
  || ! cmp -s "$suite_project/.claude/settings.json" "$just_entry_project/.claude/settings.json"; then
  record_failure 'entry-point parity: direct and managed Just updates produced different managed bytes'
fi

# The installed suite remains the trusted transaction engine. A valid archive cannot
# replace that engine with an executable of its own before the reviewed write boundary.
hostile_installer_tree="$fixture_root/hostile-installer-src/do-work-upstream"
hostile_installer_tarball="$fixture_root/hostile-installer.tar.gz"
archive_installer_marker="$fixture_root/archive-installer-ran"
build_suite_tree "$hostile_installer_tree"
printf '%s\n' '#!/usr/bin/env bash' \
  ': > "$ARCHIVE_INSTALLER_MARKER"' \
  'exit 91' > "$hostile_installer_tree/tools/install-do-work-suite.sh"
chmod +x "$hostile_installer_tree/tools/install-do-work-suite.sh"
printf '%s\n' '#!/usr/bin/env bash' \
  ': > "$ARCHIVE_INSTALLER_MARKER"' \
  'exit 92' > "$hostile_installer_tree/skills/do-work/tools/replace-text-section.sh"
chmod +x "$hostile_installer_tree/skills/do-work/tools/replace-text-section.sh"
archive_suite_tree "$hostile_installer_tree" "$hostile_installer_tarball"
trusted_engine_project="$fixture_root/trusted-engine-project"
build_suite_install "$trusted_engine_project"
init_project "$trusted_engine_project"
commit_project "$trusted_engine_project" 'current modular install'
export ARCHIVE_INSTALLER_MARKER="$archive_installer_marker"
run_updater "$trusted_engine_project" y "$hostile_installer_tarball"
assert_status 0 'installer trust: uses installed transaction engine'
assert_path_absent "$archive_installer_marker" \
  'installer trust: does not execute archive transaction helpers'

# Malformed and traversing suite manifests fail before a managed write.
for unsafe_case in malformed traversal; do
  unsafe_tree="$fixture_root/$unsafe_case-src/do-work-upstream"
  unsafe_tarball="$fixture_root/$unsafe_case.tar.gz"
  build_suite_tree "$unsafe_tree"
  if [ "$unsafe_case" = malformed ]; then
    printf 'source\tdestination\textra\n' > "$unsafe_tree/suite/modules.tsv"
  else
    printf 'source\tdestination\n../escape\t.claude/skills/do-work\n' \
      > "$unsafe_tree/suite/modules.tsv"
    printf '#!/usr/bin/env bash\nexit 0\n' > "$unsafe_tree/tools/validate-suite-manifest.sh"
    chmod +x "$unsafe_tree/tools/validate-suite-manifest.sh"
  fi
  archive_suite_tree "$unsafe_tree" "$unsafe_tarball"
  unsafe_project="$fixture_root/$unsafe_case-project"
  build_suite_install "$unsafe_project"
  init_project "$unsafe_project"
  commit_project "$unsafe_project" 'current modular install'
  before_head="$(git -C "$unsafe_project" status --porcelain)"
  run_updater "$unsafe_project" y "$unsafe_tarball"
  assert_status_nonzero "$unsafe_case manifest: exits non-zero"
  if [ "$unsafe_case" = malformed ]; then
    assert_output_matches 'manifest header' "$unsafe_case manifest: trusted validator reports failure"
  else
    assert_output_matches 'source traverses directories' \
      "$unsafe_case manifest: trusted validator rejects traversal despite bundled override"
  fi
  assert_output_lacks 'Continue with this one' \
    "$unsafe_case manifest: fails before the confirmation boundary"
  after_head="$(git -C "$unsafe_project" status --porcelain)"
  if [ "$before_head" != "$after_head" ]; then
    record_failure "$unsafe_case manifest: changed the managed install before validation"
  fi
  assert_path_absent "$fixture_root/escape" "$unsafe_case manifest: no traversal write"
done

# Even a valid textual destination must not escape through an existing client-side symlink.
symlink_project="$fixture_root/symlink-project"
symlink_target="$fixture_root/outside-project"
build_suite_install "$symlink_project"
rm -rf "$symlink_project/.claude/skills/do-work-board"
mkdir -p "$symlink_target"
ln -s "$symlink_target" "$symlink_project/.claude/skills/do-work-board"
init_project "$symlink_project"
commit_project "$symlink_project" 'modular install with unsafe client symlink'
run_updater "$symlink_project" y "$suite_tarball"
assert_status_nonzero 'destination symlink: exits non-zero'
assert_output_matches 'managed destination must not be a symlink|resolves outside the project' \
  'destination symlink: reports physical escape before confirmation'
assert_path_absent "$symlink_target/skills" 'destination symlink: does not write outside project'

# Dirty managed content is named, the discard is explicit, and declining is mutation-free.
dirty_project="$fixture_root/dirty-project"
build_suite_install "$dirty_project"
init_project "$dirty_project"
commit_project "$dirty_project" 'old suite'
printf '# LOCAL BOARD CUSTOMIZATION\n' > "$dirty_project/.claude/skills/do-work-board/SKILL.md"
run_updater "$dirty_project" n "$suite_tarball"
assert_status 0 'dirty cancel: exits 0'
assert_output_matches 'do-work-board/SKILL\.md' 'dirty cancel: lists managed dirty path'
assert_output_matches 'discards those changes' 'dirty cancel: warns before confirmation'
assert_output_matches 'Installation cancelled; no files were changed' 'dirty cancel: reports cancellation'
assert_file_contains "$dirty_project/.claude/skills/do-work-board/SKILL.md" \
  'LOCAL BOARD CUSTOMIZATION' 'dirty cancel: preserves declined customization'

git -C "$dirty_project" add .claude/skills/do-work-board/SKILL.md
run_updater "$dirty_project" y "$suite_tarball"
assert_status 0 'dirty confirm: exits 0'
assert_file_contains "$dirty_project/.claude/skills/do-work-board/SKILL.md" \
  'New module' 'dirty confirm: installs reviewed module after consent'
if [ -n "$(git -C "$dirty_project" diff --cached --name-only -- .claude/skills)" ]; then
  record_failure 'dirty confirm: confirmed discard left staged managed customization behind'
fi

# Fail while copying module two: module one is restored from Git and only new paths vanish.
failure_project="$fixture_root/failure-project"
build_suite_install "$failure_project"
init_project "$failure_project"
commit_project "$failure_project" 'old suite'
printf '%s\n' '#!/usr/bin/env bash' \
  'last_argument="${!#}"' \
  'case "$last_argument" in */.claude/skills/do-work-board/) if [ ! -f "$CP_FAILURE_MARKER" ]; then : > "$CP_FAILURE_MARKER"; exit 88; fi ;; esac' \
  'exec "$REAL_CP" "$@"' > "$stub_bin/cp"
chmod +x "$stub_bin/cp"
CP_FAILURE_MARKER="$fixture_root/cp-failed-once" run_updater "$failure_project" y "$suite_tarball"
rm -f "$stub_bin/cp"
assert_status_nonzero 'partial failure: exits non-zero'
assert_output_matches \
  'restored every managed path and the Git index to their exact pre-install state' \
  'partial failure: reports automatic filesystem and index recovery'
assert_file_contains "$failure_project/.claude/skills/do-work/payload.txt" \
  'old do-work payload' 'partial failure: restores first changed module'
assert_file_contains "$failure_project/.claude/skills/do-work-board/payload.txt" \
  'old do-work-board payload' 'partial failure: restores failed module'
assert_path_absent "$failure_project/.claude/skills/do-work/new-core.txt" \
  'partial failure: cleans only newly created file'
assert_file_contains "$failure_project/do-work/queue/sentinel.txt" 'queue sentinel' \
  'partial failure: preserves queue runtime'
assert_file_contains "$failure_project/kb/sentinel.txt" 'kb sentinel' \
  'partial failure: preserves KB runtime'
assert_file_contains "$failure_project/Justfile" 'project-recipe' \
  'partial failure: preserves Justfile'
assert_file_contains "$failure_project/.claude/settings.json" 'hooks' \
  'partial failure: preserves settings'
if [ -n "$(git -C "$failure_project" status --porcelain)" ]; then
  record_failure "partial failure: recovery did not return the managed worktree to HEAD: $(git -C "$failure_project" status --porcelain)"
fi

# The agent-facing path must delegate mutation to the same tested engine.
if ! grep -q 'tools/do-work-update\.sh.*--project-root' "$repo_root/skills/do-work/actions/version.md"; then
  record_failure 'entry-point parity: modular actions/version.md does not delegate to tools/do-work-update.sh'
fi

if [ "$fail_count" -gt 0 ]; then
  printf 'update-script behavior probes: %s failure(s).\n' "$fail_count" >&2
  exit 1
fi

printf 'update-script behavior probes passed.\n'
