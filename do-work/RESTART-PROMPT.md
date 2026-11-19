```text
Keep the main session as coordinator. Hand REQ-748 to an integrator (it enters with advance REQ-748). Resume the five reviewed round-one reports already on main at 160a96b5. Finish their live gallery check and hand them to the user to mark. Generate the second five through GPT-6.2-astra in T3 Code only after the user has marked all five first-round reports. Use direct repository commands for verification. Use Sol workers for other tasks.
do-work run --coordinate --fan-out 3
These commands are sufficient; everything below them is context.
```

---

## Reference

**Stop point.** The explicit handoff interrupted the trial integrator. All earlier agents finished. T3 confirmed the last task interrupted. The checkout had no conflicted merge or full-gate lock. Nothing was pushed. The sole live claim is the two-round report trial, REQ-748. The queue has no other pending request.

**Saved on main.** Five round-one bundles, the regenerated gallery and catalog, the active REQ and the baseline are in WIP commit `160a96b5563264d45567762c0e7499721b9a11d5`. The full snapshot range is `3ece20aea9daf389ed5c7849eb48c9e7e83bcbae..160a96b5563264d45567762c0e7499721b9a11d5`. The trial integrator staged the reviewed files from five branches. These branch tips were not merged as Git parents. Use the existing snapshot. Each bundle has an actual Astra execution receipt, empty parent/marks, and a recorded passing factual and render review. The generation-digest helper matched the recorded review digest on all five main files during handoff.

**First ten minutes.** The integrator enters using `skills/do-work/tools/do-work-cli.sh --format json advance REQ-748`; the read-only check at handoff named `preflight`. Read the interrupted publication section on the claim. Check the live gallery exposes all five reports with their own figure marks and copyable starters, hand the five links to the user, then pause for genuine feedback. The authoritative feedback file `do-work/report-gallery-feedback.json` did not exist at handoff. A test fixture's likes do not count. Use the same six pinned files from `do-work/runs/work-2026-11-18-191126/trial-source/gallery-input.json` in both rounds. Keep disqualified versions distinguishable; the local trial scratch in the claim records the earlier warm/cold correction. After real marks, generate and verify the next five with lineage, then finalize. No second-round reports were generated before handoff.

**Recovery.** Canonical `recover` reported success, `FINALIZATION-NONE`, `RECOVERY-TAKEOVER-AVAILABLE`, decision “takeover available; claim preserved”. The structural evidence is line 10 of `do-work/CHECKPOINT.md`, writer `t2s-Virtual-Machine.local:/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`, claimed `2026-11-18T23:53:11Z`. Use per-request advance to continue the preserved claim.

**Earlier releases.** Report/figure feedback is 0.329.0 (`742088b1`), named starters are 0.330.0 (`fdb8c240`), and feedback-guided generation is 0.331.0 (`37fb1718`). The durable integration receipts are in `do-work/runs/work-2026-11-18-191126/`. Earlier approval dialogs occurred on compound test wrappers. Existing approval settings were retained. Use direct commands and surface any genuine permission request.

**Parallelism.** Only the trial remains. Independent second-round reports can use disjoint output directories and worktrees. Three concurrent worker lanes fit the coordinator run. Serialize gallery publication and queue finalization. The binding claim explicitly requires the user's first-round marks before second-round generation.

**Files captured by the main WIP commit.**

- `ai-reports/2026-11-19_0154_req748-cpu-elapsed-audit/gallery-design.md`
- `ai-reports/2026-11-19_0154_req748-cpu-elapsed-audit/gallery-factual-review.json`
- `ai-reports/2026-11-19_0154_req748-cpu-elapsed-audit/gallery-generation.json`
- `ai-reports/2026-11-19_0154_req748-cpu-elapsed-audit/index.html`
- `ai-reports/2026-11-19_0154_req748-cpu-elapsed-audit/presentation.css`
- `ai-reports/2026-11-19_0154_req748-cpu-elapsed-audit/presentation.js`
- `ai-reports/2026-11-19_0154_req748-latency-ranking/gallery-design.md`
- `ai-reports/2026-11-19_0154_req748-latency-ranking/gallery-factual-review.json`
- `ai-reports/2026-11-19_0154_req748-latency-ranking/gallery-generation.json`
- `ai-reports/2026-11-19_0154_req748-latency-ranking/index.html`
- `ai-reports/2026-11-19_0154_req748-latency-ranking/presentation.css`
- `ai-reports/2026-11-19_0154_req748-latency-ranking/presentation.js`
- `ai-reports/2026-11-19_0154_req748-replicate-variability/gallery-design.md`
- `ai-reports/2026-11-19_0154_req748-replicate-variability/gallery-factual-review.json`
- `ai-reports/2026-11-19_0154_req748-replicate-variability/gallery-generation.json`
- `ai-reports/2026-11-19_0154_req748-replicate-variability/index.html`
- `ai-reports/2026-11-19_0154_req748-replicate-variability/presentation.css`
- `ai-reports/2026-11-19_0154_req748-replicate-variability/presentation.js`
- `ai-reports/2026-11-19_0154_req748-tool-invocation-counts/gallery-design.md`
- `ai-reports/2026-11-19_0154_req748-tool-invocation-counts/gallery-factual-review.json`
- `ai-reports/2026-11-19_0154_req748-tool-invocation-counts/gallery-generation.json`
- `ai-reports/2026-11-19_0154_req748-tool-invocation-counts/index.html`
- `ai-reports/2026-11-19_0154_req748-tool-invocation-counts/presentation.css`
- `ai-reports/2026-11-19_0154_req748-tool-invocation-counts/presentation.js`
- `ai-reports/2026-11-19_0154_req748-warm-cold-labels/gallery-design.md`
- `ai-reports/2026-11-19_0154_req748-warm-cold-labels/gallery-factual-review.json`
- `ai-reports/2026-11-19_0154_req748-warm-cold-labels/gallery-generation.json`
- `ai-reports/2026-11-19_0154_req748-warm-cold-labels/index.html`
- `ai-reports/2026-11-19_0154_req748-warm-cold-labels/presentation.css`
- `ai-reports/2026-11-19_0154_req748-warm-cold-labels/presentation.js`
- `ai-reports/catalog.json`
- `ai-reports/index.html`
- `do-work/working/REQ-748-first-two-round-report-trial.md`
- `do-work/working/baseline.json`
- `docs/report-gallery/index.html`

**Worktree survey.** ACTIVE denotes retained trial resources even though all workers are stopped. FOREIGN covers earlier-session work outside this trial. Cleanup is optional and deferred.

| Verdict | Checkout | Branch and tip | Disposition |
| --- | --- | --- | --- |
| ACTIVE | `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` | `main` · `160a96b5` | Keep for the open trial. Working tree clean. |
| FOREIGN | `/private/tmp/ur185-round7-recheck/worktree` | `detached` · `58f81933` | Preserve this earlier checkout. Existing local changes: M  docs/report-gallery/index.html. |
| ACTIVE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-0121cd34-af4a77dd992b15f9` | `t3/req748-round-one-report` · `a51cb141` | WIP snapshot of `ai-reports/2026-11-19_0154_req748-replicate-variability/gallery-factual-review.json`, `ai-reports/2026-11-19_0154_req748-replicate-variability/gallery-generation.json`, `ai-reports/2026-11-19_0154_req748-replicate-variability/index.html`. The reviewed page is also in the main snapshot. |
| REMOVABLE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-185bb2f0-a4b63cdde966883f` | `t3/test-ai-report-consumer` · `82a0eefe` | After rechecking the archived claim, clean status and ancestry: `git worktree remove /Users/t2/.t3/worktrees/skill-do-work2/t3-185bb2f0-a4b63cdde966883f`. |
| ACTIVE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-2c85c9be-ecaf06073373ac5a` | `t3/req748-round-one-report-2` · `e6aaa3b7` | WIP snapshot of `ai-reports/2026-11-19_0154_req748-tool-invocation-counts/gallery-factual-review.json`, `ai-reports/2026-11-19_0154_req748-tool-invocation-counts/gallery-generation.json`, `ai-reports/2026-11-19_0154_req748-tool-invocation-counts/index.html`. The reviewed page is also in the main snapshot. |
| REMOVABLE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-2dd71bad-24becc3c26aa1645` | `t3/blind-evaluation-logistics` · `e26f4962` | After rechecking the archived claim, clean status and ancestry: `git worktree remove /Users/t2/.t3/worktrees/skill-do-work2/t3-2dd71bad-24becc3c26aa1645`. |
| REMOVABLE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-2e195558-650cc849445a4539` | `t3/implement-cli-report-gallery` · `683dcf1d` | After rechecking the archived claim, clean status and ancestry: `git worktree remove /Users/t2/.t3/worktrees/skill-do-work2/t3-2e195558-650cc849445a4539`. |
| FOREIGN | `/Users/t2/.t3/worktrees/skill-do-work2/t3-56b0e749-05f74d2d71cb1e03` | `t3/add-report-gallery-examples` · `08109e2e` | Preserve this earlier checkout. |
| REMOVABLE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-577dc790-5b54304e8da8b52c` | `t3/test-feedback-starter-interaction` · `e6c2229e` | After rechecking the archived claim, clean status and ancestry: `git worktree remove /Users/t2/.t3/worktrees/skill-do-work2/t3-577dc790-5b54304e8da8b52c`. |
| REMOVABLE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-706b61e5-35adaa5f9e352122` | `t3/feedback-report-acceptance` · `f3840b08` | After rechecking the archived claim, clean status and ancestry: `git worktree remove /Users/t2/.t3/worktrees/skill-do-work2/t3-706b61e5-35adaa5f9e352122`. |
| REMOVABLE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-71a70ed1-90b9d3ce4cc36f6c` | `t3/req-744-comparison-acceptance` · `82a0eefe` | After rechecking the archived claim, clean status and ancestry: `git worktree remove /Users/t2/.t3/worktrees/skill-do-work2/t3-71a70ed1-90b9d3ce4cc36f6c`. |
| REMOVABLE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-7a614596-6fcdc8c82ce1f54c` | `t3/upgrade-gallery-build` · `fe0b06b7` | After rechecking the archived claim, clean status and ancestry: `git worktree remove /Users/t2/.t3/worktrees/skill-do-work2/t3-7a614596-6fcdc8c82ce1f54c`. |
| FOREIGN | `/Users/t2/.t3/worktrees/skill-do-work2/t3-7f608bfc-d20c197b760615fa` | `t3/gallery-design-selection` · `f006ae49` | Preserve this earlier checkout. |
| ACTIVE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-7fb8af4f-3df426ac4570c2f1` | `t3/req748-latency-ranking-report` · `5d9f4007` | WIP snapshot of `ai-reports/2026-11-19_0154_req748-latency-ranking/gallery-factual-review.json`, `ai-reports/2026-11-19_0154_req748-latency-ranking/gallery-generation.json`, `ai-reports/2026-11-19_0154_req748-latency-ranking/index.html`. The reviewed page is also in the main snapshot. |
| FOREIGN | `/Users/t2/.t3/worktrees/skill-do-work2/t3-84cfbb98-7393f63a98b9a7c0` | `t3/add-gallery-examples` · `ba07b1ed` | Preserve this earlier checkout. |
| REMOVABLE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-887ad5f6-99fe91700f1dbf6e` | `t3/ai-report-kind-acceptance` · `82a0eefe` | After rechecking the archived claim, clean status and ancestry: `git worktree remove /Users/t2/.t3/worktrees/skill-do-work2/t3-887ad5f6-99fe91700f1dbf6e`. |
| REMOVABLE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-8a8286f3-d401a8cc593cf9c3` | `t3/ur185-round7-gate-recheck` · `3dae1dc9` | After rechecking the archived claim, clean status and ancestry: `git worktree remove /Users/t2/.t3/worktrees/skill-do-work2/t3-8a8286f3-d401a8cc593cf9c3`. |
| REMOVABLE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-903dfe0f-9c784110da458af1` | `t3/implement-agent-forms-hooks` · `a8031e04` | After rechecking the archived claim, clean status and ancestry: `git worktree remove /Users/t2/.t3/worktrees/skill-do-work2/t3-903dfe0f-9c784110da458af1`. |
| REMOVABLE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-9b22b92f-0ac626ce5c13de96` | `t3/build-report-gallery-samples` · `813c0517` | After rechecking the archived claim, clean status and ancestry: `git worktree remove /Users/t2/.t3/worktrees/skill-do-work2/t3-9b22b92f-0ac626ce5c13de96`. |
| FOREIGN | `/Users/t2/.t3/worktrees/skill-do-work2/t3-a5de53e3-5b37f50fbf62596b` | `t3/req-744-acceptance-fixtures` · `82a0eefe` | Preserve this earlier checkout. Existing local changes: ?? do-work/runs/work-2026-11-18-153357/REQ-744-fixture-handback.md. |
| ACTIVE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-af3c5181-2030a073fd42b0c3` | `t3/verify-evidence-report-candidate` · `f3840b08` | Keep for the open trial. Working tree clean. |
| ACTIVE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-b62fadd9-321f6ddf24c8d44f` | `t3/fact-check-req748-trial` · `d2ac8480` | Keep for the open trial. Working tree clean. |
| REMOVABLE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-c608b8cc-2ee895741fc0d8aa` | `t3/implement-req-744-ai-report` · `5bdc38b2` | After rechecking the archived claim, clean status and ancestry: `git worktree remove /Users/t2/.t3/worktrees/skill-do-work2/t3-c608b8cc-2ee895741fc0d8aa`. |
| FOREIGN | `/Users/t2/.t3/worktrees/skill-do-work2/t3-cfabdc05-3285aedf59f7f49d` | `t3/req-750-browser-gate-tests` · `3dc5bc22` | Preserve this earlier checkout. |
| REMOVABLE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-d86889b6-4eb969cb02064aaa` | `t3/prep-local-final-report-generator` · `8e1ac984` | After rechecking the archived claim, clean status and ancestry: `git worktree remove /Users/t2/.t3/worktrees/skill-do-work2/t3-d86889b6-4eb969cb02064aaa`. |
| ACTIVE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-e1a00437-34bc521d7f89eb40` | `t3/req748-round-one-report-1` · `1a2d5fed` | WIP snapshot of `ai-reports/2026-11-19_0154_req748-warm-cold-labels/gallery-factual-review.json`, `ai-reports/2026-11-19_0154_req748-warm-cold-labels/gallery-generation.json`, `ai-reports/2026-11-19_0154_req748-warm-cold-labels/index.html`. The reviewed page is also in the main snapshot. |
| ACTIVE | `/Users/t2/.t3/worktrees/skill-do-work2/t3-ee16fbde-fdc0811501ff7d3b` | `t3/req748-cpu-elapsed-audit` · `34fe2b6a` | WIP snapshot of `ai-reports/2026-11-19_0154_req748-cpu-elapsed-audit/gallery-factual-review.json`, `ai-reports/2026-11-19_0154_req748-cpu-elapsed-audit/gallery-generation.json`, `ai-reports/2026-11-19_0154_req748-cpu-elapsed-audit/index.html`. The reviewed page is also in the main snapshot. |
