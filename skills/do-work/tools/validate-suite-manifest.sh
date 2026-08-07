#!/usr/bin/env bash
# Validates the complete do-work suite in an extracted archive or staging tree.
set -euo pipefail

fail() {
  printf 'suite manifest: %s\n' "$*" >&2
  exit 1
}

if [ "$#" -ne 2 ] || [ "$1" != '--root' ]; then
  fail 'usage: validate-suite-manifest.sh --root <archive-root>'
fi

archive_root="$2"
[ -d "$archive_root" ] || fail "root does not exist: $archive_root"
archive_root="$(cd "$archive_root" && pwd -P)"
version_file="$archive_root/VERSION"
manifest_file="$archive_root/suite/modules.tsv"

[ -f "$version_file" ] && [ ! -L "$version_file" ] \
  || fail 'VERSION must be a regular file in the suite root'
[ "$(wc -l < "$version_file" | tr -d ' ')" = '1' ] \
  || fail 'VERSION must contain exactly one newline-terminated line'
suite_version="$(sed -n '1p' "$version_file")"
printf '%s\n' "$suite_version" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$' \
  || fail 'VERSION must be a plain semantic version (X.Y.Z)'

[ -f "$manifest_file" ] && [ ! -L "$manifest_file" ] \
  || fail 'suite/modules.tsv must be a regular file in the suite root'
IFS= read -r manifest_header < "$manifest_file" \
  || fail 'suite/modules.tsv is empty'
[ "$manifest_header" = $'source\tdestination' ] \
  || fail 'manifest header must be exactly: source<TAB>destination'

row_count=0
seen_sources=''
seen_destinations=''
seen_do_work=''
seen_board=''
seen_knowledge=''
seen_toolbox=''
line_number=1

while IFS= read -r manifest_line || [ -n "$manifest_line" ]; do
  line_number=$((line_number + 1))
  [ -n "$manifest_line" ] || fail "line $line_number is blank"
  case "$manifest_line" in
    *$'\r'*) fail "line $line_number contains a carriage return" ;;
  esac

  manifest_destination="${manifest_line#*$'\t'}"
  [ "$manifest_destination" != "$manifest_line" ] \
    || fail "line $line_number must contain one tab"
  manifest_source="${manifest_line%%$'\t'*}"
  case "$manifest_destination" in
    *$'\t'*) fail "line $line_number has an unknown extra column" ;;
  esac
  [ -n "$manifest_source" ] && [ -n "$manifest_destination" ] \
    || fail "line $line_number has an empty column"

  case "/$manifest_source/" in
    */../*|*/./*) fail "line $line_number source traverses directories" ;;
  esac
  case "/$manifest_destination/" in
    */../*|*/./*) fail "line $line_number destination traverses directories" ;;
  esac
  case "$manifest_source" in
    /*) fail "line $line_number source must be relative" ;;
  esac
  case "$manifest_destination" in
    /*) fail "line $line_number destination must be relative" ;;
  esac

  case "|$seen_sources|" in
    *"|$manifest_source|"*) fail "line $line_number duplicates source $manifest_source" ;;
  esac
  case "|$seen_destinations|" in
    *"|$manifest_destination|"*) fail "line $line_number duplicates destination $manifest_destination" ;;
  esac
  seen_sources="$seen_sources|$manifest_source"
  seen_destinations="$seen_destinations|$manifest_destination"

  case "$manifest_source" in
    skills/do-work)
      expected_destination='.claude/skills/do-work'
      seen_do_work=1
      ;;
    skills/do-work-board)
      expected_destination='.claude/skills/do-work-board'
      seen_board=1
      ;;
    skills/do-work-knowledge)
      expected_destination='.claude/skills/do-work-knowledge'
      seen_knowledge=1
      ;;
    skills/do-work-toolbox)
      expected_destination='.claude/skills/do-work-toolbox'
      seen_toolbox=1
      ;;
    *) fail "line $line_number declares unexpected source $manifest_source" ;;
  esac
  [ "$manifest_destination" = "$expected_destination" ] \
    || fail "line $line_number maps $manifest_source outside its required destination"

  module_root="$archive_root/$manifest_source"
  [ -d "$module_root" ] && [ ! -L "$module_root" ] \
    || fail "$manifest_source must be a real directory, not a missing path or symlink"
  module_physical_root="$(cd "$module_root" && pwd -P)"
  case "$module_physical_root" in
    "$archive_root"/skills/*) ;;
    *) fail "$manifest_source resolves outside the suite root" ;;
  esac
  [ -s "$module_root/SKILL.md" ] && [ ! -L "$module_root/SKILL.md" ] \
    || fail "$manifest_source/SKILL.md must be a non-empty regular file"

  row_count=$((row_count + 1))
done < <(sed '1d' "$manifest_file")

[ "$row_count" -eq 4 ] || fail "manifest must contain exactly four modules (found $row_count)"
[ -n "$seen_do_work" ] && [ -n "$seen_board" ] \
  && [ -n "$seen_knowledge" ] && [ -n "$seen_toolbox" ] \
  || fail 'manifest must contain all four required do-work modules'

printf 'suite manifest valid: v%s (%s modules)\n' "$suite_version" "$row_count"
