#!/usr/bin/env bash
# REQ-656 GREEN probe: ai-report revise is documented as a new sibling bundle that supersedes
# the prior one (revise form, -rev<N> name, rev block, supersedes meta, catalog regeneration,
# prompt-injection load, Collision-Safe Publication), the table-of-contents rule sits in Report
# Design Rules, the path-last rule sits in Output Format, routing/help/guide name the form, and
# the immutability rule stays byte-identical. Prose REQ: the assertions are the fixed tokens the
# builder brief names, plus the shipped-package reference contract. Reads paths relative to the
# tree it runs in. Expected to FAIL at base 7fbeff30 (no revise form exists yet).
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
toolbox="$repo_root/skills/do-work-toolbox"
failures=0

fail() {
  printf 'FAIL: %s\n' "$1" >&2
  failures=$((failures + 1))
}

# section_text <file text> <exact heading line>: the heading line and its body, up to the next
# heading of the same or a higher level.
section_text() {
  local heading_level="${2%% *}"
  awk -v heading="$2" -v level="${#heading_level}" '
    $0 == heading { inside = 1; print; next }
    inside && /^#+ / { split($0, marks, " "); if (length(marks[1]) <= level) { inside = 0 } }
    inside { print }
  ' <<<"$1"
}

# require_in <label> <text> <fixed token>
require_in() {
  grep -Fq -- "$3" <<<"$2" || fail "$1 lacks: $3"
}

action_text="$(cat "$toolbox/actions/ai-report.md")"
reference_text="$(cat "$toolbox/actions/ai-report-reference.md")"
guide_text="$(cat "$toolbox/docs/ai-report-guide.md")"
skill_text="$(cat "$toolbox/SKILL.md")"
help_text="$(cat "$toolbox/actions/help.md")"
shared_reference_text="$(cat "$toolbox/actions/completed-work-presentation-reference.md")"

# The action: one revise subsection holding the whole procedure.
require_in 'ai-report.md' "$action_text" 'ai-report revise <dir|latest> [what changed]'
revise_section="$(section_text "$action_text" '### Revise form')"
[ -n "$revise_section" ] || fail 'ai-report.md lacks the heading: ### Revise form'
for required_token in \
  'latest' \
  'catalog.json' \
  'prompt-injection.md' \
  'Collision-Safe Publication' \
  '-rev<N>' \
  'supersedes' \
  '<meta name="ai-report-supersedes" content="<prior folder>">' \
  'ai-report-kind' \
  'rev-N' \
  'Changed' \
  'Still to do' \
  'superseded_by' \
  'Report Design Rules'
do
  require_in 'revise subsection' "$revise_section" "$required_token"
done

# Path-last rule: every invocation that writes a bundle ends with the path and a file:// link.
output_section="$(section_text "$action_text" '## Output Format')"
require_in 'ai-report.md Output Format' "$output_section" 'file://'

# Reference: the length-keyed table-of-contents rule inside Report Design Rules.
design_rules_section="$(section_text "$reference_text" '## Report Design Rules (Step 5)')"
require_in 'Report Design Rules' "$design_rules_section" 'table of contents'
require_in 'Report Design Rules' "$design_rules_section" 'six top-level sections'
require_in 'Report Design Rules' "$design_rules_section" '1,500 words'

# User-facing surfaces: guide, routing row, help.
require_in 'ai-report-guide.md' "$guide_text" 'ai-report revise'
routing_row="$(grep -F -- './actions/ai-report.md' <<<"$skill_text")" || true
for routing_phrase in 'revise the report' 'update the report'; do
  require_in 'SKILL.md ai-report routing rows' "$routing_row" "\`$routing_phrase\`"
done
require_in 'help.md' "$help_text" 'ai-report revise'

# The immutability rule and the sibling rule stay byte-identical (REQ constraint).
grep -Fxq -- 'Before creating any output directory or file, derive the complete final path and test whether it already exists. Existing artifacts are immutable: never delete, truncate, merge into, rename, migrate, or overwrite them.' <<<"$shared_reference_text" \
  || fail 'completed-work reference immutability line changed'
# shellcheck disable=SC2016
grep -Fxq -- 'If the preferred path exists, choose a fresh sibling by appending the first available numeric suffix (`-2`, `-3`, and so on), then use that one path consistently for the whole artifact. Create directories only after the collision check. A failed or partial run never grants permission to reuse an existing path; report the partial path so the user can inspect it.' <<<"$shared_reference_text" \
  || fail 'completed-work reference sibling-suffix line changed'

# Every relative link and path the shipped files name still resolves.
reference_contract_output="$(bash "$repo_root/_dev/tests/shipped-package-reference-contract.sh" 2>&1)" \
  || fail "shipped-package-reference-contract.sh failed: $reference_contract_output"

if [ "$failures" -ne 0 ]; then
  printf 'REQ-656 probe: %s failure(s)\n' "$failures" >&2
  exit 1
fi
printf 'REQ-656 probe passed.\n'
