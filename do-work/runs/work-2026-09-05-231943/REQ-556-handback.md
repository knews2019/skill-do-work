# Builder hand-back — REQ-556 (cut the debug-artifact rule prose that `do-work-cli qualify` already enforces)

**Branch:** `worktree-agent-REQ-556-cut-debug-artifact-prose`
**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-556-cut-debug-artifact-prose`
**Base:** `d3ceca3` · **Branch head:** `2eeb880f7881e50c928471c396e21b670324d44a`
**One commit.** Nothing pushed. Nothing written in the main checkout except this file.

Done. The debug-artifact and P-A-U-honesty rule is now stated once, in code. Four prose
restatements are gone, one pointer sentence names the finding codes, and a counting lock-in
holds the mention count at 2. The two mentions that are not restatements survived — the
count came out **at** the ceiling, not under it.

## File manifest

| File | Verb | What |
|---|---|---|
| `skills/do-work/actions/work.md` | modified | four edits: one four-word trim, two deletions, one added sentence |
| `skills/do-work/actions/review-work.md` | modified | one Red Flags bullet deleted |
| `skills/do-work/actions/work-reference.md` | modified | one anti-rationalization table row deleted |
| `_dev/tests/audit-lockins.sh` | modified | one counting assertion block added before the footer |

`git show --stat`: 4 files changed, 36 insertions(+), 6 deletions(-). Prose net delta is
**−4 lines**, not the exploration's forecast of ≈ −6, because the added sentence was appended
to an existing paragraph and costs no line of its own.

Not touched: `internal/corehelpers/checks.go`, the three other shipped files carrying the
byte-identical P-A-U template payload, any action file outside the three named,
`_dev/tests/contract-regressions.sh`, `_dev/tests/contracts/probe-lanes.sh`.

## The edits, by their text

**E1 — `work.md`, P-A-U phasing pointer bullet (trim).**
`[UNIFY] runs \`git diff --stat\`, runs native linters, verifies no debug artifacts, and lists each file checked`
→ `[UNIFY] runs \`git diff --stat\`, runs native linters, and lists each file checked`.
The bullet's own preamble calls this block "pointers — the underlying rules live in the loaded
crew-members files", and the real UNIFY instruction ships in the template payload inside every
REQ file. The four words were a third copy.

**E2 — `work.md`, Common Rationalizations row deleted.**
`| "P-A-U is bookkeeping — I'll just tick the boxes" | … | A checked \`[UNIFY]\` over a diff containing \`console.log\` is a false claim the qualifier will catch |`
The row's stated consequence is wrong: a `console.log` in the diff produces
`QUALIFY-LIBRARY-OUTPUT` (or `-OUTPUT-RELOCATED` / `-REPORTER-OUTPUT`), not
`QUALIFY-DEBUG-ARTIFACT`, and two of those three are warnings.

**E3 — `work.md`, Red Flags bullet deleted.**
`- All P-A-U checkboxes marked complete but diff contains \`console.log\`, \`debugger\`, or \`TODO\` (debug artifacts)`
Restates an error-severity check in the Red Flags list of the orchestrator who runs that check.

**E4 — `work.md`, one sentence added** at the end of the paragraph under
`### Qualification and Testing Judgment` (the one ending "…merely because a grep cannot see
their consumers."). It names `QUALIFY-DEBUG-ARTIFACT`, `QUALIFY-PAU-UNCHECKED`,
`QUALIFY-UNIFY-DISARMED` and the separate `QUALIFY-LIBRARY-OUTPUT` /
`-OUTPUT-RELOCATED` / `-REPORTER-OUTPUT` family, and ends "Judge those findings; do not restate
the rule they enforce." It deliberately contains neither `console.log` nor the phrase `debug
artifacts`, so it consumes no headroom under the lock-in.

**E5 — `review-work.md`, Red Flags bullet deleted.**
`- Builder checked all P-A-U boxes but the diff contains \`console.log\`, \`debugger\`, or TODO/FIXME`
Duplicates the weakest half of the hygiene bullet at line 106 and carries none of its extra content.

**E6 — `work-reference.md`, one table row deleted.**
`| "The builder checked the UNIFY box" | Read the actual diff for debug artifacts | A checked box is a claim, not a fact |`
The heading `## Qualification Anti-Rationalization Table (Step 6.3)` and the other five rows are
untouched — `_dev/tests/contracts/core-checks.sh:718` pins that heading.

## The two mentions that survived, and why the count is 2 rather than lower

```
skills/do-work/actions/review-work.md:106  Diff hygiene — no debug artifacts — console.log/print
                                           lines no contract reads … or temporary files left
                                           behind. **Protect lessons learned** — …
skills/do-work/actions/review-work.md:374  - [ ] **[UNIFY]:** (Agent: … Verify no debug
                                           artifacts in diff. …)
```

Line 106 is a read the canonical `qualify` never makes: standalone review runs on a historical
commit `qualify` never saw, `review-work.md` never mentions `qualify` or `advance` at all, and
temporary files left behind and comment intent have no Go equivalent anywhere in `checks.go`.

Line 374 is emitted template payload inside a fenced block, byte-identical in four shipped
files. Verified after the cut: the exact UNIFY sentence appears in `review-work.md`,
`capture-reference.md`, `sample-archived-req.md` and `do-work-toolbox/actions/code-review.md`,
and `sort -u` over all four copies still returns **1** distinct string.

## Red / green / red — recorded verbatim

**Baseline, before the lock-in existed** (proves the file was green to start with):

```
$ bash _dev/tests/audit-lockins.sh
Audit lock-in regressions passed.
EXIT=0
```

**RED — the new block alone, against the base revision's uncut prose.** The three action files
were still byte-identical to `d3ceca3` (`git diff --stat -- skills/` was empty at this point);
only `_dev/tests/audit-lockins.sh` carried the new assertion:

```
$ bash _dev/tests/audit-lockins.sh
FAIL: 7 debug-artifact rule mentions across work.md, review-work.md and work-reference.md; ceiling is 2 (do-work-cli qualify owns the rule)
EXIT=1
```

**GREEN — after the six prose edits:**

```
$ bash _dev/tests/audit-lockins.sh
Audit lock-in regressions passed.
EXIT=0
```

**RED again — one restatement pasted back, once into each of the three files** (the same line
each time, appended and then removed; each file's sha256 was compared before and after and
matched):

```
=== paste one restatement back into skills/do-work/actions/work.md ===
FAIL: 3 debug-artifact rule mentions across work.md, review-work.md and work-reference.md; ceiling is 2 (do-work-cli qualify owns the rule)
EXIT=1
restored: e87744dfbe3d == e87744dfbe3d
=== paste one restatement back into skills/do-work/actions/review-work.md ===
FAIL: 3 debug-artifact rule mentions across work.md, review-work.md and work-reference.md; ceiling is 2 (do-work-cli qualify owns the rule)
EXIT=1
restored: 5c1279b6741f == 5c1279b6741f
=== paste one restatement back into skills/do-work/actions/work-reference.md ===
FAIL: 3 debug-artifact rule mentions across work.md, review-work.md and work-reference.md; ceiling is 2 (do-work-cli qualify owns the rule)
EXIT=1
restored: d14e5b9b363e == d14e5b9b363e
```

**Fourth proof — the missing-file property**, run against the post-cut tree with
`work-reference.md` moved aside and moved back (sha256 identical after restore):

```
FAIL: debug-artifact prose lock-in cannot read skills/do-work/actions/work-reference.md; the file moved and the lock-in is dead
EXIT=1
```

A rename fails loudly instead of counting zero and going green.

## Test commands and exit lines

| Command | Result | Exit |
|---|---|---|
| `bash _dev/tests/audit-lockins.sh` | `Audit lock-in regressions passed.` | 0 |
| `bash _dev/tests/action-shell-blocks.sh` | `Shell-block lint passed: 74 fenced blocks and 33 shipped shell files; ShellCheck enabled.` | 0 |
| `bash _dev/tests/contract-regressions.sh` | `Contract regression checks passed.` | 0 |
| `bash _dev/tests/prescribed-shell-canonicalization.sh` | `Prescribed shell primitive canonicalization checks passed.` | 0 |
| `env DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh` | `staged skills contract: PASS` | 0 |

All five in the sanitized environment the brief specifies. The full canonical gate was not run;
the five targeted probes cover the change and the machine is shared.

`contract-regressions.sh` matters most here because it runs `_dev/tests/contracts/core-checks.sh`,
which is what pins the `work-reference.md` table heading and forbids restoring a `### Step 6.3:`
heading to `work.md`. It stayed green, so nothing deleted was load-bearing for a pin.

Additional checks: `bash -n` clean on the modified script; `shellcheck -S style` reports five
findings, all pre-existing in blocks from REQ-550/551/554 that this REQ did not touch, none in
the new block. No added line in the diff carries `TODO`, `FIXME` or `debugger` as a whole word,
so this REQ's own `qualify` gate cannot go red on its own deletions.

## Where the request and the exploration disagreed — the exploration won every time

| # | The request says | HEAD says | What I did |
|---|---|---|---|
| 1 | nine mentions across the three files (work.md 5, review-work.md 3, work-reference.md 1) | **seven** (3 / 3 / 1) | re-baselined to 7; the RED above is 7, not 9 |
| 2 | Reproduce at audited commit `dc8a64e3` | `git cat-file -t dc8a64e3` → `fatal: Not a valid object name`; this clone does not contain it | the captured RED cannot be replayed; I proved RED against `d3ceca3` instead and did not manufacture old lines to reach 9 |
| 3 | "Keep one sentence in `work.md` Step 6.3 naming the three finding codes" | `rg 'QUALIFY-' skills/ --glob '!*.go'` → zero hits; the sentence does not exist in shipped prose | it is an **addition**, not a retention. I wrote it, and the replay case that earns it is below |
| 4 | put it under Step 6.3 | `### Step 6.3:` no longer exists in `work.md`, and `core-checks.sh:701` FAILs if it is restored | put it under `### Qualification and Testing Judgment` |
| 5 | `console.log` ⇐ `QUALIFY-DEBUG-ARTIFACT` | `checks.go:392` matches only `\b(debugger\|TODO\|FIXME)\b`; `console.log` and `print(` go through `checks.go:377` → `QUALIFY-LIBRARY-OUTPUT` / `-OUTPUT-RELOCATED` / `-REPORTER-OUTPUT` | the added sentence states the real taxonomy; I did not copy the request's mapping |
| 6 | preserve the "judge entry-point or dynamic-wiring exceptions" deferral | that belongs to `QUALIFY-NEW-FILE-UNWIRED` (`checks.go:347`), a different check, whose prose home is out of scope | left it alone, unchanged |
| 7 | "`contract-regressions.sh` may pin some of these sentences; delete the matching predicates in the same commit" | it pins **none** of them — grep for `console\.log`, `debug artifact`, `UNIFY`, `Anti-Rationalization`, `Red Flags` and the three filenames all return 0 | **the constraint has no referent.** No predicate was deleted. Its 77-line self-ratchet is also exactly at its ceiling, so nothing could have been added there either |
| 8 | expected net line delta −15 | that forecast assumed nine sites | actual prose delta is **−4** |
| 9 | "if queued REQ-510 already removed the `work-reference.md` site, skip it" | REQ-510 is `completed` and the row was still there verbatim | REQ-556 owns the row; it is cut here |
| 10 | lock-in limit "≤ 3 after this REQ (today 9)" | today was 7, and 2 is reachable without cutting either non-restatement | ceiling set to **2**, inside the request's ≤ 3 |

**Why the added sentence is earned rather than a violation of "no replay case, no addition"
(`crew-members/maintenance.md:26`):** `work.md` invokes the CLI with `--format json` at two
places, and `checks.go:451-473` is the *text* renderer only. In JSON the orchestrator reads the
raw `finding.Code` token. Before this change, an orchestrator seeing the literal string
`QUALIFY-DEBUG-ARTIFACT` had zero shipped prose routing it to an owner. Fails without the
sentence, passes with it.

## Anchors that had moved

Located every edit by its text. One anchor had moved:

- `skills/do-work/actions/work-reference.md` — the anti-rationalization row the exploration
  records at **604** is at **612** at `d3ceca3`, and its table heading moved 598 → 606.
  `REQ-486` (collapsible UR progress summaries) added the composed-summary row table and its
  "Four rows carry a judgment…" paragraph above it, pushing this section down 8 lines.

Everything else held: `work.md` 292 / 574 / 583 and `review-work.md` 106 / 374 / 494 are exactly
where the exploration put them, and so is the `### Qualification and Testing Judgment` heading
region.

One thing that looked like drift and is not: `git diff`'s hunk header for the
`work-reference.md` deletion prints "Four rows carry a judgment no typed record makes for
you…". That is git's function-context heuristic reaching back to REQ-486's unrelated paragraph
at line 515, which counts rows in the **composed exit summary** table, not in the
anti-rationalization table. No row-count prose anywhere refers to the table I edited — I grepped
every reference to `Qualification Anti-Rationalization` in the repo and there are three: the
`work.md` pointer, the heading itself, and `core-checks.sh:718`.

## Sweep for restatements outside the three declared files

Per `prime-action-files.md` Traps, family `alternate-writer-contract-drift`. After the cut, the
only shipped hits for `console.log` or `debug artifacts` outside the three files are:

- the four-way P-A-U template payload (`capture-reference.md:104`,
  `sample-archived-req.md:33`, `do-work-toolbox/actions/code-review.md:312`) — deliberately
  identical, out of scope;
- `do-work-toolbox` material on a different contract: `inspect.md:150` and `:288`,
  `quick-wins.md:69`, and `crew-members/debugging.md:70` (mirrored into the toolbox).

None is a restatement of the qualification rule this REQ cut, and none is in scope.

## Decisions

**D-01 — ceiling is 2, not 3.** The exploration offered 3 as the fallback if E1's four words
were kept. Keeping them was not defensible: the bullet is explicitly a pointer block, and the
real instruction ships in the template payload. Setting the ceiling at 2 leaves the lock-in with
zero headroom, which is the point — a fourth mention of any wording is a regression.
Reversible, low reach. DECIDE & STATE.

**D-02 — the added sentence names the output-primitive family as three codes, not "a pair".**
The exploration's paste-ready sentence said "a separate pair"; `checks.go` has three
(`QUALIFY-LIBRARY-OUTPUT`, `QUALIFY-OUTPUT-RELOCATED`, `QUALIFY-REPORTER-OUTPUT`), plus
`QUALIFY-DEBUG-ARTIFACT-RELOCATED` beside the main code. Shipping "a pair" would have put a
fresh wrong count into the prose this REQ exists to de-duplicate. I wrote the codes instead of a
count. DECIDE & STATE.

**D-03 — the sentence was appended to an existing paragraph, not added as its own line.** It
belongs to the same judgment the paragraph already describes, and a standalone line would read
as a new rule rather than a pointer. This is why the line delta is −4 rather than −6.
DECIDE & STATE.

## Lesson evidence

Read whole-satellite: `_dev/primes/prime-action-files.md`, `_dev/primes/prime-shell-commands.md`,
`_dev/primes/lessons-action-files.md`, `_dev/primes/lessons-shell-commands.md`.

The REQ's frontmatter recorded `lessons-action-files.md` as **dropped for budget** (3968 tokens
against a 2000-token limit, `slugged: partial` so no targeted form). I read it anyway, because
the prime's own Read-first condition is "modifying any action file" and that is exactly what
this REQ does. It paid for itself: the `alternate-writer-contract-drift` family (REQ-566,
REQ-531, REQ-498) is what drove the outside-scope sweep recorded above, and REQ-262's rule
against putting a count in a rule sentence is what drove D-02.

From `lessons-shell-commands.md`, three entries bear directly on the lock-in and all three are
satisfied: REQ-257 (a lock-in that enumerates known spellings reads as coverage while being
unable to fail — this one counts), REQ-298 (a fallback in the same domain as a legitimate result
launders a failure into data — hence the missing-file FAIL), and REQ-234 (a derived count that
almost reproduces the remembered figure is the trap — 7 was re-derived, not adjusted toward 9).

No required-lesson path was missing.

Crew rules read: `general.md`, `shared-principles.md`, `coding-guardrails.md` (§ 5 Naming for
Reach — the block's identifiers are `debug_rule_mention_ceiling`, `debug_rule_scanned_files`,
`debug_rule_mention_count`, `debug_rule_file`, `debug_rule_file_hits`, all two-word and
plain-text findable), `communication-style.md`, `anti-slop.md`.

## Not done here — for finalization

**This commit is a release** (`_dev/primes/prime-releases.md:5` — any commit changing shipped
files under `skills/`). It needs a `CHANGELOG.md` entry, a `skills/do-work/VERSION` bump (at
`0.303.10` on this branch), and the byte-identical mirror to `skills/do-work/CHANGELOG.md`
enforced by `shipped-package-reference-contract.sh`. None of those paths is in this REQ's
`write_set`, and the brief assigns them to finalization. **Not done.**

## Discovered tasks

**DT-01 — `work-reference.md` still cites a retired step number in three places.**
Lines 422 and 433 cite "Step 6.3", and the table heading at line 606 carries it too, but
`work.md` has had no `### Step 6.3:` heading for some time. The heading cannot simply be
reworded: `_dev/tests/contracts/core-checks.sh:718` greps for the exact string
`## Qualification Anti-Rationalization Table (Step 6.3)`, so the heading and the contract have
to change in one commit. Impact: negligible, reader-confusion only. Not fixed here — out of
this REQ's scope and it needs its own lock-in change.

**DT-02 — the request's `## Constraints` carried an instruction with no referent.**
"`_dev/tests/contract-regressions.sh` may pin some of these sentences; delete the matching
predicates in the same commit" — it pins none of them, and the file is at its own 77-line
ceiling so nothing could be added there either. Worth noting for whoever writes the next
audit-derived REQ: that constraint appears to have been carried forward by template rather than
checked.

## Acceptance criteria

- [x] The rule is stated once, in code, with the action files pointing at it — `work.md` now
      names the finding codes and says to judge them rather than restate them
- [x] The two mentions that are not restatements survive (`review-work.md:106` and `:374`)
- [x] A lock-in fails if a restatement returns, counted rather than name-listed — proved red
      from all three files independently
- [x] A renamed or missing target file fails the assertion loudly instead of counting zero —
      proved
- [x] The gate is green — five suites, all exit 0
