#!/usr/bin/env bash
# Behavioral contract for the four-module suite manifest validator.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
validator="$repo_root/tools/validate-suite-manifest.sh"
fixture_root="$(mktemp -d)"
fail_count=0

cleanup_fixture() {
  rm -rf "$fixture_root"
}
trap cleanup_fixture EXIT

write_manifest() {
  local archive_root="$1"
  shift
  mkdir -p "$archive_root/suite"
  printf 'source\tdestination\n' > "$archive_root/suite/modules.tsv"
  while [ "$#" -gt 0 ]; do
    printf '%s\t%s\n' "$1" "$2" >> "$archive_root/suite/modules.tsv"
    shift 2
  done
}

make_valid_fixture() {
  local archive_root="$1" module_name
  printf '0.184.0\n' > "$archive_root/VERSION"
  write_manifest "$archive_root" \
    skills/do-work .claude/skills/do-work \
    skills/do-work-board .claude/skills/do-work-board \
    skills/do-work-knowledge .claude/skills/do-work-knowledge \
    skills/do-work-toolbox .claude/skills/do-work-toolbox
  for module_name in do-work do-work-board do-work-knowledge do-work-toolbox; do
    mkdir -p "$archive_root/skills/$module_name"
    printf '# %s\n' "$module_name" > "$archive_root/skills/$module_name/SKILL.md"
  done
  mkdir -p "$archive_root/skills/do-work/actions"
  printf '0.184.0\n' > "$archive_root/skills/do-work/VERSION"
  printf '# Version Action\n\n**Current version**: 0.184.0\n' \
    > "$archive_root/skills/do-work/actions/version.md"
}

clone_valid_fixture() {
  local case_name="$1"
  local case_root="$fixture_root/$case_name"
  mkdir -p "$case_root"
  cp -R "$fixture_root/valid/." "$case_root/"
  printf '%s\n' "$case_root"
}

run_validator() {
  local archive_root="$1"
  validator_output="$(bash "$validator" --root "$archive_root" 2>&1)"
  validator_status=$?
}

expect_pass() {
  local archive_root="$1" case_name="$2"
  run_validator "$archive_root"
  if [ "$validator_status" -ne 0 ]; then
    printf 'FAIL: %s — expected pass, got %s. Output:\n%s\n' \
      "$case_name" "$validator_status" "$validator_output" >&2
    fail_count=$((fail_count + 1))
  fi
}

expect_fail() {
  local archive_root="$1" case_name="$2"
  run_validator "$archive_root"
  if [ "$validator_status" -eq 0 ]; then
    printf 'FAIL: %s — expected rejection. Output:\n%s\n' \
      "$case_name" "$validator_output" >&2
    fail_count=$((fail_count + 1))
  fi
}

if [ ! -f "$validator" ]; then
  printf 'FAIL: tools/validate-suite-manifest.sh is missing.\n' >&2
  exit 1
fi

mkdir -p "$fixture_root/valid"
make_valid_fixture "$fixture_root/valid"
valid_before="$(find "$fixture_root/valid" -type f -exec cksum {} \; | sort)"
expect_pass "$fixture_root/valid" 'canonical four-module fixture'
valid_after="$(find "$fixture_root/valid" -type f -exec cksum {} \; | sort)"
if [ "$valid_before" != "$valid_after" ]; then
  printf 'FAIL: validation modified the canonical fixture.\n' >&2
  fail_count=$((fail_count + 1))
fi

case_root="$(clone_valid_fixture unknown-columns)"
printf 'source\tdestination\textra\n' > "$case_root/suite/modules.tsv"
expect_fail "$case_root" 'unknown columns'

case_root="$(clone_valid_fixture absolute-source)"
sed 's#skills/do-work#/#' "$case_root/suite/modules.tsv" > "$case_root/suite/modules.tsv.next"
mv "$case_root/suite/modules.tsv.next" "$case_root/suite/modules.tsv"
expect_fail "$case_root" 'absolute source'

case_root="$(clone_valid_fixture traversal)"
sed 's#skills/do-work#skills/../do-work#' "$case_root/suite/modules.tsv" > "$case_root/suite/modules.tsv.next"
mv "$case_root/suite/modules.tsv.next" "$case_root/suite/modules.tsv"
expect_fail "$case_root" 'traversing source'

case_root="$(clone_valid_fixture duplicate-source)"
sed 's#skills/do-work-board#skills/do-work#' "$case_root/suite/modules.tsv" > "$case_root/suite/modules.tsv.next"
mv "$case_root/suite/modules.tsv.next" "$case_root/suite/modules.tsv"
expect_fail "$case_root" 'duplicate source'

case_root="$(clone_valid_fixture duplicate-destination)"
sed 's#\.claude/skills/do-work-board#.claude/skills/do-work#' "$case_root/suite/modules.tsv" > "$case_root/suite/modules.tsv.next"
mv "$case_root/suite/modules.tsv.next" "$case_root/suite/modules.tsv"
expect_fail "$case_root" 'duplicate destination'

case_root="$(clone_valid_fixture destination-escape)"
sed 's#\.claude/skills/do-work-board#skills/do-work-board#' "$case_root/suite/modules.tsv" > "$case_root/suite/modules.tsv.next"
mv "$case_root/suite/modules.tsv.next" "$case_root/suite/modules.tsv"
expect_fail "$case_root" 'destination outside .claude/skills'

case_root="$(clone_valid_fixture incomplete)"
sed '/do-work-toolbox/d' "$case_root/suite/modules.tsv" > "$case_root/suite/modules.tsv.next"
mv "$case_root/suite/modules.tsv.next" "$case_root/suite/modules.tsv"
expect_fail "$case_root" 'incomplete module set'

case_root="$(clone_valid_fixture unexpected-module)"
sed 's#do-work-toolbox#do-work-extra#g' "$case_root/suite/modules.tsv" > "$case_root/suite/modules.tsv.next"
mv "$case_root/suite/modules.tsv.next" "$case_root/suite/modules.tsv"
expect_fail "$case_root" 'unexpected module mapping'

case_root="$(clone_valid_fixture missing-skill)"
rm "$case_root/skills/do-work-board/SKILL.md"
expect_fail "$case_root" 'missing SKILL.md'

case_root="$(clone_valid_fixture empty-skill)"
: > "$case_root/skills/do-work-board/SKILL.md"
expect_fail "$case_root" 'empty SKILL.md'

case_root="$(clone_valid_fixture symlink-escape)"
mkdir -p "$fixture_root/outside-skill"
printf '# outside\n' > "$fixture_root/outside-skill/SKILL.md"
rm -rf "$case_root/skills/do-work-board"
ln -s "$fixture_root/outside-skill" "$case_root/skills/do-work-board"
expect_fail "$case_root" 'symlink source escape'

case_root="$(clone_valid_fixture invalid-version)"
printf 'v1\n' > "$case_root/VERSION"
expect_fail "$case_root" 'invalid suite version'

case_root="$(clone_valid_fixture mismatched-core-version)"
printf '0.184.1\n' > "$case_root/skills/do-work/VERSION"
expect_fail "$case_root" 'core VERSION differs from suite VERSION'

case_root="$(clone_valid_fixture mismatched-action-version)"
sed 's/0\.184\.0/9.9.9/' "$case_root/skills/do-work/actions/version.md" \
  > "$case_root/skills/do-work/actions/version.md.next"
mv "$case_root/skills/do-work/actions/version.md.next" \
  "$case_root/skills/do-work/actions/version.md"
expect_fail "$case_root" 'runtime action version differs from suite VERSION'

if [ "$fail_count" -gt 0 ]; then
  exit 1
fi

printf 'Suite manifest contract probes passed.\n'
