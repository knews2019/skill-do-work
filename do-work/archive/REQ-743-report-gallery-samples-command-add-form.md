---
id: REQ-743
title: 'Report gallery: sample format, four shipped samples, list and page command, add form'
status: completed-with-issues
heavy_verified_at: 2026-11-18T16:12:15Z
heavy_verified_revision: 764f2d45dc9fc51a63e0f10a5ac73abffb7669f0
commit: 764f2d45dc9fc51a63e0f10a5ac73abffb7669f0
kb_status: pending
created_at: 2026-11-18T15:29:51Z
user_request: UR-183
domain: frontend
prime_files: [_dev/primes/prime-action-files.md]
required_lessons: [skills/do-work/tools/lessons-do-work-update.md]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: true
maintenance: false
batch: report-sample-gallery
related: [REQ-744]
claimed_at: 2026-11-18T15:33:40Z
route: C
estimate:
  p50_active_minutes: 95
  confidence: low
  calculated_at: 2026-11-18T15:35:51Z
  basis:
    - Route C
    - 18-file write set
    - 10 new files
    - 3 subsystems involved
    - 10 acceptance criteria
    - browser evidence
    - persistence changes
    - cross-route regression gates
    - full-suite verification
planning_at: 2026-11-18T15:36:41Z
write_set: [".gitattributes", "README.md", "_dev/tests/install-suite-behavior.sh", "_dev/tests/staged-skills-contract.sh", "_dev/tests/update-script-behavior.sh", "docs/report-gallery/index.html", "skills/do-work-toolbox/SKILL.md", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/actions/help.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work-toolbox/docs/report-gallery/parts-over-time/example.html", "skills/do-work-toolbox/docs/report-gallery/parts-over-time/sample.md", "skills/do-work-toolbox/docs/report-gallery/rhythm-over-cycle/example.html", "skills/do-work-toolbox/docs/report-gallery/rhythm-over-cycle/sample.md", "skills/do-work-toolbox/docs/report-gallery/two-measures-side-by-side/example.html", "skills/do-work-toolbox/docs/report-gallery/two-measures-side-by-side/sample.md", "skills/do-work-toolbox/docs/report-gallery/verdict-with-checks/example.html", "skills/do-work-toolbox/docs/report-gallery/verdict-with-checks/sample.md", "skills/do-work/actions/version.md", "skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go", "skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go", "skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/mutation.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go"]
builder_handback_at: 2026-11-18T15:52:00Z
integration_at: 2026-11-18T15:52:01Z
review_at: 2026-11-18T15:56:21Z
remediation_at: 2026-11-18T15:57:04Z
re_review_at: 2026-11-18T16:00:56Z
completed_at: 2026-11-18T16:12:16Z
release_at: 2026-11-18T16:12:16Z
---
# Report gallery: sample format, four shipped samples, list and page command, add form

## What
Give the detailed AI report workflow a gallery of samples. A sample is a short record of one way to present evidence, with a worked example. The toolbox ships four samples, a project adds its own, one command lists them for the report writer and builds a page for the user, and a new `ai-report gallery add` form saves a sample the user found somewhere.

## Why
The reader should be able to understand the main finding in most reports. The user also wants one page that shows every kind of report the suite can make.

## Decisions Already Made
The user settled these on 2026-11-18. Do not reopen them.

- **Project samples live in `docs/report-gallery/`.** The user accepted this folder on one condition: it is deployed into a project by `just run-do-work-update`. A new root folder was rejected because `folder-structure.md` allows only the listed root folders. `ai-reports/gallery/` was rejected because `reportIndexWalk` (`report_index.go`) treats every folder under `ai-reports/` as a report.
- **Shipped samples live in `skills/do-work-toolbox/docs/report-gallery/<sample-name>/`.** They reach a project as part of the toolbox files on every update. Nothing copies them into the project.
- **Four rebuilt patterns ship on the first day** (list below).
- **Adding a sample is an agent form, not a page control.** A page opened from disk cannot write into the project.

## Detailed Requirements

### Sample record
- One folder per sample, named with a two-part or longer kebab name that plain-text search finds.
- `sample.md` holds: a title, one `use_when` sentence that names the reader's question or likely confusion, the source (a link, a path, or `shipped`), the date added, and a short **Borrow** list of two or three moves.
- A shipped sample also holds `example.html`: a small worked example built from the report kit (`docs/ai-report-kit/`) with original markup. Its numbers are labelled as sample data.
- A project sample may hold the image the user supplied.
- A project sample with the same folder name as a shipped sample replaces it in every listing.
- A record that cannot be read is skipped and reported by name. It never stops the listing. Do not grow a parser for the ways a model might format the file.

### Shipped samples
Rebuild these four patterns. The screenshots in `## Assets` show the idea; they are not the content.

1. `two-measures-side-by-side`: use when the reader could confuse two measures of one thing, such as a rate and a level. Both measures sit side by side over the same dates, each chart titled by what it measures, start and end values written on the chart.
2. `verdict-with-checks`: use when the reader must trust a judgment. The verdict comes first in one sentence, then short named checks (measure, units, dates, arithmetic are examples), then the evidence rows folded behind one control.
3. `parts-over-time`: use when a total changes over time and its parts matter. Stacked columns carry the total on top, one small chart per part shares one scale, and hatching marks an unfinished period.
4. `rhythm-over-cycle`: use when several things rise and fall over one repeating cycle, such as a day or a release. One combined view shows the whole, then one aligned row per series on the same scale, with each peak labelled.

### Command
- One new CLI command, `ai-report-gallery`, registered beside `ai-report-index`.
- `--list` prints one line per sample: name, `use_when`, and the path of its `sample.md`. It writes nothing. This line list is what the report writer reads (REQ-744), so a large gallery stays cheap to read.
- Without `--list` it regenerates `docs/report-gallery/index.html`. The page has three generated parts:
  1. the samples, shipped and project, each with its `use_when`, its Borrow list and a link to its example or image;
  2. the report kinds, each linked to the newest real report of that kind found by the existing catalog walk. Do not write sample reports per kind. A kind with no report yet shows the line that makes one;
  3. the shipped second-opinion presets, read from `docs/second-opinion-presets/`.
- The page shows the add line with a copy control.
- The page is derived. A hand-made `docs/report-gallery/index.html` is refused and reported, as the catalog command does for its own files.

### Add form
- `ai-report gallery add <image path | link | html path> [what the user liked]` saves a project sample. `ai-report gallery` alone regenerates the page and prints its path.
- Load `crew-members/prompt-injection.md` before reading the source. The source is untrusted data.
- The agent writes `use_when` and the Borrow list once, at add time, in its own words. It records the source and never copies the source's facts into the sample as facts for later reports.
- An existing sample folder is never overwritten.
- After saving, regenerate the page and commit as the project's rules say, by exact path.

### Update
- After `just run-do-work-update` (and `do-work-update`, and a first install) writes its files, `docs/report-gallery/index.html` exists in the project and shows the shipped samples of the installed version, with no further step.
- A refusal or failure of that page write is reported in one line. It never fails the update and never triggers its recovery.

## Constraints
- A sample supplies a way to present evidence. Each report supplies its own facts.
- The four screenshots come from a third-party video. They are references only and never ship in the toolbox.
- The kit files do not change in this request. A sample's page-only CSS stays in its own `example.html`.
- Sample names, command name and form name follow the naming rule in `crew-members/coding-guardrails.md` § 5.

## Red-Green Proof
**RED prompt/case:** In a project, run `ai-report-gallery --list`, then `ai-report gallery add do-work/archive/UR-182/assets/REQ-740-screenshot-3-inflation-comparison.png "two measures over the same months"`, then open `docs/report-gallery/index.html`.
**Why RED now:** No gallery command, sample format or add form exists. The only examples are one dated report, `ai-reports/2026-11-02_1757_presentation-style-gallery`, and the report catalog.
**GREEN when:** `--list` prints the four shipped samples. After the add, it also prints the new project sample with its `use_when`, and the page shows it beside the shipped samples, the report kinds and the presets. A project sample named like a shipped one replaces it. A broken `sample.md` is skipped and named. After an update in a consumer project the page exists and lists the shipped samples.
**Validation:** User confirmed the inflation example during the UR-182 capture; the listing and update checks are inferred from the user's answers on 2026-11-18.

## Assets
All four are in `do-work/archive/UR-182/assets/`.
- `REQ-740-screenshot-1-daily-activity.png`: time across a day as one combined band chart, then six aligned activity rows on one scale, with labelled peaks. Source of `rhythm-over-cycle`.
- `REQ-740-screenshot-2-nobel-decades.png`: prizes by decade as stacked columns with totals, small charts per category, hatching for the unfinished decade. Source of `parts-over-time`.
- `REQ-740-screenshot-3-inflation-comparison.png`: the annual inflation rate beside the price level over the same months. Source of `two-measures-side-by-side`.
- `REQ-740-screenshot-4-inflation-evidence.png`: a verdict, six named checks, and evidence behind one control. Source of `verdict-with-checks`.

## Dependencies
REQ-744 (the report writer uses the gallery) needs this request finished first.

## Builder Guidance
The folder, the four samples, the command name and the add form are fixed. The field names inside `sample.md`, the page's look, and how the update triggers the page write are the builder's choice. Start from `report_index.go` for the generated-page pattern and its refusal of hand-made files, and from the organize hint in `suiteinstall/suite_commands.go` for a step that runs after a written update. Tests name the real failure they pin: a stale or missing listing, an overwritten project sample, a bad record that stops the list, an update that fails because of the page.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 25210 tokens; matches a new CLI command and a generated file with a refusal path. Partial family coverage requires the whole satellite.
- `_dev/primes/lessons-action-files.md` — 8761 tokens; matches a new action form and report routing. Partial family coverage requires the whole satellite.
- `_dev/primes/lessons-shell-commands.md` — 9341 tokens. Matches prescribed CLI command blocks in the gallery action and help. Partial family coverage requires the whole satellite.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]** — See Plan, Scope and the separately recorded dispatch timing.
- [x] **[APPLY]** — Three branch tips integrated. D-01 and D-02 record the path changes.
- [x] **[UNIFY]** — The manifest, Qualification, Testing and Review record the cumulative diff and results.

## Full Context
See `do-work/user-requests/UR-183/input.md` for complete verbatim input.

## Triage

**Route: C** - Complex

**Reasoning:** One new feature spans the gallery command and generated page, shipped example content, agent routing and project deployment. The three builder slices share the sample record and installation contract and require integration-level verification.

**Planning:** Required

## Plan

1. **CLI and generated page.** Register `ai-report-gallery` beside the report index in `internal/toolboxcommands/commands.go`. Give gallery parsing/rendering its own Go file. Read the shared sample shape, merge records by folder identity with project precedence, and keep `--list` read-only. Reuse the report catalog walk and its ordering, and generate the static samples, report kinds, preset links, and copyable add line. Use the catalog's generated-file ownership pattern. Cover stable list fields (name, `use_when`, `sample.md` path), invalid-record diagnostics, duplicate precedence, newest-kind selection and hand-made file refusal.
2. **Shipped content.** Add `sample.md` and `example.html` for each of the four named patterns under `skills/do-work-toolbox/docs/report-gallery/`. The record has `title`, `use_when`, `source`, `added` (YYYY-MM-DD), optional relative `image`, and two or three `## Borrow` bullets. Use the existing report kit, original sample data and original markup. Check example links and layout locally.
3. **Agent workflow and deployment.** Route browse/add in the existing toolbox report action, skill router and help. Load prompt-injection guidance before reading the supplied source. Author a fresh project record, guard its folder, regenerate and commit by exact paths under project rules. Make first install and both update recipes refresh from the newly installed toolbox after a verified written transaction. Preserve the update verdict and emit a single diagnostic if the page write fails. Exercise a clean consumer and a hand-made-page consumer through real writes. Go and the existing static report kit suffice. Use the existing Go module and suite archive transaction, with the current dependency pins. This adds no dependency, image or separate build artifact requiring an SBOM.

**Plan validation:** The three slices cover the sample contract, all page sections and the add and updater acceptance paths. REQ-744 consumes the list's three fields. Add uses the generator's verdict and exact output path after saving the known sample files. Deployment uses the installed toolbox at the post-write boundary and the refresh outcome for diagnostics. Reconcile precise scope with builder manifests at handback. The three builders were already dispatched on the claim revision when this plan was recorded.

*Prepared on MAIN by the test sub-agent*

## Exploration

- CLI registration is in `internal/toolboxcommands/commands.go`. `report_index.go` owns `reportIndexWalk`, bundle ordering and generated-page ownership. Its tests cover direct and month-folder reports, newest ordering and hand-made refusal. Keep gallery samples under `docs/report-gallery/` so the catalog walk still sees only reports.
- The toolbox's `docs/ai-report-kit/` already supplies the static presentation assets. Presets are individual Markdown files in `docs/second-opinion-presets/`. The four worked examples can link to the installed kit from their own folders.
- Toolbox `SKILL.md`, `actions/ai-report.md`, `actions/help.md` and `docs/ai-report-guide.md` carry the existing report verbs. The add form takes an untrusted input and loads `crew-members/prompt-injection.md` before reading it.
- `internal/suiteinstall/suite_commands.go` already has post-transaction install/update organize narration. An update invokes the install transaction in-process and the old installed binary runs the first upgrade that delivers this feature. Verify that the real upgrade reads the new toolbox and reaches the refresh, and that failure reporting is once per command and outside recovery.
- Required lessons refreshed against `do-work/lessons-index.md`. Read the existing updater satellite (748 tokens). Its in-process narration lesson applies directly. The action, CLI and shell satellites match but have partial family coverage and exceed the remaining budget. The shell match is for the new prescribed CLI blocks in the action and help. The UI spec informs local-file rendering, keyboard access and narrow layouts.

*Explored on MAIN by the test sub-agent*

## Scope

**Files I will touch:**

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go` — CLI registration.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go` — sample discovery, list and generated page.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go` — request-specific CLI contract tests.
- `skills/do-work-toolbox/docs/report-gallery/two-measures-side-by-side/sample.md` — record.
- `skills/do-work-toolbox/docs/report-gallery/two-measures-side-by-side/example.html` — paired measures.
- `skills/do-work-toolbox/docs/report-gallery/verdict-with-checks/sample.md` — record.
- `skills/do-work-toolbox/docs/report-gallery/verdict-with-checks/example.html` — verdict and evidence disclosure.
- `skills/do-work-toolbox/docs/report-gallery/parts-over-time/sample.md` — record.
- `skills/do-work-toolbox/docs/report-gallery/parts-over-time/example.html` — totals and components.
- `skills/do-work-toolbox/docs/report-gallery/rhythm-over-cycle/sample.md` — record.
- `skills/do-work-toolbox/docs/report-gallery/rhythm-over-cycle/example.html` — aggregate and aligned series.
- `skills/do-work-toolbox/SKILL.md` — gallery routing.
- `skills/do-work-toolbox/actions/ai-report.md` — browse and add forms.
- `skills/do-work-toolbox/actions/help.md` — discoverability.
- `skills/do-work-toolbox/docs/ai-report-guide.md` — user-facing usage.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go` — post-write gallery refresh.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go` — install/update hook outcomes.

- `README.md` — integrated file required by the handback (D-01).
- `_dev/tests/install-suite-behavior.sh` — integrated file required by the handback (D-01).
- `_dev/tests/update-script-behavior.sh` — integrated file required by the handback (D-01).
- `skills/do-work/actions/version.md` — integrated file required by the handback (D-01).
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` — integrated file required by the handback (D-01).
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go` — integrated file required by the handback (D-01).
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/mutation.go` — integrated file required by the handback (D-01).
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go` — integrated file required by the handback (D-01).

- `.gitattributes` — exclude the checkout gallery from consumer archives.
- `docs/report-gallery/index.html` — this repository’s user-openable result.

- `_dev/tests/staged-skills-contract.sh` — confirm the new repository-owned gallery is export-ignored while continuing to reject other retired docs.

**Out of scope:** edits to the report kit, third-party screenshots in the shipped package, and REQ-744 report-writer matching.

**Acceptance criteria (restated from REQ):**

- [x] Four shipped records have the shared metadata and two or three Borrow bullets, plus original worked examples using the kit and labelled sample data.
- [x] Each worked example demonstrates the named pattern and renders with its installed kit links.
- [x] Read-only `--list` yields one line per sample with name, `use_when` and an existing `sample.md` path.
- [x] Project records override matching shipped folders, and invalid records are named and skipped.
- [x] Generated page includes examples/images, newest real report by kind or creation line, shipped presets and copyable add line.
- [x] A hand-made gallery page is reported and preserved byte-for-byte.
- [x] The agent add form reads the source under the untrusted-input rules, creates a fresh project record and keeps source facts out of later reports.
- [x] The form leaves existing sample folders intact, refreshes the page and commits by exact path under the consumer's rules. Browse prints the generated path.
- [ ] First install and real writes through both updater recipes create a page with the installed version's shipped samples.
- [ ] Refresh refusal or failure is reported once and leaves the successful update and its recovery state intact.

Declared on MAIN after builder dispatch. Judge any divergence against the three handback manifests and the acceptance criteria.

## Pre-Flight

**Git:** MAIN stayed at `fe0b06b74be7bedcafc39056f86c6da9be544c4b`. The pre-existing current-run files and the requested REQ evidence were under `do-work/`. No implementation paths were dirty.

**Tests baseline:** The focused catalog and install/update probe launched and passed all four named Go tests. `advance` recorded exit 0 in `do-work/working/baseline.json`; the run directory has a copy.

**Dependencies:** Go and ShellCheck cleared the canonical gate's toolchain floors. Both Go test stages executed.

**Canonical gate:** `bash _dev/tests/maintainer-verify.sh` exited 0 on the first attempt, 2026-11-18T15:37:54Z to 15:40:26Z (152 s). The request-bound `preflight` and `green-gate` records are both `satisfied`. The gate record is for the exact argv and MAIN revision above.

**Ordering:** The builders started around 15:34 UTC. At 15:35:40 UTC all three worktrees and MAIN were observed at the same claim revision. Planning was stamped at 15:36:41 UTC. Scope and gate followed dispatch. This is the pre-integration baseline. Only the approximate dispatch time was available here.

**Evidence:** `do-work/runs/work-2026-11-18-153357/REQ-743-preparation.md`.

*Checked by work action*

## Decisions

D-01. Accepted the CLI registration assertion, separate page template, checked publish helper, install fixture, both deployment probes, README and version action. These files exercise or implement the declared command and deployment paths. The update transaction itself stays unchanged. These are handback deltas against the original 18 paths.

D-02. Corrected the scope mirror to the tested handler and installation fixture in place of the speculative update-transaction test path. Added the checkout page and its export exclusion. An already-current retry now creates a missing gallery. It leaves existing files and declined updates alone.

D-03. The staged-skills lane recognized only the previous maintainer docs paths. Add the gallery to the checked export exclusions, reuse that same list for the legacy-runtime scan, and correct the README install FAQ and guide.

## Implementation Summary

**Files changed:**
- `.gitattributes` (modified)
- `README.md` (modified)
- `_dev/tests/install-suite-behavior.sh` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)
- `_dev/tests/update-script-behavior.sh` (modified)
- `docs/report-gallery/index.html` (new)
- `skills/do-work-toolbox/SKILL.md` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/actions/help.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)
- `skills/do-work-toolbox/docs/report-gallery/parts-over-time/example.html` (new)
- `skills/do-work-toolbox/docs/report-gallery/parts-over-time/sample.md` (new)
- `skills/do-work-toolbox/docs/report-gallery/rhythm-over-cycle/example.html` (new)
- `skills/do-work-toolbox/docs/report-gallery/rhythm-over-cycle/sample.md` (new)
- `skills/do-work-toolbox/docs/report-gallery/two-measures-side-by-side/example.html` (new)
- `skills/do-work-toolbox/docs/report-gallery/two-measures-side-by-side/sample.md` (new)
- `skills/do-work-toolbox/docs/report-gallery/verdict-with-checks/example.html` (new)
- `skills/do-work-toolbox/docs/report-gallery/verdict-with-checks/sample.md` (new)
- `skills/do-work/actions/version.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/mutation.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery.go` (new)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go` (new)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go` (new)

**What was done:** Merged the gallery command, four sample records and worked examples, the agent forms, and the suite refresh hook. The generated page consumes project and toolbox records, the catalog walk and installed preset files. The CLI and updater share ownership-checked publication.

## Initial qualification

On the first integrated range, 8ea2e758..68353189, the typed qualify gate was satisfied. The three "unwired" file warnings correspond to a registered Go command, its template called from the writer, and Go-discovered tests. The scope gate reported one declared test file absent. The handler tests and first-install fixture exercise this transaction through its public command. The original declaration remains in the trail.

## Initial review

**Initial verdict: Request changes.** The fresh install, four examples, local Copy and project add passed. Both real 0.324.3 Just entry points accepted the update and left the gallery absent. A hand-made page survived the transition silently. The independent upgrade researcher reproduced those outcomes and confirmed that retrying against the same archive also left the gallery missing. The last behavior is repairable within the installed binary's ownership boundary. The initial probe's path-count check also counts archive diff output. Count only actual gallery warning lines on re-test. Main's generated artifact belongs in the checkout and outside the distribution archive.


## Prior qualification

Final cumulative range 8ea2e758..7a3fdd37. Both the qualification and scope-drift gates are satisfied. The three static warnings identify Go package files by filename. The command table calls handleAIReportGallery, the writer calls renderReportGalleryPage, and go test discovered all five gallery tests. P-A-U covers the preflight plan, three landed branches and merged-source review. The explicit page export exclusion keeps the checkout catalog and real report links outside packaged consumer data. The update hook runs in the installed binary after its successful transaction. Its second-invocation repair for a missing page reads locally installed samples.

## Prior testing

The preflight gate passed in 152 seconds. The combined build passed its first run in 119 seconds. The final candidate passed on the first attempt in 133 seconds with both Go test stages executed, ShellCheck, gofmt, vet and the repository contracts. The typed test-gate recorded satisfied scope-drift, run-blocked-check and green-gate records. The final focused probe required twelve named passes, including four controls from preflight. The gallery RED was the unregistered command and missing page. The repair RED was the page still missing on a same-version retry. Both now pass.

The clean consumer listed four records before add and five after saving the inflation-comparison image. The project sample has one reader question, three original Borrow moves and no source figures. The collision attempt preserved its bytes. Only the record, image and page entered the add commit. Both modes preferred the project override and named the malformed folder. The newest nested proposal and creation line for an absent kind were verified through the built command. The checkout gallery opened in Chrome from disk and all its local links resolved. All four kit examples and the saved image opened. Enter on Copy put the exact command on the macOS clipboard. The checkout and installed gallery fit at 390 px and 1280 px.

Both Just recipes were exercised from a clean 0.324.3 installation and a gallery-aware installation. The historical first invocation installed the suite but omitted the gallery and any warning. The equal-version retry created the missing page. Gallery-aware file-writing updates and fresh installs created the page or preserved a hand-made file with one warning. Cancellation, failed write, existing-page preservation and stable rollback were also checked in the Go command tests. The candidate distribution includes four sample records and excludes the main checkout's generated page and reference screenshots.

## Remediation

The initial review requested changes after confirming a repeat of the old Just recipe still failed to create the page. Fixed that path in f0d67478, verified a failing native test before the fix and a passing one afterward, and reran qualification, canonical testing and consumer checks. The installed previous process owns the first transaction. Its manifest permits four managed module destinations and has no project-page callback. A candidate-only change cannot perform the first invocation's additional write within that contract. Terminal judgment is completed-with-issues.

## Qualification

On 8ea2e758..764f2d45 the canonical qualification and scope-drift gates were both satisfied. The three static filename warnings have the same Go registration, template and test evidence recorded above. The packaging repair asserts the export-ignore attribute for each repository-only docs directory and uses the identical directory set when scanning for retired runtime files. A checked archive contains the four toolbox samples and omits the repository gallery. The README and gallery guide now distinguish suite-managed writes, derived page writes and the previous updater's silent preservation of a hand-made file.

## Testing

The gate after the packaging repair passed on attempt one in 124 seconds. The typed focused and repository gates passed on 764f2d45dc9fc51a63e0f10a5ac73abffb7669f0. All twelve named focused cases passed. The earlier deployment, clipboard and collision evidence applies to the same gallery runtime. The real browser, git archive contents and historical two-invocation transition were rechecked on this final tree.

## Review

**Overall: 83%** | 2026-11-18T16:07:52Z

| Dimension | Score |
| --- | --- |
| Requirements | 80% |
| Code quality | 95% |
| Test adequacy | 97% |
| Scope | 100% |
| Risk | Medium, limited to first deployment |
| Acceptance | Partial |

**Important finding.** A consumer upgrading from 0.324.3 must invoke the updater again or run `ai-report gallery` before opening a new gallery. That first updater also preserves a hand-made page silently. Both behaviors were reproduced with the prior installed executable. — impact-user-visible → report only

**Minor findings:** None outstanding.
**Acceptance:** Partial across the real local file, clipboard, project add, fresh install and both generations of updater. The two deployment criteria remain open specifically for the historical first invocation.
**Restatement sweep:** Checked report forms, help, report guide, release prompt, version action, native command registration, consumer probes and both managed Just routes. The guide and README now describe the transition. There were no earlier successful REQs in this wave.
**Suggested testing:** No additional local check is needed to establish the limitation. A published prebuilt binary on other operating systems was not exercised.
**Review method:** Sole-integrator review per the dispatch. The independently prepared upgrade analysis was corroborated with clean consumers.

The first heavy run identified an outdated repository-only docs classification. The packaging change keeps the archive assertion, uses its directory set when scanning for retired runtime, and fixes the install FAQ’s stale write-boundary claim. No remaining packaging or gallery-code finding was identified on re-review.

## Lessons Learned

- Seed a rollout check with the previous release's real Just recipe, launcher and executable. A version-advancing test driven by the new binary only establishes behavior for later upgrades.
- Count stderr lines beginning with the actual warning. A full managed diff can contain copies of the warning text or the page path.
- A repeat update can safely fill a missing derived page after the previous executable has installed the new renderer. A file already present during that equal-version check remains in place.

## Orientation

[MAP CHANGED] The toolbox now exposes a shared sample record, a cheap listing for report writers and a local gallery for readers. Browse and add route through the existing report action. The updater calls the same generator on installed samples, outside the suite transaction. The action prime still resolves its relevant report and command references.

## Heavy Verification Plan

Base `8ea2e7585d9e627f0c0d3150f00fdb2bdf0fa667`. Target `764f2d45dc9fc51a63e0f10a5ac73abffb7669f0`.

- `queue-kanban-javascript` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`. Reasons: coverage is uncertain for: .gitattributes, docs/report-gallery/index.html.
- `queue-kanban-browser` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`. Reasons: coverage is uncertain for: .gitattributes, docs/report-gallery/index.html.
- `do-work-cli-integrations` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`. Reasons: coverage is uncertain for: .gitattributes, docs/report-gallery/index.html.
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. Reasons: coverage is uncertain for: .gitattributes, docs/report-gallery/index.html.
- `updater` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`. Reasons: coverage is uncertain for: .gitattributes, docs/report-gallery/index.html.
- `installer` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`. Reasons: coverage is uncertain for: .gitattributes, docs/report-gallery/index.html.

## Heavy Verification Result

Target and execution revision `764f2d45dc9fc51a63e0f10a5ac73abffb7669f0`. All six lanes ran at the target. Chrome was explicitly supplied to the browser lane. On the prior revision, staged-skills rejected the checkout artifact, and the initial browser discovery skipped until the installed Chrome was named.

- queue-kanban-javascript — exit 0, executed, 6 s, fingerprint_uncertain.
- queue-kanban-browser — exit 0, executed, 70 s, fingerprint_uncertain.
- do-work-cli-integrations — exit 0, executed, 54 s, fingerprint_uncertain.
- staged-skills — exit 0, executed, 30 s, fingerprint_uncertain.
- updater — exit 0, executed, 67 s, fingerprint_uncertain.
- installer — exit 0, executed, 27 s, fingerprint_uncertain.

## Timing

Observed 2026-11-18T15:54:39Z to 2026-11-18T16:07:38Z: 12m 59s total, 6m 16s attributed across 3 events, 6m 43s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 6m 16s | 3 |

Slowest command: verification-gate / final integrated candidate gate, 2m 13s, exit 0, bash (2 argv tokens).
