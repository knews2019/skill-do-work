#!/usr/bin/env bash
# End-to-end probe for the fast gate's per-stage evidence reuse. It drives the
# shipped run_stage_with_evidence out of maintainer-verify.sh against a synthetic
# repository and the real do-work-cli decision, revocation and recording
# commands. Only the stage command itself is substituted, so what runs is
# observable: the wrapper is the boundary under test, and the budget runner it
# calls has its own probe in run-go-tests-with-budget-behavior.sh.
set -uo pipefail

# A test of the reuse wrapper must not read the caller's preference for that
# wrapper. Both of these names short-circuit run_stage_with_evidence in
# maintainer-verify.sh before it decides anything, and the measurement protocol
# runs the whole gate as DO_WORK_FAST_STAGE_REUSE=off, so an inherited value
# would switch off the exact boundary under test: every case below would then
# report the caller's environment as nine failures of the code. The probe clears
# both names for its own runs instead of inheriting either.
unset DO_WORK_FAST_STAGE_REUSE MAINTAINER_VERIFY_SELFTEST_LOG

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
gate_script="${DO_WORK_TEST_MAINTAINER_VERIFY:-$repo_root/_dev/tests/maintainer-verify.sh}"
probe_temporary_directory="${TMPDIR:-/tmp}"
# TMPDIR routinely ends in a slash here, and a doubled separator survives every
# path concatenation while comparing unequal to the same path taken from `pwd`.
probe_fixture_root="$(mktemp -d "${probe_temporary_directory%/}/fast-stage-reuse.XXXXXX")" || exit 1
trap 'rm -rf -- "$probe_fixture_root"' EXIT
failure_count=0

project_root="$probe_fixture_root/project"
stage_run_log="$probe_fixture_root/stage-runs.log"
mkdir -p \
  "$project_root/_dev/tests" \
  "$project_root/skills/do-work/tools" \
  "$project_root/module-alpha" \
  "$project_root/do-work/queue" \
  "$project_root/do-work/runs"

# alpha-stage covers do-work/ the way the shipped queue-kanban stage does: that
# stage builds its board from the real tree, so a change there must force it.
# The stage's own toolchain probe reads a file under do-work/runs/, which the
# manifest excludes from every seal, so the probe cases move the probe output and
# nothing else.
cat > "$project_root/_dev/tests/fast-stages.json" <<'MANIFEST'
{
  "schema_version": 1,
  "stages": [
    {
      "id": "alpha-stage",
      "argv": ["run-go-tests-with-budget.sh", "module-alpha", "./..."],
      "coverage": [
        {"kind": "subtree", "path": "module-alpha"},
        {"kind": "subtree", "path": "do-work"}
      ],
      "fingerprint": {
        "toolchain_probes": [["cat", "do-work/runs/toolchain.txt"]]
      }
    }
  ],
  "non_stage_coverage": [],
  "seal_exclusions": [
    {"kind": "subtree", "path": "do-work/runs"},
    {"kind": "exact", "path": "do-work/test-durations.tsv"}
  ]
}
MANIFEST
printf 'alpha source 1\n' > "$project_root/module-alpha/source.txt"
printf 'toolchain 1\n' > "$project_root/do-work/runs/toolchain.txt"
printf -- '---\nid: REQ-001\n---\n' > "$project_root/do-work/queue/REQ-001-fixture.md"
# The wrapper reaches the CLI through the skill's canonical launcher path. The
# shim forwards to the real launcher in this repository so the command under
# test is the real one, and so the fixture does not pay a fresh Go link for a
# copied module directory.
printf '#!/usr/bin/env bash\nexec bash %q "$@"\n' \
  "$repo_root/skills/do-work/tools/do-work-cli.sh" > "$project_root/skills/do-work/tools/do-work-cli.sh"
git -C "$project_root" init -q
git -C "$project_root" config user.email fixture@example.invalid
git -C "$project_root" config user.name Fixture
git -C "$project_root" add -A
git -C "$project_root" commit -qm 'fast-stage reuse fixture'

# The gate script ends in a command dispatch that would run a whole verification
# on source. Strip it, and refuse to continue if the strip matched nothing —
# a copy that still dispatches would run the real gate from inside this probe.
sed '/^if \[ "\$#" -gt 2 \]; then$/,$d' "$gate_script" > "$project_root/_dev/tests/maintainer-verify.sh"
if cmp -s "$gate_script" "$project_root/_dev/tests/maintainer-verify.sh"; then
  printf 'FAIL: the gate command dispatch was not found; the probe cannot source the gate safely.\n' >&2
  exit 1
fi
# shellcheck source=/dev/null
source "$project_root/_dev/tests/maintainer-verify.sh"
# The sourced gate turns on errexit. This probe reports every failing case in one
# run instead of stopping at the first, so restore its own options.
set +e
set -uo pipefail

stage_exit_status=0
stage_behavior=normal
stage_output=''
stage_status=0
stage_ran=no

# Substitutes the stage command with one whose execution is observable. It is
# called exactly where the budget runner is called in the real gate.
run_budgeted_go_tests() {
  printf 'ran %s\n' "$1" >> "$stage_run_log"
  if [ "$stage_behavior" = 'interrupted' ]; then
    # A real terminating signal, so the status the wrapper returns is the status
    # a killed stage really produces. The child dies immediately and owns no
    # lifetime beyond this call.
    bash -c 'kill -TERM $$'
    return
  fi
  return "$stage_exit_status"
}

run_stage() {
  : > "$stage_run_log"
  stage_status=0
  stage_output="$(run_stage_with_evidence alpha-stage "$project_root/module-alpha" ./... 2>&1)" || stage_status=$?
  stage_ran=no
  if [ -s "$stage_run_log" ]; then
    stage_ran=yes
  fi
}

expect_case() {
  local case_name="$1" expect_ran="$2" expect_status="$3" expect_line="$4"
  if [ "$stage_ran" != "$expect_ran" ] || [ "$stage_status" -ne "$expect_status" ] || \
    ! grep -qF -- "$expect_line" <<<"$stage_output"; then
    # The wanted line is printed beside the output because all three fields are
    # asserted and only one of them may be wrong: a message that shows a
    # matching ran, a matching status and an empty output reads as a passing
    # case, and hides that the missing disposition line is the whole failure.
    printf 'FAIL: %s ran=%s (want %s) status=%s (want %s) output=<%s> want-line=<%s>\n' \
      "$case_name" "$stage_ran" "$expect_ran" "$stage_status" "$expect_status" \
      "$(tr '\n' ' ' <<<"$stage_output")" "$expect_line" >&2
    failure_count=$((failure_count + 1))
  fi
}

change_covered_input() {
  printf 'alpha source %s\n' "$1" > "$project_root/module-alpha/source.txt"
}

run_stage
expect_case 'first run executes' yes 0 'EXECUTING (no_prior_evidence)'

run_stage
expect_case 'unchanged inputs reuse' no 0 'REUSED (fingerprint_match, recorded '

change_covered_input 2
run_stage
expect_case 'a covered input change executes' yes 0 'EXECUTING (fingerprint_mismatch)'

# The failure this catches: a do-work/-only edit reused stale evidence, so the
# gate printed "Maintainer verification passed." and exited 0 while the stage's
# own test failed on that same tree.
printf -- '---\nid: REQ-002\n---\n' > "$project_root/do-work/queue/REQ-002-fixture.md"
run_stage
expect_case 'a queue-tree change executes the stage that reads it' yes 0 'EXECUTING (fingerprint_mismatch)'

# The failure this catches: the stage appends this log itself while it runs, and
# the recorded fingerprint is the pre-run one, so a seal over the log could never
# match the evidence it authorized and the stage would never reuse again.
printf 'probe\tprobe\t0.00\t0\n' >> "$project_root/do-work/test-durations.tsv"
run_stage
expect_case "the gate's own duration log alone still reuses" no 0 'REUSED (fingerprint_match, recorded '

change_covered_input 3
stage_exit_status=9
run_stage
expect_case 'a failing stage returns its exact status' yes 9 'EXECUTING (fingerprint_mismatch)'

stage_exit_status=0
run_stage
expect_case 'a failed stage left no reusable evidence' yes 0 'EXECUTING (no_prior_evidence)'

change_covered_input 4
stage_behavior=interrupted
run_stage
expect_case 'an interrupted stage returns its signal status' yes 143 'EXECUTING (fingerprint_mismatch)'

stage_behavior=normal
run_stage
expect_case 'an interrupted stage left no reusable evidence' yes 0 'EXECUTING (no_prior_evidence)'

mv "$project_root/_dev/tests/fast-stages.json" "$probe_fixture_root/fast-stages.json.aside"
run_stage
expect_case 'an unusable decision executes rather than skips' yes 0 'EXECUTING (decision_unavailable)'
mv "$probe_fixture_root/fast-stages.json.aside" "$project_root/_dev/tests/fast-stages.json"

[ "$failure_count" -eq 0 ] || exit 1
printf 'Fast-stage evidence reuse probes passed.\n'
