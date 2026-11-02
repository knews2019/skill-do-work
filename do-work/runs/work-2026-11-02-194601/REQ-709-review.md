## Review: REQ-709 (recover ignores do-work notes that no request record can own)

**Approve**: the fix does what the REQ asks. Record-location dirt still refuses. The open points are test gaps and one duplicated layout list, not defects.
Route A (standard-depth review, as briefed) | merge range `c5ae2d8e..5a1d9387`, builder commit `4634893c`

### What's built
- `recover` (and `recover-finalization --discover`) now leaves alone any file under `do-work/` that is outside the request-record locations: queue, working, archive, user-requests, `.req-reservations`, runs, `CHECKPOINT.md`, `calibration-log.tsv`. An inbox note or a handoff note no longer stops recovery, and `do-work commit` sees the note as an ordinary unowned addition.
- Dirt in a record location still refuses with the same code and message. The refusal now names only the record path, not the note.

### Decisions / risks for you
- None to decide. D-03 (other non-record files such as `do-work/lessons-index.md` also stop being evidence) is a real behavior change: a legacy tail with a dirty `lessons-index.md` used to refuse and now settles, leaving the index dirty for `do-work commit`. I ran this case and confirmed it. This matches the user's capture answer "Any file no record can own", and no byte is lost. It also matches how lesson satellites outside `do-work/` were already treated.
- D-01 (keep `CHECKPOINT.md` and `calibration-log.tsv` as evidence) is correct and required. Discovery has semantic proofs for both files (`checkpointRemovalProves`, `calibrationAppendProves`), and the foreign-claim refusal depends on the checkpoint being evidence. Dropping them would change real finalization outcomes.

### Findings

**Important:**
- None.

**Minor:**
- F1. Two entry points changed outcome, and no test pins either one. The claim-only topology and `--assume-sole-releaser` used to refuse when a note sat beside them. Now they proceed. My scratch runs: claim-only + note went from refused (blocked: checkpoint, note, queue and working REQ) to success. Sole-releaser legacy tail + note went from refused (blocked: note) to `cleanup_complete` with commit paths `CHECKPOINT.md`, `calibration-log.tsv`, `owned.txt`. Both new outcomes are correct. Also, the legacy-tail subtests (`finalization_unowned_notes_test.go:30-58`) throw away the owned-path set that `seedSemanticLegacyTail` returns. So "settles exactly as today" is pinned only by the phase and the note being absent, not by an exact `commit_paths` match. The prime trap `guard-on-one-entry-path` asks for the condition to be tested at every entry point. — impact-negligible → report only
- F2. The record-location list at `finalization_discovery.go:1487` repeats the `do-work/` layout that the repository walker encodes with its own literals (`repository_model.go:201`, `:240`, `:253`). There is no shared declaration. If a record folder is added to the walker but not to this list, recovery treats it as a note and lets its dirt through: the failure is fail-open. The 8-row table test pins today's set but cannot detect a new folder. I agree with the builder's Discovered Task. I restamp it because fixing it creates one declaration of record locations that both readers use. — impact-rule-change → report only
- F3. A REQ-shaped file outside the record folders, for example an untracked `do-work/REQ-800-stray.md`, used to stop `recover`. Now `recover` returns success, and `do-work commit` would commit the file as an ordinary unowned change (verified in scratch, old vs new). This follows the REQ's literal location rule and D-03's "every other file". `doctor` still reports it as `STRAY-REQUEST`. — impact-negligible → report only

**Nit:**
- F4. `finalization_discovery.go:163-167`: the `owners` map is built from Implementation Summary paths. It now includes `do-work/` note paths, which the old predicate excluded. Notes never reach the loop that reads `owners`, because of the `dirty` filter at lines 135-147. So keeping a note out of a request's commit group now depends on that one filter. The handoff subtest pins it, and the builder's differential A proves it. — impact-negligible → report only
- F5. REQ `## Qualification` → Judgment calls the hand-back's "Klarna check" a stray word. It is the Klarna Test that `crew-members/shared-principles.md:16` requires (check the user's intended outcome), and the builder applied it correctly. Only the durable record is wrong. — impact-negligible → report only

### Requirements Checklist

- [x] A file under `do-work/` outside the record folders is not finalization evidence. Recovery neither refuses on it nor associates it with a request: delivered (dirty filter + narrowed predicate. The handoff subtest proves a note named in a summary is still not associated).
- [x] The rule is keyed on where request records live and stated once, with no note folders listed: delivered. `sharedFinalizationPath` is the only place in recovery that lists the locations, and `unownedDoWorkNote` derives from it. The prefix sibling `do-work/working-notes.md` is correctly treated as a note. (F2 covers the walker duplicate.)
- [x] Real unfinished finalization work next to a note still refuses or settles exactly as today: delivered. The outcome matches the same tree without the note in every case I ran: 8 record locations, foreign checkpoint line, legacy tail, sole-releaser, claim-only, and worktree renames in both directions. The only difference is that the note is gone from `blocked_paths`.
- [x] `do-work commit` groups the note as an ordinary change with no owning request: delivered. The inventory row is pinned by a test. I verified association end-to-end: `protected-inventory associate` reports `ASSOCIATION-UNOWNED` for the note even when a claimed REQ's summary names it.
- [x] UR-162 (recover should not treat untracked do-work notes as unfinished finalization work): covered by the above, with the "any file no record can own" scope.

### Acceptance Testing

**Result: Pass** (stages covered: implementation and integration)
- Ran the briefed commands from `skills/do-work/tools/do-work-cli`. The three finalization tests passed (13 subtests plus 2 Finder-metadata subtests). `go test -run 'TestRecover' ./internal/lifecycleadvance/` passed (7 tests, including `TestRecoverLeavesAnUntrackedInboxNoteForCommit`).
- Confirmed RED myself: put the pre-change `finalization_discovery.go` from `c5ae2d8e` in a scratch copy (`/tmp/req709-review/cli-old`). All 12 new finalization subtests and the lifecycleadvance test fail there.
- Full `internal/finalization` (55.5 s) and `internal/lifecycleadvance` (26.5 s) packages pass on the merged code. Two `corehelpers` inventory tests failed only in the scratch copy, because it lacks the sibling `checks/` script. Both pass in the repo.
- End-to-end with a built binary on a scratch repo (claimed working REQ whose summary names `do-work/inbox/note.md`, plus that untracked note): new `recover` returns success with `finalization_passed: true` and the claim preserved. Then `protected-inventory start` prints `A do-work/inbox/note.md`, and `associate` reports it unowned. The pre-change binary on the same repo refuses with `FINALIZATION-DISCOVERY-AMBIGUOUS` and names the note.
- Old-vs-new scratch differentials: claim-only + note, sole-releaser + note, stray REQ, legacy tail + `lessons-index.md`, worktree rename queue→inbox and inbox→queue. Results are under F1, F3 and D-03 above. Both renames refuse before and after. Only the note side drops out of `blocked_paths`.
- Rename rows reach discovery only through intent-to-add (`git add -N`), because staged state refuses first. No do-work tool uses intent-to-add. One theoretical difference needs `git add -N` and `--assume-sole-releaser` together: a UR file renamed into a note folder. Before, the whole rename was attributed. Now only the deletion side is committed, and the note copy stays on disk. I do not record it as a finding.
- Journaled (non-discovered) finalization: code read only. `commitSafety` skips the shared-path check when `journal.Discovered` is false (`finalization_apply.go:520-522`), so the narrowed predicate cannot reach it.

### Suggested Additional Testing

- Deployment: unassessed. After the release that ships this CLI, confirm the installed binary's `recover` passes with an untracked `do-work/inbox/` note.
- Live acceptance: unassessed. In this checkout, leave one untracked inbox note, run `do-work commit`, and confirm it commits the note as an unassociated group, without a `.git/info/exclude` workaround.
- Regression: add claim-only + note and `--assume-sole-releaser` + note cases that assert the exact outcome, and compare `commit_paths` against the fixture's owned set in the legacy-tail subtests (F1).

### Scores (on the record, not the headline)

**Overall: 93%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All four REQ requirements and the UR delivered |
| Code Quality | 92% | Small, clear, condition in a doc comment. F4 coupling |
| Test Adequacy | 80% | Real RED-GREEN, both halves proven by differential, negative table. Two changed entry points unpinned (F1) |
| Scope | 100% | Three files, as stated. No prose change needed |
| Risk | Low | Recovery now leaves more paths alone. A future record folder would fail open (F2). Exact-path commits mean no byte is lost |
| Acceptance | Pass | Implementation and integration: focused and full-package tests, end-to-end CLI on a scratch repo |

### Follow-ups created
- None (5 findings report only)

### Restatement sweep
Redefined element: which `do-work/` paths count as finalization evidence (the meaning of `sharedFinalizationPath`, which used to be "every path under `do-work/`"). I grepped `skills/` prose and Go comments, `README.md` and `skills/*/docs/` (the repo has no top-level `docs/`). No stale restatement found. Checked and still true: `actions/work-reference.md:721` (dirt no REQ owns refuses), `actions/commit.md:51`, `next-steps.md:32`, `actions/run-with-recovery.md:10,47` (historical), `actions/fan-out-reference.md:57` (a different contract: who writes `do-work/` from worktrees), `docs/prescribed-shell-primitives.md:67` (inventory, unchanged), `corehelpers/inventory.go:259` (ownership: `AssociateProjectPaths` still never owns a `do-work/` path), `prime-do-work-cli.md:63` (0.266.4 chokepoint lesson. D-02 explains why this fix belongs in finalization: both inventory readers that refuse are covered by the one predicate).

### Self-validation
- I checked error paths and entry points beyond the builder's tests (claim-only, sole-releaser, renames, stray REQ, lessons-index), and checked every old outcome against the pre-change code, not by reading alone.
- I did not run the journaled path or `--take-over`. The journaled path cannot reach the predicate. `--take-over` uses the same strict discovery as plain `recover`.
- No new Important finding came out of this pass.

---

Proposed block for the REQ (orchestrator stamps the timestamp):

```markdown
## Review

**Overall: 93%** | <timestamp>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 80% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1. Claim-only topology and `--assume-sole-releaser` changed from refusing to proceeding when a note sits beside them (correct, verified old vs new in scratch), but no test pins either entry. The legacy-tail subtests (`finalization_unowned_notes_test.go:30-58`) also discard the fixture's owned-path set, so they check only the phase and that the note is absent, not exact `commit_paths`. — impact-negligible → report only
- F2. The record-location list (`finalization_discovery.go:1487`) repeats the walker's layout (`repository_model.go:201,240,253`) with no shared declaration. A record folder added to the walker alone would be treated as a note, so its dirt would pass recovery (fail-open). Agrees with the builder's Discovered Task, restamped. — impact-rule-change → report only
- F3. A REQ-shaped file outside the record folders (for example untracked `do-work/REQ-800-stray.md`) no longer stops `recover`, and `do-work commit` would commit it as an unowned change. This follows the REQ's location rule and D-03. `doctor` still reports `STRAY-REQUEST`. — impact-negligible → report only
- F4 (Nit). `owners` (`finalization_discovery.go:163-167`) now includes `do-work/` note paths from Implementation Summaries. Keeping a note out of a request's group rests on the `dirty` filter alone (pinned by the handoff subtest). — impact-negligible → report only
- F5 (Nit). `## Qualification` calls the hand-back's "Klarna check" a stray word. It is the Klarna Test that `crew-members/shared-principles.md:16` requires, and the builder applied it correctly. — impact-negligible → report only
**Acceptance:** Pass — implementation and integration: briefed focused tests and full `finalization` and `lifecycleadvance` packages green, RED confirmed on pre-change code, and an end-to-end built CLI on a scratch repo shows `recover` succeeding and `associate` reporting the note unowned. Deployment and live acceptance unassessed.
**Restatement sweep:** redefined which `do-work/` paths are finalization evidence (`sharedFinalizationPath`). Swept `skills/` prose and Go comments, `README.md`, `skills/*/docs/`. No stale restatement.
**Suggested testing:** 3 items
**Follow-ups created:** None (5 findings report only)

*Reviewed by review-work action*
```

### Re-check (5a1d9387..6e2ae9a4)

**F1 is closed.** The range is test-only: builder commit `6582e1d4` changes one file, `finalization_unowned_notes_test.go` (+61/-15). No production code changed. Per the orchestrator, F5 was corrected in the REQ's Qualification. F2, F3 and F4 stay report only.

- **Sole-releaser entry is pinned.** The new subtest runs `--discover --assume-sole-releaser` on `seedSoleReleaserLegacyTail` beside the inbox note. It asserts `cleanup_complete`, exact `commit_paths` `[do-work/CHECKPOINT.md do-work/calibration-log.tsv owned.txt]`, and exact sole-releaser attributed paths. These values match what my scratch run produced, and the attributed set matches the existing REQ-499 test. It also asserts the note is still `??`.
- **Claim-only entry is pinned.** The new subtest uses the `TestRecoverFinalizationAcceptsCoherentClaimOnlyTopology` fixture plus the note. The claim leaves the checkpoint, queue and working paths dirty, so `coherentClaimOnlyTopology` is really reached. It asserts success, no finalizations, exactly one `FINALIZATION-NONE` finding, a byte-identical working REQ, and the note still `??`.
- **Exact `commit_paths` checks added.** Both earlier legacy-tail subtests now go through `assertSettledAroundNote`. The semantic tail is compared against the fixture's own `ownedPaths`, and the handoff tail against `[req-750.txt]`. This replaces the weaker "note not in commit paths" loop. The builder's variant run (dirty filter off) shows the exact check catches a note joining the group on its own.
- **RED and GREEN confirmed independently.** I put the merged test file into my pre-change scratch copy (`/tmp/req709-review/cli-old`): all 5 subtests of `TestRecoverFinalizationIgnoresNotesNoRequestRecordOwns` fail there, including both new ones. On `6e2ae9a4` in the repo, the briefed command passes all 5 + 9 subtests. I did not run the full packages, as asked.
- **New findings:** none. The claim-only subtest copies about 10 setup lines from the older claim-only test instead of sharing a helper. D-10 records why it asserts exactly one `FINALIZATION-NONE` finding. Neither is worth a finding.

**Updated scores:** Requirements 100%, Code Quality 92%, Test Adequacy 90% (was 80%), Scope 100%, Risk Low, Acceptance Pass. **Overall: 95%** (was 93%). Verdict unchanged: Approve.

Proposed replacements in the `## Review` block above: `**Overall: 95%**`, `| Test Adequacy | 90% |`, the F1 line gets the suffix `(closed in 6582e1d4 / merge 6e2ae9a4: sole-releaser and claim-only subtests plus exact commit_paths, RED-GREEN re-confirmed)`, the F5 line gets `(corrected in Qualification)`, and `**Follow-ups created:** None (5 findings report only)` stays.
