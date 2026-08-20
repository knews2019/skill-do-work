---
id: REQ-301
title: "Let qualify tell a moved line from an added one"
status: completed
created_at: 2026-08-20T08:37:00Z
status_changed_at: 2026-08-20T11:38:27Z
claimed_at: 2026-08-20T12:28:39Z
completed_at: 2026-08-20T12:38:40Z
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-20T12:28:39Z
  basis:
    - Route B
    - 2-file write set
    - shell gate with a design choice between two candidate fixes
user_request: UR-056
addendum_to: REQ-258
domain: general
impact: impact-user-visible
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: [REQ-263, REQ-264]
maintenance: false
write_set:
- skills/do-work/tools/checks/qualify.sh
- _dev/tests/prescribed-shell-cases/qualify.sh
---

# Let Qualify Tell a Moved Line From an Added One

## What

`skills/do-work/tools/checks/qualify.sh` runs `git diff` with no rename or copy detection (`-M`/`-C`) and greps `^+` for debug artifacts. Relocated text therefore reads as newly added, so **every REQ that moves code fails the `[UNIFY]` debug-artifact audit on markers that already existed**.

REQ-258 hit this: it FAILed on four `TODO` strings that are deliberate fixture data in the REQ-254 `qualify` cases and are byte-identical in `git show HEAD:_dev/tests/prescribed-shell-scripts-behavior.sh`. The REQ recorded the override with that evidence and proceeded.

The failure mode is the dangerous direction. A FAIL that is *usually* a false positive on this class of REQ teaches builders to un-check `[UNIFY]` or wave the FAIL away — which is precisely the reflex the audit exists to prevent. A gate that cries wolf on a whole category of change is worse than one that is quiet, because it erodes the response to the true positives.

**Candidate fixes** (pick one, do not build both):
- Enable copy/rename detection for the debug-artifact scan (`git diff -C --find-copies-harder`) so moved hunks are not read as additions.
- Or subtract from the flagged set any line that appears verbatim in the pre-change tree.

The first is git doing the work; the second is cheaper to reason about but can mask a genuinely re-added marker elsewhere in the same file. Whichever ships needs a lock-in case in `_dev/tests/prescribed-shell-cases/qualify.sh` proving that a *moved* marker passes while a *fresh* one still FAILs — the existing REQ-254 case already pins the second half.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-258: `tools/checks/qualify.sh` has no rename/copy detection, so every code-relocation REQ gets a false `[UNIFY]` FAIL on pre-existing debug markers inside the moved text. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

**Answered [2026-08-20]:** User approved at full scope via `do-work clarify`, in the 2026-08-20 remaining-work decision pass, which rated this "Moved code causes false failures, training builders to ignore a real safety gate. Approve after 263/264." The `depends_on: [REQ-263, REQ-264]` gate already encodes that ordering.

## Context

Discovered during REQ-258 (splitting the prescribed-shell behavior suite into per-script files). See that REQ's `## Qualification` section for the worked example and the evidence used to override the FAIL.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Reproduce the false FAIL first. Take candidate 2 (subtract lines already in the pre-change tree) over candidate 1 (`git diff -C`), because git's rename detection is file-level and a hunk moved within one file is a relocation too. Answer the REQ's stated masking risk two ways: test occurrence COUNT rather than presence, and downgrade rather than drop.
- [x] **[APPLY]:** `qualify.sh` — `count_matching_lines_in_tree`, `partition_matched_lines_by_relocation`, and all four flag sites rewired to split fresh from relocated. `_dev/tests/prescribed-shell-cases/qualify.sh` — three cases.
- [x] **[UNIFY]:** `git diff --stat` reviewed; `shellcheck --severity=warning` clean on both files; `maintainer-verify.sh` exits 0; suite 85 → 88 cases. Per-file: **`qualify.sh`** — all four flag sites go through one helper, so the marker and output-primitive scans cannot drift apart; `git grep -e` marks the pattern explicitly so content starting with `-` is data not options (prime, REQ-208); the counting helper's `git grep | awk` carries `|| true` inside the substitution because `git grep` exits 1 on the legitimate zero-match case; no `grep -q` at the end of any new pipe. **`prescribed-shell-cases/qualify.sh`** — the duplicate fixture is byte-identical to the base line on purpose (a near-copy pinned nothing and was caught by mutation testing); the block cleans up its own directory and REQ files so nothing leaks into later cases.

---

## Triage

**Route: B** - Medium

**Reasoning:** The defect is precisely stated, but the REQ offers two candidate fixes, says pick one, and names a risk in the one that is cheaper to reason about. Choosing needed reading how git's rename detection actually scopes and testing the risk rather than accepting the summary of it.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**The false FAIL, reproduced before any edit.** A fixture repo whose suite file carries two deliberate marker payloads, split into per-case files with the text byte-identical:

```
FAIL: [UNIFY] is checked but untracked file tests/alpha.sh carries debug artifacts — un-check it and flag:
  3:  printf "# TODO: alpha fixture\n" > f
FAIL: [UNIFY] is checked but untracked file tests/beta.sh carries debug artifacts — un-check it and flag:
  3:  printf "# FIXME: beta fixture\n" > g
exit=1
```

`git show HEAD:tests/suite.sh` carries both lines verbatim. Two FAILs, zero new markers.

**The reproduction fires through the path REQ-263 just added**, which makes this more urgent rather than less: a relocation into files the builder has not staged yet trips the untracked walk immediately, before any commit exists to inspect.

**Why candidate 1 was not chosen.** `git diff -C --find-copies-harder` detects **file-level** renames and copies. It would have handled REQ-258's shape (one file split into seventeen) and nothing else — a hunk moved inside one file, or between two files that both already existed, produces no rename header and stays a plain add. Those are relocations too, so keying on file topology would have fixed the instance and left the class. The condition worth keying on is the *content*.

**Why the REQ's stated risk in candidate 2 is real — and was caught by a test, not by reasoning.** The first implementation used presence at base as the test: does this exact text exist anywhere in the pre-change tree? That immediately broke REQ-263's own `same-diff-exit` case, whose fixture appends a second `print(raw_text)` to a file the tree already had one in. Presence cannot tell a **move** from a **duplicate**. That is precisely the masking the REQ predicted, and it arrived within minutes of the naive rule.

**The fix for it is a count, not a heuristic.** A move removes one occurrence and adds one, so the total is unchanged. A duplicate raises it, wherever it lands. So the test became: relocated iff the text occurred at base **and** its occurrence count did not increase. The post-change count passes `--untracked` so a relocation into a not-yet-staged file is counted on both sides — REQ-258's shape, and the whole case for this fix.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/checks/qualify.sh` (modify) — occurrence counting, the fresh/relocated partition, and all four flag sites
- `_dev/tests/prescribed-shell-cases/qualify.sh` (modify) — the moved-marker case, the duplicate-marker case, and the fresh-beside-moved case

**Files I will NOT touch:**
- `skills/do-work/actions/work.md` — Step 6.3's prose describes the checks, not their diff mechanics; nothing there says a moved line is an added one
- `tools/checks/scope-drift.sh` — a sibling check with no artifact scan

**Acceptance criteria (restated from REQ):**
- [ ] A relocated marker no longer FAILs the `[UNIFY]` audit
- [ ] A fresh marker still FAILs — the existing REQ-254 case pins this half and must keep passing
- [ ] The masking risk the REQ names is answered, not accepted
- [ ] Exactly one of the two candidate fixes is built
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0

## Pre-Flight

**Git:** ✓ working tree clean outside `do-work/`
**Tests baseline:** ✓ passing (`maintainer-verify.sh` exit 0 at `abf3554`)
**Dependencies:** ✓ go1.26.1, ShellCheck 0.11.0, `just` present

*Checked by work action*

## Decisions

- **D-01 — Candidate 2 (content already in the pre-change tree), not candidate 1 (`git diff -C`).** DECIDE & STATE. Git's rename and copy detection is file-level, so it sees a file split and misses a hunk moved within a file or between two pre-existing files. Keying on content covers every relocation shape at the cost of two `git grep`s per flagged line — and flagged lines are rare, usually zero. The REQ permits either; this one fixes the class rather than the instance.

- **D-02 — The relocation test is occurrence COUNT, not presence.** DECIDE & STATE, and the substantive answer to the risk the REQ raised. Presence at base cannot distinguish a move from a duplicate; the count can, because a move is occurrence-preserving by construction. This was not foreseen and then designed around — the naive presence rule was written first, and REQ-263's `same-diff-exit` case failed within one run because its fixture duplicates a line the tree already had. The test found the flaw the REQ had warned about in prose.

- **D-03 — Relocated lines are downgraded and named, never dropped.** DECIDE & STATE. Even with the count rule, nothing the scans found is discarded: a relocated line prints under a WARN that says the text already exists in the pre-change tree and asks the reader to confirm the move was intended. This is what makes the subtraction safe to trust — the failure mode of a subtracting gate is silence, and there is none here.

- **D-04 — Applied to all four flag sites, not just the marker scan.** DECIDE & STATE. The REQ's text is about `TODO`-class markers, but the condition ("a moved line is not an added line") is indifferent to which regex matched, and the output-primitive scans read `^+` out of the same diffs. Keying on the condition rather than on the reported instance follows `CLAUDE.md`'s "state conditions, not lists"; all four sites share one helper so they cannot drift apart.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/checks/qualify.sh` (modified)
- `_dev/tests/prescribed-shell-cases/qualify.sh` (modified)

**What was done:** Added `count_matching_lines_in_tree` (summing `git grep -c` per file, `--untracked` for the post-change side, `do-work/` excluded on both) and `partition_matched_lines_by_relocation`, which splits any block of matched lines into fresh and relocated by comparing a line's occurrence count at the ownership base against its count now. All four artifact flag sites — diff-stream markers, per-file output primitives, and the two untracked equivalents — now FAIL only on the fresh partition and emit a separate named WARN for the relocated one. Three cases added: a move that must WARN, a byte-identical duplicate that must still FAIL, and a fresh marker beside a moved one in the same file where the FAIL names only the fresh line.

## Qualification

**One FAIL, overridden with evidence; `scope-drift.sh` clean.** Nine flagged lines, all in the two declared files, and **all of them genuinely new text** — which is why this REQ's own fix does not and must not suppress them:

- Six are **fixture payloads** in `_dev/tests/prescribed-shell-cases/qualify.sh` — the marker strings the new cases write into a throwaway repo for the checker to detect. The marker *is* the test input.
- Two are **explanatory comment prose** in `skills/do-work/tools/checks/qualify.sh` (`:275`, `:300`) that name the marker class while describing the fix.
- One is the `'  # TODO: brand new leftover'` payload of the fresh-beside-moved case, which exists precisely so a genuinely new marker still FAILs.

`git diff --name-only -- . ':(exclude)do-work/'` returns only the two declared files. `[UNIFY]` stays checked.

**This is the second consecutive occurrence of the class REQ-263's hand-back predicted, and it is confirmed to be a different class from the one this REQ fixes.** REQ-301 subtracts text whose occurrence count did not grow; these lines are new text whose count went from zero to one, so the correct verdict is exactly the FAIL it gave. Any REQ that adds a marker-bearing fixture to this case file will keep hitting it. Recorded rather than filed, per the 2026-08-20 stopping policy — but it has now happened twice in three REQs, so it is surfaced in the hand-back with a recommendation rather than left in the trail.

**Judgment checks:** *(2) Substantive* — two new helpers and four rewired flag sites in the script; 12 → 15 cases. *(3) Requirements traced* — every acceptance criterion has a named case below, including the "exactly one candidate" one (D-01 records which and why the other was rejected). *(6) Flowing* — the counting helper reads real revisions through `git grep`; the mutation tests could not discriminate the count rule from the presence rule if it were stubbed.

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-cases/qualify.sh`, then `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ `qualify: 15 cases, 0 failures`; ✓ `maintainer-verify.sh` exit 0; suite total 85 → **88 named script cases across 17 per-script files**

**Red-green validation:**

- `REQ-970` moved marker: ✗ before, two FAILs on markers byte-identical in `git show HEAD:tests/suite.sh`, exit 1 → ✓ after, two WARNs naming each moved line as "relocated, not added", exit 0
- `REQ-971` byte-identical duplicate: ✓ FAILs both before and after — the guard on this REQ's own risk. Under the naive presence rule it passed, which is how the flaw was found
- `REQ-972` fresh marker beside a moved one in the same file: ✗ before, everything pooled into one FAIL → ✓ after, the FAIL names only `# TODO: brand new leftover` while the moved marker is reported separately as a WARN
- REQ-254's `REQ-952` fresh-marker case: ✓ passes unchanged — the evidence that the second half of the requirement was not weakened

**Mutation testing — each new case proven able to fail, and each mutation fails a *different* case:**

| Mutation | Result |
|---|---|
| relocation awareness removed (pre-REQ-301 behavior) | `REQ-970` fails (2 assertions) and `REQ-972` fails on the dropped moved marker |
| the count test reverted to bare presence-at-base | `REQ-971` fails, naming occurrence count as the missing test — **and** REQ-263's `REQ-953` fails, which is how the flaw surfaced in the first place |

**A fixture bug caught by that mutation pass, worth recording.** The duplicate case first wrote `> h` where the base line wrote `> g`. Not byte-identical, so it FAILed as *fresh text* under both rules and discriminated nothing — it looked green and pinned nothing. Fixed to be byte-identical, and it now fails under the presence rule as intended. A case that passes for the wrong reason is the failure mode `_dev/primes/prime-shell-commands.md` § Lessons (REQ-257) describes.

**New tests added:** three cases in `_dev/tests/prescribed-shell-cases/qualify.sh` — `REQ-970`, `REQ-971`, `REQ-972`.

**Existing tests updated:** none. All twelve prior cases pass unchanged.

*Verified by work action*

## Review

**Overall: 93%** | 2026-08-20T12:35:50Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 97% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** 3 (report only)
- Two `git grep` invocations per flagged line make the cost proportional to findings rather than to repository size. Flagged lines are normally zero and the cap on printed lines is ten, so the practical cost is nil — but a change that flagged hundreds of lines would run hundreds of tree greps, and nothing bounds that.
- The relocation test is whole-line and exact. A line moved *and* re-indented, or moved with a variable renamed, counts as fresh and FAILs. That is the conservative direction and is the right default, but it means a reformatting-plus-relocation REQ still needs the evidence override.
- `count_matching_lines_in_tree` takes the revision as an empty string to mean "working tree". It reads clearly at both call sites, but an empty-string sentinel is the kind of argument a later caller passes by accident; a named second function would not have that edge.

**Restatement sweep:** the diff changes when a gate fails, so the sweep asked who states that condition. Nothing outside the script describes the artifact scans' add-versus-move semantics — `work.md:446` glosses the check as "debug artifacts in the diff", which stays true (a relocated line is still in the diff; it just no longer fails), and that gloss's precision was already reported as a Minor under REQ-263. `_dev/primes/prime-shell-commands.md`'s REQ-258 lesson says to "expect qualify to FAIL every relocation REQ, because it reads a moved line as an added one" — **that lesson is now stale by design**, and is updated in this REQ's prime-link write rather than left to mislead the next reader.

**Acceptance:** Pass — all five criteria met with named cases, including "exactly one candidate fix" (D-01) and the masking risk answered by construction rather than accepted (D-02, D-03), each with its own failing-mutation proof.

**Suggested testing:** 1 item
- No case covers a relocation in **worktree dispatch mode** (`DO_WORK_DIFF_RANGE` set), where the base is the range's lower bound rather than `HEAD`. The code path is shared and the base resolution is one branch above it, so it is expected to hold; unpinned.

**Code Quality 90%:** four flag sites now each carry a five-line fresh/relocated block, so the same shape is written four times. Extracting a reporting helper would have collapsed them, but each carries a differently-worded message naming its own scan and file, and threading that through a helper would have been more indirection than the repetition costs.

**Follow-ups created:** None; **sweeps appended to:** None

## Lessons Learned

**What worked:** Letting an existing test kill the first design. The naive rule — "text already in the tree is relocated" — is what the REQ itself suggests, and it survived exactly one suite run before REQ-263's `same-diff-exit` case failed it. The REQ had *predicted* that risk in prose ("can mask a genuinely re-added marker"); the test is what turned the prediction into a specific, fixable defect and pointed at occurrence count as the answer. Prose warned; the fixture diagnosed.

**What didn't:** The duplicate-marker case, first time around. It wrote `> h` where the base wrote `> g`, so the text was not byte-identical, it FAILed as fresh under both the old and new rules, and it looked green while discriminating nothing. Only the mutation pass exposed it. When a case exists to pin *why* something fails, the fixture has to differ from the base in exactly the one dimension under test and in no other.

**Worth knowing:** A subtracting gate's failure mode is silence, so this one never subtracts — a relocated line is downgraded to a named WARN and still printed. That pattern is worth copying: when a check learns to excuse something, make it say what it excused. Also: `_dev/primes/prime-shell-commands.md`'s REQ-258 lesson told readers to expect this FAIL on every relocation REQ. It was true when written and is now wrong, which is what the prime-link write in this REQ corrects.

## Orientation

The debug-artifact gate can now tell a line that moved from a line that was written. It compares how many times a flagged line's exact text occurs in the tree before and after the change: unchanged means relocated, so it warns and names the line instead of failing; increased means genuinely added, so it fails as before. A fresh marker sitting beside a moved one in the same file still fails on its own, and nothing found is ever silently dropped. That removes the false FAIL every code-relocation REQ used to hit — the habituation risk this REQ was raised for. Lives in the Step 6.3 qualification gate (`_dev/primes/prime-shell-commands.md`).

Not `[MAP CHANGED]` — one gate's add-versus-move judgment got correct; no contract, field, or caller moved. Staleness spot-check on `_dev/primes/prime-shell-commands.md`: every referenced path resolves, but its REQ-258 lesson ("expect qualify to FAIL every relocation REQ") **is made stale by this change** and the prime-link write below supersedes it rather than stacking a second entry that contradicts the first.

## Prime-Link Write (Step 8 substep 7)

`_dev/primes/prime-shell-commands.md` — staged with this REQ but deliberately **not** listed in `## Implementation Summary`, per work.md Step 9 ("prime files must also be staged — they are part of the REQ's lifecycle changes even though they aren't listed in the Implementation Summary's `Files changed`"). Three writes:

1. This REQ's lesson entry added.
2. **REQ-258's lesson entry amended.** Its second clause told readers to expect this FAIL on every relocation REQ. True when written, wrong now. Superseded in place rather than contradicted by a newer entry sitting above it.
3. **Backfill, stated plainly:** REQ-263's and REQ-264's lesson entries, which their own Step 8 substep 7 should have written and did not. Their links resolve to archived REQs, and leaving the prime missing two lessons to keep a commit boundary tidy would trade real institutional memory for bookkeeping neatness. Recorded here and in the commit message so the backfill is not mistaken for this REQ's own work.
