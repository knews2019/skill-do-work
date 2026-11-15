# Hand-back — REQ-736 (the ai-report judge line in toolbox help shows --page)

## Commits
- `ab48d0a0` [REQ-736] Show the --page option on the ai-report judge line in toolbox help

## Files
- `skills/do-work-toolbox/actions/help.md` (line 18 only)

## Scope changes
none

## Red-green evidence
- Base `547d4cc5`: probe exit 1, `FAIL: help.md judge line does not show [--page <file.html>]`.
- Tip `ab48d0a0`: probe exit 0, `PASS: REQ-736 help.md judge line shows --page`.

## Verification
1. GREEN probe from the worktree: PASS (exit 0).
2. `git diff --stat HEAD~1`: `skills/do-work-toolbox/actions/help.md | 2 +-`, 1 file, 1 insertion, 1 deletion.
3. `bash _dev/tests/shipped-package-reference-contract.sh`: exit 0, `shipped package reference contract: PASS`.
4. `bash _dev/tests/staged-skills-contract.sh`: exit 2 with `staged-skill contract probes are heavy-only; run _dev/tests/maintainer-verify.sh --heavy after user permission.` The script refuses a standalone run; it ran no check. The brief forbids the builder from running `maintainer-verify.sh`, so the integrator's heavy lane must cover it. The change is one prose line inside a `text` fence, so no staged-skill check is expected to react.

## P-A-U
- PLAN: replace help.md line 18 usage with `ai-report judge <dir> [--page <file.html>]`, keep the description after two spaces (D-01).
- APPLY: one exact-string replacement, one commit by exact path.
- UNIFY: `git diff --stat HEAD~1` shows only help.md, 1+/1-; checked the diff holds no other line and no debug text; probe and reference contract pass.

## Decisions
- D-02 (DECIDE & STATE): kept the description unchanged ("Render-check a report bundle (widths, themes, links)"). Adding "or one named page" would make the overflowed line longer, and the usage already shows the option.

## Discovered Tasks
- `_dev/tests/staged-skills-contract.sh` named as a builder verify step in the brief, but it refuses a standalone run (exit 2, heavy-only). Briefs should drop it or name the heavy lane for the integrator.

## Integration seams
none. Only help.md line 18; REQ-735 touches only `crew-members/` files.

## Proposed CHANGELOG entry
**Toolbox Help Shows the ai-report judge --page Option**
What changes for an installed copy: `do-work-toolbox help` lists `ai-report judge <dir> [--page <file.html>]`, so the single-page check added in 0.319.0 is visible from the menu. Help text only; no behaviour changes.

## Proposed lesson
none

## Lessons read
- `_dev/primes/lessons-action-files.md` `[family: restated-mechanism-unchecked]` bullets (REQ-639, REQ-688, REQ-689). Applied: read `ai-report.md:36` and `ai-report-guide.md:80` before restating the form.
