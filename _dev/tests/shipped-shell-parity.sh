#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
command_map="$repo_root/_dev/tests/fixtures/shipped-shell-command-map.tsv"
cli_launcher="$repo_root/skills/do-work/tools/do-work-cli.sh"
failure_count=0

# The pre-cutover per-script cases are the behavioral oracle. They assert exact statuses,
# streams, Git/index state, private-state cleanup, and filesystem bytes. Running them after
# cutover proves the canonical command still presents the retained shell contract.
bash "$repo_root/_dev/tests/prescribed-shell-scripts-behavior.sh"
bash "$repo_root/_dev/tests/shipped-shell-thinness.sh"

while IFS=$'\t' read -r retained_path canonical_command fixture_owner; do
  [[ "$retained_path" == \#* || -z "$retained_path" ]] && continue
  if [[ ! -f "$repo_root/$fixture_owner" ]]; then
    printf 'FAIL: %s has no preserved legacy fixture owner: %s\n' "$retained_path" "$fixture_owner" >&2
    failure_count=$((failure_count + 1))
  fi
  [[ "$canonical_command" == '@launcher' ]] && continue
  command_output="$(bash "$cli_launcher" --format json "$canonical_command" --__parity-registration-probe__ 2>&1 || true)"
  if grep -Fq 'unknown command' <<<"$command_output"; then
    printf 'FAIL: %s maps to unregistered canonical command %s.\n' "$retained_path" "$canonical_command" >&2
    failure_count=$((failure_count + 1))
  fi
done < "$command_map"

[[ "$failure_count" -eq 0 ]] || exit 1
printf 'Shipped shell parity checks passed against preserved legacy fixtures.\n'
