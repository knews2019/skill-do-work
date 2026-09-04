---
id: REQ-537
title: 'Tier the maintainer gate: fast by default, --heavy on request'
status: completed
created_at: 2026-09-03T14:49:02Z
user_request: UR-104
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-537, REQ-538, REQ-539, REQ-540, REQ-541, REQ-542]
batch: two-tier-gate
write_set:
  - _dev/tests/maintainer-verify.sh
  - _dev/tests/contract-regressions.sh
  - _dev/tests/probe-batch.sh
claimed_at: 2026-09-03T15:02:38Z
completed_at: 2026-09-03T16:30:42Z
commit: 8d9d1bb
kb_status: promoted
kb_entry: REQ-537-tier-the-maintainer-gate-fast-by-default.md
---

# Tier the Maintainer Gate: Fast by Default, --heavy on Request

## What

`bash _dev/tests/maintainer-verify.sh` with no flags becomes the fast tier and the only canonical gate: toolchain floors, ShellCheck, gofmt, `go vet` for both modules, the fast aggregate, queue-kanban tests without JavaScript or browser probes, and do-work-cli tests. `--heavy` adds the board's JavaScript probes (strict marker on, zero-probe guard on), the browser lane when an engine is present, and the heavy aggregate probes. `--heavy-surfaces` prints the repo-relative path globs whose changes warrant a heavy run, so REQ-541 can consult it mechanically.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-shell-commands.md` (the listed prime), `_dev/primes/prime-releases.md`, `CLAUDE.md` and `crew-members/communication-style.md`. Approach: one tier variable read from a single argument parse, two branches at the board lane, one browser branch guarded the same way the Node branch already was, and a heavy-surface list held in the script with its board-test half derived rather than hand-listed.
- [x] **[APPLY]:** One file changed for the mechanism — `_dev/tests/maintainer-verify.sh`, +182/-117. The two other files in `write_set` (`_dev/tests/contract-regressions.sh`, `_dev/tests/probe-batch.sh`) were deliberately left untouched (D-02). Version and changelog files are the release half of the same commit.
- [x] **[UNIFY]:** `git show --stat 8d9d1bb` reads six files: the gate script, `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`. ShellCheck and the script's own `--self-test` both exit 0; no debug artifact, no leftover marker, no commented-out lane in the diff.

## Why

"then we could have a heavy weight test, but that should be executed only with user permission". Today's 301 s gate runs on every step; the maintainer decided the fast tier is the gate and the heavy tier is theirs alone (D1, D2).

## Context

- Today: lint and vet about 10 s, aggregate 112 s, queue-kanban 149 s, do-work-cli 35 s wall. The JavaScript probes are 48 s of the board run and the two re-exec meta-tests another 58 s (REQ-538 owns those).
- `_dev/tests/gate-runner.sh` and every `skills/` caller invoke the bare script and never pass `--heavy` (D2). The gate-evidence record is keyed by exact argv, so fast and heavy evidence are already distinct.
- The shimmed `--self-test` (`run_self_test`, `assert_success_stages`) enumerates stages per run; it must cover both tiers, and the strict-marker mutation check stays.

## Detailed Requirements

- Default run: floors, ShellCheck, gofmt, both `go vet`, fast aggregate, queue-kanban tests with JavaScript probes skipped through REQ-538's knob (until REQ-538 lands, pass the marker off and let probes run), do-work-cli tests.
- `--heavy`: everything the default runs plus JavaScript probes with `QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1`, the browser lane when present, and the heavy aggregate probes (REQ-539 classifies them; until then the aggregate stays whole).
- `--heavy-surfaces`: prints globs, one per line: the board's `web/**` and its JavaScript probe test files, installer and updater sources under `skills/do-work/tools/` and `tools/`, and anything else you classify heavy. Keep the list in the script, not in prose.
- Self-test: expected stage lists per tier; a `--heavy` fixture run; the marker-mutation check unchanged.
- No caller under `skills/` or the gate runner passes `--heavy`.
- Supersedes REQ-519's path-scoped `--changed` lane: the fast tier is the per-REQ check for every path, so no path scoping is built. If the fast tier cannot reach 120 s without it, record that in the commit body and capture path scoping separately; do not add it here silently.

## Constraints

- Land in place, not through `do-work run`; one integrating commit with version bump and changelog entry; prove it with one `bash _dev/tests/gate-runner.sh --once`.
- Delete before you add; every deleted test is listed in the commit body with the failure it pinned and why it no longer earns its cost. No new sentence pins, no new prose that walks a shell sequence.
- Never touch another session's claimed file under `do-work/working/`; stage explicit paths.

## Red-Green Proof
**RED prompt/case:** `bash _dev/tests/maintainer-verify.sh --heavy` and `bash _dev/tests/maintainer-verify.sh --heavy-surfaces`.
**Why RED now:** both exit 2 with the usage line; the default run executes the JavaScript probes and the full aggregate.
**GREEN when:** the default run exits 0 in under 120 s once REQ-538 to REQ-540 land (under 200 s on its own), `--heavy` exits 0 and executes the JavaScript probes, `--heavy-surfaces` prints globs, and `--self-test` exits 0.
**Validation:** User confirmed (D1, D2)

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes a shipped-style shell harness and its self-test.

## Full Context
See `do-work/user-requests/UR-104/input.md` for complete verbatim input.

---

## Landing Mode

Landed in place, not through `do-work run`, per this REQ's own Constraints: the pipeline runs the maintainer gate several times per REQ, and this REQ changes that gate. So there is no `## Triage`, no `## Plan` and no independent `## Review` section — the maintainer's own integrating commit and one `gate-runner.sh --once` are the whole proof this REQ was allowed.

## Implementation Summary

**Files changed** (all in commit `8d9d1bb`):
- `_dev/tests/maintainer-verify.sh` (modified, +182/-117) — the tier mechanism, `--heavy`, `--heavy-surfaces`, and the two-tier self-test.
- `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md` (modified) — 0.269.0 to 0.270.0.
- `CHANGELOG.md` and `skills/do-work/CHANGELOG.md` (modified) — the "0.270.0 — Two-Tier Maintainer Gate" entry, byte-identical mirror.
- `_dev/tests/contract-regressions.sh`, `_dev/tests/probe-batch.sh` — in `write_set`, deliberately unchanged (D-02).

**What was done:** a bare `bash _dev/tests/maintainer-verify.sh` is the fast tier and the only gate any caller runs. `--heavy` is the same run plus `QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1` on the board package and the strict browser-behavior lane. `--heavy-surfaces` prints the path globs whose changes warrant a heavy run, one per line, for REQ-541 to consult mechanically. The self-test gained per-tier expected stage lists, a `--heavy` fixture run and a no-Node `--heavy` fixture run; its three near-identical fixture blocks were folded into `run_success_fixture` and its failure loop into `assert_failure_stage`.

**What it does not do — read this before reading anything above as a speedup.** The two tiers are 2 s apart (629 s fast, 631 s heavy). The tier mechanism currently separates almost nothing, for two reasons that are both outside this REQ:

- The board's JavaScript probes are **not skippable yet**. REQ-538 (cutting the queue-kanban suite to its fast core) owns the marker knob that would make them `t.Skip`; until it lands, the default tier runs them too and only the strict zero-probe guard differs between the tiers.
- The do-work-cli suite is **REQ-540's** (cutting the do-work-cli tests down to transactions and recovery). Nothing here shortens it.

And the browser lane SKIPs in both tiers on this container, because no browser engine is installed (`google-chrome`, `google-chrome-stable`, `chromium`, `chromium-browser`, `chrome` are all absent from `PATH` and `QUEUE_KANBAN_BROWSER` is unset).

**The REQ's 120 s target was not met and no path scoping was added.** That was this REQ's own instruction, not an omission: "If the fast tier cannot reach 120 s without it, record that in the commit body and capture path scoping separately; do not add it here silently." It is recorded in `8d9d1bb`'s commit body. **REQ-537's value is therefore contingent on REQ-538 and REQ-540 landing.** What it delivers today is the seam those two REQs plug into, plus `--heavy-surfaces` for REQ-541 — not a faster gate.

## Decisions

- **D-01 — the heavy surface list is half static globs, half derived.** DECIDE & STATE. The static half is five directory-and-suffix globs, deliberately wider than the minimum: an unnecessary heavy run costs minutes, a missed one costs the coverage. The derived half is the nine board test files that carry the Node and browser probes, computed by grepping the four probe entry points those files call rather than hand-listed, because those nine filenames share no common token — per `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale, a hand list would silently under-report the first time a probe moved to a new file. **Risk and its answer:** the grep runs in a process substitution whose exit status nothing can read, so an absent or renamed entry point would read as "no heavy paths"; the derivation therefore fails loudly on zero matches instead of printing a short list.
- **D-02 — `_dev/tests/contract-regressions.sh` and `_dev/tests/probe-batch.sh` left untouched despite being in `write_set`.** DECIDE & STATE. The aggregate stays whole until REQ-539 (cutting the contract file and splitting the aggregate) classifies its probes, so a tier environment variable there today would be configuration with no consumer. **Delete before you add:** the honest move is to not write it.
- **D-03 — the no-Node fixture run became a no-Node `--heavy` run.** DECIDE & STATE. Node presence no longer selects a lane in the default tier, so the explicit-skip assertion that fixture carried has a subject only under `--heavy`. Coverage moved, it did not shrink; the assertion is still made, against the tier where it means something.
- **D-04 — the self-test's failure runs now clear `QUEUE_KANBAN_BROWSER`** the way its success runs already did. DECIDE & STATE. A maintainer with that variable exported would otherwise get a different failure-stage list from the same script, which is the kind of environment leak the self-test exists to catch.

## Discovered Tasks

This REQ's UR forbids new REQs for findings ("no new REQs for findings", `do-work/user-requests/UR-104/input.md`), so each item below is folded onto an existing destination per `actions/work-reference.md` → **Discovered Tasks Classification (Step 8)** and the Fold-First Rule, not minted as a REQ.

- **impact-rule-change — `--heavy` and `--heavy-surfaces` are documented nowhere a maintainer reads, except the changelog entry that introduced them.** Verified rather than assumed: `grep -rn "maintainer-verify" CLAUDE.md _dev/primes/*.md` is empty — `CLAUDE.md` does not mention the maintainer gate at all and has no `## Verify` section, so the gap is wider than "one stale section". The only prose naming the gate outside `do-work/` is `skills/do-work/tools/do-work-cli/prime-do-work-cli.md:79` ("run the unpiped `_dev/tests/maintainer-verify.sh` from the repository root"), which names the bare command and neither flag. Both maintainer-facing homes are outside this REQ's `write_set`. Prose-only with no root-cause match in the queue — REQ-541 owns the *pipeline* consulting `--heavy-surfaces`, which is a different root cause from documenting the flags for a human — so it lands on `do-work/prose-backlog.md` per the Fold-First Rule's prose-only destination.
- **impact-user-visible — the prescribed-shell interrupt case leaves its backend running, and asserts nothing about it.** Two `imagegen` fixture processes are alive in this container with `PPID 1`, for 2 h 05 m and 2 h 07 m, under `/tmp/prescribed-shell-scripts.*/image-interrupt-bin/imagegen`. Both **predate** REQ-543's teardown fix (`1cc3beb`, 14:43), so they are residue of the old leader-first teardown rather than evidence of a gap in the new one — stated that way because the two read very differently. What is durable is the fixture: `_dev/tests/prescribed-shell-cases/generate-report-image.sh:132-151` sends `kill -TERM` to the wrapper pid only, and then asserts the private staging file is gone, the old target is unchanged, and the status is non-zero — never that the backend stopped. So this whole class is invisible to the gate in either direction. Same root cause as REQ-546's F1 and F11 (a teardown that leaves a live descendant, and the test whose name promises that invariant not asserting it), whose `write_set` already carries the Go-side counterpart `internal/toolboxcommands/report_image_process_test.go` — appended there as an instance rather than captured separately. `prime-shell-commands.md`'s "Cleanup and interruption are separate trap contracts" is the rule the fixture is on the wrong side of.

**Not a discovered task, recorded as evidence in `## Testing` instead:** "nothing under `skills/` and nothing but `gate-runner.sh` passes `--heavy`" is a Detailed Requirement of this REQ that verified green, not a finding.

## Testing

**Gate evidence (the one run this REQ's Constraints require):**

```
gate 8d9d1bb53d2654f79fe89ab445f2af2eca738877 green 641s
```

`bash _dev/tests/gate-runner.sh --once` exit 0, recorded through `record-green-gate`. Re-checked while closing this record, not taken from the run report:

```
$ do-work-cli check-green-gate --at-revision 8d9d1bb -- bash _dev/tests/maintainer-verify.sh
gate evidence: state=exact_revision_match matches=true basis=exact_revision
  recorded=8d9d1bb53d2654f79fe89ab445f2af2eca738877 target=8d9d1bb53d2654f79fe89ab445f2af2eca738877
```

**Measured wall-clock, all three tiers exit 0:**

| Tier | Wall-clock | Exit |
|------|-----------|------|
| default (fast) | 629 s | 0 |
| `--heavy` | 631 s | 0 |
| `--self-test` | under 1 s (0.73 s re-measured at closure) | 0 |

**Where the 629 s goes.** The recorded gate run's own log (`$TMPDIR/do-work-gate-runs/8d9d1bb…log`, 641 s wall for the same fast tier) is the primary source, so these are read numbers, not estimates:

- `queue-kanban` — a **single package test at 467.4 s**, 73% of that run. REQ-538's territory.
- `do-work-cli` — **110.9 s** of in-test time summed across 26 packages (the heaviest: `finalization` 27.2 s, `publication` 21.6 s, `archivefetch` 15.2 s). REQ-540's territory.
- everything else — the aggregate contract suite plus the lint, `gofmt`, `go vet` and toolchain-floor lanes — about **63 s**.

The briefing this closure worked from cited 463.7 s and ~103 s for the same two lanes, from the separate standalone fast-tier run whose log is not retained; the figures above are the ones that can still be read off disk. Same conclusion either way: one package is roughly three quarters of the gate.

**`--heavy` demonstrably ran the JavaScript probes.** It ran with `QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1`, and the board package's `TestMain` (`skills/do-work-board/tools/queue-kanban/generate_test.go:196-205`) fails a green run that recorded zero probes. Green under the marker is therefore proof the probes executed, not just that they were selected. Node is present (`v22.22.2`), which is what makes that branch reachable here.

**No test or assertion was deleted.** Three near-identical fixture blocks were folded into `run_success_fixture` and the failure loop into `assert_failure_stage`; every assertion those blocks made is still made, with the tier arriving as a call argument. The strict-marker mutation check (`maintainer-verify.sh:422`) is byte-unchanged.

**Mutation evidence — three mutations against the new logic, each caught by `--self-test`:**

| Mutation | Result |
|----------|--------|
| default tier takes the heavy board branch | caught |
| `--heavy` selects the fast tier | caught |
| `--heavy-surfaces` prints nothing | caught |

**`--heavy-surfaces` output verified at closure** — 5 static globs plus 9 derived board test files, exit 0 in 12 ms.

**Detailed Requirement "no caller under `skills/` or the gate runner passes `--heavy`":** green. `grep -rn -- "--heavy"` across the tree returns only this REQ's own record, the four sibling queue REQs that reference the tier, `do-work/CHECKPOINT.md`, `UR-104`'s input, and the changelog. No `.sh` file passes it, `gate-runner.sh` included.

*Verified in place by the maintainer session; no independent review ran, per `## Landing Mode`.*

## Lessons Learned

**What worked:**
- Deriving the heavy board-test list instead of writing it down. The nine files share no filename token, so a hand list would have looked right and been wrong the first time a probe moved. Making the derivation fail loudly on zero matches is what keeps it honest, because the grep sits in a process substitution whose status nothing reads.
- Recording the miss in the commit body instead of reaching for path scoping. The REQ pre-authorized exactly one of those two responses, and the other would have shipped a `--changed` lane that REQ-538 through REQ-540 make unnecessary.

**What didn't:**
- The REQ's own sequencing. REQ-537 delivers the tier seam, but every second the seam could save belongs to REQ-538 and REQ-540, so it lands as a 2 s difference and a mechanism. Shipping the knob before the things it switches is a defensible order for review size, and it means this REQ cannot be judged on its own numbers — the changelog entry says so, and so does this record.

**Worth knowing:**
- **A green strict-marker run is the probe-execution proof; a green plain run is not.** The board package's `TestMain` refuses exit 0 when the marker is set and the probe counter is zero. That is why `--heavy` can assert "the probes ran" while the default tier can only assert "the package passed".
- One package can be three quarters of a test gate. 467 s of a 641 s run came from `queue-kanban` alone, which is the whole reason the two-tier split cannot pay off before REQ-538.
- The gate runner's own log under `$TMPDIR/do-work-gate-runs/<revision>.log` retains per-package times long after the run. It is the cheapest place to re-derive a gate breakdown without paying for the gate again.

## Orientation

`bash _dev/tests/maintainer-verify.sh` with no arguments is now the one canonical gate, and it is the fast tier. `--heavy` adds the board's strict JavaScript marker and the browser lane; `--heavy-surfaces` prints the globs that warrant a heavy run. Nothing under `skills/` and not the gate runner passes `--heavy` — it is the maintainer's flag, by hand.

Today the two tiers differ by 2 s, so treat this as a seam and not a speedup until REQ-538 (queue-kanban's marker knob) and REQ-540 (the do-work-cli suite) land.

`prime_files`: `_dev/primes/prime-shell-commands.md` — its § Closed Enumerations Go Stale is the rule D-01 applies, and its trap list still matches the script's shape.

**[MAP CHANGED]** — the gate has a tier variable where it previously had one unconditional path, and the heavy-surface list is a machine-readable output of the script rather than prose anywhere. Why it matters: every future decision about what belongs in the fast gate now has one place to be expressed, and REQ-541 can ask "is this diff heavy?" without a document to keep in sync.
