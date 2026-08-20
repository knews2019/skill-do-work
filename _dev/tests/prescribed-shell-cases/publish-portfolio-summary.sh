#!/usr/bin/env bash
# Fixture execution proofs for publish-portfolio-summary.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

portfolio_root="$fixture_root/portfolio"
mkdir -p "$portfolio_root/deliverables/portfolio-snapshots"
portfolio_source="$portfolio_root/retained-summary.md"
portfolio_canonical="$portfolio_root/deliverables/portfolio-summary.md"
portfolio_candidate="$portfolio_root/deliverables/portfolio-snapshots/portfolio-summary-20260815T120000Z.md"
printf 'new portfolio bytes\n' > "$portfolio_source"

# publish-portfolio-summary: canonical-only publication atomically refreshes only
# the canonical path and retains the source used for publication.
printf 'old canonical\n' > "$portfolio_canonical"
portfolio_canonical_output="$($toolbox_scripts/publish-portfolio-summary.sh --canonical-only "$portfolio_source" "$portfolio_canonical" 2>/dev/null)" \
  || fail_case 'publish-portfolio-summary canonical-only case returned nonzero'
cmp -s "$portfolio_source" "$portfolio_canonical" \
  || fail_case 'publish-portfolio-summary canonical-only case changed the retained bytes'
[ -f "$portfolio_source" ] \
  || fail_case 'publish-portfolio-summary canonical-only case consumed the retained source'
[ "$portfolio_canonical_output" = "$portfolio_canonical" ] \
  || fail_case 'publish-portfolio-summary canonical-only case did not report the canonical path'
find "$portfolio_root/deliverables/portfolio-snapshots" -type f -print -quit | grep -q . \
  && fail_case 'publish-portfolio-summary canonical-only case created a snapshot'

# publish-portfolio-summary: the preservation branch publishes a snapshot first, then
# refreshes canonical from the same verified bytes — but as an independent file. Same
# bytes is the requirement; same inode is the defect, because a shared inode makes the
# immutable snapshot follow every later in-place edit of canonical (REQ-205, was asserted
# the other way round by REQ-199 before durable immutability was tested).
printf 'prior canonical\n' > "$portfolio_canonical"
portfolio_snapshot_output="$($toolbox_scripts/publish-portfolio-summary.sh --with-snapshot "$portfolio_source" "$portfolio_canonical" "$portfolio_candidate" 2>/dev/null)" \
  || fail_case 'publish-portfolio-summary snapshot-success case returned nonzero'
[ "$portfolio_snapshot_output" = "$(printf '%s\n%s' "$portfolio_canonical" "$portfolio_candidate")" ] \
  || fail_case 'publish-portfolio-summary snapshot-success case did not report both published paths'
cmp -s "$portfolio_source" "$portfolio_candidate" && cmp -s "$portfolio_source" "$portfolio_canonical" \
  || fail_case 'publish-portfolio-summary snapshot-success case did not preserve byte identity'
[ -f "$portfolio_candidate" ] && [ -f "$portfolio_canonical" ] \
  || fail_case 'publish-portfolio-summary snapshot-success case did not publish two regular files'
[ "$portfolio_candidate" -ef "$portfolio_canonical" ] \
  && fail_case 'publish-portfolio-summary snapshot-success case aliased the snapshot to the canonical inode'
printf 'canonical mutated after publication\n' > "$portfolio_canonical"
[ "$(cat "$portfolio_candidate")" = 'new portfolio bytes' ] \
  || fail_case 'publish-portfolio-summary snapshot-success case let a later canonical edit rewrite the snapshot'
cp "$portfolio_source" "$portfolio_canonical"

# publish-portfolio-summary: an occupied candidate remains immutable and advances
# to the first numeric suffix without cleaning any prior snapshot.
portfolio_collision_candidate="$portfolio_root/deliverables/portfolio-snapshots/portfolio-summary-20260815T130000Z.md"
portfolio_collision_suffix="$portfolio_root/deliverables/portfolio-snapshots/portfolio-summary-20260815T130000Z-2.md"
printf 'occupied collision\n' > "$portfolio_collision_candidate"
portfolio_collision_output="$($toolbox_scripts/publish-portfolio-summary.sh --with-snapshot "$portfolio_source" "$portfolio_canonical" "$portfolio_collision_candidate" 2>/dev/null)" \
  || fail_case 'publish-portfolio-summary collision case returned nonzero'
[ "$(cat "$portfolio_collision_candidate")" = 'occupied collision' ] \
  || fail_case 'publish-portfolio-summary collision case changed the occupant'
cmp -s "$portfolio_source" "$portfolio_collision_suffix" \
  || fail_case 'publish-portfolio-summary collision case did not publish the numeric suffix'
[ "$portfolio_collision_output" = "$(printf '%s\n%s' "$portfolio_canonical" "$portfolio_collision_suffix")" ] \
  || fail_case 'publish-portfolio-summary collision case reported the wrong suffix'
[ "$(cat "$portfolio_candidate")" = 'new portfolio bytes' ] \
  || fail_case 'publish-portfolio-summary collision case cleaned an unrelated prior snapshot'

# publish-portfolio-summary: exclusive snapshot failure leaves the prior canonical
# unchanged and never leaks the private verified copy.
portfolio_ln_failure_bin="$fixture_root/portfolio-ln-failure-bin"
mkdir -p "$portfolio_ln_failure_bin"
printf '%s\n' '#!/usr/bin/env bash' 'exit 1' > "$portfolio_ln_failure_bin/ln"
chmod +x "$portfolio_ln_failure_bin/ln"
portfolio_failure_candidate="$portfolio_root/deliverables/portfolio-snapshots/portfolio-summary-20260815T140000Z.md"
printf 'stable before snapshot failure\n' > "$portfolio_canonical"
PATH="$portfolio_ln_failure_bin:$PATH" \
  "$toolbox_scripts/publish-portfolio-summary.sh" --with-snapshot "$portfolio_source" "$portfolio_canonical" "$portfolio_failure_candidate" >/dev/null 2>&1 \
  && fail_case 'publish-portfolio-summary snapshot-failure case returned success'
[ "$(cat "$portfolio_canonical")" = 'stable before snapshot failure' ] \
  || fail_case 'publish-portfolio-summary snapshot-failure case changed the prior canonical'
[ ! -e "$portfolio_failure_candidate" ] \
  || fail_case 'publish-portfolio-summary snapshot-failure case left a published snapshot'
find "$portfolio_root/deliverables" -name '.portfolio-summary.md.publishing.*' -print -quit | grep -q . \
  && fail_case 'publish-portfolio-summary snapshot-failure case leaked private bytes'

# publish-portfolio-summary: a later canonical replacement failure retains the
# already-published snapshot while preserving the prior canonical.
portfolio_mv_failure_bin="$fixture_root/portfolio-mv-failure-bin"
mkdir -p "$portfolio_mv_failure_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'if [ "$2" = "$PORTFOLIO_FAIL_CANONICAL" ]; then exit 1; fi' \
  'exec "$PORTFOLIO_REAL_MV" "$@"' \
  > "$portfolio_mv_failure_bin/mv"
chmod +x "$portfolio_mv_failure_bin/mv"
portfolio_late_failure_candidate="$portfolio_root/deliverables/portfolio-snapshots/portfolio-summary-20260815T150000Z.md"
printf 'stable before canonical failure\n' > "$portfolio_canonical"
PORTFOLIO_FAIL_CANONICAL="$portfolio_canonical" \
PORTFOLIO_REAL_MV="$(command -v mv)" \
PATH="$portfolio_mv_failure_bin:$PATH" \
  "$toolbox_scripts/publish-portfolio-summary.sh" --with-snapshot "$portfolio_source" "$portfolio_canonical" "$portfolio_late_failure_candidate" >/dev/null 2>&1 \
  && fail_case 'publish-portfolio-summary canonical-failure case returned success'
[ "$(cat "$portfolio_canonical")" = 'stable before canonical failure' ] \
  || fail_case 'publish-portfolio-summary canonical-failure case changed the prior canonical'
cmp -s "$portfolio_source" "$portfolio_late_failure_candidate" \
  || fail_case 'publish-portfolio-summary canonical-failure case did not retain the published snapshot'
find "$portfolio_root/deliverables" -name '.portfolio-summary.md.publishing.*' -print -quit | grep -q . \
  && fail_case 'publish-portfolio-summary canonical-failure case leaked private bytes'

# publish-portfolio-summary: `ln` links *into* a directory operand instead of colliding
# with it, so a snapshot candidate occupied by a directory must advance to the numeric
# suffix and leave no private file nested inside the directory.
portfolio_directory_candidate="$portfolio_root/deliverables/portfolio-snapshots/portfolio-summary-20260815T160000Z.md"
portfolio_directory_suffix="$portfolio_root/deliverables/portfolio-snapshots/portfolio-summary-20260815T160000Z-2.md"
mkdir -p "$portfolio_directory_candidate"
printf 'occupant\n' > "$portfolio_directory_candidate/occupant.txt"
printf 'stable before directory candidate\n' > "$portfolio_canonical"
portfolio_directory_output="$($toolbox_scripts/publish-portfolio-summary.sh --with-snapshot "$portfolio_source" "$portfolio_canonical" "$portfolio_directory_candidate" 2>/dev/null)" \
  || fail_case 'publish-portfolio-summary snapshot-directory case returned nonzero instead of advancing'
[ -d "$portfolio_directory_candidate" ] && [ "$(ls -A "$portfolio_directory_candidate")" = occupant.txt ] \
  || fail_case 'publish-portfolio-summary snapshot-directory case nested a private file inside the occupying directory'
cmp -s "$portfolio_source" "$portfolio_directory_suffix" \
  || fail_case 'publish-portfolio-summary snapshot-directory case did not advance to the numeric suffix'
[ "$portfolio_directory_output" = "$(printf '%s\n%s' "$portfolio_canonical" "$portfolio_directory_suffix")" ] \
  || fail_case 'publish-portfolio-summary snapshot-directory case reported the wrong published paths'

# publish-portfolio-summary: `mv` moves *into* a directory operand, so a canonical path
# occupied by a directory must fail closed — never advance, never publish inside it, and
# never leave the private copy nested there.
portfolio_directory_canonical="$portfolio_root/deliverables/canonical-as-directory.md"
mkdir -p "$portfolio_directory_canonical"
printf 'canonical occupant\n' > "$portfolio_directory_canonical/occupant.txt"
"$toolbox_scripts/publish-portfolio-summary.sh" --canonical-only "$portfolio_source" "$portfolio_directory_canonical" >/dev/null 2>&1 \
  && fail_case 'publish-portfolio-summary canonical-directory case reported success'
[ -d "$portfolio_directory_canonical" ] && [ "$(ls -A "$portfolio_directory_canonical")" = occupant.txt ] \
  || fail_case 'publish-portfolio-summary canonical-directory case did not leave the occupying directory unchanged'
find "$portfolio_root/deliverables" -name '.canonical-as-directory.md.publishing.*' -print -quit | grep -q . \
  && fail_case 'publish-portfolio-summary canonical-directory case leaked private bytes'

prescribed_shell_finish
