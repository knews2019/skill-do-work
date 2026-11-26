---
id: REQ-764
title: 'Add Product Launch Brief to Sample Designs'
status: completed
created_at: 2026-11-26T19:38:31Z
user_request: UR-194
domain: ui-design
prime_files: []
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: false
maintenance: false
claimed_at: 2026-11-26T19:44:59Z
route: B
estimate:
  p50_active_minutes: 45
  confidence: medium
  calculated_at: 2026-11-26T19:47:07Z
  basis:
    - Route B
    - 5-file write set
    - 4 new files
    - 2 subsystems involved
    - 7 acceptance criteria
    - browser evidence
    - cross-route regression gates
write_set:
  - skills/do-work-toolbox/docs/report-gallery/product-launch-brief/sample.md
  - skills/do-work-toolbox/docs/report-gallery/product-launch-brief/example.html
  - skills/do-work-toolbox/docs/report-gallery/product-launch-brief/first-ferry.webp
  - skills/do-work-toolbox/docs/report-gallery/product-launch-brief/preview.png
  - docs/report-gallery/index.html
dispatch_at: 2026-11-26T19:45:22Z
builder_handback_at: 2026-11-26T19:54:36Z
integration_at: 2026-11-26T19:57:15Z
review_at: 2026-11-26T20:03:21Z
kb_status: pending
commit: c8593fae07a841a502f5976218789af06063c2ac
heavy_verified_at: 2026-11-26T20:09:09Z
heavy_verified_revision: c8593fae07a841a502f5976218789af06063c2ac
completed_at: 2026-11-26T20:09:36Z
release_at: 2026-11-26T20:09:36Z
---
# Add Product Launch Brief to Sample Designs

## What

Add a reusable report named **Product Launch Brief** to the report gallery’s **Sample designs** section. Let readers try a product example, then follow that same example through the people, work, and decisions needed to launch it. Give the sample the same preview, open-report, and reuse options as its gallery neighbors.

## Context

The design reference is the password-protected German proposal for a daily spot-the-difference game. It follows a collectible card from its illustrated front and article excerpt through the player’s collection to the daily editorial and design workflow. Use original public-facing imagery and copy for the gallery sample. [Design reference](https://zingy-treacle-7f4c65.netlify.app/)

The discussion suggested a compact version with one flippable card, a collection view showing that card, a production workflow, and three launch questions. A thumbnail showing both faces would make the interaction discoverable. Suggested gallery question: “What will readers experience, and what would it take for our team to run it?” These are design suggestions for the builder.

## Red-Green Proof

**RED prompt/case:** Find Product Launch Brief under Sample designs and open the report.
**Why RED now:** The gallery has no sample by that name.
**GREEN when:** The named entry opens a self-contained public sample. A reader can inspect the product example, recognize it in the operating workflow, and see which decisions remain before launch. The gallery provides a working preview and the usual design reuse prompt.
**Validation:** Inferred during capture

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the applicable agent rules and record the approach.
- [x] **[APPLY]:** Build the sample and its gallery entry.
- [x] **[UNIFY]:** Review the diff and verify the gallery entry, report, and reuse controls.

*Source: [Original request](input.md)*

## Addendum (2026-11-26)

The user approved the design guidance discussed after capture. [Follow-up request](../../user-requests/UR-195/input.md)

### Worked example

- Place the launch premise beside the product example on wide screens. Stack them in reading order on phones, with the example in reach before the operating plan.
- Use a readable editorial type hierarchy, with prose around 65 characters per line. Repeat the example’s title and imagery in the collection and workflow so readers recognize the same item.
- Give the collectible a distinct outline. Present the workflow as ordered steps with the responsible person beside each one. Set the launch questions as a short list. Use one consistent accent for controls and focus.
- If using the proposed two-sided card, give it a labeled button operable by touch and keyboard. Keep the flip brief and user-triggered. Switch faces immediately when reduced motion is enabled. Show both faces together in the static gallery preview.
- Keep the example self-contained and use original local assets. Label the assumptions behind any illustrative figures.

### Reusable guidance for “Use this design”

- Let the reader try one representative part of the product.
- Follow that same example through the people and steps needed to deliver it.
- End with the launch decisions those steps expose.

### Verification

Open the sample at desktop and phone widths. Follow the same item from the demo into the workflow and identify a launch question that arises from one of its steps. For a two-sided card, check both directions by keyboard and touch, verify visible keyboard focus, and verify that reduced motion switches faces immediately. Check that the gallery preview communicates the interaction before the reader opens the report.

## Addendum (mid-run)

The user shared design preferences for the worked example and requested an independent review of the instructions.

> ```
> ^ see here the advice that I got for good design
> ```

> ```
> also use Claude Opus to review your design work instructions
> ```

The screenshot favors whitespace, useful icons, and concise language. Claude Opus 5.5 reviewed the builder instructions. The coordinator used the review to clarify the two-sided card, the preview showing both faces, and the reusable presentation moves. The active builder received both updates. See the run’s [updated design brief](https://github.com/knews2019/skill-do-work/blob/main/do-work/runs/work-2026-11-26-194510/REQ-764-brief.md).

## Triage

**Route: B** - Medium

**Reasoning:** The requested sample, layout, interaction and gallery controls are specified. The remaining discovery is where to register the sample and how the existing gallery generates its preview, report link and reuse prompt.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

- The gallery discovers `sample.md` records by folder slug from the shipped toolbox catalog and the project gallery. Project records override matching shipped slugs. Use the shipped catalog for this reusable example and the repo's generated gallery to expose it here.
- The existing renderer uses a sample's `example.html` for Open report, prefers `preview.png` for the preview, and derives Use this design from the record's title and Borrow bullets. Run `ai-report-gallery` to regenerate `docs/report-gallery/index.html`. Its existing controls require no renderer edits.
- Follow the optional UI component spec for native controls and keyboard access. Use an original local illustration and repeat the same example in the collection and editorial workflow. The approved addendum sets the card behavior and responsive verification.
- The lessons index has no entry for gallery sample content or the expected paths. The request has no captured required lessons.
- The focused baseline covers gallery sample discovery, preview/report links, starter-prompt resolution, and the installed gallery rebuild. The builder handback should carry browser evidence for the new sample itself.

*Prepared in the main checkout*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/docs/report-gallery/product-launch-brief/sample.md` (new) — title, reader question and reusable design guidance.
- `skills/do-work-toolbox/docs/report-gallery/product-launch-brief/example.html` (new) — standalone launch brief and card interaction.
- `skills/do-work-toolbox/docs/report-gallery/product-launch-brief/first-ferry.webp` (new) — original local artwork for the example, reconciled from the provisional PNG path at integration.
- `skills/do-work-toolbox/docs/report-gallery/product-launch-brief/preview.png` (new) — gallery preview showing both faces together.
- `docs/report-gallery/index.html` (regenerate) — publish the new card alongside existing Sample designs.

**Files I will NOT touch:** Gallery renderer and shared templates, unrelated samples and reports, and the private design reference. Release metadata belongs to integration.

**Acceptance criteria (restated from REQ):**
- [x] Sample designs lists Product Launch Brief with a preview that opens at full size, a working Open report link and the usual Use this design starter.
- [x] The sample runs with original copy and local artwork, and labels assumptions behind any illustrative figures.
- [x] The launch premise sits beside the product example on desktop. Phone layouts put the example before the operating plan and keep prose readable.
- [x] Readers try one collectible card and recognize its title and artwork in the collection and production workflow.
- [x] Ordered workflow steps name their owners. The brief ends with a short set of decisions that follow from the workflow.
- [x] The labeled card control works in both directions by touch and keyboard, uses a consistent control and focus accent, and changes faces immediately under reduced motion. The static preview shows both faces.
- [x] The reusable starter explains how to show an example, follow it through the delivery workflow and surface launch decisions.

## Pre-Flight

**Git:** Clean outside `do-work/` at claim revision `f9aeb2fd`. The generated gallery at that revision has no Product Launch Brief entry.

**Tests baseline:** The run's `REQ-764-focused-test.sh` launched and passed all five Go tests for gallery discovery, links, starter resolution and the installed gallery rebuild. The working baseline records exit 0.

**Dependencies:** Preflight reported no dependency warnings.

**Canonical repository gate:** `bash _dev/tests/maintainer-verify.sh` exited 0 on its first run at `f9aeb2fd`. The request-bound `advance` record saved the green baseline for that revision.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `docs/report-gallery/index.html` (modified) — regenerated with the new sample, preview, report link and reusable starter.
- `skills/do-work-toolbox/docs/report-gallery/product-launch-brief/sample.md` (new) — shipped catalog record and Borrow guidance.
- `skills/do-work-toolbox/docs/report-gallery/product-launch-brief/example.html` (new) — offline City Postcards example, native flip button, collection, workflow and decisions.
- `skills/do-work-toolbox/docs/report-gallery/product-launch-brief/first-ferry.webp` (new) — original AI-generated illustration reused in the example.
- `skills/do-work-toolbox/docs/report-gallery/product-launch-brief/preview.png` (new) — authentic offline render with JavaScript disabled, both faces in full, and a visible illustration credit.

**What was done:** The sample carries The First Ferry from an interactive trial through the collection and the newsroom workflow to three linked launch decisions. Integration refined the catalog question to name the product audience, made Borrow usable for products such as booking tools, and added visible AI credits beside the artwork before recapturing the preview.

**Source:** Builder `b61df115c5eb714de0d35aeba33451f9ba5971b7`. Integrated implementation range `f5d6bcca..c8593fae`. The builder’s dispatch was accepted at 19:45:22 UTC and its handback landed at 19:54:36 UTC on 2026-11-26, an observed interval of 9m14s. The lifecycle fields record both instants. The work action's timing recorder uses the invocation time for its endpoint, so the late integration did not add a builder-work event.

*Recorded by work action*

## Qualification

**Passed** — the request-bound qualification and scope-drift records for `f5d6bcca..c8593fae` are both satisfied. The first scope comparison identified the provisional `card-art.png` path. The approved scope allowed another format, so integration reconciled that placeholder to the one actual artwork file, `first-ferry.webp`, in both Scope and `write_set`. The five-file boundary held. All implementation files exist and are reachable through the generated gallery and sample. The range passed `git diff --check`.

*Qualified by work action*

## Testing

**Focused tests:** `advance REQ-764 --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh -- --probe-file do-work/runs/work-2026-11-26-194510/REQ-764-focused-test.sh` returned a satisfied focused-test record with raw exit 0. The five tests exercise gallery discovery, asset and report paths, starter resolution, and gallery regeneration by the installed suite.

**Red-green validation:** The gallery at `f5d6bcca` has no Product Launch Brief entry. Chrome opened the new entry and actual report from the merged gallery at 1440 px and 390 px in both themes. All 28 repository gallery checks passed. The focused offline browser script verified no external requests, the two-face rendered preview, visible keyboard focus, Enter and Space operation, touch in both directions, one accessible face at a time, immediate reduced-motion switches, the owner and decision anchors, and the gallery’s Preview, Open report and Use this design affordances.

**Named design:** Both `ai-report-gallery --resolve 'Product Launch Brief'` and `ai-report-gallery --resolve product-launch-brief` resolved the shipped catalog record. Its three Borrow moves cover selecting a representative trial, tracing it through owners and steps, and linking the decisions. They also fit choosing an appointment time and tracing the resulting booking.

**Fresh install:** Installed a Git archive of the merged implementation into a clean consumer repository. The installed CLI resolved the title to the shipped sample, the generated consumer gallery showed the named entry, and all four installed sample files matched the merged bytes.

**Repository gate:** `bash _dev/tests/maintainer-verify.sh` exited 0 at `c8593fae` in 142 seconds. Contract checks, 427 board Go tests, and 1,017 CLI Go tests passed. `advance` returned a satisfied request-bound green-gate record.

**Heavy verification plan:**
- Range: `f5d6bccabfe22e8bc28254947b901309b4add7c7..c8593fae07a841a502f5976218789af06063c2ac`
- `queue-kanban-javascript`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage is uncertain for: docs/report-gallery/index.html.
- `queue-kanban-browser`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — coverage is uncertain for: docs/report-gallery/index.html.
- `do-work-cli-integrations`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage is uncertain for: docs/report-gallery/index.html.
- `staged-skills`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — coverage is uncertain for: docs/report-gallery/index.html.
- `updater`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage is uncertain for: docs/report-gallery/index.html.
- `installer`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage is uncertain for: docs/report-gallery/index.html.

*Verified by work action*

## Review

**Overall: 98%** | 2026-11-26T20:03:21Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None.

**Minor findings:** None.

**Acceptance:** Pass — local merged sample, catalog selection, gallery preview/Open report/Use this design, accessible two-way interaction, content handoff, and installer checks. The report was exercised offline in Chrome. Touch used browser emulation. The password-protected reference was unavailable to the builder, so the sample was checked against the approved brief and design review.

**Restatement sweep:** Nothing redefined. The shipped sample adds a catalog entry. Its card and two-face preview are specific to City Postcards. Existing gallery and shared design rules are outside the diff.

**Suggested testing:** 0 items.

**Follow-ups created:** None.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Resolving the shipped title exposed the exact copy future authors receive. Applying that copy to a booking example checked that the representative trial and its delivery workflow carry over to a different product.

**What didn't:** A footer-only AI credit fell below the cropped gallery preview. Putting the credit beside the artwork carried provenance through to the offline capture.

**Worth knowing:** Disabling JavaScript renders both card faces in the actual sample. Capture that state to demonstrate the interactive feature in a static preview.

## Orientation

Sample designs now includes Product Launch Brief. A worked City Postcards trial shows how to carry a product example through an operating workflow and into launch decisions. The gallery’s existing design selection resolves the reusable guidance by name.

## Heavy Verification Plan

- Range: `f5d6bccabfe22e8bc28254947b901309b4add7c7..c8593fae07a841a502f5976218789af06063c2ac`
- `queue-kanban-javascript`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage is uncertain for: docs/report-gallery/index.html.
- `queue-kanban-browser`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — coverage is uncertain for: docs/report-gallery/index.html.
- `do-work-cli-integrations`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage is uncertain for: docs/report-gallery/index.html.
- `staged-skills`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — coverage is uncertain for: docs/report-gallery/index.html.
- `updater`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage is uncertain for: docs/report-gallery/index.html.
- `installer`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage is uncertain for: docs/report-gallery/index.html.

## Heavy Verification Result

Target and execution revision `c8593fae07a841a502f5976218789af06063c2ac`.

- `queue-kanban-javascript` executed in 7 seconds, exit 0.
- `queue-kanban-browser` executed in 77 seconds, exit 0.
- `do-work-cli-integrations` executed in 72 seconds, exit 0.
- `staged-skills` executed in 37 seconds, exit 0.
- `updater` executed in 71 seconds, exit 0.
- `installer` executed in 33 seconds, exit 0.

## Timing

Observed 2026-11-26T20:00:14Z to 2026-11-26T20:08:36Z: 8m 22s total, 7m 23s attributed across 2 events, 59s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 7m 23s | 2 |

Slowest command: verification-gate / six selected heavy verification lanes, 5m 01s, exit 0, do-work-cli.sh (20 argv tokens).

## Decisions

- D-01 — Use the builder’s original WebP illustration for The First Ferry. It is the single local artwork file permitted by Scope. The City Postcards story and collection are fictional, and the weekday schedule is labeled as a pilot assumption.
- D-02 — Use the gallery question to name the people who use the product. Use Borrow to ask what a report reader can try to understand a launch decision, then trace one recognizable result through the workflow. A booking example can use a chosen time slot and the resulting appointment.
- D-03 — Keep the card flip and the two-face preview in City Postcards. Put a visible AI credit beside the illustration in the trial, collection and workflow, including in the gallery’s rendered preview.
