---
title: "Lessons from REQ-451: Make confirmation input interruptible"
type: source-summary
topic_cluster: queue-orchestration-and-lifecycle
sources: [raw/processed/2026-09-04/REQ-451-make-confirmation-input-interruptible.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-451: Make confirmation input interruptible

Part of the [[concept-queue-task-lifecycle]] cluster.

## What the REQ was about

Make install and update confirmation input cancellation-aware so `SIGINT`, `SIGHUP`, and `SIGTERM` cannot leave the process waiting forever at the prompt. Before writes begin, signal handling must exit with the documented signal status without waiting on recovery that the blocked input path itself prevents.

## Solution summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go` (modified)
- `_dev/tests/install-suite-behavior.sh` (modified)

## What worked

- Separating pre-write confirmation cancellation from the existing write-started recovery owner kept the fix small and made real-process signal tests decisive.
- **What did not work:** The first canonical verification was held by date-bound fixtures outside this request. Repairing those tests in a standalone commit restored the gate without widening REQ-451.

## Worth knowing

- A buffered result channel is necessary because a generic `io.Reader` may complete after the caller returns; post-write recovery remains owned by the unchanged transaction boundary.

## Back-reference

See `do-work/archive/REQ-451-make-confirmation-input-interruptible.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `21036776`.
