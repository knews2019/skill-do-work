#!/usr/bin/env bash
# REQ-687 GREEN probe: ai-report --kind proposal|root-cause is documented with its fixed
# shape, slug, meta tag, mockup label and capture lines, while the default completed-work
# gate stays byte-identical. Reads paths relative to the tree it runs in. Prose REQ: the
# assertions are fixed tokens the builder brief names, plus the shipped-reference contract.
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
toolbox="$repo_root/skills/do-work-toolbox"
failures=0

fail() {
  printf 'FAIL: %s\n' "$1" >&2
  failures=$((failures + 1))
}

# require_text <file> <fixed text> — the file contains the text.
require_text() {
  local file_text
  file_text="$(cat "$1")"
  grep -Fq -- "$2" <<<"$file_text" || fail "${1#"$repo_root/"} lacks: $2"
}

action_file="$toolbox/actions/ai-report.md"
reference_file="$toolbox/actions/ai-report-reference.md"
shared_reference_file="$toolbox/actions/completed-work-presentation-reference.md"
guide_file="$toolbox/docs/ai-report-guide.md"

# The action: the flag, both kinds, the default gate kept, and the shell-guide pointer that
# prescribed-shell-canonicalization.sh requires.
require_text "$action_file" '--kind'
require_text "$action_file" 'root-cause'
require_text "$action_file" 'Proposal and Root-Cause Kinds'
require_text "$action_file" 'unfinished or unsuccessful'
require_text "$action_file" '../../do-work/docs/prescribed-shell-primitives.md'

# The reference: one section holds every rule of the two kinds.
require_text "$reference_file" '## Proposal and Root-Cause Kinds'
require_text "$reference_file" '<skill-root>/../do-work/tools/do-work-cli.sh'
reference_text="$(cat "$reference_file")"
kinds_section="$(awk '/^## Proposal and Root-Cause Kinds/{inside=1; print; next} inside && /^## /{inside=0} inside' <<<"$reference_text")"
for required_token in \
  'Terminal-Success Target Resolution' \
  'Safety Load Order' \
  'Collision-Safe Publication' \
  'smallest-change option' \
  'MOCKUP — proposal' \
  'generated/' \
  '<meta name="ai-report-kind" content="<kind>">' \
  'yyyy-mm-dd_hhmm_<kind>-<topic>' \
  'do-work capture-request:'
do
  grep -Fq -- "$required_token" <<<"$kinds_section" || fail "kinds section lacks: $required_token"
done

# Part order: each template names its parts in bold, in the REQ's order. Compares the
# character offset of each part's first occurrence inside the section.
check_part_order() {
  local template_label="$1" previous_offset=-1 part_name part_needle text_before
  shift
  for part_name in "$@"; do
    part_needle="**$part_name**"
    text_before="${kinds_section%%"$part_needle"*}"
    if [ "$text_before" = "$kinds_section" ]; then
      fail "$template_label template lacks part $part_needle"
      return
    fi
    if [ "${#text_before}" -le "$previous_offset" ]; then
      fail "$template_label template lists $part_needle out of order"
    fi
    previous_offset="${#text_before}"
  done
}
check_part_order proposal 'Decision' 'Options' 'Recommendation' 'Evidence ledger' 'Limits' 'Open questions' 'Capture lines'
check_part_order root-cause 'What happened' 'Why' 'What to change'

# The default completed-work gate stays byte-identical (REQ constraint; its lines 20 and 30).
# The backticks below are literal Markdown, not command substitution.
shared_reference_text="$(cat "$shared_reference_file")"
# shellcheck disable=SC2016
grep -Fxq -- 'Normalize every status under `../../do-work/actions/work-reference.md` → **Schema Read Contract**, then apply its **Terminal-success status set**. Accept `completed` and `completed-with-issues`; the latter is successful but its recorded issues must remain visible in the presentation. Reject `cancelled`, `failed`, and every unfinished status.' <<<"$shared_reference_text" \
  || fail 'completed-work reference terminal-success line changed'
# shellcheck disable=SC2016
grep -Fxq -- 'Never fall back to `do-work/queue/`, `do-work/working/`, or an active `do-work/user-requests/` body to make an unfinished target appear complete.' <<<"$shared_reference_text" \
  || fail 'completed-work reference no-fallback line changed'

# User-facing surfaces: guide, routing, help.
require_text "$guide_file" '--kind proposal'
require_text "$guide_file" '--kind root-cause'
skill_text="$(cat "$toolbox/SKILL.md")"
routing_row="$(grep -F -- './actions/ai-report.md' <<<"$skill_text")" || true
for routing_phrase in 'proposal report' 'root cause report' 'options report'; do
  grep -Fq -- "\`$routing_phrase\`" <<<"$routing_row" || fail "ai-report routing row lacks: $routing_phrase"
done
require_text "$toolbox/actions/help.md" 'ai-report --kind'

# Every relative link and path the shipped files name still resolves.
reference_contract_output="$(bash "$repo_root/_dev/tests/shipped-package-reference-contract.sh" 2>&1)" \
  || fail "shipped-package-reference-contract.sh failed: $reference_contract_output"

if [ "$failures" -ne 0 ]; then
  printf 'REQ-687 probe: %s failure(s)\n' "$failures" >&2
  exit 1
fi
printf 'REQ-687 probe passed.\n'
