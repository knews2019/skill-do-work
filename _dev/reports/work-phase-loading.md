# Work phase loading validation

Baseline: `afbc8be3` (0.315.2). Release: 0.315.3.

The finding is accepted: the former work action required all 603 lines before any
phase, including planning on Route A. The smallest change here extracts contiguous
procedures into adjacent files; schemas, canonical CLI ownership, review-work and
reference sections remain unchanged. No new action or runtime mechanism is added.

## Measurement method

Run `python3 _dev/tests/instruction-reading-volume.py` from this repository. It reads
the baseline through `git show` and the candidate from disk, counting LF bytes for
lines and actual UTF-8 bytes for volume. No token conversion is assumed. Each file
is counted once in one context; this is a deterministic reading walkthrough, not
a measurement of model attention or repeated reads by independent agents.

The measured boundary includes the skill router, applicable phase procedures,
`review-work.md`, the cited sections of `work-reference.md`, and the crew rules
loaded by these walkthroughs. The script prints every exact file/section and its
UTF-8 byte count for both baseline and candidate, alongside distinct file counts.
A reference file counts as one file even when several sections are read.

Reference selection matches exact heading names (or their step-suffix-free names)
in the loaded procedures and crew rules, then follows nested named references to
a fixed point. Headings inside fenced templates are ignored. Parent sections
subsume children and repeated reads count once. This conservative static model
includes named conditional recovery/follow-up sections in loaded procedures;
it does not claim those branches executed in a successful application change.
A pointer to a separate work phase alone does not load that phase.

Builds count `general`, `coding-guardrails`, `shared-principles`,
`communication-style`, `prompt-injection` and `anti-slop` crews. Review-only counts
the same set except implementation-only `general`: shared principles and coding
checks apply to review, communication is the installed session contract,
prompt-injection guards UR ingestion, and anti-slop applies to the review report.
These are general-domain walkthroughs without TDD, security, supply-chain or other
conditional domain crews. Model agents' repeated reads are not multiplied.

Project primes, lesson satellites, task artifacts, source code and references
outside `work-reference.md` remain outside the measured boundary. These are
expanded instruction-surface totals, not total task context. Byte volume spots
unnecessary reading in unrelated phases; it cannot establish instruction quality
or justify deleting required rules. No latency, cost or model-quality inference
is made.

| Walkthrough | Before files | After files | Before lines | After lines | Before bytes | After bytes |
|---|---:|---:|---:|---:|---:|---:|
| Simple fix | 10 | 15 | 2394 | 2306 | 301147 | 286689 |
| Complex change | 10 | 18 | 2394 | 2445 | 301147 | 301715 |
| Review only | 8 | 8 | 1106 | 1106 | 103725 | 103717 |

Work entry point: 603 lines / 94,224 bytes → 218 lines / 23,645 bytes.
The current candidate includes the prime-byte-budget crew change and the earlier
communication edit as well as the phase extraction. The complex walkthrough grows
by 568 bytes; review-only falls by 8 bytes. The simple walkthrough saves 14,458
bytes under this expanded, symmetric boundary. These are actual revision-to-disk
comparisons, not an attribution of every byte to the phase extraction. There is no measured time, cost or quality claim.

## Walkthroughs and check coverage

Paths below are relative to `skills/do-work/actions/`. Follow named sections in
`work-reference.md` when cited, including nested references. This is a static
instruction walkthrough, not three executed application changes.

1. **Simple fix:** a known single-file bug, Route A, ordinary review, no selected
   heavy lanes. `SKILL.md` → `work.md` → `work-selection-reference.md` validates
   arguments and schema, runs recover before advance, consumes committed claims,
   handles addenda and Open Questions, and chooses A. The shared mechanical gate
   loop retains estimation. Step 3.7 loads a matching bug-fix spec. Step 4 writes
   the skip note without loading planning. Step 5 still reads captured lessons,
   consults the index, applies the capture-reference budget and reports dropped or
   missing entries. It skips only exploration and scope/pre-flight for A.
   `work-implementation-reference.md` loads mandatory and condition-matched crews,
   primes and touch-matched satellites; retains branch isolation, landed-handback
   reuse, P-A-U, captured proof/TDD, decisions, scope boundaries, merge range and
   mandatory Implementation Summary. `work-testing-reference.md` retains actual
   diff/live-flow/requirement qualification, two re-qualification attempts,
   focused tests, the 30-second budget, canonical repository gate and its one
   retry, attribution/deferral, TDD validation and Testing evidence.
   `work-review-reference.md` dispatches `review-work.md`, preserves acceptance
   gates, remediation, lessons, orientation, prime links and KB consent.
   `work-finalization-reference.md` performs follow-up routing, failure/terminal
   judgment, release judgment, exact manifest/commit verification, cleanup and
   checkpoint/loop. No heavy companion is loaded on this successful path.
2. **Complex change:** a multi-component feature, Route C, selected heavy lanes.
   Follow the same shared spine. Before implementation, additionally load
   `work-planning-reference.md`: domain/supply-chain/prime context, requirements
   coverage, orphan-task/scope checks, consumer fields and append-only timestamp.
   After mandatory lesson discovery, load `work-exploration-reference.md`: find
   patterns and tests, preserve prime lessons, declare Scope and mirror write_set,
   collect baseline evidence, and apply repository-gate repair/deferral and saved
   range resumption. After testing and review, `work-heavy-testing-reference.md`
   preserves claimed holds, queue continuation, exhaustion drain, plan-drift and
   dirty-tree refusals, evidence reuse, and green/red/skipped dispositions. Green
   drains finalize in the same turn. Finalization then follows the same owner as A.
3. **Review only:** `SKILL.md` → `review-work.md`, standalone mode. It finds the
   archived target, reads shared principles, REQ and original UR, obtains the
   merge-aware commit diff, checks each requirement and Scope, reviews code/tests,
   runs the conditional restatement sweep, acceptance stages and additional-test
   suggestions, reports scores and evidence, self-validates, captures lessons,
   impact-classifies and folds critical follow-ups, appends Review and commits its
   exact standalone paths. It loads anti-slop before the report and uses its
   existing reference/crew loaders. It never loads selection, planning, building,
   pipeline testing or pipeline finalization merely because the work router exists.

## Recovery coverage

The complete Error Handling table remains in the always-read work entry point.
Its paths therefore remain discoverable even if a phase fails before its next load.

| Trigger | Discoverable owner / disposition |
|---|---|
| Interrupted or stale claim, archive collision, dependency hold | Selection Step 1 → work-reference Crash Recovery, Schema Read Contract and Composed Exit Summary |
| Finalization set aside versus shared dirt | Selection Step 1 → Commit & Metadata-Commit Procedure and Stuck Runs Hand Off to Judgment; per-REQ exclusion, whole-run refusal only for shared dirt |
| Invalid input or unsupported concurrency | Selection Input; existing fan-out reference isolation/fallback contract |
| Plan failure / explore failure | Entry Error Handling → Failure Classification / builder explores with reduced context |
| Missing external precondition or operator wait | Entry Error Handling → Failure Classification Environment row; preserves blocked-flip and routine-operator completion exception |
| Missing lesson file | Entry Step 5 and implementation loader: record and proceed; index budget never substitutes for touched-prime lessons |
| Landed builder handback after interruption | Implementation → fan-out reference; consume, never redispatch |
| Qualification or focused-test regression | Testing: bounded re-qualification, three fix attempts, debugging/testing crews from attempt two |
| Repository gate failure, no-op repair, saved range | Testing and exploration → work-reference Repository Gate Deferral and Resumption; second-run evidence, attribution, fail-safe refusal, validated no-op and saved-range proof |
| Failed ordinary acceptance or failed reviewer | Review and entry Error Handling: one remediation then completed-with-issues; ordinary reviewer failure noted |
| Selected evaluation loop missing/stale approval | Review → ai-evaluation-loop-reference; unresolved stays unapproved, no ordinary fallback |
| Heavy red, skipped lane, drift, dirty tree | Step 7.7 → heavy-testing companion; withdraw commit on red, retain hold on skipped/refused, drain at exhaustion |
| Finalization/hook failure | Entry Error Handling → finalization companion and work-reference; exact continuation retry, no bypass or manual tail |
| Empty queue with pending answers, blocked or held work | Entry Error Handling, Step 10 → Composed Exit Summary, clarify/stakeholder answers or heavy drain |
| Standalone no diff / no commit or invalid repair exception | review-work Step 1: report nothing to review or typed refusal; never invoke builder recovery |
| Standalone archived UR fallback / review findings | review-work Steps 3/10: archived UR stays closed, critical follow-ups fold into queue, other findings report only |

## Preservation and review

The replay script accounts for every nonblank original work-action procedure line.
Only the duplicate orchestrator checklist, two repeated shared rules and the old
skip-note wording are excluded. Checklist obligations route to their existing
numbered steps above; the skip-note condition remains explicit at Step 4. Phase
bodies retain their text and directory depth, so their relative links retain their
meaning. Existing action citations to work steps enter through the phase pointers.
The schema/template reference and review-work remain byte-unchanged; the expanded
measurement now counts the reference sections their loaders name.

The existing contract tests now inspect the authoritative phase owners. Their
negative mutations still remove/comment/negate/displace crew loads. Review found
that a single-step companion exposed an old test extractor accepting a displaced
load after a level-two heading; its boundary now stops at any level 1–3 heading.
The audit scans all phase companions, so moving a duplicate does not hide it.

Release metadata is applied by the canonical `release --manifest` transaction,
including root/core VERSION, version action and byte-identical changelogs.
Consumers use `just do-work-update` after publication to upstream main.

## Validation results

- `python3 _dev/tests/instruction-reading-volume.py`: all three read sets measured;
  preservation replay passed.
- `bash _dev/tests/maintainer-verify.sh`: exit 0. Contract suite, source/install
  reference checks, mirrored changelog check, shell-block lint, ShellCheck, gofmt,
  vet and uncached Go tests passed (427 board tests and 980 CLI tests).
- `git diff --cached --check`: clean after trimming companion trailing blank lines.
- Final self-review covered extraction boundaries, Route A lesson discovery,
  phase ordering, recovery routing, standalone review isolation, preserved
  source-of-truth ownership and release/updater packaging. The cited Scope
  heading and moved shell-guide test owner were corrected before the passing run.
- No heavy test lane was run: this change relocates instruction prose and adjusts
  its source-location checks; executable product code is unchanged.
