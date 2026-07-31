---
id: REQ-065
title: "Confirm: update HANDOFF.md's commit-hash write-back guidance"
status: completed
status_changed_at: 2026-07-31T09:31:57Z
completed_at: 2026-07-31T09:35:42Z
created_at: 2026-07-30T22:01:38Z
user_request: UR-010
addendum_to: REQ-062
domain: general
discovered_during: REQ-062
---

# Confirm: update HANDOFF.md's commit-hash write-back guidance

## What

A `[low]` discovery from REQ-062's restatement sweep. `do-work/HANDOFF.md:35` currently tells this repo's future sessions:

> `do-work/` is git-excluded via `.git/info/exclude`: stage nothing under it; write `commit:` hashes directly into archived REQs; skip metadata commits.

That advice predates `tools/checks/record-commit-hash.sh`. "Write `commit:` hashes directly" now describes exactly the free-form edit the script was built to replace — and the script handles a git-excluded `do-work/` correctly on its own (it detects the untracked path, states which git-dependent guards it skipped, and still runs every content guard).

## Why this is a question and not just a fix

`do-work/HANDOFF.md` was outside REQ-062's declared Scope, and it is a **local maintainer note** for this repo — it is git-excluded and never ships to consumers. So the change is low-stakes but it is also yours to make: the file records how *you* want sessions in this repo to behave.

## What Would Change

One line in `do-work/HANDOFF.md`. Nothing shipped, no version bump, no changelog entry.

## Open Questions

- [x] Should HANDOFF.md tell this repo's sessions to use `tools/checks/record-commit-hash.sh` for the `commit:` write-back instead of editing archived REQs by hand? → Confirmed: Yes — point it at the script, keep the "skip metadata commits" part.
  Recommended: **Yes — point it at the script, keep the "skip metadata commits" part.** The script is safe on an untracked path and dogfooding it here is how regressions in it get caught early.
  Value: this repo stops being the one place that does the write-back the unguarded way.
  Risk: none meaningful — the script is read-only toward git and leaves the file untouched on any anomaly. Fully reversible (one line).
  Also: **Leave it as is** (the note is accurate about *staging*, and the hand edit has not caused a problem in this repo yet); or **delete the write-back clause entirely** and let the standard Step 9 procedure speak for itself.

  → **User chose the third option: delete the write-back clause entirely**, with an explanation left in its place so a future session doesn't reinvent a local variant.

## Implementation

Rewrote the single bullet at `do-work/HANDOFF.md:35`.

**Removed:** `write commit: hashes directly into archived REQs; skip metadata commits` — a local procedure competing with the shipped one.

**Kept:** the git-exclusion fact and `stage nothing under it` (an environment fact, not a procedure).

**Added, per the user's "explain it as well":**

1. An explicit statement that this note carries *no* local variant of the write-back, and why the old wording was wrong — it described the free-form edit `tools/checks/record-commit-hash.sh` was built to replace, the one that blanked six archived REQs in a consumer repo. Without the *why*, a future session reads the deletion as an oversight and helpfully restores it.
2. A pointer to the single source of truth: `actions/work-reference.md` → **Commit & Metadata-Commit Procedure (Step 9)**, named as the only instruction.
3. The two real consequences of git-exclusion, bounded as "only these two" so they don't grow into a local variant again: the script runs fine on an untracked path but prints `INFO: … is not tracked by git` and skips its size-floor and diff guards while running every content guard (`tools/checks/record-commit-hash.sh:290`); and there is no metadata commit to make, since `git commit` has nothing to stage.
4. A trap the old note didn't cover: `--verify` mode reports `FAIL: … is not tracked by git` here **by design** (`tools/checks/record-commit-hash.sh:99`). A session that follows Step 9 literally will hit that FAIL, and Step 9 says a `FAIL:` line means stop, not retry — so the note has to name this one as expected before it costs someone a debugging detour.

**Not done:** no version bump, no changelog entry, no commit. `do-work/` is git-excluded in this repo, so the edited file is untracked and there is nothing to stage — the same fact the bullet documents.

*Resolved via `do-work clarify`; user selected the third option after questioning the value of the `commit:` field itself.*

## Lessons Learned

- The `commit:` field is redundant by construction in any repo do-work has run in: Step 9 mandates commit titles of the form `[{id}] {title} (Route {route})`, so `git log --grep='\[REQ-NNN\]'` recovers the hash, and the grep survives history rewrites that invalidate a recorded hash. Its one non-redundant use is disambiguating the `--no-ff` merge commit from the builder's branch commit in worktree dispatch mode. Worth knowing before anyone invests further in that field — the guard script earns its keep by protecting the *file*, not by protecting the hash.
- When deleting a local instruction that competes with a shipped one, delete it *and* say why in place. A bare deletion is indistinguishable from an omission, and the next session restores it.
