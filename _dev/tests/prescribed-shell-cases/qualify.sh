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

prescribed_shell_finish
