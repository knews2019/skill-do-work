#!/usr/bin/env bash
# Behavioral regression probes for the retained core SessionStart launcher and Go authority.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
fixture_root="$(mktemp -d)" || exit 1
trap 'rm -rf "$fixture_root"' EXIT
failure_count=0

# One physical skill root serves every scenario that needs the module. The launcher
# resolves its module directory physically and then asks the Go toolchain for the tool
# binary, so a per-case copy of the module is a per-case absolute path, a distinct Go
# build-cache action id, and a fresh link of a byte-identical binary. Sharing the root
# reaches one physical module by construction. It is safe because the hook, the launcher
# and the module are immutable in every scenario below: that the launcher writes nothing
# into the module tree is pinned separately by do-work-cli-launcher-behavior.sh, and
# assert_shared_skill_root_unchanged re-checks it here after every case.
shared_skill_root="$fixture_root/shared-skill"
shared_skill_digest=''

install_shared_skill_root() {
  mkdir -p "$shared_skill_root/hooks" "$shared_skill_root/actions" "$shared_skill_root/tools"
  cp "$repo_root/skills/do-work/hooks/session-start.sh" "$shared_skill_root/hooks/"
  cp "$repo_root/skills/do-work/tools/do-work-cli.sh" "$shared_skill_root/tools/"
  cp -R "$repo_root/skills/do-work/tools/do-work-cli" "$shared_skill_root/tools/"
  rm -f "$shared_skill_root/tools/do-work-cli/do-work-cli"
}

# The shared hook and module trees, content-addressed. actions/ is excluded because
# version.md is the per-case banner input every case rewrites.
shared_skill_immutable_digest() {
  find "$shared_skill_root/hooks" "$shared_skill_root/tools" -type f -exec shasum -a 256 {} + |
    LC_ALL=C sort | shasum -a 256
}

# The invariant that makes sharing safe. A case that writes into the shared tree fails
# here, naming itself, instead of silently corrupting whichever case runs next.
assert_shared_skill_root_unchanged() {
  local case_name="$1" current_digest
  current_digest="$(shared_skill_immutable_digest)"
  if [ "$current_digest" != "$shared_skill_digest" ]; then
    printf 'FAIL: %s modified the shared skill tree; the scenarios can no longer share it.\n' "$case_name" >&2
    failure_count=$((failure_count + 1))
  fi
}

# Every case helper takes its banner input as a required argument. With a shared root a
# case that forgot to write version.md would silently inherit the previous case's, which
# reads as a pass; a required positional makes the omission an unbound-variable error.
install_banner_input() {
  local version_text="$1"
  printf '%b' "$version_text" > "$shared_skill_root/actions/version.md"
}

install_shared_skill_root
shared_skill_digest="$(shared_skill_immutable_digest)"

run_banner_case() {
  local case_name="$1" version_text="$2" queue_count="$3" expected="$4"
  local case_root="$fixture_root/$case_name" request_index=0 output='' status=0
  mkdir -p "$case_root/project/do-work/queue"
  install_banner_input "$version_text"
  while [ "$request_index" -lt "$queue_count" ]; do
    request_index=$((request_index + 1))
    printf -- '---\nid: REQ-%03d\n---\n' "$request_index" > "$case_root/project/do-work/queue/REQ-$(printf '%03d' "$request_index")-fixture.md"
  done
  output="$(CLAUDE_PROJECT_DIR="$case_root/project" bash "$shared_skill_root/hooks/session-start.sh" 2>"$case_root/stderr")" || status=$?
  if [ "$status" -ne 0 ] || [ "$output" != "$expected" ] || [ -s "$case_root/stderr" ]; then
    printf 'FAIL: %s status=%s output=<%s> stderr=<%s>\n' "$case_name" "$status" "$output" "$(tr '\n' ' ' < "$case_root/stderr")" >&2
    failure_count=$((failure_count + 1))
  fi
  assert_shared_skill_root_unchanged "$case_name"
}

run_banner_case valid '**Current version**: 9.8.7\n' 2 "do-work v9.8.7 loaded. 2 pending REQ(s). Say 'do-work help' for commands."
run_banner_case missing '' 0 "do-work vunknown loaded. 0 pending REQ(s). Say 'do-work help' for commands."
run_banner_case reformatted '**Version**: 9.8.7\n' 0 "do-work vunknown loaded. 0 pending REQ(s). Say 'do-work help' for commands."
run_banner_case multiple '**Current version**: 1\n**Current version**: 2\n' 0 "do-work v1
2 loaded. 0 pending REQ(s). Say 'do-work help' for commands."

run_tail_case() {
  local case_name="$1" version_text="$2" request_id="$3" tail_kind="$4"
  local case_root="$fixture_root/$case_name" journal_root='' output='' expected='' before='' after='' status=0
  mkdir -p "$case_root/project/do-work/queue"
  install_banner_input "$version_text"
  git -C "$case_root/project" init -q
  if [ "$tail_kind" = journal ]; then
    journal_root="$(git -C "$case_root/project" rev-parse --git-path do-work-finalization)"
    case "$journal_root" in
      /*) ;;
      *) journal_root="$case_root/project/$journal_root" ;;
    esac
    mkdir -p "$journal_root"
    printf '{"phase":"release_applied","manifest":{"request_id":"%s"}}\n' "$request_id" > "$journal_root/$request_id.json"
    before="$(shasum -a 256 "$journal_root/$request_id.json" | cut -d' ' -f1)"
  else
    mkdir -p "$case_root/project/do-work/archive"
    printf -- '---\nid: %s\nstatus: completed\ncreated_at: 2026-09-01T00:00:00Z\ncommit:\n---\n## Implementation Summary\n- `src/x.go` (modified)\n\n## Qualification\nVerified\n' "$request_id" > "$case_root/project/do-work/archive/$request_id-tail.md"
    before="$(shasum -a 256 "$case_root/project/do-work/archive/$request_id-tail.md" | cut -d' ' -f1)"
  fi
  output="$(CLAUDE_PROJECT_DIR="$case_root/project" bash "$shared_skill_root/hooks/session-start.sh" 2>"$case_root/stderr")" || status=$?
  expected="do-work v9.8.7 loaded. 0 pending REQ(s). Say 'do-work help' for commands.
do-work: unfinished finalization for $request_id — 'do-work run' resumes it; 'do-work run-with-recovery' if this checkout is the only writer."
  if [ "$tail_kind" = journal ]; then
    after="$(shasum -a 256 "$journal_root/$request_id.json" | cut -d' ' -f1)"
  else
    after="$(shasum -a 256 "$case_root/project/do-work/archive/$request_id-tail.md" | cut -d' ' -f1)"
  fi
  if [ "$status" -ne 0 ] || [ "$output" != "$expected" ] || [ -s "$case_root/stderr" ] || [ "$before" != "$after" ]; then
    printf 'FAIL: %s status=%s output=<%s> stderr=<%s> before=%s after=%s\n' "$case_name" "$status" "$output" "$(tr '\n' ' ' < "$case_root/stderr")" "$before" "$after" >&2
    failure_count=$((failure_count + 1))
  fi
  assert_shared_skill_root_unchanged "$case_name"
}

run_tail_case unfinished-journal '**Current version**: 9.8.7\n' REQ-710 journal
run_tail_case archived-without-commit '**Current version**: 9.8.7\n' REQ-711 archive

# Retained-vs-Go housekeeping matrix: the banner stays exact while JSON carries
# the same deletion evidence, and unavailable Git preserves coordination state.
authority_root="$fixture_root/authority"
mkdir -p "$authority_root/project/do-work/working" "$authority_root/project/do-work/.req-reservations"
install_banner_input '**Current version**: 9.8.7\n'
git -C "$authority_root/project" init -q
git -C "$authority_root/project" config user.email fixture@example.invalid
git -C "$authority_root/project" config user.name Fixture
printf -- '---\nid: REQ-203\nstatus: claimed\n---\n' > "$authority_root/project/do-work/working/REQ-203-landed.md"
git -C "$authority_root/project" add do-work/working/REQ-203-landed.md
git -C "$authority_root/project" commit -qm 'land request'
touch "$authority_root/project/do-work/.req-reservations/REQ-000203"
authority_json="$(CLAUDE_PROJECT_DIR="$authority_root/project" bash "$shared_skill_root/tools/do-work-cli.sh" --repo-root "$authority_root/project" --format json session-start --skill-root "$shared_skill_root")"
if [ -e "$authority_root/project/do-work/.req-reservations/REQ-000203" ] || ! grep -q '"kind": "deleted"' <<<"$authority_json" || ! grep -q '"protocol_output": "do-work v9.8.7 loaded.' <<<"$authority_json"; then
  printf 'FAIL: committed housekeeping did not retain typed deletion plus exact protocol output.\n' >&2
  failure_count=$((failure_count + 1))
fi
printf -- '---\nid: REQ-204\nstatus: claimed\n---\n' > "$authority_root/project/do-work/working/REQ-204-uncommitted.md"
touch "$authority_root/project/do-work/.req-reservations/REQ-000204"
authority_binary="$(go tool -C "$shared_skill_root/tools/do-work-cli" -n do-work-cli)"
unavailable_json="$(PATH="$authority_root/no-git" "$authority_binary" --repo-root "$authority_root/project" --format json session-start --skill-root "$shared_skill_root")"
if [ ! -e "$authority_root/project/do-work/.req-reservations/REQ-000204" ] || ! grep -q 'RESERVATION-GIT-AUTHORITY-UNAVAILABLE' <<<"$unavailable_json" || grep -q 'REQ-000204.*deleted' <<<"$unavailable_json"; then
  printf 'FAIL: Git-unavailable SessionStart did not fail closed with typed evidence.\n' >&2
  failure_count=$((failure_count + 1))
fi
assert_shared_skill_root_unchanged authority-matrix

# Launcher failure matrix: retained SessionStart propagates both too-old Go and
# build failure with actionable diagnostics and no false banner.
run_launcher_failure_case() {
  local case_name="$1" version_text="$2" go_body="$3" expected="$4"
  local case_root="$fixture_root/$case_name" fake_bin="$fixture_root/$case_name-bin" status=0
  mkdir -p "$case_root/project" "$fake_bin"
  install_banner_input "$version_text"
  printf '%b' "$go_body" > "$fake_bin/go"
  chmod +x "$fake_bin/go"
  PATH="$fake_bin:/usr/bin:/bin" CLAUDE_PROJECT_DIR="$case_root/project" bash "$shared_skill_root/hooks/session-start.sh" >"$case_root/stdout" 2>"$case_root/stderr" || status=$?
  if [ "$status" -ne 2 ] || [ -s "$case_root/stdout" ] || ! grep -q "$expected" "$case_root/stderr"; then
    printf 'FAIL: %s launcher boundary status=%s stderr=<%s>.\n' "$case_name" "$status" "$(tr '\n' ' ' < "$case_root/stderr")" >&2
    failure_count=$((failure_count + 1))
  fi
  assert_shared_skill_root_unchanged "$case_name"
}
run_launcher_failure_case go-too-old '**Current version**: 9.8.7\n' '#!/usr/bin/env bash\nprintf "go version go1.24.0 fixture\\n"\n' 'Go 1.25.0 or newer is required'
run_launcher_failure_case go-build-failure '**Current version**: 9.8.7\n' '#!/usr/bin/env bash\nif [ "${1:-}" = version ]; then printf "go version go1.25.0 fixture\\n"; exit 0; fi\nexit 1\n' 'could not build the command; nothing was run'

missing_root="$fixture_root/missing-launcher"
mkdir -p "$missing_root/skill/hooks" "$missing_root/project"
cp "$repo_root/skills/do-work/hooks/session-start.sh" "$missing_root/skill/hooks/"
missing_status=0
CLAUDE_PROJECT_DIR="$missing_root/project" bash "$missing_root/skill/hooks/session-start.sh" >"$missing_root/stdout" 2>"$missing_root/stderr" || missing_status=$?
if [ "$missing_status" -eq 0 ] || ! grep -q 'canonical launcher is missing' "$missing_root/stderr"; then
  printf 'FAIL: core SessionStart did not stop actionably when the canonical launcher was absent.\n' >&2
  failure_count=$((failure_count + 1))
fi

if grep -Eq 'cleanup-req-reservations\.sh|repair-req-timestamps\.sh|sed |awk |find ' "$repo_root/skills/do-work/hooks/session-start.sh"; then
  printf 'FAIL: retained SessionStart path contains domain logic instead of a thin launcher.\n' >&2
  failure_count=$((failure_count + 1))
fi

[ "$failure_count" -eq 0 ] || exit 1
printf 'SessionStart hook behavior probes passed.\n'
