---
id: REQ-274
title: Retire the "the SessionStart hook exits nonzero" framing where it is still stated
status: completed
created_at: 2026-08-18T23:38:35Z
status_changed_at: 2026-08-19T13:45:20Z
claimed_at: 2026-08-20T23:40:01Z
completed_at: 2026-08-20T23:45:29Z
kb_status: pending
commit: 0efefa6
user_request: UR-056
addendum_to: REQ-267
domain: general
route: B
review_generated: true
sweep: true
sweep_key: repairer-hook-failure-framing-restated
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec: bug-fix
depends_on: []
maintenance: true
write_set:
- _dev/primes/prime-shell-commands.md
---

# Retire the "the SessionStart Hook Exits Nonzero" Framing Where It Is Still Stated

## What

A timestamp-repairer failure **does not fail the SessionStart hook**. `skills/do-work/hooks/session-start.sh:59` runs the script under `|| true`, deliberately, with a comment saying that on a tripped guard the script's failure lines *are* the audit trail and must reach the banner. `report_failure` writes to stdout, the hook captures it into `REPAIR_SUMMARY` and echoes it, and the hook exits **0** — verified by running the real hook against a wedge fixture.

The real consequence is still bad: the script exits 1 on every run and prints a `FAILED to repair …` line into every session's start banner, permanently, with no self-heal. But three live maintainer docs state the *mechanism* as a nonzero hook exit, and one of them is a standing decision rationale:

- `do-work/RESTART-PROMPT.md:39` — "can wedge the SessionStart hook into exiting nonzero every session"
- `do-work/CHECKPOINT.md:76` — REQ-255 D-04's rationale: "refusal would have made the SessionStart hook exit nonzero every session". **The decision is still right; the argument for it names a mechanism that does not exist.**
- `_dev/primes/prime-shell-commands.md:34` — "the fuzz found the shape that wedges the hook"

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [x] `do-work/RESTART-PROMPT.md` — **already corrected.** Both files were rewritten wholesale at the end of the session that
  created this REQ, and the orchestrator wrote the true framing rather than reproducing a claim it knew to be false. Verify
  rather than assume: re-grep before deciding this instance is closed.
- [x] `do-work/CHECKPOINT.md` — **already corrected**, same rewrite. REQ-255 D-04's rationale no longer appears there at all,
  because the checkpoint was replaced by the next session's; if that decision's reasoning is still wanted in a durable place,
  its archived REQ is the home, and archived REQs are immutable.
- [x] `_dev/primes/prime-shell-commands.md:34` — **still live.** This is the one that matters: a prime is loaded and acted on. — **fixed** (the line had drifted to :40). REQ-255's lesson link now reads "wedges the repairer into printing an unhealable failure line into every session's banner — the hook itself exits 0 either way". The lesson's substance is unchanged; only the mechanism it named was false.

## Note added after the instance list was written

Two of the three instances sat in `do-work/` state files that the session rewrites at its end. Rather than deliberately re-writing a known-false sentence to preserve a tidy sweep, the orchestrator wrote the correct framing into both. **This REQ is therefore mostly about the prime**, plus the sweep that checks nothing else restates the old mechanism. Do not read the two ticked boxes as work done by this REQ's builder — read them as instances that expired, and re-grep to confirm.

## Requirements

- The rule is **stated once**: a repairer failure prints into the session banner; it does not fail the hook. Nothing restates it on the old mechanism.
- Where the false mechanism carried a decision's argument (REQ-255 D-04), the decision keeps standing on the true consequence — a permanent banner failure line with no self-heal supports it just as well. **Do not silently reverse a decision while correcting its reason.**
- Sweep the primitive rather than fixing three lines: grep for every statement about what a repairer failure does to the hook, in any spelling, and check each against `session-start.sh`.
- Archived REQs under `do-work/archive/` are immutable record and stay untouched.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Context

REQ-267's independent review, Important finding 1 (gate: rule-change). Created `pending-answers` per the generation-≥2 cascade stop, since REQ-267 is itself `review_generated: true`.

This is REQ-257's lesson arriving one REQ later, and it is worth stating plainly: **a claim inherited rather than re-derived can be right in its conclusion and false in its mechanism, and the false mechanism travels.** This one travelled far enough to set REQ-267's own approval framing — the orchestrator repeated it to the maintainer twice before a builder checked the hook.

## Open Questions

- [x] REQ-267's review found that three live maintainer docs state a repairer failure as making the SessionStart hook exit nonzero, which it does not — including a standing decision rationale that rests on it. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the affected decisions are all still correct, so accept the wrong mechanism in the record.

---

## Triage

**Route: B** - Medium

**Reasoning:** The correction itself is one line, but the REQ's own Requirements make the sweep the deliverable — find every statement of what a repairer failure does to the hook, in any spelling, and check each against `session-start.sh`. That is discovery work, and the instance list is explicitly labelled a capture-time hypothesis to re-verify rather than a work order.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**The claim was re-derived, not inherited.** That is this REQ's own lesson, so accepting its premise on the strength of the REQ text would have been the defect committing itself. Two probes against the real hook:

1. The refusal shape (`created_at: 9999-99-99T99:99:99Z`) exits **0** and prints nothing — refusals are counted but voicing is opt-in through `timestamp_repair_voice_refusals`, exactly as `scripts/repair-req-timestamps.sh:108-117` documents. This shape does not reproduce the wedge at all.
2. A stub repairer that prints `do-work: FAILED to repair …` and exits **1**, wrapped by the **real** `hooks/session-start.sh`: the hook exits **0** and the FAILED line reaches the banner.

So the REQ is right on both halves: the hook does not exit nonzero, and the consequence is a permanent unhealable banner line. The mechanism is structural — line 59's `|| true` — so no repairer exit status can change it.

**Sweep of the primitive** (`wedg`, plus `hook` within 60 characters of `exit|fail|nonzero`, across every tracked `.md` and `.sh`, excluding `do-work/archive/` as immutable record):

| Site | Verdict |
|---|---|
| `_dev/primes/prime-shell-commands.md:40` (REQ-255 lesson link) | **False mechanism, live, fixed.** The line had drifted from the :34 the REQ recorded. |
| `skills/do-work/scripts/repair-req-timestamps.sh:108-117` | **Already true.** Says "print the same unhealable line into every session's start banner forever" and never claims a nonzero hook exit. This is the rule's canonical home. |
| `_dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh:197-200` | **Already true.** "a nonzero exit here would print a FAILED line into every session that no one can ever heal" — names the script's exit and the banner consequence, correctly. |
| `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh:299` | **Not this class.** "wedge the repair permanently" is about the repair, not the hook. |
| `skills/do-work/hooks/session-start.sh:54-56` | **Already true.** States why `|| true` is used instead of discarding. |
| `do-work/RESTART-PROMPT.md` | **Instance expired — verified by re-grep, not assumed.** Zero matches for hook/wedge/nonzero. |
| `do-work/CHECKPOINT.md` | **Instance expired — verified by re-grep.** REQ-255 D-04's rationale is absent; the only `hook` match is the unrelated filename `install-memory-hooks`. |
| `do-work/runs/work-2026-08-18-230100/manifest.md:16` and `REQ-267-brief.md:80` | **False mechanism, deliberately not edited** — see Decisions D-01. |

*Exploration run inline by the orchestrator*

## Scope

**Files I will touch:**
- `_dev/primes/prime-shell-commands.md` (modify) — correct REQ-255's lesson link to name the true consequence

**Files I will NOT touch:**
- `do-work/RESTART-PROMPT.md`, `do-work/CHECKPOINT.md` — instances expired; re-grep confirms both are clean.
- `skills/do-work/scripts/repair-req-timestamps.sh`, `_dev/tests/prescribed-shell-cases/*`, `skills/do-work/hooks/session-start.sh` — already state the rule truly. Editing them to "fix" a sweep would be the additive reflex `crew-members/maintenance.md` warns against.
- `do-work/runs/**` — historical run record (D-01).
- `do-work/archive/**` — immutable by the REQ's own Requirements.

**Acceptance criteria (restated from the REQ):**
1. The rule is stated once and nothing restates it on the old mechanism.
2. REQ-255 D-04's decision keeps standing on the true consequence; the decision is not reversed while its reason is corrected.
3. The primitive was swept in any spelling, each hit checked against `session-start.sh`.
4. Archived REQs untouched.
5. `bash _dev/tests/maintainer-verify.sh` exits 0.

## Decisions

- **D-01** (ESCALATE): Two run artifacts under `do-work/runs/work-2026-08-18-230100/` still carry the false mechanism verbatim ("can wedge the SessionStart hook into failing every session"). Builder chose **not** to edit them. Reasoning: a run directory is the record of what a past dispatch actually decided and why, in the same class as an archived REQ, which this REQ's own Requirements protect as immutable. Rewriting it would erase the evidence of how far the false claim travelled — which is the REQ's stated point. They are also not loaded by anything: no action, prime, or crew file reads `do-work/runs/**` as guidance. Value: the record stays honest and the sweep's finding is reported rather than buried. Risk: a future reader browsing that run directory reads the false mechanism. Mitigated because the correction now lives in the prime that *is* loaded, and because the run artifact names REQ-267 and REQ-274, which lead to this correction.
- **D-02** (DECIDE & STATE): Three swept sites already state the rule truly and were left alone. Reasoning: `crew-members/maintenance.md` § 1 — the first move is subtraction, and here the correct move was smaller still: nothing. A sweep's output is a verdict per site, not an edit per site.
- **D-03** (DECIDE & STATE): The REQ's `write_set` listed `do-work/RESTART-PROMPT.md` and `do-work/CHECKPOINT.md`; both were dropped when re-grep confirmed the instances had expired. `## Scope` is the source and `write_set` its mirror, so the field was rewritten to the single file actually touched.

## Implementation Summary

**What was done:** Swept every live statement of what a timestamp-repairer failure does to the SessionStart hook and found exactly one still on the false mechanism — REQ-255's lesson link in `_dev/primes/prime-shell-commands.md`, the one site the REQ predicted would matter because a prime is loaded and acted on. Corrected it to name the true consequence. The other seven swept sites were verified individually: three already stated the rule truly, two instances had expired, one was a different claim, and two run artifacts were deliberately left as historical record.

**Files changed:**
- `_dev/primes/prime-shell-commands.md` (modified) — REQ-255's lesson link now reads "wedges the repairer into printing an unhealable failure line into every session's banner — the hook itself exits 0 either way".

**Tests touched:** none. No behavior changed; `session-start.sh` and the repairer are untouched, which is the point — the code was always right and only the docs were wrong.

## Qualification

Passed — 1 file verified, 5 acceptance criteria traced, P-A-U confirmed.

- **[UNIFY] audit:** `git diff --stat` shows one project file changed (`_dev/primes/prime-shell-commands.md`, 1 insertion / 1 deletion). No debug artifacts — the diff is a single prose clause. No linter applies to Markdown; `maintainer-verify.sh` exits 0 and covers this file's citation contract (the lesson link's target path is unchanged and still resolves).
- **Substantive:** the changed clause replaces a false mechanism with the verified one.
- **Requirements traced:** AC1 → the sweep table shows one canonical statement and no surviving restatement; AC2 → D-04's decision is untouched and its true consequence supports it (see below); AC3 → the sweep table; AC4 → `do-work/archive/` excluded from every grep and edit; AC5 → verify exits 0.
- **Flowing:** not applicable — no data path.

**On AC2 specifically.** REQ-255 D-04 chose refusal over repair for shapes the recognizer cannot safely derive. Its recorded argument was "refusal would have made the SessionStart hook exit nonzero every session". The true consequence — refusal prints an unhealable `FAILED`/refusal line into every session banner forever, with no self-heal — supports the same decision at the same strength: both are permanent, both are unattended, both are visible to the user every session. The decision is not reversed and was not re-opened. Its archived REQ is untouched, per AC4; the correction lives in the live prime that restates it.

## Testing

- `bash _dev/tests/maintainer-verify.sh` — exit 0. Baseline was green (`launched: true`, `exit_status: 0`), so this is a clean no-regression comparison.
- **Behavioral re-derivation** (the evidence the correction rests on, not a regression test):

| Probe | Repairer exit | Hook exit | Banner |
|---|---|---|---|
| refusal shape (`9999-99-99T99:99:99Z`) | 0 | 0 | silent (voicing is opt-in) |
| stub repairer printing `FAILED to repair …` | 1 | **0** | FAILED line present |

Red-green validation is omitted: this is a documentation correction with no behavioral change. The table above is the evidence that the *new* wording is true, which is the analogous proof for a prose fix.

## Review

**Overall: 94%** — Acceptance: Pass

### Requirements Check

| Requirement | Status |
|---|---|
| Rule stated once; nothing restates the old mechanism | ✅ One live site fixed; canonical home already correct |
| REQ-255 D-04 stands on the true consequence, not silently reversed | ✅ Argued explicitly in Qualification |
| Sweep the primitive in any spelling, check each against `session-start.sh` | ✅ Eight sites, verdict each |
| Archived REQs untouched | ✅ Excluded from every grep and edit |
| `maintainer-verify.sh` exits 0 | ✅ |

### Findings

**Important — none.**

**Minor:**

- **M1:** The REQ recorded the live instance at `prime-shell-commands.md:34`; it was at :40. Line numbers in a capture-time instance list go stale as the file grows, which is why the REQ correctly told the builder to re-grep. No action — the REQ already handled this. Noted because it is the second REQ this session where a capture-time list was narrower or staler than reality.
- **M2:** Two `do-work/runs/**` artifacts still carry the false mechanism. Reported, not fixed, with the reason in D-01. A reader who disagrees with treating run directories as immutable would call this an unfixed instance.

**Nit:**

- **N1:** The corrected lesson line is now long enough to wrap several times in a terminal. Prime lesson links are already long by convention here, so it is in keeping, but the clause "— the hook itself exits 0 either way" is the part carrying the correction and could stand alone if this line is ever split.

### Restatement Sweep

Redefined element: what a repairer failure does to the SessionStart hook. The sweep for this element **is** the REQ's deliverable, so it ran as the implementation rather than as a separate review pass — the full eight-site table is in `## Exploration`. Re-run after the edit to confirm no new restatement was introduced: the corrected line is the only occurrence of the true framing added, and it names the mechanism in the same words as the canonical home in `scripts/repair-req-timestamps.sh:108-117`.

### Acceptance Testing

The deliverable is a corrected sentence, so acceptance is: does the new sentence match observed behavior? The two probes in `## Testing` were run against the real hook and confirm both halves of the new wording — the repairer's failure line reaches the banner, and the hook exits 0. Nothing was accepted on the REQ's authority.

### Scores (on the record — not the headline)

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A (documentation correction; behavioral probes are the proof) |
| Scope Discipline | 100% |
| Risk | None |
| Acceptance | Pass |

### Follow-up REQs Created

None. M2 is a recorded decision, not an open task; M1 and N1 are observations.

*Reviewed in orchestrated mode by the work orchestrator*

## Lessons Learned

**What worked:** Refusing to inherit the REQ's premise. This REQ exists because a claim travelled three documents deep without anyone re-deriving it, so accepting "the hook exits 0, verified against a wedge fixture" on the REQ's word would have reproduced the exact failure being fixed. Re-deriving it also produced something the REQ did not have: the refusal shape it names does *not* reproduce the wedge (voicing is opt-in), so the proof needed a stub repairer that actually fails. The REQ's conclusion was right and its named fixture would not have demonstrated it.

**What didn't:** Two attempts at reproducing a real repairer failure failed before the stub worked — the refusal path is silent by design, and an unwritable-directory trip does nothing when the session runs as root. Worth knowing for any future probe of this script: root defeats every permission-based failure path in it, so simulate the failure at the seam instead of trying to provoke it.

**Worth knowing:** The rule's canonical home is the comment block at `skills/do-work/scripts/repair-req-timestamps.sh:108-117`, and it has been correct all along. Every false restatement was downstream of it in a *lesson link* or a *run artifact* — narrative surfaces, not contract surfaces. That is the pattern worth watching: the contract text stayed true while the stories told about it drifted, and the stories are what the next reader reaches first.

## Orientation

The prime that carries this repo's shell-command lessons no longer tells its reader that a timestamp-repairer failure kills the SessionStart hook. It does not — the hook wraps the repairer in `|| true` on purpose so the failure lines reach the session banner. Lives in `_dev/primes/prime-shell-commands.md`, indexing the shell/hook subsystem whose contract home is `skills/do-work/scripts/repair-req-timestamps.sh`. Leaf change: no code, no contract, no renamed concept.

Prime staleness spot-check: `_dev/primes/prime-shell-commands.md` — the edited lesson link's target still resolves (green citation contract), and its other referenced paths are unaffected.
