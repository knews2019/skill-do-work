## Review: REQ-703

**Approve**: the presets, the free-text direction, the directed names and the page marks all work as the Plan designed them. One prompt-wording gap can make bold directed runs fail the backlink check.
Route C | merge `dd3639bf` (builder commit `f06e82e0`), range `f03fbc11..dd3639bf`

### What's built
- `ai-report-second-opinion` accepts `--preset <name>` and `--editor-direction '<text>'` in the Session prepare, Targeted and Batch forms. `--finalize` refuses both flags and restores the direction from the prepare record.
- A project preset wins over the shipped one. An unknown name lists the presets that exist. A symlinked preset file is refused. Directed pages are named `index-<cli>-<model>[-unverified]-<preset|custom>[-<slug>][-N].html`. They carry two meta tags and a marker-wrapped header block. Rows record `preset`, `preset_file` and `editor_direction`. Plain runs keep their names, and their pages get no marks.
- Still missing: the release files (version, changelogs, version.md) and the lesson bullet. Finalization writes these, so this review does not cover them.

### Decisions / risks for you
- **D-04 (ESCALATE, kept on purpose):** a project preset that is not committed gets refused by the outside-bundle Git check. The probe confirmed the refusal text: `outside-bundle changes: ai-reports/second-opinion-presets/fresh.md`. The user guide says this plainly: "commit a new one before using it, because the run refuses uncommitted changes outside the bundle". The agent-facing action covers it only through its general rule "Start with no Git changes outside that bundle" (F3). Value: the one check that catches backend writes outside the bundle stays whole. Risk: the first run after you write or edit a preset is refused until you commit it. The refusal names the file.
- **F1 below** asks whether the direction may override the launcher's backlink and output rules. The UR's D3 exempts only the file BOUNDARY. The launcher still enforces the backlink anyway.

### Findings

**Important:**
- F1: the prompt's DIRECTION paragraph says the direction "overrides ROLE, GOAL, every RULE and PROCESS, everything except BOUNDARY" (`skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md:17`). Read literally, that sentence also lets a direction override OUTPUT, BACKLINK and FINAL REPLY. The launcher still requires the backlink (`secondOpinionBacklink`, `report_second_opinion.go:793`). So a `bold` run ("Be bold with layout") whose editor opens the page with a hero section gets refused with "new page must open its body with Original report: …". The page stays staged, so nothing is lost, but the run fails. A direction can also suppress the `Variation slug:` line, which costs only the slug. Smallest fix: "…everything except OUTPUT, BACKLINK and BOUNDARY, which still apply in full", plus the matching words in `ai-report-guide.md:107` ("except the file boundary and the launcher's backlink"). — impact-user-visible → report only

**Minor:**
- F2: stale restatement. `skills/do-work-toolbox/docs/ai-report-guide.md:91` still says "`ai-report index` lists every directory under `ai-reports/` as a report bundle". Since D-03 the walk skips `second-opinion-presets/` (`report_index.go:94`). The code comment was updated, the guide was not. — impact-negligible → report only
- F3: the agent-facing action (`skills/do-work-toolbox/actions/ai-report.md:73`) names `preset_file` but never says a new or edited project preset must be committed first. Only the guide says it. The general "Start with no Git changes outside that bundle" rule (`ai-report.md:71`) covers it, and the refusal names the path. — impact-negligible → report only
- F4: masked-surface edge in the strip (`report_second_opinion.go:1469-1474`). The comment surface masks raw-text spans found on the unmasked page. So a real comment that contains `<script>`, such as `<!-- <script> -->`, masks everything up to the next real `</script>`. A copied header block in that range is not stripped, and the run then fails the backlink check. The probe confirmed this. `secondOpinionBanner` (`:1439-1446`) has the same pattern, so the problem is older than this REQ. — impact-negligible → report only
- F5: the meta strip regex only matches when `name` is the first attribute (`report_second_opinion.go:1491`). If an editor rewrites a copied meta as `<meta content="x" name="ai-report-second-opinion-preset">`, the tag survives and the page carries two preset metas. The probe confirmed this. It matters more once REQ-704 (linked iteration tree) makes base pages carry metas. — impact-negligible → report only
- F6: the finalize path trusts the sidecar record. `record.Direction.Preset` is not checked again against `^[a-z0-9][a-z0-9-]*$` before it reaches `secondOpinionNames` (`report_second_opinion.go:606`, `:640`). The preservation check exempts the sidecar, so the session editor could change it. This is the same class as the existing unchecked `record.Model`, and `rootedPublishFile` keeps writes inside the repository. Fix: run the preset regex and `secondOpinionModelPattern` on the record in finalize. — impact-negligible → report only
- F7: two test gaps. Nothing tests the symlinked-preset refusal (D-02), and nothing tests the strip-then-backlink order with a copied header block in the publication tail (D-09). The reviewer probe confirmed both behave correctly. — impact-negligible → report only

**Nit:**
- F8: `secondOpinionVariationSlug` takes the first `Variation slug:` line (`report_second_opinion.go:909`), but the prompt says "end with" that line. Taking the last match would survive an earlier quoted mention. The Plan's D-08 chose first-match. — impact-negligible → report only
- F9: D-10 places the header block at offset 0 when the page has no body tag, which puts it before `<!doctype html>` and switches the page to quirks mode. This can only happen on unverified pages, which get no judge. — impact-negligible → report only
- F10: `ai-report.md:73` says the variation goes "after the model in both names". For the unverified form the code puts it after `-unverified`, which the REQ shape `[-unverified]-<preset>` also says. Spell out the unverified example. — impact-negligible → report only

### Requirements Checklist
- [x] D1: both flags work alone and together in Session prepare, Targeted and Batch. Preset text comes first, then free text. `--finalize` refuses both flags and reads them from the record. Delivered.
- [x] D2: shipped `bold` and `visual` presets exist, the project preset wins, an unknown name lists the presets found, and the row has `preset` and `preset_file`. Delivered. The name regex, the symlink refusal and repository-relative `preset_file` were all probed.
- [x] D3: the prompt has a DIRECTION paragraph, a plain run gets `None.`, and the direction overrides everything except BOUNDARY. Delivered as worded. The wording problem is F1.
- [x] Empty or whitespace-only direction and a direction containing `{{` are refused with clear messages. Delivered.
- [x] Meta tags plus a marker-wrapped header block, HTML-escaped and idempotent. Delivered. The probe confirmed `-->` and `<script>` are escaped, idempotence holds with and without `</head>`/`<body>`, and a plain decoration strips back to the exact original bytes.
- [x] Directed name with `custom` and a cleaned slug, reply on the same stem. Delivered.
- [x] `AIReportSecondOpinionRun` gains three omitempty fields. Delivered.
- [x] One prompt helper replaces both fills and refuses a directed run on a template without the placeholder. Delivered.
- [x] Reduced-motion rule added to Rule 8. Delivered.
- [x] Plain run keeps its names and adds no marks. A finished page is never deleted, because the backlink-fail path keeps the stage. The editor boundary is unchanged. Delivered.
- [x] Docs: `ai-report.md`, the guide and `SKILL.md` show the flags and the `second opinion preset` trigger. Delivered.
- [ ] One release: N/A to this review, because the integrator's finalization writes it.
- [x] Old prepare records stay readable. Delivered: `omitzero` works with the Go 1.24 module floor, and the zero value means a plain run.

### Acceptance Testing

**Result: Pass** (Implementation and Integration stages)
- `go test -count=1 -run '^(TestSecondOpinionPresetResolutionAndArgumentRules|TestSecondOpinionDirectedRunNamesAndMarksPage|TestAIReportIndexSkipsSecondOpinionPresetsFolder)$' ./internal/toolboxcommands`: 3 PASS, 1.1 s.
- Regression selection `SecondOpinion|AIReportIndex|TestHandlersRegisterCanonicalToolboxCommands`: ok, 15.5 s. `resultmodel.TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence`: ok.
- Reviewer probes ran against a scratch copy of the module, so the project tree was not touched. They checked:
  - the copied-header strip, after which the backlink passes;
  - idempotence with no body and with no head;
  - escaping of `-->` and `<script>`;
  - a meta inside a script surviving a plain run;
  - the uncommitted-preset refusal (D-04);
  - the symlinked-file refusal;
  - a symlinked folder, which is refused as untracked;
  - the uppercase-name refusal and the `../x` refusal;
  - slug cleaning: punctuation only gives an empty slug, more than 4 words are cut to 4, more than 40 characters are cut, and Unicode letters are dropped.
  
  Two probes failed, as recorded in F4 and F5.
- Gate: the integrator's maintainer-verify run at `dd3639bf` was green. This review did not re-run it, and it did not run the heavy lanes.

### Suggested Additional Testing
- Live acceptance (unassessed): run `ai-report second-opinion --backend codex --preset bold latest` against a real bundle. Confirm the `-bold-<slug>` name, that the header block is readable, that the judge passes, and that the banner is inserted.
- Deployment (unassessed): after install, confirm `.claude/skills/do-work-toolbox/docs/second-opinion-presets/` is present. The `staged-skills` heavy lane covers this.
- Manual: open a directed page. Check that the blank line between preset text and free text shows (`pre-line`) and that the header block reads well in light and dark themes.
- Edge: in a real `bold` run where the editor restructures the opening, watch for the F1 backlink refusal.
- Session form in a real agent: pass `--editor-direction` with quotes and a newline, then check the meta `content` and the header text.

### Domain Review (backend)
| Criterion | Result |
|-----------|--------|
| Handles invalid input | Pass. Every bad flag value is a typed usage error, not a crash: unknown preset, bad name, empty or `{{` direction, symlinked file, empty preset file. |
| Auth enforced | N/A. Local CLI. |
| No data leaks | Pass. `preset_file` is repository-relative when the file is inside the repository. An absolute path appears only for a preset outside it. |
| Idempotency | Pass. Decoration is idempotent, and re-publishing allocates a new `-N` name. |
| Existing tests pass | Pass. The regression selection is green. |

### Anti-bloat (items the REQ and Plan did not name)
Six items, all keep, none cut. No decorative tests: each new test has a `Pins:` comment that names a real failure, and the catalog-skip test reproduces the RED the builder recorded.
1. Test helper `secondOpinionFixturePresets`: keep. Two tests use it, and it keeps the fixture hermetic.
2. Option fields `preset` and `editorDirection` on `secondOpinionOptions`: keep. The parser needs them.
3. `secondOpinionVariationSlug` as its own function, extracted from D-08's prose: keep. It is directly testable and is the seam for REQ-704 (linked iteration tree).
4. `secondOpinionDecorate` as its own function, where D-09 said "in the publication tail": keep. It is the strip and decorate seam that REQ-704 extends.
5. D-15's `presets found: none` fallback: keep. It is one line and avoids a dangling colon.
6. `aria-label` and the `white-space:pre-line` style on the header aside (D-12): keep. They make the label accessible and keep the blank line visible.

### Scores (on the record, not the headline)

**Overall: 93%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Every REQ and UR item for this REQ was delivered. The release belongs to finalization. |
| Code Quality | 88% | Clean and follows the house marker pattern. F1 prompt scope, F4 and F5 strip edges, F6 record trust. |
| Test Adequacy | 85% | RED and GREEN recorded for all three tests. Symlink refusal and copied-header publication are untested (F7). |
| Scope | 100% | 11 of 18 declared paths touched. The 7 untouched are finalization's. |
| Risk | Low | No data loss. The backlink-fail path keeps the staged page. |
| Acceptance | Pass | Implementation and Integration stages. |

### Follow-ups created
None (10 findings report only)

---

## Review

**Overall: 93%** | <REVIEW_AT>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 88% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token, the durable audit record the judgment mandates):**
- F1: the prompt DIRECTION paragraph (`ai-report-second-opinion-prompt.md:17`) says "everything except BOUNDARY", so a direction may override OUTPUT, BACKLINK and FINAL REPLY. The launcher still enforces the backlink, so bold restructures can be refused. Fix: "except OUTPUT, BACKLINK and BOUNDARY". — impact-user-visible → report only

**Minor findings:**
- F2: stale `ai-report-guide.md:91` "lists every directory under `ai-reports/`"; the presets folder is now skipped. — impact-negligible → report only
- F3: `ai-report.md:73` does not say a project preset must be committed first (D-04); only the guide does. — impact-negligible → report only
- F4: strip masked surface (`report_second_opinion.go:1469`) lets a comment containing `<script>` hide a copied header; the same pattern already exists in the banner. — impact-negligible → report only
- F5: the meta strip (`:1491`) misses copied metas whose `name` is not the first attribute. — impact-negligible → report only
- F6: finalize does not re-check `record.Direction.Preset` from the editable sidecar (`:606`); same class as the existing `record.Model`. — impact-negligible → report only
- F7: no tests for the symlinked-preset refusal or for copied-header strip-then-backlink in publication. — impact-negligible → report only
- Nits F8 to F10: first-match slug, offset-0 header before the doctype, unverified-name wording in `ai-report.md:73`. — impact-negligible → report only

**Acceptance:** Pass. Implementation and Integration stages: 3 new tests plus the 15.5 s regression selection green, and reviewer probes on a scratch copy. Deployment and Live acceptance unassessed.
**Restatement sweep:** redefined: published second-opinion name shape (directed `-<slot>[-<slug>]` segment), prompt placeholder set (`{{EDITOR_DIRECTION}}`), `--finalize` argument rule, catalog-walk bundle rule (presets folder skipped), opening of a published directed page (header block before the backlink). Stale: `ai-report-guide.md:91` (F2). All other restatements agree: `ai-report.md:71-73`, `ai-report-guide.md:72-74,107,111`, `ai-report-reference.md:148`, `SKILL.md`, launcher test `:252`.
**Suggested testing:** 5 items
**Follow-ups created:** None (10 findings report only)

*Reviewed by review-work action*
