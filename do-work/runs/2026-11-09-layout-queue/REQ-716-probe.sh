#!/usr/bin/env bash
# GREEN probe for REQ-716: the five root installer copies in tools/ are deleted (maintainer's
# choice, D-04); the bootstrap runs the staged installer in skills/do-work/tools/.
# Checks whichever tree it runs in. Expected to FAIL at base 792f981c (root tools/ still exists).
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root"

fail() { printf 'REQ-716 probe FAIL: %s\n' "$*" >&2; exit 1; }
script_names='install-do-work-suite fetch-upstream-archive do-work-cli-preamble replace-text-section validate-suite-manifest'

# 1. Root tools/ is gone, no copy moved to cmd/, and the five staged scripts remain tracked.
[ ! -e tools ] || fail 'a root tools/ path still exists'
tracked_under_tools="$(git ls-files -- tools)"
[ -z "$tracked_under_tools" ] || fail "tracked files remain under root tools/: $tracked_under_tools"
for script_name in $script_names; do
  [ ! -e "cmd/$script_name.sh" ] || fail "cmd/$script_name.sh exists; the maintainer chose deletion, not a move"
  git ls-files --error-unmatch -- "skills/do-work/tools/$script_name.sh" >/dev/null 2>&1 \
    || fail "skills/do-work/tools/$script_name.sh is not tracked"
  # D-05: the root-copy fallback (a sibling skills/do-work/tools/ one level up) is dead code now.
  if grep -qF '../skills/do-work/tools/' "skills/do-work/tools/$script_name.sh"; then
    fail "skills/do-work/tools/$script_name.sh still probes the deleted root-copy location"
  fi
done

# 2. The printed bootstrap runs the staged installer, and the README block is byte-identical to it.
bootstrap_command="$(bash skills/do-work/tools/install-do-work-suite.sh --print-bootstrap-command)"
# shellcheck disable=SC2016 # the literal $bootstrap_tmp text is what the bootstrap must contain.
grep -qF '"$bootstrap_tmp/source/skills/do-work/tools/install-do-work-suite.sh"' <<<"$bootstrap_command" \
  || fail 'the printed bootstrap does not run source/skills/do-work/tools/install-do-work-suite.sh'
if grep -qF 'source/tools/' <<<"$bootstrap_command"; then fail 'the printed bootstrap still names source/tools/'; fi
readme_bootstrap_command="$(awk '
  /^## Installation$/ {in_installation=1; next}
  in_installation && /^```bash$/ {in_command=1; next}
  in_command && /^```$/ {exit}
  in_command {print}
' README.md)"
[ "$readme_bootstrap_command" = "$bootstrap_command" ] || fail 'README installation block differs from --print-bootstrap-command output'

# 3. folder-structure-check ran (its summary line is present) and no longer reports root tools/.
layout_output="$(bash skills/do-work/tools/do-work-cli.sh --repo-root "$root" --format text folder-structure-check 2>&1)" || true
grep -qE '^folder-structure-check: ' <<<"$layout_output" || fail "folder-structure-check did not run: $layout_output"
if grep -qE 'FOLDER-STRUCTURE-ROOT-ENTRY +tools/' <<<"$layout_output"; then fail 'folder-structure-check still reports root tools/'; fi

# 4. No maintainer surface still names the root copies: the mirror check is deleted, the heavy
#    lanes and the release guard name no root tools/ path.
# shellcheck disable=SC2016 # the literal shell text of the deleted mirror loop.
if grep -qF '"$repo_root"/tools/*.sh' _dev/tests/staged-skills-contract.sh; then fail 'staged-skills-contract.sh still runs the root mirror check'; fi
if grep -qE '"path": "tools/' _dev/tests/heavy-lanes.json; then fail 'heavy-lanes.json still names a root tools/ path'; fi
if grep -qF '"tools"' skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go; then
  fail 'the release guard still counts root tools as a shipped root'
fi
go_output="$(cd skills/do-work/tools/do-work-cli && go test -count=1 -v \
  -run '^TestReleaseRefusedWhenTheImplementationShipsNothing$' ./internal/finalization/ 2>&1)" || fail "release-guard test failed: $go_output"
grep -qF -- '--- PASS: TestReleaseRefusedWhenTheImplementationShipsNothing ' <<<"$go_output" \
  || fail "release-guard test did not PASS: $go_output"
go_output="$(cd skills/do-work/tools/do-work-cli && go test -count=1 -v \
  -run '^TestRepositoryManifestNamesEveryLaneScopedMaintainerEntryPoint$' ./internal/heavyverification/ 2>&1)" || fail "heavy-lane coverage test failed: $go_output"
grep -qF -- '--- PASS: TestRepositoryManifestNamesEveryLaneScopedMaintainerEntryPoint ' <<<"$go_output" \
  || fail "heavy-lane coverage test did not PASS: $go_output"

# 5. The fast contract suites that read these scripts pass.
suite_output="$(bash _dev/tests/suite-manifest-contract.sh 2>&1)" || fail "suite-manifest-contract failed: $suite_output"
grep -qF 'Suite manifest contract probes passed.' <<<"$suite_output" || fail 'suite-manifest-contract printed no pass line'
suite_output="$(bash _dev/tests/contracts/replace-text-section.sh 2>&1)" || fail "replace-text-section contract failed: $suite_output"
grep -qF 'replace-text-section contract probes passed.' <<<"$suite_output" || fail 'replace-text-section contract printed no pass line'
suite_output="$(bash _dev/tests/audit-lockins.sh 2>&1)" || fail "audit-lockins failed: $suite_output"
grep -qF 'Audit lock-in regressions passed.' <<<"$suite_output" || fail 'audit-lockins printed no pass line'

printf 'REQ-716 probe PASS\n'
