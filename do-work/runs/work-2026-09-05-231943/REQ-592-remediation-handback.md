# REQ-592 Remediation Hand-Back

REQ-592 is the change that stops the fast maintainer gate from reusing stale test
evidence when files under `do-work/` change. The independent review scored it 74%
(Test Adequacy 55%) and found two mutations that survived the test suite. Both are
closed.

- **Branch:** `worktree-agent-REQ-592-seal-the-do-work-tree-into-both-fast-gate-stages`
- **Branch head:** `f67b97af3f300509652c3dcff9aaf1bc94512957` (one commit on top of `1d2c719`)
- **Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-592-seal-the-do-work-tree-into-both-fast-gate-stages`
- **Not pushed. Nothing staged or committed in `/home/user/skill-do-work`.**

## Files changed

Two of the four in-scope files. The other two needed no change for these findings.

- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go` — comments only (F3a, F3b, F3c)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go` — the new real-manifest test (F1) and the new table case plus its fixture file (F2)
- `_dev/tests/fast-stages.json` — unchanged. It was restored byte-identical after every ablation; `git diff` on it is empty.
- `_dev/tests/fast-stage-reuse-behavior.sh` — unchanged. F2 is a Go-table finding; the shell probe drives one stage against a fixture manifest and adding a duplicate case there would pin nothing the Go case does not.

`heavy_run.go` was not touched.

## F1 — the shipped manifest is now read by a test

**The gap.** Every existing case ran a fixture manifest. Nothing read the real
`_dev/tests/fast-stages.json`, so the data half of the fix could be reverted and the
suite stayed green.

**What was added.** `TestShippedFastStageManifestSealsTheQueuePathsItsStagesRead` in
`fast_stage_evidence_test.go`. It locates the repository root the way the heavy lane's
`heavy_maintainer_tree_test.go` does (`filepath.Abs("../../../../../..")`), reads and
decodes the real manifest, then asserts by name:

1. the `do-work-cli-fast-tests` stage covers `do-work/archive/UR-003/input.md`;
2. the `queue-kanban-fast-tests` stage covers it too (this is the other half of the
   manifest fix — its `do-work` subtree rule — and without this line that rule could be
   deleted with the suite still green);
3. no `seal_exclusions` entry matches that path (an exclusion beats coverage, so one
   that reached the path would undo both rules above without touching either);
4. `non_stage_coverage` does not claim that path is unread;
5. the file itself still exists, so the coverage rules cannot go vacuous silently.

**Self-containment.** The maintainer `_dev/` tree is export-ignored; the do-work-cli
module ships. The test skips when the `_dev/tests` **directory** is absent — that is the
"this is an installed copy" signal. It keys on the directory, not on the manifest file, so
inside a maintainer checkout a renamed, deleted or undecodable manifest is a failure here
rather than a silent skip.

**One tension a maintainer should decide on, flagged not resolved.** `.gitattributes` says
"Keep every maintainer-tree read in that file [`heavy_maintainer_tree_test.go`]; a new one
elsewhere breaks consumers again." This test is a new maintainer-tree read outside that
file. It does not break consumers, because the reason given there — those tests *cannot* be
made self-contained, so they hard-fail in an install — does not apply to a
skip-guarded read. But the instruction is now literally violated. Two honest ways to close
it: widen that `.gitattributes` note to permit skip-guarded reads, or move this test into
`heavy_maintainer_tree_test.go`. The second was not available to me — that file is outside
this REQ's four-file scope.

**Red proof (four ablations, each restored afterwards):**

| Ablation | Result |
|---|---|
| `_dev/tests/fast-stages.json` restored to its pre-fix content (`git show fce57fcc:...`) | FAIL, three assertions fired: do-work-cli does not cover the path, queue-kanban does not cover the path, `non_stage_coverage` claims it is unread |
| Delete the `exact` coverage rule on `do-work/archive/UR-003/input.md` | FAIL: "stage do-work-cli-fast-tests does not cover do-work/archive/UR-003/input.md" |
| Delete `queue-kanban-fast-tests`' `do-work` subtree coverage | FAIL: "stage queue-kanban-fast-tests does not cover do-work/archive/UR-003/input.md" |
| Broaden the `do-work/runs` seal exclusion to `do-work` | FAIL: "a seal exclusion matches do-work/archive/UR-003/input.md" |

## F2 — the tracked half of the exclusion guard is now pinned

**The gap.** Replacing the guard at `fast_stage_evidence.go:220` (the `git ls-files --cached`
loop) with a bare `if path == "" { continue }` left everything green, because every excluded
path a case mutated was untracked. In production that branch carries 231 tracked files under
`do-work/runs` and 162 under `do-work/.req-reservations`.

**What was added.** The fixture repository now commits
`do-work/runs/work-fixture/prior-notes.md` — under the excluded subtree, tracked, and no
stage's toolchain probe input, so the only thing that can move a fingerprint when it changes
is the seal itself. New table case `a tracked run-log file changed still reuses` mutates it
and expects `reused/fingerprint_match` on `beta-stage`, the fixture stage that covers
`do-work`.

**Red proof.** With the tracked-loop guard replaced by `if path == "" { continue }`, running
the whole decision table gives exactly one failure:

```
--- FAIL: TestFastStageReuseDecisionTable/a_tracked_run-log_file_changed_still_reuses
    fast_stage_evidence_test.go:463: decision = executed/fingerprint_mismatch, want reused/fingerprint_match
```

That it is the *only* failure is the point: it is the sole pin on that branch.

## F3 — three comments corrected

- **(a) `workingTreeSeals` doc.** It claimed to seal every tracked and untracked path the
  stage covers. It now says "MINUS every path `SealExclusions` names — an excluded path is
  skipped even when the stage's own coverage matches it."
- **(b) The admission condition.** Now reads "no stage reads its BYTES, NOT ITS EXISTENCE",
  with a paragraph naming the existence-level read that sits outside the test: the board
  stats every repo-relative path mentioned in any REQ or UR body
  (`skills/do-work-board/tools/queue-kanban/filementions.go`), and those mentions reach
  `do-work/runs` and `do-work/deliverables`. It states plainly that today's entries are safe
  only because no fast-stage assertion reads that map.
- **(c) The Closed Enumerations second half.** The struct doc now names
  `skills/do-work-board/tools/queue-kanban/walk.go`'s `isSkippedSection` as the reason
  `do-work/runs`, `do-work/deliverables` and `do-work/.req-reservations` pass the condition,
  and tells the next reader to read that function before adding an entry and to re-check
  these three whenever it changes. Verified against `walk.go:191-198`, which prunes exactly
  `deliverables`, `runs`, `assets` and dotted directory names.

Comments carry no ablation; they are documentation, not behaviour.

## F4 — the Qualification claim about the fixture probe-input move is wrong; correct the record

**I took the "leave the move, correct the claim" option.** Here is why the other option is
not available.

The REQ's `## Qualification` says the two fixture probe inputs were moved out of the newly
sealed tree because otherwise `toolchain probe output changed` would pass "by the file's own
byte seal moving rather than by the probe's output changing", and `toolchain probe cannot run`
would reach `fingerprint_uncertain` "through a missing seal input instead of a failing probe".
The reviewer reverted the move and both cases still passed. That is correct, and the
mechanism the reviewer gives is correct: both cases run `alpha-stage`, whose coverage is
`module-alpha` only, so a `do-work/` path was never sealed into that stage either way.

Making those two cases depend on the move is not achievable with the assertions they carry.
Both assert a disposition and a reason code, and both reasons are identical under either
mechanism: a probe output change and a probe input seal change both yield
`executed/fingerprint_mismatch`, and a removed probe input yields
`executed/fingerprint_uncertain` whether the seal step or the probe step fails first. Nothing
observable at the decision boundary distinguishes them, so a test written to depend on the
move would have to assert on internals — which is worse than leaving the claim corrected in
prose.

**What to correct in the REQ record:** the Qualification bullet beginning "The two fixture
probe inputs were moved out of the newly-sealed tree" overstates what the move achieved. The
move is harmless and keeps the fixture's intent legible (a probe input in an excluded subtree
cannot double as a seal input for a stage that *does* cover the tree, which `beta-stage`
does), but it does not restore anything the two named `alpha-stage` cases test, because
nothing was broken there.

## Follow-ups, not built (as instructed)

1. **A fixture stage with an `exact` coverage rule inside another stage's covered subtree.**
   The do-work-cli half of the fix has that shape and no fixture reproduces it. F1 now covers
   the real instance of it directly against the shipped manifest, which is the case that
   matters; a fixture version would be a second pin on the same behaviour.
2. **A drift ratchet over do-work-cli path literals** — scan the do-work-cli test corpus for
   path literals resolving under `do-work/` and assert the set equals that stage's declared
   coverage. New mechanism, not a test; belongs in its own REQ.
3. **The `.gitattributes` tension described under F1** — decide whether skip-guarded
   maintainer-tree reads are allowed outside `heavy_maintainer_tree_test.go`, and either
   widen the note or move the new test.

Everything else the review reported stays as filed, including the heavy-lane
`heavy_run.go:28-31` stale comment, which was deliberately left alone.

## Verification

Every command below ran through the prescribed environment
(`env -u NODE_OPTIONS -u GIT_CONFIG_COUNT -u GIT_CONFIG_KEY_* -u GIT_CONFIG_VALUE_*
GIT_CONFIG_GLOBAL=<scratchpad>/gitconfig-gate QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium`),
with every exit status read directly from `$?` and never through a pipe.

| Check | Result |
|---|---|
| `go -C skills/do-work/tools/do-work-cli test -count=1 ./internal/heavyverification/` | `ok` 12.98s, exit 0 |
| `bash _dev/tests/fast-stage-reuse-behavior.sh` | `Fast-stage evidence reuse probes passed.` exit 0 |
| `gofmt -l` on the package | no output |
| `go vet ./...` | no output |
| `bash _dev/tests/maintainer-verify.sh` with `DO_WORK_FAST_STAGE_REUSE=off` | `Maintainer verification passed.` exit 0, 82s wall |

The gate ran once, not twice. Both new tests were confirmed to actually execute rather than
skip, with `-v`:

```
--- PASS: TestFastStageReuseDecisionTable/a_tracked_run-log_file_changed_still_reuses (0.04s)
--- PASS: TestShippedFastStageManifestSealsTheQueuePathsItsStagesRead (0.00s)
```

## Still open from the review, unchanged by this work

The release obligation. Two shipped files under `skills/` changed again in this commit, and
no version or changelog moved. `/VERSION`, `/skills/do-work/VERSION` and
`/skills/do-work/actions/version.md:5` all still read `0.303.10`, and both changelogs need an
entry. That stays with the finalizer, and the review's sequencing caution still applies:
REQ-591 (the change that introduced fast-stage evidence reuse) is claimed and unreleased with
its shipped changes already committed, so this version number follows from that one.
