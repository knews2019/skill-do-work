# REQ-592 hand-back — seal the do-work tree into both fast gate stages

**Branch:** `worktree-agent-REQ-592-seal-the-do-work-tree-into-both-fast-gate-stages`
**Branch head:** `1a04355aec897545275d2bdb6778151f672e0e10`
**Commits:** one — `1a04355 [REQ-592] seal the do-work tree into both fast gate stages`

## Implementation Summary

The fast gate's seal skipped every `do-work/` path in both its tracked and its untracked loop, so a
`do-work/`-only edit reused stale evidence. The manifest now states which stage reads what, and one
new concept covers the single case coverage cannot express.

Four files changed, exactly the declared write set.

- `_dev/tests/fast-stages.json` — `do-work` becomes the queue-kanban stage's coverage (that stage
  builds its board from the real tree). `do-work/archive/UR-003/input.md` becomes an exact coverage
  rule on the do-work-cli stage, its one real read there. `non_stage_coverage` becomes `[]`. A new
  `seal_exclusions` list holds `do-work/runs`, `do-work/deliverables`, `do-work/.req-reservations`
  and `do-work/test-durations.tsv`.
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go` — new
  `SealExclusions` manifest field, new `fastStageSealExcludesPath` predicate, validation of the new
  list in `decodeFastStageManifest`, and the two `queueStatePrefix` guards replaced by the predicate.
  The exclusion test stays where the old guard was, before the `stageCovered` test, which is what
  lets an exclusion beat a stage's own coverage. Three comments that now said the opposite of the
  behaviour were corrected. `queueStatePrefix` is still used by `heavy_run.go` and
  `heavy_evidence.go`, so nothing became dead.
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go` —
  fixture mirrors the shipped manifest's shape (only `beta-stage` covers `do-work`), probe inputs
  moved into the excluded subtree, the duration log made Git-ignored the way the real one is. The
  `queue state changed` case became two cases; three new cases were added.
- `_dev/tests/fast-stage-reuse-behavior.sh` — same fixture move, `queue state alone still reuses`
  became `a queue-tree change executes the stage that reads it`, plus a new duration-log reuse case.

Acceptance criteria, all met and shown below: a `do-work/` change forces the stage that reads it;
`non_stage_coverage` is empty and that emptiness was verified against the board walk rather than
assumed; both assertions that pinned the old behaviour were rewritten in place with the failure each
now catches named in a comment; `do-work/test-durations.tsv` still does not invalidate its own stage.

## Decisions

**D-01. `seal_exclusions` states a condition, not a list.** The Go struct doc gives the admission
test — a path the gate or the orchestrator writes WHILE a gate runs, and whose bytes no stage reads —
and says the manifest entries are only today's set of paths passing it. Reason: this is the shape
`_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale warns about, so the next person
adding a churn directory gets a test to apply instead of a list to copy.

**D-02. `do-work/deliverables` is in the exclusion list although it does not exist in this
checkout.** It is a real destination — toolbox actions write reports there while work is running —
and `walk.go:192` prunes it from the board walk by name, so it satisfies the condition exactly as
`do-work/runs` does. A coverage rule matches strings and never touches the filesystem, so an entry
for an absent directory costs nothing.

**D-03. The Go fixture's duration log is Git-ignored, the shell fixture's is not.** The real file is
ignored, so the Go fixture reproduces that shape. The shell fixture leaves it untracked and
non-ignored, which exercises the other untracked branch. Both are sealed without the exclusion, so
both cases bite (proved by the ablation below).

**D-04. One commit, not several.** The manifest, the implementation and both rewritten assertions
have to land together; splitting them would leave a red commit on the branch.

**D-05. A decoding case for an invalid seal exclusion was added.** The new validation loop was
otherwise untested. Without it a typo'd `kind` decodes, silently matches nothing, and the excluded
file gets sealed — which turns reuse off for that stage forever. Ablation proof below.

## Red-Green Evidence

All runs used the sanitized environment from the brief. Gate exit status was captured directly, never
through a pipe.

### RED 1 — gate level, the request's own proof (unmodified tree)

```
bash _dev/tests/maintainer-verify.sh          # run 1, warm: EXIT=0 WALL=86s
printf '\n' >> do-work/archive/UR-003/input.md
# 5608 -> 5609 bytes
bash _dev/tests/maintainer-verify.sh          # run 2
```

Run 2 output:

```
EXIT=0 WALL=33s
85:maintainer-verify: stage queue-kanban-fast-tests: REUSED (fingerprint_match, recorded 2026-09-05T23:28:16Z; per-file budget verdict inherited from that run)
88:maintainer-verify: stage do-work-cli-fast-tests: REUSED (fingerprint_match, recorded 2026-09-05T23:28:41Z; per-file budget verdict inherited from that run)
90:Maintainer verification passed.
```

On that same tree, the stage's own test:

```
go test -short -count=1 -run TestDiscoverRepositoryAcceptsProductionLegacyArchiveInputClass ./internal/repositorymodel/

--- FAIL: TestDiscoverRepositoryAcceptsProductionLegacyArchiveInputClass (0.00s)
    repository_model_test.go:407: production legacy fixture changed size: got 5609 bytes
FAIL
FAIL	github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel	0.002s
FAIL
```

### RED 2 — unit level, assertion failure against the unchanged guard

Both rewritten assertions were written first, with the fixture manifests carrying only the coverage
change and no `seal_exclusions` key, so the decoder accepted them and the failure is the assertion
itself rather than a decode error.

```
go test -count=1 -run TestFastStageReuseDecisionTable ./internal/heavyverification/

--- FAIL: TestFastStageReuseDecisionTable (1.80s)
    --- FAIL: TestFastStageReuseDecisionTable/queue_state_changed_forces_the_stage_that_reads_it (0.06s)
        fast_stage_evidence_test.go:442: decision = reused/fingerprint_match, want executed/fingerprint_mismatch
FAIL
FAIL	github.com/knews2019/skill-do-work/do-work-cli/internal/heavyverification	1.799s
FAIL
```

### RED 3 — end-to-end probe, same state

```
bash _dev/tests/fast-stage-reuse-behavior.sh

FAIL: a queue-tree change executes the stage that reads it ran=no (want yes) status=0 (want 0) output=<maintainer-verify: stage alpha-stage: REUSED (fingerprint_match, recorded 2026-09-05T23:32:08Z; per-file budget verdict inherited from that run) > want-line=<EXECUTING (fingerprint_mismatch)>
EXIT=1
```

### GREEN 1 — unit level and probe

```
go test -count=1 -run TestFastStageReuseDecisionTable ./internal/heavyverification/
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/heavyverification	1.372s

bash _dev/tests/fast-stage-reuse-behavior.sh
Fast-stage evidence reuse probes passed.
EXIT=0
```

### GREEN 2 — gate level, the same sequence that was red

Warm run first (both stages executed, exit 0, 75s), then the same one-newline edit:

```
printf '\n' >> do-work/archive/UR-003/input.md
# 5609 do-work/archive/UR-003/input.md
bash _dev/tests/maintainer-verify.sh

EXIT=1 WALL=76s
85:maintainer-verify: stage queue-kanban-fast-tests: EXECUTING (fingerprint_mismatch)
89:maintainer-verify: stage do-work-cli-fast-tests: EXECUTING (fingerprint_mismatch)
2689:    repository_model_test.go:407: production legacy fixture changed size: got 5609 bytes
2690:--- FAIL: TestDiscoverRepositoryAcceptsProductionLegacyArchiveInputClass (0.00s)
3593:go-test budget: FAIL module=.../skills/do-work/tools/do-work-cli wall=21s exit=1
```

`grep -c 'Maintainer verification passed'` on that log returns `0`. The gate now fails on the tree it
used to pass.

### GREEN 3 — the duration log still reuses

Reset the fixture file, re-warm the gate (exit 0), then touch only the log:

```
# test-durations.tsv 61325 -> 61344 bytes
bash _dev/tests/maintainer-verify.sh

EXIT=0 WALL=31s
85:maintainer-verify: stage queue-kanban-fast-tests: REUSED (fingerprint_match, recorded 2026-09-05T23:38:49Z; per-file budget verdict inherited from that run)
88:maintainer-verify: stage do-work-cli-fast-tests: REUSED (fingerprint_match, recorded 2026-09-05T23:39:11Z; per-file budget verdict inherited from that run)
90:Maintainer verification passed.
```

### GREEN 4 — per-stage separation holds

A `do-work/`-only edit outside the do-work-cli stage's one covered file (`do-work/CHECKPOINT.md`):

```
EXIT=0 WALL=56s
85:maintainer-verify: stage queue-kanban-fast-tests: EXECUTING (fingerprint_mismatch)
89:maintainer-verify: stage do-work-cli-fast-tests: REUSED (fingerprint_match, recorded 2026-09-05T23:39:11Z; per-file budget verdict inherited from that run)
91:Maintainer verification passed.
```

This is the reason coverage was used instead of deleting `non_stage_coverage`: the stage that reads
the tree runs, the stage that does not keeps its reuse.

### Ablation — the three new lock-in cases genuinely bite

`seal_exclusions` was removed from both fixture manifests, keeping everything else:

```
--- FAIL: TestFastStageReuseDecisionTable/the_gate's_own_duration_log_still_reuses (0.03s)
        fast_stage_evidence_test.go:442: decision = executed/fingerprint_mismatch, want reused/fingerprint_match
    --- FAIL: TestFastStageReuseDecisionTable/a_run-log_write_still_reuses (0.03s)
        fast_stage_evidence_test.go:442: decision = executed/fingerprint_mismatch, want reused/fingerprint_match

FAIL: the gate's own duration log alone still reuses ran=yes (want no) status=0 (want 0) output=<maintainer-verify: stage alpha-stage: EXECUTING (fingerprint_mismatch) > want-line=<REUSED (fingerprint_match, recorded >
EXIT=1
```

And with the seal-exclusion validation loop removed from the decoder:

```
--- FAIL: TestFastStageManifestDecodingRefusesAmbiguity/unsupported_seal_exclusion_kind (0.00s)
        fast_stage_evidence_test.go:581: decoding accepted unsupported seal exclusion kind
```

Both fixtures and the decoder were restored immediately after each ablation, and the final gate run
below was made on the restored state.

## Testing

- `gofmt -l .` in the do-work-cli module — no output.
- `go vet ./...` in the do-work-cli module — clean.
- `go test -count=1 ./internal/heavyverification/` — `ok`, 12.4s.
- `shellcheck _dev/tests/fast-stage-reuse-behavior.sh` (ShellCheck 0.11.0) — exit 0.
- `bash _dev/tests/fast-stage-reuse-behavior.sh` — `Fast-stage evidence reuse probes passed.`
- Full canonical gate on the committed state, run in the worktree:

```
EXIT=0 WALL=75s
85:maintainer-verify: stage queue-kanban-fast-tests: EXECUTING (fingerprint_mismatch)
89:maintainer-verify: stage do-work-cli-fast-tests: EXECUTING (fingerprint_mismatch)
91:maintainer-verify: gate wall 75s
92:Maintainer verification passed.
```

Both fast stages executed in that run rather than reusing, so the whole suite really ran against the
changed code. No test from a prior REQ broke; the two that changed are the two the request named.

## Discovered Tasks

- **impact-critical** — The heavy lane has the same false-green shape for an *uncommitted tracked*
  `do-work/` edit. Its dirty-tree refusal exempts `do-work/` (`heavy_run.go:221` and `:225` both test
  `!strings.HasPrefix(..., queueStatePrefix)` before recording an offending path) and its untracked
  seal skips the tree too (`heavy_evidence.go:481`), while its committed seal only sees HEAD objects.
  The one tree the refusal skips is the tree the committed seal cannot see. REQ-592's own text says
  the heavy exclusion is safe "because it refuses a dirty tree", which is the claim that does not
  hold. Out of this request's write set; the fix shape is the one just built for the fast gate.
- **impact-noncritical** — The fast-stage seal is byte-level, but part of the queue-kanban stage's
  real dependency on `do-work/` is existence-level: `collectRepoFileMentions`
  (`skills/do-work-board/tools/queue-kanban/filementions.go:35-56`) stats every repo-relative path
  mentioned in any REQ or UR body, and `generate.go:713` runs it on every live board build. Those
  mentions reach `do-work/runs/` and `do-work/deliverables/`, so creating or deleting a mentioned
  file in an excluded subtree flips a boolean in the shipped board JSON without moving any seal. No
  current fast-stage assertion reads that map, so nothing is failing today. → report only
- **impact-noncritical** — The queue-kanban fast stage will now rarely reuse during a drain: with
  `do-work` as its coverage it seals about 730 tracked files, and every REQ claim, move or archive
  touches one of them. This is correct — the stage really reads those bytes — and it is the whole
  cost of the change, but it means REQ-591's reuse feature mostly benefits the do-work-cli stage from
  now on. Measured: a queue-only edit re-runs queue-kanban in 56s where a fully warm gate is 31s.
  → report only
- **impact-noncritical** — This commit changes a shipped file under `skills/`, so
  `_dev/primes/prime-releases.md` makes it a release: it needs a version bump, a `CHANGELOG.md` entry
  and the byte-identical `skills/do-work/CHANGELOG.md` mirror. That is Step 9 / finalization work, not
  the builder's, and doing it here would have broken the four-file scope. → report only

## Declared but not touched

- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go`,
  `heavy_evidence.go`, `_dev/tests/heavy-lanes.json` — the heavy lane's own `do-work/` exemption is a
  real gap, filed above as the impact-critical discovered task rather than fixed here.
- `_dev/tests/maintainer-verify.sh` — the stage wiring did not change, only what the manifest
  declares.
- `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, version files — release step, see the discovered
  task above.
- No fifth file was needed. The scope held at exactly four.
