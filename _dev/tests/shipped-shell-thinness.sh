#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
command_map="$repo_root/_dev/tests/fixtures/shipped-shell-command-map.tsv"
failure_count=0
mapped_count=0

while IFS=$'\t' read -r retained_path canonical_command _; do
  [[ "$retained_path" == \#* || -z "$retained_path" ]] && continue
  mapped_count=$((mapped_count + 1))
  retained_file="$repo_root/$retained_path"
  if [[ ! -x "$retained_file" ]]; then
    printf 'FAIL: retained shell path is missing or not executable: %s\n' "$retained_path" >&2
    failure_count=$((failure_count + 1))
    continue
  fi
  if ! grep -Fq '# do-work-cli compatibility launcher:' "$retained_file"; then
    printf 'FAIL: retained shell path still contains a legacy implementation: %s\n' "$retained_path" >&2
    failure_count=$((failure_count + 1))
  fi
  if grep -Eiq '(^|[^[:alnum:]_])(python([23])?|jq)([^[:alnum:]_]|$)' "$retained_file"; then
    printf 'FAIL: retained launcher contains a Python/jq implementation dependency: %s\n' "$retained_path" >&2
    failure_count=$((failure_count + 1))
  fi
  if [[ "$canonical_command" != '@launcher' && "$retained_path" != */install-do-work-suite.sh ]] \
    && ! grep -Fq 'do-work-cli.sh' "$retained_file"; then
    printf 'FAIL: retained path does not route through do-work-cli.sh: %s\n' "$retained_path" >&2
    failure_count=$((failure_count + 1))
  fi
done < "$command_map"

if [[ "$mapped_count" -ne 41 ]]; then
  printf 'FAIL: shipped shell map has %s entries, expected 41.\n' "$mapped_count" >&2
  failure_count=$((failure_count + 1))
fi
tracked_count="$(git -C "$repo_root" ls-files 'skills/**/*.sh' 'tools/*.sh' | wc -l | tr -d ' ')"
if [[ "$tracked_count" -ne "$mapped_count" ]]; then
  printf 'FAIL: shipped shell map has %s entries but Git tracks %s shipped paths.\n' "$mapped_count" "$tracked_count" >&2
  failure_count=$((failure_count + 1))
fi

[[ "$failure_count" -eq 0 ]] || exit 1
printf 'Shipped shell thinness checks passed for %s retained paths.\n' "$mapped_count"
