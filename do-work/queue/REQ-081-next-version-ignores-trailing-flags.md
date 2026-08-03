---
id: REQ-081
title: next-version ignores flags placed after the bump size and silently bumps the calling repo
status: pending
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
write_set: [tools/queue-kanban/main.go, tools/queue-kanban/release_test.go, actions/work.md, _dev/tests/contract-regressions.sh]
---

# `next-version` Ignores Flags Placed After the Bump Size and Silently Bumps the Calling Repo

## What

`queue-kanban next-version` takes the bump size as a positional argument and `--repo-root` /
`--version-file` as flags. Go's `flag.FlagSet.Parse` stops at the first non-flag argument, so every
flag placed *after* the positional is discarded. The invocation the skill itself prescribes
(`actions/work.md:603`) puts `--repo-root` last — so it writes the **calling** repo's version file
instead of the requested one, exits 0, and reports the bump as successful.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
