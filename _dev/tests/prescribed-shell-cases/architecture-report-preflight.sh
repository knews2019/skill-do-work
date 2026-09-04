#!/usr/bin/env bash
# Fixture execution proofs for architecture-report-preflight.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

preflight_script="$fixture_root/architecture-report-preflight.sh"
cat > "$preflight_script" <<EOF
#!/usr/bin/env bash
set -euo pipefail
repository_probe='.'
if [[ "\${1:-}" == --scan && "\$#" -eq 2 ]]; then
  repository_probe="\$2"
elif [[ "\${1:-}" == --publish && "\$#" -eq 3 ]]; then
  repository_probe="\$(dirname "\$2")"
fi
repository_root="\$(git -C "\$repository_probe" rev-parse --show-toplevel 2>/dev/null || pwd -P)"
DO_WORK_COMPATIBILITY_SHIM=1 exec bash "$repo_root/skills/do-work/tools/do-work-cli.sh" --repo-root "\$repository_root" --format text architecture-report-preflight "\$@"
EOF
chmod +x "$preflight_script"

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
  printf '<!doctype html>\n<html><head>\n<meta name="architecture-report-verified-at" content="%s">\n</head><body>Architecture</body></html>\n' "$watermark_hash" \
    > "$bundle_path/index.html"
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
printf '<!doctype html><title>First</title><p>first bytes</p>\n' > "$publish_repo/drafts/first.html"
printf '<!doctype html><title>Second</title><p>second bytes</p>\n' > "$publish_repo/drafts/second.html"
publish_candidate='ai-reports/2026-08-26_1530_architecture-report'
first_published="$(cd "$publish_repo" && "$preflight_script" --publish drafts/first.html "$publish_candidate")" \
  || fail_case 'architecture-report-preflight first publication returned nonzero'
[ "$first_published" = "$publish_candidate/index.html" ] \
  || fail_case "architecture-report-preflight first publication landed on $first_published"
second_published="$(cd "$publish_repo" && "$preflight_script" --publish drafts/second.html "$publish_candidate")" \
  || fail_case 'architecture-report-preflight second publication returned nonzero'
[ "$second_published" = 'ai-reports/2026-08-26_1530_architecture-report-2/index.html' ] \
  || fail_case "architecture-report-preflight second publication landed on $second_published"
cmp -s "$publish_repo/drafts/first.html" "$publish_repo/$first_published" \
  || fail_case 'architecture-report-preflight second publication overwrote the first report'
cmp -s "$publish_repo/drafts/second.html" "$publish_repo/$second_published" \
  || fail_case 'architecture-report-preflight suffixed publication does not carry the second draft'

[ ! -e "$publish_repo/$publish_candidate/architecture-report.md" ] \
  || fail_case 'architecture-report-preflight published a Markdown companion'
[ "$(find "$publish_repo/$publish_candidate" -type f | wc -l | tr -d ' ')" = 1 ] \
  || fail_case 'architecture-report-preflight single-file bundle contains extra output'

# architecture-report-preflight: an occupied candidate is stepped over, never published into.
# `mkdir -p` would succeed here and nest this run's report inside the previous run's bundle.
occupied_repo="$fixture_root/occupied"
fixture_repo_init "$occupied_repo"
printf 'seed\n' > "$occupied_repo/README.md"
fixture_repo_commit_all "$occupied_repo" base
mkdir -p "$occupied_repo/drafts"
printf 'draft bytes\n' > "$occupied_repo/drafts/report.html"
seed_report_bundle "$occupied_repo/$publish_candidate" 'aaaaaaa'
occupied_published="$(cd "$occupied_repo" && "$preflight_script" --publish drafts/report.html "$publish_candidate")" \
  || fail_case 'architecture-report-preflight occupied-candidate publication returned nonzero'
[ "$occupied_published" = 'ai-reports/2026-08-26_1530_architecture-report-2/index.html' ] \
  || fail_case "architecture-report-preflight occupied-candidate publication landed on $occupied_published"
grep -q 'content="aaaaaaa"' "$occupied_repo/$publish_candidate/index.html" \
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
seed_report_bundle "$ordering_repo/ai-reports/notadate_architecture-report" 'ccccccc'
ordering_record="$(cd "$ordering_repo" && "$preflight_script" --scan ai-reports)" \
  || fail_case 'architecture-report-preflight ordering scan returned nonzero'
[ "$(scan_field "$ordering_record" prior_report)" = 'ai-reports/2026-01-01_0900_architecture-report-10/index.html' ] \
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
printf '<!doctype html><title>Architecture</title><p>No watermark here.</p>\n' \
  > "$unreadable_repo/ai-reports/2026-01-01_0900_architecture-report/index.html"
mkdir -p "$unreadable_repo/ai-reports/2026-01-02_0900_architecture-report"
unreadable_record="$(cd "$unreadable_repo" && "$preflight_script" --scan ai-reports)" \
  || fail_case 'architecture-report-preflight unreadable-watermark scan returned nonzero'
[ "$(scan_field "$unreadable_record" prior_report)" = 'ai-reports/2026-01-01_0900_architecture-report/index.html' ] \
  || fail_case 'architecture-report-preflight selected a missing HTML file instead of the prior report'
[ "$(scan_field "$unreadable_record" prior_hash)" = 'unreadable' ] \
  || fail_case 'architecture-report-preflight collapsed a missing watermark to an empty hash'
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
(cd "$usage_repo" && "$preflight_script" --publish drafts/absent.html "$publish_candidate" >/dev/null 2>&1)
[ "$?" -eq 2 ] || fail_case 'architecture-report-preflight --publish with a missing draft did not exit 2'
[ ! -e "$usage_repo/$publish_candidate" ] \
  || fail_case 'architecture-report-preflight created a bundle from a missing draft'

# architecture-report-preflight: legacy Markdown and unfinished bundles never become HTML baselines.
legacy_repo="$fixture_root/legacy"
fixture_repo_init "$legacy_repo"
printf 'seed\n' > "$legacy_repo/README.md"
fixture_repo_commit_all "$legacy_repo" base
legacy_bundle="$legacy_repo/ai-reports/2026-08-26_1709_architecture-report"
mkdir -p "$legacy_bundle"
printf 'verified-at: 0123456 · legacy Markdown\n' > "$legacy_bundle/architecture-report.md"
legacy_record="$(cd "$legacy_repo" && "$preflight_script" --scan ai-reports)" \
  || fail_case 'architecture-report-preflight legacy-only scan returned nonzero'
[ -z "$(scan_field "$legacy_record" prior_report)" ] \
  || fail_case 'architecture-report-preflight selected a Markdown-only prior report'
seed_report_bundle "$legacy_repo/ai-reports/2026-01-01_0900_architecture-report" '0123456'
mkdir -p "$legacy_repo/ai-reports/2026-08-27_0900_architecture-report"
legacy_record="$(cd "$legacy_repo" && "$preflight_script" --scan ai-reports)" \
  || fail_case 'architecture-report-preflight mixed-format scan returned nonzero'
[ "$(scan_field "$legacy_record" prior_report)" = 'ai-reports/2026-01-01_0900_architecture-report/index.html' ] \
  || fail_case 'architecture-report-preflight chose a newer Markdown or unfinished bundle over HTML'

# architecture-report-preflight: a failed copy cannot expose partial HTML as a prior baseline.
partial_repo="$fixture_root/partial"
fixture_repo_init "$partial_repo"
printf 'seed\n' > "$partial_repo/README.md"
fixture_repo_commit_all "$partial_repo" base
mkdir -p "$partial_repo/drafts" "$partial_repo/bin"
printf '<!doctype html><title>Complete</title>\n' > "$partial_repo/drafts/report.html"
cat > "$partial_repo/bin/cp" <<'SH'
#!/usr/bin/env bash
printf '<!doctype html><title>Partial' > "$2"
exit 1
SH
chmod +x "$partial_repo/bin/cp"
(cd "$partial_repo" && PATH="$partial_repo/bin:$PATH" "$preflight_script" --publish drafts/report.html "$publish_candidate" >/dev/null 2>&1)
[ "$?" -ne 0 ] || fail_case 'architecture-report-preflight reported a failed copy as published'
[ ! -e "$partial_repo/$publish_candidate/index.html" ] \
  || fail_case 'architecture-report-preflight exposed partial HTML after a failed copy'
partial_record="$(cd "$partial_repo" && "$preflight_script" --scan ai-reports)" \
  || fail_case 'architecture-report-preflight partial-publication scan returned nonzero'
[ -z "$(scan_field "$partial_record" prior_report)" ] \
  || fail_case 'architecture-report-preflight selected a failed publication as a baseline'

prescribed_shell_finish
