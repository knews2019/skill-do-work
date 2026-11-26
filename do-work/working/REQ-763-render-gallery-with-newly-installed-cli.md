---
id: REQ-763
title: 'Render the static gallery with the newly installed CLI after an upgrade'
status: claimed
created_at: 2026-11-26T11:26:10Z
user_request: UR-193
addendum_to: REQ-743
domain: backend
suggested_spec: bug-fix
prime_files: [skills/do-work/tools/prime-do-work-update.md]
required_lessons: [skills/do-work/tools/lessons-do-work-update.md]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: true
maintenance: false
claimed_at: 2026-11-26T11:32:00Z
---

# Render the static gallery with the newly installed CLI after an upgrade

## What

When `just run-do-work-update` installs a different suite version, generate `docs/report-gallery/index.html` with the newly installed renderer during that invocation.

## Prior Implementation

REQ-743 (generate the report gallery after a suite update) added the refresh in `a8031e04`. In `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go`, the update calls the gallery handler compiled into the already-running CLI. It reads the new sample data but can still write the old layout. The current renderer places repository reports first when the project contains reports.

## Review Evidence

Finding F1, P2, pasted review. Surface-cost assessment was N/A for a direct fix. The reviewer reported the following in the consumer project:

> Opening docs/report-gallery/index.html directly still shows the ten sample designs before the repository reports and hides 30 reports under “Browse other reports.” The new report card was added using the old layout, so the checked-in page and the live gallery disagree. Regenerate the page with the updated CLI and confirm that “Repository reports” is the first section, as specified in the gallery guide.

## Red-Green Proof

**RED prompt/case:** Start a fixture project with a repository report, sample designs, and an installed version whose renderer puts samples first. Use that version’s update entry point to upgrade once to the reports-first version. Read the generated static HTML immediately.
**Why RED now:** The running updater can write the samples-first layout even after installing the reports-first renderer. The existing suite update test checks the installed sample content and misses this difference.
**GREEN when:** An automated version-changing upgrade regression test observes “Repository reports” as the first content section, with the fixture report in it and “Sample designs” later, after that single invocation.
**Validation:** User confirmed.

## AI Execution State (P-A-U Loop)

- [ ] **[PLAN]:** Read the listed prime and record the implementation approach.
- [ ] **[APPLY]:** Implement the planned change and its regression test.
- [ ] **[UNIFY]:** Run the regression and review the changed files.

*Source: Fix F1 so a version-changing suite update generates the static gallery with the newly installed renderer. Verify an upgrade from a samples-first version puts Repository reports first*
