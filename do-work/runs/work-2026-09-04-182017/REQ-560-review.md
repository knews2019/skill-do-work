# Review: REQ-560

**Approve** — foreign dirt no longer stops or gets committed by the pipeline; the refusal the REQ promised to keep is intact, proved by an independent probe, not assumed.
Route B | merge range `6adb8b9..cb3a831` (merge `cb3a831`, builder commit `4a2cee2`)

REQ-560 is "hand-back and finalize check cleanliness only on the REQ's own paths" — a path another session or the user left dirty is no longer a reason to stop the run, and the pipeline no longer authors a commit to hold somebody else's bytes.

### What's built

- `commitSafety` in `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go` now applies its shared-remainder refusal (`FINALIZATION-AMBIGUOUS-SHARED-STATE`) only when `journal.Discovered` is true. A journaled `finalize --manifest` declares its exact write set, so uncommitted paths it never declared are skipped instead of refused.
- Hand-back step 0's third category changed from "stop and surface" to "leave alone and name" in both `skills/do-work/actions/work-reference.md` (canonical) and `skills/do-work/actions/work.md` Step 6 (condensed copy). A foreign staged path is taken out of the index with `git restore --staged` rather than committed.
- The one sentence that authorized a pipeline-authored "preserve" commit was rewritten to forbid it.
- One new test, `TestFinalizeIgnoresForeignTreeDirtOutsideTheManifest`, pins the narrowing.

Still missing: nothing the REQ required, but see F5 — the CLI itself names no skipped foreign path, so requirement 5's "one line per foreign path" is carried entirely by the prose lane.

### Decisions / risks for you

- **D-01 is correct, and I checked its premise rather than taking it.** The builder deviated from the REQ's literal sentence ("a foreign staged path is unstaged only if this run staged it, otherwise it is left and named") and unstages every foreign staged path. Its argument is that the hand-back sequence uses whole-index commits, so a foreign entry left in the index is not left — it is committed. I read `internal/gittransaction/exact_commit.go:37-68`: `CommitExactPaths` refuses a non-empty index, then runs `git add -A -- <dirty declared paths>` followed by a plain `git commit`, which takes the whole index. Step 0's own prose likewise ends in a plain `git commit`. The premise holds. The REQ's `## What` forbids both stopping and committing, so unstaging is the only remaining outcome, and it costs the other writer a staging flag rather than bytes. Faithful to intent, not an overreach.
- **The safety property is real but narrower than the hand-back sentence suggests.** In the journaled lane the surviving refusal inside `commitSafety` fires only for an allowlisted row classified `X`/`XD` (the secret-name class from `internal/corehelpers/inventory.go:163`). An ordinary dirty declared path does not refuse there — and must not, because that is exactly what finalize is about to commit. The REQ's own paths are still bounded by four untouched guards: the `FINALIZATION-DIRTY-INDEX` refusal at `finalization_apply.go:458`, the same empty-index refusal in `CommitExactPaths`, the `expected_request_sha256` / `expected_checkpoint_sha256` preimage checks in `finalization_prepare.go:51-57`, and the committed-paths equality check at `exact_commit.go:92`. Net: the promise holds at the transaction level. Detail in F6.
- **Discovery keeps its refusal in full.** `--discover` is the only place `Discovered: true` is set (`finalization_discovery.go:1601`), and the new `continue` fires only when it is false. Nothing inside the `if journal.Discovered` branch changed. Verified by probe, below.

### Findings

**Important:**

- F1 — `skills/do-work/actions/work-reference.md:490` claims `git restore --staged -- <path>` "leaves the file's bytes exactly as they are so nothing the other writer typed is lost". That is false for a path staged and then modified again (porcelain `MM`). I reproduced it with git 2.43.0: index content `staged-version`, worktree content `worktree-version`, after `git restore --staged` the status is ` M` and `staged-version` exists only as a dangling blob recoverable through `git fsck --lost-found`. The working-tree bytes survive; a distinct staged snapshot does not. In a rule whose whole purpose is protecting another writer's work, an agent reading this sentence will believe the operation is always lossless and will not warn the owner. `work.md:332` carries the same claim in shorter form ("the file's bytes are untouched"). — impact-user-visible → report only
- F2 — No test anywhere pins the branch this REQ deliberately kept. `grep -rn "AMBIGUOUS-SHARED-STATE" --include=*.go` returns exactly one hit, the production string at `finalization_apply.go:508`; zero test files reference it, before or after this change. The REQ pinned the narrowing (`TestFinalizeIgnoresForeignTreeDirtOutsideTheManifest`) but left the discovery refusal and the declared-path refusal untested, so a future edit can delete either silently. I verified both by hand (see Acceptance Testing); a checked-in test would keep that verification. — impact-rule-change → report only

**Minor:**

- F3 — `work.md:332`'s condensed copy says a third-category path "is another session's"; the canonical `work-reference.md:490` says it "belongs to another session **or to the user**". The user's own uncommitted work is precisely the class the originating incident was about (a hand-edited `do-work/calibration-log.tsv`). The two texts are a hand-maintained duplicate pair with no mechanical check, so this is exactly where drift starts. — impact-negligible → report only
- F4 — `work-reference.md:490` says only a dirty path this REQ owns stops, listing "its working REQ file and **checkpoint entry**". Finalize's actual guard is a whole-file `expected_checkpoint_sha256` digest (`finalization_prepare.go:55`), so any foreign edit to `do-work/CHECKPOINT.md` — including another session's entry — refuses. The prose promises per-entry granularity the code does not have. — impact-negligible → report only
- F5 — Requirement 5 ("one line in the run's progress output per foreign path left alone") is delivered only in the step-0 prose lane. `finalize --manifest` now `continue`s past foreign dirt with no finding, no path list, and no informational output, where it previously named those paths in the refusal. At Step 9 the maintainer learns nothing from the CLI about what was skipped. In practice step 0 already named the same paths earlier in the run, which is why this is Minor rather than a missed requirement. — impact-negligible → report only
- F6 — The Implementation Summary states "The declared-path check that fires first is untouched, so a dirty path the REQ owns still refuses." That check is `row.Classification == "X" && allowed[row.Path]`, which fires only for the secret-name class, not for ordinary dirt on a declared path. The sentence reads broader than the code. The behaviour is correct; the hand-back description of it is imprecise. — impact-negligible → report only
- F7 — `finalization_req560_test.go:66-68` accepts either `" M "` or `"M  "` for the foreign tracked path. `"M  "` means staged, which would mean the path was *not* left alone. Only `" M "` is the correct expectation; the `diff-tree` assertion below catches the worst case, so the looseness does not hide a real defect today. — impact-negligible → report only

**Nit:**

- F8 — The touched line swaps the inline `row.Path == "do-work" || strings.HasPrefix(row.Path, "do-work/")` for the pre-existing `sharedFinalizationPath(row.Path)` helper (`finalization_discovery.go:1554`). I confirmed the helper is byte-identical in behaviour and predates this REQ, so it is a behaviour-free de-duplication on a line already being edited. Strictly it is one edit the REQ did not ask for. — impact-negligible → report only

### Requirements Checklist

- [x] Step 7 step 0: third category becomes "leave alone and name"; stage and allow categories stay; only a dirty path the REQ itself owns stops — delivered in both `work-reference.md:490` and `work.md:332`
- [x] Index clean of the REQ's own paths, not every path — delivered as "Step 0 ends with the index holding nothing but this REQ's own stage set"
- [~] Foreign staged path unstaged only if this run staged it, otherwise left and named — **deliberate deviation**, documented as D-01, judged correct (see Decisions above). Every foreign staged path is unstaged. The literal clause is not implemented; the intent it serves is
- [x] Manifest exact-allowlist validation unchanged — `finalization_prepare.go:135-148` is untouched by the diff; it still refuses a `commit_paths` that omits any planned lifecycle or release target
- [x] A foreign modified or untracked path outside the manifest never refuses the transaction — delivered and proved (Acceptance Testing, A1)
- [x] The tree-dirt check narrowed to declared paths, pinned by a test in that package — `finalization_req560_test.go`, red-green reproduced independently (A2)
- [x] Every sentence authorizing a pipeline-authored "unrelated work" or "preserve" commit removed — `grep -rni "unrelated-work commit|preserve commit|its own unrelated|as its own commit" skills/ _dev/` returns nothing outside `CHANGELOG.md` history. The one live sentence at `work-reference.md:355` now forbids the commit instead of authorizing it
- [~] One line per foreign path left alone in the progress output — delivered in prose for step 0; the CLI emits none at Step 9 (F5)
- [x] Claim transaction's dirty-target refusal unchanged — `internal/gittransaction/` and the claim path are outside the diff
- [x] Discovery (`recover-finalization --discover`) keeps its original refusal — proved (A3)

### Acceptance Testing

**Result: Pass**

Go 1.26.1. `gofmt -l` empty, `go vet ./...` clean, `go test -count=1 ./...` green across the whole module in the real repository tree. Everything below ran against a scratchpad copy so the working tree stayed untouched.

- **A1 — the merged change would have let the REQ-559 finalize proceed.** The new test builds exactly the incident shape: a modified tracked `do-work/calibration-log.tsv` and an untracked `do-work/audits/maintainability-draft.md`, neither declared in the manifest. `finalize --manifest` returns success with a primary commit, both foreign paths are still dirty afterwards, and `diff-tree` on the primary commit contains neither. Untracked `do-work/` files belonging to other REQs are the same class as the three that blocked REQ-559 (`REQ-515-handback.md`, `REQ-559-review.md`, `REQ-560-handback.md`).
- **A2 — RED reproduced independently.** I deleted only the three new lines (`if !journal.Discovered { continue }`) from a copy and re-ran the test. It fails with `Outcome: rolled_back`, `Code: FINALIZATION-AMBIGUOUS-SHARED-STATE`, `AffectedPaths: [do-work/audits/maintainability-draft.md, do-work/calibration-log.tsv]` — the same code and shape as the live REQ-559 refusal. Restoring the guard makes it pass. The red-green claim in `## Testing` is honest.
- **A3 — the declared-path refusal and the discovery refusal both still fire.** I called `commitSafety` directly on a fixture holding one dirty declared path and one dirty foreign path:
  - journaled (`Discovered: false`), declared path secret-classified: `FINALIZATION-AMBIGUOUS-SHARED-STATE`, paths `[do-work/credentials-note.md]` — the REQ's own path still stops, the foreign path is no longer named.
  - discovered (`Discovered: true`): `FINALIZATION-AMBIGUOUS-SHARED-STATE`, paths `[do-work/credentials-note.md, do-work/other-session-draft.md]` — discovery's refusal is intact and still names the unattributed shared path. Not narrowed.
  - journaled, ordinary declared dirty path: no refusal — correct, that is the path finalize is about to commit.
  So the answer to "would it still have refused had one of REQ-559's own declared paths been dirty" is: yes for the protected/secret class inside `commitSafety`, and yes at the transaction level through the four untouched guards named under Decisions. It does not refuse on ordinary dirt on a declared path, which is the pre-existing and correct behaviour.
- **A4 — `git restore --staged` behaviour checked against real git** (2.43.0), three cases: a newly staged untracked file returns to `??` with bytes intact; a staged tracked modification returns to ` M` with worktree bytes intact; a staged-then-modified path loses its staged snapshot from the index (F1).
- **A5 — Restatement Sweep.** The redefined elements are `FINALIZATION-AMBIGUOUS-SHARED-STATE`'s trigger condition and hand-back step 0's third category. Consumers grepped: `next-steps.md:32` and `run-with-recovery.md:10` both route this refusal to `do-work run-with-recovery`, which resolves the discovery lane — still accurate. `work-reference.md:355-356`'s stuck-run bullets are consistent with the new rule. Only `work.md` and `work-reference.md` restate step 0 (`grep -rln "allow but never stage|cached-name guard"`), and both were updated. `lessons-do-work-cli.md:5` and `prime-do-work-cli.md:25` describe discovery, unchanged. `CHANGELOG.md:3607` records the old rule as history and is correct to leave.
- **P-A-U:** all three boxes `[x]`. **Scope:** four files declared, four changed, no drift. **Diff hygiene:** no debug artifacts; the new comment documents why the refusal was narrowed and cites the REQ, which is the kind of comment to keep.

### Suggested Additional Testing

- Add a test asserting `recover-finalization --discover` still refuses on an unattributed shared path. It is the branch this REQ deliberately kept and nothing pins it (F2).
- Add a test asserting a journaled finalize still refuses on a protected/secret-classified declared path, so the surviving half of `commitSafety` is locked in too.
- Exercise the real step-0 sequence once with a concurrent session holding a staged-then-modified `do-work/` path, and confirm what the maintainer is told about the dropped staged snapshot (F1).
- Watch the next multi-REQ wave for whether the progress output actually names each skipped foreign path once. That is prose-enforced, so only observation confirms it.

### Scores (on the record — not the headline)

**Overall: 88%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 92% | All ten traced; one deliberate documented deviation (D-01) and one requirement carried only by the prose lane (F5) |
| Code Quality | 88% | Three-line change keyed on a condition, not a path list; comment explains why. Prose overstates the safety of unstaging (F1) |
| Test Adequacy | 75% | Red-green reproduced independently; the preserved refusal branches have no test at all (F2), one loose assertion (F7) |
| Scope | 98% | Exactly the four declared files; one behaviour-free helper swap on a touched line (F8) |
| Risk | Low | A safety refusal was removed from the journaled lane by design. Four untouched guards still bound what the transaction commits, and discovery's refusal is unchanged |
| Acceptance | Pass | Module green; A1-A5 above |

Computed: (92 + 88 + 75 + 98) / 4 = 88.25 → 88%. Acceptance Pass, no Critical risk, so no modifier applies. Nit findings carry no weight.

### Follow-ups created

None (8 findings report only)

---

**Verdict: Pass**

Acceptance reasoning: every requirement in `## Detailed Requirements` traces to the diff. The one literal deviation, D-01, was tested against its own premise — `CommitExactPaths` and step 0 both end in whole-index commits, so leaving a foreign staged entry in the index means committing it, which the REQ's `## What` forbids in the same sentence that forbids stopping. Unstaging is the only outcome left, and it is faithful to intent. The manifest's exact-allowlist validation is untouched by the diff and still refuses a `commit_paths` missing any planned target. Discovery's refusal is unchanged and I fired it. The declared-path refusal is unchanged and I fired it. The "stop and surface" and clean-index sentences were changed in both members of the duplicate pair and now say the same thing apart from one cosmetic omission (F3). No sentence authorizing a pipeline-authored preserve commit survives in the shipped tree. The eight findings are all report-only: none is critical, none blocks, and the two Important ones — an overstated safety claim in prose (F1) and an untested preserved branch (F2) — are documentation and coverage gaps rather than defects in the shipped behaviour.

*Reviewed by review-work action — 2026-09-04T19:01:16Z*
