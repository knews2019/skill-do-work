# Version Action

> **Part of the do-work skill.** Handles version reporting, update checks, and work recaps. User-facing walkthrough: [`docs/version-guide.md`](../docs/version-guide.md).

**Current version**: 0.107.0

**Upstream**: https://raw.githubusercontent.com/knews2019/skill-do-work/main/actions/version.md

## When to Use

**Use when:**
- The user asks "what version", "release notes", "what's new", or "history" → version + last 5 changelog entries.
- The user asks to "update", "check for updates", or "is there a newer version" → update flow.
- The user asks for a "recap" of recent work → recap flow across archive + active URs.

**Do NOT use when:**
- The user wants to see all changelog entries (more than 5) — point them at `CHANGELOG.md` directly instead of loading the full file.
- The user wants to *install* the skill fresh — that's the README install command, not this action.
- The install is global (under `~/.claude/skills/` etc.) and the user wants an update — refuse the auto-update per the preflight in Step 2 below, and redirect them.

## Input

This is a state-based action — the user's phrasing selects one of three response modes (each has its own section below):

- **Version request** — "what version", "version", "what's new", "release notes", "what's changed", "updates", "history" → report the current version + last 5 changelog entries.
- **Update check** — "update", "check for updates", "is there a newer version" → compare local against upstream and offer to apply.
- **Recap** — "recap" (dispatched with `mode: recap`) → summarize recent work across the archive and active URs.

## Responding to Version Requests

When user asks "what version", "version", "what's new", "release notes", "what's changed", "updates", or "history":

1. Report the version shown above
2. **Show last 5 skill releases**:
   - Read the first ~80 lines of `CHANGELOG.md` in the skill's root directory (same level as `SKILL.md`) — do NOT load the full file
   - Extract the 5 most recent version entries (split at `## ` headings, take first 5 blocks)
   - Reverse so newest is at the bottom (right where the user's eyes are)
   - Print them after the version number

## Responding to Update Checks

When user asks "check for updates", "update", or "is there a newer version":

1. **Fetch upstream**: Use your environment's web fetch capability to get the raw version.md from the upstream URL above
2. **Extract remote version**: Look for `**Current version**:` in the fetched content
3. **Compare versions**: Use semantic versioning comparison
4. **Report result** using the format below

### Report Format

**If update available** (remote > local):

1. **Tell the user**: `Update available: v{remote} (you have v{local}).`
2. **Preflight: confirm this is a project-local install, not a global one.** The update must overwrite the copy inside the current project — never a user-wide / global install.
   - Resolve the absolute path of the skill's root directory (where `SKILL.md` lives). Call this `<skill-root>`.
   - **Refuse to auto-update if `<skill-root>` is under any user-wide skills location**, including but not limited to:
     - `~/.claude/skills/...`
     - `~/.gemini/skills/...`
     - `~/.cursor/skills/...`
     - `~/.config/*/skills/...`
     - anything else under `$HOME` that isn't also inside the current project's root (`<project-root>`, resolved in the next bullet).
   - Resolve the current project's root with the repo's standard fallback: `git -C <invocation-dir> rev-parse --show-toplevel 2>/dev/null || pwd`. The `|| pwd` matches what `actions/install.md` uses — `git` is **optional** for the consuming project, so a non-git project resolves `<project-root>` to the invocation directory instead of being blocked (which would contradict the non-git install handling in Steps 3–4 below). The global-location refusal in the bullet above still applies in the non-git case — a skill under `~/.claude/skills/…` is rejected regardless of git. If `<skill-root>` is **not** a descendant of `<project-root>`, stop and report:
     ```
     Skill is installed at <skill-root>, which is outside the current project (<project-root>).
     Refusing to update a global/shared install from here. Either:
       - cd into the project that owns <skill-root> and re-run, or
       - install the skill locally inside this project (e.g. <project-root>/.claude/skills/do-work/) and re-run.
     ```
     Do NOT proceed. Do NOT suggest the curl command.
   - Only continue once `<skill-root>` is confirmed to live inside `<project-root>` (the git root, or the invocation directory when the project isn't a git repo) and not under any global skills location.
3. **Check for local changes** to shipped skill files at `<skill-root>`:
   - **Scope the check to skill-owned files only.** Ignore `do-work/` (queue data, archives, deliverables) — those are generated at runtime and should never block an update.
   - If `<skill-root>` is a git repo, run `git -C <skill-root> status --porcelain -- SKILL.md actions/ crew-members/ prompts/ interviews/ specs/ docs/ hooks/ tools/ CLAUDE.md AGENTS.md CHANGELOG.md README.md next-steps.md` (listing every shipped editable path) and check for uncommitted changes. Any dirty file in these paths will be clobbered by the tar extraction in step 5 if you proceed. (Previous archive files `CHANGELOG-2026-spring.md` and `CHANGELOG-pre-0.50.md` were removed in 0.76.0 — tarball-installed copies that want pre-0.65 release notes can browse them at commit `bf15fe2` on GitHub; git-cloned copies can `git show bf15fe2:CHANGELOG-2026-spring.md` locally.)
   - **Also catch _committed_ customizations before extraction** (git-repo installs) and local edits in non-git installs with a fresh upstream tarball diff. `git status --porcelain` only reports _uncommitted_ edits; a customization committed locally (including an edit to `actions/version.md` itself) otherwise looks clean. Before any destructive write (no pre-clean, no delete, no extraction yet), download the upstream tarball once, extract it to a temporary fresh upstream tree, and diff the current install against that tree:
     ```bash
     # Deterministic paths, not mktemp: Steps 3, 5, and 6 run as SEPARATE shell
     # invocations (a user-confirmation gate sits between them), and shell
     # variables do not survive across invocations — each block re-derives the
     # same paths from scratch.
     UPDATE_TMP="${TMPDIR:-/tmp}/do-work-update"
     UPSTREAM_TARBALL="$UPDATE_TMP/upstream.tar.gz"
     FRESH_UPSTREAM="$UPDATE_TMP/fresh"
     rm -rf "$UPDATE_TMP"
     mkdir -p "$FRESH_UPSTREAM"
     curl -fsSL https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz -o "$UPSTREAM_TARBALL" \
       || { echo "Upstream tarball download failed; aborting before any destructive write."; rm -rf "$UPDATE_TMP"; exit 1; }
     tar xzf "$UPSTREAM_TARBALL" -C "$FRESH_UPSTREAM" --strip-components=1 \
       --exclude='_dev' --exclude='do-work' --exclude='ai-reports' \
       --exclude='.vscode' --exclude='decisions'
     SHIPPED_PATHS=(SKILL.md actions crew-members prompts interviews specs docs hooks tools CLAUDE.md AGENTS.md CHANGELOG.md README.md next-steps.md)
     for shipped_path in "${SHIPPED_PATHS[@]}"; do
       diff -ru --new-file "$FRESH_UPSTREAM/$shipped_path" "<skill-root>/$shipped_path" | grep -v 'tools/queue-kanban/queue-kanban' || true
     done
     ```
     This diff includes legitimate upstream release changes, so don't treat every hunk as a blocker. Scan it before overwriting: current-side additions, local rewrites, or files present only in `<skill-root>` are committed/non-git customizations that would be clobbered (a file present only on the current side could instead be one upstream *removed* this release rather than a local addition — when unsure, surface it rather than assume). Surface them to the user and require explicit confirmation before proceeding. If the diff is only the expected upstream update, leave `$UPDATE_TMP` on disk for Steps 5-6; do not re-download a different archive. (The `grep -v 'tools/queue-kanban/queue-kanban'` filter drops the one expected noise line — the compiled, gitignored binary exists only on the current side and would surface as a phantom "Binary files … differ" customization on every update. Do NOT reach for `diff -x queue-kanban` here: `-x` matches basenames of files *and directories*, so it silently excluded the entire `tools/queue-kanban/` source tree from this check — real customizations to `model.go`/`board.js` sailed through invisible.)
   - **If any shipped skill files are dirty / have local modifications**: Stop and warn the user. List the modified files and ask for explicit confirmation before proceeding. Do NOT auto-update.
   - **If no local customizations are present**: Proceed to step 4 (snapshot + pre-clean) then step 5 (extract).
4. **Snapshot for rollback, then pre-clean discoverable directories.** First make the overwrite recoverable: a git-repo install already is (Step 3 confirmed a clean tree, so `git -C <skill-root> restore <file>` undoes any clobber after the fact); for a **non-git** install, copy the tree first — `cp -R <skill-root> <skill-root>.preupdate-bak`. Then pre-clean. `prompts/` and `interviews/` are upstream-controlled — their contents are owned by this skill, not the consuming project. `do-work prompts list` and `do-work interview list` glob `prompts/*.md` and `interviews/*.md`, so any upstream-removed file that stays on disk will still appear as a live workflow. The dirty check in step 3 has already confirmed these are clean, so removing the tracked `.md` files here is safe and the tar extraction will restore them fresh:
   ```
   find <skill-root>/prompts -maxdepth 1 -name '*.md' ! -name 'README.md' -delete
   find <skill-root>/interviews -maxdepth 1 -name '*.md' -delete
   ```
   Do NOT delete files in `prompts/` or `interviews/` subdirectories — only the top-level `.md` files are globbed. Do NOT touch `do-work/` or any other runtime directory.
5. **Run the update in place at `<skill-root>`** (the project-local path confirmed in step 2). Reuse the exact tarball downloaded and diffed in Step 3 so the reviewed bytes are the bytes extracted. This block runs in a fresh shell (Step 4 and a user gate sit between it and Step 3), so it re-derives the tarball path itself and refuses to improvise if the file is gone. `cd` into `<skill-root>` so the extraction cannot land in a global directory by mistake:
   ```bash
   UPDATE_TMP="${TMPDIR:-/tmp}/do-work-update"
   UPSTREAM_TARBALL="$UPDATE_TMP/upstream.tar.gz"
   test -s "$UPSTREAM_TARBALL" || { echo "Reviewed tarball missing at $UPSTREAM_TARBALL — go back to Step 3 and re-run the download + diff. Do NOT re-download here."; exit 1; }
   cd <skill-root> && tar xzf "$UPSTREAM_TARBALL" --strip-components=1 --exclude='_dev' --exclude='do-work' --exclude='ai-reports' --exclude='.vscode' --exclude='decisions'
   ```
   Never substitute a fresh `curl | tar` if the tarball is missing — extracting bytes that were never diffed defeats the entire Step 3 customization review.
   **Note:** tar extraction adds and overwrites files but does not delete files removed upstream. The `--exclude` flags (`_dev`, `do-work`, `ai-reports`, `.vscode`, `decisions`) keep the upstream repo's own dev tooling, queue/archive, sample reports, editor settings, and design-decision ADRs from landing in this install — belt-and-suspenders with the repo's `.gitattributes export-ignore`, which already strips all of them (plus the dev dotfiles `.gitignore`/`.gitattributes`) from the tarball (the flags also cover older tarballs built before that file existed). For non-discoverable directories (`actions/`, `crew-members/`, `specs/`, `docs/`) leftovers are harmless — the skill only loads files it references by name. For `prompts/` and `interviews/`, the pre-clean step above is what prevents ghost entries; if you skipped it, run `do-work prompts list` and `do-work interview list` after updating and delete anything that looks obsolete. Never delete `do-work/` (runtime state).
6. **Verify, then audit the overwrite**: Read `<skill-root>/actions/version.md` again and confirm the local version now matches the remote version. Then compare the post-update install to the same fresh upstream tree from Step 3 by re-running the `SHIPPED_PATHS` diff loop — re-derive `UPDATE_TMP`/`FRESH_UPSTREAM` exactly as the Step 3 block does (this is another fresh shell; the variables are gone). It should now be empty except for user-approved customizations that were deliberately re-applied. For a git install, `git -C <skill-root> diff -- <the shipped paths from Step 3>` or `git -C <skill-root> status` is still useful for the commit, but it is no longer the customization detector. For a non-git install, keep the `<skill-root>.preupdate-bak` snapshot until this audit passes, then delete the snapshot and `$UPDATE_TMP`.
7. **Report result**: `Updated to v{remote} at <skill-root>.`

Do NOT just print the curl command and ask the user to run it. You are the agent — run it yourself.

**If up to date** (local >= remote):

```
You're up to date (v{local})
```

**If fetch fails**:

```
Couldn't check for updates.
```

Attempt the update anyway using the curl command above (still respecting the preflight location check in step 2 and the dirty-tree check in step 3 — refuse if the install is global). If that also fails, report the error and provide the manual command as a fallback:

```
To manually update, cd into the **project-local** skill root (where SKILL.md lives inside *this* project — NOT ~/.claude/skills/, ~/.gemini/skills/, or any other global skills directory) and run:

cd <project-root>/path/to/skill-do-work
curl -sL https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz | tar xz --strip-components=1 --exclude='_dev' --exclude='do-work' --exclude='ai-reports' --exclude='.vscode' --exclude='decisions'

Or visit: https://github.com/knews2019/skill-do-work
```

## Responding to Recap Requests

When user asks "recap":

1. **Archive source** (`do-work/archive/UR-*/`): Read as before — title from `input.md`, REQs from `REQ-*.md` files inside each UR folder.
2. **Active source** (`do-work/user-requests/UR-*/`): Read `input.md` for the title. For REQs, scan `do-work/queue/REQ-*.md` files whose `user_request:` frontmatter field matches the UR id (e.g., `user_request: UR-143`). Also check `do-work/working/` for claimed REQs belonging to the UR.
3. **Merge**: Combine both lists, deduplicate by UR id (archive version wins if both exist), sort by UR number descending, take top 5.
4. **Label each UR**:
   - No label if fully archived
   - `(pending)` if the UR has any pending REQs
   - `(completed, awaiting archive)` if all its REQs are completed/done but the UR isn't archived yet
5. **Format as a "Recent Work" section**:
   ```
   ## Recent Work

   UR-144 — Block-level improved translation for ZH pairs
     REQ-361 — Block-level improved translation
   UR-143 — Model selector thinking variants (completed, awaiting archive)
     REQ-360 — Model selector thinking variants
   UR-142 — Quality-Score-Driven Repair Loop (completed, awaiting archive)
     REQ-359 — Quality-Score-Driven Repair Loop
   UR-011 — Dark mode implementation
     REQ-043 — Theme store setup
     REQ-044 — Settings panel toggle
   ```
   One line per UR, one indented line per REQ. No descriptions, no scores, no file lists.
6. **If no archive exists AND no active URs found**: Print `No completed work yet.` and skip this section.

## Red Flags

- The update flow is about to `cd` into a path under `~/.claude/skills/`, `~/.gemini/skills/`, or anywhere outside the current project's root (`<project-root>` — the git root, or the invocation directory for non-git projects) — STOP. Global installs must never be auto-updated from here.
- Remote version fetched from the upstream URL is empty, malformed, or older than local — abort; don't "update" backwards.
- The dirty-tree check reported modifications but the update proceeded anyway — user's local customizations will be clobbered.
- Recap lists the same UR twice (once from archive, once from active) — the dedup step was skipped; archive version should win.
- Version reported doesn't match the `**Current version**:` line at the top of this file — caching or path confusion; re-read the file from disk.

## Verification Checklist

- [ ] Version output shows the local version, the last 5 changelog entries, newest at the bottom.
- [ ] Update flow refused to proceed for global installs (Step 2 preflight).
- [ ] Update flow refused to proceed when shipped files had uncommitted changes (Step 3 dirty check), unless user explicitly confirmed.
- [ ] Update flow pre-cleaned `prompts/*.md` and `interviews/*.md` (Step 4) before tar extraction.
- [ ] Post-update verification re-read `actions/version.md` and confirmed the local version matches remote.
- [ ] Recap merged archive + active sources, deduped by UR id, kept the archive version on conflicts.
- [ ] Recap output is one line per UR + one indented line per REQ — no scores, no file lists, no descriptions.
