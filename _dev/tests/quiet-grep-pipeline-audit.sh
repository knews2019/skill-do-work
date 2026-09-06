#!/usr/bin/env bash
# quiet-grep-pipeline-audit.sh — repository-wide guard: no tracked shell file may decide on a
# quiet grep that is fed from a pipeline. Scans every file `git ls-files -- '*.sh'` reports,
# which is the same enumeration the ShellCheck lane uses, and pins the scanner itself against a
# two-sided fixture of named shapes.
#
# Exit 0: every tracked shell file is clean and the scanner catches every shape it must.
# Exit 1: at least one FAIL line above.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
failure_count=0

fixture_directory="$(mktemp -d "${TMPDIR:-/tmp}/quiet-grep-pipeline-audit.XXXXXX")"
cleanup_fixture_directory() {
  rm -rf -- "$fixture_directory"
}
trap cleanup_fixture_directory EXIT

# The scanner and its contract comment live in quiet-grep-pipeline-scanner.sh, sourced here and by
# action-shell-blocks.sh (which runs it over shipped Markdown fences). The fixture below is what pins
# the scanner, so a change to the shared body fails here by shape name.
# shellcheck source=_dev/tests/quiet-grep-pipeline-scanner.sh
source "$repo_root/_dev/tests/quiet-grep-pipeline-scanner.sh"

# The fixture is two-sided and every shape is asserted by name, so a failure says which shape
# was lost rather than that a count moved. Each fixture line carries a unique marker token in
# its pattern and the checks read the markers out of the scan output, so adding a shape does not
# renumber the table.
#
# Both heredocs interpolate the pipe and hash characters. That is what keeps this file's own
# bytes free of the shape its own scanner looks for, so the audit scans itself clean with no
# path exemption — a checked-in fixture file would carry the literal offending bytes and would
# then need one.
write_must_flag_fixture() {
  local fixture_path="$1" pipe_character='|' hash_character='#'
  # `\\` inside the unquoted heredoc emits one literal backslash, which is the continuation
  # shape; a lone `\` would be eaten as a heredoc line continuation instead.
  cat > "$fixture_path" <<FIXTURE
tar tzf archive.tgz $pipe_character
  grep -q flag-pipe-at-end-of-line
tar tzf archive.tgz $pipe_character grep --quiet flag-long-option-quiet
tar tzf archive.tgz $pipe_character grep --silent flag-long-option-silent
tar tzf archive.tgz $pipe_character LC_ALL=C grep -q flag-locale-prefixed-reader
tar tzf archive.tgz $pipe_character command grep -q flag-command-prefixed-reader
tar tzf archive.tgz $pipe_character \\
  grep -q flag-backslash-continuation
tar tzf archive.tgz $pipe_character /usr/bin/grep -q flag-absolute-path-reader
tar tzf archive.tgz $pipe_character grep -F -q flag-quiet-option-not-first
tar tzf archive.tgz $pipe_character grep -m 1 flag-short-max-count
tar tzf archive.tgz $pipe_character grep --max-count=1 flag-long-max-count
tar tzf archive.tgz $pipe_character egrep -q flag-egrep-reader
tar tzf archive.tgz $pipe_character
$hash_character a note between the pipe and the reader
  grep -q flag-comment-line-after-pipe
tar tzf archive.tgz $pipe_character \\
$hash_character a note between the continuation and the reader
  grep -q flag-comment-line-after-continuation
tar tzf archive.tgz $pipe_character grep -q flag-trailing-note-on-the-line $hash_character an ordinary note
tar tzf archive.tgz $pipe_character

  grep -q flag-blank-line-after-pipe
tar tzf archive.tgz $pipe_character \\

  grep -q flag-blank-line-after-continuation
tar tzf archive.tgz $pipe_character $hash_character a note after the pipe on this line
  grep -q flag-note-after-pipe-same-line
tar tzf archive.tgz $pipe_character grep -m1 flag-short-max-count-no-space
tar tzf archive.tgz $pipe_character grep -qm1 flag-bundled-quiet-max-count
FIXTURE
}

write_must_not_flag_fixture() {
  local fixture_path="$1" pipe_character='|' hash_character='#'
  cat > "$fixture_path" <<FIXTURE
$hash_character tar tzf archive.tgz $pipe_character grep -q keep-whole-line-comment
grep -q keep-file-argument-reader -- "\$archive_listing_path"
grep -q keep-herestring-reader <<<"\$archive_listing_text"
[ -z "\$archive_listing_text" ] || grep -q keep-logical-or-is-not-a-pipe <<<"\$archive_listing_text"
printf '%s\n' "\$archive_listing_text" $pipe_character grep -c keep-counting-reader-runs-to-eof
printf '%s\n' "\$archive_listing_text" $pipe_character quiet_grep_report -q keep-helper-command-name
grep -q 'keep-quoted-pipe-in-the-pattern${pipe_character}second-branch' -- "\$archive_listing_path"
FIXTURE
}

# Marker token, then a tab, then the plain-English name reported when the shape is lost.
must_flag_shapes=(
  'flag-pipe-at-end-of-line	a pipe at end of line with the reader on the next line'
  'flag-long-option-quiet	grep --quiet'
  'flag-long-option-silent	grep --silent'
  'flag-locale-prefixed-reader	a locale-prefixed reader (LC_ALL=C grep -q)'
  'flag-command-prefixed-reader	a command-prefixed reader (command grep -q)'
  'flag-backslash-continuation	a backslash continuation before the reader'
  'flag-absolute-path-reader	an absolute-path reader (/usr/bin/grep -q)'
  'flag-quiet-option-not-first	grep -F -q, where -q is not the first option'
  'flag-short-max-count	grep -m 1'
  'flag-long-max-count	grep --max-count=1'
  'flag-egrep-reader	egrep -q'
  'flag-comment-line-after-pipe	a comment line between the pipe and the reader'
  'flag-comment-line-after-continuation	a comment line between the backslash continuation and the reader'
  'flag-trailing-note-on-the-line	a trailing # note on the offending line itself'
  'flag-blank-line-after-pipe	a blank line between the pipe and the reader'
  'flag-blank-line-after-continuation	a blank line between the backslash continuation and the reader'
  'flag-note-after-pipe-same-line	a note after the pipe on the same line, with the reader on the next'
  'flag-short-max-count-no-space	grep -m1, with no space before the count'
  'flag-bundled-quiet-max-count	grep -qm1, quiet bundled with a count'
)

must_not_flag_shapes=(
  'keep-whole-line-comment	a whole-line comment carrying the offending text'
  'keep-file-argument-reader	grep -q reading a file argument, which has no writer to kill'
  'keep-herestring-reader	grep -q reading a herestring, which has no writer to kill'
  'keep-logical-or-is-not-a-pipe	a logical or, which is not a pipe'
  'keep-counting-reader-runs-to-eof	grep -c, a reader that runs to EOF'
  'keep-helper-command-name	a helper command whose name merely ends in grep'
  'keep-quoted-pipe-in-the-pattern	a pipe inside a quoted pattern on a file-reading grep'
)

must_flag_fixture_path="$fixture_directory/must-flag-shapes.sh"
write_must_flag_fixture "$must_flag_fixture_path"
if ! must_flag_scan_output="$(quiet_grep_pipeline_offenders "$must_flag_fixture_path")"; then
  printf 'FAIL: the quiet-grep scanner could not read its own must-flag fixture.\n' >&2
  failure_count=$((failure_count + 1))
else
  lost_shape_names=()
  for must_flag_shape in "${must_flag_shapes[@]}"; do
    must_flag_marker="${must_flag_shape%%	*}"
    case "$must_flag_scan_output" in
      *"$must_flag_marker"*) ;;
      *) lost_shape_names+=("${must_flag_shape#*	}") ;;
    esac
  done
  if [ "${#lost_shape_names[@]}" -gt 0 ]; then
    printf 'FAIL: the quiet-grep scanner no longer catches %s of %s ordinary spellings of the pipeline it exists to forbid:\n' \
      "${#lost_shape_names[@]}" "${#must_flag_shapes[@]}" >&2
    printf '  no longer caught: %s\n' "${lost_shape_names[@]}" >&2
    failure_count=$((failure_count + 1))
  fi
fi

must_not_flag_fixture_path="$fixture_directory/must-not-flag-shapes.sh"
write_must_not_flag_fixture "$must_not_flag_fixture_path"
if ! must_not_flag_scan_output="$(quiet_grep_pipeline_offenders "$must_not_flag_fixture_path")"; then
  printf 'FAIL: the quiet-grep scanner could not read its own must-not-flag fixture.\n' >&2
  failure_count=$((failure_count + 1))
else
  wrongly_flagged_names=()
  for must_not_flag_shape in "${must_not_flag_shapes[@]}"; do
    must_not_flag_marker="${must_not_flag_shape%%	*}"
    case "$must_not_flag_scan_output" in
      *"$must_not_flag_marker"*) wrongly_flagged_names+=("${must_not_flag_shape#*	}") ;;
    esac
  done
  if [ "${#wrongly_flagged_names[@]}" -gt 0 ]; then
    printf 'FAIL: the quiet-grep scanner flags %s safe shape(s), so its findings cannot be trusted:\n' \
      "${#wrongly_flagged_names[@]}" >&2
    printf '  wrongly flagged: %s\n' "${wrongly_flagged_names[@]}" >&2
    failure_count=$((failure_count + 1))
  elif [ -n "$must_not_flag_scan_output" ]; then
    printf 'FAIL: the quiet-grep scanner reported an unnamed finding on its must-not-flag fixture:\n%s\n' \
      "$must_not_flag_scan_output" >&2
    failure_count=$((failure_count + 1))
  fi
fi

# The walk is a condition, not a file list: every tracked shell file, the same enumeration the
# ShellCheck lane runs. A file added tomorrow is scanned without anything being registered.
tracked_shell_files=()
while IFS= read -r -d '' tracked_shell_file; do
  tracked_shell_files+=("$tracked_shell_file")
done < <(git -C "$repo_root" ls-files -z -- '*.sh')
if [ "${#tracked_shell_files[@]}" -eq 0 ]; then
  printf 'FAIL: git reported no tracked shell files, so the quiet-grep audit scanned nothing.\n' >&2
  failure_count=$((failure_count + 1))
fi

offending_file_count=0
for tracked_shell_file in ${tracked_shell_files[@]+"${tracked_shell_files[@]}"}; do
  if ! offending_lines="$(quiet_grep_pipeline_offenders "$repo_root/$tracked_shell_file")"; then
    printf 'FAIL: could not scan %s for quiet greps fed from a pipeline.\n' "$tracked_shell_file" >&2
    failure_count=$((failure_count + 1))
    continue
  fi
  [ -n "$offending_lines" ] || continue
  printf 'FAIL: %s decides on a quiet grep fed from a pipeline; under pipefail the writer'"'"'s SIGPIPE death is reported as grep'"'"'s verdict:\n' \
    "$tracked_shell_file" >&2
  printf '%s\n' "$offending_lines" | sed 's|^|  |' >&2
  offending_file_count=$((offending_file_count + 1))
  failure_count=$((failure_count + 1))
done

if [ "$failure_count" -gt 0 ]; then
  printf 'quiet-grep pipeline audit: %s failure(s), %s of %s tracked shell files carry the shape.\n' \
    "$failure_count" "$offending_file_count" "${#tracked_shell_files[@]}" >&2
  exit 1
fi
printf 'quiet-grep pipeline audit passed (%s tracked shell files, %s must-flag and %s must-not-flag shapes).\n' \
  "${#tracked_shell_files[@]}" "${#must_flag_shapes[@]}" "${#must_not_flag_shapes[@]}"
