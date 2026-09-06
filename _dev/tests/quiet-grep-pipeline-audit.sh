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

# `writer | grep -q PATTERN` under `set -o pipefail` reports the writer's SIGPIPE death as
# grep's verdict: the assertion is then wrong in both directions — a positive matcher misses a
# pattern that is present, and a negative matcher fails to flag one it should. It is silent
# below roughly 36 KB of producer output and certain above about 200 KB, so a check written this
# way passes for years and then fails, or passes forever while asserting nothing (REQ-593).
#
# quiet_grep_pipeline_offenders prints every logical line in the named file that feeds an
# early-leaving grep from a pipeline. The defect is a condition, not a spelling: a reader that
# can exit before its writer is done, downstream of a pipe, under pipefail. So the scan works on
# the condition's three ingredients instead of on the one spelling that was converted —
# ordinary shell writes the same defect fourteen different ways in the must-flag fixture below
# alone, and a regex naming only today's shape goes stale the first time someone writes another
# (prime-shell-commands.md, "Closed Enumerations Go Stale").
#   * logical lines, not physical: a pipeline continued after a trailing `|` or `\` is one
#     command, and blank lines, comment lines and a trailing note after the pipe are all skipped
#     the way bash skips them — after a pipe bash keeps looking for the next command past blank
#     and comment lines, so `producer |` / `# note` / `grep -q x` and `producer | # note` /
#     `grep -q x` are each one running pipeline;
#   * any pipe stage after the first, so `command grep`, `LC_ALL=C grep` and `/usr/bin/grep`
#     count;
#   * any early-leaving option: -q, bundled -Eq, --quiet, --silent, -m/--max-count.
# A `#` at the start of a logical line suppresses exactly what bash suppresses, which is why moving
# an offender there deletes the check rather than hiding it. A `#` AFTER an open pipe does not: bash
# keeps reading, so that note is stripped above rather than treated as an end.
#
# WHAT THIS DOES NOT CATCH, stated so its silence is not read as coverage. The reader set is
# grep/egrep/fgrep: `rg -q`, `head -1`, `sed -n '1p;q'`, `awk '/x/{exit}'` and `read` are the same
# defect with a different early-leaving reader and are invisible here. A pipeline assembled at
# runtime — a variable command, `eval`, a deferred here-doc body — is invisible to any source scan.
# The parser is textual, and its three known false-positive shapes are all LOUD — a gate failure, never
# a silent pass. It splits on a bare `|`, so a `|` inside a quoted pattern makes a phantom stage; it has
# no here-doc awareness, so a here-doc body spelling the shape reads as an offender; and the comment
# skip is unconditional, so a comment that TERMINATES a backslash-continued command joins that command
# to the next one into a phantom pipeline. None has an instance in the tree today, and the first is
# pinned as a must-not-flag shape below.
quiet_grep_pipeline_offenders() {
  # awk rather than a grep pipeline so the scan reports its own failure instead of an
  # unreadable file reading as a clean file.
  awk '
    function feeds_an_early_leaving_grep(command_text,   pipe_stages, stage_count, stage_index) {
      gsub(/\|\|/, " OR ", command_text)
      stage_count = split(command_text, pipe_stages, "|")
      for (stage_index = 2; stage_index <= stage_count; stage_index++) {
        if (pipe_stages[stage_index] ~ /(^|[^-[:alnum:]_])(e|f)?grep([[:space:]]|$)/ \
          && pipe_stages[stage_index] ~ /(^|[[:space:]])(-[A-Za-z]*[qm][A-Za-z0-9]*|--quiet|--silent|--max-count)([[:space:]=]|$)/) {
          return 1
        }
      }
      return 0
    }
    /^[[:space:]]*#/ { next }
    joined_command != "" && /^[[:space:]]*$/ { next }
    {
      physical_line = $0
      sub(/[[:space:]]+$/, "", physical_line)
      # A trailing note after the pipe is a bash comment, so the pipeline stays open across it.
      sub(/\|[[:space:]]*#.*$/, "|", physical_line)
      if (joined_command == "") { first_line_number = NR }
      if (physical_line ~ /\|$/ || physical_line ~ /\\$/) {
        sub(/\\$/, "", physical_line)
        joined_command = joined_command physical_line
        next
      }
      joined_command = joined_command physical_line
      if (feeds_an_early_leaving_grep(joined_command)) { printf "%d: %s\n", first_line_number, joined_command }
      joined_command = ""
    }
    END {
      if (joined_command != "" && feeds_an_early_leaving_grep(joined_command)) {
        printf "%d: %s\n", first_line_number, joined_command
      }
    }
  ' "$1"
}

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
