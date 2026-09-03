# Prime: Releases

Read this before committing a change to shipped files. Maintainer-only files (`CLAUDE.md`, `_dev/`, `do-work/`, `decisions/`) commit without any of it.

A commit that changes shipped files under `skills/`, `tools/`, or `suite/` is a release. Bump size, version mirrors, the entry format, and the finalize transaction are decided per `skills/do-work/actions/work-reference.md` → Changelog Entry Procedure (Step 9); what follows are the house rules that procedure applies.

- **The changelog title says what was delivered.** A reader scanning only headings should know what changed ("Board View Filters", not "The Fine Sieve"). No whimsical codenames. Verify the title is not already used by an earlier entry. Keep the entry brief, newest on top, lead with value not implementation. Every version gets an entry.
- **Repository-only dated history uses canonical links** of the form `https://github.com/knews2019/skill-do-work/blob/main/...`, because the installed core package does not carry those sidecars.
- **The installed changelog mirror is byte-identical.** Copy root `CHANGELOG.md` to `skills/do-work/CHANGELOG.md` after the entry and any history-link edits; `_dev/tests/shipped-package-reference-contract.sh` enforces it.
