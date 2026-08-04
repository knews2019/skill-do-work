---
source_type: req_lesson
req_id: REQ-081
req_path: do-work/archive/UR-016/REQ-081-next-version-ignores-trailing-flags.md
date: 2026-08-03
domain: general
module: tools/queue-kanban
tags: [queue-kanban, next-version, ignores, flags, after]
---

# Lessons from REQ-081: next-version ignores flags placed after the bump size and silently bumps the calling repo

## What the REQ was about

`queue-kanban next-version` takes the bump size as a positional argument and `--repo-root` /
`--version-file` as flags. Go's `flag.FlagSet.Parse` stops at the first non-flag argument, so every
flag placed *after* the positional is discarded. The invocation the skill itself prescribes
(`actions/work.md:603`) puts `--repo-root` last — so it writes the **calling** repo's version file
instead of the requested one, exits 0, and reports the bump as successful.

## Solution summary

Extracted the argument handling into `parseNextVersionArguments(args []string) (nextVersionArguments, error)` — a pure function with no `os.Exit`, which is what makes any of this assertable. It parses **twice**: the first `Parse` consumes leading flags and halts on the bump size, `Arg(0)` is taken, then `Args()[1:]` is parsed again to pick up everything after it. The doc comment says why this over lifting the positional out by index: the positional is not always at index 0, so an index-based lift needs its own mini-parser, whereas double-parsing keeps `flag.FlagSet` the single authority on what a flag looks like (`--flag=value`, `-flag value`, `--`). The FlagSet moved to `ContinueOnError` with output discarded so errors return instead of exiting; the command renders them and a usage line, exiting 2.

## What worked

- **Reproducing the RED end-to-end before writing anything, and again after.** Builder Guidance called this out specifically, and it earned its keep: the unit test alone would not have caught the original bug, because the bug was in how parsed values reached the resolver. The throwaway-repo probe is what proves the *right file* got written.
- **Deliberately restoring the broken behaviour inside the new function to watch the tests fail.** Six failures, naming the exact assertions the REQ predicted. A test written after a fix, never observed red, proves only that it agrees with the code.

## What didn't work

- **Requirement 5's expected answer was wrong, and the REQ said so in advance.** It predicted "none" and asked for the check to be recorded anyway "because a future positional would reintroduce this silently." The audit instead found the same defect *today* on five flags-only subcommands. Running an audit whose answer you have already guessed is still worth doing — that is the whole reason requirement 5 exists, and it was right.

## Worth knowing

- **`flag.FlagSet.Parse` halting at the first non-flag argument is not a positional-only hazard.** On a flags-only subcommand a stray token placed first halts the parse and every flag after it is discarded, silently. `NArg()`/`Args()` must be checked even where no positional is expected.
- **`os.Exit` in a command function is an untestability boundary, and it is where bugs hide.** Every error path in `runNextVersionCommand` exited, so the whole argument-handling path had zero coverage while the surrounding release logic had plenty. The tell is a test file that only ever calls the *inner* helpers.
- **A contract assertion on a prescribed command must pin its argument order, not its name.** The existing check asserted `queue-kanban next-version` appears in `actions/work.md` — true throughout, including while the documented invocation silently wrote the wrong repo.

## Back-reference

See `do-work/archive/UR-016/REQ-081-next-version-ignores-trailing-flags.md` for the full REQ — triage, implementation, review, and lessons. Commit `84d79c1`.
