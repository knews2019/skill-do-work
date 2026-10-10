## Review: REQ-696

**Approve** — all six stale wording sites now match the code and the core routing table. The restatement sweep found two more stale "stuck → forensics" lines outside the REQ. Both are report only.
Route A | merge range 8636a267..92a2dca9 (builder commit e8381dc0)

### What's built
- forensics.md:10 and README.md:176 send "how long has in-flight work been quiet" to status and keep forensics for failed or broken work. clarify.md:106 states the full fence rule. The capture.md addendum example uses a bare three-backtick fence. The fixture header and the contract fail message name the validate-feedback forward-row exception.
- Nothing is missing. The six files equal the `write_set`.

### Decisions / risks for you
- None.

### Findings

**Minor:**
- F1 `skills/do-work-toolbox/actions/tutorial.md:260-261`: `"Something seems stuck"` → `do-work forensics`. This is the same stale routing that F9 of REQ-690 (run-status action) fixed in forensics.md and README.md, and core `SKILL.md:40` now routes `stuck` to `./actions/status.md`. Forensics still reports `STUCK-WORK`, so a user who follows it still gets an answer. Fix: route `"Something seems stuck"` to `do-work status`. — impact-negligible → report only
- F2 `skills/do-work/docs/roadmap-guide.md:48`: `Suspect something is *broken* or *stuck* → do-work forensics`. Same class as F1. Fix: `Suspect something is *broken* → do-work forensics; wonder whether in-flight work is *stuck* → do-work status`. — impact-negligible → report only

**Nit:**
- F3 Builder find 1, `skills/do-work/actions/capture-reference.md:196`: it cites clarify.md Step 4 and then restates the fence rule. The builder is right that this is a duplicate. After this REQ, both texts and both Go writers agree: `containedOutsideBytes` (`publication_manifest.go:107-121`) and `containedOutsideText` (`state_apply.go:1052-1070`) both use longest run + 1, a minimum of 3 and no info string. Nothing is wrong today. The only cost is drift risk later. Not impact-critical. — impact-negligible → report only
- F4 Builder find 2, `skills/do-work/docs/forensics-guide.md:3` and `skills/do-work/actions/forensics.md:5`: they keep "stuck" as something forensics detects. That is true. Forensics maps doctor's `STUCK-WORK` finding (forensics.md:49, :54) and lists it in the guide's table (forensics-guide.md:13). Triage D-02 was right to keep both lines, and no change is needed. Not impact-critical. — impact-negligible → report only

### Requirements Checklist

- [x] REQ-690 F9, `forensics.md:10`: reviewer's exact text. The pointer is true because `actions/status.md:10` lists "is it stuck" and `SKILL.md:40` routes `stuck` there — delivered
- [x] REQ-690 F9, `README.md:176`: reviewer's exact text. `do-work status` is a real trigger (SKILL.md:40), and status reports minutes since last activity (status-guide.md:3) — delivered
- [x] REQ-688 F2, `clarify.md:106`: reviewer's exact text. It matches both Go writers (+1, minimum 3, the `> ` + fence line has no info string) and `capture-reference.md:196` — delivered
- [x] REQ-688 F3, `capture.md:136,138`: bare `> ` plus three backticks (checked with `cat -vet`). The enclosing fence is a three-backtick fenced block that opens at `capture.md:131` with the info string `markdown`. A line that starts with `> ` cannot close it, because a closing fence allows only up to three leading spaces. The example renders whole — delivered
- [x] REQ-692 M2, fixture header line 1: reviewer's exact sentence. Line 3 is still the tab-separated header. The parser skips `#` lines (contract.sh:366-370). The exception it names matches core `SKILL.md:36` (`validate-feedback`, `triage feedback` → `../do-work-toolbox/actions/validate-feedback.md`) — delivered
- [x] REQ-692 M2, `staged-skills-contract.sh:792`: reviewer's exact text. The check counts only `` `./actions/$public_action.md` `` in the core routing section, so the message ("through ./actions/ … a ../$sibling_owner/ forward row is allowed") now states exactly what the check tests. `$sibling_owner` is bound at :776. The logic is unchanged — delivered
- [x] Constraint: prose only, no new test — delivered

### Acceptance Testing

**Result: Pass** (implementation and integration stages)
- `bash do-work/runs/work-2026-10-10-192201/REQ-696-probe.sh` on main: exit 0, `REQ-696 GREEN probe: ok`.
- `bash -n _dev/tests/staged-skills-contract.sh`: exit 0. `git diff 8636a267..92a2dca9 --check`: exit 0.
- The repository gate on main at 92a2dca9 exited 0 in 139 s (the integrator's record in `## Testing`). I did not re-run it.

### Suggested Additional Testing

- Integration, heavy lane: the `staged-skills` heavy lane is unassessed in this review, and the integrator's heavy drain owns it. Run `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` to exercise the edited fixture header and the line-792 message.

### Restatement sweep (narrow: two rules only)

- Fence rule ("one backtick longer than the longest backtick run, minimum three, no info string"): I grepped `skills/` and `README.md` for `longest backtick`, `backtick run`, `longer than the longest`, `longer fence`, `four-backtick`, and contained fences of 4+ backticks or with an info string (`^\s*> `{3,}[A-Za-z]`, `^\s*> `{4,}`). The only restatements are `clarify.md:106` and `capture-reference.md:196`, and both agree. No contained example still carries an info string or an oversized fence. Stale: none.
- "Stuck" routing: I grepped `skills/` and `README.md` for `stuck`. Stale: `tutorial.md:260-261` (F1) and `roadmap-guide.md:48` (F2). These lines say only that forensics *detects* stuck work, and that is true: `forensics.md:5,49,159`, `forensics-guide.md:3,13`, `roadmap.md:5,124,278`, `roadmap-guide.md:5`, `stray-check.md:16`, `stray-check-guide.md:5,59`, `work.md:580`.
- This is not the wave's last integration, so no inherited sibling elements were swept.

### Scores (on the record — not the headline)

**Overall: 100%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | 6/6 sites, reviewers' exact text |
| Code Quality | 100% | Each new line is true against the code it describes |
| Test Adequacy | N/A | Prose plus one message string. The REQ forbids a new test, and the probe plus the gate stand in |
| Scope | 100% | Files changed equal `write_set`. D-05 and D-06 are recorded |
| Risk | None | Message text and a `#` comment only. The check logic is unchanged |
| Acceptance | Pass | Probe and syntax checks, gate on the merged tree |

### Follow-ups created
- None (4 findings report only)

## Review

**Overall: 100%** | <integrator stamps>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** F1 `skills/do-work-toolbox/actions/tutorial.md:260-261` routes "Something seems stuck" to `do-work forensics`, while core SKILL.md:40 routes `stuck` to status. Fix: point it at `do-work status` — impact-negligible → report only; F2 `skills/do-work/docs/roadmap-guide.md:48` routes "broken or stuck" to `do-work forensics`. Fix: broken → forensics, stuck in-flight work → status — impact-negligible → report only; F3 (Nit, builder find 1) `capture-reference.md:196` restates the fence rule after citing clarify.md. It now agrees with clarify.md:106 and both Go writers, so it is a duplicate only and not stale — impact-negligible → report only; F4 (Nit, builder find 2) `forensics-guide.md:3` and `forensics.md:5` keep "stuck" as something forensics detects (doctor `STUCK-WORK`, forensics.md:49). This is true, D-02 stands, and no change is needed — impact-negligible → report only. Anti-bloat: 0 helpers, options, files or tests that the REQ did not name.
**Acceptance:** Pass — implementation and integration stages: the GREEN probe passes on main, `bash -n` and `git diff --check` are clean, and the repository gate exited 0 at 92a2dca9. The staged-skills heavy lane is left to the integrator's drain.
**Restatement sweep:** redefined the containment fence rule (clarify.md:106; consistent: capture-reference.md:196, publication_manifest.go:107, state_apply.go:1052) and the "stuck" routing wording (stale: tutorial.md:260-261 F1, roadmap-guide.md:48 F2; consistent detection-only wording: forensics.md:5,49, forensics-guide.md:3,13, roadmap.md:5,124,278, stray-check.md:16, stray-check-guide.md:5,59)
**Suggested testing:** 1 item
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*
