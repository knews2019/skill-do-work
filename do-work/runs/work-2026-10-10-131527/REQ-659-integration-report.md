# REQ-659 integration report (frontmatter set and req append-section)

- Release 0.305.103, "do-work-cli Writes REQ Stamps and Sections and Refuses What the Schema Forbids". Archive `do-work/archive/REQ-659-frontmatter-set-and-req-append-section.md` (flat, UR-145 stays open). Finalization commit `4f533248`: cleanup_complete, no blocked paths or reason codes, first try.
- Range: pre 1660977c. First merge a56cd9c4. Review-fix re-merge 8a33190c (recorded as `commit:`). Integrator fix fd8c934d on the builder branch. No conflicts.
- Scope: builder's edit to `corehelpers/commands_test.go` (handler count 21 to 22) accepted (D-17) and declared in Scope and `write_set`.
- Review 91%, Pass, no important findings. I fixed F1 (a `--from` body with its own `## ` heading gave a wrong refusal or wrote a stray section; now `SECTION-BODY-HAS-HEADING`, one assertion), F4 and F5 (two stale stamp-rule sentences in `work-reference.md`) (D-18). Delta re-review 94%, Pass, all three closed. Report only: F2 symlinked-dir path check (negligible), F3 `set` accepts any `status` value (rule-change), F6 D-12 CRLF note wrong (corrected in D-18), F7 small duplication (negligible), N1 a body with an open fence or comment hides later sections (user-visible; fix = refuse unless the visible section count grows by one).
- Gate (`DO_WORK_FAST_STAGE_REUSE=off`) exit 0 at a56cd9c4 (151 s, load 2.22) and 8a33190c (149 s, load 4.48). Probe exit 0 both times; green record at 8a33190c via `record-green-gate`.
- Heavy lanes at 8a33190c, all exit 0, executed: do-work-cli-integrations 68 s, staged-skills 39 s, updater 67 s, installer 30 s. Earlier drain at a56cd9c4 also green.
- Timing events: handback-merge 1, verification-gate 4, review 2; folded into `## Timing`. Builder-work skipped (hand-back had landed).
- The merged commands wrote this REQ's stamps and six of its sections; no command started refusing in a new way after the merge.
- Lesson: one `alternate-writer-contract-drift` bullet in `lessons-do-work-cli.md` (flat-path GitHub link; REQ-658 must re-point it when it closes UR-145); index row 19697 tokens. `contract-regressions.sh` passes after finalization.
- Cleanup: builder worktree removed, branch deleted with `-d`, pruned; drain checkout removed.
- Left dirty: siblings' untracked `REQ-*-handback.md` files (not mine) and this report (untracked, for the coordinator).
- Report-only discovered tasks: REQ-658 should reuse `ResolveTarget`, not a new resolver; `fan-out-reference.md:161` does not name `frontmatter set` (hand path valid).
