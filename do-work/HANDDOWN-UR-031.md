# HANDDOWN — UR-031: Four-Skill Distribution Follow-ups

**Written:** 2026-08-10T10:45:25Z, at v0.186.12 after REQ-160 metadata commit `4edf0d8`.
**Stop condition:** The user-authorized REQ-158/REQ-159/REQ-160 serial run is complete. Do not claim REQ-161 or REQ-162 until the user clarifies them: REQ-161 was flipped to `pending` by separate workspace activity after the user's three-REQ approval, while REQ-162 remains `pending-answers`.

## Completed in this session

| REQ | Outcome | Review | Version | Implementation | Metadata |
|---|---|---:|---:|---|---|
| REQ-158 | Shared offset-preserving Markdown rendered-region classification | 63%, Partial | 0.186.10 | `47b71fd` | `2f1efc6` |
| REQ-159 | Multiline ordinary-quote and triple-backtick Just collision state | 84%, Partial | 0.186.11 | `6ba3a27` | `d19989c` |
| REQ-160 | Occurrence-complete retired-command matching | 98%, Pass | 0.186.12 | `3d8613a` | `4edf0d8` |

Each REQ completed its normal claim, fresh Plan/Explore/Builder/Review contexts, RED/GREEN, owner qualification, independent testing/review, Lessons/Orientation, archive, patch release, implementation commit, and provenance write-back. Metadata-only commit `698324a` corrected REQ-159's initially omitted P-A-U lifecycle section; its recorded implementation hash remains `6ba3a27`.

## What the three REQs delivered

### REQ-158

- Centralized shipped-Markdown visibility decisions in one length-preserving classifier before structural links, reference definitions, and bare first-party URL discovery.
- Added effective-column indentation and top-level paragraph context, escape-parity-aware backtick delimiters, and escaped inline/reference link masking.
- Preserved publication topology, raw/blob/path containment policy, source offsets, and distribution contracts.
- Review reproduced remaining escaped closing-delimiter and list-paragraph variants and created consent-gated REQ-161.

### REQ-159

- Made the root and shipped `replace-text-section.sh` helpers retain raw single-quote, cooked double-quote, and exact triple-backtick state across physical lines.
- Added Just-parseable positives, exact sorted real-collision controls, and byte-preserving pre-mutation rejection checks.
- Kept both helper copies byte-identical and preserved triple-string, managed-span, reserved-name, installer-ordering, and transaction behavior.
- Just 1.46 also accepts ordinary single-backtick multiline commands; that explicitly out-of-scope form remains consent-gated as REQ-162.

### REQ-160

- Changed the test-only retired-command matcher to carry source spans and exempt only the exact approved queue-board branding/test-reference occurrences.
- Continued eligible install/setup candidate evaluation after boundary-invalid longer triggers and restored historical `install-<target>` prefix semantics only inside the recurrence guard.
- Preserved the 186-row historical inventory, 585 direct-boundary negatives, current sibling routes, prime fingerprints, and repaired live surfaces without restoring a runtime alias.
- Independent review passed a 34-case adversarial matrix, focused/full suites, and preservation sweeps; no follow-up was created.

## Live queue at handdown

| REQ | Status | Dependency / action |
|---|---|---|
| REQ-161 | `pending` | Review follow-up from REQ-158. Separate workspace activity flipped it after the user's clarification had already named only REQ-158/159/160. Preserve it, but obtain explicit user direction before claim/implementation. |
| REQ-162 | `pending-answers` | Review/discovered follow-up from REQ-159. One consent question must be answered before it can become claimable. |

`do-work/working/` contains no REQ. UR-031 remains open with 24 completed members and two unresolved follow-ups; completed REQs remain loose in `do-work/archive/` until every member is terminally resolved.

## Live risks and consent boundary

- REQ-161 covers two reproduced Markdown-classification gaps: escaped closing-bracket/opening-parenthesis link forms can leak hidden first-party URLs, and four-space continuations inside bullet/ordered-list paragraphs can hide live links.
- REQ-162 covers a reproduced Just grammar gap: ordinary single-backtick commands may span physical lines, but reserved-looking column-zero payload inside them can still false-collide.
- Both are non-critical review-depth follow-ups. The prior instruction forbids auto-running pending-answers descendants without clarification; REQ-161's externally written `pending` status is not evidence that the user approved this newly created descendant.

## Verification evidence

- REQ-158 and its requalification captured exact RED/GREEN fixtures for effective indentation/paragraph context, escaped backticks, escaped inline links, and escaped reference definitions.
- REQ-159 RED reported all five literal payload names plus an extra hidden control; GREEN retained exact real collision diagnostics and byte preservation through Just-parseable fixtures.
- REQ-160 aggregate RED reproduced seven occurrence/fallback misses; GREEN passed 34 independent adversarial cases, 186 exact row identities, and 585 direct-boundary negatives.
- Post-v0.186.12 `bash _dev/tests/contract-regressions.sh` — PASS, including suite manifest, shipped reference, record-hash/blanked-REQ, updater, staged-skills, and installer probes.
- Focused staged-skills, suite-manifest, `bash -n`, warning-level ShellCheck, changelog/version mirror identity, fixture hash/count, protected-file hashes, and `git diff --check` — PASS.
- Provenance read-back records exactly `47b71fd`, `6ba3a27`, and `3d8613a` in archived REQ-158/159/160 respectively; REQ-160's metadata patch was immediately verified as a one-line `commit:` insertion.

## Cleanup result

- Passes 0–6 completed. No terminal REQ is stranded in queue/working, no misplaced request tree/archive folder exists, no consumed run scratch exists, no orphan `worktree-agent-*` worktree/branch exists, and the blanked-REQ scanner found no blanked or unparseable REQ/UR file.
- UR membership was derived from `user_request:` across queue/working/archive: 24 completed, REQ-161 `pending`, and REQ-162 `pending-answers`. Cleanup correctly left UR-031 open and made no structural move.
- No files moved during cleanup, so no documentation link repoint or cleanup-only commit was required.

## Preserved approved files

The incoming approved content remains byte-identical:

- `decisions/records/adr-019-four-skill-suite-contract.md` — SHA-256 `2d5a54bc9435f8643f4d30e332c37426fbacb15442503eba338cb1f9ab11b282`.
- `do-work/user-requests/UR-031/input.md` — SHA-256 `ed156a18dc11f4f367e80d0e1cca8dbf676dffaae8030622df214e8070bab160`.

Separate workspace activity committed those already-approved bytes during the run (`d9a384d` and `0982c84`); this session did not rewrite them.

## Copy-paste restart prompt

```text
do-work clarify

Resume UR-031 from /Users/t2/Desktop/e1-experimental-repos/skill-do-work2. Read do-work/CHECKPOINT.md and do-work/HANDDOWN-UR-031.md first. Clarify REQ-161 and REQ-162 explicitly before running either. REQ-161 currently says pending only because separate workspace activity changed it after the user's prior approval named REQ-158/159/160; do not treat that as consent for the newly created follow-up. Preserve decisions/records/adr-019-four-skill-suite-contract.md and do-work/user-requests/UR-031/input.md byte-identically.
```
