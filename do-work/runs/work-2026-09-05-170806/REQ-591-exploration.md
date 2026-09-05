# REQ-591 exploration — reduce repeated setup and unaffected reruns in the fast gate

Read-only exploration. No tracked file was modified. All measurements redirected
`DO_WORK_TEST_DURATION_LOG` into the scratchpad, so `do-work/test-durations.tsv` is untouched.

## Measurement conditions (applies to every number below)

- Machine: Apple M4 Max (Virtual), arm64, 8 logical CPUs, Darwin 25.5.0.
- Toolchain: `go1.26.1 darwin/arm64`, ShellCheck 0.11.0, GOCACHE warm at 24 GB.
- Repository HEAD `e0bdf8bf92d7c2b9663cff64a3d85a5220216991`; working tree carries only the
  claimed REQ-591 file plus three untracked hand-back files under `do-work/runs/`.
- Load average during measurement 2.5–4.9 (1-minute), i.e. a **loaded window** — other Claude
  agents in this session were active. No `maintainer-verify.sh` process was running (checked
  with `ps` before each measurement). Every number is therefore an upper bound, not an
  uncontended baseline. Requirement 5 asks for an uncontended comparison; that has to be
  redone by the builder, ideally in a detached worktree.
- The whole gate was **not** run end to end. Each stage was run exactly once, individually,
  which sums to about one gate's worth of work.

---

## 1. What the fast gate runs, and what each stage costs

### 1.1 Ordered stage list (default/fast tier, `_dev/tests/maintainer-verify.sh` with no argument)

Entry point is `run_verification fast` (`_dev/tests/maintainer-verify.sh:500`). Everything in
`run_verification` is **strictly serial**.

| # | Stage | What it does | Code |
|---|---|---|---|
| 0 | Preamble (script scope, runs before any tier is chosen) | Resolves `repo_root`, builds `test_run_id`, counts `other_gate_processes` via `ps`+`awk`, exports `DO_WORK_TEST_*` | `maintainer-verify.sh:13–31` |
| 1 | Required-command presence | `command -v` for `git go shellcheck bash`; returns 1 if any is missing | `:514–519` |
| 2 | Go version floor | `go version`, compares against `minimum_go_version=go1.26.1` with `version_at_least` | `:521–528` |
| 3 | gofmt resolution | `go env GOROOT`, requires `$GOROOT/bin/gofmt` executable. Deliberately not PATH gofmt | `:529–540` |
| 4 | ShellCheck version floor | `shellcheck --version`, parses the `version: ` line, floor `0.11.0` | `:542–555` |
| 5 | Tracked shell enumeration | `git -C repo ls-files -z -- '*.sh'`, refuses on empty | `:557–563` |
| 6 | **ShellCheck lint** | `shellcheck --severity=warning -- <92 tracked .sh files>` from `repo_root` | `:564–567` |
| 7 | Tracked Go enumeration | `git ls-files -z -- '*.go'`, keeps only paths that exist on disk | `:569–576` |
| 8 | **gofmt formatting check** | `$GOROOT/bin/gofmt -l -- <274 files>`; verdict is emptiness of stdout, never exit status | `:577–589` |
| 9 | **Aggregate contract suite** | `DO_WORK_MAINTAINER_TIER=fast bash _dev/tests/contract-regressions.sh` | `:591–593` |
| 10 | queue-kanban `go vet ./...` | in `skills/do-work-board/tools/queue-kanban` | `:595–599` |
| 11 | **queue-kanban fast tests** | `QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off DO_WORK_GO_TEST_EXCLUDE_PREFIXES=TestJavaScriptBehavior,TestBrowserBehavior` → `run-go-tests-with-budget.sh <module> ./...` | `:605–614` |
| 12 | do-work-cli `go vet ./...` | in `skills/do-work/tools/do-work-cli` | `:650–654` |
| 13 | **do-work-cli fast tests** | `run-go-tests-with-budget.sh <module> -short ./...` | `:655–661` |
| 14 | Success line | `Maintainer verification passed.` | `:663` |

Fast tier never runs: the strict JavaScript lane, the strict browser lane,
`DO_WORK_HEAVY_TESTS=1`, or the three heavy probes (staged-skills, updater, installer).
`run-go-tests-with-budget.sh` forces `-count=1`, so Go's own test-result cache is defeated by
design (`run-go-tests-with-budget.sh:41`).

`--self-test` (`run_self_test`, `:373–498`) is a separate mode that builds a shimmed fixture
repository, asserts exactly nine stages run exactly once per tier, injects a failure at each
of eleven stages, and mutation-tests the strict-JavaScript marker. It is not part of the
default tier.

### 1.2 How probes are registered and launched

Three files:

- **`_dev/tests/contract-regressions.sh`** (77 lines, with a hard ratchet ceiling of 77 lines
  self-checked at `:17–22`). It:
  1. Validates the tier (`fast`/`heavy`), else exit 2.
  2. Sources `test-duration-log.sh`.
  3. Runs **four owner contracts serially and in-process** via `run_contract_file`, which
     `source`s each inside a subshell and times it with `date +%s` (whole seconds):
     `contracts/core-checks.sh`, `contracts/queue-kanban.sh`,
     `contracts/replace-text-section.sh`, `contracts/recovery-set-aside.sh`.
  4. Sources `probe-batch.sh` (defines `launch_probe`/`collect_probes`).
  5. Sources `contracts/probe-lanes.sh` (registers and **launches** every probe).
  6. Calls `collect_probes` (a bare `wait`, then prints outputs in launch order).

- **`_dev/tests/contracts/probe-lanes.sh`** (59 lines) is registration only. `register_probe`
  checks the script is executable (counts a failure if not) and calls `launch_probe`.
  **Fast tier registers 11 probes**, in this order:
  `suite-manifest-contract.sh`, `shipped-package-reference-contract.sh`,
  `action-shell-blocks.sh`, `session-start-hook-behavior.sh`,
  `prescribed-shell-canonicalization.sh`, `defensive-surface-audit.sh`,
  `audit-lockins.sh`, `do-work-cli-launcher-behavior.sh`,
  `p50-estimator-determinism.sh`, `select-simple-reqs-behavior.sh`,
  `run-go-tests-with-budget-behavior.sh`.
  Heavy tier additionally builds one shared `do-work-cli` binary
  (`go build -ldflags='-s -w'`), exports `DO_WORK_TEST_DO_WORK_CLI_BINARY`, and registers
  `staged-skills-contract.sh`, `update-script-behavior.sh`, `install-suite-behavior.sh`.
  Fast tier prints the `SKIP:` line for those three instead.

- **`_dev/tests/probe-batch.sh`** (66 lines) owns execution. `launch_probe` writes the failure
  message to `$probe_batch_root/<name>.message`, turns **job control on** (`set -m`, so the
  child gets its own process group — the update-script probe's interrupt case sends SIGINT),
  and backgrounds a subshell that runs the scripts in order, times each with `date +%s`,
  prints `test-file duration: <basename> <n>s (limit …)`, records a duration row, and enforces
  the per-file budget in fast tier. Output is redirected to `<name>.out`.

**Concurrency: unbounded.** All 11 (fast) or 14 (heavy) probes are launched at once with no
worker cap, no `GOMAXPROCS` bound and no semaphore, onto 8 cores. `collect_probes` is a bare
`wait`. This matters for requirement 5: the *recorded* per-probe durations in
`do-work/test-durations.tsv` are 11-way contended numbers, not costs.

Ordering inside the aggregate lane is therefore: **~10 s of serial owner contracts, then a
~14 s parallel probe burst.** The four owner contracts do not participate in the batch.

### 1.3 Stage-by-stage cost (measured individually, this machine, loaded window)

| Stage | Wall | CPU (user+sys) | Notes |
|---|---|---|---|
| ShellCheck lint (92 files) | **3.70 s** | 3.51 + 0.12 | single-threaded |
| gofmt `-l` (274 files) | **0.07 s** | 0.29 + 0.06 | already negligible |
| queue-kanban `go vet ./...` | **0.18 s** warm | 0.14 + 0.38 | Go caches vet results; cold is seconds |
| do-work-cli `go vet ./...` | **0.20 s** warm | 0.16 + 0.41 | same |
| Aggregate contract suite | **23.94 s** | 9.08 + 9.12 | ≈10 s serial contracts + ≈14 s parallel probes |
| queue-kanban fast tests | **21.64 s** | 14.47 + 6.41 | 396 tests, slowest file `generate_test.go` 10.46 s |
| do-work-cli `-short` tests | **51.56 s** | **66.54 + 99.05** | 773 tests, slowest file `finalization_recovery_test.go` 20.40 s |
| version probes + `git ls-files` ×2 | ~0.3 s | — | |
| **Fast gate total (sum of stages)** | **≈ 102 s** | | |

The do-work-cli lane is 50 % of gate wall time and ~85 % of gate CPU, and its CPU is
**sys-dominated** (99 s sys vs 67 s user) — that is process spawning (git subprocesses), which
is exactly what REQ-574 identified.

Individual fast probes, run **serially and alone** (uncontended cost):

| Probe | Alone | As recorded in the 11-way batch |
|---|---|---|
| `session-start-hook-behavior.sh` | **12.90 s** | 14–25 s |
| `prescribed-shell-canonicalization.sh` | 5.03 s | 9–18 s |
| `action-shell-blocks.sh` | 3.86 s | 6–14 s |
| `suite-manifest-contract.sh` | 2.17 s | 4–10 s |
| `p50-estimator-determinism.sh` | 1.38 s | 3–8 s |
| `shipped-package-reference-contract.sh` | 1.21 s | 5–12 s |
| `run-go-tests-with-budget-behavior.sh` | 1.08 s | 2 s |
| `do-work-cli-launcher-behavior.sh` | 0.93 s | 2 s |
| `select-simple-reqs-behavior.sh` | 0.84 s | 3–6 s |
| `audit-lockins.sh` | 0.70 s | 3 s |
| `defensive-surface-audit.sh` | 0.02 s | 0 s |
| **serial sum of the 10 non-session-start probes** | **17.22 s** | |

Owner contracts, run alone: `core-checks.sh` 5.25 s, `replace-text-section.sh` 8.21 s,
`queue-kanban.sh` 0.05 s, `recovery-set-aside.sh` 0.04 s. Serial total 13.55 s.

**Conclusion for the aggregate lane:** `session-start-hook-behavior.sh` alone (12.9 s) is
longer than all ten other probes put together (17.2 s serial, ~3–4 s if perfectly parallel).
It is the lane's critical path. The second cost is that the two slow owner contracts (13.5 s
of serial work) do not join the parallel batch.

---

## 2. `_dev/tests/session-start-hook-behavior.sh` — the named setup-reuse candidate

142 lines. `set -uo pipefail`, one `mktemp -d` fixture root, `trap rm -rf EXIT`.

### 2.1 How many times it copies the CLI module

`install_core_fixture` (`:10–17`) is called **9 times**:

| # | Caller | Line |
|---|---|---|
| 1–4 | `run_banner_case valid / missing / reformatted / multiple` | `:36–40` via `:22` |
| 5–6 | `run_tail_case unfinished-journal / archived-without-commit` | `:77–78` via `:45` |
| 7 | authority matrix (`install_core_fixture "$authority_root"`) | `:83` |
| 8–9 | `run_launcher_failure_case go-too-old / go-build-failure` | `:123–124` via `:112` |

Each call does:
```
cp  skills/do-work/hooks/session-start.sh   -> case_root/skill/hooks/
cp  skills/do-work/tools/do-work-cli.sh     -> case_root/skill/tools/
cp -R skills/do-work/tools/do-work-cli      -> case_root/skill/tools/     (3.4 MB, 212 files)
rm -f case_root/skill/tools/do-work-cli/do-work-cli
```
A tenth root (`missing-launcher`, `:126–128`) copies only `session-start.sh` — no module.

### 2.2 Instrumented timing (a copy of the script run from the scratchpad, repo_root repointed)

```
PHASE setup                        0.00s
PHASE 4-banner-cases               6.37s
PHASE 2-tail-cases                 3.14s
PHASE authority-matrix             1.71s
PHASE 2-launcher-failure-cases     0.51s
PHASE missing-launcher             0.01s
PHASE grep-check                   0.01s
COPIES 9   total-copy-seconds 0.95
```
Whole script standalone: **12.90 s** wall (`/usr/bin/time -p`), of which the nine `cp -R`
calls are **0.95 s total (0.106 s each)**.

**The copy is not the cost.** The cost is that each copy is a *new absolute path*, and
`do-work-cli.sh` resolves `script_dir` with `pwd -P` (`do-work-cli.sh:7`) and then runs
`go tool -C "$module_dir" -n do-work-cli` (`:52`). A distinct module directory produces a
distinct Go build-cache action, so every fixture root pays a fresh link:

```
cp -R of the module                      real 0.09s
go tool -n in a NEW copy   (path A)      real 1.50s  user 3.52s  sys 1.06s  -> .../7f9e0e96…-d/do-work-cli
go tool -n again, same copy              real 0.08s  user 0.09s  sys 0.22s  -> same 7f9e0e96…
go tool -n in a SECOND new copy (path B) real 1.08s  user 3.45s  sys 0.98s  -> .../799560195…-d/do-work-cli
go tool -n in the real repo module       real 0.07s  user 0.07s  sys 0.19s  -> .../44a1bcff…-d/do-work-cli
```
Three different directories, three different cache action IDs, ~1.1–1.5 s wall and ~4.5 s CPU
each; repeats within one directory are 0.08 s. Seven of the nine roots reach a real
toolchain, so ≈ 7 × 1.5 s ≈ 10.5 s of the 12.9 s is repeated linking of a byte-identical
binary.

(Note the contrast: six Go test packages each run `go build -o <own temp> ../../cmd/do-work-cli`
from the *same* source directory, and warm repeats there cost only 0.33 s — same action ID,
cache hit. The penalty here is specific to copying the module to a new path.)

### 2.3 What each scenario needs, and what is immutable

| Scenario | Boundary under test | Reads | Writes |
|---|---|---|---|
| 4 banner cases (`:19–40`) | Version parsing from `actions/version.md` + pending-REQ count + exact banner bytes + empty stderr + status 0 | real `session-start.sh` → real `do-work-cli.sh` → real `go tool` | `skill/actions/version.md` (**different per case — this is the input under test**), `project/do-work/queue/REQ-*.md` |
| 2 tail cases (`:42–78`) | Unfinished-finalization tail line, and that the journal/archive file is **not mutated** (sha256 before/after) | same real launcher path | `project/` git repo (`git init -q`), journal JSON or archive REQ, `version.md` |
| authority matrix (`:82–105`) | Committed vs uncommitted working REQ → typed `deleted` housekeeping vs `RESERVATION-GIT-AUTHORITY-UNAVAILABLE` fail-closed; `PATH=<nonexistent>` removes Git | `do-work-cli.sh` once, then `go tool -C … -n` and the binary directly | `project/` git repo + commit, `do-work/working/`, `do-work/.req-reservations/`, `version.md` |
| 2 launcher-failure cases (`:109–124`) | Too-old Go and build failure → exit 2, empty stdout, actionable stderr | **fake `go` on PATH**, no real build | fake-bin dir, `version.md` |
| missing-launcher (`:126–134`) | Absent `tools/do-work-cli.sh` → nonzero + "canonical launcher is missing" | hooks only | — |
| static grep (`:136–139`) | Retained hook is a thin launcher (no `sed`/`awk`/`find`/domain scripts) | real repo file | — |

**Immutable in every scenario:** `skill/hooks/session-start.sh`, `skill/tools/do-work-cli.sh`,
and the whole `skill/tools/do-work-cli/` module. That the launcher writes nothing into the
module tree is already pinned by a separate probe —
`_dev/tests/do-work-cli-launcher-behavior.sh:118` asserts the module directory is still empty
after a run — so sharing it is not an unverified assumption.

**Must stay per-scenario:** `skill/actions/version.md` (it *is* the banner input), the entire
`project/` tree (queue files, git repository, journal, archive, reservations), the fake-`go`
PATH directory, and the `missing-launcher` root (it must lack `tools/do-work-cli.sh`).

### 2.4 What can concretely be shared

Build **one** skill root under the fixture root holding `hooks/`, `tools/do-work-cli.sh` and
`tools/do-work-cli/`, and give each scenario its own `project/`. Two shapes, both preserving
the real launcher path:

- **Shared root, per-case `version.md` rewrite.** Cases already run strictly serially, so
  rewriting `<shared-root>/actions/version.md` per case is safe. Simplest diff.
- **Per-case skill root whose `tools/do-work-cli` is a symlink to one shared copy.**
  `do-work-cli.sh` resolves with `pwd -P`, so the symlink collapses to the shared physical
  path and hits the same cached link action. Keeps `version.md` genuinely per-case.

Either collapses 7 fresh link actions to 1. Expected: **12.9 s → ~3.5 s standalone**, and the
aggregate contract lane's parallel burst drops from ~14 s to ~9 s (next critical path becomes
`prescribed-shell-canonicalization.sh`), i.e. the lane goes from ~24 s to ~19 s.

Coverage that must be retained explicitly, per requirement 1: the real launcher path (banner
and tail cases still exec `session-start.sh` → `do-work-cli.sh` → `go tool`), the two build
failure/too-old cases (they need their own fake-`go` PATH, not their own module copy), the
missing-launcher case (own root), and the installed-layout shape the `--skill-root` argument
exercises in the authority matrix.

---

## 3. Other repeated-setup hot paths, ranked

| Rank | Location | The repetition | Repetitions | Rough cost | Notes |
|---|---|---|---|---|---|
| **R1** | `_dev/tests/session-start-hook-behavior.sh:10–17` | `install_core_fixture` copies the 3.4 MB CLI module to a fresh path | 9 (7 reach a real toolchain) | **~10.5 s of 12.9 s** | Section 2. Fast tier. Highest single win in the aggregate lane. |
| **R2** | `internal/lifecycleadvance/*_test.go` | 13 distinct `runAdvanceGit(t, root, "init", "-q")` sites, each followed by two `git config` calls; no `TestMain` template; 37 top-level tests; **zero `t.Parallel()`** | ~13 helper sites × N cases | package sums 41 s of file time (`finalization_gate_test.go` 14.9 s, `queue_commands_test.go` 11.3 s, `evidence_gates_test.go` 8.2 s, `recovery_commands_test.go` 7.1 s) | Exactly the REQ-574 shape, untouched. |
| **R3** | `internal/heavyverification/*_test.go` | Three separate repository builders (`newLaneEvidenceRepository:76`, `newHeavyRunRepository:61`, `newRuntimeEvidenceRepository:66`) plus inline `git init` at `heavy_verification_test.go:182,345`; 49 top-level tests; zero `t.Parallel()` | many | package sums ~49 s (`heavy_reuse_regression_test.go` 17.3 s, `heavy_maintainer_tree_test.go` 11.0 s) | Newest, largest uncovered package. |
| **R4** | `internal/knowledgecommands/*_test.go` | 11 `git init` sites, 57 top-level tests, zero `t.Parallel()` | 11+ | `memory_commands_test.go` 7.7 s, `bkb_init_test.go` 4.8 s | |
| **R5** | `_dev/tests/contracts/core-checks.sh` | 10 `mktemp -d` probe directories, 6 of them `git init -q` + config + commits, built inline | 6 repos | 5.25 s alone | Runs **serially** before the probe batch, so it is on the lane's critical path twice over. |
| **R6** | `internal/requeststate`, `internal/gittransaction`, `internal/suiteinstall`, `internal/cleanup`, `internal/corehelpers` | one repository builder each, no `TestMain` template, zero `t.Parallel()` | per test | 13.7 s / 10.3 s / 15.7 s / 14.4 s / 21.7 s of file time | Lower per-package density than R2/R3. |
| **R7** | `_dev/tests/update-script-behavior.sh:126–133` | `make_suite_fixture` copies the CLI module into each suite fixture (same fresh-path link penalty as R1) | 3 static `cp -R` sites + loops | heavy tier only | Out of fast-gate scope but the same class. |
| **R8** | `_dev/tests/install-suite-behavior.sh` | 11 `cp -R` sites copying `skills/`, `suite/` and whole project trees | 11+ | heavy tier only | |
| **R9** | `_dev/tests/suite-manifest-contract.sh:44–50` | `clone_valid_fixture` copies a small synthetic archive per case | per case | 2.17 s total | Cheap; the copied tree is tiny. |
| **R10** | Six test packages independently `go build`/`go tool -n` the CLI | `corehelpers`, `finalization`, `lifecycleadvance`, `nextselection`, `publication`, `suiteinstall` | 6 | ~0.33 s each warm (**cache hit — same source dir**) | Measured, not assumed. Not worth attacking; the heavy tier's shared `DO_WORK_TEST_DO_WORK_CLI_BINARY` has no Go consumer, and adding one would perturb the heavy fingerprint (see §5.4). |

Microbenchmark supporting R2–R6, 40 iterations each on this machine:

```
40 × (git init -q + 2 git config)   wall 1.179s   cpu 0.453u + 0.597s   → 29.5 ms each
40 × cp -R of a prepared template   wall 0.304s   cpu 0.030u + 0.267s   →  7.6 ms each
```
≈ 22 ms saved per repository. `os.CopyFS` (REQ-574's actual mechanism) is cheaper still than
shelling out to `cp`.

---

## 4. What REQ-574 already did

`do-work/archive/UR-115/REQ-574-bring-do-work-cli-test-files-under-the-30s-budget.md`, merged
`50569e88c8f1f5234cbdfaf0efaede671d72b13c`, range `982e94f0..50569e88`, 4 files, all test code.

**Techniques:**

1. **Baseline repository built once in `TestMain`, copied per fixture.**
   `finalization_commands_test.go` gained a `TestMain` that builds one initialized, configured,
   empty repository; `newFinalizationRepository` copies it instead of running `git init` + two
   `git config` (55 calls). `publication_commands_test.go` gained a `TestMain` +
   `buildDeferGateRepositoryTemplate` that builds both defer-gate baselines once;
   `newDeferGateRepository` copies the matching template, eliminating seven git commands × 30
   calls. Gotcha recorded: `git init --template=` omits `.git/hooks`, which one test writes
   into, so both `TestMain`s recreate it empty.
2. **Resolve an executable once instead of per case.** `runRetainedInventory` in
   `corehelpers/inventory_test.go` now resolves the CLI with a `sync.Once`
   (`resolveRetainedInventoryExecutable:222`) and invokes it directly with the argv and
   environment `uncommitted-inventory.sh` would have used, instead of two shells and two
   Go-toolchain probes per case. Exactly one caller was deliberately left on the full shim so
   the launcher contract keeps an end-to-end test.
3. **Concurrency where fixtures are independent** — 15 `t.Parallel()` calls in
   `internal/finalization`; tests reassigning package-level hooks (`afterFinalizationPhase`,
   `enumerateTrackedReleasePaths`) were kept serial.

**Recorded result:** module wall 65 s → 61 s, 772 tests both times, worst-file headroom
3.77 s → 5.66 s against the 30 s limit, no assertion changed, no file split. It explicitly
declined to fix the two-concurrent-gates case, calling that a scheduling decision.

**Files REQ-574 covered:** `internal/corehelpers/inventory_test.go`,
`internal/finalization/finalization_commands_test.go`,
`internal/publication/publication_commands_test.go`,
`internal/publication/defer_gate_test.go` (and, transitively, every caller of those two
package fixtures — `finalization_recovery_test.go`, `finalization_req499_test.go`).

**Currently-slow files REQ-574 did NOT cover** (from the most recent full-gate rows in
`do-work/test-durations.tsv`, run `20260905T194912Z-40611`, `other_gate_processes=1`):

| File | Recorded | Has `TestMain` template? | `t.Parallel()`? |
|---|---|---|---|
| `internal/heavyverification/heavy_reuse_regression_test.go` | 17.30 s | no | no |
| `internal/lifecycleadvance/finalization_gate_test.go` | 14.92 s | no | no |
| `internal/requeststate/state_apply_test.go` | 13.67 s | no | no |
| `internal/lifecycleadvance/queue_commands_test.go` | 11.29 s | no | no |
| `internal/heavyverification/heavy_maintainer_tree_test.go` | 10.97 s | no | no |
| `internal/suiteinstall/install_transaction_test.go` | 10.50 s | no | no |
| `internal/gittransaction/git_transaction_test.go` | 10.28 s | no | no |
| `internal/lifecycleadvance/evidence_gates_test.go` | 8.17 s | no | no |
| `internal/cleanup/cleanup_apply_test.go` | 8.11 s | no | no |
| `internal/knowledgecommands/memory_commands_test.go` | 7.68 s | no | no |
| `internal/heavyverification/heavy_evidence_test.go` | 7.67 s | no | no |
| `internal/lifecycleadvance/recovery_commands_test.go` | 7.12 s | no | no |

`t.Parallel()` appears **only** in `internal/finalization` (15 uses). Every other package is
fully serial.

Also worth flagging: in that same recorded run,
`internal/finalization/finalization_req499_test.go` was **37.93 s** — over the 30 s per-file
budget — and `finalization_recovery_test.go` 33.27 s. Under my own uncontended run just now
the whole module measured 51.56 s with slowest file 20.40 s. The budget breach in the log is
contention, not regression, which is precisely the "two gates share the machine" case REQ-574
declined to fix.

---

## 5. Existing evidence and reuse mechanisms, and their real invalidation contracts

### 5.1 Green-gate evidence — `internal/gateevidence/gate_evidence.go`

**Storage.** `storedGateEvidence` (`:27–35`): schema version 1, repository identity, gate
command argv, argv SHA-256, provenance, exit status, recorded revision. Written to
`<git-common-dir>/do-work-green-gates/<sha256(argv)>.json`, mode 0600, in a 0700 directory
(`resolveEvidenceContext:189–221`, `publishEvidenceRecord:238–268`,
`ensurePrivateDirectory:270–292`). Linked worktrees share the store (git *common* dir).

**What it keys on** (`RecordGreenGate:44`, `evaluateEvidenceRecord:107`):
1. Repository identity = symlink-resolved, cleaned Git common directory.
2. `sha256(json.Marshal(argv))` **plus** a full argv element-wise compare (`equalArgv:413`).
3. The recorded commit revision.

**What invalidates it**, in the order `evaluateEvidenceRecord` checks:
- `GateEvidenceMissing` — no record file.
- `GateEvidenceInvalidRecord` — wrong schema version, wrong provenance, non-zero stored exit
  status, empty recorded revision, unreadable/non-regular/wrong-permission file, or the file
  changing identity/size/mtime during the read (`readEvidenceRecord:294–330`).
- `GateEvidenceDifferentRepository` — identity mismatch.
- `GateEvidenceDifferentArgv` — digest or argv mismatch.
- `GateEvidenceRecordedRevisionMissing` — the recorded commit no longer resolves
  (`commitResolvesExactly:345`).
- `GateEvidenceExactRevisionMatch` → **match**, basis `exact_revision`.
- `GateEvidenceRecordedRevisionNotAncestor` — recorded revision is not an ancestor of target
  (`recordedRevisionIsAncestor:356`).
- `GateEvidenceInvalidated` — intervening commits touch anything outside `_dev/gate-runs/`
  (`interveningCommitsAreGateLogs:371`, using `git diff-tree --no-commit-id --name-only -r -m -z`).
- `GateEvidenceLogDescendantMatch` → **match**, basis `gate_log_only_descendant`.

**What it does NOT cover — the gap this REQ has to close.** There is **no timestamp field and
no expiry**, and, decisively, **no working-tree input at all**. A green recorded at HEAD stays
green at that HEAD no matter what is uncommitted, untracked, or changed in the toolchain.
`RecordGreenGate` only ever writes `GateExitStatus: 0`; a non-zero direct status is refused
with `GATE-EVIDENCE-NOT-GREEN` and no record (`gate_commands.go:39–54`).

**Writer:** `_dev/tests/gate-runner.sh:60–67` — runs the gate once per new HEAD, then
`do-work-cli record-green-gate --gate-exit-status 0 -- bash _dev/tests/maintainer-verify.sh`.
**Readers:** `internal/lifecycleadvance/evidence_gates.go:198–201` (the pipeline's
already-green claim path) and `internal/repairvalidation/already_green.go:214`
(`CheckGreenGateAtRevision` at the record's own revision).

### 5.2 Heavy-lane fingerprint — `internal/heavyverification/heavy_evidence.go`

`laneFingerprint(repositoryRoot, lane, manifest, committedTree)` (`:288–363`) digests a
byte-stable `fingerprintDocument` struct (`:264–271`, deliberately a struct not a map) over:

1. **`command_argv`** — the lane's exact argv from the manifest.
2. **Covered committed files** — for every entry of `git ls-tree -r -z --full-tree <revision>`
   (`readCommittedTree:200`), the `mode type objectid path` quadruple is included when the lane
   covers the path **or** the manifest classifies it nowhere. That second clause is explicit:
   "An unclassified input forces all lanes in the planner. Include it here too, or reuse could
   silently undo that conservative selection" (`:301–303`). Any entry that is not `100644`/
   `100755` (symlink, submodule, unsupported mode) is a **hard error**, not a weaker
   fingerprint (`:304–306`). A lane covering no committed path at the revision is an error
   (`:310–314`).
3. **Untracked regular-file bytes** — `fingerprintUntrackedFiles:469–503` runs
   `git ls-files -z --others --exclude-standard` and again with `--ignored`, skips anything
   under `do-work/` (`queueStatePrefix`), and seals `untracked <perm> <sha256> <path>` for
   lane-covered paths (and, for non-ignored files, any manifest-unclassified path). A
   non-regular untracked input is an error.
4. **Bounded toolchain probes** — every `fingerprint.toolchain_probes` argv is run through
   `runFingerprintProbe` with a **5-second** timeout inside its own process group
   (`heavy_fingerprint_probe.go:14–37`); on timeout the group is terminated and the result is
   an error. Output is SHA-256'd. A lane with `Fingerprint == nil` is **never reusable**
   (`:289–291`); `validate()` requires at least one probe (`:235–248`).
   In practice the probe is `_dev/tests/heavy-runtime-fingerprint.py <module> bash git python3 [node]`,
   which seals native-binary bytes for each named tool plus `go`, the module-selected
   `GOROOT/bin/go`, the full `go env -json` (minus `GOGCCFLAGS`), `go list -m -json all`, the
   full `git config --null --list`, and `os.uname()`. It **refuses** (exit 1, "fingerprint
   uncertain") on any of: a non-native executable wrapper, a set `BASH_ENV`/`ENV`/`PYTHONPATH`/
   `PYTHONHOME`/`NODE_OPTIONS`/`LD_PRELOAD`/`LD_LIBRARY_PATH`/`DYLD_INSERT_LIBRARIES`/
   `DYLD_LIBRARY_PATH`, an unexpected `GIT_CONFIG_*`, non-empty `GOFLAGS` or a `GOWORK` other
   than `""`/`off`, a multi-token `CC`/`CXX` under `CGO_ENABLED=1`, an unrecognized Git config
   key, or a local (versionless) module `replace`.
5. **The entire inherited environment** — `:333–356` unions the declared names with **every**
   name in `os.Environ()`, sorts, and seals `{name, set, sha256(value)}` for each. The comment
   states why: "A hand-maintained subset cannot prove that an omitted selector or tool setting
   stayed unchanged."

### 5.3 `reused` vs `executed`

`decideLaneReuse(store, lane, fingerprint, fingerprintError, evaluatedAt)` (`:375–409`) returns
one disposition and exactly one reason. Every path except the last is `executed`:

| Reason constant | Condition |
|---|---|
| `fingerprint_uncertain` | the fingerprint could not be computed (probe failed/timed out, unsupported tree mode, unreadable untracked input, `readCommittedTree` error) |
| `reuse_disabled` | `LaneRunRequest.EvidenceReuse` is false — the **zero value**, so reuse is never accidental (`heavy_run.go:141`) |
| `evidence_store_unavailable` | `store == nil` |
| `evidence_unusable` | read error, or the record fails any of: schema version 2, repository identity, lane id, argv equality, `ExitStatus != 0`, `Skipped`, empty fingerprint, empty execution revision, unparseable `RecordedAt`, or a `RecordedAt` in the future |
| `no_prior_evidence` | no record file |
| `fingerprint_mismatch` | digests differ |
| `evidence_expired` | `evaluatedAt - recordedAt > 4h` |
| `fingerprint_match` | → **`reused`** |

Expiry and fingerprint equality are checked **independently**, so either alone forces
execution, and the reported reason names the first condition that failed.

**Four-hour expiry:** `laneEvidenceMaximumAge = 4 * time.Hour` (`:25`). Its comment: "a ceiling,
never a guarantee. It is checked independently of the fingerprint, so either condition alone
forces a rerun, and a reuse never refreshes the record that authorized it." That last clause is
enforced at `heavy_run.go:172–174`: the stored `RecordedAt` is `executedAt`, the instant the
lane actually ran, so reuse cannot extend the window.

**Surrounding guarantees in `heavy_run.go`:**
- `refuseDirtyTrackedTree:205–234` — `git status --porcelain -z --untracked-files=no`; any
  tracked change outside `do-work/` refuses with `HEAVY-RUN-DIRTY-TREE`. Rename/copy source
  paths are checked too.
- `verifyLaneRevision:188–200` runs the dirty check **and** a HEAD compare before and after
  every lane, catching a lane that commits its own inputs.
- `InvalidateLaneEvidence:426–442` deletes the prior success **before** any execution attempt,
  so a failed, skipped, interrupted or unlaunchable run cannot leave an older green reusable;
  an inaccessible store refuses the whole run (`HEAVY-RUN-EVIDENCE-INVALIDATION`).
- Only a green, unskipped lane with a determinable fingerprint is ever stored
  (`recordSuccessfulLane:356–372`). A skip is detected from the `SKIP:` line prefix on the
  lane's own output, not a lane-name list (`laneSkipWatcher`).
- A reused lane reports `WallSeconds: 0` plus `EvidenceRevision` and `EvidenceRecordedAt`, "so
  a reader never mistakes an inherited duration for a measured one" (`reusedLaneExecution:338`).
- Interruption breaks the loop and leaves remaining lanes **absent from the run** rather than
  reported (`:177–181`).
- CLI surface: `run-heavy-verification [--manifest] [--lane …] [--lane-timeout-seconds]
  [--no-evidence-reuse]`; **reuse is on by default** (`heavy_commands.go:185`).

### 5.4 Direct consequences for this REQ

- **The green-gate record cannot be extended into an input-aware fast-gate cache without
  adding working-tree and toolchain inputs.** It has no timestamp, no dirty-tree awareness and
  no toolchain seal. Requirement 2 names uncommitted changes explicitly.
- **`RunLanes` cannot be the fast gate's execution path as it stands**: `refuseDirtyTrackedTree`
  refuses on any tracked modification outside `do-work/`, and the fast gate's whole purpose is
  to run on a dirty tree during development. What *is* reusable is the layer below it —
  `coverageRule`, `laneFingerprint`, `laneEvidenceStore`, `decideLaneReuse`, the disposition and
  reason constants, and the four-hour ceiling.
- **The fingerprint seals committed object ids, not working-tree bytes.** For a dirty-tree fast
  gate, a covered path that differs from HEAD must either have its working-tree bytes sealed
  (the machinery for that already exists in `fingerprintUntrackedFiles`, which hashes real file
  contents) or force execution.
- **Do not export `DO_WORK_TEST_DO_WORK_CLI_BINARY` from the fast gate.** The heavy fingerprint
  seals it explicitly (`heavy-runtime-fingerprint.py:86` as `supplied_cli`, and
  `heavy_maintainer_tree_test.go:94,105,133,165` pins it as a fingerprint-affecting setting).
  Setting it in the fast tier changes every heavy lane's environment seal.

---

## 6. Where a per-stage reuse decision could live

### 6.1 The insertion point

`run_verification()` in `_dev/tests/maintainer-verify.sh` is a flat sequence of shell blocks.
The smallest honest change is a helper — call it `run_stage_with_evidence <stage-id> -- <argv>`
— that asks one CLI command for a decision, executes on anything but a match, records on a
green exit, and prints one line per stage saying `executed`/`reused` and the reason. Backed by
a fast-stage manifest (`_dev/tests/fast-stages.json`) with the same schema as
`_dev/tests/heavy-lanes.json` and the same `coverageRule` matcher, so the "unclassified path
forces everything" rule comes for free.

Bootstrap constraint to name: the decision helper is `do-work-cli`, and `do-work-cli` is one of
the things under test. A warm `go build` of it is 0.33 s, so the overhead is acceptable, but a
red CLI must fail the gate rather than silently disable reuse.

Second insertion point, independent of the first: `probe-batch.sh` could consult the same
decision per probe, since it already owns launch, timing and status collection.

### 6.2 Smallest honest fingerprint per expensive stage

| Stage | Minimum honest coverage | Verdict |
|---|---|---|
| ShellCheck lint (3.7 s) | `shellcheck` binary bytes + `--version` output; the exact ordered list of tracked `*.sh` paths; the **working-tree** bytes of each (not the committed blob). | **Safe.** Closed, cheap (92 files), fully determinable. |
| gofmt (0.07 s) | `$GOROOT/bin/gofmt` bytes + the working-tree bytes of 274 tracked `*.go`. | Safe but **not worth it** — the fingerprint would cost more than the check. |
| `go vet` ×2 (0.2 s warm) | — | **Do nothing.** Go already caches vet results content-addressed. |
| queue-kanban fast tests (21.6 s) | Every file under `skills/do-work-board/tools/queue-kanban` (source, testdata, `go.mod`/`go.sum`), tracked-modified and untracked alike; the `heavy-runtime-fingerprint.py` seal for that module; the full environment (the `QUEUE_KANBAN_*` selectors decide which tests run at all); the gate scripts that build the selection (`maintainer-verify.sh`, `run-go-tests-with-budget.sh`). | **Safe**, and this is precisely "undo `-count=1` under a stricter fingerprint than Go's own test cache". Check first whether any test reads outside its module. |
| do-work-cli fast tests (51.6 s) | Same shape for `skills/do-work/tools/do-work-cli`, **plus** the `git` binary and Git configuration — these tests spawn real git constantly. `heavy-runtime-fingerprint.py <module> bash git python3` already seals exactly that. Plus `-short` in the argv seal. | **Safe.** Largest single win: half the gate's wall time, 85 % of its CPU. |
| Aggregate suite — narrow probes (`run-go-tests-with-budget-behavior.sh`, `do-work-cli-launcher-behavior.sh`, `p50-estimator-determinism.sh`, `select-simple-reqs-behavior.sh`, `defensive-surface-audit.sh`, `audit-lockins.sh`) | The specific scripts and shipped files each reads. | Declarable — but these are already the **cheap** ones (0.02–1.4 s each). Low value. |
| Aggregate suite — wide probes (`shipped-package-reference-contract.sh`, `action-shell-blocks.sh`, `prescribed-shell-canonicalization.sh`, `staged-skills-contract.sh`) | Effectively the whole tracked tree. `shipped-package-reference-contract.sh` resolves shipped markdown link targets against real repository paths (`:2092–2224`), and a link can point anywhere including `do-work/archive/`. | **Only safe under a whole-tree fingerprint.** A narrow declared coverage here would be a false green. |
| `session-start-hook-behavior.sh` (12.9 s) | The hook, the launcher, the whole CLI module, the Go toolchain, `git`. | Reusable in principle, but **fix the setup first** (§2). After that it is ~3.5 s and reuse buys little. |
| Owner contracts (`core-checks.sh` 5.25 s, `replace-text-section.sh` 8.21 s) | The shipped scripts each covers, plus the contract file itself. | Declarable; medium value. Moving these two into the parallel batch is a separate, cheaper win — but the REQ says added concurrency alone is not acceptance. |

### 6.3 "Unknown impact selects the broader verification", concretely

- A tracked path changed at HEAD, or dirty in the working tree, that **no stage rule
  classifies** → run **every** stage. (Existing precedent: `laneFingerprint:301–303` and
  `selectHeavyLanes` `plan.Uncertain`.)
- An untracked, non-ignored regular file a stage could read → include its bytes; if it is not
  a regular file → **error, execute**.
- A toolchain probe that fails, times out at 5 s, or reports "fingerprint uncertain" (opaque
  wrapper, `LD_PRELOAD`-family variable set, unknown Git config key, `GOFLAGS`/`GOWORK`, local
  module replace) → `fingerprint_uncertain`, **execute**.
- Evidence store unreadable, wrong permissions, or the record changed during the read →
  `evidence_unusable`, **execute**, and invalidate before running.
- Any environment variable added, removed or changed → mismatch → **execute**. (The whole
  environment is in the seal by design.)
- A skipped, failed or interrupted stage → **never** records evidence, and its prior record is
  deleted before the attempt. A partial gate run supplies no evidence for any stage.

### 6.4 Where reuse is NOT safe

- **Anything depending on unsealed host state.** The browser lane has no fingerprint
  declaration for stated reasons ("browser discovery, profiles, fonts, and shared runtime
  assets are not a complete, stable closure we can currently seal",
  `heavy-runtime-fingerprint.py:4–6`). It is heavy-only, so it does not affect the fast tier —
  but the same reasoning bars any future fast stage with that shape.
- **`_dev/tests/gate-runner.sh` itself** — it decides whether to run at all; caching its
  decision is caching the cache.
- **The wide markdown/reference contracts** under any coverage narrower than the whole tree
  (§6.2).
- **A stage that reads the real `do-work/` queue tree.** I checked all 15 fast-tier
  scripts: none does — every `do-work/queue|working|archive` reference in
  `session-start-hook-behavior.sh`, `select-simple-reqs-behavior.sh` and
  `contracts/core-checks.sh` is a fixture path under a `mktemp -d` root, and
  `shipped-package-reference-contract.sh`'s matches are in-test string literals. This must be
  **re-verified** if a stage is added, because `queueStatePrefix` deliberately excludes
  `do-work/` from both the dirty-tree refusal and the untracked seal.
- **Whole-gate reuse keyed on HEAD alone**, i.e. the current green-gate record. That is the
  false green requirement 2 exists to prevent.

---

## 7. How the gate's own tests are structured

Behavior tests for the gate scripts: `_dev/tests/run-go-tests-with-budget-behavior.sh` (57
lines), `_dev/tests/session-start-hook-behavior.sh` (142), plus
`_dev/tests/do-work-cli-launcher-behavior.sh` (129),
`_dev/tests/prescribed-shell-scripts-behavior.sh`, `_dev/tests/select-simple-reqs-behavior.sh`,
`_dev/tests/update-script-behavior.sh`, `_dev/tests/install-suite-behavior.sh`. The gate's own
stage list is covered separately by `maintainer-verify.sh --self-test`.

### The harness shape a new reuse-decision test should copy

1. Shebang `#!/usr/bin/env bash`, then `set -euo pipefail` (for a fail-fast probe) or
   `set -uo pipefail` with a `failure_count` tally (for a probe that reports several failures
   in one run — the session-start probe's shape).
2. `repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"`.
3. Allow the unit under test to be overridden by an env variable so the probe can be pointed at
   a mutated copy: `budget_runner="${DO_WORK_TEST_BUDGET_RUNNER:-$repo_root/_dev/tests/run-go-tests-with-budget.sh}"`.
4. `fixture_root="$(mktemp -d "${TMPDIR:-/tmp}/<name>.XXXXXX")"` + `trap 'rm -rf -- "$fixture_root"' EXIT`.
5. Build the smallest synthetic input that can exercise both the positive and the negative case.
6. One small `run_*` helper that sets the environment the real caller sets, so each case is one
   line.
7. Assert with explicit `FAIL:` lines on stderr naming the case, then exit non-zero.
8. Close with a single success line on stdout.
9. Register in `_dev/tests/contracts/probe-lanes.sh` with `register_probe <name> <path> '<failure message>'`.
   Note `register_probe` requires the script to be **executable** (`chmod +x`) or it counts a
   failure on the spot.
10. Budget: the new probe must finish in under 30 s in the fast tier
    (`probe-batch.sh:38–43`).

### Verbatim example case — `_dev/tests/run-go-tests-with-budget-behavior.sh:25–55`

```bash
run_fixture() {
  local run_id="$1"
  local test_pattern="$2"
  DO_WORK_TEST_DURATION_LOG="$duration_log" \
    DO_WORK_TEST_RUN_ID="$run_id" \
    DO_WORK_TEST_OTHER_GATE_PROCESSES=0 \
    DO_WORK_TEST_REPO_ROOT="$fixture_root" \
    DO_WORK_TEST_ENFORCE_BUDGET=yes \
    DO_WORK_TEST_FILE_BUDGET_SECONDS=30 \
    bash "$budget_runner" "$fixture_root" -run "$test_pattern" ./...
}

run_fixture passing-run '^TestPassingFile$' >/dev/null
failing_status=0
run_fixture failing-run '^TestFailingFile$' >/dev/null 2>&1 || failing_status=$?
if [ "$failing_status" -ne 1 ]; then
  printf 'FAIL: failing Go fixture returned %s; expected the original go test status 1.\n' \
    "$failing_status" >&2
  exit 1
fi

for run_id in passing-run failing-run; do
  if ! awk -F '\t' -v run_id="$run_id" '
    $1 == run_id && $2 == "budget_test.go" && $3 ~ /^[0-9]+([.][0-9]+)?$/ && $4 == "0" { matches++ }
    END { exit(matches == 1 ? 0 : 1) }
  ' "$duration_log"; then
    printf 'FAIL: %s did not log exactly one attributed budget_test.go duration row.\n' \
      "$run_id" >&2
    exit 1
  fi
done
```

That is the exact shape for the reuse-decision counter-cases: one helper, one line per case,
positive first, then each negative (relevant source changed, fixture changed, gate script
changed, runtime input changed, uncommitted change present, evidence expired, evidence
unreadable), each asserting both the **disposition** and the **reason code**.

For the "toolchain/runtime input changed" counter-case there is a ready-made pattern:
`do-work-cli-launcher-behavior.sh:28–56` writes a fake `go` on PATH driven by
`FAKE_GO_VERSION` / `FAKE_GO_BUILD_FAIL` / `FAKE_GO_TOOL_EXIT`, and
`maintainer-verify.sh:130–256` (`write_command_shim`) is a fuller shim covering `go`,
`gofmt`, `shellcheck`, `git`, `bash` and `node` behind one script, logging one line per stage.
On the Go side, `internal/heavyverification/heavy_reuse_regression_test.go` and
`heavy_evidence_test.go` already hold the equivalent table-driven reuse/counter-case tests for
lanes; a fast-stage decision built on the same functions belongs beside them.

---

## 8. Measurement facts

### 8.1 How `run-go-tests-with-budget.sh` measures per-file duration, and what it sums

`_dev/tests/run-go-tests-with-budget.sh:66–127` (embedded Python):

1. Builds `test_file_by_name` by walking `**/*_test.go` under the module and regexing
   `^func\s+(Test\w+)\s*\(` (`:81–84`). A test whose function name it cannot find is
   attributed to the literal bucket `"<unknown test file>"`.
2. Reads `go test -json` events and keeps `Elapsed` **only** for events where
   `Action in {"pass","fail"}`, `Test` is set, and **`"/" not in Test`** (`:92–93`) — subtests
   are excluded, so a parent's number is simply how long its subtests took end to end.
3. Sums those elapsed values per source file (`:95–98`).
4. Appends one row per file to `DO_WORK_TEST_DURATION_LOG`:
   `run_id \t <module-relative>/<file> \t %.2f \t other_gate_processes` (`:99–111`).
5. Fails when any file's sum is `>= budget_seconds`, but **only** when
   `DO_WORK_TEST_ENFORCE_BUDGET == "yes"` **and** `test_status == 0` (`:78`) — a red test run
   never also reports a budget breach.
6. Prints one summary line: `go-test budget: module=… wall=…s tests=N slowest-file=F:Xs limit=…`.

So the recorded number is **summed top-level-test elapsed time**, not process wall time and not
CPU. `wall=` in the summary line is separate, from `date +%s` around the `go test` call
(`:20,43`), and is not written to the log. Shell probes are timed differently: whole seconds
from `date +%s` in `probe-batch.sh:33–37` and `contract-regressions.sh:38–46`.

### 8.2 Is whole-gate wall time recorded anywhere?

**No.** `_dev/tests/gate-runner.sh:67,69` computes it and prints one line —
`gate <revision> green|red <seconds>s <log path> [recorded|…]` — to **stdout**, and
`skills/do-work/hooks/session-start.sh:30` launches the runner with
`nohup … >/dev/null 2>&1 &`. The per-revision file under `$TMPDIR/do-work-gate-runs/<rev>.log`
holds the gate's own output, which contains per-file and per-module lines but no total. No
gate-run directory existed on this machine at exploration time. The green-gate evidence record
carries no duration either. **Whole-gate wall time has to be added by this REQ**, which
requirement 5 asks for anyway.

### 8.3 Process-tree CPU on macOS

Obtainable, with one caveat. macOS has no cgroups, and `ps` cannot see a reaped process, but
`getrusage(RUSAGE_CHILDREN)` aggregates every waited-for descendant, and both of these expose
it:

- bash's `time` keyword with `TIMEFORMAT='%3R %3U %3S'` — verified working here: the
  do-work-cli lane reported `wall=51.555s user=66.543s sys=99.050s`, plainly aggregated across
  children on an 8-core machine.
- `/usr/bin/time -l <command>` — same rusage, plus max RSS and page-fault counters.

Caveat: `/usr/bin/time -p` under some shells under-reports grandchildren; a first measurement
of `session-start-hook-behavior.sh` gave `real 12.90 user 0.32 sys 1.35`, which is not the real
tree cost. Use the bash `time` keyword, and note that a detached process (the `nohup`'d
gate-runner) escapes the accounting entirely.

Recommended shape for the before/after comparison:
`bash -c 'TIMEFORMAT="%3R %3U %3S"; time bash _dev/tests/maintainer-verify.sh'` in a detached
worktree, with `GOMAXPROCS` and the probe-batch concurrency both fixed.

### 8.4 What `other_gate_processes` counts, and what it misses

`_dev/tests/maintainer-verify.sh:20–27`:
```bash
/bin/ps -Ao pid=,comm=,args= | /usr/bin/awk -v own_pid="$$" '
  $1 != own_pid && $2 ~ /(^|\/)(ba)?sh$/ && $0 ~ /\/maintainer-verify\.sh( |$)/ { count++ }
  END { print count + 0 }'
```
It counts processes whose `comm` is `sh` or `bash` (or a path ending in one) and whose full
argv contains `/maintainer-verify.sh` followed by a space or end-of-line, excluding its own
PID. The value is sampled **once at gate start**, exported as
`DO_WORK_TEST_OTHER_GATE_PROCESSES`, and stamped on every row of the run.

**It misses:**
- Every competing load that is not another `maintainer-verify.sh` — other repositories' gates,
  bare `go test` runs, agent sessions, editors, browsers, indexers. This is the dominant blind
  spot; my measurements ran at load average 2.5–4.9 with `other_gate_processes` == 0.
- Gates that **start after** this one begins — the single sample can never see them, yet they
  contend for the whole remaining run.
- The gate's own children (probes, `go test`, `git`) — excluded by design.
- A gate launched under a shell whose `comm` is not `sh`/`bash` (`zsh`, `dash`) or through a
  wrapper that does not put the path in argv.
- Anything that changes mid-run: CPU thermal state, memory pressure, a cold vs warm GOCACHE.

It also does not record `GOMAXPROCS`, the machine's core count, or the probe-batch fan-out, so
a row cannot be normalized after the fact. Requirement 5's "recorded toolchain/cache state and
no competing expensive gate" is not satisfiable from this column alone.

---

## 9. Primes — every entry bearing on this work

### 9.1 `_dev/primes/prime-shell-commands.md` (declared by the REQ)

Relevant to shell fixtures and the launcher path:

> **Cleanup and interruption are separate trap contracts.** Keep temporary-file cleanup on
> `EXIT`; a HUP/INT/TERM handler that only cleans and returns resumes the interrupted
> workflow, so terminating signals must exit with their conventional status and let `EXIT`
> perform cleanup. For background work the shell intends to kill later, follow
> [Testing → Background Work and Synthetic Load]…: the child owns its lifetime bound, because
> parent-side cleanup cannot survive SIGKILL.

> **Shell state does not survive between prescribed command blocks.** … a variable defined in
> one block — especially a `mktemp` random path — expands empty in the next … Blocks must
> re-derive what they need from deterministic paths and guard-check that inherited artifacts
> actually exist.

> **A checking tool that reports findings on stdout while exiting zero makes any
> exit-status-shaped lane decorative.** `gofmt -l` lists every unformatted file and still exits
> 0 … Read the *emptiness of the output* as the verdict for tools of this shape, and capture
> that output in a `local`-declared-then-assigned pair — `local name; name="$(...)"` — because
> `local name="$(...)"` takes `local`'s own exit status and masks a genuine crash from `set -e`.

> **Unchecked Exit Status Reads as Content.** A command or process substitution whose exit
> status is discarded while only its content is judged lets a tool that never ran read as a
> tool that found nothing. … *if this command were missing entirely, would the value it
> produces be distinguishable from a legitimate one?*
> **Directly load-bearing here:** a reuse decision computed from a fingerprint helper whose
> failure collapses to an empty string would read as "nothing changed" and produce a false
> green. This is the same defect shape requirement 2 forbids.

> **Closed Enumerations Go Stale.** When a rule applies "whenever X happens" … state the
> trigger *condition* in the rule's canonical home and mark any caller/value list as
> illustrative, not exhaustive.
> **Applies to:** a hand-maintained list of "stages eligible for reuse" or "inputs that matter"
> will go stale. The existing code already obeys this — the skip detector keys on the `SKIP:`
> prefix rather than a lane-name list, and the environment seal takes **all** of `os.Environ()`
> rather than a declared subset.

> **Every Flag on a Shipped Script Needs a Non-Test Caller.** Before adding an option to a
> script under `skills/*/tools/` or `skills/*/scripts/`, name the caller that is not the
> script's own test suite.
> **Applies to:** a `--no-reuse` / `--force` flag on anything shipped needs a real caller.
> `_dev/tests/` is maintainer-side and export-ignored, so a flag added there is exempt.

> **`git status --porcelain` collapses wholly-untracked directories** into a single `?? dir/`
> row … use `--untracked-files=all` or `git ls-files --others --exclude-standard`.
> **Applies to:** `refuseDirtyTrackedTree` uses `--untracked-files=no` and so is unaffected,
> but any new dirty-input check written for the fast gate must not enumerate untracked inputs
> through plain `--porcelain`. `fingerprintUntrackedFiles` already uses the `ls-files` form.

> **`git show --name-only` prints the commit header and message before the file list** … Use
> `git diff-tree --no-commit-id --name-only -r -m <commit>`.
> **Applies to:** `interveningCommitsAreGateLogs` already uses the correct form.

Also: "When a review finds a bug in prescribed-command logic, **grep the same primitive across
all actions before calling it fixed** — these patterns are usually copy-pasted, so the fix is
rarely local."

### 9.2 `_dev/primes/lessons-shell-commands.md`

Bearing on fixtures, migration parity and background lifetime:

> [family: legacy-fixture-implementation-shape] **REQ-420: fixtures that intercept shell
> internals preserve the implementation being retired; capture observable status, output,
> argv, filesystem, and Git effects before cutover.**
> **Directly load-bearing:** the shared-fixture change in §2 must keep asserting the same
> observable effects (status, exact stdout bytes, empty stderr, before/after file hashes),
> not merely that the same commands were reached.

> **REQ-423: EXIT cleanup and signal control flow are separate contracts** — a cleanup-only
> HUP/INT/TERM trap returns to the workflow, and the regression needs a valid success fallback
> present so that continuation cannot hide behind an already failing route.

> **REQ-325: a handle published one command after the launch is not a handle a trap can rely
> on** — traps run *between* commands, so `cmd & pid=$!` has a window the cleanup sees empty;
> and a file created before its first HUP/INT/TERM trap is a file no trap owns, because a
> default-action signal never runs the EXIT handler.
> **Applies to:** `probe-batch.sh:20–52` backgrounds probes under `set -m`; anything added
> there inherits this hazard.

> **REQ-186: give identical baseline child invocations one required owner.**
> **Applies directly to §2 and §3:** nine copies of one module and thirteen `git init` sites
> are identical baseline setups without a single owner.

> **REQ-187: keep one maintainer command inventory and close aggregate self-test recursion with
> a fixture-only mode.**
> **Applies to:** `MAINTAINER_VERIFY_SELFTEST_LOG` is that fixture-only mode
> (`maintainer-verify.sh:74–80` bypasses the budget runner entirely). A reuse layer must not
> break that seam, and `--self-test` currently asserts each of nine stages runs **exactly
> once** (`assert_success_stages:295–330`) — adding a reuse decision changes that count and
> that test has to be updated deliberately, not loosened.

> **REQ-193: lock complete shell-contract predicates so deletion or negation cannot survive a
> broad regex.**

> **REQ-243 / REQ-263: run every stated RED against pre-change code before writing anything** —
> "a 'RED' that was already red is the cheapest signal the work is half done."

> **REQ-217: only `git archive` honors `export-ignore`; a mirrored script must probe both
> sibling depths** — relevant because `_dev/` is export-ignored and a fast-stage manifest under
> `_dev/tests/` must never be cited by anything shipped.

> **REQ-216: macOS bash 3.2 makes empty-array expansion fatal under `set -u`; use `set --` for
> optional arguments.**

> **REQ-298: whenever a fallback value is in the same DOMAIN as a legitimate result, the
> failure has been laundered into data.**

### 9.3 Entries the REQ dropped for budget (recovered here)

From `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`:

> **[family: fixture-cost-is-subprocess-spawning] REQ-574:** when a test file overruns the
> gate's per-file duration budget, **measure before rewriting** — the cost here was spawning,
> not asserting. … Three traps sit on that path: `run-go-tests-with-budget.sh` sums `Elapsed`
> only for tests whose name has no `/`, so subtests never count; `git init --template=` omits
> `.git/hooks` along with the sample hooks, and a test that writes a real hook then fails on a
> missing directory; and `t.Setenv` cannot be traded for a per-child `Cmd.Env` when the code
> under test shells out itself, because the in-process call resolves its own binary from the
> process PATH. **The file that is slow is usually not the file to edit** — a package's fixture
> helper lives in one file and fixing it there speeds up every caller … Whole-module gains read
> smaller than per-package ones because `go test ./...` runs packages concurrently and hands
> freed CPU straight to a sibling.

> **[family: background-worker-self-bound] 2026-09-05:** background work must stop itself if
> its parent dies. During the review of REQ-581 six unbounded shell loops were added for a
> synthetic-load check. The parent died before its trailing `kill`; the incident report records
> roughly **4.5 hours of orphaned load and 24 CPU-hours consumed**. An `EXIT` trap still cannot
> handle SIGKILL. … Use the smallest justified load experiment and label its measurement window
> as loaded.

> **[family: shipped-module-test-self-containment] 0.294.2:** a test that reads the maintainer
> tree cannot ship. Two heavy-lane contract tests walked six levels up to
> `_dev/tests/heavy-lanes.json`; in an installed copy that walk lands on `<project>/.claude/`.
> … Every maintainer-tree read now sits in one file that `.gitattributes` export-ignores.
> **Directly load-bearing:** a `_dev/tests/fast-stages.json` read from Go must follow the same
> pattern, and the check is to simulate the install (copy `skills/` under a `.claude/` and run
> the module's tests), not to read the diff.

> **[family: smoke-vs-characterization] REQ-414 / REQ-415:** migration gates must compare
> retained and replacement statuses, ordered facts, actions, paths, and effects; a smoke matrix
> can stay green while destructive cleanup and inventory semantics diverge.

From `skills/do-work/crew-members/testing.md § Background Work and Synthetic Load` (cited by
the prime and by the REQ's Constraints):

> Whenever a shell backgrounds a process it intends to kill later, the background process must
> carry its own time bound. … A trailing `kill` or an `EXIT` trap is convenient early cleanup,
> never the only lifetime limit … **Synthetic load is an experiment, not a routine extra test.**
> Name the timing-sensitive failure it probes, choose the smallest worker count, duration and
> repeat count that answer it, and stop when that evidence is collected. Record those limits
> with the result. Do not overlap it with another session's gate or benchmark, and **do not use
> measurements from its load window as normal queue-performance evidence.**
>
> ```bash
> ( load_deadline=$((SECONDS + 120)); while (( SECONDS < load_deadline )); do :; done ) &
> ```

The REQ's Constraints already say synthetic load is unnecessary here. The applicable half is
the last sentence: my numbers above come from a loaded window and are labelled as such.

---

## 10. Recommended order of work (my reading, not a decision)

1. **R1 — share the skill root in `session-start-hook-behavior.sh`.** Smallest diff, largest
   aggregate-lane win, explicitly named by the REQ, zero assertion change. 12.9 s → ~3.5 s,
   aggregate lane ~24 s → ~19 s. This alone satisfies acceptance criterion 1.
2. **R2/R3 — apply the REQ-574 `TestMain`-template technique to `internal/lifecycleadvance`
   and `internal/heavyverification`.** Those are the two densest uncovered packages (13 and 5+
   `git init` sites, 41 s and 49 s of file time, zero `t.Parallel()`). Both fixture helpers
   live in one file per package, so the REQ-574 scope lesson applies: edit the helper, not the
   slow file.
3. **Then the reuse decision**, on the do-work-cli test stage first (51.6 s wall, 165 s CPU —
   half the gate and 85 % of its CPU), built on `coverageRule` + `laneFingerprint` +
   `laneEvidenceStore` + `decideLaneReuse`, with working-tree bytes added to the seal so it is
   honest on a dirty tree. Do **not** extend the green-gate record, which has no tree, toolchain
   or time input.
4. Add whole-gate wall and process-tree CPU recording (§8.2, §8.3) before the before/after
   comparison, and take the baseline in a detached worktree with a fixed worker limit.

Not recommended as acceptance evidence, though it is a real ~10 s win: moving the four serial
owner contracts into the parallel probe batch. Requirement 5 rules out concurrency alone as
acceptance, and it would make the aggregate lane's contention worse before reuse lands.
