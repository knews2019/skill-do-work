---
id: REQ-603
status: pending
domain: general
created_at: 2026-09-06T08:19:05Z
user_request: UR-105
review_generated: true
impact: impact-user-visible
effort_estimate: effort-moderate
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-shell-commands.md]
tdd: true
depends_on: [REQ-597]
related: [REQ-597, REQ-601]
write_set: [skills/do-work/scripts/protected-inventory.sh, skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go, skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go, skills/do-work/actions/commit.md, skills/do-work-toolbox/actions/inspect.md, skills/do-work/docs/prescribed-shell-primitives.md]
title: 'Let the protected-inventory launcher pass global flags, and stop its shim discarding the text it prepared'
---

# Let the Protected-Inventory Launcher Pass Global Flags, and Stop Its Shim Discarding the Text It Prepared

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

Two defects in one wrapper, both measured by REQ-597's builders while correcting the prose that
describes it (evidence and fixtures in `do-work/runs/work-2026-09-05-231943/REQ-597-handback.md`).

**The launcher cannot pass a global flag.** `skills/do-work/scripts/protected-inventory.sh:6` is
`exec bash .../do-work-cli.sh --format text protected-inventory "$@"`: everything a caller passes lands
after the command token, so `--repo-root` is rejected as `unknown option` and the runtime takes the
current directory as the root (`commandruntime/command_runtime.go:103-107`). Run from a subdirectory,
`start` prints the same rows but `associate` exits 2 with nothing on either stream. Neither
`commit.md` nor `inspect.md` can run the wrapper from anywhere but the project root, and both now say
so because REQ-597 had no better answer in prose.

**The compatibility shim discards the text it prepared.** `internal/corehelpers/inventory.go:445-456`
replaces the result text unconditionally, so the `NO-DO-WORK-DIR` line prepared at `:205`, the
`PARSE-FAILED` line at `:212` and a walk error's `HELPER-USAGE` finding at `:215` never reach a
caller. Measured silent exit 2 for: a missing `do-work/` directory, an unmatched backtick in an
Implementation Summary line, and an unreadable REQ file. Both callers read every exit 2 as "skip REQ
tracing", which is the reading that let `inspect.md`'s broken blocks ship unnoticed.

## Why

A silent exit 2 that means "skip" is indistinguishable from a silent exit 2 that means "the tool could
not run", and the callers have been treating the second as the first. The wrapper prepares the text that
would tell them apart and then throws it away.

## Detailed Requirements

- The launcher passes global flags through. `tools/checks/associate-files.sh:10-17` shows the
  translating shape; after the change `protected-inventory.sh --repo-root <root> start` (or the flag
  after the mode, whichever the CLI's own convention is; say which) works from any directory. State the
  chosen convention in the guide's launcher sentence.
- The shim keeps what the handlers prepared: `NO-DO-WORK-DIR`, `PARSE-FAILED` and a walk error's
  finding reach stdout or stderr as the handler intended. One test per swallowed outcome, red before
  the change, asserting the text is present.
- Then the callers: `commit.md` and `inspect.md` stop reading a printed finding as the skip condition.
  Exit 2 with a finding is an error to report; a genuine "nothing to associate against" is whatever the
  handler prints for it. Rewrite the exit-2 sentences in both files from the new behaviour, measured.
- Three neighbouring prose defects in the same sentences, found on the same pass: `commit.md:67` tells
  a re-run to append to the retained quarantine, but `start` replaces it (`inventory.go:393`) and only
  `associate` unions (`:420-425`), so say which command a re-run uses; `commit.md:61`'s "exit 2 means
  this is not a git repo" also covers a `git status` failure (`:368`) and a quarantine write failure
  (`:394`); `associate` after `start --dry-run` exits 2 as not-started (`:412-415`, `:392-396`),
  which neither action states.
- The guide's sentences REQ-597 wrote about the current-directory dependence are updated to the new
  behaviour in the same commit; a sentence describing the old launcher must not outlive it.

## Constraints

- Every prose sentence is written from a measured run, not from the code alone (REQ-596's rule for this
  file class). Keep the fixtures REQ-597 left under the scratchpad if they still exist; rebuild them if not.
- Shipped Go and shipped shell: a release.

## Open Questions

None.
