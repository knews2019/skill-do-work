#!/usr/bin/env bash
# quiet-grep-pipeline-scanner.sh — the one definition of the quiet-grep pipeline scanner.
# Sourced by quiet-grep-pipeline-audit.sh, whose two-sided fixture pins every shape it must
# and must not flag, and by action-shell-blocks.sh, which runs it over every shipped Markdown
# shell fence and shipped shell file under skills/. One body, so the two walks cannot disagree
# about what the shape is. This file runs nothing itself.
#
# `writer | grep -q PATTERN` under `set -o pipefail` reports the writer's SIGPIPE death as
# grep's verdict: the assertion is then wrong in both directions — a positive matcher misses a
# pattern that is present, and a negative matcher fails to flag one it should. Timing decides,
# not content, with the flip size depending on how the writer flushes, so a check written this
# way passes for years and then fails, or passes forever while asserting nothing (REQ-593).
#
# quiet_grep_pipeline_offenders prints every logical line in the named file that feeds an
# early-leaving grep from a pipeline. The defect is a condition, not a spelling: a reader that
# can exit before its writer is done, downstream of a pipe, under pipefail. So the scan works on
# the condition's three ingredients instead of on the one spelling that was converted —
# ordinary shell writes the same defect many different ways in the must-flag fixture in
# quiet-grep-pipeline-audit.sh alone, and a regex naming only today's shape goes stale the first
# time someone writes another
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
# pinned as a must-not-flag shape in quiet-grep-pipeline-audit.sh.
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
