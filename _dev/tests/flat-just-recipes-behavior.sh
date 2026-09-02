#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "$0")/../.." && pwd)"
template="$repo_root/skills/do-work-board/justfile.template"
root_justfile="$repo_root/justfile"

expected_recipes='architecture-report-preflight audit-metrics bkb-init bkb-lint-structure bkb-status do-work-answer do-work-cancel do-work-capture-files do-work-claim do-work-cleanup do-work-complete do-work-defer-gate do-work-doctor do-work-fail do-work-next do-work-release do-work-unblock do-work-update dream-scan generate-report-image generate-report-image-batch install-last30days interview-export interview-ingest interview-list interview-reset interview-status interview-versions kanban-static kanban-summary memory-audit memory-bootstrap memory-forget memory-recall memory-remember memory-status publish-portfolio-summary run-do-work-update run-kanban run-kanban-cli do-work-note'
expected_recipes="$(tr ' ' '\n' <<<"$expected_recipes" | LC_ALL=C sort | paste -sd' ' -)"

for justfile_path in "$template" "$root_justfile"; do
  actual_recipes="$(just --justfile "$justfile_path" --summary | tr ' ' '\n' | grep -v '^maintainer-verify$' | LC_ALL=C sort | paste -sd' ' -)"
  if [ "$actual_recipes" != "$expected_recipes" ]; then
    printf 'FAIL: %s recipe inventory differs\n got: %s\nwant: %s\n' "$justfile_path" "$actual_recipes" "$expected_recipes" >&2
    exit 1
  fi
done

fixture_root="$(mktemp -d)"
trap 'rm -rf "$fixture_root"' EXIT
mkdir -p "$fixture_root/.claude/skills/do-work/tools"
cp "$template" "$fixture_root/justfile"
cat > "$fixture_root/.claude/skills/do-work/tools/do-work-cli.sh" <<'SH'
#!/bin/sh
printf '%s\0' "$@" >> "$DO_WORK_RECIPE_ARGV_LOG"
SH
chmod +x "$fixture_root/.claude/skills/do-work/tools/do-work-cli.sh"
argv_log="$fixture_root/argv.log"
: > "$argv_log"

# Actual-template shell boundary: each hostile value must remain one inert argv
# byte sequence and neither shell-substitution spelling may execute.
substitution_sentinel="$fixture_root/command-substitution-ran"
backtick_sentinel="$fixture_root/backtick-substitution-ran"
hostile_substitution="\$(touch \"$substitution_sentinel\")"
hostile_backtick="\`touch \"$backtick_sentinel\"\`"
hostile_arguments=(
  'space value'
  "single ' quote"
  'double " quote'
  '$HOME literal dollar'
  "$hostile_substitution"
  "$hostile_backtick"
  $'tab\tvalue'
  $'line one\nline two'
)
DO_WORK_RECIPE_ARGV_LOG="$argv_log" just --justfile "$fixture_root/justfile" \
  do-work-answer "${hostile_arguments[@]}" >/dev/null 2>&1
if [ -e "$substitution_sentinel" ] || [ -e "$backtick_sentinel" ]; then
  printf 'FAIL: actual managed do-work-answer recipe executed hostile shell substitution\n' >&2
  exit 1
fi
hostile_expected="$fixture_root/hostile.expected"
printf '%s\0' --repo-root "$fixture_root" answer "${hostile_arguments[@]}" > "$hostile_expected"
if ! cmp -s "$hostile_expected" "$argv_log"; then
  printf 'FAIL: actual managed do-work-answer recipe did not preserve hostile argv bytes\n' >&2
  exit 1
fi
: > "$argv_log"

recipe_arguments() {
  case "$1" in
    interview-status|interview-export|interview-ingest|interview-reset|interview-versions) printf '%s\n' 'fixture-template' ;;
    memory-remember) printf '%s\n' 'remembered text' ;;
    memory-forget) printf '%s\n' 'forget query' ;;
    memory-bootstrap) printf '%s\n' 'manifest.json' ;;
  esac
}

for recipe_name in $expected_recipes; do
  case "$recipe_name" in
    run-kanban|run-kanban-cli|kanban-static|kanban-summary)
      just --justfile "$fixture_root/justfile" --dry-run "$recipe_name" >/dev/null 2>&1
      ;;
    *)
      recipe_arg="$(recipe_arguments "$recipe_name")"
      if [ -n "$recipe_arg" ]; then
        DO_WORK_RECIPE_ARGV_LOG="$argv_log" just --justfile "$fixture_root/justfile" "$recipe_name" "$recipe_arg" >/dev/null 2>&1
      else
        DO_WORK_RECIPE_ARGV_LOG="$argv_log" just --justfile "$fixture_root/justfile" "$recipe_name" >/dev/null 2>&1
      fi
      ;;
  esac
done

printf 'flat Just recipe behavior passed (41 definitions)\n'
