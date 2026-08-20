#!/usr/bin/env bash
# Fixture execution proofs for install-last30days.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

# install-last30days: a SKILL.md-only tree fails check, is repaired from a
# complete fixture, and receives the full subtree plus ignore/Python guarantees.
upstream_repo="$fixture_root/last30days-upstream"
fixture_repo_init "$upstream_repo"
mkdir -p "$upstream_repo/skills/last30days/scripts" "$upstream_repo/skills/last30days/support"
printf '# Last30Days\n' > "$upstream_repo/skills/last30days/SKILL.md"
printf 'runtime\n' > "$upstream_repo/skills/last30days/scripts/last30days.py"
printf 'support\n' > "$upstream_repo/skills/last30days/support/data.txt"
fixture_repo_commit_all "$upstream_repo" fixture
last_project="$fixture_root/last-project"
fixture_repo_init "$last_project"
mkdir -p "$last_project/.claude/skills/last30days"
printf '# Sentinel only\n' > "$last_project/.claude/skills/last30days/SKILL.md"
python_bin="$fixture_root/python-bin"
mkdir -p "$python_bin"
printf '%s\n' '#!/usr/bin/env bash' 'exit 0' > "$python_bin/python3.12"
chmod +x "$python_bin/python3.12"
PATH="$python_bin:$PATH" "$toolbox_scripts/install-last30days.sh" check "$last_project" >/dev/null 2>&1 \
  && fail_case 'install-last30days sentinel-only check accepted a missing runtime script'
last_output="$(TMPDIR="$fixture_root" PATH="$python_bin:$PATH" "$toolbox_scripts/install-last30days.sh" install "$last_project" "$upstream_repo" 2>&1)" \
  || fail_case "install-last30days complete-source repair returned nonzero: $last_output"
[ "$(cat "$last_project/.claude/skills/last30days/scripts/last30days.py")" = runtime ] \
  && [ "$(cat "$last_project/.claude/skills/last30days/support/data.txt")" = support ] \
  || fail_case 'install-last30days complete-source repair omitted part of the source subtree'
git -C "$last_project" check-ignore -q .claude/skills/last30days/SKILL.md || fail_case 'install-last30days fixture-source case did not resolve the sibling exclude helper'
PATH="$python_bin:$PATH" "$toolbox_scripts/install-last30days.sh" check "$last_project" >/dev/null 2>&1 \
  || fail_case 'install-last30days repaired tree failed the complete runtime/ignore/Python check'

# install-last30days: reject an incomplete source before publishing any final tree.
incomplete_upstream_repo="$fixture_root/last30days-incomplete-upstream"
fixture_repo_init "$incomplete_upstream_repo"
mkdir -p "$incomplete_upstream_repo/skills/last30days"
printf '# Incomplete\n' > "$incomplete_upstream_repo/skills/last30days/SKILL.md"
fixture_repo_commit_all "$incomplete_upstream_repo" fixture
incomplete_source_project="$fixture_root/last-incomplete-source-project"
mkdir -p "$incomplete_source_project"
TMPDIR="$fixture_root" PATH="$python_bin:$PATH" "$toolbox_scripts/install-last30days.sh" install "$incomplete_source_project" "$incomplete_upstream_repo" >/dev/null 2>&1 \
  && fail_case 'install-last30days incomplete-source case returned success'
[ ! -e "$incomplete_source_project/.claude/skills/last30days" ] \
  || fail_case 'install-last30days incomplete-source case published a final tree'

# install-last30days: a partial copy failure leaves no final tree or private path.
copy_failure_bin="$fixture_root/copy-failure-bin"
mkdir -p "$copy_failure_bin"
printf '%s\n' '#!/usr/bin/env bash' 'while [ "$#" -gt 1 ]; do shift; done' 'copy_destination="$1"' 'mkdir -p "$copy_destination/scripts"' 'printf partial > "$copy_destination/SKILL.md"' 'exit 1' > "$copy_failure_bin/cp"
chmod +x "$copy_failure_bin/cp"
copy_failure_project="$fixture_root/last-copy-failure-project"
mkdir -p "$copy_failure_project"
TMPDIR="$fixture_root" PATH="$copy_failure_bin:$python_bin:$PATH" "$toolbox_scripts/install-last30days.sh" install "$copy_failure_project" "$upstream_repo" >/dev/null 2>&1 \
  && fail_case 'install-last30days copy-failure case returned success'
[ ! -e "$copy_failure_project/.claude/skills/last30days" ] \
  || fail_case 'install-last30days copy-failure case left a final tree'
find "$copy_failure_project/.claude/skills" -name '.last30days.*' -print -quit | grep -q . \
  && fail_case 'install-last30days copy-failure case leaked private paths'

# install-last30days: a publication failure with no prior destination removes any
# simulated partial final tree.
publication_failure_bin="$fixture_root/publication-failure-bin"
mkdir -p "$publication_failure_bin"
printf '%s\n' '#!/usr/bin/env bash' 'publication_source="$1"' 'publication_destination="$2"' 'case "$publication_source:$publication_destination" in *.last30days.staging.*:*/.claude/skills/last30days) mkdir -p "$publication_destination"; printf partial > "$publication_destination/SKILL.md"; exit 1 ;; esac' 'exec "$LAST30DAYS_REAL_MV" "$@"' > "$publication_failure_bin/mv"
chmod +x "$publication_failure_bin/mv"
publication_failure_project="$fixture_root/last-publication-failure-project"
mkdir -p "$publication_failure_project"
TMPDIR="$fixture_root" LAST30DAYS_REAL_MV="$(command -v mv)" PATH="$publication_failure_bin:$python_bin:$PATH" \
  "$toolbox_scripts/install-last30days.sh" install "$publication_failure_project" "$upstream_repo" >/dev/null 2>&1 \
  && fail_case 'install-last30days publication-failure case returned success'
[ ! -e "$publication_failure_project/.claude/skills/last30days" ] \
  || fail_case 'install-last30days publication-failure case left a new final tree'
find "$publication_failure_project/.claude/skills" -name '.last30days.*' -print -quit | grep -q . \
  && fail_case 'install-last30days publication-failure case leaked private paths'

# install-last30days: a replacement publication failure restores an existing
# incomplete tree byte-for-byte and removes private staging/backup paths.
replacement_failure_project="$fixture_root/last-replacement-failure-project"
mkdir -p "$replacement_failure_project/.claude/skills/last30days/legacy"
printf original-sentinel > "$replacement_failure_project/.claude/skills/last30days/SKILL.md"
printf original-byte > "$replacement_failure_project/.claude/skills/last30days/legacy/data.bin"
cp -R "$replacement_failure_project/.claude/skills/last30days" "$fixture_root/last30days-original-snapshot"
TMPDIR="$fixture_root" LAST30DAYS_REAL_MV="$(command -v mv)" PATH="$publication_failure_bin:$python_bin:$PATH" \
  "$toolbox_scripts/install-last30days.sh" install "$replacement_failure_project" "$upstream_repo" >/dev/null 2>&1 \
  && fail_case 'install-last30days replacement-failure case returned success'
diff -r "$fixture_root/last30days-original-snapshot" "$replacement_failure_project/.claude/skills/last30days" >/dev/null 2>&1 \
  || fail_case 'install-last30days replacement-failure case did not restore the prior tree byte-for-byte'
find "$replacement_failure_project/.claude/skills" -name '.last30days.*' -print -quit | grep -q . \
  && fail_case 'install-last30days replacement-failure case leaked private paths'

# install-last30days: interruption after the prior tree moves to backup restores it
# and removes both staging and backup paths.
backup_interrupt_bin="$fixture_root/backup-interrupt-bin"
mkdir -p "$backup_interrupt_bin"
printf '%s\n' '#!/usr/bin/env bash' 'interrupt_destination="$2"' 'case "$interrupt_destination" in */.last30days.backup.*/previous) "$LAST30DAYS_REAL_MV" "$@" || exit $?; kill -TERM "$PPID"; exit 0 ;; esac' 'exec "$LAST30DAYS_REAL_MV" "$@"' > "$backup_interrupt_bin/mv"
chmod +x "$backup_interrupt_bin/mv"
backup_interrupt_project="$fixture_root/last-backup-interrupt-project"
mkdir -p "$backup_interrupt_project/.claude/skills/last30days/legacy"
printf interrupt-sentinel > "$backup_interrupt_project/.claude/skills/last30days/SKILL.md"
printf interrupt-byte > "$backup_interrupt_project/.claude/skills/last30days/legacy/data.bin"
cp -R "$backup_interrupt_project/.claude/skills/last30days" "$fixture_root/last30days-interrupt-snapshot"
TMPDIR="$fixture_root" LAST30DAYS_REAL_MV="$(command -v mv)" PATH="$backup_interrupt_bin:$python_bin:$PATH" \
  "$toolbox_scripts/install-last30days.sh" install "$backup_interrupt_project" "$upstream_repo" >/dev/null 2>&1
backup_interrupt_status=$?
[ "$backup_interrupt_status" -eq 143 ] \
  || fail_case "install-last30days backup-interruption case returned $backup_interrupt_status instead of 143"
diff -r "$fixture_root/last30days-interrupt-snapshot" "$backup_interrupt_project/.claude/skills/last30days" >/dev/null 2>&1 \
  || fail_case 'install-last30days backup-interruption case did not restore the prior tree byte-for-byte'
find "$backup_interrupt_project/.claude/skills" -name '.last30days.*' -print -quit | grep -q . \
  && fail_case 'install-last30days backup-interruption case leaked private paths'

# install-last30days: the target can reappear between the backup `mv` and the
# publication `mv`, and `mv` then nests the staging tree inside it and exits zero. The
# mv shim below recreates the target in exactly that window; publication must fail
# closed, leave the reappeared tree byte-for-byte, and keep the prior tree recoverable.
publication_collision_bin="$fixture_root/publication-collision-bin"
mkdir -p "$publication_collision_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'case "${2:-}" in */.claude/skills/last30days) mkdir -p "$2"; printf owned-by-someone-else > "$2/keep.txt" ;; esac' \
  'exec /bin/mv "$@"' \
  > "$publication_collision_bin/mv"
chmod +x "$publication_collision_bin/mv"
publication_collision_project="$fixture_root/last-publication-collision-project"
publication_collision_target="$publication_collision_project/.claude/skills/last30days"
mkdir -p "$publication_collision_target/legacy"
printf collision-sentinel > "$publication_collision_target/SKILL.md"
printf collision-byte > "$publication_collision_target/legacy/data.bin"
TMPDIR="$fixture_root" PATH="$publication_collision_bin:$python_bin:$PATH" \
  "$toolbox_scripts/install-last30days.sh" install "$publication_collision_project" "$upstream_repo" >/dev/null 2>&1 \
  && fail_case 'install-last30days publication-collision case returned success'
[ "$(cat "$publication_collision_target/keep.txt" 2>/dev/null)" = owned-by-someone-else ] \
  || fail_case 'install-last30days publication-collision case did not preserve the reappeared target byte-for-byte'
[ "$(ls -A "$publication_collision_target" 2>/dev/null)" = keep.txt ] \
  || fail_case 'install-last30days publication-collision case left its staging tree nested inside the reappeared target'
find "$publication_collision_project/.claude/skills" -name '.last30days.staging.*' -print -quit | grep -q . \
  && fail_case 'install-last30days publication-collision case leaked private staging'
[ -s "$(find "$publication_collision_project/.claude/skills" -path '*/.last30days.backup.*/previous/SKILL.md' -print -quit)" ] \
  || fail_case 'install-last30days publication-collision case did not leave the prior tree recoverable at its backup path'

prescribed_shell_finish
