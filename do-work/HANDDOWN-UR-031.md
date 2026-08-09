# HANDDOWN — UR-031: Four-Skill Distribution Follow-ups

**Written:** 2026-08-09T19:26:50Z, at v0.186.8 after REQ-156 metadata commit `654f777`.
**Stop condition:** User requested a safe stop after REQ-156. The queue loop stopped intentionally before REQ-157; no claimed REQ or partial project edit remains.

## Completed in this session

| REQ | Outcome | Version | Implementation | Metadata |
|---|---|---:|---|---|
| REQ-154 | Markdown reference parser boundary hardening; review 67%, Partial | 0.186.6 | `22551dc` | `d987e42` |
| REQ-155 | Exact manual Stop-hook object path and empty-wrapper semantics; review 100%, Pass | 0.186.7 | `c1f8e21` | `916ec52` |
| REQ-156 | Raw/cooked Just triple-string collision scanning; review 82%, Partial | 0.186.8 | `db9cd11` | `654f777` |

Each REQ completed the normal claim, fresh plan/explore/builder/review where required, RED/GREEN, qualification, independent tests, lessons/orientation, archive, patch release, implementation commit, and provenance metadata commit lifecycle. Root and installed-core versions/changelogs are synchronized at 0.186.8.

## Live queue at handdown

| REQ | Status | Dependency / action |
|---|---|---|
| REQ-157 | `pending` | Dependency-ready and next in numeric order. Test-only complete retired-alias vocabulary; never restore runtime aliases or flag ordinary prose. |
| REQ-158 | `pending-answers` | Addendum to REQ-154. Needs user clarification before the rendered-region classification follow-up can enter work. |
| REQ-159 | `pending-answers` | Addendum to REQ-156. Needs user clarification before ordinary multiline quotes/triple-backticks can enter work. |

`do-work/working/` contains no REQ. UR-031 remains open because these three members are unresolved; completed REQs remain loose in `do-work/archive/` until the UR can close.

## Live risks and review follow-ups

- REQ-158: the shipped Markdown guard's independent masking/extraction/bare-URL paths still disagree at escaped first-party links, effective-column indented code, four-space paragraph continuations, and escaped backticks. Current release contracts pass, but valid docs can still false-fail or a broken link can be missed in those bounded cases.
- REQ-159: the Just collision scanner now handles `'''` and `"""`, but valid ordinary multiline single/double strings and triple-backtick literals can still false-collide on column-zero reserved-looking payload. Rejection is pre-mutation and byte-preserving; real collisions remain detected.
- REQ-157: review of the command-restatement sweep found the recurrence guard's alias vocabulary incomplete. The accepted boundary is a comprehensive test-only historical vocabulary, with no runtime compatibility layer and no generic prose ban.

## Verification evidence

- `bash _dev/tests/contract-regressions.sh` — PASS after each REQ and after v0.186.8 release.
- `bash _dev/tests/install-suite-behavior.sh` — PASS.
- `bash _dev/tests/staged-skills-contract.sh` — PASS.
- `bash _dev/tests/suite-manifest-contract.sh` — PASS.
- Paired replacer and installer byte identities — PASS.
- `bash -n` and `shellcheck -S warning` on each changed shell surface — PASS.
- Just fixture/template parsing, changelog identity, and `git diff --check` — PASS.
- `queue-kanban verify` — no findings (its legacy version/changelog probe reported an explicit modular-path skip only).
- Cleanup Passes 0–6 — no stranded terminal REQs, misplaced request data, orphan worktrees/branches, or blanked/unparseable request files.

## Preserved dirty files

These were not staged into REQ or handdown commits:

- `decisions/records/adr-019-four-skill-suite-contract.md` — pre-existing user edit, SHA-256 `2d5a54bc9435f8643f4d30e332c37426fbacb15442503eba338cb1f9ab11b282`.
- `do-work/user-requests/UR-031/input.md` — pre-existing user/clarification edit, SHA-256 `ed156a18dc11f4f367e80d0e1cca8dbf676dffaae8030622df214e8070bab160`.
- `do-work/queue/REQ-157-complete-retired-core-alias-guard.md` — approved clarification update for the next normal REQ lifecycle; leave unstaged until REQ-157 is claimed.

## Copy-paste restart prompt

```text
do-work run

Resume UR-031 from /Users/t2/Desktop/e1-experimental-repos/skill-do-work2. Read do-work/CHECKPOINT.md and do-work/HANDDOWN-UR-031.md first. Process the live claimable queue serially; REQ-157 is next. Do not auto-run REQ-158, REQ-159, or any other pending-answers follow-up unless the user clarifies it first. Preserve the existing edits in decisions/records/adr-019-four-skill-suite-contract.md and do-work/user-requests/UR-031/input.md. Treat the approved REQ-157 queue edit as its normal lifecycle input. Follow skills/do-work/SKILL.md -> skills/do-work/actions/work.md exactly, with fresh plan/explore/builder/review contexts, per-REQ tests/review/version/changelog/implementation+metadata commits, context wipes, cleanup, and a refreshed handdown.
```
