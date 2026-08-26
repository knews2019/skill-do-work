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

# architecture-report-preflight: a first run reports no prior baseline and the unsuffixed candidate.
first_run_repo="$fixture_root/first-run"
fixture_repo_init "$first_run_repo"
printf 'seed\n' > "$first_run_repo/README.md"
fixture_repo_commit_all "$first_run_repo" base
first_run_record="$(cd "$first_run_repo" && "$preflight_script" --scan docs)" \
  || fail_case 'architecture-report-preflight first-run scan returned nonzero'
first_run_date="$(date -u +%Y%m%d)"
[ "$(scan_field "$first_run_record" report_candidate)" = "docs/architecture-report_${first_run_date}.md" ] \
  || fail_case 'architecture-report-preflight first-run candidate is not the unsuffixed dated path'
[ -z "$(scan_field "$first_run_record" prior_report)" ] \
  || fail_case 'architecture-report-preflight first-run reported a prior report'
[ -z "$(scan_field "$first_run_record" prior_hash)" ] \
  || fail_case 'architecture-report-preflight first-run reported a prior hash'
[ "$(scan_field "$first_run_record" prior_hash_resolves)" = 'n/a' ] \
  || fail_case 'architecture-report-preflight first-run did not mark resolution as not applicable'
first_run_head="$(git -C "$first_run_repo" rev-parse --short HEAD)"
[ "$(scan_field "$first_run_record" head_hash)" = "$first_run_head" ] \
  || fail_case 'architecture-report-preflight first-run head_hash does not match HEAD'

# architecture-report-preflight: publication never modifies an existing report and escalates the suffix.
publish_repo="$fixture_root/publish"
fixture_repo_init "$publish_repo"
printf 'seed\n' > "$publish_repo/README.md"
fixture_repo_commit_all "$publish_repo" base
mkdir -p "$publish_repo/drafts"
printf 'first bytes\n' > "$publish_repo/drafts/first.md"
printf 'second bytes\n' > "$publish_repo/drafts/second.md"
publish_candidate='docs/architecture-report_20260826.md'
first_published="$(cd "$publish_repo" && "$preflight_script" --publish drafts/first.md "$publish_candidate")" \
  || fail_case 'architecture-report-preflight first publication returned nonzero'
[ "$first_published" = "$publish_candidate" ] \
  || fail_case "architecture-report-preflight first publication landed on $first_published"
second_published="$(cd "$publish_repo" && "$preflight_script" --publish drafts/second.md "$publish_candidate")" \
  || fail_case 'architecture-report-preflight second publication returned nonzero'
[ "$second_published" = 'docs/architecture-report_20260826_2.md' ] \
  || fail_case "architecture-report-preflight second publication landed on $second_published"
[ "$(cat "$publish_repo/$publish_candidate")" = 'first bytes' ] \
  || fail_case 'architecture-report-preflight second publication overwrote the first report'
[ "$(cat "$publish_repo/$second_published")" = 'second bytes' ] \
  || fail_case 'architecture-report-preflight suffixed publication does not carry the second draft'

# architecture-report-preflight: a directory squatting the candidate path is stepped over, not published into.
directory_candidate_repo="$fixture_root/directory-candidate"
fixture_repo_init "$directory_candidate_repo"
printf 'seed\n' > "$directory_candidate_repo/README.md"
fixture_repo_commit_all "$directory_candidate_repo" base
mkdir -p "$directory_candidate_repo/docs/architecture-report_20260826.md" "$directory_candidate_repo/drafts"
printf 'draft bytes\n' > "$directory_candidate_repo/drafts/report.md"
directory_published="$(cd "$directory_candidate_repo" && "$preflight_script" --publish drafts/report.md "$publish_candidate")" \
  || fail_case 'architecture-report-preflight directory-candidate publication returned nonzero'
[ "$directory_published" = 'docs/architecture-report_20260826_2.md' ] \
  || fail_case "architecture-report-preflight directory-candidate publication landed on $directory_published"
[ -z "$(find "$directory_candidate_repo/docs/architecture-report_20260826.md" -type f)" ] \
  || fail_case 'architecture-report-preflight nested a report inside the squatting directory'

# architecture-report-preflight: the newest prior report is chosen numerically, not lexically.
ordering_repo="$fixture_root/ordering"
fixture_repo_init "$ordering_repo"
printf 'seed\n' > "$ordering_repo/README.md"
fixture_repo_commit_all "$ordering_repo" base
ordering_head="$(git -C "$ordering_repo" rev-parse --short HEAD)"
mkdir -p "$ordering_repo/docs"
printf 'verified-at: %s\n' 'aaaaaaa' > "$ordering_repo/docs/architecture-report_20260101_2.md"
printf 'verified-at: %s\n' "$ordering_head" > "$ordering_repo/docs/architecture-report_20260101_10.md"
printf 'not a report\n' > "$ordering_repo/docs/architecture-report_notadate.md"
ordering_record="$(cd "$ordering_repo" && "$preflight_script" --scan docs)" \
  || fail_case 'architecture-report-preflight ordering scan returned nonzero'
[ "$(scan_field "$ordering_record" prior_report)" = 'docs/architecture-report_20260101_10.md' ] \
  || fail_case 'architecture-report-preflight chose the prior report lexically instead of numerically'
[ "$(scan_field "$ordering_record" prior_hash)" = "$ordering_head" ] \
  || fail_case 'architecture-report-preflight did not read the prior watermark hash'
[ "$(scan_field "$ordering_record" prior_hash_resolves)" = 'yes' ] \
  || fail_case 'architecture-report-preflight did not resolve a real prior hash'

# architecture-report-preflight: an unparseable watermark reports `unreadable`, never an empty hash.
unreadable_repo="$fixture_root/unreadable"
fixture_repo_init "$unreadable_repo"
printf 'seed\n' > "$unreadable_repo/README.md"
fixture_repo_commit_all "$unreadable_repo" base
mkdir -p "$unreadable_repo/docs"
printf '# Architecture Report\n\nNo watermark here.\n' > "$unreadable_repo/docs/architecture-report_20260101.md"
unreadable_record="$(cd "$unreadable_repo" && "$preflight_script" --scan docs)" \
  || fail_case 'architecture-report-preflight unreadable-watermark scan returned nonzero'
[ "$(scan_field "$unreadable_record" prior_hash)" = 'unreadable' ] \
  || fail_case 'architecture-report-preflight collapsed an unreadable watermark to an empty hash'
[ "$(scan_field "$unreadable_record" prior_hash_resolves)" = 'no' ] \
  || fail_case 'architecture-report-preflight claimed an unreadable watermark resolves'

# architecture-report-preflight: a watermark hash absent from this repository resolves as no.
unresolvable_repo="$fixture_root/unresolvable"
fixture_repo_init "$unresolvable_repo"
printf 'seed\n' > "$unresolvable_repo/README.md"
fixture_repo_commit_all "$unresolvable_repo" base
mkdir -p "$unresolvable_repo/docs"
printf 'verified-at: 0123456 · 2026-01-01\n' > "$unresolvable_repo/docs/architecture-report_20260101.md"
unresolvable_record="$(cd "$unresolvable_repo" && "$preflight_script" --scan docs)" \
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
  || fail_case 'architecture-report-preflight published a report from a missing draft'

prescribed_shell_finish
