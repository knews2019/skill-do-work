---
id: UR-058
title: Render every verify finding on the HTML board
created_at: 2026-08-19T13:47:06Z
requests: [REQ-284, REQ-285]
word_count: 1538
---

# Render Every Verify Finding on the HTML Board

## Summary

An upstream suggestion authored against v0.212.25 in a consumer repo: `queue-kanban verify` finds
sixteen categories of queue and process breakage, the HTML board renders three, and the board is the
surface a human actually looks at. The document asks for the whole finding set to reach the board,
names five implementation hazards (C1-C5), proposes a code shape, lists tests, and offers a lesser
alternative.

The document was triaged with `do-work validate-feedback` before capture. Every line number in it
resolves correctly against the tree at 0.212.25, and two load-bearing claims were reproduced
independently (the mtime cache's behaviour, and the probe cost).

## Extracted Requests

| Finding | Verdict | Captured as |
|---|---|---|
| F1 — render all sixteen categories on the board | Accept | REQ-284 (producer) + REQ-285 (client) |
| F2 — C1, the mtime cache cannot see time or git state | Diagnosis accepted, 5s TTL rejected | REQ-284 — serve computes findings fresh per request, no TTL. Maintainer confirmed: the dataset is small and the tool runs on fast machines |
| F3 — C2, probe cost is ~30ms | Accept (measured independently: 0.13s summary vs 0.17s verify on this repo) | No work — evidence only |
| F4 — C3, three categories would render three times | Accept | REQ-284 (suppression in the Go producer) |
| F5 — C4, `generate` publishes a shareable snapshot | Accept, remedy replaced | REQ-284 — no `scope` enum; no absolute path in the emitted JSON at all, pinned by a `generate_test.go` assertion |
| F6 — C5, permanent unfixable version-probe skip | Accept, remedy amended | **Already queued as REQ-282** (UR-057, captured minutes earlier from a different upstream document by the same consumer) — not re-captured |
| F7 — suggested code shape | Accept, minus the `scope` parameter | REQ-284, REQ-285 |
| F8 — tests worth adding | Accept, two cases amended | REQ-284, REQ-285, REQ-286 |
| F9 — lesser alternative, a `kanban-verify` just recipe | Push back | Not captured. Superseded by queued REQ-283 (UR-057), which routes `verify` through `do-work-board` with no Justfile surface |

## Batch Constraints

- `verify` keeps its non-zero exit code. It is the CI/hook contract and a strip in a browser cannot
  replace it. `main.go` is unchanged by this batch.
- Neither surface starts repairing anything. `do-work cleanup` still owns fixes and still asks first.
- The board tool's write-surface count stays at three. Nothing in this batch writes pipeline state.
- Precedent: REQ-214 made `verify` surface the board's completion anomalies. This batch is the same
  argument in the other direction, and the same never-silent principle governs it — a skipped probe
  must still render, because silence reads as "checked and clean".

## Triage corrections carried forward

Three premise errors in the document, verified against the code, that the builders must not inherit:

1. "The board never runs a probe that touches git" is false. `buildBoard` calls `lookupGitCommitDate`,
   which shells out to `git log -1 --format=%cI` once per terminal ticket carrying a commit hash
   (`model.go:1310`). The board already pays a per-completed-REQ git cost.
2. The gap table's "future-dated timestamps — verify: no" is wrong. `appendClaimFindings` reports a
   future-dated `claimed_at` under `claim-needs-attention`. What verify lacks is the general
   `FutureTimestampFields` sweep across all fields, not the claimed-at case.
3. C5 understates the defect. The skip fires in the suite repo too, not only in consumer installs —
   there is no `actions/` at this repo's root since the four-skill split, and `verify` exposes no
   `--version-file` override. All three release categories are dead by default everywhere.

One path slip: the installer is `tools/install-do-work-suite.sh`, not `do-work/tools/…`.

## Decisions taken with the maintainer

- **F2 — no cache.** Findings are computed on every `/board-data.js` request. The maintainer's reason: the
  dataset is small and this tool runs on fast machines. The mtime cache on the board build stays.
- **F5 — nothing that leaves the machine carries an absolute path.** Stated by the maintainer as the
  general rule: a path that resolves here means something else, or nothing, elsewhere. It replaces the
  upstream `scope` enum, applies to both callers, and is pinned by a test rather than by prose.
- **F9 — no Justfile recipe.** REQ-283 already delivers the reachability half through the skill that owns
  the tool. Nothing was queued for F9, so there was nothing to cancel.

## Duplicate found at capture

F6 is already in the queue. UR-057 — a separate four-defect report from the same consumer repo, captured
into this queue at 13:42:45Z, minutes before this one — carries it as D3 → REQ-282, with a wider scope
(it also found that the probes are dead in the maintainer repo) and an explicit decline of suite-root
discovery. Nothing was captured for F6 here.

Three UR-057 REQs write `skills/do-work-board/tools/queue-kanban/verify.go`: REQ-280, REQ-281, and
REQ-282. REQ-284 writes it too, in a disjoint region. The overlap is declared in REQ-284's write set and
Dependencies section; no ordering dependency was invented.

## Full Verbatim Input

do-work validate-feedback: # Upstream suggestion for `knews2019/skill-do-work` — render every `verify` finding on the HTML board

**How to use this file:** paste everything below the horizontal rule into a Claude Code session opened in a
clone of `knews2019/skill-do-work`. It is written to be actionable cold, with no reference back to the
consumer repo it was authored in. Observed against **v0.212.25**; all line numbers are from that tag.

---

## Request

`queue-kanban verify` finds sixteen categories of queue and process breakage. The HTML board renders three
of them. The board is the surface a human actually looks at; `verify` is a command nobody runs unless an
agent runs it for them. Please move the whole finding set onto the board, so that seeing the board means
seeing the problems.

Two things stay as they are: `verify` keeps its non-zero exit code (it is the CI/hook contract, and a strip
in a browser cannot replace it), and neither surface starts repairing anything — `do-work cleanup` still owns
fixes and still asks first.

## What each surface shows today

Measured on a consumer repo with 34 pending REQs, 2 claimed, 2 completion anomalies.

| Finding category | `verify` | Board |
| --- | --- | --- |
| `completion-anomaly` | yes | yes — data-warnings banner **and** the anomalies strip |
| `duplicate-req-id` | yes | yes — data-warnings banner |
| `stray-req-file` | yes | yes — data-warnings banner |
| `claim-needs-attention` (3h) | yes | **no** |
| `checkpoint-names-missing-req` | yes | **no** |
| `stranded-finished-req` | yes | **no** |
| `assigned-elsewhere-claimed-here` | yes | **no** |
| `ur-archived-with-live-member` | yes | **no** |
| `merged-worktree-leftover` | yes | **no** |
| `unmerged-worktree-leftover` | yes | **no** |
| `worktree-merge-state-undetermined` | yes | **no** |
| `worktree-wrote-queue-state` | yes | **no** |
| `worktree-committed-queue-state` | yes | **no** |
| version-vs-CHANGELOG agreement (3 categories) | yes | **no** |
| future-dated timestamps | **no** | yes — data-warnings banner |
| dangling deps, column, testing, schema | **no** | yes — data-warnings banner |

Neither is a superset. `verify` forwards only the duplicate-id warnings out of `board.Warnings`
(`verify.go` → `appendDuplicateRequestIdFindings`), and the board never runs a probe that touches git.

The gap that prompted this: two REQs sat claimed for 13h against a 3h threshold. The board printed
`claimed Aug 19, 00:28 UTC · 13h 01m` on both cards and flagged neither. The age was on screen; the rule
was not applied to it. Nobody ran `verify`, so nobody knew.

## Where the code is

- `verify.go:105` `runVerifyProbes(repoRootOverride, now)` — builds a board, runs ten `append*Findings`
  probes, returns `VerifyReport{Findings, SkippedProbes}`. Read-only by contract.
- `main.go:340` `runVerifyCommand` — the only caller. Prints and exits `report.ExitCode()`.
- `model.go:403` — `detectCompletionAnomaly` runs inside `buildBoard`, which is why anomalies are the one
  category that already reaches every view.
- `generate.go:69` `generatedBoardData` — the JSON the client renders; `Warnings []string` at line 80.
- `web/board-cards.js:387-406` — the data-warnings banner, a flat `<li>` per string.
- `web/board-cards.js:412-431` and `web/template.html:137-146` — the completion-anomalies strip: a titled
  section with a count, a hint line, and per-item cards. This is the shape to copy.
- `serve.go:33-44` and `:304-310` — the mtime cache: each request stats the `do-work` tree and rebuilds only
  when a file's mtime moved.

## Five things that will bite, found while reading the code

### C1 — the mtime cache cannot see two of the inputs (the real design problem)

`serve.go`'s cache fingerprints the `do-work` tree. Two probe families depend on state outside it:

- `claim-needs-attention` is **time**-dependent. A claim crosses 3h with no file changing. Under the current
  cache the board would never rebuild, so the finding that motivated this whole request is exactly the one a
  naive implementation fails to show.
- The five worktree probes and the release probes read git and files outside `do-work/`.

So findings must not be computed inside the mtime-gated path. Compute them per `/board-data.js` request with
their own short TTL (5s is plenty — the page is manually reloaded), or extend the fingerprint with a coarse
time bucket plus a git-state stamp. The TTL is simpler and cannot go stale in a way that hides a finding.

### C2 — cost is not the objection it looks like

Measured on the consumer repo, best of three, warm:

| Command | Wall |
| --- | --- |
| `queue-kanban summary` (board build only) | 0.08s |
| `queue-kanban verify` (board build + all probes) | 0.11s |

~30ms for the probe set. The git surface is 5 short subprocesses baseline (`rev-parse --git-dir`,
`worktree list --porcelain`, `branch --list`, two `rev-parse`), plus ~3 per `worktree-agent-*` branch —
normally zero. A 5s TTL keeps even a reload-heavy session under one probe run per reload.

### C3 — three categories would render three times

`completion-anomaly`, `duplicate-req-id` and `stray-req-file` already reach the board as warnings, and
anomalies additionally get their own strip and a per-card `ANOMALY` badge. Appending the `verify` findings to
`Warnings` verbatim would print the same anomaly prose a third time, in near-identical words. Pick one owner
per category: the board keeps rendering the three it already has, and the new strip renders the other
thirteen. That suppression belongs in the Go producer, not in JS — the client should get a list it can render blindly.

### C4 — `generate` publishes a shareable snapshot

`kanban-static` writes `build/queue-kanban-board/index.html` for sharing. Worktree findings carry local
absolute paths, branch names, and `git status` output; `assigned-elsewhere-claimed-here` carries agent
names. Either omit the git-derived categories from `generate` output, or emit them with paths reduced to
basenames. The queue-derived findings (stale claims, stranded REQs, ghost checkpoints, archived-UR members)
are safe in a snapshot and are the ones a shared board most wants.

### C5 — a consumer repo gets one permanent unfixable skip

`release.go:26` — `defaultVersionFileRelativePath = "actions/version.md"` — is resolved against the
**consumer's** repo root, where it can never exist. Every consumer's `verify` run ends with:

```
- skipped version-vs-changelog probes: no version file readable at <consumer>/actions/version.md
```

Today an agent reads that once and ignores it. On the board it becomes a permanent line of noise on every
project that installs the suite, which is the fastest way to teach people to stop reading the strip. Skipped
probes must still render — `verify.go`'s own comment on the integration-ref skip is right that silence reads
as "checked and clean" — so the fix is to detect that the probe is inapplicable rather than failed: when
`actions/version.md` is absent **and** the repo root is not the suite repo, drop the probe from the report
instead of listing it as skipped. Worth fixing before, or with, this change.

## Suggested shape

1. **`verify.go`** — split `runVerifyProbes` into a board-taking core, e.g.
   `collectVerifyFindings(repoRoot, board, now, scope)` plus the current thin wrapper that builds the board
   first. `generate`/`serve` already hold a `*Board`; without this the board gets built twice per request.
   `scope` selects `all` (CLI, serve) vs `queueOnly` (static generate) for C4.
2. **`generate.go`** — add to `generatedBoardData`:
   ```go
   VerifyFindings []generatedVerifyFinding `json:"verifyFindings,omitempty"`
   VerifySkipped  []string                 `json:"verifySkipped,omitempty"`
   ```
   with `{category, detail, remedy, fixable}` per finding, category-suppression for C3 applied here.
3. **`web/template.html` + `board-cards.js` + `board.css`** — a second strip modeled on
   `board-anomalies` (`template.html:137`): title, count, one card per finding, category as the badge, the
   remedy as the muted second line, skipped probes in a collapsed footer. Reuse the `.board-anomalies-*`
   rules rather than writing a parallel palette; the anomaly strip's visual weight is already calibrated.
4. **`serve.go`** — compute findings on the `/board-data.js` path with a 5s TTL, outside `cachedBoardData`
   (C1). Keep them out of the markdown/copy payload unless you want them in a paste.
5. **`main.go`** — unchanged. `verify` stays the exit-code surface.

## Tests worth adding

Mirroring the existing files:

- `verify_test.go` — a stale claim crossing the threshold with **no file mtime change** produces a finding on
  a second collection with an advanced `now`. This is C1's regression test and the one that matters.
- `generate_test.go` — the three board-owned categories do not appear in `verifyFindings`; the other ten do.
- `generate_test.go` — `queueOnly` scope emits no absolute filesystem path anywhere in the JSON (C4).
- `board_live_test.go` — `/board-data.js` carries findings and the TTL serves a second request from cache.
- `release_test.go` — a repo without `actions/version.md` that is not the suite repo reports neither a
  finding nor a skip (C5).

## Lesser alternative, if the board change is not wanted

Add a `kanban-verify` recipe to `do-work-board/justfile.template` (the block already ships `run-kanban`,
`run-kanban-cli`, `kanban-static`, `kanban-summary`, `run-do-work-update`), plus the name in the
post-install assertion loop at `do-work/tools/install-do-work-suite.sh:287`. Collision reservation needs no
change — `replace-text-section.sh:325` derives reserved names from the template itself.

```just
# Read-only queue + git-state check — stale claims, worktree leftovers, broken completion stamps (exits non-zero on findings)
kanban-verify:
    cd "{{justfile_directory()}}/.claude/skills/do-work-board/tools/queue-kanban" && go build -o queue-kanban . && ./queue-kanban verify --repo-root "{{justfile_directory()}}"
```

This makes the check reachable by hand but does not solve the actual problem: it is still a thing someone has
to decide to run, and the reason the 13h claims went unnoticed is that nobody decided to.

## Version and changelog

Left to your release ritual — this file proposes no version number.

---
*Captured: 2026-08-19T13:47:06Z*
