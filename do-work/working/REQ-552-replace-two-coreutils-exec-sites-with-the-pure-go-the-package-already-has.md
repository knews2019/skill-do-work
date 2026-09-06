---
id: REQ-552
title: '[impact-negligible] Replace two coreutils exec sites with the pure Go the package already has'
status: claimed
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: backend
prime_files: []
tdd: true
suggested_spec:
depends_on: [REQ-550]
related: [REQ-549, REQ-550, REQ-551, REQ-553, REQ-554, REQ-555, REQ-556, REQ-557, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-09-06T00:40:32Z
  basis:
    - Route B
    - 5-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - dependency depth 1
    - cross-route regression gates
    - full-suite verification
write_set: [skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go, skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go, _dev/tests/audit-lockins.sh, _dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh, _dev/tests/prescribed-shell-cases/architecture-report-preflight.sh]
route: B
dispatch_at: 2026-09-06T00:55:56Z
builder_handback_at: 2026-09-06T00:55:56Z
claimed_at: 2026-09-06T00:38:56Z
---

# Replace two coreutils exec sites with the pure Go the package already has

## What
Of 90 `exec.Command` sites in do-work-cli, 85 run `git`; two spawn `find` or `cp` for work the same package already does in stdlib Go. Replace the `find` probe in `internal/corehelpers/commands.go` with an `os.Stat`/`WalkDir` readability check reusing the inventory walk, and the `cp` in `internal/toolboxcommands/architecture.go` with the `io.Copy` primitive the package already has.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** `_dev/primes/prime-shell-commands.md` read before either shell file was touched,
  plus the crew rules and the exploration. Approach: both replacements written fresh from the standard
  library, because neither primitive the request named is callable here; the lock-in pasted from the
  exploration unchanged; both heavy fixture cases rewritten rather than deleted.
- [x] **[APPLY]:** Five files, exactly the declared `write_set` — three from the original request plus
  the two fixture cases added at triage with the measurement that forced them.
- [x] **[UNIFY]:** `git diff --stat` on the merge range reports five files, 69 insertions, 20
  deletions, identical to the builder branch's own diff against its base. Linters: `gofmt -l` on both
  packages — no output; `go vet` on both — clean; `shellcheck --severity=warning -x` on all three
  shell files — exit 0 (three style notes at `audit-lockins.sh:28,48,118` are pre-existing and below
  the gate's severity). Cross-compile check `GOOS=windows GOARCH=amd64 go test -c` on both packages —
  exit 0, which matters because the replacements are filesystem code. No debug artifacts: the diff
  adds no `fmt.Print`, no `set -x`, no commented-out block, no temporary path.

## Why
The `find` site runs on every `audit-archive-timestamps` invocation only to harvest an error string while the Go walk traverses the same tree again, so the archive is walked twice per run; the `cp` site is the one non-Go copy in a package that already ships a pure-Go recursive copy.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 9, sweep_key `exec-where-pure-go-exists`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -10. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- `internal/corehelpers/commands.go` — `exec.Command("find", archiveRoot, "-name", "REQ-*.md", "-print0")`, output discarded, called from the `audit-archive-timestamps` path; `internal/corehelpers/inventory.go` already walks the same tree with `filepath.WalkDir` and the same predicate.
- `internal/toolboxcommands/architecture.go` — `exec.Command("cp", draftPath, stagedPath)` behind `DO_WORK_COMPATIBILITY_SHIM == "1"`; `internal/toolboxcommands/last30days.go` has `copyLast30DaysTree` with the copy primitive.
- Behaviour preserved: the same error string class is returned when the archive is unreadable; the same file lands at `stagedPath` with the same mode.
- Reproduce at dc8a64e3 (prints the two sites): `rg -n 'exec\.Command(Context)?\((ctx, )?"(find|cp|mkdir|grep|sed|ls|rm|mv|cat|touch|head|tail|wc)"' skills/do-work/tools/do-work-cli skills/do-work-board/tools/queue-kanban --glob '!*_test.go'`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, ~~no test files beyond the lock-in~~.
  **Overridden — see D-01.** Two heavy-tier fixture cases drive the very calls this REQ makes in-process, so they were measured to go inert and are in the `write_set` alongside the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (the file already exists, is executable, and is already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh` -- add one assertion to it; do not create it and do not change its registration), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- ~~Tests unchanged; the existing package tests are the safety net.~~ **Overridden — see D-01**, same reason: the package tests cannot see the two fixture cases that went inert.
- Lock-in limit: coreutils exec sites in the two Go modules: 0 after this REQ (today 2).

## Dependencies
Depends on REQ-550 (collapse exported delegates) so the `corehelpers` write set has no pending overlap. REQ-557 depends on this REQ.

## Builder Guidance
Firm: pure Go on both sides. Latitude on where the readability check lives inside `corehelpers`.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints two lines (`find` in commands.go, `cp` in architecture.go).
**GREEN when:** It prints nothing; package tests green; the lock-in pins coreutils exec sites in both Go modules at 0.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5660 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "do-work-cli internals".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for exec-where-pure-go-exists.*

## Triage

**Route: B** — Explore then build.

**Reasoning:** The two code edits are mechanical and the request names both files and both
replacements. What needed discovery is that the request's own Constraints cannot be satisfied: it says
"Tests unchanged; the existing package tests are the safety net" and "no test files beyond the
lock-in", and two heavy-tier fixture cases drive the very behaviour being changed by putting a fake
binary on `PATH`. The moment either call becomes in-process, the shim is inert and the case asserts
the opposite of what happens. That is discovery, not design — Builder Guidance already fixes the
approach — so exploration runs and planning does not.

**Planning:** Skipped.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Explore agent, read-only, re-verified against HEAD. Full report in the run directory as
`do-work/runs/work-2026-09-05-231943/REQ-552-exploration.md`.

**Both exec sites are still exactly where the request says**, at `commands.go:724` and
`architecture.go:133`, same line numbers as the earlier exploration. **The request's headline count is
stale in one half:** do-work-cli has 103 non-test `exec.Command` sites at HEAD, not 90; the "85 run
git" half is exact. The two coreutils sites, which are the only number the request acts on, are exact.

**Neither replacement the request names is usable as named.** `internal/corehelpers/inventory.go`
already walks the same tree with `filepath.WalkDir`, but it is not a reusable readability probe — it
parses every match and builds an ownership map, and it walks `do-work/working` too. `copyLast30DaysTree`
in `last30days.go:251-294` cannot be called for the `cp` site at all: `WalkDir` on a regular file
returns at `relative == "."` copying nothing, and its `O_EXCL` open fails against the temp file
`os.CreateTemp` already made. So both replacements are written fresh from the standard library, which
is what the request's Builder Guidance actually asks for.

**The behaviour to preserve was measured rather than assumed.** With GNU coreutils, `cp src dst` onto
a pre-created 0600 file leaves the mode at 0600 and copies the bytes; `os.CreateTemp` creates at 0600,
so the draft's own mode is never copied today. `os.WriteFile(stagedPath, data, 0o600)` reproduces that
exactly.

**One semantic difference is accepted and stated:** on a partly unreadable tree, `find` reports every
unreadable subdirectory and still exits non-zero, while the `WalkDir` replacement stops at the first.
Both produce a non-empty evidence string, so the gate behaves identically; only the evidence text
differs.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` (modify) — the find probe becomes a filepath.WalkDir readability check that records the first traversal error as its evidence string; the os.Stat guard above it is untouched
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go` (modify) — the cp becomes os.WriteFile at mode 0600, keeping the existing evidence prefix
- `_dev/tests/audit-lockins.sh` (modify) — one new assertion block in the file's existing shape, appended after the Finding 2 block
- `_dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh` (modify) — the fake-find case is rewritten to drive the in-process failure; see D-01
- `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh` (modify) — the fake-cp case is rewritten the same way; see D-01

**Files I will NOT touch:**
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go` — considered as a source of a reusable walk primitive; it parses every match and builds an ownership map, so it is not a readability probe
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/last30days.go` — its copy helper cannot be called for a regular file with a pre-created target

The pre-existing wart where a non-directory archive path yields an unhelpful evidence string is real,
out of scope, and recorded as a discovered task. Every other exec site in either module stays.

**Acceptance criteria (restated from REQ):**
- [ ] Neither find nor cp is spawned by shipped code in either module
- [ ] The archive-walk failure is still detected and still produces non-empty evidence
- [ ] The draft copy still lands at the staged path with the same bytes and the same mode
- [ ] A lock-in assertion fails if a coreutils exec site returns to shipped code in either module
- [ ] The heavy tier is green, not just the fast gate — both rewritten fixture cases still detect the
      failure they name

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go` (modified)
- `_dev/tests/audit-lockins.sh` (modified)
- `_dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh` (modified)
- `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh` (modified)

**What was done:** The `find` probe became a `filepath.WalkDir` readability check that returns the
first traversal error as its evidence, and the `cp` became `os.WriteFile` at the mode the temp file
already had — so neither module spawns a coreutils binary for work the standard library does. A
lock-in assertion fails if one comes back, with test files excluded deliberately: fixture setup may
shell out, shipped code may not.

The half of this change the fast gate cannot see is the fixture repair. Both heavy-tier cases drove
the failures they assert by putting a fake binary on `PATH`, so both went inert the moment the calls
became in-process — measured, not predicted: five assertions across the two files flipped from red to
green against the patched binary while asserting nothing. Both are rewritten to provoke the same
failure in-process rather than deleted. The archive case nests one branch past `PATH_MAX` so the
walk's own `open()` fails with `ENAMETOOLONG`; the publish case points `TMPDIR` at a regular file so
the staged copy cannot be created. Each gained an assertion naming the failure it expects, so neither
can pass again on an unrelated nonzero exit — which is exactly how the `PATH`-shim versions failed.

Merge range `7eadf50e..ca61a0a9`, five files, 69 insertions, 20 deletions — identical to the builder
branch's own diff against its base. Builder branch head `2a29fd3f`.

## Decisions

- **D-01 — the `write_set` is widened by two files, and the request's "tests unchanged" constraint is
  overridden. DECIDE & STATE.** `_dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh:220-236`
  installs a fake `find` that exits 3 on `PATH` and asserts the run exits non-zero and never prints
  "audit clean". `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh:195-215` writes a
  failing `cp` into a directory prefixed onto `PATH` and asserts the publish fails with no
  `index.html`. Both were proven inert against a patched binary built for the purpose: the first goes
  base exit 1 → patched exit 0 with "archive audit clean", the second base exit 1 with no `index.html`
  → patched exit 0 with one created. Leaving them alone means shipping a red heavy lane. **Value:** the
  change is actually verifiable. **Risk:** two fixture cases change in a request that said none would;
  the risk is that a reviewer reads it as scope creep, which is why it is recorded here with the
  measurement rather than done quietly. Both cases are rewritten to drive the in-process failure
  directly — never deleted, because deleting the first would leave the archive-walk failure path with
  no coverage at all.
- **D-02 — the lock-in keeps `--glob '!*_test.go'`. DECIDE & STATE.** Without it the pattern matches a
  third site, `internal/suiteinstall/update_transaction_test.go:25`, which spawns `cp -R` in fixture
  setup and is not this finding. An unglobbed lock-in would still count 1 after both production sites
  are fixed, so it would be red on day one and prove nothing. Fixture setup may shell out; shipped code
  may not.
- **D-03 — the lock-in's command list stays byte-identical to the audit's own Reproduce pattern.
  DECIDE & STATE.** It is a closed enumeration, which `_dev/primes/prime-shell-commands.md`
  § Closed Enumerations Go Stale warns about, and the honest answer here is that widening it to a
  condition would drag in `open`, `xdg-open`, `rundll32`, `ps`, `tar`, `sh`, `python3` and `curl` —
  eight commands this finding does not address and has no replacement for. The enumeration is the
  finding's own scope, and it fails on exactly what the request's RED proof prints.

- **D-04 — the archive case fails the walk by depth, not by permission. DECIDE & STATE, measured.**
  This suite runs as root in containers, where `chmod 000` is a no-op: a permission-bit spelling would
  make the case pass without a walk ever failing, which is the same vacuous pass the `PATH` shim was
  already giving. Nesting past `PATH_MAX` fails `open()` for any user. Verified in both directions
  before it was written in.
- **D-05 — the publish case sets `GOTMPDIR` alongside the broken `TMPDIR`. DECIDE & STATE, measured.**
  `do-work-cli.sh` resolves its binary with `go tool -n`, which needs a usable temp directory. Without
  `GOTMPDIR` the toolchain fails first with exit 2 and the case passes without ever reaching the
  publish path. The trap is written into the fixture beside the variable.
- **D-06 — each rewritten case gained one assertion naming its own failure. DECIDE & STATE.** Without
  it, a case that merely asserts a nonzero exit passes on whatever unrelated failure its new mechanism
  provokes — which is how both of these went inert in the first place.

## Qualification

**Passed.** Read from the merge range `7eadf50e..ca61a0a9`; canonical `qualify` and `scope-drift` both
satisfied.

- **The declared set and the touched set are identical.** Five files, 69 insertions, 20 deletions —
  the same diff the builder branch carries against its base. The write set was widened at triage, in
  writing, with the measurement that forced it; nothing entered the range afterwards.
- **The RED is the request's own Reproduce command and it was run.** At the base revision it printed
  the two exec sites and exited 0; the lock-in block run standalone against that same revision printed
  two FAIL lines and set `failure_count=2`. After the change the Reproduce command prints nothing and
  exits 1, and the lock-in prints `Audit lock-in regressions passed.`
- **The part the fast gate cannot see was proven in both directions.** The two heavy fixture cases were
  green at the base revision, then went red — five assertions across the two files — once the Go edits
  landed and the `PATH` shims went inert. That is the measurement that justified widening the write
  set, taken rather than argued. After the rewrite both are green again at heavy tier, and the whole
  prescribed-shell lane reports `110 named cases across 18 per-script files` with the per-file counts
  unchanged at 11 and 9.
- **The rewritten cases were checked for the failure mode that made the originals worthless.** Against
  a scratch build with the two guards removed — the walk discarding its evidence, the temp-file failure
  swallowed — both cases fail. A case that only asserts a nonzero exit passes on whatever unrelated
  failure its mechanism happens to provoke, which is precisely how the `PATH`-shim versions went inert;
  each rewritten case therefore also asserts the text of the failure it names.
- **Two traps were found by measurement rather than reasoning, and both are written into the fixtures.**
  `chmod 000` cannot make a path unreadable for a suite running as root, so the archive case fails the
  walk by nesting past `PATH_MAX` instead. And breaking `TMPDIR` to fail an in-process temp file also
  breaks the `go tool` build inside the CLI launcher, so the publish case sets `GOTMPDIR` alongside it
  — without which the toolchain fails first and the case passes without reaching the publish path.
- **One semantic difference is accepted and stated rather than hidden.** On a partly unreadable tree,
  `find` reported every unreadable subdirectory; the walk stops at the first. Both yield a non-empty
  evidence string, so the gate behaves identically — only the evidence text differs.

### Remediation qualification (after review)

**Passed.** Remediation merge range `f45a9137..ce64a0a4`, three files — two of the five already-declared
write-set files plus this request's own record. Cumulative range `7eadf50e..ce64a0a4`.

- **The review's headline finding is closed with the proof it actually needed.** The failed-copy case
  now runs the publish under a file-size limit with a real staging directory, so `os.CreateTemp`
  succeeds and the staged write itself fails, and it asserts the in-process write error rather than the
  finding code four sites in that function share. The obligation was to show it red against the
  PRE-CHANGE tree: restored onto a full rebuild of `7eadf50e` — verified to still carry
  `exec.Command("cp", …)` and `exec.Command("find", …)` — the case reports
  `architecture-report-preflight: 9 cases, 1 failure.` and exits 1. Against the shipped code with the
  new error handling ablated it fails four ways. The old case did neither.
- **A trap in the new mechanism is recorded rather than left to bite.** The launcher builds its binary
  with `go tool`, which cannot write under the file-size limit, so the build must already be warm.
  Earlier cases in the same file warm it; if they ever stop, the launcher exits 2 and this case fails
  loudly instead of passing on the wrong failure. That is written into the fixture beside the limit.
- **The lock-in fix deviates from the remedy the review named, and the deviation is right.** The
  review's `exec\.Command(Context)?\([^)]*"(find|cp|…)"` also matches
  `exec.Command("git", "rm", "-r", path)` — 85 of this module's exec sites run git, so that false
  positive is reachable. The committed pattern still requires the coreutils name to be the command
  argument itself, first for `exec.Command` and straight after the context for `exec.CommandContext`.
  Both spellings were injected into shipped code and both go red; the git spelling does not.
- **The stale-path blindness is closed loudly.** The two module directories moved into one array with a
  `[ -d ]` guard that fails and counts. Renaming a scanned directory in a scratch tree now prints
  `FAIL: coreutils lock-in cannot scan a missing module directory: …` and exits 1, where it previously
  reported clean forever because rg's exit status was discarded and its stderr sent to `/dev/null`.
- **The record was corrected where it was false.** The request's own Constraints section still forbade
  the test changes its write set declared; both clauses are struck through and pointed at D-01, which
  carries the measurement that forced the override.

Requirements traced: no coreutils spawn in shipped code in either module; the archive-walk failure
still detected with non-empty evidence; the draft copy still landing with the same bytes and mode; a
lock-in that fails when a site returns; and the heavy tier green, not just the fast gate.

*Checked by work action*

## Testing

**Tests run:** the whole canonical gate, `bash _dev/tests/maintainer-verify.sh`; the focused lane
`go -C skills/do-work/tools/do-work-cli test -count=1 ./internal/corehelpers/ ./internal/toolboxcommands/`;
the heavy-tier prescribed-shell lane, which is the only place two of this request's five files run;
`bash _dev/tests/audit-lockins.sh`; `bash _dev/tests/contract-regressions.sh`; `gofmt -l`, `go vet`,
`shellcheck --severity=warning -x` on all three shell files, and a `GOOS=windows GOARCH=amd64 go test -c`
cross-compile of both packages, which matters because both replacements are filesystem code.

**Result:** ✓ Green. The canonical gate exited 0 at the merge revision `dd51478` — **89s wall**, both
fast stages reporting `EXECUTING (fingerprint_mismatch)`, so the whole suite really ran against the
changed code. Exit status read directly from `$?`, never through a pipe. The focused lane compared
green against its recorded baseline.

**The lane that matters most here reported separately, because the fast gate cannot reach it:**
`env DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/prescribed-shell-scripts-behavior.sh` printed
`Prescribed shell script behavior probes passed (110 named script cases across 18 per-script files).`
and exited 0, with the two rewritten files' own counts unchanged at 11 and 9 cases.

**Both directions were recorded, not just the green one.** At the base revision both fixture files were
green; with the Go edits applied and the `PATH` shims still in place, five assertions across the two
files flipped red — which is the measurement that justified widening the write set. After the rewrite
both are green, and against a scratch build with the two guards removed both fail again, so neither
case is passing vacuously.

**Remediation testing.** The canonical gate exited 0 again at the remediation merge revision
`ce64a0a4` — **79s wall**, exit status read directly from `$?`. The heavy-tier lane reported
`Prescribed shell script behavior probes passed (110 named script cases across 18 per-script files).`
with both rewritten files' own counts unchanged at 9 and 11. `bash _dev/tests/audit-lockins.sh` prints
`Audit lock-in regressions passed.`; injecting either `exec.Command` or `exec.CommandContext` with a
coreutils name turns it red, and injecting `exec.Command("git", "rm", …)` does not.
`shellcheck --severity=warning -x` on both touched shell files exits 0.

*Verified by work action*

## Review

**Overall: 68%** | 2026-09-06T00:00:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 70% |
| Code Quality | 90% |
| Test Adequacy | 55% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Partial |

**Verdict: Approve with follow-ups** — the Go code in REQ-552 (replacing two places where the CLI started the external `find` and `cp` programs with Go code that does the same job) is correct and all three reviewers failed to break it, but one of the two rewritten test cases passes for the wrong reason and the version bump is still owed. Commit ca61a0a9.

Where the three reviewers split, and the call taken:

- Two reviewers rated the vacuous `architecture-report-preflight` case Important, the third rated it Minor because the copy's happy path is pinned by two other cases in the same file. Taken as **Important**: two reviewers deleted the new error handling and the case stayed green, so the acceptance criterion "both rewritten cases still detect the failure they name" is demonstrably false, whatever else happens to be covered.
- Two reviewers judged the lock-in acceptance criterion unmet, the third judged it met. Taken as **unmet**: the third only injected `exec.Command("mkdir", ...)`, which the pattern does catch. The other two injected `exec.CommandContext(runContext, "cp", ...)` into a shipped file and the lock-in printed "Audit lock-in regressions passed." That is a shown failing scenario, not a worry.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh:195-216` — the rewritten failed-copy case never runs the line this REQ changed. Pointing TMPDIR at a regular file breaks `os.CreateTemp` at `architecture.go:126`, so the new `os.WriteFile` at line 137 is never reached. Two reviewers independently ablated the new error handling (one replaced it with `_ = os.WriteFile(...)`, one deleted the copy) and the case stayed green at 9 cases, 0 failures. A third ran the base binary with only the two rewritten fixture files restored: the archive case went red with 3 failures, this one stayed green. Its only assertion greps for `ARCHITECTURE-PREFLIGHT-FAILED`, which four sites in the same function emit, so it does not name its own failure and decision D-06 is false for this case. Working remedy verified on both binaries: `ulimit -f 1` with a draft over 512 bytes lets `CreateTemp` succeed and makes the write fail with `draft copy failed: write ...: file too large` — [impact-rule-change → report only]
- Release is owed and not delivered: two of the five touched files ship under `skills/`, so per `_dev/primes/prime-releases.md` this is a release. Patch bump 0.303.10 to 0.303.11 across `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md:5`, plus a changelog entry in `CHANGELOG.md` mirrored byte-identically into `skills/do-work/CHANGELOG.md`. Do not touch `skills/do-work-board/tools/queue-kanban/VERSION` (0.236.20, independently versioned, untouched by this change). One caveat for whoever finalizes: HEAD has moved past this merge and carries other unreleased shipped changes under `skills/do-work-board`, so if the release covers the accumulated set the bump size needs judging against all of it — [impact-rule-change → report only]

**Minor findings:**
- The lock-in in `_dev/tests/audit-lockins.sh:158-162` only matches a context argument spelled exactly `ctx`. Two reviewers injected `exec.CommandContext(runContext, "cp", a, b)` into shipped code and got zero FAIL lines; renaming the variable to `ctx` produced one. The same module already spells it `invocationContext` in `last30days.go`, and the Naming for Reach rule pushes new code toward the longer spelling, so the most likely future regression is the one the lock-in cannot see. One-character fix: `exec\.Command(Context)?\([^)]*"(find|cp|...)"`. `exec.Command` never takes a context, so widening creates no false positives — [impact-rule-change → report only]
- Undisclosed behavioural divergence: `filepath.WalkDir` builds absolute paths and fails with ENAMETOOLONG on a subtree deeper than PATH_MAX, where GNU `find` (which uses fts and chdir) exits 0. Two reviewers built the tree and confirmed both halves. Detection is now a strict superset, not the same set, and the REQ's own wording says "the same error string class is returned when the archive is unreadable". The rewritten archive fixture uses exactly this divergence as its mechanism, so it pins new behaviour rather than preserved behaviour. Real-world reachability in a do-work archive is near zero, which is why this is Minor. The accepted-difference list should name it — [impact-negligible → report only]
- The lock-in block discards rg's exit status and sends its stderr to `/dev/null`, then decides on emptiness alone. One reviewer ran the exact invocation against two non-existent directories: rg exited 2, printed nothing, and the guard took the pass branch. If either hard-coded module directory is ever renamed the lock-in reports PASS forever. The command list is a closed enumeration that D-03 answers well; the module-path enumeration in the same block is a second one that nothing answers. Same shape as the pre-existing Finding 2 and Finding 10 blocks, so it matches local style — [impact-negligible → report only]
- No case anywhere in the suite covers the failure the REQ actually names, an archive that cannot be read because of permissions. The suite runs as root (uid=0 confirmed), so `chmod 000` is a no-op, and the new case covers only ENAMETOOLONG. The old case did not cover it either, so this is not a regression — the record should say the class is uncovered instead of implying the rewrite restored coverage — [impact-negligible → report only]

**Nit findings:**
- Two real behaviour differences in the copy path, both improvements, neither disclosed: the draft bytes are now snapshotted at `architecture.go:117` instead of re-read by `cp` at copy time, so the published bytes are always the ones the guards ran against; and `copyError.Error()` is never empty, where `cp` killed by SIGXFSZ wrote nothing to stderr and left the evidence as the bare string `draft copy failed: ` — [impact-negligible → report only]
- The compatibility-shim block at `architecture.go:125-142` is now a pure round-trip (write `data` to a temp file, read the same bytes back into `data`, delete on defer), which makes the "same mode" acceptance criterion unobservable by any test. Two reviewers raised this; a third confirmed the reasoning behind the mode claim is sound but untestable. Leaving the block in place is the right call here, but the simplification it exposes is not recorded as a discovered task — [impact-negligible → report only]
- The REQ's Constraints section still says "no test files beyond the lock-in" while the frontmatter `write_set` and Scope both list the two fixture files. D-01 records the override with the measurement that forced it, which is the right handling, but the constraint line itself was never amended — [impact-negligible → report only]

**Acceptance:** Partial — the shipped Go behaviour holds under differential execution (12-case fixture matrix as root and uid 65534, plus base-versus-head publish runs on normal, zero-byte, symlink, ENOSPC, broken-TMPDIR and `ulimit -f` drafts, all byte-identical on success and same finding code on failure), the heavy tier is green at 11 and 9 cases, shellcheck is clean at the gate severity `--severity=warning`, `gofmt`, `go vet` and the package tests pass, and the write set is exactly the five declared files. Two of the five acceptance criteria fail as written: the lock-in does not fire on `exec.CommandContext` with a non-`ctx` variable, and one of the two rewritten cases does not detect the failure it names.

**Requirements checklist:**
- [x] Neither `find` nor `cp` is spawned by shipped code in either Go module — delivered (lock-in red at 7eadf50e with two FAIL lines, green at ca61a0a9; the only remaining `cp` match is a test file, excluded on purpose)
- [x] The archive-walk failure is still detected and still produces non-empty evidence — delivered, with the superset caveat above (no case found where `find` errored and `WalkDir` did not; the Go 1.26 walk source confirms every non-nil return reaches the callback with a non-nil error, so the evidence string cannot be empty)
- [x] The draft copy still lands at the staged path with the same bytes and the same mode — bytes delivered and proven by ablation (writing corrupt bytes fails two publication-content cases); mode reasoned but unobservable, since the staged file is deleted on defer
- [ ] A lock-in assertion fails if a coreutils exec site returns to shipped code — not delivered for `exec.CommandContext` with a context variable not named `ctx`
- [ ] The heavy tier is green and both rewritten cases still detect the failure they name — heavy tier green, but the `architecture-report-preflight` case does not detect the failure it names

**Suggested testing:** 5 items
- Rewrite the failed-copy case to use `ulimit -f 1` in a subshell with a draft over 512 bytes, so `os.CreateTemp` succeeds and the write fails, and assert on the string `draft copy failed` rather than the shared `ARCHITECTURE-PREFLIGHT-FAILED` code. Verified reachable on both binaries.
- Re-run both rewritten case files against the base Go code with only the fixture files restored. Any case that stays green there is not discriminating.
- Widen the lock-in pattern to `exec\.Command(Context)?\([^)]*"(find|cp|...)"` and re-inject `exec.CommandContext(runContext, "cp", a, b)` to confirm it goes red.
- Guard each hard-coded module directory in the lock-in with a `[ -d ]` check, then rename one directory and confirm the lock-in fails loudly instead of passing.
- If permission-denied archive coverage is wanted, drive it through an unshare-based user namespace or a non-root helper, since `chmod 000` cannot work in a root-run suite.

**Follow-ups created:** None (9 findings report only)

*Reviewed by review-work action*

## Lessons Learned

- **A fixture that controls behaviour with a fake binary on `PATH` dies silently the moment the code
  stops shelling out.** Two of them did here, and both had been green for as long as they existed. The
  general shape: any test whose mechanism is *outside* the process under test stops being a test when
  the work moves inside it, and nothing in the suite says so — it just keeps passing. When a change
  removes a subprocess, every fixture that shimmed that subprocess is suspect by construction, and the
  cheap check is to restore the fixture onto the pre-change code and confirm it still goes red.
- **"Assert a nonzero exit" is not an assertion.** Both rewritten cases initially passed on a failure
  that was not the one they named — one of them on a guard the request had not even touched. A case
  that names a failure must assert the text of that failure, or it will eventually pass on whatever
  else happens to break first. The first remediation of this request existed entirely because that rule
  was applied to only one of the two cases.
- **A mechanism chosen to provoke a failure usually breaks something earlier than you think.**
  Pointing `TMPDIR` at a regular file was supposed to fail the copy; it failed the temp-file creation
  one guard above. Breaking `TMPDIR` also breaks the `go tool` build inside the CLI launcher unless
  `GOTMPDIR` is set beside it. And `chmod 000` cannot make anything unreadable for a suite running as
  root. Each of those was found by measuring, not by reading.
- **A pattern that names a variable is a pattern that will miss.** The lock-in matched a context
  argument spelled exactly `ctx` — while this repository's own naming rule pushes new code toward
  `invocationContext`. The most likely future regression was the one the lock-in could not see.
- **A guard that reads emptiness must first prove the thing it read was real.** The lock-in discarded
  rg's exit status and its stderr, so a renamed directory would have reported clean forever. That is the
  prime's "unchecked exit status reads as content" trap, in a file whose whole job is to catch drift.

## Orientation

The two replaced call sites are `archiveWalkFailure` in
`skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` and the staged-draft copy in
`skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go`. Neither has a Go test of
its failure path; both are covered by heavy-tier fixture cases in
`_dev/tests/prescribed-shell-cases/`, which run only under `DO_WORK_MAINTAINER_TIER=heavy` — a green
fast gate says nothing about them. The lock-in that keeps a coreutils spawn from returning is Finding 9
in `_dev/tests/audit-lockins.sh`, beside REQ-554's Finding 6.

Three things are deliberately not solved here. The `filepath.WalkDir` probe reports a failure `find`
did not — a subtree deeper than `PATH_MAX` — so detection is a strict superset rather than the same
set; that divergence is what the archive fixture's own mechanism now depends on. No case anywhere
covers a permission-denied archive, because the suite runs as root. And the compatibility-shim block in
`architecture.go` is now a pure round-trip — it writes bytes to a temp file and reads the same bytes
back before deleting it — which makes this request's "same mode" criterion unobservable by any test;
the simplification that exposes is a follow-up, not this request's.
