---
id: REQ-763
title: 'Render the static gallery with the newly installed CLI after an upgrade'
status: completed
route: B
heavy_verified_at: 2026-11-26T11:59:27Z
heavy_verified_revision: 086aa28d13f71da13b29a672f480ea9d8acc9a1d
commit: 086aa28d13f71da13b29a672f480ea9d8acc9a1d
review_at: 2026-11-26T11:55:34Z
integration_at: 2026-11-26T11:52:34Z
dispatch_at: 2026-11-26T11:42:20Z
write_set: [skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go, skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.gohtml, skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go, skills/do-work-toolbox/docs/ai-report-guide.md]
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-11-26T11:35:29Z
  basis: [Route B, 5-file write set, 2 subsystems involved, 2 acceptance criteria, cross-route regression gates]
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
completed_at: 2026-11-26T11:59:34Z
release_at: 2026-11-26T11:59:34Z
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

- [x] **[PLAN]:** Keep the update process trusted and render the static page using the installed suite’s declarative layout. Verify the new layout wins across a real version change.
- [x] **[APPLY]:** Integrated the layout template, updater handoff, migration instructions and upgrade regression.
- [x] **[UNIFY]:** Reviewed the implementation manifest and diff. Both affected packages and go vet passed in the builder checkout.

*Source: Fix F1 so a version-changing suite update generates the static gallery with the newly installed renderer. Verify an upgrade from a samples-first version puts Repository reports first*

## Open Questions

- [x] “The samples-first updater renders with its own compiled template after installing the new files. A change in the incoming release cannot change that first invocation. May I adjust the acceptance test to cover future upgrades?” → Yes. Fix future upgrades and document the one-time gallery refresh for existing installs.

On 2026-11-26 the user confirmed that the automated one-invocation proof covers future upgrades. Existing installations need instructions for one explicit gallery refresh.

## Triage

**Route: B** — The failure and expected output are known. Explore how to select the installed renderer while preserving the updater’s trusted-engine boundary.

## Plan

Planning not required — Route B, exploration-guided implementation.

## Exploration

The original samples-first binary embeds its template. The updater’s trusted-engine regression forbids launching downloaded shell or Go code. Future updaters can read the installed gallery layout as data and render it with the existing Go view model. Keep the ownership marker stable and the user’s migration instructions paste-ready.

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.gohtml`
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go`
- `skills/do-work-toolbox/docs/ai-report-guide.md`

Adjust the upgrade handler or fixture helpers if the data handoff needs them.

Acceptance requires distinguishable old and new layouts, a repository report shown first following one future upgrade, and one documented refresh command for consumers of the old updater.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (25,337 estimated tokens). The fixture must exercise the running updater across versions. Partial family coverage precludes selecting individual entries under the 2,000-token budget. The 847-token updater satellite directly covers this first-upgrade issue and was read in full.

## Historical Replay

A temporary consumer installed from commit `d4f842d3` (0.332.1) upgraded successfully to 0.332.6 through its own `just run-do-work-update`. The generated page’s first section was `samples`. A single `bash .claude/skills/do-work/tools/do-work-cli.sh --repo-root . ai-report-gallery` invocation changed the first section to `generated-designs` (Repository reports). The test used a local archive server and a fixture report.

## Pre-Flight

The focused suite-install and gallery tests passed on the claim revision. The first repository gate exited 1 on the 30-second test-file budget (unrelated second-opinion tests, 33.45s). Its direct retry passed in 153s with the slowest CLI test file at 21.53s. Both advance preflight gate records were satisfied before implementation.

## Implementation Summary

The old updater called its own compiled renderer after replacing the suite. The update now loads the installed layout as data from the verified core skill and renders it with the running CLI’s report data. The ownership header and static page publisher stay with that CLI. The migration section documents the confirmed legacy refresh.

- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.gohtml` (created)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_resolve_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go` (modified)

## Qualification

The merge passed diff hygiene and the implementation check. The guide, Go renderer, extracted layout and regression test match the initial scope. Six additional files wire the layout selection through the updater, supply the layout in fixture archives, pass the embedded layout explicitly to existing callers, or test preservation of a previous page when loading fails. All ten changed paths support the requested behavior. Reviewed the complete merge range and checked the static gallery page in a temporary consumer.

A future-updater fixture compiled against a samples-first layout ran a single `just run-do-work-update` and produced Repository reports first, containing the report link, with Sample designs later. No extra gallery command ran in this proof.

## Testing

**Red-green validation:** The builder ran `TestVersionChangingUpdateRendersInstalledGalleryLayout` against the old refresh behavior. It installed the target version but the first section still contained “Sample designs” and the target version’s sample. With the installed-layout handoff the same test passed in 2.23 seconds. The merged focused probe passed all four selected tests, including preservation of the previous gallery for unreadable and incompatible layouts. A separate replay through the public Just recipe passed with a samples-first caller and reports-first incoming suite.

**Repository gate:** `_dev/tests/maintainer-verify.sh` passed on the merge revision in 144 seconds. Its Go suites ran 427 board tests and 1,017 CLI tests.

## Review

**Acceptance: Pass.** Independent review of the merged range found no actionable issues. Requirements 100%, quality 98%, tests 100%, scope 98%, overall 99%. The additional six files are accounted for in Qualification. Three focused tests passed again during review. A template requiring fields unavailable in the running CLI can produce a warning and leave the previous page intact. Regenerating with the new CLI uses its updated data model.

## Lessons Learned

**What worked:** Distinct source and destination layouts exposed the failure even when the updated sample data was already visible in the old layout.
**What didn't:** Running the same compiled renderer with two versions of sample data could not distinguish the layouts.
**Worth knowing:** An incoming release cannot make an unmodified old updater load a template it never read. Test both the historical migration command and the first invocation from a template-aware updater.

## Orientation

A version-changing suite update now uses the installed template to lay out the static gallery. The running CLI remains responsible for collecting report data and publishing the generated page. The updater prime’s file references still resolve.

## Heavy Verification Plan

- Base revision `ff2d66495d925e0780e43a0d4d545daeddea817c`
- Target revision `086aa28d13f71da13b29a672f480ea9d8acc9a1d`
- `do-work-cli-integrations` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`. skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.gohtml matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_resolve_test.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go matched subtree skills/do-work/tools/do-work-cli.
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. skills/do-work-toolbox/docs/ai-report-guide.md matched subtree skills; skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go matched subtree skills; skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go matched subtree skills; skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go matched subtree skills; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go matched subtree skills; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback.go matched subtree skills; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go matched subtree skills; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.gohtml matched subtree skills; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_resolve_test.go matched subtree skills; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go matched subtree skills.
- `updater` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`. skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.gohtml matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_resolve_test.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go matched subtree skills/do-work/tools/do-work-cli.
- `installer` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`. skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_feedback.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.gohtml matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_resolve_test.go matched subtree skills/do-work/tools/do-work-cli; skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go matched subtree skills/do-work/tools/do-work-cli.

## Heavy Verification Result

Target and execution revision `086aa28d13f71da13b29a672f480ea9d8acc9a1d`.

- `do-work-cli-integrations` executed in 69 seconds, exit 0.
- `staged-skills` executed in 38 seconds, exit 0.
- `updater` executed in 73 seconds, exit 0.
- `installer` executed in 39 seconds, exit 0.
