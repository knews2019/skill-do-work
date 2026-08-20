# Version Action

> **Part of the do-work skill.** Handles version reporting, update checks, and work recaps. User-facing walkthrough: [`docs/version-guide.md`](../docs/version-guide.md).

**Current version**: 0.219.0

**Upstream**: https://raw.githubusercontent.com/knews2019/skill-do-work/main/skills/do-work/actions/version.md

## When to Use

**Use when:**
- The user asks "what version", "release notes", "what's new", or "history" → version + last 5 changelog entries.
- The user asks to "update", "check for updates", or asks whether there is a newer version → shared update engine.
- The user asks for a "recap" of recent work → recap flow across archive + active URs.

**Do NOT use when:**
- The user wants every changelog entry — point them at `CHANGELOG.md`.
- The user wants a fresh install — use the README installation path.
- The installed skill is outside the current project — the shared engine refuses global/shared updates.

## Input

The user's phrasing selects one mode:

- **Version request** — "what version", "version", "what's new", "release notes", "what's changed", "updates", "history".
- **Update check** — "update", "check for updates", "is there a newer version".
- **Recap** — "recap" (dispatched with `mode: recap`).

## Responding to Version Requests

1. Report the version shown above.
2. Read only the first ~80 lines of `CHANGELOG.md` at the skill root.
3. Extract the five newest `## ` release blocks, reverse them, and print them with the newest at the bottom.

## Responding to Update Checks

All update discovery, review, confirmation, mutation, byte verification, and recovery is owned by `tools/do-work-update.sh`. This action and `just run-do-work-update` must execute that same engine; do not duplicate its archive or overwrite logic here.

1. Resolve `<skill-root>` as the directory containing this action's parent `SKILL.md`.
2. Resolve `<project-root>` from the invocation directory with the repository fallback:
   ```bash
   PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
   ```
3. Execute the engine yourself and relay its output:
   ```bash
   bash "<skill-root>/tools/do-work-update.sh" --project-root "$PROJECT_ROOT"
   ```
4. Do not download a second archive, perform a second diff, add another confirmation, or fall back to a direct `curl | tar` mutation. If the engine refuses or fails, report that exact failure; it is the authoritative safety boundary.

The engine already tries two routes before giving up: the anonymous tarball over HTTP, then a shallow git clone repacked with `git archive`. Its failure message names the outcome of both. When it reports that neither route reached the host — a corporate proxy, a blocked domain, a rate limiter that outlasts the retries — the supported escape hatch is `DO_WORK_UPSTREAM_URL`, which points both the updater and the installer at a different archive URL for that invocation:

```bash
DO_WORK_UPSTREAM_URL=https://example.internal/do-work/main.tar.gz bash "<skill-root>/tools/do-work-update.sh" --project-root "$PROJECT_ROOT"
```

Relay that suggestion when the engine reports a total fetch failure. Editing a vendored file is never the answer.

The installed manifest validator approves all four modules first, then the updater delegates that same downloaded archive to its installed full-suite installer. The installer owns the reviewed confirmation, refreshes the managed Just section, composes core hook settings, verifies every installed byte, and restores every changed managed path on failure. The update does not mutate `do-work/`, `kb/`, application files, bytes outside the marked Just section, unrelated settings entries, or any other project configuration.

## Responding to Recap Requests

When the user asks for a recap:

1. **Archive source** (`do-work/archive/UR-*/`): read the title from `input.md` and REQs from the `REQ-*.md` files inside each UR folder.
2. **Active source** (`do-work/user-requests/UR-*/`): read the title from `input.md`. Scan `do-work/queue/REQ-*.md` and `do-work/working/REQ-*.md` for matching `user_request:` frontmatter.
3. **Merge**: combine both lists, deduplicate by UR id (archive wins), sort by UR number descending, and take the newest five.
4. **Label each UR**:
   - no label when fully archived;
   - `(pending)` when any REQ is pending or claimed;
   - `(completed, awaiting archive)` when every REQ is terminally resolved but the UR is still active.
5. **Format**:
   ```text
   ## Recent Work

   UR-144 — Block-level improved translation for ZH pairs
     REQ-361 — Block-level improved translation
   UR-143 — Model selector thinking variants (completed, awaiting archive)
     REQ-360 — Model selector thinking variants
   ```
   Use one line per UR and one indented line per REQ. Include no descriptions, scores, or file lists.
6. If no archive and no active URs exist, print `No completed work yet.` and omit the section.

## Red Flags

- The update path is about to run any mutation other than `tools/do-work-update.sh --project-root …` — stop; the two entry points have drifted.
- The engine reports the skill outside the project or the project is not a Git worktree — stop rather than improvising a global or unrecoverable update.
- The engine reports a malformed suite archive, unsafe manifest, older upstream version, or failed recovery — report the failure; do not bypass it.
- Recap lists the same UR from active and archive sources — keep only the archive version.
- The version reported differs from the `**Current version**:` line above — re-read this file from disk.

## Verification Checklist

- [ ] Version output includes the local version and five releases with the newest last.
- [ ] Update mode executed `<skill-root>/tools/do-work-update.sh --project-root <project-root>` rather than reproducing update steps.
- [ ] No second download, diff, confirmation, or extraction occurred in the action path.
- [ ] Global/shared installs and non-Git projects were refused by the engine.
- [ ] Recap merged, deduplicated, sorted, and labeled active/archive URs correctly.
