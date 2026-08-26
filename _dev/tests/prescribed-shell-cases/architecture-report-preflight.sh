#!/usr/bin/env bash
# Fixture execution proofs for architecture-report-preflight.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

preflight_script="$toolbox_scripts/architecture-report-preflight.sh"

# Reads one `key=value` line out of a captured --scan record.
scan_field() {
  local record_text="$1"
  local field_name="$2"
  local record_line

  while IFS= read -r record_line; do
    case "$record_line" in
      "$field_name"=*) printf '%s' "${record_line#*=}"; return 0 ;;
    esac
  done <<< "$record_text"
  return 1
}

# Seeds a published report bundle carrying a watermark, without going through --publish.
seed_report_bundle() {
  local bundle_path="$1"
  local watermark_hash="$2"

  mkdir -p "$bundle_path"
  printf 'verified-at: %s · 2026-01-01 · 0.1.0 · prior: none\n' "$watermark_hash" \
    > "$bundle_path/architecture-report.md"
}

# architecture-report-preflight: a first run reports no prior baseline and the unsuffixed candidate.
first_run_repo="$fixture_root/first-run"
fixture_repo_init "$first_run_repo"
printf 'seed\n' > "$first_run_repo/README.md"
fixture_repo_commit_all "$first_run_repo" base
first_run_record="$(cd "$first_run_repo" && "$preflight_script" --scan ai-reports)" \
  || fail_case 'architecture-report-preflight first-run scan returned nonzero'
first_run_slug="$(scan_field "$first_run_record" report_slug)"
case "$first_run_slug" in
  [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]_[0-9][0-9][0-9][0-9]) ;;
  *) fail_case "architecture-report-preflight slug is not yyyy-mm-dd_hhmm: $first_run_slug" ;;
esac
[ "$(scan_field "$first_run_record" report_candidate)" = "ai-reports/${first_run_slug}_architecture-report" ] \
  || fail_case 'architecture-report-preflight first-run candidate is not the unsuffixed bundle path'
[ -z "$(scan_field "$first_run_record" prior_report)" ] \
  || fail_case 'architecture-report-preflight first-run reported a prior report'
[ -z "$(scan_field "$first_run_record" prior_hash)" ] \
  || fail_case 'architecture-report-preflight first-run reported a prior hash'
[ "$(scan_field "$first_run_record" prior_hash_resolves)" = 'n/a' ] \
  || fail_case 'architecture-report-preflight first-run did not mark resolution as not applicable'
[ "$(scan_field "$first_run_record" head_hash)" = "$(git -C "$first_run_repo" rev-parse --short HEAD)" ] \
  || fail_case 'architecture-report-preflight first-run head_hash does not match HEAD'

# architecture-report-preflight: publication writes a bundle and never modifies an existing one.
publish_repo="$fixture_root/publish"
fixture_repo_init "$publish_repo"
printf 'seed\n' > "$publish_repo/README.md"
fixture_repo_commit_all "$publish_repo" base
mkdir -p "$publish_repo/drafts"
printf 'first bytes\n' > "$publish_repo/drafts/first.md"
printf 'second bytes\n' > "$publish_repo/drafts/second.md"
publish_candidate='ai-reports/2026-08-26_1530_architecture-report'
first_published="$(cd "$publish_repo" && "$preflight_script" --publish drafts/first.md "$publish_candidate")" \
  || fail_case 'architecture-report-preflight first publication returned nonzero'
[ "$first_published" = "$publish_candidate/architecture-report.md" ] \
  || fail_case "architecture-report-preflight first publication landed on $first_published"
second_published="$(cd "$publish_repo" && "$preflight_script" --publish drafts/second.md "$publish_candidate")" \
  || fail_case 'architecture-report-preflight second publication returned nonzero'
[ "$second_published" = 'ai-reports/2026-08-26_1530_architecture-report-2/architecture-report.md' ] \
  || fail_case "architecture-report-preflight second publication landed on $second_published"
[ "$(cat "$publish_repo/$first_published")" = 'first bytes' ] \
  || fail_case 'architecture-report-preflight second publication overwrote the first report'
[ "$(cat "$publish_repo/$second_published")" = 'second bytes' ] \
  || fail_case 'architecture-report-preflight suffixed publication does not carry the second draft'

# architecture-report-preflight: an occupied candidate is stepped over, never published into.
# `mkdir -p` would succeed here and nest this run's report inside the previous run's bundle.
occupied_repo="$fixture_root/occupied"
fixture_repo_init "$occupied_repo"
printf 'seed\n' > "$occupied_repo/README.md"
fixture_repo_commit_all "$occupied_repo" base
mkdir -p "$occupied_repo/drafts"
printf 'draft bytes\n' > "$occupied_repo/drafts/report.md"
seed_report_bundle "$occupied_repo/$publish_candidate" 'aaaaaaa'
occupied_published="$(cd "$occupied_repo" && "$preflight_script" --publish drafts/report.md "$publish_candidate")" \
  || fail_case 'architecture-report-preflight occupied-candidate publication returned nonzero'
[ "$occupied_published" = 'ai-reports/2026-08-26_1530_architecture-report-2/architecture-report.md' ] \
  || fail_case "architecture-report-preflight occupied-candidate publication landed on $occupied_published"
[ "$(head -c 12 "$occupied_repo/$publish_candidate/architecture-report.md")" = 'verified-at:' ] \
  || fail_case 'architecture-report-preflight overwrote the occupied bundle report'

# architecture-report-preflight: the newest prior report is chosen numerically, not lexically.
ordering_repo="$fixture_root/ordering"
fixture_repo_init "$ordering_repo"
printf 'seed\n' > "$ordering_repo/README.md"
fixture_repo_commit_all "$ordering_repo" base
ordering_head="$(git -C "$ordering_repo" rev-parse --short HEAD)"
seed_report_bundle "$ordering_repo/ai-reports/2026-01-01_0900_architecture-report-2" 'aaaaaaa'
seed_report_bundle "$ordering_repo/ai-reports/2026-01-01_0900_architecture-report-10" "$ordering_head"
seed_report_bundle "$ordering_repo/ai-reports/2026-01-01_0800_architecture-report" 'bbbbbbb'
mkdir -p "$ordering_repo/ai-reports/2026-08-26_1530_UR-042-something-else"
mkdir -p "$ordering_repo/ai-reports/notadate_architecture-report"
ordering_record="$(cd "$ordering_repo" && "$preflight_script" --scan ai-reports)" \
  || fail_case 'architecture-report-preflight ordering scan returned nonzero'
[ "$(scan_field "$ordering_record" prior_report)" = 'ai-reports/2026-01-01_0900_architecture-report-10/architecture-report.md' ] \
  || fail_case 'architecture-report-preflight chose the prior report lexically instead of numerically'
[ "$(scan_field "$ordering_record" prior_hash)" = "$ordering_head" ] \
  || fail_case 'architecture-report-preflight did not read the prior watermark hash'
[ "$(scan_field "$ordering_record" prior_hash_resolves)" = 'yes' ] \
  || fail_case 'architecture-report-preflight did not resolve a real prior hash'

# architecture-report-preflight: an unparseable or absent watermark reports `unreadable`, never empty.
unreadable_repo="$fixture_root/unreadable"
fixture_repo_init "$unreadable_repo"
printf 'seed\n' > "$unreadable_repo/README.md"
fixture_repo_commit_all "$unreadable_repo" base
mkdir -p "$unreadable_repo/ai-reports/2026-01-01_0900_architecture-report"
printf '# Architecture Report\n\nNo watermark here.\n' \
  > "$unreadable_repo/ai-reports/2026-01-01_0900_architecture-report/architecture-report.md"
mkdir -p "$unreadable_repo/ai-reports/2026-01-02_0900_architecture-report"
unreadable_record="$(cd "$unreadable_repo" && "$preflight_script" --scan ai-reports)" \
  || fail_case 'architecture-report-preflight unreadable-watermark scan returned nonzero'
[ "$(scan_field "$unreadable_record" prior_hash)" = 'unreadable' ] \
  || fail_case 'architecture-report-preflight collapsed a missing report file to an empty hash'
[ "$(scan_field "$unreadable_record" prior_hash_resolves)" = 'no' ] \
  || fail_case 'architecture-report-preflight claimed an unreadable watermark resolves'

# architecture-report-preflight: a watermark hash absent from this repository resolves as no.
unresolvable_repo="$fixture_root/unresolvable"
fixture_repo_init "$unresolvable_repo"
printf 'seed\n' > "$unresolvable_repo/README.md"
fixture_repo_commit_all "$unresolvable_repo" base
seed_report_bundle "$unresolvable_repo/ai-reports/2026-01-01_0900_architecture-report" '0123456'
unresolvable_record="$(cd "$unresolvable_repo" && "$preflight_script" --scan ai-reports)" \
  || fail_case 'architecture-report-preflight unresolvable-hash scan returned nonzero'
[ "$(scan_field "$unresolvable_record" prior_hash)" = '0123456' ] \
  || fail_case 'architecture-report-preflight did not read a syntactically valid prior hash'
[ "$(scan_field "$unresolvable_record" prior_hash_resolves)" = 'no' ] \
  || fail_case 'architecture-report-preflight claimed an absent commit resolves'

# architecture-report-preflight: usage errors and a missing draft exit 2 without publishing.
usage_repo="$fixture_root/usage"
fixture_repo_init "$usage_repo"
printf 'seed\n' > "$usage_repo/README.md"
fixture_repo_commit_all "$usage_repo" base
(cd "$usage_repo" && "$preflight_script" >/dev/null 2>&1)
[ "$?" -eq 2 ] || fail_case 'architecture-report-preflight with no verb did not exit 2'
(cd "$usage_repo" && "$preflight_script" --scan >/dev/null 2>&1)
[ "$?" -eq 2 ] || fail_case 'architecture-report-preflight --scan without a directory did not exit 2'
(cd "$usage_repo" && "$preflight_script" --publish drafts/absent.md "$publish_candidate" >/dev/null 2>&1)
[ "$?" -eq 2 ] || fail_case 'architecture-report-preflight --publish with a missing draft did not exit 2'
[ ! -e "$usage_repo/$publish_candidate" ] \
  || fail_case 'architecture-report-preflight created a bundle from a missing draft'

prescribed_shell_finish
