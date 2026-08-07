#!/usr/bin/env bash
# Hermetic behavioral probes for the legacy-to-suite bridge updater.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
update_script="$repo_root/tools/do-work-update.sh"
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
  printf 'FAIL: bridge updater, suite installer, manifest validator, and section replacer must exist and be executable.\n' >&2
  exit 1
fi

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

build_legacy_install() {
  local project_path="$1"
  local install_path="$project_path/.claude/skills/do-work"
  mkdir -p "$install_path/actions" "$install_path/prompts" "$install_path/interviews" \
    "$install_path/docs" "$install_path/tools" "$project_path/do-work/queue" "$project_path/kb"
  printf '# do-work\n\nLegacy install.\n' > "$install_path/SKILL.md"
  printf '# Version Action\n\n**Current version**: 0.0.1\n' > "$install_path/actions/version.md"
  printf '# stale prompt\n' > "$install_path/prompts/stale.md"
  printf '# stale interview\n' > "$install_path/interviews/stale.md"
  printf '# old guide\n' > "$install_path/docs/guide.md"
  cp "$update_script" "$install_path/tools/do-work-update.sh"
  cp "$manifest_validator" "$install_path/tools/validate-suite-manifest.sh"
  cp "$suite_installer" "$install_path/tools/install-do-work-suite.sh"
  cp "$section_replacer" "$install_path/tools/replace-text-section.sh"
  chmod +x "$install_path/tools/"*.sh
  printf 'queue sentinel\n' > "$project_path/do-work/queue/sentinel.txt"
  printf 'kb sentinel\n' > "$project_path/kb/sentinel.txt"
  printf 'project recipe\n' > "$project_path/Justfile"
  mkdir -p "$project_path/.claude"
  printf '{"hooks":{}}\n' > "$project_path/.claude/settings.json"
}

build_root_legacy_install() {
  local project_path="$1"
  mkdir -p "$project_path/actions" "$project_path/prompts" "$project_path/interviews" \
    "$project_path/docs" "$project_path/tools" "$project_path/do-work/queue" "$project_path/kb"
  printf '# do-work\n\nRoot fallback install.\n' > "$project_path/SKILL.md"
  printf '# Version Action\n\n**Current version**: 0.0.1\n' > "$project_path/actions/version.md"
  printf '# old guide\n' > "$project_path/docs/guide.md"
  cp "$update_script" "$project_path/tools/do-work-update.sh"
  cp "$manifest_validator" "$project_path/tools/validate-suite-manifest.sh"
  cp "$suite_installer" "$project_path/tools/install-do-work-suite.sh"
  cp "$section_replacer" "$project_path/tools/replace-text-section.sh"
  chmod +x "$project_path/tools/"*.sh
  printf 'application sentinel\n' > "$project_path/app.txt"
  printf 'queue sentinel\n' > "$project_path/do-work/queue/sentinel.txt"
  printf 'kb sentinel\n' > "$project_path/kb/sentinel.txt"
  printf 'project recipe\n' > "$project_path/Justfile"
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

build_legacy_archive() {
  local archive_path="$1" tree_root="$fixture_root/legacy-src/do-work-upstream"
  mkdir -p "$tree_root/actions" "$tree_root/prompts" "$tree_root/interviews" \
    "$tree_root/docs" "$tree_root/tools"
  printf '# do-work\n\nLegacy upstream.\n' > "$tree_root/SKILL.md"
  printf '# Version Action\n\n**Current version**: 0.0.2\n' > "$tree_root/actions/version.md"
  printf '# fresh prompt\n' > "$tree_root/prompts/fresh.md"
  printf '# fresh interview\n' > "$tree_root/interviews/fresh.md"
  printf '# new guide\n' > "$tree_root/docs/guide.md"
  cp "$update_script" "$tree_root/tools/do-work-update.sh"
  cp "$manifest_validator" "$tree_root/tools/validate-suite-manifest.sh"
  cp "$suite_installer" "$tree_root/tools/install-do-work-suite.sh"
  cp "$section_replacer" "$tree_root/tools/replace-text-section.sh"
  chmod +x "$tree_root/tools/"*.sh
  tar czf "$archive_path" -C "$fixture_root/legacy-src" do-work-upstream
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

legacy_tarball="$fixture_root/legacy.tar.gz"
suite_tree="$fixture_root/suite-src/do-work-upstream"
suite_tarball="$fixture_root/suite.tar.gz"
build_legacy_archive "$legacy_tarball"
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

# Capability discovery is exact, standalone, and side-effect free.
capability_probe="$fixture_root/capability-project"
mkdir -p "$capability_probe"
before_capability="$(find "$capability_probe" -print)"
probe_output="$(cd "$capability_probe" && bash "$update_script" --capabilities 2>&1)"
probe_status=$?
assert_status 0 'capabilities: exits 0'
if [ "$probe_output" != 'suite-layout-v2' ]; then
  record_failure "capabilities: expected exact suite-layout-v2, got: $probe_output"
fi
after_capability="$(find "$capability_probe" -print)"
if [ "$before_capability" != "$after_capability" ]; then
  record_failure 'capabilities: modified the filesystem'
fi

# Legacy archives still update through the bridge and preserve all project-owned surfaces.
legacy_project="$fixture_root/legacy-project"
build_legacy_install "$legacy_project"
init_project "$legacy_project"
commit_project "$legacy_project" 'legacy install'
run_updater "$legacy_project" y "$legacy_tarball"
assert_status 0 'legacy update: exits 0'
assert_output_matches 'Updated to v0\.0\.2' 'legacy update: reports version'
assert_file_contains "$legacy_project/.claude/skills/do-work/docs/guide.md" 'new guide' \
  'legacy update: installs reviewed bytes'
assert_path_absent "$legacy_project/.claude/skills/do-work/prompts/stale.md" \
  'legacy update: removes stale managed content'
assert_file_contains "$legacy_project/do-work/queue/sentinel.txt" 'queue sentinel' \
  'legacy update: preserves queue runtime'
assert_file_contains "$legacy_project/kb/sentinel.txt" 'kb sentinel' \
  'legacy update: preserves KB runtime'
assert_file_contains "$legacy_project/Justfile" 'project recipe' \
  'legacy update: preserves Justfile'
assert_file_contains "$legacy_project/.claude/settings.json" 'hooks' \
  'legacy update: preserves settings'
if [ "$(wc -l < "$CURL_CALL_LOG" | tr -d ' ')" != 1 ]; then
  record_failure 'legacy update: expected exactly one archive download'
fi

# The Just recipe's repository-root fallback uses the same engine without managing app files.
root_project="$fixture_root/root-project"
build_root_legacy_install "$root_project"
init_project "$root_project"
commit_project "$root_project" 'root fallback install'
: > "$CURL_CALL_LOG"
probe_output="$(printf 'y\n' | FAKE_TARBALL="$legacy_tarball" \
  bash "$root_project/tools/do-work-update.sh" --project-root "$root_project" 2>&1)"
probe_status=$?
assert_status 0 'root fallback: exits 0'
assert_file_contains "$root_project/actions/version.md" '0\.0\.2' \
  'root fallback: advances managed skill files'
assert_file_contains "$root_project/app.txt" 'application sentinel' \
  'root fallback: preserves application file'
assert_file_contains "$root_project/Justfile" 'project recipe' \
  'root fallback: preserves project Justfile'

# A valid future archive installs all four modules as one reviewed transaction.
suite_project="$fixture_root/suite-project"
build_legacy_install "$suite_project"
cat > "$suite_project/Justfile" <<'JUST'
custom-before:
    echo keep

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
JUST
cat > "$suite_project/.claude/settings.json" <<'JSON'
{
  "custom": "keep",
  "hooks": {
    "SessionStart": [
      {"hooks": [{"type": "command", "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/memory-session-start.sh\""}]}
    ],
    "Stop": [
      {"hooks": [{"type": "command", "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/memory-stop-capture.sh\""}]}
    ]
  }
}
JSON
init_project "$suite_project"
commit_project "$suite_project" 'bridge install'
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
  'suite update: migrates the managed Just section'
assert_file_contains "$suite_project/Justfile" '\.claude/skills/do-work-board/tools/queue-kanban' \
  'suite update: points board recipes at the board sibling'
assert_file_contains "$suite_project/Justfile" '^custom-before:$' \
  'suite update: preserves custom Just content'
if grep -q '\.claude/skills/do-work/tools/queue-kanban' "$suite_project/Justfile"; then
  record_failure 'suite update: retained the legacy board recipe path'
fi
assert_file_contains "$suite_project/.claude/settings.json" \
  '\.claude/skills/do-work-knowledge/hooks/memory-session-start\.sh' \
  'suite update: migrates the legacy memory SessionStart path'
assert_file_contains "$suite_project/.claude/settings.json" \
  '\.claude/skills/do-work-knowledge/hooks/memory-stop-capture\.sh' \
  'suite update: migrates the legacy memory Stop path'
assert_file_contains "$suite_project/.claude/settings.json" '"custom"[[:space:]]*:[[:space:]]*"keep"' \
  'suite update: preserves unrelated settings'
if grep -q '\.claude/skills/do-work/hooks/memory-' "$suite_project/.claude/settings.json"; then
  record_failure 'suite update: retained a legacy memory hook path'
fi
if [ "$(wc -l < "$CURL_CALL_LOG" | tr -d ' ')" != 1 ]; then
  record_failure 'suite update: expected exactly one archive download'
fi

# The installed bridge remains the trusted transaction engine. A valid archive cannot
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
commit_project "$trusted_engine_project" 'bridge install'
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
  build_legacy_install "$unsafe_project"
  init_project "$unsafe_project"
  commit_project "$unsafe_project" 'bridge install'
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
build_root_legacy_install "$symlink_project"
mkdir -p "$symlink_target"
ln -s "$symlink_target" "$symlink_project/.claude"
init_project "$symlink_project"
commit_project "$symlink_project" 'root bridge with unsafe client symlink'
: > "$CURL_CALL_LOG"
probe_output="$(printf 'y\n' | FAKE_TARBALL="$suite_tarball" \
  bash "$symlink_project/tools/do-work-update.sh" --project-root "$symlink_project" 2>&1)"
probe_status=$?
assert_status_nonzero 'destination symlink: exits non-zero'
assert_output_matches 'resolves outside the project' \
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
assert_output_matches 'restored every managed path to its exact pre-install state' \
  'partial failure: reports automatic recovery'
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
if ! grep -q 'tools/do-work-update\.sh.*--project-root' "$repo_root/actions/version.md"; then
  record_failure 'entry-point parity: actions/version.md does not delegate to tools/do-work-update.sh'
fi

if [ "$fail_count" -gt 0 ]; then
  printf 'update-script behavior probes: %s failure(s).\n' "$fail_count" >&2
  exit 1
fi

printf 'update-script behavior probes passed.\n'
