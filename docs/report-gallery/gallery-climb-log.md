# Blind gallery climb

Round 0 began at `eb5ddd58` with 79/100, 4/4 reader needs matched and all seven gates green across four views. The [frozen baseline](baseline/README.md) contains the original score and evidence.

| Round | Candidate change | Gates | Incumbent | Candidate | Reader matches (incumbent / candidate) | Decision |
| --- | --- | --- | ---: | ---: | --- | --- |
| [1](climb/round-1/decision.json) | Align the release dossier name and subject with the board client report | [28/28](climb/round-1/views/gates.json) | 74 | 75 | 4/4 · 4/4 | Revert. Total gain was +1. |
| [2](climb/round-2/decision.json) | Label the duration-label alternatives as a proposal and say what it compares | [28/28](climb/round-2/views/gates.json) | 72 | 71 | 4/4 · 4/4 | Revert. Total change was −1. |
| [3](climb/round-3/decision.json) | Group the cards by kind and retire the duplicate kind list | [28/28](climb/round-3/views/gates.json) | 77 | 81 | 4/4 · 4/4 | Keep. Attention improved by 4 with every other score steady. |
| [4](climb/round-4/decision.json) | Put both diagnostic reports under “Diagnose a problem” | [28/28](climb/round-4/views/gates.json) | 72 | 75 | 4/4 · 4/4 | Keep. Choice gained 2 and Attention gained 1. |
| [5](climb/round-5/candidate.patch) | Distinguish the Tinka proposal by its subject | [G2, G5, G6 failed in all four views](climb/round-5/views/gates.json) | unscored | unscored | — | Revert. The rename also changed the source URL. |
| [6](climb/round-6/decision.json) | Rename and describe the review canvas while preserving the source URL | [28/28](climb/round-6/views/gates.json) | 78 | 79 | 4/4 · 4/4 | Revert. Total gain was +1. |
| [7](climb/round-7/decision.json) | Lead the cards with their subjects so readers can distinguish reports before reading the previews | [28/28](climb/round-7/gate.log) | 73 | 74 | 4/4 · 4/4 | Revert. Total gain was +1. |

Round 3 qualified on score and became the comparison incumbent. Its Proposals-only classification of the duration-label report left just one visible root-cause example, so that intermediate HTML was ineligible for final integration. Round 4 restored two visible examples in the diagnosis group while qualifying against Round 3.

## Endpoint

The climb stopped after the third consecutive revert (Rounds 5–7). Across seven rounds, two qualified, four reverted on score and one reverted unscored after a gate failure.

Round 4 at `58f81933` is the final incumbent, at **75/100** on its qualifying pair, with **4/4** needs matched and **28/28** gate checks green. The 90-point goal was unmet. [Scored HTML](index.html) · [Desktop light](climb/round-4/views/1440-light.png) · [Desktop dark](climb/round-4/views/1440-dark.png) · [Phone light](climb/round-4/views/390-light.png) · [Phone dark](climb/round-4/views/390-dark.png) · [Gates](climb/round-4/views/gates.json).

Normal regeneration reproduced the adjudicated HTML byte for byte, including the stray closing div called out by the Round 6 evaluator. Both diagnostic examples are visible under “Diagnose a problem.” The duration-labels report also presents design alternatives in its preview and source report. The final gallery has 3 completed-work, 3 proposal, 2 diagnosis and 2 architecture cards. [Parity and endpoint checks](climb/final-endpoint-check.log) · [Frozen baseline checksums](climb/frozen-baseline-check.log).
