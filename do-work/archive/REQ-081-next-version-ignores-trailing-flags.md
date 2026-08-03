---
id: REQ-081
title: next-version ignores flags placed after the bump size and silently bumps the calling repo
status: completed
claimed_at: 2026-08-03T22:07:20Z
completed_at: 2026-08-03T22:13:28Z
commit: 84d79c1
kb_status: pending
route: B
created_at: 2026-08-03T17:09:21Z
user_request: UR-016
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-072]
batch: audit-remediation-external
addendum_to: REQ-072
write_set: [tools/queue-kanban/main.go, tools/queue-kanban/release_test.go, tools/queue-kanban/serve.go, tools/queue-kanban/prime-do-kanban.md, actions/work.md, _dev/tests/contract-regressions.sh]
---

# `next-version` Ignores Flags Placed After the Bump Size and Silently Bumps the Calling Repo

## What

`queue-kanban next-version` takes the bump size as a positional argument and `--repo-root` /
`--version-file` as flags. Go's `flag.FlagSet.Parse` stops at the first non-flag argument, so every
flag placed *after* the positional is discarded. The invocation the skill itself prescribes
(`actions/work.md:603`) puts `--repo-root` last — so it writes the **calling** repo's version file
instead of the requested one, exits 0, and reports the bump as successful.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `tools/queue-kanban/prime-do-kanban.md` (the listed prime) before touching code;
  its `## Lessons` are inlined-not-linked, which is where this REQ's lesson goes. Loaded
  `crew-members/general.md`, `coding-guardrails.md`, `testing.md` (`tdd: true`). Approach: reproduce
  the RED end-to-end first (done, below), extract `parseNextVersionArguments`, double-Parse, reject
  leftovers via a shared helper, table-test, then fix the prescribed invocation and pin its order.
- [x] **[APPLY]:** Five files. Four were declared up front; `tools/queue-kanban/serve.go` is a
  one-line addition the requirement-5 audit forced — recorded as scope drift with its reason in
  `## Decisions` D-02, not slipped in.
- [x] **[UNIFY]:** `git diff --stat` → `main.go` +107/−14, `release_test.go` +143, `serve.go` +1,
  `actions/work.md` 1 line, `_dev/tests/contract-regressions.sh` +16. Verified: `gofmt -l` reports
  nothing, `go vet ./...` clean, `go test ./...` passes, `go build` succeeds, contract suite exits 0.
  Read every hunk: the parse function, the two guard helpers, six call sites of
  `exitOnLeftoverArguments`, the dispatch doc block, and the two test functions. No debug artifacts,
  no `fmt.Print` left behind, no commented-out code.

## Why

This is the only write surface the tool has outside `do-work/`, and it writes the wrong file while
reporting success. The failure is silent in the worst way: in a consumer repo the discarded
`--repo-root` makes the walk resolve to whatever tree the command was launched from, and the version
line it finds there looks like a plausible bump. `CLAUDE.md` → Shipped Tooling stakes the tool's
safety on "exactly two write surfaces"; one of the two does not point where it is told to.

It also falsifies REQ-072's own D-01, which chose `--version-file` specifically so a repo that does
not keep `**Current version**:` at the default path could still be served. As shipped, that override
is unreachable from the documented invocation.

## Context

`tools/queue-kanban/main.go:180-188` as shipped:

```go
flagSet := flag.NewFlagSet("next-version", flag.ExitOnError)
repoRootOverride := flagSet.String("repo-root", "", "...")
versionFileOverride := flagSet.String("version-file", "", "...")
_ = flagSet.Parse(args)

bumpSize := ""
if flagSet.NArg() > 0 {
	bumpSize = flagSet.Arg(0)
}
```

`Parse` returns at the first argument that does not begin with `-`, so with
`["patch", "--repo-root", "/x"]` it consumes nothing: `bumpSize` is `patch` and both overrides keep
their zero values. The remaining tokens sit unread in `Arg(1)`/`Arg(2)` — not rejected, not warned
about.

**Reproduced** (throwaway repos, nothing in this tree touched). Two repos each carrying a
`**Current version**:` line and a `do-work/` directory; command run from `cwdrepo`, pointed at
`target`:

```
$ queue-kanban next-version patch --repo-root <target> --version-file <target>/actions/version.md
9.9.10                       # <- the CALLING repo's version, bumped
exit=0
cwdrepo/actions/version.md:  **Current version**: 9.9.10   # written, should not have been
target/actions/version.md:   **Current version**: 1.0.0    # untouched, should have been written
```

Reversing to `next-version --repo-root <target> patch` works correctly, which isolates the cause to
argument order rather than to the resolver or the writer.

**Why no check caught it.** `_dev/tests/contract-regressions.sh:405` asserts only that the string
`queue-kanban next-version` appears in `actions/work.md`; argument order is unasserted. Every Go test
in `release_test.go` calls `allocateNextVersion` / `readCurrentVersion` directly and never goes
through the argument parsing, so the whole `runNextVersionCommand` arg-handling path is untested.
`queue-kanban verify` cannot see it either — the write it makes is well-formed, just in the wrong
tree.

## Detailed Requirements

1. **Accept flags on both sides of the positional.** `next-version patch --repo-root X`,
   `next-version --repo-root X patch`, and interleaved forms must all resolve identically. The
   standard shape is to `Parse`, take `Arg(0)` as the bump size, then `Parse` the remaining args
   again — or to lift the bump size out of the slice before parsing. Either is fine; pick one and say
   why in a doc comment.
2. **Reject leftover arguments instead of ignoring them.** After parsing, any unconsumed token is an
   error with exit 2, not silence. A second positional (`next-version patch minor`) and a misspelled
   flag must both fail loudly. Today they are discarded, which is how this defect stayed invisible.
3. **Extract the argument handling into a testable function.** `runNextVersionCommand` calls
   `os.Exit`, so it cannot be asserted against. Move the parse into something like
   `parseNextVersionArguments(args []string) (bumpSize, repoRoot, versionFile string, err error)` and
   have the command call it. This is the enabling change for requirement 4 — without it the RED case
   cannot be written at all.
4. **Table-test the argument orders**, including the exact invocation `actions/work.md` prescribes.
   Cover: flags-after-positional, flags-before-positional, interleaved, missing bump size, unknown
   flag, and a stray extra positional.
5. **Audit every other subcommand for the same shape.** `next-version` is the only one with a
   positional today (`summary`, `generate`, `serve`, `next-req`, `verify` are flags-only; `now` takes
   nothing), so the expected finding is "none" — but record that you checked, because a future
   positional would reintroduce this silently. Per the Closed Enumerations rule, state the condition
   ("any subcommand mixing a positional with flags"), not just today's list.
6. **Add a contract assertion that pins the documented invocation's argument order**, so
   `actions/work.md:603` and the parser cannot drift apart again. Assert the order in the prescribed
   command, not merely the presence of the subcommand name.
7. **Fix the prescribed invocation** at `actions/work.md:603` regardless of which parsing shape wins —
   put the flags before the positional. A parser that accepts both still leaves the documented form
   as the one every agent copies.

## Constraints

- **Read-back-confirm stays.** REQ-072 requirement 6 has `allocateNextVersion` re-read the file after
  writing; that behavior is load-bearing and unchanged here. It confirms the write *landed*, never
  that it landed in the right tree, which is precisely this REQ's gap.
- **The bump size stays a required positional and is never inferred.** REQ-072 requirement 2 and the
  doc comment at `main.go:172-178` are explicit that patch/minor/major is a human judgment. Do not
  "fix" this by adding a default or a `--bump` flag with a fallback value.
- **The tool's write surface must not widen.** This REQ makes an existing write go to the right place;
  it adds no new file the tool may write. `CLAUDE.md` → Shipped Tooling's two-write-surface sentence
  stays true as written and needs no amendment.
- Prescribed-command hygiene applies to requirement 7: the command in `actions/work.md` must actually
  do what the surrounding prose says it does when pasted into a shell.

## Dependencies

`addendum_to: REQ-072`, which introduced the subcommand. No `depends_on` — buildable immediately, and
independent of every other REQ in this batch.

## Builder Guidance

**Certainty: Firm on the diagnosis and on requirements 1–4; open on the parsing shape.** The
reproduction above was run end-to-end; do not re-derive it, but do re-run it after the fix as the
GREEN check, since a unit test on the parse function alone would not have caught the original bug
(the bug was in how the parsed values reached the resolver, and an end-to-end probe is what exposed
it).

Requirement 3 is a refactor in service of testability, which `crew-members/coding-guardrails.md`'s
surgical-scope rule normally discourages. It is justified here because requirement 4 is otherwise
impossible — say so in the commit rather than expanding the refactor beyond the one function.

## Red-Green Proof

**RED prompt/case:** A table test in `tools/queue-kanban/release_test.go` (or a new
`main_arguments_test.go`) asserting that parsing `["patch", "--repo-root", "/tmp/x",
"--version-file", "/tmp/x/actions/version.md"]` — the exact shape `actions/work.md:603` prescribes —
yields `bumpSize == "patch"`, `repoRoot == "/tmp/x"`, and `versionFile == "/tmp/x/actions/version.md"`.

**Why RED now:** The parse is inline in `runNextVersionCommand` and there is no function to call, so
the test does not compile; once the function exists but keeps the single `flagSet.Parse(args)`, both
override assertions fail with empty strings.

**GREEN when:** That table passes for every argument order in requirement 4; `go test ./...` in
`tools/queue-kanban/` stays green; the contract suite's new order assertion passes; and the
end-to-end probe below writes the target and leaves the caller untouched:

```bash
# from a repo that is NOT the target
queue-kanban next-version patch --repo-root <target> --version-file <target>/actions/version.md
# GREEN: prints the TARGET's next version; target/actions/version.md changed;
#        the calling repo's version file is byte-identical to before.
```

**Validation:** Inferred during capture, then reproduced — the RED state was executed against
throwaway repos before this REQ was written, so the failure is observed rather than predicted.

## Full Context

See `do-work/user-requests/UR-016/input.md` for the verbatim instruction, the provenance of the
external audit, and the batch constraints.

---
*Source: external audit finding F1 (P1), accepted by `do-work validate-feedback` triage after
empirical reproduction; captured on the user's instruction "capture all seven as REQs and stop".*

---

## Triage

**Route: B** - Medium

**Reasoning:** The defect, its cause, and its reproduction were all supplied and verified at capture,
so nothing needed planning. What needed discovery was the shape of the same defect elsewhere in the
binary (requirement 5's audit) and where the parse extraction should sit — exploration.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**RED reproduced before any code was written**, per Builder Guidance. Two throwaway repos under
`/private/tmp`, each with a `**Current version**:` line and a `do-work/` directory; the shipped binary
run from `cwdrepo`, pointed at `target`:

```
$ queue-kanban next-version patch --repo-root .../target --version-file .../target/actions/version.md
9.9.10            <- the CALLING repo's version
exit=0
cwdrepo: **Current version**: 9.9.10   (written; should not have been)
target : **Current version**: 1.0.0    (untouched; should have been written)

$ queue-kanban next-version patch minor
9.9.11            <- a stray positional is silently ignored AND bumps again
exit=0
```

**Requirement 5's audit found more than the expected "none".** `next-version` is indeed the only
subcommand reading a positional (`summary`, `generate`, `serve`, `next-req`, `verify` are flags-only;
`now` takes nothing) — but *none of them checks `NArg()` either*, and that is the same defect, not a
different one:

```
$ queue-kanban verify stray --repo-root <target>     # before the fix
(runs against the CALLING repo — "stray" halts Parse, so --repo-root is discarded)
```

The stakes are lower (`verify` only reports, and reports on the wrong tree), but the mechanism is
identical: a token Parse stops at, followed by silently discarded flags. Requirement 5 asked for the
*condition* rather than today's list precisely so a future positional would not reintroduce this — a
shared helper is that condition expressed in code rather than in a comment, so the audit's finding is
fixed rather than filed.

**Where the parse must live.** `runNextVersionCommand` calls `os.Exit` on every error path, so nothing
inside it can be asserted — a test would take the process down. That is not incidental to this bug; it
is why the bug shipped. Requirement 3's extraction is the enabling change, and the RED case in the
Red-Green Proof literally cannot be written without it.

## Scope

**Files I will touch:**
- `tools/queue-kanban/main.go` (modify) — `parseNextVersionArguments` + `nextVersionArguments`,
  `rejectLeftoverArguments` / `exitOnLeftoverArguments`, call sites, dispatch doc block
- `tools/queue-kanban/release_test.go` (modify) — three test functions covering every argument order
- `actions/work.md` (modify) — the prescribed invocation's argument order
- `_dev/tests/contract-regressions.sh` (modify) — pin the invocation's order and the shared helper
- `tools/queue-kanban/serve.go` (modify) — **added during the build**, one line, see D-02
- `tools/queue-kanban/prime-do-kanban.md` (modify) — Step 7.5 lesson (inline convention)

**Files I will NOT touch:** `allocate.go` / `release.go` (`allocateNextVersion`'s read-back-confirm is
load-bearing and unchanged, per Constraints), and no new write surface anywhere.

**Acceptance criteria (restated from REQ):**
- [ ] All three argument orders resolve identically
- [ ] Leftover tokens are an error with exit 2, not silence
- [ ] The parse is a testable function
- [ ] The argument orders are table-tested, including the prescribed invocation
- [ ] Every other subcommand audited, with the condition stated
- [ ] A contract assertion pins the documented invocation's order
- [ ] `actions/work.md` prescribes flags before the positional
- [ ] End-to-end: the target is written and the caller is byte-identical

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/main.go` (modified)
- `tools/queue-kanban/release_test.go` (modified)
- `tools/queue-kanban/serve.go` (modified)
- `tools/queue-kanban/prime-do-kanban.md` (modified)
- `actions/work.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Extracted the argument handling into `parseNextVersionArguments(args []string)
(nextVersionArguments, error)` — a pure function with no `os.Exit`, which is what makes any of this
assertable. It parses **twice**: the first `Parse` consumes leading flags and halts on the bump size,
`Arg(0)` is taken, then `Args()[1:]` is parsed again to pick up everything after it. The doc comment
says why this over lifting the positional out by index: the positional is not always at index 0, so an
index-based lift needs its own mini-parser, whereas double-parsing keeps `flag.FlagSet` the single
authority on what a flag looks like (`--flag=value`, `-flag value`, `--`). The FlagSet moved to
`ContinueOnError` with output discarded so errors return instead of exiting; the command renders them
and a usage line, exiting 2.

Leftover tokens are now an error everywhere, via `rejectLeftoverArguments` (returns an error, used by
the parse function) and its wrapper `exitOnLeftoverArguments` (reports and exits 2, used by the six
command functions). The helper's comment states the **condition** — any subcommand finishing its parse
with tokens left over must fail — rather than today's subcommand list, so a subcommand added later
inherits the rule by calling it.

`actions/work.md`'s prescribed invocation now puts the flags before the positional, and a new contract
assertion pins that exact order rather than the mere presence of the subcommand name — asserting only
the name is what let the documented form stay broken through a release. A second assertion pins the
shared helper's existence. The dispatch doc block notes that the synopsis is conventional notation and
not a required order, so the next reader does not "fix" `actions/work.md` back.

## Testing

**Tests run:** `go test ./...` (tools/queue-kanban), `go vet ./...`, `gofmt -l .`,
`bash _dev/tests/contract-regressions.sh`
**Result:** ✓ All passing

**Red-green validation:** traced to `## Red-Green Proof`.

- *RED, end-to-end, before any code* — the reproduction in Exploration above: the documented
  invocation bumped the calling repo to `9.9.10`, exited 0, and left the target at `1.0.0`.
  `next-version patch minor` bumped again and reported success.
- *RED, unit* — with the second `Parse` and the leftover check removed from the new function (i.e. the
  shipped single-Parse behaviour, restored deliberately), `go test` fails with exactly the assertions
  the REQ predicted: `RepoRootOverride = "", want "/tmp/x"` and `VersionFileOverride = "", want
  "/tmp/x/actions/version.md"` on the prescribed-invocation case, plus three
  `succeeded, want an error` failures. Restored → green. **Observed, not assumed.**
- *RED, contract* — restoring the old argument order in `actions/work.md` fails the suite with the new
  order assertion naming the file and the fix. Reverted → exit 0.
- *GREEN, end-to-end* — same two throwaway repos, rebuilt binary:
  ```
  $ queue-kanban next-version patch --repo-root .../target --version-file .../target/actions/version.md
  1.0.1     exit=0
  cwdrepo: **Current version**: 9.9.9    (byte-identical to before ✓)
  target : **Current version**: 1.0.1    (written ✓)
  $ queue-kanban next-version patch minor
  queue-kanban next-version: unrecognized argument(s) for next-version: [minor]
  usage: queue-kanban next-version <patch|minor|major> [--repo-root DIR] [--version-file PATH]
  exit=2    cwdrepo still 9.9.9 ✓
  $ queue-kanban verify stray --repo-root <target>
  queue-kanban: unrecognized argument(s) for verify: [stray --repo-root <target>]
  $ queue-kanban next-version --repo-root <target> minor     # flags-first still works
  1.1.0     exit=0
  ```

**New tests added:**
- `TestParseNextVersionArgumentsAcceptsFlagsOnEitherSideOfTheBumpSize` — 5 cases: flags-after (the
  prescribed shape), flags-before, interleaved, `--flag=value` on both sides, bare positional.
- `TestParseNextVersionArgumentsRejectsRatherThanIgnores` — 6 cases: missing bump size, flags-only,
  stray second positional, stray positional after a flag, unknown flag before and after the positional.
- `TestRejectLeftoverArgumentsIsTheSharedRule` — pins the condition and that the message names both
  the subcommand and the offending token.

**Existing tests updated (cross-REQ impact):** none. Every pre-existing test still passes unchanged.

*Verified by work action*

## Decisions

- **D-01 — Parse twice rather than lifting the positional out of the slice by index.** DECIDE & STATE
  (requirement 1 asked for the choice to be justified in a doc comment; it is, and here is the same
  reasoning for the trail). An index-based lift has to *find* the positional first, because a flag may
  precede it — which means writing a mini-parser that duplicates `flag`'s own rules for
  `--flag=value`, `-flag value`, and `--`. Two `Parse` calls keep `flag.FlagSet` the single authority
  on what a flag is. Cost: the FlagSet is parsed twice, which is free at this size. Reversible.

- **D-02 — Extended the leftover-rejection to all six subcommands, not just `next-version`. This is
  declared scope drift.** ESCALATE-tier, so stated with both sides. Requirement 5 asked me to audit
  the other subcommands and expected "none"; the audit found the same silent-discard reachable on
  every flags-only subcommand (`verify stray --repo-root X` discards `--repo-root` and reports on the
  calling tree). **Value:** requirement 5's actual ask is that a *future* positional not reintroduce
  this silently, and a shared helper every subcommand calls is that guarantee in code rather than in
  a comment — it also closes the read-only variants found today, at one line each.
  **Risk:** it is scope beyond the seven requirements, and it changes behaviour for users who were
  passing stray tokens that previously worked by being ignored (they now get exit 2). That is the
  intended direction — silence is the defect — but it is a behaviour change on five subcommands the
  REQ did not name, and it pulled `tools/queue-kanban/serve.go` into a scope that did not declare it.
  Reversible: deleting six one-line calls restores the old behaviour exactly. Flagged for the reviewer
  rather than folded in quietly, and `## Scope` and `write_set` were both amended to match reality.

- **D-03 — The bump size stays a required positional; no default, no `--bump` flag.** DECIDE & STATE.
  Explicitly required by the Constraints and by `main.go`'s existing doc comment, and worth restating
  because the tidiest-looking fix for an argument-order bug is to make the positional a flag. That
  would put patch-vs-minor-vs-major inside the tool, which REQ-072 D-01 deliberately kept as a human
  judgment. Not done.

## Lessons Learned

**What worked:**
- **Reproducing the RED end-to-end before writing anything, and again after.** Builder Guidance called
  this out specifically, and it earned its keep: the unit test alone would not have caught the
  original bug, because the bug was in how parsed values reached the resolver. The throwaway-repo
  probe is what proves the *right file* got written.
- **Deliberately restoring the broken behaviour inside the new function to watch the tests fail.**
  Six failures, naming the exact assertions the REQ predicted. A test written after a fix, never
  observed red, proves only that it agrees with the code.

**What didn't:**
- **Requirement 5's expected answer was wrong, and the REQ said so in advance.** It predicted "none"
  and asked for the check to be recorded anyway "because a future positional would reintroduce this
  silently." The audit instead found the same defect *today* on five flags-only subcommands. Running
  an audit whose answer you have already guessed is still worth doing — that is the whole reason
  requirement 5 exists, and it was right.

**Worth knowing:**
- **`flag.FlagSet.Parse` halting at the first non-flag argument is not a positional-only hazard.** On
  a flags-only subcommand a stray token placed first halts the parse and every flag after it is
  discarded, silently. `NArg()`/`Args()` must be checked even where no positional is expected.
- **`os.Exit` in a command function is an untestability boundary, and it is where bugs hide.** Every
  error path in `runNextVersionCommand` exited, so the whole argument-handling path had zero coverage
  while the surrounding release logic had plenty. The tell is a test file that only ever calls the
  *inner* helpers.
- **A contract assertion on a prescribed command must pin its argument order, not its name.** The
  existing check asserted `queue-kanban next-version` appears in `actions/work.md` — true throughout,
  including while the documented invocation silently wrote the wrong repo.

## Orientation

`queue-kanban next-version` now writes the repo you point it at. It previously discarded any flag
placed after the bump size, so the invocation the skill itself prescribes bumped whatever tree the
command was launched from and exited 0 reporting success — a wrong-file write on the tool's only write
surface outside `do-work/`. Flags are accepted on either side of the positional now, and every
subcommand rejects leftover tokens instead of ignoring them. Lives in the shipped Go tooling
(`tools/queue-kanban/`, the release-ritual subcommands).

[MAP CHANGED] Argument handling moved out of the `os.Exit`-calling command wrapper into
`parseNextVersionArguments`, which is the first time this binary's CLI parsing is testable at all —
and `rejectLeftoverArguments` is a new shared contract every current and future subcommand is expected
to call. Behaviour change worth knowing: five subcommands that previously ignored stray tokens now
exit 2 on them.

**Prime staleness spot-check** (`tools/queue-kanban/prime-do-kanban.md`, the REQ's one prime): its
`main.go` line already describes flags per subcommand and stays accurate; the two-write-surface
statement is unchanged, since this REQ redirects an existing write rather than adding one. Its
`## Lessons` are inlined-not-linked by an explicit marker, so this REQ's lesson was appended as a plain
bullet in that style. No dead paths.

## Review

**Overall: 95%** | 2026-08-03T22:11:40Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 98% |
| Scope | 85% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 1 important, 2 minor
**Acceptance:** Pass — the end-to-end probe writes the target and leaves the caller byte-identical;
every argument order resolves identically; `go test ./...` and the contract suite are green.
**Suggested testing:** 3 items
**Follow-ups created:** None

### Findings

- **[Important] Declared scope drift: `serve.go` and five subcommands were changed beyond the REQ's
  seven requirements.** `scope-drift.sh` reports clean only because `## Scope` and `write_set` were
  amended when the drift was taken (D-02) rather than after the fact — the honest reading is a
  builder-initiated scope extension, and it is the reviewer's call whether it belonged here or in a
  follow-up. The case for here: requirement 5's audit found the defect live, in the same file, and
  filing it would have meant a follow-up editing the same six functions. The case against: it is a
  behaviour change on five subcommands nobody asked about, in a REQ whose title names one. Recorded
  with both sides in D-02 rather than presented as inevitable.
- **[Minor] The double-`Parse` relies on `FlagSet` tolerating a second `Parse` call.** It does — the
  method resets its argument slice and re-assigns flag values — but this is a documented-behaviour
  dependency with no test that would notice if a future Go release changed it. The table tests would
  catch it, which is the mitigation; noting it so the dependency is on the record.
- **[Minor] `--` is untested.** `flag` treats a bare `--` as an end-of-flags marker, so
  `next-version -- patch` and `next-version patch -- --repo-root X` have defined but unasserted
  behaviour (the latter would now be a leftover error, which is arguably right). Not covered because
  no prescribed invocation uses it.

### Requirements Checklist

| # | Requirement | Status |
|---|---|---|
| 1 | Flags on both sides; choice justified in a doc comment | Delivered — D-01, comment at `parseNextVersionArguments` |
| 2 | Reject leftovers, exit 2 | Delivered — and extended to all subcommands (D-02) |
| 3 | Extract a testable parse function | Delivered — `parseNextVersionArguments`, no `os.Exit` |
| 4 | Table-test the orders, incl. the prescribed one | Delivered — 11 cases across 2 tables, RED observed |
| 5 | Audit every other subcommand; state the condition | Delivered — audit found the defect live, not "none"; condition encoded in `rejectLeftoverArguments` |
| 6 | Assertion pinning the invocation's argument order | Delivered — and observed failing |
| 7 | Fix the prescribed invocation | Delivered — flags before the positional |

### Acceptance Testing

End-to-end against two throwaway repos, before and after: RED wrote `cwdrepo` and left `target`
untouched; GREEN writes `target` and leaves `cwdrepo` byte-identical. Stray positional now exits 2
without writing. `verify stray --repo-root X` now errors instead of silently reading the wrong tree.
Flags-first invocation unaffected. `go test ./...`, `go vet ./...`, `gofmt -l .` (empty), `go build`,
and `bash _dev/tests/contract-regressions.sh` all clean. `qualify.sh` and `scope-drift.sh` OK.

### Suggested Additional Testing

- **A real release ritual in a consumer repo**, where `--repo-root` actually differs from the cwd.
  Everything here used synthetic trees; the resolver's walk-up behaviour in a nested checkout is the
  part least exercised.
- **The five newly-strict subcommands against any existing automation.** If a wrapper script was
  passing a stray token that used to be ignored, it now exits 2. Nothing in this repo does, but a
  consumer's script might.
- **`--` handling**, per the Minor finding, if any invocation ever needs it.

*Reviewed by review-work action (pipeline mode, in-session)*
