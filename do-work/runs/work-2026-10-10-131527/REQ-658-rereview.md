## Re-review: REQ-658 (do-work-cli `finalize --auto-manifest`), remediation delta

**Approve.** F1 to F4 from the first review (`REQ-658-review.md`) are closed. The delta `457608cc..fe788bf4` adds no dead code and no wrong refusal code, and the shipped prose now matches the code. Three new findings, all negligible and report only.

Delta: 5 files, +37/-34, one builder commit `e09c25d8` and the re-merge `fe788bf4`.

### Evidence

- `go test -C skills/do-work/tools/do-work-cli -count=1 -run AutoManifest ./internal/finalization/`: ok (3.3 s).
- `go vet ./internal/finalization/`: clean. The `context` import was removed together with the deleted non-emit branch, so nothing unused is left.
- `bash _dev/tests/contracts/core-checks.sh`: exit 0, "core-checks contract probes passed."
- Reviewer-only tests in a scratch copy of the module (`scratchpad/int-658/rereview/cli/internal/finalization/zz_rereview_test.go`; nothing added to ROOT): 3 of 3 pass.
  - No `--emit`: code `FINALIZATION-AUTO-INPUT`, exit 1, evidence "--emit <path> is required: pass the emitted manifest to advance --finalization-manifest". `.git/do-work-finalization/` is not created and `do-work/working/REQ-760.md` stays in place.
  - Staged index: code `FINALIZATION-INDEX-NOT-EMPTY`, and `.git/do-work-finalization/` does not exist afterwards.
  - Same id in `working/` and `archive/`: code `FINALIZATION-REQUEST-NOT-WORKING` with evidence `REQUEST-AMBIGUOUS: REQ-760 resolves to 2 repository records`. With the archive copy only: "REQ-760 is in archive, not do-work/working/".

### Item by item

- **F1 (non-emit mode skipped the `advance` phase gate): closed.** `judgedInputProblem` refuses an empty `EmitPath` (`finalization_auto_manifest.go:94-95`). This check runs before any file read, Git call or `journalLocations`. The non-dry `prepareManifestJournal` + `advanceJournal` branch is deleted, and the only `prepareManifestJournal` calls left in the file are dry (`:186`, `:201`). The new assertion in `TestFinalizeAutoManifestRefusesWithoutAMessageFile` is a real pin: under the old code that call returned success, so the `OutcomeRefused` check would fail. D-21 in the REQ records the departure from requirements 4 and 5 and the reason for it.
- **F2 (a second REQ-id resolver): closed.** `requeststate.ResolveTarget(snapshot, requestID, "")` plus a `TreeSection != "working"` check (`:150-157`) is the same shape as `corehelpers/request_writers.go:196-204`. An ambiguous id now refuses before the dry core runs.
- **F3 (prose framed mechanical fields as authored): closed.** `work.md:462` and `work-reference.md:719` now name the judged inputs, then the fields the command fills, and say the command never finalizes. Step 8 item 5 shrank from 1080 to 945 characters. Step 9 (`work.md:467`, "the Step 8 `--emit` file or a hand-built one") is still consistent with Step 8. `lessons-do-work-cli.md:127` now names `--emit` as required and says the phase gate still applies. A repository-wide search for `auto-manifest` finds no other shipped text, and no shipped text describes the old direct mode.
- **F4 (index refusal created `.git/do-work-finalization/`): closed.** The staged-index check now runs before `journalLocations` (`:126-144`). This is the same order as `prepareManifestJournal` (`finalization_prepare.go:66-71`), so an index refusal behaves the same in both commands.
- F5 to F7 stay report only, as the brief says. D-21 corrects the D-16 wording (F5).

### New findings

- N1. No repository test pins the F4 fix. `assertNothingWritten` (`finalization_auto_manifest_test.go:43-54`) calls `journalLocations` itself, which runs `MkdirAll`, so it cannot see whether the command created `.git/do-work-finalization/`. Only the scratch test above proves F4. Fix if wanted: check `os.Lstat(<git-dir>/do-work-finalization)` in the staged-path test before calling `assertNothingWritten`.
- N2. `FINALIZATION-REQUEST-NOT-WORKING` now also wraps `REQUEST-NOT-FOUND`, `REQUEST-AMBIGUOUS` and `REQUEST-IDENTITY-MISMATCH` (`:151-153`), and it drops the paths that `StateRefusal` carries (for an identity mismatch, that is the request path). The inner code is still in the evidence, so the refusal can still be read and still exits 1. D-17 in the REQ ("zero or several working files") and its last sentence about a "non-emit real core" refusal are now stale, and D-21 does not correct D-17. Both are REQ-record text, not shipped text.
- N3. `work-reference.md:719` lists the judged inputs as "terminal transition and status, the commit message, any release payload, the extra … paths, and the provenance mode". It leaves out the failure type and failure error, which `work.md:462` covers as "terminal/failure values". The list is incomplete, not wrong.

### Re-review (delta 457608cc..fe788bf4) | 2026-10-10T15:26:38Z

- Overall: 93% | Acceptance: Pass (focused AutoManifest tests ok; core-checks.sh exit 0; scratch tests prove missing `--emit` → `FINALIZATION-AUTO-INPUT` exit 1 with nothing written, the index refusal creates no `.git/do-work-finalization/`, and a duplicate id refuses `REQUEST-AMBIGUOUS`)
- F1: closed (`--emit` is required at `finalization_auto_manifest.go:94`, the non-emit finalize branch and the `context` import are deleted, and the new assertion fails under the old code). F2: closed (`ResolveTarget` + working-section check at `:150-157`). F3: closed (`work.md:462`, `work-reference.md:719` and `lessons-do-work-cli.md:127` state the judged/mechanical split and "never finalizes", Step 9 is consistent, no other shipped text describes the direct mode). F4: closed (the index check runs before `journalLocations`, proven in a scratch repo).
- N1 F4 is not pinned by a repository test, because `assertNothingWritten` creates `.git/do-work-finalization/` itself through `journalLocations` — impact-negligible → report only
- N2 `FINALIZATION-REQUEST-NOT-WORKING` now wraps not-found, ambiguous and identity-mismatch refusals and drops their paths; REQ decision D-17 still describes the old resolver and the deleted non-emit refusal — impact-negligible → report only
- N3 `work-reference.md:719` judged-input list omits the failure type and failure error (`work.md:462` includes them) — impact-negligible → report only
