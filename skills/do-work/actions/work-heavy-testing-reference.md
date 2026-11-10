# Work phase: heavy-testing

> Companion to [work.md](./work.md); core queue procedures. Read when its phase routing selects this file. Existing step numbers remain stable.

### Step 7.7: Heavy-Test Hold and Drain

**Heavy-test hold:** after Steps 7 and 7.5 pass (or a selected AI-Evaluation-Loop returns unresolved, which remains unapproved), a REQ whose Testing section recorded selected heavy lanes does **not** run those lanes and does not hold the queue loop open. Record the fast commands and durations plus the typed lane plan, then:

1. land the implementation commit and record its full hash in the REQ's canonical `commit:` field;
2. append `## Heavy Verification Plan` with the exact base/target revisions and each selected lane's id, argv, and reasons; and
3. recompute selection and continue every unrelated runnable REQ.

The REQ stays `claimed` in `do-work/working/`: no status change, no `## Open Questions` machine line, no move to `do-work/queue/`, and no checkpoint edit. The landed `commit:` is what makes its source ready for dependent REQs while it waits.

**Heavy-lane drain (at queue exhaustion):** when no claimable pending REQ remains and at least one held claimed REQ exists — including one routed here by `recover`'s `RECOVERY-CLAIM-HELD-FOR-HEAVY-LANES` finding — this loop runs the held lanes itself. No permission is asked. Per held REQ, recompute `plan-heavy-verification` at its stored base and `commit:` and refuse any drift against the stored `## Heavy Verification Plan`. Plan drift, or a stored `historical-revalidation` plan, leaves that REQ held and is a typed finding for a human, never a hand edit. Union the selected lane ids across every remaining REQ and run them once at HEAD, repeating `--lane` for each unioned id:

```bash
<skill-root>/tools/do-work-cli.sh --repo-root "<project-root>" --format json \
  run-heavy-verification --manifest _dev/tests/heavy-lanes.json --lane "<lane-id>"
```

The runner refuses an unknown lane and a dirty tracked tree outside `do-work/` before anything executes; `HEAVY-RUN-DIRTY-TREE` is a typed condition to resolve, never a hand edit. It reports each lane's exit status, skip state, wall seconds, and `disposition`, and a red or skipped lane is a warning finding rather than a command failure.

A lane whose deterministic fingerprint — its argv, its covered and globally unclassified committed paths, regular untracked input bytes, bounded toolchain probes, and the entire inherited environment — still matches a successful record no older than four hours is reported `reused` instead of executed; every other lane is `executed` with the exact reason it was not reused (`fingerprint_mismatch`, `evidence_expired`, `no_prior_evidence`, `fingerprint_uncertain`). Execution revokes the previous success before launch; an invalidation failure refuses the run. Fingerprint uncertainty (including an opaque browser runtime or a symlink input) always executes. Shipped lane argv explicitly isolates system/global Git configuration; custom manifest commands retain their declared environment. Report each lane's disposition alongside its result: a `reused` green was measured by an earlier run, not this one. Add `--no-evidence-reuse` when a lane must execute regardless — a suspected flake, or an input outside its declared fingerprint.

Then dispose of each held REQ from the subset of lanes it selected, with no `answer` transaction in any outcome:

- **Green** — every lane that REQ selected is present in the run, exited 0, and was not skipped. Write `heavy_verified_at: <now>` (current UTC instant — Timestamp rule, `actions/work-reference.md`), `heavy_verified_revision: <runner's execution revision>`, and a `## Heavy Verification Result` section naming the target revision, the execution revision, and one line per lane onto the claimed record, then run Steps 8 and 9 for it in this same turn.
- **Red** — any selected lane exited non-zero. Delete `commit:` and the `## Heavy Verification Plan` section from the record so dependents stop building against withdrawn work, then enter Step 7's remediation path now and re-hold at this step after a fresh review. A selected AI-Evaluation-Loop retains its counters and stopping limits; exhausted review remains unresolved.
- **Skipped** — any selected lane was skipped. The REQ stays claimed and held, its `HEAVY-RUN-LANE-SKIPPED` finding names the lane, the next drain retries it, and the composed exit summary names it if the run ends while it is still held.

A `commit:` drained once this run is never drained again.
