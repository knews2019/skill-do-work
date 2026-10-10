---
id: UR-153
title: 'Simplify the 2026-10-09 upstream-suggestion captures: fold twenty REQs into eleven'
created_at: 2026-10-10T13:05:57Z
requests: [REQ-687, REQ-688, REQ-689, REQ-690, REQ-691, REQ-692]
word_count: 680
---
# Simplify the 2026-10-09 Upstream-Suggestion Captures: Fold Twenty REQs Into Eleven

## Summary
The maintainer asked whether the twenty pending REQs captured on 2026-10-09 from five upstream suggestion reports (UR-144 to UR-149) should be simplified, and chose the aggressive option through the ask tool. Fifteen REQs were cancelled through `do-work abandon` (commit 5e2752d5) and recaptured here as six: REQ-654 and REQ-661 (UR-144), REQ-662 to REQ-667 (UR-146), REQ-668 (UR-147), REQ-669 (UR-148) and REQ-670 to REQ-674 (UR-149). Five REQs were left as captured: REQ-655, REQ-657, REQ-658, REQ-659 and REQ-660 (claimed, waiting for the next session), and REQ-656 (pending). The cancelled originals keep their full bodies in the archive; each new REQ names its original and the source report, so no detail is lost.

What changed, per the findings in the verbatim input: the validate-feedback chain and its routing row became one REQ and the installer retire line was dropped as operator work (F1, F2); the six coordinate-mode REQs became one prose REQ, the two Go commands that would have checked whether a lock's owner pid is dead were dropped because REQ-073 rules out a PID check in shipped code, and the stall check now reads `do-work status --watch` (F3, F4, F5); `do-work status` lost its `--fix` part and `do-work trace` lost its Notion browser adapter and new UR body section (F6); `ai-report` lost the `options` kind and `capture-files init` became the source report's own `--example` fallback (F7). REQ-656's edge to REQ-655 was left in place because REQ-655 is already being built.

## Extracted Requests
| REQ | Title |
|-----|-------|
| REQ-687 | ai-report --kind proposal and root-cause writes a decision-first brief for unfinished work |
| REQ-688 | capture-files --example prints one valid manifest with payload templates, and the capture-reference fence example is fixed |
| REQ-689 | do-work run --coordinate: the mandatory coordinator shape, its run rules, a prose preflight, a stall loop, and a handoff that resumes coordinated |
| REQ-690 | do-work status action and do-work-cli run-status report one class, ETA and remedy per open REQ |
| REQ-691 | do-work trace action reports how much of an outside spec is captured and built, with one dated verdict per ask |
| REQ-692 | validate-feedback --capture [--run]: file-path input, wrong-repo check, Discuss questions, one capture, verify, optional run, and a core do-work route |

## Batch Constraints
- No new frontmatter field and no new status anywhere in the batch.
- No PID check, lock helper, heartbeat or liveness probe in shipped Go code (REQ-073); liveness judgments stay in the action layer, on this machine, labelled as such.
- Each REQ is its own release per `_dev/primes/prime-releases.md`.
- REQ-690, REQ-691 and REQ-692 each add a row to `skills/do-work/SKILL.md`'s routing table; whichever lands later rebases its row. No `depends_on` edges in this batch.
- The cancelled originals under `do-work/archive/` are the full-detail record; builders read them when a new body summarises.

## Full Verbatim Input
> ```
> User message (2026-10-10, session skill-do-work2-b1):
> 
> see the following REQs, should we simplify them?
> 
> # the reqs
> 
> [The user pasted the board rendering of twenty pending REQ files: REQ-654, REQ-655, REQ-656, REQ-657 (UR-144), REQ-658, REQ-659, REQ-660, REQ-661 (UR-145), REQ-662, REQ-663, REQ-664, REQ-665, REQ-666, REQ-667 (UR-146), REQ-668 (UR-147), REQ-669 (UR-148), REQ-670, REQ-671, REQ-672, REQ-673, REQ-674 (UR-149). The pasted bodies are the queue files as of commit 901198c5; the cancelled ones are archived unchanged under do-work/archive/ (REQ-654, REQ-661 at the root; the rest under archive/UR-146, UR-147, UR-148, UR-149).]
> 
> Assistant verdict, as presented to the user:
> 
> Yes. Twenty REQs for five features is too many, and two of them contradict a third. One fact first: the pasted status is stale. Eleven of the twenty (654, 655, 657 to 661, 667 to 670, 674) were claimed by the other live session about 15 seconds before I looked, with no builder worktrees yet.
> 
> Findings:
> - F1, validate-feedback chain (670 to 673). One flag, one 154-line action file, four serial REQs and four releases. Fold into one REQ.
> - F2, installer retire line (674 C5b). Shipped Go plus a closed name list, to save a one-time git rm in nine repos. That is operator work. Keep only the routing row and fold it into F1.
> - F3, coordinate batch prose (662, 666, 667). All three edit the same two prose files. One REQ.
> - F4, PID checks in Go (663, 665). Both print whether a lock's owner process is dead, as tested Go commands. REQ-073 rules out a PID check (verify.go:1123), and 668 in the same queue refuses it for that reason. 665 also parses the prose handoff, which your no-grammar ruling warns against, on three observed cases. Cancel 665. Reduce 663 to a prose checklist (board disk level, run-policy file) inside F3.
> - F5, stall check overlaps status (664 vs 668). Both read liveness. 664 becomes one paragraph in F3: loop `do-work status --watch`, act on the quiet and stale rows. It then depends on 668.
> - F6, the two big ones. 668: drop P5 --fix, a wrapper around a remedy line already printed. 669: drop the Notion adapter (site-specific CSS selectors in a platform-agnostic skill) and the new UR body section; the ticket id goes in the UR title.
> - F7, ai-report batch. 654: drop the `options` kind, which its own assumption says equals `proposal`. 656: drop the edge to 655 and regenerate the catalog only if one exists. 661: take the report's own fallback (--example plus the fence fix), since it admits the thinnest evidence. 657 stays `later`.
> 
> 658, 659 and 660 stay as they are. The counts behind them are large and they are programs replacing prose.
> 
> Question 1 (ask tool): "Which simplification level should be applied?"
> Options: "Aggressive (Recommended)" (all of F1 to F7; 20 REQs become about 11) / "Baseline (pending only)" / "Leave as is".
> Answer: Aggressive (Recommended).
> 
> User, mid-turn: "phandoff, also check the REQs they will be simplified" (typed into the parallel session skill-do-work2-77, screenshot shared) and "the parallel session now is stopped". The parallel session released REQ-654, 661, 667, 668, 669, 670, 674 back to the queue (recover --take-over, commits 2b109946 to 61e02428) and wrote its handoff (c644ee47). REQ-655, 657, 658, 659, 660 stay claimed in do-work/working/ for the next session.
> 
> Question 2 (ask tool): "Cancel these 15 pending REQs and recapture them as 6? ... Untouched: 655, 656, 657, 658, 659, 660."
> Options: "Yes, cancel 15 and recapture 6 (Recommended)" (R1 ai-report --kind proposal|root-cause; R2 capture-files --example plus the fence-example fix; R3 run --coordinate: rules, flag, prose preflight, stall loop, handoff and teardown; R4 do-work status minus --fix; R5 do-work trace minus the Notion adapter and PM Ticket section; R6 validate-feedback --capture [--run] with the core routing row, minus the installer check) / "Cancel 15, but keep 674's installer check" / "Stop here, cancel nothing".
> Answer: Yes, cancel 15 and recapture 6 (Recommended).
> 
> Decision recorded while drafting: REQ-656 is left untouched after all. Its only proposed change was dropping the depends_on edge to REQ-655, and REQ-655 is already claimed and being built, so the edge costs nothing.
> ```
