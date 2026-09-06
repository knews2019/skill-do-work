#!/usr/bin/env bash
# Audit lock-in regressions. Pinned after maintainability audits.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
failure_count=0

# Finding 10: exported one-line delegates with no production caller (REQ-550)
delegate_findings="$(
  for d in "$repo_root/skills/do-work/tools/do-work-cli" "$repo_root/skills/do-work-board/tools/queue-kanban"; do
    find "$d" -name '*.go' ! -name '*_test.go' -exec awk '
      FNR==1{f=FILENAME}
      /^func [A-Z]/{
        fn=$0; l=FNR; getline a; getline b;
        if (b ~ /^}$/ && a ~ /^\t(return )?[a-zA-Z0-9_.]+\(/) {
          match(fn,/func [A-Za-z0-9_]+/);
          print substr(fn,RSTART+5,RLENGTH-5) "\t" f ":" l
        }
      }' {} \;
  done | while IFS=$'\t' read -r name loc; do
    [ -z "$name" ] && continue
    if [ "$name" = "AssociateProjectPaths" ] || [ "$name" = "ApplyTimestampPlan" ] || \
       [ "$name" = "DownloadAtomic" ] || [ "$name" = "CheckGreenGate" ]; then
      printf '%s\t%s\n' "$name" "$loc"
    else
      target_line="$(sed -n "${loc##*:}"'{n;p;}' "${loc%%:*}")"
      if echo "$target_line" | grep -Eq '^\t(return )?[a-z][a-zA-Z0-9_]*\('; then
        prod=$(rg -c "\b$name\(" --glob '*.go' --glob '!*_test.go' "$repo_root/skills/" 2>/dev/null | grep -v "${loc%%:*}" | wc -l | tr -d ' ')
        [ "$prod" -eq 0 ] && printf '%s\t%s\n' "$name" "$loc"
      fi
    fi
  done
)"

if [ -n "$delegate_findings" ]; then
  while IFS=$'\t' read -r name loc; do
    [ -z "$name" ] && continue
    printf 'FAIL: %s (%s) is an exported one-line delegate with no production caller\n' "$name" "$loc" >&2
    failure_count=$((failure_count + 1))
  done <<< "$delegate_findings"
fi

# Finding 5: toolbox-shims-no-callers (REQ-551)
callerless_shims="$(
  for f in "$repo_root"/skills/do-work-toolbox/scripts/*.sh; do
    [ -e "$f" ] || continue
    b=$(basename "$f")
    n=$(rg -l --fixed-strings "$b" "$repo_root/skills" "$repo_root/tools" "$repo_root/suite" "$repo_root/README.md" "$repo_root/CLAUDE.md" "$repo_root/_dev/primes" --glob '!*CHANGELOG*' | grep -v "/$b$" | wc -l | tr -d ' ')
    [ "$n" -eq 0 ] && echo "$f"
  done
)"
if [ -n "$callerless_shims" ]; then
  while IFS= read -r f; do
    [ -z "$f" ] && continue
    printf 'FAIL: caller-less toolbox shell shim found: %s\n' "$f" >&2
    failure_count=$((failure_count + 1))
  done <<< "$callerless_shims"
fi

# Shipped shell delegating check (REQ-551 companion)
non_delegating="$(
  find "$repo_root/skills" "$repo_root/tools" -name '*.sh' | while read -r f; do
    rg -q do-work-cli "$f" || echo "NON-DELEGATING: $f"
  done
)"
if [ -n "$non_delegating" ]; then
  while IFS= read -r line; do
    [ -z "$line" ] && continue
    printf 'FAIL: %s\n' "$line" >&2
    failure_count=$((failure_count + 1))
  done <<< "$non_delegating"
fi

# Finding 8: dead-path-pointers-in-records (REQ-549)
# Topic-index `sources:` lists and the primes are read as live routing, so a token that
# matches no tracked file costs every reader a search. Dated records keep their citations
# as written and are not scanned: decisions/records/, decisions/audits/,
# decisions/imported-specs/ and decisions/log.md cite the tree as it was.
tracked_file_list="$(cd "$repo_root" && git ls-files)"

report_token_if_dead() {
  local candidate_token="$1"
  local citing_file="$2"
  local resolvable_token="$candidate_token"

  # Not a repo-relative file claim, so not checked: a glob pattern, a token rooted
  # somewhere the reader resolves from instead (leading "/"), or a bare `do-work/...`
  # token, which is the consuming project's queue state. Same carve-outs the shipped
  # reference contract draws — see _dev/primes/prime-action-files.md Cross-Referencing.
  case "$candidate_token" in
    '' | *' '* | *'*'* | /* | do-work/*) return 0 ;;
  esac
  while [ "${resolvable_token#../}" != "$resolvable_token" ]; do
    resolvable_token="${resolvable_token#../}"
  done
  # Substring match, the same reading `git ls-files | grep -F` gives: a token is live when
  # some tracked path ends with it. Matched in-shell because `grep -q` closes the pipe
  # early and `pipefail` would then read the writer's SIGPIPE as a miss.
  case "$tracked_file_list" in
    *"$resolvable_token"*) return 0 ;;
  esac
  printf '%s\t%s\n' "$candidate_token" "${citing_file#"$repo_root/"}"
}

dead_routing_tokens="$(
  for index_file in "$repo_root"/decisions/topics/*.md; do
    [ -e "$index_file" ] || continue
    awk '/^sources:/ { in_sources = 1; next }
         /^[a-z_]+:/ { in_sources = 0 }
         in_sources && /^  - / { sub(/^  - /, ""); print }' "$index_file" |
      while IFS= read -r source_entry; do
        report_token_if_dead "$source_entry" "$index_file"
      done
  done
  # A backticked token carrying a directory and a file extension is a path claim.
  for prime_file in "$repo_root"/_dev/primes/*.md; do
    [ -e "$prime_file" ] || continue
    grep -o '`[^`]*`' "$prime_file" | tr -d '`' | grep -E '/.*\.[a-z]{2,4}$' |
      while IFS= read -r prime_citation; do
        report_token_if_dead "$prime_citation" "$prime_file"
      done
  done
)"

if [ -n "$dead_routing_tokens" ]; then
  while IFS=$'\t' read -r dead_token citing_file; do
    [ -z "$dead_token" ] && continue
    printf 'FAIL: %s cites %s, which matches no tracked file\n' "$citing_file" "$dead_token" >&2
    failure_count=$((failure_count + 1))
  done <<< "$dead_routing_tokens"
fi

# Finding 2: cli-launcher-preamble-copied (REQ-553)
# The two exempt paths are named in full, not by trailing path shape: a third copy at
# skills/do-work-toolbox/tools/do-work-cli-preamble.sh is exactly the hand-rolled preamble
# this pins at zero, and a suffix filter would read it as the exempt file.
hand_rolled_preambles="$(
  rg -l --glob '*.sh' 'for cli_candidate in|^launcher_arguments=\(--format text\)$' \
    "$repo_root/skills" "$repo_root/tools" 2>/dev/null \
    | grep -vxF -e "$repo_root/tools/do-work-cli-preamble.sh" \
                -e "$repo_root/skills/do-work/tools/do-work-cli-preamble.sh"
)"
if [ -n "$hand_rolled_preambles" ]; then
  while IFS= read -r f; do
    [ -z "$f" ] && continue
    printf 'FAIL: hand-rolled do-work-cli launcher preamble outside the preamble pair: %s\n' "$f" >&2
    failure_count=$((failure_count + 1))
  done <<< "$hand_rolled_preambles"
fi

# Finding 9: exec-where-pure-go-exists (REQ-552)
# A coreutils subprocess spawned for work the same Go module already does in the standard
# library. Test files are excluded on purpose: fixture setup may shell out, shipped code may
# not — suiteinstall/update_transaction_test.go:25 spawns `cp -R` and is not this finding.
# The command list is the audit's own Reproduce pattern; widening it would drag in
# open/xdg-open/rundll32/ps/tar/sh/python3/curl, which this finding does not address.
# The pattern accepts any context expression, not the literal variable name `ctx`: the same
# module already spells one `invocationContext`, and Naming for Reach pushes new code toward
# the longer form, so an `exec.CommandContext(runContext, "cp", …)` was the regression this
# lock-in could not see. It still requires the coreutils name to be the command argument
# itself — first for exec.Command, straight after the context for exec.CommandContext — so a
# legitimate `exec.Command("git", "rm", …)` is not a false positive.
coreutils_module_directories=(
  "$repo_root/skills/do-work/tools/do-work-cli"
  "$repo_root/skills/do-work-board/tools/queue-kanban"
)
# rg prints nothing and exits 2 on a missing directory, which is indistinguishable from a
# clean scan once only emptiness is read. Ask the question that has a real answer first, so a
# renamed module directory fails here instead of passing forever.
for coreutils_module_directory in "${coreutils_module_directories[@]}"; do
  [ -d "$coreutils_module_directory" ] && continue
  printf 'FAIL: coreutils lock-in cannot scan a missing module directory: %s\n' \
    "${coreutils_module_directory#"$repo_root/"}" >&2
  failure_count=$((failure_count + 1))
done
coreutils_exec_sites="$(
  rg -n 'exec\.Command(Context\([^,]+,|\()\s*"(find|cp|mkdir|grep|sed|ls|rm|mv|cat|touch|head|tail|wc)"' \
    "${coreutils_module_directories[@]}" \
    --glob '!*_test.go' 2>/dev/null
)"
if [ -n "$coreutils_exec_sites" ]; then
  while IFS= read -r coreutils_site; do
    [ -z "$coreutils_site" ] && continue
    printf 'FAIL: coreutils spawned where the module already has pure Go: %s\n' "${coreutils_site#"$repo_root/"}" >&2
    failure_count=$((failure_count + 1))
  done <<< "$coreutils_exec_sites"
fi

# Finding 6: commit-inspect-shared-body (REQ-554)
# The inventory tag legend, the secret-shaped basename patterns, the four file-reading
# bullets, the association semantics and both manual "do it by hand" fallbacks live once, in
# skills/do-work/docs/prescribed-shell-primitives.md. The assertion below scans each action
# file on its own for a load-bearing sentence out of each moved passage, so a paste back into
# either action alone fails.
#
# It replaced a difflib count of lines identical in both actions, held under a ceiling of 30,
# which was wrong in both directions. It could not see the return it was named for: the
# metric counts only lines present in BOTH files, so pasting the tag legend and all four
# reading bullets back into commit.md alone left the count at 30 and exited 0, and restoring
# the whole pre-move commit.md lowered it to 29. And it fired where nothing was duplicated:
# a sequence alignment rematches when a line is removed, so deleting one row of inspect.md's
# own flow diagram scored 33, deleting any single fenced-block opening line scored 32, and
# adding one identical heading plus one body line to both actions scored 32. Headroom would
# only have hidden those trips inside the slack while catching nothing the phrase list below
# misses, so the count is gone rather than re-baselined. Finding new cross-file duplication
# is the maintainability audit's job; this file pins what the audit already found.
#
# What this does not catch: a paraphrase. A moved passage reworded in both actions passes
# every check here, and the audit is what finds that too.
shared_body_action_files=(
  "$repo_root/skills/do-work/actions/commit.md"
  "$repo_root/skills/do-work-toolbox/actions/inspect.md"
)
# A renamed action must fail here rather than scan clean. rg exits 2 on a missing path and
# the status read below reports that, but naming the path first says which one moved.
for shared_body_action_file in "${shared_body_action_files[@]}"; do
  [ -f "$shared_body_action_file" ] && continue
  printf 'FAIL: shared-body lock-in cannot scan a missing action file: %s\n' \
    "${shared_body_action_file#"$repo_root/"}" >&2
  failure_count=$((failure_count + 1))
done
# One sentence per moved passage: two legend rows, the case-insensitivity rule, the first and
# last reading bullets, and two association rules. Illustrative of the moved body, not an
# inventory of it — the guide's own section is the inventory.
moved_shared_body_phrases=(
  '- **M** — modified (a renamed path is tagged M too'
  '- **XD** — deleted secret-shaped path'
  'Secret-shaped matching is case-insensitive'
  '- **Modified files**: Read the `git diff` for each file.'
  '- **Deleted secret-shaped files (`XD`)**: Note only the path'
  '- **Conflict resolution:** a path claimed by two REQs'
  '- **Partial matches count.**'
)
for shared_body_action_file in "${shared_body_action_files[@]}"; do
  [ -f "$shared_body_action_file" ] || continue
  for moved_shared_body_phrase in "${moved_shared_body_phrases[@]}"; do
    moved_phrase_matches="$(rg -n --fixed-strings -- "$moved_shared_body_phrase" \
      "$shared_body_action_file" 2>/dev/null)"
    moved_phrase_scan_status=$?
    if [ "$moved_phrase_scan_status" -gt 1 ]; then
      printf 'FAIL: could not scan %s for moved shared prose (rg exit %s).\n' \
        "${shared_body_action_file#"$repo_root/"}" "$moved_phrase_scan_status" >&2
      failure_count=$((failure_count + 1))
      continue
    fi
    [ -n "$moved_phrase_matches" ] || continue
    while IFS= read -r moved_phrase_site; do
      [ -z "$moved_phrase_site" ] && continue
      printf 'FAIL: %s:%s restates prose that is canonical in skills/do-work/docs/prescribed-shell-primitives.md#protected-inventory-fallbacks; cite that section instead.\n' \
        "${shared_body_action_file#"$repo_root/"}" "${moved_phrase_site%%:*}" >&2
      failure_count=$((failure_count + 1))
    done <<< "$moved_phrase_matches"
  done
done

# The two 207-word inventory fallbacks differed by two relative-path fixups, so no whole-line
# comparison of the two actions ever scored them. This is what pins them at zero.
# rg's own exit status is read rather than a piped count: an awk total prints 0 both when
# nothing matched and when the scan never ran, which prime-shell-commands.md
# § Unchecked Exit Status Reads as Content bans. rg exit 1 is no matches; 2 or more is a
# scan failure and is reported as one. The glob must be '**/actions/*.md' — a single '*'
# does not cross a directory separator, so '*/actions/*.md' matches nothing under skills/.
manual_fallback_matches="$(rg -n --fixed-strings 'If the script is missing or will not run' \
  "$repo_root/skills" --glob '**/actions/*.md' 2>/dev/null)"
manual_fallback_scan_status=$?
if [ "$manual_fallback_scan_status" -gt 1 ]; then
  printf 'FAIL: could not scan shipped actions for manual fallback sentences (rg exit %s).\n' \
    "$manual_fallback_scan_status" >&2
  failure_count=$((failure_count + 1))
elif [ -n "$manual_fallback_matches" ]; then
  while IFS= read -r fallback_site; do
    [ -z "$fallback_site" ] && continue
    printf 'FAIL: manual "do it by hand" fallback remains in a shipped action: %s — the protected-inventory by-hand procedures live in skills/do-work/docs/prescribed-shell-primitives.md#protected-inventory-fallbacks.\n' \
      "${fallback_site#"$repo_root/"}" >&2
    failure_count=$((failure_count + 1))
  done <<< "$manual_fallback_matches"
fi

# Finding 1: qualify-debug-artifact-prose-restated (REQ-556)
# do-work-cli qualify owns the debug-artifact and P-A-U-honesty rule
# (QUALIFY-DEBUG-ARTIFACT, QUALIFY-PAU-UNCHECKED, QUALIFY-UNIFY-DISARMED in
# skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go), so the action files carry
# one pointer instead of a second copy of the rule.
#
# Matches are counted, not lines. review-work.md carries two matched strings on one physical
# line, so a line count turned a pure reflow of that bullet -- splitting it after "no debug
# artifacts --" with no word changed -- into a FAIL that claimed a restatement had returned.
# rg -o emits one row per match, so the number moves only when the words do. rg's own exit
# status is read rather than folded into the count: exit 1 is "no matches", 2 or more is a
# scan that never ran, which _dev/primes/prime-shell-commands.md
# -> Unchecked Exit Status Reads as Content bans reading as zero. No pipeline, so nothing
# can die upstream of a count that still looks legitimate.
#
# The pin is exact, not a ceiling. Above it, a restatement came back. Below it, one of the
# two mentions that are not restatements was cut: review-work.md's standalone-review hygiene
# bullet (a read qualify never makes, because standalone review sees a diff qualify never
# saw) and the emitted P-A-U template payload, which is byte-identical in four shipped files.
# A deliberate change to how often these files name debug markers moves this pin in the same
# commit, with the reason in the commit message.
#
# The pattern set is the marker vocabulary checks.go:24 matches plus the two spellings the
# cut prose used, matched case-insensitively and singular-or-plural, so a reworded paste-back
# is caught and not only a byte-identical one -- capitalised "Debug artifacts", singular
# "debug artifact" and a bullet that keeps only the marker words all hit. The three marker
# words are split across adjacent quoted strings for the reason checks.go splits the same
# three: a literal here would make qualify flag this file's own diff as a debug artifact.
# What this does not catch: a restatement written in neither vocabulary, such as a
# rationalization row about leftover print statements. Finding those is the audit's job.
debug_rule_mention_pin=3
debug_rule_marker_pattern='console\.log|debug artifacts?|\b(debug''ger|TO''DO|FIX''ME)\b'
debug_rule_scanned_files=(
  "$repo_root/skills/do-work/actions/work.md"
  "$repo_root/skills/do-work/actions/review-work.md"
  "$repo_root/skills/do-work/actions/work-reference.md"
)
debug_rule_mention_count=0
debug_rule_mention_sites=()
for debug_rule_file in "${debug_rule_scanned_files[@]}"; do
  if [ ! -f "$debug_rule_file" ]; then
    printf 'FAIL: debug-artifact prose lock-in cannot read %s; the file moved and the lock-in is dead\n' \
      "${debug_rule_file#"$repo_root/"}" >&2
    failure_count=$((failure_count + 1))
    continue
  fi
  debug_rule_matches="$(rg -n -o -i -e "$debug_rule_marker_pattern" -- "$debug_rule_file" 2>/dev/null)"
  debug_rule_scan_status=$?
  if [ "$debug_rule_scan_status" -gt 1 ]; then
    printf 'FAIL: could not scan %s for debug-artifact rule prose (rg exit %s).\n' \
      "${debug_rule_file#"$repo_root/"}" "$debug_rule_scan_status" >&2
    failure_count=$((failure_count + 1))
    continue
  fi
  [ -n "$debug_rule_matches" ] || continue
  while IFS= read -r debug_rule_match_row; do
    [ -z "$debug_rule_match_row" ] && continue
    debug_rule_mention_sites+=("${debug_rule_file#"$repo_root/"}:${debug_rule_match_row%%:*} (${debug_rule_match_row#*:})")
    debug_rule_mention_count=$((debug_rule_mention_count + 1))
  done <<< "$debug_rule_matches"
done
if [ "$debug_rule_mention_count" -ne "$debug_rule_mention_pin" ]; then
  if [ "$debug_rule_mention_count" -gt "$debug_rule_mention_pin" ]; then
    printf 'FAIL: %s debug-artifact rule mentions across work.md, review-work.md and work-reference.md; the pin is %s and do-work-cli qualify owns the rule, so a new mention restates it. Sites:\n' \
      "$debug_rule_mention_count" "$debug_rule_mention_pin" >&2
  else
    printf 'FAIL: the debug-artifact rule mention count fell to %s across work.md, review-work.md and work-reference.md; the pin is %s. The two mentions that are not restatements must survive: review-work.md standalone-review hygiene bullet, and the emitted P-A-U template payload. Sites still present:\n' \
      "$debug_rule_mention_count" "$debug_rule_mention_pin" >&2
  fi
  if [ "${#debug_rule_mention_sites[@]}" -gt 0 ]; then
    for debug_rule_mention_site in "${debug_rule_mention_sites[@]}"; do
      printf '  %s\n' "$debug_rule_mention_site" >&2
    done
  fi
  failure_count=$((failure_count + 1))
fi

# Companion pin for the same finding: work.md names QUALIFY-* codes so the reader can judge
# the findings instead of re-deriving the rule. The prose calls its list illustrative, so
# completeness is deliberately not pinned -- but every code it does name must still exist in
# checks.go, or the pointer that replaced the restated rule is itself the stale prose.
# The lookup is a substring match, so a code that is a prefix of a longer one still finds it
# (QUALIFY-DEBUG-ARTIFACT inside QUALIFY-DEBUG-ARTIFACT-RELOCATED). That pair is emitted by
# one check and only ever moves together, so the looser match costs nothing here.
qualify_code_source="$repo_root/skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go"
if [ ! -f "$qualify_code_source" ]; then
  printf 'FAIL: cannot read %s; the QUALIFY-* codes named in the action files cannot be checked\n' \
    "${qualify_code_source#"$repo_root/"}" >&2
  failure_count=$((failure_count + 1))
else
  for debug_rule_file in "${debug_rule_scanned_files[@]}"; do
    [ -f "$debug_rule_file" ] || continue
    named_qualify_codes="$(rg -n -o -e 'QUALIFY(-[A-Z]+)+' -- "$debug_rule_file" 2>/dev/null)"
    named_qualify_scan_status=$?
    if [ "$named_qualify_scan_status" -gt 1 ]; then
      printf 'FAIL: could not scan %s for QUALIFY-* code names (rg exit %s).\n' \
        "${debug_rule_file#"$repo_root/"}" "$named_qualify_scan_status" >&2
      failure_count=$((failure_count + 1))
      continue
    fi
    [ -n "$named_qualify_codes" ] || continue
    while IFS= read -r named_qualify_row; do
      [ -z "$named_qualify_row" ] && continue
      named_qualify_code="${named_qualify_row#*:}"
      rg -q --fixed-strings -- "$named_qualify_code" "$qualify_code_source" 2>/dev/null
      named_qualify_lookup_status=$?
      [ "$named_qualify_lookup_status" -eq 0 ] && continue
      if [ "$named_qualify_lookup_status" -gt 1 ]; then
        printf 'FAIL: could not search %s for %s (rg exit %s).\n' \
          "${qualify_code_source#"$repo_root/"}" "$named_qualify_code" "$named_qualify_lookup_status" >&2
      else
        printf 'FAIL: %s:%s names %s, which no longer exists in %s; the pointer is stale prose.\n' \
          "${debug_rule_file#"$repo_root/"}" "${named_qualify_row%%:*}" "$named_qualify_code" \
          "${qualify_code_source#"$repo_root/"}" >&2
      fi
      failure_count=$((failure_count + 1))
    done <<< "$named_qualify_codes"
  done
fi

# Finding 7: stale-shell-ownership-prose (REQ-555)
# The executable-homes table exists to route a reader to the home that OWNS a mechanic. When the
# mechanics moved into Go, nine rows still named the retained shell launchers, and one sentence still
# said a six-line launcher orchestrates two other scripts.
#
# Both halves are keyed on the CONDITION, not on the spelling the first version of this check happened
# to delete — a review put the identical defect back seven ways past a negative "no backticked .sh in
# cell two" test, including the table's own `…` house style, an unbackticked path, a shim moved into
# the Mechanics cell, one leading space, an inserted sub-heading, and an emptied table.
#
#   Route rows: every table is found by its own header row rather than by the heading above it, so a
#   sub-heading or a second table cannot hide one. A row is an offender when it names a `.sh` path that
#   is itself a do-work-cli launcher — read from the file, not guessed from the name — so a genuinely
#   shell-owned route would pass and a renamed shim would not. A table with no rows fails too, because
#   an empty table is as silent as a missing one.
#
#   The orchestration claim: any line that names the protected-inventory launcher AND one of the two
#   check launchers it was said to drive. A rewording ("that orchestrates", "coordinates", "a wrapper
#   that orchestrates") is the same false claim, and all three used to pass a fixed-string test.
shell_ownership_guide="$repo_root/skills/do-work/docs/prescribed-shell-primitives.md"
if [ ! -f "$shell_ownership_guide" ]; then
  printf 'FAIL: cannot read %s; the executable-homes table cannot be checked\n' \
    "${shell_ownership_guide#"$repo_root/"}" >&2
  failure_count=$((failure_count + 1))
else
  # Resolves a route path the way the guide writes them: `../../<pkg>/...` is a sibling package under
  # skills/, anything else is relative to the do-work package.
  resolve_guide_route_path() {
    case "$1" in
      ../../*) printf '%s/skills/%s\n' "$repo_root" "${1#../../}" ;;
      *) printf '%s/skills/do-work/%s\n' "$repo_root" "$1" ;;
    esac
  }

  route_row_findings="$(awk '
    /^[[:space:]]*\|[[:space:]]*Canonical executable route[[:space:]]*\|/ {
      inside_route_table = 1; saw_route_table = 1; route_row_count = 0; next
    }
    inside_route_table && /^[[:space:]]*\|[[:space:]]*:?-/ { next }
    inside_route_table && !/^[[:space:]]*\|/ {
      if (route_row_count == 0) printf "emptytable\t%s\t\n", FNR
      inside_route_table = 0
    }
    inside_route_table {
      route_row_count++
      row = $0
      while (match(row, /[A-Za-z0-9_.\/-]+\.sh/)) {
        candidate = substr(row, RSTART, RLENGTH)
        row = substr(row, RSTART + RLENGTH)
        if (candidate ~ /do-work-cli\.sh$/) continue
        printf "candidate\t%s\t%s\n", FNR, candidate
      }
    }
    END {
      if (!saw_route_table) exit 3
      if (inside_route_table && route_row_count == 0) printf "emptytable\t%s\t\n", FNR
    }
  ' "$shell_ownership_guide")"
  route_row_scan_status=$?
  if [ "$route_row_scan_status" -ne 0 ]; then
    printf 'FAIL: no "| Canonical executable route |" table header remains in %s (awk exit %s); the route ratchet cannot run.\n' \
      "${shell_ownership_guide#"$repo_root/"}" "$route_row_scan_status" >&2
    failure_count=$((failure_count + 1))
  elif [ -n "$route_row_findings" ]; then
    while IFS=$'\t' read -r route_finding_kind route_finding_line route_finding_path; do
      [ -z "$route_finding_kind" ] && continue
      if [ "$route_finding_kind" = 'emptytable' ]; then
        printf 'FAIL: %s: the executable-homes table ending at line %s has no rows; an empty table names no home.\n' \
          "${shell_ownership_guide#"$repo_root/"}" "$route_finding_line" >&2
        failure_count=$((failure_count + 1))
        continue
      fi
      resolved_route_path="$(resolve_guide_route_path "$route_finding_path")"
      if [ ! -f "$resolved_route_path" ]; then
        printf 'FAIL: %s:%s names %s as a canonical route and no such file exists.\n' \
          "${shell_ownership_guide#"$repo_root/"}" "$route_finding_line" "$route_finding_path" >&2
        failure_count=$((failure_count + 1))
        continue
      fi
      rg -q --fixed-strings -- 'do-work-cli.sh' "$resolved_route_path" 2>/dev/null
      route_launcher_status=$?
      if [ "$route_launcher_status" -eq 0 ]; then
        printf 'FAIL: %s:%s routes owned mechanics to %s, which is itself a do-work-cli launcher; the row names the subcommand that owns them.\n' \
          "${shell_ownership_guide#"$repo_root/"}" "$route_finding_line" "$route_finding_path" >&2
        failure_count=$((failure_count + 1))
      elif [ "$route_launcher_status" -gt 1 ]; then
        printf 'FAIL: could not read %s to decide whether it is a launcher (rg exit %s).\n' \
          "$route_finding_path" "$route_launcher_status" >&2
        failure_count=$((failure_count + 1))
      fi
    done <<< "$route_row_findings"
  fi

  orchestration_claim_lines="$(awk '
    /scripts\/protected-inventory\.sh/ &&
      (/tools\/checks\/uncommitted-inventory\.sh/ || /tools\/checks\/associate-files\.sh/) { print FNR }
  ' "$shell_ownership_guide")"
  orchestration_scan_status=$?
  if [ "$orchestration_scan_status" -ne 0 ]; then
    printf 'FAIL: could not scan %s for the orchestration claim (awk exit %s).\n' \
      "${shell_ownership_guide#"$repo_root/"}" "$orchestration_scan_status" >&2
    failure_count=$((failure_count + 1))
  elif [ -n "$orchestration_claim_lines" ]; then
    while IFS= read -r orchestration_claim_line; do
      [ -z "$orchestration_claim_line" ] && continue
      printf 'FAIL: %s:%s names the protected-inventory launcher beside a check launcher it was said to drive; the do-work-cli subcommand owns the whole check and launches neither.\n' \
        "${shell_ownership_guide#"$repo_root/"}" "$orchestration_claim_line" >&2
      failure_count=$((failure_count + 1))
    done <<< "$orchestration_claim_lines"
  fi
fi

# Finding 4: per-req-duplicate-go-helpers (REQ-557)
# Six helper names carried fifteen definitions across the CLI's internal packages, each
# duplicate sitting in a package that already imported a package holding an earlier copy.
# Three of the names had copies that DISAGREED: one uniqueSorted dropped empty strings, the
# two compareSemver bodies were inverted, and physicalPath carried two contracts.
#
# The union of the old names and the canonical names that replaced them is pinned at exactly
# seven — a floor as well as a ceiling. A ceiling alone catches a seventh copy but not a
# helper quietly moved back into a per-REQ file under a name this pattern no longer sees; a
# floor alone catches the reverse. Both halves matter, so the count is compared for equality.
#
# This name list is the pinned set itself, so it is exhaustive by construction rather than
# illustrative: adding a name to it means changing the expected count in the same edit.
#
# Both path resolvers are counted and both are legitimate. REQ-557 decided NOT to merge them:
# suiteinstall's existingPhysicalPath treats absence as an error, and that refusal is the
# existence check for the installed skill root, while knowledgecommands' physicalPath walks
# missing ancestors. Two names, two contracts, one definition each — which is why the pinned
# total is seven and not six.
shared_helper_definition_root="$repo_root/skills/do-work/tools/do-work-cli/internal"
shared_helper_expected_definitions=7
shared_helper_definitions="$(rg -n --glob '*.go' --glob '!*_test.go' \
  '^func (uniqueSorted|UniqueSortedStrings|subtractPaths|SubtractStringValues|requestIDLess|RequestIDLess|firstError|FirstNonNilError|compareSemver|CompareSemanticVersions|physicalPath|existingPhysicalPath)\(' \
  "$shared_helper_definition_root")"
shared_helper_scan_status=$?
# rg's status is read, not its output: exit 1 (no match) and exit 2 (could not search) both
# produce empty output, and a ratchet that only judged the text would read a scan that never
# ran as "no duplicates found".
if [ "$shared_helper_scan_status" -gt 1 ]; then
  printf 'FAIL: could not scan %s for shared-helper definitions (rg exit %s); the duplicate-helper ratchet did not run.\n' \
    "${shared_helper_definition_root#"$repo_root/"}" "$shared_helper_scan_status" >&2
  failure_count=$((failure_count + 1))
elif [ "$shared_helper_scan_status" -eq 1 ]; then
  printf 'FAIL: no definition of any shared helper name remains under %s; REQ-557 pinned exactly %s.\n' \
    "${shared_helper_definition_root#"$repo_root/"}" "$shared_helper_expected_definitions" >&2
  failure_count=$((failure_count + 1))
else
  shared_helper_definition_count="$(printf '%s\n' "$shared_helper_definitions" | wc -l | tr -d ' ')"
  if [ "$shared_helper_definition_count" -ne "$shared_helper_expected_definitions" ]; then
    printf 'FAIL: %s definitions of the shared helper names under %s; REQ-557 pinned exactly %s, one per canonical name plus the two deliberately separate path resolvers:\n' \
      "$shared_helper_definition_count" "${shared_helper_definition_root#"$repo_root/"}" \
      "$shared_helper_expected_definitions" >&2
    printf '%s\n' "$shared_helper_definitions" | sed 's|^|  |' >&2
    failure_count=$((failure_count + 1))
  fi
fi
# Finding 7 (second half): shell machinery credited to a Go-owned route (REQ-595)
# Every route in the executable-homes table is a do-work-cli subcommand, so a Mechanics cell that
# names shell machinery is describing an implementation that no longer runs. That is how the
# run-blocked-check cell came to claim "GNU timeout selection and isolated stock-Bash process-group"
# for a command that looks up no timeout binary and builds its group with Go's Setpgid.
#
# WHAT THIS DOES NOT CATCH, stated so its silence is not read as coverage. Two of the three false
# cells REQ-595 fixed named no shell at all: install-memory-hooks credited itself with verification
# and rollback that the knowledge actions still perform by hand, and record-timing-event credited
# itself with a fold that fold-timing-summary owns. A cell that claims another owner's work is
# indistinguishable from a true one by any text scan. Only the shell-machinery shape is guarded here.
#
# The term list is illustrative of the condition, not a closed set: a word that can only describe a
# shell implementation. Measured against the table as it stands: zero matches. Against the table
# before REQ-595: one, naming the cell and its line.
shell_machinery_rows="$(awk '
  /^[[:space:]]*\|[[:space:]]*Canonical executable route[[:space:]]*\|/ { inside_route_table = 1; next }
  inside_route_table && /^[[:space:]]*\|[[:space:]]*:?-/ { next }
  inside_route_table && !/^[[:space:]]*\|/ { inside_route_table = 0 }
  inside_route_table {
    split($0, table_cells, "|")
    mechanics_cell = tolower(table_cells[3])
    if (mechanics_cell ~ /bash|subshell|job control|set -m|gnu |[^a-z]shell[^a-z]/) {
      printf "%s\t%s\n", FNR, table_cells[3]
    }
  }
' "$shell_ownership_guide")"
shell_machinery_scan_status=$?
if [ "$shell_machinery_scan_status" -ne 0 ]; then
  printf 'FAIL: could not scan %s for shell machinery in the Mechanics column (awk exit %s).\n' \
    "${shell_ownership_guide#"$repo_root/"}" "$shell_machinery_scan_status" >&2
  failure_count=$((failure_count + 1))
elif [ -n "$shell_machinery_rows" ]; then
  while IFS=$'\t' read -r shell_machinery_line shell_machinery_cell; do
    [ -z "$shell_machinery_line" ] && continue
    printf 'FAIL: %s:%s credits a do-work-cli subcommand with shell machinery (%s); the route is Go and owns no shell.\n' \
      "${shell_ownership_guide#"$repo_root/"}" "$shell_machinery_line" \
      "$(printf '%s' "$shell_machinery_cell" | sed 's/^ *//;s/ *$//')" >&2
    failure_count=$((failure_count + 1))
  done <<< "$shell_machinery_rows"
fi


if [ "$failure_count" -gt 0 ]; then
  exit 1
fi

printf 'Audit lock-in regressions passed.\n'

