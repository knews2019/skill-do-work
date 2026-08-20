#!/usr/bin/env bash
# Fixture execution proofs for audit-archive-timestamps.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

# audit-archive-timestamps: under --fix, a future stamp in a committed archived REQ
# (inside an archived UR folder, proving the recursive scan) is rewritten to the
# introducing commit's author time and the correction logs the sourcing commit hash.
audit_fix_project="$fixture_root/audit-fix-project"
fixture_repo_init "$audit_fix_project"
mkdir -p "$audit_fix_project/do-work/archive/UR-901"
printf -- '---\nid: REQ-901\nstatus: completed\ncreated_at: 2026-08-10T09:00:00Z\nclaimed_at: 2026-08-10T10:00:00Z\ncompleted_at: 2093-05-05T05:05:05Z\n---\nbody\n' \
  > "$audit_fix_project/do-work/archive/UR-901/REQ-901-future.md"
git -C "$audit_fix_project" add -A
GIT_AUTHOR_DATE='2026-08-12T10:00:00Z' GIT_COMMITTER_DATE='2026-08-12T10:05:00Z' \
  git -C "$audit_fix_project" commit -qm fixture
audit_fix_output="$("$core_scripts/audit-archive-timestamps.sh" --fix "$audit_fix_project")" \
  || fail_case 'audit-archive-timestamps fixing case returned nonzero'
grep -q '^completed_at: 2026-08-12T10:00:00Z$' "$audit_fix_project/do-work/archive/UR-901/REQ-901-future.md" \
  || fail_case 'audit-archive-timestamps fixing case did not rewrite the stamp to the introducing commit author time'
printf '%s' "$audit_fix_output" \
  | grep -Eq 'REQ-901-future\.md completed_at: 2093-05-05T05:05:05Z -> 2026-08-12T10:00:00Z \(commit [0-9a-f]{7} author time\)' \
  || fail_case 'audit-archive-timestamps fixing case did not log the correction with its sourcing commit hash'

# audit-archive-timestamps: without --fix the default run only reports — the pending
# correction prints as `would repair`, the archived file keeps its bytes, and the
# exit code is nonzero so a caller can gate on findings.
audit_report_project="$fixture_root/audit-report-project"
fixture_repo_init "$audit_report_project"
mkdir -p "$audit_report_project/do-work/archive"
printf -- '---\nid: REQ-902\nstatus: completed\ncreated_at: 2026-08-10T09:00:00Z\ncompleted_at: 2093-05-05T05:05:05Z\n---\nbody\n' \
  > "$audit_report_project/do-work/archive/REQ-902-future.md"
git -C "$audit_report_project" add -A
GIT_AUTHOR_DATE='2026-08-12T10:00:00Z' GIT_COMMITTER_DATE='2026-08-12T10:05:00Z' \
  git -C "$audit_report_project" commit -qm fixture
cp "$audit_report_project/do-work/archive/REQ-902-future.md" "$fixture_root/audit-report-before.md"
audit_report_output="$("$core_scripts/audit-archive-timestamps.sh" "$audit_report_project")" \
  && fail_case 'audit-archive-timestamps report-only case exited zero with a correction pending'
printf '%s' "$audit_report_output" \
  | grep -q 'would repair do-work/archive/REQ-902-future.md completed_at: 2093-05-05T05:05:05Z -> 2026-08-12T10:00:00Z' \
  || fail_case 'audit-archive-timestamps report-only case did not print the pending correction'
cmp -s "$fixture_root/audit-report-before.md" "$audit_report_project/do-work/archive/REQ-902-future.md" \
  || fail_case 'audit-archive-timestamps report-only case changed bytes without --fix'

# audit-archive-timestamps: a clean committed archive passes through byte-identical
# with exit zero — and queue scope is never scanned: a wrong queue stamp belongs to
# the hook-run repairer, not the archive audit.
audit_clean_project="$fixture_root/audit-clean-project"
fixture_repo_init "$audit_clean_project"
mkdir -p "$audit_clean_project/do-work/archive" "$audit_clean_project/do-work/queue"
printf -- '---\nid: REQ-903\nstatus: completed\ncreated_at: 2026-08-01T09:00:00Z\nclaimed_at: 2026-08-02T09:00:00Z\ncompleted_at: 2026-08-03T09:00:00Z\n---\nbody\n' \
  > "$audit_clean_project/do-work/archive/REQ-903-clean.md"
printf -- '---\nid: REQ-904\nstatus: pending\ncreated_at: 2093-06-06T06:06:06Z\n---\nbody\n' \
  > "$audit_clean_project/do-work/queue/REQ-904-queue.md"
git -C "$audit_clean_project" add -A
GIT_AUTHOR_DATE='2026-08-04T09:00:00Z' GIT_COMMITTER_DATE='2026-08-04T09:00:00Z' \
  git -C "$audit_clean_project" commit -qm fixture
cp "$audit_clean_project/do-work/archive/REQ-903-clean.md" "$fixture_root/audit-clean-before.md"
cp "$audit_clean_project/do-work/queue/REQ-904-queue.md" "$fixture_root/audit-queue-before.md"
audit_clean_output="$("$core_scripts/audit-archive-timestamps.sh" --fix "$audit_clean_project")" \
  || fail_case 'audit-archive-timestamps clean-archive case returned nonzero'
printf '%s' "$audit_clean_output" | grep -q 'archive audit clean' \
  || fail_case 'audit-archive-timestamps clean-archive case did not report a clean audit'
cmp -s "$fixture_root/audit-clean-before.md" "$audit_clean_project/do-work/archive/REQ-903-clean.md" \
  || fail_case 'audit-archive-timestamps clean-archive case changed a clean archived file'
cmp -s "$fixture_root/audit-queue-before.md" "$audit_clean_project/do-work/queue/REQ-904-queue.md" \
  || fail_case 'audit-archive-timestamps clean-archive case wrote into queue scope'

# audit-archive-timestamps: an impossible ordering in a committed archived REQ is
# clamped so created_at <= claimed_at <= completed_at — the introducing commit's
# author time precedes the anchor here, so both later fields land on the clamp floor.
audit_order_project="$fixture_root/audit-order-project"
fixture_repo_init "$audit_order_project"
mkdir -p "$audit_order_project/do-work/archive"
printf -- '---\nid: REQ-905\nstatus: completed\ncreated_at: 2026-08-10T12:00:00Z\nclaimed_at: 2026-08-01T09:00:00Z\ncompleted_at: 2026-08-03T10:00:00Z\n---\nbody\n' \
  > "$audit_order_project/do-work/archive/REQ-905-order.md"
git -C "$audit_order_project" add -A
GIT_AUTHOR_DATE='2026-08-05T08:00:00Z' GIT_COMMITTER_DATE='2026-08-05T08:00:00Z' \
  git -C "$audit_order_project" commit -qm fixture
audit_order_output="$("$core_scripts/audit-archive-timestamps.sh" --fix "$audit_order_project")" \
  || fail_case 'audit-archive-timestamps ordering case returned nonzero'
grep -q '^claimed_at: 2026-08-10T12:00:00Z$' "$audit_order_project/do-work/archive/REQ-905-order.md" \
  || fail_case 'audit-archive-timestamps ordering case did not clamp claimed_at up to created_at'
grep -q '^completed_at: 2026-08-10T12:00:00Z$' "$audit_order_project/do-work/archive/REQ-905-order.md" \
  || fail_case 'audit-archive-timestamps ordering case did not clamp completed_at up to the repaired claimed_at'
printf '%s' "$audit_order_output" | grep -q 'clamped to 2026-08-10T12:00:00Z' \
  || fail_case 'audit-archive-timestamps ordering case did not log the clamp'

# audit-archive-timestamps: an archived defect whose introducing commit cannot be
# blamed (an uncommitted file) is reported and left byte-identical — replacements
# derive from git alone, and the file mtime is never consulted as a fallback.
audit_blameless_project="$fixture_root/audit-blameless-project"
fixture_repo_init "$audit_blameless_project"
mkdir -p "$audit_blameless_project/do-work/archive"
printf -- '---\nid: REQ-906\nstatus: completed\ncreated_at: 2093-07-07T07:07:07Z\n---\nbody\n' \
  > "$audit_blameless_project/do-work/archive/REQ-906-untracked.md"
TZ=UTC touch -m -t 202608101200.00 "$audit_blameless_project/do-work/archive/REQ-906-untracked.md"
cp "$audit_blameless_project/do-work/archive/REQ-906-untracked.md" "$fixture_root/audit-blameless-before.md"
audit_blameless_output="$("$core_scripts/audit-archive-timestamps.sh" --fix "$audit_blameless_project")" \
  && fail_case 'audit-archive-timestamps blameless case exited zero on an unrepairable defect'
printf '%s' "$audit_blameless_output" | grep -q 'FAILED to repair' \
  || fail_case 'audit-archive-timestamps blameless case did not report the unrepairable defect'
printf '%s' "$audit_blameless_output" | grep -q 'file mtime' \
  && fail_case 'audit-archive-timestamps blameless case offered the file mtime as a source'
cmp -s "$fixture_root/audit-blameless-before.md" "$audit_blameless_project/do-work/archive/REQ-906-untracked.md" \
  || fail_case 'audit-archive-timestamps blameless case changed the file it could not derive for'

# audit-archive-timestamps: the widened shapes (space-separated, quoted, CRLF,
# BOM) repair through the archive scan too — the fix lives in the sourced
# library, and this pins the shared-fix-reaches-both-tools property instead of
# assuming it (REQ-255; REQ-247 review).
audit_shapes_project="$fixture_root/audit-shapes-project"
fixture_repo_init "$audit_shapes_project"
mkdir -p "$audit_shapes_project/do-work/archive"
printf -- '---\nid: REQ-907\nstatus: completed\ncreated_at: 2093-01-01 00:00:00\n---\nbody\n' \
  > "$audit_shapes_project/do-work/archive/REQ-907-space.md"
printf -- '---\nid: REQ-908\nstatus: completed\ncreated_at: "2093-01-01 00:00:00"\n---\nbody\n' \
  > "$audit_shapes_project/do-work/archive/REQ-908-quoted-space.md"
printf -- '---\r\nid: REQ-909\r\nstatus: completed\r\ncreated_at: 2093-03-03T03:03:03Z\r\n---\r\nbody\r\n' \
  > "$audit_shapes_project/do-work/archive/REQ-909-crlf.md"
printf -- '\xef\xbb\xbf---\nid: REQ-910\nstatus: completed\ncreated_at: 2093-04-04T04:04:04Z\n---\nbody\n' \
  > "$audit_shapes_project/do-work/archive/REQ-910-bom.md"
git -C "$audit_shapes_project" add -A
GIT_AUTHOR_DATE='2026-08-12T10:00:00Z' GIT_COMMITTER_DATE='2026-08-12T10:05:00Z' \
  git -C "$audit_shapes_project" commit -qm fixture
audit_shapes_output="$("$core_scripts/audit-archive-timestamps.sh" --fix "$audit_shapes_project")" \
  || fail_case 'audit-archive-timestamps widened-shapes case returned nonzero'
for widened_fixture in REQ-907-space REQ-908-quoted-space REQ-910-bom; do
  grep -q '^created_at: 2026-08-12T10:00:00Z$' "$audit_shapes_project/do-work/archive/$widened_fixture.md" \
    || fail_case "audit-archive-timestamps widened-shapes case did not repair $widened_fixture through the archive scan"
done
grep -q $'^created_at: 2026-08-12T10:00:00Z\r$' "$audit_shapes_project/do-work/archive/REQ-909-crlf.md" \
  || fail_case 'audit-archive-timestamps widened-shapes case did not repair the CRLF file (or dropped the CR)'
[ "$(head -c 3 "$audit_shapes_project/do-work/archive/REQ-910-bom.md")" = "$(printf '\xef\xbb\xbf')" ] \
  || fail_case 'audit-archive-timestamps widened-shapes case did not keep the BOM bytes in place'
printf '%s' "$audit_shapes_output" \
  | grep -q 'REQ-907-space.md created_at: 2093-01-01 00:00:00 -> 2026-08-12T10:00:00Z' \
  || fail_case 'audit-archive-timestamps widened-shapes case did not report the full old value in the audit line'

# audit-archive-timestamps: the refused and duplicate-key shapes hold through
# the archive scan too — a calendar-impossible stamp and a numeric-offset stamp
# stay byte-identical (and are not defects), while a duplicated anchor repairs
# on its effective (last) occurrence from the introducing commit's author time
# (REQ-255; REQ-257). The refusals belong to the sourced library, so this is
# what fails if the auditor ever grows its own recognizer beside the shared one.
audit_parity_project="$fixture_root/audit-parity-project"
fixture_repo_init "$audit_parity_project"
mkdir -p "$audit_parity_project/do-work/archive"
printf -- '---\nid: REQ-911\nstatus: completed\ncreated_at: 9999-99-99T99:99:99Z\n---\nbody\n' \
  > "$audit_parity_project/do-work/archive/REQ-911-impossible.md"
printf -- '---\nid: REQ-912\nstatus: completed\ncreated_at: 2026-08-10T12:00:00Z\nclaimed_at: 2026-08-11T12:00:00Z\nclaimed_at: 2026-08-01T09:00:00Z\n---\nbody\n' \
  > "$audit_parity_project/do-work/archive/REQ-912-duplicate-anchor.md"
printf -- '---\nid: REQ-913\nstatus: completed\ncreated_at: 2093-01-01T00:00:00+02:00\n---\nbody\n' \
  > "$audit_parity_project/do-work/archive/REQ-913-offset.md"
git -C "$audit_parity_project" add -A
GIT_AUTHOR_DATE='2026-08-12T10:00:00Z' GIT_COMMITTER_DATE='2026-08-12T10:05:00Z' \
  git -C "$audit_parity_project" commit -qm fixture
cp "$audit_parity_project/do-work/archive/REQ-911-impossible.md" "$fixture_root/audit-impossible-before.md"
cp "$audit_parity_project/do-work/archive/REQ-913-offset.md" "$fixture_root/audit-offset-before.md"
audit_parity_output="$("$core_scripts/audit-archive-timestamps.sh" --fix "$audit_parity_project")" \
  || fail_case 'audit-archive-timestamps refusal-parity case returned nonzero'
cmp -s "$fixture_root/audit-impossible-before.md" "$audit_parity_project/do-work/archive/REQ-911-impossible.md" \
  || fail_case 'audit-archive-timestamps refusal-parity case erased a calendar-impossible stamp in the archive'
cmp -s "$fixture_root/audit-offset-before.md" "$audit_parity_project/do-work/archive/REQ-913-offset.md" \
  || fail_case 'audit-archive-timestamps refusal-parity case repaired a numeric-offset stamp in the archive — the offset refusal is the sourced library one and must reach every tool built on it'
grep -q '^claimed_at: 2026-08-12T10:00:00Z$' "$audit_parity_project/do-work/archive/REQ-912-duplicate-anchor.md" \
  || fail_case 'audit-archive-timestamps refusal-parity case did not repair the effective (last) anchor occurrence'
grep -q '^claimed_at: 2026-08-11T12:00:00Z$' "$audit_parity_project/do-work/archive/REQ-912-duplicate-anchor.md" \
  || fail_case 'audit-archive-timestamps refusal-parity case rewrote the shadowed first occurrence'
printf '%s' "$audit_parity_output" | grep -E '(would repair|repaired) ' | grep -qE 'REQ-911|REQ-913' \
  && fail_case 'audit-archive-timestamps refusal-parity case logged a refused stamp as a correction'

# audit-archive-timestamps: the fence and padding shapes reach the archive scan
# too, because both live in the sourced extractor and recognizer rather than in
# the queue/working scan (REQ-267, both instances through the second scope). An
# unterminated fence is refused whole here as well; a stamp padded inside its
# quotes is repaired from the introducing commit's author time.
audit_shape_project="$fixture_root/audit-shape-project"
fixture_repo_init "$audit_shape_project"
mkdir -p "$audit_shape_project/do-work/archive"
printf -- '---\nid: REQ-914\nstatus: completed\ncreated_at: 2093-01-01T00:00:00Z' \
  > "$audit_shape_project/do-work/archive/REQ-914-unterminated.md"
printf -- '---\nid: REQ-915\nstatus: completed\ncreated_at: "2093-01-01 00:00:00 "\n---\nbody\n' \
  > "$audit_shape_project/do-work/archive/REQ-915-padded-quote.md"
git -C "$audit_shape_project" add -A
GIT_AUTHOR_DATE='2026-08-12T10:00:00Z' GIT_COMMITTER_DATE='2026-08-12T10:05:00Z' \
  git -C "$audit_shape_project" commit -qm fixture
cp "$audit_shape_project/do-work/archive/REQ-914-unterminated.md" "$fixture_root/audit-unterminated-before.md"
audit_shape_output="$("$core_scripts/audit-archive-timestamps.sh" --fix "$audit_shape_project")" \
  || fail_case 'audit-archive-timestamps shape-parity case returned nonzero'
cmp -s "$fixture_root/audit-unterminated-before.md" "$audit_shape_project/do-work/archive/REQ-914-unterminated.md" \
  || fail_case 'audit-archive-timestamps shape-parity case rewrote an archived file whose fence never closes'
grep -q '^created_at: 2026-08-12T10:00:00Z$' "$audit_shape_project/do-work/archive/REQ-915-padded-quote.md" \
  || fail_case 'audit-archive-timestamps shape-parity case refused a padded quoted instant in the archive'
printf '%s' "$audit_shape_output" | grep -E '(would repair|repaired) ' | grep -q 'REQ-914' \
  && fail_case 'audit-archive-timestamps shape-parity case logged a correction for the refused unterminated file'

# audit-archive-timestamps: a value the recognizer refused is voiced, and the summary
# stops calling the run clean. The refusal itself is right and permanent, so the exit
# status stays 0 — a nonzero exit here would print a FAILED line into every session
# that no one can ever heal, which is the REQ-267 wedge (REQ-268 instance 1).
audit_voiced_project="$fixture_root/audit-voiced-project"
fixture_repo_init "$audit_voiced_project"
mkdir -p "$audit_voiced_project/do-work/archive"
printf -- '---\nid: REQ-916\nstatus: completed\ncreated_at: 9999-99-99T99:99:99Z\n---\nbody\n' \
  > "$audit_voiced_project/do-work/archive/REQ-916-impossible.md"
fixture_repo_commit_all "$audit_voiced_project" fixture
cp "$audit_voiced_project/do-work/archive/REQ-916-impossible.md" "$fixture_root/audit-voiced-before.md"
audit_voiced_output="$("$core_scripts/audit-archive-timestamps.sh" "$audit_voiced_project")" \
  || fail_case 'audit-archive-timestamps voiced-refusal case exited nonzero on a permanent refusal'
printf '%s' "$audit_voiced_output" | grep -q 'REQ-916-impossible.md' \
  || fail_case 'audit-archive-timestamps voiced-refusal case did not name the file holding the refused value'
printf '%s' "$audit_voiced_output" | grep -q 'refused' \
  || fail_case 'audit-archive-timestamps voiced-refusal case did not say the value was refused'
printf '%s' "$audit_voiced_output" | grep -q 'audit clean' \
  && fail_case 'audit-archive-timestamps voiced-refusal case still reported clean for a value it never inspected'
cmp -s "$fixture_root/audit-voiced-before.md" "$audit_voiced_project/do-work/archive/REQ-916-impossible.md" \
  || fail_case 'audit-archive-timestamps voiced-refusal case changed the file it refused'

# audit-archive-timestamps: a file walk that could not complete is a failure, never a
# clean answer with a count of zero. The exit status of the walk is what decides, not
# the number of iterations it managed (REQ-268 instance 2).
audit_walk_project="$fixture_root/audit-walk-project"
fixture_repo_init "$audit_walk_project"
mkdir -p "$audit_walk_project/do-work/archive"
printf -- '---\nid: REQ-917\nstatus: completed\ncreated_at: 2026-08-10T12:00:00Z\n---\nbody\n' \
  > "$audit_walk_project/do-work/archive/REQ-917-clean.md"
fixture_repo_commit_all "$audit_walk_project" fixture
audit_walk_bin="$fixture_root/audit-walk-bin"
mkdir -p "$audit_walk_bin"
printf '%s\n' '#!/usr/bin/env bash' 'exit 3' > "$audit_walk_bin/find"
chmod +x "$audit_walk_bin/find"
audit_walk_output="$(PATH="$audit_walk_bin:$PATH" "$core_scripts/audit-archive-timestamps.sh" "$audit_walk_project" 2>&1)" \
  && fail_case 'audit-archive-timestamps failed-walk case exited zero after the walk failed'
printf '%s' "$audit_walk_output" | grep -q 'audit clean' \
  && fail_case 'audit-archive-timestamps failed-walk case reported clean for an archive it never scanned'

# audit-archive-timestamps: without its shared library the auditor inspects nothing, so
# it must say so rather than counting files it never read. A lone copy of the script is
# the real shape of this — a partial install, or a moved sibling (REQ-268 sweep).
audit_lone_project="$fixture_root/audit-lone-project"
fixture_repo_init "$audit_lone_project"
mkdir -p "$audit_lone_project/do-work/archive"
printf -- '---\nid: REQ-918\nstatus: completed\ncreated_at: 2093-01-01T00:00:00Z\n---\nbody\n' \
  > "$audit_lone_project/do-work/archive/REQ-918-future.md"
fixture_repo_commit_all "$audit_lone_project" fixture
audit_lone_dir="$fixture_root/audit-lone-dir"
fixture_repo_clone_script "$core_scripts/audit-archive-timestamps.sh" "$audit_lone_dir/audit-archive-timestamps.sh"
audit_lone_output="$("$audit_lone_dir/audit-archive-timestamps.sh" "$audit_lone_project" 2>&1)" \
  && fail_case 'audit-archive-timestamps missing-library case exited zero with its shared machinery absent'
printf '%s' "$audit_lone_output" | grep -q 'audit clean' \
  && fail_case 'audit-archive-timestamps missing-library case reported clean without its detection predicate'

prescribed_shell_finish
