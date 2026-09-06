# Prime: Releases

Read this before committing a change to shipped files. Maintainer-only files (`CLAUDE.md`, `_dev/`, `do-work/`, `decisions/`) commit without any of it.

A commit that changes shipped files under `skills/`, `tools/`, or `suite/` is a release, and only such a commit is: since 0.305.31 the finalizer refuses a release whose implementation changes nothing under a declared module root, `suite/` or `tools/` beyond release metadata (`RELEASE-WITHOUT-SHIPPED-CHANGE`), after four versions were bumped for `_dev/tests`-only changes. Bump size, version mirrors, the entry format, and the finalize transaction are decided per `skills/do-work/actions/work-reference.md` → Changelog Entry Procedure (Step 9); what follows are the house rules that procedure applies.

- **The changelog title says what was delivered.** A reader scanning only headings should know what changed ("Board View Filters", not "The Fine Sieve"). No whimsical codenames. Verify the title is not already used by an earlier entry. Keep the entry brief, newest on top, lead with value not implementation. Every version gets an entry.
- **Repository-only dated history uses canonical links** of the form `https://github.com/knews2019/skill-do-work/blob/main/...`, because the installed core package does not carry those sidecars.
- **The installed changelog mirror is byte-identical.** Copy root `CHANGELOG.md` to `skills/do-work/CHANGELOG.md` after the entry and any history-link edits; `_dev/tests/shipped-package-reference-contract.sh` enforces it.

## Traps

- [family: canonical-link-outlives-its-target] Closing a UR moves archived REQ paths → sweep canonical history links in the same change; a valid URL prefix does not prove its target exists.

## Stakes

- Release ownership and mirrored metadata
  Req: release only shipped changes; keep every owned version and changelog mirror coherent.
  Value: consumers receive an identifiable, reproducible suite update.
  Risk: stale mirrors misidentify installed content; incorrect ownership can rewrite another package. See `skills/do-work/tools/do-work-cli/internal/publication/release.go` and `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go`.

## Lessons

See [`lessons-releases.md`](lessons-releases.md) — read it before changing release ownership, history links or mirrors.
