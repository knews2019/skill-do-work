# REQ-234 Hand-Back — Stop the Shell Behavior Suite Counting Its Own Cases

**Branch:** `worktree-agent-REQ-234-stop-shell-behavior-count`
**Commit:** `88b083a` — `[REQ-234] Derive the shell behavior suite's case count at run time`

## File Manifest

| Action | File | What changed |
|---|---|---|
| modified | `_dev/tests/prescribed-shell-scripts-behavior.sh` | Replaced the hand-maintained `(47 named script cases)` literal with a count grepped out of the file at run time, and stated the counting rule — one case = one `# <script-name>: <claim>` header comment at column zero — in a comment beside it. Added one `suite_file` variable next to the existing `repo_root` derivation so the count is readable from any caller's directory. |

Nothing outside the declared write set was touched. `git status --porcelain` showed exactly one modified path.

## RED Evidence

The literal, and every count actually derivable from the file:

```
== RED: the literal on the closing line ==
1099:printf 'Prescribed shell script behavior probes passed (47 named script cases).\n'

== RED: counts actually derivable from the file ==
case-header comments (^# <script-name>: ): 42
all column-zero comment lines:              115
fail_case invocations:                      179
distinct "<script> <name> case" phrases:     46

literal claims 47; no count above equals 47.
```

The 46 is the closest near-miss, so I probed it specifically to see whether the
original convention was recoverable. It is not: several named units in the suite
label their assertions `replay` rather than `case` and so fall outside that count
entirely —

```
== fail_case messages with no ' case' token (sample) ==
ai-report all-failed batch replay did not wait for every launched job
ai-report interrupted batch replay leaked invocation-private staging
ai-report mixed batch replay returned nonzero
ai-report publish-collision replay leaked invocation-private staging
install-last30days sentinel-only check accepted a missing runtime script
run-blocked-check process-tree cleanup killed the test runner group
```

There are 23 such distinct messages across 8 named families. No rule over this file yields
47, which independently confirms REQ-229's D-02: the convention that produced 45
is unrecoverable.

## GREEN Evidence

This is the deletion branch of the Finding-Closure Ratchet — the named surface is
the literal, and it is gone:

```
== GREEN: no literal remains on the closing line ==
1106:printf 'Prescribed shell script behavior probes passed (%s named script cases).\n' "$named_case_count"

== GREEN: no hardcoded 47 anywhere in the file ==
(no match — the literal is gone)

== the counting rule still counts exactly the 42 case headers (new comment does not self-match) ==
42
```

Suite run, from the repo root:

```
Prescribed shell script behavior probes passed (42 named script cases).
exit=0
```

And from an unrelated cwd, to prove the `suite_file` derivation is not cwd-sensitive:

```
$ cd <scratchpad> && bash <abs-path>/_dev/tests/prescribed-shell-scripts-behavior.sh
Prescribed shell script behavior probes passed (42 named script cases).
exit=0
```

ShellCheck at warning level over the file: exit 0, no output.

## The New Closing Lines, Verbatim

```bash
# One case is one fixture block, and every block opens with a header comment of the shape
# `<script-name>: <what it proves>` at column zero. That shape is the definition, and the
# count below is that shape grepped out of this file at run time — so the reported number
# and the file cannot disagree, and nothing here is a remembered figure.
named_case_count="$(grep -cE '^# [a-z0-9][a-z0-9-]*: ' "$suite_file")"
printf 'Prescribed shell script behavior probes passed (%s named script cases).\n' "$named_case_count"
```

Supporting line, added beside the existing `repo_root` derivation at the top:

```bash
# Absolute path to this file, so the closing count can read it from any caller's directory.
suite_file="$repo_root/_dev/tests/${BASH_SOURCE[0]##*/}"
```

## maintainer-verify

```
Maintainer verification passed.
MAINTAINER_VERIFY_EXIT=0
```

Run unpiped from the worktree root, with `echo $?` on its own line. Its output
includes the suite's new closing line, `Prescribed shell script behavior probes
passed (42 named script cases).`

## Decisions

- **D-01: Computed the count rather than dropping the claim.** Both branches close
  the REQ, so I tested the drop branch against `maintenance.md`'s own standard —
  *"if removing something would have to be restored next week, it was foundation,
  not bloat."* It fails that test. The repo's history shows the number in active
  use as evidence of suite growth, four separate times, by hand and in shipped
  release notes: `REQ-220`'s hand-back (`exit 0 (44 named script cases)`),
  `REQ-221`'s hand-back (`45 named script cases`), and two `CHANGELOG.md` entries
  (`41 → 44 named script cases`, `27 → 29 named script cases`). That is a standing
  maintainer habit reaching for "how big is this suite", not decoration. Deleting
  the count would remove something demonstrably wanted; computing it serves the
  same want and makes it true for the first time.

  The delete-before-you-add reflex is still honored, because the *addition* here
  is not a new rule — `maintenance.md`'s worked example for this exact pattern
  ("a stale enumeration") prescribes generalizing the enumeration to a trigger
  condition so there is nothing to keep in sync, which is precisely what a
  run-time count is. What got deleted is the hand-maintained number and the
  obligation to remember it.

- **D-02: Defined one case as one case-header comment (42), not one `fail_case`
  name (46).** The brief warned that a fragile grep silently disagreeing with a
  human count is worse than no count, so I checked both candidate definitions
  against the file. The header-comment shape matches 42 lines with **zero false
  positives** — I listed all 42 and every one is a genuine case header; the
  convention is already 100% uniform across the file, so the rule is descriptive
  of what is there, not imposed on it. It is also what a human counts scanning
  the file: one fixture block, one header. The `fail_case`-name definition is the
  fragile one — it would silently miss the 23 `replay`-labelled assertions shown
  in the RED evidence, disagreeing with a human count in exactly the way the brief
  forbids. Rejected on that basis.

- **D-03: The reported number changes 47 → 42, deliberately.** The REQ authorizes
  picking a definition, and any honest definition moves the number, since 47
  matched none. Worth knowing when reading old hand-backs: counts recorded before
  this commit are on the old, undefined convention and are not comparable with
  counts after it. I did not edit any historical record to match — those are
  history, and rewriting them would be worse than the discontinuity.

- **D-04: Added one variable rather than using `${BASH_SOURCE[0]}` inline.** The
  suite's only two `cd` calls are both inside subshells, so a bare
  `${BASH_SOURCE[0]}` would work today — but it would be the file's single
  cwd-sensitive read, and a future top-level `cd` in this fixture-heavy suite
  would break it silently. One line closes that, reusing the `repo_root` already
  derived two lines above.

## Integration Seams

None. The change is confined to one file, nothing sources or parses its closing
line (`staged-skills-contract.sh:183` invokes the suite and consumes only its exit
status), and no shared file needs a line from me.

## Discovered Tasks

- **[low]** `_dev/tests/prescribed-shell-scripts-behavior.sh` names its assertion
  units two ways: most say `<script> <name> case ...`, but 23 distinct messages
  across 8 families (`ai-report` all-failed / interrupted / mixed /
  publish-collision batch replays, `generate-report-image-batch` usage-error,
  `install-last30days` complete-source / repaired-tree / sentinel-only checks,
  `run-blocked-check` process-tree cleanup) say `replay`, `check`, `cleanup`, or
  the plural `cases` instead of the singular `case` token. This split
  vocabulary is what made the original count unrecoverable, and it is why the
  count now keys on the header comment rather than the assertion label. Not fixed
  here: renaming 23 assertion messages is outside this REQ's write-set intent and
  would touch test text, which the REQ forbids weakening. Worth a follow-up only
  if someone wants the assertion labels to be countable too.
