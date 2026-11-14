---
id: REQ-732
title: '[impact-negligible] Allow docs/run-prompts/ in the staged-skills contract and move the run prompts there'
status: claimed
created_at: 2026-11-14T22:47:58Z
user_request: UR-176
domain: general
prime_files: []
tdd: false
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
related: [REQ-726, REQ-727, REQ-728, REQ-729, REQ-730, REQ-731]
batch: tinka-follow-ups
write_set: ["_dev/tests/staged-skills-contract.sh", ".gitattributes", "docs/run-prompts/*.md", "_dev/run-prompts/*.md"]
claimed_at: 2026-11-14T23:28:20Z
---
# Allow docs/run-prompts/ in the Staged-Skills Contract and Move the Run Prompts There
## What
Let the staged-skills contract accept `docs/run-prompts/` under the root `docs` folder, keep that folder out of the shipped package, and move the prompt files from `_dev/run-prompts/` there.
## Why
The maintainer's global rule says a prompt for another session goes to `docs/run-prompts/<date>-<slug>.md`. In this repository that file broke a check, so the rule could not be followed here.
## Context
- `_dev/tests/staged-skills-contract.sh` (the loops near lines 75 and 99 to 114, message "legacy root runtime must be retired at modular cutover") allows only `docs/kb/`, `docs/decisions/` and `docs/changelog/` under the root `docs` folder.
- On 2026-11-14 the file `docs/run-prompts/2026-11-14-tinka-canvas-suggestion.md` turned the heavy lane `staged-skills` red five commits after it was added, and was removed in commit 00568685. The fast gate did not catch it.
- `.gitattributes` marks the three allowed folders `export-ignore`, and the contract asserts those lines.
- REQ-724 (move the maintainer's global standing rules into the skill) is queued and names this question as open. It does not depend on this request.
## Detailed Requirements
- Add `docs/run-prompts` to both lists in the contract: the export-ignore assertion and the allowed paths of the leftover scan.
- Add the matching `export-ignore` line to `.gitattributes`, with one comment line in the style of its neighbours.
- Move every file from `_dev/run-prompts/` to `docs/run-prompts/` with `git mv`, and remove the empty folder.
- Sweep the repository for links or paths that name `_dev/run-prompts/` in files that are not archived records, and repoint them. Archived request records keep their text.
## Constraints
- Nothing under `skills/` or `suite/` changes, so this is no release.
- The contract keeps refusing every other new file under the root `docs` folder.
## Dependencies
None. The other requests of this batch can run in any order.
## Red-Green Proof
**RED prompt/case:** With a file `docs/run-prompts/probe.md` present, run `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`.
**Why RED now:** The lane fails with the legacy-root message that names the probe file.
**GREEN when:** The same lane passes with that file present and still fails with a file `docs/other/probe.md`; both probe files are removed afterwards; the moved prompts are in `docs/run-prompts/`.
**Validation:** Inferred during capture. The proof comes from the prompt file the maintainer asked to capture; it was not confirmed as a separate question.
## Open Questions
- [x] Where do run prompts live in this repository: `docs/run-prompts/` or `_dev/run-prompts/`? → Allow docs/run-prompts/ (Recommended).
  Note (2026-11-14): The maintainer chose to make the global rule work here as in every other project. The accepted cost: one more allowed folder under `docs` to keep out of the shipped package. Keeping `_dev/run-prompts/` with a line in `CLAUDE.md` is out of scope.
## Builder Guidance
- Certainty: firm. The proof uses the heavy lane `staged-skills`, which the normal heavy drain of a run executes.
- Scope: the smallest change. Two list entries in the contract, one line in `.gitattributes`, and the moves.
## Full Context
See `do-work/user-requests/UR-176/input.md` for complete verbatim input.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-run-prompts-folder-rule.md*
