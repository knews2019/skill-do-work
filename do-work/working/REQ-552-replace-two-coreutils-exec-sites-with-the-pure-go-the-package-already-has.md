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
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (the file already exists, is executable, and is already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh` -- add one assertion to it; do not create it and do not change its registration), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- Tests unchanged; the existing package tests are the safety net.
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

*Verified by work action*
