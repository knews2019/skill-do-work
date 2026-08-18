---
id: REQ-225
title: State verified-exact-publication once as a condition in the shipped shell guide
status: completed
status_changed_at: 2026-08-17T21:09:46Z
completed_at: 2026-08-18T00:20:49Z
claimed_at: 2026-08-18T00:07:07Z
domain: general
created_at: 2026-08-17T21:02:00Z
user_request: UR-042
addendum_to: REQ-220
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
maintenance: true
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-08-18T00:11:00Z
  basis:
    - Route B
    - 1-file write set
    - 4 acceptance criteria
    - cross-route regression gates
write_set:
- skills/do-work/docs/prescribed-shell-primitives.md
- _dev/tests/prescribed-shell-canonicalization.sh
---

# Discovered Task: State Verified-Exact-Publication Once as a Condition

## What

`skills/do-work/docs/prescribed-shell-primitives.md` is the canonical shipped guide for shell used across do-work actions. It states the "a rename onto an occupied destination nests instead of colliding" rule **only** inside its `## Portfolio summary publication` section, phrased as a property of one script. Restate it once as a condition that applies to any publication, and have the per-script sections point at that statement instead of each carrying (or omitting) their own copy.

## Context

Found while implementing REQ-220. The same defect has now been fixed in four separate places — `publish-portfolio-summary.sh` (REQ-199/205), the `ai-report` prescribed batch block (REQ-204), and `generate-report-image.sh` plus `install-last30days.sh` (REQ-220). Each fix was local, and each of the last three was found by a review sweep rather than by reading the guide, because the guide never says the rule in a form that would make a reader check their own publication against it.

`CLAUDE.md` § *State conditions, not lists* names this exact failure: when a rule applies "whenever X happens", key it on the condition, because a hand-maintained per-script list goes stale as the set grows. `_dev/primes/prime-shell-commands.md` § *Closed Enumerations Go Stale* records the same lesson from four independent defects.

## Requirements

- State the condition once — something to the effect of *any publication whose destination could be occupied verifies the path it actually wrote, rather than reading the rename's exit status as proof* — in a location that applies to every shipped publication helper, not inside one script's section.
- Have the existing per-script sections reference that statement rather than restating it. The `## Portfolio summary publication` section's script-specific consequences (snapshot candidates advance by numeric suffix, the canonical path fails closed) stay where they are — those are policy, not the shared primitive.
- Do not change any script behavior. This REQ is documentation only; all four scripts already implement the rule.
- Keep the shipped reference contract green (`_dev/tests/shipped-package-reference-contract.sh`) — this is a shipped file, so it may not cite `_dev/` paths.

## Open Questions

- [x] Should the shared shell guide state the publication-nesting rule once as a general condition, instead of describing it inside one script's section? → Confirmed: Yes, add to queue
  *[2026-08-17] User confirmed via `do-work clarify`. Consent given for this cascade-depth-two follow-up to run another autonomous cycle. Scope stays documentation-only: no script behavior changes, and the portfolio-summary section's script-specific policy stays where it is.*
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — leave the rule stated per script.
  Why this is yours: nothing is broken. All four affected scripts already verify their publications, so this changes no behavior and fixes no defect — it is a judgment call about how the shipped guide is organized, and reorganizing a canonical shipped document is a taste decision rather than a repair. It is also worth knowing that REQ-220 was itself a review follow-up two generations deep, so the cascade-depth rule requires your consent before another autonomous cycle. The argument for doing it: this defect class has been found four times by review sweep and zero times by someone reading the guide, which is the concrete cost of the rule living inside one script's section. The argument against: the guide is deliberately organized by executable home, and a cross-cutting section cuts against that organizing principle.

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is stated exactly (one shared condition, per-script sections point at it) and the write set names the file, but *which* sections currently restate or omit the rule — and which shipped helpers the condition would actually cover — has to be discovered before the wording can be right.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**The guide** (`skills/do-work/docs/prescribed-shell-primitives.md`, 104 lines, 11 `##` sections). Three sections describe a publication:

- `## Atomic download publication` — describes `scripts/atomic-download.sh`. **Omits** the rule entirely.
- `## Portfolio summary publication` — states the rule inline, mixed into one sentence with `publish-portfolio-summary.sh`'s own policy (snapshot candidate advances by numeric suffix, canonical fails closed).
- `## Report image batch publication` — states the rule inline a second time, mixed with `generate-report-image-batch.sh`'s policy (discard only its own stage, exit nonzero).

Two more helpers implement the rule with no section of their own — `generate-report-image.sh` and `install-last30days.sh` — and appear only as rows in `## Shipped executable homes`. So the rule is stated twice, implemented five times, and omitted from the one section that covers two helpers.

**The verification each script actually performs.** Confirmed by reading the source, not the table: `publish-portfolio-summary.sh` (snapshot and canonical, lines 102-163), `generate-report-image.sh` (lines 112-117), `generate-report-image-batch.sh` (lines 169-173), and `install-last30days.sh` (lines 98-103) each probe a `nested_*` path after the rename and fail or advance. All four match the REQ's premise.

**Two shipped publication helpers do not.** `scripts/atomic-download.sh:44` publishes with a bare `mv` and returns; `scripts/capture-screenshot.sh:33` publishes with a bare `ln` and returns. Neither checks the path it wrote. Reproduced — see `## Discovered Tasks`.

**Constraints that bound the wording.**
- `_dev/tests/prescribed-shell-canonicalization.sh` ratchets ten `##` headings in this guide by exact match; none may be renamed or removed.
- The same test fails any *other* shipped markdown that restates canonical rationale. New prose inside the canonical guide is exempt by construction — the guide is skipped by that loop.
- `_dev/tests/shipped-package-reference-contract.sh` forbids shipped files from citing `_dev/` paths. Shipped docs also use REQ ids only as illustrative examples, never as citations of this repo's own history, so no REQ id may appear in the new text.
- `_dev/tests/action-shell-blocks.sh` shellchecks fenced blocks; this change adds no fenced block.

**Where the condition goes.** Immediately before `## Atomic download publication`, so it sits ahead of all three publication sections and reads as governing them. Placing it inside any one of them would repeat the defect the REQ describes.

*Explored by work action (inline, serial mode)*

## Scope

**Files I will touch:**
- `skills/do-work/docs/prescribed-shell-primitives.md` (modify) — add the shared condition as its own section; replace the two inline restatements with pointers; give the atomic-download section the pointer it never had.
- `_dev/tests/prescribed-shell-canonicalization.sh` (modify) — add the new heading to the existing exact-match heading ratchet (D-01).

**Files I will NOT touch:** every shipped publication helper (`scripts/atomic-download.sh`, `scripts/capture-screenshot.sh`, `../do-work-toolbox/scripts/*`) — the REQ is documentation-only and the maintainer's recorded consent fixes that boundary. The two unverified helpers go to `## Discovered Tasks`, not into this diff.

**Acceptance criteria (restated from REQ):**
- [ ] The condition is stated once, in a location that applies to every shipped publication helper rather than inside one script's section.
- [ ] The existing per-script sections reference that statement instead of restating it.
- [ ] `## Portfolio summary publication` keeps its script-specific consequences (snapshot candidates advance by numeric suffix, the canonical path fails closed).
- [ ] No script behavior changes.
- [ ] `_dev/tests/shipped-package-reference-contract.sh` stays green.

## Implementation Summary

**Files changed:**
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)
- `_dev/tests/prescribed-shell-canonicalization.sh` (modified)

**What was done:** Added a `## Verified exact publication` section to the shipped guide, placed ahead of all three publication sections, stating once that a publication whose destination could be occupied verifies the path it actually wrote and that the rename's or link's exit status is not that proof. Replaced the two inline restatements — in `## Portfolio summary publication` and `## Report image batch publication` — with in-file anchor references to that section, keeping each helper's own policy sentence intact. Gave `## Atomic download publication`, which had no copy of the rule at all, a reference to the same section. Added the new heading to the exact-match heading ratchet in `_dev/tests/prescribed-shell-canonicalization.sh`. No shipped script was modified.

## Decisions

- **D-01**: Extended the declared write set from one file to two, adding `_dev/tests/prescribed-shell-canonicalization.sh`. Reasoning: this REQ moves rationale out of per-script sections into one shared section and leaves three in-file anchor links pointing at it, which is precisely the shape the existing ten-heading exact-match ratchet in that test already protects for the other ten sections. Proven against a replay case rather than assumed — with the heading absent from the ratchet, deleting the whole shared section and stranding all three anchors exits 0 and the suite notices nothing; with the heading added, the same deletion exits 1 naming the missing heading (`maintenance.md` § 3). One line added to a mechanism that already exists, not a new mechanism. DECIDE & STATE — reversible, and the extension is recorded here and mirrored into `write_set` as Step 6 requires.
- **D-02**: `## Atomic download publication` states that neither `atomic-download.sh` nor the screenshot install sharing its mechanics makes the verification check yet, instead of asserting that it does. Reasoning: the requirement is that the per-script sections reference the shared condition, and the only two readings available for this section were an accurate one and a false one — both helpers publish with a bare `mv`/`ln` and return. ESCALATE. **Value:** the guide is the contract agents read before invoking these helpers, and a section that silently implies compliance is how this defect class stayed unfound through four fixes; the sentence deletes itself when the follow-up lands. **Risk:** a shipped document now names a defect in shipped code, which reads oddly to anyone who reads the guide before the queue. Fully reversible — one sentence, removed when the follow-up REQ closes.

## Discovered Tasks

- **[critical]** `skills/do-work/scripts/atomic-download.sh` and `skills/do-work/scripts/capture-screenshot.sh` publish without verifying the path they actually wrote — the same defect class already fixed in four other helpers, found here for the fifth and sixth time. Reproduced, not inferred:
  - `atomic-download.sh:44` runs `mv "$download_path" "$target_path"`. With the target path occupied by a directory, the download nests as `<target>/<target-basename>.download.XXXXXX`, the script clears `download_path` so the cleanup trap spares it, and it exits 0. Reproduced with a `file://` source: exit 0, target still a directory, private file abandoned inside it.
  - `capture-screenshot.sh:33` runs `ln "$copy_path" "$destination_path"`. With the destination occupied by a directory, `ln` nests instead of refusing, so the documented no-clobber install silently does not happen and exit is 0. Under `--staged` this compounds into data loss: the success path then runs `rm "$source_path"` and destroys the staged screenshot — the only copy the capture dispatch holds — while the destination never receives it. Reproduced: exit 0, staged source deleted, destination still a directory holding an orphaned `.copying.XXXXXX` file.

  Out of scope here by the maintainer's recorded consent on this REQ ("no script behavior changes"). The fix pattern is already written four times in this repo (`publish-portfolio-summary.sh:102-163`, `generate-report-image.sh:112-117`, `generate-report-image-batch.sh:169-173`, `install-last30days.sh:98-103`): probe the nested path after the rename, remove only your own nested artifact, and fail closed. `capture-screenshot.sh` must additionally not delete the staged source on a failed publication.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` (canonical baseline: ShellCheck 0.11.0 over 50 tracked shell files, the full aggregate contract suite, queue-kanban vet/tests/strict-JS lane, audit-metrics vet/tests)
**Result:** ✓ All passing — exit 0, same as the pre-flight baseline recorded before implementation

**Red-green validation:**
- `_dev/tests/prescribed-shell-canonicalization.sh` heading ratchet (D-01): ✗ before → ✓ after. RED was established by mutation rather than assumed: with the section present but the heading absent from the ratchet, deleting the entire `## Verified exact publication` section and stranding all three in-file anchors exits **0** and the suite reports nothing. With the heading added, the same deletion exits **1** with `FAIL: prescribed-shell guide is missing heading: ## Verified exact publication`. The guide was restored from a byte copy afterwards and the suite re-run green.
- `_dev/tests/shipped-package-reference-contract.sh` (REQ acceptance criterion): ✓ PASS, run standalone as well as inside the aggregate suite.

**New tests added:** none — one entry added to an existing exact-match ratchet.

**Existing tests updated (cross-REQ impact):** `_dev/tests/prescribed-shell-canonicalization.sh` — one heading added to the required-heading list it has enforced since the canonicalization campaign. No existing assertion was changed, weakened, or removed.

**Environment note:** this checkout had no `shellcheck`, no `just`, and Go 1.24.7 against the required exactly-`go1.26.1`, so `maintainer-verify.sh` could not run at all on arrival. ShellCheck 0.11.0, just 1.43.0, and Go 1.26.1 were installed to a session-local directory and prepended to `PATH`; nothing in the repository was changed to accommodate the gap, and the baseline was confirmed green before any implementation edit.

*Verified by work action*

## Lessons Learned

**What worked:** Writing the condition down did the job it was captured to do, immediately and on its own author. Reading the guide's three publication sections against a single stated condition — rather than script by script — surfaced two more helpers that never verify their publication, one of them with a data-loss path. Four earlier fixes each took a review sweep to find; this one took a paragraph. Mutation-testing the ratchet before adding it was also worth the two extra runs: it turned "this looks like it should be pinned" into a demonstrated 0-then-1.

**What didn't:** The first instinct was to write `## Atomic download publication` as pointing at the shared condition the way the other two sections do, which reads well and would have been false — that helper publishes with a bare `mv` and returns. Asserting compliance because the sentence is symmetric is the same shape of error as trusting a rename's exit status. The section states the gap instead (D-02).

**Worth knowing:** `ln` refusing an occupied destination is a guarantee about *files* only; on an occupied directory it nests and exits zero, exactly like `mv`. Any helper whose no-clobber promise rests on bare `ln` has half the guarantee it appears to have. The compounding case is worse than the nesting itself: `capture-screenshot.sh --staged` reads that zero exit as permission to delete the staged source, so a publication that never happened destroys the only copy of the payload.

## Orientation

Anyone writing or reviewing a shipped publication helper now has one place that says what a publication owes its caller — `skills/do-work/docs/prescribed-shell-primitives.md` § **Verified exact publication** — instead of inferring it from whichever script's section they happened to open. The rule is keyed on the condition (a destination that could be occupied) rather than on a list of helpers, so a helper added later is covered without anyone remembering to extend anything. **[MAP CHANGED]** — the shipped shell guide was organized strictly by executable home, and this is its first section that is a shared primitive rather than one script's contract; per-script sections now reference it instead of restating it. Staleness spot-check on the REQ's prime `_dev/primes/prime-shell-commands.md`: every referenced path still resolves, and its § *Closed Enumerations Go Stale* now has a fifth instance rather than a contradiction — the prime is not stale.

## Review

**Overall: 95%** | 2026-08-18T00:20:04Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- `skills/do-work-toolbox/actions/present-work.md:137` restates the canonical container-not-a-collision rationale in full, which the mandatory restatement sweep surfaced as stale the moment this REQ gave that rationale a single home — gate: rule-change → rerouted pending-answers as REQ-230 (new sweep, `sweep_key: caller-doc-restates-canonical-publication-rationale`; REQ-225 carries `review_generated: true`, so the generation-≥2 depth stop applies)
- `atomic-download.sh` and `capture-screenshot.sh` publish without verifying the path they wrote, the second with a data-loss path under `--staged` — gate: user-visible → REQ-229 created `status: pending` (critical pierce). Provenance: raised by the builder in `## Discovered Tasks` and confirmed independently here by re-running both reproductions, not a second discovery.

**Minor findings:** 1 (report only) — the change touches `_dev/tests/prescribed-shell-canonicalization.sh`, a `.sh` file, while the REQ says "do not change any script behavior." The requirement's own sentence scopes that to the publication helpers ("all four scripts already implement the rule"), and no shipped script was touched; recorded here so the reading is auditable rather than assumed.

**Acceptance:** Pass — all five restated acceptance criteria verified against the file, and the canonical baseline exits 0.
**Suggested testing:** 0 items
**Follow-ups created:** REQ-229, REQ-230; **sweeps appended to:** None

*Reviewed by review-work action*
