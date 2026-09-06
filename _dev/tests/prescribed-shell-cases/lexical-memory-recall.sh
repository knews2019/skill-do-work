#!/usr/bin/env bash
# Fixture execution proofs for lexical-memory-recall.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

# lexical-memory-recall: apostrophes and command syntax remain data while attribution is emitted.
memory_root="$fixture_root/memory with space"
mkdir -p "$memory_root/logs"
printf '%s\n' '# Memory' '## Decisions' 'Use cobalt release trains.' > "$memory_root/working-memory.md"
recall_output="$($knowledge_scripts/lexical-memory-recall.sh "$memory_root" "cobalt'; touch $fixture_root/injected #")" || fail_case 'lexical-memory-recall raw-query case returned nonzero'
[ ! -e "$fixture_root/injected" ] || fail_case 'lexical-memory-recall raw-query case executed query text'
grep -q $'working memory\t## Decisions\tUse cobalt' <<<"$recall_output" || fail_case 'lexical-memory-recall raw-query case omitted attribution'

prescribed_shell_finish
