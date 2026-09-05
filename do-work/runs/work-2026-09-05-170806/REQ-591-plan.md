# Implementation plan — REQ-591: reduce repeated setup and unaffected reruns in the fast gate

Base revision at planning time: `e0bdf8bf92d7c2b9663cff64a3d85a5220216991`. Machine: 8 logical
CPUs, `go1.26.1 darwin/arm64`, ShellCheck 0.11.0. Load average during planning 4.99 — a loaded
window, which is why the measurement protocol in §5 is written the way it is.

Every claim below was checked against source, not taken from the exploration report. Where the
exploration's reasoning was wrong or incomplete, §6 says so.

---

## 1. Scope judgment (answer first)

**Both parts belong in this REQ, sequenced. Task count is 5.**

Part two is the whole of acceptance criterion 2 and most of criteria 3 and 4. Splitting it out
would leave REQ-591 unable to meet its own acceptance, which is a capture-level change (rewriting
the criteria) rather than a builder decision. The REQ's Builder Guidance already asks for one
coherent request with small verified increments, and the five tasks below are exactly that.

Five tasks is the signal line, not past it. If the orchestrator wants smaller increments, the
clean seam is **inside T2**, between the evidence engine (fingerprint + decision, pure Go) and the
three CLI commands that expose it. The seam between part one and part two is T1/T2 — a split there
hands a follow-up REQ acceptance criterion 2 and the second half of criteria 3 and 4.

---

## 2. Files to modify, in order

### T1 — `_dev/tests/session-start-hook-behavior.sh` (part one, the measurable win)

**What changes.** Replace the nine per-case `install_core_fixture` calls with **one shared skill
root** built once, plus a per-case `project/` tree and a per-case rewrite of
`<shared-root>/actions/version.md`.

Concretely:

- New `install_shared_skill_root` runs once before any case and builds
  `$fixture_root/shared-skill/{hooks/session-start.sh, tools/do-work-cli.sh, tools/do-work-cli/, actions/}`,
  with `rm -f tools/do-work-cli/do-work-cli` exactly as today.
- `run_banner_case`, `run_tail_case`, the authority matrix and `run_launcher_failure_case` each
  keep their own `$case_root/project/` (git repo, queue files, journal, archive, reservations) and
  their own `$fake_bin` where they have one, and point `CLAUDE_PROJECT_DIR` at that project while
  running `$shared_skill_root/hooks/session-start.sh`.
- Every case writes `$shared_skill_root/actions/version.md` before it runs. `version.md` becomes a
  **required** parameter of each helper — no case can run without setting it, so a later case can
  never silently inherit the previous case's banner input. The cases are strictly serial today and
  stay serial; nothing in this file backgrounds anything.
- `missing_root` keeps its own root: it must *lack* `tools/do-work-cli.sh`. Cost is one `cp`.
- The static `grep` check at the end is untouched — it reads the real repository file.

**Why this shape and not the symlink shape.** The exploration offered a second option: per-case
skill roots whose `tools/do-work-cli` is a symlink into one shared module, justified by
`do-work-cli.sh` resolving with `pwd -P`. **That justification is wrong.** `do-work-cli.sh:7`
applies `pwd -P` to the directory holding the launcher (`.../skill/tools`), then *appends*
`/do-work-cli` (`:8`); the appended component is never canonicalized by the shell. Whether the
symlink collapses to one Go build-cache action depends on `go tool -C` chdir'ing and Go's
`os.Getwd()` falling back to the physical path — plausible, unverified, and a silent regression to
1-link-per-root if it does not hold. The shared-root shape does not depend on it at all: nine
consumers reach one physical module directory by construction.

**Expected effect.** Seven fresh `go tool -n` link actions (~1.1–1.5 s wall, ~4.5 s CPU each)
collapse to one. Standalone probe 12.90 s → ~3.5 s; aggregate contract lane ~24 s → ~19 s.

**What proves coverage did not weaken.** Per scenario, the assertion that still fails if the
boundary breaks:

| Scenario | Boundary | Assertion that catches a break |
|---|---|---|
| 4 banner cases (`:36–40`) | real `session-start.sh` → `do-work-cli.sh` → `go tool` → binary; version parsing; pending count | exact stdout byte compare + `status -ne 0` + `[ -s stderr ]`. A broken launcher path exits 2 with stderr; a leaked `version.md` makes `missing`/`reformatted` print `v9.8.7` instead of `vunknown`. |
| 2 tail cases (`:77–78`) | unfinished-finalization tail; journal/archive file not mutated | exact two-line stdout compare + `shasum -a 256` before/after. A shared `project/` would let case 2 see case 1's journal and the exact compare fails. |
| authority matrix (`:82–105`) | installed layout via `--skill-root`; committed → typed `deleted`; Git-unavailable → `RESERVATION-GIT-AUTHORITY-UNAVAILABLE`, state preserved | reservation file absent/present checks + `grep -q '"kind": "deleted"'` + `grep -q '"protocol_output": "do-work v9.8.7 loaded.'` + the negative `grep -q 'REQ-000204.*deleted'`. Also `go tool -C "$shared/tools/do-work-cli" -n do-work-cli` must still succeed, which fails loudly if the shared tree is not a real module. |
| 2 launcher-failure cases (`:123–124`) | fake `go` on PATH; too-old Go and build failure | `status -ne 2`, `[ -s stdout ]`, `grep -q "$expected" stderr`. `$fake_bin` stays per case; a leaked PATH would let the real toolchain build and the status would be 0. |
| missing-launcher (`:126–134`) | absent canonical launcher | `missing_status -eq 0` fails the case, plus `grep -q 'canonical launcher is missing'`. Own root, so a shared root would make it exit 0. |
| static grep (`:136–139`) | hook stays a thin launcher | unchanged, reads the real repo file. |

**One new assertion, added deliberately** (this strengthens coverage, it does not replace any):
after every case, assert the shared module directory contains no files the fixture did not put
there — `[ -z "$(ls -A "$shared_skill_root/tools/do-work-cli/do-work-cli" 2>/dev/null)" ]` and a
`find`-based compare of the tree against its post-install state. This is the invariant that makes
sharing safe. `_dev/tests/do-work-cli-launcher-behavior.sh:118–122` already pins the same property
for the launcher in isolation; this pins it for the shared fixture, so a future case that writes
into the shared tree fails here rather than corrupting a later case silently.

**Other probes and owner contracts — what is worth doing and what is not.**

Do, in this REQ: nothing beyond T1. Rationale per candidate:

- `_dev/tests/contracts/core-checks.sh` (5.25 s, six inline `git init` repositories) — real, but it
  runs *serially* before the probe batch, so its saving is real gate wall time. Deferred anyway:
  the fixture shapes there are heterogeneous (different commits per probe), so a template needs
  six templates, and the win is ~2 s against T1's ~9.4 s. Name it as a follow-up.
- `internal/lifecycleadvance` (13 `git init` sites, 41 s of file time) and
  `internal/heavyverification` (three repository builders, ~49 s) are the two densest REQ-574-shaped
  targets. **Not in this REQ.** They are pure REQ-574 repetition, the REQ says build on that work
  rather than repeat it, and part two removes the whole `do-work-cli` stage from a matching-input
  run — which subsumes their contribution for the case this REQ is optimizing. Doing both would
  make the before/after comparison unable to separate setup savings from avoided execution, which
  acceptance criterion 3 explicitly requires.
- `_dev/tests/update-script-behavior.sh` and `install-suite-behavior.sh` carry the same
  fresh-path link penalty as T1 but are **heavy-tier only** — out of the fast gate's scope.
- Moving the four serial owner contracts into the parallel probe batch (~10 s) is deliberately
  **not done**: acceptance criterion 3 rules out greater concurrency alone, and it would worsen
  contention in the lane T1 just improved.

### T2 — Fast-stage evidence engine and command surface (Go)

Files, in order:

1. **`skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go`** (new,
   written first — `tdd: true`). Table-driven positive + counter-cases; see §4.
2. **`skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go`** (new).
   - `fastStageManifest` / `manifestFastStage` decoded by the **same** `decodeManifest` shape
     (`DisallowUnknownFields`, trailing-value refusal, schema version, duplicate-id refusal,
     per-rule `coverageRule.validate()`), reusing `coverageRule`, `laneCoversPath` and
     `laneFingerprintInputs.validate()` unchanged.
   - `workingTreeFileSeals(repositoryRoot, coverage, manifest)` — generalizes the existing
     `fingerprintUntrackedFiles` (`heavy_evidence.go:469–504`) to **all** files, tracked and
     untracked, sealed from **working-tree bytes**:
     - `git ls-files -z --cached` for tracked paths; `git ls-files -z --others --exclude-standard`
       and the same with `--ignored` for untracked, exactly as today.
     - a path is sealed when the stage covers it **or** the manifest classifies it nowhere
       (the `heavy_evidence.go:301–303` rule, kept verbatim in intent).
     - `do-work/` (`queueStatePrefix`) is skipped, as today.
     - `os.Lstat` failure (a tracked path deleted from the worktree), a non-regular file (symlink,
       socket, directory-in-place-of-file), or a read error is an **error**, never a weaker seal.
       An error becomes `fingerprint_uncertain` and the stage executes.
   - `fastStageFingerprint(repositoryRoot, stage, manifest)` — digests a `fastStageFingerprintDocument`
     struct (byte-stable, not a map) over: schema version, stage id, exact argv, the sorted seals
     above, every declared toolchain probe's output through the existing
     `runFingerprintProbe(..., 5*time.Second)` (`heavy_fingerprint_probe.go`), and the **whole**
     `os.Environ()` sealed name-by-name as `laneFingerprint:333–356` does. A stage with no
     `fingerprint` block is never reusable. A stage covering no path at all is an error.
   - `fastStageEvidenceStore` — same `laneEvidenceStore` mechanics (git **common** dir, 0700
     directory, 0600 records, `atomicfile`, identity/size/mtime re-check on read) under a
     **distinct directory name** `do-work-fast-stages`, with a distinct `fastStageSchemaVersion`
     and a `StageID` field. This separation is load-bearing: a fast-gate green computed on a dirty
     tree must never be readable as heavy-lane evidence, which is attributable to HEAD.
   - `decideFastStageReuse(...)` — the same eight-way disposition/reason ladder as
     `decideLaneReuse` (`heavy_evidence.go:375–409`), reusing the existing reason constants and the
     4-hour `laneEvidenceMaximumAge`, checked independently of the fingerprint.
3. **`skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_commands.go`** (new) and
   its registration in the CLI's command table (same file the heavy commands are registered in).
   Three commands, each with one job:
   - `decide-fast-stage --manifest <path> --stage <id>` — pure, no writes. Default output is one
     line: `<disposition> <reason> <fingerprint-or-dash> <recorded-at-or-dash> <evidence-revision-or-dash>`.
     `--format json` carries the same fields as a typed `resultmodel.CommandResult`.
   - `invalidate-fast-stage --stage <id>` — deletes any prior record; refuses (non-zero) if it
     cannot prove the deletion, mirroring `InvalidateLaneEvidence`'s contract
     (`heavy_run.go:149–151`).
   - `record-fast-stage --manifest <path> --stage <id> --fingerprint <sha256> --stage-exit-status 0`
     — **recomputes** the fingerprint and refuses if it differs from the supplied one, then writes
     the record stamped with the instant the stage ran. A non-zero `--stage-exit-status` is refused
     with a typed code and writes nothing, mirroring `record-green-gate`'s `GATE-EVIDENCE-NOT-GREEN`
     (`gate_commands.go:39–54`).
4. **`skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_commands_test.go`** (new) —
   argument parsing, refusals, and the recompute-mismatch refusal.

### T3 — `_dev/tests/fast-stages.json` (new manifest)

Same schema family as `_dev/tests/heavy-lanes.json`. Two stages:

```
do-work-cli-fast-tests
  argv: ["env","GIT_CONFIG_NOSYSTEM=1","GIT_CONFIG_GLOBAL=/dev/null",
         "bash","_dev/tests/run-go-tests-with-budget.sh",
         "skills/do-work/tools/do-work-cli","-short","./..."]
  coverage: subtree skills/do-work/tools/do-work-cli
            subtree _dev/tests
            + the verified maintainer-tree paths its tests read (see below)
  fingerprint.toolchain_probes:
      [["env","GIT_CONFIG_NOSYSTEM=1","GIT_CONFIG_GLOBAL=/dev/null","python3",
        "_dev/tests/heavy-runtime-fingerprint.py",
        "skills/do-work/tools/do-work-cli","bash","git","python3"]]

queue-kanban-fast-tests
  argv: ["env","QUEUE_KANBAN_JAVASCRIPT_PROBES=off","QUEUE_KANBAN_BROWSER_PROBES=off",
         "DO_WORK_GO_TEST_EXCLUDE_PREFIXES=TestJavaScriptBehavior,TestBrowserBehavior",
         "bash","_dev/tests/run-go-tests-with-budget.sh",
         "skills/do-work-board/tools/queue-kanban","./..."]
  coverage: subtree skills/do-work-board/tools/queue-kanban
            subtree _dev/tests
            + verified extras
  fingerprint.toolchain_probes: same helper, module = queue-kanban, tools bash git python3

non_stage_coverage: [ {"kind":"subtree","path":"do-work"} ]
```

Two things this shape buys:

- The `QUEUE_KANBAN_*` selectors move from an inline `env`-prefixed subshell in
  `maintainer-verify.sh:606–609` into the **sealed argv**. Today those selectors decide which tests
  run at all and would not appear in a fingerprint computed by a separate decision process. Putting
  them in argv removes that divergence and makes the executed command and the sealed command
  provably the same string.
- `non_stage_coverage` declares **only** `do-work/`. Everything else in the repository that no
  stage covers is unclassified, so it is sealed into **every** stage and any change to it forces
  both stages to execute. That is the conservative default requirement 2 demands; narrowing it
  later needs evidence, not convenience.

**The coverage extras are a real verification step, not a guess.** Sixteen `*_test.go` files across
the two modules reference repository-root-relative paths. Most are string literals written into
temp-repo fixtures, but at least four are genuine maintainer-tree reads —
`heavy_maintainer_tree_test.go` (`_dev/tests/heavy-lanes.json`, `heavy-runtime-fingerprint.py`,
`maintainer-verify.sh`, `run-go-tests-with-budget.sh`, `test-duration-log.sh`, `fixture-repo.sh`,
`skills/do-work/actions/version.md`, `skills/do-work/hooks/hooks.json`,
`skills/do-work/tools/do-work-cli.sh`, `skills/do-work/agent-instructions.template.md`,
`skills/do-work-board/justfile.template`, `README.md`),
`resultmodel/result_model_test.go` (`_dev/tests/heavy-lanes.json`),
`cmd/do-work-cli/gate_evidence_integration_test.go`, `defer_gate_test.go`,
`heavy_verification_test.go` and `lifecycle_timing_test.go` (`_dev/tests/maintainer-verify.sh`).
T3 enumerates these by grep, declares each in the owning stage's coverage, and T4 adds the test
that keeps the declaration honest.

### T4 — `_dev/tests/maintainer-verify.sh` wiring, self-test guard, manifest pinning

1. New helper `run_stage_with_evidence <stage-id> -- <argv...>` placed beside
   `run_budgeted_go_tests`. Behavior, in order:
   - If `MAINTAINER_VERIFY_SELFTEST_LOG` is set → run the argv directly and return. Same bypass
     seam `run_budgeted_go_tests:74–80` already uses.
   - Call `decide-fast-stage` through `bash "$repo_root/skills/do-work/tools/do-work-cli.sh"`, from
     `$repo_root`, with `env -u DO_WORK_TEST_RUN_ID -u DO_WORK_TEST_OTHER_GATE_PROCESSES
     -u DO_WORK_TEST_DURATION_LOG` (see D-07). Capture **status and output separately**
     (`local decision_line; decision_line="$(...)" || decision_status=$?`) — a decision command that
     could not run must read as "execute", never as "nothing changed"
     (`prime-shell-commands.md` § Unchecked Exit Status Reads as Content).
   - Any non-zero status, an unparseable line, or a disposition that is not exactly `reused` →
     **execute**, with reason `decision_unavailable` where the command itself failed.
   - `reused` → print
     `maintainer-verify: stage <id>: REUSED (fingerprint_match, recorded <ts>, revision <rev>)`
     and return 0.
   - Otherwise print `maintainer-verify: stage <id>: EXECUTING (<reason>)`, call
     `invalidate-fast-stage` (a failure here fails the gate — an unrevocable prior green is the
     false-green shape requirement 2 forbids), run the argv capturing its exact status, and on
     status 0 only, call `record-fast-stage`. Return the stage's **exact** status.
2. Replace the two inline Go-test stages (`:604–610` board default tier, `:659` CLI `-short`) with
   `run_stage_with_evidence` calls carrying the manifest argv. The heavy tier and the
   `--heavy-lane` paths are untouched: they keep `run_budgeted_go_tests` and the existing
   heavy-lane evidence.
3. Add a closing summary line before `Maintainer verification passed.`:
   `maintainer-verify: gate wall <n>s stages executed=<n> reused=<n>` from a `date +%s` pair. One
   line, no new stage, and it is what makes a selective run reviewable (requirement 3) and gives
   §5's protocol an in-band cross-check on the external timing.
4. `--self-test` stays at exactly nine stages, each exactly once. The bypass in step 1 is what
   keeps that assertion literal rather than loosened —
   `_dev/primes/lessons-shell-commands.md` (REQ-187) says that count must be changed deliberately,
   and here it is deliberately **not** changed: `--self-test` proves the gate's stage list, and the
   reuse decision is proved by its own probe (T5) where it can be tested far more thoroughly.
   Add one self-test assertion that the bypass is active — that the nine-stage run recorded no
   `decide-fast-stage` invocation — so a future change cannot reach nine stages *through* reuse.
5. **`skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_maintainer_tree_test.go`** —
   add fast-stage assertions to this file, **not a new one**. It is already the single
   export-ignored maintainer-tree reader (`.gitattributes:62`), and the
   `shipped-module-test-self-containment` lesson (0.294.2) says every maintainer-tree read sits in
   one file. New assertions:
   - every `_dev/tests/fast-stages.json` stage argv appears literally in
     `_dev/tests/maintainer-verify.sh`;
   - every out-of-module repository path literal in each module's `*_test.go` sources is matched by
     that stage's coverage rules — the drift guard for T3's enumeration;
   - `non_stage_coverage` declares `do-work/` and nothing else;
   - every stage declares a `fingerprint` block (a stage without one can never reuse, so a silent
     omission is a silent disabling).

### T5 — `_dev/tests/fast-stage-reuse-behavior.sh` (new probe) + `_dev/tests/contracts/probe-lanes.sh`

End-to-end shell probe; details in §4. Registered with `register_probe` — note it must be
`chmod +x` or `register_probe` counts a failure on the spot. Must finish under 30 s
(`probe-batch.sh:38–43`). `_dev/tests/contract-regressions.sh` is **not** touched: its 77-line
ratchet ceiling is self-checked at `:17–22`, and registration happens in `probe-lanes.sh`.

---

## 3. Architectural decisions

**D-01 — Share one skill root; do not use the symlink shape. DECIDE & STATE.**
The symlink shape's win rests on Go canonicalizing a path component the shell never canonicalized;
the shared-root shape reaches one physical module by construction and needs no such assumption.

**D-02 — `version.md` becomes a required parameter of every case helper. DECIDE & STATE.**
A shared root makes "forgot to set the banner input" an ordering bug that reads as a pass; making
it a required argument makes it a syntax error instead.

**D-03 — Extend `internal/heavyverification` rather than add a package. DECIDE & STATE.**
`coverageRule`, `decodeManifest`, the evidence store mechanics, the disposition/reason constants,
the 4-hour ceiling and `runFingerprintProbe` are all unexported there. A new package would either
duplicate them or force a wide export surface, and the REQ forbids a second verification platform.
The package name then under-describes its contents; renaming it is a separate mechanical change and
is not done here.

**D-04 — The fast-stage fingerprint seals working-tree bytes, not committed object ids. DECIDE & STATE.**
`laneFingerprint` can seal committed objects only because `RunLanes` first refuses a dirty tree
(`heavy_run.go:99`, `205–234`). The fast gate exists to run on a dirty tree, so committed objects
would be a false green the moment anything is uncommitted — which requirement 2 names explicitly.
The mechanism already exists: `fingerprintUntrackedFiles` hashes real file bytes; T2 generalizes it
to tracked paths.

**D-05 — Fast-stage records live in a separate key space (`do-work-fast-stages`, own schema
version, `stage_id` field). DECIDE & STATE.**
A fast green is computed on a possibly-dirty tree and is not attributable to a revision. If the two
key spaces could cross, a dirty-tree fast green could authorize a heavy lane, which is a strictly
weaker guarantee standing in for a stronger one.

**D-06 — Reusable stages are exactly the two Go test stages. DECIDE & STATE.**
They are 73 s of a ~102 s gate and ~85 % of its CPU. Excluded, with reasons:
`go vet` ×2 (Go already caches vet results content-addressed — 0.2 s warm; caching a cache);
gofmt (0.07 s — the fingerprint would cost more than the check);
the version/floor probes (they *are* the toolchain check; caching them caches what detects a
toolchain change); the aggregate contract suite and every probe inside it (the wide markdown and
reference probes such as `shipped-package-reference-contract.sh` resolve link targets against
arbitrary repository paths, so anything narrower than a whole-tree seal is a false green — and the
narrow probes are already 0.02–1.4 s); `gate-runner.sh` (it decides whether to run at all).
ShellCheck lint (3.7 s, a genuinely closed 92-file input set) is a legitimate future candidate and
is named as a follow-up rather than done here — it is 3.6 % of the gate against three more CLI
round trips.

**D-07 — Three named environment variables are removed from the decision child's environment, in
the shell, with the reason stated at the call site. ESCALATE-adjacent, decided: DECIDE & STATE.**
`maintainer-verify.sh:20–29` exports `DO_WORK_TEST_RUN_ID` (a timestamp-plus-PID label),
`DO_WORK_TEST_OTHER_GATE_PROCESSES` (a `ps` sample) and `DO_WORK_TEST_DURATION_LOG` (an output
path) before any tier runs. All three change or could change between runs, and the fingerprint
seals the whole environment — so with them present, **reuse could never fire**. They are removed
from the decision child only, in the shell, so the Go rule "seal every variable, no exceptions"
stays absolute and the exclusion is one greppable line rather than a hand-maintained list inside
the fingerprint. `DO_WORK_TEST_ENFORCE_BUDGET` and `DO_WORK_TEST_FILE_BUDGET_SECONDS` are **not**
excluded: they change the verdict.
*Value:* reuse becomes possible at all. *Risk:* the three excluded variables label and route the
duration log rather than deciding any assertion; if that ever stops being true the exclusion is
wrong, which is why it is one visible line at the call site and not buried in Go.

**D-08 — A reused stage inherits its per-file budget verdict along with its pass/fail verdict; the
reuse line says so. ESCALATE.**
A reused stage writes no duration rows and enforces no per-file budget for that run. The budget
exists to catch tests getting slower as they are added — an input-determined property, which the
fingerprint covers. It does **not** cover a breach caused purely by contention, and five gate runs
during this work run failed on exactly that. So reuse suppresses a check that has recently been
producing false failures, which is arguably a feature and arguably a coverage loss.
*Value:* the gate stops failing on other sessions' load for a stage whose inputs are provably
unchanged. *Risk:* a real contention problem goes unreported for up to four hours on a matching
tree; reversible by disabling reuse for that stage in the manifest. **Recommendation: accept, and
print the inherited verdict** — `REUSED (... recorded <ts>)` names the run whose budget verdict is
standing in.

**D-09 — `record-fast-stage` recomputes the fingerprint and refuses on mismatch. DECIDE & STATE.**
The fast-tier analogue of `verifyLaneRevision`'s before-and-after check (`heavy_run.go:157–159`).
It catches a stage that modified its own inputs while running, which would otherwise record a green
against a tree that no longer exists. Costs one extra ~0.4 s fingerprint per green stage.

**D-10 — `--self-test` keeps its nine-stages-exactly-once assertion; the reuse wrapper is bypassed
under `MAINTAINER_VERIFY_SELFTEST_LOG`. DECIDE & STATE.**
`assert_success_stages:277–309` proves the gate's *stage list*. The self-test shim's `bash` case
accepts only `contract-regressions.sh` (`:247–252`), so a decision call would exit 64 there anyway.
Reuse gets its own probe (T5) where the shim would have limited it. The lesson from REQ-187 says
this count must change deliberately — it deliberately does not change.

**D-11 — Reuse is on by default in the fast gate, and there is no opt-out flag on any shipped
script. DECIDE & STATE.**
`prime-shell-commands.md` § Every Flag on a Shipped Script Needs a Non-Test Caller: a
`--no-reuse` flag on `skills/do-work/tools/*` would have no non-test caller. `_dev/tests/` is
maintainer-side and export-ignored, so the measurement protocol's forced-execution runs use an
environment variable read only by `maintainer-verify.sh` (`DO_WORK_FAST_STAGE_REUSE=off`), which is
exempt.

---

## 4. Testing approach

### T1 — no new test file

The proof is that every existing assertion in `session-start-hook-behavior.sh` still passes
byte-identically (the table in §2 names which assertion catches which broken boundary), plus the one
new shared-tree-unchanged assertion. The RED for T1 is a **measurement**, not a failing assertion:
record the standalone wall time and the copy count before the change (12.90 s, 9 copies) and after.
`lessons-shell-commands.md` (REQ-243/263) warns that a RED which was already red means the work is
half done — here there is no behavioral RED to have, and saying so is more honest than inventing
one.

### T2 — `fast_stage_evidence_test.go`, copying `heavy_reuse_regression_test.go`'s table shape

Each case asserts **both** the disposition and the exact reason code. Positive first, then every
counter-case:

| Case | Expected |
|---|---|
| nothing changed | `reused` / `fingerprint_match` |
| covered source file modified **and committed** | `executed` / `fingerprint_mismatch` |
| covered source file modified **and left uncommitted** | `executed` / `fingerprint_mismatch` |
| covered testdata/fixture file changed | `executed` / `fingerprint_mismatch` |
| a `_dev/tests` gate script changed | `executed` / `fingerprint_mismatch` |
| new untracked non-ignored file under coverage | `executed` / `fingerprint_mismatch` |
| a path no stage classifies changed | `executed` / `fingerprint_mismatch` |
| a `do-work/` file changed | `reused` / `fingerprint_match` (the declared exemption, pinned) |
| an environment variable added, removed, or changed | `executed` / `fingerprint_mismatch` |
| toolchain probe output differs | `executed` / `fingerprint_mismatch` |
| toolchain probe exits non-zero, or exceeds its 5 s bound | `executed` / `fingerprint_uncertain` |
| a tracked covered path deleted from the worktree | `executed` / `fingerprint_uncertain` |
| a covered path replaced by a symlink | `executed` / `fingerprint_uncertain` |
| stage declares no `fingerprint` block | `executed` / `fingerprint_uncertain` |
| record older than 4 h | `executed` / `evidence_expired` |
| record with non-zero exit status, or `skipped`, or empty fingerprint | `executed` / `evidence_unusable` |
| record from a different repository identity | `executed` / `evidence_unusable` |
| record stamped in the future | `executed` / `evidence_unusable` |
| a **heavy-lane** record written at the fast stage's key | `executed` / `evidence_unusable` |
| evidence directory group- or world-readable | `executed` / `evidence_unusable` |
| store unreadable | `executed` / `evidence_unusable` |

Plus, for the requirement-3 named case: a change confined to
`skills/do-work-board/tools/queue-kanban` yields `queue-kanban-fast-tests` → `executed` and
`do-work-cli-fast-tests` → `reused`.

The environment and probe counter-cases have a ready pattern:
`heavy_maintainer_tree_test.go:94,105,133,165` already drives fingerprint-affecting settings, and
`do-work-cli-launcher-behavior.sh:28–56` shows the fake-`go`-on-PATH shape for a runtime change.

### T5 — `_dev/tests/fast-stage-reuse-behavior.sh`, copying `run-go-tests-with-budget-behavior.sh`

Exactly that file's shape: `set -euo pipefail`, `repo_root` from `BASH_SOURCE`, the unit overridable
by an environment variable, one `mktemp -d` fixture root with `trap 'rm -rf -- "$fixture_root"' EXIT`,
one small `run_*` helper so each case is one line, explicit `FAIL:` lines on stderr, one success
line on stdout. Fixture: a synthetic git repository with a two-file "module", a synthetic
`fast-stages.json`, and a trivial stage command that touches a marker file so execution is
observable. Cases:

1. first run **executes**, prints the `EXECUTING` line, touches the marker, records evidence;
2. second run **reuses**, prints the `REUSED` line naming the recorded timestamp, and does **not**
   touch the marker;
3. a failing stage returns its **exact** non-zero status and records nothing; the next run reports
   `no_prior_evidence`;
4. a stage killed by a signal mid-run records nothing and the next run reports `no_prior_evidence`
   (the interruption case; the child carries its own time bound per
   `crew-members/testing.md` § Background Work and Synthetic Load);
5. editing a covered file between runs forces execution;
6. editing only a `do-work/` file between runs still reuses;
7. a `decide-fast-stage` that cannot run (unreadable manifest) makes the gate **execute**, not skip.

### Regression gates for the whole REQ

`bash _dev/tests/maintainer-verify.sh --self-test`, `gofmt -l`, `go vet ./...` in both modules, and
the canonical `bash _dev/tests/maintainer-verify.sh` in a detached worktree.

---

## 5. Measurement protocol

Requirement 5 in full: fixed worker limit, recorded toolchain/cache state, no competing expensive
gate, no synthetic load, wall **and** process-tree CPU, smallest bounded comparison, no benchmarking
by saturation.

**Where.** A detached worktree (`git worktree add --detach`), never the shared main tree — other
sessions dirty it, and five gate runs during this work run failed on per-file budgets purely from
contention. Pre-create the `do-work/test-durations.tsv` header in the worktree, or the duration
logger writes into a file the worktree does not have.

**Revisions.** `B` = the REQ's implementation base (the commit the builder branches from). `A1` = the
merge of T1 only (isolates setup savings). `A2` = the merge of the whole REQ (adds avoided
execution). Three revisions, not two, because acceptance criterion 3 requires the evidence to
separate setup savings from avoided execution.

**Fixed worker limit.** `GOMAXPROCS=4` exported for every run at every revision. Record it. The
fast-tier probe batch fan-out is 11 and is **unchanged by this work** — no concurrency is added
anywhere, which is also what keeps "greater concurrency alone" out of the acceptance claim.

**Recorded state, captured once per session and re-checked at each revision:**
`go version`, `go env GOCACHE`, `du -sh "$(go env GOCACHE)"`, `shellcheck --version`,
`sw_vers -productVersion`, `sysctl -n hw.logicalcpu`, `git rev-parse HEAD`, `git status --porcelain`.

**Valid window, and what invalidates a measurement.** Immediately before and immediately after each
run:

- `sysctl -n vm.loadavg` — the 1-minute figure must be **≤ 2.0** at both samples (8 logical CPUs);
- `/bin/ps -Ao pid=,comm=,args=` must show **no** other `maintainer-verify.sh`, no `go test`, and no
  sibling agent gate.

If either sample fails, or the two 1-minute figures differ by more than 1.0, the run is **discarded
and repeated** — never adjusted, never averaged in. Every recorded number carries its two load
samples. Today's load has ranged 2 to 59, so an invalid window is the expected case and the protocol
has to say what to do about it rather than assume one.

**Cache warming.** One discarded warm-up run at each revision before the measured ones, so no
revision pays a cold link the other does not.

**Commands.** `/usr/bin/time -p` under-reports grandchildren here (it gave `user 0.32 sys 1.35` for a
probe whose real tree cost was several CPU-seconds). Use the bash `time` keyword, which reports
`getrusage(RUSAGE_CHILDREN)` aggregated across every waited-for descendant:

```bash
bash -c 'TIMEFORMAT="%3R %3U %3S"; time bash _dev/tests/maintainer-verify.sh'
```

Note that the `nohup`'d gate-runner escapes this accounting entirely; it must not be running.

**Repetitions and ordering.** Three measured runs per condition, interleaved
`B, A1, A2, B, A1, A2, B, A1, A2` so any thermal or background drift spreads across conditions
rather than landing on one. Report the **median** wall and the median `user+sys`. Hard cap of 12
measured runs total; if three valid repetitions per condition cannot be obtained inside that cap,
record what was obtained and state that the quiet window was not available — an explicitly
incomplete measurement, which requirement 5 permits, is honest where an averaged noisy one is not.

**Four conditions, because reuse has two states:**

| Condition | What it isolates |
|---|---|
| `B` | baseline |
| `A1` | setup savings from T1 alone |
| `A2` cold evidence (`DO_WORK_FAST_STAGE_REUSE=off`, or a first run after invalidation) | that reuse machinery adds no measurable cost when nothing is reused |
| `A2` warm evidence (a second consecutive run on an unchanged tree) | avoided execution |

Contention is separated by the recorded load samples and by every condition coming from one window.

**Recorded alongside:** the existing per-file rows in `do-work/test-durations.tsv` for the runs that
executed, and the new `gate wall <n>s stages executed=<n> reused=<n>` line for every run.

---

## 6. Where this plan corrects or extends the exploration

- **The symlink shape's justification is wrong.** `do-work-cli.sh:7` applies `pwd -P` to the
  launcher's directory and appends `/do-work-cli` unresolved (`:8`). The exploration's claim that
  "the symlink collapses to the shared physical path" via `pwd -P` does not follow. D-01 avoids
  depending on it.
- **`DO_WORK_TEST_RUN_ID` would have made reuse impossible.** It is a timestamp-plus-PID exported at
  `maintainer-verify.sh:20` before any tier runs, and the fingerprint seals the whole environment.
  The exploration's §6 sketch did not surface this; D-07 handles it.
- **The `QUEUE_KANBAN_*` selectors are passed as environment, not argv** (`:606–609`), so a
  decision computed in a separate process would not see them. T3 moves them into the sealed argv.
- **The board stage's `-short`-equivalent is an exclusion list, not a flag.** The CLI stage's
  `-short` is argv and is sealed; the board stage's selection is `DO_WORK_GO_TEST_EXCLUDE_PREFIXES`,
  which the same argv move seals.
- **`_dev/tests/heavy-lanes.json` declares `non_heavy_coverage` including `suffix-under . .md` and
  `suffix-under . .go`** — i.e. every markdown and Go file in the repository is "classified". The
  fast-stage manifest deliberately does **not** copy that: it declares only `do-work/`, so every
  other unclassified path forces both stages.
- **`do-work/test-durations.tsv` is gitignored** (`.gitignore:3`) and `do-work/` is skipped by the
  untracked seal anyway, so the duration log cannot invalidate its own stage's fingerprint.
- **Genuine maintainer-tree reads exist in the module tests** (four files at minimum), so a
  module-only coverage declaration would be a false green. T3 enumerates them; T4 pins them.

---

## 7. Task count

**Five.**

1. T1 — share one skill root in `_dev/tests/session-start-hook-behavior.sh`.
2. T2 — fast-stage evidence engine and three CLI commands in `internal/heavyverification`, tests first.
3. T3 — `_dev/tests/fast-stages.json`, with its coverage extras enumerated from the module tests.
4. T4 — gate wiring, summary line, self-test bypass guard, maintainer-tree pinning assertions.
5. T5 — `_dev/tests/fast-stage-reuse-behavior.sh` and its registration, plus the §5 measurement.

T2 is the largest and splits cleanly at the engine/command seam if smaller increments are wanted.

---

## 8. Requirement coverage check

| Requirement / criterion | Covered by |
|---|---|
| R1 reduce repeated setup in measured hot paths; immutable material shared; writable state isolated; real launcher, build failures, missing tools and installed layout retained | **T1** (§2 table names the retained assertion per scenario) |
| R2 avoid repeating verification whose complete relevant inputs are unchanged; use existing mechanisms; source, transitive deps, fixtures, gate scripts, configuration, toolchain/runtime, **including uncommitted changes**; unknown/incomplete/unverifiable → broader verification | **T2** (working-tree seal, D-04; unclassified-forces-all; `fingerprint_uncertain` on any undeterminable input), **T3** (coverage + toolchain probes) |
| R3 record which checks executed, which reused, and why; preserve failure and interruption statuses; skipped/failed/incomplete supplies no evidence; board-only change reuses CLI evidence only after proof | **T4** (per-stage `EXECUTING`/`REUSED` lines with reason, summary line, exact status return, invalidate-before-execute, record only on zero), **T2**/**T5** (interruption and failure cases; the board-only case is an explicit test) |
| R4 keep useful failure coverage; rollback/interrupted-recovery/locking/concurrent-writer/cleanup tests preserved; no silent assertion removal, no mock for a real boundary, no move to the heavy tier | **T1** (no assertion changed, one added), **T4** (heavy tier and `--heavy-lane` untouched; `--self-test` nine-stage count unchanged); no test is deleted, split, renamed or moved anywhere in this plan |
| R5 before/after on comparable revisions, fixed worker limit, recorded toolchain/cache state, no competing gate, no synthetic load, wall **and** process-tree CPU, smallest bounded comparison | **§5** (three revisions, four conditions, `GOMAXPROCS=4`, recorded state, load-gated window, bash `time` keyword, 3 interleaved repetitions, 12-run cap) |
| AC1 repeated setup measurably reduced in ≥1 hot path, assertions retained, writable fixtures independent | **T1** |
| AC2 matching-input case reuses; source/dependency/fixture/gate-script/runtime change invalidates; uncertain impact and dirty relevant inputs cannot reuse stale evidence; positive and negative tests | **T2** test table (§4), **T5** end-to-end probe |
| AC3 end-to-end duration improves under recorded comparable conditions; evidence separates setup savings, avoided execution and contention; no one-off noisy sample, no rename/split, no raised timeout, no concurrency-only claim | **§5** four conditions; no file is renamed or split; no timeout is raised; probe fan-out unchanged at 11 |
| AC4 fast-tier per-file budget and heavy-tier policy still satisfied; required correctness checks pass; rollback/recovery/locking/cleanup tests still detect their failures; reused results visibly distinguishable | **T4** (`REUSED` line; heavy policy untouched), **D-08** (the budget verdict is inherited, stated, and printed — this is the one place the plan trades a check, and it is escalated rather than assumed) |

**Nothing in the REQ is left uncovered.** The single point where a check is weakened rather than
preserved is D-08 (a reused stage enforces no per-file budget for that run), and it is escalated
with its value, its risk and its reversal.
