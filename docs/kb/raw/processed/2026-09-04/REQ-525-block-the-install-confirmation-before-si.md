---
source_type: req_lesson
req_id: REQ-525
req_path: do-work/archive/REQ-525-block-the-install-confirmation-before-signalling-it.md
date: 2026-09-03
domain: testing
module: skills/do-work/tools/do-work-cli
tags: [testing, block, install, confirmation]
---

# Lessons from REQ-525: Block the install confirmation before signalling it

## What the REQ was about

Make `TestBuiltInstallAndUpdateExit130WhenSignalsInterruptBlockedConfirmation` wait until the installer is actually blocked on its confirmation prompt before delivering the signal. Today it signals on a timer, so under parallel load the installer can finish rendering its managed-install diff and exit before the signal is handled, and the assertion reads `exit = <nil>, want 130`.

## Solution summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go` (modified)

## What worked

- Making the diagnosis a separate, write-nothing pass with explicit permission to stop and report a product defect. The builder used it: it refuted the captured premise *and* my own hypothesis, then declined to touch the test. Had it been told simply to fix the flake, the honest move — reporting that the test was right — would have looked like failing the task.
- `GOMAXPROCS=1` on the child, and running the control in both directions. The pinned neutered tree fails 9 of 9 every run; the unpinned one fails 0 of 9 across three. Without that second half, the pin would have looked like a preference instead of the thing that makes these lock-ins.
- Letting the builder reject the review's proposed test. That shape provably does not reproduce the bug, and finding out cost one build rather than a shipped decorative test.

## What didn't work

- My capture. I filed this as `effort-mechanical` test synchronization and wrote that the test "signals on a timer"; it polls for the prompt text, and the defect was a production signal-handling bug that made `Ctrl-C` report success. The triage was wrong in kind, not in size.
- My hypothesis. I proposed the window between writing the prompt bytes and entering the read. A 500 ms delay past that point still fails, with the read parked in its select — refuted, not merely unconfirmed.
- The first fix, which took 130 unconditionally and so reported exit 130 with empty stdout for a complete verified install. One exit owner was necessary and not sufficient: the owner also has to know what the work concluded.
- Four sightings across three subtests before I pulled this forward. The evidence was in the captured stderr of every one of them — `Install this complete four-skill suite? [y/N] Installation cancelled; no files were changed.` — sitting in my own gate logs, saying plainly that main had taken the declined branch.

## Worth knowing

- `_dev/tests/install-suite-behavior.sh:650` already asserted exit 130 from this path, inside the canonical gate, and was latently flaky the same way. A decision that looks like a judgment call is worth grepping for first: the contract may already be written down.
- `exec.CommandContext`'s `watchCtx` injects `context.Canceled` whenever `Cancel()` succeeds, so cancelling before the parent reaps a subprocess converts a successful validation into a failure. That is why signalling from inside a stub reproduces the mid-write case rather than the post-verify one.
- `exec` puts an `*os.File` child descriptor back into blocking mode via `Fd()`, so a narration pipe's buffer size has to be measured on a throwaway pipe rather than assumed.
- An interrupted `update-suite` still leaks its extracted upstream tree, because `os.Exit` in a defer skips the caller's `RemoveAll`. Scratch under `TMPDIR`, no repository effect, and it resolves with the `ExitCodeOverride` plumbing.

## Back-reference

See `do-work/archive/REQ-525-block-the-install-confirmation-before-signalling-it.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `8c06caa`.
