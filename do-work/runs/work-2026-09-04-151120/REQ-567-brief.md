# Builder Brief — REQ-567

Implement REQ-567 (repair the shipped lesson links to archived UR paths) on branch `worktree-agent-REQ-567-repair-shipped-lesson-links` in worktree `/Users/t2/Desktop/e1-experimental-repos/worktree-agent-REQ-567-repair-shipped-lesson-links`.

Read the committed request at `do-work/working/REQ-567-repair-shipped-lesson-links-to-archived-ur-paths.md`, plus `skills/do-work/crew-members/general.md`, `skills/do-work/crew-members/coding-guardrails.md`, `skills/do-work/crew-members/communication-style.md`, `skills/do-work/crew-members/backend.md`, and `skills/do-work/crew-members/testing.md`.

Write boundary: modify only `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`. Do not write any `do-work/` path in the builder worktree. Preserve the three family markers, labels, and lesson prose; replace only the REQ-491, REQ-492, and REQ-493 destinations with their actual `do-work/archive/UR-095/` paths and verified `#lessons-learned` anchors.

TDD evidence: the unchanged canonical gate already failed with fingerprint `sha256:3af85b84722557f94ddfd466fc32136086fb5fed306e478bd344f689902472ff` and exactly those three broken links. After the edit run the focused shipped-package reference contract, then `bash _dev/tests/maintainer-verify.sh`. Review the exact diff, run `git diff --check`, and commit the single-file change on the builder branch.

Write the complete handback to the main-tree path `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-04-151120/REQ-567-handback.md` using `apply_patch`. Include branch/commit, file manifest, RED/GREEN evidence, commands and results, P-A-U APPLY/UNIFY evidence, required lessons read/missing, decisions, discovered tasks, and any integration seam. Return only a one-line status after the file exists.
