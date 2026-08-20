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

prescribed_shell_finish
