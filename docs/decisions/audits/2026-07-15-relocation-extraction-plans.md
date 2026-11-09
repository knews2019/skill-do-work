# Relocation Extraction Plans (plan-only — no extraction performed)

Produced by REQ-019 (UR-003 harness-bloat cleanup). Evidence basis:
`2026-07-15-harness-bloat-audit-phase1-2.md` RELOCATE bucket. Inbound-reference
lists below are grep-verified against the repo at the time of writing, not
recalled. Nothing in this document has been executed — each plan is ready to
run as its own maintenance REQ when the user green-lights it.

Shared context for all three plans:

- **Router impact is small but real:** each package owns a routing-table row, an
  Action Dispatch row, a help-menu block (`actions/help.md`), and a next-steps
  block — all four must go in the same commit as the extraction.
- **`do-work update` handles removals cleanly:** the update flow pre-cleans
  shipped directories before tarball extraction (`actions/version.md`), so
  files removed upstream disappear from consumer installs on their next update.
  The upstream diff step surfaces current-side-only files with a
  "surface rather than assume" rule — each extraction's release notes must name
  every removed path so users can confirm the diff is the expected removal, not
  a lost local customization.
- **Version bump:** each extraction is a **minor** bump with a changelog entry
  listing removed paths + the sibling install command.

---

## Plan A — prompts library → sibling repo `skill-prompt-library`

**What moves** (~19,844 words):

| File | Words |
|---|---|
| prompts/README.md (index) | 942 |
| prompts/prompt-kit-step0..step6 (7 files) | 7,781 |
| prompts/business-vendor-strategic-sort.md | 744 |
| prompts/economics-inference-stress-test.md | 865 |
| prompts/economics-saas-repricing-exposure.md | 883 |
| prompts/tech-inference-architecture-decision.md | 941 |
| prompts/tech-infrastructure-compute-geography-risk.md | 854 |
| prompts/weekly-structural-diff-original.md | 756 |
| prompts/architecture-decisions-log_create-or-expand.md | 2,470 |
| prompts/dark-code-kit_{audit,comprehension-gate,context-layer-generator}.md | 3,608 |

**D-01 (builder decision, user's final call at extraction time):** the whole
library moves, including the 4 dev-adjacent prompts (ADR log, dark-code kit).
Rationale: one clean seam; the runner already supports project-local prompt
dirs, so dev-adjacent prompts stay installable per-project. The 4 files are
marked *may stay* — keeping them means the shipped `prompts/` dir survives with
5 files and the runner needs no default-path change.

**What stays in do-work:** `actions/prompts.md` (the runner, 2,271 words) —
generic list/show/run machinery. Its resolution order changes to:
(1) project-local `prompts/` (already supported, consent-gated),
(2) the sibling library's install location. `docs/prompts-guide.md` stays with
an updated "where prompts come from" section.

**Seams to cut in the same commit** (grep-verified inbound refs):
SKILL.md routing row 22 Notes (mention library is separate), Action Dispatch
row (stays — runner remains routed), `actions/help.md` prompt-library block
(reword: "run prompts from a project or installed library"),
`next-steps.md` prompts block, `actions/version.md` SHIPPED_PATHS (`prompts`
entry stays only if the dir ships non-empty — remove if fully extracted),
`_dev/tests/contract-regressions.sh` (`prompts/README.md` in
active_runtime_docs; retired-prompt `assert_file_missing` — both need updating),
`crew-members/prompt-injection.md` JIT caller list ("prompts run" — marked
illustrative, still worth updating), README.md feature list.

**Migration note (release entry text):** "The prompt library now lives at
`<owner>/skill-prompt-library` — install it per-project with
`do-work install prompt-library` *(new install target to add at extraction
time)* or clone it anywhere and point `do-work prompts` at it. `do-work update`
removes the old `prompts/` contents automatically (pre-clean); your own
project-local prompts are untouched."

**Effort:** small. Zero wiring into the pipeline (audit-verified); the runner
change and test updates are the whole job.

---

## Plan B — interview subsystem → sibling repo `skill-interview`

**What moves** (~12,615 words):

| File | Words |
|---|---|
| actions/interview.md | 4,344 |
| actions/interview-reference.md | 4,257 |
| interviews/work-operating-model.md | 2,535 |
| crew-members/interviewer.md | 664 |
| docs/interview-guide.md | 815 |

**What stays in do-work:** nothing functional. A routing tombstone in the
routing table's Notes for one release cycle ("interview moved to
skill-interview — install: …") is optional but recommended, since `do-work
interview` is muscle memory for existing users; remove the tombstone after one
minor version.

**Seams to cut** (grep-verified): SKILL.md routing row 21 + dispatch row +
argument-hint token; `actions/help.md` Interviews block; `next-steps.md`
interview block; `actions/version.md` SHIPPED_PATHS + update-glob + pre-clean
lines that name `interviews`; `actions/bkb.md` (interview→bkb ingest handoff —
becomes "if skill-interview is installed, its export can be ingested"; bkb
ingest already treats the inbox generically, so this is a doc-line change, not
logic); `actions/prompts.md` + `actions/version.md` incidental mentions;
`crew-members/clear-questions.md` JIT_CONTEXT caller list (illustrative);
`docs/bkb-guide.md` cross-mention; README.md; `AGENTS`/CLAUDE.md structure tree
(`interviews/` line).

**Coupling verdict (audit):** self-contained; the only functional touchpoint is
the optional bkb-ingest handoff, which is data-shaped (drops a document in the
inbox) and survives extraction unchanged on the bkb side.

**Migration note:** "The interview framework moved to `<owner>/skill-interview`.
Existing interview state (`do-work/interview/<template>/`) is *your* data and is
not touched by `do-work update`; the new skill reads the same layout. Install:
clone `skill-interview` into your skills directory; the `interview` verb now
belongs to it."

**Effort:** medium — many small seams, but all mechanical; the subsystem's own
files move unmodified.

---

## Plan C — bkb + dream → sibling repo `skill-knowledge-base`

**What moves** (~14,938 words):

| File | Words |
|---|---|
| actions/bkb.md | 6,711 |
| actions/bkb-reference.md | 2,182 |
| actions/dream.md | 3,471 |
| docs/bkb-guide.md | 1,573 |
| docs/dream-guide.md | 1,001 |

bkb and dream travel together: dream's default resolution targets the wiki
layouts bkb builds (`./kb/wiki`, `./knowledge-base/wiki`), and both share the
memory-hygiene mission.

**What stays in do-work:** `actions/kb-lessons-handoff.md` (2,146 words) — this
is the deliberate boundary. It is do-work's *outbound* offer (promote a REQ's
Lessons Learned into whatever KB exists) and already degrades exactly right
when no `kb/` is present: defers to `kb_status: pending`, points the user at
KB setup, never blocks archival. Post-extraction its "run `do-work bkb init`"
pointer becomes "install skill-knowledge-base and run its init." The
`kb_status`/`kb_entry` REQ frontmatter fields and the Schema Read Contract enum
stay — they belong to the handoff, not to bkb.

**Seams to cut** (grep-verified): SKILL.md routing rows 20 + 29, both dispatch
rows, argument-hint tokens; cleanup routing row 12's carve-out note (names
dream — reword to "memory/wiki phrasings belong to skill-knowledge-base if
installed"); `actions/help.md` Knowledge-base block; `next-steps.md` bkb +
dream blocks; `actions/kb-lessons-handoff.md` bkb-init pointer (as above);
`actions/interview.md`/`interview-reference.md` ingest mentions (moot if Plan B
runs first); `actions/work.md`/`work-reference.md`/`review-work.md`/
`sample-archived-req.md`/`roadmap.md`/`tutorial.md` kb_status-related mentions
(STAY — handoff-owned); `crew-members/prompt-injection.md` + `anti-slop.md`
JIT caller lists (illustrative); README.md; CLAUDE.md "Lessons → Knowledge Base
Handoff" section (reworded, stays — the handoff stays).

**Migration note:** "bkb and dream moved to `<owner>/skill-knowledge-base`.
Your `kb/` directory is project data — nothing in `do-work update` touches it;
the new skill operates on the same layout. The Lessons→KB handoff stays in
do-work and now points at the new skill's init. REQs with `kb_status: pending`
remain valid and are picked up by the new skill's triage."

**Effort:** medium-high — most seams of the three (the handoff boundary must be
reworded carefully), but the handoff's existing graceful degradation means no
pipeline logic changes.

---

## Recommended sequence

1. **Plan A** (prompts) — zero pipeline coupling, immediate ~19.8k-word payload
   drop, proves the removal-release-notes mechanics on the safest package.
2. **Plan B** (interview) — self-contained, exercises the routing-tombstone
   pattern.
3. **Plan C** (bkb+dream) — run last; benefits from A/B having settled the
   mechanics, and from Plan B removing the interview→bkb mention first.

Combined effect if all three run: ~47.4k words leave the repo (~21% of the
shipped skill), the routing table drops 3 rows, and the help menu loses two
blocks — with the work pipeline untouched except for one reworded pointer in
kb-lessons-handoff.md.
