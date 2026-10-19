#!/usr/bin/env bash
# REQ-699 focused probe: the REQ's acceptance greps plus the two shipped-file contract tests.
set -u
root="${DO_WORK_PROBE_ROOT:-$(pwd)}"
cd "$root" || exit 2
f=skills/do-work-toolbox/actions/ai-report.md
g=skills/do-work-toolbox/docs/ai-report-guide.md
fail=0
if grep -n "Do not commit" "$f"; then echo "FAIL: 'Do not commit' still present"; fail=1; fi
grep -q "Commit as the project's rules say" "$f" || { echo "FAIL: step 7 rule missing"; fail=1; }
grep -q "one finding per" "$f" || { echo "FAIL: revise step 5 plural missing"; fail=1; }
if grep -n -i "link finding remains\|only finding is the expected link\|only the expected prior-bundle link finding)" "$f"; then echo "FAIL: singular link-finding phrase left"; fail=1; fi
if grep -n "left to you" "$g"; then echo "FAIL: guide still says committing is left to you"; fail=1; fi
bash _dev/tests/shipped-package-reference-contract.sh || fail=1
bash _dev/tests/contract-regressions.sh || fail=1
echo "probe exit $fail"
exit $fail
