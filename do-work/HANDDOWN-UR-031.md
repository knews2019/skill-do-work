# HANDDOWN — UR-031: Four-Skill Distribution Follow-ups

**Written:** 2026-08-09T20:44:31Z, at v0.186.9 after REQ-157 metadata commit `6010d81`.
**Stop condition:** The serial claimable queue is exhausted. REQ-158, REQ-159, and REQ-160 are all `pending-answers`; none may be claimed or auto-run without explicit user clarification.

## Completed in this session

| REQ | Outcome | Version | Implementation | Metadata |
|---|---|---:|---|---|
| REQ-157 | Complete test-only retired core alias inventory; review 73%, Partial | 0.186.9 | `1f7a245` | `6010d81` |

REQ-157 completed the normal claim, fresh plan/explore/builder/review contexts, RED/GREEN, one bounded qualification remediation, independent qualification/testing, review, lessons/orientation, archive, patch release, implementation commit, and provenance metadata commit lifecycle.

## What REQ-157 delivered

- Added `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv` as the single test-only historical inventory reconstructed from the deleted router, moved-command shim, and install-normalization history.
- Covered 186 concrete rows: 117 direct aliases, 22 install targets across `install <target>`, `install-<target>`, and `setup <target>`, plus bare `install`, `setup`, and `install-` heads.
- Replaced the former canonical-action-plus-seven-samples guard with fixture validation, longest-first exact-boundary matching, full root/module row mutations, row-identity checks, branding/prose/current-command controls, and history/archive/fixture exclusions.
- Qualification caught and repaired an over-broad em-dash exemption before review: `Deprecated — do-work kanban` is now a positive regression while the exact board product title remains accepted.
- Runtime routers, help, actions, and all 15 REQ-153 repaired live surfaces remain unchanged; legacy aliases exist only in the test fixture and archived history.

## Live queue at handdown

| REQ | Status | Dependency / action |
|---|---|---|
| REQ-158 | `pending-answers` | Addendum to REQ-154. Needs user clarification before rendered-region classification work can enter the claimable queue. |
| REQ-159 | `pending-answers` | Addendum to REQ-156. Needs user clarification before ordinary multiline quote/triple-backtick work can enter the claimable queue. |
| REQ-160 | `pending-answers` | Addendum to REQ-157. Needs user clarification before occurrence-complete retired-alias matching can enter the claimable queue. |

`do-work/working/` contains no REQ. UR-031 remains open with 21 completed members and three unresolved consent-gated follow-ups; completed REQs remain loose in `do-work/archive/` until every member is terminally resolved.

## Live risks and review follow-ups

- REQ-158: the shipped Markdown guard's independent masking/extraction/bare-URL paths still disagree at escaped first-party links, effective-column indented code, four-space paragraph continuations, and escaped backticks.
- REQ-159: the Just collision scanner handles triple-single and triple-double strings, but ordinary multiline single/double strings and triple-backtick literals can still false-collide on column-zero reserved-looking payload.
- REQ-160: the retired-command inventory is complete, but two matcher edge classes remain. Whole-line test-reference exemptions can hide a second real `queue board` invocation, and longest-prefix selection does not retry a shorter historical install/setup head or interpret the former `install-<target>` route as a prefix after a longer candidate misses its boundary.

All three are non-critical and intentionally consent-gated by the generation-depth rule. A future `do-work run` must continue to skip them until the user answers via `do-work clarify` or edits their status deliberately.

## Verification evidence

- REQ-157 RED: the old partial guard recognized 0/4 required representative aliases (`kanban`, `recall`, `code review`, `describe changes`).
- `bash _dev/tests/staged-skills-contract.sh` — PASS after implementation, qualification remediation, independent review, and v0.186.9 release.
- `bash _dev/tests/contract-regressions.sh` — PASS after implementation and after the release bump.
- `bash _dev/tests/install-suite-behavior.sh`, `suite-manifest-contract.sh`, and `shipped-package-reference-contract.sh` — PASS.
- `bash -n` and `shellcheck -S warning` on `_dev/tests/staged-skills-contract.sh` — PASS.
- Scope drift, fixture integrity/history counts, unique sibling routes, updater-prime fingerprints, REQ-153 repair preservation, changelog identity, and `git diff --check` — PASS.
- `queue-kanban verify` — no findings; its legacy version/changelog probe reported only the known modular-path skip.
- Commit-hash write-back verification — archived REQ-157 is nonblank and HEAD's metadata patch is exactly `commit: 1f7a245`.

## Cleanup result

- Passes 0–6 completed. No terminal REQ is stranded in queue/working, no misplaced request data or archive folder exists, no consumed run scratch exists, no orphan `worktree-agent-*` worktree/branch exists, and the blanked-REQ scanner found no blanked or unparseable REQ/UR file.
- UR-031 membership was derived from `user_request:` across queue/working/archive: 21 completed, three `pending-answers`. Cleanup correctly left the UR open and made no structural commit.
- `skills/do-work/` is the intentional installed core package and contains no request-state subtrees.

## Preserved dirty files

These remain unstaged and byte-identical to the incoming handdown:

- `decisions/records/adr-019-four-skill-suite-contract.md` — SHA-256 `2d5a54bc9435f8643f4d30e332c37426fbacb15442503eba338cb1f9ab11b282`.
- `do-work/user-requests/UR-031/input.md` — SHA-256 `ed156a18dc11f4f367e80d0e1cca8dbf676dffaae8030622df214e8070bab160`.

The approved REQ-157 clarification edit entered its normal archived lifecycle in implementation commit `1f7a245`; it is no longer an outstanding dirty queue edit.

## Copy-paste restart prompt

```text
do-work clarify

Resume UR-031 from /Users/t2/Desktop/e1-experimental-repos/skill-do-work2. Read do-work/CHECKPOINT.md and do-work/HANDDOWN-UR-031.md first. Present REQ-158, REQ-159, and REQ-160 as the complete pending-answers set. Do not flip any of them to pending and do not run do-work until the user explicitly answers each consent question. Preserve the existing edits in decisions/records/adr-019-four-skill-suite-contract.md and do-work/user-requests/UR-031/input.md.
```
