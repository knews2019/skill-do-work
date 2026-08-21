---
id: REQ-298
title: "[impact-rule-change] Review fix: sweep the unchecked-exit-status primitive across every shipped script"
status: completed
created_at: 2026-08-19T19:43:58Z
status_changed_at: 2026-08-20T08:22:51Z
claimed_at: 2026-08-21T02:30:39Z
completed_at: 2026-08-21T02:48:28Z
kb_status: pending
commit:
route: C
user_request: UR-056
addendum_to: REQ-268
domain: general
review_generated: true
sweep: true
sweep_key: unchecked-exit-status-read-as-content
impact: impact-rule-change
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
depends_on: []
maintenance: false
write_set: []
---

# Review Fix: Sweep the Unchecked-Exit-Status Primitive Across Every Shipped Script

## What

REQ-268 closed five instances of one condition — a command or process substitution whose
exit status is discarded while only its content is judged, so a tool that never ran reads
as a tool that found nothing — and stated that condition in the headers of the two scripts
it touched. Its independent review then found the same primitive **inherited verbatim from
a third script that neither REQ-268 nor its Requirements named**:
`skills/do-work/tools/checks/record-commit-hash.sh` is the file REQ-268's own header cites
as its mandatory guard-style template, and it carries the pattern itself (for example
`head_blob_bytes="$(git cat-file -s … 2>/dev/null || true)"`). The copy direction was
template → copies, so fixing the copies and leaving the template means the next script
written from it re-imports the defect.

**Done means the class cannot recur:** the condition is stated once where shipped shell is
governed, and every shipped script that takes a substitution and judges only its content
either checks the status or says in place why the content alone is sufficient. Patching N
sites one at a time is what this REQ exists to avoid.

## Reproduced During Clarify (2026-08-20)

The REQ was written from the review's static reading. Driving the real script proved the
severity is higher than "the template still carries the pattern", and narrowed which of the
four sites actually matters.

**One site fails OPEN, and it is the incident guard itself** —
`skills/do-work/tools/checks/record-commit-hash.sh:466`:

```bash
head_blob_bytes="$(git cat-file -s "HEAD:$tracked_full_name" 2>/dev/null || true)"
if [ -n "$head_blob_bytes" ] && [ "$head_blob_bytes" -gt 0 ] && ...
```

`|| true` collapses *there is no blob in HEAD* and *git could not answer* into one empty
string, and the `-n` test then skips the truncation floor for both. Against a 10,178-byte
archived REQ truncated to 85 bytes — the shape of the original incident — with everything
real except a failing `git cat-file -s`:

```
control (real git):        FAIL: ... is 85 bytes on disk but 10178 bytes in HEAD —
                           content was lost BEFORE this script ran.        exit 1
cat-file -s failing:       OK: ... all content guards passed.
                           Now stage and commit exactly this file: ...     exit 0
```

It wrote the hash into the remnant and printed the instruction to commit it. That is the
exact failure this script exists to prevent.

**Three sites fail SAFE** (`:189`, `:194`, `:605`): a dead command yields `0`, the
comparison mismatches, and the run stops with a false FAIL. Real instances of the condition,
but they cost a confusing message rather than a lost file. They need the same treatment for
the class's sake; they are not urgent.

## Instances

- [x] `skills/do-work/tools/checks/record-commit-hash.sh:466` — the truncation floor, the
  one reproduced fail-open. Fix it the way `repair-req-timestamps.sh` was fixed under
  REQ-268: ask `git cat-file -e` whether the blob exists, treat an existing blob whose size
  will not read as a failure, and keep the genuine no-blob case as the skip it is. One probe
  in `_dev/tests/record-commit-hash-guards.sh`, which already has the fixture repo for it.
- [x] `record-commit-hash.sh:189`, `:194`, `:605` — the three fail-safe sites. Same
  treatment for the class, or an in-place comment recording that each was judged and why
  the content alone is sufficient there. Either is fine; silence is not.
- [x] Every other shipped script under `skills/*/scripts/` and `skills/*/tools/` — the
  instance list is a sample, not the scope. `_dev/primes/prime-shell-commands.md`
  § Closed Enumerations Go Stale applies: key the rule on the condition and mark this
  list illustrative.
- [x] `_dev/primes/prime-shell-commands.md` itself — the condition belongs in the trap
  list so it is loaded before shell is written, rather than restated per script header.

## Context

Found during the independent review of REQ-268 (finding I1, Important, gate:
`impact-rule-change` — it changes a rule applying across every shipped script rather than
fixing one visible defect). REQ-268 fixed all five sites inside its own two-file write set,
including the three the review surfaced, so nothing named there is left open; what remains
is the spread beyond that write set, which is genuinely a different REQ.

Created `pending-answers` per the generation-≥2 cascade depth stop: REQ-268 itself carries
`review_generated: true`.

## Requirements

- The condition is stated once, where shipped shell is governed, rather than re-derived per
  script header.
- `record-commit-hash.sh`'s own substitutions either check their status or state why the
  content is sufficient on its own.
- A mechanical check that finds the shape, so a new script cannot reintroduce it silently —
  or, if no check can be written without false positives, a stated review convention with
  the reason recorded.

  > **Recorded dissent, overruled by the user at clarify (2026-08-20).** I argued this
  > requirement should be dropped: discarding an exit status is *correct* at most of the
  > sites it appears — `grep -c` on a no-match is the common one, and this very script has
  > two of those with comments explaining why — so a `|| true`-shaped check flags mostly
  > legitimate code and gets muted within a week. My recommendation was to state the
  > condition in `_dev/primes/prime-shell-commands.md` instead, where it is loaded before
  > shell is written. The user chose the full scope. **Build it as written**, and take the
  > escape hatch the requirement already contains honestly: if the check cannot be made to
  > separate the correct uses from the defects, say so with the evidence and fall back to
  > the stated convention — that is a finding, not a failure to deliver.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Red-Green Proof

**RED prompt/case:** A probe in `_dev/tests/record-commit-hash-guards.sh` — **not** the
prescribed-shell behavior suite — that drives `record-commit-hash.sh` against a
pre-truncated archived REQ with `git cat-file -s` failing on PATH, and asserts it does not
report `all content guards passed`. The guards file is this script's dedicated behavioral
suite and already builds the throwaway git repo the fixture needs; REQ-276 hit exactly this
misdirection in its own `write_set` and recorded it as D-01.
**Why RED now:** Reproduced — see **Reproduced During Clarify** above. The unfixed script
prints `OK: ... all content guards passed`, writes the hash into an 85-byte remnant of a
10,178-byte file, and exits 0.
**GREEN when:** That case passes and the full suite still exits 0.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure
Ratchet (Step 6.5)**.

## Open Questions

- [x] REQ-268's review found that the unchecked-exit-status pattern it fixed in two
  timestamp scripts was copied from `record-commit-hash.sh` — the script whose guard style
  those two are required to imitate — so the template still carries the defect and any new
  guarded-edit script written from it inherits it. Fixing that means touching the script
  that guards the last write every REQ passes through, plus a sweep of the other shipped
  scripts and a new rule in the shell prime. Should I process this as a new task?
  → **Yes, add to queue — the FULL scope as written**, including the repo-wide sweep and the
  mechanical check. Chosen over the narrower "just the reproduced bug" option I recommended.

  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

  Why this is yours rather than mine: the fix is safe in itself, but the file it centres on
  is the one that once truncated six archived REQs to zero bytes when edited freely, and
  the sweep's real cost is the third requirement — deciding whether a mechanical check for
  this shape can be written without flagging the many places where discarding a status is
  correct. That is a judgment about how much guard machinery this repo wants, not a
  technical unknown.

**Answered [2026-08-20]:** User approved via `do-work clarify`, at full scope, after asking
how the fix would work and being shown the reproduction below and the four sites.

---

## AI Execution State (P-A-U Loop)

<!-- Added by the builder: capture minted this REQ without the section, which left
     tools/checks/qualify.sh Check 4's box audit DISARMED rather than passing. -->

- [x] **[PLAN]:** Reproduced the RED before writing anything; approach in `## Decisions`.
- [x] **[APPLY]:** Five files, all within the sweep's scope.
- [x] **[UNIFY]:** `git diff --stat` reviewed per file. `shellcheck` clean on all three edited scripts and both edited suites; `bash -n` clean. Diff grepped for debug artifacts — none. Every mutation reverted, confirmed by a clean suite run.

## Triage

**Route: C** - Complex

**Reasoning:** A reproduced fail-open in the repo's own incident guard, a sweep across every shipped script, a rule to state at its home, and a mechanical check whose feasibility the REQ itself disputes. The check's viability could not be decided without measuring it against the real corpus first.

**Planning:** Required

## Plan

1. Reproduce the RED against the real script before touching it — the REQ reproduced it during clarify and that claim is not inherited.
2. Fix `:466` by asking `git cat-file -e` first, so *no blob* and *git could not answer* stop being the same empty string.
3. Judge `:189`, `:194`, `:605` — fix where the shape is the same, comment where the content genuinely suffices. Silence is not an option the REQ allows.
4. Sweep every shipped script for the condition, not the instance list.
5. State the condition in `_dev/primes/prime-shell-commands.md`.
6. Attempt the mechanical check. **Measure its false-positive rate on the real corpus before committing to it**, and take the REQ's escape hatch honestly if the measurement says to.

*Planned inline by the orchestrator*

## The Sweep, and What It Measured

**Every shipped script under `skills/*/scripts/` and `skills/*/tools/`** was scanned for a command substitution whose exit status is discarded into a value indistinguishable from a legitimate result (`|| true`, `|| echo 0`, `|| echo ""`).

**15 sites, in 8 files. Every one of them is a correct use.**

| Site | Why the discard is right |
|---|---|
| `generate-report-image{,-batch}.sh` ×3, `run-blocked-check.sh` ×3 — `ps -o pgid=`/`-o stat=` | an empty answer means *the process is gone*, which is the information the caller wants |
| `do-work-update.sh:58`, `install-do-work-suite.sh:62` — `git rev-parse --show-toplevel` | empty means *this is not a repository*, a real and expected state |
| `preflight.sh:47` — `mktemp` | empty is checked on the very next line and reported as a WARN |
| `blanked-req-scan.sh:247,249` — internal resolvers | empty is normalized to `"-"` two lines later |
| `do-work-update.sh:111` — `read_action_version` | empty means *not installed* |
| `install-do-work-suite.sh:350` — `stat -c … \|\| stat -f … \|\| true` | the GNU/BSD fallback chain; the final `true` means *no stat worked* |
| `blanked-req-scan.sh:165,256` — `git cat-file -s` | display only — **improved anyway**, see D-03 |

**Beyond `record-commit-hash.sh`, the sweep found no second instance of the defect.** The primitive had been copied template → copies, and REQ-268 fixed the copies; this REQ fixes the template, and there is nowhere else it landed.

## Decisions

- **D-01** (ESCALATE): **The mechanical check is narrow, and that is the REQ's escape hatch taken with the evidence it asks for.** The broad check — any substitution with a collapsing fallback — was built and measured: **15 flagged, 0 defects.** A check with a 100% false-positive rate on the shipped corpus would be muted within a week, which is exactly the dissent the REQ records and overrules. What separates the defect from the correct use is semantic — *does a guard then make a safety decision on the value, where the collapsed value silently satisfies the safe branch* — and that is not greppable.
  So the check pins the one query where the distinction is unambiguous: **`git cat-file -s` may not discard its status into a number or an empty string.** Its answer is only ever a size guard; a display-only caller may fall back to `'?'`, which no reader can mistake for a size. Zero false positives on the corpus, and it catches the exact incident shape — verified by reinstating it.
  **This is a partial delivery of the requirement as written, reported as a finding rather than presented as complete.** The general condition is carried by the prime section instead, which is where the overruled recommendation wanted it — so the outcome is both: a real check on the provable case, and the stated convention for the rest.
- **D-02** (DECIDE & STATE): `:189` got the same structural fix as `:466` rather than a comment. It fails safe, but it fails safe with the *wrong message* — "the committed content is not what was verified" sends the reader hunting a pre-commit hook that was never involved. A guard that stops for the wrong stated reason costs the same debugging hour as one that does not stop.
- **D-03** (DECIDE & STATE): `:194` and `:605` got judgment comments, and `:605`'s existing comment **named the wrong command**: it said `|| true` absorbs grep's no-match, when under `pipefail` it is absorbing `diff`'s exit 1 — which is the *normal* case there, since the line exists to count how two files differ. The `|| true` is load-bearing for the success path and cannot be removed. Recording that mattered more than the fix would have.
- **D-04** (DECIDE & STATE): `blanked-req-scan.sh:165`'s `|| echo 0` became `|| echo '?'`, matching its sibling at `:256`. Display-only, so not a safety defect — but a `0` there reads as "the recoverable version is empty too", which is the opposite of the truth and would talk a reader out of a restore that would have worked.

## Implementation Summary

**What was done:** Reproduced the fail-open in the repo's own incident guard, fixed it by asking whether the blob exists before asking its size, judged the three fail-safe sites in place, swept every shipped script and measured the result, stated the condition where shipped shell is governed, and added the narrow mechanical check the corpus can actually support.

**Files changed:**
- `skills/do-work/tools/checks/record-commit-hash.sh` (modified) — the truncation floor now refuses when an existing blob will not size and skips only a genuine no-blob; `:189` given the same treatment; `:194` and `:605` carry judgment comments, one correcting which command's status is absorbed.
- `skills/do-work/tools/checks/blanked-req-scan.sh` (modified) — both `git cat-file -s` display sites fall back to `'?'` rather than `0`, with the reason in place.
- `_dev/primes/prime-shell-commands.md` (modified) — new § *Unchecked Exit Status Reads as Content*: the condition, the measurement showing why discarding is usually correct, and three ways out in preference order.
- `_dev/tests/contract-regressions.sh` (modified) — the narrow `git cat-file -s` check, carrying the measurement that explains its narrowness.
- `_dev/tests/record-commit-hash-guards.sh` (modified) — the REQ's captured RED as a permanent probe, plus a control and a no-blob case.

**Tests touched:** one behavioral probe (four assertions plus a control and a negative case) and one static check. No existing assertion changed meaning.

## Qualification

Passed — 5 files verified, 4 acceptance criteria traced, P-A-U confirmed.

- **[UNIFY] audit:** `shellcheck` clean on every edited file; `bash -n` clean. No debug artifacts. `maintainer-verify.sh` exits 0.
- **Substantive:** the fix changes real behaviour — proven by driving the real script against a real truncated file with a stubbed git.
- **Requirements traced:** "stated once where shipped shell is governed" → the prime section; "`record-commit-hash.sh`'s substitutions check their status or state why not" → all four judged, two fixed two commented; "a mechanical check, or a stated convention with the reason recorded" → D-01, both delivered with the measurement; "verify exits 0" → it does.
- **Flowing:** the guard actually probes git and changes its verdict — shown by the control/stub pair.

## Testing

- `bash _dev/tests/contract-regressions.sh` — passes, including the guards suite it invokes. `bash _dev/tests/maintainer-verify.sh` — exit 0. `shellcheck` clean on all five edited files.

**Red-green validation** — the REQ's captured RED, reproduced first on the untouched script:

| | Before the fix | After |
|---|---|---|
| 13,900-byte archived REQ truncated to 57 bytes, everything real except a failing `git cat-file -s` | `OK: … all content guards passed`, hash written into the 57-byte remnant, operator told to commit it | `FAIL: … exists in HEAD but its size could not be read`, **nothing written**, exit 1 |
| control, real git | `FAIL: … 57 bytes on disk but 13900 bytes in HEAD` | unchanged |
| a file with no blob in HEAD (new file) | skipped | **still skipped** — the fix must not turn "this file is new" into a failure |

All four are now permanent assertions in `record-commit-hash-guards.sh`, including the control (so a stubbed run that proves nothing cannot pass) and the no-blob case (so the fix cannot over-refuse).

**Both new checks were proven against the unfixed script**, by reverting the fix and re-running:

| Check | On the unfixed script |
|---|---|
| behavioral probe | **FAIL** ×2 — "the guard failed OPEN", "exited 0 with an unreadable HEAD size" |
| static `git cat-file -s` check | **FAIL**, naming the file and line |

The static check was also proven to catch a *display-site* regression: reverting `blanked-req-scan.sh:170` to `|| echo 0` fails it.

**The broad check's measurement** (D-01's evidence): built, run over every shipped script, **15 flagged and 0 defects** — the table in `## The Sweep` names each and why it is correct.

## Review

**Overall: 90%** — Acceptance: Pass

### Requirements Check

| Requirement | Status |
|---|---|
| The condition is stated once, where shipped shell is governed | ✅ new prime section |
| `record-commit-hash.sh`'s substitutions check their status or say why the content suffices | ✅ four sites, two fixed, two commented |
| A mechanical check — or a stated convention with the reason recorded | ⚠️ **both, and the split is a finding** — narrow check on the provable case, convention for the rest, with the measurement (D-01) |
| `maintainer-verify.sh` exits 0 | ✅ |

### Findings

**Important:**

- **I1 (the delivery gap, reported not hidden):** The requirement asked for a check that finds *the shape*. What shipped finds one query's instance of the shape. The measurement is the reason — 15 correct uses, 0 defects, on the broad form — and the REQ's own escape hatch sanctions exactly this outcome *provided it comes with the evidence*, which `## The Sweep` and D-01 carry. Flagging it as Important so the partial scope is a decision the user sees rather than a detail in a table.

**Minor:**

- **M1:** The narrow check is keyed on `git cat-file -s`, which is a spelling — the thing REQ-293 spent its whole scope removing elsewhere. It is defensible here because the *condition* now lives in the prime and the check is explicitly the provable subset, but a reader who has just read REQ-293 will notice the tension, and they are not wrong.
- **M2:** The sweep covered `skills/*/scripts/` and `skills/*/tools/` — shipped scripts. Shell embedded in action-file code blocks was not swept; `_dev/tests/action-shell-blocks.sh` governs those separately, and the REQ's scope line names scripts.

**Nit:**

- **N1:** The behavioral probe builds a 200-line fixture REQ to get a size ratio that trips the floor. A shorter file with a lower ratio would do, but 200 lines matches the incident's shape and costs milliseconds.

### Restatement Sweep

Redefined element: what a discarded exit status means, and where that rule lives.

- `record-commit-hash.sh`'s own header cites its guard style as the template other scripts copy — re-read: still accurate, and now the template no longer carries the defect it was propagating.
- `repair-req-timestamps.sh` and `audit-archive-timestamps.sh` — the two scripts REQ-268 fixed by copying the *correct* pattern. Grepped both: neither carries the collapsing shape, and neither restates the rule in a way the new prime section contradicts.
- `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale — the neighbouring rule, and the one this REQ's instance list explicitly defers to. The new section is keyed on the condition and marks its 15-site table as a measurement rather than a list to maintain.
- `_dev/tests/action-shell-blocks.sh` — checked whether it states an overlapping rule about exit status in action-file shell: it does not.

No stale restatement remains. One was *corrected* as part of the work: `:605`'s comment named grep where it meant diff (D-03).

### Acceptance Testing

The script was driven, not read: a real git repo, a real 13,900-byte archived REQ, a real truncation to 57 bytes, and a stub shadowing exactly one git subcommand. Before the fix it wrote the hash into the remnant; after it refuses and writes nothing. The control and the no-blob case bound the fix on both sides.

### Scores (on the record — not the headline)

| Dimension | Score |
|---|---|
| Requirements | 85% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope Discipline | 95% |
| Risk | Low |
| Acceptance | Pass |

Requirements 85% for I1 — three of four fully met, the fourth delivered in the split form its own escape hatch defines. Acceptance is Pass rather than Partial because the escape hatch is part of the requirement, not a deviation from it, and the evidence it demands is present.

### Follow-up REQs Created

None. I1 is a scope outcome the REQ pre-authorised, recorded here for the user's judgment rather than queued as work.

## Lessons Learned

**What worked:** Building the broad check *first*, running it, and counting. The REQ recorded a dissent that predicted false positives and a user decision to build it anyway with an evidence-based escape hatch — and the only honest way to use that hatch is to have built the thing and measured it. Fifteen flagged, zero defects is a number; "I think it would flag too much" is an opinion, and the REQ was explicit that only the former would do.

**What didn't:** The first attempt to run `record-commit-hash-guards.sh` directly failed with "must exist and be executable" against a file that plainly was. The suite tests the *canonical* `tools/checks/` path, and `contract-regressions.sh` invokes it through a `sed` that repoints it at the shipped copy. Worth knowing before debugging a phantom permissions problem: in this repo a `_dev/tests/*.sh` file is not necessarily runnable on its own terms, because nothing auto-discovers them and some are rewritten at their call site.

**Worth knowing:** The sharpest form of this defect is that **the fallback value looked like an answer**. `|| true` on a size query produces `""`, and `""` is what "there is no blob" also produces — so the guard could not tell a missing file from a broken tool, and the safe-looking branch was the wrong one for one of them. The general shape to watch: whenever a fallback value is in the same domain as a legitimate result, the failure has been laundered into data. `'?'` is safe precisely because it is not a number.

## Orientation

The script that exists to stop archived REQs being truncated could itself be talked into blessing a truncated file, by nothing more exotic than git failing to answer one question — and it is the template two other scripts were written from. It now refuses when it cannot read what it needs, and the condition behind it is stated where shell is written rather than re-derived per script header. Lives in the core check tooling (`skills/do-work/tools/checks/`) with its rule in `_dev/primes/prime-shell-commands.md`.

**[MAP CHANGED]** — `_dev/primes/prime-shell-commands.md` gains a section every author of shipped shell now loads before writing any, and the contract suite gains a check that fails a build on one specific misuse. The prime is the durable half: the check covers the provable subset by design, and the section says so.

Prime staleness spot-check: `_dev/primes/prime-shell-commands.md` — edited here; its referenced paths still resolve and its Closed Enumerations rule is honored by the new section rather than contradicted.
