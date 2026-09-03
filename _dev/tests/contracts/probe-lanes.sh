#!/usr/bin/env bash
# Register standalone behavior files. probe-batch.sh owns execution and collection.
# shellcheck disable=SC2154 # This file is sourced by contract-regressions.sh.

register_probe() {
  local probe_name="$1"
  local probe_script="$2"
  local failure_message="$3"
  if [ ! -x "$probe_script" ]; then
    printf 'FAIL: %s is missing or not executable.\n' "${probe_script#"$repo_root/"}" >&2
    fail_count=$((fail_count + 1))
    return
  fi
  launch_probe "$probe_name" "$failure_message" "$probe_script"
}

register_probe suite_manifest_probe "$repo_root/_dev/tests/suite-manifest-contract.sh" \
  'suite manifest contract probes failed (see the FAIL lines above).'
register_probe shipped_reference_probe "$repo_root/_dev/tests/shipped-package-reference-contract.sh" \
  'shipped package reference contract failed (see the FAIL lines above).'
register_probe action_shell_probe "$repo_root/_dev/tests/action-shell-blocks.sh" \
  'shipped shell-block lint failed (see the attributed FAIL lines above).'
register_probe session_start_probe "$repo_root/_dev/tests/session-start-hook-behavior.sh" \
  'SessionStart hook behavior probes failed (see the fixture FAIL lines above).'
register_probe prescribed_shell_probe "$repo_root/_dev/tests/prescribed-shell-canonicalization.sh" \
  'prescribed shell primitive canonicalization failed (see the attributed FAIL lines above).'
register_probe defensive_surface_probe "$repo_root/_dev/tests/defensive-surface-audit.sh" \
  'defensive-surface exact deletion regression failed (see the attributed FAIL lines above).'
register_probe do_work_cli_launcher_probe "$repo_root/_dev/tests/do-work-cli-launcher-behavior.sh" \
  'do-work-cli launcher behavior probes failed (see the FAIL lines above).'
register_probe p50_estimator_probe "$repo_root/_dev/tests/p50-estimator-determinism.sh" \
  'P50 estimator probes failed (see the FAIL lines above).'
register_probe select_simple_reqs_probe "$repo_root/_dev/tests/select-simple-reqs-behavior.sh" \
  'cheaper-model selector probes failed (see the FAIL lines above).'
register_probe go_test_budget_probe "$repo_root/_dev/tests/run-go-tests-with-budget-behavior.sh" \
  'Go test budget behavior probes failed (see the FAIL lines above).'

if [ "$verification_tier" = heavy ]; then
  shared_cli_binary="$probe_batch_root/do-work-cli"
  if ! (
    cd "$repo_root/skills/do-work/tools/do-work-cli"
    go build -ldflags='-s -w' -o "$shared_cli_binary" ./cmd/do-work-cli
  ); then
    printf 'FAIL: could not build the shared heavy-probe do-work-cli.\n' >&2
    fail_count=$((fail_count + 1))
  else
    export DO_WORK_TEST_DO_WORK_CLI_BINARY="$shared_cli_binary"
    register_probe staged_skills_probe "$repo_root/_dev/tests/staged-skills-contract.sh" \
      'staged skills contract probes failed (see the FAIL lines above).'
    register_probe update_script_probe "$repo_root/_dev/tests/update-script-behavior.sh" \
      'update-script behavior probes failed (see the FAIL lines above).'
    register_probe suite_installer_probe "$repo_root/_dev/tests/install-suite-behavior.sh" \
      'suite installer behavior probes failed (see the FAIL lines above).'
  fi
else
  printf 'SKIP: staged skills, updater, and installer probes require maintainer-verify.sh --heavy.\n'
fi
