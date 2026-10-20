#!/usr/bin/env bash
# GREEN checks for REQ-703 (integrator test gate): second opinions accept --preset and --editor-direction, a directed
# page gets its preset/slug name, meta tags and header block, a plain run keeps today's name, and the report catalog
# skips the project preset folder. Also checks the shipped prompt placeholder, the reduced-motion rule and the two
# shipped presets. Names only this REQ's three tests; a skip is not a pass.
# Expected to FAIL at the pre-build baseline (the tests and shipped files do not exist yet). Under 30 s on a quiet machine.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module="$root/skills/do-work/tools/do-work-cli"
docs="$root/skills/do-work-toolbox/docs"
prompt_text="$(cat "$docs/ai-report-second-opinion-prompt.md")"
grep -q -- '{{EDITOR_DIRECTION}}' <<<"$prompt_text" || { echo "shipped prompt has no {{EDITOR_DIRECTION}} placeholder"; exit 1; }
grep -q -- 'prefers-reduced-motion' <<<"$prompt_text" || { echo "shipped prompt rule 8 does not mention prefers-reduced-motion"; exit 1; }
grep -q -- 'Variation slug:' <<<"$prompt_text" || { echo "shipped prompt FINAL REPLY does not ask for a Variation slug line"; exit 1; }
for preset in bold visual; do
  preset_file="$docs/second-opinion-presets/$preset.md"
  [ -s "$preset_file" ] || { echo "shipped preset missing or empty: $preset_file"; exit 1; }
  preset_text="$(cat "$preset_file")"
  if grep -q -- '{{' <<<"$preset_text"; then echo "shipped preset contains a {{ placeholder: $preset_file"; exit 1; fi
done
probe_tests=(TestSecondOpinionPresetResolutionAndArgumentRules TestSecondOpinionDirectedRunNamesAndMarksPage TestAIReportIndexSkipsSecondOpinionPresetsFolder)
pattern="^($(IFS='|'; echo "${probe_tests[*]}"))\$"
output="$(go test -C "$module" -count=1 -v -run "$pattern" ./internal/toolboxcommands 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for probe_test in "${probe_tests[@]}"; do
  grep -q -- "^--- PASS: $probe_test " <<<"$output" || { printf '%s\n' "$output"; echo "$probe_test did not pass (skipped or missing)"; exit 1; }
done
echo "REQ-703 probe: shipped prompt and presets checked; preset, directed-run and catalog-skip tests passed"
