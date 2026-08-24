#!/usr/bin/env bash
# Fixture execution proofs for qualify.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

# qualify: an output line added to a file that owns its process exit is the file's own
# reporting, not a debug artifact — the scan passes it and names the reason (REQ-254; the
# scan once FAILed a checker's only success line, REQ-244's remediation). Both output
# primitives (print(, console.log) get the same treatment: the condition is the class,
# not the token. Serial and range modes are both exercised — each reads a different diff.
core_checks="$repo_root/skills/do-work/tools/checks"
qualify_repo="$fixture_root/qualify-repo"
fixture_repo_init "$qualify_repo"
mkdir -p "$qualify_repo/src" "$qualify_repo/do-work"
printf '%s\n' '#!/usr/bin/env bash' 'python3 - <<PY' 'import sys' 'if sites_missing:' \
  '    raise SystemExit("site check failed")' 'PY' > "$qualify_repo/site-checker.sh"
printf '%s\n' 'if (writeFailed) {' '  process.exit(1);' '}' > "$qualify_repo/report-writer.js"
printf '%s\n' 'def parse_value(raw_text):' '    return raw_text.strip()' > "$qualify_repo/src/value_parser.py"
fixture_repo_commit_all "$qualify_repo" base
qualify_base_commit="$(git -C "$qualify_repo" rev-parse HEAD)"
printf '%s\n' '## Implementation Summary' \
  '- `site-checker.sh` (modified) — success line' \
  '- `report-writer.js` (modified) — success line' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-950-reporting.md"
printf 'print("site check: all sites cited")\n' >> "$qualify_repo/site-checker.sh"
printf 'console.log("report written");\n' >> "$qualify_repo/report-writer.js"
qualify_reporting_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-950-reporting.md 2>&1)" \
  || fail_case 'qualify reporter-output serial case FAILed a check file on its own success line'
printf '%s' "$qualify_reporting_output" | grep -q 'reporting' \
  || fail_case 'qualify reporter-output serial case did not name the pass reason (reporting)'
fixture_repo_commit_all "$qualify_repo" reporting
qualify_reporting_commit="$(git -C "$qualify_repo" rev-parse HEAD)"
(cd "$qualify_repo" && DO_WORK_DIFF_RANGE="$qualify_base_commit..$qualify_reporting_commit" \
  "$core_checks/qualify.sh" do-work/REQ-950-reporting.md >/dev/null 2>&1) \
  || fail_case 'qualify reporter-output range case FAILed a check file on its own success line'

# qualify: an output primitive added to a file that never ends its own process still FAILs
# as leftover instrumentation, and the FAIL names the file and the reason (REQ-254).
printf '%s\n' '## Implementation Summary' \
  '- `src/value_parser.py` (modified) — parsing tweak' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-951-instrumentation.md"
printf 'print(raw_text)\n' >> "$qualify_repo/src/value_parser.py"
qualify_instrumentation_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-951-instrumentation.md 2>&1)" \
  && fail_case 'qualify instrumentation serial case passed a debug print in library code'
printf '%s' "$qualify_instrumentation_output" | grep -q 'instrumentation' \
  || fail_case 'qualify instrumentation serial case did not name the FAIL reason (instrumentation)'
printf '%s' "$qualify_instrumentation_output" | grep -q 'src/value_parser.py' \
  || fail_case 'qualify instrumentation serial case did not name the offending file'
fixture_repo_commit_all "$qualify_repo" instrumentation
qualify_instrumentation_commit="$(git -C "$qualify_repo" rev-parse HEAD)"
(cd "$qualify_repo" && DO_WORK_DIFF_RANGE="$qualify_reporting_commit..$qualify_instrumentation_commit" \
  "$core_checks/qualify.sh" do-work/REQ-951-instrumentation.md >/dev/null 2>&1) \
  && fail_case 'qualify instrumentation range case passed a debug print in library code'

# qualify: the reporter exemption covers output primitives only — an unfinished-work
# marker (TODO) added to the same reporter file still FAILs (REQ-254; pins the hole a
# file-level exemption would open).
printf '%s\n' '## Implementation Summary' \
  '- `site-checker.sh` (modified) — regex note' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-952-marker.md"
printf '# TODO: tighten the site regex\n' >> "$qualify_repo/site-checker.sh"
qualify_marker_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-952-marker.md 2>&1)" \
  && fail_case 'qualify unfinished-marker case let a reporter file carry a fresh TODO'
printf '%s' "$qualify_marker_output" | grep -q 'debug artifacts' \
  || fail_case 'qualify unfinished-marker case did not report the TODO as a debug artifact'
git -C "$qualify_repo" checkout -q -- site-checker.sh

# qualify: ownership is judged at the BASE revision, so an exit idiom added in the same
# diff as a debug print does not retroactively make library code a reporter (REQ-263;
# REQ-254's review reproduced the downgrade to WARN).
printf '%s\n' '## Implementation Summary' \
  '- `src/value_parser.py` (modified) — parsing tweak' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-953-same-diff-exit.md"
printf '%s\n' 'print(raw_text)' '' 'if _run_as_script:' '    import sys' '    sys.exit(0)' \
  >> "$qualify_repo/src/value_parser.py"
qualify_same_diff_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-953-same-diff-exit.md 2>&1)" \
  && fail_case 'qualify same-diff-exit case let an exit added in this diff excuse a debug print in library code'
printf '%s' "$qualify_same_diff_output" | grep -q 'never ends its own process' \
  || fail_case 'qualify same-diff-exit case did not FAIL on the base-revision ownership verdict'
git -C "$qualify_repo" checkout -q -- src/value_parser.py

# qualify: a brand-NEW checker whose prints and whose exit arrive together is still a
# reporter — it has no base revision to be judged at, so it is judged on its own content.
# This is the case that makes the categorical rule ("exit added in this diff => FAIL")
# wrong, and it must stay a WARN (REQ-263).
printf '%s\n' '## Implementation Summary' \
  '- `new-site-checker.sh` (new) — brand new check' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-954-new-checker.md"
printf '%s\n' '#!/usr/bin/env bash' 'if [ -z "${1:-}" ]; then' '  exit 2' 'fi' \
  'print_report() { :; }' 'console.log("new checker ok")' \
  > "$qualify_repo/new-site-checker.sh"
(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-954-new-checker.md >/dev/null 2>&1) \
  || fail_case 'qualify new-checker case FAILed a brand-new checker on its own reporting line'
rm -f "$qualify_repo/new-site-checker.sh"

# qualify: docstring prose is not ownership — a library file whose docstring merely says
# "exit 1 on failure" still FAILs on an added debug print. The bare `exit N` form must be
# statement-shaped, not prose (REQ-263).
printf '%s\n' '"""Render values.' '' 'The caller should exit 1 on failure.' '"""' '' \
  'def render_value(raw_text):' '    return raw_text' > "$qualify_repo/src/value_renderer.py"
fixture_repo_commit_all "$qualify_repo" renderer
printf '%s\n' '## Implementation Summary' \
  '- `src/value_renderer.py` (modified) — render tweak' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-955-docstring-exit.md"
printf 'print("debug", raw_text)\n' >> "$qualify_repo/src/value_renderer.py"
qualify_docstring_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-955-docstring-exit.md 2>&1)" \
  && fail_case 'qualify docstring-exit case let docstring prose "exit 1 on failure" pass as process-exit ownership'
printf '%s' "$qualify_docstring_output" | grep -q 'src/value_renderer.py' \
  || fail_case 'qualify docstring-exit case did not name the offending file'

# qualify: the DOCUMENTED RESIDUAL of that narrowing, pinned so the boundary is a stated
# limit rather than a surprise. Prose with a word before the idiom is now correctly
# rejected ("the caller should exit 1 on failure", "On a bad row, exit 1"); what survives
# is the narrower shape of a docstring line consisting of nothing BUT the idiom and its
# status, which is indistinguishable from a shell statement without a language parser.
# If this case ever starts failing, the probe got sharper and the residual note in
# qualify.sh should be updated to match (REQ-250: pin every documented limitation with a
# fixture that can fail).
printf '%s\n' '"""Render totals.' '' 'On a bad row:' '' '    exit 1' '"""' '' \
  'def render_total(raw_text):' '    return raw_text' > "$qualify_repo/src/total_renderer.py"
fixture_repo_commit_all "$qualify_repo" total-renderer
printf '%s\n' '## Implementation Summary' \
  '- `src/total_renderer.py` (modified) — total tweak' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-956-residual.md"
printf 'print("debug", raw_text)\n' >> "$qualify_repo/src/total_renderer.py"
(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-956-residual.md >/dev/null 2>&1) \
  || fail_case 'qualify documented-residual case now FAILs — the exit probe got sharper than its stated limit; update the residual note in qualify.sh'
git -C "$qualify_repo" checkout -q -- src/total_renderer.py src/value_renderer.py

# qualify: the WARN branch prints the matched lines exactly as the FAIL branch does, so
# "confirm from the diff" does not cost a second manual dig (REQ-263).
printf '%s\n' '## Implementation Summary' \
  '- `site-checker.sh` (modified) — another success line' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-957-warn-lines.md"
printf 'print("site check: %%d sites" %% site_total)\n' >> "$qualify_repo/site-checker.sh"
qualify_warn_lines_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-957-warn-lines.md 2>&1)" \
  || fail_case 'qualify warn-lines case FAILed a reporter file on its own success line'
printf '%s' "$qualify_warn_lines_output" | grep -q 'site check: %d sites' \
  || fail_case 'qualify warn-lines case WARNed without printing the matched line — the FAIL branch prints its lines and the WARN must too'
git -C "$qualify_repo" checkout -q -- site-checker.sh

# qualify: both artifact scans see an UNTRACKED, non-ignored file in serial mode. Step 6.3
# runs before the REQ's commit, so a new source file the builder never staged appears in
# neither `git diff` nor `git diff --staged` and used to be scanned by neither half
# (REQ-263 addendum, from the 2026-08-20 consumer review). A correctly-ignored file in the
# same tree stays unscanned — the ignore filter is what `--exclude-standard` buys.
printf '%s\n' 'ignored-helper.js' > "$qualify_repo/.gitignore"
fixture_repo_commit_all "$qualify_repo" ignore-rules
printf '%s\n' '## Implementation Summary' \
  '- `src/untracked_helper.js` (new) — helper' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-958-untracked.md"
printf '%s\n' 'function helper(raw) {' '  console.log("debug", raw);' \
  '  // TODO: handle the empty case' '  return raw;' '}' \
  > "$qualify_repo/src/untracked_helper.js"
printf '%s\n' 'console.log("ignored debug");' '// TODO: ignored marker' \
  > "$qualify_repo/ignored-helper.js"
qualify_untracked_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-958-untracked.md 2>&1)" \
  && fail_case 'qualify untracked case passed an untracked file carrying a console.log and a TODO'
printf '%s' "$qualify_untracked_output" | grep -q 'src/untracked_helper.js' \
  || fail_case 'qualify untracked case did not name the untracked file'
printf '%s' "$qualify_untracked_output" | grep -q 'debug artifacts' \
  || fail_case 'qualify untracked case missed the TODO — the unfinished-marker scan still does not read untracked files'
printf '%s' "$qualify_untracked_output" | grep -q 'leftover instrumentation' \
  || fail_case 'qualify untracked case missed the console.log — the output-primitive scan still does not read untracked files'
printf '%s' "$qualify_untracked_output" | grep -q 'ignored-helper.js' \
  && fail_case 'qualify untracked case scanned a correctly-ignored file — --exclude-standard is the ignore filter and must be honored'

# qualify: worktree dispatch mode reads committed work, so an untracked file left lying in
# the tree is not its business and must stay unscanned (REQ-263 addendum).
# The previous case's untracked offenders are removed BEFORE this commit: leaving them
# would sweep them into the range and FAIL this case on the prior case's fixture.
rm -f "$qualify_repo/src/untracked_helper.js" "$qualify_repo/ignored-helper.js"
qualify_untracked_base="$(git -C "$qualify_repo" rev-parse HEAD)"
printf 'const parsed = 1;\n' >> "$qualify_repo/report-writer.js"
fixture_repo_commit_all "$qualify_repo" range-untracked
qualify_untracked_range_commit="$(git -C "$qualify_repo" rev-parse HEAD)"
printf '%s\n' '## Implementation Summary' \
  '- `report-writer.js` (modified) — parse tweak' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-959-range-untracked.md"
printf '%s\n' 'console.log("stray");' '// TODO: stray marker' > "$qualify_repo/src/stray_helper.js"
qualify_range_untracked_output="$(cd "$qualify_repo" \
  && DO_WORK_DIFF_RANGE="$qualify_untracked_base..$qualify_untracked_range_commit" \
  "$core_checks/qualify.sh" do-work/REQ-959-range-untracked.md 2>&1)" \
  || fail_case 'qualify range-untracked case FAILed on an untracked file that worktree dispatch mode must not read'
printf '%s' "$qualify_range_untracked_output" | grep -q 'stray_helper.js' \
  && fail_case 'qualify range-untracked case scanned an untracked file in worktree dispatch mode — that mode reads committed work only'
rm -f "$qualify_repo/src/stray_helper.js" "$qualify_repo/src/untracked_helper.js" "$qualify_repo/ignored-helper.js"

# qualify: a REQ carrying NO P-A-U section leaves Check 4 disarmed rather than passed —
# every [UNIFY]-gated FAIL keys on the box's state, so a missing box satisfies all of them
# by absence. The run must say so instead of printing a bare OK (REQ-264; this is the shape
# REQ-254's own qualification "Passed" with).
printf '%s\n' '## Implementation Summary' \
  '- `src/value_parser.py` (modified) — parsing tweak' \
  > "$qualify_repo/do-work/REQ-960-no-pau.md"
printf 'print(raw_text)\n' >> "$qualify_repo/src/value_parser.py"
qualify_no_pau_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-960-no-pau.md 2>&1)"
printf '%s' "$qualify_no_pau_output" | grep -q 'DISARMED' \
  || fail_case 'qualify no-P-A-U case printed no disarmed-audit warning — a REQ with no P-A-U section still passes Check 4 silently'
git -C "$qualify_repo" checkout -q -- src/value_parser.py

# qualify: the vacuity guard on that warning — a REQ that DOES carry the section must not
# get it, or the warning would fire on every run and mean nothing (REQ-264).
printf '%s' "$qualify_reporting_output" | grep -q 'DISARMED' \
  && fail_case 'qualify no-P-A-U warning fired on a REQ that carries the section — the check is keying on the wrong thing'

# qualify: a MOVED debug-artifact marker is not an added one (REQ-301). Every scan reads
# `^+` out of a diff, so relocated text used to read exactly like written text and every
# code-relocation REQ FAILed the audit on markers that already existed — REQ-258 hit it on
# four deliberate fixture TODO strings and had to override the FAIL with evidence. The move
# is downgraded to a named WARN, never dropped.
mkdir -p "$qualify_repo/relocation"
printf '%s\n' '#!/usr/bin/env bash' 'case_alpha() {' '  printf "# TODO: alpha fixture\n" > f' '}' \
  'case_beta() {' '  printf "# FIXME: beta fixture\n" > g' '}' > "$qualify_repo/relocation/suite.sh"
fixture_repo_commit_all "$qualify_repo" relocation-base
printf '%s\n' '## Implementation Summary' \
  '- `relocation/suite.sh` (modified) — now dispatches' \
  '- `relocation/alpha.sh` (new) — alpha case' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-970-relocated-marker.md"
printf '%s\n' '#!/usr/bin/env bash' 'case_alpha() {' '  printf "# TODO: alpha fixture\n" > f' '}' \
  > "$qualify_repo/relocation/alpha.sh"
printf '%s\n' '#!/usr/bin/env bash' 'source relocation/alpha.sh' \
  'case_beta() {' '  printf "# FIXME: beta fixture\n" > g' '}' > "$qualify_repo/relocation/suite.sh"
qualify_relocated_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-970-relocated-marker.md 2>&1)" \
  || fail_case 'qualify relocated-marker case FAILed a marker that was moved, not added — every code-relocation REQ trips this'
printf '%s' "$qualify_relocated_output" | grep -q 'relocated, not added' \
  || fail_case 'qualify relocated-marker case did not name the moved marker — a relocation must be downgraded to a named WARN, never dropped'

# qualify: the masking risk the REQ names, pinned. A marker DUPLICATED rather than moved
# raises its occurrence count, so it is a genuine addition and must still FAIL — presence
# at base is not the test, count is (REQ-301).
printf '%s\n' '## Implementation Summary' \
  '- `relocation/gamma.sh` (new) — duplicate, not a move' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-971-duplicated-marker.md"
# Byte-identical to the line still sitting in relocation/suite.sh, deliberately: the whole
# point is that the text is unchanged and only its OCCURRENCE COUNT went up. A near-copy
# would FAIL as fresh text under either rule and pin nothing.
printf '%s\n' '#!/usr/bin/env bash' 'case_gamma() {' '  printf "# FIXME: beta fixture\n" > g' '}' \
  > "$qualify_repo/relocation/gamma.sh"
qualify_duplicated_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-971-duplicated-marker.md 2>&1)" \
  && fail_case 'qualify duplicated-marker case excused a marker that was copied rather than moved — mere presence at base was used as the test instead of occurrence count'
printf '%s' "$qualify_duplicated_output" | grep -q 'relocation/gamma.sh' \
  || fail_case 'qualify duplicated-marker case did not name the file that gained the extra copy'
rm -f "$qualify_repo/relocation/gamma.sh"

# qualify: a fresh marker sitting beside a moved one in the SAME file still FAILs, and the
# FAIL names only the fresh line while the moved one is reported separately (REQ-301). This
# is what makes the downgrade safe: nothing is pooled and nothing is lost.
printf '%s\n' '## Implementation Summary' \
  '- `relocation/alpha.sh` (new) — alpha case plus a leftover' \
  '' '## AI Execution State' \
  '- [x] **[PLAN]** done' '- [x] **[APPLY]** done' '- [x] **[UNIFY]** done' \
  > "$qualify_repo/do-work/REQ-972-fresh-beside-moved.md"
printf '%s\n' '  # TODO: brand new leftover' >> "$qualify_repo/relocation/alpha.sh"
qualify_mixed_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-972-fresh-beside-moved.md 2>&1)" \
  && fail_case 'qualify fresh-beside-moved case passed a genuinely new marker because a moved one shared its file'
printf '%s' "$qualify_mixed_output" | grep -q 'brand new leftover' \
  || fail_case 'qualify fresh-beside-moved case did not name the fresh marker in its FAIL'
printf '%s' "$qualify_mixed_output" | grep -q 'relocated, not added' \
  || fail_case 'qualify fresh-beside-moved case dropped the moved marker instead of reporting it alongside the FAIL'
rm -rf "$qualify_repo/relocation" "$qualify_repo/do-work/REQ-97"*.md

# qualify: a P-A-U section that keeps PLAN and APPLY but drops UNIFY is disarmed too. The
# artifact FAILs all key on a CHECKED [UNIFY] line, so a REQ with no UNIFY box at all
# satisfies them by absence exactly as a sectionless one does — counting all three boxes
# left this hole, and a PR review found it (REQ-264 remediation).
printf '%s\n' '## Implementation Summary' \
  '- `src/value_parser.py` (modified) — parsing tweak' \
  '' '## AI Execution State (P-A-U Loop)' \
  '- [x] **[PLAN]:** done' '- [x] **[APPLY]:** done' \
  > "$qualify_repo/do-work/REQ-961-no-unify-box.md"
printf 'print(raw_text)\n' >> "$qualify_repo/src/value_parser.py"
qualify_no_unify_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-961-no-unify-box.md 2>&1)"
printf '%s' "$qualify_no_unify_output" | grep -q 'DISARMED' \
  || fail_case 'qualify no-UNIFY-box case printed no disarmed-audit warning — a REQ keeping PLAN and APPLY but dropping UNIFY still passes Check 4 silently'
printf '%s' "$qualify_no_unify_output" | grep -q 'no \[UNIFY\] box' \
  || fail_case 'qualify no-UNIFY-box case did not distinguish a missing UNIFY box from a missing section — the remedies differ'
git -C "$qualify_repo" checkout -q -- src/value_parser.py

# qualify: ownership follows a RENAMED file back to its base path. A renamed file does not
# exist at the base under its new path, so a bare probe fell through to the post-change
# working copy and read the file as brand new — handing a rename that also adds an exit
# idiom and a debug print the very reporter exemption the base-revision rule denies
# (REQ-263 remediation, found by a PR review).
printf '%s\n' 'def render_row(raw_text):' '    return raw_text' > "$qualify_repo/src/row_renderer.py"
fixture_repo_commit_all "$qualify_repo" row-renderer
printf '%s\n' '## Implementation Summary' \
  '- `src/row_renderer_v2.py` (modified) — renamed and tweaked' \
  '' '## AI Execution State (P-A-U Loop)' \
  '- [x] **[PLAN]:** done' '- [x] **[APPLY]:** done' '- [x] **[UNIFY]:** done' \
  > "$qualify_repo/do-work/REQ-962-renamed-ownership.md"
git -C "$qualify_repo" mv src/row_renderer.py src/row_renderer_v2.py
printf '%s\n' 'print("debug", raw_text)' '' 'if True:' '    import sys' '    sys.exit(0)' \
  >> "$qualify_repo/src/row_renderer_v2.py"
qualify_renamed_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-962-renamed-ownership.md 2>&1)" \
  && fail_case 'qualify renamed-ownership case let a rename plus a same-change exit idiom excuse a debug print — ownership was judged on post-change content instead of the base path'
printf '%s' "$qualify_renamed_output" | grep -q 'src/row_renderer_v2.py' \
  || fail_case 'qualify renamed-ownership case did not name the renamed file'
git -C "$qualify_repo" checkout -q -- . 2>/dev/null || true
git -C "$qualify_repo" reset -q 2>/dev/null || true

# qualify: the relocation test compares WHOLE LINES, not substrings. Replacing a long marker
# with a short one that is a substring of it leaves a substring count unchanged on both
# sides, so a brand-new marker read as relocated and the WARN even claimed its exact text
# already existed. `git grep` has no -x, so the exactness comes from piping its fixed-string
# prefilter through real `grep -c -x -F` (REQ-301 remediation, found by a PR review).
printf '%s\n' 'def parse_total(raw_text):' '    # TODO remove the deprecated total parser' \
  '    return raw_text' > "$qualify_repo/src/total_parser.py"
fixture_repo_commit_all "$qualify_repo" total-parser
printf '%s\n' '## Implementation Summary' \
  '- `src/total_parser.py` (modified) — shorten the note' \
  '' '## AI Execution State (P-A-U Loop)' \
  '- [x] **[PLAN]:** done' '- [x] **[APPLY]:** done' '- [x] **[UNIFY]:** done' \
  > "$qualify_repo/do-work/REQ-973-substring-not-relocation.md"
printf '%s\n' 'def parse_total(raw_text):' '    # TODO' '    return raw_text' \
  > "$qualify_repo/src/total_parser.py"
qualify_substring_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-973-substring-not-relocation.md 2>&1)" \
  && fail_case 'qualify substring-not-relocation case excused a brand-new marker because its text is a substring of a line that was removed — the relocation test is counting substrings instead of whole lines'
printf '%s' "$qualify_substring_output" | grep -q 'debug artifacts' \
  || fail_case 'qualify substring-not-relocation case did not FAIL the fresh marker as a debug artifact'
git -C "$qualify_repo" checkout -q -- src/total_parser.py

# qualify: ownership resolves a rename whose destination is UNTRACKED. A plain `mv` (not
# `git mv`) leaves the old path a tracked deletion and the destination untracked, so it is in
# neither the working nor the staged diff and no rename can be paired from them. Without a
# post-change private index the moved file falls through to its own content and takes the
# reporter exemption the base-revision rule denies (REQ-263 remediation, second round, found
# by a PR review).
printf '%s\n' 'def render_cell(raw_text):' '    return raw_text' > "$qualify_repo/src/cell_renderer.py"
fixture_repo_commit_all "$qualify_repo" cell-renderer
printf '%s\n' '## Implementation Summary' \
  '- `src/cell_renderer_v2.py` (modified) — moved with plain mv and tweaked' \
  '' '## AI Execution State (P-A-U Loop)' \
  '- [x] **[PLAN]:** done' '- [x] **[APPLY]:** done' '- [x] **[UNIFY]:** done' \
  > "$qualify_repo/do-work/REQ-974-untracked-rename-ownership.md"
mv "$qualify_repo/src/cell_renderer.py" "$qualify_repo/src/cell_renderer_v2.py"
printf '%s\n' 'print("debug", raw_text)' '' 'if True:' '    import sys' '    sys.exit(0)' \
  >> "$qualify_repo/src/cell_renderer_v2.py"
qualify_untracked_rename_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-974-untracked-rename-ownership.md 2>&1)" \
  && fail_case 'qualify untracked-rename case let a plain mv plus a same-change exit idiom excuse a debug print — the rename destination was untracked, so ownership fell through to post-change content'
printf '%s' "$qualify_untracked_rename_output" | grep -q 'never ends its own process' \
  || fail_case 'qualify untracked-rename case did not reach the base-revision ownership verdict for an untracked rename destination'
# The private index must not disturb the real one: the move is still unstaged afterwards.
git -C "$qualify_repo" diff --cached --name-only | grep -q 'cell_renderer' \
  && fail_case 'qualify untracked-rename case left the fixture repo index modified — the rename probe must build its post-change tree in a PRIVATE index, never the real one'
rm -f "$qualify_repo/src/cell_renderer_v2.py"
git -C "$qualify_repo" checkout -q -- . 2>/dev/null || true

# qualify: the rename probe's private index must not add objects to the real repository
# either. `git add -A` writes a blob for every file it stages, so with only GIT_INDEX_FILE
# redirected this read-only probe left an unreachable copy of every untracked file inside
# the repository, surviving the temporary index it was staged into (found by a PR review).
printf '%s\n' 'def render_grid(raw_text):' '    return raw_text' > "$qualify_repo/src/grid_renderer.py"
fixture_repo_commit_all "$qualify_repo" grid-renderer
printf '%s\n' '## Implementation Summary' \
  '- `src/grid_renderer_v2.py` (modified) — moved with plain mv and tweaked' \
  '' '## AI Execution State (P-A-U Loop)' \
  '- [x] **[PLAN]:** done' '- [x] **[APPLY]:** done' '- [x] **[UNIFY]:** done' \
  > "$qualify_repo/do-work/REQ-975-private-object-database.md"
mv "$qualify_repo/src/grid_renderer.py" "$qualify_repo/src/grid_renderer_v2.py"
printf '%s\n' 'print("debug", raw_text)' '' 'if True:' '    import sys' '    sys.exit(0)' \
  >> "$qualify_repo/src/grid_renderer_v2.py"
# An untracked file the probe would stage: its blob is what used to leak into the repository.
awk 'BEGIN { for (line_number = 1; line_number <= 2000; line_number++)
  print "untracked-payload-line", line_number }' > "$qualify_repo/uncommitted-payload.txt"
qualify_object_count_before="$(find "$qualify_repo/.git/objects" -type f | wc -l | tr -d '[:space:]')"
(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-975-private-object-database.md >/dev/null 2>&1) || true
qualify_object_count_after="$(find "$qualify_repo/.git/objects" -type f | wc -l | tr -d '[:space:]')"
[ "$qualify_object_count_before" = "$qualify_object_count_after" ] \
  || fail_case "qualify private-object-database case left $((qualify_object_count_after - qualify_object_count_before)) new object(s) in the fixture repository — the rename probe must redirect GIT_OBJECT_DIRECTORY as well as GIT_INDEX_FILE, or every staged untracked file outlives the probe as an unreachable blob"
rm -f "$qualify_repo/uncommitted-payload.txt" "$qualify_repo/src/grid_renderer_v2.py"
git -C "$qualify_repo" checkout -q -- . 2>/dev/null || true

# qualify: the untracked artifact scans read text files only. A binary asset carrying the
# bytes `TODO` or `print(` is not source, and on greps that report a binary match on stdout
# (BSD, GNU before 3.5) it became an unreadable FAIL line and then reached the ownership
# probe, which reads the whole file into a variable. On GNU grep >= 3.5 the diagnostic goes
# to stderr instead, which the scans never redirect — so pre-fix it still surfaced in
# qualify's own output, naming a file the reader cannot inspect. Asserting the path appears
# nowhere in the combined output pins both flavors; the text half pins the guard against
# over-skipping real source.
printf '%s\n' '## Implementation Summary' \
  '- `src/new_helper.py` (new) — untracked helper carrying a debug print' \
  '' '## AI Execution State (P-A-U Loop)' \
  '- [x] **[PLAN]:** done' '- [x] **[APPLY]:** done' '- [x] **[UNIFY]:** done' \
  > "$qualify_repo/do-work/REQ-976-binary-untracked-scan.md"
printf 'def build_row(raw_text):\n    print(raw_text)\n    return raw_text\n' > "$qualify_repo/src/new_helper.py"
mkdir -p "$qualify_repo/assets"
printf 'ICON\000 TODO \000 print( \000\n' > "$qualify_repo/assets/thumbnail.bin"
qualify_binary_scan_output="$(cd "$qualify_repo" && "$core_checks/qualify.sh" do-work/REQ-976-binary-untracked-scan.md 2>&1)" \
  && fail_case 'qualify binary-untracked-scan case passed a debug print in an untracked text file — the binary guard must skip binaries only'
printf '%s' "$qualify_binary_scan_output" | grep -q 'src/new_helper.py' \
  || fail_case 'qualify binary-untracked-scan case stopped scanning untracked TEXT files — the binary guard is over-skipping real source'
printf '%s' "$qualify_binary_scan_output" | grep -q 'thumbnail.bin' \
  && fail_case 'qualify binary-untracked-scan case reported an untracked BINARY asset as a source artifact — the scans must skip binaries before reading them'
rm -rf "$qualify_repo/assets" "$qualify_repo/src/new_helper.py"
git -C "$qualify_repo" checkout -q -- . 2>/dev/null || true

# qualify: every path on a multi-path Implementation Summary bullet reaches Check 1.
# Matching first tokens made the old first-token parser pass this missing later path.
qualify_summary_repo="$fixture_root/qualify-summary-repo"
fixture_repo_init "$qualify_summary_repo"
mkdir -p "$qualify_summary_repo/src" "$qualify_summary_repo/do-work"
printf 'base\n' > "$qualify_summary_repo/src/first.txt"
printf 'base\n' > "$qualify_summary_repo/root-file.txt"
fixture_repo_commit_all "$qualify_summary_repo" base
printf 'change\n' >> "$qualify_summary_repo/src/first.txt"
printf '%s\n' '## Implementation Summary' \
  '- `src/first.txt`, `src/missing-later.txt` (modified) — paired updates' \
  '' '## AI Execution State (P-A-U Loop)' \
  '- [x] **[PLAN]:** done' '- [x] **[APPLY]:** done' '- [x] **[UNIFY]:** done' \
  > "$qualify_summary_repo/do-work/REQ-980-multi-path-mismatch.md"
qualify_multi_path_mismatch_output="$(cd "$qualify_summary_repo" && "$core_checks/qualify.sh" do-work/REQ-980-multi-path-mismatch.md 2>&1)" \
  && fail_case 'qualify multi-path mismatch case passed because only the first Implementation Summary path reached Check 1'
printf '%s' "$qualify_multi_path_mismatch_output" | grep -q 'listed (modified) but not on disk: src/missing-later.txt' \
  || fail_case 'qualify multi-path mismatch case did not name the missing later path'
git -C "$qualify_summary_repo" checkout -q -- .

# qualify: a matching multi-path list stays clean, and backticked spans on a prose-only
# bullet never become phantom file claims.
printf 'change\n' >> "$qualify_summary_repo/src/first.txt"
printf 'change\n' >> "$qualify_summary_repo/root-file.txt"
printf '%s\n' '## Implementation Summary' \
  '- `src/first.txt`, `root-file.txt` (modified) — paired updates' \
  '- Notes mention `phantom-file.txt` and `sort -u`, but claim no file.' \
  '' '## AI Execution State (P-A-U Loop)' \
  '- [x] **[PLAN]:** done' '- [x] **[APPLY]:** done' '- [x] **[UNIFY]:** done' \
  > "$qualify_summary_repo/do-work/REQ-981-multi-path-match.md"
qualify_multi_path_match_output="$(cd "$qualify_summary_repo" && "$core_checks/qualify.sh" do-work/REQ-981-multi-path-match.md 2>&1)" \
  || fail_case 'qualify matching multi-path case rejected two real modified paths'
printf '%s' "$qualify_multi_path_match_output" | grep -q 'phantom-file.txt' \
  && fail_case 'qualify matching multi-path case treated a backticked prose span as a file claim'
git -C "$qualify_summary_repo" checkout -q -- .

# qualify: root-level filenames are paths even without a slash. A slash-only filter would
# silently drop this missing second item and let the first path carry the check.
printf 'change\n' >> "$qualify_summary_repo/src/first.txt"
printf '%s\n' '## Implementation Summary' \
  '- `src/first.txt`, `missing-root.txt` (modified) — paired updates' \
  '' '## AI Execution State (P-A-U Loop)' \
  '- [x] **[PLAN]:** done' '- [x] **[APPLY]:** done' '- [x] **[UNIFY]:** done' \
  > "$qualify_summary_repo/do-work/REQ-982-root-path.md"
qualify_root_path_output="$(cd "$qualify_summary_repo" && "$core_checks/qualify.sh" do-work/REQ-982-root-path.md 2>&1)" \
  && fail_case 'qualify root-path case dropped a filename-only later path'
printf '%s' "$qualify_root_path_output" | grep -q 'listed (modified) but not on disk: missing-root.txt' \
  || fail_case 'qualify root-path case did not name the missing root-level filename'
git -C "$qualify_summary_repo" checkout -q -- .

# qualify: later (new) paths reach the wiring half as well as the existence check.
printf 'export const first = 1;\n' > "$qualify_summary_repo/src/first_new.js"
printf 'export const laterWidget = 1;\n' > "$qualify_summary_repo/src/later_widget.js"
printf '%s\n' '## Implementation Summary' \
  '- `src/first_new.js`, `src/later_widget.js` (new) — paired helpers' \
  '' '## AI Execution State (P-A-U Loop)' \
  '- [x] **[PLAN]:** done' '- [x] **[APPLY]:** done' '- [x] **[UNIFY]:** done' \
  > "$qualify_summary_repo/do-work/REQ-983-later-wiring.md"
qualify_later_wiring_output="$(cd "$qualify_summary_repo" && "$core_checks/qualify.sh" do-work/REQ-983-later-wiring.md 2>&1)" \
  || fail_case 'qualify later-wiring case rejected two real new paths'
printf '%s' "$qualify_later_wiring_output" | grep -q 'new) file has no static reference anywhere: src/later_widget.js' \
  || fail_case 'qualify later-wiring case never sent the later path through Check 5'
rm -f "$qualify_summary_repo/src/first_new.js" "$qualify_summary_repo/src/later_widget.js"

# qualify: an unmatched backtick on a path-led Summary bullet is malformed input, not a
# partial list whose already-valid first path may qualify the whole REQ.
printf 'change\n' >> "$qualify_summary_repo/src/first.txt"
printf '%s\n' '## Implementation Summary' \
  '- `src/first.txt`, `unterminated.txt (modified) — malformed pair' \
  '' '## AI Execution State (P-A-U Loop)' \
  '- [x] **[PLAN]:** done' '- [x] **[APPLY]:** done' '- [x] **[UNIFY]:** done' \
  > "$qualify_summary_repo/do-work/REQ-984-unmatched-summary.md"
qualify_unmatched_summary_output="$(cd "$qualify_summary_repo" && "$core_checks/qualify.sh" do-work/REQ-984-unmatched-summary.md 2>&1)" \
  && fail_case 'qualify unmatched-summary case accepted a partial Implementation Summary path list'
printf '%s' "$qualify_unmatched_summary_output" | grep -q 'FAIL:.*unmatched backtick' \
  || fail_case 'qualify unmatched-summary case did not fail loudly on the malformed path-led bullet'

prescribed_shell_finish
