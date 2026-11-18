# REQ-744 integration handback

Released **0.326.0** locally on main. REQ-744 is **completed** after independent acceptance. Step 9 returned `cleanup_complete` with no blockers. The report writer now uses gallery samples to choose how to present the evidence, records its choice, and includes fitting second-opinion prompts.

[Archived request](../../archive/UR-183/REQ-744-report-writer-ranks-gallery-samples.md) · [Finalization receipt](REQ-744-finalization.json)

| Revision | Commit |
| --- | --- |
| Integration base | `84f371e577923ed1e156581bf1f8886cf01d5432` |
| Builder candidate | `5bdc38b2dd2c86e92658f73ffc73d86fb97881f8` |
| Main merge and verified implementation | `bbc1e5d0db0edfb0462b46b51af43b83cd809ac1` |
| Canonical release and archive | `49c2ab2792fd9442c13040924e2547cd74d9e8f4` |

The main merge contained the four scoped files and matched the independently installed candidate. The do-work merge helper rejected the T3 branch namespace, so integration used the documented Git merge sequence. Qualification and scope passed. The five-test focused probe and the canonical gate passed, including 428 board and 995 CLI tests in 121 seconds. The selected `staged-skills` lane executed and passed in 31 seconds. The release version, changelog mirrors, and shipped references passed after finalization. [Merge](REQ-744-merge.log) · [Qualification](REQ-744-qualification.json) · [Focused and canonical testing](REQ-744-testing.json) · [Heavy result](REQ-744-heavy-run.json) · [Release verification](REQ-744-release-verification.log)

- **Ordinary comparison and prompt.** The installed report showed 90 → 160 successes and 90% → 80% success rate in aligned, separately labelled charts. It recorded `two-measures-side-by-side` and offered `verdict-with-checks`. The displayed command matched the handback and prepared with the shipped `80-20` preset and the exact direction. Prepare left the original bundle byte-identical and published no sibling. Both four-view render passes were green. [Comparison acceptance](REQ-744-comparison-acceptance.md) · [Prepare receipt](REQ-744-comparison-prepare-verified.json)
- **No fit, revise, and editor control.** The spelling report recorded `none` and kept a content-specific layout and a teaching prompt. Its revision ranked the evidence again, retained the prior bundle bytes, and had only the expected prior-bundle link finding. The link resolved on disk. The explicit-page control prepared the requested editor direction. The revision used text-only evidence. The planned comparison-evidence revision was not run separately. [No-match and revise acceptance](REQ-744-no-match-acceptance.md)
- **Proposal and root-cause.** Both metrics briefs matched the paired measures, retained their kind-specific sections, and offered a valid preset and single direction for the runner-up. A wording-only proposal exercised `none`. The writer corrected a percent-label wrap spotted during visual review. All three final briefs passed the four-view judge. [Kind acceptance](REQ-744-kind-acceptance.md)

The wave-end sweep found one minor stale restatement. The kit CSS header at line 17 still describes one alternate layout. That source comment has no rendering or command effect and is recorded as report only. The action, reference, skeleton, and guide carry the current contract.

Finalization closed UR-183 and moved REQ-743 into its archive folder byte for byte. The shipped lesson link now resolves to that path. UR-184 and all four queued REQs remain byte-identical. The accepted receipts, run evidence and focused baseline are recorded locally. T3 worktrees remain available. The coordinator owns the run manifest and Step 10.
