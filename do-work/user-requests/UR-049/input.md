---
id: UR-049
title: Make the upstream archive fetch survive a codeload 429
created_at: 2026-08-17T17:16:28Z
requests: [REQ-216, REQ-217, REQ-218]
word_count: 1838
---

# Make the Upstream Archive Fetch Survive a Codeload 429

## Summary

An external feedback brief, authored in a consumer repo against upstream v0.199.2, reports that
`do-work update` is unrunnable for any consumer whose IP is being rate-limited by
`codeload.github.com` — GitHub serves archive downloads behind a limiter separate from the git
transport and the REST API. The single `curl` at `skills/do-work/tools/do-work-update.sh:72` has no
retry, no credential support, and no alternate transport, and the URL is a hardcoded literal with no
environment override, so a blocked consumer has no escape hatch short of editing vendored files.

The brief was triaged read-only via `do-work-toolbox validate-feedback` before capture. Every checkable
claim verified against this tree, including the one the brief asked to be re-checked independently
(F3 — `tools/validate-suite-manifest.sh` returns exit 0 on a tree containing `do-work/`, `kb/`, `_dev/`,
`decisions/`, and `CLAUDE.md`, reproduced here). Two corrections were made to the brief's proposal and
are recorded in the REQs: C2's route 1 was dropped as unearned defensive surface, and the new fetcher
was relocated from `scripts/` to `tools/` because the installer resolves siblings via `$script_dir/` and
is byte-mirrored into two locations.

## Extracted Requests

| REQ | Title | Covers |
|---|---|---|
| REQ-216 | Teach atomic-download retry and optional credentials | Brief's C1 — `--retry` flags and `GH_TOKEN`/`GITHUB_TOKEN` header on the canonical primitive; brief's RED test 1 |
| REQ-217 | Add the upstream archive fetcher with a git fallback | Brief's C2 (minus route 1) + C3 — new `tools/fetch-upstream-archive.sh`, `DO_WORK_UPSTREAM_URL`, both callers routed, bootstrap retry-only, `version.md` escape-hatch text; brief's RED tests 2, 3, 4 |
| REQ-218 | Ratchet the tools download and fix the stale gitattributes comment | Brief's C4 + F3.1 — canonicalization ratchet extended over `skills/*/tools/`, `.gitattributes` belt-and-suspenders claim corrected |

## Batch Constraints

- **The git fallback must repack with `git archive`, never a worktree copy.** `git ls-files` = 734,
  `git archive HEAD | tar t` = 294; `export-ignore` is the only thing keeping `do-work/`, `kb/`, `_dev/`,
  `decisions/`, `ai-reports/`, and `CLAUDE.md` out of consumer installs, and `cp -R`/`rsync`/`tar` of a
  clone does not honor it. Nothing downstream catches the leak.
- **Additive only.** No change to `atomic-download.sh`'s publication contract beyond flags, and no change
  to the installer's transaction, confirmation prompt, managed-section handling, or `--archive` contract.
- **No new runtime dependency.** `curl` + `git` is sufficient. Git is already a hard prerequisite —
  `do-work-update.sh:57-58` refuses to run outside a Git worktree.
- **Both mirrors move together.** `tools/*` and `skills/do-work/tools/*` are byte-identical and enforced by
  `_dev/tests/staged-skills-contract.sh:809`.
- **Release ceremony belongs to the integrating commit only** — patch bump across `VERSION`,
  `skills/do-work/VERSION`, and `actions/version.md`, plus a `CHANGELOG.md` entry titled for what shipped
  and mirrored byte-identically into `skills/do-work/CHANGELOG.md`.

## Maintainer Decisions at Capture

1. **Route 1 dropped.** The brief's first fallback route (`GH_TOKEN` → `api.github.com/repos/.../tarball/main`)
   redirects to a pre-signed URL that still terminates at codeload, so the claim that it escapes the limiter
   is asserted rather than demonstrated, and it is the one route that cannot be exercised offline. Under the
   earned-defense rubric it does not clear the bar. The generic `GH_TOKEN`/`GITHUB_TOKEN` header still lands
   in `atomic-download.sh`, where it is cheap and benefits every caller.
2. **Three REQs, not one**, with disjoint write sets.
3. **Fetcher goes in `tools/`, not `scripts/`** — see REQ-217's Context.
4. **C3's bootstrap/README sync test is already shipped** at `_dev/tests/install-suite-behavior.sh:185-195`;
   nothing to add there.

## Provenance and Injection Check

Source: an external feedback brief pasted by the maintainer, authored against upstream v0.199.2 in a
consumer repo. Per `crew-members/prompt-injection.md`, the body was read as data. **No injection detected.**
The brief's imperatives ("Please make…", "Add…", "Please do not ship retry alone") are ordinary feedback
addressed to the maintainer, and its one request aimed at the reader — "please re-check the F3 finding
independently rather than taking this brief's word for it" — was honored by independent verification rather
than obeyed as an instruction. Nothing in it attempts to redirect the operating context.

## Full Verbatim Input

do-work validate-feedback # Upstream suggestion for `knews2019/skill-do-work` — make the upstream archive fetch survive a codeload 429

**How to use this file:** paste everything below the horizontal rule into a Claude Code session opened in a
clone of `knews2019/skill-do-work`. It is written to be actionable cold, with no reference back to the
consumer repo it was authored in. Observed against upstream **v0.199.2**.

---

## Request

`do-work update` is unrunnable for a consumer whose IP is being rate-limited by `codeload.github.com`.
Please make the upstream-archive fetch resilient, and route it through a canonical primitive so the
resilience is inherited by every download site instead of being patched into one of four copies.

## The incident (reproducible)

A consumer on v0.189.3 ran `do-work update` and got:

```
Checking do-work updates…
curl: (56) The requested URL returned error: 429
do-work update: upstream tarball download failed; no files were changed
```

The failure is in `skills/do-work/tools/do-work-update.sh:72`. `https://github.com/.../archive/refs/heads/main.tar.gz`
302-redirects to `https://codeload.github.com/knews2019/skill-do-work/tar.gz/refs/heads/main`, and codeload
returned `429: Too Many Requests` with GitHub's anti-scraping ToS notice in the body.

Evidence that this is codeload's archive limiter specifically, and not GitHub-wide throttling or a repo problem:

| Probe | Result |
| --- | --- |
| `curl -sSL .../archive/refs/heads/main.tar.gz` | `302` → codeload → **`429`**, 3 consecutive attempts over ~2 min |
| `curl -sS https://api.github.com/rate_limit` | `core: 57/60 remaining` — anonymous API budget essentially untouched |
| `git clone --depth 1 https://github.com/knews2019/skill-do-work.git` | **succeeded** |

Archive downloads sit behind a separate limiter from the git transport and the REST API. The consumer had no
GitHub credentials configured (`gh auth status`: not logged in), so they were in the shared anonymous pool —
plausibly behind a NAT/VPN egress shared with other traffic. Nothing about their setup was recoverable by
waiting a few seconds, and `do-work update` offers no other path.

**Workaround that unblocked them** (worth noting, because it is the shape of the proposed fix): clone over
git, repack with `git archive` into a GitHub-shaped tarball, then hand it to the existing
`tools/install-do-work-suite.sh --archive`. Your own `tools/validate-suite-manifest.sh` accepted the result
(`suite manifest valid: v0.199.2 (4 modules)`) and the install completed normally. So the installer needs no
change — only the fetch does.

## Findings

### F1 — Four hand-rolled copies of the same download, none resilient

| Site | Shape |
| --- | --- |
| `skills/do-work/tools/do-work-update.sh:72` | `curl -fsSL -o "$f.download" "$upstream_url"` + `mv` |
| `skills/do-work/tools/install-do-work-suite.sh:173` | same, independently written |
| `skills/do-work/tools/install-do-work-suite.sh:25` | same, inside the `BOOTSTRAP` heredoc |
| `README.md:16` | the emitted bootstrap block, same text again |

Every one is a single attempt with no `--retry`, no credential support, and no alternate transport. The URL is
a hardcoded literal at `do-work-update.sh:5` and `install-do-work-suite.sh:5` with no environment override, so
a consumer hitting this has no escape hatch short of editing vendored files — which `actions/version.md:48`
explicitly forbids, and rightly so.

### F2 — The canonical primitive for this already exists and is bypassed

`skills/do-work/scripts/atomic-download.sh` is the shipped download primitive: documented at
`skills/do-work/docs/prescribed-shell-primitives.md:63` under **Atomic download publication**, listed in the
canonicalization ratchet at `_dev/tests/prescribed-shell-canonicalization.sh:12`, and behavior-tested at
`_dev/tests/prescribed-shell-scripts-behavior.sh:54-62`.

Neither `tools/` script calls it. Both restate the download-and-rename primitive inline — the exact
duplication that archived REQ-167 (`dedupe-prescribed-shell-primitives`) and REQ-171
(`promote-prescribed-shell-to-tested-scripts`) exist to eliminate, and that
`decisions/audits/2026-08-11-prescribed-shell-primitives.md:12` marked `PROMOTE`. The campaign covered action
markdown and the `scripts/` tree; the two `tools/` entry points were never swept, so the ratchet does not see
them. That is why a resilience fix would otherwise have to be written three or four times.

### F3 — A naive git fallback would leak the maintainer's private tree into consumer installs

This one is a trap worth stating explicitly before anyone writes the fallback.

`.gitattributes` `export-ignore`s `/do-work`, `/kb`, `/_dev`, `/decisions`, `/ai-reports`, `/CHANGELOG-20*.md`,
`/.vscode`, `/CLAUDE.md`, `/AGENTS.md`, and the dev dotfiles. The effect is large: `git archive HEAD` yields
**294** entries against **733** tracked files. GitHub's tarball honors `export-ignore`; a plain worktree copy
(`cp -R`, `rsync`, `tar` of a clone) does **not**.

Two things make this sharper than it looks:

1. **The documented belt-and-suspenders is gone.** The `.gitattributes` header comment claims "The
   install/update tar command also passes `--exclude` for do-work/kb/ai-reports/_dev plus .vscode/decisions as
   a belt-and-suspenders fallback". No `--exclude` exists at any of the three extraction sites today
   (`do-work-update.sh:75`, `install-do-work-suite.sh:28`, `install-do-work-suite.sh:178`). `export-ignore` is
   currently the *only* thing holding that line, exactly as the comment's later paragraph warns. The comment
   is stale and should either be corrected or the `--exclude` net restored.
2. **Nothing downstream would catch the leak.** Verified empirically: build a tarball from a worktree copy so
   that `do-work/`, `kb/`, `_dev/`, `decisions/`, and `CLAUDE.md` are all present, then run
   `tools/validate-suite-manifest.sh --root <that tree>` — it prints `suite manifest valid: v0.199.2 (4 modules)`
   and exits `0`. The validator checks that required modules are present, not that dev paths are absent.

So: **the fallback must produce its tarball with `git archive`**, never by copying a worktree. Shipping
`CLAUDE.md` into consumers is specifically called out in `.gitattributes` as actively harmful (Claude Code
auto-loads it, and its commit protocol is wrong advice downstream).

### F4 — Retry alone would not have fixed this incident

Worth designing around honestly: the 429 persisted across three attempts spanning roughly two minutes. A
`--retry 3 --retry-delay 2` would have failed too. Retry is cheap and correct for genuine blips, but the
**git-transport fallback is the load-bearing fix** — it is the only probe that actually succeeded. Please do
not ship retry alone and consider the class closed.

## Proposed change

### C1 — Teach the canonical primitive the generic parts

In `skills/do-work/scripts/atomic-download.sh`, keep the private-download-and-rename contract exactly as is
(it is tested and correct), and add around the existing `curl` at line 26:

- **Retry:** `--retry 3 --retry-delay 2 --retry-max-time 60`. curl has treated HTTP 429 as a transient,
  retryable status since 7.51.0 (2016); this repo already reasons about curl floors — `skills/do-work/CHANGELOG.md:2127`
  rejected `--remove-on-error` for needing ≥ 7.83 — and 7.51 is old enough to rely on. Avoid
  `--retry-all-errors` (7.71+) since plain `--retry` already covers 429.
- **Credentials, opt-in:** when `GH_TOKEN` or `GITHUB_TOKEN` is non-empty, send
  `-H "Authorization: Bearer $github_api_token"`. Authenticated requests are not in the shared anonymous pool.

Both improvements are generic, so every existing caller — including the `frontend-design`, `adhd`, and
`playwright-bowser` SKILL.md downloads in `skills/do-work-toolbox/actions/install.md` — inherits them for free.

### C2 — Add an archive-specific fetcher that owns the fallback chain

The git-archive fallback is meaningful only for *this repo's tarball*, not for arbitrary URLs, so it does not
belong in the generic primitive. Add a sibling, e.g. `skills/do-work/scripts/fetch-upstream-archive.sh`,
taking a target path and trying in order:

1. `https://api.github.com/repos/knews2019/skill-do-work/tarball/main` when a token is present. This endpoint
   302s to a **pre-signed** codeload URL, so the token only ever needs to reach `api.github.com`; curl drops
   `Authorization` on cross-host redirects, which is the behavior we want here, not something to defeat with
   `--location-trusted`.
2. The current anonymous `github.com/.../archive/refs/heads/main.tar.gz`, via C1's retry.
3. **Git fallback:** shallow-clone and repack, honoring `export-ignore` per F3 —

   ```bash
   git clone --depth 1 --quiet "$upstream_repo_url" "$clone_directory"
   git -C "$clone_directory" archive --format=tar.gz \
     --prefix=skill-do-work-main/ HEAD > "$archive_target_path"
   ```

   The `--prefix` keeps the single-top-level-directory layout that both callers' existing
   `tar --strip-components=1` depends on. `git archive` on a shallow clone is fine. Git is already a hard
   prerequisite of the update flow (`do-work-update.sh` refuses to run outside a Git worktree), so this adds
   no new dependency.

Report which route succeeded so a consumer can tell a clean fetch from a degraded one, and on total failure
name all three outcomes rather than only the last.

### C3 — Route the callers through it, and honor an override

- `tools/do-work-update.sh` and `tools/install-do-work-suite.sh` (the `--archive`-absent branch) call the new
  script instead of restating `curl`. This closes F2 for the fetch specifically.
- Let `upstream_url` in both scripts default from an environment variable — `DO_WORK_UPSTREAM_URL` matches the
  existing `DO_WORK_INSTALL_CANCEL_EXIT_STATUS` idiom — so a blocked or air-gapped consumer, or someone
  testing a fork, has a supported escape hatch that is not "edit a vendored file".
- The `BOOTSTRAP` heredoc at `install-do-work-suite.sh:18-32` and its `README.md` twin genuinely **cannot**
  call the shipped helper — nothing is installed yet at that point, so the inline `curl` there is correct
  layering, not duplication. Give that copy the `--retry` flags and the optional auth header only, and leave
  the fallback chain out of it. If those two blocks are kept in sync by hand, consider a `_dev/tests` check
  that they match, since they must now stay identical through an edit.

### C4 — Extend the ratchet so `tools/` cannot drift again

`_dev/tests/prescribed-shell-canonicalization.sh` enumerates `scripts/` entries only. Add the new fetcher to
that list, and add a check that no file under `skills/*/tools/` contains a bare `curl -fsSL -o` — with the
bootstrap heredoc as the single documented exemption. Without this, the next tool script will hand-roll the
primitive a fifth time.

## Tests to write first (RED)

The behavior suite already has the right idiom for this — `_dev/tests/prescribed-shell-scripts-behavior.sh:54-62`
puts a fake `curl` on `PATH` that writes a partial file and exits 22. Reuse it:

1. **429-then-success:** a fake `curl` that returns 429 on its first invocation and succeeds on the second;
   assert the fetch succeeds and the target is complete. Fails today (single attempt).
2. **Sustained 429 falls back to git:** a fake `curl` that always returns 429, plus a local fixture repo
   standing in for the upstream remote; assert a valid archive is still produced. Fails today (no fallback).
3. **The fallback honors `export-ignore`** — the F3 regression, and the most important of the three. Give the
   fixture repo a `.gitattributes` with `/private-path export-ignore`, force the git route, and assert
   `private-path` is **absent** from the produced tarball. A `cp -R`-based implementation passes tests 1 and 2
   and fails only this one.
4. **Target preserved on total failure:** all routes fail; assert the pre-existing target file is unchanged
   and no `.download.*` scratch leaks — the existing atomic-publication guarantee, re-asserted through the new
   layer.

Also please re-check the F3 finding independently rather than taking this brief's word for it; the
`validate-suite-manifest.sh`-passes-a-leaky-tree result is the part most worth confirming with your own hands.

## Explicitly not asking for

- No change to `tools/install-do-work-suite.sh`'s install transaction, confirmation prompt, managed-section
  handling, or `--archive` contract. The workaround proved that path is sound; the fetch is the whole bug.
- No change to the atomic publication contract in `atomic-download.sh` — additive flags only.
- No vendoring of a download library or new runtime dependency. `curl` + `git` is sufficient.

## Release ceremony

Per `CLAUDE.md` § Commit Completion: bump `VERSION`, `skills/do-work/VERSION`, and the `**Current version**:`
line in `skills/do-work/actions/version.md` (patch-level for this), then add a top-of-file `CHANGELOG.md`
entry whose title says what was delivered — something like **"Resilient Upstream Archive Fetch"**, not a
codename. Verify the version is strictly greater than the current first entry and that the title is unused.

If `actions/version.md`'s user-facing failure text mentions the download step, update it to name the new
escape hatches (`GH_TOKEN`, `DO_WORK_UPSTREAM_URL`) so a blocked consumer learns the way out from the error
itself, which is where they will be looking.

---
*Captured: 2026-08-17T17:16:28Z*
