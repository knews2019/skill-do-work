#!/usr/bin/env bash
# Hermetic behavioral probes for the current four-module suite updater.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
# shellcheck source=_dev/tests/fixture-repo.sh
source "$repo_root/_dev/tests/fixture-repo.sh"
update_script="$repo_root/skills/do-work/tools/do-work-update.sh"
manifest_validator="$repo_root/tools/validate-suite-manifest.sh"
suite_installer="$repo_root/tools/install-do-work-suite.sh"
section_replacer="$repo_root/tools/replace-text-section.sh"
upstream_fetcher_source="$repo_root/tools/fetch-upstream-archive.sh"
do_work_cli_launcher="$repo_root/skills/do-work/tools/do-work-cli.sh"
do_work_cli_module="$repo_root/skills/do-work/tools/do-work-cli"
fail_count=0

# The updater, installer, validator, replacer and fetcher are all launchers over do-work-cli
# now, so the command is built once here and copied into every fixture below.
if ! (cd "$do_work_cli_module" && go build -o do-work-cli ./cmd/do-work-cli); then
  printf 'FAIL: could not pre-build do-work-cli for the update fixtures.\n' >&2
  exit 1
fi
touch "$do_work_cli_module/do-work-cli"

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

fixture_root="${DO_WORK_TEST_FIXTURE_ROOT:-$(mktemp -d)}"
cleanup_fixture() {
	if [ -n "${archive_server_pid:-}" ]; then
		kill "$archive_server_pid" 2>/dev/null || true
		wait "$archive_server_pid" 2>/dev/null || true
	fi
  [ -z "${DO_WORK_TEST_FIXTURE_ROOT:-}" ] || return 0
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

# copy_do_work_cli_module places the Go command beside a fixture's launchers, with the binary
# pre-built and its mtime pinned past every source. The launcher rebuilds whenever a source
# is newer, and these fixtures must run without `go` deciding the outcome (REQ-407 C10).
copy_do_work_cli_module() {
  local tools_directory="$1"
  cp "$do_work_cli_launcher" "$tools_directory/do-work-cli.sh"
  chmod +x "$tools_directory/do-work-cli.sh"
  mkdir -p "$tools_directory/do-work-cli"
  cp -R "$do_work_cli_module/." "$tools_directory/do-work-cli/"
  find "$tools_directory/do-work-cli" -type f \( -name '*.go' -o -name 'go.mod' -o -name 'go.sum' \) \
    -exec touch -t 200001010000 {} +
  touch "$tools_directory/do-work-cli/do-work-cli"
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
  cp "$upstream_fetcher_source" \
    "$project_path/.claude/skills/do-work/tools/fetch-upstream-archive.sh"
  copy_do_work_cli_module "$project_path/.claude/skills/do-work/tools"
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
  cp "$upstream_fetcher_source" "$tree_root/tools/fetch-upstream-archive.sh"
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
  cp "$upstream_fetcher_source" \
    "$tree_root/skills/do-work/tools/fetch-upstream-archive.sh"
  copy_do_work_cli_module "$tree_root/skills/do-work/tools"
  chmod +x "$tree_root/skills/do-work/tools/"*.sh
  cp "$repo_root/skills/do-work/hooks/hooks.json" "$tree_root/skills/do-work/hooks/hooks.json"
  cp "$repo_root/skills/do-work/agent-instructions.template.md" \
    "$tree_root/skills/do-work/agent-instructions.template.md"
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
export PATH="$stub_bin:$PATH"
export REAL_CP="$real_cp"
HTTP_CALL_LOG="$fixture_root/http-calls.log"
archive_server_source="$fixture_root/archive-server.go"
archive_server_binary="$fixture_root/archive-server"
cat > "$archive_server_source" <<'GOEOF'
package main
import("fmt";"net";"net/http";"os";"strconv")
func main(){payload,err:=os.ReadFile(os.Args[1]);if err!=nil{panic(err)};listener,err:=net.Listen("tcp","127.0.0.1:0");if err!=nil{panic(err)};status:=200;if len(os.Args)>4{status,_=strconv.Atoi(os.Args[4])};if err:=os.WriteFile(os.Args[2],[]byte("http://"+listener.Addr().String()),0600);err!=nil{panic(err)};handler:=http.HandlerFunc(func(response http.ResponseWriter,request *http.Request){handle,_:=os.OpenFile(os.Args[3],os.O_CREATE|os.O_APPEND|os.O_WRONLY,0600);fmt.Fprintln(handle,request.URL.Path);handle.Close();response.WriteHeader(status);_,_=response.Write(payload)});panic(http.Serve(listener,handler))}
GOEOF
if ! go build -o "$archive_server_binary" "$archive_server_source"; then
  printf 'FAIL: could not build the in-process archive fixture server.\n' >&2
  exit 1
fi
archive_server_pid=''
archive_server_url=''
start_archive_server() {
  local tarball_path="$1" status_code="${2:-200}" url_file="$fixture_root/archive-server-url"
  : > "$HTTP_CALL_LOG"
  rm -f "$url_file"
  "$archive_server_binary" "$tarball_path" "$url_file" "$HTTP_CALL_LOG" "$status_code" &
  archive_server_pid=$!
  wait_count=0
  while [ ! -s "$url_file" ] && [ "$wait_count" -lt 100 ]; do sleep 0.02; wait_count=$((wait_count + 1)); done
  archive_server_url="$(cat "$url_file" 2>/dev/null)"
  [ -n "$archive_server_url" ] || { record_failure 'archive fixture server did not start'; return 1; }
}
stop_archive_server() {
  [ -n "$archive_server_pid" ] || return 0
  kill "$archive_server_pid" 2>/dev/null || true
  wait "$archive_server_pid" 2>/dev/null || true
  archive_server_pid=''
}

run_updater() {
  local project_path="$1" answer_text="$2" tarball_path="$3"
  start_archive_server "$tarball_path" || return
  probe_output="$(printf '%s\n' "$answer_text" \
    | DO_WORK_UPSTREAM_URL="$archive_server_url/archive/refs/heads/main.tar.gz" \
      bash "$project_path/.claude/skills/do-work/tools/do-work-update.sh" \
        --project-root "$project_path" 2>&1)"
  probe_status=$?
  stop_archive_server
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
updated_recipe_summary="$(just --justfile "$suite_project/Justfile" --summary | tr ' ' '\n' | grep -v '^custom-before$' | LC_ALL=C sort | paste -sd' ' -)"
template_recipe_summary="$(just --justfile "$repo_root/skills/do-work-board/justfile.template" --summary | tr ' ' '\n' | LC_ALL=C sort | paste -sd' ' -)"
if [ "$updated_recipe_summary" != "$template_recipe_summary" ]; then
  record_failure 'suite update: managed recipe surface differs from the board-owned template'
fi
assert_file_contains "$suite_project/.claude/settings.json" \
  '\.claude/skills/do-work-knowledge/hooks/memory-session-start\.sh' \
  'suite update: preserves the current memory SessionStart path'
assert_file_contains "$suite_project/.claude/settings.json" \
  '\.claude/skills/do-work-knowledge/hooks/memory-stop-capture\.sh' \
  'suite update: preserves the current memory Stop path'
assert_file_contains "$suite_project/.claude/settings.json" '"custom"[[:space:]]*:[[:space:]]*"keep"' \
  'suite update: preserves unrelated settings'
if [ "$(wc -l < "$HTTP_CALL_LOG" | tr -d ' ')" != 1 ]; then
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
start_archive_server "$suite_tarball"
probe_output="$(cd "$just_entry_project" && printf 'y\n' \
  | DO_WORK_UPSTREAM_URL="$archive_server_url/archive/refs/heads/main.tar.gz" just run-do-work-update 2>&1)"
probe_status=$?
stop_archive_server
assert_status 0 'entry-point parity: managed Just updater exits 0'
# The built do-work-cli binary is build output, not suite content: Go embeds the build
# directory in it, so two projects that build it under different paths hold different bytes
# for identical sources. It is filtered out of the diff's OUTPUT rather than with `diff -x`,
# because -x matches the basenames of directories too and `do-work-cli` is also the module
# directory's name — excluding by name would blind this check to the entire Go source tree.
parity_diff_status=0
parity_report="$(diff -qr "$suite_project/.claude/skills" "$just_entry_project/.claude/skills" 2>&1)" \
  || parity_diff_status=$?
if [ "$parity_diff_status" -gt 1 ]; then
  record_failure "entry-point parity: managed skill trees could not be compared: $parity_report"
else
  parity_differences="$(printf '%s\n' "$parity_report" \
    | grep -v 'tools/do-work-cli/do-work-cli ' | grep -v '^$' || true)"
  if [ -n "$parity_differences" ]; then
    record_failure "entry-point parity: direct and managed Just updates produced different managed bytes: $parity_differences"
  fi
fi
if ! cmp -s "$suite_project/Justfile" "$just_entry_project/Justfile" \
  || ! cmp -s "$suite_project/.claude/settings.json" "$just_entry_project/.claude/settings.json"; then
  record_failure 'entry-point parity: direct and managed Just updates produced different managed configuration bytes'
fi
# Both entry points must have installed a runnable command, whoever built it.
for parity_project in "$suite_project" "$just_entry_project"; do
  if [ ! -x "$parity_project/.claude/skills/do-work/tools/do-work-cli/do-work-cli" ]; then
    record_failure "entry-point parity: ${parity_project##*/} has no executable do-work-cli after the update"
  fi
done

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
# The trusted-engine guarantee now covers the Go command as well: a hostile archive that
# plants its own do-work-cli.sh and its own Go source must have neither built nor run.
printf '%s\n' '#!/usr/bin/env bash' \
  ': > "$ARCHIVE_INSTALLER_MARKER"' \
  'exit 93' > "$hostile_installer_tree/skills/do-work/tools/do-work-cli.sh"
chmod +x "$hostile_installer_tree/skills/do-work/tools/do-work-cli.sh"
mkdir -p "$hostile_installer_tree/skills/do-work/tools/do-work-cli/cmd/do-work-cli"
printf '%s\n' \
  'package main' \
  'import "os"' \
  'func main() { _ = os.WriteFile(os.Getenv("ARCHIVE_INSTALLER_MARKER"), nil, 0o644); os.Exit(94) }' \
  > "$hostile_installer_tree/skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go"
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

# --- Upstream archive fetcher (REQ-217) ---------------------------------
# The reported incident was a sustained 429 across ~2 minutes: retry alone does not
# close it, and `git clone` was the only probe that succeeded. These three cases pin
# the fallback, the export-ignore guarantee, and the preserved-target failure path.
archive_fetcher="$repo_root/tools/fetch-upstream-archive.sh"
if [ ! -x "$archive_fetcher" ]; then
  record_failure 'upstream fetcher: tools/fetch-upstream-archive.sh must exist and be executable'
else
  fetch_root="$fixture_root/upstream-fetch"
  mkdir -p "$fetch_root"

  # A local repository standing in for the upstream remote, carrying one path that
  # .gitattributes marks export-ignore.
  upstream_fixture_repo="$fetch_root/upstream-repo"
  fixture_repo_init "$upstream_fixture_repo"
  mkdir -p "$upstream_fixture_repo/private-path"
  mkdir -p "$upstream_fixture_repo/suite"
  printf 'maintainer only\n' > "$upstream_fixture_repo/private-path/notes.md"
  printf 'shipped\n' > "$upstream_fixture_repo/VERSION"
  printf 'core\tskills/do-work\t.claude/skills/do-work\n' > "$upstream_fixture_repo/suite/modules.tsv"
  printf 'default branch\n' > "$upstream_fixture_repo/default-branch-marker.txt"
  printf '/private-path export-ignore\n' > "$upstream_fixture_repo/.gitattributes"
  fixture_repo_commit_all "$upstream_fixture_repo" 'upstream fixture'
  git -C "$upstream_fixture_repo" checkout -qb requested-branch
  rm "$upstream_fixture_repo/default-branch-marker.txt"
  printf 'requested branch\n' > "$upstream_fixture_repo/requested-branch-marker.txt"
  fixture_repo_commit_all "$upstream_fixture_repo" 'requested branch fixture'
  git -C "$upstream_fixture_repo" checkout -q main
  git -C "$upstream_fixture_repo" checkout -qb feature/fix
  printf 'slashed branch\n' > "$upstream_fixture_repo/slashed-branch-marker.txt"
  fixture_repo_commit_all "$upstream_fixture_repo" 'slashed branch fixture'
  git -C "$upstream_fixture_repo" checkout -q main

  # Case 1: a host that answers 429 forever must not stop the fetch — the git route wins.
  start_archive_server "$suite_tarball" 429
  rate_limited_url="$archive_server_url"

  fallback_archive="$fetch_root/fallback.tar.gz"
  fallback_report="$(bash "$archive_fetcher" "$fallback_archive" \
    "$rate_limited_url/archive/refs/heads/main.tar.gz" "$upstream_fixture_repo" 2>/dev/null)" \
    || record_failure 'upstream fetcher: sustained rate limiting was not survived by the git route'
  if ! tar tzf "$fallback_archive" >/dev/null 2>&1; then
    record_failure 'upstream fetcher: the git route did not produce a readable archive'
  fi
  case "$fallback_report" in
    *git*) ;;
    *) record_failure "upstream fetcher: the succeeding route was not reported (got: $fallback_report)" ;;
  esac

  # Case 2: the git route must honor export-ignore. A cp -R / rsync / tar-the-clone
  # implementation passes cases 1 and 3 and fails only here.
  if tar tzf "$fallback_archive" 2>/dev/null | grep -q 'private-path'; then
    record_failure 'upstream fetcher: the git route shipped an export-ignore path into the archive'
  fi
  if ! tar tzf "$fallback_archive" 2>/dev/null | grep -q 'VERSION'; then
    record_failure 'upstream fetcher: the git route omitted tracked, shipped content'
  fi
  fallback_top_level="$(tar tzf "$fallback_archive" 2>/dev/null | sed -n '1s|/.*||p')"
  if [ -z "$fallback_top_level" ] \
    || tar tzf "$fallback_archive" 2>/dev/null | grep -qv "^$fallback_top_level/"; then
    record_failure 'upstream fetcher: the git route did not produce a single top-level directory'
  fi

  # A canonical branch URL must select that exact ref. Default HEAD is retained only
  # when the URL does not match the branch grammar, and a missing requested ref fails.
  requested_branch_archive="$fetch_root/requested-branch.tar.gz"
  probe_output="$(bash "$archive_fetcher" "$requested_branch_archive" \
    "$rate_limited_url/archive/refs/heads/requested-branch.tar.gz" \
    "$upstream_fixture_repo" 2>&1)"
  probe_status=$?
  assert_status 0 'upstream fetcher: requested branch exits 0'
  if ! tar tzf "$requested_branch_archive" 2>/dev/null | grep -q 'requested-branch-marker\.txt'; then
    record_failure 'upstream fetcher: requested branch archive omitted its marker'
  fi
  if tar tzf "$requested_branch_archive" 2>/dev/null | grep -q 'default-branch-marker\.txt'; then
    record_failure 'upstream fetcher: requested branch archive substituted default HEAD'
  fi
  requested_branch_extract="$fetch_root/requested-branch-extract"
  mkdir -p "$requested_branch_extract"
  tar xzf "$requested_branch_archive" -C "$requested_branch_extract" --strip-components=1
  assert_file_contains "$requested_branch_extract/VERSION" 'shipped' \
    'upstream fetcher: ordinary branch survives one-component extraction at the root'
  assert_file_contains "$requested_branch_extract/suite/modules.tsv" 'skills/do-work' \
    'upstream fetcher: ordinary branch manifest survives one-component extraction at the root'
  assert_file_contains "$requested_branch_extract/requested-branch-marker.txt" 'requested branch' \
    'upstream fetcher: ordinary branch content survives one-component extraction at the root'

  slashed_branch_archive="$fetch_root/slashed-branch.tar.gz"
  probe_output="$(bash "$archive_fetcher" "$slashed_branch_archive" \
    "$rate_limited_url/archive/refs/heads/feature/fix.tar.gz" \
    "$upstream_fixture_repo" 2>&1)"
  probe_status=$?
  assert_status 0 'upstream fetcher: slashed branch exits 0'
  slashed_branch_extract="$fetch_root/slashed-branch-extract"
  mkdir -p "$slashed_branch_extract"
  tar xzf "$slashed_branch_archive" -C "$slashed_branch_extract" --strip-components=1
  assert_file_contains "$slashed_branch_extract/VERSION" 'shipped' \
    'upstream fetcher: slashed branch VERSION survives one-component extraction at the root'
  assert_file_contains "$slashed_branch_extract/suite/modules.tsv" 'skills/do-work' \
    'upstream fetcher: slashed branch manifest survives one-component extraction at the root'
  assert_file_contains "$slashed_branch_extract/slashed-branch-marker.txt" 'slashed branch' \
    'upstream fetcher: slashed branch content survives one-component extraction at the root'
  assert_path_absent "$slashed_branch_extract/fix" \
    'upstream fetcher: slashed branch does not leave a branch-name directory after stripping'

  missing_branch_archive="$fetch_root/missing-branch.tar.gz"
  probe_output="$(bash "$archive_fetcher" "$missing_branch_archive" \
    "$rate_limited_url/archive/refs/heads/missing-branch.tar.gz" \
    "$upstream_fixture_repo" 2>&1)"
  probe_status=$?
  assert_status_nonzero 'upstream fetcher: missing requested branch fails'
  assert_path_absent "$missing_branch_archive" \
    'upstream fetcher: missing requested branch publishes no archive'

  default_branch_archive="$fetch_root/default-branch.tar.gz"
  probe_output="$(bash "$archive_fetcher" "$default_branch_archive" \
    "$rate_limited_url/unparseable.tar.gz" "$upstream_fixture_repo" 2>&1)"
  probe_status=$?
  assert_status 0 'upstream fetcher: unparseable URL uses default HEAD'
  if ! tar tzf "$default_branch_archive" 2>/dev/null | grep -q 'default-branch-marker\.txt'; then
    record_failure 'upstream fetcher: unparseable URL archive omitted the default marker'
  fi
  if tar tzf "$default_branch_archive" 2>/dev/null | grep -q 'requested-branch-marker\.txt'; then
    record_failure 'upstream fetcher: unparseable URL selected a non-default branch'
  fi

  # A trapped interruption must terminate before the valid Git fallback can publish.
  bash_path="$(command -v bash)"
  for signal_case in HUP:129 INT:130 TERM:143; do
    signal_name="${signal_case%%:*}"
    expected_signal_status="${signal_case##*:}"
    interrupted_archive="$fetch_root/interrupted-$signal_name.tar.gz"
    probe_output="$(FETCH_SIGNAL="$signal_name" FETCHER_PATH="$archive_fetcher" \
      ARCHIVE_PATH="$interrupted_archive" UPSTREAM_REPO="$upstream_fixture_repo" HTTP_BASE="$rate_limited_url" \
      "$bash_path" -c '
        bash() {
          kill -s "$FETCH_SIGNAL" "$$"
          return 1
        }
        export -f bash
        exec "$1" "$FETCHER_PATH" "$ARCHIVE_PATH" \
          "$HTTP_BASE/archive/refs/heads/main.tar.gz" "$UPSTREAM_REPO"
      ' _ "$bash_path" 2>&1)"
    probe_status=$?
    assert_status "$expected_signal_status" \
      "upstream fetcher: $signal_name exits with its conventional status"
    assert_path_absent "$interrupted_archive" \
      "upstream fetcher: $signal_name publishes no archive"
    assert_output_lacks 'upstream archive fetched' \
      "upstream fetcher: $signal_name reports no success"
  done

  # Case 3: when every route fails, the pre-existing target survives untouched and no
  # private scratch is left behind.
  preserved_target="$fetch_root/preserved.tar.gz"
  printf 'previously downloaded archive\n' > "$preserved_target"
  bash "$archive_fetcher" "$preserved_target" \
    "$rate_limited_url/archive/refs/heads/main.tar.gz" "$fetch_root/no-such-repo" >/dev/null 2>&1 \
    && record_failure 'upstream fetcher: total failure reported success'
  if [ "$(cat "$preserved_target")" != 'previously downloaded archive' ]; then
    record_failure 'upstream fetcher: total failure changed the pre-existing target'
  fi
  if find "$fetch_root" -name 'preserved.tar.gz.download.*' -print -quit | grep -q .; then
    record_failure 'upstream fetcher: total failure leaked atomic-download scratch'
  fi
  if find "$fetch_root" -name 'preserved.tar.gz.fetching.*' -print -quit | grep -q .; then
    record_failure 'upstream fetcher: total failure leaked git-route scratch'
  fi
  total_failure_report="$(bash "$archive_fetcher" "$preserved_target" \
    "$rate_limited_url/archive/refs/heads/main.tar.gz" "$fetch_root/no-such-repo" 2>&1 >/dev/null || true)"
  case "$total_failure_report" in
    *'HTTP route'*'Git route'*) ;;
    *) record_failure 'upstream fetcher: total failure did not name both route outcomes' ;;
  esac
  case "$total_failure_report" in
    *DO_WORK_UPSTREAM_URL*) ;;
    *) record_failure 'upstream fetcher: total failure did not name the DO_WORK_UPSTREAM_URL escape hatch' ;;
  esac
  stop_archive_server
fi

# REQ-414 D-07 moved the install/update characterization seam to direct Go HTTP fetching.
# The install and update transactions moved into Go, so the two callers that used to shell
# out to the fetcher now call the archivefetch package. These assertions follow the behaviour
# to the file that owns it rather than being retired: the point is still that neither caller
# reaches for a bare curl, and that both honour the documented override.
install_transaction_source="$repo_root/skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go"
update_transaction_source="$repo_root/skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go"
archive_fetch_source="$repo_root/skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch.go"
for upstream_caller in "$install_transaction_source" "$update_transaction_source"; do
  if ! grep -q 'archivefetch\.FetchArchive' "$upstream_caller"; then
    record_failure "upstream fetcher: ${upstream_caller##*/} does not delegate its fetch to the shared fetch package"
  fi
  if grep -q '"curl"' "$upstream_caller"; then
    record_failure "upstream fetcher: ${upstream_caller##*/} reaches for curl directly"
  fi
done
if ! grep -q 'DO_WORK_UPSTREAM_URL' "$archive_fetch_source"; then
  record_failure 'upstream fetcher: archive_fetch.go does not honor DO_WORK_UPSTREAM_URL'
fi
if ! grep -q 'http\.NewRequestWithContext' "$archive_fetch_source" \
  || grep -Eq 'exec\.CommandContext\([^\n]*(curl|atomic-download)' "$archive_fetch_source"; then
  record_failure 'upstream fetcher: archive_fetch.go does not own HTTP transfer directly in Go'
fi
# The public shell entry points must still exist as launchers over the command.
for suite_launcher in \
  "$repo_root/skills/do-work/tools/do-work-update.sh" \
  "$repo_root/tools/install-do-work-suite.sh"; do
  if ! grep -q 'do-work-cli' "$suite_launcher"; then
    record_failure "upstream fetcher: ${suite_launcher##*/} no longer delegates to do-work-cli"
  fi
done

# The agent-facing path must delegate mutation to the same tested engine.
if ! grep -q 'tools/do-work-cli\.sh.*update-suite' "$repo_root/skills/do-work/actions/version.md"; then
  record_failure 'entry-point parity: modular actions/version.md does not delegate to the canonical update-suite command'
fi

if [ "$fail_count" -gt 0 ]; then
  printf 'update-script behavior probes: %s failure(s).\n' "$fail_count" >&2
  exit 1
fi

printf 'update-script behavior probes passed.\n'
