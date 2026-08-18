---
id: REQ-229
title: Verify the published path in the download and screenshot helpers
status: pending
domain: general
created_at: 2026-08-18T00:18:44Z
user_request: UR-042
addendum_to: REQ-225
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
maintenance: false
write_set:
- skills/do-work/scripts/atomic-download.sh
- skills/do-work/scripts/capture-screenshot.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
---

# Discovered Task: Verify the Published Path in the Download and Screenshot Helpers

## What

`skills/do-work/scripts/atomic-download.sh` and `skills/do-work/scripts/capture-screenshot.sh` publish their payload and return without checking the path they actually wrote. When the destination path is occupied by a directory, `mv` and `ln` nest the payload inside it and exit zero, so both helpers report a success that did not happen. Close the same defect class already closed in four sibling helpers.

## Context

Found while implementing REQ-225, which states the rule once in the shipped guide as `## Verified exact publication`. Writing the condition down is what surfaced these two — the same defect had been found four times before, every time by a review sweep and never by reading the guide, which is the argument REQ-225 was captured on. Both instances below were reproduced, not inferred.

**`atomic-download.sh:44`** runs `mv "$download_path" "$target_path"`. With the target occupied by a directory the download nests as `<target>/<target-basename>.download.XXXXXX`, line 48 then clears `download_path` so the EXIT trap spares it, and the script exits 0. Reproduced against a `file://` source: exit 0, target still a directory, private file abandoned inside it. Callers include the suite installer's `SKILL.md` downloads and `tools/fetch-upstream-archive.sh`, both of which read the exit status as proof the file landed.

**`capture-screenshot.sh:33`** runs `ln "$copy_path" "$destination_path"`. `ln` refuses an occupied *file* destination, which is where the no-clobber guarantee comes from, but nests on an occupied *directory* and exits zero. Under `--staged` this compounds into data loss: the success path continues to `rm "$source_path"` at line 45 and destroys the staged screenshot — the only copy the capture dispatch holds — while the destination never receives it. Reproduced: exit 0, staged source deleted, destination still a directory holding an orphaned `.copying.XXXXXX` file.

## Requirements

- Both helpers verify the path they actually wrote after publishing, and fail closed when the write nested instead of publishing.
- Each helper removes only its own nested artifact and leaves the occupying directory exactly as it was, matching the four helpers that already do this (`publish-portfolio-summary.sh:102-163`, `generate-report-image.sh:112-117`, `generate-report-image-batch.sh:169-173`, `install-last30days.sh:98-103`).
- `capture-screenshot.sh` must not delete the staged source when publication did not happen — the staged source is preserved on every other failure path in that script and this one is the exception.
- Do not weaken either helper's existing guarantees: `atomic-download.sh` keeps rename-on-success and failure preservation, `capture-screenshot.sh` keeps byte verification and the no-clobber refusal on an occupied file.
- Delete the temporary sentence in `skills/do-work/docs/prescribed-shell-primitives.md` § Atomic download publication that records these two helpers as not yet making the check.

## Red-Green Proof

**RED prompt/case:** Two new cases in `_dev/tests/prescribed-shell-scripts-behavior.sh`: (1) invoke `atomic-download.sh` with a `file://` source and a target path occupied by a directory; (2) invoke `capture-screenshot.sh --staged` with a destination path occupied by a directory.
**Why RED now:** Both exit 0 today. Case 1 leaves the target a directory with a `.download.XXXXXX` file inside it. Case 2 additionally deletes the staged source, so the screenshot is gone and was never published.
**GREEN when:** Both cases exit nonzero, the occupying directory is unchanged and holds no nested private artifact, and case 2's staged source is still present.
**Validation:** Reproduced by hand during REQ-225; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] Auto-approved: critical severity (data-loss risk on `capture-screenshot.sh --staged`, silent false success in both). → Added to queue immediately.
